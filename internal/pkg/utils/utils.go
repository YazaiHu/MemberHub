package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 使用 bcrypt 加密密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword 验证密码
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateUUID 生成UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateOrderNo 生成订单号
func GenerateOrderNo(prefix string) string {
	return fmt.Sprintf("%s%s%06d", prefix, time.Now().Format("20060102150405"), RandomInt(999999))
}

// GenerateCouponCode 生成优惠券码
func GenerateCouponCode() string {
	timestamp := time.Now().Unix()
	random := RandomInt(9999)
	str := fmt.Sprintf("%d%04d", timestamp, random)
	return MD5(str)[:12]
}

// MD5 计算MD5哈希
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// RandomInt 生成随机整数 [0, max]
func RandomInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max+1)))
	return int(n.Int64())
}

// RandomString 生成随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

// InArray 判断元素是否在数组中
func InArray[T comparable](item T, array []T) bool {
	for _, v := range array {
		if v == item {
			return true
		}
	}
	return false
}

// Pointer 返回值的指针
func Pointer[T any](v T) *T {
	return &v
}

// Deref 解引用指针，如果为nil则返回零值
func Deref[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}
	return *v
}

// Min 返回最小值
func Min[T int | int64 | float64](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max 返回最大值
func Max[T int | int64 | float64](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// CalculatePoints 计算消费积分
// amount: 消费金额（分）
// ratio: 积分比例（每元返多少积分）
func CalculatePoints(amount int64, ratio float64) int64 {
	// 转换为元
	amountYuan := float64(amount) / 100.0
	// 计算积分
	points := amountYuan * ratio
	return int64(points)
}

// GetStartOfDay 获取当天开始时间
func GetStartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// GetEndOfDay 获取当天结束时间
func GetEndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999999999, t.Location())
}
