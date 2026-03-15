package recharge

import (
	"time"
)

// MemberBalance 会员余额汇总表
type MemberBalance struct {
	ID            int64     `gorm:"column:id;primaryKey" json:"id"`
	UserID        int64     `gorm:"column:user_id;uniqueIndex;not null" json:"user_id"`
	Balance       int64     `gorm:"column:balance;default:0" json:"balance"`             // 余额（分）
	FrozenBalance int64     `gorm:"column:frozen_balance;default:0" json:"frozen_balance"` // 冻结金额（分）
	TotalRecharge int64     `gorm:"column:total_recharge;default:0" json:"total_recharge"` // 累计充值（分）
	TotalConsume  int64     `gorm:"column:total_consume;default:0" json:"total_consume"`   // 累计消费（分）
	Version       int       `gorm:"column:version;default:0" json:"version"`             // 乐观锁
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt     time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (MemberBalance) TableName() string {
	return "member_balance"
}

// BalanceTransaction 余额流水表
type BalanceTransaction struct {
	ID            int64     `gorm:"column:id;primaryKey" json:"id"`
	TransactionNo string    `gorm:"column:transaction_no;uniqueIndex;size:64;not null" json:"transaction_no"`
	UserID        int64     `gorm:"column:user_id;index;not null" json:"user_id"`
	Amount        int64     `gorm:"column:amount;not null" json:"amount"` // 金额（分，正数增加，负数减少）
	Type          int8      `gorm:"column:type;not null;index" json:"type"` // 1-充值, 2-消费, 3-退款, 4-手动调整
	Source        string    `gorm:"column:source;size:50;not null" json:"source"`
	RefID         *int64    `gorm:"column:ref_id;index" json:"ref_id,omitempty"`
	Remark        string    `gorm:"column:remark;size:255" json:"remark,omitempty"`
	OperatorID    *int64    `gorm:"column:operator_id" json:"operator_id,omitempty"`
	CreatedAt     time.Time `gorm:"column:created_at;index" json:"created_at"`
}

// TableName 指定表名
func (BalanceTransaction) TableName() string {
	return "balance_transactions"
}

// 余额变动类型常量
const (
	BalanceTypeRecharge = 1 // 充值
	BalanceTypeConsume  = 2 // 消费
	BalanceTypeRefund   = 3 // 退款
	BalanceTypeAdjust   = 4 // 手动调整
)

// RechargePromotion 充值活动表
type RechargePromotion struct {
	ID             int64      `gorm:"column:id;primaryKey" json:"id"`
	Title          string     `gorm:"column:title;size:100;not null" json:"title"`
	Description    string     `gorm:"column:description;size:255" json:"description,omitempty"`
	RechargeAmount int        `gorm:"column:recharge_amount;not null" json:"recharge_amount"` // 充值金额（元）
	BonusAmount    int        `gorm:"column:bonus_amount;not null" json:"bonus_amount"`       // 赠送金额（元）
	BonusPoints    int        `gorm:"column:bonus_points;default:0" json:"bonus_points"`      // 赠送积分
	StartTime      *time.Time `gorm:"column:start_time" json:"start_time,omitempty"`
	EndTime        *time.Time `gorm:"column:end_time" json:"end_time,omitempty"`
	UserLimit      *int       `gorm:"column:user_limit" json:"user_limit,omitempty"`  // 每人限购次数
	TotalLimit     *int       `gorm:"column:total_limit" json:"total_limit,omitempty"` // 总限购次数
	Status         int8       `gorm:"column:status;default:1;index" json:"status"`    // 0-停用, 1-启用
	SortOrder      int        `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt      time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (RechargePromotion) TableName() string {
	return "recharge_promotions"
}

// RechargeOrder 充值订单表
type RechargeOrder struct {
	ID             int64      `gorm:"column:id;primaryKey" json:"id"`
	OrderNo        string     `gorm:"column:order_no;uniqueIndex;size:64;not null" json:"order_no"`
	UserID         int64      `gorm:"column:user_id;index;not null" json:"user_id"`
	PromotionID    *int64     `gorm:"column:promotion_id;index" json:"promotion_id,omitempty"`
	RechargeAmount int        `gorm:"column:recharge_amount;not null" json:"recharge_amount"` // 充值金额（分）
	BonusAmount    int        `gorm:"column:bonus_amount;default:0" json:"bonus_amount"`      // 赠送金额（分）
	BonusPoints    int        `gorm:"column:bonus_points;default:0" json:"bonus_points"`      // 赠送积分
	TotalAmount    int        `gorm:"column:total_amount;not null" json:"total_amount"`       // 实际到账金额（分）
	PayAmount      int        `gorm:"column:pay_amount;not null" json:"pay_amount"`           // 支付金额（分）
	PayMethod      int8       `gorm:"column:pay_method;default:1" json:"pay_method"`          // 1-微信支付
	TransactionID  string     `gorm:"column:transaction_id;size:64;index" json:"transaction_id,omitempty"` // 微信支付交易号
	Status         int8       `gorm:"column:status;default:1;index" json:"status"`            // 1-待支付, 2-已支付, 3-已取消, 4-已退款
	PaidAt         *time.Time `gorm:"column:paid_at" json:"paid_at,omitempty"`
	ExpiredAt      *time.Time `gorm:"column:expired_at" json:"expired_at,omitempty"`
	CreatedAt      time.Time  `gorm:"column:created_at;index" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (RechargeOrder) TableName() string {
	return "recharge_orders"
}

// 订单状态常量
const (
	OrderStatusPending   = 1 // 待支付
	OrderStatusPaid      = 2 // 已支付
	OrderStatusCancelled = 3 // 已取消
	OrderStatusRefunded  = 4 // 已退款
)

// 支付方式常量
const (
	PayMethodWeChatPay = 1 // 微信支付
)

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	PromotionID *int64 `json:"promotion_id"` // 活动ID，可选
	Amount      int    `json:"amount"`       // 充值金额（元），如果不参加活动则必填
}

// WeChatPayRequest 微信支付请求
type WeChatPayRequest struct {
	AppID      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	Body       string `xml:"body"`
	OutTradeNo string `xml:"out_trade_no"`
	TotalFee   int    `xml:"total_fee"`
	SpbillCreateIP string `xml:"spbill_create_ip"`
	NotifyURL  string `xml:"notify_url"`
	TradeType  string `xml:"trade_type"`
	OpenID     string `xml:"openid"`
}

// WeChatPayResponse 微信支付响应
type WeChatPayResponse struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppID      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	PrepayID   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
}

// WeChatPayNotify 微信支付回调通知
type WeChatPayNotify struct {
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
