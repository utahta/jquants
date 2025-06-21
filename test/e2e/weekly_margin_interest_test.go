//go:build e2e
// +build e2e

package e2e

import (
	"testing"
	"time"

	"github.com/utahta/jquants"
)

// TestWeeklyMarginInterestEndpoint は/markets/weekly_margin_interestエンドポイントの完全なテスト
func TestWeeklyMarginInterestEndpoint(t *testing.T) {
	t.Run("GetWeeklyMarginInterest_ByCode", func(t *testing.T) {
		// トヨタ自動車の信用取引週末残高を取得
		friday := getRecentFriday()

		interests, err := jq.WeeklyMarginInterest.GetWeeklyMarginInterestByCodeAndDateRange("7203", friday, friday)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get weekly margin interest: %v", err)
		}

		if len(interests) == 0 {
			t.Skip("No weekly margin interest data available")
		}

		// 各週末残高データを詳細に検証
		for i, interest := range interests {
			// 基本情報の検証
			if interest.Code != "72030" && interest.Code != "7203" {
				t.Errorf("Interest[%d]: Code = %v, want 72030 or 7203", i, interest.Code)
			}
			if interest.Date == "" {
				t.Errorf("Interest[%d]: Date is empty", i)
			}

			// 日付フォーマットの検証（YYYY-MM-DD形式）
			if len(interest.Date) != 10 || interest.Date[4] != '-' || interest.Date[7] != '-' {
				t.Errorf("Interest[%d]: Date format invalid = %v, want YYYY-MM-DD", i, interest.Date)
			}

			// IssueTypeの検証（1, 2, 3のいずれか）
			if interest.IssueType != "" {
				validIssueTypes := map[string]bool{"1": true, "2": true, "3": true}
				if !validIssueTypes[interest.IssueType] {
					t.Errorf("Interest[%d]: IssueType = %v, want one of [1, 2, 3]", i, interest.IssueType)
				}
			}

			// 日付が金曜日かチェック
			date, err := time.Parse("2006-01-02", interest.Date)
			if err == nil && date.Weekday() != time.Friday {
				t.Logf("Interest[%d]: Date %s is not Friday (%s)", i, interest.Date, date.Weekday().String())
			}

			// 信用買い残高の検証
			if interest.LongMarginTradeVolume < 0 {
				t.Errorf("Interest[%d]: LongMarginTradeVolume = %v, want >= 0",
					i, interest.LongMarginTradeVolume)
			}
			if interest.LongStandardizedMarginTradeVolume < 0 {
				t.Errorf("Interest[%d]: LongStandardizedMarginTradeVolume = %v, want >= 0",
					i, interest.LongStandardizedMarginTradeVolume)
			}

			// 信用売り残高の検証
			if interest.ShortMarginTradeVolume < 0 {
				t.Errorf("Interest[%d]: ShortMarginTradeVolume = %v, want >= 0",
					i, interest.ShortMarginTradeVolume)
			}
			if interest.ShortStandardizedMarginTradeVolume < 0 {
				t.Errorf("Interest[%d]: ShortStandardizedMarginTradeVolume = %v, want >= 0",
					i, interest.ShortStandardizedMarginTradeVolume)
			}

			// 一般信用残高の検証
			if interest.LongNegotiableMarginTradeVolume < 0 {
				t.Errorf("Interest[%d]: LongNegotiableMarginTradeVolume = %v, want >= 0",
					i, interest.LongNegotiableMarginTradeVolume)
			}
			if interest.ShortNegotiableMarginTradeVolume < 0 {
				t.Errorf("Interest[%d]: ShortNegotiableMarginTradeVolume = %v, want >= 0",
					i, interest.ShortNegotiableMarginTradeVolume)
			}

			// バランス計算
			marginBuyRatio := 0.0
			marginSellRatio := 0.0
			if interest.LongMarginTradeVolume > 0 && interest.ShortMarginTradeVolume > 0 {
				total := interest.LongMarginTradeVolume + interest.ShortMarginTradeVolume
				marginBuyRatio = interest.LongMarginTradeVolume / total * 100
				marginSellRatio = interest.ShortMarginTradeVolume / total * 100
			}

			// 最初の3件の詳細ログ
			if i < 3 {
				t.Logf("Interest[%d]: Date=%s, Code=%s", i, interest.Date, interest.Code)
				t.Logf("  Long Margin: %.0f shares (Std: %.0f, Neg: %.0f)",
					interest.LongMarginTradeVolume,
					interest.LongStandardizedMarginTradeVolume,
					interest.LongNegotiableMarginTradeVolume)
				t.Logf("  Short Margin: %.0f shares (Std: %.0f, Neg: %.0f)",
					interest.ShortMarginTradeVolume,
					interest.ShortStandardizedMarginTradeVolume,
					interest.ShortNegotiableMarginTradeVolume)
				if marginBuyRatio > 0 || marginSellRatio > 0 {
					t.Logf("  Long/Short Ratio: %.1f%% / %.1f%%", marginBuyRatio, marginSellRatio)
				}
			}
		}

		t.Logf("Retrieved %d weekly margin interest records", len(interests))
	})

	t.Run("GetWeeklyMarginInterest_ByDate", func(t *testing.T) {
		// 特定の金曜日の全銘柄データを取得
		friday := getRecentFriday()

		params := jquants.WeeklyMarginInterestParams{
			Date: friday,
		}

		resp, err := jq.WeeklyMarginInterest.GetWeeklyMarginInterest(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Logf("Failed to get weekly margin interest by date: %v", err)
			return
		}

		if resp == nil || len(resp.WeeklyMarginInterest) == 0 {
			t.Skip("No weekly margin interest data for the specified date")
		}

		t.Logf("Retrieved %d weekly margin interest records for %s", len(resp.WeeklyMarginInterest), friday)

		// 日付の一致確認
		for i, interest := range resp.WeeklyMarginInterest {
			if interest.Date != friday {
				t.Errorf("Interest[%d]: Date = %v, want %v", i, interest.Date, friday)
			}

			// 日付フォーマットの検証（YYYY-MM-DD形式）
			if len(interest.Date) != 10 || interest.Date[4] != '-' || interest.Date[7] != '-' {
				t.Errorf("Interest[%d]: Date format invalid = %v, want YYYY-MM-DD", i, interest.Date)
			}

			// IssueTypeの検証
			if interest.IssueType != "" {
				validIssueTypes := map[string]bool{"1": true, "2": true, "3": true}
				if !validIssueTypes[interest.IssueType] {
					t.Errorf("Interest[%d]: IssueType = %v, want one of [1, 2, 3]", i, interest.IssueType)
				}
			}

			if i >= 10 {
				break // 最初の10件のみ確認
			}
		}

		// 銘柄別の集計
		codeCount := make(map[string]int)
		totalMarginBuy := 0.0
		totalMarginSell := 0.0

		for _, interest := range resp.WeeklyMarginInterest {
			codeCount[interest.Code]++
			totalMarginBuy += interest.LongMarginTradeVolume
			totalMarginSell += interest.ShortMarginTradeVolume
		}

		t.Logf("Market summary for %s:", friday)
		t.Logf("  Unique codes: %d", len(codeCount))
		t.Logf("  Total long margin volume: %.0f shares", totalMarginBuy)
		t.Logf("  Total short margin volume: %.0f shares", totalMarginSell)

		if totalMarginBuy > 0 && totalMarginSell > 0 {
			buyRatio := totalMarginBuy / (totalMarginBuy + totalMarginSell) * 100
			t.Logf("  Market buy/sell ratio: %.1f%% / %.1f%%", buyRatio, 100-buyRatio)
		}
	})

	t.Run("GetWeeklyMarginInterest_TimeSeries", func(t *testing.T) {
		// 過去数週間のトレンド分析
		endFriday := getRecentFriday()
		startDate, _ := time.Parse("2006-01-02", endFriday)
		startDate = startDate.AddDate(0, 0, -28) // 4週間前
		startFriday := startDate.Format("2006-01-02")

		interests, err := jq.WeeklyMarginInterest.GetWeeklyMarginInterestByCodeAndDateRange("7203", startFriday, endFriday)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get time series data: %v", err)
		}

		if len(interests) == 0 {
			t.Skip("No time series data available")
		}

		t.Logf("Time series analysis for code 7203 (%s to %s):", startFriday, endFriday)

		// 週次変化の分析
		for i, interest := range interests {
			if i < 5 { // 最新5週分
				marginBalance := interest.LongMarginTradeVolume - interest.ShortMarginTradeVolume

				t.Logf("Week %s:", interest.Date)
				t.Logf("  Net margin position: %.0f shares (Long: %.0f, Short: %.0f)",
					marginBalance, interest.LongMarginTradeVolume, interest.ShortMarginTradeVolume)
				t.Logf("  Standardized margin: Long %.0f, Short %.0f shares",
					interest.LongStandardizedMarginTradeVolume,
					interest.ShortStandardizedMarginTradeVolume)
			}
		}

		// トレンド計算
		if len(interests) >= 2 {
			latest := interests[0]
			previous := interests[1]

			buyChange := latest.LongMarginTradeVolume - previous.LongMarginTradeVolume
			sellChange := latest.ShortMarginTradeVolume - previous.ShortMarginTradeVolume

			t.Logf("Weekly changes (latest vs previous):")
			t.Logf("  Long margin: %+.0f shares", buyChange)
			t.Logf("  Short margin: %+.0f shares", sellChange)
		}
	})

	t.Run("GetWeeklyMarginInterest_Pagination", func(t *testing.T) {
		// ページネーションのテスト
		friday := getRecentFriday()

		params := jquants.WeeklyMarginInterestParams{
			Date: friday,
		}

		resp, err := jq.WeeklyMarginInterest.GetWeeklyMarginInterest(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No weekly margin interest data available")
		}

		if resp == nil || len(resp.WeeklyMarginInterest) == 0 {
			t.Skip("No data available for pagination test")
		}

		firstPageCount := len(resp.WeeklyMarginInterest)
		t.Logf("First page: %d records", firstPageCount)

		if resp.PaginationKey != "" {
			// 次のページを取得
			params.PaginationKey = resp.PaginationKey
			resp2, err := jq.WeeklyMarginInterest.GetWeeklyMarginInterest(params)
			if err != nil {
				t.Fatalf("Failed to get next page: %v", err)
			}

			if resp2 != nil && len(resp2.WeeklyMarginInterest) > 0 {
				t.Logf("Second page: %d records", len(resp2.WeeklyMarginInterest))
			}
		}
	})

	t.Run("GetWeeklyMarginInterest_MarketAnalysis", func(t *testing.T) {
		// 市場全体の信用取引分析
		friday := getRecentFriday()

		params := jquants.WeeklyMarginInterestParams{
			Date: friday,
		}

		resp, err := jq.WeeklyMarginInterest.GetWeeklyMarginInterest(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No data available for market analysis")
		}

		if resp == nil || len(resp.WeeklyMarginInterest) == 0 {
			t.Skip("No data available")
		}

		// 上位銘柄の分析（信用買い残高順）
		type MarginData struct {
			Code      string
			BuyValue  float64
			SellValue float64
			NetValue  float64
		}

		var marginRanking []MarginData

		for _, interest := range resp.WeeklyMarginInterest {
			if len(marginRanking) >= 20 {
				break // 上位20銘柄のみ
			}

			netVolume := interest.LongMarginTradeVolume - interest.ShortMarginTradeVolume

			marginRanking = append(marginRanking, MarginData{
				Code:      interest.Code,
				BuyValue:  interest.LongMarginTradeVolume,
				SellValue: interest.ShortMarginTradeVolume,
				NetValue:  netVolume,
			})
		}

		t.Logf("Top margin trading stocks on %s:", friday)
		for i, data := range marginRanking {
			if i >= 10 {
				break
			}
			t.Logf("  Code %s: Long=%.0f, Short=%.0f, Net=%+.0f shares",
				data.Code, data.BuyValue, data.SellValue, data.NetValue)
		}
	})

	t.Run("GetWeeklyMarginInterest_ErrorHandling", func(t *testing.T) {
		// エラーケースのテスト

		// 存在しない銘柄コード
		interests, err := jq.WeeklyMarginInterest.GetWeeklyMarginInterestByCodeAndDateRange("99999", "2024-01-05", "2024-01-05")
		if err == nil && len(interests) > 0 {
			t.Error("Expected error or empty result for invalid code")
		}

		// 平日の日付（金曜日以外）
		monday := time.Now()
		for monday.Weekday() != time.Monday {
			monday = monday.AddDate(0, 0, -1)
		}
		mondayStr := monday.Format("2006-01-02")

		params := jquants.WeeklyMarginInterestParams{
			Date: mondayStr,
		}

		resp, err := jq.WeeklyMarginInterest.GetWeeklyMarginInterest(params)
		if err == nil && resp != nil && len(resp.WeeklyMarginInterest) > 0 {
			t.Logf("Warning: Got data for non-Friday date %s", mondayStr)
		}
	})
}
