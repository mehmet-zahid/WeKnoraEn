package middleware

import (
	"bytes"
	"context"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/Tencent/WeKnora/internal/logger"
	"github.com/Tencent/WeKnora/internal/types"
	secutils "github.com/Tencent/WeKnora/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	maxBodySize = 1024 * 10 // Maximum 10KB body content for logging
)

// loggerResponseBodyWriter custom ResponseWriter to capture response content (for logger middleware)
type loggerResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write overrides Write method, writes to both buffer and original writer
func (r loggerResponseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// sanitizeBody removes sensitive information
func sanitizeBody(body string) string {
	result := body
	// Replace common sensitive fields (JSON format)
	sensitivePatterns := []struct {
		pattern     string
		replacement string
	}{
		{`"password"\s*:\s*"[^"]*"`, `"password":"***"`},
		{`"token"\s*:\s*"[^"]*"`, `"token":"***"`},
		{`"access_token"\s*:\s*"[^"]*"`, `"access_token":"***"`},
		{`"refresh_token"\s*:\s*"[^"]*"`, `"refresh_token":"***"`},
		{`"authorization"\s*:\s*"[^"]*"`, `"authorization":"***"`},
		{`"api_key"\s*:\s*"[^"]*"`, `"api_key":"***"`},
		{`"secret"\s*:\s*"[^"]*"`, `"secret":"***"`},
		{`"apikey"\s*:\s*"[^"]*"`, `"apikey":"***"`},
		{`"apisecret"\s*:\s*"[^"]*"`, `"apisecret":"***"`},
	}

	for _, p := range sensitivePatterns {
		re := regexp.MustCompile(p.pattern)
		result = re.ReplaceAllString(result, p.replacement)
	}

	return result
}

// readRequestBody reads request body (limited size for logging, but full read for reset)
func readRequestBody(c *gin.Context) string {
	if c.Request.Body == nil {
		return ""
	}

	// Check Content-Type, only log JSON types
	contentType := c.GetHeader("Content-Type")
	if !strings.Contains(contentType, "application/json") &&
		!strings.Contains(contentType, "application/x-www-form-urlencoded") &&
		!strings.Contains(contentType, "text/") {
		return "[Non-text type, skipped]"
	}

	// Read full body content (no size limit), because we need to fully reset for subsequent handlers
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return "[Failed to read request body]"
	}

	// Reset request body, using full content, ensuring subsequent handlers can read complete data
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Body for logging (limited size)
	var logBodyBytes []byte
	if len(bodyBytes) > maxBodySize {
		logBodyBytes = bodyBytes[:maxBodySize]
	} else {
		logBodyBytes = bodyBytes
	}

	bodyStr := string(logBodyBytes)
	if len(bodyBytes) > maxBodySize {
		bodyStr += "... [Content too long, truncated]"
	}

	return sanitizeBody(bodyStr)
}

// RequestID middleware adds a unique request ID to the context
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get request ID from header or generate a new one
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		safeRequestID := secutils.SanitizeForLog(requestID)
		// Set request ID in header
		c.Header("X-Request-ID", requestID)

		// Set request ID in context
		c.Set(types.RequestIDContextKey.String(), requestID)

		// Set logger in context
		requestLogger := logger.GetLogger(c)
		requestLogger = requestLogger.WithField("request_id", safeRequestID)
		c.Set(types.LoggerContextKey.String(), requestLogger)

		// Set request ID in the global context for logging
		c.Request = c.Request.WithContext(
			context.WithValue(
				context.WithValue(c.Request.Context(), types.RequestIDContextKey, requestID),
				types.LoggerContextKey, requestLogger,
			),
		)

		c.Next()
	}
}

// Logger middleware logs request details with request ID, input and output
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Read request body (before Next, because Next will consume body)
		var requestBody string
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
			requestBody = readRequestBody(c)
		}

		// Create response body capturer
		responseBody := &bytes.Buffer{}
		responseWriter := &loggerResponseBodyWriter{
			ResponseWriter: c.Writer,
			body:           responseBody,
		}
		c.Writer = responseWriter

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get client IP and status code
		clientIP := c.ClientIP()
		statusCode := c.Writer.Status()
		method := c.Request.Method

		if raw != "" {
			path = path + "?" + raw
		}

		// Read response body
		responseBodyStr := ""
		if responseBody.Len() > 0 {
			// Check Content-Type, only log JSON types
			contentType := c.Writer.Header().Get("Content-Type")
			if strings.Contains(contentType, "application/json") ||
				strings.Contains(contentType, "text/") {
				bodyBytes := responseBody.Bytes()
				if len(bodyBytes) > maxBodySize {
					responseBodyStr = string(bodyBytes[:maxBodySize]) + "... [Content too long, truncated]"
				} else {
					responseBodyStr = string(bodyBytes)
				}
				responseBodyStr = sanitizeBody(responseBodyStr)
			} else {
				responseBodyStr = "[Non-text type, skipped]"
			}
		}

		// Build log message
		logMsg := logger.GetLogger(c)
		logMsg = logMsg.WithFields(map[string]interface{}{
			"method":      method,
			"path":        secutils.SanitizeForLog(path),
			"status_code": statusCode,
			"size":        c.Writer.Size(),
			"latency":     latency.String(),
			"client_ip":   secutils.SanitizeForLog(clientIP),
		})

		// Add request body (if any)
		if requestBody != "" {
			logMsg = logMsg.WithField("request_body", secutils.SanitizeForLog(requestBody))
		}

		// Add response body (if any)
		if responseBodyStr != "" {
			logMsg = logMsg.WithField("response_body", secutils.SanitizeForLog(responseBodyStr))
		}
		logMsg.Info()
	}
}
