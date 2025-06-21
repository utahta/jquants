//go:build e2e
// +build e2e

package e2e

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants"
)

// TestPricesAMEndpoint は/prices/prices_am（プレミアムプラン専用）エンドポイントの完全なテスト
func TestPricesAMEndpoint(t *testing.T) {
	t.Run("GetPricesAM_ByCode", func(t *testing.T) {
		// トヨタ自動車の前場四本値を取得
		resp, err := jq.PricesAM.GetPricesAMByCode("7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get prices AM: %v", err)
		}

		if resp == nil || len(resp.PricesAM) == 0 {
			t.Skip("No prices AM data available")
		}

		// 各前場四本値データを詳細に検証
		for i, quote := range resp.PricesAM {
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

			// 四本値の論理的整合性チェック
			if quote.MorningHigh != nil && quote.MorningLow != nil && *quote.MorningHigh < *quote.MorningLow {
				t.Errorf("Quote[%d]: MorningHigh (%v) < MorningLow (%v)", i, *quote.MorningHigh, *quote.MorningLow)
			}
			if quote.MorningOpen != nil && quote.MorningHigh != nil && *quote.MorningOpen > *quote.MorningHigh {
				t.Errorf("Quote[%d]: MorningOpen (%v) > MorningHigh (%v)", i, *quote.MorningOpen, *quote.MorningHigh)
			}
			if quote.MorningOpen != nil && quote.MorningLow != nil && *quote.MorningOpen < *quote.MorningLow {
				t.Errorf("Quote[%d]: MorningOpen (%v) < MorningLow (%v)", i, *quote.MorningOpen, *quote.MorningLow)
			}
			if quote.MorningClose != nil && quote.MorningHigh != nil && *quote.MorningClose > *quote.MorningHigh {
				t.Errorf("Quote[%d]: MorningClose (%v) > MorningHigh (%v)", i, *quote.MorningClose, *quote.MorningHigh)
			}
			if quote.MorningClose != nil && quote.MorningLow != nil && *quote.MorningClose < *quote.MorningLow {
				t.Errorf("Quote[%d]: MorningClose (%v) < MorningLow (%v)", i, *quote.MorningClose, *quote.MorningLow)
			}

			// 価格の妥当性チェック（負の値は通常ありえない）
			if quote.MorningOpen != nil && *quote.MorningOpen < 0 {
				t.Errorf("Quote[%d]: MorningOpen = %v, want >= 0", i, *quote.MorningOpen)
			}
			if quote.MorningHigh != nil && *quote.MorningHigh < 0 {
				t.Errorf("Quote[%d]: MorningHigh = %v, want >= 0", i, *quote.MorningHigh)
			}
			if quote.MorningLow != nil && *quote.MorningLow < 0 {
				t.Errorf("Quote[%d]: MorningLow = %v, want >= 0", i, *quote.MorningLow)
			}
			if quote.MorningClose != nil && *quote.MorningClose < 0 {
				t.Errorf("Quote[%d]: MorningClose = %v, want >= 0", i, *quote.MorningClose)
			}

			// 出来高の妥当性チェック
			if quote.MorningVolume != nil && *quote.MorningVolume < 0 {
				t.Errorf("Quote[%d]: MorningVolume = %v, want >= 0", i, *quote.MorningVolume)
			}
			if quote.MorningTurnoverValue != nil && *quote.MorningTurnoverValue < 0 {
				t.Errorf("Quote[%d]: MorningTurnoverValue = %v, want >= 0", i, *quote.MorningTurnoverValue)
			}

			// 最初の5件の詳細ログ
			if i < 5 {
				t.Logf("Quote[%d]: Date=%s, Code=%s", i, quote.Date, quote.Code)
				openStr := "nil"
				if quote.MorningOpen != nil {
					openStr = fmt.Sprintf("%.0f", *quote.MorningOpen)
				}
				highStr := "nil"
				if quote.MorningHigh != nil {
					highStr = fmt.Sprintf("%.0f", *quote.MorningHigh)
				}
				lowStr := "nil"
				if quote.MorningLow != nil {
					lowStr = fmt.Sprintf("%.0f", *quote.MorningLow)
				}
				closeStr := "nil"
				if quote.MorningClose != nil {
					closeStr = fmt.Sprintf("%.0f", *quote.MorningClose)
				}
				volumeStr := "nil"
				if quote.MorningVolume != nil {
					volumeStr = fmt.Sprintf("%.0f", *quote.MorningVolume)
				}
				turnoverStr := "nil"
				if quote.MorningTurnoverValue != nil {
					turnoverStr = fmt.Sprintf("%.0f", *quote.MorningTurnoverValue)
				}
				t.Logf("  OHLC: Open=%s, High=%s, Low=%s, Close=%s", openStr, highStr, lowStr, closeStr)
				t.Logf("  Volume: %s, TurnoverValue: %s", volumeStr, turnoverStr)
			}
		}

		t.Logf("Retrieved %d prices AM records", len(resp.PricesAM))
	})

	t.Run("GetPricesAM_All", func(t *testing.T) {
		// 全銘柄の前場四本値を取得（当日のデータのみ利用可能）
		params := jquants.PricesAMParams{}

		resp, err := jq.PricesAM.GetPricesAM(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Logf("Failed to get prices AM by date: %v", err)
			return
		}

		if resp == nil || len(resp.PricesAM) == 0 {
			t.Skip("No prices AM data for the specified date")
		}

		t.Logf("Retrieved %d prices AM records for all stocks", len(resp.PricesAM))

		// 基本的な検証（上位10件）
		for i, quote := range resp.PricesAM {
			if quote.Code == "" {
				t.Errorf("Quote[%d]: Code is empty", i)
			} else {
				// 銘柄コードは4桁または5桁
				if len(quote.Code) != 4 && len(quote.Code) != 5 {
					t.Errorf("Quote[%d]: Code length = %d, want 4 or 5", i, len(quote.Code))
				}
			}
			if quote.Date == "" {
				t.Errorf("Quote[%d]: Date is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(quote.Date) != 10 || quote.Date[4] != '-' || quote.Date[7] != '-' {
					t.Errorf("Quote[%d]: Date format invalid = %v, want YYYY-MM-DD", i, quote.Date)
				}
			}
			if i >= 10 {
				break // 最初の10件のみ確認
			}
		}
	})

	t.Run("GetPricesAM_MarketAnalysis", func(t *testing.T) {
		// 市場全体の前場四本値分析（当日のデータのみ）
		params := jquants.PricesAMParams{}

		resp, err := jq.PricesAM.GetPricesAM(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No data available for market analysis")
		}

		if resp == nil || len(resp.PricesAM) == 0 {
			t.Skip("No data available")
		}

		// 市場統計の計算
		totalVolume := float64(0)
		totalTurnover := float64(0)
		priceChanges := 0
		priceDrops := 0
		priceRises := 0

		for _, quote := range resp.PricesAM {
			if quote.MorningVolume != nil {
				totalVolume += *quote.MorningVolume
			}
			if quote.MorningTurnoverValue != nil {
				totalTurnover += *quote.MorningTurnoverValue
			}

			// 前場の価格変動分析
			if quote.MorningOpen != nil && quote.MorningClose != nil && *quote.MorningOpen > 0 && *quote.MorningClose > 0 {
				if *quote.MorningClose > *quote.MorningOpen {
					priceRises++
				} else if *quote.MorningClose < *quote.MorningOpen {
					priceDrops++
				}
				priceChanges++
			}
		}

		t.Logf("Market analysis for today's morning session:")
		t.Logf("  Total stocks: %d", len(resp.PricesAM))
		t.Logf("  Total volume: %.0f million shares", totalVolume/1000000)
		t.Logf("  Total turnover: %.0f billion yen", totalTurnover/1000000000)

		if priceChanges > 0 {
			riseRatio := float64(priceRises) / float64(priceChanges) * 100
			dropRatio := float64(priceDrops) / float64(priceChanges) * 100
			t.Logf("  Price movements: %.1f%% rise, %.1f%% drop", riseRatio, dropRatio)
		}
	})

	t.Run("GetPricesAM_Pagination", func(t *testing.T) {
		// ページネーションのテスト（当日のデータのみ）
		params := jquants.PricesAMParams{}

		resp, err := jq.PricesAM.GetPricesAM(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No prices AM data available")
		}

		if resp == nil || len(resp.PricesAM) == 0 {
			t.Skip("No data available for pagination test")
		}

		firstPageCount := len(resp.PricesAM)
		t.Logf("First page: %d records", firstPageCount)

		if resp.PaginationKey != "" {
			// 次のページを取得
			params.PaginationKey = resp.PaginationKey
			resp2, err := jq.PricesAM.GetPricesAM(params)
			if err != nil {
				t.Fatalf("Failed to get next page: %v", err)
			}

			if resp2 != nil && len(resp2.PricesAM) > 0 {
				t.Logf("Second page: %d records", len(resp2.PricesAM))
			}
		}
	})

	t.Run("GetPricesAM_ErrorHandling", func(t *testing.T) {
		// エラーケースのテスト

		// 存在しない銘柄コード
		resp, err := jq.PricesAM.GetPricesAMByCode("99999")
		if err == nil && resp != nil && len(resp.PricesAM) > 0 {
			t.Error("Expected error or empty result for invalid code")
		}

		// 無効なコード
		params := jquants.PricesAMParams{
			Code: "invalid-code",
		}

		resp, err = jq.PricesAM.GetPricesAM(params)
		if err == nil && resp != nil && len(resp.PricesAM) > 0 {
			t.Error("Expected error or empty result for invalid code")
		}
	})
}
