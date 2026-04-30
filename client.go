package velafi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const defaultBaseURL = "https://api.velafi.com"

type Client struct {
	apiKey    string
	apiSecret string
	baseURL   string
	httpClient *http.Client

	mu                 sync.Mutex
	token              string
	expiresAt          time.Time
	tokenRefreshBuffer time.Duration
}

func NewClient(apiKey, apiSecret string, opts ...Option) *Client {
	c := &Client{
		apiKey:             apiKey,
		apiSecret:          apiSecret,
		baseURL:            defaultBaseURL,
		httpClient:         &http.Client{Timeout: 30 * time.Second},
		tokenRefreshBuffer: 5 * time.Minute,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

type apiResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

func (c *Client) do(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}

	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-BH-TOKEN", c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.httpClient.Do(req)
}

func (c *Client) doJSON(ctx context.Context, method, path string, body any, result any) error {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("velafi: marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	resp, err := c.do(ctx, method, path, bodyReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("velafi: read response body: %w", err)
	}

	var apiResp apiResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		if resp.StatusCode >= 400 {
			return &Error{
				HTTPStatus: resp.StatusCode,
				Message:    http.StatusText(resp.StatusCode),
			}
		}
		return fmt.Errorf("velafi: unmarshal response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return &Error{
			HTTPStatus: resp.StatusCode,
			Code:       apiResp.Code,
			Message:    apiResp.Message,
			RequestID:  resp.Header.Get("X-Request-Id"),
		}
	}

	if apiResp.Code != 0 {
		return &Error{
			HTTPStatus: resp.StatusCode,
			Code:       apiResp.Code,
			Message:    apiResp.Message,
			RequestID:  resp.Header.Get("X-Request-Id"),
		}
	}

	if result != nil && len(apiResp.Data) > 0 {
		if err := json.Unmarshal(apiResp.Data, result); err != nil {
			return fmt.Errorf("velafi: unmarshal data: %w", err)
		}
	}

	return nil
}

func (c *Client) get(ctx context.Context, path string, result any) error {
	return c.doJSON(ctx, http.MethodGet, path, nil, result)
}

func (c *Client) post(ctx context.Context, path string, body any, result any) error {
	return c.doJSON(ctx, http.MethodPost, path, body, result)
}

func (c *Client) put(ctx context.Context, path string, body any, result any) error {
	return c.doJSON(ctx, http.MethodPut, path, body, result)
}

func (c *Client) delete(ctx context.Context, path string) error {
	return c.doJSON(ctx, http.MethodDelete, path, nil, nil)
}
