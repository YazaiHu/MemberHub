package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/messaging/rabbitmq"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/redis"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/wechat"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/handler"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/handler/admin"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/handler/miniapp"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/middleware"
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

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 创建路由
	router := setupRouter()

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// 启动服务器
	go func() {
		logger.Info("Server starting", zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

func setupRouter() *gin.Engine {
	router := gin.New()

	// 全局中间件
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})

	// 初始化处理器
	miniappAuthHandler := miniapp.NewAuthHandler()
	miniappMemberHandler := miniapp.NewMemberHandler()
	miniappPointsHandler := miniapp.NewPointsHandler()
	miniappRechargeHandler := miniapp.NewRechargeHandler()
	miniappCouponHandler := miniapp.NewCouponHandler()
	miniappStoreHandler := miniapp.NewStoreHandler()
	miniappPromotionHandler := miniapp.NewPromotionHandler()

	adminAuthHandler := admin.NewAuthHandler()
	adminMemberHandler := admin.NewMemberHandler()
	adminPointsHandler := admin.NewPointsHandler()
	adminRechargeHandler := admin.NewRechargeHandler()
	adminCouponHandler := admin.NewCouponHandler()
	adminStoreHandler := admin.NewStoreHandler()
	adminPromotionHandler := admin.NewPromotionHandler()

	callbackHandler := handler.NewCallbackHandler()

	// API路由组
	api := router.Group("/api")
	{
		// 小程序端接口
		miniappGroup := api.Group("/miniapp")
		{
			// 认证相关（无需token）
			auth := miniappGroup.Group("/auth")
			{
				auth.POST("/login", miniappAuthHandler.Login)
				auth.POST("/refresh", miniappAuthHandler.RefreshToken)
			}

			// 需要认证的接口
			authenticated := miniappGroup.Group("", middleware.MemberAuth())
			{
				// 会员中心
				authenticated.GET("/profile", miniappMemberHandler.GetProfile)
				authenticated.PUT("/profile", miniappMemberHandler.UpdateProfile)
				authenticated.POST("/profile/bind-phone", miniappMemberHandler.BindPhone)
				authenticated.GET("/asset", miniappMemberHandler.GetAsset)

				// 积分
				points := authenticated.Group("/points")
				{
					points.GET("/balance", miniappPointsHandler.GetBalance)
					points.GET("/history", miniappPointsHandler.GetHistory)
					points.GET("/exchange/rules", miniappPointsHandler.ListExchangeRules)
					points.POST("/exchange", miniappPointsHandler.Exchange)
					points.GET("/exchange/records", miniappPointsHandler.GetExchangeRecords)
				}

				// 充值
				rechargeGroup := authenticated.Group("/recharge")
				{
					rechargeGroup.GET("/promotions", miniappRechargeHandler.GetPromotions)
					rechargeGroup.POST("/create-order", miniappRechargeHandler.CreateOrder)
					rechargeGroup.GET("/balance", miniappRechargeHandler.GetBalance)
					rechargeGroup.GET("/balance/history", miniappRechargeHandler.GetBalanceHistory)
					rechargeGroup.GET("/records", miniappRechargeHandler.GetRechargeRecords)
				}

				// 优惠券
				coupons := authenticated.Group("/coupons")
				{
					coupons.GET("/templates", miniappCouponHandler.ListTemplates)
					coupons.POST("/receive", miniappCouponHandler.ReceiveCoupon)
					coupons.GET("/available", miniappCouponHandler.GetAvailableCoupons)
					coupons.GET("/all", miniappCouponHandler.GetAllCoupons)
				}

				// 门店
				stores := authenticated.Group("/stores")
				{
					stores.GET("", miniappStoreHandler.ListStores)
					stores.POST("/nearby", miniappStoreHandler.GetNearbyStores)
				}

				// 特价商品
				promotions := authenticated.Group("/promotions")
				{
					promotions.GET("/weekly", miniappPromotionHandler.GetWeeklyPromotions)
					promotions.GET("", miniappPromotionHandler.GetActivePromotions)
				}
			}
		}

		// 管理后台接口
		adminGroup := api.Group("/admin")
		{
			// 管理员登录（无需token）
			authGroup := adminGroup.Group("/auth")
			{
				authGroup.POST("/login", adminAuthHandler.Login)
				authGroup.POST("/refresh", adminAuthHandler.RefreshToken)
			}

			// 需要管理员认证的接口
			authenticated := adminGroup.Group("", middleware.AdminAuth())
			{
				// 会员管理
				members := authenticated.Group("/members")
				{
					members.GET("", adminMemberHandler.ListMembers)
					members.GET("/:id", adminMemberHandler.GetMember)
					members.GET("/:id/asset", adminMemberHandler.GetMemberAsset)
					members.POST("/:id/disable", adminMemberHandler.DisableMember)
					members.POST("/:id/enable", adminMemberHandler.EnableMember)
				}

				// 门店管理
				stores := authenticated.Group("/stores")
				{
					stores.GET("", adminStoreHandler.ListStores)
					stores.GET("/:id", adminStoreHandler.GetStore)
					stores.POST("", adminStoreHandler.CreateStore)
					stores.PUT("", adminStoreHandler.UpdateStore)
					stores.DELETE("/:id", adminStoreHandler.DeleteStore)
				}

				// 积分管理
				pointsGroup := authenticated.Group("/points")
				{
					pointsGroup.GET("/rules", adminPointsHandler.ListExchangeRules)
					pointsGroup.POST("/rules", adminPointsHandler.CreateRule)
					pointsGroup.PUT("/rules/:id", adminPointsHandler.UpdateRule)
					pointsGroup.POST("/exchange/verify", adminPointsHandler.VerifyExchange)
					pointsGroup.GET("/exchange/records", adminPointsHandler.ListExchangeRecords)
				}

				// 会员积分和余额调整
				members.POST("/:id/adjust-points", adminPointsHandler.AdjustPoints)
				members.POST("/:id/adjust-balance", adminRechargeHandler.AdjustBalance)

				// 充值管理
				rechargeGroup := authenticated.Group("/recharge")
				{
					rechargeGroup.GET("/promotions", adminRechargeHandler.ListPromotions)
					rechargeGroup.POST("/promotions", adminRechargeHandler.CreatePromotion)
					rechargeGroup.PUT("/promotions/:id", adminRechargeHandler.UpdatePromotion)
					rechargeGroup.GET("/orders", adminRechargeHandler.ListOrders)
				}

				// 优惠券管理
				couponGroup := authenticated.Group("/coupons")
				{
					couponGroup.GET("/templates", adminCouponHandler.ListTemplates)
					couponGroup.GET("/templates/:id", adminCouponHandler.GetTemplate)
					couponGroup.POST("/templates", adminCouponHandler.CreateTemplate)
					couponGroup.PUT("/templates", adminCouponHandler.UpdateTemplate)
					couponGroup.POST("/push", adminCouponHandler.CreatePushTask)
					couponGroup.GET("/push-tasks", adminCouponHandler.ListPushTasks)
					couponGroup.GET("/push-tasks/:id", adminCouponHandler.GetPushTask)
				}

				// 特价商品管理
				promotionGroup := authenticated.Group("/promotions")
				{
					promotionGroup.GET("", adminPromotionHandler.ListProducts)
					promotionGroup.GET("/:id", adminPromotionHandler.GetProduct)
					promotionGroup.POST("", adminPromotionHandler.CreateProduct)
					promotionGroup.PUT("", adminPromotionHandler.UpdateProduct)
					promotionGroup.DELETE("/:id", adminPromotionHandler.DeleteProduct)
				}
			}
		}

		// 支付回调
		callback := api.Group("/callback")
		{
			callback.POST("/wechat-pay", callbackHandler.WeChatPayCallback)
		}
	}

	return router
}
