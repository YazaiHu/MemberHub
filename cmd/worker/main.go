package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/messaging/rabbitmq"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/redis"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/wechat"
	"github.com/YazaiHu/MemberHub/internal/interfaces/worker"
	"github.com/YazaiHu/MemberHub/internal/pkg/config"
	"github.com/YazaiHu/MemberHub/internal/pkg/logger"
)

func main() {
	// 加载配置
	configPath := "configs/config.dev.yaml"
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		configPath = path
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化日志
	if err := logger.InitLogger(&cfg.Log); err != nil {
		log.Fatal("Failed to init logger:", err)
	}
	defer logger.Sync()

	// 初始化数据库
	if err := mysql.InitDB(&cfg.Database); err != nil {
		logger.Fatal("Failed to init database", zap.Error(err))
	}
	defer mysql.Close()

	// 初始化Redis
	if err := redis.InitRedis(&cfg.Redis); err != nil {
		logger.Fatal("Failed to init redis", zap.Error(err))
	}
	defer redis.Close()

	// 初始化RabbitMQ
	if err := rabbitmq.InitRabbitMQ(&cfg.RabbitMQ); err != nil {
		logger.Fatal("Failed to init rabbitmq", zap.Error(err))
	}
	defer rabbitmq.Close()

	// 初始化微信SDK
	if err := wechat.InitWeChat(&cfg.WeChat); err != nil {
		logger.Fatal("Failed to init wechat", zap.Error(err))
	}

	logger.Info("Worker service starting...")

	// 启动所有Worker
	errChan := make(chan error, 3)

	// 1. 优惠券推送Worker
	couponPushWorker := worker.NewCouponPushWorker()
	go func() {
		logger.Info("Starting CouponPushWorker")
		if err := couponPushWorker.Start(); err != nil {
			errChan <- fmt.Errorf("CouponPushWorker error: %w", err)
		}
	}()

	// 2. 微信通知Worker
	wechatNotifyWorker := worker.NewWechatNotifyWorker()
	go func() {
		logger.Info("Starting WechatNotifyWorker")
		if err := wechatNotifyWorker.Start(); err != nil {
			errChan <- fmt.Errorf("WechatNotifyWorker error: %w", err)
		}
	}()

	// 3. 积分过期Worker（定时任务）
	pointsExpireWorker := worker.NewPointsExpireWorker()
	go func() {
		logger.Info("Starting PointsExpireWorker")
		if err := pointsExpireWorker.Start(); err != nil {
			errChan <- fmt.Errorf("PointsExpireWorker error: %w", err)
		}
	}()

	logger.Info("All workers started successfully")

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		logger.Error("Worker error occurred", zap.Error(err))
	case sig := <-quit:
		logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
	}

	// 停止Worker
	logger.Info("Shutting down workers...")
	wechatNotifyWorker.Stop()

	logger.Info("Workers stopped")
}
