//go:build e2e
// +build e2e

package e2e

import (
	"testing"

	"github.com/utahta/jquants"
)

// TestDailyMarginInterestEndpoint は/markets/margin-alertエンドポイントの完全なテスト
func TestDailyMarginInterestEndpoint(t *testing.T) {
	t.Run("GetDailyMarginInterest_ByCode", func(t *testing.T) {
		// 日々公表銘柄として指定されやすい銘柄の残高を取得
		// 注：日々公表銘柄に指定されていない銘柄はデータがない可能性がある
		params := jquants.DailyMarginInterestParams{
			Code: "13260", // 日々公表銘柄として一般的な銘柄
		}

		resp, err := jq.DailyMarginInterest.GetDailyMarginInterest(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get daily margin interest: %v", err)
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No daily margin interest data available for the specified code")
		}

		// 各日々公表データを詳細に検証
		for i, interest := range resp.Data {
			// 基本情報の検証
			if interest.Code == "" {
				t.Errorf("Interest[%d]: Code is empty", i)
			}
			if interest.PubDate == "" {
				t.Errorf("Interest[%d]: PubDate is empty", i)
			}

			// 日付フォーマットの検証（YYYY-MM-DD形式）
			if len(interest.PubDate) != 10 || interest.PubDate[4] != '-' || interest.PubDate[7] != '-' {
				t.Errorf("Interest[%d]: PubDate format invalid = %v, want YYYY-MM-DD", i, interest.PubDate)
			}

			// 申込日の検証
			if interest.AppDate == "" {
				t.Errorf("Interest[%d]: AppDate is empty", i)
			}

			// 売合計信用残高の検証
			if interest.ShrtOut < 0 {
				t.Errorf("Interest[%d]: ShrtOut = %v, want >= 0", i, interest.ShrtOut)
			}

			// 買合計信用残高の検証
			if interest.LongOut < 0 {
				t.Errorf("Interest[%d]: LongOut = %v, want >= 0", i, interest.LongOut)
			}

			// 取組比率の検証（売残/買残×100）
			if interest.SLRatio < 0 {
				t.Errorf("Interest[%d]: SLRatio = %v, want >= 0", i, interest.SLRatio)
			}

			// 一般信用残高の検証
			if interest.ShrtNegOut < 0 {
				t.Errorf("Interest[%d]: ShrtNegOut = %v, want >= 0", i, interest.ShrtNegOut)
			}
			if interest.LongNegOut < 0 {
				t.Errorf("Interest[%d]: LongNegOut = %v, want >= 0", i, interest.LongNegOut)
			}

			// 制度信用残高の検証
			if interest.ShrtStdOut < 0 {
				t.Errorf("Interest[%d]: ShrtStdOut = %v, want >= 0", i, interest.ShrtStdOut)
			}
			if interest.LongStdOut < 0 {
				t.Errorf("Interest[%d]: LongStdOut = %v, want >= 0", i, interest.LongStdOut)
			}

			// 公表理由のヘルパーメソッド検証
			_ = interest.PubReason.IsRestricted()
			_ = interest.PubReason.IsDailyPublication()
			_ = interest.PubReason.IsMonitoring()
			_ = interest.PubReason.IsRestrictedByJSF()
			_ = interest.PubReason.IsPrecautionByJSF()
			_ = interest.PubReason.IsUnclearOrSecOnAlert()

			// 最初の3件の詳細ログ
			if i < 3 {
				t.Logf("Interest[%d]: PubDate=%s, Code=%s, AppDate=%s",
					i, interest.PubDate, interest.Code, interest.AppDate)
				t.Logf("  Long: %.0f (Std: %.0f, Neg: %.0f)",
					interest.LongOut, interest.LongStdOut, interest.LongNegOut)
				t.Logf("  Short: %.0f (Std: %.0f, Neg: %.0f)",
					interest.ShrtOut, interest.ShrtStdOut, interest.ShrtNegOut)
				t.Logf("  SL Ratio: %.2f%%", interest.SLRatio)
				t.Logf("  TSE Margin Reg: %s", interest.TSEMrgnRegCls)
			}
		}

		t.Logf("Retrieved %d daily margin interest records", len(resp.Data))
	})

	t.Run("GetDailyMarginInterest_ByDate", func(t *testing.T) {
		// 特定日の全銘柄データを取得
		date := getTestDate()

		params := jquants.DailyMarginInterestParams{
			Date: date,
		}

		resp, err := jq.DailyMarginInterest.GetDailyMarginInterest(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Logf("Failed to get daily margin interest by date: %v", err)
			return
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No daily margin interest data for the specified date")
		}

		t.Logf("Retrieved %d daily margin interest records for %s", len(resp.Data), date)

		// 最初の10件を検証
		for i, interest := range resp.Data {
			if i >= 10 {
				break
			}

			if interest.Code == "" {
				t.Errorf("Interest[%d]: Code is empty", i)
			}

			// 各公表理由フラグを確認
			reasons := []string{}
			if interest.PubReason.IsRestricted() {
				reasons = append(reasons, "Restricted")
			}
			if interest.PubReason.IsDailyPublication() {
				reasons = append(reasons, "DailyPublication")
			}
			if interest.PubReason.IsMonitoring() {
				reasons = append(reasons, "Monitoring")
			}
			if interest.PubReason.IsRestrictedByJSF() {
				reasons = append(reasons, "RestrictedByJSF")
			}
			if interest.PubReason.IsPrecautionByJSF() {
				reasons = append(reasons, "PrecautionByJSF")
			}

			t.Logf("Interest[%d]: Code=%s, Reasons=%v", i, interest.Code, reasons)
		}

		// ページネーションキーの確認
		if resp.PaginationKey != "" {
			t.Logf("Pagination key present: %s", resp.PaginationKey)
		}
	})

	t.Run("GetDailyMarginInterestByCode_Convenience", func(t *testing.T) {
		// 便利メソッドのテスト
		interests, err := jq.DailyMarginInterest.GetDailyMarginInterestByCode("13260")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get daily margin interest by code: %v", err)
		}

		if len(interests) == 0 {
			t.Skip("No daily margin interest data available")
		}

		t.Logf("Retrieved %d daily margin interest records for code", len(interests))

		// 全データが同じ銘柄コードか確認（5桁コードの場合もある）
		firstCode := interests[0].Code
		for i, interest := range interests {
			if interest.Code != firstCode {
				t.Errorf("Interest[%d]: Code = %v, want %v", i, interest.Code, firstCode)
			}
		}
	})

	t.Run("GetDailyMarginInterestByDate_Convenience", func(t *testing.T) {
		// 便利メソッドのテスト（日付指定）
		date := getTestDate()

		interests, err := jq.DailyMarginInterest.GetDailyMarginInterestByDate(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get daily margin interest by date: %v", err)
		}

		if len(interests) == 0 {
			t.Skip("No daily margin interest data available for the date")
		}

		t.Logf("Retrieved %d daily margin interest records for %s", len(interests), date)

		// 統計情報
		totalLong := 0.0
		totalShort := 0.0
		restrictedCount := 0
		dailyPubCount := 0

		for _, interest := range interests {
			totalLong += interest.LongOut
			totalShort += interest.ShrtOut
			if interest.PubReason.IsRestricted() {
				restrictedCount++
			}
			if interest.PubReason.IsDailyPublication() {
				dailyPubCount++
			}
		}

		t.Logf("Market summary for %s:", date)
		t.Logf("  Total long margin: %.0f shares", totalLong)
		t.Logf("  Total short margin: %.0f shares", totalShort)
		t.Logf("  Restricted stocks: %d", restrictedCount)
		t.Logf("  Daily publication stocks: %d", dailyPubCount)
	})

	t.Run("GetDailyMarginInterestByCodeAndDateRange", func(t *testing.T) {
		// 期間指定のテスト
		to := getTestDate()
		from := "20250601" // 約2週間分

		interests, err := jq.DailyMarginInterest.GetDailyMarginInterestByCodeAndDateRange("13260", from, to)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get daily margin interest by code and date range: %v", err)
		}

		if len(interests) == 0 {
			t.Skip("No daily margin interest data available for the date range")
		}

		t.Logf("Retrieved %d records from %s to %s", len(interests), from, to)

		// トレンド分析
		if len(interests) >= 2 {
			latest := interests[0]
			oldest := interests[len(interests)-1]

			longChange := latest.LongOut - oldest.LongOut
			shortChange := latest.ShrtOut - oldest.ShrtOut

			t.Logf("Trend analysis (%s to %s):", oldest.PubDate, latest.PubDate)
			t.Logf("  Long margin change: %+.0f shares", longChange)
			t.Logf("  Short margin change: %+.0f shares", shortChange)
		}
	})

	t.Run("GetDailyMarginInterest_HelperMethods", func(t *testing.T) {
		// ヘルパーメソッドのテスト
		params := jquants.DailyMarginInterestParams{
			Date: getTestDate(),
		}

		resp, err := jq.DailyMarginInterest.GetDailyMarginInterest(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No data available")
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No data available")
		}

		// ヘルパーメソッドの動作確認
		interest := resp.Data[0]

		// 前日比の取得（「-」の場合は0を返す）
		shortOutChg := interest.GetShortOutChgValue()
		longOutChg := interest.GetLongOutChgValue()
		shortOutRatio := interest.GetShortOutRatioValue()
		longOutRatio := interest.GetLongOutRatioValue()

		t.Logf("Helper method results for %s:", interest.Code)
		t.Logf("  ShortOutChg: %.0f", shortOutChg)
		t.Logf("  LongOutChg: %.0f", longOutChg)
		t.Logf("  ShortOutRatio: %.2f%%", shortOutRatio)
		t.Logf("  LongOutRatio: %.2f%%", longOutRatio)
	})

	t.Run("GetDailyMarginInterest_TSERegulationConstants", func(t *testing.T) {
		// 東証信用貸借規制区分定数のテスト
		params := jquants.DailyMarginInterestParams{
			Date: getTestDate(),
		}

		resp, err := jq.DailyMarginInterest.GetDailyMarginInterest(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No data available")
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No data available")
		}

		// 規制区分の集計
		regCounts := make(map[string]int)
		for _, interest := range resp.Data {
			regCounts[interest.TSEMrgnRegCls]++
		}

		t.Logf("TSE Margin Regulation distribution:")
		for reg, count := range regCounts {
			regName := getTSERegulationName(reg)
			t.Logf("  %s (%s): %d stocks", reg, regName, count)
		}
	})

	t.Run("GetDailyMarginInterest_ErrorCases", func(t *testing.T) {
		// パラメータなしはエラー
		params := jquants.DailyMarginInterestParams{}

		_, err := jq.DailyMarginInterest.GetDailyMarginInterest(params)
		if err == nil {
			t.Error("Expected error for empty parameters")
		}

		// 存在しない銘柄コード（エラーまたは空結果）
		params = jquants.DailyMarginInterestParams{
			Code: "99999",
		}

		resp, err := jq.DailyMarginInterest.GetDailyMarginInterest(params)
		if err == nil && resp != nil && len(resp.Data) > 0 {
			t.Error("Expected empty result for invalid code")
		}
	})
}

// getTSERegulationName は東証信用貸借規制区分コードから名称を取得する
func getTSERegulationName(code string) string {
	names := map[string]string{
		jquants.TSEMarginRegulationNone:              "規制なし",
		jquants.TSEMarginRegulationCautionForNew:     "新規建て注意喚起",
		jquants.TSEMarginRegulationCautionForSelling: "売り注意喚起",
		jquants.TSEMarginRegulationCautionForBuying:  "買い注意喚起",
		jquants.TSEMarginRegulationRestrictedNew:     "新規建て規制",
		jquants.TSEMarginRegulationRestrictedSelling: "売り規制",
		jquants.TSEMarginRegulationRestrictedBuying:  "買い規制",
	}
	if name, ok := names[code]; ok {
		return name
	}
	return "不明"
}
