package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestAPIError_Error は Error() の文字列書式を固定する。
// 既存の利用者がメッセージに依存している可能性があるため、書式は変更しない。
func TestAPIError_Error(t *testing.T) {
	tests := []struct {
		name string
		err  *APIError
		want string
	}{
		{
			name: "rate limit",
			err:  &APIError{StatusCode: 429, Body: `{"message":"Too Many Requests"}`},
			want: `API error: status=429, body={"message":"Too Many Requests"}`,
		},
		{
			name: "forbidden",
			err:  &APIError{StatusCode: 403, Body: `{"message":"Forbidden"}`},
			want: `API error: status=403, body={"message":"Forbidden"}`,
		},
		{
			name: "empty body",
			err:  &APIError{StatusCode: 500, Body: ""},
			want: "API error: status=500, body=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestAPIError_As(t *testing.T) {
	base := &APIError{StatusCode: 429, Body: `{"message":"Too Many Requests"}`}

	tests := []struct {
		name string
		err  error
	}{
		{name: "direct", err: base},
		{name: "wrapped once", err: fmt.Errorf("failed to get daily quotes: %w", base)},
		{name: "wrapped twice", err: fmt.Errorf("outer: %w", fmt.Errorf("failed to get daily quotes: %w", base))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var apiErr *APIError
			if !errors.As(tt.err, &apiErr) {
				t.Fatalf("errors.As failed for %v", tt.err)
			}
			if apiErr.StatusCode != 429 {
				t.Errorf("StatusCode = %d, want 429", apiErr.StatusCode)
			}
			if apiErr.Body != base.Body {
				t.Errorf("Body = %q, want %q", apiErr.Body, base.Body)
			}
			if !IsRateLimitExceeded(tt.err) {
				t.Error("IsRateLimitExceeded = false, want true")
			}
		})
	}
}

// TestStatusCode_MisleadingMessage はボディに 401 や 429 という数字列が含まれていても、
// 実際のステータスコードで判定されることを確認する。
func TestStatusCode_MisleadingMessage(t *testing.T) {
	// 銘柄コードやメッセージ本文に、他のステータスコードと一致する数字列が含まれるケース
	body := `{"message":"invalid code 4290 (401k related), see status 429 doc"}`
	err := fmt.Errorf("failed to get daily quotes: %w", &APIError{StatusCode: http.StatusBadRequest, Body: body})

	// 文字列マッチであれば誤判定するボディであることを確認する
	if !strings.Contains(err.Error(), "429") || !strings.Contains(err.Error(), "401") {
		t.Fatal("test fixture must contain misleading 401/429 digits")
	}

	code, ok := StatusCode(err)
	if !ok || code != http.StatusBadRequest {
		t.Errorf("StatusCode() = (%d, %v), want (400, true)", code, ok)
	}
	if IsRateLimitExceeded(err) {
		t.Error("IsRateLimitExceeded = true for a 400 response containing '429'")
	}
	if IsAuthError(err) {
		t.Error("IsAuthError = true for a 400 response containing '401'")
	}
	if IsServerError(err) {
		t.Error("IsServerError = true for a 400 response")
	}
}

func TestErrorPredicates(t *testing.T) {
	tests := []struct {
		name            string
		err             error
		wantCode        int
		wantOK          bool
		wantRateLimited bool
		wantAuth        bool
		wantServer      bool
	}{
		{name: "400", err: &APIError{StatusCode: 400}, wantCode: 400, wantOK: true},
		{name: "401", err: &APIError{StatusCode: 401}, wantCode: 401, wantOK: true, wantAuth: true},
		{name: "403", err: &APIError{StatusCode: 403}, wantCode: 403, wantOK: true, wantAuth: true},
		{name: "429", err: &APIError{StatusCode: 429}, wantCode: 429, wantOK: true, wantRateLimited: true},
		{name: "500", err: &APIError{StatusCode: 500}, wantCode: 500, wantOK: true, wantServer: true},
		{name: "503", err: &APIError{StatusCode: 503}, wantCode: 503, wantOK: true, wantServer: true},
		// 210 (No Content) はデータ無しを示すステータスで、5xx ではない
		{name: "210", err: &APIError{StatusCode: 210}, wantCode: 210, wantOK: true},
		{name: "nil", err: nil},
		{name: "context canceled", err: context.Canceled},
		{name: "non-API error", err: errors.New("API error: status=429, body=oops")},
		{name: "wrapped non-API error", err: fmt.Errorf("failed to send request: %w", context.DeadlineExceeded)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, ok := StatusCode(tt.err)
			if code != tt.wantCode || ok != tt.wantOK {
				t.Errorf("StatusCode() = (%d, %v), want (%d, %v)", code, ok, tt.wantCode, tt.wantOK)
			}
			if got := IsRateLimitExceeded(tt.err); got != tt.wantRateLimited {
				t.Errorf("IsRateLimitExceeded() = %v, want %v", got, tt.wantRateLimited)
			}
			if got := IsAuthError(tt.err); got != tt.wantAuth {
				t.Errorf("IsAuthError() = %v, want %v", got, tt.wantAuth)
			}
			if got := IsServerError(tt.err); got != tt.wantServer {
				t.Errorf("IsServerError() = %v, want %v", got, tt.wantServer)
			}
		})
	}
}

func TestClient_DoRequest_APIError(t *testing.T) {
	const body = `{"message":"Too Many Requests"}`
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTooManyRequests)
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	c := NewClient("test-api-key")
	c.baseURL = server.URL

	var resp struct{ Message string }
	err := c.DoRequest(context.Background(), http.MethodGet, "/test", nil, &resp)
	if err == nil {
		t.Fatal("expected an error")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("errors.As failed for %v", err)
	}
	if apiErr.StatusCode != http.StatusTooManyRequests {
		t.Errorf("StatusCode = %d, want 429", apiErr.StatusCode)
	}
	if apiErr.Body != body {
		t.Errorf("Body = %q, want %q", apiErr.Body, body)
	}
	if !IsRateLimitExceeded(err) {
		t.Error("IsRateLimitExceeded = false, want true")
	}
	// 既存の利用者のためにメッセージ書式を維持する
	want := "API error: status=429, body=" + body
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}

func TestClient_DoRequest_APIErrorWithCache(t *testing.T) {
	var callCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(`{"message":"Forbidden"}`))
	}))
	defer server.Close()

	c := NewClient("test-api-key", WithCache())
	c.baseURL = server.URL

	// キャッシュ有効時（singleflight経由）でも同じ型が返る
	var resp struct{ Message string }
	err := c.DoRequest(context.Background(), http.MethodGet, "/test", nil, &resp)
	if !IsAuthError(err) {
		t.Fatalf("IsAuthError = false for %v", err)
	}

	// エラーレスポンスはキャッシュされない
	if c.CacheSize() != 0 {
		t.Errorf("expected error response not to be cached, got %d entries", c.CacheSize())
	}
	if err := c.DoRequest(context.Background(), http.MethodGet, "/test", nil, &resp); !IsAuthError(err) {
		t.Errorf("IsAuthError = false for %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}

	// DoRequestNoCache も同様
	if err := c.DoRequestNoCache(context.Background(), http.MethodGet, "/test", nil, &resp); !IsAuthError(err) {
		t.Errorf("IsAuthError = false for %v", err)
	}
}

func TestClient_DoRequest_NonAPIError(t *testing.T) {
	// ステータスは200だがボディがJSONではない場合、APIErrorにはならない
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("not json"))
	}))
	defer server.Close()

	c := NewClient("test-api-key")
	c.baseURL = server.URL

	var resp struct{ Message string }
	err := c.DoRequest(context.Background(), http.MethodGet, "/test", nil, &resp)
	if err == nil {
		t.Fatal("expected an error")
	}
	if code, ok := StatusCode(err); ok {
		t.Errorf("StatusCode() = (%d, true), want ok=false for a decode error", code)
	}
}
