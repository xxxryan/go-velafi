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
		if r.URL.Path != "/v1/countries" {
			t.Errorf("path = %q, want /v1/countries", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code":    0,
			"message": "success",
			"data":    []map[string]any{{"code": "US", "name": "United States"}},
		})
	})

	countries, err := c.ListCountries(context.Background())
	if err != nil {
		t.Fatalf("ListCountries() error = %v", err)
	}
	if len(countries) != 1 {
		t.Fatalf("len = %d, want 1", len(countries))
	}
	if countries[0].Code != "US" {
		t.Errorf("Code = %q, want %q", countries[0].Code, "US")
	}
}

func TestListFiatCurrencies(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/currencies/fiat" {
			t.Errorf("path = %q, want /v1/currencies/fiat", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code":    0,
			"message": "success",
			"data":    []map[string]any{{"code": "USD", "name": "US Dollar"}},
		})
	})

	currencies, err := c.ListFiatCurrencies(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(currencies) != 1 || currencies[0].Code != "USD" {
		t.Errorf("unexpected result: %+v", currencies)
	}
}

func TestListCryptoCurrencies(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/currencies/crypto" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code":    0,
			"message": "success",
			"data":    []map[string]any{{"code": "BTC", "name": "Bitcoin"}},
		})
	})

	currencies, err := c.ListCryptoCurrencies(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(currencies) != 1 || currencies[0].Code != "BTC" {
		t.Errorf("unexpected: %+v", currencies)
	}
}

func TestListFiatCryptoPairs(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/pairs/fiat-crypto" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code":    0,
			"message": "success",
			"data":    []map[string]any{{"fromCurrency": "USD", "toCurrency": "BTC"}},
		})
	})

	pairs, err := c.ListFiatCryptoPairs(context.Background())
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if len(pairs) != 1 || pairs[0].FromCurrency != "USD" {
		t.Errorf("unexpected: %+v", pairs)
	}
}
