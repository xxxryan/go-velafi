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
			t.Errorf("path = %q", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("country") != "Mexico" || q.Get("from") != "MXN" || q.Get("to") != "USDT" {
			t.Errorf("query = %v", q)
		}
		if q.Get("createQuoteId") != "true" {
			t.Errorf("createQuoteId = %q", q.Get("createQuoteId"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200, "msg": "SUCCESS",
			"data": map[string]any{"price": "17.84", "quoteId": "abc123"},
		})
	})

	quote, err := c.GetCryptoQuote(context.Background(), &CryptoQuoteParams{
		Country: "Mexico", From: "MXN", To: "USDT", CreateQuoteID: true,
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if quote.Price != "17.84" {
		t.Errorf("Price = %q", quote.Price)
	}
	if quote.QuoteID != "abc123" {
		t.Errorf("QuoteID = %q", quote.QuoteID)
	}
}
