package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

const (
	BaseURL = "https://api.jquants.com/v2"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
	// Cache related fields
	cache        map[string][]byte
	cacheEnabled bool
	mu           sync.RWMutex
	sf           singleflight.Group
}

// ClientOption is a function type for configuring Client settings.
type ClientOption func(*Client)

// WithCache enables caching functionality.
// Cache is only applied to GET requests.
func WithCache() ClientOption {
	return func(c *Client) {
		c.cacheEnabled = true
		c.cache = make(map[string][]byte)
	}
}

// NewClient creates a new Client.
// Options can be specified to enable features such as caching.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: BaseURL,
		apiKey:  apiKey,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// NewClientFromEnv creates a Client using the API key from the JQUANTS_API_KEY environment variable.
// Returns an error if the environment variable is not set.
// Options can be specified to enable features such as caching.
func NewClientFromEnv(opts ...ClientOption) (*Client, error) {
	apiKey := os.Getenv("JQUANTS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("JQUANTS_API_KEY environment variable is not set")
	}
	return NewClient(apiKey, opts...), nil
}

func (c *Client) DoRequest(method, path string, body interface{}, result interface{}) error {
	cacheKey := method + ":" + path

	if method == http.MethodGet && c.cacheEnabled {
		// Check cache first
		c.mu.RLock()
		if cached, ok := c.cache[cacheKey]; ok {
			c.mu.RUnlock()
			if result != nil {
				if err := json.Unmarshal(cached, result); err != nil {
					return fmt.Errorf("failed to decode cached response: %w", err)
				}
			}
			return nil
		}
		c.mu.RUnlock()

		// Use singleflight to deduplicate concurrent requests for the same key
		respBody, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
			// Check cache again (another goroutine may have cached it)
			c.mu.RLock()
			if cached, ok := c.cache[cacheKey]; ok {
				c.mu.RUnlock()
				return cached, nil
			}
			c.mu.RUnlock()

			// Perform HTTP request
			data, err := c.doHTTPRequest(method, path, body)
			if err != nil {
				return nil, err
			}

			// Store in cache
			c.mu.Lock()
			c.cache[cacheKey] = data
			c.mu.Unlock()

			return data, nil
		})

		if err != nil {
			return err
		}

		if result != nil {
			if err := json.Unmarshal(respBody.([]byte), result); err != nil {
				return fmt.Errorf("failed to decode response: %w", err)
			}
		}
		return nil
	}

	// Non-cached request (cache disabled or non-GET method)
	respBody, err := c.doHTTPRequest(method, path, body)
	if err != nil {
		return err
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// doHTTPRequest performs the actual HTTP request.
func (c *Client) doHTTPRequest(method, path string, body interface{}) ([]byte, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		req.Header.Set("x-api-key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// ClearCache clears the cache.
func (c *Client) ClearCache() {
	if !c.cacheEnabled {
		return
	}
	c.mu.Lock()
	c.cache = make(map[string][]byte)
	c.mu.Unlock()
}

// CacheSize returns the number of cache entries.
func (c *Client) CacheSize() int {
	if !c.cacheEnabled {
		return 0
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.cache)
}
