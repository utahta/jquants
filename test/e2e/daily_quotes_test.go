//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"testing"
	"time"

	"github.com/utahta/jquants"
)

// TestDailyQuotesEndpoint は/prices/daily_quotesエンドポイントの完全なテスト
func TestDailyQuotesEndpoint(t *testing.T) {
	t.Run("GetDailyQuotes_ByCodeAndDateRange", func(t *testing.T) {
		// 過去10日間のトヨタ自動車の株価を取得
		to := getTestDate()
		fromTime, _ := time.Parse("20060102", to)
		fromTime = fromTime.AddDate(0, 0, -10)
		from := fromTime.Format("20060102")

		params := jquants.DailyQuotesParams{
			Code: "7203",
			From: from,
			To:   to,
		}

		resp, err := jq.Quotes.GetDailyQuotes(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get daily quotes: %v", err)
		}

		if resp == nil || len(resp.DailyQuotes) == 0 {
			t.Skip("No quotes data available for the specified date range")
		}

		// 各日次データを詳細に検証
		for i, quote := range resp.DailyQuotes {
			// 基本情報の検証
			if quote.Code != "72030" && quote.Code != "7203" {
				t.Errorf("Quote[%d]: Code = %v, want 72030 or 7203", i, quote.Code)
			}
			if quote.Date == "" {
				t.Errorf("Quote[%d]: Date is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(quote.Date) != 10 || quote.Date[4] != '-' || quote.Date[7] != '-' {
					t.Errorf("Quote[%d]: Date format invalid = %v, want YYYY-MM-DD", i, quote.Date)
				}
			}

			// 四本値の検証
			if quote.Open == nil {
				t.Errorf("Quote[%d]: Open is nil", i)
			} else if *quote.Open <= 0 {
				t.Errorf("Quote[%d]: Open = %v, want > 0", i, *quote.Open)
			}

			if quote.High == nil {
				t.Errorf("Quote[%d]: High is nil", i)
			} else if *quote.High <= 0 {
				t.Errorf("Quote[%d]: High = %v, want > 0", i, *quote.High)
			}

			if quote.Low == nil {
				t.Errorf("Quote[%d]: Low is nil", i)
			} else if *quote.Low <= 0 {
				t.Errorf("Quote[%d]: Low = %v, want > 0", i, *quote.Low)
			}

			if quote.Close == nil {
				t.Errorf("Quote[%d]: Close is nil", i)
			} else if *quote.Close <= 0 {
				t.Errorf("Quote[%d]: Close = %v, want > 0", i, *quote.Close)
			}

			// 四本値の論理的整合性チェック
			if quote.High != nil && quote.Low != nil && *quote.High < *quote.Low {
				t.Errorf("Quote[%d]: High (%v) < Low (%v)", i, *quote.High, *quote.Low)
			}
			if quote.Open != nil && quote.High != nil && *quote.Open > *quote.High {
				t.Errorf("Quote[%d]: Open (%v) > High (%v)", i, *quote.Open, *quote.High)
			}
			if quote.Open != nil && quote.Low != nil && *quote.Open < *quote.Low {
				t.Errorf("Quote[%d]: Open (%v) < Low (%v)", i, *quote.Open, *quote.Low)
			}
			if quote.Close != nil && quote.High != nil && *quote.Close > *quote.High {
				t.Errorf("Quote[%d]: Close (%v) > High (%v)", i, *quote.Close, *quote.High)
			}
			if quote.Close != nil && quote.Low != nil && *quote.Close < *quote.Low {
				t.Errorf("Quote[%d]: Close (%v) < Low (%v)", i, *quote.Close, *quote.Low)
			}

			// 出来高の検証
			if quote.Volume == nil {
				t.Errorf("Quote[%d]: Volume is nil", i)
			} else if *quote.Volume < 0 {
				t.Errorf("Quote[%d]: Volume = %v, want >= 0", i, *quote.Volume)
			}

			// 売買代金の検証
			if quote.TurnoverValue == nil {
				t.Errorf("Quote[%d]: TurnoverValue is nil", i)
			} else if *quote.TurnoverValue < 0 {
				t.Errorf("Quote[%d]: TurnoverValue = %v, want >= 0", i, *quote.TurnoverValue)
			}

			// ストップ高・ストップ安フラグの検証
			if quote.UpperLimit != "0" && quote.UpperLimit != "1" {
				t.Errorf("Quote[%d]: UpperLimit = %v, want 0 or 1", i, quote.UpperLimit)
			}
			if quote.LowerLimit != "0" && quote.LowerLimit != "1" {
				t.Errorf("Quote[%d]: LowerLimit = %v, want 0 or 1", i, quote.LowerLimit)
			}

			// 調整係数の検証
			if quote.AdjustmentFactor <= 0 {
				t.Errorf("Quote[%d]: AdjustmentFactor = %v, want > 0", i, quote.AdjustmentFactor)
			}

			// 調整後四本値の検証
			if quote.AdjustmentOpen == nil {
				t.Errorf("Quote[%d]: AdjustmentOpen is nil", i)
			}
			if quote.AdjustmentHigh == nil {
				t.Errorf("Quote[%d]: AdjustmentHigh is nil", i)
			}
			if quote.AdjustmentLow == nil {
				t.Errorf("Quote[%d]: AdjustmentLow is nil", i)
			}
			if quote.AdjustmentClose == nil {
				t.Errorf("Quote[%d]: AdjustmentClose is nil", i)
			}
			if quote.AdjustmentVolume == nil {
				t.Errorf("Quote[%d]: AdjustmentVolume is nil", i)
			}

			if i == 0 {
				// 最初のデータの詳細ログ
				t.Logf("Latest quote: Date=%s, O=%v, H=%v, L=%v, C=%v, V=%v",
					quote.Date,
					safeDeref(quote.Open),
					safeDeref(quote.High),
					safeDeref(quote.Low),
					safeDeref(quote.Close),
					safeDeref(quote.Volume))
			}
		}

		t.Logf("Retrieved %d quotes for date range %s to %s", len(resp.DailyQuotes), from, to)

		// データの並び順を確認
		if len(resp.DailyQuotes) > 1 {
			isAscending := resp.DailyQuotes[0].Date < resp.DailyQuotes[1].Date
			t.Logf("Data order: %s", ternary(isAscending, "ascending", "descending"))
		}
	})

	t.Run("GetDailyQuotes_ByDate", func(t *testing.T) {
		// 特定日の全銘柄データを取得
		date := getTestDate()

		params := jquants.DailyQuotesParams{
			Date: date,
		}

		resp, err := jq.Quotes.GetDailyQuotes(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Logf("No quotes data for date %s: %v", date, err)
			return
		}

		if resp == nil || len(resp.DailyQuotes) == 0 {
			t.Skip("No quotes data available for the specified date")
		}

		t.Logf("Found %d quotes for date %s", len(resp.DailyQuotes), date)

		// 最初の10件を検証
		for i, quote := range resp.DailyQuotes {
			if i >= 10 {
				break
			}

			// 日付の検証（APIはYYYY-MM-DD形式で返す）
			expectedDate := getTestDateFormatted()
			if quote.Date != expectedDate {
				t.Errorf("Quote[%d]: Date = %v, want %v", i, quote.Date, expectedDate)
			}
			if quote.Code == "" {
				t.Errorf("Quote[%d]: Code is empty", i)
			}
		}

		// ページネーションキーの確認
		if resp.PaginationKey != "" {
			t.Logf("Pagination key present: %s", resp.PaginationKey)
		}
	})

	t.Run("GetDailyQuotes_Pagination", func(t *testing.T) {
		// ページネーションのテスト
		date := getTestDate()

		params := jquants.DailyQuotesParams{
			Date: date,
		}

		resp, err := jq.Quotes.GetDailyQuotes(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No quotes data available")
		}

		if resp == nil || len(resp.DailyQuotes) == 0 {
			t.Skip("No quotes data available")
		}

		firstPageCount := len(resp.DailyQuotes)
		t.Logf("First page: %d quotes", firstPageCount)

		if resp.PaginationKey != "" {
			// 次のページを取得
			params.PaginationKey = resp.PaginationKey
			resp2, err := jq.Quotes.GetDailyQuotes(params)
			if err != nil {
				t.Fatalf("Failed to get next page: %v", err)
			}

			if resp2 != nil && len(resp2.DailyQuotes) > 0 {
				t.Logf("Second page: %d quotes", len(resp2.DailyQuotes))

				// 異なるデータであることを確認
				if resp2.DailyQuotes[0].Code == resp.DailyQuotes[0].Code {
					t.Error("Second page might contain duplicate data")
				}
			}
		}
	})

	t.Run("GetDailyQuotesByCode_Convenience", func(t *testing.T) {
		// 便利メソッドのテスト（過去30日分）
		quotes, err := jq.Quotes.GetDailyQuotesByCode("7203", 30)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get daily quotes by code: %v", err)
		}

		if len(quotes) == 0 {
			t.Skip("No quotes data available")
		}

		t.Logf("Retrieved %d quotes for code 7203", len(quotes))

		// 全データが同じ銘柄コードか確認
		for i, quote := range quotes {
			if quote.Code != "72030" && quote.Code != "7203" {
				t.Errorf("Quote[%d]: Code = %v, want 72030 or 7203", i, quote.Code)
			}
		}

		// 日付の連続性をチェック（営業日ベース）
		if len(quotes) > 1 {
			prevDate := quotes[0].Date
			gaps := 0
			for i := 1; i < len(quotes); i++ {
				currDate := quotes[i].Date
				// 簡易的なギャップチェック（実際には祝日も考慮が必要）
				prevTime, _ := time.Parse("2006-01-02", prevDate)
				currTime, _ := time.Parse("2006-01-02", currDate)
				daysDiff := int(prevTime.Sub(currTime).Hours() / 24)
				if daysDiff > 3 { // 週末を考慮して3日以上のギャップ
					gaps++
				}
				prevDate = currDate
			}
			if gaps > 0 {
				t.Logf("Found %d date gaps (might include holidays)", gaps)
			}
		}
	})

	t.Run("GetDailyQuotes_ErrorCases", func(t *testing.T) {
		// 無効な銘柄コード
		params := jquants.DailyQuotesParams{
			Code: "99999", // 存在しない銘柄コード
			Date: getTestDate(),
		}

		resp, err := jq.Quotes.GetDailyQuotes(params)
		if err == nil && len(resp.DailyQuotes) > 0 {
			t.Error("Expected no data for invalid code")
		}

		// 未来の日付
		futureDate := time.Now().AddDate(0, 0, 7).Format("20060102")
		params = jquants.DailyQuotesParams{
			Code: "7203",
			Date: futureDate,
		}

		resp, err = jq.Quotes.GetDailyQuotes(params)
		if err == nil && len(resp.DailyQuotes) > 0 {
			t.Error("Expected no data for future date")
		}
	})
}

// safeDeref は*float64を安全に文字列に変換するヘルパー関数
func safeDeref(p *float64) string {
	if p == nil {
		return "nil"
	}
	return fmt.Sprintf("%.2f", *p)
}

// ternary は三項演算子の代替
func ternary(cond bool, trueVal, falseVal string) string {
	if cond {
		return trueVal
	}
	return falseVal
}
