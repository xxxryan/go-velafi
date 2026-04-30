package velafi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient_Defaults(t *testing.T) {
	c := NewClient("key", "secret")
	if c.apiKey != "key" {
		t.Errorf("apiKey = %q, want %q", c.apiKey, "key")
	}
	if c.apiSecret != "secret" {
		t.Errorf("apiSecret = %q, want %q", c.apiSecret, "secret")
	}
	if c.baseURL != "https://api.velafi.com" {
		t.Errorf("baseURL = %q, want %q", c.baseURL, "https://api.velafi.com")
	}
	if c.tokenRefreshBuffer != 5*time.Minute {
		t.Errorf("tokenRefreshBuffer = %v, want %v", c.tokenRefreshBuffer, 5*time.Minute)
	}
}

func TestNewClient_WithOptions(t *testing.T) {
	custom := &http.Client{Timeout: 30 * time.Second}
	c := NewClient("key", "secret",
		WithBaseURL("https://custom.api.velafi.com"),
		WithHTTPClient(custom),
		WithTokenRefreshBuffer(10*time.Minute),
	)
	if c.baseURL != "https://custom.api.velafi.com" {
		t.Errorf("baseURL = %q, want %q", c.baseURL, "https://custom.api.velafi.com")
	}
	if c.httpClient != custom {
		t.Error("httpClient should be the custom client")
	}
	if c.tokenRefreshBuffer != 10*time.Minute {
		t.Errorf("tokenRefreshBuffer = %v, want %v", c.tokenRefreshBuffer, 10*time.Minute)
	}
}

func TestDoJSON_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-BH-TOKEN") != "test-token" {
			t.Errorf("missing or wrong X-BH-TOKEN header: %q", r.Header.Get("X-BH-TOKEN"))
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{"name": "TestCountry"},
		})
	}))
	defer srv.Close()

	c := NewClient("key", "secret", WithBaseURL(srv.URL))
	c.token = "test-token"
	c.expiresAt = time.Now().Add(1 * time.Hour)

	var result struct {
		Name string `json:"name"`
	}
	err := c.doJSON(context.Background(), http.MethodGet, "/v2/base/countrys", nil, &result)
	if err != nil {
		t.Fatalf("doJSON() error = %v", err)
	}
	if result.Name != "TestCountry" {
		t.Errorf("result.Name = %q, want %q", result.Name, "TestCountry")
	}
}

func TestDoJSON_BusinessError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"code": 10001,
			"msg":  "invalid parameter",
		})
	}))
	defer srv.Close()

	c := NewClient("key", "secret", WithBaseURL(srv.URL))
	c.token = "test-token"
	c.expiresAt = time.Now().Add(1 * time.Hour)

	var result any
	err := c.doJSON(context.Background(), http.MethodGet, "/v2/test", nil, &result)
	if err == nil {
		t.Fatal("doJSON() should return error for non-200 code")
	}
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error should be *Error, got %T", err)
	}
	if apiErr.Code != 10001 {
		t.Errorf("Code = %d, want 10001", apiErr.Code)
	}
	if apiErr.Message != "invalid parameter" {
		t.Errorf("Message = %q, want %q", apiErr.Message, "invalid parameter")
	}
}

func TestDoJSON_HTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"code": 50000,
			"msg":  "internal error",
		})
	}))
	defer srv.Close()

	c := NewClient("key", "secret", WithBaseURL(srv.URL))
	c.token = "test-token"
	c.expiresAt = time.Now().Add(1 * time.Hour)

	var result any
	err := c.doJSON(context.Background(), http.MethodGet, "/v2/test", nil, &result)
	if err == nil {
		t.Fatal("doJSON() should return error for HTTP 500")
	}
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Fatalf("error should be *Error, got %T", err)
	}
	if apiErr.HTTPStatus != 500 {
		t.Errorf("HTTPStatus = %d, want 500", apiErr.HTTPStatus)
	}
}
