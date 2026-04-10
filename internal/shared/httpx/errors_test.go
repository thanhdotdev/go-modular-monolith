package httpx

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestWriteErrorAddsOriginalErrorToContextForUnhandledErrors(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	WriteError(context, errors.New("database connection failed"))

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}

	if !strings.Contains(recorder.Body.String(), "\"code\":\"internal_error\"") {
		t.Fatalf("expected internal_error response, got %s", recorder.Body.String())
	}

	if len(context.Errors) != 1 {
		t.Fatalf("expected 1 gin error, got %d", len(context.Errors))
	}

	if !strings.Contains(context.Errors.String(), "database connection failed") {
		t.Fatalf("expected original error in gin context, got %s", context.Errors.String())
	}
}

func TestWriteErrorDoesNotAddMappedErrorsToContext(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	WriteError(
		context,
		errNoRows,
		ErrorMapping{
			Err:    errNoRows,
			Status: http.StatusNotFound,
			Code:   "not_found",
		},
	)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}

	if len(context.Errors) != 0 {
		t.Fatalf("expected 0 gin errors, got %d", len(context.Errors))
	}
}

var errNoRows = errors.New("no rows")
