package service

// PaymentChannel 支付渠道配置
type PaymentChannel struct {
	Key         string `json:"key"`          // 渠道标识：alipay, wxpay, epusdt
	DisplayName string `json:"display_name"` // 显示名称：支付宝、微信支付、USDT
	EpayType    string `json:"epay_type"`    // 网关参数：epay, alipay, wxpay 等
	Icon        string `json:"icon"`         // 图标标识
	Enabled     bool   `json:"enabled"`      // 是否启用
	SortOrder   int    `json:"sort_order"`   // 排序
}

type SystemSettings struct {
	RegistrationEnabled bool
	EmailVerifyEnabled  bool

	SmtpHost     string
	SmtpPort     int
	SmtpUsername string
	SmtpPassword string
	SmtpFrom     string
	SmtpFromName string
	SmtpUseTLS   bool

	TurnstileEnabled   bool
	TurnstileSiteKey   string
	TurnstileSecretKey string

	SiteName     string
	SiteLogo     string
	SiteSubtitle string
	ApiBaseUrl   string
	ContactInfo  string
	DocUrl       string

	DefaultConcurrency int
	DefaultBalance     float64

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
	RegistrationEnabled bool
	EmailVerifyEnabled  bool
	TurnstileEnabled    bool
	TurnstileSiteKey    string
	SiteName            string
	SiteLogo            string
	SiteSubtitle        string
	ApiBaseUrl          string
	ContactInfo         string
	DocUrl              string
	Version             string
}
