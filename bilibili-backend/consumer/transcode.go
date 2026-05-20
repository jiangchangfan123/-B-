package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"bilibili-backend/config"
	"bilibili-backend/dao"
	"bilibili-backend/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

// StartTranscodeConsumer 启动转码消费者
func StartTranscodeConsumer(db *gorm.DB) {
	if utils.RabbitCh == nil {
		log.Println("[Transcode] RabbitMQ not connected, consumer not started")
		return
	}

	queue := config.C.RabbitMQ.Queue
	msgs, err := utils.RabbitCh.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Printf("[Transcode] Failed to start consumer: %v", err)
		return
	}

	videoDao := dao.NewVideoDao(db)

	go func() {
		for msg := range msgs {
			processTranscodeTask(msg, videoDao)
		}
	}()
	log.Println("[Transcode] Consumer started, listening queue:", queue)
}

func processTranscodeTask(msg amqp.Delivery, videoDao *dao.VideoDao) {
	var task utils.TranscodeMessage
	if err := json.Unmarshal(msg.Body, &task); err != nil {
		log.Printf("[Transcode] Invalid message: %v", err)
		msg.Nack(false, false)
		return
	}

	log.Printf("[Transcode] Processing video_id=%d", task.VideoID)

	// 1. 准备临时目录
	tmpDir := filepath.Join(os.TempDir(), "transcode", fmt.Sprintf("%d", task.VideoID))
	os.MkdirAll(tmpDir, 0755)
	defer os.RemoveAll(tmpDir)

	// 从 InputURL 解析输入原片的 object key
	inputKey := extractKeyFromURL(task.InputURL)
	if inputKey == "" {
		log.Printf("[Transcode] Failed to extract input key from URL: %s", task.InputURL)
		_ = videoDao.UpdateTranscodeStatus(task.VideoID, 3, "", 2)
		msg.Nack(false, false)
		return
	}

	inputPath := filepath.Join(tmpDir, "input"+filepath.Ext(inputKey))
	outputPath := filepath.Join(tmpDir, "output_480p.mp4")

	// 2. 从 MinIO 下载原片
	ctx := context.Background()
	if err := utils.MinIOFGetObject(ctx, task.Bucket, inputKey, inputPath); err != nil {
		log.Printf("[Transcode] Download failed: %v", err)
		_ = videoDao.UpdateTranscodeStatus(task.VideoID, 3, "", 2)
		msg.Nack(false, false)
		return
	}

	// 3. 检查 FFmpeg
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		log.Println("[Transcode] FFmpeg not found, skipping transcode")
		_ = videoDao.UpdateTranscodeStatus(task.VideoID, 3, "", 2)
		msg.Nack(false, false)
		return
	}

	// 4. FFmpeg 转码
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-vf", "scale=-2:480",
		"-c:v", "libx264",
		"-preset", "fast",
		"-crf", "23",
		"-c:a", "copy",
		"-movflags", "+faststart",
		"-y", outputPath,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("[Transcode] FFmpeg failed: %v", err)
		_ = videoDao.UpdateTranscodeStatus(task.VideoID, 3, "", 2)
		msg.Nack(false, false)
		return
	}

	// 5. 上传转码后文件到 MinIO
	outputKey := fmt.Sprintf("videos/480p/%d_480p.mp4", task.VideoID)
	if err := utils.MinIOUploadFile(ctx, task.Bucket, outputKey, outputPath, "video/mp4"); err != nil {
		log.Printf("[Transcode] Upload failed: %v", err)
		_ = videoDao.UpdateTranscodeStatus(task.VideoID, 3, "", 2)
		msg.Nack(false, false)
		return
	}

	transcodedURL := utils.MinIOGetURL(ctx, task.Bucket, outputKey)

	// 6. 更新数据库
	if err := videoDao.UpdateTranscodeStatus(task.VideoID, 2, transcodedURL, 1); err != nil {
		log.Printf("[Transcode] DB update failed: %v", err)
		msg.Nack(false, false)
		return
	}

	log.Printf("[Transcode] Success video_id=%d, url=%s", task.VideoID, transcodedURL)
	msg.Ack(false)
}

func extractKeyFromURL(url string) string {
	// http://host/bucket/key
	// key 可能包含多个 / ，用 SplitN 保留完整 key
	parts := strings.SplitN(url, "/", 5)
	if len(parts) < 5 {
		return ""
	}
	return parts[4]
}
