package miniapp

import (
	"strconv"

	"github.com/gin-gonic/gin"

	couponApp "github.com/YazaiHu/MemberHub/internal/application/coupon"
	"github.com/YazaiHu/MemberHub/internal/domain/coupon"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/middleware"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// CouponHandler 优惠券处理器
type CouponHandler struct {
	couponService *couponApp.Service
}

// NewCouponHandler 创建优惠券处理器
func NewCouponHandler() *CouponHandler {
	couponRepo := coupon.NewRepository()
	return &CouponHandler{
		couponService: couponApp.NewService(couponRepo),
	}
}

// ListTemplates 优惠券模板列表
// @Summary 优惠券模板列表
// @Tags 小程序-优惠券
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]coupon.CouponTemplate}
// @Router /api/miniapp/coupons/templates [get]
func (h *CouponHandler) ListTemplates(c *gin.Context) {
	templates, err := h.couponService.GetActiveTemplates(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, templates)
}

// ReceiveCouponRequest 领取优惠券请求
type ReceiveCouponRequest struct {
	TemplateID int64 `json:"template_id" binding:"required"`
}

// ReceiveCoupon 领取优惠券
// @Summary 领取优惠券
// @Tags 小程序-优惠券
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body ReceiveCouponRequest true "优惠券信息"
// @Success 200 {object} response.Response{data=coupon.UserCoupon}
// @Router /api/miniapp/coupons/receive [post]
func (h *CouponHandler) ReceiveCoupon(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req ReceiveCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	userCoupon, err := h.couponService.ReceiveCoupon(c.Request.Context(), userID, req.TemplateID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "coupon received successfully", userCoupon)
}

// GetAvailableCoupons 获取可用优惠券
// @Summary 获取可用优惠券
// @Tags 小程序-优惠券
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/miniapp/coupons/available [get]
func (h *CouponHandler) GetAvailableCoupons(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	coupons, total, err := h.couponService.GetAvailableCoupons(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, coupons, total, page, pageSize)
}

// GetAllCoupons 获取所有优惠券
// @Summary 获取所有优惠券
// @Tags 小程序-优惠券
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/miniapp/coupons/all [get]
func (h *CouponHandler) GetAllCoupons(c *gin.Context) {
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	coupons, total, err := h.couponService.GetAllCoupons(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, coupons, total, page, pageSize)
}
