package oidc

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"crypto/rand"
)

// OIDCClient handles OIDC protocol operations
type OIDCClient struct {
	httpClient *http.Client
}

// NewOIDCClient creates a new OIDC client
func NewOIDCClient() *OIDCClient {
	return &OIDCClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// OIDCConfig contains OIDC client configuration
type OIDCConfig struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

// OIDCProviderConfig represents the OIDC provider's configuration
type OIDCProviderConfig struct {
	Issuer                string `json:"issuer"`
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
	JWKSURI               string `json:"jwks_uri"`
	UserinfoEndpoint      string `json:"userinfo_endpoint,omitempty"`
}

// OIDCTokenResponse represents the token response from the provider
type OIDCTokenResponse struct {
	AccessToken  string `json:"access_token"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// IDTokenClaims represents claims from an ID Token
type IDTokenClaims struct {
	Sub               string   `json:"sub"`
	Email             string   `json:"email"`
	EmailVerified     bool     `json:"email_verified"`
	Name              string   `json:"name,omitempty"`
	PreferredUsername string   `json:"preferred_username,omitempty"`
	Iss               string   `json:"iss"`
	Aud               AudClaim `json:"aud"`
	Exp               int64    `json:"exp"`
	Iat               int64    `json:"iat"`

	// Forum-specific fields
	ID             string                 `json:"id,omitempty"`              // User unique ID
	Username       string                 `json:"username,omitempty"`        // Forum username
	AvatarTemplate string                 `json:"avatar_template,omitempty"` // Avatar URL template
	Active         bool                   `json:"active,omitempty"`          // Account active status
	TrustLevel     int                    `json:"trust_level,omitempty"`     // Trust level (0-4)
	Silenced       bool                   `json:"silenced,omitempty"`        // Silenced status
	ExternalIDs    map[string]interface{} `json:"external_ids,omitempty"`    // External ID associations
	APIKey         string                 `json:"api_key,omitempty"`         // API access key
}

// AudClaim handles aud field which can be string or []string
type AudClaim []string

// UnmarshalJSON implements json.Unmarshaler for AudClaim
func (a *AudClaim) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*a = []string{s}
		return nil
	}

	// Try to unmarshal as array
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		*a = arr
		return nil
	}

	return fmt.Errorf("aud must be string or array of strings")
}

// String returns the first audience value
func (a AudClaim) String() string {
	if len(a) > 0 {
		return a[0]
	}
	return ""
}

// DiscoverOIDCConfig discovers OIDC configuration from the issuer's .well-known endpoint
func (c *OIDCClient) DiscoverOIDCConfig(ctx context.Context, issuerURL string) (*OIDCProviderConfig, error) {
	wellKnownURL := strings.TrimRight(issuerURL, "/") + "/.well-known/openid-configuration"

	req, err := http.NewRequestWithContext(ctx, "GET", wellKnownURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create discovery request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to discover OIDC configuration: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("OIDC discovery returned status %d: %s", resp.StatusCode, string(body))
	}

	var config OIDCProviderConfig
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode OIDC configuration: %w", err)
	}

	return &config, nil
}

// ExchangeCode exchanges an authorization code for tokens
func (c *OIDCClient) ExchangeCode(ctx context.Context, config *OIDCConfig, providerConfig *OIDCProviderConfig, code, codeVerifier, redirectURI string) (*OIDCTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)
	data.Set("code_verifier", codeVerifier)

	req, err := http.NewRequestWithContext(ctx, "POST", providerConfig.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create token request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token exchange failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp OIDCTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode token response: %w", err)
	}

	return &tokenResp, nil
}

// RefreshToken refreshes an access token using a refresh token
func (c *OIDCClient) RefreshToken(ctx context.Context, config *OIDCConfig, providerConfig *OIDCProviderConfig, refreshToken string) (*OIDCTokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", config.ClientID)
	data.Set("client_secret", config.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", providerConfig.TokenEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create refresh request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp OIDCTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to decode refresh response: %w", err)
	}

	return &tokenResp, nil
}

// ParseIDToken parses an ID Token and extracts claims (without signature verification)
func (c *OIDCClient) ParseIDToken(idToken string) (*IDTokenClaims, error) {
	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format")
	}

	// Decode payload (second part)
	payload := parts[1]

	// Add padding if necessary
	switch len(payload) % 4 {
	case 2:
		payload += "=="
	case 3:
		payload += "="
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode JWT payload: %w", err)
	}

	var claims IDTokenClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse JWT claims: %w", err)
	}

	// Basic expiration check
	if claims.Exp > 0 && time.Now().Unix() > claims.Exp {
		return nil, fmt.Errorf("ID token has expired")
	}

	return &claims, nil
}

// GenerateState generates a random state string for CSRF protection
func GenerateState() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate state: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateCodeVerifier generates a PKCE code verifier
func GenerateCodeVerifier() (string, error) {
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate code verifier: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateCodeChallenge generates a PKCE code challenge from a verifier
func GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64URLEncode(hash[:])
}

// base64URLEncode encodes data in base64 URL encoding without padding
func base64URLEncode(data []byte) string {
	return base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(data)
}

// BuildAuthorizationURL builds the OIDC authorization URL
func BuildAuthorizationURL(authEndpoint, clientID, redirectURI, state, codeChallenge string, scopes []string) string {
	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("state", state)
	params.Set("code_challenge", codeChallenge)
	params.Set("code_challenge_method", "S256")

	if len(scopes) > 0 {
		params.Set("scope", strings.Join(scopes, " "))
	} else {
		// Default OIDC scopes
		params.Set("scope", "openid email profile")
	}

	return authEndpoint + "?" + params.Encode()
}
