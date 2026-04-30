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
		t.Fatal("token should not be empty after ensureToken")
	}
	if c.expiresAt.IsZero() {
		t.Fatal("expiresAt should be set")
	}
	t.Logf("token acquired, expires at %s", c.expiresAt.Format(time.RFC3339))
}

func TestIntegration_ListCountries(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	countries, err := c.ListCountries(ctx)
	if err != nil {
		t.Fatalf("ListCountries() error = %v", err)
	}

	if len(countries) == 0 {
		t.Fatal("expected at least one country")
	}
	t.Logf("got %d countries, first: %+v", len(countries), countries[0])
}

func TestIntegration_ListFiatCurrencies(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	currencies, err := c.ListFiatCurrencies(ctx)
	if err != nil {
		t.Fatalf("ListFiatCurrencies() error = %v", err)
	}

	if len(currencies) == 0 {
		t.Fatal("expected at least one fiat currency")
	}
	t.Logf("got %d fiat currencies, first: %+v", len(currencies), currencies[0])
}

func TestIntegration_ListCryptoCurrencies(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	currencies, err := c.ListCryptoCurrencies(ctx)
	if err != nil {
		t.Fatalf("ListCryptoCurrencies() error = %v", err)
	}

	if len(currencies) == 0 {
		t.Fatal("expected at least one crypto currency")
	}
	t.Logf("got %d crypto currencies, first: %+v", len(currencies), currencies[0])
}

func TestIntegration_ListBuySymbols(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	symbols, err := c.ListBuySymbols(ctx)
	if err != nil {
		t.Fatalf("ListBuySymbols() error = %v", err)
	}

	t.Logf("got %d buy symbols", len(symbols))
	if len(symbols) > 0 {
		t.Logf("first: %s/%s/%s", symbols[0].Country, symbols[0].Fiat, symbols[0].Crypto)
	}
}

func TestIntegration_ListSellSymbols(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	symbols, err := c.ListSellSymbols(ctx)
	if err != nil {
		t.Fatalf("ListSellSymbols() error = %v", err)
	}

	t.Logf("got %d sell symbols", len(symbols))
}

func TestIntegration_ListFiatFiatSymbols(t *testing.T) {
	c := integrationClient(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	symbols, err := c.ListFiatFiatSymbols(ctx)
	if err != nil {
		t.Fatalf("ListFiatFiatSymbols() error = %v", err)
	}

	t.Logf("got %d fiat-fiat symbols", len(symbols))
}
