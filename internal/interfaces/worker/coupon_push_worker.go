package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"

	couponApp "github.com/YazaiHu/MemberHub/internal/application/coupon"
	"github.com/YazaiHu/MemberHub/internal/domain/coupon"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/messaging/rabbitmq"
	"github.com/YazaiHu/MemberHub/internal/pkg/logger"
)

// CouponPushWorker 优惠券推送Worker
type CouponPushWorker struct {
	couponService *couponApp.Service
}

// NewCouponPushWorker 创建优惠券推送Worker
func NewCouponPushWorker() *CouponPushWorker {
	couponRepo := coupon.NewRepository()
	return &CouponPushWorker{
		couponService: couponApp.NewService(couponRepo),
	}
}

// CouponPushMessage 优惠券推送消息
type CouponPushMessage struct {
	TaskID int64 `json:"task_id"`
}

// Start 启动Worker
func (w *CouponPushWorker) Start() error {
	logger.Info("Starting CouponPushWorker")

	// 获取消费者通道
	msgs, err := rabbitmq.Consume("queue.coupon.push", "coupon-push-worker")
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

	logger.Info("CouponPushWorker started successfully")
	<-forever

	return nil
}

// handleMessage 处理单条消息
func (w *CouponPushWorker) handleMessage(d amqp.Delivery) {
	ctx := context.Background()

	// 解析消息
	var msg CouponPushMessage
	if err := json.Unmarshal(d.Body, &msg); err != nil {
		logger.Error("Failed to unmarshal message",
			zap.Error(err),
			zap.String("body", string(d.Body)),
		)
		d.Nack(false, false) // 拒绝消息，不重新入队
		return
	}

	logger.Info("Processing coupon push task",
		zap.Int64("task_id", msg.TaskID),
	)

	// 获取推送任务
	task, err := w.couponService.GetPushTask(ctx, msg.TaskID)
	if err != nil {
		logger.Error("Failed to get push task",
			zap.Int64("task_id", msg.TaskID),
			zap.Error(err),
		)
		d.Nack(false, true) // 拒绝消息，重新入队
		return
	}

	// 检查任务状态
	if task.Status != coupon.TaskStatusPending {
		logger.Warn("Task is not pending",
			zap.Int64("task_id", msg.TaskID),
			zap.Int8("status", task.Status),
		)
		d.Ack(false) // 确认消息
		return
	}

	// 更新任务状态为执行中
	now := time.Now()
	task.Status = coupon.TaskStatusProcessing
	task.StartedAt = &now
	if err := w.couponService.UpdatePushTask(ctx, task); err != nil {
		logger.Error("Failed to update task status",
			zap.Int64("task_id", msg.TaskID),
			zap.Error(err),
		)
		d.Nack(false, true)
		return
	}

	// 解析目标用户
	var targetUserIDs []int64
	if task.TargetType == coupon.TargetTypeSpecified {
		if err := json.Unmarshal([]byte(task.TargetUsers), &targetUserIDs); err != nil {
			logger.Error("Failed to unmarshal target users",
				zap.Int64("task_id", msg.TaskID),
				zap.Error(err),
			)
			w.markTaskFailed(ctx, task, err.Error())
			d.Ack(false)
			return
		}
	} else if task.TargetType == coupon.TargetTypeAll {
		// TODO: 获取所有用户ID（这里简化处理）
		logger.Warn("TargetTypeAll not fully implemented")
	}

	// 批量发放优惠券
	successCount := 0
	failCount := 0

	for _, userID := range targetUserIDs {
		// 发放优惠券
		_, err := w.couponService.IssueCouponToUser(ctx, userID, task.TemplateID)
		if err != nil {
			logger.Error("Failed to issue coupon",
				zap.Int64("task_id", msg.TaskID),
				zap.Int64("user_id", userID),
				zap.Error(err),
			)
			failCount++
			continue
		}

		successCount++

		// 如果需要发送微信消息
		if task.SendWechatMsg == 1 {
			// 发送微信模板消息（发送到通知队列）
			notifyMsg := WechatNotifyMessage{
				UserID:     userID,
				TemplateID: "coupon_received",
				Data: map[string]interface{}{
					"title":   task.Title,
					"content": "您有一张新的优惠券",
				},
			}
			if err := w.publishWechatNotify(notifyMsg); err != nil {
				logger.Error("Failed to publish wechat notify",
					zap.Int64("user_id", userID),
					zap.Error(err),
				)
			}
		}

		// 每100个打印一次进度
		if (successCount+failCount)%100 == 0 {
			logger.Info("Coupon push progress",
				zap.Int64("task_id", msg.TaskID),
				zap.Int("success_count", successCount),
				zap.Int("fail_count", failCount),
				zap.Int("total", task.TotalCount),
			)
		}
	}

	// 更新任务完成状态
	completedAt := time.Now()
	task.Status = coupon.TaskStatusCompleted
	task.CompletedAt = &completedAt
	task.SuccessCount = successCount
	task.FailCount = failCount

	if err := w.couponService.UpdatePushTask(ctx, task); err != nil {
		logger.Error("Failed to update task completion",
			zap.Int64("task_id", msg.TaskID),
			zap.Error(err),
		)
		d.Nack(false, true)
		return
	}

	logger.Info("Coupon push task completed",
		zap.Int64("task_id", msg.TaskID),
		zap.Int("success_count", successCount),
		zap.Int("fail_count", failCount),
	)

	// 确认消息
	d.Ack(false)
}

// markTaskFailed 标记任务失败
func (w *CouponPushWorker) markTaskFailed(ctx context.Context, task *coupon.CouponPushTask, reason string) {
	task.Status = coupon.TaskStatusCancelled
	if err := w.couponService.UpdatePushTask(ctx, task); err != nil {
		logger.Error("Failed to mark task as failed",
			zap.Int64("task_id", task.ID),
			zap.Error(err),
		)
	}
}

// publishWechatNotify 发布微信通知消息
func (w *CouponPushWorker) publishWechatNotify(msg WechatNotifyMessage) error {
	body, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	return rabbitmq.Publish("exchange.notify", "notify.wechat", body)
}
