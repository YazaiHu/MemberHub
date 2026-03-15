package promotion

import (
	"context"

	"github.com/YazaiHu/MemberHub/internal/domain/promotion"
)

// Service 特价商品应用服务
type Service struct {
	promotionRepo promotion.Repository
}

// NewService 创建特价商品应用服务
func NewService(promotionRepo promotion.Repository) *Service {
	return &Service{
		promotionRepo: promotionRepo,
	}
}

// ListProducts 特价商品列表
func (s *Service) ListProducts(ctx context.Context, storeID *int64, status *int8, page, pageSize int) ([]*promotion.PromotionProduct, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.promotionRepo.ListProducts(ctx, storeID, status, offset, pageSize)
}

// GetProduct 获取特价商品详情
func (s *Service) GetProduct(ctx context.Context, productID int64) (*promotion.PromotionProduct, error) {
	return s.promotionRepo.GetProductByID(ctx, productID)
}

// CreateProduct 创建特价商品
func (s *Service) CreateProduct(ctx context.Context, product *promotion.PromotionProduct) error {
	return s.promotionRepo.CreateProduct(ctx, product)
}

// UpdateProduct 更新特价商品
func (s *Service) UpdateProduct(ctx context.Context, product *promotion.PromotionProduct) error {
	return s.promotionRepo.UpdateProduct(ctx, product)
}

// DeleteProduct 删除特价商品
func (s *Service) DeleteProduct(ctx context.Context, productID int64) error {
	return s.promotionRepo.DeleteProduct(ctx, productID)
}

// GetActiveProducts 获取有效的特价商品
func (s *Service) GetActiveProducts(ctx context.Context, storeID *int64) ([]*promotion.PromotionProduct, error) {
	return s.promotionRepo.GetActiveProducts(ctx, storeID)
}

// GetWeeklyProducts 获取本周特价商品
func (s *Service) GetWeeklyProducts(ctx context.Context, storeID *int64) ([]*promotion.PromotionProduct, error) {
	return s.promotionRepo.GetWeeklyProducts(ctx, storeID)
}
