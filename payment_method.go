package velafi

import (
	"context"
	"strconv"
)

func (c *Client) GetPaymentTemplate(ctx context.Context, paymentID int) (PaymentTemplate, error) {
	path := "/v2/payments/templates" + buildQuery(map[string]string{
		"paymentId": strconv.Itoa(paymentID),
	})
	var result PaymentTemplate
	err := c.get(ctx, path, &result)
	return result, err
}

func (c *Client) GetPaymentTemplateMeta(ctx context.Context, paymentID int) (*PaymentTemplateMeta, error) {
	path := "/v2/payments/templates/metamessage" + buildQuery(map[string]string{
		"paymentId": strconv.Itoa(paymentID),
	})
	var result PaymentTemplateMeta
	err := c.get(ctx, path, &result)
	return &result, err
}

func (c *Client) AddPaymentMethod(ctx context.Context, params *AddPaymentMethodParams) (*AddPaymentMethodResult, error) {
	var result AddPaymentMethodResult
	err := c.post(ctx, "/v2/payments", params, &result)
	return &result, err
}

func (c *Client) ListPaymentMethods(ctx context.Context, params *ListPaymentMethodsParams) (*PaymentMethodList, error) {
	q := map[string]string{
		"country": params.Country,
		"fiat":    params.Fiat,
	}
	if params.Status > 0 {
		q["status"] = strconv.Itoa(params.Status)
	}
	if params.MerchantID > 0 {
		q["merchantId"] = strconv.Itoa(params.MerchantID)
	}
	if params.CurrentPage > 0 {
		q["currentPage"] = strconv.Itoa(params.CurrentPage)
	}
	if params.PageSize > 0 {
		q["pageSize"] = strconv.Itoa(params.PageSize)
	}
	var result PaymentMethodList
	err := c.get(ctx, "/v2/payments"+buildQuery(q), &result)
	return &result, err
}

func (c *Client) DeletePaymentMethod(ctx context.Context, userPaymentID int) error {
	return c.delete(ctx, "/v2/payments/"+strconv.Itoa(userPaymentID))
}
