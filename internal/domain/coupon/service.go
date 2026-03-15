package coupon

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/redis"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
	"github.com/YazaiHu/MemberHub/internal/pkg/utils"
)

// Service 优惠券领域服务
type Service struct {
	repo Repository
}

// NewService 创建优惠券服务
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// ReceiveCoupon 用户领取优惠券
func (s *Service) ReceiveCoupon(ctx context.Context, userID, templateID int64) (*UserCoupon, error) {
	// 1. 获取优惠券模板
	template, err := s.repo.GetTemplateByID(ctx, templateID)
	if err != nil {
		return nil, err
	}

	// 2. 检查模板状态
	if template.Status != 1 {
		return nil, errors.ErrCouponNotFound.WithMessage("coupon template is not active")
	}

	// 3. 检查时间范围
	now := time.Now()
	if template.StartTime != nil && now.Before(*template.StartTime) {
		return nil, errors.ErrInvalidParams.WithMessage("coupon not available yet")
	}
	if template.EndTime != nil && now.After(*template.EndTime) {
		return nil, errors.ErrInvalidParams.WithMessage("coupon has expired")
	}

	// 4. 检查用户领取限制
	count, err := s.repo.CountUserCouponsByTemplate(ctx, userID, templateID)
	if err != nil {
		return nil, err
	}
	if count >= template.PerUserLimit {
		return nil, errors.ErrCouponLimitExceeded
	}

	// 5. 分布式锁防止重复领取
	lockKey := fmt.Sprintf("lock:coupon:receive:%d:%d", userID, templateID)
	lock, acquired, err := redis.AcquireLock(lockKey, 5*time.Second)
	if err != nil {
		return nil, err
	}
	if !acquired {
		return nil, errors.ErrDuplicateOperation
	}
	defer lock.Release()

	// 6. 使用事务
	var userCoupon *UserCoupon
	err = mysql.Transaction(func(tx *gorm.DB) error {
		// 6.1 原子减库存
		if err := s.repo.DecrTemplateQuantity(ctx, templateID); err != nil {
			return err
		}

		// 6.2 生成优惠券
		couponCode := utils.GenerateCouponCode()
		expiredAt := time.Now().AddDate(0, 0, template.ValidDays)

		userCoupon = &UserCoupon{
			CouponCode:       couponCode,
			UserID:           userID,
			TemplateID:       templateID,
			Name:             template.Name,
			Type:             template.Type,
			DiscountType:     template.DiscountType,
			DiscountValue:    template.DiscountValue,
			MinAmount:        template.MinAmount,
			ApplicableStores: template.ApplicableStores,
			Status:           CouponStatusUnused,
			ExpiredAt:        expiredAt,
		}

		if err := s.repo.CreateUserCoupon(ctx, userCoupon); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return userCoupon, nil
}

// UseCoupon 使用优惠券
func (s *Service) UseCoupon(ctx context.Context, couponID, orderID int64) error {
	coupon, err := s.repo.GetUserCouponByID(ctx, couponID)
	if err != nil {
		return err
	}

	// 检查状态
	if coupon.Status != CouponStatusUnused {
		return errors.ErrCouponUsed
	}

	// 检查是否过期
	if time.Now().After(coupon.ExpiredAt) {
		return errors.ErrCouponExpired
	}

	// 更新状态
	now := time.Now()
	coupon.Status = CouponStatusUsed
	coupon.UsedAt = &now
	coupon.UsedOrderID = &orderID

	return s.repo.UpdateUserCoupon(ctx, coupon)
}

// CreatePushTask 创建推送任务
func (s *Service) CreatePushTask(ctx context.Context, req *PushCouponRequest, createdBy int64) (*CouponPushTask, error) {
	// 验证模板
	template, err := s.repo.GetTemplateByID(ctx, req.TemplateID)
	if err != nil {
		return nil, err
	}

	if template.Status != 1 {
		return nil, errors.ErrInvalidParams.WithMessage("template is not active")
	}

	// 生成任务号
	taskNo := utils.GenerateOrderNo("CPT")

	// 构造任务
	task := &CouponPushTask{
		TaskNo:        taskNo,
		TemplateID:    req.TemplateID,
		Title:         req.Title,
		TargetType:    req.TargetType,
		Status:        TaskStatusPending,
		SendWechatMsg: 0,
		CreatedBy:     createdBy,
	}

	if req.SendWechatMsg {
		task.SendWechatMsg = 1
	}

	// 处理目标用户
	if req.TargetType == TargetTypeSpecified {
		if len(req.TargetUserIDs) == 0 {
			return nil, errors.ErrInvalidParams.WithMessage("target_user_ids is required")
		}

		// 序列化用户ID列表
		targetUsers, _ := json.Marshal(req.TargetUserIDs)
		task.TargetUsers = string(targetUsers)
		task.TotalCount = len(req.TargetUserIDs)
	}

	// 创建任务
	if err := s.repo.CreatePushTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// IssueCouponToUser 给指定用户发放优惠券（批量推送使用）
func (s *Service) IssueCouponToUser(ctx context.Context, userID, templateID int64) (*UserCoupon, error) {
	// 获取模板
	template, err := s.repo.GetTemplateByID(ctx, templateID)
	if err != nil {
		return nil, err
	}

	// 生成优惠券（不检查领取限制，批量推送时管理员主动发放）
	couponCode := utils.GenerateCouponCode()
	expiredAt := time.Now().AddDate(0, 0, template.ValidDays)

	userCoupon := &UserCoupon{
		CouponCode:       couponCode,
		UserID:           userID,
		TemplateID:       templateID,
		Name:             template.Name,
		Type:             template.Type,
		DiscountType:     template.DiscountType,
		DiscountValue:    template.DiscountValue,
		MinAmount:        template.MinAmount,
		ApplicableStores: template.ApplicableStores,
		Status:           CouponStatusUnused,
		ExpiredAt:        expiredAt,
	}

	if err := s.repo.CreateUserCoupon(ctx, userCoupon); err != nil {
		return nil, err
	}

	return userCoupon, nil
}

// GetUserAvailableCoupons 获取用户可用优惠券
func (s *Service) GetUserAvailableCoupons(ctx context.Context, userID int64, page, pageSize int) ([]*UserCoupon, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	status := int8(CouponStatusUnused)
	offset := (page - 1) * pageSize
	return s.repo.ListUserCoupons(ctx, userID, &status, offset, pageSize)
}

// GetUserAllCoupons 获取用户所有优惠券
func (s *Service) GetUserAllCoupons(ctx context.Context, userID int64, page, pageSize int) ([]*UserCoupon, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.repo.ListUserCoupons(ctx, userID, nil, offset, pageSize)
}

// MarkExpiredCoupons 标记过期优惠券（定时任务调用）
func (s *Service) MarkExpiredCoupons(ctx context.Context) (int64, error) {
	return s.repo.MarkExpiredCoupons(ctx)
}
