package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/infrastructure/errors"
	"github.com/Wei-Shaw/sub2api/internal/pkg/oidc"
)

// OIDCSSOService handles SSO authentication using OIDC for user login
type OIDCSSOService struct {
	sessionStore   *OIDCSessionStore
	oidcClient     *oidc.OIDCClient
	settingService *SettingService
	userRepo       UserRepository
	authService    *AuthService
	stopCh         chan struct{}
}

// NewOIDCSSOService creates a new OIDC SSO service
func NewOIDCSSOService(
	settingService *SettingService,
	userRepo UserRepository,
	authService *AuthService,
) *OIDCSSOService {
	service := &OIDCSSOService{
		sessionStore:   NewOIDCSessionStore(),
		oidcClient:     oidc.NewOIDCClient(),
		settingService: settingService,
		userRepo:       userRepo,
		authService:    authService,
		stopCh:         make(chan struct{}),
	}

	// Start cleanup goroutine
	go service.sessionStore.StartCleanup(service.stopCh)

	return service
}

// Stop stops the SSO service
func (s *OIDCSSOService) Stop() {
	close(s.stopCh)
}

// OIDCSessionStore manages OIDC sessions
type OIDCSessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*OIDCSession
}

// OIDCSession represents an OIDC authentication session
type OIDCSession struct {
	State        string
	CodeVerifier string
	RedirectURI  string
	CreatedAt    time.Time
}

// NewOIDCSessionStore creates a new session store
func NewOIDCSessionStore() *OIDCSessionStore {
	return &OIDCSessionStore{
		sessions: make(map[string]*OIDCSession),
	}
}

// Set stores a session
func (store *OIDCSessionStore) Set(sessionID string, session *OIDCSession) {
	store.mu.Lock()
	defer store.mu.Unlock()
	store.sessions[sessionID] = session
}

// Get retrieves a session
func (store *OIDCSessionStore) Get(sessionID string) (*OIDCSession, bool) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	session, ok := store.sessions[sessionID]
	return session, ok
}

// Delete removes a session
func (store *OIDCSessionStore) Delete(sessionID string) {
	store.mu.Lock()
	defer store.mu.Unlock()
	delete(store.sessions, sessionID)
}

// StartCleanup periodically cleans up expired sessions
func (store *OIDCSessionStore) StartCleanup(stopCh chan struct{}) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			store.cleanup()
		case <-stopCh:
			return
		}
	}
}

func (store *OIDCSessionStore) cleanup() {
	store.mu.Lock()
	defer store.mu.Unlock()

	now := time.Now()
	for id, session := range store.sessions {
		if now.Sub(session.CreatedAt) > 30*time.Minute {
			delete(store.sessions, id)
		}
	}
}

// OIDCAuthURLResult contains the authorization URL and session info
type OIDCAuthURLResult struct {
	AuthURL   string `json:"auth_url"`
	SessionID string `json:"session_id"`
}

// GenerateAuthURL generates an OIDC authorization URL
func (s *OIDCSSOService) GenerateAuthURL(ctx context.Context) (*OIDCAuthURLResult, error) {
	// Get SSO configuration
	ssoConfig, err := s.getSSOConfig(ctx)
	if err != nil {
		return nil, err
	}

	// Discover OIDC provider configuration
	providerConfig, err := s.oidcClient.DiscoverOIDCConfig(ctx, ssoConfig.IssuerURL)
	if err != nil {
		return nil, infraerrors.BadRequest("OIDC_DISCOVERY_FAILED",
			fmt.Sprintf("Failed to discover OIDC configuration: %v", err))
	}

	// Generate PKCE values
	state, err := oidc.GenerateState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate state: %w", err)
	}

	codeVerifier, err := oidc.GenerateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %w", err)
	}

	codeChallenge := oidc.GenerateCodeChallenge(codeVerifier)

	// Generate session ID
	sessionID, err := oidc.GenerateState()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	// Store session
	session := &OIDCSession{
		State:        state,
		CodeVerifier: codeVerifier,
		RedirectURI:  ssoConfig.RedirectURI,
		CreatedAt:    time.Now(),
	}
	s.sessionStore.Set(sessionID, session)

	// Build authorization URL
	authURL := oidc.BuildAuthorizationURL(
		providerConfig.AuthorizationEndpoint,
		ssoConfig.ClientID,
		ssoConfig.RedirectURI,
		state,
		codeChallenge,
		nil, // Use default scopes
	)

	return &OIDCAuthURLResult{
		AuthURL:   authURL,
		SessionID: sessionID,
	}, nil
}

// ExchangeCodeAndCreateUser exchanges an authorization code and creates/logs in a user
func (s *OIDCSSOService) ExchangeCodeAndCreateUser(ctx context.Context, code, state, sessionID string) (string, *User, bool, error) {
	// Retrieve session
	session, ok := s.sessionStore.Get(sessionID)
	if !ok {
		return "", nil, false, infraerrors.BadRequest("SESSION_EXPIRED",
			"SSO session not found or expired, please try again")
	}

	// Verify state
	if session.State != state {
		return "", nil, false, infraerrors.BadRequest("INVALID_STATE",
			"Invalid OAuth state, possible CSRF attack")
	}

	// Clean up session
	defer s.sessionStore.Delete(sessionID)

	// Get SSO configuration
	ssoConfig, err := s.getSSOConfig(ctx)
	if err != nil {
		return "", nil, false, err
	}

	// Discover OIDC provider configuration
	providerConfig, err := s.oidcClient.DiscoverOIDCConfig(ctx, ssoConfig.IssuerURL)
	if err != nil {
		return "", nil, false, infraerrors.BadRequest("OIDC_DISCOVERY_FAILED",
			fmt.Sprintf("Failed to discover OIDC configuration: %v", err))
	}

	// Exchange code for tokens
	tokenResp, err := s.oidcClient.ExchangeCode(ctx,
		&oidc.OIDCConfig{
			IssuerURL:    ssoConfig.IssuerURL,
			ClientID:     ssoConfig.ClientID,
			ClientSecret: ssoConfig.ClientSecret,
			RedirectURI:  session.RedirectURI,
		},
		providerConfig,
		code,
		session.CodeVerifier,
		session.RedirectURI,
	)
	if err != nil {
		return "", nil, false, infraerrors.BadRequest("TOKEN_EXCHANGE_FAILED",
			fmt.Sprintf("Failed to exchange authorization code: %v", err))
	}

	// Parse ID Token to get user claims
	claims, err := s.oidcClient.ParseIDToken(tokenResp.IDToken)
	if err != nil {
		return "", nil, false, infraerrors.BadRequest("INVALID_ID_TOKEN",
			fmt.Sprintf("Failed to parse ID token: %v", err))
	}

	// Serialize complete claims to JSON for storage (save all fields regardless of what they are)
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		log.Printf("Warning: Failed to serialize claims to JSON: %v", err)
		claimsJSON = []byte("{}")
	}

	// Log all available information from claims
	log.Printf("[SSO] Parsed ID Token claims: Sub=%s, Email=%s, ID=%s, Username=%s, Name=%s, PreferredUsername=%s, TrustLevel=%d, Silenced=%t, Active=%t, AvatarTemplate=%s",
		claims.Sub, claims.Email, claims.ID, claims.Username, claims.Name, claims.PreferredUsername, claims.TrustLevel, claims.Silenced, claims.Active, claims.AvatarTemplate)
	log.Printf("[SSO] Full claims JSON: %s", string(claimsJSON))

	// Validate silenced status (Discourse forum requirement)
	if claims.Silenced {
		return "", nil, false, infraerrors.Forbidden("USER_SILENCED",
			"Your account is silenced and cannot register")
	}

	// Validate trust level (Discourse forum requirement)
	minTrustLevelStr, _ := s.settingService.GetSetting(ctx, SettingKeySSOMinTrustLevel)
	minTrustLevel := 0
	if minTrustLevelStr != "" {
		fmt.Sscanf(minTrustLevelStr, "%d", &minTrustLevel)
	}
	if claims.TrustLevel < minTrustLevel {
		return "", nil, false, infraerrors.Forbidden("TRUST_LEVEL_TOO_LOW",
			fmt.Sprintf("Minimum trust level required: %d, your level: %d", minTrustLevel, claims.TrustLevel))
	}

	// Use ID field as email (for Discourse forums), fallback to standard email claim
	email := claims.ID
	if email == "" {
		email = claims.Email
	}

	log.Printf("[SSO] Using email/ID field: %s (from claims.ID=%s, claims.Email=%s)", email, claims.ID, claims.Email)

	// Validate email
	if email == "" {
		return "", nil, false, infraerrors.BadRequest("NO_EMAIL",
			"ID token does not contain email or id claim")
	}

	// Validate email domain if configured (only if email contains @)
	if strings.Contains(email, "@") && !s.validateEmailDomain(ctx, email) {
		domain := strings.Split(email, "@")[1]
		return "", nil, false, infraerrors.Forbidden("EMAIL_DOMAIN_NOT_ALLOWED",
			fmt.Sprintf("Email domain %s is not allowed to login via SSO", domain))
	}

	log.Printf("[SSO] Checking if user exists with email: %s", email)

	// Check if user exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil && !errors.Is(err, ErrUserNotFound) {
		return "", nil, false, fmt.Errorf("failed to check user existence: %w", err)
	}

	if existingUser != nil {
		log.Printf("[SSO] User already exists: ID=%d, Email=%s, HasPassword=%t", existingUser.ID, existingUser.Email, existingUser.PasswordHash != "")
	} else {
		log.Printf("[SSO] User does not exist, will create new user")
	}

	var user *User
	isNewUser := false

	if existingUser != nil {
		// User exists - check if it's a password user
		if existingUser.PasswordHash != "" {
			return "", nil, false, infraerrors.Conflict("EMAIL_EXISTS_WITH_PASSWORD",
				"This email is registered with password login, please use password to login")
		}

		// SSO user exists, use it
		user = existingUser
	} else {
		// User doesn't exist - check if auto-create is enabled
		autoCreate, _ := s.settingService.GetBoolSetting(ctx, SettingKeySSOAutoCreateUser)
		if !autoCreate {
			return "", nil, false, infraerrors.Forbidden("SSO_AUTO_CREATE_DISABLED",
				"SSO auto-create is disabled, please contact administrator")
		}

		// Create new user
		defaultBalance := s.settingService.GetDefaultBalance(ctx)
		defaultConcurrency := s.settingService.GetDefaultConcurrency(ctx)

		// Use Username from claims (Discourse forum), fallback to other fields
		username := claims.Username
		if username == "" {
			username = claims.PreferredUsername
		}
		if username == "" && strings.Contains(email, "@") {
			username = strings.Split(email, "@")[0]
		}
		if username == "" {
			username = email // Last resort
		}

		log.Printf("[SSO] Creating new user with: Email=%s, Username=%s, Avatar=%s, DefaultBalance=%.2f, DefaultConcurrency=%d",
			email, username, claims.AvatarTemplate, defaultBalance, defaultConcurrency)
		log.Printf("[SSO] SSOData length: %d bytes", len(claimsJSON))

		user = &User{
			Email:        email,
			Username:     username,
			Avatar:       claims.AvatarTemplate,
			SSOData:      string(claimsJSON),
			PasswordHash: "", // Empty password hash marks SSO user
			Role:         RoleUser,
			Balance:      defaultBalance,
			Concurrency:  defaultConcurrency,
			Status:       StatusActive,
		}

		log.Printf("[SSO] Calling userRepo.Create with user: %+v", user)

		if err := s.userRepo.Create(ctx, user); err != nil {
			log.Printf("[SSO] Failed to create user: %v", err)
			return "", nil, false, fmt.Errorf("failed to create user: %w", err)
		}

		log.Printf("[SSO] User created successfully: ID=%d", user.ID)

		isNewUser = true
		log.Printf("Created new SSO user: %s (sub: %s, trust_level: %d)", email, claims.Sub, claims.TrustLevel)
	}

	// Generate JWT token
	token, err := s.authService.GenerateToken(user)
	if err != nil {
		return "", nil, false, fmt.Errorf("failed to generate token: %w", err)
	}

	return token, user, isNewUser, nil
}

// getSSOConfig retrieves SSO configuration from settings
func (s *OIDCSSOService) getSSOConfig(ctx context.Context) (*SSOConfig, error) {
	issuerURL, _ := s.settingService.GetSetting(ctx, SettingKeySSOIssuerURL)
	clientID, _ := s.settingService.GetSetting(ctx, SettingKeySSOClientID)
	clientSecret, _ := s.settingService.GetSetting(ctx, SettingKeySSOClientSecret)
	redirectURI, _ := s.settingService.GetSetting(ctx, SettingKeySSORedirectURI)

	if issuerURL == "" || clientID == "" || clientSecret == "" || redirectURI == "" {
		return nil, infraerrors.ServiceUnavailable("SSO_NOT_CONFIGURED",
			"SSO is not properly configured")
	}

	return &SSOConfig{
		IssuerURL:    issuerURL,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
	}, nil
}

// validateEmailDomain validates if an email domain is allowed
func (s *OIDCSSOService) validateEmailDomain(ctx context.Context, email string) bool {
	allowedDomainsJSON, _ := s.settingService.GetSetting(ctx, SettingKeySSOAllowedDomains)
	if allowedDomainsJSON == "" || allowedDomainsJSON == "[]" {
		return true // No restriction
	}

	var allowedDomains []string
	if err := json.Unmarshal([]byte(allowedDomainsJSON), &allowedDomains); err != nil {
		log.Printf("Failed to parse allowed domains: %v", err)
		return true // On error, allow
	}

	if len(allowedDomains) == 0 {
		return true
	}

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	domain := strings.ToLower(parts[1])
	for _, allowed := range allowedDomains {
		if strings.ToLower(allowed) == domain {
			return true
		}
	}

	return false
}

// SSOConfig contains SSO configuration
type SSOConfig struct {
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}
