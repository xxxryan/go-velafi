package velafi

import (
	"context"
	"net/url"
)

func (c *Client) GetBuySymbols(ctx context.Context) ([]BuySellSymbol, error) {
	var result []BuySellSymbol
	err := c.get(ctx, "/v2/base/buy/symbols", &result)
	return result, err
}

func (c *Client) GetSellSymbols(ctx context.Context) ([]BuySellSymbol, error) {
	var result []BuySellSymbol
	err := c.get(ctx, "/v2/base/sell/symbols", &result)
	return result, err
}

func (c *Client) GetBuyPayments(ctx context.Context, country, fiat, crypto string) (*BuySellPayments, error) {
	path := "/v2/base/buy/payments" + buildQuery(map[string]string{
		"country": country, "fiat": fiat, "crypto": crypto,
	})
	var result BuySellPayments
	err := c.get(ctx, path, &result)
	return &result, err
}

func (c *Client) GetSellPayments(ctx context.Context, country, fiat, crypto string) (*BuySellPayments, error) {
	path := "/v2/base/sell/payments" + buildQuery(map[string]string{
		"country": country, "fiat": fiat, "crypto": crypto,
	})
	var result BuySellPayments
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
