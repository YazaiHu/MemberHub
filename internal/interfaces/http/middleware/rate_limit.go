package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/redis"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// RateLimit 限流中间件
// limit: 时间窗口内允许的请求数
// window: 时间窗口（秒）
func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用IP作为限流key
		key := fmt.Sprintf("rate_limit:%s:%s", c.Request.URL.Path, c.ClientIP())

		// 检查限流
		allowed, err := redis.RateLimiter(key, limit, window)
		if err != nil {
			// 限流失败不影响业务
			c.Next()
			return
		}

		if !allowed {
			response.Error(c, errors.ErrTooManyRequests)
			c.Abort()
			return
		}

		c.Next()
	}
}

// RateLimitByUser 按用户限流
func RateLimitByUser(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			c.Next()
			return
		}

		// 使用用户ID作为限流key
		key := fmt.Sprintf("rate_limit:%s:user:%d", c.Request.URL.Path, userID)

		// 检查限流
		allowed, err := redis.RateLimiter(key, limit, window)
		if err != nil {
			c.Next()
			return
		}

		if !allowed {
			response.Error(c, errors.ErrTooManyRequests)
			c.Abort()
			return
		}

		c.Next()
	}
}
