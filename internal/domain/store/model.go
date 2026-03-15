package store

import (
	"time"
)

// Store 门店
type Store struct {
	ID           int64     `gorm:"column:id;primaryKey" json:"id"`
	Name         string    `gorm:"column:name;size:100;not null" json:"name"`
	Address      string    `gorm:"column:address;size:255;not null" json:"address"`
	Phone        string    `gorm:"column:phone;size:20" json:"phone,omitempty"`
	Latitude     float64   `gorm:"column:latitude;type:decimal(10,7)" json:"latitude"`
	Longitude    float64   `gorm:"column:longitude;type:decimal(10,7)" json:"longitude"`
	OpeningHours string    `gorm:"column:opening_hours;size:100" json:"opening_hours,omitempty"` // 营业时间
	Description  string    `gorm:"column:description;size:500" json:"description,omitempty"`
	Status       int8      `gorm:"column:status;default:1;index" json:"status"` // 0-停用, 1-启用
	SortOrder    int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (Store) TableName() string {
	return "stores"
}

// 门店状态常量
const (
	StoreStatusDisabled = 0 // 停用
	StoreStatusEnabled  = 1 // 启用
)

// NearbyStoreRequest 附近门店查询请求
type NearbyStoreRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Radius    int     `json:"radius"` // 半径（公里），默认10公里
}

// StoreDistance 门店距离信息
type StoreDistance struct {
	Store    *Store  `json:"store"`
	Distance float64 `json:"distance"` // 距离（公里）
}
