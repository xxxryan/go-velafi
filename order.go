package velafi

import (
	"context"
	"strconv"
)

func (c *Client) CreateFiatCryptoOrder(ctx context.Context, params *CreateFiatCryptoOrderParams) (*Order, error) {
	var result Order
	err := c.post(ctx, "/v1/orders/fiat-crypto", params, &result)
	return &result, err
}

func (c *Client) CreateCryptoFiatOrder(ctx context.Context, params *CreateCryptoFiatOrderParams) (*Order, error) {
	var result Order
	err := c.post(ctx, "/v1/orders/crypto-fiat", params, &result)
	return &result, err
}

func (c *Client) CreateFiatFiatOrder(ctx context.Context, params *CreateFiatFiatOrderParams) (*Order, error) {
	var result Order
	err := c.post(ctx, "/v1/orders/fiat-fiat", params, &result)
	return &result, err
}

func (c *Client) ConfirmOrder(ctx context.Context, orderID string) (*OrderConfirmation, error) {
	var result OrderConfirmation
	err := c.post(ctx, "/v1/orders/"+orderID+"/confirm", nil, &result)
	return &result, err
}

func (c *Client) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	var result Order
	err := c.get(ctx, "/v1/orders/"+orderID, &result)
	return &result, err
}

func (c *Client) ListOrders(ctx context.Context, params *ListOrdersParams) (*OrderList, error) {
	q := map[string]string{
		"merchantId": params.MerchantID,
		"status":     params.Status,
		"type":       params.Type,
		"fromDate":   params.FromDate,
		"toDate":     params.ToDate,
	}
	if params.Page > 0 {
		q["page"] = strconv.Itoa(params.Page)
	}
	if params.Limit > 0 {
		q["limit"] = strconv.Itoa(params.Limit)
	}

	var result OrderList
	err := c.get(ctx, "/v1/orders"+buildQuery(q), &result)
	return &result, err
}

func (c *Client) UploadInvoiceDocuments(ctx context.Context, orderID string, fileIDs []string) (*InvoiceDocuments, error) {
	body := struct {
		FileIDs []string `json:"fileIds"`
	}{FileIDs: fileIDs}

	var result InvoiceDocuments
	err := c.post(ctx, "/v1/orders/"+orderID+"/invoice-documents", body, &result)
	return &result, err
}
