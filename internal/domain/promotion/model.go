package promotion

import (
	"time"
)

// PromotionProduct 特价商品
type PromotionProduct struct {
	ID            int64      `gorm:"column:id;primaryKey" json:"id"`
	StoreID       int64      `gorm:"column:store_id;index;not null" json:"store_id"`
	Name          string     `gorm:"column:name;size:100;not null" json:"name"`
	Description   string     `gorm:"column:description;size:500" json:"description,omitempty"`
	OriginalPrice int        `gorm:"column:original_price;not null" json:"original_price"` // 原价（分）
	PromotionPrice int       `gorm:"column:promotion_price;not null" json:"promotion_price"` // 特价（分）
	Stock         int        `gorm:"column:stock;not null" json:"stock"` // 库存
	ImageURL      string     `gorm:"column:image_url;size:255" json:"image_url,omitempty"`
	StartTime     time.Time  `gorm:"column:start_time;not null;index" json:"start_time"`
	EndTime       time.Time  `gorm:"column:end_time;not null;index" json:"end_time"`
	Status        int8       `gorm:"column:status;default:1;index" json:"status"` // 0-停用, 1-启用
	SortOrder     int        `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (PromotionProduct) TableName() string {
	return "promotion_products"
}

// 特价商品状态常量
const (
	PromotionStatusDisabled = 0 // 停用
	PromotionStatusEnabled  = 1 // 启用
)

// IsActive 判断特价商品是否在有效期内
func (p *PromotionProduct) IsActive() bool {
	now := time.Now()
	return p.Status == PromotionStatusEnabled &&
		!now.Before(p.StartTime) &&
		!now.After(p.EndTime) &&
		p.Stock > 0
}

// DiscountPercent 计算折扣百分比
func (p *PromotionProduct) DiscountPercent() int {
	if p.OriginalPrice == 0 {
		return 0
	}
	return int(float64(p.PromotionPrice) / float64(p.OriginalPrice) * 100)
}

// SaveAmount 计算节省金额（分）
func (p *PromotionProduct) SaveAmount() int {
	return p.OriginalPrice - p.PromotionPrice
}

// PromotionProductWithStore 带门店信息的特价商品
type PromotionProductWithStore struct {
	*PromotionProduct
	StoreName string `json:"store_name,omitempty"`
}
