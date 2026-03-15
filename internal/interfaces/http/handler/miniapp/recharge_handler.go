package miniapp

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

// RechargeHandler 充值处理器
type RechargeHandler struct {
	rechargeService *rechargeApp.Service
}

// NewRechargeHandler 创建充值处理器
func NewRechargeHandler() *RechargeHandler {
	rechargeRepo := recharge.NewRepository()
	pointsRepo := points.NewRepository()
	return &RechargeHandler{
		rechargeService: rechargeApp.NewService(rechargeRepo, pointsRepo),
	}
}

// GetPromotions 获取充值活动列表
// @Summary 获取充值活动列表
// @Tags 小程序-充值
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]recharge.RechargePromotion}
// @Router /api/miniapp/recharge/promotions [get]
func (h *RechargeHandler) GetPromotions(c *gin.Context) {
	promotions, err := h.rechargeService.GetActivePromotions(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, promotions)
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	PromotionID *int64 `json:"promotion_id"` // 活动ID，可选
	Amount      int    `json:"amount"`       // 充值金额（元），不参加活动时必填
}

// CreateOrder 创建充值订单
// @Summary 创建充值订单
// @Tags 小程序-充值
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateOrderRequest true "订单信息"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Router /api/miniapp/recharge/create-order [post]
func (h *RechargeHandler) CreateOrder(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	// 创建订单
	domainReq := &recharge.CreateOrderRequest{
		PromotionID: req.PromotionID,
		Amount:      req.Amount,
	}

	order, err := h.rechargeService.CreateOrder(c.Request.Context(), userID, domainReq)
	if err != nil {
		response.Error(c, err)
		return
	}

	// TODO: 调用微信支付统一下单API，获取支付参数
	// 这里暂时返回订单信息，实际应该返回微信支付参数
	payParams := map[string]interface{}{
		"order_no":   order.OrderNo,
		"pay_amount": order.PayAmount,
		// 实际应该包含微信支付参数：
		// "appId": ...,
		// "timeStamp": ...,
		// "nonceStr": ...,
		// "package": ...,
		// "signType": ...,
		// "paySign": ...,
	}

	response.Success(c, gin.H{
		"order":       order,
		"pay_params":  payParams,
	})
}

// GetBalance 获取余额
// @Summary 获取余额
// @Tags 小程序-充值
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=recharge.MemberBalance}
// @Router /api/miniapp/recharge/balance [get]
func (h *RechargeHandler) GetBalance(c *gin.Context) {
	userID := middleware.GetUserID(c)

	balance, err := h.rechargeService.GetUserBalance(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, balance)
}

// GetBalanceHistory 获取余额历史
// @Summary 获取余额历史
// @Tags 小程序-充值
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/miniapp/recharge/balance/history [get]
func (h *RechargeHandler) GetBalanceHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	transactions, total, err := h.rechargeService.GetBalanceHistory(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, transactions, total, page, pageSize)
}

// GetRechargeRecords 获取充值记录
// @Summary 获取充值记录
// @Tags 小程序-充值
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/miniapp/recharge/records [get]
func (h *RechargeHandler) GetRechargeRecords(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	records, total, err := h.rechargeService.GetRechargeRecords(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, records, total, page, pageSize)
}
