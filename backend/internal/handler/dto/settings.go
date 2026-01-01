package dto

// PaymentChannel 支付渠道配置
type PaymentChannel struct {
	Key         string `json:"key"`          // 渠道标识：alipay, wxpay, epusdt
	DisplayName string `json:"display_name"` // 显示名称：支付宝、微信支付、USDT
	EpayType    string `json:"epay_type"`    // 网关参数：epay, alipay, wxpay 等
	Icon        string `json:"icon"`         // 图标标识
	Enabled     bool   `json:"enabled"`      // 是否启用
	SortOrder   int    `json:"sort_order"`   // 排序
}

// SystemSettings represents the admin settings API response payload.
type SystemSettings struct {
	RegistrationEnabled bool `json:"registration_enabled"`
	EmailVerifyEnabled  bool `json:"email_verify_enabled"`

	SmtpHost     string `json:"smtp_host"`
	SmtpPort     int    `json:"smtp_port"`
	SmtpUsername string `json:"smtp_username"`
	SmtpPassword string `json:"smtp_password,omitempty"`
	SmtpFrom     string `json:"smtp_from_email"`
	SmtpFromName string `json:"smtp_from_name"`
	SmtpUseTLS   bool   `json:"smtp_use_tls"`

	TurnstileEnabled   bool   `json:"turnstile_enabled"`
	TurnstileSiteKey   string `json:"turnstile_site_key"`
	TurnstileSecretKey string `json:"turnstile_secret_key,omitempty"`

	SiteName     string `json:"site_name"`
	SiteLogo     string `json:"site_logo"`
	SiteSubtitle string `json:"site_subtitle"`
	ApiBaseUrl   string `json:"api_base_url"`
	ContactInfo  string `json:"contact_info"`
	DocUrl       string `json:"doc_url"`

	DefaultConcurrency int     `json:"default_concurrency"`
	DefaultBalance     float64 `json:"default_balance"`

	// SSO设置
	SSOEnabled           bool     `json:"sso_enabled"`
	PasswordLoginEnabled bool     `json:"password_login_enabled"`
	SSOIssuerURL         string   `json:"sso_issuer_url"`
	SSOClientID          string   `json:"sso_client_id"`
	SSOClientSecret      string   `json:"sso_client_secret,omitempty"`
	SSORedirectURI       string   `json:"sso_redirect_uri"`
	SSOAllowedDomains    []string `json:"sso_allowed_domains"`
	SSOAutoCreateUser    bool     `json:"sso_auto_create_user"`

	// 易支付设置
	EpayEnabled     bool   `json:"epay_enabled"`
	EpayApiURL      string `json:"epay_api_url"`
	EpayMerchantID  string `json:"epay_merchant_id"`
	EpayMerchantKey string `json:"epay_merchant_key,omitempty"`
	EpayNotifyURL   string `json:"epay_notify_url"`
	EpayReturnURL   string `json:"epay_return_url"`

	// 支付渠道配置
	PaymentChannels []PaymentChannel `json:"payment_channels"`
}

type PublicSettings struct {
	RegistrationEnabled bool   `json:"registration_enabled"`
	EmailVerifyEnabled  bool   `json:"email_verify_enabled"`
	TurnstileEnabled    bool   `json:"turnstile_enabled"`
	TurnstileSiteKey    string `json:"turnstile_site_key"`
	SiteName            string `json:"site_name"`
	SiteLogo            string `json:"site_logo"`
	SiteSubtitle        string `json:"site_subtitle"`
	ApiBaseUrl          string `json:"api_base_url"`
	ContactInfo         string `json:"contact_info"`
	DocUrl              string `json:"doc_url"`
	Version             string `json:"version"`
}
