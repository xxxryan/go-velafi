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

func TestListCountries(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/base/countrys" {
			t.Errorf("path = %q, want /v2/base/countrys", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": []map[string]any{{"country": "United States", "abbr": "US"}},
		})
	})

	countries, err := c.ListCountries(context.Background())
	if err != nil {
		t.Fatalf("ListCountries() error = %v", err)
	}
	if len(countries) != 1 {
		t.Fatalf("len = %d, want 1", len(countries))
	}
	if countries[0].Country != "United States" {
		t.Errorf("Country = %q, want %q", countries[0].Country, "United States")
	}
	if countries[0].Abbr != "US" {
		t.Errorf("Abbr = %q, want %q", countries[0].Abbr, "US")
	}
}

func TestListFiatCurrencies(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/base/fiats" {
			t.Errorf("path = %q, want /v2/base/fiats", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": []map[string]any{{"country": "US", "fiat": "USD", "accuracy": 2}},
		})
	})

	currencies, err := c.ListFiatCurrencies(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(currencies) != 1 {
		t.Fatalf("len = %d, want 1", len(currencies))
	}
	if currencies[0].Fiat != "USD" {
		t.Errorf("Fiat = %q, want %q", currencies[0].Fiat, "USD")
	}
	if currencies[0].Accuracy != 2 {
		t.Errorf("Accuracy = %d, want 2", currencies[0].Accuracy)
	}
}

func TestListCryptoCurrencies(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/base/cryptos" {
			t.Errorf("path = %q, want /v2/base/cryptos", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": []map[string]any{{"crypto": "BTC", "accuracy": 8}},
		})
	})

	currencies, err := c.ListCryptoCurrencies(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(currencies) != 1 {
		t.Fatalf("len = %d, want 1", len(currencies))
	}
	if currencies[0].Crypto != "BTC" {
		t.Errorf("Crypto = %q, want %q", currencies[0].Crypto, "BTC")
	}
	if currencies[0].Accuracy != 8 {
		t.Errorf("Accuracy = %d, want 8", currencies[0].Accuracy)
	}
}

func TestListBuySymbols(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/base/buy/symbols" {
			t.Errorf("path = %q, want /v2/base/buy/symbols", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": []map[string]any{{"country": "US", "fiat": "USD", "crypto": "BTC", "accuracy": 8}},
		})
	})

	symbols, err := c.ListBuySymbols(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(symbols) != 1 {
		t.Fatalf("len = %d, want 1", len(symbols))
	}
	if symbols[0].Country != "US" {
		t.Errorf("Country = %q, want %q", symbols[0].Country, "US")
	}
	if symbols[0].Crypto != "BTC" {
		t.Errorf("Crypto = %q, want %q", symbols[0].Crypto, "BTC")
	}
}

func TestListBuyPayments(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/base/buy/payments" {
			t.Errorf("path = %q, want /v2/base/buy/payments", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("country") != "US" {
			t.Errorf("country = %q, want US", q.Get("country"))
		}
		if q.Get("fiat") != "USD" {
			t.Errorf("fiat = %q, want USD", q.Get("fiat"))
		}
		if q.Get("crypto") != "BTC" {
			t.Errorf("crypto = %q, want BTC", q.Get("crypto"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"country": "US",
				"fiat":    "USD",
				"crypto":  "BTC",
				"paymentList": []map[string]any{
					{"paymentId": 1, "fiatFee": "2.5", "paymentType": 1},
				},
			},
		})
	})

	payments, err := c.ListBuyPayments(context.Background(), "US", "USD", "BTC")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if payments.Country != "US" {
		t.Errorf("Country = %q, want %q", payments.Country, "US")
	}
	if len(payments.PaymentList) != 1 {
		t.Fatalf("len(PaymentList) = %d, want 1", len(payments.PaymentList))
	}
	if payments.PaymentList[0].PaymentID != 1 {
		t.Errorf("PaymentID = %d, want 1", payments.PaymentList[0].PaymentID)
	}
	if payments.PaymentList[0].FiatFee != "2.5" {
		t.Errorf("FiatFee = %q, want %q", payments.PaymentList[0].FiatFee, "2.5")
	}
}
