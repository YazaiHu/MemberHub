package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	couponApp "github.com/YazaiHu/MemberHub/internal/application/coupon"
	"github.com/YazaiHu/MemberHub/internal/domain/coupon"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/middleware"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// CouponHandler 优惠券管理处理器
type CouponHandler struct {
	couponService *couponApp.Service
}

// NewCouponHandler 创建优惠券管理处理器
func NewCouponHandler() *CouponHandler {
	couponRepo := coupon.NewRepository()
	return &CouponHandler{
		couponService: couponApp.NewService(couponRepo),
	}
}

// ListTemplates 优惠券模板列表
// @Summary 优惠券模板列表
// @Tags 管理端-优惠券
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/coupons/templates [get]
func (h *CouponHandler) ListTemplates(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	templates, total, err := h.couponService.ListTemplates(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, templates, total, page, pageSize)
}

// GetTemplate 优惠券模板详情
// @Summary 优惠券模板详情
// @Tags 管理端-优惠券
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Success 200 {object} response.Response{data=coupon.CouponTemplate}
// @Router /api/admin/coupons/templates/{id} [get]
func (h *CouponHandler) GetTemplate(c *gin.Context) {
	templateID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	template, err := h.couponService.GetTemplate(c.Request.Context(), templateID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, template)
}

// CreateTemplateRequest 创建优惠券模板请求
type CreateTemplateRequest struct {
	Name             string  `json:"name" binding:"required"`
	Type             int8    `json:"type" binding:"required,oneof=1 2 3"`
	DiscountType     int8    `json:"discount_type" binding:"required,oneof=1 2"`
	DiscountValue    int     `json:"discount_value" binding:"required,min=1"`
	MinAmount        int     `json:"min_amount"`
	TotalQuantity    int     `json:"total_quantity" binding:"required,min=1"`
	PerUserLimit     int     `json:"per_user_limit" binding:"required,min=1"`
	ValidDays        int     `json:"valid_days" binding:"required,min=1"`
	ApplicableStores string  `json:"applicable_stores,omitempty"`
	Description      string  `json:"description,omitempty"`
	StartTime        *string `json:"start_time,omitempty"`
	EndTime          *string `json:"end_time,omitempty"`
	Status           int8    `json:"status" binding:"oneof=0 1"`
	SortOrder        int     `json:"sort_order"`
}

// CreateTemplate 创建优惠券模板
// @Summary 创建优惠券模板
// @Tags 管理端-优惠券
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateTemplateRequest true "优惠券模板信息"
// @Success 200 {object} response.Response{data=coupon.CouponTemplate}
// @Router /api/admin/coupons/templates [post]
func (h *CouponHandler) CreateTemplate(c *gin.Context) {
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	template := &coupon.CouponTemplate{
		Name:             req.Name,
		Type:             req.Type,
		DiscountType:     req.DiscountType,
		DiscountValue:    req.DiscountValue,
		MinAmount:        req.MinAmount,
		TotalQuantity:    req.TotalQuantity,
		PerUserLimit:     req.PerUserLimit,
		ValidDays:        req.ValidDays,
		ApplicableStores: req.ApplicableStores,
		Description:      req.Description,
		Status:           req.Status,
		SortOrder:        req.SortOrder,
	}

	// 解析时间
	if req.StartTime != nil && *req.StartTime != "" {
		// 时间解析逻辑由应用层或领域层处理
	}
	if req.EndTime != nil && *req.EndTime != "" {
		// 时间解析逻辑由应用层或领域层处理
	}

	if err := h.couponService.CreateTemplate(c.Request.Context(), template); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "template created successfully", template)
}

// UpdateTemplateRequest 更新优惠券模板请求
type UpdateTemplateRequest struct {
	ID               int64   `json:"id" binding:"required"`
	Name             string  `json:"name" binding:"required"`
	Type             int8    `json:"type" binding:"required,oneof=1 2 3"`
	DiscountType     int8    `json:"discount_type" binding:"required,oneof=1 2"`
	DiscountValue    int     `json:"discount_value" binding:"required,min=1"`
	MinAmount        int     `json:"min_amount"`
	TotalQuantity    int     `json:"total_quantity" binding:"required,min=1"`
	RemainingQuantity int    `json:"remaining_quantity" binding:"required,min=0"`
	PerUserLimit     int     `json:"per_user_limit" binding:"required,min=1"`
	ValidDays        int     `json:"valid_days" binding:"required,min=1"`
	ApplicableStores string  `json:"applicable_stores,omitempty"`
	Description      string  `json:"description,omitempty"`
	StartTime        *string `json:"start_time,omitempty"`
	EndTime          *string `json:"end_time,omitempty"`
	Status           int8    `json:"status" binding:"oneof=0 1"`
	SortOrder        int     `json:"sort_order"`
}

// UpdateTemplate 更新优惠券模板
// @Summary 更新优惠券模板
// @Tags 管理端-优惠券
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body UpdateTemplateRequest true "优惠券模板信息"
// @Success 200 {object} response.Response
// @Router /api/admin/coupons/templates [put]
func (h *CouponHandler) UpdateTemplate(c *gin.Context) {
	var req UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	template := &coupon.CouponTemplate{
		ID:               req.ID,
		Name:             req.Name,
		Type:             req.Type,
		DiscountType:     req.DiscountType,
		DiscountValue:    req.DiscountValue,
		MinAmount:        req.MinAmount,
		TotalQuantity:    req.TotalQuantity,
		RemainingQuantity: req.RemainingQuantity,
		PerUserLimit:     req.PerUserLimit,
		ValidDays:        req.ValidDays,
		ApplicableStores: req.ApplicableStores,
		Description:      req.Description,
		Status:           req.Status,
		SortOrder:        req.SortOrder,
	}

	if err := h.couponService.UpdateTemplate(c.Request.Context(), template); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "template updated successfully", nil)
}

// CreatePushTaskRequest 创建推送任务请求
type CreatePushTaskRequest struct {
	TemplateID    int64   `json:"template_id" binding:"required"`
	Title         string  `json:"title" binding:"required"`
	TargetType    int8    `json:"target_type" binding:"required,oneof=1 2 3"`
	TargetUserIDs []int64 `json:"target_user_ids,omitempty"`
	SendWechatMsg bool    `json:"send_wechat_msg"`
}

// CreatePushTask 创建推送任务
// @Summary 创建推送任务
// @Tags 管理端-优惠券
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreatePushTaskRequest true "推送任务信息"
// @Success 200 {object} response.Response{data=coupon.CouponPushTask}
// @Router /api/admin/coupons/push [post]
func (h *CouponHandler) CreatePushTask(c *gin.Context) {
	adminID := middleware.GetUserID(c)

	var req CreatePushTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	pushReq := &coupon.PushCouponRequest{
		TemplateID:    req.TemplateID,
		Title:         req.Title,
		TargetType:    req.TargetType,
		TargetUserIDs: req.TargetUserIDs,
		SendWechatMsg: req.SendWechatMsg,
	}

	task, err := h.couponService.CreatePushTask(c.Request.Context(), pushReq, adminID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "push task created successfully", task)
}

// ListPushTasks 推送任务列表
// @Summary 推送任务列表
// @Tags 管理端-优惠券
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/coupons/push-tasks [get]
func (h *CouponHandler) ListPushTasks(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	tasks, total, err := h.couponService.ListPushTasks(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, tasks, total, page, pageSize)
}

// GetPushTask 推送任务详情
// @Summary 推送任务详情
// @Tags 管理端-优惠券
// @Produce json
// @Security BearerAuth
// @Param id path int true "任务ID"
// @Success 200 {object} response.Response{data=coupon.CouponPushTask}
// @Router /api/admin/coupons/push-tasks/{id} [get]
func (h *CouponHandler) GetPushTask(c *gin.Context) {
	taskID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	task, err := h.couponService.GetPushTask(c.Request.Context(), taskID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, task)
}
