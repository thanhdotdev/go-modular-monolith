package httpx

import "github.com/gin-gonic/gin"

type Meta struct {
	RequestID string `json:"requestId,omitempty"`
}

type SuccessResponse struct {
	Data any  `json:"data"`
	Meta Meta `json:"meta,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
	Meta  Meta      `json:"meta,omitempty"`
}

func OK(c *gin.Context, data any) {
	JSON(c, 200, data)
}

func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, SuccessResponse{
		Data: data,
		Meta: Meta{
			RequestID: RequestID(c),
		},
	})
}

func Error(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{
		Error: ErrorBody{
			Code:    code,
			Message: message,
		},
		Meta: Meta{
			RequestID: RequestID(c),
		},
	})
}
