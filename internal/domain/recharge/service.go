package recharge

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/YazaiHu/MemberHub/internal/domain/points"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/redis"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
	"github.com/YazaiHu/MemberHub/internal/pkg/utils"
)

// Service 充值领域服务
type Service struct {
	repo         Repository
	pointsRepo   points.Repository
	pointsService *points.Service
}

// NewService 创建充值服务
func NewService(repo Repository, pointsRepo points.Repository) *Service {
	return &Service{
		repo:         repo,
		pointsRepo:   pointsRepo,
		pointsService: points.NewService(pointsRepo),
	}
}

// CreateOrder 创建充值订单
func (s *Service) CreateOrder(ctx context.Context, userID int64, req *CreateOrderRequest) (*RechargeOrder, error) {
	var promotion *RechargePromotion
	var rechargeAmount int
	var bonusAmount int
	var bonusPoints int

	// 如果指定了活动ID
	if req.PromotionID != nil {
		var err error
		promotion, err = s.repo.GetPromotionByID(ctx, *req.PromotionID)
		if err != nil {
			return nil, err
		}

		// 检查活动状态
		if promotion.Status != 1 {
			return nil, errors.ErrPromotionExpired.WithMessage("promotion is not active")
		}

		// 检查时间范围
		now := time.Now()
		if promotion.StartTime != nil && now.Before(*promotion.StartTime) {
			return nil, errors.ErrPromotionExpired.WithMessage("promotion not started yet")
		}
		if promotion.EndTime != nil && now.After(*promotion.EndTime) {
			return nil, errors.ErrPromotionExpired.WithMessage("promotion has ended")
		}

		// 检查用户限购次数
		if promotion.UserLimit != nil {
			count, err := s.repo.CountUserRechargeByPromotion(ctx, userID, *req.PromotionID)
			if err != nil {
				return nil, err
			}
			if count >= *promotion.UserLimit {
				return nil, errors.ErrConflict.WithMessage("user recharge limit exceeded")
			}
		}

		rechargeAmount = promotion.RechargeAmount * 100 // 转换为分
		bonusAmount = promotion.BonusAmount * 100
		bonusPoints = promotion.BonusPoints
	} else {
		// 不参加活动，直接充值
		if req.Amount <= 0 {
			return nil, errors.ErrInvalidAmount.WithMessage("amount must be greater than 0")
		}
		rechargeAmount = req.Amount * 100 // 转换为分
		bonusAmount = 0
		bonusPoints = 0
	}

	// 创建订单
	orderNo := utils.GenerateOrderNo("RC")
	expiredAt := time.Now().Add(30 * time.Minute) // 30分钟后过期

	order := &RechargeOrder{
		OrderNo:        orderNo,
		UserID:         userID,
		PromotionID:    req.PromotionID,
		RechargeAmount: rechargeAmount,
		BonusAmount:    bonusAmount,
		BonusPoints:    bonusPoints,
		TotalAmount:    rechargeAmount + bonusAmount,
		PayAmount:      rechargeAmount,
		PayMethod:      PayMethodWeChatPay,
		Status:         OrderStatusPending,
		ExpiredAt:      &expiredAt,
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	return order, nil
}

// ProcessPaymentCallback 处理支付回调（关键：保证事务原子性）
func (s *Service) ProcessPaymentCallback(ctx context.Context, orderNo, transactionID string, paidAmount int) error {
	// 使用分布式锁防止重复处理
	lockKey := fmt.Sprintf("lock:recharge:callback:%s", orderNo)
	lock, acquired, err := redis.AcquireLock(lockKey, 10*time.Second)
	if err != nil {
		return err
	}
	if !acquired {
		return errors.ErrDuplicateOperation.WithMessage("callback is being processed")
	}
	defer lock.Release()

	// 使用数据库事务保证原子性
	return mysql.Transaction(func(tx *gorm.DB) error {
		// 1. 获取订单
		order, err := s.repo.GetOrderByNo(ctx, orderNo)
		if err != nil {
			return err
		}

		// 2. 幂等性检查：订单已支付
		if order.Status == OrderStatusPaid {
			return nil // 已处理，直接返回成功
		}

		// 3. 检查订单状态
		if order.Status != OrderStatusPending {
			return errors.ErrInvalidParams.WithMessage("order status invalid")
		}

		// 4. 检查支付金额
		if paidAmount != order.PayAmount {
			return errors.ErrInvalidAmount.WithMessage("paid amount mismatch")
		}

		// 5. 检查订单是否过期
		if order.ExpiredAt != nil && time.Now().After(*order.ExpiredAt) {
			return errors.ErrOrderExpired
		}

		// 6. 幂等性检查：流水表中是否已存在
		transactionNo := orderNo // 使用订单号作为流水号
		_, err = s.repo.GetTransactionByNo(ctx, transactionNo)
		if err == nil {
			// 流水已存在，说明已处理过，直接返回成功
			return nil
		}
		if err != errors.ErrNotFound {
			return err
		}

		// 7. 更新订单状态
		now := time.Now()
		order.Status = OrderStatusPaid
		order.TransactionID = transactionID
		order.PaidAt = &now

		if err := s.repo.UpdateOrder(ctx, order); err != nil {
			return fmt.Errorf("update order failed: %w", err)
		}

		// 8. 增加用户余额（带乐观锁）
		balance, err := s.repo.GetOrCreateBalance(ctx, order.UserID)
		if err != nil {
			return fmt.Errorf("get balance failed: %w", err)
		}

		// 增加余额 = 充值金额 + 赠送金额
		balance.Balance += int64(order.TotalAmount)
		balance.TotalRecharge += int64(order.RechargeAmount)

		if err := s.repo.UpdateBalance(ctx, balance); err != nil {
			return fmt.Errorf("update balance failed: %w", err)
		}

		// 9. 记录余额流水
		balanceTransaction := &BalanceTransaction{
			TransactionNo: transactionNo,
			UserID:        order.UserID,
			Amount:        int64(order.TotalAmount),
			Type:          BalanceTypeRecharge,
			Source:        "recharge",
			RefID:         &order.ID,
			Remark:        fmt.Sprintf("充值：¥%.2f", float64(order.RechargeAmount)/100),
		}

		if err := s.repo.CreateTransaction(ctx, balanceTransaction); err != nil {
			return fmt.Errorf("create balance transaction failed: %w", err)
		}

		// 10. 赠送积分（如果有）
		if order.BonusPoints > 0 {
			if err := s.pointsService.AddPoints(
				ctx,
				order.UserID,
				int64(order.BonusPoints),
				points.PointsTypeRecharge,
				"recharge_bonus",
				&order.ID,
				fmt.Sprintf("充值赠送：%d积分", order.BonusPoints),
				0, // 不设置过期时间
			); err != nil {
				return fmt.Errorf("add bonus points failed: %w", err)
			}
		}

		return nil
	})
}

// CancelOrder 取消订单
func (s *Service) CancelOrder(ctx context.Context, orderID int64) error {
	order, err := s.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return err
	}

	// 只能取消待支付的订单
	if order.Status != OrderStatusPending {
		return errors.ErrInvalidParams.WithMessage("only pending orders can be cancelled")
	}

	order.Status = OrderStatusCancelled
	return s.repo.UpdateOrder(ctx, order)
}

// GetUserBalance 获取用户余额
func (s *Service) GetUserBalance(ctx context.Context, userID int64) (*MemberBalance, error) {
	return s.repo.GetOrCreateBalance(ctx, userID)
}

// GetBalanceTransactions 获取余额流水
func (s *Service) GetBalanceTransactions(ctx context.Context, userID int64, page, pageSize int) ([]*BalanceTransaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.repo.ListTransactions(ctx, userID, offset, pageSize)
}

// AddBalance 增加余额（手动调整）
func (s *Service) AddBalance(ctx context.Context, userID int64, amount int64, remark string, operatorID int64) error {
	if amount <= 0 {
		return errors.ErrInvalidParams.WithMessage("amount must be positive")
	}

	return mysql.Transaction(func(tx *gorm.DB) error {
		// 获取余额
		balance, err := s.repo.GetOrCreateBalance(ctx, userID)
		if err != nil {
			return err
		}

		// 增加余额
		balance.Balance += amount

		if err := s.repo.UpdateBalance(ctx, balance); err != nil {
			return err
		}

		// 记录流水
		transactionNo := utils.GenerateOrderNo("BT")
		transaction := &BalanceTransaction{
			TransactionNo: transactionNo,
			UserID:        userID,
			Amount:        amount,
			Type:          BalanceTypeAdjust,
			Source:        "manual_adjust",
			Remark:        remark,
			OperatorID:    &operatorID,
		}

		return s.repo.CreateTransaction(ctx, transaction)
	})
}

// DeductBalance 扣减余额（消费、手动调整）
func (s *Service) DeductBalance(ctx context.Context, userID int64, amount int64, balanceType int8, source string, refID *int64, remark string) error {
	if amount <= 0 {
		return errors.ErrInvalidParams.WithMessage("amount must be positive")
	}

	return mysql.Transaction(func(tx *gorm.DB) error {
		// 获取余额
		balance, err := s.repo.GetBalance(ctx, userID)
		if err != nil {
			return err
		}

		// 检查余额是否足够
		if balance.Balance < amount {
			return errors.ErrInvalidParams.WithMessage("insufficient balance")
		}

		// 扣减余额
		balance.Balance -= amount
		if balanceType == BalanceTypeConsume {
			balance.TotalConsume += amount
		}

		if err := s.repo.UpdateBalance(ctx, balance); err != nil {
			return err
		}

		// 记录流水（负数）
		transactionNo := utils.GenerateOrderNo("BT")
		transaction := &BalanceTransaction{
			TransactionNo: transactionNo,
			UserID:        userID,
			Amount:        -amount, // 负数表示扣减
			Type:          balanceType,
			Source:        source,
			RefID:         refID,
			Remark:        remark,
		}

		return s.repo.CreateTransaction(ctx, transaction)
	})
}
