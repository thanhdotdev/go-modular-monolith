package customerhttp

import (
	"context"
	"net/http"
	"net/http/httptest"
	customerapplication "project-example/internal/modules/customer/application"
	customerdomain "project-example/internal/modules/customer/domain"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type stubUseCase struct {
	customer *customerapplication.CustomerDTO
	err      error
}

func (s stubUseCase) GetCustomer(_ context.Context, _ string) (*customerapplication.CustomerDTO, error) {
	if s.err != nil {
		return nil, s.err
	}

	return s.customer, nil
}

func TestGetCustomerReturnsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	handler := NewHandler(stubUseCase{err: customerdomain.ErrCustomerNotFound})
	handler.RegisterRoutes(router.Group("/api/v1"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/customers/cus-404", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), "\"code\":\"customer_not_found\"") {
		t.Fatalf("expected customer_not_found code in body, got %s", recorder.Body.String())
	}
}
