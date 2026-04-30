package velafi

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
)

const (
	EventOrderFiatCryptoCreated     = "order.fiat_crypto.created"
	EventOrderFiatCryptoCompleted   = "order.fiat_crypto.completed"
	EventOrderFiatCryptoFailed      = "order.fiat_crypto.failed"
	EventOrderFiatFiatCreated       = "order.fiat_fiat.created"
	EventOrderFiatFiatCompleted     = "order.fiat_fiat.completed"
	EventOrderFiatFiatFailed        = "order.fiat_fiat.failed"
	EventFundingReceived            = "funding.received"
	EventFundingCompleted           = "funding.completed"
	EventStablecoinPaymentCreated   = "stablecoin.payment.created"
	EventStablecoinPaymentCompleted = "stablecoin.payment.completed"
)

type EventHandlerFunc func(ctx context.Context, event *WebhookEvent) error

type WebhookEvent struct {
	Event     string          `json:"event"`
	Timestamp string          `json:"timestamp"`
	RawData   json.RawMessage `json:"data"`
}

type WebhookHandler struct {
	secret   string
	handlers map[string]EventHandlerFunc
}

func NewWebhookHandler(secret string) *WebhookHandler {
	return &WebhookHandler{
		secret:   secret,
		handlers: make(map[string]EventHandlerFunc),
	}
}

func (h *WebhookHandler) On(event string, handler EventHandlerFunc) {
	h.handlers[event] = handler
}

func (h *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("X-Velafi-Signature")
	if !VerifyWebhookSignature(h.secret, body, signature) {
		http.Error(w, "invalid signature", http.StatusForbidden)
		return
	}

	var event WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	handler, ok := h.handlers[event.Event]
	if !ok {
		w.WriteHeader(http.StatusOK)
		return
	}

	if err := handler(r.Context(), &event); err != nil {
		http.Error(w, "handler error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func VerifyWebhookSignature(secret string, body []byte, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}
