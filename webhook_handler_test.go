package velafi

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func signPayload(secret string, body []byte) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyWebhookSignature(t *testing.T) {
	secret := "test-secret"
	body := []byte(`{"event":"order.fiat_crypto.completed","data":{}}`)
	sig := signPayload(secret, body)

	if !VerifyWebhookSignature(secret, body, sig) {
		t.Error("valid signature should return true")
	}
	if VerifyWebhookSignature(secret, body, "invalid-sig") {
		t.Error("invalid signature should return false")
	}
	if VerifyWebhookSignature("wrong-secret", body, sig) {
		t.Error("wrong secret should return false")
	}
}

func TestWebhookHandler_ValidEvent(t *testing.T) {
	secret := "my-secret"
	handler := NewWebhookHandler(secret)

	var receivedEvent *WebhookEvent
	handler.On(EventOrderFiatCryptoCompleted, func(ctx context.Context, event *WebhookEvent) error {
		receivedEvent = event
		return nil
	})

	payload := map[string]any{
		"event":     "order.fiat_crypto.completed",
		"timestamp": "2026-04-30T10:00:00Z",
		"data":      map[string]any{"orderId": "ord-123"},
	}
	body, _ := json.Marshal(payload)
	sig := signPayload(secret, body)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("X-Velafi-Signature", sig)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if receivedEvent == nil {
		t.Fatal("handler was not called")
	}
	if receivedEvent.Event != "order.fiat_crypto.completed" {
		t.Errorf("Event = %q", receivedEvent.Event)
	}
}

func TestWebhookHandler_InvalidSignature(t *testing.T) {
	handler := NewWebhookHandler("my-secret")
	handler.On(EventOrderFiatCryptoCompleted, func(ctx context.Context, event *WebhookEvent) error {
		t.Error("handler should not be called for invalid signature")
		return nil
	})

	body := []byte(`{"event":"order.fiat_crypto.completed","data":{}}`)
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("X-Velafi-Signature", "bad-signature")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want 403", w.Code)
	}
}

func TestWebhookHandler_UnregisteredEvent(t *testing.T) {
	handler := NewWebhookHandler("my-secret")

	payload := map[string]any{
		"event":     "order.fiat_fiat.completed",
		"timestamp": "2026-04-30T10:00:00Z",
		"data":      map[string]any{},
	}
	body, _ := json.Marshal(payload)
	sig := signPayload("my-secret", body)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("X-Velafi-Signature", sig)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 for unregistered event", w.Code)
	}
}

func TestWebhookHandler_HandlerError(t *testing.T) {
	secret := "my-secret"
	handler := NewWebhookHandler(secret)

	handler.On(EventFundingReceived, func(ctx context.Context, event *WebhookEvent) error {
		return errors.New("processing failed")
	})

	payload := map[string]any{
		"event":     "funding.received",
		"timestamp": "2026-04-30T10:00:00Z",
		"data":      map[string]any{"fundingId": "f-1"},
	}
	body, _ := json.Marshal(payload)
	sig := signPayload(secret, body)

	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(body))
	req.Header.Set("X-Velafi-Signature", sig)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want 500", w.Code)
	}
}
