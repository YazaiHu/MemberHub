package points

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/redis"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
	"github.com/YazaiHu/MemberHub/internal/pkg/utils"
)

// Service 积分领域服务
type Service struct {
	repo Repository
}

// NewService 创建积分服务
func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// AddPoints 增加积分（消费赠送、充值赠送、手动调整）
func (s *Service) AddPoints(ctx context.Context, userID int64, points int64, pointsType int8, source string, refID *int64, remark string, expireDays int) error {
	if points <= 0 {
		return errors.ErrInvalidParams.WithMessage("points must be positive")
	}

	// 使用事务
	return mysql.Transaction(func(tx *gorm.DB) error {
		// 1. 获取或创建会员积分记录
		memberPoints, err := s.repo.GetOrCreateMemberPoints(ctx, userID)
		if err != nil {
			return err
		}

		// 2. 创建积分流水
		transactionNo := utils.GenerateOrderNo("PT")
		var expireAt *time.Time
		if expireDays > 0 {
			t := time.Now().AddDate(0, 0, expireDays)
			expireAt = &t
		}

		transaction := &PointsTransaction{
			TransactionNo: transactionNo,
			UserID:        userID,
			Points:        points,
			Type:          pointsType,
			Source:        source,
			RefID:         refID,
			Remark:        remark,
			ExpireAt:      expireAt,
			Status:        1,
		}

		if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
			return err
		}

		// 3. 更新汇总表（带乐观锁）
		memberPoints.TotalPoints += points
		memberPoints.AvailablePoints += points

		if err := s.repo.UpdateMemberPoints(ctx, memberPoints); err != nil {
			return err
		}

		return nil
	})
}

// DeductPoints 扣减积分（兑换扣减、过期扣减）
func (s *Service) DeductPoints(ctx context.Context, userID int64, points int64, pointsType int8, source string, refID *int64, remark string) error {
	if points <= 0 {
		return errors.ErrInvalidParams.WithMessage("points must be positive")
	}

	return mysql.Transaction(func(tx *gorm.DB) error {
		// 1. 获取会员积分
		memberPoints, err := s.repo.GetMemberPoints(ctx, userID)
		if err != nil {
			return err
		}

		// 2. 检查积分是否足够
		if memberPoints.AvailablePoints < points {
			return errors.ErrInsufficientPoints
		}

		// 3. 创建积分流水（负数表示扣减）
		transactionNo := utils.GenerateOrderNo("PT")
		transaction := &PointsTransaction{
			TransactionNo: transactionNo,
			UserID:        userID,
			Points:        -points, // 负数
			Type:          pointsType,
			Source:        source,
			RefID:         refID,
			Remark:        remark,
			Status:        1,
		}

		if err := s.repo.CreateTransaction(ctx, transaction); err != nil {
			return err
		}

		// 4. 更新汇总表（带乐观锁）
		memberPoints.AvailablePoints -= points

		if pointsType == PointsTypeExchange {
			memberPoints.UsedPoints += points
		} else if pointsType == PointsTypeExpire {
			memberPoints.ExpiredPoints += points
		}

		if err := s.repo.UpdateMemberPoints(ctx, memberPoints); err != nil {
			return err
		}

		return nil
	})
}

// ExchangePoints 积分兑换（高并发场景）
func (s *Service) ExchangePoints(ctx context.Context, userID, ruleID int64) (*ExchangeRecord, error) {
	// 1. 获取兑换规则
	rule, err := s.repo.GetRuleByID(ctx, ruleID)
	if err != nil {
		return nil, err
	}

	// 2. 检查规则状态
	if rule.Status != 1 {
		return nil, errors.ErrPointsRuleNotFound.WithMessage("exchange rule is unavailable")
	}

	// 3. 检查时间范围
	now := time.Now()
	if rule.StartTime != nil && now.Before(*rule.StartTime) {
		return nil, errors.ErrInvalidParams.WithMessage("exchange not started yet")
	}
	if rule.EndTime != nil && now.After(*rule.EndTime) {
		return nil, errors.ErrInvalidParams.WithMessage("exchange has ended")
	}

	// 4. 检查用户今日兑换次数
	todayCount, err := s.repo.CountUserExchangeToday(ctx, userID, ruleID)
	if err != nil {
		return nil, err
	}
	if todayCount >= rule.UserDailyLimit {
		return nil, errors.ErrExchangeLimitExceeded.WithMessage("daily exchange limit exceeded")
	}

	// 5. 检查用户总兑换次数
	if rule.UserTotalLimit != nil {
		totalCount, err := s.repo.CountUserExchangeTotal(ctx, userID, ruleID)
		if err != nil {
			return nil, err
		}
		if totalCount >= *rule.UserTotalLimit {
			return nil, errors.ErrExchangeLimitExceeded.WithMessage("total exchange limit exceeded")
		}
	}

	// 6. 分布式锁防止重复提交
	lockKey := fmt.Sprintf("lock:exchange:%d:%d", userID, ruleID)
	lock, acquired, err := redis.AcquireLock(lockKey, 5*time.Second)
	if err != nil {
		return nil, err
	}
	if !acquired {
		return nil, errors.ErrDuplicateOperation
	}
	defer lock.Release()

	// 7. Redis原子扣库存（快速失败）
	if rule.RemainingStock != nil {
		remaining, err := DecrStockInRedis(ruleID)
		if err != nil || remaining < 0 {
			// 回滚Redis库存
			IncrStockInRedis(ruleID)
			return nil, errors.ErrStockNotEnough
		}
	}

	// 8. 使用事务完成兑换
	var record *ExchangeRecord
	err = mysql.Transaction(func(tx *gorm.DB) error {
		// 8.1 扣减数据库库存
		if rule.RemainingStock != nil {
			if err := s.repo.DecrRuleStock(ctx, ruleID); err != nil {
				return err
			}
		}

		// 8.2 扣减用户积分
		if err := s.DeductPoints(ctx, userID, int64(rule.PointsCost), PointsTypeExchange, "exchange", &ruleID, fmt.Sprintf("兑换：%s", rule.Title)); err != nil {
			return err
		}

		// 8.3 创建兑换记录
		recordNo := utils.GenerateOrderNo("EX")
		verifyCode := utils.RandomString(8)

		record = &ExchangeRecord{
			RecordNo:    recordNo,
			UserID:      userID,
			RuleID:      ruleID,
			StoreID:     rule.StoreID,
			PointsCost:  rule.PointsCost,
			ProductName: rule.ProductName,
			Status:      ExchangeStatusPending,
			VerifyCode:  verifyCode,
		}

		if err := s.repo.CreateRecord(ctx, record); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		// 事务失败，回滚Redis库存
		if rule.RemainingStock != nil {
			IncrStockInRedis(ruleID)
		}
		return nil, err
	}

	return record, nil
}

// VerifyExchange 核销兑换
func (s *Service) VerifyExchange(ctx context.Context, verifyCode string, verifierID int64) error {
	// 1. 获取兑换记录
	record, err := s.repo.GetRecordByVerifyCode(ctx, verifyCode)
	if err != nil {
		return err
	}

	// 2. 检查状态
	if record.Status != ExchangeStatusPending {
		return errors.ErrInvalidParams.WithMessage("record already verified or cancelled")
	}

	// 3. 更新记录
	now := time.Now()
	record.Status = ExchangeStatusVerified
	record.VerifiedAt = &now
	record.VerifiedBy = &verifierID

	return s.repo.UpdateRecord(ctx, record)
}

// CancelExchange 取消兑换（退还积分）
func (s *Service) CancelExchange(ctx context.Context, recordID int64, reason string) error {
	// 1. 获取兑换记录
	record, err := s.repo.GetRecordByID(ctx, recordID)
	if err != nil {
		return err
	}

	// 2. 检查状态
	if record.Status != ExchangeStatusPending {
		return errors.ErrInvalidParams.WithMessage("only pending records can be cancelled")
	}

	// 3. 使用事务
	return mysql.Transaction(func(tx *gorm.DB) error {
		// 3.1 更新兑换记录
		now := time.Now()
		record.Status = ExchangeStatusCancelled
		record.CancelledAt = &now
		record.CancelReason = reason

		if err := s.repo.UpdateRecord(ctx, record); err != nil {
			return err
		}

		// 3.2 退还积分
		if err := s.AddPoints(ctx, record.UserID, int64(record.PointsCost), PointsTypeAdjust, "cancel_exchange", &recordID, fmt.Sprintf("取消兑换退还：%s", record.ProductName), 0); err != nil {
			return err
		}

		// 3.3 退还库存
		rule, err := s.repo.GetRuleByID(ctx, record.RuleID)
		if err == nil && rule.RemainingStock != nil {
			// 增加数据库库存
			*rule.RemainingStock++
			s.repo.UpdateRule(ctx, rule)

			// 增加Redis库存
			IncrStockInRedis(record.RuleID)
		}

		return nil
	})
}

// GetUserPoints 获取用户积分
func (s *Service) GetUserPoints(ctx context.Context, userID int64) (*MemberPoints, error) {
	return s.repo.GetOrCreateMemberPoints(ctx, userID)
}

// GetPointsTransactions 获取积分流水
func (s *Service) GetPointsTransactions(ctx context.Context, userID int64, page, pageSize int) ([]*PointsTransaction, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.repo.ListTransactions(ctx, userID, offset, pageSize)
}
