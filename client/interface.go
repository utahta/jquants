package client

// HTTPClient defines the interface for making HTTP requests
type HTTPClient interface {
	DoRequest(method, path string, body interface{}, result interface{}) error
}
