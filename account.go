package velafi

import "context"

func (c *Client) GetAccountDetails(ctx context.Context) (*AccountDetails, error) {
	var result AccountDetails
	err := c.get(ctx, "/v2/user/account", &result)
	return &result, err
}
