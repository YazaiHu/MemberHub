package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/YazaiHu/MemberHub/internal/pkg/config"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// Claims JWT载荷
type Claims struct {
	UserID   int64  `json:"user_id"`
	UserType string `json:"user_type"` // member, admin
	StoreID  int64  `json:"store_id,omitempty"` // 店铺管理员的店铺ID
	RoleCode string `json:"role_code,omitempty"` // 角色代码
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID int64, userType string, storeID int64, roleCode string) (string, error) {
	cfg := config.GetConfig()

	// 根据用户类型选择过期时间
	expireTime := cfg.JWT.ExpireTime
	if userType == "admin" {
		expireTime = cfg.JWT.AdminExpireTime
	}

	now := time.Now()
	claims := &Claims{
		UserID:   userID,
		UserType: userType,
		StoreID:  storeID,
		RoleCode: roleCode,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    cfg.JWT.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expireTime)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWT.Secret))
}

// ParseToken 解析JWT token
func ParseToken(tokenString string) (*Claims, error) {
	cfg := config.GetConfig()

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(cfg.JWT.Secret), nil
	})

	if err != nil {
		return nil, errors.ErrInvalidToken.Wrap(err)
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrInvalidToken
}

// RefreshToken 刷新token
func RefreshToken(tokenString string) (string, error) {
	claims, err := ParseToken(tokenString)
	if err != nil {
		return "", err
	}

	// 生成新token
	return GenerateToken(claims.UserID, claims.UserType, claims.StoreID, claims.RoleCode)
}
