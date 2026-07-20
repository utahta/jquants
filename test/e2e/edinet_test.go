//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"strings"
	"testing"
)

// TestEdinetMajorShareholdersEndpoint は/edinet/major-shareholdersエンドポイントのテスト
// Standardプラン以上で利用可能。
// 有価証券報告書の提出は6月に集中し単日指定では空になりやすいため、銘柄指定で取得する。
func TestEdinetMajorShareholdersEndpoint(t *testing.T) {
	t.Run("GetMajorShareholders_ByCode", func(t *testing.T) {
		docs, err := jq.EdinetMajorShareholders.GetMajorShareholdersByCode(context.Background(), "7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get major shareholders: %v", err)
		}
		if len(docs) == 0 {
			t.Fatal("No major shareholders filings for 7203 (data since 2016-06 should exist)")
		}

		for i, doc := range docs {
			if !strings.HasPrefix(doc.DocId, "S") {
				t.Errorf("Doc[%d]: DocId = %v, want S-prefixed", i, doc.DocId)
			}
			if doc.Code != "72030" {
				t.Errorf("Doc[%d]: Code = %v, want 72030", i, doc.Code)
			}
			if doc.SubDate == "" || doc.PerSt == "" || doc.PerEn == "" {
				t.Errorf("Doc[%d]: SubDate/PerSt/PerEn should not be empty", i)
			}
			if len(doc.Hldrs) == 0 {
				t.Errorf("Doc[%d]: Hldrs is empty", i)
			}
			for j, h := range doc.Hldrs {
				if h.HldrName == "" {
					t.Errorf("Doc[%d].Hldrs[%d]: HldrName is empty", i, j)
				}
				if h.ShsRatio <= 0 || h.ShsRatio > 1 {
					t.Errorf("Doc[%d].Hldrs[%d]: ShsRatio = %v, want (0, 1]", i, j, h.ShsRatio)
				}
			}
		}
		t.Logf("Retrieved %d major shareholders filings for 7203", len(docs))
	})
}

// TestEdinetCrossShareholdingsEndpoint は/edinet/cross-shareholdingsエンドポイントのテスト
// Standardプラン以上で利用可能。
func TestEdinetCrossShareholdingsEndpoint(t *testing.T) {
	t.Run("GetCrossShareholdings_ByCode", func(t *testing.T) {
		docs, err := jq.EdinetCrossShareholdings.GetCrossShareholdingsByCode(context.Background(), "7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get cross shareholdings: %v", err)
		}
		if len(docs) == 0 {
			t.Fatal("No cross shareholdings filings for 7203 (data since 2020-03 should exist)")
		}

		reports := 0
		issues := 0
		for i, doc := range docs {
			if !strings.HasPrefix(doc.DocId, "S") {
				t.Errorf("Doc[%d]: DocId = %v, want S-prefixed", i, doc.DocId)
			}
			if doc.Code != "72030" {
				t.Errorf("Doc[%d]: Code = %v, want 72030", i, doc.Code)
			}
			if doc.HasReport() {
				reports++
				issues += len(doc.Report.Spec) + len(doc.Report.Deem)
				if doc.Report.HldrName == "" {
					t.Errorf("Doc[%d]: Report.HldrName is empty", i)
				}
			}
		}
		if reports == 0 {
			t.Error("No document has a Report block for 7203")
		}
		t.Logf("Retrieved %d filings (%d with Report block, %d issues) for 7203",
			len(docs), reports, issues)
	})
}

// TestEdinetLargeVolumeShareholdersEndpoint は/edinet/large-volume-shareholdersエンドポイントのテスト
// Standardプラン以上で利用可能。
func TestEdinetLargeVolumeShareholdersEndpoint(t *testing.T) {
	t.Run("GetLargeVolumeShareholders_ByDate", func(t *testing.T) {
		date := getTestDateFormatted()
		docs, err := jq.EdinetLargeVolumeShareholders.GetLargeVolumeShareholdersByDate(context.Background(), date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get large volume shareholders: %v", err)
		}
		if len(docs) == 0 {
			t.Skipf("No large volume filings on %s", date)
		}

		changeReports := 0
		withTotals := 0
		for i, doc := range docs {
			if !strings.HasPrefix(doc.DocId, "S") {
				t.Errorf("Doc[%d]: DocId = %v, want S-prefixed", i, doc.DocId)
			}
			if doc.SubDate != date {
				t.Errorf("Doc[%d]: SubDate = %v, want %v", i, doc.SubDate, date)
			}
			// 合計欄はnullの書類がある。値がある場合のみ範囲を検証する
			if doc.TotalShsRatio != nil {
				withTotals++
				if *doc.TotalShsRatio <= 0 || *doc.TotalShsRatio > 1 {
					t.Errorf("Doc[%d]: TotalShsRatio = %v, want (0, 1]", i, *doc.TotalShsRatio)
				}
			}
			if len(doc.Hldrs) == 0 {
				t.Errorf("Doc[%d]: Hldrs is empty", i)
			}
			for j, h := range doc.Hldrs {
				if h.ShsRatio < 0 || h.ShsRatio > 1 {
					t.Errorf("Doc[%d].Hldrs[%d]: ShsRatio = %v, want [0, 1]", i, j, h.ShsRatio)
				}
			}
			if doc.IsChangeReport() {
				changeReports++
			}
		}
		t.Logf("Retrieved %d large volume filings (%d change reports, %d with totals) for %s",
			len(docs), changeReports, withTotals, date)
	})
}
