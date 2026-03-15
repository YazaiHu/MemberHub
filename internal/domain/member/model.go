package member

import (
	"time"
)

// User 用户模型
type User struct {
	ID             int64      `gorm:"column:id;primaryKey" json:"id"`
	OpenID         string     `gorm:"column:openid;uniqueIndex;size:64;not null" json:"openid"`
	UnionID        string     `gorm:"column:unionid;size:64" json:"unionid,omitempty"`
	Nickname       string     `gorm:"column:nickname;size:100" json:"nickname,omitempty"`
	Avatar         string     `gorm:"column:avatar;size:255" json:"avatar,omitempty"`
	Phone          string     `gorm:"column:phone;size:20;index" json:"phone,omitempty"`
	Gender         int8       `gorm:"column:gender;default:0" json:"gender"` // 0-未知, 1-男, 2-女
	Birthday       *time.Time `gorm:"column:birthday;type:date" json:"birthday,omitempty"`
	Status         int8       `gorm:"column:status;default:1;index" json:"status"` // 0-禁用, 1-正常
	RegisterTime   time.Time  `gorm:"column:register_time;default:CURRENT_TIMESTAMP" json:"register_time"`
	LastLoginTime  *time.Time `gorm:"column:last_login_time" json:"last_login_time,omitempty"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Admin 管理员模型
type Admin struct {
	ID            int64      `gorm:"column:id;primaryKey" json:"id"`
	Username      string     `gorm:"column:username;uniqueIndex;size:50;not null" json:"username"`
	Password      string     `gorm:"column:password;size:255;not null" json:"-"` // 不返回给前端
	RealName      string     `gorm:"column:real_name;size:50" json:"real_name,omitempty"`
	Phone         string     `gorm:"column:phone;size:20" json:"phone,omitempty"`
	Email         string     `gorm:"column:email;size:100" json:"email,omitempty"`
	StoreID       *int64     `gorm:"column:store_id;index" json:"store_id,omitempty"`
	RoleID        *int64     `gorm:"column:role_id;index" json:"role_id,omitempty"`
	Status        int8       `gorm:"column:status;default:1;index" json:"status"` // 0-禁用, 1-正常
	LastLoginTime *time.Time `gorm:"column:last_login_time" json:"last_login_time,omitempty"`
	LastLoginIP   string     `gorm:"column:last_login_ip;size:50" json:"last_login_ip,omitempty"`
	CreatedAt     time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at" json:"updated_at"`

	// 关联
	Role *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admins"
}

// Role 角色模型
type Role struct {
	ID          int64     `gorm:"column:id;primaryKey" json:"id"`
	Name        string    `gorm:"column:name;size:50;not null" json:"name"`
	Code        string    `gorm:"column:code;uniqueIndex;size:50;not null" json:"code"`
	Description string    `gorm:"column:description;size:255" json:"description,omitempty"`
	Status      int8      `gorm:"column:status;default:1" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限模型
type Permission struct {
	ID          int64     `gorm:"column:id;primaryKey" json:"id"`
	Name        string    `gorm:"column:name;size:50;not null" json:"name"`
	Code        string    `gorm:"column:code;uniqueIndex;size:100;not null" json:"code"`
	Resource    string    `gorm:"column:resource;size:100;not null" json:"resource"`
	Action      string    `gorm:"column:action;size:50;not null" json:"action"`
	Description string    `gorm:"column:description;size:255" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联模型
type RolePermission struct {
	ID           int64     `gorm:"column:id;primaryKey" json:"id"`
	RoleID       int64     `gorm:"column:role_id;not null;index" json:"role_id"`
	PermissionID int64     `gorm:"column:permission_id;not null;index" json:"permission_id"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserAsset 用户资产概览
type UserAsset struct {
	UserID          int64 `json:"user_id"`
	AvailablePoints int64 `json:"available_points"`    // 可用积分
	Balance         int64 `json:"balance"`             // 余额（分）
	CouponCount     int   `json:"coupon_count"`        // 可用优惠券数量
	TotalRecharge   int64 `json:"total_recharge"`      // 累计充值（分）
	TotalConsume    int64 `json:"total_consume"`       // 累计消费（分）
}
