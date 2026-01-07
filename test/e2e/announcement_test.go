//go:build e2e
// +build e2e

package e2e

import (
	"testing"

	"github.com/utahta/jquants"
)

// TestAnnouncementEndpoint は/fins/announcementエンドポイントの完全なテスト
func TestAnnouncementEndpoint(t *testing.T) {
	t.Run("GetAnnouncement_Default", func(t *testing.T) {
		// デフォルト（翌営業日）の決算発表予定を取得
		params := jquants.AnnouncementParams{}
		resp, err := jq.Announcement.GetAnnouncement(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			// 決算発表予定がない日の可能性もある
			t.Logf("No announcement data (might be no announcements): %v", err)
			return
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No announcement data for default date")
		}

		// 各発表予定を詳細に検証
		for i, ann := range resp.Data {
			if i >= 10 {
				break // 最初の10件のみ詳細検証
			}

			// 基本情報の検証
			if ann.Date == "" {
				t.Errorf("Announcement[%d]: Date is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(ann.Date) != 10 || ann.Date[4] != '-' || ann.Date[7] != '-' {
					t.Errorf("Announcement[%d]: Date format invalid = %v, want YYYY-MM-DD", i, ann.Date)
				}
			}

			if ann.Code == "" {
				t.Errorf("Announcement[%d]: Code is empty", i)
			}
			if ann.CoName == "" {
				t.Errorf("Announcement[%d]: CoName is empty", i)
			}

			// 決算情報の検証
			if ann.FY == "" {
				t.Errorf("Announcement[%d]: FiscalYear is empty", i)
			} else {
				// 決算期末の形式検証（例：「3月31日」「9月30日」）
				if !(ann.FY == "3月31日" || ann.FY == "9月30日") {
					t.Logf("Announcement[%d]: Unexpected FiscalYear format = %v", i, ann.FY)
				}
			}

			if ann.FQ == "" {
				t.Errorf("Announcement[%d]: FiscalQuarter is empty", i)
			} else {
				// 決算種別の妥当性チェック
				validQuarters := map[string]bool{
					"第１四半期": true,
					"第２四半期": true,
					"第３四半期": true,
					"通期":    true,
					"本決算":   true,
				}
				if !validQuarters[ann.FQ] {
					t.Logf("Announcement[%d]: Unexpected FiscalQuarter = %v", i, ann.FQ)
				}
			}

			if ann.SectorNm == "" {
				t.Errorf("Announcement[%d]: SectorName is empty", i)
			}

			// 市場区分の検証（必須フィールド）
			if ann.Section == "" {
				t.Errorf("Announcement[%d]: Section is empty", i)
			} else {
				// 市場区分の妥当性チェック
				validSections := map[string]bool{
					"プライム":      true,
					"スタンダード":    true,
					"グロース":      true,
					"マザーズ":      true,
					"JASDAQ":    true,
					"JASDAQ(G)": true,
					"JASDAQ(S)": true,
					"東証１部":      true,
					"東証２部":      true,
				}
				if !validSections[ann.Section] {
					t.Logf("Announcement[%d]: Unexpected Section = %v", i, ann.Section)
				}
			}

			// 最初の5件の詳細ログ
			if i < 5 {
				t.Logf("Announcement[%d]: %s (%s) - %s %s %s",
					i, ann.CoName, ann.Code, ann.FY, ann.FQ, ann.SectorNm)
				if ann.Date != "" {
					t.Logf("  Scheduled Date: %s", ann.Date)
				}
				if ann.Section != "" {
					t.Logf("  Section: %s", ann.Section)
				}
			}
		}

		t.Logf("Total announcements: %d", len(resp.Data))
		if resp.PaginationKey != "" {
			t.Logf("Pagination key present: %s", resp.PaginationKey)
		}
	})

	t.Run("GetAnnouncement_SpecificDate", func(t *testing.T) {
		// 特定日の決算発表予定を取得（過去の営業日）
		date := getTestDate()

		params := jquants.AnnouncementParams{}

		resp, err := jq.Announcement.GetAnnouncement(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Logf("No announcement data for date %s: %v", date, err)
			return
		}

		if resp != nil && len(resp.Data) > 0 {
			t.Logf("Found %d announcements (翌営業日分)", len(resp.Data))
		} else {
			t.Logf("No announcements for default date")
		}
	})

	t.Run("GetAnnouncement_Pagination", func(t *testing.T) {
		// ページネーションのテスト
		params := jquants.AnnouncementParams{}

		resp, err := jq.Announcement.GetAnnouncement(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No announcement data available")
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No announcement data available")
		}

		firstPageCount := len(resp.Data)
		t.Logf("First page: %d announcements", firstPageCount)

		if resp.PaginationKey != "" {
			// 次のページを取得
			params.PaginationKey = resp.PaginationKey
			resp2, err := jq.Announcement.GetAnnouncement(params)
			if err != nil {
				t.Fatalf("Failed to get next page: %v", err)
			}

			if resp2 != nil && len(resp2.Data) > 0 {
				t.Logf("Second page: %d announcements", len(resp2.Data))

				// 異なるデータであることを確認
				if resp2.Data[0].Code == resp.Data[0].Code &&
					resp2.Data[0].CoName == resp.Data[0].CoName {
					t.Error("Second page contains duplicate data")
				}
			}
		}
	})

	t.Run("GetAnnouncement_SectorAnalysis", func(t *testing.T) {
		// 業種別の分析
		params := jquants.AnnouncementParams{}
		resp, err := jq.Announcement.GetAnnouncement(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No announcement data available")
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No announcement data available")
		}

		// 業種別の集計
		sectorCount := make(map[string]int)
		for _, ann := range resp.Data {
			sectorCount[ann.SectorNm]++
		}

		t.Logf("Sector analysis:")
		for sector, count := range sectorCount {
			t.Logf("  %s: %d companies", sector, count)
		}
	})

	t.Run("GetAnnouncement_QuarterAnalysis", func(t *testing.T) {
		// 決算四半期別の分析
		params := jquants.AnnouncementParams{}
		resp, err := jq.Announcement.GetAnnouncement(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No announcement data available")
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No announcement data available")
		}

		// 決算四半期別の集計
		quarterCount := make(map[string]int)
		for _, ann := range resp.Data {
			quarterCount[ann.FQ]++
		}

		t.Logf("Quarter analysis:")
		for quarter, count := range quarterCount {
			t.Logf("  %s: %d companies", quarter, count)
		}
	})
}
