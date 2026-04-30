package velafi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestCreateFiatCryptoOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/order/fiat_to_crypto" {
			t.Errorf("path = %q, want /v2/order/fiat_to_crypto", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["country"] != "US" {
			t.Errorf("country = %v, want US", body["country"])
		}
		if body["fiat"] != "USD" {
			t.Errorf("fiat = %v, want USD", body["fiat"])
		}
		if body["crypto"] != "BTC" {
			t.Errorf("crypto = %v, want BTC", body["crypto"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"orderId": 1001,
			},
		})
	})

	result, err := c.CreateFiatCryptoOrder(context.Background(), &CreateFiatCryptoOrderParams{
		Country:    "US",
		Crypto:     "BTC",
		Fiat:       "USD",
		FiatAmount: "5000",
		PaymentID:  1,
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if result.OrderID != 1001 {
		t.Errorf("OrderID = %d, want 1001", result.OrderID)
	}
}

func TestGetOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/order/detail" {
			t.Errorf("path = %q, want /v2/order/detail", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		q := r.URL.Query()
		if q.Get("orderId") != "1001" {
			t.Errorf("orderId = %q, want 1001", q.Get("orderId"))
		}
		if q.Get("orderType") != "fiat_to_crypto" {
			t.Errorf("orderType = %q, want fiat_to_crypto", q.Get("orderType"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"orderId":     1001,
				"orderType":   "fiat_to_crypto",
				"orderStatus": 1,
				"country":     "US",
				"fiat":        "USD",
				"crypto":      "BTC",
				"fiatAmount":  "5000",
				"orderPrice":  "50000.00",
			},
		})
	})

	order, err := c.GetOrder(context.Background(), &GetOrderParams{
		OrderID:   1001,
		OrderType: "fiat_to_crypto",
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if order.OrderID != 1001 {
		t.Errorf("OrderID = %d, want 1001", order.OrderID)
	}
	if order.OrderType != "fiat_to_crypto" {
		t.Errorf("OrderType = %q, want fiat_to_crypto", order.OrderType)
	}
	if order.FiatAmount != "5000" {
		t.Errorf("FiatAmount = %q, want 5000", order.FiatAmount)
	}
}

func TestListOrders(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/orders" {
			t.Errorf("path = %q, want /v2/orders", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("orderType") != "fiat_to_crypto" {
			t.Errorf("orderType = %q, want fiat_to_crypto", q.Get("orderType"))
		}
		if q.Get("currentPage") != "1" {
			t.Errorf("currentPage = %q, want 1", q.Get("currentPage"))
		}
		if q.Get("pageSize") != "20" {
			t.Errorf("pageSize = %q, want 20", q.Get("pageSize"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"total":       1,
				"currentPage": 1,
				"size":        20,
				"record":      []map[string]any{{"orderId": 1001, "orderType": "fiat_to_crypto", "orderStatus": 1}},
			},
		})
	})

	list, err := c.ListOrders(context.Background(), &ListOrdersParams{
		OrderType:   "fiat_to_crypto",
		CurrentPage: 1,
		PageSize:    20,
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if list.Total != 1 {
		t.Errorf("Total = %d, want 1", list.Total)
	}
	if list.CurrentPage != 1 {
		t.Errorf("CurrentPage = %d, want 1", list.CurrentPage)
	}
	if len(list.Record) != 1 {
		t.Fatalf("len(Record) = %d, want 1", len(list.Record))
	}
}
