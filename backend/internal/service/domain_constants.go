package service

// Status constants
const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
	StatusError    = "error"
	StatusUnused   = "unused"
	StatusUsed     = "used"
	StatusExpired  = "expired"
)

// Role constants
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Platform constants
const (
	PlatformAnthropic   = "anthropic"
	PlatformOpenAI      = "openai"
	PlatformGemini      = "gemini"
	PlatformAntigravity = "antigravity"
)

// Account type constants
const (
	AccountTypeOAuth      = "oauth"       // OAuth类型账号（full scope: profile + inference）
	AccountTypeSetupToken = "setup-token" // Setup Token类型账号（inference only scope）
	AccountTypeApiKey     = "apikey"      // API Key类型账号
)

// Redeem type constants
const (
	RedeemTypeBalance      = "balance"
	RedeemTypeConcurrency  = "concurrency"
	RedeemTypeSubscription = "subscription"
)

// Admin adjustment type constants
const (
	AdjustmentTypeAdminBalance     = "admin_balance"     // 管理员调整余额
	AdjustmentTypeAdminConcurrency = "admin_concurrency" // 管理员调整并发数
)

// Group subscription type constants
const (
	SubscriptionTypeStandard     = "standard"     // 标准计费模式（按余额扣费）
	SubscriptionTypeSubscription = "subscription" // 订阅模式（按限额控制）
)

// Subscription status constants
const (
	SubscriptionStatusActive    = "active"
	SubscriptionStatusExpired   = "expired"
	SubscriptionStatusSuspended = "suspended"
)

// Setting keys
const (
	// 注册设置
	SettingKeyRegistrationEnabled = "registration_enabled" // 是否开放注册
	SettingKeyEmailVerifyEnabled  = "email_verify_enabled" // 是否开启邮件验证

	// 邮件服务设置
	SettingKeySmtpHost     = "smtp_host"      // SMTP服务器地址
	SettingKeySmtpPort     = "smtp_port"      // SMTP端口
	SettingKeySmtpUsername = "smtp_username"  // SMTP用户名
	SettingKeySmtpPassword = "smtp_password"  // SMTP密码（加密存储）
	SettingKeySmtpFrom     = "smtp_from"      // 发件人地址
	SettingKeySmtpFromName = "smtp_from_name" // 发件人名称
	SettingKeySmtpUseTLS   = "smtp_use_tls"   // 是否使用TLS

	// Cloudflare Turnstile 设置
	SettingKeyTurnstileEnabled   = "turnstile_enabled"    // 是否启用 Turnstile 验证
	SettingKeyTurnstileSiteKey   = "turnstile_site_key"   // Turnstile Site Key
	SettingKeyTurnstileSecretKey = "turnstile_secret_key" // Turnstile Secret Key

	// OEM设置
	SettingKeySiteName     = "site_name"     // 网站名称
	SettingKeySiteLogo     = "site_logo"     // 网站Logo (base64)
	SettingKeySiteSubtitle = "site_subtitle" // 网站副标题
	SettingKeyApiBaseUrl   = "api_base_url"  // API端点地址（用于客户端配置和导入）
	SettingKeyContactInfo  = "contact_info"  // 客服联系方式
	SettingKeyDocUrl       = "doc_url"       // 文档链接

	// 默认配置
	SettingKeyDefaultConcurrency = "default_concurrency" // 新用户默认并发量
	SettingKeyDefaultBalance     = "default_balance"     // 新用户默认余额

	// 管理员 API Key
	SettingKeyAdminApiKey = "admin_api_key" // 全局管理员 API Key（用于外部系统集成）

	// SSO设置
	SettingKeySSOEnabled           = "sso_enabled"             // 是否启用SSO登录
	SettingKeyPasswordLoginEnabled = "password_login_enabled"  // 是否启用密码登录
	SettingKeySSOIssuerURL         = "sso_issuer_url"          // OIDC Issuer URL
	SettingKeySSOClientID          = "sso_client_id"           // OIDC Client ID
	SettingKeySSOClientSecret      = "sso_client_secret"       // OIDC Client Secret
	SettingKeySSORedirectURI       = "sso_redirect_uri"        // OIDC Redirect URI
	SettingKeySSOAllowedDomains    = "sso_allowed_domains"     // 允许的邮箱域名（JSON数组）
	SettingKeySSOAutoCreateUser    = "sso_auto_create_user"    // 是否自动创建用户
	SettingKeySSOMinTrustLevel     = "sso_min_trust_level"     // 最小信任等级（0-4）

	// Gemini 配额策略（JSON）
	SettingKeyGeminiQuotaPolicy = "gemini_quota_policy"
)

// Admin API Key prefix (distinct from user "sk-" keys)
const AdminApiKeyPrefix = "admin-"

// 易支付设置
const (
	SettingKeyEpayEnabled        = "epay_enabled"        // 是否启用易支付
	SettingKeyEpayApiURL         = "epay_api_url"        // 易支付API地址
	SettingKeyEpayMerchantID     = "epay_merchant_id"    // 商户ID
	SettingKeyEpayMerchantKey    = "epay_merchant_key"   // 商户密钥（加密存储）
	SettingKeyEpayNotifyURL      = "epay_notify_url"     // 异步回调URL
	SettingKeyEpayReturnURL      = "epay_return_url"     // 同步回调URL
	SettingKeyPaymentChannels    = "payment_channels"    // 支付渠道配置（JSON数组）
)

// GetDefaultPaymentChannels 获取默认支付渠道配置
func GetDefaultPaymentChannels() []PaymentChannel {
	return []PaymentChannel{
		{
			Key:         "alipay",
			DisplayName: "支付宝",
			EpayType:    "epay",
			Icon:        "alipay",
			Enabled:     true,
			SortOrder:   1,
		},
		{
			Key:         "wxpay",
			DisplayName: "微信支付",
			EpayType:    "epay",
			Icon:        "wechat",
			Enabled:     true,
			SortOrder:   2,
		},
		{
			Key:         "epusdt",
			DisplayName: "USDT",
			EpayType:    "epay",
			Icon:        "usdt",
			Enabled:     true,
			SortOrder:   3,
		},
	}
}
