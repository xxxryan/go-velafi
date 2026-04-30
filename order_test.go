package velafi

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestCreateFiatCryptoOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orders/fiat-crypto" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["merchantId"] != "m-1" {
			t.Errorf("merchantId = %v", body["merchantId"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code":    0,
			"message": "success",
			"data": map[string]any{
				"orderId":      "ord-123",
				"status":       "created",
				"fromCurrency": "USD",
				"toCurrency":   "BTC",
				"fromAmount":   "5000",
				"toAmount":     "0.1",
			},
		})
	})

	order, err := c.CreateFiatCryptoOrder(context.Background(), &CreateFiatCryptoOrderParams{
		MerchantID:      "m-1",
		FromCurrency:    "USD",
		ToCurrency:      "BTC",
		FromAmount:      "5000",
		PaymentMethodID: "pm-1",
		ToAddress:       "bc1q...",
		Network:         "Bitcoin",
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if order.OrderID != "ord-123" {
		t.Errorf("OrderID = %q", order.OrderID)
	}
}

func TestCreateCryptoFiatOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orders/crypto-fiat" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{"orderId": "ord-456", "status": "created", "fromCurrency": "BTC", "toCurrency": "USD", "fromAmount": "0.1", "toAmount": "5000"},
		})
	})

	order, err := c.CreateCryptoFiatOrder(context.Background(), &CreateCryptoFiatOrderParams{
		MerchantID: "m-1", FromCurrency: "BTC", ToCurrency: "USD", FromAmount: "0.1", Network: "Bitcoin", PaymentMethodID: "pm-1",
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if order.OrderID != "ord-456" {
		t.Errorf("OrderID = %q", order.OrderID)
	}
}

func TestCreateFiatFiatOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orders/fiat-fiat" {
			t.Errorf("path = %q", r.URL.Path)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{"orderId": "ord-789", "status": "created", "fromCurrency": "EUR", "toCurrency": "USD", "fromAmount": "1000", "toAmount": "1100"},
		})
	})

	order, err := c.CreateFiatFiatOrder(context.Background(), &CreateFiatFiatOrderParams{
		MerchantID: "m-1", FromCurrency: "EUR", ToCurrency: "USD", FromAmount: "1000", FromPaymentMethodID: "pm-from", ToPaymentMethodID: "pm-to",
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if order.OrderID != "ord-789" {
		t.Errorf("OrderID = %q", order.OrderID)
	}
}

func TestConfirmOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orders/ord-123/confirm" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{"orderId": "ord-123", "status": "confirmed", "confirmedAt": "2026-04-30T10:00:00Z"},
		})
	})

	result, err := c.ConfirmOrder(context.Background(), "ord-123")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if result.Status != "confirmed" {
		t.Errorf("Status = %q", result.Status)
	}
}

func TestGetOrder(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orders/ord-123" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{"orderId": "ord-123", "status": "completed", "fromCurrency": "USD", "toCurrency": "BTC", "fromAmount": "5000", "toAmount": "0.1"},
		})
	})

	order, err := c.GetOrder(context.Background(), "ord-123")
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if order.OrderID != "ord-123" {
		t.Errorf("OrderID = %q", order.OrderID)
	}
}

func TestListOrders(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orders" {
			t.Errorf("path = %q", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("merchantId") != "m-1" {
			t.Errorf("merchantId = %q", q.Get("merchantId"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{
				"total": 1, "page": 1, "limit": 20,
				"items": []map[string]any{{"orderId": "ord-123", "status": "completed", "fromCurrency": "USD", "toCurrency": "BTC", "fromAmount": "5000", "toAmount": "0.1"}},
			},
		})
	})

	list, err := c.ListOrders(context.Background(), &ListOrdersParams{MerchantID: "m-1"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if list.Total != 1 {
		t.Errorf("Total = %d", list.Total)
	}
	if len(list.Items) != 1 {
		t.Fatalf("len(Items) = %d", len(list.Items))
	}
}

func TestUploadInvoiceDocuments(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/orders/ord-123/invoice-documents" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q", r.Method)
		}
		body, _ := io.ReadAll(r.Body)
		if !strings.Contains(string(body), "file-a") {
			t.Errorf("body missing fileIds: %s", body)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 0, "message": "success",
			"data": map[string]any{
				"orderId": "ord-123",
				"documents": []map[string]any{
					{"fileId": "file-a", "filename": "invoice.pdf", "uploadedAt": "2026-04-30T10:00:00Z"},
				},
			},
		})
	})

	result, err := c.UploadInvoiceDocuments(context.Background(), "ord-123", []string{"file-a"})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if result.OrderID != "ord-123" {
		t.Errorf("OrderID = %q", result.OrderID)
	}
	if len(result.Documents) != 1 {
		t.Errorf("len(Documents) = %d", len(result.Documents))
	}
}
