package store

import (
	"context"

	"github.com/YazaiHu/MemberHub/internal/domain/store"
)

// Service 门店应用服务
type Service struct {
	storeRepo store.Repository
}

// NewService 创建门店应用服务
func NewService(storeRepo store.Repository) *Service {
	return &Service{
		storeRepo: storeRepo,
	}
}

// ListStores 门店列表
func (s *Service) ListStores(ctx context.Context, status *int8, page, pageSize int) ([]*store.Store, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.storeRepo.ListStores(ctx, status, offset, pageSize)
}

// GetStore 获取门店详情
func (s *Service) GetStore(ctx context.Context, storeID int64) (*store.Store, error) {
	return s.storeRepo.GetStoreByID(ctx, storeID)
}

// CreateStore 创建门店
func (s *Service) CreateStore(ctx context.Context, st *store.Store) error {
	return s.storeRepo.CreateStore(ctx, st)
}

// UpdateStore 更新门店
func (s *Service) UpdateStore(ctx context.Context, st *store.Store) error {
	return s.storeRepo.UpdateStore(ctx, st)
}

// DeleteStore 删除门店
func (s *Service) DeleteStore(ctx context.Context, storeID int64) error {
	return s.storeRepo.DeleteStore(ctx, storeID)
}

// GetActiveStores 获取启用的门店
func (s *Service) GetActiveStores(ctx context.Context) ([]*store.Store, error) {
	return s.storeRepo.GetActiveStores(ctx)
}

// GetNearbyStores 获取附近门店
func (s *Service) GetNearbyStores(ctx context.Context, latitude, longitude float64, radius int) ([]*store.StoreDistance, error) {
	return s.storeRepo.GetNearbyStores(ctx, latitude, longitude, radius)
}
