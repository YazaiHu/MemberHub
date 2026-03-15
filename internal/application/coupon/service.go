package coupon

import (
	"context"

	"github.com/YazaiHu/MemberHub/internal/domain/coupon"
)

// Service 优惠券应用服务
type Service struct {
	couponService *coupon.Service
	couponRepo    coupon.Repository
}

// NewService 创建优惠券应用服务
func NewService(couponRepo coupon.Repository) *Service {
	return &Service{
		couponService: coupon.NewService(couponRepo),
		couponRepo:    couponRepo,
	}
}

// GetActiveTemplates 获取有效的优惠券模板
func (s *Service) GetActiveTemplates(ctx context.Context) ([]*coupon.CouponTemplate, error) {
	return s.couponRepo.GetActiveTemplates(ctx)
}

// ReceiveCoupon 领取优惠券
func (s *Service) ReceiveCoupon(ctx context.Context, userID, templateID int64) (*coupon.UserCoupon, error) {
	return s.couponService.ReceiveCoupon(ctx, userID, templateID)
}

// GetAvailableCoupons 获取可用优惠券
func (s *Service) GetAvailableCoupons(ctx context.Context, userID int64, page, pageSize int) ([]*coupon.UserCoupon, int64, error) {
	return s.couponService.GetUserAvailableCoupons(ctx, userID, page, pageSize)
}

// GetAllCoupons 获取所有优惠券
func (s *Service) GetAllCoupons(ctx context.Context, userID int64, page, pageSize int) ([]*coupon.UserCoupon, int64, error) {
	return s.couponService.GetUserAllCoupons(ctx, userID, page, pageSize)
}

// UseCoupon 使用优惠券
func (s *Service) UseCoupon(ctx context.Context, couponID, orderID int64) error {
	return s.couponService.UseCoupon(ctx, couponID, orderID)
}

// --- 管理员功能 ---

// ListTemplates 优惠券模板列表
func (s *Service) ListTemplates(ctx context.Context, page, pageSize int) ([]*coupon.CouponTemplate, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.couponRepo.ListTemplates(ctx, nil, offset, pageSize)
}

// CreateTemplate 创建优惠券模板
func (s *Service) CreateTemplate(ctx context.Context, template *coupon.CouponTemplate) error {
	// 初始化剩余数量
	template.RemainingQuantity = template.TotalQuantity
	return s.couponRepo.CreateTemplate(ctx, template)
}

// UpdateTemplate 更新优惠券模板
func (s *Service) UpdateTemplate(ctx context.Context, template *coupon.CouponTemplate) error {
	return s.couponRepo.UpdateTemplate(ctx, template)
}

// GetTemplate 获取优惠券模板
func (s *Service) GetTemplate(ctx context.Context, templateID int64) (*coupon.CouponTemplate, error) {
	return s.couponRepo.GetTemplateByID(ctx, templateID)
}

// CreatePushTask 创建推送任务
func (s *Service) CreatePushTask(ctx context.Context, req *coupon.PushCouponRequest, createdBy int64) (*coupon.CouponPushTask, error) {
	return s.couponService.CreatePushTask(ctx, req, createdBy)
}

// ListPushTasks 推送任务列表
func (s *Service) ListPushTasks(ctx context.Context, page, pageSize int) ([]*coupon.CouponPushTask, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.couponRepo.ListPushTasks(ctx, offset, pageSize)
}

// GetPushTask 获取推送任务
func (s *Service) GetPushTask(ctx context.Context, taskID int64) (*coupon.CouponPushTask, error) {
	return s.couponRepo.GetPushTaskByID(ctx, taskID)
}

// UpdatePushTask 更新推送任务
func (s *Service) UpdatePushTask(ctx context.Context, task *coupon.CouponPushTask) error {
	return s.couponRepo.UpdatePushTask(ctx, task)
}

// IssueCouponToUser 给指定用户发放优惠券（批量推送使用）
func (s *Service) IssueCouponToUser(ctx context.Context, userID, templateID int64) (*coupon.UserCoupon, error) {
	return s.couponService.IssueCouponToUser(ctx, userID, templateID)
}
