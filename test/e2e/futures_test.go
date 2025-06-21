//go:build e2e
// +build e2e

package e2e

import (
	"testing"

	"github.com/utahta/jquants"
)

// TestFuturesEndpoint は/derivatives/futures（プレミアムプラン専用）エンドポイントの完全なテスト
func TestFuturesEndpoint(t *testing.T) {
	t.Run("GetFutures_ByDate", func(t *testing.T) {
		// 最近の営業日の先物四本値を取得
		date := getTestDate()

		futures, err := jq.Futures.GetFuturesByDate(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get futures: %v", err)
		}

		if len(futures) == 0 {
			t.Skip("No futures data available")
		}

		// 各先物データを詳細に検証
		for i, future := range futures {
			// 基本情報の検証
			if future.Date == "" {
				t.Errorf("Future[%d]: Date is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(future.Date) != 10 || future.Date[4] != '-' || future.Date[7] != '-' {
					t.Errorf("Future[%d]: Date format invalid = %v, want YYYY-MM-DD", i, future.Date)
				}
				// 日付の一致確認
				expectedDate := getTestDateFormatted()
				if future.Date != expectedDate {
					t.Errorf("Future[%d]: Date = %v, want %v", i, future.Date, expectedDate)
				}
			}

			if future.Code == "" {
				t.Errorf("Future[%d]: Code is empty", i)
			}

			// 限月の検証
			if future.ContractMonth == "" {
				t.Errorf("Future[%d]: ContractMonth is empty", i)
			} else {
				// 限月フォーマットの検証（YYYY-MM形式）
				if len(future.ContractMonth) != 7 || future.ContractMonth[4] != '-' {
					t.Errorf("Future[%d]: ContractMonth format invalid = %v, want YYYY-MM", i, future.ContractMonth)
				}
			}

			// 先物商品区分の検証
			if future.DerivativesProductCategory == "" {
				t.Errorf("Future[%d]: DerivativesProductCategory is empty", i)
			}

			// 緊急取引証拠金発動区分の検証
			if future.EmergencyMarginTriggerDivision != "001" && future.EmergencyMarginTriggerDivision != "002" {
				t.Errorf("Future[%d]: EmergencyMarginTriggerDivision = %v, want 001 or 002", i, future.EmergencyMarginTriggerDivision)
			}

			// 日通し四本値の論理的整合性チェック
			if future.WholeDayHigh < future.WholeDayLow {
				t.Errorf("Future[%d]: WholeDayHigh (%v) < WholeDayLow (%v)", i, future.WholeDayHigh, future.WholeDayLow)
			}
			if future.WholeDayOpen > 0 && future.WholeDayHigh > 0 && future.WholeDayOpen > future.WholeDayHigh {
				t.Errorf("Future[%d]: WholeDayOpen (%v) > WholeDayHigh (%v)", i, future.WholeDayOpen, future.WholeDayHigh)
			}
			if future.WholeDayOpen > 0 && future.WholeDayLow > 0 && future.WholeDayOpen < future.WholeDayLow {
				t.Errorf("Future[%d]: WholeDayOpen (%v) < WholeDayLow (%v)", i, future.WholeDayOpen, future.WholeDayLow)
			}

			// 日中四本値の論理的整合性チェック
			if future.DaySessionHigh < future.DaySessionLow {
				t.Errorf("Future[%d]: DaySessionHigh (%v) < DaySessionLow (%v)", i, future.DaySessionHigh, future.DaySessionLow)
			}

			// ナイトセッション四本値のチェック（データがある場合）
			if future.HasNightSession() {
				nightHigh := future.GetNightSessionHigh()
				nightLow := future.GetNightSessionLow()
				if nightHigh != nil && nightLow != nil && *nightHigh < *nightLow {
					t.Errorf("Future[%d]: NightSessionHigh (%v) < NightSessionLow (%v)", i, *nightHigh, *nightLow)
				}
			}

			// 出来高・建玉の妥当性チェック
			if future.Volume < 0 {
				t.Errorf("Future[%d]: Volume = %v, want >= 0", i, future.Volume)
			}
			if future.OpenInterest < 0 {
				t.Errorf("Future[%d]: OpenInterest = %v, want >= 0", i, future.OpenInterest)
			}
			if future.TurnoverValue < 0 {
				t.Errorf("Future[%d]: TurnoverValue = %v, want >= 0", i, future.TurnoverValue)
			}

			// 先物商品区分の検証（有効な値かチェック）
			validCategories := map[string]bool{
				"TOPIXF":   true, // TOPIX先物
				"TOPIXMF":  true, // ミニTOPIX先物
				"MOTF":     true, // マザーズ先物
				"NKVIF":    true, // 日経平均VI先物
				"NKYDF":    true, // 日経平均・配当指数先物
				"NK225F":   true, // 日経225先物
				"NK225MF":  true, // 日経225mini先物
				"JN400F":   true, // JPX日経インデックス400先物
				"REITF":    true, // 東証REIT指数先物
				"DJIAF":    true, // NYダウ先物
				"JGBLF":    true, // 長期国債先物
				"NK225MCF": true, // 日経225マイクロ先物
				"TOA3MF":   true, // TONA3ヶ月金利先物
			}
			if !validCategories[future.DerivativesProductCategory] {
				t.Logf("Future[%d]: Unknown DerivativesProductCategory: %s", i, future.DerivativesProductCategory)
			}

			// 最初の5件の詳細ログ
			if i < 5 {
				t.Logf("Future[%d]: Date=%s, Code=%s", i, future.Date, future.Code)
				t.Logf("  Category: %s, Month: %s", future.DerivativesProductCategory, future.ContractMonth)
				t.Logf("  WholeDay OHLC: Open=%.1f, High=%.1f, Low=%.1f, Close=%.1f",
					future.WholeDayOpen, future.WholeDayHigh, future.WholeDayLow, future.WholeDayClose)
				t.Logf("  DaySession OHLC: Open=%.1f, High=%.1f, Low=%.1f, Close=%.1f",
					future.DaySessionOpen, future.DaySessionHigh, future.DaySessionLow, future.DaySessionClose)
				t.Logf("  Volume: %.0f, OpenInterest: %.0f, Turnover: %.0f",
					future.Volume, future.OpenInterest, future.TurnoverValue)

				// ナイトセッションデータの表示
				if future.HasNightSession() {
					nightOpen := future.GetNightSessionOpen()
					nightHigh := future.GetNightSessionHigh()
					nightLow := future.GetNightSessionLow()
					nightClose := future.GetNightSessionClose()
					t.Logf("  NightSession OHLC: Open=%.1f, High=%.1f, Low=%.1f, Close=%.1f",
						*nightOpen, *nightHigh, *nightLow, *nightClose)

					// ギャップの計算
					if gap := future.GetDayNightGap(); gap != nil {
						t.Logf("  Day-Night Gap: %.1f", *gap)
					}
				} else {
					t.Logf("  NightSession: No data")
				}

				// 前場データの表示
				if future.HasMorningSession() {
					morningOpen := future.GetMorningSessionOpen()
					morningHigh := future.GetMorningSessionHigh()
					morningLow := future.GetMorningSessionLow()
					morningClose := future.GetMorningSessionClose()
					t.Logf("  MorningSession OHLC: Open=%.1f, High=%.1f, Low=%.1f, Close=%.1f",
						*morningOpen, *morningHigh, *morningLow, *morningClose)
				}

				// 2016年7月19日以降の追加情報
				if future.SettlementPrice != nil {
					t.Logf("  Settlement: %.1f", *future.SettlementPrice)
				}
				if future.LastTradingDay != nil && *future.LastTradingDay != "" {
					// LastTradingDayのフォーマット検証
					if len(*future.LastTradingDay) != 10 || (*future.LastTradingDay)[4] != '-' || (*future.LastTradingDay)[7] != '-' {
						t.Errorf("Future[%d]: LastTradingDay format invalid = %v, want YYYY-MM-DD", i, *future.LastTradingDay)
					}
					t.Logf("  Last Trading Day: %s", *future.LastTradingDay)
				}
				if future.SpecialQuotationDay != nil && *future.SpecialQuotationDay != "" {
					// SpecialQuotationDayのフォーマット検証
					if len(*future.SpecialQuotationDay) != 10 || (*future.SpecialQuotationDay)[4] != '-' || (*future.SpecialQuotationDay)[7] != '-' {
						t.Errorf("Future[%d]: SpecialQuotationDay format invalid = %v, want YYYY-MM-DD", i, *future.SpecialQuotationDay)
					}
					t.Logf("  SQ Day: %s", *future.SpecialQuotationDay)
				}
				if future.VolumeOnlyAuction != nil {
					// Volume(OnlyAuction)の妥当性チェック
					if *future.VolumeOnlyAuction < 0 {
						t.Errorf("Future[%d]: VolumeOnlyAuction = %v, want >= 0", i, *future.VolumeOnlyAuction)
					}
					t.Logf("  Volume (Only Auction): %.0f", *future.VolumeOnlyAuction)
				}
				if future.CentralContractMonthFlag != nil {
					// 中心限月フラグの検証
					if *future.CentralContractMonthFlag != "0" && *future.CentralContractMonthFlag != "1" {
						t.Errorf("Future[%d]: CentralContractMonthFlag = %v, want 0 or 1", i, *future.CentralContractMonthFlag)
					}
				}
				if future.IsCentralContractMonth() {
					t.Logf("  Central Contract Month: Yes")
				}
				if future.IsEmergencyMarginTriggered() {
					t.Logf("  Emergency Margin Triggered: Yes")
				}
			}
		}

		t.Logf("Retrieved %d futures records", len(futures))
	})

	t.Run("GetFutures_ByCategory", func(t *testing.T) {
		// 特定カテゴリ（日経225先物）のみを取得
		date := getTestDate()

		futures, err := jq.Futures.GetFuturesByCategory(date, "NK225F")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get NK225F futures: %v", err)
		}

		if len(futures) == 0 {
			t.Skip("No NK225F futures data available")
		}

		// 全て日経225先物か確認
		for i, future := range futures {
			if future.DerivativesProductCategory != "NK225F" {
				t.Errorf("Future[%d]: Expected NK225F but got %s", i, future.DerivativesProductCategory)
			}
		}

		t.Logf("Retrieved %d NK225F futures", len(futures))
	})

	t.Run("GetCentralContractMonthFutures", func(t *testing.T) {
		// 中心限月のみを取得
		date := getTestDate()

		futures, err := jq.Futures.GetCentralContractMonthFutures(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get central contract month futures: %v", err)
		}

		if len(futures) == 0 {
			t.Skip("No central contract month futures data available")
		}

		// 全て中心限月か確認
		for i, future := range futures {
			if !future.IsCentralContractMonth() {
				t.Errorf("Future[%d]: Expected central contract month but got CentralContractMonthFlag = %v",
					i, future.CentralContractMonthFlag)
			}
		}

		t.Logf("Retrieved %d central contract month futures", len(futures))
	})

	t.Run("GetFutures_ContractMonthAnalysis", func(t *testing.T) {
		// 限月別の分析
		date := getTestDate()

		futures, err := jq.Futures.GetFuturesByDate(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No futures data available for contract month analysis")
		}

		if len(futures) == 0 {
			t.Skip("No futures data available")
		}

		// 限月別の集計
		monthlyData := make(map[string]struct {
			count       int
			totalVolume float64
			totalOI     float64
			categories  map[string]int
		})

		for _, future := range futures {
			data := monthlyData[future.ContractMonth]
			if data.categories == nil {
				data.categories = make(map[string]int)
			}
			data.count++
			data.totalVolume += future.Volume
			data.totalOI += future.OpenInterest
			data.categories[future.DerivativesProductCategory]++
			monthlyData[future.ContractMonth] = data
		}

		t.Logf("Contract month analysis:")
		for month, data := range monthlyData {
			t.Logf("  %s: %d contracts, Volume=%.0f, OI=%.0f",
				month, data.count, data.totalVolume, data.totalOI)

			for category, count := range data.categories {
				t.Logf("    %s: %d contracts", category, count)
			}
		}
	})

	t.Run("GetFutures_CategoryAnalysis", func(t *testing.T) {
		// 商品区分別の分析
		date := getTestDate()

		futures, err := jq.Futures.GetFuturesByDate(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No futures data available for category analysis")
		}

		if len(futures) == 0 {
			t.Skip("No futures data available")
		}

		// 商品区分別の集計
		categoryData := make(map[string]struct {
			count       int
			totalVolume float64
			totalOI     float64
			contracts   []string
		})

		for _, future := range futures {
			data := categoryData[future.DerivativesProductCategory]
			data.count++
			data.totalVolume += future.Volume
			data.totalOI += future.OpenInterest
			data.contracts = append(data.contracts, future.ContractMonth)
			categoryData[future.DerivativesProductCategory] = data
		}

		t.Logf("Category analysis:")
		for category, data := range categoryData {
			t.Logf("  %s: %d contracts, Volume=%.0f, OI=%.0f",
				category, data.count, data.totalVolume, data.totalOI)
		}
	})

	t.Run("GetFutures_PriceAnalysis", func(t *testing.T) {
		// 価格・ボラティリティ分析
		date := getTestDate()

		futures, err := jq.Futures.GetFuturesByDate(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No futures data available for price analysis")
		}

		if len(futures) == 0 {
			t.Skip("No futures data available")
		}

		totalVolume := float64(0)
		totalTurnover := float64(0)
		rangeSum := float64(0)
		validPriceCount := 0

		for _, future := range futures {
			totalVolume += future.Volume
			totalTurnover += future.TurnoverValue

			// 日通しレンジの計算
			if future.WholeDayHigh > 0 && future.WholeDayLow > 0 {
				rangeSum += future.GetWholeDayRange()
				validPriceCount++
			}
		}

		t.Logf("Price analysis for %s:", date)
		t.Logf("  Total contracts: %d", len(futures))
		t.Logf("  Total volume: %.0f contracts", totalVolume)
		t.Logf("  Total turnover: %.0f yen", totalTurnover)

		if validPriceCount > 0 {
			avgRange := rangeSum / float64(validPriceCount)
			t.Logf("  Average daily range: %.2f points", avgRange)
		}
	})

	t.Run("GetFutures_SessionAnalysis", func(t *testing.T) {
		// セッション分析（ナイト・日中・前場）
		date := getTestDate()

		futures, err := jq.Futures.GetFuturesByDate(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No futures data available for session analysis")
		}

		if len(futures) == 0 {
			t.Skip("No futures data available")
		}

		nightSessionCount := 0
		morningSessionCount := 0
		gapAnalysis := make([]float64, 0)

		for _, future := range futures {
			if future.HasNightSession() {
				nightSessionCount++

				// ギャップ分析
				if gap := future.GetDayNightGap(); gap != nil {
					gapAnalysis = append(gapAnalysis, *gap)
				}
			}

			if future.HasMorningSession() {
				morningSessionCount++
			}
		}

		t.Logf("Session analysis:")
		t.Logf("  Contracts with night session: %d/%d", nightSessionCount, len(futures))
		t.Logf("  Contracts with morning session: %d/%d", morningSessionCount, len(futures))

		if len(gapAnalysis) > 0 {
			avgGap := float64(0)
			for _, gap := range gapAnalysis {
				avgGap += gap
			}
			avgGap /= float64(len(gapAnalysis))
			t.Logf("  Average day-night gap: %.2f points", avgGap)
		}
	})

	t.Run("GetFutures_EmergencyMarginAnalysis", func(t *testing.T) {
		// 緊急取引証拠金発動区分の分析
		date := getTestDate()

		params := jquants.FuturesParams{
			Date: date,
		}

		resp, err := jq.Futures.GetFutures(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No futures data available for emergency margin analysis")
		}

		if resp == nil || len(resp.Futures) == 0 {
			t.Skip("No futures data available")
		}

		normalCount := 0
		emergencyCount := 0

		for _, future := range resp.Futures {
			if future.IsEmergencyMarginTriggered() {
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

	t.Run("GetFutures_ErrorHandling", func(t *testing.T) {
		// エラーケースのテスト

		// 未来の日付
		futures, err := jq.Futures.GetFuturesByDate("2030-01-01")
		if err == nil && len(futures) > 0 {
			t.Error("Expected error or empty result for future date")
		}

		// 無効な日付形式
		futures, err = jq.Futures.GetFuturesByDate("invalid-date")
		if err == nil && len(futures) > 0 {
			t.Error("Expected error or empty result for invalid date format")
		}

		// 空の日付（必須パラメータ）
		params := jquants.FuturesParams{
			Date: "",
		}

		_, err = jq.Futures.GetFutures(params)
		if err == nil {
			t.Error("Expected error for missing required date parameter")
		}

		// 無効なカテゴリ
		futures, err = jq.Futures.GetFuturesByCategory(getTestDate(), "INVALID_CATEGORY")
		if err == nil && len(futures) > 0 {
			t.Error("Expected error or empty result for invalid category")
		}
	})
}
