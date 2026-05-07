package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	velafi "github.com/xxxryan/go-velafi"
)

func main() {
	apiKey := os.Getenv("VELAFI_API_KEY")
	apiSecret := os.Getenv("VELAFI_API_SECRET")
	merchantIDStr := os.Getenv("VELAFI_MERCHANT_ID")
	if apiKey == "" || apiSecret == "" || merchantIDStr == "" {
		log.Fatal("set VELAFI_API_KEY, VELAFI_API_SECRET, VELAFI_MERCHANT_ID")
	}
	merchantID, err := strconv.Atoi(merchantIDStr)
	if err != nil {
		log.Fatalf("invalid VELAFI_MERCHANT_ID: %v", err)
	}

	fiat := envOrDefault("VELAFI_FIAT", "MXN")
	crypto := envOrDefault("VELAFI_CRYPTO", "USDT")
	fiatAmount := envOrDefault("VELAFI_FIAT_AMOUNT", "500")
	cryptoAmount := envOrDefault("VELAFI_CRYPTO_AMOUNT", "50")

	client := velafi.NewClient(apiKey, apiSecret, velafi.WithSandbox())
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	fmt.Println("========== On-Ramp: MXN → USDT ==========")
	onRamp(ctx, client, merchantID, fiat, crypto, fiatAmount)

	fmt.Println("\n========== Off-Ramp: USDT → MXN ==========")
	offRamp(ctx, client, merchantID, fiat, crypto, cryptoAmount)
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func findSymbol(symbols []velafi.BuySellSymbol, fiat, crypto string) *velafi.BuySellSymbol {
	for _, s := range symbols {
		if s.Fiat == fiat && s.Crypto == crypto {
			return &s
		}
	}
	return nil
}

func onRamp(ctx context.Context, client *velafi.Client, merchantID int, fiat, crypto, fiatAmount string) {
	// Step 1: 确认交易对
	fmt.Print("[1] GetBuySymbols... ")
	symbols, err := client.GetBuySymbols(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	sym := findSymbol(symbols, fiat, crypto)
	if sym == nil {
		log.Fatalf("交易对 %s/%s 不存在", fiat, crypto)
	}
	fmt.Printf("OK (%s %s/%s)\n", sym.Country, sym.Fiat, sym.Crypto)

	// Step 2: 查询可用付款方式
	fmt.Print("[2] GetBuyPayments... ")
	payments, err := client.GetBuyPayments(ctx, sym.Country, sym.Fiat, sym.Crypto)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if len(payments.PaymentList) == 0 {
		log.Fatal("无可用付款方式")
	}
	for _, p := range payments.PaymentList {
		fmt.Printf("\n    paymentId=%d fee=%.2f type=%d trench=%s", p.PaymentID, p.FiatFee, p.PaymentType, p.Trench)
	}
	payment := payments.PaymentList[0]
	// 沙盒中自动模式可能未激活，优先选手动 (paymentType=0)
	for _, p := range payments.PaymentList {
		if p.PaymentType == 0 {
			payment = p
			break
		}
	}
	fmt.Printf("\n    → 选择 paymentId=%d\n", payment.PaymentID)

	// Step 3: 获取报价 (quoteId 10秒有效，立即下单)
	fmt.Print("[3] GetCryptoQuote... ")
	quote, err := client.GetCryptoQuote(ctx, &velafi.CryptoQuoteParams{
		Country: sym.Country, From: sym.Fiat, To: sym.Crypto, CreateQuoteID: true,
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("OK (price=%s, quoteId=%s)\n", quote.Price, quote.QuoteID)

	// Step 4: 创建 on-ramp 订单
	fmt.Printf("[4] CreateFiatToCryptoOrder (amount=%s %s)... ", fiatAmount, fiat)
	order, err := client.CreateFiatToCryptoOrder(ctx, &velafi.CreateFiatToCryptoOrderParams{
		Country:    sym.Country,
		Crypto:     sym.Crypto,
		Fiat:       sym.Fiat,
		FiatAmount: fiatAmount,
		PaymentID:  payment.PaymentID,
		MerchantID: merchantID,
		QuoteID:    quote.QuoteID,
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("OK (orderId=%s)\n", order.OrderID)

	// Step 5: 查询订单状态
	fmt.Print("[5] GetOrderDetail... ")
	detail, err := client.GetOrderDetail(ctx, &velafi.GetOrderDetailParams{
		OrderID: order.OrderID.String(), OrderType: "fiat_to_crypto",
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("OK\n")
	fmt.Printf("    orderId:     %s\n", detail.OrderID)
	fmt.Printf("    orderStatus: %d (%s)\n", detail.OrderStatus, velafi.OrderStatus(detail.OrderStatus))
	fmt.Printf("    fiatAmount:  %s %s\n", detail.FiatAmount, fiat)
	fmt.Printf("    crypto:      %s\n", detail.Crypto)
	if len(detail.PaymentInfo) > 0 {
		fmt.Printf("    paymentInfo: %s\n", string(detail.PaymentInfo))
	}
}

func offRamp(ctx context.Context, client *velafi.Client, merchantID int, fiat, crypto, cryptoAmount string) {
	// Step 1: 确认交易对
	fmt.Print("[1] GetSellSymbols... ")
	symbols, err := client.GetSellSymbols(ctx)
	if err != nil {
		log.Fatalf("%v", err)
	}
	sym := findSymbol(symbols, fiat, crypto)
	if sym == nil {
		log.Fatalf("交易对 %s/%s 不存在", crypto, fiat)
	}
	fmt.Printf("OK (%s %s/%s)\n", sym.Country, sym.Crypto, sym.Fiat)

	// Step 2: 查询 payment template
	fmt.Print("[2] GetSellPayments... ")
	payments, err := client.GetSellPayments(ctx, sym.Country, sym.Fiat, sym.Crypto)
	if err != nil {
		log.Fatalf("%v", err)
	}
	if len(payments.PaymentList) == 0 {
		log.Fatal("无可用收款方式")
	}
	payment := payments.PaymentList[0]
	fmt.Printf("OK (paymentId=%d)\n", payment.PaymentID)

	fmt.Printf("    GetPaymentTemplate(paymentId=%d)... ", payment.PaymentID)
	tpl, err := client.GetPaymentTemplate(ctx, payment.PaymentID)
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("OK\n")
	for k, v := range tpl {
		fmt.Printf("    %s: %s\n", k, v)
	}

	// Step 3: 创建用户收款方式（根据 template 字段填充）
	fmt.Print("[3] AddPaymentMethod... ")
	fieldJSON := make(map[string]any)
	for k := range tpl {
		if strings.Contains(strings.ToLower(k), "optativo") || strings.Contains(strings.ToLower(k), "optional") {
			continue
		}
		if strings.Contains(strings.ToLower(k), "nombre") || strings.Contains(strings.ToLower(k), "name") {
			fieldJSON[k] = "John Doe"
		} else {
			fieldJSON[k] = "1234567890123456"
		}
	}
	addResult, err := client.AddPaymentMethod(ctx, &velafi.AddPaymentMethodParams{
		MerchantID: merchantID,
		PaymentID:  payment.PaymentID,
		Country:    sym.Country,
		Fiat:       sym.Fiat,
		RealName:   "John Doe",
		FieldJSON:  fieldJSON,
	})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		fmt.Println("    (收款方式添加失败，跳过后续步骤)")
		return
	}
	fmt.Printf("OK (userPaymentId=%d, status=%d)\n", addResult.ID, addResult.Status)

	// Step 4: 获取报价
	fmt.Print("[4] GetCryptoQuote (off-ramp)... ")
	quote, err := client.GetCryptoQuote(ctx, &velafi.CryptoQuoteParams{
		Country: sym.Country, From: sym.Crypto, To: sym.Fiat, CreateQuoteID: true,
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("OK (price=%s)\n", quote.Price)

	// Step 5: 创建 off-ramp 订单
	fmt.Printf("[5] CreateCryptoToFiatOrder (amount=%s %s)... ", cryptoAmount, crypto)
	order, err := client.CreateCryptoToFiatOrder(ctx, &velafi.CreateCryptoToFiatOrderParams{
		Country:       sym.Country,
		Crypto:        sym.Crypto,
		Fiat:          sym.Fiat,
		CryptoAmount:  cryptoAmount,
		UserPaymentID: addResult.ID,
		MerchantID:    merchantID,
		QuoteID:       quote.QuoteID,
	})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		fmt.Println("    (off-ramp 需要 crypto 余额，沙盒可能余额为 0)")
		return
	}
	fmt.Printf("OK (orderId=%s)\n", order.OrderID)

	// Step 6: 查询订单状态
	fmt.Print("[6] GetOrderDetail... ")
	detail, err := client.GetOrderDetail(ctx, &velafi.GetOrderDetailParams{
		OrderID: order.OrderID.String(), OrderType: "crypto_to_fiat",
	})
	if err != nil {
		log.Fatalf("%v", err)
	}
	fmt.Printf("OK (status=%d %s)\n", detail.OrderStatus, velafi.OrderStatus(detail.OrderStatus))

	// Step 7: 确认订单 (如果状态 eligible)
	if detail.OrderStatus == int(velafi.OrderStatusPaid) {
		fmt.Print("[7] ConfirmOrder... ")
		if err := client.ConfirmOrder(ctx, &velafi.ConfirmOrderParams{
			OrderID: 0, OrderType: "crypto_to_fiat",
		}); err != nil {
			fmt.Printf("失败: %v\n", err)
		} else {
			fmt.Println("OK")
		}
	}
}
