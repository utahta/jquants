//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants"
)

// TestOptionsEndpoint は/derivatives/options（プレミアムプラン専用）エンドポイントの完全なテスト
func TestOptionsEndpoint(t *testing.T) {
	t.Run("GetOptions_ByDate", func(t *testing.T) {
		// 最近の営業日の個別株オプションを取得
		date := getTestDate()

		params := jquants.OptionsParams{
			Date: date,
		}

		resp, err := jq.Options.GetOptions(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get options: %v", err)
		}

		if resp == nil || len(resp.Options) == 0 {
			t.Skip("No options data available")
		}

		// 各オプションデータを詳細に検証
		for i, option := range resp.Options {
			// 基本情報の検証
			if option.Date == "" {
				t.Errorf("Option[%d]: Date is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(option.Date) != 10 || option.Date[4] != '-' || option.Date[7] != '-' {
					t.Errorf("Option[%d]: Date format invalid = %v, want YYYY-MM-DD", i, option.Date)
				}
				// 日付の一致確認
				expectedDate := getTestDateFormatted()
				if option.Date != expectedDate {
					t.Errorf("Option[%d]: Date = %v, want %v", i, option.Date, expectedDate)
				}
			}

			if option.Code == "" {
				t.Errorf("Option[%d]: Code is empty", i)
			}

			// オプション商品区分の検証
			if option.DerivativesProductCategory == "" {
				t.Errorf("Option[%d]: DerivativesProductCategory is empty", i)
			}

			// 原資産銘柄コードの検証（有価証券オプション以外は"-"）
			if option.UnderlyingSSO == "" {
				t.Errorf("Option[%d]: UnderlyingSSO is empty", i)
			}

			// 限月の検証
			if option.ContractMonth == "" {
				t.Errorf("Option[%d]: ContractMonth is empty", i)
			} else {
				// 限月フォーマットの検証（YYYY-MM形式またはYYYY-WW形式）
				if len(option.ContractMonth) != 7 || option.ContractMonth[4] != '-' {
					t.Errorf("Option[%d]: ContractMonth format invalid = %v, want YYYY-MM or YYYY-WW", i, option.ContractMonth)
				}
				// NK225MWE（日経225miniオプション）の場合はYYYY-WW形式
				if option.DerivativesProductCategory == "NK225MWE" {
					// 週番号の妥当性チェック（01-53）
					weekStr := option.ContractMonth[5:]
					week := 0
					_, err := fmt.Sscanf(weekStr, "%d", &week)
					if err != nil || week < 1 || week > 53 {
						t.Errorf("Option[%d]: Invalid week number in ContractMonth = %v", i, option.ContractMonth)
					}
				}
			}

			// 権利行使価格の妥当性チェック
			if option.StrikePrice <= 0 {
				t.Errorf("Option[%d]: StrikePrice = %v, want > 0", i, option.StrikePrice)
			}

			// プットコール区分の検証
			if option.PutCallDivision != "1" && option.PutCallDivision != "2" {
				t.Errorf("Option[%d]: PutCallDivision = %v, want 1 or 2", i, option.PutCallDivision)
			}

			// 日通し四本値の論理的整合性チェック
			if option.WholeDayHigh < option.WholeDayLow {
				t.Errorf("Option[%d]: WholeDayHigh (%v) < WholeDayLow (%v)", i, option.WholeDayHigh, option.WholeDayLow)
			}
			if option.WholeDayOpen > 0 && option.WholeDayHigh > 0 && option.WholeDayOpen > option.WholeDayHigh {
				t.Errorf("Option[%d]: WholeDayOpen (%v) > WholeDayHigh (%v)", i, option.WholeDayOpen, option.WholeDayHigh)
			}
			if option.WholeDayOpen > 0 && option.WholeDayLow > 0 && option.WholeDayOpen < option.WholeDayLow {
				t.Errorf("Option[%d]: WholeDayOpen (%v) < WholeDayLow (%v)", i, option.WholeDayOpen, option.WholeDayLow)
			}
			if option.WholeDayClose > 0 && option.WholeDayHigh > 0 && option.WholeDayClose > option.WholeDayHigh {
				t.Errorf("Option[%d]: WholeDayClose (%v) > WholeDayHigh (%v)", i, option.WholeDayClose, option.WholeDayHigh)
			}
			if option.WholeDayClose > 0 && option.WholeDayLow > 0 && option.WholeDayClose < option.WholeDayLow {
				t.Errorf("Option[%d]: WholeDayClose (%v) < WholeDayLow (%v)", i, option.WholeDayClose, option.WholeDayLow)
			}

			// 日中四本値の論理的整合性チェック
			if option.DaySessionHigh < option.DaySessionLow {
				t.Errorf("Option[%d]: DaySessionHigh (%v) < DaySessionLow (%v)", i, option.DaySessionHigh, option.DaySessionLow)
			}

			// 価格の妥当性チェック（負の値は通常ありえない）
			if option.WholeDayOpen < 0 {
				t.Errorf("Option[%d]: WholeDayOpen = %v, want >= 0", i, option.WholeDayOpen)
			}
			if option.WholeDayHigh < 0 {
				t.Errorf("Option[%d]: WholeDayHigh = %v, want >= 0", i, option.WholeDayHigh)
			}
			if option.WholeDayLow < 0 {
				t.Errorf("Option[%d]: WholeDayLow = %v, want >= 0", i, option.WholeDayLow)
			}
			if option.WholeDayClose < 0 {
				t.Errorf("Option[%d]: WholeDayClose = %v, want >= 0", i, option.WholeDayClose)
			}

			// 出来高・建玉の妥当性チェック
			if option.Volume < 0 {
				t.Errorf("Option[%d]: Volume = %v, want >= 0", i, option.Volume)
			}
			if option.OpenInterest < 0 {
				t.Errorf("Option[%d]: OpenInterest = %v, want >= 0", i, option.OpenInterest)
			}

			// 最初の5件の詳細ログ
			if i < 5 {
				optionType := "Call"
				if option.IsPut() {
					optionType = "Put"
				}

				t.Logf("Option[%d]: Date=%s, Code=%s", i, option.Date, option.Code)
				t.Logf("  Underlying: %s, Type: %s, Strike: %.0f, Month: %s",
					option.UnderlyingSSO, optionType, option.StrikePrice, option.ContractMonth)
				t.Logf("  WholeDay OHLC: Open=%.1f, High=%.1f, Low=%.1f, Close=%.1f",
					option.WholeDayOpen, option.WholeDayHigh, option.WholeDayLow, option.WholeDayClose)
				t.Logf("  DaySession OHLC: Open=%.1f, High=%.1f, Low=%.1f, Close=%.1f",
					option.DaySessionOpen, option.DaySessionHigh, option.DaySessionLow, option.DaySessionClose)
				t.Logf("  Volume: %.0f, OpenInterest: %.0f, Turnover: %.0f",
					option.Volume, option.OpenInterest, option.TurnoverValue)

				// ナイトセッションデータの表示
				if option.HasNightSession() {
					nightOpen := option.GetNightSessionOpen()
					nightHigh := option.GetNightSessionHigh()
					nightLow := option.GetNightSessionLow()
					nightClose := option.GetNightSessionClose()
					if nightOpen != nil && nightHigh != nil && nightLow != nil && nightClose != nil {
						t.Logf("  NightSession OHLC: Open=%.1f, High=%.1f, Low=%.1f, Close=%.1f",
							*nightOpen, *nightHigh, *nightLow, *nightClose)
					}
				} else {
					t.Logf("  NightSession: No data")
				}

				// 価格・リスク情報（2016年7月19日以降）
				if option.SettlementPrice != nil {
					t.Logf("  Settlement: %.1f", *option.SettlementPrice)
				}
				if option.TheoreticalPrice != nil {
					t.Logf("  Theoretical: %.3f", *option.TheoreticalPrice)
				}
				if option.ImpliedVolatility != nil {
					t.Logf("  IV: %.2f%%", *option.ImpliedVolatility)
				}
				// 2016年7月19日以降の追加フィールド
				if option.LastTradingDay != nil && *option.LastTradingDay != "" {
					// LastTradingDayのフォーマット検証
					if len(*option.LastTradingDay) != 10 || (*option.LastTradingDay)[4] != '-' || (*option.LastTradingDay)[7] != '-' {
						t.Errorf("Option[%d]: LastTradingDay format invalid = %v, want YYYY-MM-DD", i, *option.LastTradingDay)
					}
					t.Logf("  Last Trading Day: %s", *option.LastTradingDay)
				}
				if option.SpecialQuotationDay != nil && *option.SpecialQuotationDay != "" {
					// SpecialQuotationDayのフォーマット検証
					if len(*option.SpecialQuotationDay) != 10 || (*option.SpecialQuotationDay)[4] != '-' || (*option.SpecialQuotationDay)[7] != '-' {
						t.Errorf("Option[%d]: SpecialQuotationDay format invalid = %v, want YYYY-MM-DD", i, *option.SpecialQuotationDay)
					}
					t.Logf("  SQ Day: %s", *option.SpecialQuotationDay)
				}
				if option.VolumeOnlyAuction != nil {
					// Volume(OnlyAuction)の妥当性チェック
					if *option.VolumeOnlyAuction < 0 {
						t.Errorf("Option[%d]: VolumeOnlyAuction = %v, want >= 0", i, *option.VolumeOnlyAuction)
					}
					t.Logf("  Volume (Only Auction): %.0f", *option.VolumeOnlyAuction)
				}
				if option.CentralContractMonthFlag != nil {
					// 中心限月フラグの検証
					if *option.CentralContractMonthFlag != "0" && *option.CentralContractMonthFlag != "1" {
						t.Errorf("Option[%d]: CentralContractMonthFlag = %v, want 0 or 1", i, *option.CentralContractMonthFlag)
					}
					if *option.CentralContractMonthFlag == "1" {
						t.Logf("  Central Contract Month: Yes")
					}
				}
				if option.BaseVolatility != nil {
					t.Logf("  Base Volatility: %.2f%%", *option.BaseVolatility)
				}
				if option.UnderlyingPrice != nil {
					t.Logf("  Underlying Price: %.2f", *option.UnderlyingPrice)
				}
				if option.InterestRate != nil {
					t.Logf("  Interest Rate: %.4f%%", *option.InterestRate)
				}
			}
		}

		t.Logf("Retrieved %d options records", len(resp.Options))
	})

	t.Run("GetSecurityOptionsByCode", func(t *testing.T) {
		// トヨタ自動車の有価証券オプションを取得
		date := getTestDate()
		options, err := jq.Options.GetSecurityOptionsByCode(date, "7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get security options by code: %v", err)
		}

		if len(options) == 0 {
			t.Skip("No security options data available for code 7203")
		}

		// 全て指定銘柄のオプションか確認（有価証券オプションでは銘柄コードが設定される）
		for i, option := range options {
			if !option.IsSecurityOption() {
				t.Errorf("Option[%d]: Expected security option but UnderlyingSSO = %v", i, option.UnderlyingSSO)
			}
			// 有価証券オプションの場合、UnderlyingSSOに銘柄コードが設定される
			if option.UnderlyingSSO != "7203" && option.UnderlyingSSO != "72030" {
				t.Errorf("Option[%d]: UnderlyingSSO = %v, want 7203 or 72030",
					i, option.UnderlyingSSO)
			}
		}

		t.Logf("Retrieved %d security options for code 7203", len(options))
	})

	t.Run("GetCallOptionsByCode", func(t *testing.T) {
		// トヨタ自動車のコールオプションのみを取得（フィルタリングで実現）
		date := getTestDate()
		allOptions, err := jq.Options.GetSecurityOptionsByCode(date, "7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get security options: %v", err)
		}

		if len(allOptions) == 0 {
			t.Skip("No security options data available for code 7203")
		}

		// コールオプションのみをフィルタリング
		var callOptions []jquants.Option
		for _, option := range allOptions {
			if option.IsCall() {
				callOptions = append(callOptions, option)
			}
		}

		if len(callOptions) == 0 {
			t.Skip("No call options data available for code 7203")
		}

		// 全てコールオプションか確認
		for i, option := range callOptions {
			if !option.IsCall() {
				t.Errorf("Option[%d]: Expected call option but got PutCallDivision = %v", i, option.PutCallDivision)
			}
			if option.UnderlyingSSO != "7203" && option.UnderlyingSSO != "72030" {
				t.Errorf("Option[%d]: UnderlyingSSO = %v, want 7203 or 72030",
					i, option.UnderlyingSSO)
			}
		}

		t.Logf("Retrieved %d call options for code 7203", len(callOptions))
	})

	t.Run("GetPutOptionsByCode", func(t *testing.T) {
		// トヨタ自動車のプットオプションのみを取得（フィルタリングで実現）
		date := getTestDate()
		allOptions, err := jq.Options.GetSecurityOptionsByCode(date, "7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get security options: %v", err)
		}

		if len(allOptions) == 0 {
			t.Skip("No security options data available for code 7203")
		}

		// プットオプションのみをフィルタリング
		var putOptions []jquants.Option
		for _, option := range allOptions {
			if option.IsPut() {
				putOptions = append(putOptions, option)
			}
		}

		if len(putOptions) == 0 {
			t.Skip("No put options data available for code 7203")
		}

		// 全てプットオプションか確認
		for i, option := range putOptions {
			if !option.IsPut() {
				t.Errorf("Option[%d]: Expected put option but got PutCallDivision = %v", i, option.PutCallDivision)
			}
			if option.UnderlyingSSO != "7203" && option.UnderlyingSSO != "72030" {
				t.Errorf("Option[%d]: UnderlyingSSO = %v, want 7203 or 72030",
					i, option.UnderlyingSSO)
			}
		}

		t.Logf("Retrieved %d put options for code 7203", len(putOptions))
	})

	t.Run("GetOptions_OptionChainAnalysis", func(t *testing.T) {
		// オプションチェーンの分析（有価証券オプション）
		date := getTestDate()
		options, err := jq.Options.GetSecurityOptionsByCode(date, "7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No options data available for option chain analysis")
		}

		if len(options) == 0 {
			t.Skip("No options data available")
		}

		// 権利行使価格別の分類
		strikeMap := make(map[float64]struct {
			calls []jquants.Option
			puts  []jquants.Option
		})

		for _, option := range options {
			entry := strikeMap[option.StrikePrice]
			if option.IsCall() {
				entry.calls = append(entry.calls, option)
			} else if option.IsPut() {
				entry.puts = append(entry.puts, option)
			}
			strikeMap[option.StrikePrice] = entry
		}

		t.Logf("Option chain analysis for code 7203:")
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
				totalCallVolume := float64(0)
				totalPutVolume := float64(0)

				for _, call := range entry.calls {
					totalCallOI += call.OpenInterest
					totalCallVolume += call.Volume
				}
				for _, put := range entry.puts {
					totalPutOI += put.OpenInterest
					totalPutVolume += put.Volume
				}

				t.Logf("    OI - Calls: %.0f, Puts: %.0f", totalCallOI, totalPutOI)
				t.Logf("    Volume - Calls: %.0f, Puts: %.0f", totalCallVolume, totalPutVolume)
			}

			count++
		}
	})

	t.Run("GetOptions_UnderlyingAssetAnalysis", func(t *testing.T) {
		// 原資産別の分析
		date := getTestDate()

		params := jquants.OptionsParams{
			Date: date,
		}

		resp, err := jq.Options.GetOptions(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No options data available for underlying asset analysis")
		}

		if resp == nil || len(resp.Options) == 0 {
			t.Skip("No options data available")
		}

		// 原資産別の集計
		underlyingData := make(map[string]struct {
			count       int
			totalVolume float64
			totalOI     float64
			callCount   int
			putCount    int
		})

		for _, option := range resp.Options {
			underlying := option.UnderlyingSSO
			if option.IsSecurityOption() {
				// 有価証券オプションの場合は銘柄コードを使用
				underlying = option.UnderlyingSSO
			} else {
				// 指数オプションの場合は商品区分を使用
				underlying = option.DerivativesProductCategory
			}

			data := underlyingData[underlying]
			data.count++
			data.totalVolume += option.Volume
			data.totalOI += option.OpenInterest

			if option.IsCall() {
				data.callCount++
			} else if option.IsPut() {
				data.putCount++
			}

			underlyingData[underlying] = data
		}

		t.Logf("Underlying asset analysis:")
		count := 0
		for underlying, data := range underlyingData {
			if count >= 10 { // 上位10銘柄のみ
				break
			}

			t.Logf("  %s: %d options (%d calls, %d puts)",
				underlying, data.count, data.callCount, data.putCount)
			t.Logf("    Volume: %.0f, OI: %.0f", data.totalVolume, data.totalOI)

			count++
		}
	})

	t.Run("GetOptions_Pagination", func(t *testing.T) {
		// ページネーションのテスト
		date := getTestDate()

		params := jquants.OptionsParams{
			Date: date,
		}

		resp, err := jq.Options.GetOptions(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No options data available")
		}

		if resp == nil || len(resp.Options) == 0 {
			t.Skip("No data available for pagination test")
		}

		firstPageCount := len(resp.Options)
		t.Logf("First page: %d records", firstPageCount)

		if resp.PaginationKey != "" {
			// 次のページを取得
			params.PaginationKey = resp.PaginationKey
			resp2, err := jq.Options.GetOptions(params)
			if err != nil {
				t.Fatalf("Failed to get next page: %v", err)
			}

			if resp2 != nil && len(resp2.Options) > 0 {
				t.Logf("Second page: %d records", len(resp2.Options))
			}
		}
	})

	t.Run("GetOptions_ErrorHandling", func(t *testing.T) {
		// エラーケースのテスト

		// 存在しない銘柄コード
		date := getTestDate()
		options, err := jq.Options.GetSecurityOptionsByCode(date, "99999")
		if err == nil && len(options) > 0 {
			t.Error("Expected error or empty result for invalid code")
		}

		// 未来の日付
		params := jquants.OptionsParams{
			Date: "2030-01-01",
		}

		resp, err := jq.Options.GetOptions(params)
		if err == nil && resp != nil && len(resp.Options) > 0 {
			t.Error("Expected error or empty result for future date")
		}

		// 無効な日付形式
		params = jquants.OptionsParams{
			Date: "invalid-date",
		}

		resp2, err := jq.Options.GetOptions(params)
		if err == nil && resp2 != nil && len(resp2.Options) > 0 {
			t.Error("Expected error or empty result for invalid date format")
		}
	})
}
