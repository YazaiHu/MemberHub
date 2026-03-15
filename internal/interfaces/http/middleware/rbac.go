package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// RequireRole 角色权限中间件
func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCode := GetRoleCode(c)

		// 超级管理员拥有所有权限
		if roleCode == "super_admin" {
			c.Next()
			return
		}

		// 检查角色是否在允许列表中
		for _, role := range allowedRoles {
			if roleCode == role {
				c.Next()
				return
			}
		}

		response.Error(c, errors.ErrForbidden.WithMessage("insufficient permissions"))
		c.Abort()
	}
}

// RequireStoreAccess 店铺访问权限中间件
// 店铺管理员只能访问自己的店铺数据
func RequireStoreAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		roleCode := GetRoleCode(c)

		// 超级管理员可以访问所有店铺
		if roleCode == "super_admin" {
			c.Next()
			return
		}

		// 店铺管理员只能访问自己的店铺
		if roleCode == "store_admin" {
			storeID := GetStoreID(c)
			if storeID == 0 {
				response.Error(c, errors.ErrForbidden.WithMessage("no store assigned"))
				c.Abort()
				return
			}

			// 从请求参数中获取目标店铺ID
			targetStoreID := c.GetInt64("store_id")
			if targetStoreID == 0 {
				// 如果没有指定店铺ID，默认使用管理员的店铺ID
				c.Set("store_id", storeID)
				c.Next()
				return
			}

			// 检查是否访问自己的店铺
			if targetStoreID != storeID {
				response.Error(c, errors.ErrForbidden.WithMessage("cannot access other store data"))
				c.Abort()
				return
			}
		}

		c.Next()
	}
}
