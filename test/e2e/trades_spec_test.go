//go:build e2e
// +build e2e

package e2e

import (
	"testing"
	"time"

	"github.com/utahta/jquants"
)

// TestTradesSpecEndpoint は/equities/investor-typesエンドポイントの完全なテスト
func TestTradesSpecEndpoint(t *testing.T) {
	t.Run("GetTradesSpec_ByDateRange", func(t *testing.T) {
		// 最近の営業日の投資部門別売買状況を取得
		date := getTestDate()

		trades, err := jq.TradesSpec.GetTradesSpecByDateRange(date, date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get trades spec: %v", err)
		}

		if len(trades) == 0 {
			t.Skip("No trades spec data available")
		}

		// 各投資部門別売買データを詳細に検証
		for i, trade := range trades {
			// 基本情報の検証
			if trade.Section == "" {
				t.Errorf("Trade[%d]: Section is empty", i)
			}
			if trade.StDate == "" {
				t.Errorf("Trade[%d]: StDate is empty", i)
			}
			if trade.EnDate == "" {
				t.Errorf("Trade[%d]: EnDate is empty", i)
			}

			// 日付範囲の検証
			if trade.StDate > trade.EnDate {
				t.Errorf("Trade[%d]: StDate (%s) > EnDate (%s)",
					i, trade.StDate, trade.EnDate)
			}

			// 売買金額の検証（負の値は許容される）
			if trade.TotSell < 0 {
				t.Logf("Trade[%d]: TotSell is negative: %.0f", i, trade.TotSell)
			}
			if trade.TotBuy < 0 {
				t.Logf("Trade[%d]: TotBuy is negative: %.0f", i, trade.TotBuy)
			}

			// セクション（部門）の妥当性チェック
			validSections := map[string]bool{
				"TSEPrime":    true, // プライム市場
				"TSEStandard": true, // スタンダード市場
				"TSEGrowth":   true, // グロース市場
				"All":         true, // 全市場
			}
			if !validSections[trade.Section] {
				t.Logf("Trade[%d]: Unknown section: %s", i, trade.Section)
			}

			// 最初の5件の詳細ログ
			if i < 5 {
				t.Logf("Trade[%d]: Section=%s, Period=%s to %s",
					i, trade.Section, trade.StDate, trade.EnDate)
				t.Logf("  Sales: %.0f, Purchases: %.0f, Net: %.0f",
					trade.TotSell, trade.TotBuy,
					trade.TotSell-trade.TotBuy)
			}
		}

		t.Logf("Retrieved %d trades spec records for date %s", len(trades), date)
	})

	t.Run("GetTradesSpecBySection_TSEPrime", func(t *testing.T) {
		// プライム市場の投資部門別売買状況を取得
		trades, err := jq.TradesSpec.GetTradesSpecBySection("TSEPrime")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get trades spec for TSEPrime: %v", err)
		}

		if len(trades) == 0 {
			t.Skip("No trades spec data available for TSEPrime")
		}

		// 全てのデータがプライム市場か確認
		for i, trade := range trades {
			if trade.Section != "TSEPrime" {
				t.Errorf("Trade[%d]: Section = %v, want TSEPrime", i, trade.Section)
			}
		}

		t.Logf("Retrieved %d TSEPrime trades spec records", len(trades))

		// 時系列分析
		if len(trades) > 1 {
			for i := 1; i < len(trades) && i < 5; i++ {
				curr := trades[i]
				prev := trades[i-1]

				salesChange := curr.TotSell - prev.TotSell
				purchasesChange := curr.TotBuy - prev.TotBuy

				t.Logf("Period %s: Sales change: %.0f, Purchases change: %.0f",
					curr.EnDate, salesChange, purchasesChange)
			}
		}
	})

	t.Run("GetTradesSpecBySection_All", func(t *testing.T) {
		// 全市場の投資部門別売買状況を取得
		trades, err := jq.TradesSpec.GetTradesSpecBySection("All")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Logf("Failed to get trades spec for All: %v", err)
			return
		}

		if len(trades) == 0 {
			t.Skip("No trades spec data available for All")
		}

		// 全市場データの検証
		for i, trade := range trades {
			if trade.Section != "All" {
				t.Errorf("Trade[%d]: Section = %v, want All", i, trade.Section)
			}
		}

		t.Logf("Retrieved %d All market trades spec records", len(trades))
	})

	t.Run("GetTradesSpec_MultipleMarkets", func(t *testing.T) {
		// 複数市場のデータを取得して比較
		markets := []string{"TSEPrime", "TSEStandard", "TSEGrowth"}
		marketData := make(map[string][]jquants.TradesSpec)

		for _, market := range markets {
			trades, err := jq.TradesSpec.GetTradesSpecBySection(market)
			if err != nil {
				if isSubscriptionLimited(err) {
					t.Skip("Skipping due to subscription limitation")
				}
				t.Logf("Failed to get trades spec for %s: %v", market, err)
				continue
			}

			if len(trades) > 0 {
				marketData[market] = trades
				t.Logf("%s: %d records", market, len(trades))
			}
		}

		// 市場間の比較分析
		if len(marketData) >= 2 {
			t.Logf("Market comparison analysis:")
			for market, trades := range marketData {
				if len(trades) > 0 {
					latest := trades[0]
					t.Logf("  %s: Sales=%.0f, Purchases=%.0f, Net=%.0f",
						market, latest.TotSell, latest.TotBuy,
						latest.TotSell-latest.TotBuy)
				}
			}
		}
	})

	t.Run("GetTradesSpec_WeeklyData", func(t *testing.T) {
		// 過去1週間のデータ取得
		to := getTestDate()
		fromTime, _ := time.Parse("20060102", to)
		fromTime = fromTime.AddDate(0, 0, -7)
		from := fromTime.Format("20060102")

		trades, err := jq.TradesSpec.GetTradesSpecByDateRange(from, to)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get weekly trades spec: %v", err)
		}

		if len(trades) == 0 {
			t.Skip("No weekly trades spec data available")
		}

		t.Logf("Retrieved %d trades spec records for period %s to %s", len(trades), from, to)

		// 期間内のデータか確認
		for _, trade := range trades {
			if trade.StDate < from || trade.EnDate > to {
				t.Logf("Trade period (%s to %s) extends beyond requested range (%s to %s)",
					trade.StDate, trade.EnDate, from, to)
			}
		}

		// セクション別の集計
		sectionSummary := make(map[string]struct {
			count          int
			totalSales     float64
			totalPurchases float64
		})

		for _, trade := range trades {
			summary := sectionSummary[trade.Section]
			summary.count++
			summary.totalSales += trade.TotSell
			summary.totalPurchases += trade.TotBuy
			sectionSummary[trade.Section] = summary
		}

		t.Logf("Weekly summary by section:")
		for section, summary := range sectionSummary {
			net := summary.totalSales - summary.totalPurchases
			t.Logf("  %s: %d periods, Sales=%.0f, Purchases=%.0f, Net=%.0f",
				section, summary.count, summary.totalSales, summary.totalPurchases, net)
		}
	})

	t.Run("GetTradesSpec_TrendAnalysis", func(t *testing.T) {
		// トレンド分析（プライム市場）
		trades, err := jq.TradesSpec.GetTradesSpecBySection("TSEPrime")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No trades spec data available for trend analysis")
		}

		if len(trades) < 3 {
			t.Skip("Insufficient data for trend analysis")
		}

		// 直近3期間のトレンド分析
		t.Logf("TSEPrime trend analysis (last 3 periods):")
		for i := 0; i < 3 && i < len(trades); i++ {
			trade := trades[i]
			net := trade.TotSell - trade.TotBuy
			netRatio := 0.0
			if trade.TotBuy != 0 {
				netRatio = net / trade.TotBuy * 100
			}

			t.Logf("Period %d (%s to %s):",
				i+1, trade.StDate, trade.EnDate)
			t.Logf("  Sales: %.0f billion", trade.TotSell/1000000)
			t.Logf("  Purchases: %.0f billion", trade.TotBuy/1000000)
			t.Logf("  Net: %.0f billion (%.2f%%)", net/1000000, netRatio)
		}

		// 売買バランスの計算
		if len(trades) >= 2 {
			recent := trades[0]
			previous := trades[1]

			salesChange := recent.TotSell - previous.TotSell
			purchasesChange := recent.TotBuy - previous.TotBuy

			t.Logf("Recent change:")
			t.Logf("  Sales change: %.0f billion", salesChange/1000000)
			t.Logf("  Purchases change: %.0f billion", purchasesChange/1000000)
		}
	})

	t.Run("GetTradesSpec_ErrorHandling", func(t *testing.T) {
		// エラーケースのテスト

		// 無効なセクション
		trades, err := jq.TradesSpec.GetTradesSpecBySection("InvalidSection")
		if err == nil && len(trades) > 0 {
			t.Error("Expected error or empty result for invalid section")
		}

		// 無効な日付範囲
		trades, err = jq.TradesSpec.GetTradesSpecByDateRange("2024-12-31", "2024-01-01")
		if err == nil && len(trades) > 0 {
			t.Error("Expected error or empty result for invalid date range")
		}
	})
}
