package miniapp

import (
	"github.com/gin-gonic/gin"

	storeApp "github.com/YazaiHu/MemberHub/internal/application/store"
	"github.com/YazaiHu/MemberHub/internal/domain/store"
	"github.com/YazaiHu/MemberHub/internal/interfaces/http/response"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// StoreHandler 门店处理器
type StoreHandler struct {
	storeService *storeApp.Service
}

// NewStoreHandler 创建门店处理器
func NewStoreHandler() *StoreHandler {
	storeRepo := store.NewRepository()
	return &StoreHandler{
		storeService: storeApp.NewService(storeRepo),
	}
}

// ListStores 门店列表
// @Summary 门店列表
// @Tags 小程序-门店
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]store.Store}
// @Router /api/miniapp/stores [get]
func (h *StoreHandler) ListStores(c *gin.Context) {
	stores, err := h.storeService.GetActiveStores(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, stores)
}

// NearbyStoresRequest 附近门店请求
type NearbyStoresRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Radius    int     `json:"radius"` // 半径（公里），默认10公里
}

// GetNearbyStores 获取附近门店
// @Summary 附近门店
// @Tags 小程序-门店
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body NearbyStoresRequest true "位置信息"
// @Success 200 {object} response.Response{data=[]store.StoreDistance}
// @Router /api/miniapp/stores/nearby [post]
func (h *StoreHandler) GetNearbyStores(c *gin.Context) {
	var req NearbyStoresRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, errors.ErrInvalidParams.Wrap(err))
		return
	}

	stores, err := h.storeService.GetNearbyStores(c.Request.Context(), req.Latitude, req.Longitude, req.Radius)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, stores)
}
