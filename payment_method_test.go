package velafi

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func TestGetPaymentTemplate(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payments/templates" {
			t.Errorf("path = %q, want /v2/payments/templates", r.URL.Path)
		}
		if r.URL.Query().Get("paymentId") != "1" {
			t.Errorf("paymentId = %q, want 1", r.URL.Query().Get("paymentId"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]string{
				"accountNumber": "",
				"bankName":      "",
				"routingNumber": "",
			},
		})
	})

	tpl, err := c.GetPaymentTemplate(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if _, ok := tpl["accountNumber"]; !ok {
		t.Error("template should contain accountNumber key")
	}
	if _, ok := tpl["bankName"]; !ok {
		t.Error("template should contain bankName key")
	}
}

func TestGetPaymentTemplateMeta(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payments/templates/metamessage" {
			t.Errorf("path = %q, want /v2/payments/templates/metamessage", r.URL.Path)
		}
		if r.URL.Query().Get("paymentId") != "1" {
			t.Errorf("paymentId = %q, want 1", r.URL.Query().Get("paymentId"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"id":          1,
				"paymentName": "Bank Transfer",
				"paymentType": 1,
				"trench":      "local",
				"fieldList": []map[string]any{
					{
						"index":      1,
						"indexCode":  "accountNumber",
						"textType":   1,
						"title":      "Account Number",
						"promptText": "Enter your account number",
						"minLimit":   5,
						"maxLimit":   20,
						"isOptional": false,
						"isAccount":  true,
					},
				},
			},
		})
	})

	meta, err := c.GetPaymentTemplateMeta(context.Background(), 1)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if meta.ID != 1 {
		t.Errorf("ID = %d, want 1", meta.ID)
	}
	if meta.PaymentName != "Bank Transfer" {
		t.Errorf("PaymentName = %q, want %q", meta.PaymentName, "Bank Transfer")
	}
	if len(meta.FieldList) != 1 {
		t.Fatalf("len(FieldList) = %d, want 1", len(meta.FieldList))
	}
	if meta.FieldList[0].IndexCode != "accountNumber" {
		t.Errorf("IndexCode = %q, want %q", meta.FieldList[0].IndexCode, "accountNumber")
	}
}

func TestAddPaymentMethod(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payments" {
			t.Errorf("path = %q, want /v2/payments", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %q, want POST", r.Method)
		}
		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["country"] != "US" {
			t.Errorf("country = %v, want US", body["country"])
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"id":     101,
				"status": 1,
			},
		})
	})

	result, err := c.AddPaymentMethod(context.Background(), &AddPaymentMethodParams{
		MerchantID: 1,
		PaymentID:  1,
		Country:    "US",
		Fiat:       "USD",
		RealName:   "John Doe",
		FieldJSON:  map[string]any{"accountNumber": "123456"},
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if result.ID != 101 {
		t.Errorf("ID = %d, want 101", result.ID)
	}
	if result.Status != 1 {
		t.Errorf("Status = %d, want 1", result.Status)
	}
}

func TestListPaymentMethods(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payments" {
			t.Errorf("path = %q, want /v2/payments", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %q, want GET", r.Method)
		}
		q := r.URL.Query()
		if q.Get("country") != "US" {
			t.Errorf("country = %q, want US", q.Get("country"))
		}
		if q.Get("fiat") != "USD" {
			t.Errorf("fiat = %q, want USD", q.Get("fiat"))
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": map[string]any{
				"currentPage": 1,
				"size":        20,
				"total":       1,
				"record": []map[string]any{
					{
						"id":                101,
						"merchantId":        1,
						"country":           "US",
						"fiat":              "USD",
						"paymentId":         1,
						"paymentMethodName": "Bank Transfer",
						"realName":          "John Doe",
						"status":            1,
						"hasRefundAccount":  0,
						"fieldList":         map[string]string{"accountNumber": "123456"},
						"createTime":        1714500000000,
					},
				},
			},
		})
	})

	list, err := c.ListPaymentMethods(context.Background(), &ListPaymentMethodsParams{
		Country: "US",
		Fiat:    "USD",
	})
	if err != nil {
		t.Fatalf("error = %v", err)
	}
	if list.Total != 1 {
		t.Errorf("Total = %d, want 1", list.Total)
	}
	if len(list.Record) != 1 {
		t.Fatalf("len(Record) = %d, want 1", len(list.Record))
	}
	if list.Record[0].ID != 101 {
		t.Errorf("ID = %d, want 101", list.Record[0].ID)
	}
	if list.Record[0].PaymentMethodName != "Bank Transfer" {
		t.Errorf("PaymentMethodName = %q, want %q", list.Record[0].PaymentMethodName, "Bank Transfer")
	}
}

func TestDeletePaymentMethod(t *testing.T) {
	c := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v2/payments/101" {
			t.Errorf("path = %q, want /v2/payments/101", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("method = %q, want DELETE", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"code": 200,
			"msg":  "SUCCESS",
			"data": nil,
		})
	})

	err := c.DeletePaymentMethod(context.Background(), 101)
	if err != nil {
		t.Fatalf("error = %v", err)
	}
}
