package member

import (
	"context"
	"fmt"

	"github.com/YazaiHu/MemberHub/internal/domain/member"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/persistence/mysql"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
)

// Service 会员服务
type Service struct {
	memberRepo member.Repository
}

// NewService 创建会员服务
func NewService(memberRepo member.Repository) *Service {
	return &Service{
		memberRepo: memberRepo,
	}
}

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(ctx context.Context, userID int64) (*member.User, error) {
	return s.memberRepo.GetUserByID(ctx, userID)
}

// UpdateUserInfo 更新用户信息
func (s *Service) UpdateUserInfo(ctx context.Context, user *member.User) error {
	// 检查用户是否存在
	existUser, err := s.memberRepo.GetUserByID(ctx, user.ID)
	if err != nil {
		return err
	}

	// 更新允许的字段
	existUser.Nickname = user.Nickname
	existUser.Avatar = user.Avatar
	existUser.Gender = user.Gender
	existUser.Birthday = user.Birthday

	return s.memberRepo.UpdateUser(ctx, existUser)
}

// BindPhone 绑定手机号
func (s *Service) BindPhone(ctx context.Context, userID int64, phone string) error {
	// 检查手机号是否已被其他用户绑定
	existUser, err := s.memberRepo.GetUserByPhone(ctx, phone)
	if err == nil && existUser.ID != userID {
		return errors.ErrConflict.WithMessage("phone number already bound to another user")
	}

	// 获取用户
	user, err := s.memberRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	// 更新手机号
	user.Phone = phone
	return s.memberRepo.UpdateUser(ctx, user)
}

// GetUserAsset 获取用户资产概览
func (s *Service) GetUserAsset(ctx context.Context, userID int64) (*member.UserAsset, error) {
	db := mysql.GetDB()
	asset := &member.UserAsset{
		UserID: userID,
	}

	// 查询积分
	var pointsResult struct {
		AvailablePoints int64
	}
	if err := db.WithContext(ctx).Table("member_points").
		Select("IFNULL(available_points, 0) as available_points").
		Where("user_id = ?", userID).
		Scan(&pointsResult).Error; err != nil {
		return nil, fmt.Errorf("query points failed: %w", err)
	}
	asset.AvailablePoints = pointsResult.AvailablePoints

	// 查询余额
	var balanceResult struct {
		Balance       int64
		TotalRecharge int64
		TotalConsume  int64
	}
	if err := db.WithContext(ctx).Table("member_balance").
		Select("IFNULL(balance, 0) as balance, IFNULL(total_recharge, 0) as total_recharge, IFNULL(total_consume, 0) as total_consume").
		Where("user_id = ?", userID).
		Scan(&balanceResult).Error; err != nil {
		return nil, fmt.Errorf("query balance failed: %w", err)
	}
	asset.Balance = balanceResult.Balance
	asset.TotalRecharge = balanceResult.TotalRecharge
	asset.TotalConsume = balanceResult.TotalConsume

	// 查询可用优惠券数量
	var couponCount int64
	if err := db.WithContext(ctx).Table("user_coupons").
		Where("user_id = ? AND status = 1 AND expired_at > NOW()", userID).
		Count(&couponCount).Error; err != nil {
		return nil, fmt.Errorf("query coupon count failed: %w", err)
	}
	asset.CouponCount = int(couponCount)

	return asset, nil
}

// ListUsers 用户列表（管理员使用）
func (s *Service) ListUsers(ctx context.Context, page, pageSize int, status *int8) ([]*member.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	return s.memberRepo.ListUsers(ctx, offset, pageSize, status)
}

// DisableUser 禁用用户
func (s *Service) DisableUser(ctx context.Context, userID int64) error {
	user, err := s.memberRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Status = 0
	return s.memberRepo.UpdateUser(ctx, user)
}

// EnableUser 启用用户
func (s *Service) EnableUser(ctx context.Context, userID int64) error {
	user, err := s.memberRepo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.Status = 1
	return s.memberRepo.UpdateUser(ctx, user)
}
