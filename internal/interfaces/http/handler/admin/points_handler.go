package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	pointsApp "github.com/YazaiHu/MemberHub/internal/application/points"
	"github.com/YazaiHu/MemberHub/internal/domain/points"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/middleware"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// PointsHandler 积分管理处理器
type PointsHandler struct {
	pointsService *pointsApp.Service
}

// NewPointsHandler 创建积分管理处理器
func NewPointsHandler() *PointsHandler {
	pointsRepo := points.NewRepository()
	return &PointsHandler{
		pointsService: pointsApp.NewService(pointsRepo),
	}
}

// AdjustPointsRequest 调整积分请求
type AdjustPointsRequest struct {
	Points int64  `json:"points" binding:"required"` // 正数增加，负数扣减
	Remark string `json:"remark" binding:"required"`
}

// AdjustPoints 手动调整积分
// @Summary 手动调整积分
// @Tags 管理后台-积分
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "会员ID"
// @Param body body AdjustPointsRequest true "调整信息"
// @Success 200 {object} response.Response
// @Router /api/admin/members/{id}/adjust-points [post]
func (h *PointsHandler) AdjustPoints(c *gin.Context) {
	adminID := middleware.GetUserID(c)
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	var req AdjustPointsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	if err := h.pointsService.ManualAdjustPoints(c.Request.Context(), userID, req.Points, req.Remark, adminID); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "points adjusted successfully", nil)
}

// ListExchangeRules 兑换规则列表
// @Summary 兑换规则列表
// @Tags 管理后台-积分
// @Produce json
// @Security BearerAuth
// @Param store_id query int false "门店ID"
// @Param status query int false "状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/points/rules [get]
func (h *PointsHandler) ListExchangeRules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var storeID *int64
	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		id, _ := strconv.ParseInt(storeIDStr, 10, 64)
		storeID = &id
	}

	// 管理员可以查看所有状态的规则
	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = new(int8)
		*status = int8(s)
	}

	// TODO: 店铺管理员只能查看自己店铺的规则
	roleCode := middleware.GetRoleCode(c)
	if roleCode == "store_admin" {
		adminStoreID := middleware.GetStoreID(c)
		storeID = &adminStoreID
	}

	rules, total, err := h.pointsService.ListExchangeRules(c.Request.Context(), storeID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, rules, total, page, pageSize)
}

// CreateRuleRequest 创建兑换规则请求
type CreateRuleRequest struct {
	StoreID        int64   `json:"store_id" binding:"required"`
	Title          string  `json:"title" binding:"required"`
	Description    string  `json:"description"`
	PointsCost     int     `json:"points_cost" binding:"required,gt=0"`
	ProductType    int8    `json:"product_type" binding:"required,oneof=1 2 3"`
	ProductName    string  `json:"product_name" binding:"required"`
	ProductImage   string  `json:"product_image"`
	TotalStock     *int    `json:"total_stock"`
	DailyLimit     *int    `json:"daily_limit"`
	UserDailyLimit int     `json:"user_daily_limit" binding:"required,gt=0"`
	UserTotalLimit *int    `json:"user_total_limit"`
	StartTime      *string `json:"start_time"` // YYYY-MM-DD HH:MM:SS
	EndTime        *string `json:"end_time"`
	Status         int8    `json:"status" binding:"oneof=0 1"`
	SortOrder      int     `json:"sort_order"`
}

// CreateRule 创建兑换规则
// @Summary 创建兑换规则
// @Tags 管理后台-积分
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateRuleRequest true "规则信息"
// @Success 200 {object} response.Response{data=points.ExchangeRule}
// @Router /api/admin/points/rules [post]
func (h *PointsHandler) CreateRule(c *gin.Context) {
	var req CreateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	rule := &points.ExchangeRule{
		StoreID:        req.StoreID,
		Title:          req.Title,
		Description:    req.Description,
		PointsCost:     req.PointsCost,
		ProductType:    req.ProductType,
		ProductName:    req.ProductName,
		ProductImage:   req.ProductImage,
		TotalStock:     req.TotalStock,
		RemainingStock: req.TotalStock,
		DailyLimit:     req.DailyLimit,
		UserDailyLimit: req.UserDailyLimit,
		UserTotalLimit: req.UserTotalLimit,
		Status:         req.Status,
		SortOrder:      req.SortOrder,
	}

	// TODO: 解析时间字段

	if err := h.pointsService.CreateExchangeRule(c.Request.Context(), rule); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "rule created successfully", rule)
}

// UpdateRuleRequest 更新兑换规则请求（同CreateRuleRequest）
type UpdateRuleRequest = CreateRuleRequest

// UpdateRule 更新兑换规则
// @Summary 更新兑换规则
// @Tags 管理后台-积分
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "规则ID"
// @Param body body UpdateRuleRequest true "规则信息"
// @Success 200 {object} response.Response
// @Router /api/admin/points/rules/{id} [put]
func (h *PointsHandler) UpdateRule(c *gin.Context) {
	ruleID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	var req UpdateRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	rule := &points.ExchangeRule{
		ID:             ruleID,
		StoreID:        req.StoreID,
		Title:          req.Title,
		Description:    req.Description,
		PointsCost:     req.PointsCost,
		ProductType:    req.ProductType,
		ProductName:    req.ProductName,
		ProductImage:   req.ProductImage,
		TotalStock:     req.TotalStock,
		RemainingStock: req.TotalStock,
		DailyLimit:     req.DailyLimit,
		UserDailyLimit: req.UserDailyLimit,
		UserTotalLimit: req.UserTotalLimit,
		Status:         req.Status,
		SortOrder:      req.SortOrder,
	}

	if err := h.pointsService.UpdateExchangeRule(c.Request.Context(), rule); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "rule updated successfully", nil)
}

// VerifyExchangeRequest 核销请求
type VerifyExchangeRequest struct {
	VerifyCode string `json:"verify_code" binding:"required"`
}

// VerifyExchange 核销兑换
// @Summary 核销兑换
// @Tags 管理后台-积分
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body VerifyExchangeRequest true "核销信息"
// @Success 200 {object} response.Response
// @Router /api/admin/points/exchange/verify [post]
func (h *PointsHandler) VerifyExchange(c *gin.Context) {
	adminID := middleware.GetUserID(c)

	var req VerifyExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	if err := h.pointsService.VerifyExchange(c.Request.Context(), req.VerifyCode, adminID); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "verified successfully", nil)
}

// ListExchangeRecords 兑换记录列表
// @Summary 兑换记录列表
// @Tags 管理后台-积分
// @Produce json
// @Security BearerAuth
// @Param user_id query int false "会员ID"
// @Param store_id query int false "门店ID"
// @Param status query int false "状态"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/points/exchange/records [get]
func (h *PointsHandler) ListExchangeRecords(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// TODO: 实现完整的查询逻辑
	// var userID *int64
	// if userIDStr := c.Query("user_id"); userIDStr != "" {
	// 	id, _ := strconv.ParseInt(userIDStr, 10, 64)
	// 	userID = &id
	// }

	// var storeID *int64
	// if storeIDStr := c.Query("store_id"); storeIDStr != "" {
	// 	id, _ := strconv.ParseInt(storeIDStr, 10, 64)
	// 	storeID = &id
	// }

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		s, _ := strconv.Atoi(statusStr)
		status = new(int8)
		*status = int8(s)
	}

	// 店铺管理员只能查看自己店铺的记录
	// roleCode := middleware.GetRoleCode(c)
	// if roleCode == "store_admin" {
	// 	adminStoreID := middleware.GetStoreID(c)
	// 	storeID = &adminStoreID
	// }

	// TODO: 调用仓储方法获取记录列表
	// records, total, err := h.pointsService.ListExchangeRecords(ctx, userID, storeID, status, offset, limit)
	_ = status // Suppress unused variable warning for now

	response.SuccessWithPage(c, []interface{}{}, 0, page, pageSize)
}
