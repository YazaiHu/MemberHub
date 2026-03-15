package miniapp

import (
	"strconv"

	"github.com/gin-gonic/gin"

	promotionApp "github.com/YazaiHu/MemberHub/internal/application/promotion"
	"github.com/YazaiHu/MemberHub/internal/domain/promotion"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
)

// PromotionHandler 特价商品处理器
type PromotionHandler struct {
	promotionService *promotionApp.Service
}

// NewPromotionHandler 创建特价商品处理器
func NewPromotionHandler() *PromotionHandler {
	promotionRepo := promotion.NewRepository()
	return &PromotionHandler{
		promotionService: promotionApp.NewService(promotionRepo),
	}
}

// GetWeeklyPromotions 本周特价商品
// @Summary 本周特价商品
// @Tags 小程序-特价商品
// @Produce json
// @Security BearerAuth
// @Param store_id query int false "门店ID（可选，不传则返回所有门店）"
// @Success 200 {object} response.Response{data=[]promotion.PromotionProduct}
// @Router /api/miniapp/promotions/weekly [get]
func (h *PromotionHandler) GetWeeklyPromotions(c *gin.Context) {
	var storeID *int64
	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		id, err := strconv.ParseInt(storeIDStr, 10, 64)
		if err == nil {
			storeID = &id
		}
	}

	products, err := h.promotionService.GetWeeklyProducts(c.Request.Context(), storeID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, products)
}

// GetActivePromotions 当前有效特价商品
// @Summary 当前有效特价商品
// @Tags 小程序-特价商品
// @Produce json
// @Security BearerAuth
// @Param store_id query int false "门店ID（可选，不传则返回所有门店）"
// @Success 200 {object} response.Response{data=[]promotion.PromotionProduct}
// @Router /api/miniapp/promotions [get]
func (h *PromotionHandler) GetActivePromotions(c *gin.Context) {
	var storeID *int64
	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		id, err := strconv.ParseInt(storeIDStr, 10, 64)
		if err == nil {
			storeID = &id
		}
	}

	products, err := h.promotionService.GetActiveProducts(c.Request.Context(), storeID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, products)
}
