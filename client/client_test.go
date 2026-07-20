package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	t.Run("without options", func(t *testing.T) {
		c := NewClient("test-api-key")
		if c.apiKey != "test-api-key" {
			t.Errorf("expected apiKey to be 'test-api-key', got '%s'", c.apiKey)
		}
		if c.cacheEnabled {
			t.Error("expected cacheEnabled to be false by default")
		}
		if c.cache != nil {
			t.Error("expected cache to be nil by default")
		}
	})

	t.Run("with cache option", func(t *testing.T) {
		c := NewClient("test-api-key", WithCache())
		if !c.cacheEnabled {
			t.Error("expected cacheEnabled to be true")
		}
		if c.cache == nil {
			t.Error("expected cache to be initialized")
		}
	})
}

func TestNewClientFromEnv(t *testing.T) {
	t.Run("with cache option", func(t *testing.T) {
		t.Setenv("JQUANTS_API_KEY", "test-api-key")
		c, err := NewClientFromEnv(WithCache())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !c.cacheEnabled {
			t.Error("expected cacheEnabled to be true")
		}
		if c.cache == nil {
			t.Error("expected cache to be initialized")
		}
	})
}

func TestClient_DoRequest_CacheHitMiss(t *testing.T) {
	var callCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "hello"})
	}))
	defer server.Close()

	c := NewClient("test-api-key", WithCache())
	c.baseURL = server.URL

	type response struct {
		Message string `json:"message"`
	}

	// First request (cache miss)
	var resp1 response
	err := c.DoRequest(context.Background(), http.MethodGet, "/test", nil, &resp1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp1.Message != "hello" {
		t.Errorf("expected message 'hello', got '%s'", resp1.Message)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	// Second request (cache hit)
	var resp2 response
	err = c.DoRequest(context.Background(), http.MethodGet, "/test", nil, &resp2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp2.Message != "hello" {
		t.Errorf("expected message 'hello', got '%s'", resp2.Message)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call (cached), got %d", callCount)
	}

	// Different path is a separate cache entry
	var resp3 response
	err = c.DoRequest(context.Background(), http.MethodGet, "/test2", nil, &resp3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls, got %d", callCount)
	}
}

func TestClient_DoRequestNoCache(t *testing.T) {
	var callCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]int{"count": callCount})
	}))
	defer server.Close()

	c := NewClient("test-api-key", WithCache())
	c.baseURL = server.URL

	type response struct {
		Count int `json:"count"`
	}

	// DoRequestNoCache always hits the server and never populates the cache
	var resp1, resp2 response
	if err := c.DoRequestNoCache(context.Background(), http.MethodGet, "/signed", nil, &resp1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := c.DoRequestNoCache(context.Background(), http.MethodGet, "/signed", nil, &resp2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (cache bypassed), got %d", callCount)
	}
	if resp1.Count != 1 || resp2.Count != 2 {
		t.Errorf("expected fresh responses 1 and 2, got %d and %d", resp1.Count, resp2.Count)
	}
	if c.CacheSize() != 0 {
		t.Errorf("expected empty cache after DoRequestNoCache requests, got %d entries", c.CacheSize())
	}

	// An existing cache entry is not read by DoRequestNoCache
	var resp3, resp4 response
	if err := c.DoRequest(context.Background(), http.MethodGet, "/signed", nil, &resp3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := c.DoRequestNoCache(context.Background(), http.MethodGet, "/signed", nil, &resp4); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp4.Count != 4 {
		t.Errorf("expected DoRequestNoCache to bypass existing cache entry, got count %d", resp4.Count)
	}
}

// plainHTTPClient is an HTTPClient implementation without NoCacheRequester
// support, emulating a downstream custom client.
type plainHTTPClient struct {
	calls int
}

func (p *plainHTTPClient) DoRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	p.calls++
	return nil
}

func TestDoRequestNoCache_FallbackForPlainClient(t *testing.T) {
	p := &plainHTTPClient{}

	// NoCacheRequesterを実装しないクライアントでは通常のDoRequestにフォールバックする
	err := DoRequestNoCache(context.Background(), p, http.MethodGet, "/signed", nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.calls != 1 {
		t.Errorf("expected fallback to DoRequest exactly once, got %d calls", p.calls)
	}
}

func TestClient_DoRequest_CacheDisabled(t *testing.T) {
	var callCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "hello"})
	}))
	defer server.Close()

	c := NewClient("test-api-key") // no cache
	c.baseURL = server.URL

	type response struct {
		Message string `json:"message"`
	}

	// First request
	var resp1 response
	err := c.DoRequest(context.Background(), http.MethodGet, "/test", nil, &resp1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	// Second request (called again since no cache)
	var resp2 response
	err = c.DoRequest(context.Background(), http.MethodGet, "/test", nil, &resp2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (no cache), got %d", callCount)
	}
}

func TestClient_DoRequest_POSTNotCached(t *testing.T) {
	var callCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "hello"})
	}))
	defer server.Close()

	c := NewClient("test-api-key", WithCache())
	c.baseURL = server.URL

	type response struct {
		Message string `json:"message"`
	}

	// First POST request
	var resp1 response
	err := c.DoRequest(context.Background(), http.MethodPost, "/test", nil, &resp1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	// Second POST request (POST is not cached)
	var resp2 response
	err = c.DoRequest(context.Background(), http.MethodPost, "/test", nil, &resp2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (POST not cached), got %d", callCount)
	}
}

func TestClient_DoRequest_CacheWithQueryParams(t *testing.T) {
	var callCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"path": r.URL.RequestURI()})
	}))
	defer server.Close()

	c := NewClient("test-api-key", WithCache())
	c.baseURL = server.URL

	type response struct {
		Path string `json:"path"`
	}

	// Different query parameters are separate cache entries
	var resp1 response
	err := c.DoRequest(context.Background(), http.MethodGet, "/test?code=7203", nil, &resp1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 1 {
		t.Errorf("expected 1 call, got %d", callCount)
	}

	var resp2 response
	err = c.DoRequest(context.Background(), http.MethodGet, "/test?code=9984", nil, &resp2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (different query params), got %d", callCount)
	}

	// Same query parameters result in cache hit
	var resp3 response
	err = c.DoRequest(context.Background(), http.MethodGet, "/test?code=7203", nil, &resp3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if callCount != 2 {
		t.Errorf("expected 2 calls (cached), got %d", callCount)
	}
}

func TestClient_ClearCache(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "hello"})
	}))
	defer server.Close()

	c := NewClient("test-api-key", WithCache())
	c.baseURL = server.URL

	// Populate cache
	var resp struct{ Message string }
	_ = c.DoRequest(context.Background(), http.MethodGet, "/test1", nil, &resp)
	_ = c.DoRequest(context.Background(), http.MethodGet, "/test2", nil, &resp)

	if c.CacheSize() != 2 {
		t.Errorf("expected cache size 2, got %d", c.CacheSize())
	}

	// Clear cache
	c.ClearCache()

	if c.CacheSize() != 0 {
		t.Errorf("expected cache size 0 after clear, got %d", c.CacheSize())
	}
}

func TestClient_CacheSize_Disabled(t *testing.T) {
	c := NewClient("test-api-key") // no cache
	if c.CacheSize() != 0 {
		t.Errorf("expected cache size 0 when disabled, got %d", c.CacheSize())
	}
}

func TestClient_ClearCache_Disabled(t *testing.T) {
	c := NewClient("test-api-key") // no cache
	// Verify it does not panic
	c.ClearCache()
}

func TestClient_DoRequest_ConcurrentAccess(t *testing.T) {
	var callCount int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&callCount, 1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "hello"})
	}))
	defer server.Close()

	c := NewClient("test-api-key", WithCache())
	c.baseURL = server.URL

	type response struct {
		Message string `json:"message"`
	}

	// Send requests concurrently
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var resp response
			err := c.DoRequest(context.Background(), http.MethodGet, "/test", nil, &resp)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		}()
	}
	wg.Wait()

	// With singleflight, concurrent requests for the same key result in only 1 API call
	if atomic.LoadInt64(&callCount) != 1 {
		t.Errorf("expected 1 call (singleflight), got %d", callCount)
	}
	if c.CacheSize() != 1 {
		t.Errorf("expected cache size 1, got %d", c.CacheSize())
	}
}

func TestClient_DoRequest_ContextCancellation(t *testing.T) {
	var callCount int64
	release := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&callCount, 1)
		<-release
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "hello"})
	}))
	defer server.Close()

	c := NewClient("test-api-key", WithCache())
	c.baseURL = server.URL

	type response struct {
		Message string `json:"message"`
	}

	// Start a request and cancel its context while the request is in flight
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := make(chan error, 1)
	go func() {
		var resp response
		errCh <- c.DoRequest(ctx, http.MethodGet, "/test", nil, &resp)
	}()

	// Wait until the request reaches the server
	deadline := time.Now().Add(5 * time.Second)
	for atomic.LoadInt64(&callCount) == 0 {
		if time.Now().After(deadline) {
			t.Fatal("request did not reach the server")
		}
		time.Sleep(time.Millisecond)
	}
	cancel()

	// The canceled caller returns immediately with the context error
	select {
	case err := <-errCh:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("expected context.Canceled, got %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("canceled caller did not return")
	}

	// The detached flight keeps running and populates the cache
	close(release)
	deadline = time.Now().Add(5 * time.Second)
	for c.CacheSize() == 0 {
		if time.Now().After(deadline) {
			t.Fatal("flight did not populate the cache after cancellation")
		}
		time.Sleep(time.Millisecond)
	}

	// A subsequent request is served from the cache without a new API call
	var resp response
	if err := c.DoRequest(context.Background(), http.MethodGet, "/test", nil, &resp); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Message != "hello" {
		t.Errorf("expected message 'hello', got '%s'", resp.Message)
	}
	if got := atomic.LoadInt64(&callCount); got != 1 {
		t.Errorf("expected 1 call, got %d", got)
	}
}
