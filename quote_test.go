package velafi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetCryptoQuote(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/user/crypto-quote" {
			t.Errorf("path = %q, want /v2/user/crypto-quote", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		q := r.URL.Query()
		if q.Get("country") != "US" {
			t.Errorf("country = %q, want US", q.Get("country"))
		}
		if q.Get("from") != "USD" {
			t.Errorf("from = %q, want USD", q.Get("from"))
		}
		if q.Get("to") != "BTC" {
			t.Errorf("to = %q, want BTC", q.Get("to"))
		}
		if q.Get("createQuoteId") != "true" {
			t.Errorf("createQuoteId = %q, want true", q.Get("createQuoteId"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"price":   "50000.00",
				"quoteId": "q-abc-123",
			},
		})
	})

	quote, err := c.GetCryptoQuote(context.Background(), &CryptoQuoteParams{
		Country:       "US",
		From:          "USD",
		To:            "BTC",
		CreateQuoteID: true,
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if quote.Price != "50000.00" {
		t.Errorf("Price = %q, want %q", quote.Price, "50000.00")
	}
	if quote.QuoteID != "q-abc-123" {
		t.Errorf("QuoteID = %q, want %q", quote.QuoteID, "q-abc-123")
	}
}

func TestGetFiatQuote(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/user/fiat-quote" {
			t.Errorf("path = %q, want /v2/user/fiat-quote", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("onRampCountry") != "US" {
			t.Errorf("onRampCountry = %q, want US", q.Get("onRampCountry"))
		}
		if q.Get("onRampFiat") != "USD" {
			t.Errorf("onRampFiat = %q, want USD", q.Get("onRampFiat"))
		}
		if q.Get("offRampCountry") != "GB" {
			t.Errorf("offRampCountry = %q, want GB", q.Get("offRampCountry"))
		}
		if q.Get("offRampFiat") != "GBP" {
			t.Errorf("offRampFiat = %q, want GBP", q.Get("offRampFiat"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"price":   "0.79",
				"quoteId": "q-fiat-456",
			},
		})
	})

	quote, err := c.GetFiatQuote(context.Background(), &FiatQuoteParams{
		OnRampCountry:  "US",
		OnRampFiat:     "USD",
		OffRampCountry: "GB",
		OffRampFiat:    "GBP",
		CreateQuoteID:  true,
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if quote.Price != "0.79" {
		t.Errorf("Price = %q, want %q", quote.Price, "0.79")
	}
	if quote.QuoteID != "q-fiat-456" {
		t.Errorf("QuoteID = %q, want %q", quote.QuoteID, "q-fiat-456")
	}
}
