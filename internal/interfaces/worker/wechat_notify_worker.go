package worker

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/messaging/rabbitmq"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/wechat"
	"github.com/YazaiHu/MemberHub/internal/pkg/logger"
)

// WechatNotifyWorker 微信通知Worker
type WechatNotifyWorker struct {
	rateLimiter *time.Ticker // 限流器（20条/秒）
}

// NewWechatNotifyWorker 创建微信通知Worker
func NewWechatNotifyWorker() *WechatNotifyWorker {
	return &WechatNotifyWorker{
		rateLimiter: time.NewTicker(50 * time.Millisecond), // 20条/秒
	}
}

// WechatNotifyMessage 微信通知消息
type WechatNotifyMessage struct {
	UserID     int64                  `json:"user_id"`
	TemplateID string                 `json:"template_id"`
	Data       map[string]interface{} `json:"data"`
	Page       string                 `json:"page,omitempty"`
}

// Start 启动Worker
func (w *WechatNotifyWorker) Start() error {
	logger.Info("Starting WechatNotifyWorker")

	// 获取消费者通道
	msgs, err := rabbitmq.Consume("queue.notify.wechat", "wechat-notify-worker")
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	// 处理消息
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			w.handleMessage(d)
		}
	}()

	logger.Info("WechatNotifyWorker started successfully")
	<-forever

	return nil
}

// handleMessage 处理单条消息
func (w *WechatNotifyWorker) handleMessage(d amqp.Delivery) {
	// 限流
	<-w.rateLimiter.C

	// 解析消息
	var msg WechatNotifyMessage
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		logger.Error("Failed to unmarshal message",
			zap.Error(err),
			zap.String("body", string(d.Body)),
		)
		d.Nack(false, false) // 拒绝消息，不重新入队
		return
	}

	logger.Info("Sending wechat notification",
		zap.Int64("user_id", msg.UserID),
		zap.String("template_id", msg.TemplateID),
	)

	// 发送微信模板消息
	if err := w.sendTemplateMessage(msg); err != nil {
		logger.Error("Failed to send wechat notification",
			zap.Int64("user_id", msg.UserID),
			zap.String("template_id", msg.TemplateID),
			zap.Error(err),
		)

		// 重试3次
		if d.Headers != nil {
			retryCount, ok := d.Headers["x-retry-count"].(int32)
			if !ok {
				retryCount = 0
			}

			if retryCount < 3 {
				// 设置重试次数并重新入队
				d.Headers["x-retry-count"] = retryCount + 1
				d.Nack(false, true)
				return
			}
		}

		// 超过重试次数，丢弃消息
		logger.Error("Message discarded after max retries",
			zap.Int64("user_id", msg.UserID),
		)
		d.Nack(false, false)
		return
	}

	logger.Info("Wechat notification sent successfully",
		zap.Int64("user_id", msg.UserID),
	)

	// 确认消息
	d.Ack(false)
}

// sendTemplateMessage 发送微信模板消息
func (w *WechatNotifyWorker) sendTemplateMessage(msg WechatNotifyMessage) error {
	// TODO: 实际调用微信API发送模板消息
	// 这里需要：
	// 1. 获取用户的 openid
	// 2. 构造模板消息数据
	// 3. 调用微信API

	// 使用微信SDK发送模板消息（框架实现）
	templateMsg := &wechat.TemplateMessage{
		ToUser:     fmt.Sprintf("user_%d", msg.UserID), // TODO: 实际应该从数据库获取用户的openid
		TemplateID: msg.TemplateID,
		Page:       msg.Page,
		Data:       msg.Data,
	}

	return wechat.SendTemplateMessage(templateMsg)
}

// Stop 停止Worker
func (w *WechatNotifyWorker) Stop() {
	if w.rateLimiter != nil {
		w.rateLimiter.Stop()
	}
}
