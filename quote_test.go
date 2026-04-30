package velafi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetCryptoFiatQuote(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/quotes/crypto-fiat" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		q := r.URL.Query()
		if q.Get("merchantId") != "m-1" {
			t.Errorf("merchantId = %q", q.Get("merchantId"))
		}
		if q.Get("fromCurrency") != "BTC" {
			t.Errorf("fromCurrency = %q", q.Get("fromCurrency"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code":    0,
			"message": "success",
			"data": map[string]any{
				"quoteId":      "q-123",
				"fromCurrency": "BTC",
				"toCurrency":   "USD",
				"fromAmount":   "0.1",
				"toAmount":     "5000",
				"exchangeRate": "50000",
				"fee":          "10",
				"expiredAt":    "2026-04-30T10:05:00Z",
			},
		})
	})

	quote, err := c.GetCryptoFiatQuote(context.Background(), &CryptoFiatQuoteParams{
		MerchantID:      "m-1",
		FromCurrency:    "BTC",
		ToCurrency:      "USD",
		FromAmount:      "0.1",
		PaymentMethodID: "pm-1",
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if quote.QuoteID != "q-123" {
		t.Errorf("QuoteID = %q, want %q", quote.QuoteID, "q-123")
	}
	if quote.ToAmount != "5000" {
		t.Errorf("ToAmount = %q, want %q", quote.ToAmount, "5000")
	}
}

func TestGetFiatFiatQuote(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/quotes/fiat-fiat" {
			t.Errorf("path = %q", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("fromPaymentMethodId") != "pm-from" {
			t.Errorf("fromPaymentMethodId = %q", q.Get("fromPaymentMethodId"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code":    0,
			"message": "success",
			"data": map[string]any{
				"quoteId":      "q-456",
				"fromCurrency": "EUR",
				"toCurrency":   "USD",
				"fromAmount":   "1000",
				"toAmount":     "1100",
				"exchangeRate": "1.1",
				"fee":          "5",
				"expiredAt":    "2026-04-30T10:05:00Z",
			},
		})
	})

	quote, err := c.GetFiatFiatQuote(context.Background(), &FiatFiatQuoteParams{
		MerchantID:          "m-1",
		FromCurrency:        "EUR",
		ToCurrency:          "USD",
		FromAmount:          "1000",
		FromPaymentMethodID: "pm-from",
		ToPaymentMethodID:   "pm-to",
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if quote.QuoteID != "q-456" {
		t.Errorf("QuoteID = %q", quote.QuoteID)
	}
}
