package velafi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestGenerateSignature(t *testing.T) {
	c := NewClient("key", "test-secret")

	ts, sig := c.generateSignature()

	if ts == "" {
		t.Fatal("timestamp should not be empty")
	}
	if _, err := strconv.ParseInt(ts, 10, 64); err != nil {
		t.Fatalf("timestamp should be parseable as int64: %v", err)
	}

	message := "timestamp=" + ts
	mac := hmac.New(sha256.New, []byte("test-secret"))
	mac.Write([]byte(message))
	expected := hex.EncodeToString(mac.Sum(nil))

	if sig != expected {
		t.Errorf("signature = %q, want %q", sig, expected)
	}
}

func TestEnsureToken_FetchesOnFirstCall(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/token/generate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(404)
			return
		}
		if r.Header.Get("X-BH-APIKEY") != "test-key" {
			t.Errorf("X-BH-APIKEY = %q, want %q", r.Header.Get("X-BH-APIKEY"), "test-key")
		}
		if r.URL.Query().Get("timestamp") == "" {
			t.Error("missing timestamp query param")
		}
		if r.URL.Query().Get("signature") == "" {
			t.Error("missing signature query param")
		}

		expireTime := strconv.FormatInt(time.Now().Add(1*time.Hour).UnixMilli(), 10)
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"token":      "fresh-token-123",
				"expireTime": expireTime,
			},
		})
	}))
	defer srv.Close()

	c := NewClient("test-key", "test-secret", WithBaseURL(srv.URL))

	err := c.ensureToken(context.Background())
	if err != nil {
		t.Fatalf("ensureToken() error = %v", err)
	}
	if c.token != "fresh-token-123" {
		t.Errorf("token = %q, want %q", c.token, "fresh-token-123")
	}
	if c.expiresAt.IsZero() {
		t.Error("expiresAt should be set")
	}
}

func TestEnsureToken_ReusesValidToken(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		expireTime := strconv.FormatInt(time.Now().Add(1*time.Hour).UnixMilli(), 10)
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"token":      "token-" + strconv.Itoa(callCount),
				"expireTime": expireTime,
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", "secret", WithBaseURL(srv.URL))

	_ = c.ensureToken(context.Background())
	_ = c.ensureToken(context.Background())

	if callCount != 1 {
		t.Errorf("token endpoint called %d times, want 1", callCount)
	}
}

func TestEnsureToken_RefreshesExpiredToken(t *testing.T) {
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		expireTime := strconv.FormatInt(time.Now().Add(1*time.Hour).UnixMilli(), 10)
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"token":      "token-" + strconv.Itoa(callCount),
				"expireTime": expireTime,
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", "secret", WithBaseURL(srv.URL))
	c.token = "old-token"
	c.expiresAt = time.Now().Add(-1 * time.Minute)

	err := c.ensureToken(context.Background())
	if err != nil {
		t.Fatalf("ensureToken() error = %v", err)
	}
	if callCount != 1 {
		t.Errorf("token endpoint called %d times, want 1", callCount)
	}
	if !strings.HasPrefix(c.token, "token-") {
		t.Errorf("token should have been refreshed, got %q", c.token)
	}
}

func TestEnsureToken_ConcurrentSafety(t *testing.T) {
	var mu sync.Mutex
	callCount := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		callCount++
		mu.Unlock()
		time.Sleep(50 * time.Millisecond)
		expireTime := strconv.FormatInt(time.Now().Add(1*time.Hour).UnixMilli(), 10)
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"token":      "concurrent-token",
				"expireTime": expireTime,
			},
		})
	}))
	defer srv.Close()

	c := NewClient("key", "secret", WithBaseURL(srv.URL))

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = c.ensureToken(context.Background())
		}()
	}
	wg.Wait()

	mu.Lock()
	defer mu.Unlock()
	if callCount != 1 {
		t.Errorf("token endpoint called %d times, want 1 (concurrent safety violated)", callCount)
	}
}
