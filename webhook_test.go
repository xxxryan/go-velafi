package velafi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateWebhook(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/webhook" {
			t.Errorf("path = %q, want /v2/webhook", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["eventType"] != "order.fiat_crypto.completed" {
			t.Errorf("eventType = %v, want order.fiat_crypto.completed", body["eventType"])
		}
		if body["url"] != "https://example.com/webhook" {
			t.Errorf("url = %v, want https://example.com/webhook", body["url"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"webhookId": "wh-abc-123",
				"eventType": "order.fiat_crypto.completed",
				"url":       "https://example.com/webhook",
				"status":    1,
				"publicKey": "pk-test-key",
			},
		})
	})

	wh, err := c.CreateWebhook(context.Background(), &CreateWebhookParams{
		EventType: "order.fiat_crypto.completed",
		URL:       "https://example.com/webhook",
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if wh.WebhookID != "wh-abc-123" {
		t.Errorf("WebhookID = %q, want %q", wh.WebhookID, "wh-abc-123")
	}
	if wh.EventType != "order.fiat_crypto.completed" {
		t.Errorf("EventType = %q, want %q", wh.EventType, "order.fiat_crypto.completed")
	}
	if wh.URL != "https://example.com/webhook" {
		t.Errorf("URL = %q, want %q", wh.URL, "https://example.com/webhook")
	}
	if wh.Status != 1 {
		t.Errorf("Status = %d, want 1", wh.Status)
	}
	if wh.PublicKey != "pk-test-key" {
		t.Errorf("PublicKey = %q, want %q", wh.PublicKey, "pk-test-key")
	}
}

func TestListWebhooks(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/webhooks" {
			t.Errorf("path = %q, want /v2/webhooks", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		q := r.URL.Query()
		if q.Get("status") != "1" {
			t.Errorf("status = %q, want 1", q.Get("status"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": []map[string]any{
				{
					"webhookId": "wh-abc-123",
					"eventType": "order.fiat_crypto.completed",
					"url":       "https://example.com/webhook",
					"status":    1,
				},
			},
		})
	})

	webhooks, err := c.ListWebhooks(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(webhooks) != 1 {
		t.Fatalf("len = %d, want 1", len(webhooks))
	}
	if webhooks[0].WebhookID != "wh-abc-123" {
		t.Errorf("WebhookID = %q, want %q", webhooks[0].WebhookID, "wh-abc-123")
	}
	if webhooks[0].EventType != "order.fiat_crypto.completed" {
		t.Errorf("EventType = %q, want %q", webhooks[0].EventType, "order.fiat_crypto.completed")
	}
}
