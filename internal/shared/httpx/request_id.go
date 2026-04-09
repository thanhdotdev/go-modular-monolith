package httpx

import "github.com/gin-gonic/gin"

const (
	HeaderRequestID = "X-Request-ID"
	contextKey      = "request_id"
)

func SetRequestID(c *gin.Context, requestID string) {
	c.Set(contextKey, requestID)
	c.Writer.Header().Set(HeaderRequestID, requestID)
}

func RequestID(c *gin.Context) string {
	requestID, ok := contextValue[string](c, contextKey)
	if ok {
		return requestID
	}

	return c.Writer.Header().Get(HeaderRequestID)
}

func contextValue[T any](c *gin.Context, key string) (T, bool) {
	var zero T

	value, ok := c.Get(key)
	if !ok {
		return zero, false
	}

	typed, ok := value.(T)
	if !ok {
		return zero, false
	}

	return typed, true
}
