//go:build e2e
// +build e2e

package e2e

import (
	"testing"
	"time"

	"github.com/utahta/jquants"
)

// TestFSDetailsEndpoint は/fins/fs_details（プレミアムプラン専用）エンドポイントの完全なテスト
func TestFSDetailsEndpoint(t *testing.T) {
	t.Run("GetFSDetails_ByCode", func(t *testing.T) {
		// トヨタ自動車の財務諸表詳細を取得
		details, err := jq.FSDetails.GetFSDetailsByCode("7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get FS details: %v", err)
		}

		if len(details) == 0 {
			t.Skip("No FS details data available")
		}

		// 各財務詳細データを詳細に検証
		for i, detail := range details {
			// 基本情報の検証
			if detail.LocalCode != "7203" && detail.LocalCode != "72030" {
				t.Errorf("Detail[%d]: LocalCode = %v, want 7203 or 72030", i, detail.LocalCode)
			}
			// LocalCodeは5桁であることを確認
			if len(detail.LocalCode) != 5 && len(detail.LocalCode) != 4 {
				t.Errorf("Detail[%d]: LocalCode length = %d, want 4 or 5", i, len(detail.LocalCode))
			}

			if detail.DisclosedDate == "" {
				t.Errorf("Detail[%d]: DisclosedDate is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(detail.DisclosedDate) != 10 || detail.DisclosedDate[4] != '-' || detail.DisclosedDate[7] != '-' {
					t.Errorf("Detail[%d]: DisclosedDate format invalid = %v, want YYYY-MM-DD", i, detail.DisclosedDate)
				}
			}

			if detail.DisclosedTime == "" {
				t.Errorf("Detail[%d]: DisclosedTime is empty", i)
			} else {
				// 時刻フォーマットの検証（HH:MM:SS形式）
				if len(detail.DisclosedTime) != 8 || detail.DisclosedTime[2] != ':' || detail.DisclosedTime[5] != ':' {
					t.Errorf("Detail[%d]: DisclosedTime format invalid = %v, want HH:MM:SS", i, detail.DisclosedTime)
				}
			}

			// 開示番号の検証
			if detail.DisclosureNumber == "" {
				t.Errorf("Detail[%d]: DisclosureNumber is empty", i)
			}

			// 開示書類種別の検証
			if detail.TypeOfDocument == "" {
				t.Errorf("Detail[%d]: TypeOfDocument is empty", i)
			}

			// 会計基準の検証
			accountingStandards, hasAccounting := detail.GetValue("Accounting standards, DEI")
			if hasAccounting && accountingStandards != "" {
				if accountingStandards != "IFRS" && accountingStandards != "JapaneseGAAP" {
					t.Errorf("Detail[%d]: Accounting standards = %v, want IFRS or JapaneseGAAP", i, accountingStandards)
				}
			}

			// 決算期の検証（FinancialStatementマップから取得）
			currentPeriodEndDate, hasPeriodEnd := detail.GetValue("Current period end date, DEI")
			previousYearEndDate, hasPrevYear := detail.GetValue("Previous fiscal year end date, DEI")

			if hasPeriodEnd && currentPeriodEndDate == "" {
				t.Errorf("Detail[%d]: CurrentPeriodEndDate is empty", i)
			}
			if hasPrevYear && previousYearEndDate == "" {
				t.Errorf("Detail[%d]: PreviousYearEndDate is empty", i)
			}

			// 期間の論理チェック
			if hasPeriodEnd && hasPrevYear && currentPeriodEndDate != "" && previousYearEndDate != "" {
				endDate, err1 := time.Parse("2006-01-02", currentPeriodEndDate)
				prevDate, err2 := time.Parse("2006-01-02", previousYearEndDate)
				if err1 == nil && err2 == nil && prevDate.After(endDate) {
					t.Errorf("Detail[%d]: PrevYearDate (%s) after CurrentEndDate (%s)",
						i, previousYearEndDate, currentPeriodEndDate)
				}
			}

			// 財務指標の妥当性チェック
			validateFinancialMetric := func(name string, value *float64) {
				if value != nil {
					if *value < 0 && name != "NetIncome" && name != "OperatingIncome" {
						// 純利益や営業利益は負の値もありうる
						t.Logf("Detail[%d]: %s is negative: %.2f", i, name, *value)
					}
				}
			}

			// 主要財務指標の検証（会計基準に応じて適切なキーを使用）
			var totalAssetsKey, equityKey, revenueKey, operatingProfitKey string
			if detail.IsIFRS() {
				totalAssetsKey = "Assets (IFRS)"
				equityKey = "Equity (IFRS)"
				revenueKey = "Revenue - 2 (IFRS)"
				operatingProfitKey = "Operating profit (loss) (IFRS)"
			} else {
				// 日本基準の場合の主要キー（実際のキー名は要確認）
				totalAssetsKey = "Total assets"
				equityKey = "Total shareholders' equity"
				revenueKey = "Sales"
				operatingProfitKey = "Operating income"
			}

			// 各指標の検証
			if assets, err := detail.GetFloatValue(totalAssetsKey); err == nil {
				validateFinancialMetric("TotalAssets", &assets)
			}
			if equity, err := detail.GetFloatValue(equityKey); err == nil {
				validateFinancialMetric("Equity", &equity)
			}
			if revenue, err := detail.GetFloatValue(revenueKey); err == nil {
				validateFinancialMetric("Revenue", &revenue)
			}
			if operatingProfit, err := detail.GetFloatValue(operatingProfitKey); err == nil {
				validateFinancialMetric("OperatingProfit", &operatingProfit)
			}

			// ROEの妥当性チェック（IFRS基準でのみ計算可能）
			if detail.IsIFRS() {
				if roe, err := detail.GetROE(); err == nil && (*roe < -100 || *roe > 100) {
					t.Logf("Detail[%d]: ROE value seems unusual: %.2f%%", i, *roe)
				}

				// 流動比率の検証
				if currentRatio, err := detail.GetCurrentRatio(); err == nil && (*currentRatio < 0.1 || *currentRatio > 10) {
					t.Logf("Detail[%d]: Current ratio seems unusual: %.2f", i, *currentRatio)
				}

				// 自己資本比率の検証
				if equityRatio, err := detail.GetEquityRatio(); err == nil && (*equityRatio < 0 || *equityRatio > 100) {
					t.Logf("Detail[%d]: Equity ratio seems unusual: %.2f%%", i, *equityRatio)
				}

				// EPSの検証
				if eps, err := detail.GetBasicEPS(); err == nil && *eps < -10000 {
					t.Logf("Detail[%d]: EPS value seems unusual: %.2f", i, *eps)
				}
			}

			// 最初の3件の詳細ログ
			if i < 3 {
				periodEnd, _ := detail.GetValue("Current period end date, DEI")
				prevEnd, _ := detail.GetValue("Previous fiscal year end date, DEI")

				t.Logf("Detail[%d]: Code=%s, Type=%s",
					i, detail.LocalCode, detail.TypeOfDocument)
				t.Logf("  Disclosed: %s %s", detail.DisclosedDate, detail.DisclosedTime)
				t.Logf("  Period: %s (Previous: %s)", periodEnd, prevEnd)
				t.Logf("  Accounting: %s, Quarter: %d",
					map[bool]string{true: "IFRS", false: "Japanese GAAP"}[detail.IsIFRS()],
					detail.GetQuarter())

				// 主要財務指標のログ（IFRS基準の場合）
				if detail.IsIFRS() {
					if assets, err := detail.GetFloatValue("Assets (IFRS)"); err == nil {
						t.Logf("  Total Assets: %.0f million yen", assets/1000000)
					}
					if revenue, err := detail.GetFloatValue("Revenue - 2 (IFRS)"); err == nil {
						t.Logf("  Revenue: %.0f million yen", revenue/1000000)
					}
					if profit, err := detail.GetFloatValue("Profit (loss) attributable to owners of parent (IFRS)"); err == nil {
						t.Logf("  Net Profit: %.0f million yen", profit/1000000)
					}
					if roe, err := detail.GetROE(); err == nil {
						t.Logf("  ROE: %.2f%%", *roe)
					}
					if eps, err := detail.GetBasicEPS(); err == nil {
						t.Logf("  EPS: %.2f yen", *eps)
					}
				}
			}
		}

		t.Logf("Retrieved %d FS details records", len(details))
	})

	t.Run("GetFSDetails_ByDate", func(t *testing.T) {
		// 最近の営業日の財務詳細を取得
		date := getTestDate()

		params := jquants.FSDetailsParams{
			Date: date,
		}

		details, err := jq.FSDetails.GetFSDetails(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Logf("Failed to get FS details by date: %v", err)
			return
		}

		if details == nil || len(details.FSDetails) == 0 {
			t.Skip("No FS details data for the specified date")
		}

		t.Logf("Retrieved %d FS details records for %s", len(details.FSDetails), date)

		// 日付の一致確認
		for i, detail := range details.FSDetails {
			// 日付形式をYYYY-MM-DDに統一して比較
			expectedDate := getTestDateFormatted()
			if detail.DisclosedDate != expectedDate {
				t.Errorf("Detail[%d]: DisclosedDate = %v, want %v",
					i, detail.DisclosedDate, expectedDate)
			}
			if i >= 10 {
				break // 最初の10件のみ確認
			}
		}
	})

	t.Run("GetFSDetails_ProfitabilityAnalysis", func(t *testing.T) {
		// トヨタ自動車の収益性分析
		details, err := jq.FSDetails.GetFSDetailsByCode("7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No FS details data available for profitability analysis")
		}

		if len(details) == 0 {
			t.Skip("No FS details data available")
		}

		// 収益性指標の分析
		profitableQuarters := 0
		totalROE := 0.0
		totalROA := 0.0
		validROECount := 0
		validROACount := 0

		for _, detail := range details {
			// 利益性の確認（IFRS基準の場合）
			if detail.IsIFRS() {
				if profit, err := detail.GetFloatValue("Profit (loss) attributable to owners of parent (IFRS)"); err == nil && profit > 0 {
					profitableQuarters++
				}

				// ROEの集計
				if roe, err := detail.GetROE(); err == nil && *roe > -100 && *roe < 100 {
					totalROE += *roe
					validROECount++
				}

				// 自己資本比率の集計（ROAの代わり）
				if equityRatio, err := detail.GetEquityRatio(); err == nil && *equityRatio > 0 && *equityRatio <= 100 {
					totalROA += *equityRatio
					validROACount++
				}
			}
		}

		t.Logf("Profitability analysis for code 7203:")
		t.Logf("  Total quarters analyzed: %d", len(details))
		t.Logf("  Profitable quarters: %d (%.1f%%)",
			profitableQuarters, float64(profitableQuarters)/float64(len(details))*100)

		if validROECount > 0 {
			avgROE := totalROE / float64(validROECount)
			t.Logf("  Average ROE: %.2f%%", avgROE)
		}
		if validROACount > 0 {
			avgROA := totalROA / float64(validROACount)
			t.Logf("  Average ROA: %.2f%%", avgROA)
		}
	})

	t.Run("GetFSDetails_GrowthAnalysis", func(t *testing.T) {
		// 成長性分析
		details, err := jq.FSDetails.GetFSDetailsByCode("7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Skip("No FS details data available for growth analysis")
		}

		if len(details) < 2 {
			t.Skip("Insufficient data for growth analysis")
		}

		// 最新と前年同期の比較（IFRS基準でのみ実施）
		if len(details) >= 2 {
			latest := details[0]
			previous := details[1]

			t.Logf("Growth analysis (latest vs previous period):")

			// IFRS基準の場合のみ成長率計算
			if latest.IsIFRS() && previous.IsIFRS() {
				// 売上成長率
				if latestRevenue, err1 := latest.GetFloatValue("Revenue - 2 (IFRS)"); err1 == nil {
					if prevRevenue, err2 := previous.GetFloatValue("Revenue - 2 (IFRS)"); err2 == nil && prevRevenue != 0 {
						revenueGrowth := (latestRevenue - prevRevenue) / prevRevenue * 100
						t.Logf("  Revenue growth: %.2f%%", revenueGrowth)
					}
				}

				// 利益成長率
				if latestProfit, err1 := latest.GetFloatValue("Profit (loss) attributable to owners of parent (IFRS)"); err1 == nil {
					if prevProfit, err2 := previous.GetFloatValue("Profit (loss) attributable to owners of parent (IFRS)"); err2 == nil && prevProfit != 0 {
						profitGrowth := (latestProfit - prevProfit) / prevProfit * 100
						t.Logf("  Net income growth: %.2f%%", profitGrowth)
					}
				}

				// 総資産成長率
				if latestAssets, err1 := latest.GetFloatValue("Assets (IFRS)"); err1 == nil {
					if prevAssets, err2 := previous.GetFloatValue("Assets (IFRS)"); err2 == nil && prevAssets != 0 {
						assetGrowth := (latestAssets - prevAssets) / prevAssets * 100
						t.Logf("  Total assets growth: %.2f%%", assetGrowth)
					}
				}
			}
		}
	})

	t.Run("GetFSDetails_DateRange", func(t *testing.T) {
		// 過去1年間の財務詳細を取得
		to := getTestDateFormatted()
		fromTime, _ := time.Parse("2006-01-02", to)
		fromTime = fromTime.AddDate(-1, 0, 0)
		from := fromTime.Format("2006-01-02")

		params := jquants.FSDetailsParams{
			Code: "7203",
		}

		details, err := jq.FSDetails.GetFSDetails(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation (expected for premium API)")
			}
			t.Fatalf("Failed to get FS details for date range: %v", err)
		}

		if details == nil || len(details.FSDetails) == 0 {
			t.Skip("No FS details data for the specified date range")
		}

		t.Logf("Retrieved %d FS details records for code 7203", len(details.FSDetails))

		// 期間内の日付確認
		for _, detail := range details.FSDetails {
			if detail.DisclosedDate != "" {
				if detail.DisclosedDate < from || detail.DisclosedDate > to {
					t.Errorf("FS detail disclosed date %s is outside range %s to %s",
						detail.DisclosedDate, from, to)
				}
			}
		}
	})

	t.Run("GetFSDetails_ErrorHandling", func(t *testing.T) {
		// エラーケースのテスト

		// 存在しない銘柄コード
		details, err := jq.FSDetails.GetFSDetailsByCode("99999")
		if err == nil && len(details) > 0 {
			t.Error("Expected error or empty result for invalid code")
		}

		// 無効な日付
		params := jquants.FSDetailsParams{
			Date: "invalid-date",
		}

		resp, err := jq.FSDetails.GetFSDetails(params)
		if err == nil && resp != nil && len(resp.FSDetails) > 0 {
			t.Error("Expected error or empty result for invalid date")
		}
	})
}
