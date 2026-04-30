package velafi

import (
	"net/http"
	"time"
)

type Option func(*Client)

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func WithSandbox() Option {
	return func(c *Client) {
		c.baseURL = SandboxBaseURL
	}
}

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

func WithTokenRefreshBuffer(d time.Duration) Option {
	return func(c *Client) {
		c.tokenRefreshBuffer = d
	}
}
