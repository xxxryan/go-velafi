package velafi

import "encoding/json"

// Basic Configuration

type Country struct {
	Country string `json:"country"`
	Abbr    string `json:"abbr"`
}

type FiatCurrency struct {
	Country  string `json:"country"`
	Fiat     string `json:"fiat"`
	Accuracy int    `json:"accuracy"`
}

type CryptoCurrency struct {
	Crypto   string `json:"crypto"`
	Accuracy int    `json:"accuracy"`
}

type BuySellSymbol struct {
	Country  string `json:"country"`
	Fiat     string `json:"fiat"`
	Crypto   string `json:"crypto"`
	Accuracy int    `json:"accuracy"`
}

type FiatFiatSymbol struct {
	OnRampCountry  string `json:"onRampCountry"`
	OnRampFiat     string `json:"onRampFiat"`
	OffRampCountry string `json:"offRampCountry"`
	OffRampFiat    string `json:"offRampFiat"`
	Accuracy       int    `json:"accuracy"`
}

type PaymentOption struct {
	PaymentID   int    `json:"paymentId"`
	FiatFee     string `json:"fiatFee"`
	PaymentType int    `json:"paymentType"`
	Trench      string `json:"trench,omitempty"`
}

type BuySellPayments struct {
	Country     string          `json:"country"`
	Fiat        string          `json:"fiat"`
	Crypto      string          `json:"crypto"`
	PaymentList []PaymentOption `json:"paymentList"`
}

type FiatFiatPayments struct {
	OnRampCountry   string          `json:"onRampCountry"`
	OnRampFiat      string          `json:"onRampFiat"`
	OffRampCountry  string          `json:"offRampCountry"`
	OffRampFiat     string          `json:"offRampFiat"`
	PaymentListFrom []PaymentOption `json:"paymentListFrom"`
	PaymentListTo   []PaymentOption `json:"paymentListTo"`
}

// File Upload

type UploadFileParams struct {
	BusinessType string
	FilePaths    []string
}

type UploadedFile struct {
	FileName    string `json:"fileName"`
	FileType    string `json:"fileType"`
	FileURL     string `json:"fileUrl"`
	TempFileURL string `json:"tempFileUrl"`
}

// Quote

type CryptoQuoteParams struct {
	Country       string
	From          string
	To            string
	CreateQuoteID bool
}

type FiatQuoteParams struct {
	OnRampCountry  string
	OnRampFiat     string
	OffRampCountry string
	OffRampFiat    string
	CreateQuoteID  bool
}

type Quote struct {
	Price   string `json:"price"`
	QuoteID string `json:"quoteId,omitempty"`
}

// Order

type CreateFiatCryptoOrderParams struct {
	Country      string `json:"country"`
	Crypto       string `json:"crypto"`
	Fiat         string `json:"fiat"`
	FiatAmount   string `json:"fiatAmount"`
	PaymentID    int    `json:"paymentId"`
	ClientID     string `json:"clientId,omitempty"`
	MerchantID   int    `json:"merchantId,omitempty"`
	Remark       string `json:"remark,omitempty"`
	DepositAlias string `json:"depositAlias,omitempty"`
	QuoteID      string `json:"quoteId,omitempty"`
}

type CreateCryptoFiatOrderParams struct {
	Country       string `json:"country"`
	Crypto        string `json:"crypto"`
	Fiat          string `json:"fiat"`
	CryptoAmount  string `json:"cryptoAmount"`
	UserPaymentID int    `json:"userPaymentId"`
	ClientID      string `json:"clientId,omitempty"`
	MerchantID    int    `json:"merchantId,omitempty"`
	Remark        string `json:"remark,omitempty"`
	QuoteID       string `json:"quoteId,omitempty"`
}

type CreateFiatFiatOrderParams struct {
	ClientID          string `json:"clientId,omitempty"`
	OnRampCountry     string `json:"onRampCountry"`
	OnRampMerchantID  int    `json:"onRampMerchantId,omitempty"`
	OnRampFiat        string `json:"onRampFiat"`
	OnRampFiatAmount  string `json:"onRampFiatAmount"`
	OnRampPaymentID   int    `json:"onRampPaymentId"`
	DepositAlias      string `json:"depositAlias,omitempty"`
	OffRampCountry    string `json:"offRampCountry"`
	OffRampMerchantID int    `json:"offRampMerchantId,omitempty"`
	OffRampFiat       string `json:"offRampFiat"`
	OffRampPaymentID  int    `json:"offRampPaymentId"`
	Remark            string `json:"remark,omitempty"`
	QuoteID           string `json:"quoteId,omitempty"`
}

type CreateOrderResult struct {
	OrderID int64 `json:"orderId"`
}

type GetOrderParams struct {
	OrderID   int64
	OrderType string // fiat_to_crypto, crypto_to_fiat, fiat_to_fiat
}

type Order struct {
	OrderID       int64           `json:"orderId"`
	ClientID      string          `json:"clientId,omitempty"`
	MerchantID    int             `json:"merchantId,omitempty"`
	PaymentID     int             `json:"paymentId,omitempty"`
	UserPaymentID int             `json:"userPaymentId,omitempty"`
	Country       string          `json:"country,omitempty"`
	Crypto        string          `json:"crypto,omitempty"`
	Fiat          string          `json:"fiat,omitempty"`
	OrderType     string          `json:"orderType"`
	OrderPrice    string          `json:"orderPrice,omitempty"`
	CryptoAmount  string          `json:"cryptoAmount,omitempty"`
	FiatAmount    string          `json:"fiatAmount,omitempty"`
	FiatFee       string          `json:"fiatFee,omitempty"`
	OrderStatus   int             `json:"orderStatus"`
	TraceNumber   string          `json:"traceNumber,omitempty"`
	PaymentInfo   json.RawMessage `json:"paymentInfo,omitempty"`
	FailCode      string          `json:"failCode,omitempty"`
	FailReason    string          `json:"failReason,omitempty"`
	CreateTime    int64           `json:"createTime,omitempty"`
	CompletedTime int64           `json:"completedTime,omitempty"`
}

type ListOrdersParams struct {
	CurrentPage int
	PageSize    int
	StartTime   string
	EndTime     string
	OrderType   string
	OrderStatus int
}

type OrderList struct {
	Size        int               `json:"size"`
	CurrentPage int               `json:"currentPage"`
	Total       int               `json:"total"`
	Record      []json.RawMessage `json:"record"`
}

// Payment Method

type PaymentTemplate map[string]string

type PaymentTemplateMetaField struct {
	Index      int    `json:"index"`
	IndexCode  string `json:"indexCode"`
	TextType   int    `json:"textType"`
	Title      string `json:"title"`
	PromptText string `json:"promptText"`
	MinLimit   int    `json:"minLimit"`
	MaxLimit   int    `json:"maxLimit"`
	IsOptional bool   `json:"isOptional"`
	IsAccount  bool   `json:"isAccount"`
	ExtendInfo string `json:"extendInfo,omitempty"`
}

type PaymentTemplateMeta struct {
	ID          int                        `json:"id"`
	PaymentName string                     `json:"paymentName"`
	PaymentType int                        `json:"paymentType"`
	Trench      string                     `json:"trench"`
	FieldList   []PaymentTemplateMetaField `json:"fieldList"`
}

type AddPaymentMethodParams struct {
	MerchantID int            `json:"merchantId"`
	PaymentID  int            `json:"paymentId"`
	Country    string         `json:"country"`
	Fiat       string         `json:"fiat"`
	RealName   string         `json:"realName"`
	FieldJSON  map[string]any `json:"fieldJson"`
	Remark     string         `json:"remark,omitempty"`
}

type AddPaymentMethodResult struct {
	ID         int    `json:"id"`
	Status     int    `json:"status"`
	FailReason string `json:"failReason,omitempty"`
}

type ListPaymentMethodsParams struct {
	Country     string
	Status      int
	Fiat        string
	MerchantID  int
	CurrentPage int
	PageSize    int
}

type PaymentMethod struct {
	ID                int               `json:"id"`
	MerchantID        int               `json:"merchantId"`
	Country           string            `json:"country"`
	Fiat              string            `json:"fiat"`
	PaymentID         int               `json:"paymentId"`
	PaymentMethodName string            `json:"paymentMethodName"`
	RealName          string            `json:"realName"`
	Status            int               `json:"status"`
	HasRefundAccount  int               `json:"hasRefundAccount"`
	FieldList         map[string]string `json:"fieldList"`
	Remark            string            `json:"remark,omitempty"`
	CreateTime        int64             `json:"createTime"`
}

type PaymentMethodList struct {
	CurrentPage int             `json:"currentPage"`
	Size        int             `json:"size"`
	Total       int             `json:"total"`
	Record      []PaymentMethod `json:"record"`
}

// Webhook

type CreateWebhookParams struct {
	EventType string `json:"eventType"`
	URL       string `json:"url"`
}

type Webhook struct {
	WebhookID string `json:"webhookId"`
	EventType string `json:"eventType"`
	URL       string `json:"url"`
	Status    int    `json:"status"`
	PublicKey string `json:"publicKey,omitempty"`
}
