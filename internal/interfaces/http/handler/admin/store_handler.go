package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"

	storeApp "github.com/YazaiHu/MemberHub/internal/application/store"
	"github.com/YazaiHu/MemberHub/internal/domain/store"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// StoreHandler 门店管理处理器
type StoreHandler struct {
	storeService *storeApp.Service
}

// NewStoreHandler 创建门店管理处理器
func NewStoreHandler() *StoreHandler {
	storeRepo := store.NewRepository()
	return &StoreHandler{
		storeService: storeApp.NewService(storeRepo),
	}
}

// ListStores 门店列表
// @Summary 门店列表
// @Tags 管理端-门店
// @Produce json
// @Security BearerAuth
// @Param status query int false "状态" Enums(0, 1)
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} response.Response{data=response.PageData}
// @Router /api/admin/stores [get]
func (h *StoreHandler) ListStores(c *gin.Context) {
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

	stores, total, err := h.storeService.ListStores(c.Request.Context(), status, page, pageSize)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithPage(c, stores, total, page, pageSize)
}

// GetStore 门店详情
// @Summary 门店详情
// @Tags 管理端-门店
// @Produce json
// @Security BearerAuth
// @Param id path int true "门店ID"
// @Success 200 {object} response.Response{data=store.Store}
// @Router /api/admin/stores/{id} [get]
func (h *StoreHandler) GetStore(c *gin.Context) {
	storeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	st, err := h.storeService.GetStore(c.Request.Context(), storeID)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, st)
}

// CreateStoreRequest 创建门店请求
type CreateStoreRequest struct {
	Name         string  `json:"name" binding:"required"`
	Address      string  `json:"address" binding:"required"`
	Phone        string  `json:"phone"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
	OpeningHours string  `json:"opening_hours"`
	Description  string  `json:"description"`
	Status       int8    `json:"status" binding:"oneof=0 1"`
	SortOrder    int     `json:"sort_order"`
}

// CreateStore 创建门店
// @Summary 创建门店
// @Tags 管理端-门店
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body CreateStoreRequest true "门店信息"
// @Success 200 {object} response.Response{data=store.Store}
// @Router /api/admin/stores [post]
func (h *StoreHandler) CreateStore(c *gin.Context) {
	var req CreateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	st := &store.Store{
		Name:         req.Name,
		Address:      req.Address,
		Phone:        req.Phone,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		OpeningHours: req.OpeningHours,
		Description:  req.Description,
		Status:       req.Status,
		SortOrder:    req.SortOrder,
	}

	if err := h.storeService.CreateStore(c.Request.Context(), st); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "store created successfully", st)
}

// UpdateStoreRequest 更新门店请求
type UpdateStoreRequest struct {
	ID           int64   `json:"id" binding:"required"`
	Name         string  `json:"name" binding:"required"`
	Address      string  `json:"address" binding:"required"`
	Phone        string  `json:"phone"`
	Latitude     float64 `json:"latitude" binding:"required"`
	Longitude    float64 `json:"longitude" binding:"required"`
	OpeningHours string  `json:"opening_hours"`
	Description  string  `json:"description"`
	Status       int8    `json:"status" binding:"oneof=0 1"`
	SortOrder    int     `json:"sort_order"`
}

// UpdateStore 更新门店
// @Summary 更新门店
// @Tags 管理端-门店
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body UpdateStoreRequest true "门店信息"
// @Success 200 {object} response.Response
// @Router /api/admin/stores [put]
func (h *StoreHandler) UpdateStore(c *gin.Context) {
	var req UpdateStoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	st := &store.Store{
		ID:           req.ID,
		Name:         req.Name,
		Address:      req.Address,
		Phone:        req.Phone,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		OpeningHours: req.OpeningHours,
		Description:  req.Description,
		Status:       req.Status,
		SortOrder:    req.SortOrder,
	}

	if err := h.storeService.UpdateStore(c.Request.Context(), st); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "store updated successfully", nil)
}

// DeleteStore 删除门店
// @Summary 删除门店
// @Tags 管理端-门店
// @Produce json
// @Security BearerAuth
// @Param id path int true "门店ID"
// @Success 200 {object} response.Response
// @Router /api/admin/stores/{id} [delete]
func (h *StoreHandler) DeleteStore(c *gin.Context) {
	storeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	if err := h.storeService.DeleteStore(c.Request.Context(), storeID); err != nil {
		response.Error(c, err)
		return
	}

	response.SuccessWithMessage(c, "store deleted successfully", nil)
}
