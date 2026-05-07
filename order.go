package velafi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (c *Client) CreateFiatToCryptoOrder(ctx context.Context, params *CreateFiatToCryptoOrderParams) (*CreateOrderResult, error) {
	var result CreateOrderResult
	err := c.post(ctx, "/v2/order/fiat_to_crypto", params, &result)
	return &result, err
}

func (c *Client) CreateCryptoToFiatOrder(ctx context.Context, params *CreateCryptoToFiatOrderParams) (*CreateOrderResult, error) {
	var result CreateOrderResult
	err := c.post(ctx, "/v2/order/crypto_to_fiat", params, &result)
	return &result, err
}

func (c *Client) GetOrderDetail(ctx context.Context, params *GetOrderDetailParams) (*Order, error) {
	path := "/v2/order/detail" + buildQuery(map[string]string{
		"orderId":   params.OrderID,
		"orderType": params.OrderType,
	})
	var result Order
	err := c.get(ctx, path, &result)
	return &result, err
}

func (c *Client) ConfirmOrder(ctx context.Context, params *ConfirmOrderParams) error {
	return c.post(ctx, "/v2/order/confirm", params, nil)
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

func (c *Client) UploadOrderInvoice(ctx context.Context, orderID string, orderType string, filePaths []string) error {
	if err := c.ensureToken(ctx); err != nil {
		return err
	}

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	if err := writer.WriteField("orderId", orderID); err != nil {
		return fmt.Errorf("velafi: write orderId: %w", err)
	}
	if err := writer.WriteField("orderType", orderType); err != nil {
		return fmt.Errorf("velafi: write orderType: %w", err)
	}

	for _, fp := range filePaths {
		f, err := os.Open(fp)
		if err != nil {
			return fmt.Errorf("velafi: open file %s: %w", fp, err)
		}
		part, err := writer.CreateFormFile("files", filepath.Base(fp))
		if err != nil {
			f.Close()
			return fmt.Errorf("velafi: create form file: %w", err)
		}
		if _, err := io.Copy(part, f); err != nil {
			f.Close()
			return fmt.Errorf("velafi: copy file: %w", err)
		}
		f.Close()
	}
	writer.Close()

	reqURL := c.baseURL + "/v2/order/invoice"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, &buf)
	if err != nil {
		return fmt.Errorf("velafi: create request: %w", err)
	}
	req.Header.Set("X-BH-TOKEN", c.token)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("velafi: read response: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return fmt.Errorf("velafi: unmarshal response: %w", err)
	}
	if resp.StatusCode >= 400 || apiResp.Code != 200 {
		return &Error{HTTPStatus: resp.StatusCode, Code: apiResp.Code, Message: apiResp.Msg}
	}
	return nil
}
