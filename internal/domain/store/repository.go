package store

import (
	"context"
	"fmt"
	"math"

	"gorm.io/gorm"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// Repository 门店仓储接口
type Repository interface {
	// Store CRUD
	CreateStore(ctx context.Context, store *Store) error
	GetStoreByID(ctx context.Context, storeID int64) (*Store, error)
	UpdateStore(ctx context.Context, store *Store) error
	DeleteStore(ctx context.Context, storeID int64) error
	ListStores(ctx context.Context, status *int8, offset, limit int) ([]*Store, int64, error)
	GetActiveStores(ctx context.Context) ([]*Store, error)
	GetNearbyStores(ctx context.Context, latitude, longitude float64, radius int) ([]*StoreDistance, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository 创建门店仓储
func NewRepository() Repository {
	return &repository{
		db: mysql.GetDB(),
	}
}

// CreateStore 创建门店
func (r *repository) CreateStore(ctx context.Context, store *Store) error {
	if err := r.db.WithContext(ctx).Create(store).Error; err != nil {
		return fmt.Errorf("create store failed: %w", err)
	}
	return nil
}

// GetStoreByID 获取门店
func (r *repository) GetStoreByID(ctx context.Context, storeID int64) (*Store, error) {
	var store Store
	if err := r.db.WithContext(ctx).Where("id = ?", storeID).First(&store).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get store failed: %w", err)
	}
	return &store, nil
}

// UpdateStore 更新门店
func (r *repository) UpdateStore(ctx context.Context, store *Store) error {
	if err := r.db.WithContext(ctx).Save(store).Error; err != nil {
		return fmt.Errorf("update store failed: %w", err)
	}
	return nil
}

// DeleteStore 删除门店（软删除，实际是更新状态）
func (r *repository) DeleteStore(ctx context.Context, storeID int64) error {
	result := r.db.WithContext(ctx).Model(&Store{}).
		Where("id = ?", storeID).
		Update("status", StoreStatusDisabled)

	if result.Error != nil {
		return fmt.Errorf("delete store failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrNotFound
	}

	return nil
}

// ListStores 门店列表
func (r *repository) ListStores(ctx context.Context, status *int8, offset, limit int) ([]*Store, int64, error) {
	var stores []*Store
	var total int64

	query := r.db.WithContext(ctx).Model(&Store{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count stores failed: %w", err)
	}

	if err := query.Order("sort_order DESC, id DESC").Offset(offset).Limit(limit).Find(&stores).Error; err != nil {
		return nil, 0, fmt.Errorf("list stores failed: %w", err)
	}

	return stores, total, nil
}

// GetActiveStores 获取启用的门店
func (r *repository) GetActiveStores(ctx context.Context) ([]*Store, error) {
	var stores []*Store
	if err := r.db.WithContext(ctx).
		Where("status = ?", StoreStatusEnabled).
		Order("sort_order DESC, id DESC").
		Find(&stores).Error; err != nil {
		return nil, fmt.Errorf("get active stores failed: %w", err)
	}
	return stores, nil
}

// GetNearbyStores 获取附近门店
// 使用球面距离公式（Haversine formula）计算距离
func (r *repository) GetNearbyStores(ctx context.Context, latitude, longitude float64, radius int) ([]*StoreDistance, error) {
	if radius <= 0 {
		radius = 10 // 默认10公里
	}

	// 查询所有启用的门店
	stores, err := r.GetActiveStores(ctx)
	if err != nil {
		return nil, err
	}

	// 计算距离并过滤
	var nearbyStores []*StoreDistance
	for _, store := range stores {
		distance := calculateDistance(latitude, longitude, store.Latitude, store.Longitude)
		if distance <= float64(radius) {
			nearbyStores = append(nearbyStores, &StoreDistance{
				Store:    store,
				Distance: distance,
			})
		}
	}

	// 按距离排序（冒泡排序，因为数据量不大）
	for i := 0; i < len(nearbyStores); i++ {
		for j := i + 1; j < len(nearbyStores); j++ {
			if nearbyStores[i].Distance > nearbyStores[j].Distance {
				nearbyStores[i], nearbyStores[j] = nearbyStores[j], nearbyStores[i]
			}
		}
	}

	return nearbyStores, nil
}

// calculateDistance 计算两个经纬度之间的距离（公里）
// 使用 Haversine 公式
func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371.0 // 地球半径（公里）

	// 转换为弧度
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	// Haversine 公式
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadius * c
	return math.Round(distance*100) / 100 // 保留两位小数
}
