package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// CORSConfig CORS 配置
type CORSConfig struct {
	// AllowedOrigins 允许的源列表，如果包含 "*" 则允许所有源（但不支持凭证）
	AllowedOrigins []string
	// AllowCredentials 是否允许携带凭证（cookies、authorization headers）
	AllowCredentials bool
}

// DefaultCORSConfig 返回默认 CORS 配置
// 生产环境应通过配置文件指定具体的允许源
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: false, // 使用通配符时不支持凭证
	}
}

// CORS 跨域中间件（使用默认配置）
func CORS() gin.HandlerFunc {
	return CORSWithConfig(DefaultCORSConfig())
}

// CORSWithConfig 使用自定义配置的 CORS 中间件
func CORSWithConfig(config CORSConfig) gin.HandlerFunc {
	// 构建允许的源集合用于快速查找
	allowedOriginsSet := make(map[string]bool)
	allowAll := false
	for _, origin := range config.AllowedOrigins {
		if origin == "*" {
			allowAll = true
		} else {
			allowedOriginsSet[strings.ToLower(origin)] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 确定是否允许该源
		var allowedOrigin string
		if allowAll {
			// 允许所有源时，如果请求包含 Origin，返回该 Origin（更安全）
			// 如果不包含凭证，可以使用 *
			if config.AllowCredentials && origin != "" {
				// 注意：Allow-Credentials: true 与 Allow-Origin: * 不兼容
				// 此时应返回具体的 origin（如果允许）或拒绝
				allowedOrigin = origin
			} else if origin != "" {
				allowedOrigin = origin
			} else {
				allowedOrigin = "*"
			}
		} else if origin != "" && allowedOriginsSet[strings.ToLower(origin)] {
			allowedOrigin = origin
		}

		// 只有当源被允许时才设置 CORS 头
		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)

			// 只有在非通配符模式下才设置 Credentials
			if config.AllowCredentials && allowedOrigin != "*" {
				c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-API-Key, x-api-key, x-goog-api-key, anthropic-version, anthropic-beta")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
			c.Writer.Header().Set("Access-Control-Max-Age", "86400") // 预检请求缓存 24 小时
		}

		// 处理预检请求
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
