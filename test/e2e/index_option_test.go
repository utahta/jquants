//go:build e2e
// +build e2e

package e2e

import (
	"testing"

	"github.com/utahta/jquants"
)

// TestIndexOptionEndpoint は日経225オプションエンドポイントの完全なテスト
func TestIndexOptionEndpoint(t *testing.T) {
	t.Run("GetIndexOptions_ByDate", func(t *testing.T) {
		// 最近の営業日の日経225オプションを取得
		date := getTestDate()

		params := jquants.IndexOptionParams{
			Date: date,
		}

		resp, err := jq.IndexOption.GetIndexOptions(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get index options: %v", err)
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No index options data available")
		}

		// 各オプションデータを詳細に検証
		expectedDate := getTestDateFormatted() // YYYY-MM-DD形式
		for i, option := range resp.Data {
			// 基本情報の検証
			if option.Date == "" {
				t.Errorf("Option[%d]: Date is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(option.Date) != 10 || option.Date[4] != '-' || option.Date[7] != '-' {
					t.Errorf("Option[%d]: Date format invalid = %v, want YYYY-MM-DD", i, option.Date)
				}
				if option.Date != expectedDate {
					t.Errorf("Option[%d]: Date = %v, want %v", i, option.Date, expectedDate)
				}
			}

			if option.Code == "" {
				t.Errorf("Option[%d]: Code is empty", i)
			}

			// 限月の検証
			if option.CM == "" {
				t.Errorf("Option[%d]: ContractMonth is empty", i)
			} else {
				// 限月フォーマットの検証（YYYY-MM形式）
				if len(option.CM) != 7 || option.CM[4] != '-' {
					t.Errorf("Option[%d]: ContractMonth format invalid = %v, want YYYY-MM", i, option.CM)
				}
			}

			// 権利行使価格の妥当性チェック
			if option.Strike <= 0 {
				t.Errorf("Option[%d]: StrikePrice = %v, want > 0", i, option.Strike)
			}

			// プットコール区分の検証
			if option.PCDiv != "1" && option.PCDiv != "2" {
				t.Errorf("Option[%d]: PutCallDivision = %v, want 1 or 2", i, option.PCDiv)
			}

			// 緊急取引証拠金発動区分の検証
			if option.EmMrgnTrgDiv != "001" && option.EmMrgnTrgDiv != "002" {
				t.Errorf("Option[%d]: EmergencyMarginTriggerDivision = %v, want 001 or 002", i, option.EmMrgnTrgDiv)
			}

			// 日通し四本値の論理的整合性チェック
			if option.H < option.L {
				t.Errorf("Option[%d]: WholeDayHigh (%v) < WholeDayLow (%v)", i, option.H, option.L)
			}
			if option.O > 0 && option.H > 0 && option.O > option.H {
				t.Errorf("Option[%d]: WholeDayOpen (%v) > WholeDayHigh (%v)", i, option.O, option.H)
			}

			// 日中四本値の論理的整合性チェック
			if option.AH < option.AL {
				t.Errorf("Option[%d]: DaySessionHigh (%v) < DaySessionLow (%v)", i, option.AH, option.AL)
			}

			// 価格の妥当性チェック（負の値は通常ありえない）
			if option.O < 0 {
				t.Errorf("Option[%d]: WholeDayOpen = %v, want >= 0", i, option.O)
			}
			if option.H < 0 {
				t.Errorf("Option[%d]: WholeDayHigh = %v, want >= 0", i, option.H)
			}
			if option.L < 0 {
				t.Errorf("Option[%d]: WholeDayLow = %v, want >= 0", i, option.L)
			}
			if option.C < 0 {
				t.Errorf("Option[%d]: WholeDayClose = %v, want >= 0", i, option.C)
			}

			// 出来高・建玉の妥当性チェック
			if option.Vo < 0 {
				t.Errorf("Option[%d]: Volume = %v, want >= 0", i, option.Vo)
			}
			if option.OI < 0 {
				t.Errorf("Option[%d]: OpenInterest = %v, want >= 0", i, option.OI)
			}
			if option.Va < 0 {
				t.Errorf("Option[%d]: TurnoverValue = %v, want >= 0", i, option.Va)
			}

			// 最初の5件の詳細ログ
			if i < 5 {
				optionType := "Call"
				if option.PCDiv == "1" {
					optionType = "Put"
				}

				t.Logf("Option[%d]: Date=%s, Code=%s", i, option.Date, option.Code)
				t.Logf("  Type: %s, Strike: %.0f, Month: %s",
					optionType, option.Strike, option.CM)
				t.Logf("  WholeDay OHLC: Open=%.1f, High=%.1f, Low=%.1f, Close=%.1f",
					option.O, option.H, option.L, option.C)
				t.Logf("  DaySession OHLC: Open=%.1f, High=%.1f, Low=%.1f, Close=%.1f",
					option.AO, option.AH, option.AL, option.AC)
				t.Logf("  Volume: %.0f, OpenInterest: %.0f, Turnover: %.0f",
					option.Vo, option.OI, option.Va)

				// ナイトセッションデータの表示
				if option.HasNightSession() {
					t.Logf("  NightSession OHLC: Open=%.1f, High=%.1f, Low=%.1f, Close=%.1f",
						*option.EO, *option.EH,
						*option.EL, *option.EC)
				} else {
					t.Logf("  NightSession: No data")
				}

				// 価格・リスク情報（2016年7月19日以降）
				if option.Settle != nil {
					t.Logf("  Settlement: %.1f", *option.Settle)
				}
				if option.Theo != nil {
					t.Logf("  Theoretical: %.3f", *option.Theo)
				}
				if option.IV != nil {
					t.Logf("  IV: %.2f%%", *option.IV)
				}
				// 2016年7月19日以降の追加フィールド
				if option.LTD != "" {
					// LastTradingDayのフォーマット検証
					if len(option.LTD) != 10 || option.LTD[4] != '-' || option.LTD[7] != '-' {
						t.Errorf("Option[%d]: LastTradingDay format invalid = %v, want YYYY-MM-DD", i, option.LTD)
					}
					t.Logf("  Last Trading Day: %s", option.LTD)
				}
				if option.SQD != "" {
					// SpecialQuotationDayのフォーマット検証
					if len(option.SQD) != 10 || option.SQD[4] != '-' || option.SQD[7] != '-' {
						t.Errorf("Option[%d]: SpecialQuotationDay format invalid = %v, want YYYY-MM-DD", i, option.SQD)
					}
					t.Logf("  SQ Day: %s", option.SQD)
				}
				if option.VoOA != nil {
					// Volume(OnlyAuction)の妥当性チェック
					if *option.VoOA < 0 {
						t.Errorf("Option[%d]: VolumeOnlyAuction = %v, want >= 0", i, *option.VoOA)
					}
					t.Logf("  Volume (Only Auction): %.0f", *option.VoOA)
				}
				if option.BaseVol != nil {
					// BaseVolatilityの妥当性チェック
					if *option.BaseVol < 0 || *option.BaseVol > 200 {
						t.Errorf("Option[%d]: BaseVolatility = %v, want 0-200", i, *option.BaseVol)
					}
					t.Logf("  Base Volatility: %.2f%%", *option.BaseVol)
				}
				if option.UnderPx != nil {
					// UnderlyingPriceの妥当性チェック
					if *option.UnderPx <= 0 {
						t.Errorf("Option[%d]: UnderlyingPrice = %v, want > 0", i, *option.UnderPx)
					}
					t.Logf("  Underlying Price: %.2f", *option.UnderPx)
				}
				if option.IR != nil {
					// InterestRateの妥当性チェック
					if *option.IR < -10 || *option.IR > 10 {
						t.Errorf("Option[%d]: InterestRate = %v, want -10 to 10", i, *option.IR)
					}
					t.Logf("  Interest Rate: %.4f%%", *option.IR)
				}
			}
		}

		t.Logf("Retrieved %d index options records", len(resp.Data))
	})

	t.Run("GetCallOptions", func(t *testing.T) {
		// コールオプションのみを取得
		date := getTestDate()

		options, err := jq.IndexOption.GetCallOptions(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get call options: %v", err)
		}

		if len(options) == 0 {
			t.Skip("No call options data available")
		}

		// 全てコールオプションか確認
		for i, option := range options {
			if !option.IsCall() {
				t.Errorf("Option[%d]: Expected call option but got PutCallDivision = %v", i, option.PCDiv)
			}
		}

		t.Logf("Retrieved %d call options", len(options))
	})

	t.Run("GetPutOptions", func(t *testing.T) {
		// プットオプションのみを取得
		date := getTestDate()

		options, err := jq.IndexOption.GetPutOptions(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get put options: %v", err)
		}

		if len(options) == 0 {
			t.Skip("No put options data available")
		}

		// 全てプットオプションか確認
		for i, option := range options {
			if !option.IsPut() {
				t.Errorf("Option[%d]: Expected put option but got PutCallDivision = %v", i, option.PCDiv)
			}
		}

		t.Logf("Retrieved %d put options", len(options))
	})

	t.Run("GetOptionChain", func(t *testing.T) {
		// オプションチェーンの分析
		date := getTestDate()

		options, err := jq.IndexOption.GetOptionChain(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No index options data available for option chain analysis")
		}

		if len(options) == 0 {
			t.Skip("No index options data available")
		}

		// 権利行使価格別の分類
		strikeMap := make(map[float64]struct {
			calls []jquants.IndexOption
			puts  []jquants.IndexOption
		})

		for _, option := range options {
			entry := strikeMap[option.Strike]
			if option.IsCall() {
				entry.calls = append(entry.calls, option)
			} else if option.IsPut() {
				entry.puts = append(entry.puts, option)
			}
			strikeMap[option.Strike] = entry
		}

		t.Logf("Option chain analysis for %s:", date)
		t.Logf("  Unique strike prices: %d", len(strikeMap))

		// 代表的な権利行使価格でのコール・プット分析
		count := 0
		for strike, entry := range strikeMap {
			if count >= 5 { // 最初の5つの権利行使価格のみ
				break
			}

			callCount := len(entry.calls)
			putCount := len(entry.puts)

			t.Logf("  Strike %.0f: %d calls, %d puts", strike, callCount, putCount)

			// 建玉の合計
			if callCount > 0 && putCount > 0 {
				totalCallOI := float64(0)
				totalPutOI := float64(0)

				for _, call := range entry.calls {
					totalCallOI += call.OI
				}
				for _, put := range entry.puts {
					totalPutOI += put.OI
				}

				t.Logf("    Open Interest - Calls: %.0f, Puts: %.0f", totalCallOI, totalPutOI)
			}

			count++
		}
	})

	t.Run("GetIndexOptions_ContractMonthAnalysis", func(t *testing.T) {
		// 限月別の分析
		date := getTestDate()

		params := jquants.IndexOptionParams{
			Date: date,
		}

		resp, err := jq.IndexOption.GetIndexOptions(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No index options data available for contract month analysis")
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No index options data available")
		}

		// 限月別の集計
		monthlyData := make(map[string]struct {
			count       int
			totalVolume float64
			totalOI     float64
			calls       int
			puts        int
		})

		for _, option := range resp.Data {
			data := monthlyData[option.CM]
			data.count++
			data.totalVolume += option.Vo
			data.totalOI += option.OI

			if option.IsCall() {
				data.calls++
			} else if option.IsPut() {
				data.puts++
			}

			monthlyData[option.CM] = data
		}

		t.Logf("Contract month analysis:")
		for month, data := range monthlyData {
			t.Logf("  %s: %d options (%d calls, %d puts), Volume=%.0f, OI=%.0f",
				month, data.count, data.calls, data.puts, data.totalVolume, data.totalOI)
		}
	})

	t.Run("GetIndexOptions_EmergencyMarginAnalysis", func(t *testing.T) {
		// 緊急取引証拠金発動区分の分析
		date := getTestDate()

		params := jquants.IndexOptionParams{
			Date: date,
		}

		resp, err := jq.IndexOption.GetIndexOptions(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No index options data available for emergency margin analysis")
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No index options data available")
		}

		normalCount := 0
		emergencyCount := 0

		for _, option := range resp.Data {
			if option.IsEmergencyMarginTriggered() {
				emergencyCount++
			} else {
				normalCount++
			}
		}

		t.Logf("Emergency margin analysis:")
		t.Logf("  Normal (002): %d records", normalCount)
		t.Logf("  Emergency (001): %d records", emergencyCount)

		if emergencyCount > 0 {
			t.Logf("  WARNING: Emergency margin triggered on %s", date)
		}
	})

	t.Run("GetIndexOptions_Pagination", func(t *testing.T) {
		// ページネーションのテスト
		date := getTestDate()

		params := jquants.IndexOptionParams{
			Date: date,
		}

		resp, err := jq.IndexOption.GetIndexOptions(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No index options data available")
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No data available for pagination test")
		}

		firstPageCount := len(resp.Data)
		t.Logf("First page: %d records", firstPageCount)

		if resp.PaginationKey != "" {
			// 次のページを取得
			params.PaginationKey = resp.PaginationKey
			resp2, err := jq.IndexOption.GetIndexOptions(params)
			if err != nil {
				t.Fatalf("Failed to get next page: %v", err)
			}

			if resp2 != nil && len(resp2.Data) > 0 {
				t.Logf("Second page: %d records", len(resp2.Data))
			}
		}
	})

	t.Run("GetIndexOptions_ErrorHandling", func(t *testing.T) {
		// エラーケースのテスト

		// 未来の日付
		params := jquants.IndexOptionParams{
			Date: "2030-01-01",
		}

		resp, err := jq.IndexOption.GetIndexOptions(params)
		if err == nil && resp != nil && len(resp.Data) > 0 {
			t.Error("Expected error or empty result for future date")
		}

		// 無効な日付形式
		params = jquants.IndexOptionParams{
			Date: "invalid-date",
		}

		resp, err = jq.IndexOption.GetIndexOptions(params)
		if err == nil && resp != nil && len(resp.Data) > 0 {
			t.Error("Expected error or empty result for invalid date format")
		}

		// 空の日付（必須パラメータ）
		params = jquants.IndexOptionParams{
			Date: "",
		}

		_, err = jq.IndexOption.GetIndexOptions(params)
		if err == nil {
			t.Error("Expected error for missing required date parameter")
		}
	})
}
