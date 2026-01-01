package handler

import (
	"log"
	"net/http"

	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/gin-gonic/gin"
)

// PaymentHandler 支付回调处理器
type PaymentHandler struct {
	paymentService *service.PaymentService
	settingService *service.SettingService
}

// NewPaymentHandler 创建支付回调处理器
func NewPaymentHandler(paymentService *service.PaymentService, settingService *service.SettingService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		settingService: settingService,
	}
}

// HandleEpayNotify 处理易支付异步回调
// POST /api/v1/payment/notify/epay
func (h *PaymentHandler) HandleEpayNotify(c *gin.Context) {
	// 1. 解析回调参数（表单或查询参数）
	params := make(map[string]string)

	// 优先从表单获取
	if c.Request.Method == http.MethodPost {
		if err := c.Request.ParseForm(); err == nil {
			for key, values := range c.Request.PostForm {
				if len(values) > 0 {
					params[key] = values[0]
				}
			}
		}
	}

	// 如果表单为空，从查询参数获取
	if len(params) == 0 {
		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				params[key] = values[0]
			}
		}
	}

	log.Printf("[Payment] Received epay notify: order_no=%s, trade_no=%s, trade_status=%s",
		params["out_trade_no"], params["trade_no"], params["trade_status"])

	// 2. 处理回调
	if err := h.paymentService.HandlePaymentCallback(c.Request.Context(), params); err != nil {
		log.Printf("[Payment] Epay notify failed: %v", err)

		// 根据错误类型返回不同状态码
		if err == service.ErrInvalidSign {
			c.String(http.StatusBadRequest, "fail")
			return
		}
		if err == service.ErrOrderLocked {
			// 订单正在处理中，告诉易支付稍后重试
			c.String(http.StatusTooManyRequests, "processing")
			return
		}

		// 其他错误返回失败
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	// 3. 返回成功（易支付要求返回"success"字符串）
	c.String(http.StatusOK, "success")
}

// HandleEpayReturn 处理易支付同步回调（可选）
// GET /api/v1/payment/return/epay
func (h *PaymentHandler) HandleEpayReturn(c *gin.Context) {
	// 提取回调参数
	params := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) > 0 {
			params[key] = values[0]
		}
	}

	orderNo := params["out_trade_no"]
	if orderNo == "" {
		// 没有订单号，重定向到充值页面
		c.Redirect(http.StatusFound, "/recharge")
		return
	}

	// 验证签名（可选但推荐）
	// 获取易支付配置
	settings, err := h.settingService.GetAllSettings(c.Request.Context())
	if err == nil && settings.EpayEnabled && settings.EpayMerchantKey != "" {
		// 验证签名
		epayClient := service.NewEpayClient(settings.EpayApiURL, settings.EpayMerchantID, settings.EpayMerchantKey)
		if !epayClient.VerifySign(params) {
			log.Printf("[Payment] Return callback invalid signature for order %s", orderNo)
			// 签名无效，仍然重定向但记录日志
		}
	}

	// 同步回调只用于页面跳转，不做业务处理
	// 重定向到前端支付结果页面
	redirectURL := "/recharge/result?order_no=" + orderNo
	c.Redirect(http.StatusFound, redirectURL)
}

// GetPaymentChannels 获取可用支付渠道
// GET /api/v1/payment/channels
func (h *PaymentHandler) GetPaymentChannels(c *gin.Context) {
	channels, err := h.settingService.GetPaymentChannels(c.Request.Context())
	if err != nil {
		log.Printf("[Payment] Failed to get payment channels: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get payment channels"})
		return
	}

	// 转换为前端需要的格式
	type ChannelResponse struct {
		Key         string `json:"key"`
		DisplayName string `json:"display_name"`
		Icon        string `json:"icon"`
	}

	result := make([]ChannelResponse, 0, len(channels))
	for _, ch := range channels {
		result = append(result, ChannelResponse{
			Key:         ch.Key,
			DisplayName: ch.DisplayName,
			Icon:        ch.Icon,
		})
	}

	response.Success(c, result)
}
