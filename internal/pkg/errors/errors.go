package errors

import (
	"errors"
	"fmt"
)

// 业务错误码
const (
	CodeSuccess           = 0
	CodeInvalidParams     = 400
	CodeUnauthorized      = 401
	CodeForbidden         = 403
	CodeNotFound          = 404
	CodeConflict          = 409
	CodeTooManyRequests   = 429
	CodeInternalError     = 500
	CodeServiceUnavailable = 503
)

// 业务错误
var (
	// 通用错误
	ErrInvalidParams     = New(CodeInvalidParams, "invalid parameters")
	ErrUnauthorized      = New(CodeUnauthorized, "unauthorized")
	ErrForbidden         = New(CodeForbidden, "forbidden")
	ErrNotFound          = New(CodeNotFound, "resource not found")
	ErrConflict          = New(CodeConflict, "resource conflict")
	ErrTooManyRequests   = New(CodeTooManyRequests, "too many requests")
	ErrInternalError     = New(CodeInternalError, "internal server error")
	ErrServiceUnavailable = New(CodeServiceUnavailable, "service unavailable")

	// 认证相关
	ErrInvalidToken      = New(CodeUnauthorized, "invalid token")
	ErrTokenExpired      = New(CodeUnauthorized, "token expired")
	ErrInvalidCredentials = New(CodeUnauthorized, "invalid credentials")

	// 会员相关
	ErrUserNotFound      = New(CodeNotFound, "user not found")
	ErrUserDisabled      = New(CodeForbidden, "user is disabled")
	ErrUserExists        = New(CodeConflict, "user already exists")

	// 积分相关
	ErrInsufficientPoints = New(CodeInvalidParams, "insufficient points")
	ErrPointsRuleNotFound = New(CodeNotFound, "points exchange rule not found")
	ErrStockNotEnough     = New(CodeConflict, "stock not enough")
	ErrExchangeLimitExceeded = New(CodeConflict, "exchange limit exceeded")
	ErrPointsExpired      = New(CodeInvalidParams, "points expired")

	// 充值相关
	ErrOrderNotFound      = New(CodeNotFound, "order not found")
	ErrOrderPaid          = New(CodeConflict, "order already paid")
	ErrOrderExpired       = New(CodeInvalidParams, "order expired")
	ErrInvalidAmount      = New(CodeInvalidParams, "invalid amount")
	ErrPromotionNotFound  = New(CodeNotFound, "promotion not found")
	ErrPromotionExpired   = New(CodeInvalidParams, "promotion expired")

	// 优惠券相关
	ErrCouponNotFound     = New(CodeNotFound, "coupon not found")
	ErrCouponExpired      = New(CodeInvalidParams, "coupon expired")
	ErrCouponUsed         = New(CodeConflict, "coupon already used")
	ErrCouponSoldOut      = New(CodeConflict, "coupon sold out")
	ErrCouponLimitExceeded = New(CodeConflict, "coupon receive limit exceeded")

	// 门店相关
	ErrStoreNotFound      = New(CodeNotFound, "store not found")
	ErrStoreDisabled      = New(CodeForbidden, "store is disabled")

	// 并发控制
	ErrConcurrentUpdate   = New(CodeConflict, "concurrent update conflict, please retry")
	ErrDuplicateOperation = New(CodeConflict, "duplicate operation")

	// 微信相关
	ErrWeChatAPIFailed    = New(CodeInternalError, "wechat api call failed")
	ErrInvalidSignature   = New(CodeInvalidParams, "invalid signature")
)

// AppError 应用错误
type AppError struct {
	Code    int
	Message string
	Err     error
}

// New 创建新的应用错误
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Error 实现 error 接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Wrap 包装底层错误
func (e *AppError) Wrap(err error) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Err:     err,
	}
}

// WithMessage 添加详细消息
func (e *AppError) WithMessage(message string) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: message,
		Err:     e.Err,
	}
}

// Is 判断错误是否相同
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	return e.Code == t.Code
}

// Unwrap 解包错误
func (e *AppError) Unwrap() error {
	return e.Err
}

// IsAppError 判断是否为应用错误
func IsAppError(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr)
}

// GetCode 获取错误码
func GetCode(err error) int {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Code
	}
	return CodeInternalError
}

// GetMessage 获取错误消息
func GetMessage(err error) string {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Message
	}
	return err.Error()
}
