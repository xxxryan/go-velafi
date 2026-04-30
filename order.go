package velafi

import (
	"context"
	"strconv"
)

func (c *Client) CreateFiatCryptoOrder(ctx context.Context, params *CreateFiatCryptoOrderParams) (*CreateOrderResult, error) {
	var result CreateOrderResult
	err := c.post(ctx, "/v2/order/fiat_to_crypto", params, &result)
	return &result, err
}

func (c *Client) CreateCryptoFiatOrder(ctx context.Context, params *CreateCryptoFiatOrderParams) (*CreateOrderResult, error) {
	var result CreateOrderResult
	err := c.post(ctx, "/v2/order/crypto_to_fiat", params, &result)
	return &result, err
}

func (c *Client) CreateFiatFiatOrder(ctx context.Context, params *CreateFiatFiatOrderParams) (*CreateOrderResult, error) {
	var result CreateOrderResult
	err := c.post(ctx, "/v2/order/fiat_to_fiat", params, &result)
	return &result, err
}

func (c *Client) GetOrder(ctx context.Context, params *GetOrderParams) (*Order, error) {
	path := "/v2/order/detail" + buildQuery(map[string]string{
		"orderId":   strconv.FormatInt(params.OrderID, 10),
		"orderType": params.OrderType,
	})
	var result Order
	err := c.get(ctx, path, &result)
	return &result, err
}

func (c *Client) ListOrders(ctx context.Context, params *ListOrdersParams) (*OrderList, error) {
	q := map[string]string{
		"orderType": params.OrderType,
		"startTime": params.StartTime,
		"endTime":   params.EndTime,
	}
	if params.CurrentPage > 0 {
		q["currentPage"] = strconv.Itoa(params.CurrentPage)
	}
	if params.PageSize > 0 {
		q["pageSize"] = strconv.Itoa(params.PageSize)
	}
	if params.OrderStatus > 0 {
		q["orderStatus"] = strconv.Itoa(params.OrderStatus)
	}
	var result OrderList
	err := c.get(ctx, "/v2/orders"+buildQuery(q), &result)
	return &result, err
}
