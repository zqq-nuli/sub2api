package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeadersConfig 安全响应头配置
type SecurityHeadersConfig struct {
	// XContentTypeOptions 防止浏览器进行 MIME 类型嗅探
	XContentTypeOptions string
	// XFrameOptions 防止点击劫持攻击
	XFrameOptions string
	// XXSSProtection 启用 XSS 过滤器（旧版浏览器）
	XXSSProtection string
	// StrictTransportSecurity HSTS 配置（仅 HTTPS）
	StrictTransportSecurity string
	// ContentSecurityPolicy CSP 配置
	ContentSecurityPolicy string
	// ReferrerPolicy 控制 Referrer 信息
	ReferrerPolicy string
	// PermissionsPolicy 控制浏览器功能权限
	PermissionsPolicy string
}

// DefaultSecurityHeadersConfig 返回默认安全头配置
func DefaultSecurityHeadersConfig() SecurityHeadersConfig {
	return SecurityHeadersConfig{
		XContentTypeOptions:     "nosniff",
		XFrameOptions:           "DENY",
		XXSSProtection:          "1; mode=block",
		StrictTransportSecurity: "max-age=31536000; includeSubDomains",
		ContentSecurityPolicy:   "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self' data:; connect-src 'self' https:;",
		ReferrerPolicy:          "strict-origin-when-cross-origin",
		PermissionsPolicy:       "geolocation=(), microphone=(), camera=()",
	}
}

// SecurityHeaders 安全响应头中间件（使用默认配置）
func SecurityHeaders() gin.HandlerFunc {
	return SecurityHeadersWithConfig(DefaultSecurityHeadersConfig())
}

// SecurityHeadersWithConfig 使用自定义配置的安全响应头中间件
func SecurityHeadersWithConfig(config SecurityHeadersConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止 MIME 类型嗅探
		if config.XContentTypeOptions != "" {
			c.Header("X-Content-Type-Options", config.XContentTypeOptions)
		}

		// 防止点击劫持
		if config.XFrameOptions != "" {
			c.Header("X-Frame-Options", config.XFrameOptions)
		}

		// XSS 保护（旧版浏览器）
		if config.XXSSProtection != "" {
			c.Header("X-XSS-Protection", config.XXSSProtection)
		}

		// HSTS - 仅在 HTTPS 连接时设置
		if config.StrictTransportSecurity != "" {
			// 检查是否为 HTTPS 请求（考虑反向代理的情况）
			isHTTPS := c.Request.TLS != nil ||
				c.GetHeader("X-Forwarded-Proto") == "https" ||
				c.GetHeader("X-Forwarded-Ssl") == "on"
			if isHTTPS {
				c.Header("Strict-Transport-Security", config.StrictTransportSecurity)
			}
		}

		// 内容安全策略
		if config.ContentSecurityPolicy != "" {
			c.Header("Content-Security-Policy", config.ContentSecurityPolicy)
		}

		// Referrer 策略
		if config.ReferrerPolicy != "" {
			c.Header("Referrer-Policy", config.ReferrerPolicy)
		}

		// 权限策略
		if config.PermissionsPolicy != "" {
			c.Header("Permissions-Policy", config.PermissionsPolicy)
		}

		c.Next()
	}
}

// APISecurityHeaders 专门为 API 端点设计的安全头
// 比默认配置更严格，适用于纯 API 服务
func APISecurityHeaders() gin.HandlerFunc {
	return SecurityHeadersWithConfig(SecurityHeadersConfig{
		XContentTypeOptions:     "nosniff",
		XFrameOptions:           "DENY",
		XXSSProtection:          "0", // 现代浏览器建议禁用，CSP 更有效
		StrictTransportSecurity: "max-age=63072000; includeSubDomains; preload",
		ContentSecurityPolicy:   "default-src 'none'; frame-ancestors 'none'",
		ReferrerPolicy:          "no-referrer",
		PermissionsPolicy:       "geolocation=(), microphone=(), camera=(), payment=()",
	})
}
