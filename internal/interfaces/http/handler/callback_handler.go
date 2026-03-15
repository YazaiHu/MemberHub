package handler

import (
	"encoding/xml"
	"io"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	rechargeApp "github.com/YazaiHu/MemberHub/internal/application/recharge"
	"github.com/YazaiHu/MemberHub/internal/domain/points"
	"github.com/YazaiHu/MemberHub/internal/domain/recharge"
	"github.com/YazaiHu/MemberHub/internal/pkg/logger"
)

// CallbackHandler 回调处理器
type CallbackHandler struct {
	rechargeService *rechargeApp.Service
}

// NewCallbackHandler 创建回调处理器
func NewCallbackHandler() *CallbackHandler {
	rechargeRepo := recharge.NewRepository()
	pointsRepo := points.NewRepository()
	return &CallbackHandler{
		rechargeService: rechargeApp.NewService(rechargeRepo, pointsRepo),
	}
}

// WeChatPayCallback 微信支付回调
// @Summary 微信支付回调
// @Tags 回调
// @Accept xml
// @Produce xml
// @Param body body string true "微信回调XML"
// @Success 200 {string} string "XML"
// @Router /api/callback/wechat-pay [post]
func (h *CallbackHandler) WeChatPayCallback(c *gin.Context) {
	// 读取请求体
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logger.Error("read callback body failed", zap.Error(err))
		h.respondFail(c, "读取请求失败")
		return
	}

	logger.Info("received wechat pay callback", zap.String("body", string(body)))

	// 解析XML
	var notify recharge.WeChatPayNotify
	if err := xml.Unmarshal(body, &notify); err != nil {
		logger.Error("unmarshal callback xml failed", zap.Error(err))
		h.respondFail(c, "解析XML失败")
		return
	}

	// TODO: 验证签名
	// if !verifySignature(&notify) {
	//     logger.Error("invalid signature")
	//     h.respondFail(c, "签名验证失败")
	//     return
	// }

	// 检查返回码
	if notify.ReturnCode != "SUCCESS" {
		logger.Error("callback return code not success",
			zap.String("return_code", notify.ReturnCode),
			zap.String("return_msg", notify.ReturnMsg))
		h.respondFail(c, notify.ReturnMsg)
		return
	}

	// 检查业务结果
	if notify.ResultCode != "SUCCESS" {
		logger.Error("callback result code not success",
			zap.String("result_code", notify.ResultCode))
		h.respondFail(c, "支付失败")
		return
	}

	// 处理支付成功回调
	err = h.rechargeService.ProcessPaymentCallback(
		c.Request.Context(),
		notify.OutTradeNo,
		notify.TransactionID,
		notify.TotalFee,
	)

	if err != nil {
		logger.Error("process payment callback failed",
			zap.Error(err),
			zap.String("order_no", notify.OutTradeNo),
			zap.String("transaction_id", notify.TransactionID))
		h.respondFail(c, "处理回调失败")
		return
	}

	logger.Info("payment callback processed successfully",
		zap.String("order_no", notify.OutTradeNo),
		zap.String("transaction_id", notify.TransactionID),
		zap.Int("total_fee", notify.TotalFee))

	// 返回成功
	h.respondSuccess(c)
}

// respondSuccess 返回成功响应
func (h *CallbackHandler) respondSuccess(c *gin.Context) {
	response := `<xml>
  <return_code><![CDATA[SUCCESS]]></return_code>
  <return_msg><![CDATA[OK]]></return_msg>
</xml>`
	c.Data(200, "application/xml", []byte(response))
}

// respondFail 返回失败响应
func (h *CallbackHandler) respondFail(c *gin.Context, msg string) {
	response := `<xml>
  <return_code><![CDATA[FAIL]]></return_code>
  <return_msg><![CDATA[` + msg + `]]></return_msg>
</xml>`
	c.Data(200, "application/xml", []byte(response))
}
