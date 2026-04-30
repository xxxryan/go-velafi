package velafi

import "context"

func (c *Client) CreateWebhook(ctx context.Context, params *CreateWebhookParams) (*Webhook, error) {
	var result Webhook
	err := c.post(ctx, "/v1/webhooks", params, &result)
	return &result, err
}

func (c *Client) ListWebhooks(ctx context.Context, merchantID string) ([]Webhook, error) {
	path := "/v1/webhooks" + buildQuery(map[string]string{"merchantId": merchantID})
	var result []Webhook
	err := c.get(ctx, path, &result)
	return result, err
}

func (c *Client) UpdateWebhook(ctx context.Context, webhookID string, params *UpdateWebhookParams) (*Webhook, error) {
	var result Webhook
	err := c.put(ctx, "/v1/webhooks/"+webhookID, params, &result)
	return &result, err
}

func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	return c.delete(ctx, "/v1/webhooks/"+webhookID)
}
