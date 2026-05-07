# go-velafi

Go SDK for the [VelaFi](https://docs.velafi.com) v2 API — on-ramp/off-ramp cryptocurrency operations.

[![Go Reference](https://pkg.go.dev/badge/github.com/xfinancial/go-velafi.svg)](https://pkg.go.dev/github.com/xfinancial/go-velafi)

## Features

- Zero external dependencies (pure Go standard library)
- Automatic token management (HMAC-SHA256 signing, lazy init, auto-refresh)
- Full MXN ↔ USDT on-ramp/off-ramp support
- Sandbox and production environments
- Structured error handling with `errors.As` support

## Installation

```bash
go get github.com/xfinancial/go-velafi
```

Requires Go 1.22+

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    velafi "github.com/xfinancial/go-velafi"
)

func main() {
    client := velafi.NewClient("your-api-key", "your-api-secret",
        velafi.WithSandbox(), // use sandbox for testing
    )
    ctx := context.Background()

    // Get MXN/USDT quote
    quote, err := client.GetCryptoQuote(ctx, &velafi.CryptoQuoteParams{
        Country: "Mexico", From: "MXN", To: "USDT", CreateQuoteID: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("MXN/USDT price: %s, quoteId: %s\n", quote.Price, quote.QuoteID)
}
```

## Configuration

```go
// Production (default)
client := velafi.NewClient(apiKey, apiSecret)

// Sandbox
client := velafi.NewClient(apiKey, apiSecret, velafi.WithSandbox())

// Custom options
client := velafi.NewClient(apiKey, apiSecret,
    velafi.WithBaseURL("https://custom.api.velafi.com"),
    velafi.WithHTTPClient(myHTTPClient),
    velafi.WithTokenRefreshBuffer(10 * time.Minute),
)
```

## API Reference

### Base Configuration
| Method | Description |
|--------|-------------|
| `GetBuySymbols(ctx)` | List on-ramp trading pairs (fiat→crypto) |
| `GetSellSymbols(ctx)` | List off-ramp trading pairs (crypto→fiat) |
| `GetBuyPayments(ctx, country, fiat, crypto)` | Get available payment methods for on-ramp |
| `GetSellPayments(ctx, country, fiat, crypto)` | Get available payment methods for off-ramp |

### Quote
| Method | Description |
|--------|-------------|
| `GetCryptoQuote(ctx, params)` | Get price quote (quoteId valid for 10 seconds) |

### Order
| Method | Description |
|--------|-------------|
| `CreateFiatToCryptoOrder(ctx, params)` | Create on-ramp order |
| `CreateCryptoToFiatOrder(ctx, params)` | Create off-ramp order |
| `GetOrderDetail(ctx, params)` | Query single order status |
| `ListOrders(ctx, params)` | Query order list with filters |
| `ConfirmOrder(ctx, params)` | Confirm order (accelerate off-ramp release) |
| `UploadOrderInvoice(ctx, orderID, orderType, files)` | Upload compliance documents |

### Payment Method
| Method | Description |
|--------|-------------|
| `GetPaymentTemplate(ctx, paymentID)` | Get payment method field template |
| `GetPaymentTemplateMeta(ctx, paymentID)` | Get detailed field metadata with validation rules |
| `AddPaymentMethod(ctx, params)` | Register a payout method (bank account) |
| `ListPaymentMethods(ctx, params)` | List registered payment methods |
| `DeletePaymentMethod(ctx, userPaymentID)` | Remove a payment method |
| `SetRefundAccount(ctx, params)` | Set refund destination account |

### Account
| Method | Description |
|--------|-------------|
| `GetAccountDetails(ctx)` | Get account info (KYC status, monthly limits) |

### Merchant
| Method | Description |
|--------|-------------|
| `ListMerchantAccounts(ctx, params)` | List merchant fiat accounts with balances |
| `ActivateMerchantAccount(ctx, params)` | Activate a fiat channel (e.g., MXN CLABE) |
| `GetPendingFundAccount(ctx, params)` | Get deposit account info for funding |
| `ClaimPendingFund(ctx, params)` | Withdraw available balance |
| `GetFundingRecords(ctx, params)` | Query funding history |
| `GetFiatFlow(ctx, params)` | Query fiat flow ledger |
| `GetCryptoBalance(ctx, merchantID)` | Query crypto asset balances |
| `GetCryptoDepositAddress(ctx, params)` | Get crypto deposit address |

### Webhook
| Method | Description |
|--------|-------------|
| `CreateWebhook(ctx, params)` | Subscribe to event notifications |
| `ListWebhooks(ctx, status)` | List webhook subscriptions |
| `UpdateWebhook(ctx, webhookID, params)` | Update webhook URL or status |

## Error Handling

```go
order, err := client.CreateFiatToCryptoOrder(ctx, params)
if err != nil {
    var apiErr *velafi.Error
    if errors.As(err, &apiErr) {
        fmt.Printf("API error: code=%d, message=%s\n", apiErr.Code, apiErr.Message)
    }

    if velafi.IsUnauthorized(err) {
        // token expired or invalid credentials
    }
}
```

## Order Status

```go
const (
    OrderStatusReviewing      = 10 // Under compliance review
    OrderStatusRFIRequested   = 11 // Additional documents required
    OrderStatusRFIUploaded    = 12 // Documents uploaded, awaiting review
    OrderStatusApproved       = 30 // Approved
    OrderStatusDocsPending    = 31 // Supporting documents needed
    OrderStatusPendingPayment = 40 // Awaiting user payment
    OrderStatusPaid           = 50 // Payment confirmed
    OrderStatusReleased       = 60 // Completed (crypto/fiat released)
    OrderStatusCanceled       = 70 // Canceled
    OrderStatusRefunded       = 72 // Refunded
    OrderStatusRefunding      = 73 // Refund in progress
)
```

---

## Use Case: User On-Ramp (MXN → USDT)

**Business scenario**: User pays MXN through local bank transfer, your system receives USDT and credits the user's internal USDT balance.

### Flow Diagram

```
User                    Your System                   VelaFi
  │                          │                           │
  │ ─── Request deposit ───> │                           │
  │                          │ ── GetBuyPayments ──────> │
  │                          │ <── paymentId list ────── │
  │                          │                           │
  │                          │ ── GetCryptoQuote ──────> │
  │                          │ <── price + quoteId ───── │
  │                          │                           │
  │ <── Show amount & ────── │                           │
  │     payment info         │                           │
  │                          │                           │
  │ ─── Confirm deposit ──>  │                           │
  │                          │ ── CreateFiatToCrypto ──> │
  │                          │ <── orderId ───────────── │
  │                          │                           │
  │                          │ ── GetOrderDetail ──────> │
  │                          │ <── status + paymentInfo ─│
  │                          │                           │
  │ <── Transfer MXN to ──── │                           │
  │     CLABE account        │                           │
  │                          │                           │
  │                          │  [Webhook: status=60]     │
  │                          │ <── ORDER_WEBHOOK ─────── │
  │                          │                           │
  │                          │  Credit user USDT balance │
  │ <── Deposit confirmed ── │                           │
```

### Implementation

```go
package main

import (
    "context"
    "fmt"
    "log"

    velafi "github.com/xfinancial/go-velafi"
)

// Step 1: Get available payment methods for MXN → USDT
func getPaymentOptions(ctx context.Context, client *velafi.Client) {
    payments, err := client.GetBuyPayments(ctx, "Mexico", "MXN", "USDT")
    if err != nil {
        log.Fatal(err)
    }
    for _, p := range payments.PaymentList {
        fmt.Printf("paymentId=%d, fee=%.2f, type=%d, trench=%s\n",
            p.PaymentID, p.FiatFee, p.PaymentType, p.Trench)
    }
    // paymentType: 0=manual transfer, 1=automatic debit
    // Select a paymentId to use for order creation
}

// Step 2: Show user the price and collect deposit amount
func getQuote(ctx context.Context, client *velafi.Client) *velafi.Quote {
    quote, err := client.GetCryptoQuote(ctx, &velafi.CryptoQuoteParams{
        Country:       "Mexico",
        From:          "MXN",
        To:            "USDT",
        CreateQuoteID: true, // Lock price for 10 seconds
    })
    if err != nil {
        log.Fatal(err)
    }
    // quote.Price = "17.84" means 1 USDT = 17.84 MXN
    // quote.QuoteID is valid for 10 seconds
    return quote
}

// Step 3: Create on-ramp order (call immediately after getting quote)
func createOnRampOrder(ctx context.Context, client *velafi.Client, merchantID int, quoteID string, mxnAmount string) string {
    order, err := client.CreateFiatToCryptoOrder(ctx, &velafi.CreateFiatToCryptoOrderParams{
        Country:    "Mexico",
        Crypto:     "USDT",
        Fiat:       "MXN",
        FiatAmount: mxnAmount,     // e.g. "1000" = 1000 MXN
        PaymentID:  14,            // From Step 1
        MerchantID: merchantID,
        QuoteID:    quoteID,       // From Step 2 (optional, locks the rate)
        ClientID:   "user-123-deposit-001", // Your internal reference
    })
    if err != nil {
        log.Fatal(err)
    }
    return order.OrderID.String()
}

// Step 4: Get order details (contains payment instructions for user)
func getPaymentInfo(ctx context.Context, client *velafi.Client, orderID string) {
    detail, err := client.GetOrderDetail(ctx, &velafi.GetOrderDetailParams{
        OrderID:   orderID,
        OrderType: "fiat_to_crypto",
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Status: %s\n", velafi.OrderStatus(detail.OrderStatus))
    fmt.Printf("Amount: %s MXN\n", detail.FiatAmount)
    fmt.Printf("Fee: %s MXN\n", detail.FiatFee)

    // detail.PaymentInfo contains bank account for user to transfer to
    // e.g. {"Cuenta CLABE": "706180304649761358"}
    // Show this to the user so they can make the MXN transfer
    fmt.Printf("Payment Info: %s\n", string(detail.PaymentInfo))
}

// Step 5: Handle order completion (via polling or webhook)
func handleOrderCompletion(ctx context.Context, client *velafi.Client, orderID string) {
    detail, err := client.GetOrderDetail(ctx, &velafi.GetOrderDetailParams{
        OrderID:   orderID,
        OrderType: "fiat_to_crypto",
    })
    if err != nil {
        log.Fatal(err)
    }

    switch velafi.OrderStatus(detail.OrderStatus) {
    case velafi.OrderStatusReleased: // 60
        // VelaFi has received MXN and released USDT
        usdtAmount := detail.CryptoAmount
        fmt.Printf("Order complete! Credit user: %s USDT\n", usdtAmount)
        // >>> YOUR BUSINESS LOGIC: credit user's internal USDT balance <<<
        // db.CreditUserBalance(userID, "USDT", usdtAmount)

    case velafi.OrderStatusRFIRequested: // 11
        // Compliance requires documents — upload invoice
        client.UploadOrderInvoice(ctx, orderID, "fiat_to_crypto", []string{"/path/to/invoice.pdf"})

    case velafi.OrderStatusCanceled, velafi.OrderStatusRefunded: // 70, 72
        fmt.Printf("Order failed: %s\n", detail.FailReason)
        // >>> Notify user of failure <<<

    default:
        fmt.Printf("Order in progress: status=%d (%s)\n",
            detail.OrderStatus, velafi.OrderStatus(detail.OrderStatus))
    }
}
```

### Production Webhook Integration

In production, use webhooks instead of polling:

```go
// Register webhook (one-time setup)
func setupWebhook(ctx context.Context, client *velafi.Client) {
    wh, err := client.CreateWebhook(ctx, &velafi.CreateWebhookParams{
        EventType: "ORDER_WEBHOOK",
        URL:       "https://your-domain.com/webhooks/velafi",
    })
    if err != nil {
        log.Fatal(err)
    }
    // Store wh.PublicKey for signature verification
    fmt.Printf("Webhook created: id=%s, publicKey=%s\n", wh.WebhookID, wh.PublicKey)
}
```

Webhook payload for order completion:
```json
{
    "webhookId": "xxx",
    "eventType": "ORDER_WEBHOOK",
    "orderId": "652836514845351936",
    "clientId": "user-123-deposit-001",
    "merchantId": 15127640,
    "orderType": "BUY",
    "fiat": "MXN",
    "fiatAmount": "1000",
    "fiatFee": "2.5",
    "crypto": "USDT",
    "cryptoAmount": "55.86",
    "orderPrice": "17.84",
    "orderStatus": 60,
    "createTime": "1715000000000",
    "updateTime": "1715000060000"
}
```

### Complete On-Ramp Flow Summary

```
1. GetBuyPayments("Mexico", "MXN", "USDT")  → choose paymentId
2. GetCryptoQuote(MXN→USDT, createQuoteId)  → get price & quoteId
3. CreateFiatToCryptoOrder(amount, paymentId, quoteId)  → get orderId
4. GetOrderDetail(orderId) → show paymentInfo (CLABE) to user
5. User transfers MXN to the CLABE account
6. Poll GetOrderDetail or receive ORDER_WEBHOOK
7. When status=60: credit user's USDT balance with cryptoAmount
```

---

## Use Case: System Off-Ramp (USDT → MXN)

**Business scenario**: System holds USDT, needs to pay out MXN to user's bank account (withdrawal).

### Flow Diagram

```
User                    Your System                   VelaFi
  │                          │                           │
  │ ─── Request withdrawal ─>│                           │
  │     (amount, bank info)  │                           │
  │                          │                           │
  │                          │ ── GetSellPayments ─────> │
  │                          │ ── GetPaymentTemplate ──> │
  │                          │ ── AddPaymentMethod ────> │
  │                          │ <── userPaymentId ─────── │
  │                          │                           │
  │                          │ ── GetCryptoQuote ──────> │
  │                          │ <── price + quoteId ───── │
  │                          │                           │
  │                          │  Debit user USDT balance  │
  │                          │                           │
  │                          │ ── CreateCryptoToFiat ──> │
  │                          │ <── orderId ───────────── │
  │                          │                           │
  │                          │  [Webhook: status=60]     │
  │                          │ <── ORDER_WEBHOOK ─────── │
  │                          │                           │
  │ <── MXN sent to bank ─── │                           │
```

### Implementation

```go
// Step 1: Register user's bank account as payout method
func registerPayoutAccount(ctx context.Context, client *velafi.Client, merchantID int) int {
    // First, check what fields are needed
    meta, err := client.GetPaymentTemplateMeta(ctx, 72) // 72 = SPEI CLABE
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Template: %s (%d fields)\n", meta.PaymentName, len(meta.FieldList))
    for _, f := range meta.FieldList {
        fmt.Printf("  %s (required=%v, type=%d)\n", f.Title, !f.IsOptional, f.TextType)
    }

    // Create payout method with user's bank info
    result, err := client.AddPaymentMethod(ctx, &velafi.AddPaymentMethodParams{
        MerchantID: merchantID,
        PaymentID:  72, // SPEI CLABE
        Country:    "Mexico",
        Fiat:       "MXN",
        RealName:   "Juan Pérez",
        FieldJSON: map[string]any{
            "Account Type":        "clabe",
            "Full Name":           "Juan Pérez",
            "Bank Account Number": "012345678901234567", // 18-digit CLABE
            "Bank Code":           "",
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    // result.ID is the userPaymentId for order creation
    // result.Status: 1=active, 2=verifying, 3=failed
    return result.ID
}

// Step 2: Check crypto balance before creating off-ramp order
func checkBalance(ctx context.Context, client *velafi.Client, merchantID int) {
    balance, err := client.GetCryptoBalance(ctx, merchantID)
    if err != nil {
        log.Fatal(err)
    }
    for _, b := range balance.Balances {
        if b.Asset == "USDT" {
            fmt.Printf("USDT: free=%s, locked=%s\n", b.Free, b.Locked)
        }
    }
}

// Step 3: Create off-ramp order
func createOffRampOrder(ctx context.Context, client *velafi.Client, merchantID, userPaymentID int, usdtAmount string) string {
    // Get quote first
    quote, err := client.GetCryptoQuote(ctx, &velafi.CryptoQuoteParams{
        Country:       "Mexico",
        From:          "USDT",
        To:            "MXN",
        CreateQuoteID: true,
    })
    if err != nil {
        log.Fatal(err)
    }
    // quote.Price = "16.89" means 1 USDT = 16.89 MXN

    // Create order immediately (quoteId valid 10 seconds)
    order, err := client.CreateCryptoToFiatOrder(ctx, &velafi.CreateCryptoToFiatOrderParams{
        Country:       "Mexico",
        Crypto:        "USDT",
        Fiat:          "MXN",
        CryptoAmount:  usdtAmount,   // e.g. "100" = 100 USDT
        UserPaymentID: userPaymentID, // From Step 1
        MerchantID:    merchantID,
        QuoteID:       quote.QuoteID,
        ClientID:      "user-123-withdraw-001",
    })
    if err != nil {
        log.Fatal(err)
    }
    return order.OrderID.String()
}

// Step 4: Confirm off-ramp order (can accelerate settlement from status 50→60)
func confirmOffRamp(ctx context.Context, client *velafi.Client, orderID int64) {
    err := client.ConfirmOrder(ctx, &velafi.ConfirmOrderParams{
        OrderID:   orderID,
        OrderType: "crypto_to_fiat",
    })
    if err != nil {
        log.Printf("Confirm failed (may not be eligible yet): %v", err)
    }
}
```

### Complete Off-Ramp Flow Summary

```
1. GetSellPayments("Mexico", "MXN", "USDT")  → get paymentId options
2. GetPaymentTemplateMeta(paymentId)  → get required fields
3. AddPaymentMethod(CLABE details)  → get userPaymentId
4. GetCryptoBalance(merchantId)  → verify sufficient USDT
5. GetCryptoQuote(USDT→MXN, createQuoteId)  → get price & quoteId
6. CreateCryptoToFiatOrder(amount, userPaymentId, quoteId)  → get orderId
7. Poll GetOrderDetail or receive ORDER_WEBHOOK
8. When status=50: optionally call ConfirmOrder to accelerate
9. When status=60: MXN has been sent to user's CLABE account
```

---

## Environment Variables

For the example program:

```bash
export VELAFI_API_KEY="your-api-key"
export VELAFI_API_SECRET="your-api-secret"
export VELAFI_MERCHANT_ID="15127640"

# Optional overrides
export VELAFI_FIAT="MXN"           # default: MXN
export VELAFI_CRYPTO="USDT"        # default: USDT
export VELAFI_FIAT_AMOUNT="500"    # default: 500
export VELAFI_CRYPTO_AMOUNT="50"   # default: 50
```

## Running the Example

```bash
go run ./example/
```

## Sandbox Notes

- Base URL: `https://api-test.velafi.com`
- No real money transfers occur
- Webhooks are **disabled** in sandbox
- KYC verification is bypassed
- Orders may require manual approval from VelaFi team to progress past compliance review (status 12)
- Off-ramp orders require pre-funded crypto balance

## License

MIT
