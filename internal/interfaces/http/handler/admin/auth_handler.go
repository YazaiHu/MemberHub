package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/YazaiHu/MemberHub/internal/application/auth"
	memberApp "github.com/YazaiHu/MemberHub/internal/application/member"
	"github.com/YazaiHu/MemberHub/internal/domain/member"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// AuthHandler 管理员认证处理器
type AuthHandler struct {
	authService *auth.Service
}

// NewAuthHandler 创建管理员认证处理器
func NewAuthHandler() *AuthHandler {
	memberRepo := member.NewRepository()
	return &AuthHandler{
		authService: auth.NewService(memberRepo),
	}
}

// Login 管理员登录
// @Summary 管理员登录
// @Tags 管理后台-认证
// @Accept json
// @Produce json
// @Param body body auth.AdminLoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=auth.AdminLoginResponse}
// @Router /api/admin/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	// 获取客户端IP
	ip := c.ClientIP()

	result, err := h.authService.AdminLogin(c.Request.Context(), &req, ip)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, result)
}

// RefreshToken 刷新token
// @Summary 刷新token
// @Tags 管理后台-认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=map[string]string}
// @Router /api/admin/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, errors.ErrUnauthorized)
		return
	}

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

// MemberHandler 会员管理处理器
type MemberHandler struct {
	memberService *memberApp.Service
}

// NewMemberHandler 创建会员管理处理器
func NewMemberHandler() *MemberHandler {
	memberRepo := member.NewRepository()
	return &MemberHandler{
		memberService: memberApp.NewService(memberRepo),
	}
}

// ListMembers 会员列表
// @Summary 会员列表
// @Tags 管理后台-会员
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param status query int false "状态: 0-禁用, 1-正常"
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/members [get]
func (h *MemberHandler) ListMembers(c *gin.Context) {
	// 获取分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// 获取状态参数
	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = new(int8)
		*status = int8(s)
	}

	users, total, err := h.memberService.ListUsers(c.Request.Context(), page, pageSize, status)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, users, total, page, pageSize)
}

// GetMember 获取会员详情
// @Summary 获取会员详情
// @Tags 管理后台-会员
// @Produce json
// @Security BearerAuth
// @Param id path int true "会员ID"
// @Success 200 {object} response.Response{data=member.User}
// @Router /api/admin/members/{id} [get]
func (h *MemberHandler) GetMember(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	user, err := h.memberService.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, user)
}

// GetMemberAsset 获取会员资产
// @Summary 获取会员资产
// @Tags 管理后台-会员
// @Produce json
// @Security BearerAuth
// @Param id path int true "会员ID"
// @Success 200 {object} response.Response{data=member.UserAsset}
// @Router /api/admin/members/{id}/asset [get]
func (h *MemberHandler) GetMemberAsset(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	asset, err := h.memberService.GetUserAsset(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, asset)
}

// DisableMember 禁用会员
// @Summary 禁用会员
// @Tags 管理后台-会员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会员ID"
// @Success 200 {object} response.Response
// @Router /api/admin/members/{id}/disable [post]
func (h *MemberHandler) DisableMember(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	if err := h.memberService.DisableUser(c.Request.Context(), userID); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "disabled successfully", nil)
}

// EnableMember 启用会员
// @Summary 启用会员
// @Tags 管理后台-会员
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会员ID"
// @Success 200 {object} response.Response
// @Router /api/admin/members/{id}/enable [post]
func (h *MemberHandler) EnableMember(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	if err := h.memberService.EnableUser(c.Request.Context(), userID); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "enabled successfully", nil)
}
