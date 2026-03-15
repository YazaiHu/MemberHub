package points

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/redis"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// Repository 积分仓储接口
type Repository interface {
	// MemberPoints
	GetOrCreateMemberPoints(ctx context.Context, userID int64) (*MemberPoints, error)
	GetMemberPoints(ctx context.Context, userID int64) (*MemberPoints, error)
	UpdateMemberPoints(ctx context.Context, points *MemberPoints) error

	// PointsTransaction
	CreateTransaction(ctx context.Context, tx *PointsTransaction) error
	GetTransactionByNo(ctx context.Context, transactionNo string) (*PointsTransaction, error)
	ListTransactions(ctx context.Context, userID int64, offset, limit int) ([]*PointsTransaction, int64, error)

	// ExchangeRule
	CreateRule(ctx context.Context, rule *ExchangeRule) error
	GetRuleByID(ctx context.Context, ruleID int64) (*ExchangeRule, error)
	UpdateRule(ctx context.Context, rule *ExchangeRule) error
	ListRules(ctx context.Context, storeID *int64, status *int8, offset, limit int) ([]*ExchangeRule, int64, error)
	DecrRuleStock(ctx context.Context, ruleID int64) error // 原子减库存

	// ExchangeRecord
	CreateRecord(ctx context.Context, record *ExchangeRecord) error
	GetRecordByID(ctx context.Context, recordID int64) (*ExchangeRecord, error)
	GetRecordByNo(ctx context.Context, recordNo string) (*ExchangeRecord, error)
	GetRecordByVerifyCode(ctx context.Context, verifyCode string) (*ExchangeRecord, error)
	UpdateRecord(ctx context.Context, record *ExchangeRecord) error
	ListRecords(ctx context.Context, userID *int64, storeID *int64, status *int8, offset, limit int) ([]*ExchangeRecord, int64, error)
	CountUserExchangeToday(ctx context.Context, userID, ruleID int64) (int, error)
	CountUserExchangeTotal(ctx context.Context, userID, ruleID int64) (int, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository 创建积分仓储
func NewRepository() Repository {
	return &repository{
		db: mysql.GetDB(),
	}
}

// GetOrCreateMemberPoints 获取或创建会员积分
func (r *repository) GetOrCreateMemberPoints(ctx context.Context, userID int64) (*MemberPoints, error) {
	var points MemberPoints
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&points).Error

	if err == gorm.ErrRecordNotFound {
		// 不存在则创建
		points = MemberPoints{
			UserID:          userID,
			TotalPoints:     0,
			AvailablePoints: 0,
			FrozenPoints:    0,
			UsedPoints:      0,
			ExpiredPoints:   0,
			Version:         0,
		}
		if err := r.db.WithContext(ctx).Create(&points).Error; err != nil {
			return nil, fmt.Errorf("create member points failed: %w", err)
		}
		return &points, nil
	}

	if err != nil {
		return nil, fmt.Errorf("get member points failed: %w", err)
	}

	return &points, nil
}

// GetMemberPoints 获取会员积分
func (r *repository) GetMemberPoints(ctx context.Context, userID int64) (*MemberPoints, error) {
	var points MemberPoints
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&points).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get member points failed: %w", err)
	}
	return &points, nil
}

// UpdateMemberPoints 更新会员积分（带乐观锁）
func (r *repository) UpdateMemberPoints(ctx context.Context, points *MemberPoints) error {
	result := r.db.WithContext(ctx).Model(&MemberPoints{}).
		Where("id = ? AND version = ?", points.ID, points.Version).
		Updates(map[string]interface{}{
			"total_points":     points.TotalPoints,
			"available_points": points.AvailablePoints,
			"frozen_points":    points.FrozenPoints,
			"used_points":      points.UsedPoints,
			"expired_points":   points.ExpiredPoints,
			"version":          gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return fmt.Errorf("update member points failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrConcurrentUpdate
	}

	points.Version++
	return nil
}

// CreateTransaction 创建积分流水
func (r *repository) CreateTransaction(ctx context.Context, tx *PointsTransaction) error {
	if err := r.db.WithContext(ctx).Create(tx).Error; err != nil {
		return fmt.Errorf("create points transaction failed: %w", err)
	}
	return nil
}

// GetTransactionByNo 根据流水号获取流水
func (r *repository) GetTransactionByNo(ctx context.Context, transactionNo string) (*PointsTransaction, error) {
	var tx PointsTransaction
	if err := r.db.WithContext(ctx).Where("transaction_no = ?", transactionNo).First(&tx).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get transaction failed: %w", err)
	}
	return &tx, nil
}

// ListTransactions 积分流水列表
func (r *repository) ListTransactions(ctx context.Context, userID int64, offset, limit int) ([]*PointsTransaction, int64, error) {
	var transactions []*PointsTransaction
	var total int64

	query := r.db.WithContext(ctx).Model(&PointsTransaction{}).Where("user_id = ? AND status = 1", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count transactions failed: %w", err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&transactions).Error; err != nil {
		return nil, 0, fmt.Errorf("list transactions failed: %w", err)
	}

	return transactions, total, nil
}

// CreateRule 创建兑换规则
func (r *repository) CreateRule(ctx context.Context, rule *ExchangeRule) error {
	if err := r.db.WithContext(ctx).Create(rule).Error; err != nil {
		return fmt.Errorf("create exchange rule failed: %w", err)
	}
	return nil
}

// GetRuleByID 获取兑换规则
func (r *repository) GetRuleByID(ctx context.Context, ruleID int64) (*ExchangeRule, error) {
	var rule ExchangeRule
	if err := r.db.WithContext(ctx).Where("id = ?", ruleID).First(&rule).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrPointsRuleNotFound
		}
		return nil, fmt.Errorf("get exchange rule failed: %w", err)
	}
	return &rule, nil
}

// UpdateRule 更新兑换规则
func (r *repository) UpdateRule(ctx context.Context, rule *ExchangeRule) error {
	if err := r.db.WithContext(ctx).Save(rule).Error; err != nil {
		return fmt.Errorf("update exchange rule failed: %w", err)
	}
	return nil
}

// ListRules 兑换规则列表
func (r *repository) ListRules(ctx context.Context, storeID *int64, status *int8, offset, limit int) ([]*ExchangeRule, int64, error) {
	var rules []*ExchangeRule
	var total int64

	query := r.db.WithContext(ctx).Model(&ExchangeRule{})
	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count rules failed: %w", err)
	}

	if err := query.Order("sort_order DESC, id DESC").Offset(offset).Limit(limit).Find(&rules).Error; err != nil {
		return nil, 0, fmt.Errorf("list rules failed: %w", err)
	}

	return rules, total, nil
}

// DecrRuleStock 原子减库存
func (r *repository) DecrRuleStock(ctx context.Context, ruleID int64) error {
	result := r.db.WithContext(ctx).Model(&ExchangeRule{}).
		Where("id = ? AND remaining_stock > 0", ruleID).
		Updates(map[string]interface{}{
			"remaining_stock": gorm.Expr("remaining_stock - 1"),
			"version":         gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return fmt.Errorf("decr rule stock failed: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.ErrStockNotEnough
	}

	return nil
}

// CreateRecord 创建兑换记录
func (r *repository) CreateRecord(ctx context.Context, record *ExchangeRecord) error {
	if err := r.db.WithContext(ctx).Create(record).Error; err != nil {
		return fmt.Errorf("create exchange record failed: %w", err)
	}
	return nil
}

// GetRecordByID 获取兑换记录
func (r *repository) GetRecordByID(ctx context.Context, recordID int64) (*ExchangeRecord, error) {
	var record ExchangeRecord
	if err := r.db.WithContext(ctx).Where("id = ?", recordID).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get exchange record failed: %w", err)
	}
	return &record, nil
}

// GetRecordByNo 根据兑换单号获取记录
func (r *repository) GetRecordByNo(ctx context.Context, recordNo string) (*ExchangeRecord, error) {
	var record ExchangeRecord
	if err := r.db.WithContext(ctx).Where("record_no = ?", recordNo).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get exchange record failed: %w", err)
	}
	return &record, nil
}

// GetRecordByVerifyCode 根据核销码获取记录
func (r *repository) GetRecordByVerifyCode(ctx context.Context, verifyCode string) (*ExchangeRecord, error) {
	var record ExchangeRecord
	if err := r.db.WithContext(ctx).Where("verify_code = ?", verifyCode).First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("get exchange record failed: %w", err)
	}
	return &record, nil
}

// UpdateRecord 更新兑换记录
func (r *repository) UpdateRecord(ctx context.Context, record *ExchangeRecord) error {
	if err := r.db.WithContext(ctx).Save(record).Error; err != nil {
		return fmt.Errorf("update exchange record failed: %w", err)
	}
	return nil
}

// ListRecords 兑换记录列表
func (r *repository) ListRecords(ctx context.Context, userID *int64, storeID *int64, status *int8, offset, limit int) ([]*ExchangeRecord, int64, error) {
	var records []*ExchangeRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&ExchangeRecord{})
	if userID != nil {
		query = query.Where("user_id = ?", *userID)
	}
	if storeID != nil {
		query = query.Where("store_id = ?", *storeID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count records failed: %w", err)
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&records).Error; err != nil {
		return nil, 0, fmt.Errorf("list records failed: %w", err)
	}

	return records, total, nil
}

// CountUserExchangeToday 统计用户今日兑换次数
func (r *repository) CountUserExchangeToday(ctx context.Context, userID, ruleID int64) (int, error) {
	var count int64
	today := time.Now().Format("2006-01-02")

	if err := r.db.WithContext(ctx).Model(&ExchangeRecord{}).
		Where("user_id = ? AND rule_id = ? AND DATE(created_at) = ? AND status != ?",
			userID, ruleID, today, ExchangeStatusCancelled).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count user exchange today failed: %w", err)
	}

	return int(count), nil
}

// CountUserExchangeTotal 统计用户总兑换次数
func (r *repository) CountUserExchangeTotal(ctx context.Context, userID, ruleID int64) (int, error) {
	var count int64

	if err := r.db.WithContext(ctx).Model(&ExchangeRecord{}).
		Where("user_id = ? AND rule_id = ? AND status != ?",
			userID, ruleID, ExchangeStatusCancelled).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count user exchange total failed: %w", err)
	}

	return int(count), nil
}

// --- Redis辅助方法 ---

// GetStockFromRedis 从Redis获取库存
func GetStockFromRedis(ruleID int64) (int, error) {
	key := fmt.Sprintf("exchange:stock:%d", ruleID)
	val, err := redis.Get(key)
	if err != nil {
		return 0, err
	}

	var stock int
	_, err = fmt.Sscanf(val, "%d", &stock)
	return stock, err
}

// SetStockToRedis 设置Redis库存
func SetStockToRedis(ruleID int64, stock int) error {
	key := fmt.Sprintf("exchange:stock:%d", ruleID)
	return redis.Set(key, fmt.Sprintf("%d", stock), 24*time.Hour)
}

// DecrStockInRedis 在Redis中原子减库存
func DecrStockInRedis(ruleID int64) (int64, error) {
	key := fmt.Sprintf("exchange:stock:%d", ruleID)
	return redis.Decr(key)
}

// IncrStockInRedis 在Redis中原子增库存（回滚用）
func IncrStockInRedis(ruleID int64) (int64, error) {
	key := fmt.Sprintf("exchange:stock:%d", ruleID)
	return redis.Incr(key)
}
