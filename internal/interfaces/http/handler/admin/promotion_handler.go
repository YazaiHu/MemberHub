package admin

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	promotionApp "github.com/YazaiHu/MemberHub/internal/application/promotion"
	"github.com/YazaiHu/MemberHub/internal/domain/promotion"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// PromotionHandler 特价商品管理处理器
type PromotionHandler struct {
	promotionService *promotionApp.Service
}

// NewPromotionHandler 创建特价商品管理处理器
func NewPromotionHandler() *PromotionHandler {
	promotionRepo := promotion.NewRepository()
	return &PromotionHandler{
		promotionService: promotionApp.NewService(promotionRepo),
	}
}

// ListProducts 特价商品列表
// @Summary 特价商品列表
// @Tags 管理端-特价商品
// @Produce json
// @Security BearerAuth
// @Param store_id query int false "门店ID"
// @Param status query int false "状态" Enums(0, 1)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/promotions [get]
func (h *PromotionHandler) ListProducts(c *gin.Context) {
	var storeID *int64
	if storeIDStr := c.Query("store_id"); storeIDStr != "" {
		id, err := strconv.ParseInt(storeIDStr, 10, 64)
		if err == nil {
			storeID = &id
		}
	}

	var status *int8
	if statusStr := c.Query("status"); statusStr != "" {
		statusVal, err := strconv.Atoi(statusStr)
		if err == nil {
			s := int8(statusVal)
			status = &s
		}
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	products, total, err := h.promotionService.ListProducts(c.Request.Context(), storeID, status, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, products, total, page, pageSize)
}

// GetProduct 特价商品详情
// @Summary 特价商品详情
// @Tags 管理端-特价商品
// @Produce json
// @Security BearerAuth
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response{data=promotion.PromotionProduct}
// @Router /api/admin/promotions/{id} [get]
func (h *PromotionHandler) GetProduct(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	product, err := h.promotionService.GetProduct(c.Request.Context(), productID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, product)
}

// CreateProductRequest 创建特价商品请求
type CreateProductRequest struct {
	StoreID        int64  `json:"store_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	OriginalPrice  int    `json:"original_price" binding:"required,min=1"`
	PromotionPrice int    `json:"promotion_price" binding:"required,min=1"`
	Stock          int    `json:"stock" binding:"required,min=0"`
	ImageURL       string `json:"image_url"`
	StartTime      string `json:"start_time" binding:"required"`
	EndTime        string `json:"end_time" binding:"required"`
	Status         int8   `json:"status" binding:"oneof=0 1"`
	SortOrder      int    `json:"sort_order"`
}

// CreateProduct 创建特价商品
// @Summary 创建特价商品
// @Tags 管理端-特价商品
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateProductRequest true "特价商品信息"
// @Success 200 {object} response.Response{data=promotion.PromotionProduct}
// @Router /api/admin/promotions [post]
func (h *PromotionHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.WithMessage("invalid start_time format"))
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.WithMessage("invalid end_time format"))
		return
	}

	// 验证时间范围
	if endTime.Before(startTime) {
		response.Error(c, errors.ErrInvalidParams.WithMessage("end_time must be after start_time"))
		return
	}

	// 验证价格
	if req.PromotionPrice >= req.OriginalPrice {
		response.Error(c, errors.ErrInvalidParams.WithMessage("promotion_price must be less than original_price"))
		return
	}

	product := &promotion.PromotionProduct{
		StoreID:        req.StoreID,
		Name:           req.Name,
		Description:    req.Description,
		OriginalPrice:  req.OriginalPrice,
		PromotionPrice: req.PromotionPrice,
		Stock:          req.Stock,
		ImageURL:       req.ImageURL,
		StartTime:      startTime,
		EndTime:        endTime,
		Status:         req.Status,
		SortOrder:      req.SortOrder,
	}

	if err := h.promotionService.CreateProduct(c.Request.Context(), product); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "promotion product created successfully", product)
}

// UpdateProductRequest 更新特价商品请求
type UpdateProductRequest struct {
	ID             int64  `json:"id" binding:"required"`
	StoreID        int64  `json:"store_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	OriginalPrice  int    `json:"original_price" binding:"required,min=1"`
	PromotionPrice int    `json:"promotion_price" binding:"required,min=1"`
	Stock          int    `json:"stock" binding:"required,min=0"`
	ImageURL       string `json:"image_url"`
	StartTime      string `json:"start_time" binding:"required"`
	EndTime        string `json:"end_time" binding:"required"`
	Status         int8   `json:"status" binding:"oneof=0 1"`
	SortOrder      int    `json:"sort_order"`
}

// UpdateProduct 更新特价商品
// @Summary 更新特价商品
// @Tags 管理端-特价商品
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body UpdateProductRequest true "特价商品信息"
// @Success 200 {object} response.Response
// @Router /api/admin/promotions [put]
func (h *PromotionHandler) UpdateProduct(c *gin.Context) {
	var req UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	// 解析时间
	startTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.WithMessage("invalid start_time format"))
		return
	}

	endTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.WithMessage("invalid end_time format"))
		return
	}

	// 验证时间范围
	if endTime.Before(startTime) {
		response.Error(c, errors.ErrInvalidParams.WithMessage("end_time must be after start_time"))
		return
	}

	// 验证价格
	if req.PromotionPrice >= req.OriginalPrice {
		response.Error(c, errors.ErrInvalidParams.WithMessage("promotion_price must be less than original_price"))
		return
	}

	product := &promotion.PromotionProduct{
		ID:             req.ID,
		StoreID:        req.StoreID,
		Name:           req.Name,
		Description:    req.Description,
		OriginalPrice:  req.OriginalPrice,
		PromotionPrice: req.PromotionPrice,
		Stock:          req.Stock,
		ImageURL:       req.ImageURL,
		StartTime:      startTime,
		EndTime:        endTime,
		Status:         req.Status,
		SortOrder:      req.SortOrder,
	}

	if err := h.promotionService.UpdateProduct(c.Request.Context(), product); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "promotion product updated successfully", nil)
}

// DeleteProduct 删除特价商品
// @Summary 删除特价商品
// @Tags 管理端-特价商品
// @Produce json
// @Security BearerAuth
// @Param id path int true "商品ID"
// @Success 200 {object} response.Response
// @Router /api/admin/promotions/{id} [delete]
func (h *PromotionHandler) DeleteProduct(c *gin.Context) {
	productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	if err := h.promotionService.DeleteProduct(c.Request.Context(), productID); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "promotion product deleted successfully", nil)
}
