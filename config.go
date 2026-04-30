package velafi

import (
	"context"
	"net/url"
)

func (c *Client) ListCountries(ctx context.Context) ([]Country, error) {
	var result []Country
	err := c.get(ctx, "/v2/base/countrys", &result)
	return result, err
}

func (c *Client) ListFiatCurrencies(ctx context.Context) ([]FiatCurrency, error) {
	var result []FiatCurrency
	err := c.get(ctx, "/v2/base/fiats", &result)
	return result, err
}

func (c *Client) ListCryptoCurrencies(ctx context.Context) ([]CryptoCurrency, error) {
	var result []CryptoCurrency
	err := c.get(ctx, "/v2/base/cryptos", &result)
	return result, err
}

func (c *Client) ListBuySymbols(ctx context.Context) ([]BuySellSymbol, error) {
	var result []BuySellSymbol
	err := c.get(ctx, "/v2/base/buy/symbols", &result)
	return result, err
}

func (c *Client) ListSellSymbols(ctx context.Context) ([]BuySellSymbol, error) {
	var result []BuySellSymbol
	err := c.get(ctx, "/v2/base/sell/symbols", &result)
	return result, err
}

func (c *Client) ListFiatFiatSymbols(ctx context.Context) ([]FiatFiatSymbol, error) {
	var result []FiatFiatSymbol
	err := c.get(ctx, "/v2/base/fiat/symbols", &result)
	return result, err
}

func (c *Client) ListBuyPayments(ctx context.Context, country, fiat, crypto string) (*BuySellPayments, error) {
	path := "/v2/base/buy/payments" + buildQuery(map[string]string{
		"country": country, "fiat": fiat, "crypto": crypto,
	})
	var result BuySellPayments
	err := c.get(ctx, path, &result)
	return &result, err
}

func (c *Client) ListSellPayments(ctx context.Context, country, fiat, crypto string) (*BuySellPayments, error) {
	path := "/v2/base/sell/payments" + buildQuery(map[string]string{
		"country": country, "fiat": fiat, "crypto": crypto,
	})
	var result BuySellPayments
	err := c.get(ctx, path, &result)
	return &result, err
}

func (c *Client) ListFiatFiatPayments(ctx context.Context, onRampCountry, onRampFiat, offRampCountry, offRampFiat string) (*FiatFiatPayments, error) {
	path := "/v2/base/fiat/payments" + buildQuery(map[string]string{
		"onRampCountry": onRampCountry, "onRampFiat": onRampFiat,
		"offRampCountry": offRampCountry, "offRampFiat": offRampFiat,
	})
	var result FiatFiatPayments
	err := c.get(ctx, path, &result)
	return &result, err
}

func buildQuery(params map[string]string) string {
	v := url.Values{}
	for key, val := range params {
		if val != "" {
			v.Set(key, val)
		}
	}
	if len(v) == 0 {
		return ""
	}
	return "?" + v.Encode()
}
