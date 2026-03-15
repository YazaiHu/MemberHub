package recharge

import (
	"context"

	"github.com/YazaiHu/MemberHub/internal/domain/points"
	"github.com/YazaiHu/MemberHub/internal/domain/recharge"
)

// Service 充值应用服务
type Service struct {
	rechargeService *recharge.Service
	rechargeRepo    recharge.Repository
}

// NewService 创建充值应用服务
func NewService(rechargeRepo recharge.Repository, pointsRepo points.Repository) *Service {
	return &Service{
		rechargeService: recharge.NewService(rechargeRepo, pointsRepo),
		rechargeRepo:    rechargeRepo,
	}
}

// GetActivePromotions 获取有效的充值活动
func (s *Service) GetActivePromotions(ctx context.Context) ([]*recharge.RechargePromotion, error) {
	return s.rechargeRepo.GetActivePromotions(ctx)
}

// ListPromotions 充值活动列表
func (s *Service) ListPromotions(ctx context.Context, page, pageSize int) ([]*recharge.RechargePromotion, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	status := int8(1) // 只显示启用的
	offset := (page - 1) * pageSize
	return s.rechargeRepo.ListPromotions(ctx, &status, offset, pageSize)
}

// CreateOrder 创建充值订单
func (s *Service) CreateOrder(ctx context.Context, userID int64, req *recharge.CreateOrderRequest) (*recharge.RechargeOrder, error) {
	return s.rechargeService.CreateOrder(ctx, userID, req)
}

// GetUserBalance 获取用户余额
func (s *Service) GetUserBalance(ctx context.Context, userID int64) (*recharge.MemberBalance, error) {
	return s.rechargeService.GetUserBalance(ctx, userID)
}

// GetBalanceHistory 获取余额历史
func (s *Service) GetBalanceHistory(ctx context.Context, userID int64, page, pageSize int) ([]*recharge.BalanceTransaction, int64, error) {
	return s.rechargeService.GetBalanceTransactions(ctx, userID, page, pageSize)
}

// GetRechargeRecords 获取充值记录
func (s *Service) GetRechargeRecords(ctx context.Context, userID int64, page, pageSize int) ([]*recharge.RechargeOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	status := int8(recharge.OrderStatusPaid) // 只显示已支付的
	offset := (page - 1) * pageSize
	return s.rechargeRepo.ListOrders(ctx, &userID, &status, offset, pageSize)
}

// ProcessPaymentCallback 处理支付回调
func (s *Service) ProcessPaymentCallback(ctx context.Context, orderNo, transactionID string, paidAmount int) error {
	return s.rechargeService.ProcessPaymentCallback(ctx, orderNo, transactionID, paidAmount)
}

// --- 管理员功能 ---

// ManualAdjustBalance 手动调整余额
func (s *Service) ManualAdjustBalance(ctx context.Context, userID int64, amount int64, remark string, operatorID int64) error {
	if amount > 0 {
		return s.rechargeService.AddBalance(ctx, userID, amount, remark, operatorID)
	} else if amount < 0 {
		return s.rechargeService.DeductBalance(ctx, userID, -amount, recharge.BalanceTypeAdjust, "manual_adjust", &operatorID, remark)
	}
	return nil
}

// CreatePromotion 创建充值活动
func (s *Service) CreatePromotion(ctx context.Context, promotion *recharge.RechargePromotion) error {
	return s.rechargeRepo.CreatePromotion(ctx, promotion)
}

// UpdatePromotion 更新充值活动
func (s *Service) UpdatePromotion(ctx context.Context, promotion *recharge.RechargePromotion) error {
	return s.rechargeRepo.UpdatePromotion(ctx, promotion)
}

// ListAllPromotions 所有充值活动列表（管理员）
func (s *Service) ListAllPromotions(ctx context.Context, page, pageSize int) ([]*recharge.RechargePromotion, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.rechargeRepo.ListPromotions(ctx, nil, offset, pageSize)
}

// ListAllOrders 所有充值订单列表（管理员）
func (s *Service) ListAllOrders(ctx context.Context, userID *int64, status *int8, page, pageSize int) ([]*recharge.RechargeOrder, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.rechargeRepo.ListOrders(ctx, userID, status, offset, pageSize)
}
