package velafi

import "encoding/json"

// Basic Configuration

type Country struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Currency struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol,omitempty"`
}

type Pair struct {
	FromCurrency string `json:"fromCurrency"`
	ToCurrency   string `json:"toCurrency"`
}

type PaymentMethodInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Currency string `json:"currency"`
	Country  string `json:"country"`
	Type     string `json:"type"`
}

// File

type UploadFileParams struct {
	FilePath string
	Purpose  string // "id_document", "selfie", "business_document"
}

type File struct {
	FileID    string `json:"fileId"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
	Size      int64  `json:"size"`
	CreatedAt string `json:"createdAt"`
}

// Quote

type CryptoFiatQuoteParams struct {
	MerchantID      string `json:"merchantId"`
	FromCurrency    string `json:"fromCurrency"`
	ToCurrency      string `json:"toCurrency"`
	FromAmount      string `json:"fromAmount,omitempty"`
	ToAmount        string `json:"toAmount,omitempty"`
	PaymentMethodID string `json:"paymentMethodId"`
}

type FiatFiatQuoteParams struct {
	MerchantID          string `json:"merchantId"`
	FromCurrency        string `json:"fromCurrency"`
	ToCurrency          string `json:"toCurrency"`
	FromAmount          string `json:"fromAmount,omitempty"`
	ToAmount            string `json:"toAmount,omitempty"`
	FromPaymentMethodID string `json:"fromPaymentMethodId"`
	ToPaymentMethodID   string `json:"toPaymentMethodId"`
}

type Quote struct {
	QuoteID      string `json:"quoteId"`
	FromCurrency string `json:"fromCurrency"`
	ToCurrency   string `json:"toCurrency"`
	FromAmount   string `json:"fromAmount"`
	ToAmount     string `json:"toAmount"`
	ExchangeRate string `json:"exchangeRate"`
	Fee          string `json:"fee"`
	ExpiredAt    string `json:"expiredAt"`
}

// Order

type CreateFiatCryptoOrderParams struct {
	MerchantID      string `json:"merchantId"`
	FromCurrency    string `json:"fromCurrency"`
	ToCurrency      string `json:"toCurrency"`
	FromAmount      string `json:"fromAmount,omitempty"`
	ToAmount        string `json:"toAmount,omitempty"`
	PaymentMethodID string `json:"paymentMethodId"`
	ToAddress       string `json:"toAddress"`
	Network         string `json:"network"`
	QuoteID         string `json:"quoteId,omitempty"`
	ClientOrderID   string `json:"clientOrderId,omitempty"`
	Memo            string `json:"memo,omitempty"`
}

type CreateCryptoFiatOrderParams struct {
	MerchantID      string `json:"merchantId"`
	FromCurrency    string `json:"fromCurrency"`
	ToCurrency      string `json:"toCurrency"`
	FromAmount      string `json:"fromAmount,omitempty"`
	ToAmount        string `json:"toAmount,omitempty"`
	Network         string `json:"network"`
	PaymentMethodID string `json:"paymentMethodId"`
	QuoteID         string `json:"quoteId,omitempty"`
	ClientOrderID   string `json:"clientOrderId,omitempty"`
}

type CreateFiatFiatOrderParams struct {
	MerchantID          string `json:"merchantId"`
	FromCurrency        string `json:"fromCurrency"`
	ToCurrency          string `json:"toCurrency"`
	FromAmount          string `json:"fromAmount,omitempty"`
	ToAmount            string `json:"toAmount,omitempty"`
	FromPaymentMethodID string `json:"fromPaymentMethodId"`
	ToPaymentMethodID   string `json:"toPaymentMethodId"`
	QuoteID             string `json:"quoteId,omitempty"`
	ClientOrderID       string `json:"clientOrderId,omitempty"`
}

type Order struct {
	OrderID             string `json:"orderId"`
	ClientOrderID       string `json:"clientOrderId,omitempty"`
	MerchantID          string `json:"merchantId,omitempty"`
	Type                string `json:"type,omitempty"`
	Status              string `json:"status"`
	FromCurrency        string `json:"fromCurrency"`
	ToCurrency          string `json:"toCurrency"`
	FromAmount          string `json:"fromAmount"`
	ToAmount            string `json:"toAmount"`
	ExchangeRate        string `json:"exchangeRate,omitempty"`
	Fee                 string `json:"fee,omitempty"`
	PaymentMethodID     string `json:"paymentMethodId,omitempty"`
	ToAddress           string `json:"toAddress,omitempty"`
	DepositAddress      string `json:"depositAddress,omitempty"`
	Network             string `json:"network,omitempty"`
	FromPaymentMethodID string `json:"fromPaymentMethodId,omitempty"`
	ToPaymentMethodID   string `json:"toPaymentMethodId,omitempty"`
	CreatedAt           string `json:"createdAt,omitempty"`
	UpdatedAt           string `json:"updatedAt,omitempty"`
}

type OrderConfirmation struct {
	OrderID     string `json:"orderId"`
	Status      string `json:"status"`
	ConfirmedAt string `json:"confirmedAt"`
}

type ListOrdersParams struct {
	MerchantID string `json:"merchantId"`
	Status     string `json:"status,omitempty"`
	Type       string `json:"type,omitempty"`
	FromDate   string `json:"fromDate,omitempty"`
	ToDate     string `json:"toDate,omitempty"`
	Page       int    `json:"page,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}

type OrderList struct {
	Total int      `json:"total"`
	Page  int      `json:"page"`
	Limit int      `json:"limit"`
	Items []*Order `json:"items"`
}

type InvoiceDocuments struct {
	OrderID   string            `json:"orderId"`
	Documents []InvoiceDocument `json:"documents"`
}

type InvoiceDocument struct {
	FileID     string `json:"fileId"`
	Filename   string `json:"filename"`
	UploadedAt string `json:"uploadedAt"`
}

// Payment Method

type ListPaymentTemplatesParams struct {
	Currency string `json:"currency,omitempty"`
	Country  string `json:"country,omitempty"`
	Type     string `json:"type,omitempty"`
}

type PaymentTemplate struct {
	TemplateID string                 `json:"templateId"`
	Name       string                 `json:"name"`
	Currency   string                 `json:"currency"`
	Country    string                 `json:"country"`
	Type       string                 `json:"type"`
	Fields     []PaymentTemplateField `json:"fields"`
}

type PaymentTemplateField struct {
	Key      string `json:"key"`
	Label    string `json:"label"`
	Type     string `json:"type"`
	Required bool   `json:"required"`
}

type PaymentTemplateMetamessage struct {
	TemplateID  string `json:"templateId"`
	Metamessage struct {
		Fields []PaymentTemplateMetaField `json:"fields"`
	} `json:"metamessage"`
}

type PaymentTemplateMetaField struct {
	Key        string `json:"key"`
	Label      string `json:"label"`
	Type       string `json:"type"`
	Required   bool   `json:"required"`
	Validation *struct {
		Pattern   string `json:"pattern,omitempty"`
		MinLength int    `json:"minLength,omitempty"`
		MaxLength int    `json:"maxLength,omitempty"`
	} `json:"validation,omitempty"`
}

type AddPaymentMethodParams struct {
	MerchantID string         `json:"merchantId"`
	TemplateID string         `json:"templateId"`
	Fields     map[string]any `json:"fields"`
}

type PaymentMethod struct {
	MethodID       string          `json:"methodId"`
	MerchantID     string          `json:"merchantId"`
	TemplateID     string          `json:"templateId"`
	Status         string          `json:"status"`
	Fields         json.RawMessage `json:"fields"`
	RefundMethodID string          `json:"refundMethodId,omitempty"`
	CreatedAt      string          `json:"createdAt"`
}

type RefundAccountResult struct {
	MethodID       string `json:"methodId"`
	RefundMethodID string `json:"refundMethodId"`
	UpdatedAt      string `json:"updatedAt"`
}

// Webhook

type CreateWebhookParams struct {
	URL        string   `json:"url"`
	Events     []string `json:"events"`
	Secret     string   `json:"secret,omitempty"`
	MerchantID string   `json:"merchantId"`
}

type UpdateWebhookParams struct {
	URL    string   `json:"url,omitempty"`
	Events []string `json:"events,omitempty"`
	Secret string   `json:"secret,omitempty"`
	Status string   `json:"status,omitempty"`
}

type Webhook struct {
	WebhookID  string   `json:"webhookId"`
	URL        string   `json:"url"`
	Events     []string `json:"events"`
	MerchantID string   `json:"merchantId"`
	Status     string   `json:"status"`
	CreatedAt  string   `json:"createdAt,omitempty"`
	UpdatedAt  string   `json:"updatedAt,omitempty"`
}

// Webhook Events (inbound payloads)

type FiatCryptoOrderEvent struct {
	OrderID       string `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
	MerchantID    string `json:"merchantId"`
	Type          string `json:"type"`
	Status        string `json:"status"`
	FromCurrency  string `json:"fromCurrency"`
	ToCurrency    string `json:"toCurrency"`
	FromAmount    string `json:"fromAmount"`
	ToAmount      string `json:"toAmount"`
	ExchangeRate  string `json:"exchangeRate"`
	Fee           string `json:"fee"`
	ToAddress     string `json:"toAddress"`
	TxHash        string `json:"txHash"`
	Network       string `json:"network"`
	CreatedAt     string `json:"createdAt"`
	CompletedAt   string `json:"completedAt"`
}

type FiatFiatOrderEvent struct {
	OrderID             string `json:"orderId"`
	ClientOrderID       string `json:"clientOrderId"`
	MerchantID          string `json:"merchantId"`
	Type                string `json:"type"`
	Status              string `json:"status"`
	FromCurrency        string `json:"fromCurrency"`
	ToCurrency          string `json:"toCurrency"`
	FromAmount          string `json:"fromAmount"`
	ToAmount            string `json:"toAmount"`
	ExchangeRate        string `json:"exchangeRate"`
	Fee                 string `json:"fee"`
	FromPaymentMethodID string `json:"fromPaymentMethodId"`
	ToPaymentMethodID   string `json:"toPaymentMethodId"`
	CreatedAt           string `json:"createdAt"`
	CompletedAt         string `json:"completedAt"`
}

type FundingRecordEvent struct {
	FundingID     string `json:"fundingId"`
	MerchantID    string `json:"merchantId"`
	Currency      string `json:"currency"`
	Amount        string `json:"amount"`
	Status        string `json:"status"`
	AccountNumber string `json:"accountNumber"`
	Reference     string `json:"reference"`
	CreatedAt     string `json:"createdAt"`
}

type StablecoinPaymentEvent struct {
	PaymentID   string `json:"paymentId"`
	MerchantID  string `json:"merchantId"`
	Currency    string `json:"currency"`
	Network     string `json:"network"`
	Amount      string `json:"amount"`
	FromAddress string `json:"fromAddress"`
	ToAddress   string `json:"toAddress"`
	TxHash      string `json:"txHash"`
	Status      string `json:"status"`
	CreatedAt   string `json:"createdAt"`
	CompletedAt string `json:"completedAt"`
}
