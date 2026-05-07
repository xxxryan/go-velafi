package velafi

import (
	"context"
	"os"
	"testing"
	"time"
)

func integrationClient(t *testing.T) *Client {
	t.Helper()

	apiKey := os.Getenv("VELAFI_API_KEY")
	apiSecret := os.Getenv("VELAFI_API_SECRET")

	if apiKey == "" || apiSecret == "" {
		t.Skip("skipping integration test: set VELAFI_API_KEY, VELAFI_API_SECRET")
	}

	return NewClient(apiKey, apiSecret, WithSandbox())
}

func TestIntegration_TokenGeneration(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err := c.ensureToken(ctx)
	if err != nil {
		t.Fatalf("ensureToken() error = %v", err)
	}
	if c.token == "" {
		t.Fatal("token should not be empty")
	}
	t.Logf("token acquired, expires at %s", c.expiresAt.Format(time.RFC3339))
}

func TestIntegration_GetBuySymbols(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	symbols, err := c.GetBuySymbols(ctx)
	if err != nil {
		t.Fatalf("GetBuySymbols() error = %v", err)
	}
	if len(symbols) == 0 {
		t.Fatal("expected at least one symbol")
	}
	t.Logf("got %d buy symbols, first: %+v", len(symbols), symbols[0])
}

func TestIntegration_GetSellSymbols(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	symbols, err := c.GetSellSymbols(ctx)
	if err != nil {
		t.Fatalf("GetSellSymbols() error = %v", err)
	}
	if len(symbols) == 0 {
		t.Fatal("expected at least one symbol")
	}
	t.Logf("got %d sell symbols, first: %+v", len(symbols), symbols[0])
}

func TestIntegration_GetBuyPayments(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payments, err := c.GetBuyPayments(ctx, "Mexico", "MXN", "USDT")
	if err != nil {
		t.Fatalf("GetBuyPayments() error = %v", err)
	}
	t.Logf("got %d buy payments for Mexico/MXN/USDT", len(payments.PaymentList))
	for i, p := range payments.PaymentList {
		t.Logf("  [%d] paymentId=%d, fee=%.2f, type=%d, trench=%s", i, p.PaymentID, p.FiatFee, p.PaymentType, p.Trench)
	}
}

func TestIntegration_GetSellPayments(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	payments, err := c.GetSellPayments(ctx, "Mexico", "MXN", "USDT")
	if err != nil {
		t.Fatalf("GetSellPayments() error = %v", err)
	}
	t.Logf("got %d sell payments for Mexico/MXN/USDT", len(payments.PaymentList))
}

func TestIntegration_GetCryptoQuote(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	quote, err := c.GetCryptoQuote(ctx, &CryptoQuoteParams{
		Country: "Mexico", From: "MXN", To: "USDT",
	})
	if err != nil {
		t.Fatalf("GetCryptoQuote() error = %v", err)
	}
	t.Logf("MXN/USDT price: %s", quote.Price)
}
