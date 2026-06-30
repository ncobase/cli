package templates

func MiddlewareInputTemplate() string {
	return `package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// InputSanitizationConfig holds input sanitization settings.
type InputSanitizationConfig struct {
	BlockedSubstrings []string
}

// DefaultSanitizationConfig returns conservative input sanitization defaults.
func DefaultSanitizationConfig() *InputSanitizationConfig {
	return &InputSanitizationConfig{
		BlockedSubstrings: []string{"\x00", "<script", "javascript:"},
	}
}

// InputSanitization rejects common dangerous request query markers.
func InputSanitization(config *InputSanitizationConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultSanitizationConfig()
	}
	return func(c *gin.Context) {
		query := strings.ToLower(c.Request.URL.RawQuery)
		for _, marker := range config.BlockedSubstrings {
			if marker != "" && strings.Contains(query, strings.ToLower(marker)) {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "request query contains blocked content"})
				return
			}
		}
		c.Next()
	}
}

// InputValidationConfig holds request input limits.
type InputValidationConfig struct {
	MaxBodyBytes        int64
	AllowedContentTypes []string
}

// DefaultValidationConfig returns production-safe request validation defaults.
func DefaultValidationConfig() *InputValidationConfig {
	return &InputValidationConfig{
		MaxBodyBytes:        10 << 20,
		AllowedContentTypes: []string{"", "application/json", "multipart/form-data", "application/x-www-form-urlencoded"},
	}
}

// InputValidation enforces request size and content type boundaries.
func InputValidation(config *InputValidationConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultValidationConfig()
	}
	return func(c *gin.Context) {
		if config.MaxBodyBytes > 0 {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.MaxBodyBytes)
			c.Writer.Header().Set("X-Max-Body-Bytes", strconv.FormatInt(config.MaxBodyBytes, 10))
		}
		contentType := strings.ToLower(strings.TrimSpace(strings.SplitN(c.ContentType(), ";", 2)[0]))
		for _, expected := range config.AllowedContentTypes {
			if contentType == expected {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusUnsupportedMediaType, gin.H{"error": "unsupported content type"})
	}
}
`
}
