//go:build e2e
// +build e2e

package e2e

import (
	"testing"
)

// TestShortSellingEndpoint は/markets/short_sellingエンドポイントの完全なテスト
func TestShortSellingEndpoint(t *testing.T) {
	t.Run("GetShortSelling_ByDate", func(t *testing.T) {
		// 最近の営業日の業種別空売り比率を取得
		date := getTestDate()
		
		shorts, err := jq.ShortSelling.GetShortSellingByDate(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get short selling data: %v", err)
		}

		if len(shorts) == 0 {
			t.Skip("No short selling data available")
		}

		// 各業種別空売りデータを詳細に検証
		for i, short := range shorts {
			// 基本情報の検証
			if short.Date == "" {
				t.Errorf("ShortSelling[%d]: Date is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(short.Date) != 10 || short.Date[4] != '-' || short.Date[7] != '-' {
					t.Errorf("ShortSelling[%d]: Date format invalid = %v, want YYYY-MM-DD", i, short.Date)
				}
				// 日付が指定した日付と一致するか確認
				expectedDate := getTestDateFormatted()
				if short.Date != expectedDate {
					t.Errorf("ShortSelling[%d]: Date = %v, want %v", i, short.Date, expectedDate)
				}
			}
			
			if short.Sector33Code == "" {
				t.Errorf("ShortSelling[%d]: Sector33Code is empty", i)
			} else {
				// 33業種コードは4桁
				if len(short.Sector33Code) != 4 {
					t.Errorf("ShortSelling[%d]: Sector33Code length = %d, want 4", i, len(short.Sector33Code))
				}
			}
			
			// 売買代金の検証（負の値は通常ありえない）
			if short.SellingExcludingShortSellingTurnoverValue < 0 {
				t.Errorf("ShortSelling[%d]: SellingExcludingShortSellingTurnoverValue = %v, want >= 0",
					i, short.SellingExcludingShortSellingTurnoverValue)
			}
			if short.ShortSellingWithRestrictionsTurnoverValue < 0 {
				t.Errorf("ShortSelling[%d]: ShortSellingWithRestrictionsTurnoverValue = %v, want >= 0",
					i, short.ShortSellingWithRestrictionsTurnoverValue)
			}
			if short.ShortSellingWithoutRestrictionsTurnoverValue < 0 {
				t.Errorf("ShortSelling[%d]: ShortSellingWithoutRestrictionsTurnoverValue = %v, want >= 0",
					i, short.ShortSellingWithoutRestrictionsTurnoverValue)
			}
			
			// 比率の計算と検証
			totalTurnover := short.SellingExcludingShortSellingTurnoverValue +
				short.ShortSellingWithRestrictionsTurnoverValue +
				short.ShortSellingWithoutRestrictionsTurnoverValue
			
			if totalTurnover > 0 {
				shortSellingRatio := (short.ShortSellingWithRestrictionsTurnoverValue +
					short.ShortSellingWithoutRestrictionsTurnoverValue) / totalTurnover * 100
				
				// 空売り比率が異常でないかチェック（通常0-50%程度）
				if shortSellingRatio > 100 {
					t.Errorf("ShortSelling[%d]: Short selling ratio too high: %.2f%%", i, shortSellingRatio)
				}
				
				// 最初の5件の詳細ログ
				if i < 5 {
					t.Logf("ShortSelling[%d]: Sector=%s, Date=%s",
						i, short.Sector33Code, short.Date)
					t.Logf("  Excl Short Selling: %.0f", short.SellingExcludingShortSellingTurnoverValue)
					t.Logf("  Short w/ Restrictions: %.0f", short.ShortSellingWithRestrictionsTurnoverValue)
					t.Logf("  Short w/o Restrictions: %.0f", short.ShortSellingWithoutRestrictionsTurnoverValue)
					t.Logf("  Total Turnover: %.0f", totalTurnover)
					t.Logf("  Short Selling Ratio: %.2f%%", shortSellingRatio)
				}
			}
		}
		
		t.Logf("Retrieved %d short selling records for date %s", len(shorts), date)
	})

	t.Run("GetShortSellingBySector", func(t *testing.T) {
		// 特定業種（輸送用機器：3050）の空売り比率データを取得
		sectorCode := "3050"
		
		shorts, err := jq.ShortSelling.GetShortSellingBySector(sectorCode)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get short selling data by sector: %v", err)
		}

		if len(shorts) == 0 {
			t.Skip("No short selling data available for sector 3050")
		}

		// 全てのデータが指定セクターか確認
		for i, short := range shorts {
			if short.Sector33Code != sectorCode {
				t.Errorf("ShortSelling[%d]: Sector33Code = %v, want %v", i, short.Sector33Code, sectorCode)
			}
			// 日付フォーマットの検証
			if short.Date != "" {
				if len(short.Date) != 10 || short.Date[4] != '-' || short.Date[7] != '-' {
					t.Errorf("ShortSelling[%d]: Date format invalid = %v, want YYYY-MM-DD", i, short.Date)
				}
			}
		}
		
		t.Logf("Retrieved %d short selling records for sector %s", len(shorts), sectorCode)
		
		// 時系列トレンド分析
		if len(shorts) > 1 {
			t.Logf("Time series analysis for sector %s:", sectorCode)
			for i := 0; i < 5 && i < len(shorts); i++ {
				short := shorts[i]
				totalTurnover := short.SellingExcludingShortSellingTurnoverValue +
					short.ShortSellingWithRestrictionsTurnoverValue +
					short.ShortSellingWithoutRestrictionsTurnoverValue
				
				if totalTurnover > 0 {
					shortSellingRatio := (short.ShortSellingWithRestrictionsTurnoverValue +
						short.ShortSellingWithoutRestrictionsTurnoverValue) / totalTurnover * 100
					
					t.Logf("  %s: %.2f%% short selling ratio (turnover: %.0f)", 
						short.Date, shortSellingRatio, totalTurnover)
				}
			}
		}
	})

	t.Run("GetShortSelling_SectorAnalysis", func(t *testing.T) {
		// 最新日の全業種の空売り比率分析
		date := getTestDate()
		
		shorts, err := jq.ShortSelling.GetShortSellingByDate(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No short selling data available for analysis")
		}

		if len(shorts) == 0 {
			t.Skip("No short selling data available")
		}

		// 業種別の空売り比率ランキング
		type SectorRatio struct {
			SectorCode string
			Ratio      float64
			Turnover   float64
		}
		
		var sectorRatios []SectorRatio
		
		for _, short := range shorts {
			totalTurnover := short.SellingExcludingShortSellingTurnoverValue +
				short.ShortSellingWithRestrictionsTurnoverValue +
				short.ShortSellingWithoutRestrictionsTurnoverValue
			
			if totalTurnover > 0 {
				shortSellingRatio := (short.ShortSellingWithRestrictionsTurnoverValue +
					short.ShortSellingWithoutRestrictionsTurnoverValue) / totalTurnover * 100
				
				sectorRatios = append(sectorRatios, SectorRatio{
					SectorCode: short.Sector33Code,
					Ratio:      shortSellingRatio,
					Turnover:   totalTurnover,
				})
			}
		}
		
		// 上位10業種の空売り比率
		t.Logf("Top 10 sectors by short selling ratio on %s:", date)
		count := 0
		for _, sector := range sectorRatios {
			if count >= 10 {
				break
			}
			t.Logf("  Sector %s: %.2f%% (turnover: %.0f million)",
				sector.SectorCode, sector.Ratio, sector.Turnover/1000000)
			count++
		}
		
		// 統計情報
		if len(sectorRatios) > 0 {
			var totalRatio, minRatio, maxRatio float64
			minRatio = sectorRatios[0].Ratio
			maxRatio = sectorRatios[0].Ratio
			
			for _, sector := range sectorRatios {
				totalRatio += sector.Ratio
				if sector.Ratio < minRatio {
					minRatio = sector.Ratio
				}
				if sector.Ratio > maxRatio {
					maxRatio = sector.Ratio
				}
			}
			
			avgRatio := totalRatio / float64(len(sectorRatios))
			t.Logf("Short selling ratio statistics:")
			t.Logf("  Average: %.2f%%", avgRatio)
			t.Logf("  Min: %.2f%%", minRatio)
			t.Logf("  Max: %.2f%%", maxRatio)
			t.Logf("  Sectors analyzed: %d", len(sectorRatios))
		}
	})

	t.Run("GetShortSelling_MultiSector", func(t *testing.T) {
		// 複数の主要業種の比較
		majorSectors := []string{"0050", "1050", "2050", "3050", "4050"} // 主要業種コード
		sectorData := make(map[string][]float64) // セクター -> 直近の空売り比率リスト
		
		for _, sectorCode := range majorSectors {
			shorts, err := jq.ShortSelling.GetShortSellingBySector(sectorCode)
			if err != nil {
				if isSubscriptionLimited(err) {
					t.Skip("Skipping due to subscription limitation")
				}
				continue
			}
			
			if len(shorts) == 0 {
				continue
			}
			
			// 直近5日分の空売り比率を計算
			ratios := make([]float64, 0, 5)
			for i := 0; i < 5 && i < len(shorts); i++ {
				short := shorts[i]
				totalTurnover := short.SellingExcludingShortSellingTurnoverValue +
					short.ShortSellingWithRestrictionsTurnoverValue +
					short.ShortSellingWithoutRestrictionsTurnoverValue
				
				if totalTurnover > 0 {
					ratio := (short.ShortSellingWithRestrictionsTurnoverValue +
						short.ShortSellingWithoutRestrictionsTurnoverValue) / totalTurnover * 100
					ratios = append(ratios, ratio)
				}
			}
			
			if len(ratios) > 0 {
				sectorData[sectorCode] = ratios
			}
		}
		
		// 業種間比較
		t.Logf("Major sectors comparison (recent 5 days avg):")
		for sectorCode, ratios := range sectorData {
			if len(ratios) > 0 {
				var sum float64
				for _, ratio := range ratios {
					sum += ratio
				}
				avg := sum / float64(len(ratios))
				t.Logf("  Sector %s: %.2f%% (based on %d days)", sectorCode, avg, len(ratios))
			}
		}
	})

	t.Run("GetShortSelling_ErrorHandling", func(t *testing.T) {
		// エラーケースのテスト
		
		// 存在しない業種コード（0000は無効）
		shorts, err := jq.ShortSelling.GetShortSellingBySector("0000")
		if err == nil && len(shorts) > 0 {
			t.Error("Expected error or empty result for invalid sector code")
		}
		
		// 未来の日付
		futureDate := "2030-01-01"
		shorts, err = jq.ShortSelling.GetShortSellingByDate(futureDate)
		if err == nil && len(shorts) > 0 {
			t.Error("Expected error or empty result for future date")
		}
	})
}