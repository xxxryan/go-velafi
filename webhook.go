package velafi

import (
	"context"
	"strconv"
)

func (c *Client) CreateWebhook(ctx context.Context, params *CreateWebhookParams) (*Webhook, error) {
	var result Webhook
	err := c.post(ctx, "/v2/webhook", params, &result)
	return &result, err
}

func (c *Client) ListWebhooks(ctx context.Context, status int) ([]Webhook, error) {
	q := map[string]string{}
	if status > 0 {
		q["status"] = strconv.Itoa(status)
	}
	var result []Webhook
	err := c.get(ctx, "/v2/webhooks"+buildQuery(q), &result)
	return result, err
}
