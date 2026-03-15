package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	pkgErrors "github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    pkgErrors.CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 成功响应（自定义消息）
func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    pkgErrors.CodeSuccess,
		Message: message,
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, err error) {
	code := pkgErrors.GetCode(err)
	message := pkgErrors.GetMessage(err)

	statusCode := codeToHTTPStatus(code)

	c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithMessage 错误响应（自定义消息）
func ErrorWithMessage(c *gin.Context, code int, message string) {
	statusCode := codeToHTTPStatus(code)

	c.JSON(statusCode, Response{
		Code:    code,
		Message: message,
	})
}

// codeToHTTPStatus 业务错误码转HTTP状态码
func codeToHTTPStatus(code int) int {
	switch code {
	case pkgErrors.CodeSuccess:
		return http.StatusOK
	case pkgErrors.CodeInvalidParams:
		return http.StatusBadRequest
	case pkgErrors.CodeUnauthorized:
		return http.StatusUnauthorized
	case pkgErrors.CodeForbidden:
		return http.StatusForbidden
	case pkgErrors.CodeNotFound:
		return http.StatusNotFound
	case pkgErrors.CodeConflict:
		return http.StatusConflict
	case pkgErrors.CodeTooManyRequests:
		return http.StatusTooManyRequests
	case pkgErrors.CodeInternalError:
		return http.StatusInternalServerError
	case pkgErrors.CodeServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		if code >= 400 && code < 500 {
			return code
		}
		return http.StatusInternalServerError
	}
}

// PageData 分页数据
type PageData struct {
	List       interface{} `json:"list"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

// SuccessWithPage 分页成功响应
func SuccessWithPage(c *gin.Context, list interface{}, total int64, page, pageSize int) {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}

	c.JSON(http.StatusOK, Response{
		Code:    pkgErrors.CodeSuccess,
		Message: "success",
		Data: PageData{
			List:       list,
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	})
}
