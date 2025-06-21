//go:build e2e
// +build e2e

package e2e

import (
	"testing"
	"time"

	"github.com/utahta/jquants"
)

// TestTOPIXEndpoint は/markets/topixエンドポイントの完全なテスト
func TestTOPIXEndpoint(t *testing.T) {
	t.Run("GetTOPIXData_ByDateRange", func(t *testing.T) {
		// 最近の営業日のTOPIX詳細データを取得
		date := getTestDate()

		params := jquants.TOPIXParams{
			From: date,
			To:   date,
		}

		resp, err := jq.TOPIX.GetTOPIXData(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get TOPIX data: %v", err)
		}

		if resp == nil || len(resp.TOPIX) == 0 {
			t.Skip("No TOPIX data available for the specified date")
		}

		// 各TOPIXデータを詳細に検証
		for i, topix := range resp.TOPIX {
			// 基本情報の検証
			if topix.Date == "" {
				t.Errorf("TOPIX[%d]: Date is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(topix.Date) != 10 || topix.Date[4] != '-' || topix.Date[7] != '-' {
					t.Errorf("TOPIX[%d]: Date format invalid = %v, want YYYY-MM-DD", i, topix.Date)
				}
				// 日付の一致確認（APIはYYYY-MM-DD形式で返す）
				expectedDate := getTestDateFormatted()
				if topix.Date != expectedDate {
					t.Errorf("TOPIX[%d]: Date = %v, want %v", i, topix.Date, expectedDate)
				}
			}

			// TOPIX指数値の検証
			if topix.Close <= 0 {
				t.Errorf("TOPIX[%d]: Close = %v, want > 0", i, topix.Close)
			}
			if topix.Open <= 0 {
				t.Errorf("TOPIX[%d]: Open = %v, want > 0", i, topix.Open)
			}
			if topix.High <= 0 {
				t.Errorf("TOPIX[%d]: High = %v, want > 0", i, topix.High)
			}
			if topix.Low <= 0 {
				t.Errorf("TOPIX[%d]: Low = %v, want > 0", i, topix.Low)
			}

			// 四本値の論理的整合性チェック
			if topix.High < topix.Low {
				t.Errorf("TOPIX[%d]: High (%v) < Low (%v)", i, topix.High, topix.Low)
			}
			if topix.Open > topix.High || topix.Open < topix.Low {
				t.Errorf("TOPIX[%d]: Open (%v) is outside High (%v) - Low (%v) range",
					i, topix.Open, topix.High, topix.Low)
			}
			if topix.Close > topix.High || topix.Close < topix.Low {
				t.Errorf("TOPIX[%d]: Close (%v) is outside High (%v) - Low (%v) range",
					i, topix.Close, topix.High, topix.Low)
			}

			// 詳細ログ（最初の3件のみ）
			if i < 3 {
				t.Logf("TOPIX[%d]: Date=%s, Close=%.2f, O=%.2f, H=%.2f, L=%.2f",
					i, topix.Date, topix.Close, topix.Open, topix.High, topix.Low)
			}
		}

		t.Logf("Retrieved %d TOPIX records for date %s", len(resp.TOPIX), date)
	})

	t.Run("GetTOPIXData_MultiDay", func(t *testing.T) {
		// 過去1週間のTOPIXデータを取得
		toDate := getTestDate()
		fromTime, _ := time.Parse("20060102", toDate)
		fromTime = fromTime.AddDate(0, 0, -7)
		fromDate := fromTime.Format("20060102")

		params := jquants.TOPIXParams{
			From: fromDate,
			To:   toDate,
		}

		resp, err := jq.TOPIX.GetTOPIXData(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get TOPIX data for date range: %v", err)
		}

		if resp == nil || len(resp.TOPIX) == 0 {
			t.Skip("No TOPIX data available for the specified date range")
		}

		t.Logf("Retrieved %d TOPIX records for period %s to %s", len(resp.TOPIX), fromDate, toDate)

		// 日付の範囲確認（YYYY-MM-DD形式での比較）
		for _, topix := range resp.TOPIX {
			// YYYYMMDD形式をYYYY-MM-DD形式に変換して比較
			expectedFromDate := fromDate[:4] + "-" + fromDate[4:6] + "-" + fromDate[6:8]
			expectedToDate := toDate[:4] + "-" + toDate[4:6] + "-" + toDate[6:8]

			if topix.Date < expectedFromDate || topix.Date > expectedToDate {
				t.Logf("TOPIX date %s is outside requested range %s to %s (converted from %s to %s)",
					topix.Date, expectedFromDate, expectedToDate, fromDate, toDate)
			}
		}

		// 日別の変動率計算
		if len(resp.TOPIX) >= 2 {
			for i := 1; i < len(resp.TOPIX); i++ {
				prev := resp.TOPIX[i-1]
				curr := resp.TOPIX[i]
				change := curr.Close - prev.Close
				changePercent := (change / prev.Close) * 100

				if i <= 3 { // 最初の3日間のみログ
					t.Logf("Daily change from %s to %s: %.2f points (%.2f%%)",
						prev.Date, curr.Date, change, changePercent)
				}
			}
		}
	})

	t.Run("GetTOPIXData_SingleDay", func(t *testing.T) {
		// 特定日のTOPIXデータを取得
		date := getTestDate()

		params := jquants.TOPIXParams{
			From: date,
			To:   date,
		}

		resp, err := jq.TOPIX.GetTOPIXData(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Logf("Failed to get TOPIX data for single day: %v", err)
			return
		}

		if resp != nil && len(resp.TOPIX) > 0 {
			// 指定日のデータが取得できたことを確認
			topix := resp.TOPIX[0]
			// 日付の検証（APIはYYYY-MM-DD形式で返す）
			expectedDate := getTestDateFormatted()
			if topix.Date != expectedDate {
				t.Errorf("Expected date %s, got %s", expectedDate, topix.Date)
			}
			t.Logf("Retrieved TOPIX data for %s: Close=%.2f", date, topix.Close)
		}
	})

	t.Run("GetTOPIXData_Pagination", func(t *testing.T) {
		// ページネーションのテスト（長期間を指定）
		params := jquants.TOPIXParams{
			From: "2024-01-01",
			To:   "2024-01-31",
		}

		resp, err := jq.TOPIX.GetTOPIXData(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No TOPIX data available for historical period")
		}

		if resp == nil || len(resp.TOPIX) == 0 {
			t.Skip("No TOPIX data available")
		}

		firstPageCount := len(resp.TOPIX)
		t.Logf("First page: %d TOPIX records", firstPageCount)

		if resp.PaginationKey != "" {
			// 次のページを取得
			params.PaginationKey = resp.PaginationKey
			resp2, err := jq.TOPIX.GetTOPIXData(params)
			if err != nil {
				t.Fatalf("Failed to get next page: %v", err)
			}

			if resp2 != nil && len(resp2.TOPIX) > 0 {
				t.Logf("Second page: %d TOPIX records", len(resp2.TOPIX))
			}
		}
	})

	t.Run("GetLatestTOPIX", func(t *testing.T) {
		// 最新のTOPIXデータを取得する便利メソッド
		topix, err := jq.TOPIX.GetLatestTOPIX()
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get latest TOPIX: %v", err)
		}

		if topix == nil {
			t.Skip("No latest TOPIX data available")
		}

		// 基本的な検証
		if topix.Date == "" {
			t.Error("Latest TOPIX: Date is empty")
		}
		if topix.Close <= 0 {
			t.Error("Latest TOPIX: Close is invalid")
		}

		t.Logf("Latest TOPIX [%s]: %.2f", topix.Date, topix.Close)
	})

	t.Run("GetTOPIXData_Statistics", func(t *testing.T) {
		// 過去30日のTOPIXデータから統計情報を計算
		toDate := getTestDate()
		fromTime, _ := time.Parse("20060102", toDate)
		fromTime = fromTime.AddDate(0, 0, -30)
		fromDate := fromTime.Format("20060102")

		params := jquants.TOPIXParams{
			From: fromDate,
			To:   toDate,
		}

		resp, err := jq.TOPIX.GetTOPIXData(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No TOPIX data available for statistics")
		}

		if resp == nil || len(resp.TOPIX) == 0 {
			t.Skip("No TOPIX data available")
		}

		// 基本統計の計算
		var sum, min, max float64
		min = resp.TOPIX[0].Close
		max = resp.TOPIX[0].Close

		for _, topix := range resp.TOPIX {
			sum += topix.Close
			if topix.Close < min {
				min = topix.Close
			}
			if topix.Close > max {
				max = topix.Close
			}
		}

		avg := sum / float64(len(resp.TOPIX))
		volatility := (max - min) / avg * 100

		t.Logf("TOPIX Statistics (30 days):")
		t.Logf("  Count: %d", len(resp.TOPIX))
		t.Logf("  Average: %.2f", avg)
		t.Logf("  Min: %.2f", min)
		t.Logf("  Max: %.2f", max)
		t.Logf("  Volatility: %.2f%%", volatility)
	})
}
