package velafi

import (
	"errors"
	"testing"
)

func TestError_Error(t *testing.T) {
	e := &Error{
		HTTPStatus: 400,
		Code:       10001,
		Message:    "invalid parameter",
	}
	got := e.Error()
	want := "velafi: 10001 invalid parameter (http 400)"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestError_ErrorWithRequestID(t *testing.T) {
	e := &Error{
		HTTPStatus: 500,
		Code:       50000,
		Message:    "internal error",
		RequestID:  "req-abc-123",
	}
	got := e.Error()
	want := "velafi: 50000 internal error (http 500, request_id=req-abc-123)"
	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestErrorsAs(t *testing.T) {
	e := &Error{HTTPStatus: 404, Code: 40400, Message: "not found"}
	var wrapped error = e

	var apiErr *Error
	if !errors.As(wrapped, &apiErr) {
		t.Fatal("errors.As should match *Error")
	}
	if apiErr.HTTPStatus != 404 {
		t.Errorf("HTTPStatus = %d, want 404", apiErr.HTTPStatus)
	}
}

func TestIsNotFound(t *testing.T) {
	e := &Error{HTTPStatus: 404, Code: 40400, Message: "not found"}
	if !IsNotFound(e) {
		t.Error("IsNotFound should return true for 404")
	}
	if IsNotFound(&Error{HTTPStatus: 400}) {
		t.Error("IsNotFound should return false for 400")
	}
	if IsNotFound(errors.New("network error")) {
		t.Error("IsNotFound should return false for non-API errors")
	}
}

func TestIsUnauthorized(t *testing.T) {
	e := &Error{HTTPStatus: 401, Code: 40100, Message: "unauthorized"}
	if !IsUnauthorized(e) {
		t.Error("IsUnauthorized should return true for 401")
	}
	if IsUnauthorized(&Error{HTTPStatus: 200}) {
		t.Error("IsUnauthorized should return false for 200")
	}
}

func TestIsRateLimited(t *testing.T) {
	e := &Error{HTTPStatus: 429, Code: 42900, Message: "rate limited"}
	if !IsRateLimited(e) {
		t.Error("IsRateLimited should return true for 429")
	}
	if IsRateLimited(&Error{HTTPStatus: 200}) {
		t.Error("IsRateLimited should return false for 200")
	}
}
