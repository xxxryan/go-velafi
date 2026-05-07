package velafi

import "encoding/json"

// FlexString handles JSON fields that may be either a string or a number.
type FlexString string

func (f *FlexString) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || string(data) == "null" {
		return nil
	}
	if data[0] == '"' {
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		*f = FlexString(s)
		return nil
	}
	*f = FlexString(data)
	return nil
}

func (f FlexString) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(f))
}

func (f FlexString) String() string {
	return string(f)
}

// Basic Configuration

type BuySellSymbol struct {
	Country  string `json:"country"`
	Fiat     string `json:"fiat"`
	Crypto   string `json:"crypto"`
	Accuracy int    `json:"accuracy"`
}

type PaymentOption struct {
	PaymentID   int     `json:"paymentId"`
	FiatFee     float64 `json:"fiatFee"`
	PaymentType int     `json:"paymentType"`
	Trench      string  `json:"trench,omitempty"`
}

type BuySellPayments struct {
	Country     string          `json:"country"`
	Fiat        string          `json:"fiat"`
	Crypto      string          `json:"crypto"`
	PaymentList []PaymentOption `json:"paymentList"`
}

// Quote

type CryptoQuoteParams struct {
	Country       string
	From          string
	To            string
	CreateQuoteID bool
}

type Quote struct {
	Price   string `json:"price"`
	QuoteID string `json:"quoteId,omitempty"`
}

// Payment Method

type PaymentTemplate map[string]string

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

// Order

type CreateFiatToCryptoOrderParams struct {
	Country      string `json:"country"`
	Crypto       string `json:"crypto"`
	Fiat         string `json:"fiat"`
	FiatAmount   string `json:"fiatAmount"`
	PaymentID    int    `json:"paymentId"`
	MerchantID   int    `json:"merchantId,omitempty"`
	ClientID     string `json:"clientId,omitempty"`
	Remark       string `json:"remark,omitempty"`
	DepositAlias string `json:"depositAlias,omitempty"`
	QuoteID      string `json:"quoteId,omitempty"`
}

type CreateCryptoToFiatOrderParams struct {
	Country       string `json:"country"`
	Crypto        string `json:"crypto"`
	Fiat          string `json:"fiat"`
	CryptoAmount  string `json:"cryptoAmount"`
	UserPaymentID int    `json:"userPaymentId"`
	MerchantID    int    `json:"merchantId,omitempty"`
	ClientID      string `json:"clientId,omitempty"`
	Remark        string `json:"remark,omitempty"`
	QuoteID       string `json:"quoteId,omitempty"`
}

type CreateOrderResult struct {
	OrderID FlexString `json:"orderId"`
}

type GetOrderDetailParams struct {
	OrderID   string
	OrderType string
}

type Order struct {
	OrderID       FlexString      `json:"orderId"`
	ClientID      string          `json:"clientId"`
	MerchantID    int             `json:"merchantId"`
	PaymentID     int             `json:"paymentId"`
	UserPaymentID int             `json:"userPaymentId"`
	Country       string          `json:"country"`
	Crypto        string          `json:"crypto"`
	Fiat          string          `json:"fiat"`
	OrderType     string          `json:"orderType"`
	OrderPrice    string          `json:"orderPrice"`
	CryptoAmount  string          `json:"cryptoAmount"`
	FiatAmount    string          `json:"fiatAmount"`
	FiatFee       string          `json:"fiatFee"`
	OrderStatus   int             `json:"orderStatus"`
	HasRefund     int             `json:"hasRefund"`
	TraceNumber   string          `json:"traceNumber"`
	PaymentInfo   json.RawMessage `json:"paymentInfo"`
	FailCode      int             `json:"failCode"`
	FailReason    string          `json:"failReason"`
	CreateTime    string          `json:"createTime"`
	CompletedTime string          `json:"completedTime"`
}

type ConfirmOrderParams struct {
	OrderID   int64  `json:"orderId"`
	OrderType string `json:"orderType"`
	Direction string `json:"direction,omitempty"`
}

// Payment Method (additional types)

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

type SetRefundAccountParams struct {
	UserPaymentID int `json:"userPaymentId"`
	MerchantID    int `json:"merchantId"`
}

// Order (additional types)

type ListOrdersParams struct {
	CurrentPage int
	PageSize    int
	StartTime   string
	EndTime     string
	OrderType   string
	OrderStatus int
}

type OrderListItem struct {
	OrderID       FlexString `json:"orderId"`
	ClientID      string     `json:"clientId"`
	MerchantID    int        `json:"merchantId"`
	PaymentID     int        `json:"paymentId"`
	Country       string     `json:"country"`
	Crypto        string     `json:"crypto"`
	Fiat          string     `json:"fiat"`
	OrderType     string     `json:"orderType"`
	OrderPrice    string     `json:"orderPrice"`
	CryptoAmount  string     `json:"cryptoAmount"`
	FiatAmount    string     `json:"fiatAmount"`
	FiatFee       string     `json:"fiatFee"`
	OrderStatus   int        `json:"orderStatus"`
	HasRefund     int        `json:"hasRefund"`
	CreateTime    string     `json:"createTime"`
	CompletedTime string     `json:"completedTime"`
}

type OrderList struct {
	Size        int             `json:"size"`
	CurrentPage int             `json:"currentPage"`
	Total       int             `json:"total"`
	Record      []OrderListItem `json:"record"`
}

// Account

type AccountDetails struct {
	UserID      int     `json:"userId"`
	Country     string  `json:"country"`
	CompanyName string  `json:"companyName"`
	KYCPassed   bool    `json:"kycPassed"`
	MonthLimit  float64 `json:"monthLimit"`
	MonthUsed   float64 `json:"monthUsed"`
}

// Merchant

type MerchantAccount struct {
	MerchantID        int    `json:"merchantId"`
	Fiat              string `json:"fiat"`
	Balance           string `json:"balance"`
	Status            int    `json:"status"`
	PaymentID         int    `json:"paymentId"`
	PaymentMethodName string `json:"paymentMethodName"`
	CreateTime        string `json:"createTime"`
	UpdateTime        string `json:"updateTime"`
	MonthBuyLimit     string `json:"monthBuyLimit"`
	MonthBuyUsed      string `json:"monthBuyUsed"`
	MonthSellLimit    string `json:"monthSellLimit"`
	MonthSellUsed     string `json:"monthSellUsed"`
	FailReason        string `json:"failReason,omitempty"`
}

type MerchantAccountList struct {
	CurrentPage int               `json:"currentPage"`
	Size        int               `json:"size"`
	Total       int               `json:"total"`
	Record      []MerchantAccount `json:"record"`
}

type ListMerchantAccountsParams struct {
	MerchantID   int
	MerchantName string
	Email        string
	Fiat         string
	Status       int
	CurrentPage  int
	PageSize     int
}

type ActivateMerchantAccountParams struct {
	MerchantID  int               `json:"merchantId"`
	Fiat        string            `json:"fiat"`
	Trench      string            `json:"trench"`
	CallbackURI string            `json:"callbackUri,omitempty"`
	FieldList   map[string]string `json:"fieldList,omitempty"`
}

type ActivateMerchantAccountResult struct {
	Fiat       string `json:"fiat"`
	Status     int    `json:"status"`
	VerifyLink string `json:"verifyLink,omitempty"`
	FailReason string `json:"failReason,omitempty"`
}

type GetPendingFundAccountParams struct {
	MerchantID   int
	PaymentID    int
	Fiat         string
	DepositAlias string
	Amount       string
}

type PendingFundAccount struct {
	Fiat              string              `json:"fiat"`
	PaymentMethodName string              `json:"paymentMethodName"`
	TxID              int64               `json:"txId,omitempty"`
	FieldList         map[string][]string `json:"fieldList"`
}

type ClaimPendingFundParams struct {
	MerchantID    int    `json:"merchantId"`
	Fiat          string `json:"fiat"`
	Amount        string `json:"amount"`
	UserPaymentID int    `json:"userPaymentId"`
	ClientID      string `json:"clientId,omitempty"`
}

type ClaimPendingFundResult struct {
	TxID          int64   `json:"txId"`
	MerchantID    int     `json:"merchantId"`
	ClientID      string  `json:"clientId"`
	Fiat          string  `json:"fiat"`
	TotalAmount   float64 `json:"totalAmount"`
	Amount        float64 `json:"amount"`
	Fee           float64 `json:"fee"`
	UserPaymentID int     `json:"userPaymentId"`
	Status        int     `json:"status"`
	CreateTime    string  `json:"createTime"`
}

type GetFundingRecordsParams struct {
	TxID        int64
	StartTime   int64
	EndTime     int64
	MerchantID  int
	Fiat        string
	Type        string // ALL, DEPOSIT, WITHDRAW, BUY_CRYPTO, REFUND, TRANSFER_IN, TRANSFER_OUT
	Status      int    // 0=all, 1=pending, 2=completed, 3=canceled, 4=refunded
	ClientID    string
	CurrentPage int
	PageSize    int
}

type FundingRecord struct {
	ID            string  `json:"id"`
	MerchantID    int     `json:"merchantId"`
	MerchantName  string  `json:"merchantName"`
	Fiat          string  `json:"fiat"`
	TotalAmount   float64 `json:"totalAmount"`
	Amount        float64 `json:"amount"`
	Fee           float64 `json:"fee"`
	UserPaymentID int     `json:"userPaymentId"`
	Type          string  `json:"type"`
	Status        int     `json:"status"`
	CreateTime    string  `json:"createTime"`
	UpdateTime    string  `json:"updateTime"`
}

type FundingRecordList struct {
	CurrentPage int             `json:"currentPage"`
	Size        int             `json:"size"`
	Total       int             `json:"total"`
	Data        []FundingRecord `json:"data"`
}

type GetFiatFlowParams struct {
	StartTime   int64
	EndTime     int64
	MerchantID  int
	Fiat        string
	TxID        int64
	FlowType    string // ALL, DEPOSIT, WITHDRAW, BUY_CRYPTO, REFUND, etc.
	CurrentPage int
	PageSize    int
}

type FiatFlowRecord struct {
	ID         int64   `json:"id"`
	MerchantID int     `json:"merchantId"`
	TxID       string  `json:"txId"`
	Category   string  `json:"category"`
	FlowType   string  `json:"flowType"`
	Fiat       string  `json:"fiat"`
	Changed    float64 `json:"changed"`
	Total      float64 `json:"total"`
	ExtID      string  `json:"extId"`
	CreateTime string  `json:"createTime"`
}

type FiatFlowList struct {
	CurrentPage int              `json:"currentPage"`
	Size        int              `json:"size"`
	Total       int              `json:"total"`
	Data        []FiatFlowRecord `json:"data"`
}

type CryptoBalance struct {
	Asset     string `json:"asset"`
	AssetID   string `json:"assetId"`
	AssetName string `json:"assetName"`
	Total     string `json:"total"`
	Free      string `json:"free"`
	Locked    string `json:"locked"`
}

type CryptoBalanceResult struct {
	CanDeposit bool            `json:"canDeposit"`
	Balances   []CryptoBalance `json:"balances"`
}

type GetCryptoDepositAddressParams struct {
	TokenID    string
	MerchantID int
	ChainType  string
}

type CryptoDepositAddress struct {
	AllowDeposit          bool   `json:"allowDeposit"`
	Address               string `json:"address"`
	AddressExt            string `json:"addressExt"`
	MinQuantity           string `json:"minQuantity"`
	NeedAddressTag        bool   `json:"needAddressTag"`
	RequiredConfirmNum    int    `json:"requiredConfirmNum"`
	CanWithdrawConfirmNum int    `json:"canWithdrawConfirmNum"`
	TokenType             string `json:"tokenType"`
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

type UpdateWebhookParams struct {
	URL    string `json:"url"`
	Status int    `json:"status"`
}

// OrderStatus represents the status of an order.
type OrderStatus int

const (
	OrderStatusReviewing      OrderStatus = 10
	OrderStatusRFIRequested   OrderStatus = 11
	OrderStatusRFIUploaded    OrderStatus = 12
	OrderStatusApproved       OrderStatus = 30
	OrderStatusDocsPending    OrderStatus = 31
	OrderStatusPendingPayment OrderStatus = 40
	OrderStatusPaid           OrderStatus = 50
	OrderStatusReleased       OrderStatus = 60
	OrderStatusCanceled       OrderStatus = 70
	OrderStatusRefunded       OrderStatus = 72
	OrderStatusRefunding      OrderStatus = 73
)

func (s OrderStatus) String() string {
	names := map[OrderStatus]string{
		OrderStatusReviewing:      "reviewing",
		OrderStatusRFIRequested:   "RFI requested",
		OrderStatusRFIUploaded:    "RFI uploaded",
		OrderStatusApproved:       "approved",
		OrderStatusDocsPending:    "docs pending",
		OrderStatusPendingPayment: "pending payment",
		OrderStatusPaid:           "paid",
		OrderStatusReleased:       "released",
		OrderStatusCanceled:       "canceled",
		OrderStatusRefunded:       "refunded",
		OrderStatusRefunding:      "refunding",
	}
	if name, ok := names[s]; ok {
		return name
	}
	return "unknown"
}
