//go:build e2e
// +build e2e

package e2e

import (
	"os"
	"strings"
	"testing"

	"github.com/utahta/jquants"
	"github.com/utahta/jquants/auth"
	"github.com/utahta/jquants/client"
)

var (
	jq *jquants.JQuantsAPI
)

func TestMain(m *testing.M) {
	// 認証設定
	c := client.NewClient()
	a := auth.NewAuth(c)

	// 環境変数から認証
	if err := a.InitFromEnv(); err != nil {
		panic("Authentication failed: " + err.Error())
	}

	// JQuantsAPI作成
	jq = jquants.NewJQuantsAPI(c)

	// テスト実行
	os.Exit(m.Run())
}

// isSubscriptionLimited はサブスクリプションプランによる制限があるかをチェックする
func isSubscriptionLimited(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	// エラーメッセージにサブスクリプション制限の文言が含まれているかチェック
	return strings.Contains(errMsg, "This API is not available on your subscription") ||
		strings.Contains(errMsg, "not available on your subscription")
}

// getTestDate はテスト用の固定営業日を取得する
func getTestDate() string {
	// 2025年6月13日（金曜日）を返す
	// これにより、テストの一貫性とデータの存在を保証
	return "20250613"
}

// getTestDateFormatted はテスト用の固定営業日をYYYY-MM-DD形式で取得する
func getTestDateFormatted() string {
	// APIレスポンスで使用されるYYYY-MM-DD形式で返す
	return "2025-06-13"
}
