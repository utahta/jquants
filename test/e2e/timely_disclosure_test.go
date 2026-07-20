//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/utahta/jquants"
)

// TestTimelyDisclosureEndpoint は/td/list, /td/files, /td/bulkエンドポイントのテスト
// TDnet/適時開示情報アドオン契約が必要。
func TestTimelyDisclosureEndpoint(t *testing.T) {
	t.Run("GetDisclosures_ByDate", func(t *testing.T) {
		date := getTestDateFormatted()
		disclosures, err := jq.TimelyDisclosure.GetDisclosuresByDate(context.Background(), date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (TDnet add-on required)")
			}
			t.Fatalf("Failed to get disclosures: %v", err)
		}
		if len(disclosures) == 0 {
			t.Skipf("No disclosures on %s", date)
		}

		for i, d := range disclosures {
			if len(d.DiscNo) != 14 {
				t.Errorf("Disclosure[%d]: DiscNo = %v, want 14 digits", i, d.DiscNo)
			}
			if d.DiscDate != date {
				t.Errorf("Disclosure[%d]: DiscDate = %v, want %v", i, d.DiscDate, date)
			}
			if d.Title == "" {
				t.Errorf("Disclosure[%d]: Title is empty", i)
			}
			if len(d.DiscItems) == 0 {
				t.Errorf("Disclosure[%d]: DiscItems is empty", i)
			}
			if i >= 20 {
				break
			}
		}
		t.Logf("Retrieved %d disclosures for %s", len(disclosures), date)

		// 先頭の開示についてファイルURLの取得も検証
		files, err := jq.TimelyDisclosure.GetDisclosureFiles(context.Background(),
			jquants.TimelyDisclosureFilesParams{DiscNo: disclosures[0].DiscNo})
		if err != nil {
			t.Fatalf("Failed to get disclosure files for %s: %v", disclosures[0].DiscNo, err)
		}
		if files.Files.PDF == "" && files.Files.SummaryPDF == "" && files.Files.XBRL == "" {
			t.Errorf("GetDisclosureFiles(%s): all file URLs are empty", disclosures[0].DiscNo)
		}
	})

	t.Run("GetBulkFile", func(t *testing.T) {
		bulk, err := jq.TimelyDisclosure.GetBulkFile(context.Background())
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (TDnet add-on required)")
			}
			t.Fatalf("Failed to get bulk file: %v", err)
		}
		if bulk.URL == "" {
			t.Error("GetBulkFile(): URL is empty")
		}
		if bulk.LastUpdated == "" {
			t.Error("GetBulkFile(): LastUpdated is empty")
		}
		t.Logf("Bulk CSV last updated: %s", bulk.LastUpdated)
	})
}
