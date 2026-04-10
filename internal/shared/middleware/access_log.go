package middleware

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"net/http"
	"project-example/internal/platform/config"
	"time"

	applogger "project-example/internal/platform/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func AccessLog(cfg config.LoggingConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		requestLogger := applogger.FromContext(c.Request.Context())
		rawPath := c.Request.URL.Path
		requestPayload, requestPayloadTruncated := captureRequestBody(c.Request, cfg)
		responseWriter := newBodyLogWriter(c.Writer, cfg)
		if responseWriter != nil {
			c.Writer = responseWriter
		}

		requestFields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", rawPath),
			zap.String("client_ip", c.ClientIP()),
		}

		if c.Request.URL.RawQuery != "" {
			requestFields = append(requestFields, zap.String("query", c.Request.URL.RawQuery))
		}

		if c.Request.ContentLength >= 0 {
			requestFields = append(requestFields, zap.Int64("request_size_bytes", c.Request.ContentLength))
		}

		appendPayloadFields(&requestFields, "request_payload", requestPayload, requestPayloadTruncated)

		requestLogger.Info(fmt.Sprintf("[REQUEST] %s", rawPath), requestFields...)

		c.Next()

		latency := time.Since(startedAt)
		fields := []zap.Field{
			zap.String("method", c.Request.Method),
			zap.String("path", rawPath),
			zap.String("raw_path", c.Request.URL.Path),
			zap.String("route", c.FullPath()),
			zap.Int("status", c.Writer.Status()),
			zap.Float64("latency_ms", float64(latency)/float64(time.Millisecond)),
			zap.Int("response_size_bytes", c.Writer.Size()),
			zap.String("client_ip", c.ClientIP()),
		}

		if c.Request.URL.RawQuery != "" {
			fields = append(fields, zap.String("query", c.Request.URL.RawQuery))
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		if responseWriter != nil {
			appendPayloadFields(
				&fields,
				"response_payload",
				responseWriter.Body(c.Writer.Header().Get("Content-Type")),
				responseWriter.Truncated(),
			)
		}

		message := fmt.Sprintf("[RESPONSE] %s", rawPath)
		if c.Writer.Status() >= 500 {
			requestLogger.Error(message, fields...)
			return
		}

		if c.Writer.Status() >= 400 {
			requestLogger.Warn(message, fields...)
			return
		}

		requestLogger.Info(message, fields...)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	buffer *limitedBuffer
}

func newBodyLogWriter(writer gin.ResponseWriter, cfg config.LoggingConfig) *bodyLogWriter {
	if !cfg.IncludeResponseBody || cfg.BodyMaxBytes <= 0 {
		return nil
	}

	return &bodyLogWriter{
		ResponseWriter: writer,
		buffer:         newLimitedBuffer(cfg.BodyMaxBytes),
	}
}

func (w *bodyLogWriter) Write(data []byte) (int, error) {
	if w.buffer != nil {
		_, _ = w.buffer.Write(data)
	}

	return w.ResponseWriter.Write(data)
}

func (w *bodyLogWriter) WriteString(value string) (int, error) {
	if w.buffer != nil {
		_, _ = w.buffer.WriteString(value)
	}

	return w.ResponseWriter.WriteString(value)
}

func (w *bodyLogWriter) Body(contentType string) string {
	if w == nil || w.buffer == nil || !shouldLogBody(contentType) {
		return ""
	}

	return w.buffer.String()
}

func (w *bodyLogWriter) Truncated() bool {
	if w == nil || w.buffer == nil {
		return false
	}

	return w.buffer.Truncated()
}

type limitedBuffer struct {
	max       int
	buffer    bytes.Buffer
	truncated bool
}

func newLimitedBuffer(max int) *limitedBuffer {
	return &limitedBuffer{max: max}
}

func (b *limitedBuffer) Write(data []byte) (int, error) {
	if b.max <= 0 || len(data) == 0 {
		return len(data), nil
	}

	remaining := b.max - b.buffer.Len()
	if remaining <= 0 {
		b.truncated = true
		return len(data), nil
	}

	if len(data) > remaining {
		_, _ = b.buffer.Write(data[:remaining])
		b.truncated = true
		return len(data), nil
	}

	return b.buffer.Write(data)
}

func (b *limitedBuffer) WriteString(value string) (int, error) {
	return b.Write([]byte(value))
}

func (b *limitedBuffer) String() string {
	return b.buffer.String()
}

func (b *limitedBuffer) Truncated() bool {
	return b.truncated
}

func captureRequestBody(request *http.Request, cfg config.LoggingConfig) (string, bool) {
	if !cfg.IncludeRequestBody || cfg.BodyMaxBytes <= 0 || request == nil || request.Body == nil {
		return "", false
	}

	if !shouldLogBody(request.Header.Get("Content-Type")) {
		return "", false
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(request.Body, int64(cfg.BodyMaxBytes)+1))
	if err != nil {
		return "", false
	}

	request.Body = io.NopCloser(io.MultiReader(bytes.NewReader(bodyBytes), request.Body))
	if len(bodyBytes) == 0 {
		return "", false
	}

	if len(bodyBytes) > cfg.BodyMaxBytes {
		return string(bodyBytes[:cfg.BodyMaxBytes]), true
	}

	return string(bodyBytes), false
}

func appendPayloadFields(fields *[]zap.Field, key, payload string, truncated bool) {
	if payload == "" {
		return
	}

	*fields = append(*fields, zap.String(key, payload))
	if truncated {
		*fields = append(*fields, zap.Bool(key+"_truncated", true))
	}
}

func shouldLogBody(contentType string) bool {
	if contentType == "" {
		return true
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	if mediaType == "application/json" ||
		mediaType == "application/xml" ||
		mediaType == "application/x-www-form-urlencoded" {
		return true
	}

	if len(mediaType) >= len("text/") && mediaType[:len("text/")] == "text/" {
		return true
	}

	return len(mediaType) > len("+json") && mediaType[len(mediaType)-len("+json"):] == "+json"
}
