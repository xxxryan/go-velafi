package velafi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateFiatToCryptoOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/order/fiat_to_crypto" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200, "msg": "SUCCESS",
			"data": map[string]any{"orderId": 1001},
		})
	})

	result, err := c.CreateFiatToCryptoOrder(context.Background(), &CreateFiatToCryptoOrderParams{
		Country: "Mexico", Crypto: "USDT", Fiat: "MXN", FiatAmount: "500", PaymentID: 14,
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if result.OrderID != "1001" {
		t.Errorf("OrderID = %q", result.OrderID)
	}
}

func TestCreateCryptoToFiatOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/order/crypto_to_fiat" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200, "msg": "SUCCESS",
			"data": map[string]any{"orderId": "2002"},
		})
	})

	result, err := c.CreateCryptoToFiatOrder(context.Background(), &CreateCryptoToFiatOrderParams{
		Country: "Mexico", Crypto: "USDT", Fiat: "MXN", CryptoAmount: "50", UserPaymentID: 91,
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if result.OrderID != "2002" {
		t.Errorf("OrderID = %q", result.OrderID)
	}
}

func TestGetOrderDetail(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/order/detail" {
			t.Errorf("path = %q", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("orderId") != "1001" || q.Get("orderType") != "fiat_to_crypto" {
			t.Errorf("query = %v", q)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200, "msg": "SUCCESS",
			"data": map[string]any{
				"orderId": "1001", "orderType": "fiat_to_crypto", "orderStatus": 40,
				"country": "Mexico", "fiat": "MXN", "crypto": "USDT", "fiatAmount": "500",
			},
		})
	})

	order, err := c.GetOrderDetail(context.Background(), &GetOrderDetailParams{
		OrderID: "1001", OrderType: "fiat_to_crypto",
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if order.OrderID != "1001" {
		t.Errorf("OrderID = %q", order.OrderID)
	}
	if order.OrderStatus != 40 {
		t.Errorf("OrderStatus = %d", order.OrderStatus)
	}
}

func TestConfirmOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/order/confirm" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200, "msg": "SUCCESS", "data": true,
		})
	})

	err := c.ConfirmOrder(context.Background(), &ConfirmOrderParams{
		OrderID: 1001, OrderType: "fiat_to_crypto",
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}
