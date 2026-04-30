package velafi

import "context"

func (c *Client) ListPaymentTemplates(ctx context.Context, params *ListPaymentTemplatesParams) ([]PaymentTemplate, error) {
	path := "/v1/payment-templates" + buildQuery(map[string]string{
		"currency": params.Currency,
		"country":  params.Country,
		"type":     params.Type,
	})
	var result []PaymentTemplate
	err := c.get(ctx, path, &result)
	return result, err
}

func (c *Client) GetPaymentTemplateMetamessage(ctx context.Context, templateID string) (*PaymentTemplateMetamessage, error) {
	var result PaymentTemplateMetamessage
	err := c.get(ctx, "/v1/payment-templates/"+templateID+"/metamessage", &result)
	return &result, err
}

func (c *Client) AddPaymentMethod(ctx context.Context, params *AddPaymentMethodParams) (*PaymentMethod, error) {
	var result PaymentMethod
	err := c.post(ctx, "/v1/payment-methods", params, &result)
	return &result, err
}

func (c *Client) GetPaymentMethod(ctx context.Context, methodID string) (*PaymentMethod, error) {
	var result PaymentMethod
	err := c.get(ctx, "/v1/payment-methods/"+methodID, &result)
	return &result, err
}

func (c *Client) DeletePaymentMethod(ctx context.Context, methodID string) error {
	return c.delete(ctx, "/v1/payment-methods/"+methodID)
}

func (c *Client) SetRefundAccount(ctx context.Context, methodID string, refundMethodID string) (*RefundAccountResult, error) {
	body := struct {
		RefundMethodID string `json:"refundMethodId"`
	}{RefundMethodID: refundMethodID}

	var result RefundAccountResult
	err := c.post(ctx, "/v1/payment-methods/"+methodID+"/refund-account", body, &result)
	return &result, err
}
