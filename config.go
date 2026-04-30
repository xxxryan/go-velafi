package velafi

import (
	"context"
	"net/url"
)

func (c *Client) ListCountries(ctx context.Context) ([]Country, error) {
	var result []Country
	err := c.get(ctx, "/v1/countries", &result)
	return result, err
}

func (c *Client) ListFiatCurrencies(ctx context.Context) ([]Currency, error) {
	var result []Currency
	err := c.get(ctx, "/v1/currencies/fiat", &result)
	return result, err
}

func (c *Client) ListCryptoCurrencies(ctx context.Context) ([]Currency, error) {
	var result []Currency
	err := c.get(ctx, "/v1/currencies/crypto", &result)
	return result, err
}

func (c *Client) ListFiatCryptoPairs(ctx context.Context) ([]Pair, error) {
	var result []Pair
	err := c.get(ctx, "/v1/pairs/fiat-crypto", &result)
	return result, err
}

func (c *Client) ListCryptoFiatPairs(ctx context.Context) ([]Pair, error) {
	var result []Pair
	err := c.get(ctx, "/v1/pairs/crypto-fiat", &result)
	return result, err
}

func (c *Client) ListFiatFiatPairs(ctx context.Context) ([]Pair, error) {
	var result []Pair
	err := c.get(ctx, "/v1/pairs/fiat-fiat", &result)
	return result, err
}

func (c *Client) ListFiatCryptoPaymentMethods(ctx context.Context) ([]PaymentMethodInfo, error) {
	var result []PaymentMethodInfo
	err := c.get(ctx, "/v1/payment-methods/fiat-crypto", &result)
	return result, err
}

func (c *Client) ListCryptoFiatPaymentMethods(ctx context.Context) ([]PaymentMethodInfo, error) {
	var result []PaymentMethodInfo
	err := c.get(ctx, "/v1/payment-methods/crypto-fiat", &result)
	return result, err
}

func (c *Client) ListFiatFiatPaymentMethods(ctx context.Context) ([]PaymentMethodInfo, error) {
	var result []PaymentMethodInfo
	err := c.get(ctx, "/v1/payment-methods/fiat-fiat", &result)
	return result, err
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
