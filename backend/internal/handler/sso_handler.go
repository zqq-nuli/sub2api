package handler

import (
	"github.com/Wei-Shaw/sub2api/internal/handler/dto"
	"github.com/Wei-Shaw/sub2api/internal/pkg/response"
	"github.com/Wei-Shaw/sub2api/internal/service"

	"github.com/gin-gonic/gin"
)

// SSOHandler handles SSO authentication requests
type SSOHandler struct {
	ssoService     *service.OIDCSSOService
	settingService *service.SettingService
}

// NewSSOHandler creates a new SSO handler
func NewSSOHandler(ssoService *service.OIDCSSOService, settingService *service.SettingService) *SSOHandler {
	return &SSOHandler{
		ssoService:     ssoService,
		settingService: settingService,
	}
}

// SSOConfigResponse represents the SSO configuration for clients
type SSOConfigResponse struct {
	SSOEnabled           bool `json:"sso_enabled"`
	PasswordLoginEnabled bool `json:"password_login_enabled"`
}

// GetSSOConfig returns SSO configuration for the login page
// GET /api/v1/auth/sso/config
func (h *SSOHandler) GetSSOConfig(c *gin.Context) {
	ssoEnabled, _ := h.settingService.GetBoolSetting(c.Request.Context(), service.SettingKeySSOEnabled)
	passwordLoginEnabled, _ := h.settingService.GetBoolSetting(c.Request.Context(), service.SettingKeyPasswordLoginEnabled)

	response.Success(c, &SSOConfigResponse{
		SSOEnabled:           ssoEnabled,
		PasswordLoginEnabled: passwordLoginEnabled,
	})
}

// GenerateAuthURL generates an SSO authorization URL
// GET /api/v1/auth/sso/authorize
func (h *SSOHandler) GenerateAuthURL(c *gin.Context) {
	// Check if SSO is enabled
	ssoEnabled, _ := h.settingService.GetBoolSetting(c.Request.Context(), service.SettingKeySSOEnabled)
	if !ssoEnabled {
		response.Forbidden(c, "SSO login is disabled")
		return
	}

	result, err := h.ssoService.GenerateAuthURL(c.Request.Context())
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	response.Success(c, result)
}

// SSOCallbackRequest represents the SSO callback request
type SSOCallbackRequest struct {
	Code      string `form:"code" binding:"required"`
	State     string `form:"state" binding:"required"`
	SessionID string `form:"session_id" binding:"required"`
}

// SSOCallbackResponse represents the SSO callback response
type SSOCallbackResponse struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	User        *dto.User `json:"user"`
	IsNewUser   bool      `json:"is_new_user"`
}

// Callback handles the SSO callback
// GET /api/v1/auth/sso/callback?code=xxx&state=xxx&session_id=xxx
func (h *SSOHandler) Callback(c *gin.Context) {
	var req SSOCallbackRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.BadRequest(c, "Invalid callback parameters: "+err.Error())
		return
	}

	// Exchange code and create/login user
	token, user, isNewUser, err := h.ssoService.ExchangeCodeAndCreateUser(
		c.Request.Context(),
		req.Code,
		req.State,
		req.SessionID,
	)
	if err != nil {
		response.ErrorFrom(c, err)
		return
	}

	// Convert to DTO
	userDTO := dto.UserFromService(user)

	response.Success(c, &SSOCallbackResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		User:        userDTO,
		IsNewUser:   isNewUser,
	})
}
