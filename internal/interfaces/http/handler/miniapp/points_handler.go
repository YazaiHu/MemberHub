package miniapp

import (
	"strconv"

	"github.com/gin-gonic/gin"

	pointsApp "github.com/YazaiHu/MemberHub/internal/application/points"
	"github.com/YazaiHu/MemberHub/internal/domain/points"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/middleware"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// PointsHandler 积分处理器
type PointsHandler struct {
	pointsService *pointsApp.Service
}

// NewPointsHandler 创建积分处理器
func NewPointsHandler() *PointsHandler {
	pointsRepo := points.NewRepository()
	return &PointsHandler{
		pointsService: pointsApp.NewService(pointsRepo),
	}
}

// GetBalance 获取积分余额
// @Summary 获取积分余额
// @Tags 小程序-积分
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=points.MemberPoints}
// @Router /api/miniapp/points/balance [get]
func (h *PointsHandler) GetBalance(c *gin.Context) {
	userID := middleware.GetUserID(c)

	memberPoints, err := h.pointsService.GetUserPoints(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, memberPoints)
}

// GetHistory 获取积分历史
// @Summary 获取积分历史
// @Tags 小程序-积分
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/miniapp/points/history [get]
func (h *PointsHandler) GetHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	transactions, total, err := h.pointsService.GetPointsHistory(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, transactions, total, page, pageSize)
}

// ListExchangeRules 兑换规则列表
// @Summary 兑换规则列表
// @Tags 小程序-积分
// @Produce json
// @Security BearerAuth
// @Param store_id query int false "门店ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/miniapp/points/exchange/rules [get]
func (h *PointsHandler) ListExchangeRules(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var storeID *int64
	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		id, _ := strconv.ParseInt(storeIDStr, 10, 64)
		storeID = &id
	}

	rules, total, err := h.pointsService.ListExchangeRules(c.Request.Context(), storeID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, rules, total, page, pageSize)
}

// ExchangeRequest 兑换请求
type ExchangeRequest struct {
	RuleID int64 `json:"rule_id" binding:"required"`
}

// Exchange 积分兑换
// @Summary 积分兑换
// @Tags 小程序-积分
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body ExchangeRequest true "兑换信息"
// @Success 200 {object} response.Response{data=points.ExchangeRecord}
// @Router /api/miniapp/points/exchange [post]
func (h *PointsHandler) Exchange(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req ExchangeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	record, err := h.pointsService.ExchangePoints(c.Request.Context(), userID, req.RuleID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "exchange successful", record)
}

// GetExchangeRecords 兑换记录列表
// @Summary 兑换记录列表
// @Tags 小程序-积分
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/miniapp/points/exchange/records [get]
func (h *PointsHandler) GetExchangeRecords(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	records, total, err := h.pointsService.GetExchangeRecords(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, records, total, page, pageSize)
}
