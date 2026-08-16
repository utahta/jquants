package jquants

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/utahta/jquants/client"
)

// TestServiceError_As はサービスメソッドがエラーをラップしても
// errors.As で *client.APIError を取り出せることを確認する。
func TestServiceError_As(t *testing.T) {
	mockClient := client.NewMockClient()
	mockClient.SetError("GET", "/equities/bars/daily?code=7203", &client.APIError{
		StatusCode: http.StatusTooManyRequests,
		Body:       `{"message":"Too Many Requests"}`,
	})
	service := NewQuotesService(mockClient)

	// ページネーションで複数回ラップされる経路を通す
	_, err := service.GetDailyQuotesByCode(context.Background(), "7203")
	if err == nil {
		t.Fatal("expected an error")
	}

	var apiErr *client.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("errors.As failed for %v", err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("StatusCode = %d, want 429", apiErr.StatusCode)
	}
	if !client.IsRateLimitExceeded(err) {
		t.Error("IsRateLimitExceeded = false, want true")
	}
	if client.IsAuthError(err) {
		t.Error("IsAuthError = true, want false")
	}
}

// TestServiceError_MisleadingMessage は銘柄コード等の数字列に引きずられず、
// ステータスコードで判定されることを確認する。
func TestServiceError_MisleadingMessage(t *testing.T) {
	mockClient := client.NewMockClient()
	// ボディに 401 / 429 という数字列を含むが、実際のステータスは400
	mockClient.SetError("GET", "/equities/bars/daily?code=4290", &client.APIError{
		StatusCode: http.StatusBadRequest,
		Body:       `{"message":"invalid code 4290, retry limit 429 exceeded for 401k"}`,
	})
	service := NewQuotesService(mockClient)

	_, err := service.GetDailyQuotes(context.Background(), DailyQuotesParams{Code: "4290"})
	if err == nil {
		t.Fatal("expected an error")
	}

	code, ok := client.StatusCode(err)
	if !ok || code != http.StatusBadRequest {
		t.Errorf("StatusCode() = (%d, %v), want (400, true)", code, ok)
	}
	if client.IsRateLimitExceeded(err) {
		t.Error("IsRateLimitExceeded = true for a 400 response containing '429'")
	}
	if client.IsAuthError(err) {
		t.Error("IsAuthError = true for a 400 response containing '401'")
	}
}
