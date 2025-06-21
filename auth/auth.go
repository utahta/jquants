package auth

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/utahta/jquants/client"
)

type LoginRequest struct {
	Email    string `json:"mailaddress"`
	Password string `json:"password"`
}

type LoginResponse struct {
	RefreshToken string `json:"refreshToken"`
}

type TokenResponse struct {
	IDToken string `json:"idToken"`
}

type Auth struct {
	client       client.HTTPClient
	refreshToken string
	idToken      string
}

func NewAuth(c client.HTTPClient) *Auth {
	return &Auth{
		client: c,
	}
}

// Login authenticates with email and password to get a refresh token
func (a *Auth) Login(email, password string) error {
	loginReq := LoginRequest{
		Email:    email,
		Password: password,
	}

	var resp LoginResponse
	if err := a.client.DoRequest("POST", "/token/auth_user", loginReq, &resp); err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	a.refreshToken = resp.RefreshToken

	// リフレッシュトークンを保存
	if err := a.saveRefreshToken(resp.RefreshToken); err != nil {
		// 保存に失敗しても警告のみ（ログイン自体は成功）
		fmt.Fprintf(os.Stderr, "Warning: failed to save refresh token: %v\n", err)
	}

	return a.GetAccessToken()
}

func (a *Auth) GetAccessToken() error {
	if a.refreshToken == "" {
		return fmt.Errorf("no refresh token available")
	}

	// クエリパラメータとして送信
	path := fmt.Sprintf("/token/auth_refresh?refreshtoken=%s", url.QueryEscape(a.refreshToken))

	var resp TokenResponse
	if err := a.client.DoRequest("POST", path, nil, &resp); err != nil {
		return fmt.Errorf("failed to get access token: %w", err)
	}

	a.idToken = resp.IDToken
	a.client.SetAccessToken(a.idToken)

	return nil
}

// GetRefreshToken returns the current refresh token
func (a *Auth) GetRefreshToken() string {
	return a.refreshToken
}

// SetRefreshToken sets the refresh token
func (a *Auth) SetRefreshToken(token string) {
	a.refreshToken = token
}

// InitFromEnv initializes auth with refresh token from environment variable or config file
func (a *Auth) InitFromEnv() error {
	// 1. まず環境変数から取得を試みる
	refreshToken := os.Getenv("JQUANTS_REFRESH_TOKEN")

	// 2. 環境変数になければ設定ファイルから読み込む
	if refreshToken == "" {
		var err error
		refreshToken, err = a.loadRefreshToken()
		if err != nil {
			// 設定ファイルもなければエラー
			return fmt.Errorf("no refresh token found: set JQUANTS_REFRESH_TOKEN environment variable")
		}
	}

	a.refreshToken = refreshToken
	return a.GetAccessToken()
}

// getConfigPath returns the path to the config file
func (a *Auth) getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	// セキュリティ: パスを構築して正規化
	configDir := filepath.Join(home, ".jquants")
	configPath := filepath.Join(configDir, "refresh_token")
	
	// パスが適切なディレクトリ内にあることを確認
	cleanPath := filepath.Clean(configPath)
	if !strings.HasPrefix(cleanPath, filepath.Clean(home)) {
		return "", fmt.Errorf("invalid config path")
	}
	
	return cleanPath, nil
}

// saveRefreshToken saves the refresh token to config file
func (a *Auth) saveRefreshToken(token string) error {
	configPath, err := a.getConfigPath()
	if err != nil {
		return err
	}
	
	// ディレクトリを作成
	// #nosec G301 -- directory permissions are appropriately restrictive (0700)
	if err := os.MkdirAll(filepath.Dir(configPath), 0700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// トークンを保存
	// #nosec G304 -- configPath is validated in getConfigPath() to ensure it's within the user's home directory
	if err := os.WriteFile(configPath, []byte(token), 0600); err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}

	return nil
}

// loadRefreshToken loads the refresh token from config file
func (a *Auth) loadRefreshToken() (string, error) {
	configPath, err := a.getConfigPath()
	if err != nil {
		return "", err
	}
	
	// #nosec G304 -- configPath is validated in getConfigPath() to ensure it's within the user's home directory
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("refresh token not found")
		}
		return "", fmt.Errorf("failed to read refresh token: %w", err)
	}

	return string(data), nil
}
