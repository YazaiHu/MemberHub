package member

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// Repository 会员仓储接口
type Repository interface {
	// User
	CreateUser(ctx context.Context, user *User) error
	GetUserByID(ctx context.Context, userID int64) (*User, error)
	GetUserByOpenID(ctx context.Context, openID string) (*User, error)
	GetUserByPhone(ctx context.Context, phone string) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	UpdateLastLoginTime(ctx context.Context, userID int64) error
	ListUsers(ctx context.Context, offset, limit int, status *int8) ([]*User, int64, error)

	// Admin
	CreateAdmin(ctx context.Context, admin *Admin) error
	GetAdminByID(ctx context.Context, adminID int64) (*Admin, error)
	GetAdminByUsername(ctx context.Context, username string) (*Admin, error)
	UpdateAdmin(ctx context.Context, admin *Admin) error
	UpdateAdminLastLogin(ctx context.Context, adminID int64, ip string) error
	ListAdmins(ctx context.Context, offset, limit int, storeID *int64, status *int8) ([]*Admin, int64, error)

	// Role
	GetRoleByID(ctx context.Context, roleID int64) (*Role, error)
	GetRoleByCode(ctx context.Context, code string) (*Role, error)
	ListRoles(ctx context.Context) ([]*Role, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository 创建会员仓储实例
func NewRepository() Repository {
	return &repository{
		db: mysql.GetDB(),
	}
}

// CreateUser 创建用户
func (r *repository) CreateUser(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("create user failed: %w", err)
	}
	return nil
}

// GetUserByID 根据ID获取用户
func (r *repository) GetUserByID(ctx context.Context, userID int64) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id failed: %w", err)
	}
	return &user, nil
}

// GetUserByOpenID 根据OpenID获取用户
func (r *repository) GetUserByOpenID(ctx context.Context, openID string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("openid = ?", openID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by openid failed: %w", err)
	}
	return &user, nil
}

// GetUserByPhone 根据手机号获取用户
func (r *repository) GetUserByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	if err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by phone failed: %w", err)
	}
	return &user, nil
}

// UpdateUser 更新用户
func (r *repository) UpdateUser(ctx context.Context, user *User) error {
	if err := r.db.WithContext(ctx).Save(user).Error; err != nil {
		return fmt.Errorf("update user failed: %w", err)
	}
	return nil
}

// UpdateLastLoginTime 更新最后登录时间
func (r *repository) UpdateLastLoginTime(ctx context.Context, userID int64) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&User{}).
		Where("id = ?", userID).
		Update("last_login_time", now).Error; err != nil {
		return fmt.Errorf("update last login time failed: %w", err)
	}
	return nil
}

// ListUsers 用户列表
func (r *repository) ListUsers(ctx context.Context, offset, limit int, status *int8) ([]*User, int64, error) {
	var users []*User
	var total int64

	query := r.db.WithContext(ctx).Model(&User{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count users failed: %w", err)
	}

	// 查询列表
	if err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("list users failed: %w", err)
	}

	return users, total, nil
}

// CreateAdmin 创建管理员
func (r *repository) CreateAdmin(ctx context.Context, admin *Admin) error {
	if err := r.db.WithContext(ctx).Create(admin).Error; err != nil {
		return fmt.Errorf("create admin failed: %w", err)
	}
	return nil
}

// GetAdminByID 根据ID获取管理员
func (r *repository) GetAdminByID(ctx context.Context, adminID int64) (*Admin, error) {
	var admin Admin
	if err := r.db.WithContext(ctx).Preload("Role").Where("id = ?", adminID).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("get admin by id failed: %w", err)
	}
	return &admin, nil
}

// GetAdminByUsername 根据用户名获取管理员
func (r *repository) GetAdminByUsername(ctx context.Context, username string) (*Admin, error) {
	var admin Admin
	if err := r.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrUserNotFound
		}
		return nil, fmt.Errorf("get admin by username failed: %w", err)
	}
	return &admin, nil
}

// UpdateAdmin 更新管理员
func (r *repository) UpdateAdmin(ctx context.Context, admin *Admin) error {
	if err := r.db.WithContext(ctx).Save(admin).Error; err != nil {
		return fmt.Errorf("update admin failed: %w", err)
	}
	return nil
}

// UpdateAdminLastLogin 更新管理员最后登录信息
func (r *repository) UpdateAdminLastLogin(ctx context.Context, adminID int64, ip string) error {
	now := time.Now()
	if err := r.db.WithContext(ctx).Model(&Admin{}).
		Where("id = ?", adminID).
		Updates(map[string]interface{}{
			"last_login_time": now,
			"last_login_ip":   ip,
		}).Error; err != nil {
		return fmt.Errorf("update admin last login failed: %w", err)
	}
	return nil
}

// ListAdmins 管理员列表
func (r *repository) ListAdmins(ctx context.Context, offset, limit int, storeID *int64, status *int8) ([]*Admin, int64, error) {
	var admins []*Admin
	var total int64

	query := r.db.WithContext(ctx).Model(&Admin{}).Preload("Role")
	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count admins failed: %w", err)
	}

	// 查询列表
	if err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&admins).Error; err != nil {
		return nil, 0, fmt.Errorf("list admins failed: %w", err)
	}

	return admins, total, nil
}

// GetRoleByID 根据ID获取角色
func (r *repository) GetRoleByID(ctx context.Context, roleID int64) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Where("id = ?", roleID).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get role by id failed: %w", err)
	}
	return &role, nil
}

// GetRoleByCode 根据代码获取角色
func (r *repository) GetRoleByCode(ctx context.Context, code string) (*Role, error) {
	var role Role
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get role by code failed: %w", err)
	}
	return &role, nil
}

// ListRoles 角色列表
func (r *repository) ListRoles(ctx context.Context) ([]*Role, error) {
	var roles []*Role
	if err := r.db.WithContext(ctx).Where("status = ?", 1).Find(&roles).Error; err != nil {
		return nil, fmt.Errorf("list roles failed: %w", err)
	}
	return roles, nil
}
