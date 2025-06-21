package auth

import (
	"os"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestAuth_GetAccessToken_WithExistingRefreshToken(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	auth := NewAuth(mockClient)
	auth.SetRefreshToken("existing-refresh-token")

	// Mock response（クエリパラメータを含むパスで設定）
	mockClient.SetResponse("POST", "/token/auth_refresh?refreshtoken=existing-refresh-token", TokenResponse{
		IDToken: "new-id-token",
	})

	// Test
	err := auth.GetAccessToken()
	if err != nil {
		t.Errorf("GetAccessToken failed: %v", err)
	}

	// Verify
	if mockClient.RequestCount != 1 {
		t.Errorf("Expected 1 request, got %d", mockClient.RequestCount)
	}

	if mockClient.AccessToken != "new-id-token" {
		t.Errorf("Expected access token to be set, got %s", mockClient.AccessToken)
	}
}

func TestAuth_TokenGetterSetter(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	auth := NewAuth(mockClient)

	// Test initial state
	if auth.GetRefreshToken() != "" {
		t.Errorf("Expected empty refresh token initially, got %s", auth.GetRefreshToken())
	}

	// Test setter
	auth.SetRefreshToken("test-token")
	if auth.GetRefreshToken() != "test-token" {
		t.Errorf("Expected refresh token to be set, got %s", auth.GetRefreshToken())
	}
}

func TestAuth_GetAccessToken_NoRefreshToken(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	auth := NewAuth(mockClient)

	// Test
	err := auth.GetAccessToken()
	if err == nil {
		t.Errorf("Expected error when no refresh token available, got nil")
	}
}

func TestAuth_InitFromEnv(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	auth := NewAuth(mockClient)

	// Test with no env variable
	os.Unsetenv("JQUANTS_REFRESH_TOKEN")
	err := auth.InitFromEnv()
	if err == nil {
		t.Errorf("Expected error when env variable not set, got nil")
	}

	// Test with env variable set
	os.Setenv("JQUANTS_REFRESH_TOKEN", "test-refresh-token-from-env")
	defer os.Unsetenv("JQUANTS_REFRESH_TOKEN")

	// Mock response（クエリパラメータを含むパスで設定）
	mockClient.SetResponse("POST", "/token/auth_refresh?refreshtoken=test-refresh-token-from-env", TokenResponse{
		IDToken: "test-id-token-from-env",
	})

	// Test
	err = auth.InitFromEnv()
	if err != nil {
		t.Errorf("InitFromEnv failed: %v", err)
	}

	// Verify
	if auth.GetRefreshToken() != "test-refresh-token-from-env" {
		t.Errorf("Expected refresh token from env, got %s", auth.GetRefreshToken())
	}

	if mockClient.AccessToken != "test-id-token-from-env" {
		t.Errorf("Expected access token to be set, got %s", mockClient.AccessToken)
	}
}

func TestAuth_Login(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	auth := NewAuth(mockClient)

	// Mock responses
	// 1. Login response
	mockClient.SetResponse("POST", "/token/auth_user", LoginResponse{
		RefreshToken: "login-refresh-token",
	})
	// 2. GetAccessToken response
	mockClient.SetResponse("POST", "/token/auth_refresh?refreshtoken=login-refresh-token", TokenResponse{
		IDToken: "login-id-token",
	})

	// Test
	err := auth.Login("test@example.com", "password123")
	if err != nil {
		t.Errorf("Login failed: %v", err)
	}

	// Verify
	if auth.GetRefreshToken() != "login-refresh-token" {
		t.Errorf("Expected refresh token from login, got %s", auth.GetRefreshToken())
	}

	if mockClient.AccessToken != "login-id-token" {
		t.Errorf("Expected access token to be set, got %s", mockClient.AccessToken)
	}

	// Verify login request
	if mockClient.LastPath != "/token/auth_refresh?refreshtoken=login-refresh-token" {
		t.Errorf("Expected auth_refresh path, got %s", mockClient.LastPath)
	}
}
