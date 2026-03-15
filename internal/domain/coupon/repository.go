package coupon

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// Repository 优惠券仓储接口
type Repository interface {
	// CouponTemplate
	CreateTemplate(ctx context.Context, template *CouponTemplate) error
	GetTemplateByID(ctx context.Context, templateID int64) (*CouponTemplate, error)
	UpdateTemplate(ctx context.Context, template *CouponTemplate) error
	ListTemplates(ctx context.Context, status *int8, offset, limit int) ([]*CouponTemplate, int64, error)
	GetActiveTemplates(ctx context.Context) ([]*CouponTemplate, error)
	DecrTemplateQuantity(ctx context.Context, templateID int64) error // 原子减库存

	// UserCoupon
	CreateUserCoupon(ctx context.Context, coupon *UserCoupon) error
	GetUserCouponByID(ctx context.Context, couponID int64) (*UserCoupon, error)
	GetUserCouponByCode(ctx context.Context, couponCode string) (*UserCoupon, error)
	UpdateUserCoupon(ctx context.Context, coupon *UserCoupon) error
	ListUserCoupons(ctx context.Context, userID int64, status *int8, offset, limit int) ([]*UserCoupon, int64, error)
	CountUserCouponsByTemplate(ctx context.Context, userID, templateID int64) (int, error)
	CountAvailableCoupons(ctx context.Context, userID int64) (int, error)
	MarkExpiredCoupons(ctx context.Context) (int64, error) // 标记过期优惠券

	// CouponPushTask
	CreatePushTask(ctx context.Context, task *CouponPushTask) error
	GetPushTaskByID(ctx context.Context, taskID int64) (*CouponPushTask, error)
	GetPushTaskByNo(ctx context.Context, taskNo string) (*CouponPushTask, error)
	UpdatePushTask(ctx context.Context, task *CouponPushTask) error
	ListPushTasks(ctx context.Context, offset, limit int) ([]*CouponPushTask, int64, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository 创建优惠券仓储
func NewRepository() Repository {
	return &repository{
		db: mysql.GetDB(),
	}
}

// CreateTemplate 创建优惠券模板
func (r *repository) CreateTemplate(ctx context.Context, template *CouponTemplate) error {
	if err := r.db.WithContext(ctx).Create(template).Error; err != nil {
		return fmt.Errorf("create coupon template failed: %w", err)
	}
	return nil
}

// GetTemplateByID 获取优惠券模板
func (r *repository) GetTemplateByID(ctx context.Context, templateID int64) (*CouponTemplate, error) {
	var template CouponTemplate
	if err := r.db.WithContext(ctx).Where("id = ?", templateID).First(&template).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCouponNotFound
		}
		return nil, fmt.Errorf("get coupon template failed: %w", err)
	}
	return &template, nil
}

// UpdateTemplate 更新优惠券模板
func (r *repository) UpdateTemplate(ctx context.Context, template *CouponTemplate) error {
	if err := r.db.WithContext(ctx).Save(template).Error; err != nil {
		return fmt.Errorf("update coupon template failed: %w", err)
	}
	return nil
}

// ListTemplates 优惠券模板列表
func (r *repository) ListTemplates(ctx context.Context, status *int8, offset, limit int) ([]*CouponTemplate, int64, error) {
	var templates []*CouponTemplate
	var total int64

	query := r.db.WithContext(ctx).Model(&CouponTemplate{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count coupon templates failed: %w", err)
	}

	if err := query.Order("sort_order DESC, id DESC").Offset(offset).Limit(limit).Find(&templates).Error; err != nil {
		return nil, 0, fmt.Errorf("list coupon templates failed: %w", err)
	}

	return templates, total, nil
}

// GetActiveTemplates 获取有效的优惠券模板
func (r *repository) GetActiveTemplates(ctx context.Context) ([]*CouponTemplate, error) {
	var templates []*CouponTemplate
	query := r.db.WithContext(ctx).Where("status = 1 AND remaining_quantity > 0")

	// 时间范围检查
	query = query.Where("(start_time IS NULL OR start_time <= NOW()) AND (end_time IS NULL OR end_time >= NOW())")

	if err := query.Order("sort_order DESC, id DESC").Find(&templates).Error; err != nil {
		return nil, fmt.Errorf("get active coupon templates failed: %w", err)
	}

	return templates, nil
}

// DecrTemplateQuantity 原子减库存
func (r *repository) DecrTemplateQuantity(ctx context.Context, templateID int64) error {
	result := r.db.WithContext(ctx).Model(&CouponTemplate{}).
		Where("id = ? AND remaining_quantity > 0", templateID).
		Update("remaining_quantity", gorm.Expr("remaining_quantity - 1"))

	if result.Error != nil {
		return fmt.Errorf("decr template quantity failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrCouponSoldOut
	}

	return nil
}

// CreateUserCoupon 创建用户优惠券
func (r *repository) CreateUserCoupon(ctx context.Context, coupon *UserCoupon) error {
	if err := r.db.WithContext(ctx).Create(coupon).Error; err != nil {
		return fmt.Errorf("create user coupon failed: %w", err)
	}
	return nil
}

// GetUserCouponByID 获取用户优惠券
func (r *repository) GetUserCouponByID(ctx context.Context, couponID int64) (*UserCoupon, error) {
	var coupon UserCoupon
	if err := r.db.WithContext(ctx).Where("id = ?", couponID).First(&coupon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCouponNotFound
		}
		return nil, fmt.Errorf("get user coupon failed: %w", err)
	}
	return &coupon, nil
}

// GetUserCouponByCode 根据优惠券码获取
func (r *repository) GetUserCouponByCode(ctx context.Context, couponCode string) (*UserCoupon, error) {
	var coupon UserCoupon
	if err := r.db.WithContext(ctx).Where("coupon_code = ?", couponCode).First(&coupon).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrCouponNotFound
		}
		return nil, fmt.Errorf("get user coupon failed: %w", err)
	}
	return &coupon, nil
}

// UpdateUserCoupon 更新用户优惠券
func (r *repository) UpdateUserCoupon(ctx context.Context, coupon *UserCoupon) error {
	if err := r.db.WithContext(ctx).Save(coupon).Error; err != nil {
		return fmt.Errorf("update user coupon failed: %w", err)
	}
	return nil
}

// ListUserCoupons 用户优惠券列表
func (r *repository) ListUserCoupons(ctx context.Context, userID int64, status *int8, offset, limit int) ([]*UserCoupon, int64, error) {
	var coupons []*UserCoupon
	var total int64

	query := r.db.WithContext(ctx).Model(&UserCoupon{}).Where("user_id = ?", userID)
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count user coupons failed: %w", err)
	}

	if err := query.Order("status ASC, expired_at ASC").Offset(offset).Limit(limit).Find(&coupons).Error; err != nil {
		return nil, 0, fmt.Errorf("list user coupons failed: %w", err)
	}

	return coupons, total, nil
}

// CountUserCouponsByTemplate 统计用户领取某模板的优惠券数量
func (r *repository) CountUserCouponsByTemplate(ctx context.Context, userID, templateID int64) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserCoupon{}).
		Where("user_id = ? AND template_id = ?", userID, templateID).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count user coupons by template failed: %w", err)
	}
	return int(count), nil
}

// CountAvailableCoupons 统计用户可用优惠券数量
func (r *repository) CountAvailableCoupons(ctx context.Context, userID int64) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&UserCoupon{}).
		Where("user_id = ? AND status = ? AND expired_at > NOW()", userID, CouponStatusUnused).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count available coupons failed: %w", err)
	}
	return int(count), nil
}

// MarkExpiredCoupons 标记过期优惠券
func (r *repository) MarkExpiredCoupons(ctx context.Context) (int64, error) {
	result := r.db.WithContext(ctx).Model(&UserCoupon{}).
		Where("status = ? AND expired_at <= NOW()", CouponStatusUnused).
		Update("status", CouponStatusExpired)

	if result.Error != nil {
		return 0, fmt.Errorf("mark expired coupons failed: %w", result.Error)
	}

	return result.RowsAffected, nil
}

// CreatePushTask 创建推送任务
func (r *repository) CreatePushTask(ctx context.Context, task *CouponPushTask) error {
	if err := r.db.WithContext(ctx).Create(task).Error; err != nil {
		return fmt.Errorf("create push task failed: %w", err)
	}
	return nil
}

// GetPushTaskByID 获取推送任务
func (r *repository) GetPushTaskByID(ctx context.Context, taskID int64) (*CouponPushTask, error) {
	var task CouponPushTask
	if err := r.db.WithContext(ctx).Where("id = ?", taskID).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get push task failed: %w", err)
	}
	return &task, nil
}

// GetPushTaskByNo 根据任务号获取
func (r *repository) GetPushTaskByNo(ctx context.Context, taskNo string) (*CouponPushTask, error) {
	var task CouponPushTask
	if err := r.db.WithContext(ctx).Where("task_no = ?", taskNo).First(&task).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get push task failed: %w", err)
	}
	return &task, nil
}

// UpdatePushTask 更新推送任务
func (r *repository) UpdatePushTask(ctx context.Context, task *CouponPushTask) error {
	if err := r.db.WithContext(ctx).Save(task).Error; err != nil {
		return fmt.Errorf("update push task failed: %w", err)
	}
	return nil
}

// ListPushTasks 推送任务列表
func (r *repository) ListPushTasks(ctx context.Context, offset, limit int) ([]*CouponPushTask, int64, error) {
	var tasks []*CouponPushTask
	var total int64

	query := r.db.WithContext(ctx).Model(&CouponPushTask{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count push tasks failed: %w", err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&tasks).Error; err != nil {
		return nil, 0, fmt.Errorf("list push tasks failed: %w", err)
	}

	return tasks, total, nil
}
