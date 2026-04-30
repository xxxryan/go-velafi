package velafi

import (
	"context"
	"strconv"
)

func (c *Client) GetCryptoQuote(ctx context.Context, params *CryptoQuoteParams) (*Quote, error) {
	q := map[string]string{
		"country": params.Country,
		"from":    params.From,
		"to":      params.To,
	}
	if params.CreateQuoteID {
		q["createQuoteId"] = strconv.FormatBool(params.CreateQuoteID)
	}
	var result Quote
	err := c.get(ctx, "/v2/user/crypto-quote"+buildQuery(q), &result)
	return &result, err
}

func (c *Client) GetFiatQuote(ctx context.Context, params *FiatQuoteParams) (*Quote, error) {
	q := map[string]string{
		"onRampCountry":  params.OnRampCountry,
		"onRampFiat":     params.OnRampFiat,
		"offRampCountry": params.OffRampCountry,
		"offRampFiat":    params.OffRampFiat,
	}
	if params.CreateQuoteID {
		q["createQuoteId"] = strconv.FormatBool(params.CreateQuoteID)
	}
	var result Quote
	err := c.get(ctx, "/v2/user/fiat-quote"+buildQuery(q), &result)
	return &result, err
}
