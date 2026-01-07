package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MockClient is a mock implementation of HTTPClient interface for testing
type MockClient struct {
	Responses    map[string]interface{}
	Errors       map[string]error
	RequestCount int
	LastMethod   string
	LastPath     string
	LastBody     interface{}
}

// NewMockClient creates a new mock client
func NewMockClient() *MockClient {
	return &MockClient{
		Responses: make(map[string]interface{}),
		Errors:    make(map[string]error),
	}
}

// DoRequest implements HTTPClient interface
func (m *MockClient) DoRequest(method, path string, body interface{}, result interface{}) error {
	m.RequestCount++
	m.LastMethod = method
	m.LastPath = path
	m.LastBody = body

	key := fmt.Sprintf("%s:%s", method, path)

	// Check if error is set for this request
	if err, ok := m.Errors[key]; ok {
		return err
	}

	// Check if response is set for this request (exact match)
	if resp, ok := m.Responses[key]; ok {
		// Marshal and unmarshal to simulate JSON encoding/decoding
		jsonData, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonData, result)
	}

	// Check if response is set for this request (prefix match for flexibility)
	for mockKey, resp := range m.Responses {
		if strings.HasPrefix(key, mockKey) {
			jsonData, err := json.Marshal(resp)
			if err != nil {
				return err
			}
			return json.Unmarshal(jsonData, result)
		}
	}

	return fmt.Errorf("no mock response set for %s", key)
}

// SetResponse sets a mock response for a specific method and path
func (m *MockClient) SetResponse(method, path string, response interface{}) {
	key := fmt.Sprintf("%s:%s", method, path)
	m.Responses[key] = response
}

// SetError sets a mock error for a specific method and path
func (m *MockClient) SetError(method, path string, err error) {
	key := fmt.Sprintf("%s:%s", method, path)
	m.Errors[key] = err
}
