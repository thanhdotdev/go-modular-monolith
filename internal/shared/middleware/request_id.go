package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	applogger "project-example/internal/platform/logger"
	"project-example/internal/shared/httpx"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader(httpx.HeaderRequestID)
		if requestID == "" {
			requestID = newRequestID()
		}

		httpx.SetRequestID(c, requestID)
		c.Request = c.Request.WithContext(applogger.WithRequestID(c.Request.Context(), requestID))
		c.Next()
	}
}

func newRequestID() string {
	buffer := make([]byte, 12)
	if _, err := rand.Read(buffer); err == nil {
		return hex.EncodeToString(buffer)
	}

	return time.Now().UTC().Format("20060102150405.000000000")
}
