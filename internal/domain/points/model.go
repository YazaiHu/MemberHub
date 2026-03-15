package points

import (
	"time"
)

// MemberPoints 会员积分汇总表
type MemberPoints struct {
	ID              int64     `gorm:"column:id;primaryKey" json:"id"`
	UserID          int64     `gorm:"column:user_id;uniqueIndex;not null" json:"user_id"`
	TotalPoints     int64     `gorm:"column:total_points;default:0" json:"total_points"`       // 累计获得积分
	AvailablePoints int64     `gorm:"column:available_points;default:0" json:"available_points"` // 可用积分
	FrozenPoints    int64     `gorm:"column:frozen_points;default:0" json:"frozen_points"`     // 冻结积分
	UsedPoints      int64     `gorm:"column:used_points;default:0" json:"used_points"`         // 已使用积分
	ExpiredPoints   int64     `gorm:"column:expired_points;default:0" json:"expired_points"`   // 已过期积分
	Version         int       `gorm:"column:version;default:0" json:"version"`                // 乐观锁版本号
	CreatedAt       time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt       time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (MemberPoints) TableName() string {
	return "member_points"
}

// PointsTransaction 积分流水表
type PointsTransaction struct {
	ID            int64      `gorm:"column:id;primaryKey" json:"id"`
	TransactionNo string     `gorm:"column:transaction_no;uniqueIndex;size:64;not null" json:"transaction_no"`
	UserID        int64      `gorm:"column:user_id;index;not null" json:"user_id"`
	Points        int64      `gorm:"column:points;not null" json:"points"`       // 正数增加，负数减少
	Type          int8       `gorm:"column:type;not null;index" json:"type"`     // 1-消费赠送, 2-充值赠送, 3-手动调整, 4-兑换扣减, 5-过期扣减
	Source        string     `gorm:"column:source;size:50;not null" json:"source"` // 来源
	RefID         *int64     `gorm:"column:ref_id;index" json:"ref_id,omitempty"`        // 关联ID
	Remark        string     `gorm:"column:remark;size:255" json:"remark,omitempty"`
	ExpireAt      *time.Time `gorm:"column:expire_at;index" json:"expire_at,omitempty"` // 过期时间
	Status        int8       `gorm:"column:status;default:1" json:"status"`      // 0-已撤销, 1-正常
	OperatorID    *int64     `gorm:"column:operator_id" json:"operator_id,omitempty"`
	CreatedAt     time.Time  `gorm:"column:created_at;index" json:"created_at"`
}

// TableName 指定表名
func (PointsTransaction) TableName() string {
	return "points_transactions"
}

// 积分变动类型常量
const (
	PointsTypeConsume  = 1 // 消费赠送
	PointsTypeRecharge = 2 // 充值赠送
	PointsTypeAdjust   = 3 // 手动调整
	PointsTypeExchange = 4 // 兑换扣减
	PointsTypeExpire   = 5 // 过期扣减
)

// ExchangeRule 积分兑换规则
type ExchangeRule struct {
	ID             int64      `gorm:"column:id;primaryKey" json:"id"`
	StoreID        int64      `gorm:"column:store_id;index;not null" json:"store_id"`
	Title          string     `gorm:"column:title;size:100;not null" json:"title"`
	Description    string     `gorm:"column:description;size:255" json:"description,omitempty"`
	PointsCost     int        `gorm:"column:points_cost;not null" json:"points_cost"`
	ProductType    int8       `gorm:"column:product_type;not null" json:"product_type"` // 1-实物, 2-优惠券, 3-服务
	ProductName    string     `gorm:"column:product_name;size:100;not null" json:"product_name"`
	ProductImage   string     `gorm:"column:product_image;size:255" json:"product_image,omitempty"`
	TotalStock     *int       `gorm:"column:total_stock" json:"total_stock,omitempty"`       // NULL表示不限
	RemainingStock *int       `gorm:"column:remaining_stock" json:"remaining_stock,omitempty"`
	DailyLimit     *int       `gorm:"column:daily_limit" json:"daily_limit,omitempty"`       // 每日限兑总数
	UserDailyLimit int        `gorm:"column:user_daily_limit;default:1" json:"user_daily_limit"` // 每人每日限兑
	UserTotalLimit *int       `gorm:"column:user_total_limit" json:"user_total_limit,omitempty"` // 每人总限兑
	StartTime      *time.Time `gorm:"column:start_time" json:"start_time,omitempty"`
	EndTime        *time.Time `gorm:"column:end_time" json:"end_time,omitempty"`
	Status         int8       `gorm:"column:status;default:1;index" json:"status"` // 0-下架, 1-上架
	SortOrder      int        `gorm:"column:sort_order;default:0" json:"sort_order"`
	Version        int        `gorm:"column:version;default:0" json:"version"` // 乐观锁
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (ExchangeRule) TableName() string {
	return "points_exchange_rules"
}

// ExchangeRecord 积分兑换记录
type ExchangeRecord struct {
	ID          int64      `gorm:"column:id;primaryKey" json:"id"`
	RecordNo    string     `gorm:"column:record_no;uniqueIndex;size:64;not null" json:"record_no"`
	UserID      int64      `gorm:"column:user_id;index;not null" json:"user_id"`
	RuleID      int64      `gorm:"column:rule_id;index;not null" json:"rule_id"`
	StoreID     int64      `gorm:"column:store_id;index;not null" json:"store_id"`
	PointsCost  int        `gorm:"column:points_cost;not null" json:"points_cost"`
	ProductName string     `gorm:"column:product_name;size:100;not null" json:"product_name"`
	Status      int8       `gorm:"column:status;default:1;index" json:"status"` // 1-待核销, 2-已核销, 3-已取消
	VerifyCode  string     `gorm:"column:verify_code;uniqueIndex;size:20" json:"verify_code,omitempty"`
	VerifiedAt  *time.Time `gorm:"column:verified_at" json:"verified_at,omitempty"`
	VerifiedBy  *int64     `gorm:"column:verified_by" json:"verified_by,omitempty"`
	CancelledAt *time.Time `gorm:"column:cancelled_at" json:"cancelled_at,omitempty"`
	CancelReason string    `gorm:"column:cancel_reason;size:255" json:"cancel_reason,omitempty"`
	CreatedAt   time.Time  `gorm:"column:created_at;index" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (ExchangeRecord) TableName() string {
	return "points_exchange_records"
}

// 兑换记录状态常量
const (
	ExchangeStatusPending  = 1 // 待核销
	ExchangeStatusVerified = 2 // 已核销
	ExchangeStatusCancelled = 3 // 已取消
)
