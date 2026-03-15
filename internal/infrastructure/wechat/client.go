package wechat

import (
	"context"
	"fmt"
	"time"

	"github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
	"github.com/silenceper/wechat/v2/miniprogram"
	"go.uber.org/zap"

	"github.com/YazaiHu/MemberHub/internal/pkg/config"
	"github.com/YazaiHu/MemberHub/internal/pkg/logger"
)

var miniProgram *miniprogram.MiniProgram

// InitWeChat 初始化微信SDK
func InitWeChat(cfg *config.WeChatConfig) error {
	wc := wechat.NewWechat()

	// 使用内存缓存
	memCache := cache.NewMemory()

	miniCfg := &miniConfig.Config{
		AppID:     cfg.AppID,
		AppSecret: cfg.AppSecret,
		Cache:     memCache,
	}

	miniProgram = wc.GetMiniProgram(miniCfg)

	return nil
}

// GetMiniProgram 获取小程序实例
func GetMiniProgram() *miniprogram.MiniProgram {
	return miniProgram
}

// Code2Session 获取用户session信息
func Code2Session(code string) (openID, sessionKey, unionID string, err error) {
	auth := miniProgram.GetAuth()
	result, err := auth.Code2Session(code)
	if err != nil {
		return "", "", "", fmt.Errorf("code2session failed: %w", err)
	}

	return result.OpenID, result.SessionKey, result.UnionID, nil
}

// GetAccessToken 获取access_token
func GetAccessToken() (string, error) {
	ctx := context.Background()
	token, err := miniProgram.GetContext().GetAccessToken()
	if err != nil {
		return "", fmt.Errorf("get access token failed: %w", err)
	}

	// 从context中获取access token（这是一个简化示例）
	_ = ctx
	return token, nil
}

// PaymentNotifyData 支付回调数据
type PaymentNotifyData struct {
	ReturnCode    string `xml:"return_code"`
	ReturnMsg     string `xml:"return_msg"`
	AppID         string `xml:"appid"`
	MchID         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"`
	Sign          string `xml:"sign"`
	ResultCode    string `xml:"result_code"`
	OpenID        string `xml:"openid"`
	TradeType     string `xml:"trade_type"`
	BankType      string `xml:"bank_type"`
	TotalFee      int    `xml:"total_fee"`
	CashFee       int    `xml:"cash_fee"`
	TransactionID string `xml:"transaction_id"`
	OutTradeNo    string `xml:"out_trade_no"`
	TimeEnd       string `xml:"time_end"`
}

// TemplateMessage 模板消息
type TemplateMessage struct {
	ToUser      string                 `json:"touser"`
	TemplateID  string                 `json:"template_id"`
	Page        string                 `json:"page,omitempty"`
	Data        map[string]interface{} `json:"data"`
	MiniprogramState string            `json:"miniprogram_state,omitempty"` // developer, trial, formal
}

// TemplateDataItem 模板消息数据项
type TemplateDataItem struct {
	Value string `json:"value"`
}

// SendTemplateMessage 发送模板消息
func SendTemplateMessage(msg *TemplateMessage) error {
	// TODO: 实际生产环境需要调用微信API
	// 这里是框架实现，Worker测试时不需要真实调用微信

	// 日志记录（实际应该调用微信API）
	logger.Info("Sending template message",
		zap.String("to_user", msg.ToUser),
		zap.String("template_id", msg.TemplateID),
		zap.String("page", msg.Page),
	)

	// 模拟成功（实际环境需要调用微信API并处理错误）
	return nil
}

// RateLimitedSender 限流发送器
type RateLimitedSender struct {
	ratePerSecond int
	ticker        *time.Ticker
}

// NewRateLimitedSender 创建限流发送器
func NewRateLimitedSender(ratePerSecond int) *RateLimitedSender {
	return &RateLimitedSender{
		ratePerSecond: ratePerSecond,
		ticker:        time.NewTicker(time.Second / time.Duration(ratePerSecond)),
	}
}

// Send 限流发送
func (s *RateLimitedSender) Send(msg *TemplateMessage) error {
	<-s.ticker.C // 等待限流
	return SendTemplateMessage(msg)
}

// Stop 停止限流器
func (s *RateLimitedSender) Stop() {
	s.ticker.Stop()
}
