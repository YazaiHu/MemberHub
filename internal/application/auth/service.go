package auth

import (
	"context"
	"fmt"

	"github.com/YazaiHu/MemberHub/internal/domain/member"
	"github.com/YazaiHu/MemberHub/internal/infrastructure/wechat"
	"github.com/YazaiHu/MemberHub/internal/pkg/errors"
	"github.com/YazaiHu/MemberHub/internal/pkg/utils"
)

// Service 认证服务
type Service struct {
	memberRepo member.Repository
}

// NewService 创建认证服务
func NewService(memberRepo member.Repository) *Service {
	return &Service{
		memberRepo: memberRepo,
	}
}

// WeChatLoginRequest 微信登录请求
type WeChatLoginRequest struct {
	Code     string `json:"code" binding:"required"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int8   `json:"gender"`
}

// WeChatLoginResponse 微信登录响应
type WeChatLoginResponse struct {
	Token      string        `json:"token"`
	UserID     int64         `json:"user_id"`
	IsNewUser  bool          `json:"is_new_user"`
	UserInfo   *member.User  `json:"user_info"`
}

// WeChatLogin 微信小程序登录
func (s *Service) WeChatLogin(ctx context.Context, req *WeChatLoginRequest) (*WeChatLoginResponse, error) {
	// 1. 调用微信API获取openid
	openID, sessionKey, unionID, err := wechat.Code2Session(req.Code)
	if err != nil {
		return nil, errors.ErrWeChatAPIFailed.Wrap(err)
	}

	_ = sessionKey // sessionKey可用于解密用户信息，这里暂不使用

	// 2. 查询用户是否存在
	user, err := s.memberRepo.GetUserByOpenID(ctx, openID)
	isNewUser := false

	if err != nil {
		if err == errors.ErrUserNotFound {
			// 用户不存在，创建新用户
			user = &member.User{
				OpenID:   openID,
				UnionID:  unionID,
				Nickname: req.Nickname,
				Avatar:   req.Avatar,
				Gender:   req.Gender,
				Status:   1,
			}

			if err := s.memberRepo.CreateUser(ctx, user); err != nil {
				return nil, err
			}

			isNewUser = true
		} else {
			return nil, err
		}
	} else {
		// 用户已存在，更新最后登录时间
		if err := s.memberRepo.UpdateLastLoginTime(ctx, user.ID); err != nil {
			// 更新登录时间失败不影响登录流程
			fmt.Printf("update last login time failed: %v\n", err)
		}

		// 更新用户信息（如果有变化）
		updated := false
		if req.Nickname != "" && user.Nickname != req.Nickname {
			user.Nickname = req.Nickname
			updated = true
		}
		if req.Avatar != "" && user.Avatar != req.Avatar {
			user.Avatar = req.Avatar
			updated = true
		}
		if req.Gender > 0 && user.Gender != req.Gender {
			user.Gender = req.Gender
			updated = true
		}

		if updated {
			if err := s.memberRepo.UpdateUser(ctx, user); err != nil {
				fmt.Printf("update user info failed: %v\n", err)
			}
		}
	}

	// 3. 检查用户状态
	if user.Status == 0 {
		return nil, errors.ErrUserDisabled
	}

	// 4. 生成JWT token
	token, err := utils.GenerateToken(user.ID, "member", 0, "member")
	if err != nil {
		return nil, err
	}

	return &WeChatLoginResponse{
		Token:     token,
		UserID:    user.ID,
		IsNewUser: isNewUser,
		UserInfo:  user,
	}, nil
}

// AdminLoginRequest 管理员登录请求
type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AdminLoginResponse 管理员登录响应
type AdminLoginResponse struct {
	Token     string        `json:"token"`
	AdminID   int64         `json:"admin_id"`
	AdminInfo *member.Admin `json:"admin_info"`
}

// AdminLogin 管理员登录
func (s *Service) AdminLogin(ctx context.Context, req *AdminLoginRequest, ip string) (*AdminLoginResponse, error) {
	// 1. 查询管理员
	admin, err := s.memberRepo.GetAdminByUsername(ctx, req.Username)
	if err != nil {
		if err == errors.ErrUserNotFound {
			return nil, errors.ErrInvalidCredentials
		}
		return nil, err
	}

	// 2. 检查状态
	if admin.Status == 0 {
		return nil, errors.ErrUserDisabled.WithMessage("admin account is disabled")
	}

	// 3. 验证密码
	if !utils.CheckPassword(req.Password, admin.Password) {
		return nil, errors.ErrInvalidCredentials
	}

	// 4. 更新最后登录信息
	if err := s.memberRepo.UpdateAdminLastLogin(ctx, admin.ID, ip); err != nil {
		fmt.Printf("update admin last login failed: %v\n", err)
	}

	// 5. 获取角色代码
	roleCode := "member"
	if admin.Role != nil {
		roleCode = admin.Role.Code
	}

	// 6. 生成JWT token
	storeID := int64(0)
	if admin.StoreID != nil {
		storeID = *admin.StoreID
	}

	token, err := utils.GenerateToken(admin.ID, "admin", storeID, roleCode)
	if err != nil {
		return nil, err
	}

	// 清空密码字段
	admin.Password = ""

	return &AdminLoginResponse{
		Token:     token,
		AdminID:   admin.ID,
		AdminInfo: admin,
	}, nil
}

// RefreshToken 刷新token
func (s *Service) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	return utils.RefreshToken(oldToken)
}
