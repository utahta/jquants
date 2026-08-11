//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"testing"

	"github.com/utahta/jquants"
)

// TestEarningsDateEndpoint は/fins/earnings-dateエンドポイントのテスト
// 全プランで利用可能。銘柄指定では予定日の公表履歴が返る。
func TestEarningsDateEndpoint(t *testing.T) {
	validQuarters := map[string]bool{"1Q": true, "2Q": true, "3Q": true, "FY": true}

	// 公表履歴のうち最後に公表されたレコード。後続のサブテストで使用する
	var latest struct {
		PubDate string
		SchDate string
	}

	t.Run("GetEarningsDates_ByCode", func(t *testing.T) {
		dates, err := jq.EarningsDate.GetEarningsDatesByCode(context.Background(), "7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get earnings dates: %v", err)
		}
		if len(dates) == 0 {
			t.Fatal("No earnings dates for 7203")
		}

		for i, d := range dates {
			if d.Code != "72030" {
				t.Errorf("EarningsDate[%d]: Code = %v, want 72030", i, d.Code)
			}
			if len(d.PubDate) != 10 || d.PubDate[4] != '-' || d.PubDate[7] != '-' {
				t.Errorf("EarningsDate[%d]: PubDate = %v, want YYYY-MM-DD", i, d.PubDate)
			}
			// 未定の場合は空文字
			if d.SchDate != "" && (len(d.SchDate) != 10 || d.SchDate[4] != '-' || d.SchDate[7] != '-') {
				t.Errorf("EarningsDate[%d]: SchDate = %v, want YYYY-MM-DD or empty", i, d.SchDate)
			}
			if !validQuarters[d.FQName] {
				t.Errorf("EarningsDate[%d]: FQName = %v, want 1Q/2Q/3Q/FY", i, d.FQName)
			}
			if len(d.FYE) != 4 {
				t.Errorf("EarningsDate[%d]: FYE = %v, want MMDD", i, d.FYE)
			}
			if d.CoName == "" {
				t.Errorf("EarningsDate[%d]: CoName is empty", i)
			}

			if d.PubDate > latest.PubDate {
				latest.PubDate = d.PubDate
				latest.SchDate = d.SchDate
			}
		}
		t.Logf("Retrieved %d earnings date records for 7203 (latest published %s -> %s)",
			len(dates), latest.PubDate, latest.SchDate)
	})

	t.Run("GetEarningsDates_ByDate", func(t *testing.T) {
		if latest.PubDate == "" {
			t.Skip("No published record available from GetEarningsDates_ByCode")
		}

		dates, err := jq.EarningsDate.GetEarningsDatesByDate(context.Background(), latest.PubDate)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get earnings dates for %s: %v", latest.PubDate, err)
		}

		// 7203の予定日が公表された日なので、その銘柄が必ず含まれる
		if !containsCode(dates, "72030") {
			t.Errorf("Earnings dates published on %s do not contain 72030", latest.PubDate)
		}
		t.Logf("Retrieved %d records published on %s", len(dates), latest.PubDate)
	})

	t.Run("GetEarningsDates_ByScheduledDate", func(t *testing.T) {
		if latest.SchDate == "" {
			t.Skip("Latest scheduled date for 7203 is undetermined")
		}

		dates, err := jq.EarningsDate.GetEarningsDatesByScheduledDate(context.Background(), latest.SchDate)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get earnings dates for scheduled date %s: %v", latest.SchDate, err)
		}

		for i, d := range dates {
			if d.SchDate != latest.SchDate {
				t.Errorf("EarningsDate[%d]: SchDate = %v, want %v", i, d.SchDate, latest.SchDate)
			}
		}

		// scheduled_dateは決算区分ごとに最後に公表されたレコードだけを返す。
		// 参照可能期間が遅延しているプランではより新しい変更が見えず、
		// 変更前の予定日で引くことになるため、72030の包含は保証されない。
		t.Logf("Retrieved %d records scheduled on %s (contains 72030: %v)",
			len(dates), latest.SchDate, containsCode(dates, "72030"))
	})
}

func containsCode(dates []jquants.EarningsDate, code string) bool {
	for _, d := range dates {
		if d.Code == code {
			return true
		}
	}
	return false
}
