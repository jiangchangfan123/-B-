package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"bilibili-backend/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitConn *amqp.Connection
var RabbitCh *amqp.Channel

// TranscodeMessage 转码任务消息结构
type TranscodeMessage struct {
	VideoID   uint64 `json:"video_id"`
	InputURL  string `json:"input_url"`  // MinIO 可访问的原片 URL
	Bucket    string `json:"bucket"`     // 输出 bucket
	ObjectKey string `json:"object_key"` // 输出文件路径
}

// InitRabbitMQ 初始化 RabbitMQ 连接
func InitRabbitMQ() error {
	cfg := config.C.RabbitMQ
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s/", cfg.Username, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(uri)
	if err != nil {
		return fmt.Errorf("rabbitmq connect failed: %w", err)
	}
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("rabbitmq channel failed: %w", err)
	}
	RabbitConn = conn
	RabbitCh = ch

	// 声明队列
	queue := config.C.RabbitMQ.Queue
	_, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("rabbitmq queue declare failed: %w", err)
	}

	log.Println("[RabbitMQ] Connected successfully, queue:", queue)
	return nil
}

// PublishTranscodeTask 发送转码任务
func PublishTranscodeTask(msg TranscodeMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	queue := config.C.RabbitMQ.Queue
	return RabbitCh.PublishWithContext(
		context.Background(),
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// ConsumeTranscodeTasks 消费转码任务
func ConsumeTranscodeTasks(handler func(amqp.Delivery)) error {
	queue := config.C.RabbitMQ.Queue
	msgs, err := RabbitCh.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack 手动确认
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return err
	}
	go func() {
		for msg := range msgs {
			handler(msg)
		}
	}()
	log.Println("[RabbitMQ] Consumer started")
	return nil
}

// CloseRabbitMQ 关闭连接
func CloseRabbitMQ() {
	if RabbitCh != nil {
		RabbitCh.Close()
	}
	if RabbitConn != nil {
		RabbitConn.Close()
	}
}
