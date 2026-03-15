package points

import (
	"context"

	"github.com/YazaiHu/MemberHub/internal/domain/points"
)

// Service 积分应用服务
type Service struct {
	pointsService *points.Service
	pointsRepo    points.Repository
}

// NewService 创建积分应用服务
func NewService(pointsRepo points.Repository) *Service {
	return &Service{
		pointsService: points.NewService(pointsRepo),
		pointsRepo:    pointsRepo,
	}
}

// GetUserPoints 获取用户积分
func (s *Service) GetUserPoints(ctx context.Context, userID int64) (*points.MemberPoints, error) {
	return s.pointsService.GetUserPoints(ctx, userID)
}

// GetPointsHistory 获取积分历史
func (s *Service) GetPointsHistory(ctx context.Context, userID int64, page, pageSize int) ([]*points.PointsTransaction, int64, error) {
	return s.pointsService.GetPointsTransactions(ctx, userID, page, pageSize)
}

// ListExchangeRules 兑换规则列表
func (s *Service) ListExchangeRules(ctx context.Context, storeID *int64, page, pageSize int) ([]*points.ExchangeRule, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	status := int8(1) // 只显示上架的
	offset := (page - 1) * pageSize
	return s.pointsRepo.ListRules(ctx, storeID, &status, offset, pageSize)
}

// ExchangePoints 积分兑换
func (s *Service) ExchangePoints(ctx context.Context, userID, ruleID int64) (*points.ExchangeRecord, error) {
	return s.pointsService.ExchangePoints(ctx, userID, ruleID)
}

// GetExchangeRecords 获取兑换记录
func (s *Service) GetExchangeRecords(ctx context.Context, userID int64, page, pageSize int) ([]*points.ExchangeRecord, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.pointsRepo.ListRecords(ctx, &userID, nil, nil, offset, pageSize)
}

// VerifyExchange 核销兑换
func (s *Service) VerifyExchange(ctx context.Context, verifyCode string, verifierID int64) error {
	return s.pointsService.VerifyExchange(ctx, verifyCode, verifierID)
}

// --- 管理员功能 ---

// ManualAdjustPoints 手动调整积分
func (s *Service) ManualAdjustPoints(ctx context.Context, userID int64, pointsAmount int64, remark string, operatorID int64) error {
	if pointsAmount > 0 {
		return s.pointsService.AddPoints(ctx, userID, pointsAmount, points.PointsTypeAdjust, "manual_adjust", &operatorID, remark, 0)
	} else if pointsAmount < 0 {
		return s.pointsService.DeductPoints(ctx, userID, -pointsAmount, points.PointsTypeAdjust, "manual_adjust", &operatorID, remark)
	}
	return nil
}

// CreateExchangeRule 创建兑换规则
func (s *Service) CreateExchangeRule(ctx context.Context, rule *points.ExchangeRule) error {
	// 初始化剩余库存
	if rule.TotalStock != nil {
		rule.RemainingStock = rule.TotalStock

		// 同步到Redis
		points.SetStockToRedis(rule.ID, *rule.RemainingStock)
	}

	return s.pointsRepo.CreateRule(ctx, rule)
}

// UpdateExchangeRule 更新兑换规则
func (s *Service) UpdateExchangeRule(ctx context.Context, rule *points.ExchangeRule) error {
	// 同步库存到Redis
	if rule.RemainingStock != nil {
		points.SetStockToRedis(rule.ID, *rule.RemainingStock)
	}

	return s.pointsRepo.UpdateRule(ctx, rule)
}

// CancelExchange 取消兑换
func (s *Service) CancelExchange(ctx context.Context, recordID int64, reason string) error {
	return s.pointsService.CancelExchange(ctx, recordID, reason)
}
