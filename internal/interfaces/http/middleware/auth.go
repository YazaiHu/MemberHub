package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
	"github.com/YazaiHu/MemberHub/internal/pkg/utils"
)

// Context Key常量
const (
	ContextKeyUserID   = "user_id"
	ContextKeyUserType = "user_type"
	ContextKeyStoreID  = "store_id"
	ContextKeyRoleCode = "role_code"
	ContextKeyClaims   = "claims"
)

// Auth JWT认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, errors.ErrUnauthorized)
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, errors.ErrInvalidToken)
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			response.Error(c, err)
			c.Abort()
			return
		}

		// 将用户信息存入context
		c.Set(ContextKeyUserID, claims.UserID)
		c.Set(ContextKeyUserType, claims.UserType)
		c.Set(ContextKeyStoreID, claims.StoreID)
		c.Set(ContextKeyRoleCode, claims.RoleCode)
		c.Set(ContextKeyClaims, claims)

		c.Next()
	}
}

// MemberAuth 会员认证中间件
func MemberAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		Auth()(c)
		if c.IsAborted() {
			return
		}

		userType, _ := c.Get(ContextKeyUserType)
		if userType != "member" {
			response.Error(c, errors.ErrForbidden.WithMessage("member access required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminAuth 管理员认证中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		Auth()(c)
		if c.IsAborted() {
			return
		}

		userType, _ := c.Get(ContextKeyUserType)
		if userType != "admin" {
			response.Error(c, errors.ErrForbidden.WithMessage("admin access required"))
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserID 从context获取用户ID
func GetUserID(c *gin.Context) int64 {
	userID, exists := c.Get(ContextKeyUserID)
	if !exists {
		return 0
	}
	return userID.(int64)
}

// GetUserType 从context获取用户类型
func GetUserType(c *gin.Context) string {
	userType, exists := c.Get(ContextKeyUserType)
	if !exists {
		return ""
	}
	return userType.(string)
}

// GetStoreID 从context获取店铺ID
func GetStoreID(c *gin.Context) int64 {
	storeID, exists := c.Get(ContextKeyStoreID)
	if !exists {
		return 0
	}
	return storeID.(int64)
}

// GetRoleCode 从context获取角色代码
func GetRoleCode(c *gin.Context) string {
	roleCode, exists := c.Get(ContextKeyRoleCode)
	if !exists {
		return ""
	}
	return roleCode.(string)
}

// GetClaims 从context获取Claims
func GetClaims(c *gin.Context) *utils.Claims {
	claims, exists := c.Get(ContextKeyClaims)
	if !exists {
		return nil
	}
	return claims.(*utils.Claims)
}
