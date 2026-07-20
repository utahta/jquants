//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/utahta/jquants"
)

// TestBulkEndpoint は/bulk/list, /bulk/getエンドポイントのテスト
func TestBulkEndpoint(t *testing.T) {
	t.Run("GetFiles_And_GetDownloadURL", func(t *testing.T) {
		resp, err := jq.Bulk.GetFiles(context.Background(), jquants.BulkListParams{
			Endpoint: "/equities/bars/daily",
		})
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get bulk file list: %v", err)
		}
		if len(resp.Data) == 0 {
			t.Skip("No bulk files available for /equities/bars/daily")
		}

		for i, f := range resp.Data {
			if f.Key == "" {
				t.Errorf("File[%d]: Key is empty", i)
			}
			if f.Size <= 0 {
				t.Errorf("File[%d]: Size = %v, want > 0", i, f.Size)
			}
			if f.LastModified == "" {
				t.Errorf("File[%d]: LastModified is empty", i)
			}
			if i >= 10 {
				break
			}
		}
		t.Logf("Retrieved %d bulk files for /equities/bars/daily", len(resp.Data))

		// 先頭ファイルのダウンロードURL取得（URLは5分有効・1回限りのためダウンロードはしない）
		url, err := jq.Bulk.GetDownloadURL(context.Background(), jquants.BulkGetParams{
			Key: resp.Data[0].Key,
		})
		if err != nil {
			t.Fatalf("Failed to get download URL for %s: %v", resp.Data[0].Key, err)
		}
		if url == "" {
			t.Errorf("GetDownloadURL(%s): URL is empty", resp.Data[0].Key)
		}
	})
}
