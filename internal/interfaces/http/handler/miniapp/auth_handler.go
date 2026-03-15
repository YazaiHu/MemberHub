package miniapp

import (
	"github.com/gin-gonic/gin"

	"github.com/YazaiHu/MemberHub/internal/application/auth"
	memberApp "github.com/YazaiHu/MemberHub/internal/application/member"
	"github.com/YazaiHu/MemberHub/internal/domain/member"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/middleware"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// AuthHandler 小程序认证处理器
type AuthHandler struct {
	authService   *auth.Service
	memberService *memberApp.Service
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	memberRepo := member.NewRepository()
	return &AuthHandler{
		authService:   auth.NewService(memberRepo),
		memberService: memberApp.NewService(memberRepo),
	}
}

// Login 微信小程序登录
// @Summary 微信小程序登录
// @Tags 小程序-认证
// @Accept json
// @Produce json
// @Param body body auth.WeChatLoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=auth.WeChatLoginResponse}
// @Router /api/miniapp/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.WeChatLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	result, err := h.authService.WeChatLogin(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// RefreshToken 刷新token
// @Summary 刷新token
// @Tags 小程序-认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/miniapp/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// 从请求头获取旧token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, errors.ErrUnauthorized)
		return
	}

	// 去掉"Bearer "前缀
	token := authHeader[7:]

	newToken, err := h.authService.RefreshToken(c.Request.Context(), token)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"token": newToken,
	})
}

// MemberHandler 会员处理器
type MemberHandler struct {
	memberService *memberApp.Service
}

// NewMemberHandler 创建会员处理器
func NewMemberHandler() *MemberHandler {
	memberRepo := member.NewRepository()
	return &MemberHandler{
		memberService: memberApp.NewService(memberRepo),
	}
}

// GetProfile 获取个人信息
// @Summary 获取个人信息
// @Tags 小程序-会员
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=member.User}
// @Router /api/miniapp/profile [get]
func (h *MemberHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	user, err := h.memberService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

// UpdateProfileRequest 更新个人信息请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int8   `json:"gender"`
	Birthday string `json:"birthday"` // YYYY-MM-DD
}

// UpdateProfile 更新个人信息
// @Summary 更新个人信息
// @Tags 小程序-会员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body UpdateProfileRequest true "个人信息"
// @Success 200 {object} response.Response
// @Router /api/miniapp/profile [put]
func (h *MemberHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	// 构造用户对象
	user := &member.User{
		ID:       userID,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
	}

	// TODO: 解析生日字段

	if err := h.memberService.UpdateUserInfo(c.Request.Context(), user); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// BindPhoneRequest 绑定手机号请求
type BindPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
}

// BindPhone 绑定手机号
// @Summary 绑定手机号
// @Tags 小程序-会员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body BindPhoneRequest true "手机号"
// @Success 200 {object} response.Response
// @Router /api/miniapp/profile/bind-phone [post]
func (h *MemberHandler) BindPhone(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req BindPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	if err := h.memberService.BindPhone(c.Request.Context(), userID, req.Phone); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// GetAsset 获取资产概览
// @Summary 获取资产概览
// @Tags 小程序-会员
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=member.UserAsset}
// @Router /api/miniapp/asset [get]
func (h *MemberHandler) GetAsset(c *gin.Context) {
	userID := middleware.GetUserID(c)

	asset, err := h.memberService.GetUserAsset(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, asset)
}
