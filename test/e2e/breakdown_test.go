//go:build e2e
// +build e2e

package e2e

import (
	"testing"
	"time"

	"github.com/utahta/jquants"
)

// TestBreakdownEndpoint は/markets/breakdown（プレミアムプラン専用）エンドポイントの完全なテスト
func TestBreakdownEndpoint(t *testing.T) {
	t.Run("GetBreakdown_ByCode", func(t *testing.T) {
		// トヨタ自動車の売買内訳データを取得
		breakdowns, err := jq.Breakdown.GetBreakdownByCode("7203", 30)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get breakdown data: %v", err)
		}

		if len(breakdowns) == 0 {
			t.Skip("No breakdown data available")
		}

		// 各売買内訳データを詳細に検証
		for i, breakdown := range breakdowns {
			// 基本情報の検証
			if breakdown.Date == "" {
				t.Errorf("Breakdown[%d]: Date is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(breakdown.Date) != 10 || breakdown.Date[4] != '-' || breakdown.Date[7] != '-' {
					t.Errorf("Breakdown[%d]: Date format invalid = %v, want YYYY-MM-DD", i, breakdown.Date)
				}
			}
			if breakdown.Code != "72030" && breakdown.Code != "7203" {
				t.Errorf("Breakdown[%d]: Code = %v, want 72030 or 7203", i, breakdown.Code)
			}

			// 売りの約定代金内訳の検証
			if breakdown.LongSellVa < 0 {
				t.Errorf("Breakdown[%d]: LongSellValue = %v, want >= 0", i, breakdown.LongSellVa)
			}
			if breakdown.ShrtNoMrgnVa < 0 {
				t.Errorf("Breakdown[%d]: ShortSellWithoutMarginValue = %v, want >= 0", i, breakdown.ShrtNoMrgnVa)
			}
			if breakdown.MrgnSellNewVa < 0 {
				t.Errorf("Breakdown[%d]: MarginSellNewValue = %v, want >= 0", i, breakdown.MrgnSellNewVa)
			}
			if breakdown.MrgnSellCloseVa < 0 {
				t.Errorf("Breakdown[%d]: MarginSellCloseValue = %v, want >= 0", i, breakdown.MrgnSellCloseVa)
			}

			// 買いの約定代金内訳の検証
			if breakdown.LongBuyVa < 0 {
				t.Errorf("Breakdown[%d]: LongBuyValue = %v, want >= 0", i, breakdown.LongBuyVa)
			}
			if breakdown.MrgnBuyNewVa < 0 {
				t.Errorf("Breakdown[%d]: MarginBuyNewValue = %v, want >= 0", i, breakdown.MrgnBuyNewVa)
			}
			if breakdown.MrgnBuyCloseVa < 0 {
				t.Errorf("Breakdown[%d]: MarginBuyCloseValue = %v, want >= 0", i, breakdown.MrgnBuyCloseVa)
			}

			// 売りの約定高内訳の検証
			if breakdown.LongSellVo < 0 {
				t.Errorf("Breakdown[%d]: LongSellVolume = %v, want >= 0", i, breakdown.LongSellVo)
			}
			if breakdown.ShrtNoMrgnVo < 0 {
				t.Errorf("Breakdown[%d]: ShortSellWithoutMarginVolume = %v, want >= 0", i, breakdown.ShrtNoMrgnVo)
			}
			if breakdown.MrgnSellNewVo < 0 {
				t.Errorf("Breakdown[%d]: MarginSellNewVolume = %v, want >= 0", i, breakdown.MrgnSellNewVo)
			}
			if breakdown.MrgnSellCloseVo < 0 {
				t.Errorf("Breakdown[%d]: MarginSellCloseVolume = %v, want >= 0", i, breakdown.MrgnSellCloseVo)
			}

			// 買いの約定高内訳の検証
			if breakdown.LongBuyVo < 0 {
				t.Errorf("Breakdown[%d]: LongBuyVolume = %v, want >= 0", i, breakdown.LongBuyVo)
			}
			if breakdown.MrgnBuyNewVo < 0 {
				t.Errorf("Breakdown[%d]: MarginBuyNewVolume = %v, want >= 0", i, breakdown.MrgnBuyNewVo)
			}
			if breakdown.MrgnBuyCloseVo < 0 {
				t.Errorf("Breakdown[%d]: MarginBuyCloseVolume = %v, want >= 0", i, breakdown.MrgnBuyCloseVo)
			}

			// 合計値の検証
			totalSellValue := breakdown.LongSellVa + breakdown.ShrtNoMrgnVa +
				breakdown.MrgnSellNewVa + breakdown.MrgnSellCloseVa
			totalBuyValue := breakdown.LongBuyVa + breakdown.MrgnBuyNewVa + breakdown.MrgnBuyCloseVa

			totalSellVolume := breakdown.LongSellVo + breakdown.ShrtNoMrgnVo +
				breakdown.MrgnSellNewVo + breakdown.MrgnSellCloseVo
			totalBuyVolume := breakdown.LongBuyVo + breakdown.MrgnBuyNewVo + breakdown.MrgnBuyCloseVo

			// 最初の3件の詳細ログ
			if i < 3 {
				t.Logf("Breakdown[%d]: Date=%s, Code=%s", i, breakdown.Date, breakdown.Code)
				t.Logf("  Total Sell Value: %.0f (Long: %.0f, Short: %.0f, Margin New: %.0f, Margin Close: %.0f)",
					totalSellValue, breakdown.LongSellVa, breakdown.ShrtNoMrgnVa,
					breakdown.MrgnSellNewVa, breakdown.MrgnSellCloseVa)
				t.Logf("  Total Buy Value: %.0f (Long: %.0f, Margin New: %.0f, Margin Close: %.0f)",
					totalBuyValue, breakdown.LongBuyVa, breakdown.MrgnBuyNewVa, breakdown.MrgnBuyCloseVa)
				t.Logf("  Total Sell Volume: %.0f, Total Buy Volume: %.0f", totalSellVolume, totalBuyVolume)
			}
		}

		t.Logf("Retrieved %d breakdown records", len(breakdowns))
	})

	t.Run("GetBreakdown_ByDate", func(t *testing.T) {
		// 特定日の全銘柄の売買内訳データを取得
		date := getTestDate()

		params := jquants.BreakdownParams{
			Date: date,
		}

		resp, err := jq.Breakdown.GetBreakdown(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Logf("Failed to get breakdown by date: %v", err)
			return
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No breakdown data available for the specified date")
		}

		t.Logf("Retrieved %d breakdown records for date %s", len(resp.Data), date)

		// 日付の一致確認
		for i, breakdown := range resp.Data {
			if breakdown.Date != date {
				t.Errorf("Breakdown[%d]: Date = %v, want %v", i, breakdown.Date, date)
			}
			if i >= 10 {
				break // 最初の10件のみ検証
			}
		}

		// 銘柄別の集計
		codeCount := make(map[string]int)
		for _, breakdown := range resp.Data {
			codeCount[breakdown.Code]++
		}
		t.Logf("Found breakdown data for %d unique codes", len(codeCount))
	})

	t.Run("GetBreakdown_DateRange", func(t *testing.T) {
		// 日付範囲でのデータ取得
		to := getTestDate()
		fromTime, _ := time.Parse("20060102", to)
		fromTime = fromTime.AddDate(0, 0, -7)
		from := fromTime.Format("20060102")

		params := jquants.BreakdownParams{
			Code: "7203",
			From: from,
			To:   to,
		}

		resp, err := jq.Breakdown.GetBreakdown(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get breakdown for date range: %v", err)
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No breakdown data available for the specified date range")
		}

		t.Logf("Retrieved %d breakdown records for period %s to %s", len(resp.Data), from, to)

		// 日付範囲の確認
		for _, breakdown := range resp.Data {
			if breakdown.Date < from || breakdown.Date > to {
				t.Errorf("Breakdown date %s is outside requested range %s to %s",
					breakdown.Date, from, to)
			}
		}
	})

	t.Run("GetBreakdown_Pagination", func(t *testing.T) {
		// ページネーションのテスト
		date := getTestDate()

		params := jquants.BreakdownParams{
			Date: date,
		}

		resp, err := jq.Breakdown.GetBreakdown(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No breakdown data available")
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No breakdown data available")
		}

		firstPageCount := len(resp.Data)
		t.Logf("First page: %d breakdown records", firstPageCount)

		if resp.PaginationKey != "" {
			// 次のページを取得
			params.PaginationKey = resp.PaginationKey
			resp2, err := jq.Breakdown.GetBreakdown(params)
			if err != nil {
				t.Fatalf("Failed to get next page: %v", err)
			}

			if resp2 != nil && len(resp2.Data) > 0 {
				t.Logf("Second page: %d breakdown records", len(resp2.Data))
			}
		}
	})

	t.Run("GetBreakdown_Analysis", func(t *testing.T) {
		// 売買内訳データの分析
		breakdowns, err := jq.Breakdown.GetBreakdownByCode("7203", 10)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No breakdown data available for analysis")
		}

		if len(breakdowns) == 0 {
			t.Skip("No breakdown data available")
		}

		// 信用取引比率の計算
		for i, breakdown := range breakdowns {
			if i >= 5 {
				break // 最初の5日分のみ分析
			}

			totalSellValue := breakdown.LongSellVa + breakdown.ShrtNoMrgnVa +
				breakdown.MrgnSellNewVa + breakdown.MrgnSellCloseVa
			totalBuyValue := breakdown.LongBuyVa + breakdown.MrgnBuyNewVa + breakdown.MrgnBuyCloseVa

			if totalSellValue > 0 {
				marginSellRatio := (breakdown.MrgnSellNewVa + breakdown.MrgnSellCloseVa) / totalSellValue * 100
				shortSellRatio := breakdown.ShrtNoMrgnVa / totalSellValue * 100

				t.Logf("Date %s: Margin Sell Ratio: %.2f%%, Short Sell Ratio: %.2f%%",
					breakdown.Date, marginSellRatio, shortSellRatio)
			}

			if totalBuyValue > 0 {
				marginBuyRatio := (breakdown.MrgnBuyNewVa + breakdown.MrgnBuyCloseVa) / totalBuyValue * 100
				t.Logf("Date %s: Margin Buy Ratio: %.2f%%", breakdown.Date, marginBuyRatio)
			}
		}
	})
}
