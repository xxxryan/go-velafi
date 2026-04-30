package velafi

import "context"

func (c *Client) GetCryptoFiatQuote(ctx context.Context, params *CryptoFiatQuoteParams) (*Quote, error) {
	path := "/v1/quotes/crypto-fiat" + buildQuery(map[string]string{
		"merchantId":      params.MerchantID,
		"fromCurrency":    params.FromCurrency,
		"toCurrency":      params.ToCurrency,
		"fromAmount":      params.FromAmount,
		"toAmount":        params.ToAmount,
		"paymentMethodId": params.PaymentMethodID,
	})
	var result Quote
	err := c.get(ctx, path, &result)
	return &result, err
}

func (c *Client) GetFiatFiatQuote(ctx context.Context, params *FiatFiatQuoteParams) (*Quote, error) {
	path := "/v1/quotes/fiat-fiat" + buildQuery(map[string]string{
		"merchantId":          params.MerchantID,
		"fromCurrency":        params.FromCurrency,
		"toCurrency":          params.ToCurrency,
		"fromAmount":          params.FromAmount,
		"toAmount":            params.ToAmount,
		"fromPaymentMethodId": params.FromPaymentMethodID,
		"toPaymentMethodId":   params.ToPaymentMethodID,
	})
	var result Quote
	err := c.get(ctx, path, &result)
	return &result, err
}
