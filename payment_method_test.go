package velafi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestListPaymentTemplates(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/payment-templates" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": []map[string]any{{
				"templateId": "tpl-1", "name": "Bank Transfer", "currency": "USD", "country": "US", "type": "bank",
				"fields": []map[string]any{{"key": "accountNumber", "label": "Account Number", "type": "string", "required": true}},
			}},
		})
	})

	templates, err := c.ListPaymentTemplates(context.Background(), &ListPaymentTemplatesParams{Currency: "USD"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(templates) != 1 || templates[0].TemplateID != "tpl-1" {
		t.Errorf("unexpected: %+v", templates)
	}
}

func TestGetPaymentTemplateMetamessage(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/payment-templates/tpl-1/metamessage" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{
				"templateId": "tpl-1",
				"metamessage": map[string]any{
					"fields": []map[string]any{{"key": "accountNumber", "label": "Account Number", "type": "string", "required": true}},
				},
			},
		})
	})

	meta, err := c.GetPaymentTemplateMetamessage(context.Background(), "tpl-1")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if meta.TemplateID != "tpl-1" {
		t.Errorf("TemplateID = %q", meta.TemplateID)
	}
}

func TestAddPaymentMethod(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/payment-methods" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{"methodId": "pm-1", "merchantId": "m-1", "templateId": "tpl-1", "status": "active", "fields": map[string]any{}, "createdAt": "2026-04-30T10:00:00Z"},
		})
	})

	pm, err := c.AddPaymentMethod(context.Background(), &AddPaymentMethodParams{
		MerchantID: "m-1", TemplateID: "tpl-1", Fields: map[string]any{"accountNumber": "123"},
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if pm.MethodID != "pm-1" {
		t.Errorf("MethodID = %q", pm.MethodID)
	}
}

func TestGetPaymentMethod(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/payment-methods/pm-1" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{"methodId": "pm-1", "merchantId": "m-1", "templateId": "tpl-1", "status": "active", "fields": map[string]any{}, "createdAt": "2026-04-30T10:00:00Z"},
		})
	})

	pm, err := c.GetPaymentMethod(context.Background(), "pm-1")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if pm.MethodID != "pm-1" {
		t.Errorf("MethodID = %q", pm.MethodID)
	}
}

func TestDeletePaymentMethod(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/payment-methods/pm-1" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success", "data": nil,
		})
	})

	err := c.DeletePaymentMethod(context.Background(), "pm-1")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}

func TestSetRefundAccount(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/payment-methods/pm-1/refund-account" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{"methodId": "pm-1", "refundMethodId": "pm-2", "updatedAt": "2026-04-30T10:00:00Z"},
		})
	})

	result, err := c.SetRefundAccount(context.Background(), "pm-1", "pm-2")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if result.RefundMethodID != "pm-2" {
		t.Errorf("RefundMethodID = %q", result.RefundMethodID)
	}
}
