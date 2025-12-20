package templates

// MiddlewareCORSTemplate generates the cors.go file
func MiddlewareCORSTemplate() string {
	return `package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ncobase/ncore/config"
)

// CORSConfig defines the config for CORS middleware.
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

// GetCORSConfig returns the CORS configuration based on the environment
func GetCORSConfig(conf *config.Config) CORSConfig {
	corsConfig := CORSConfig{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With", "X-Session-ID", "X-API-Key"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		MaxAge:           12 * 60 * 60, // 12 hours
	}

	if conf.IsProd() {
		corsConfig.AllowOrigins = []string{}
		if conf.Domain == "" {
			// Add your production domains here
			// corsConfig.AllowOrigins = []string{"https://example.com"}
		}
	} else {
		corsConfig.AllowOrigins = []string{
			"http://localhost:3000",
			"http://localhost:8080",
			"http://127.0.0.1:3000",
			"http://127.0.0.1:8080",
		}
	}

	return corsConfig
}

// CORSHandler is a middleware for handling CORS.
func CORSHandler(conf *config.Config) gin.HandlerFunc {
	corsConfig := GetCORSConfig(conf)

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if isOriginAllowed(origin, corsConfig.AllowOrigins) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join(corsConfig.AllowMethods, ","))
		c.Writer.Header().Set("Access-Control-Allow-Headers", strings.Join(corsConfig.AllowHeaders, ","))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsConfig.AllowCredentials))
		c.Writer.Header().Set("Access-Control-Expose-Headers", strings.Join(corsConfig.ExposeHeaders, ","))
		c.Writer.Header().Set("Access-Control-Max-Age", strconv.Itoa(corsConfig.MaxAge))

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isOriginAllowed checks if the origin is allowed in the list of allowed origins
func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
		if strings.HasSuffix(allowed, ":*") {
			prefix := strings.TrimSuffix(allowed, ":*")
			if strings.HasPrefix(origin, prefix) {
				return true
			}
		}
	}
	return false
}
`
}

// MiddlewareSecurityHeadersTemplate generates the security_headers.go file
func MiddlewareSecurityHeadersTemplate() string {
	return `package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// SecurityHeadersConfig holds security headers configuration
type SecurityHeadersConfig struct {
	// ContentSecurityPolicy sets CSP header
	ContentSecurityPolicy string ` + "`" + `json:"content_security_policy"` + "`" + `

	// StrictTransportSecurity sets HSTS header
	StrictTransportSecurity string ` + "`" + `json:"strict_transport_security"` + "`" + `

	// XFrameOptions sets X-Frame-Options header
	XFrameOptions string ` + "`" + `json:"x_frame_options"` + "`" + `

	// XContentTypeOptions sets X-Content-Type-Options header
	XContentTypeOptions string ` + "`" + `json:"x_content_type_options"` + "`" + `

	// XSSProtection sets X-XSS-Protection header
	XSSProtection string ` + "`" + `json:"xss_protection"` + "`" + `

	// ReferrerPolicy sets Referrer-Policy header
	ReferrerPolicy string ` + "`" + `json:"referrer_policy"` + "`" + `

	// PermissionsPolicy sets Permissions-Policy header
	PermissionsPolicy string ` + "`" + `json:"permissions_policy"` + "`" + `

	// CrossOriginEmbedderPolicy sets Cross-Origin-Embedder-Policy header
	CrossOriginEmbedderPolicy string ` + "`" + `json:"cross_origin_embedder_policy"` + "`" + `

	// CrossOriginResourcePolicy sets Cross-Origin-Resource-Policy header
	CrossOriginResourcePolicy string ` + "`" + `json:"cross_origin_resource_policy"` + "`" + `

	// CrossOriginOpenerPolicy sets Cross-Origin-Opener-Policy header
	CrossOriginOpenerPolicy string ` + "`" + `json:"cross_origin_opener_policy"` + "`" + `

	// Custom headers to add
	CustomHeaders map[string]string ` + "`" + `json:"custom_headers"` + "`" + `

	// Remove these headers from response
	RemoveHeaders []string ` + "`" + `json:"remove_headers"` + "`" + `
}

// DefaultSecurityHeadersConfig returns default security headers configuration
func DefaultSecurityHeadersConfig() *SecurityHeadersConfig {
	return &SecurityHeadersConfig{
		ContentSecurityPolicy: strings.Join([]string{
			"default-src 'self'",
			"script-src 'self' 'unsafe-inline' 'unsafe-eval'",
			"style-src 'self' 'unsafe-inline'",
			"img-src 'self' data: https:",
			"font-src 'self' data:",
			"connect-src 'self'",
			"frame-ancestors 'none'",
			"base-uri 'self'",
			"form-action 'self'",
			"object-src 'none'",
		}, "; "),
		StrictTransportSecurity:   "max-age=31536000; includeSubDomains; preload",
		XFrameOptions:             "DENY",
		XContentTypeOptions:       "nosniff",
		XSSProtection:             "1; mode=block",
		ReferrerPolicy:            "strict-origin-when-cross-origin",
		PermissionsPolicy:         "camera=(), microphone=(), geolocation=()",
		CrossOriginEmbedderPolicy: "require-corp",
		CrossOriginResourcePolicy: "cross-origin",
		CrossOriginOpenerPolicy:   "same-origin",
		CustomHeaders: map[string]string{
			"X-Powered-By":                      "", // Remove this header
			"Server":                            "", // Remove server information
			"X-Content-Duration":                "",
			"X-Robots-Tag":                      "noindex, nofollow, nosnippet, noarchive",
			"Cache-Control":                     "no-cache, no-store, must-revalidate, private",
			"Pragma":                            "no-cache",
			"Expires":                           "0",
			"X-Download-Options":                "noopen",
			"X-Permitted-Cross-Domain-Policies": "none",
		},
		RemoveHeaders: []string{
			"X-Powered-By",
			"Server",
			"X-AspNet-Version",
			"X-AspNetMvc-Version",
		},
	}
}

// SecurityHeaders provides security headers middleware
func SecurityHeaders(config *SecurityHeadersConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultSecurityHeadersConfig()
	}

	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Remove unwanted headers
		for _, header := range config.RemoveHeaders {
			c.Writer.Header().Del(header)
		}

		// Set security headers
		if config.ContentSecurityPolicy != "" {
			c.Writer.Header().Set("Content-Security-Policy", config.ContentSecurityPolicy)
		}

		if config.StrictTransportSecurity != "" {
			c.Writer.Header().Set("Strict-Transport-Security", config.StrictTransportSecurity)
		}

		if config.XFrameOptions != "" {
			c.Writer.Header().Set("X-Frame-Options", config.XFrameOptions)
		}

		if config.XContentTypeOptions != "" {
			c.Writer.Header().Set("X-Content-Type-Options", config.XContentTypeOptions)
		}

		if config.XSSProtection != "" {
			c.Writer.Header().Set("X-XSS-Protection", config.XSSProtection)
		}

		if config.ReferrerPolicy != "" {
			c.Writer.Header().Set("Referrer-Policy", config.ReferrerPolicy)
		}

		if config.PermissionsPolicy != "" {
			c.Writer.Header().Set("Permissions-Policy", config.PermissionsPolicy)
		}

		if config.CrossOriginEmbedderPolicy != "" {
			c.Writer.Header().Set("Cross-Origin-Embedder-Policy", config.CrossOriginEmbedderPolicy)
		}

		if config.CrossOriginResourcePolicy != "" {
			c.Writer.Header().Set("Cross-Origin-Resource-Policy", config.CrossOriginResourcePolicy)
		}

		if config.CrossOriginOpenerPolicy != "" {
			c.Writer.Header().Set("Cross-Origin-Opener-Policy", config.CrossOriginOpenerPolicy)
		}

		// Set custom headers
		for name, value := range config.CustomHeaders {
			if value == "" {
				c.Writer.Header().Del(name)
			} else {
				c.Writer.Header().Set(name, value)
			}
		}
	}
}

// InputSanitizationConfig holds input sanitization configuration
func DefaultSanitizationConfig() *struct{} {
	return &struct{}{}
}

// InputSanitization middleware placeholder
func InputSanitization(config *struct{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// InputValidationConfig holds input validation configuration
func DefaultValidationConfig() *struct{} {
	return &struct{}{}
}

// InputValidation middleware placeholder
func InputValidation(config *struct{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
`
}

// MiddlewareTraceTemplate generates the trace.go file
func MiddlewareTraceTemplate() string {
	return `package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/ncobase/ncore/consts"
	"github.com/ncobase/ncore/ctxutil"
	"github.com/ncobase/ncore/logging/observes"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// Trace middleware for request tracing and context setup
func Trace(c *gin.Context) {
	// Create context with Gin context embedded
	ctx := ctxutil.WithGinContext(c.Request.Context(), c)

	// Handle trace ID
	traceID := c.GetHeader(consts.TraceKey)
	if traceID == "" {
		ctx, traceID = ctxutil.EnsureTraceID(ctx)
	} else {
		ctx = ctxutil.SetTraceID(ctx, traceID)
	}

	// Set client information in context
	ctx = setClientInfoToContext(ctx, c)

	// Update request context - this is crucial!
	c.Request = c.Request.WithContext(ctx)

	// Set trace ID in Gin's context for easy access
	c.Set(ctxutil.TraceIDKey, traceID)

	// Set trace header in response
	c.Writer.Header().Set(consts.TraceKey, traceID)

	// Create OpenTelemetry tracing context
	path := c.Request.URL.Path
	if path == "" {
		path = c.FullPath()
	}
	tc := observes.NewTracingContext(ctx, path, 100)
	defer tc.End()

	tc.SetAttributes(
		attribute.String("http.method", c.Request.Method),
		attribute.String("http.path", path),
		attribute.String("trace.id", traceID),
		attribute.String("client.ip", c.ClientIP()),
		attribute.String("user.agent", c.GetHeader("User-Agent")),
	)

	ctx = context.WithValue(tc.Context(), "tracing_context", tc)
	c.Request = c.Request.WithContext(ctx)

	c.Next()

	// Update span with response status
	status := c.Writer.Status()
	tc.SetAttributes(
		attribute.Int("http.status_code", status),
	)
	tc.SetStatus(codes.Code(status), http.StatusText(status))
}

// setClientInfoToContext sets client information to context
func setClientInfoToContext(ctx context.Context, c *gin.Context) context.Context {
	// Extract client information
	ip := extractClientIP(c)
	userAgent := c.GetHeader("User-Agent")
	referer := c.GetHeader("Referer")
	sessionID := extractSessionID(c)

	// Set to context
	ctx = ctxutil.SetClientIP(ctx, ip)
	ctx = ctxutil.SetUserAgent(ctx, userAgent)
	ctx = ctxutil.SetSessionID(ctx, sessionID)
	ctx = ctxutil.SetReferer(ctx, referer)

	// Also set HTTP request for direct access
	ctx = ctxutil.SetHTTPRequest(ctx, c.Request)

	return ctx
}

// extractClientIP extracts real client IP from various headers
func extractClientIP(c *gin.Context) string {
	// Priority order for IP extraction
	headers := []string{
		"X-Forwarded-For",
		"X-Real-IP",
		"CF-Connecting-IP",
		"X-Client-IP",
		"X-Cluster-Client-IP",
	}

	for _, header := range headers {
		ip := c.GetHeader(header)
		if ip != "" && ip != "unknown" {
			// Handle X-Forwarded-For which may contain multiple IPs
			if header == "X-Forwarded-For" {
				ips := strings.Split(ip, ",")
				if len(ips) > 0 {
					cleanIP := strings.TrimSpace(ips[0])
					if cleanIP != "" && cleanIP != "unknown" {
						return cleanIP
					}
				}
			} else {
				return ip
			}
		}
	}

	// Fallback to Gin's ClientIP method
	return c.ClientIP()
}

// OtelTrace middleware for OpenTelemetry trace
func OtelTrace(c *gin.Context) {
	// Reuse Trace logic or extend as needed
	// For now, it's just a placeholder or could be same as Trace
	c.Next()
}
`
}

// MiddlewareLoggerTemplate generates the logger.go file
func MiddlewareLoggerTemplate() string {
	return `package middleware

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ncobase/ncore/ctxutil"
	"github.com/ncobase/ncore/logging/logger"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// responseWriter wraps the original responseWriter to capture response data
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write writes the data to the buffer
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

var (
	// skippedPaths is a list of paths that should be skipped for detailed logging
	skippedPaths = []string{
		"*swagger*",
		"*attachments/*",
	}

	// binaryTypes is a list of content types that should be treated as binary
	binaryTypes = []string{
		"application/octet-stream",
		"application/pdf",
		"image/",
		"audio/",
		"video/",
	}

	// Use a sync.Pool to reduce allocations
	bufferPool = sync.Pool{
		New: func() any {
			return new(bytes.Buffer)
		},
	}
)

// Logger is a middleware for logging
func Logger(c *gin.Context) {
	start := time.Now()
	ctx := ctxutil.FromGinContext(c)

	// Check if the path should be skipped
	if shouldSkipPath(c.Request, skippedPaths) {
		c.Next()
		return
	}

	// Capture request body
	var requestBody any
	if c.Request.Body != nil {
		// Skip multipart forms to avoid reading large files into memory and breaking ParseMultipartForm
		if !strings.HasPrefix(c.ContentType(), "multipart/") {
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err != nil {
				logger.Errorf(ctx, "Failed to read request body: %v", err)
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				requestBody = processBody(bodyBytes, c.ContentType(), c.Request.URL.Path)
			}
		}
	}

	// Wrap response writer
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufferPool.Put(buf)

	w := &responseWriter{body: buf, ResponseWriter: c.Writer}
	c.Writer = w

	c.Next()

	// Prepare log entry
	entry := logrus.Fields{
		"status":     c.Writer.Status(),
		"method":     c.Request.Method,
		"path":       c.Request.URL.Path,
		"query":      c.Request.URL.RawQuery,
		"ip":         c.ClientIP(),
		"latency":    time.Since(start),
		"user_agent": c.Request.UserAgent(),
	}

	if requestBody != nil {
		entry["request_body"] = requestBody
	}

	responseBody := processBody(w.body.Bytes(), w.Header().Get("Content-Type"), c.Request.URL.Path)
	if responseBody != nil {
		entry["response_body"] = responseBody
	}

	if len(c.Errors) > 0 {
		entry["error"] = c.Errors.String()
	}

	// Log request
	l := logger.WithFields(ctx, entry)
	switch {
	case c.Writer.Status() >= http.StatusInternalServerError:
		l.Error("Oops! Something went wrong on our end")
	case c.Writer.Status() >= http.StatusBadRequest && c.Writer.Status() < http.StatusInternalServerError:
		l.Warn("Request couldn't be processed")
	case c.Writer.Status() >= http.StatusMultipleChoices && c.Writer.Status() < http.StatusBadRequest:
		l.Info("Redirecting request")
	case c.Writer.Status() >= http.StatusOK && c.Writer.Status() < http.StatusMultipleChoices:
		l.Info("Request processed successfully")
	default:
		l.Info("Request completed with status: " + strconv.Itoa(c.Writer.Status()))
	}
}

// processBody processes the body of the request or response
func processBody(body []byte, contentType, _ string) any {
	if len(body) == 0 {
		return nil
	}

	if isBinaryContentType(contentType) {
		return base64.StdEncoding.EncodeToString(body)
	}

	var jsonBody any
	if json.Valid(body) {
		if err := json.Unmarshal(body, &jsonBody); err != nil {
			return string(body)
		}
		return jsonBody
	}

	return string(body)
}

// isBinaryContentType checks if the content type is a binary type
func isBinaryContentType(contentType string) bool {
	contentType = strings.ToLower(strings.TrimSpace(strings.SplitN(contentType, ";", 2)[0]))
	for _, t := range binaryTypes {
		if strings.HasPrefix(contentType, t) {
			return true
		}
	}
	return false
}

// shouldSkipPath checks if the path should be skipped
func shouldSkipPath(r *http.Request, skippedPaths []string) bool {
	// Implementation simplified
	return false
}
`
}

// MiddlewareClientInfoTemplate generates the client_info.go file
func MiddlewareClientInfoTemplate() string {
	return `package middleware

import (
	"net"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ncobase/ncore/ctxutil"
)

// ClientInfo middleware extracts and sets client information to context
func ClientInfo(c *gin.Context) {
	// Get context
	ctx := c.Request.Context()
	if _, ok := ctxutil.GetGinContext(ctx); !ok {
		ctx = ctxutil.WithGinContext(ctx, c)
	}

	// Extract and set client information
	ctx = ctxutil.SetClientInfo(ctx,
		extractRealClientIP(c),
		extractUserAgent(c),
		extractSessionID(c),
	)

	// Set HTTP request for direct access
	ctx = ctxutil.SetHTTPRequest(ctx, c.Request)

	// Update request context
	c.Request = c.Request.WithContext(ctx)

	c.Next()
}

// extractRealClientIP extracts real client IP with comprehensive proxy support
func extractRealClientIP(c *gin.Context) string {
	// Check forwarded headers in order of priority
	forwardedHeaders := []struct {
		name     string
		multiple bool // whether header can contain multiple IPs
	}{
		{"X-Forwarded-For", true},
		{"X-Real-IP", false},
		{"CF-Connecting-IP", false},
		{"X-Client-IP", false},
		{"X-Cluster-Client-IP", false},
		{"Forwarded-For", false},
		{"Forwarded", false},
	}

	for _, header := range forwardedHeaders {
		value := c.GetHeader(header.name)
		if value == "" || value == "unknown" {
			continue
		}

		if header.multiple {
			// Handle comma-separated IPs (X-Forwarded-For)
			ips := strings.Split(value, ",")
			for _, ip := range ips {
				cleanIP := strings.TrimSpace(ip)
				if isValidPublicIP(cleanIP) {
					return cleanIP
				}
			}
		} else {
			if isValidPublicIP(value) {
				return value
			}
		}
	}

	// Fallback to Gin's ClientIP method
	if clientIP := c.ClientIP(); clientIP != "" {
		return clientIP
	}

	// Last resort: extract from RemoteAddr
	if c.Request != nil && c.Request.RemoteAddr != "" {
		if host, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
			return host
		}
		return c.Request.RemoteAddr
	}

	return "unknown"
}

// extractUserAgent extracts User-Agent header
func extractUserAgent(c *gin.Context) string {
	if ua := c.GetHeader("User-Agent"); ua != "" {
		return ua
	}
	return "unknown"
}

// extractSessionID extracts session ID from various sources
func extractSessionID(c *gin.Context) string {
	// Try session cookie
	sessionCookieNames := []string{"session_id", "sessionid", "SESSIONID"}
	for _, cookieName := range sessionCookieNames {
		if sessionID, err := c.Cookie(cookieName); err == nil && sessionID != "" {
			return sessionID
		}
	}

	// Try session headers
	sessionHeaders := []string{"X-Session-ID", "X-Session-Id", "Session-ID"}
	for _, headerName := range sessionHeaders {
		if sessionID := c.GetHeader(headerName); sessionID != "" {
			return sessionID
		}
	}

	return ""
}

// isValidPublicIP checks if IP is valid and not private/reserved
func isValidPublicIP(ipStr string) bool {
	if ipStr == "" || ipStr == "unknown" {
		return false
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	// Check for private/reserved IP ranges
	privateRanges := []string{
		"10.0.0.0/8",     // RFC1918
		"172.16.0.0/12",  // RFC1918
		"192.168.0.0/16", // RFC1918
		"127.0.0.0/8",    // Loopback
		"169.254.0.0/16", // Link Local
		"224.0.0.0/4",    // Multicast
		"240.0.0.0/4",    // Reserved
		"::1/128",        // IPv6 loopback
		"fc00::/7",       // IPv6 unique local
		"fe80::/10",      // IPv6 link local
		"ff00::/8",       // IPv6 multicast
	}

	for _, rangeStr := range privateRanges {
		_, subnet, err := net.ParseCIDR(rangeStr)
		if err != nil {
			continue
		}
		if subnet.Contains(ip) {
			return false
		}
	}

	return true
}
`
}
