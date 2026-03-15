package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	rechargeApp "github.com/YazaiHu/MemberHub/internal/application/recharge"
	"github.com/YazaiHu/MemberHub/internal/domain/points"
	"github.com/YazaiHu/MemberHub/internal/domain/recharge"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/middleware"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// RechargeHandler 充值管理处理器
type RechargeHandler struct {
	rechargeService *rechargeApp.Service
}

// NewRechargeHandler 创建充值管理处理器
func NewRechargeHandler() *RechargeHandler {
	rechargeRepo := recharge.NewRepository()
	pointsRepo := points.NewRepository()
	return &RechargeHandler{
		rechargeService: rechargeApp.NewService(rechargeRepo, pointsRepo),
	}
}

// AdjustBalanceRequest 调整余额请求
type AdjustBalanceRequest struct {
	Amount int64  `json:"amount" binding:"required"` // 正数增加，负数扣减（单位：分）
	Remark string `json:"remark" binding:"required"`
}

// AdjustBalance 手动调整余额
// @Summary 手动调整余额
// @Tags 管理后台-充值
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会员ID"
// @Param body body AdjustBalanceRequest true "调整信息"
// @Success 200 {object} response.Response
// @Router /api/admin/members/{id}/adjust-balance [post]
func (h *RechargeHandler) AdjustBalance(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	var req AdjustBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	if err := h.rechargeService.ManualAdjustBalance(c.Request.Context(), userID, req.Amount, req.Remark, adminID); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "balance adjusted successfully", nil)
}

// CreatePromotionRequest 创建充值活动请求
type CreatePromotionRequest struct {
	Title          string  `json:"title" binding:"required"`
	Description    string  `json:"description"`
	RechargeAmount int     `json:"recharge_amount" binding:"required,gt=0"` // 充值金额（元）
	BonusAmount    int     `json:"bonus_amount" binding:"gte=0"`            // 赠送金额（元）
	BonusPoints    int     `json:"bonus_points" binding:"gte=0"`            // 赠送积分
	StartTime      *string `json:"start_time"`
	EndTime        *string `json:"end_time"`
	UserLimit      *int    `json:"user_limit"`
	TotalLimit     *int    `json:"total_limit"`
	Status         int8    `json:"status" binding:"oneof=0 1"`
	SortOrder      int     `json:"sort_order"`
}

// CreatePromotion 创建充值活动
// @Summary 创建充值活动
// @Tags 管理后台-充值
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreatePromotionRequest true "活动信息"
// @Success 200 {object} response.Response{data=recharge.RechargePromotion}
// @Router /api/admin/recharge/promotions [post]
func (h *RechargeHandler) CreatePromotion(c *gin.Context) {
	var req CreatePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	promotion := &recharge.RechargePromotion{
		Title:          req.Title,
		Description:    req.Description,
		RechargeAmount: req.RechargeAmount,
		BonusAmount:    req.BonusAmount,
		BonusPoints:    req.BonusPoints,
		UserLimit:      req.UserLimit,
		TotalLimit:     req.TotalLimit,
		Status:         req.Status,
		SortOrder:      req.SortOrder,
	}

	// TODO: 解析时间字段

	if err := h.rechargeService.CreatePromotion(c.Request.Context(), promotion); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "promotion created successfully", promotion)
}

// UpdatePromotionRequest 更新充值活动请求（同CreatePromotionRequest）
type UpdatePromotionRequest = CreatePromotionRequest

// UpdatePromotion 更新充值活动
// @Summary 更新充值活动
// @Tags 管理后台-充值
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "活动ID"
// @Param body body UpdatePromotionRequest true "活动信息"
// @Success 200 {object} response.Response
// @Router /api/admin/recharge/promotions/{id} [put]
func (h *RechargeHandler) UpdatePromotion(c *gin.Context) {
	promotionID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	var req UpdatePromotionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	promotion := &recharge.RechargePromotion{
		ID:             promotionID,
		Title:          req.Title,
		Description:    req.Description,
		RechargeAmount: req.RechargeAmount,
		BonusAmount:    req.BonusAmount,
		BonusPoints:    req.BonusPoints,
		UserLimit:      req.UserLimit,
		TotalLimit:     req.TotalLimit,
		Status:         req.Status,
		SortOrder:      req.SortOrder,
	}

	if err := h.rechargeService.UpdatePromotion(c.Request.Context(), promotion); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "promotion updated successfully", nil)
}

// ListPromotions 充值活动列表
// @Summary 充值活动列表
// @Tags 管理后台-充值
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/recharge/promotions [get]
func (h *RechargeHandler) ListPromotions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	promotions, total, err := h.rechargeService.ListAllPromotions(c.Request.Context(), page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, promotions, total, page, pageSize)
}

// ListOrders 充值订单列表
// @Summary 充值订单列表
// @Tags 管理后台-充值
// @Produce json
// @Security BearerAuth
// @Param user_id query int false "会员ID"
// @Param status query int false "状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/recharge/orders [get]
func (h *RechargeHandler) ListOrders(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var userID *int64
	if userIDStr := c.Query("user_id"); userIDStr != "" {
		id, _ := strconv.ParseInt(userIDStr, 10, 64)
		userID = &id
	}

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = new(int8)
		*status = int8(s)
	}

	orders, total, err := h.rechargeService.ListAllOrders(c.Request.Context(), userID, status, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, orders, total, page, pageSize)
}
