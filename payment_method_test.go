package velafi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetPaymentTemplate(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payments/templates" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.URL.Query().Get("paymentId") != "72" {
			t.Errorf("paymentId = %q", r.URL.Query().Get("paymentId"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200, "msg": "SUCCESS",
			"data": map[string]any{"Account Type": "clabe", "Full Name": "Name"},
		})
	})

	tpl, err := c.GetPaymentTemplate(context.Background(), 72)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if tpl["Account Type"] != "clabe" {
		t.Errorf("Account Type = %q", tpl["Account Type"])
	}
}

func TestAddPaymentMethod(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payments" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200, "msg": "SUCCESS",
			"data": map[string]any{"id": 456, "status": 1, "failReason": ""},
		})
	})

	result, err := c.AddPaymentMethod(context.Background(), &AddPaymentMethodParams{
		MerchantID: 3, PaymentID: 72, Country: "Mexico", Fiat: "MXN",
		RealName: "John Doe", FieldJSON: map[string]any{"Account Type": "clabe"},
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if result.ID != 456 {
		t.Errorf("ID = %d", result.ID)
	}
	if result.Status != 1 {
		t.Errorf("Status = %d", result.Status)
	}
}
