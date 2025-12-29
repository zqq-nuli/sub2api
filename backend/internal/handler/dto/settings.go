package dto

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
