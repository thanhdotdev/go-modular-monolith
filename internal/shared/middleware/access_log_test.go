package middleware

import (
	"io"
	"net/http/httptest"
	"project-example/internal/platform/config"
	"strings"
	"testing"
)

func TestCaptureRequestBodyLimitsLogPayloadAndPreservesRequestBody(t *testing.T) {
	fullBody := "abcdefghij"
	request := httptest.NewRequest("POST", "/orders", strings.NewReader(fullBody))
	request.Header.Set("Content-Type", "application/json")

	payload, truncated := captureRequestBody(request, config.LoggingConfig{
		IncludeRequestBody: true,
		BodyMaxBytes:       5,
	})

	if payload != "abcde" {
		t.Fatalf("expected payload abcde, got %q", payload)
	}

	if !truncated {
		t.Fatal("expected payload to be truncated")
	}

	remainingBody, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatalf("read preserved request body: %v", err)
	}

	if string(remainingBody) != fullBody {
		t.Fatalf("expected request body %q, got %q", fullBody, string(remainingBody))
	}
}
