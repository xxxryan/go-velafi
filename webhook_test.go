package velafi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateWebhook(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/webhooks" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{
				"webhookId": "wh-1", "url": "https://example.com/webhook",
				"events": []string{"order.fiat_crypto.completed"}, "merchantId": "m-1", "status": "active", "createdAt": "2026-04-30T10:00:00Z",
			},
		})
	})

	wh, err := c.CreateWebhook(context.Background(), &CreateWebhookParams{
		URL: "https://example.com/webhook", Events: []string{"order.fiat_crypto.completed"}, MerchantID: "m-1",
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if wh.WebhookID != "wh-1" {
		t.Errorf("WebhookID = %q", wh.WebhookID)
	}
}

func TestListWebhooks(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/webhooks" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("merchantId") != "m-1" {
			t.Errorf("merchantId = %q", r.URL.Query().Get("merchantId"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": []map[string]any{{"webhookId": "wh-1", "url": "https://example.com/webhook", "events": []string{"order.fiat_crypto.completed"}, "merchantId": "m-1", "status": "active"}},
		})
	})

	webhooks, err := c.ListWebhooks(context.Background(), "m-1")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(webhooks) != 1 || webhooks[0].WebhookID != "wh-1" {
		t.Errorf("unexpected: %+v", webhooks)
	}
}

func TestUpdateWebhook(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/webhooks/wh-1" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPut {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{"webhookId": "wh-1", "url": "https://new.example.com/webhook", "events": []string{"order.fiat_crypto.completed"}, "status": "active", "updatedAt": "2026-04-30T11:00:00Z"},
		})
	})

	wh, err := c.UpdateWebhook(context.Background(), "wh-1", &UpdateWebhookParams{URL: "https://new.example.com/webhook"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if wh.URL != "https://new.example.com/webhook" {
		t.Errorf("URL = %q", wh.URL)
	}
}

func TestDeleteWebhook(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/webhooks/wh-1" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{"code": 0, "message": "success", "data": nil})
	})

	err := c.DeleteWebhook(context.Background(), "wh-1")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}
