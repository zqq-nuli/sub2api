package service

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
