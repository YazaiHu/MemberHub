package promotion

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// Repository 特价商品仓储接口
type Repository interface {
	// PromotionProduct CRUD
	CreateProduct(ctx context.Context, product *PromotionProduct) error
	GetProductByID(ctx context.Context, productID int64) (*PromotionProduct, error)
	UpdateProduct(ctx context.Context, product *PromotionProduct) error
	DeleteProduct(ctx context.Context, productID int64) error
	ListProducts(ctx context.Context, storeID *int64, status *int8, offset, limit int) ([]*PromotionProduct, int64, error)
	GetActiveProducts(ctx context.Context, storeID *int64) ([]*PromotionProduct, error)
	GetWeeklyProducts(ctx context.Context, storeID *int64) ([]*PromotionProduct, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository 创建特价商品仓储
func NewRepository() Repository {
	return &repository{
		db: mysql.GetDB(),
	}
}

// CreateProduct 创建特价商品
func (r *repository) CreateProduct(ctx context.Context, product *PromotionProduct) error {
	if err := r.db.WithContext(ctx).Create(product).Error; err != nil {
		return fmt.Errorf("create promotion product failed: %w", err)
	}
	return nil
}

// GetProductByID 获取特价商品
func (r *repository) GetProductByID(ctx context.Context, productID int64) (*PromotionProduct, error) {
	var product PromotionProduct
	if err := r.db.WithContext(ctx).Where("id = ?", productID).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get promotion product failed: %w", err)
	}
	return &product, nil
}

// UpdateProduct 更新特价商品
func (r *repository) UpdateProduct(ctx context.Context, product *PromotionProduct) error {
	if err := r.db.WithContext(ctx).Save(product).Error; err != nil {
		return fmt.Errorf("update promotion product failed: %w", err)
	}
	return nil
}

// DeleteProduct 删除特价商品（软删除）
func (r *repository) DeleteProduct(ctx context.Context, productID int64) error {
	result := r.db.WithContext(ctx).Model(&PromotionProduct{}).
		Where("id = ?", productID).
		Update("status", PromotionStatusDisabled)

	if result.Error != nil {
		return fmt.Errorf("delete promotion product failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// ListProducts 特价商品列表
func (r *repository) ListProducts(ctx context.Context, storeID *int64, status *int8, offset, limit int) ([]*PromotionProduct, int64, error) {
	var products []*PromotionProduct
	var total int64

	query := r.db.WithContext(ctx).Model(&PromotionProduct{})

	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}

	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count promotion products failed: %w", err)
	}

	if err := query.Order("sort_order DESC, start_time DESC, id DESC").
		Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("list promotion products failed: %w", err)
	}

	return products, total, nil
}

// GetActiveProducts 获取有效的特价商品
func (r *repository) GetActiveProducts(ctx context.Context, storeID *int64) ([]*PromotionProduct, error) {
	var products []*PromotionProduct
	query := r.db.WithContext(ctx).Where("status = ? AND stock > 0", PromotionStatusEnabled)

	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}

	// 时间范围检查
	now := time.Now()
	query = query.Where("start_time <= ? AND end_time >= ?", now, now)

	if err := query.Order("sort_order DESC, start_time DESC, id DESC").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("get active promotion products failed: %w", err)
	}

	return products, nil
}

// GetWeeklyProducts 获取本周特价商品（7天内的特价）
func (r *repository) GetWeeklyProducts(ctx context.Context, storeID *int64) ([]*PromotionProduct, error) {
	var products []*PromotionProduct
	query := r.db.WithContext(ctx).Where("status = ? AND stock > 0", PromotionStatusEnabled)

	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}

	// 时间范围：从现在开始到7天后
	now := time.Now()
	weekLater := now.AddDate(0, 0, 7)
	query = query.Where("start_time <= ? AND end_time >= ?", weekLater, now)

	if err := query.Order("sort_order DESC, start_time ASC, id DESC").Find(&products).Error; err != nil {
		return nil, fmt.Errorf("get weekly promotion products failed: %w", err)
	}

	return products, nil
}
