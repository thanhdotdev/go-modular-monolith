package discounthttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	discountapplication "project-example/internal/modules/discount/application"
	discountdomain "project-example/internal/modules/discount/domain"

	"github.com/gin-gonic/gin"
)

type stubUseCase struct {
	discount *discountapplication.DiscountDTO
	err      error
}

func (s stubUseCase) GetDiscount(_ context.Context, _ string) (*discountapplication.DiscountDTO, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.discount, nil
}

func TestGetDiscountReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewHandler(stubUseCase{err: discountdomain.ErrDiscountNotFound})
	handler.RegisterRoutes(router.Group("/api/v1"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/discounts/UNKNOWN", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), "\"code\":\"discount_not_found\"") {
		t.Fatalf("expected discount_not_found code in body, got %s", recorder.Body.String())
	}
}
