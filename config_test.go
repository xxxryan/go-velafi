package velafi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	c := NewClient("key", "secret", WithBaseURL(srv.URL))
	c.token = "test-token"
	c.expiresAt = time.Now().Add(1 * time.Hour)
	return c
}

func TestGetBuySymbols(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/base/buy/symbols" {
			t.Errorf("path = %q, want /v2/base/buy/symbols", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200, "msg": "SUCCESS",
			"data": []map[string]any{{"country": "Mexico", "fiat": "MXN", "crypto": "USDT", "accuracy": 2}},
		})
	})

	symbols, err := c.GetBuySymbols(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(symbols) != 1 || symbols[0].Fiat != "MXN" {
		t.Errorf("unexpected: %+v", symbols)
	}
}

func TestGetSellSymbols(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/base/sell/symbols" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200, "msg": "SUCCESS",
			"data": []map[string]any{{"country": "Mexico", "fiat": "MXN", "crypto": "USDT", "accuracy": 2}},
		})
	})

	symbols, err := c.GetSellSymbols(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(symbols) != 1 || symbols[0].Crypto != "USDT" {
		t.Errorf("unexpected: %+v", symbols)
	}
}

func TestGetBuyPayments(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/base/buy/payments" {
			t.Errorf("path = %q", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("country") != "Mexico" || q.Get("fiat") != "MXN" || q.Get("crypto") != "USDT" {
			t.Errorf("query = %v", q)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200, "msg": "SUCCESS",
			"data": map[string]any{
				"country": "Mexico", "fiat": "MXN", "crypto": "USDT",
				"paymentList": []map[string]any{{"paymentId": 72, "fiatFee": 1.0, "paymentType": 1, "trench": "CLABE - TESORED"}},
			},
		})
	})

	payments, err := c.GetBuyPayments(context.Background(), "Mexico", "MXN", "USDT")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if payments.Country != "Mexico" {
		t.Errorf("Country = %q", payments.Country)
	}
	if len(payments.PaymentList) != 1 || payments.PaymentList[0].PaymentID != 72 {
		t.Errorf("unexpected PaymentList: %+v", payments.PaymentList)
	}
}

func TestGetSellPayments(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/base/sell/payments" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200, "msg": "SUCCESS",
			"data": map[string]any{
				"country": "Mexico", "fiat": "MXN", "crypto": "USDT",
				"paymentList": []map[string]any{{"paymentId": 14, "fiatFee": 0.0, "paymentType": 0}},
			},
		})
	})

	payments, err := c.GetSellPayments(context.Background(), "Mexico", "MXN", "USDT")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(payments.PaymentList) != 1 {
		t.Errorf("unexpected PaymentList: %+v", payments.PaymentList)
	}
}
