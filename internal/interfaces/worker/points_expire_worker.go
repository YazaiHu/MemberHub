package worker

import (
	"context"
	"time"

	"go.uber.org/zap"

	pointsApp "github.com/YazaiHu/MemberHub/internal/application/points"
	"github.com/YazaiHu/MemberHub/internal/domain/points"
	"github.com/YazaiHu/MemberHub/internal/pkg/logger"
)

// PointsExpireWorker 积分过期Worker
type PointsExpireWorker struct {
	pointsService *pointsApp.Service
}

// NewPointsExpireWorker 创建积分过期Worker
func NewPointsExpireWorker() *PointsExpireWorker {
	pointsRepo := points.NewRepository()
	return &PointsExpireWorker{
		pointsService: pointsApp.NewService(pointsRepo),
	}
}

// Start 启动Worker（定时任务，每天凌晨2点执行）
func (w *PointsExpireWorker) Start() error {
	logger.Info("Starting PointsExpireWorker")

	// 计算下次执行时间（明天凌晨2点）
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 2, 0, 0, 0, now.Location())
	duration := next.Sub(now)

	logger.Info("PointsExpireWorker scheduled",
		zap.String("next_run", next.Format("2006-01-02 15:04:05")),
	)

	// 首次延迟
	time.Sleep(duration)

	// 每24小时执行一次
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()

	for {
		w.processExpiredPoints()
		<-ticker.C
	}
}

// processExpiredPoints 处理过期积分
func (w *PointsExpireWorker) processExpiredPoints() {
	ctx := context.Background()

	logger.Info("Processing expired points")

	// TODO: 实现积分过期逻辑
	// 1. 查询所有过期的积分
	// 2. 批量标记为已过期
	// 3. 记录过期流水

	// 这里简化处理，实际应该在 points domain 中实现 MarkExpiredPoints 方法
	_ = ctx // Suppress unused warning for now
	logger.Info("Expired points processed (placeholder)")
}
