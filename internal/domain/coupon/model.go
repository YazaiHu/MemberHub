package coupon

import (
	"time"
)

// CouponTemplate 优惠券模板
type CouponTemplate struct {
	ID               int64      `gorm:"column:id;primaryKey" json:"id"`
	Name             string     `gorm:"column:name;size:100;not null" json:"name"`
	Type             int8       `gorm:"column:type;not null;index" json:"type"` // 1-满减券, 2-折扣券, 3-代金券
	DiscountType     int8       `gorm:"column:discount_type;not null" json:"discount_type"` // 1-金额, 2-折扣
	DiscountValue    int        `gorm:"column:discount_value;not null" json:"discount_value"` // 金额单位分，折扣单位0.1%
	MinAmount        int        `gorm:"column:min_amount;default:0" json:"min_amount"` // 最低消费金额（分）
	TotalQuantity    int        `gorm:"column:total_quantity;not null" json:"total_quantity"` // 总发行量
	RemainingQuantity int       `gorm:"column:remaining_quantity;not null" json:"remaining_quantity"` // 剩余数量
	PerUserLimit     int        `gorm:"column:per_user_limit;default:1" json:"per_user_limit"` // 每人限领数量
	ValidDays        int        `gorm:"column:valid_days;not null" json:"valid_days"` // 有效天数
	ApplicableStores string     `gorm:"column:applicable_stores;type:json" json:"applicable_stores,omitempty"` // 适用门店ID列表（JSON）
	Description      string     `gorm:"column:description;size:255" json:"description,omitempty"`
	StartTime        *time.Time `gorm:"column:start_time" json:"start_time,omitempty"`
	EndTime          *time.Time `gorm:"column:end_time" json:"end_time,omitempty"`
	Status           int8       `gorm:"column:status;default:1;index" json:"status"` // 0-停用, 1-启用
	SortOrder        int        `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (CouponTemplate) TableName() string {
	return "coupon_templates"
}

// 优惠券类型常量
const (
	CouponTypeFullReduction = 1 // 满减券
	CouponTypeDiscount      = 2 // 折扣券
	CouponTypeCash          = 3 // 代金券
)

// 优惠类型常量
const (
	DiscountTypeAmount  = 1 // 金额
	DiscountTypePercent = 2 // 折扣
)

// UserCoupon 用户优惠券
type UserCoupon struct {
	ID               int64      `gorm:"column:id;primaryKey" json:"id"`
	CouponCode       string     `gorm:"column:coupon_code;uniqueIndex;size:32;not null" json:"coupon_code"`
	UserID           int64      `gorm:"column:user_id;index;not null" json:"user_id"`
	TemplateID       int64      `gorm:"column:template_id;index;not null" json:"template_id"`
	Name             string     `gorm:"column:name;size:100;not null" json:"name"`
	Type             int8       `gorm:"column:type;not null" json:"type"`
	DiscountType     int8       `gorm:"column:discount_type;not null" json:"discount_type"`
	DiscountValue    int        `gorm:"column:discount_value;not null" json:"discount_value"`
	MinAmount        int        `gorm:"column:min_amount;default:0" json:"min_amount"`
	ApplicableStores string     `gorm:"column:applicable_stores;type:json" json:"applicable_stores,omitempty"`
	Status           int8       `gorm:"column:status;default:1;index" json:"status"` // 1-未使用, 2-已使用, 3-已过期
	UsedAt           *time.Time `gorm:"column:used_at" json:"used_at,omitempty"`
	UsedOrderID      *int64     `gorm:"column:used_order_id" json:"used_order_id,omitempty"`
	ExpiredAt        time.Time  `gorm:"column:expired_at;index;not null" json:"expired_at"`
	ReceivedAt       time.Time  `gorm:"column:received_at;default:CURRENT_TIMESTAMP" json:"received_at"`
	CreatedAt        time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (UserCoupon) TableName() string {
	return "user_coupons"
}

// 用户优惠券状态常量
const (
	CouponStatusUnused  = 1 // 未使用
	CouponStatusUsed    = 2 // 已使用
	CouponStatusExpired = 3 // 已过期
)

// CouponPushTask 优惠券推送任务
type CouponPushTask struct {
	ID               int64      `gorm:"column:id;primaryKey" json:"id"`
	TaskNo           string     `gorm:"column:task_no;uniqueIndex;size:64;not null" json:"task_no"`
	TemplateID       int64      `gorm:"column:template_id;index;not null" json:"template_id"`
	Title            string     `gorm:"column:title;size:100;not null" json:"title"`
	TargetType       int8       `gorm:"column:target_type;not null" json:"target_type"` // 1-全部用户, 2-指定用户, 3-条件筛选
	TargetUsers      string     `gorm:"column:target_users;type:json" json:"target_users,omitempty"` // 目标用户ID列表（JSON）
	TargetConditions string     `gorm:"column:target_conditions;type:json" json:"target_conditions,omitempty"` // 筛选条件（JSON）
	TotalCount       int        `gorm:"column:total_count;default:0" json:"total_count"`
	SuccessCount     int        `gorm:"column:success_count;default:0" json:"success_count"`
	FailCount        int        `gorm:"column:fail_count;default:0" json:"fail_count"`
	Status           int8       `gorm:"column:status;default:1;index" json:"status"` // 1-待执行, 2-执行中, 3-已完成, 4-已取消
	SendWechatMsg    int8       `gorm:"column:send_wechat_msg;default:1" json:"send_wechat_msg"` // 是否发送微信通知
	StartedAt        *time.Time `gorm:"column:started_at" json:"started_at,omitempty"`
	CompletedAt      *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
	CreatedBy        int64      `gorm:"column:created_by;not null" json:"created_by"`
	CreatedAt        time.Time  `gorm:"column:created_at;index" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (CouponPushTask) TableName() string {
	return "coupon_push_tasks"
}

// 推送任务状态常量
const (
	TaskStatusPending    = 1 // 待执行
	TaskStatusProcessing = 2 // 执行中
	TaskStatusCompleted  = 3 // 已完成
	TaskStatusCancelled  = 4 // 已取消
)

// 推送目标类型常量
const (
	TargetTypeAll       = 1 // 全部用户
	TargetTypeSpecified = 2 // 指定用户
	TargetTypeCondition = 3 // 条件筛选
)

// ReceiveCouponRequest 领取优惠券请求
type ReceiveCouponRequest struct {
	TemplateID int64 `json:"template_id" binding:"required"`
}

// PushCouponRequest 推送优惠券请求
type PushCouponRequest struct {
	TemplateID       int64   `json:"template_id" binding:"required"`
	Title            string  `json:"title" binding:"required"`
	TargetType       int8    `json:"target_type" binding:"required,oneof=1 2 3"`
	TargetUserIDs    []int64 `json:"target_user_ids,omitempty"` // TargetType=2时必填
	SendWechatMsg    bool    `json:"send_wechat_msg"`
}
