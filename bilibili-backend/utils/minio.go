package utils

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"bilibili-backend/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

// InitMinIO 初始化 MinIO 客户端
func InitMinIO() error {
	cfg := config.C.MinIO
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return fmt.Errorf("minio init failed: %w", err)
	}
	MinioClient = client

	// 创建 bucket（如果不存在）
	ctx := context.Background()
	buckets := []string{cfg.BucketVideos, cfg.BucketCovers, cfg.BucketAvatars}
	for _, bucket := range buckets {
		exists, err := client.BucketExists(ctx, bucket)
		if err != nil {
			return fmt.Errorf("minio check bucket %s failed: %w", bucket, err)
		}
		if !exists {
			if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
				return fmt.Errorf("minio create bucket %s failed: %w", bucket, err)
			}
			log.Printf("[MinIO] Bucket created: %s", bucket)
		}
	}
	log.Println("[MinIO] Connected successfully")
	return nil
}

// MinIOUploadFile 上传文件到 MinIO
func MinIOUploadFile(ctx context.Context, bucket, objectName, filePath, contentType string) error {
	_, err := MinioClient.FPutObject(ctx, bucket, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// MinIOUploadBytes 上传字节数组到 MinIO
func MinIOUploadBytes(ctx context.Context, bucket, objectName string, data []byte, contentType string) error {
	_, err := MinioClient.PutObject(ctx, bucket, objectName, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// MinIORemoveObject 删除 MinIO 上的文件
func MinIORemoveObject(ctx context.Context, bucket, objectName string) error {
	return MinioClient.RemoveObject(ctx, bucket, objectName, minio.RemoveObjectOptions{})
}

// MinIOGetURL 获取 MinIO 文件的访问 URL
func MinIOGetURL(ctx context.Context, bucket, objectName string) string {
	return fmt.Sprintf("http://%s/%s/%s", config.C.MinIO.Endpoint, bucket, objectName)
}

// MinIOFGetObject 从 MinIO 下载文件到本地
func MinIOFGetObject(ctx context.Context, bucket, objectName, filePath string) error {
	return MinioClient.FGetObject(ctx, bucket, objectName, filePath, minio.GetObjectOptions{})
}

// GetMinIOConfig 获取 MinIO 配置
func GetMinIOConfig() config.MinIOConfig {
	return config.C.MinIO
}

// GetBucketFromURL 从 MinIO URL 中提取 bucket 名
func GetBucketFromURL(url string) string {
	// URL 格式: http://endpoint/bucket/object
	parts := strings.Split(url, "/")
	if len(parts) >= 4 {
		return parts[3]
	}
	return ""
}

// GetObjectNameFromURL 从 MinIO URL 中提取 object 名
func GetObjectNameFromURL(url string) string {
	// URL 格式: http://endpoint/bucket/object
	parts := strings.Split(url, "/")
	if len(parts) >= 5 {
		return strings.Join(parts[4:], "/")
	}
	return ""
}

// SplitURL 分割 URL
func SplitURL(url string) []string {
	return strings.Split(url, "/")
}
