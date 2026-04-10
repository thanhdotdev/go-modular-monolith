package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorMapping struct {
	Err     error
	Status  int
	Code    string
	Message string
}

func WriteError(c *gin.Context, err error, mappings ...ErrorMapping) {
	if err == nil {
		return
	}

	for _, mapping := range mappings {
		if errors.Is(err, mapping.Err) {
			message := mapping.Message
			if message == "" {
				message = err.Error()
			}

			Error(c, mapping.Status, mapping.Code, message)
			return
		}
	}

	_ = c.Error(err)
	Error(c, http.StatusInternalServerError, "internal_error", "internal server error")
}
