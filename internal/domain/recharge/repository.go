package recharge

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// Repository 充值仓储接口
type Repository interface {
	// MemberBalance
	GetOrCreateBalance(ctx context.Context, userID int64) (*MemberBalance, error)
	GetBalance(ctx context.Context, userID int64) (*MemberBalance, error)
	UpdateBalance(ctx context.Context, balance *MemberBalance) error

	// BalanceTransaction
	CreateTransaction(ctx context.Context, tx *BalanceTransaction) error
	GetTransactionByNo(ctx context.Context, transactionNo string) (*BalanceTransaction, error)
	ListTransactions(ctx context.Context, userID int64, offset, limit int) ([]*BalanceTransaction, int64, error)

	// RechargePromotion
	CreatePromotion(ctx context.Context, promotion *RechargePromotion) error
	GetPromotionByID(ctx context.Context, promotionID int64) (*RechargePromotion, error)
	UpdatePromotion(ctx context.Context, promotion *RechargePromotion) error
	ListPromotions(ctx context.Context, status *int8, offset, limit int) ([]*RechargePromotion, int64, error)
	GetActivePromotions(ctx context.Context) ([]*RechargePromotion, error)

	// RechargeOrder
	CreateOrder(ctx context.Context, order *RechargeOrder) error
	GetOrderByID(ctx context.Context, orderID int64) (*RechargeOrder, error)
	GetOrderByNo(ctx context.Context, orderNo string) (*RechargeOrder, error)
	UpdateOrder(ctx context.Context, order *RechargeOrder) error
	ListOrders(ctx context.Context, userID *int64, status *int8, offset, limit int) ([]*RechargeOrder, int64, error)
	CountUserRechargeByPromotion(ctx context.Context, userID, promotionID int64) (int, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository 创建充值仓储
func NewRepository() Repository {
	return &repository{
		db: mysql.GetDB(),
	}
}

// GetOrCreateBalance 获取或创建会员余额
func (r *repository) GetOrCreateBalance(ctx context.Context, userID int64) (*MemberBalance, error) {
	var balance MemberBalance
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&balance).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在则创建
		balance = MemberBalance{
			UserID:        userID,
			Balance:       0,
			FrozenBalance: 0,
			TotalRecharge: 0,
			TotalConsume:  0,
			Version:       0,
		}
		if err := r.db.WithContext(ctx).Create(&balance).Error; err != nil {
			return nil, fmt.Errorf("create member balance failed: %w", err)
		}
		return &balance, nil
	}

	if err != nil {
		return nil, fmt.Errorf("get member balance failed: %w", err)
	}

	return &balance, nil
}

// GetBalance 获取会员余额
func (r *repository) GetBalance(ctx context.Context, userID int64) (*MemberBalance, error) {
	var balance MemberBalance
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&balance).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get member balance failed: %w", err)
	}
	return &balance, nil
}

// UpdateBalance 更新会员余额（带乐观锁）
func (r *repository) UpdateBalance(ctx context.Context, balance *MemberBalance) error {
	result := r.db.WithContext(ctx).Model(&MemberBalance{}).
		Where("id = ? AND version = ?", balance.ID, balance.Version).
		Updates(map[string]interface{}{
			"balance":        balance.Balance,
			"frozen_balance": balance.FrozenBalance,
			"total_recharge": balance.TotalRecharge,
			"total_consume":  balance.TotalConsume,
			"version":        gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return fmt.Errorf("update member balance failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrConcurrentUpdate
	}

	balance.Version++
	return nil
}

// CreateTransaction 创建余额流水
func (r *repository) CreateTransaction(ctx context.Context, tx *BalanceTransaction) error {
	if err := r.db.WithContext(ctx).Create(tx).Error; err != nil {
		return fmt.Errorf("create balance transaction failed: %w", err)
	}
	return nil
}

// GetTransactionByNo 根据流水号获取流水
func (r *repository) GetTransactionByNo(ctx context.Context, transactionNo string) (*BalanceTransaction, error) {
	var tx BalanceTransaction
	if err := r.db.WithContext(ctx).Where("transaction_no = ?", transactionNo).First(&tx).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get balance transaction failed: %w", err)
	}
	return &tx, nil
}

// ListTransactions 余额流水列表
func (r *repository) ListTransactions(ctx context.Context, userID int64, offset, limit int) ([]*BalanceTransaction, int64, error) {
	var transactions []*BalanceTransaction
	var total int64

	query := r.db.WithContext(ctx).Model(&BalanceTransaction{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count transactions failed: %w", err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("list transactions failed: %w", err)
	}

	return transactions, total, nil
}

// CreatePromotion 创建充值活动
func (r *repository) CreatePromotion(ctx context.Context, promotion *RechargePromotion) error {
	if err := r.db.WithContext(ctx).Create(promotion).Error; err != nil {
		return fmt.Errorf("create recharge promotion failed: %w", err)
	}
	return nil
}

// GetPromotionByID 获取充值活动
func (r *repository) GetPromotionByID(ctx context.Context, promotionID int64) (*RechargePromotion, error) {
	var promotion RechargePromotion
	if err := r.db.WithContext(ctx).Where("id = ?", promotionID).First(&promotion).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrPromotionNotFound
		}
		return nil, fmt.Errorf("get recharge promotion failed: %w", err)
	}
	return &promotion, nil
}

// UpdatePromotion 更新充值活动
func (r *repository) UpdatePromotion(ctx context.Context, promotion *RechargePromotion) error {
	if err := r.db.WithContext(ctx).Save(promotion).Error; err != nil {
		return fmt.Errorf("update recharge promotion failed: %w", err)
	}
	return nil
}

// ListPromotions 充值活动列表
func (r *repository) ListPromotions(ctx context.Context, status *int8, offset, limit int) ([]*RechargePromotion, int64, error) {
	var promotions []*RechargePromotion
	var total int64

	query := r.db.WithContext(ctx).Model(&RechargePromotion{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count promotions failed: %w", err)
	}

	if err := query.Order("sort_order DESC, id DESC").Offset(offset).Limit(limit).Find(&promotions).Error; err != nil {
		return nil, 0, fmt.Errorf("list promotions failed: %w", err)
	}

	return promotions, total, nil
}

// GetActivePromotions 获取有效的充值活动
func (r *repository) GetActivePromotions(ctx context.Context) ([]*RechargePromotion, error) {
	var promotions []*RechargePromotion
	query := r.db.WithContext(ctx).Where("status = 1")

	// 时间范围检查
	query = query.Where("(start_time IS NULL OR start_time <= NOW()) AND (end_time IS NULL OR end_time >= NOW())")

	if err := query.Order("sort_order DESC, id DESC").Find(&promotions).Error; err != nil {
		return nil, fmt.Errorf("get active promotions failed: %w", err)
	}

	return promotions, nil
}

// CreateOrder 创建充值订单
func (r *repository) CreateOrder(ctx context.Context, order *RechargeOrder) error {
	if err := r.db.WithContext(ctx).Create(order).Error; err != nil {
		return fmt.Errorf("create recharge order failed: %w", err)
	}
	return nil
}

// GetOrderByID 获取充值订单
func (r *repository) GetOrderByID(ctx context.Context, orderID int64) (*RechargeOrder, error) {
	var order RechargeOrder
	if err := r.db.WithContext(ctx).Where("id = ?", orderID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrOrderNotFound
		}
		return nil, fmt.Errorf("get recharge order failed: %w", err)
	}
	return &order, nil
}

// GetOrderByNo 根据订单号获取订单
func (r *repository) GetOrderByNo(ctx context.Context, orderNo string) (*RechargeOrder, error) {
	var order RechargeOrder
	if err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrOrderNotFound
		}
		return nil, fmt.Errorf("get recharge order failed: %w", err)
	}
	return &order, nil
}

// UpdateOrder 更新充值订单
func (r *repository) UpdateOrder(ctx context.Context, order *RechargeOrder) error {
	if err := r.db.WithContext(ctx).Save(order).Error; err != nil {
		return fmt.Errorf("update recharge order failed: %w", err)
	}
	return nil
}

// ListOrders 充值订单列表
func (r *repository) ListOrders(ctx context.Context, userID *int64, status *int8, offset, limit int) ([]*RechargeOrder, int64, error) {
	var orders []*RechargeOrder
	var total int64

	query := r.db.WithContext(ctx).Model(&RechargeOrder{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count orders failed: %w", err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&orders).Error; err != nil {
		return nil, 0, fmt.Errorf("list orders failed: %w", err)
	}

	return orders, total, nil
}

// CountUserRechargeByPromotion 统计用户参加活动的充值次数
func (r *repository) CountUserRechargeByPromotion(ctx context.Context, userID, promotionID int64) (int, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&RechargeOrder{}).
		Where("user_id = ? AND promotion_id = ? AND status = ?", userID, promotionID, OrderStatusPaid).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count user recharge by promotion failed: %w", err)
	}
	return int(count), nil
}
