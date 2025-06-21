//go:build e2e
// +build e2e

package e2e

import (
	"testing"
	"time"
)

// TestShortSellingPositionsEndpoint は/markets/short_selling_positionsエンドポイントの完全なテスト
func TestShortSellingPositionsEndpoint(t *testing.T) {
	t.Run("GetShortSellingPositions_ByCode", func(t *testing.T) {
		// トヨタ自動車の空売り残高報告を取得
		positions, err := jq.ShortSellingPositions.GetShortSellingPositionsByCode("7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get short selling positions: %v", err)
		}

		if len(positions) == 0 {
			t.Skip("No short selling positions data available")
		}

		// 各空売り残高報告を詳細に検証
		for i, position := range positions {
			// 基本情報の検証
			if position.Code != "72030" && position.Code != "7203" {
				t.Errorf("Position[%d]: Code = %v, want 72030 or 7203", i, position.Code)
			}
			if position.DisclosedDate == "" {
				t.Errorf("Position[%d]: DisclosedDate is empty", i)
			}
			if position.CalculatedDate == "" {
				t.Errorf("Position[%d]: CalculatedDate is empty", i)
			}
			if position.ShortSellerName == "" {
				t.Errorf("Position[%d]: ShortSellerName is empty", i)
			}
			
			// 日付の妥当性チェック
			if position.DisclosedDate < position.CalculatedDate {
				t.Errorf("Position[%d]: DisclosedDate (%s) < CalculatedDate (%s)",
					i, position.DisclosedDate, position.CalculatedDate)
			}
			
			// 残高割合の検証（通常0.5%以上で報告義務）
			if position.ShortPositionsToSharesOutstandingRatio < 0.5 {
				t.Logf("Position[%d]: Low ratio: %.2f%% (might be below reporting threshold)",
					i, position.ShortPositionsToSharesOutstandingRatio)
			}
			if position.ShortPositionsToSharesOutstandingRatio > 50 {
				t.Errorf("Position[%d]: Extremely high ratio: %.2f%%",
					i, position.ShortPositionsToSharesOutstandingRatio)
			}
			
			// 株数の検証
			if position.ShortPositionsInSharesNumber <= 0 {
				t.Errorf("Position[%d]: ShortPositionsInSharesNumber = %v, want > 0",
					i, position.ShortPositionsInSharesNumber)
			}
			if position.ShortPositionsInTradingUnitsNumber <= 0 {
				t.Errorf("Position[%d]: ShortPositionsInTradingUnitsNumber = %v, want > 0",
					i, position.ShortPositionsInTradingUnitsNumber)
			}
			
			// 住所情報の検証（オプショナル）
			if position.ShortSellerAddress == "" {
				t.Logf("Position[%d]: ShortSellerAddress is empty", i)
			}
			
			// 投資一任契約の情報確認
			hasDiscretionary := position.HasDiscretionaryInvestment()
			if hasDiscretionary {
				t.Logf("Position[%d]: Has discretionary investment contract", i)
				if position.DiscretionaryInvestmentContractorName != "" {
					t.Logf("  Contractor: %s", position.DiscretionaryInvestmentContractorName)
				}
				if position.InvestmentFundName != "" {
					t.Logf("  Fund: %s", position.InvestmentFundName)
				}
			}
			
			// 個人投資家かチェック
			if position.IsIndividual() {
				t.Logf("Position[%d]: Individual investor", i)
			}
			
			// 最初の5件の詳細ログ
			if i < 5 {
				t.Logf("Position[%d]: %s", i, position.ShortSellerName)
				t.Logf("  Disclosed: %s, Calculated: %s", position.DisclosedDate, position.CalculatedDate)
				t.Logf("  Ratio: %.2f%%, Shares: %.0f", 
					position.ShortPositionsToSharesOutstandingRatio, position.ShortPositionsInSharesNumber)
				
				// 前回からの変化
				changeRatio := position.GetPositionChangeRatio()
				if changeRatio != 0 {
					direction := "increased"
					if changeRatio < 0 {
						direction = "decreased"
					}
					t.Logf("  Position %s by %.2f%% from previous report", direction, abs(changeRatio))
				}
			}
		}
		
		t.Logf("Retrieved %d short selling positions for code 7203", len(positions))
	})

	t.Run("GetShortSellingPositions_ByDisclosedDate", func(t *testing.T) {
		// 最近の金曜日の全銘柄空売り残高報告を取得
		disclosedDate := getRecentFriday()
		
		positions, err := jq.ShortSellingPositions.GetShortSellingPositionsByDisclosedDate(disclosedDate)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			// データがない可能性もある
			t.Logf("No positions data for date %s: %v", disclosedDate, err)
			return
		}

		if len(positions) == 0 {
			t.Skipf("No short selling positions data for %s", disclosedDate)
		}

		// 日付の一致確認
		for i, position := range positions {
			if position.DisclosedDate != disclosedDate {
				t.Errorf("Position[%d]: DisclosedDate = %v, want %v", i, position.DisclosedDate, disclosedDate)
			}
			if i >= 10 {
				break // 最初の10件のみ確認
			}
		}
		
		// 銘柄別の集計
		codeCount := make(map[string]int)
		sellerCount := make(map[string]int)
		
		for _, position := range positions {
			codeCount[position.Code]++
			sellerCount[position.ShortSellerName]++
		}
		
		t.Logf("Disclosed date %s summary:", disclosedDate)
		t.Logf("  Total positions: %d", len(positions))
		t.Logf("  Unique codes: %d", len(codeCount))
		t.Logf("  Unique sellers: %d", len(sellerCount))
		
		// 最も多く報告している空売り者
		maxCount := 0
		topSeller := ""
		for seller, count := range sellerCount {
			if count > maxCount {
				maxCount = count
				topSeller = seller
			}
		}
		if topSeller != "" {
			t.Logf("  Top seller: %s (%d positions)", topSeller, maxCount)
		}
	})

	t.Run("GetShortSellingPositions_PositionAnalysis", func(t *testing.T) {
		// トヨタ自動車の残高変化分析
		positions, err := jq.ShortSellingPositions.GetShortSellingPositionsByCode("7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No positions data available for analysis")
		}

		if len(positions) < 2 {
			t.Skip("Insufficient data for position analysis")
		}

		// 残高変化の分析
		increases := 0
		decreases := 0
		noChanges := 0
		
		for _, position := range positions {
			if position.IsIncrease() {
				increases++
			} else if position.IsDecrease() {
				decreases++
			} else if position.IsNoChange() {
				noChanges++
			}
		}
		
		t.Logf("Position change analysis for code 7203:")
		t.Logf("  Increases: %d", increases)
		t.Logf("  Decreases: %d", decreases)
		t.Logf("  No changes: %d", noChanges)
		
		// 大きな変化があったポジション
		t.Logf("Significant position changes:")
		for i, position := range positions {
			if i >= 5 {
				break
			}
			
			changeRatio := position.GetPositionChangeRatio()
			if abs(changeRatio) > 10 { // 10%以上の変化
				direction := "increased"
				if changeRatio < 0 {
					direction = "decreased"
				}
				t.Logf("  %s: %s by %.2f%% (from %.2f%% to %.2f%%)",
					position.ShortSellerName, direction, abs(changeRatio),
					position.ShortPositionsInPreviousReportingRatio,
					position.ShortPositionsToSharesOutstandingRatio)
			}
		}
	})

	t.Run("GetShortSellingPositions_DateRange", func(t *testing.T) {
		// 過去1ヶ月の期間でトヨタ自動車の空売り残高を取得
		to := time.Now().Format("2006-01-02")
		from := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
		
		positions, err := jq.ShortSellingPositions.GetShortSellingPositionsByCodeAndDateRange("7203", from, to)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get positions for date range: %v", err)
		}

		if len(positions) == 0 {
			t.Skip("No positions data for the specified date range")
		}

		t.Logf("Retrieved %d positions for period %s to %s", len(positions), from, to)
		
		// 期間内の日付確認
		for _, position := range positions {
			if position.DisclosedDate < from || position.DisclosedDate > to {
				t.Errorf("Position disclosed date %s is outside range %s to %s",
					position.DisclosedDate, from, to)
			}
		}
		
		// 時系列トレンド
		if len(positions) > 1 {
			t.Logf("Monthly trend analysis:")
			prevRatio := 0.0
			for i, position := range positions {
				if i >= 5 {
					break
				}
				
				currentRatio := position.ShortPositionsToSharesOutstandingRatio
				if i > 0 {
					change := currentRatio - prevRatio
					t.Logf("  %s: %.2f%% (change: %+.2f%%)",
						position.DisclosedDate, currentRatio, change)
				} else {
					t.Logf("  %s: %.2f%%", position.DisclosedDate, currentRatio)
				}
				prevRatio = currentRatio
			}
		}
	})

	t.Run("GetShortSellingPositions_ThresholdAnalysis", func(t *testing.T) {
		// 閾値分析（報告義務の閾値分析）
		positions, err := jq.ShortSellingPositions.GetShortSellingPositionsByCode("7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No positions data available for threshold analysis")
		}

		if len(positions) == 0 {
			t.Skip("No positions data available")
		}

		// 閾値別の分類
		below1Percent := 0
		oneToFivePercent := 0
		aboveFivePercent := 0
		
		for _, position := range positions {
			ratio := position.ShortPositionsToSharesOutstandingRatio
			if ratio < 1.0 {
				below1Percent++
			} else if ratio <= 5.0 {
				oneToFivePercent++
			} else {
				aboveFivePercent++
			}
		}
		
		t.Logf("Position threshold analysis:")
		t.Logf("  Below 1%%: %d positions", below1Percent)
		t.Logf("  1%% to 5%%: %d positions", oneToFivePercent)
		t.Logf("  Above 5%%: %d positions", aboveFivePercent)
		
		// 最大ポジション
		maxRatio := 0.0
		maxSeller := ""
		for _, position := range positions {
			if position.ShortPositionsToSharesOutstandingRatio > maxRatio {
				maxRatio = position.ShortPositionsToSharesOutstandingRatio
				maxSeller = position.ShortSellerName
			}
		}
		
		if maxSeller != "" {
			t.Logf("Largest position: %s with %.2f%%", maxSeller, maxRatio)
		}
	})

	t.Run("GetShortSellingPositions_ErrorHandling", func(t *testing.T) {
		// エラーケースのテスト
		
		// 存在しない銘柄コード
		positions, err := jq.ShortSellingPositions.GetShortSellingPositionsByCode("99999")
		if err == nil && len(positions) > 0 {
			t.Error("Expected error or empty result for invalid code")
		}
		
		// 無効な日付
		positions, err = jq.ShortSellingPositions.GetShortSellingPositionsByDisclosedDate("invalid-date")
		if err == nil && len(positions) > 0 {
			t.Error("Expected error or empty result for invalid date")
		}
	})
}

// abs は絶対値を返すヘルパー関数
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// getRecentFriday は最近の金曜日を取得する
func getRecentFriday() string {
	now := time.Now()
	// 今日が金曜日でない場合は、前の金曜日を探す
	for now.Weekday() != time.Friday {
		now = now.AddDate(0, 0, -1)
	}
	// さらに1週間前にする（当週のデータは利用できない可能性があるため）
	now = now.AddDate(0, 0, -7)
	return now.Format("2006-01-02")
}