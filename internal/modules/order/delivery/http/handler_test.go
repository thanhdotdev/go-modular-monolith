package orderhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	orderapplication "project-example/internal/modules/order/application"
	orderdomain "project-example/internal/modules/order/domain"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type stubUseCase struct {
	order *orderapplication.OrderDTO
	err   error
}

func (s stubUseCase) GetOrder(_ context.Context, _ string) (*orderapplication.OrderDTO, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.order, nil
}

func TestGetOrderReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewHandler(stubUseCase{err: orderdomain.ErrOrderNotFound})
	handler.RegisterRoutes(router.Group("/api/v1"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/orders/ord-404", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), "\"code\":\"order_not_found\"") {
		t.Fatalf("expected order_not_found code in body, got %s", recorder.Body.String())
	}
}
