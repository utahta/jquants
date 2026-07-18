package client

import "context"

// HTTPClient defines the interface for making HTTP requests
type HTTPClient interface {
	DoRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error
}
