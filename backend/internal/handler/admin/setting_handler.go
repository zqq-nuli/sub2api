package admin

import (
	"fmt"

	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/oidc"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// SettingHandler 系统设置处理器
type SettingHandler struct {
	settingService   *service.SettingService
	emailService     *service.EmailService
	turnstileService *service.TurnstileService
}

// NewSettingHandler 创建系统设置处理器
func NewSettingHandler(settingService *service.SettingService, emailService *service.EmailService, turnstileService *service.TurnstileService) *SettingHandler {
	return &SettingHandler{
		settingService:   settingService,
		emailService:     emailService,
		turnstileService: turnstileService,
	}
}

// maskSensitiveString 对敏感字符串进行掩码处理
// 显示前4个字符和后4个字符，中间用 * 替代
func maskSensitiveString(s string) string {
	if s == "" {
		return ""
	}
	if len(s) <= 8 {
		return "********"
	}
	return s[:4] + "********" + s[len(s)-4:]
}

// GetSettings 获取所有系统设置
// GET /api/v1/admin/settings
func (h *SettingHandler) GetSettings(c *gin.Context) {
	settings, err := h.settingService.GetAllSettings(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// 敏感字段进行掩码处理，避免在 API 响应中泄露
	response.Success(c, dto.SystemSettings{
		RegistrationEnabled: settings.RegistrationEnabled,
		EmailVerifyEnabled:  settings.EmailVerifyEnabled,
		SmtpHost:            settings.SmtpHost,
		SmtpPort:            settings.SmtpPort,
		SmtpUsername:        settings.SmtpUsername,
		SmtpPassword:        maskSensitiveString(settings.SmtpPassword), // 掩码
		SmtpFrom:            settings.SmtpFrom,
		SmtpFromName:        settings.SmtpFromName,
		SmtpUseTLS:          settings.SmtpUseTLS,
		TurnstileEnabled:    settings.TurnstileEnabled,
		TurnstileSiteKey:    settings.TurnstileSiteKey,
		TurnstileSecretKey:  maskSensitiveString(settings.TurnstileSecretKey), // 掩码
		SiteName:            settings.SiteName,
		SiteLogo:            settings.SiteLogo,
		SiteSubtitle:        settings.SiteSubtitle,
		ApiBaseUrl:          settings.ApiBaseUrl,
		ContactInfo:         settings.ContactInfo,
		DocUrl:              settings.DocUrl,
		DefaultConcurrency:  settings.DefaultConcurrency,
		DefaultBalance:      settings.DefaultBalance,
		// SSO设置
		SSOEnabled:           settings.SSOEnabled,
		PasswordLoginEnabled: settings.PasswordLoginEnabled,
		SSOIssuerURL:         settings.SSOIssuerURL,
		SSOClientID:          settings.SSOClientID,
		SSOClientSecret:      maskSensitiveString(settings.SSOClientSecret), // 掩码
		SSORedirectURI:       settings.SSORedirectURI,
		SSOAllowedDomains:    settings.SSOAllowedDomains,
		SSOAutoCreateUser:    settings.SSOAutoCreateUser,
		// 易支付设置
		EpayEnabled:     settings.EpayEnabled,
		EpayApiURL:      settings.EpayApiURL,
		EpayMerchantID:  settings.EpayMerchantID,
		EpayMerchantKey: maskSensitiveString(settings.EpayMerchantKey), // 掩码
		EpayNotifyURL:   settings.EpayNotifyURL,
		EpayReturnURL:   settings.EpayReturnURL,
		// 支付渠道
		PaymentChannels: serviceChannelsToDTO(settings.PaymentChannels),
	})
}

// UpdateSettingsRequest 更新设置请求
type UpdateSettingsRequest struct {
	// 注册设置
	RegistrationEnabled bool `json:"registration_enabled"`
	EmailVerifyEnabled  bool `json:"email_verify_enabled"`

	// 邮件服务设置
	SmtpHost     string `json:"smtp_host"`
	SmtpPort     int    `json:"smtp_port"`
	SmtpUsername string `json:"smtp_username"`
	SmtpPassword string `json:"smtp_password"`
	SmtpFrom     string `json:"smtp_from_email"`
	SmtpFromName string `json:"smtp_from_name"`
	SmtpUseTLS   bool   `json:"smtp_use_tls"`

	// Cloudflare Turnstile 设置
	TurnstileEnabled   bool   `json:"turnstile_enabled"`
	TurnstileSiteKey   string `json:"turnstile_site_key"`
	TurnstileSecretKey string `json:"turnstile_secret_key"`

	// OEM设置
	SiteName     string `json:"site_name"`
	SiteLogo     string `json:"site_logo"`
	SiteSubtitle string `json:"site_subtitle"`
	ApiBaseUrl   string `json:"api_base_url"`
	ContactInfo  string `json:"contact_info"`
	DocUrl       string `json:"doc_url"`

	// 默认配置
	DefaultConcurrency int     `json:"default_concurrency"`
	DefaultBalance     float64 `json:"default_balance"`

	// SSO设置
	SSOEnabled           bool     `json:"sso_enabled"`
	PasswordLoginEnabled bool     `json:"password_login_enabled"`
	SSOIssuerURL         string   `json:"sso_issuer_url"`
	SSOClientID          string   `json:"sso_client_id"`
	SSOClientSecret      string   `json:"sso_client_secret"`
	SSORedirectURI       string   `json:"sso_redirect_uri"`
	SSOAllowedDomains    []string `json:"sso_allowed_domains"`
	SSOAutoCreateUser    bool     `json:"sso_auto_create_user"`

	// 易支付设置
	EpayEnabled     bool   `json:"epay_enabled"`
	EpayApiURL      string `json:"epay_api_url"`
	EpayMerchantID  string `json:"epay_merchant_id"`
	EpayMerchantKey string `json:"epay_merchant_key"`
	EpayNotifyURL   string `json:"epay_notify_url"`
	EpayReturnURL   string `json:"epay_return_url"`

	// 支付渠道配置
	PaymentChannels []dto.PaymentChannel `json:"payment_channels"`
}

// UpdateSettings 更新系统设置
// PUT /api/v1/admin/settings
func (h *SettingHandler) UpdateSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// 验证参数
	if req.DefaultConcurrency < 1 {
		req.DefaultConcurrency = 1
	}
	if req.DefaultBalance < 0 {
		req.DefaultBalance = 0
	}
	if req.SmtpPort <= 0 {
		req.SmtpPort = 587
	}

	// Turnstile 参数验证
	if req.TurnstileEnabled {
		// 检查必填字段
		if req.TurnstileSiteKey == "" {
			response.BadRequest(c, "Turnstile Site Key is required when enabled")
			return
		}
		if req.TurnstileSecretKey == "" {
			response.BadRequest(c, "Turnstile Secret Key is required when enabled")
			return
		}

		// 获取当前设置，检查参数是否有变化
		currentSettings, err := h.settingService.GetAllSettings(c.Request.Context())
		if err != nil {
			response.ErrorFrom(c, err)
			return
		}

		// 当 site_key 或 secret_key 任一变化时验证（避免配置错误导致无法登录）
		siteKeyChanged := currentSettings.TurnstileSiteKey != req.TurnstileSiteKey
		secretKeyChanged := currentSettings.TurnstileSecretKey != req.TurnstileSecretKey
		if siteKeyChanged || secretKeyChanged {
			if err := h.turnstileService.ValidateSecretKey(c.Request.Context(), req.TurnstileSecretKey); err != nil {
				response.ErrorFrom(c, err)
				return
			}
		}
	}

	settings := &service.SystemSettings{
		RegistrationEnabled: req.RegistrationEnabled,
		EmailVerifyEnabled:  req.EmailVerifyEnabled,
		SmtpHost:            req.SmtpHost,
		SmtpPort:            req.SmtpPort,
		SmtpUsername:        req.SmtpUsername,
		SmtpPassword:        req.SmtpPassword,
		SmtpFrom:            req.SmtpFrom,
		SmtpFromName:        req.SmtpFromName,
		SmtpUseTLS:          req.SmtpUseTLS,
		TurnstileEnabled:    req.TurnstileEnabled,
		TurnstileSiteKey:    req.TurnstileSiteKey,
		TurnstileSecretKey:  req.TurnstileSecretKey,
		SiteName:            req.SiteName,
		SiteLogo:            req.SiteLogo,
		SiteSubtitle:        req.SiteSubtitle,
		ApiBaseUrl:          req.ApiBaseUrl,
		ContactInfo:         req.ContactInfo,
		DocUrl:              req.DocUrl,
		DefaultConcurrency:  req.DefaultConcurrency,
		DefaultBalance:      req.DefaultBalance,
		// SSO设置
		SSOEnabled:           req.SSOEnabled,
		PasswordLoginEnabled: req.PasswordLoginEnabled,
		SSOIssuerURL:         req.SSOIssuerURL,
		SSOClientID:          req.SSOClientID,
		SSOClientSecret:      req.SSOClientSecret,
		SSORedirectURI:       req.SSORedirectURI,
		SSOAllowedDomains:    req.SSOAllowedDomains,
		SSOAutoCreateUser:    req.SSOAutoCreateUser,
		// 易支付设置
		EpayEnabled:     req.EpayEnabled,
		EpayApiURL:      req.EpayApiURL,
		EpayMerchantID:  req.EpayMerchantID,
		EpayMerchantKey: req.EpayMerchantKey,
		EpayNotifyURL:   req.EpayNotifyURL,
		EpayReturnURL:   req.EpayReturnURL,
		// 支付渠道
		PaymentChannels: dtoChannelsToService(req.PaymentChannels),
	}

	if err := h.settingService.UpdateSettings(c.Request.Context(), settings); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// 重新获取设置返回
	updatedSettings, err := h.settingService.GetAllSettings(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, dto.SystemSettings{
		RegistrationEnabled: updatedSettings.RegistrationEnabled,
		EmailVerifyEnabled:  updatedSettings.EmailVerifyEnabled,
		SmtpHost:            updatedSettings.SmtpHost,
		SmtpPort:            updatedSettings.SmtpPort,
		SmtpUsername:        updatedSettings.SmtpUsername,
		SmtpPassword:        maskSensitiveString(updatedSettings.SmtpPassword), // 掩码
		SmtpFrom:            updatedSettings.SmtpFrom,
		SmtpFromName:        updatedSettings.SmtpFromName,
		SmtpUseTLS:          updatedSettings.SmtpUseTLS,
		TurnstileEnabled:    updatedSettings.TurnstileEnabled,
		TurnstileSiteKey:    updatedSettings.TurnstileSiteKey,
		TurnstileSecretKey:  maskSensitiveString(updatedSettings.TurnstileSecretKey), // 掩码
		SiteName:            updatedSettings.SiteName,
		SiteLogo:            updatedSettings.SiteLogo,
		SiteSubtitle:        updatedSettings.SiteSubtitle,
		ApiBaseUrl:          updatedSettings.ApiBaseUrl,
		ContactInfo:         updatedSettings.ContactInfo,
		DocUrl:              updatedSettings.DocUrl,
		DefaultConcurrency:  updatedSettings.DefaultConcurrency,
		DefaultBalance:      updatedSettings.DefaultBalance,
		// SSO设置
		SSOEnabled:           updatedSettings.SSOEnabled,
		PasswordLoginEnabled: updatedSettings.PasswordLoginEnabled,
		SSOIssuerURL:         updatedSettings.SSOIssuerURL,
		SSOClientID:          updatedSettings.SSOClientID,
		SSOClientSecret:      maskSensitiveString(updatedSettings.SSOClientSecret), // 掩码
		SSORedirectURI:       updatedSettings.SSORedirectURI,
		SSOAllowedDomains:    updatedSettings.SSOAllowedDomains,
		SSOAutoCreateUser:    updatedSettings.SSOAutoCreateUser,
		// 易支付设置
		EpayEnabled:     updatedSettings.EpayEnabled,
		EpayApiURL:      updatedSettings.EpayApiURL,
		EpayMerchantID:  updatedSettings.EpayMerchantID,
		EpayMerchantKey: maskSensitiveString(updatedSettings.EpayMerchantKey), // 掩码
		EpayNotifyURL:   updatedSettings.EpayNotifyURL,
		EpayReturnURL:   updatedSettings.EpayReturnURL,
		// 支付渠道
		PaymentChannels: serviceChannelsToDTO(updatedSettings.PaymentChannels),
	})
}

// TestSmtpRequest 测试SMTP连接请求
type TestSmtpRequest struct {
	SmtpHost     string `json:"smtp_host" binding:"required"`
	SmtpPort     int    `json:"smtp_port"`
	SmtpUsername string `json:"smtp_username"`
	SmtpPassword string `json:"smtp_password"`
	SmtpUseTLS   bool   `json:"smtp_use_tls"`
}

// TestSmtpConnection 测试SMTP连接
// POST /api/v1/admin/settings/test-smtp
func (h *SettingHandler) TestSmtpConnection(c *gin.Context) {
	var req TestSmtpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if req.SmtpPort <= 0 {
		req.SmtpPort = 587
	}

	// 如果未提供密码，从数据库获取已保存的密码
	password := req.SmtpPassword
	if password == "" {
		savedConfig, err := h.emailService.GetSmtpConfig(c.Request.Context())
		if err == nil && savedConfig != nil {
			password = savedConfig.Password
		}
	}

	config := &service.SmtpConfig{
		Host:     req.SmtpHost,
		Port:     req.SmtpPort,
		Username: req.SmtpUsername,
		Password: password,
		UseTLS:   req.SmtpUseTLS,
	}

	err := h.emailService.TestSmtpConnectionWithConfig(config)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "SMTP connection successful"})
}

// SendTestEmailRequest 发送测试邮件请求
type SendTestEmailRequest struct {
	Email        string `json:"email" binding:"required,email"`
	SmtpHost     string `json:"smtp_host" binding:"required"`
	SmtpPort     int    `json:"smtp_port"`
	SmtpUsername string `json:"smtp_username"`
	SmtpPassword string `json:"smtp_password"`
	SmtpFrom     string `json:"smtp_from_email"`
	SmtpFromName string `json:"smtp_from_name"`
	SmtpUseTLS   bool   `json:"smtp_use_tls"`
}

// SendTestEmail 发送测试邮件
// POST /api/v1/admin/settings/send-test-email
func (h *SettingHandler) SendTestEmail(c *gin.Context) {
	var req SendTestEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if req.SmtpPort <= 0 {
		req.SmtpPort = 587
	}

	// 如果未提供密码，从数据库获取已保存的密码
	password := req.SmtpPassword
	if password == "" {
		savedConfig, err := h.emailService.GetSmtpConfig(c.Request.Context())
		if err == nil && savedConfig != nil {
			password = savedConfig.Password
		}
	}

	config := &service.SmtpConfig{
		Host:     req.SmtpHost,
		Port:     req.SmtpPort,
		Username: req.SmtpUsername,
		Password: password,
		From:     req.SmtpFrom,
		FromName: req.SmtpFromName,
		UseTLS:   req.SmtpUseTLS,
	}

	siteName := h.settingService.GetSiteName(c.Request.Context())
	subject := "[" + siteName + "] Test Email"
	body := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f5f5f5; margin: 0; padding: 20px; }
        .container { max-width: 600px; margin: 0 auto; background-color: #ffffff; border-radius: 8px; overflow: hidden; box-shadow: 0 2px 8px rgba(0,0,0,0.1); }
        .header { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); color: white; padding: 30px; text-align: center; }
        .content { padding: 40px 30px; text-align: center; }
        .success { color: #10b981; font-size: 48px; margin-bottom: 20px; }
        .footer { background-color: #f8f9fa; padding: 20px; text-align: center; color: #999; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>` + siteName + `</h1>
        </div>
        <div class="content">
            <div class="success">✓</div>
            <h2>Email Configuration Successful!</h2>
            <p>This is a test email to verify your SMTP settings are working correctly.</p>
        </div>
        <div class="footer">
            <p>This is an automated test message.</p>
        </div>
    </div>
</body>
</html>
`

	if err := h.emailService.SendEmailWithConfig(config, req.Email, subject, body); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Test email sent successfully"})
}

// GetAdminApiKey 获取管理员 API Key 状态
// GET /api/v1/admin/settings/admin-api-key
func (h *SettingHandler) GetAdminApiKey(c *gin.Context) {
	maskedKey, exists, err := h.settingService.GetAdminApiKeyStatus(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"exists":     exists,
		"masked_key": maskedKey,
	})
}

// RegenerateAdminApiKey 生成/重新生成管理员 API Key
// POST /api/v1/admin/settings/admin-api-key/regenerate
func (h *SettingHandler) RegenerateAdminApiKey(c *gin.Context) {
	key, err := h.settingService.GenerateAdminApiKey(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"key": key, // 完整 key 只在生成时返回一次
	})
}

// DeleteAdminApiKey 删除管理员 API Key
// DELETE /api/v1/admin/settings/admin-api-key
func (h *SettingHandler) DeleteAdminApiKey(c *gin.Context) {
	if err := h.settingService.DeleteAdminApiKey(c.Request.Context()); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Admin API key deleted"})
}

// TestSSORequest 测试SSO连接请求
type TestSSORequest struct {
	IssuerURL string `json:"issuer_url" binding:"required"`
}

// TestSSOConnection 测试SSO配置
// POST /api/v1/admin/settings/test-sso
func (h *SettingHandler) TestSSOConnection(c *gin.Context) {
	var req TestSSORequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// 使用OIDC客户端验证配置
	oidcClient := oidc.NewOIDCClient()
	config, err := oidcClient.DiscoverOIDCConfig(c.Request.Context(), req.IssuerURL)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{
		"message": "SSO configuration is valid",
		"issuer":  config.Issuer,
	})
}

// UpdateSingleSettingRequest 更新单个配置项请求
type UpdateSingleSettingRequest struct {
	Value any `json:"value" binding:"required"`
}

// UpdateSingleSetting 更新单个配置项（用于实时保存）
// PATCH /api/v1/admin/settings/:key
func (h *SettingHandler) UpdateSingleSetting(c *gin.Context) {
	key := c.Param("key")
	if key == "" {
		response.BadRequest(c, "Setting key is required")
		return
	}

	var req UpdateSingleSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	// 将value转换为字符串
	value := fmt.Sprintf("%v", req.Value)

	// 布尔值特殊处理
	if boolVal, ok := req.Value.(bool); ok {
		if boolVal {
			value = "true"
		} else {
			value = "false"
		}
	}

	// 更新单个配置
	if err := h.settingService.SetSetting(c.Request.Context(), key, value); err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Setting updated successfully"})
}

// dtoChannelsToService 将 DTO PaymentChannel 转换为 service PaymentChannel
func dtoChannelsToService(channels []dto.PaymentChannel) []service.PaymentChannel {
	if channels == nil {
		return nil
	}
	result := make([]service.PaymentChannel, len(channels))
	for i, ch := range channels {
		result[i] = service.PaymentChannel{
			Key:         ch.Key,
			DisplayName: ch.DisplayName,
			EpayType:    ch.EpayType,
			Icon:        ch.Icon,
			Enabled:     ch.Enabled,
			SortOrder:   ch.SortOrder,
		}
	}
	return result
}

// serviceChannelsToDTO 将 service PaymentChannel 转换为 DTO PaymentChannel
func serviceChannelsToDTO(channels []service.PaymentChannel) []dto.PaymentChannel {
	if channels == nil {
		return nil
	}
	result := make([]dto.PaymentChannel, len(channels))
	for i, ch := range channels {
		result[i] = dto.PaymentChannel{
			Key:         ch.Key,
			DisplayName: ch.DisplayName,
			EpayType:    ch.EpayType,
			Icon:        ch.Icon,
			Enabled:     ch.Enabled,
			SortOrder:   ch.SortOrder,
		}
	}
	return result
}
