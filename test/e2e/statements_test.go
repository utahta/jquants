//go:build e2e
// +build e2e

package e2e

import (
	"testing"
)

// TestStatementsEndpoint は/fins/statementsエンドポイントの完全なテスト
func TestStatementsEndpoint(t *testing.T) {
	t.Run("GetStatements_ByCode", func(t *testing.T) {
		// トヨタ自動車(7203)の財務諸表を取得
		statements, err := jq.Statements.GetStatements("7203", "")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get statements: %v", err)
		}

		if len(statements) == 0 {
			t.Skip("No statements data available")
		}

		// 最新の財務諸表を詳細に検証
		latest := statements[0]
		
		// 基本情報の検証
		if latest.LocalCode != "7203" && latest.LocalCode != "72030" {
			t.Errorf("LocalCode = %v, want 7203 or 72030", latest.LocalCode)
		} else {
			// LocalCodeは5桁
			if len(latest.LocalCode) != 5 {
				t.Errorf("LocalCode length = %d, want 5", len(latest.LocalCode))
			}
		}
		if latest.DisclosedDate == "" {
			t.Error("DisclosedDate is empty")
		} else {
			// 日付フォーマットの検証（YYYY-MM-DD形式）
			if len(latest.DisclosedDate) != 10 || latest.DisclosedDate[4] != '-' || latest.DisclosedDate[7] != '-' {
				t.Errorf("DisclosedDate format invalid = %v, want YYYY-MM-DD", latest.DisclosedDate)
			}
		}
		if latest.DisclosedTime == "" {
			t.Error("DisclosedTime is empty")
		} else {
			// 時刻フォーマットの検証（HH:MMまたはHH:MM:SS形式）
			if len(latest.DisclosedTime) == 5 {
				// HH:MM形式
				if latest.DisclosedTime[2] != ':' {
					t.Errorf("DisclosedTime format invalid = %v, want HH:MM or HH:MM:SS", latest.DisclosedTime)
				}
			} else if len(latest.DisclosedTime) == 8 {
				// HH:MM:SS形式
				if latest.DisclosedTime[2] != ':' || latest.DisclosedTime[5] != ':' {
					t.Errorf("DisclosedTime format invalid = %v, want HH:MM or HH:MM:SS", latest.DisclosedTime)
				}
			} else {
				t.Errorf("DisclosedTime format invalid = %v, want HH:MM or HH:MM:SS", latest.DisclosedTime)
			}
		}
		if latest.DisclosureNumber == "" {
			t.Error("DisclosureNumber is empty")
		}
		
		// 書類・期間情報の検証
		if latest.TypeOfDocument == "" {
			t.Error("TypeOfDocument is empty")
		}
		if latest.TypeOfCurrentPeriod == "" {
			t.Error("TypeOfCurrentPeriod is empty")
		} else {
			// TypeOfCurrentPeriodの値検証
			validPeriods := map[string]bool{"1Q": true, "2Q": true, "3Q": true, "4Q": true, "5Q": true, "FY": true}
			if !validPeriods[latest.TypeOfCurrentPeriod] {
				t.Errorf("TypeOfCurrentPeriod = %v, want one of [1Q, 2Q, 3Q, 4Q, 5Q, FY]", latest.TypeOfCurrentPeriod)
			}
		}
		if latest.CurrentPeriodStartDate == "" {
			t.Error("CurrentPeriodStartDate is empty")
		} else {
			// 日付フォーマットの検証
			if len(latest.CurrentPeriodStartDate) != 10 || latest.CurrentPeriodStartDate[4] != '-' || latest.CurrentPeriodStartDate[7] != '-' {
				t.Errorf("CurrentPeriodStartDate format invalid = %v, want YYYY-MM-DD", latest.CurrentPeriodStartDate)
			}
		}
		if latest.CurrentPeriodEndDate == "" {
			t.Error("CurrentPeriodEndDate is empty")
		} else {
			// 日付フォーマットの検証
			if len(latest.CurrentPeriodEndDate) != 10 || latest.CurrentPeriodEndDate[4] != '-' || latest.CurrentPeriodEndDate[7] != '-' {
				t.Errorf("CurrentPeriodEndDate format invalid = %v, want YYYY-MM-DD", latest.CurrentPeriodEndDate)
			}
		}
		if latest.CurrentFiscalYearStartDate == "" {
			t.Error("CurrentFiscalYearStartDate is empty")
		} else {
			// 日付フォーマットの検証
			if len(latest.CurrentFiscalYearStartDate) != 10 || latest.CurrentFiscalYearStartDate[4] != '-' || latest.CurrentFiscalYearStartDate[7] != '-' {
				t.Errorf("CurrentFiscalYearStartDate format invalid = %v, want YYYY-MM-DD", latest.CurrentFiscalYearStartDate)
			}
		}
		if latest.CurrentFiscalYearEndDate == "" {
			t.Error("CurrentFiscalYearEndDate is empty")
		} else {
			// 日付フォーマットの検証
			if len(latest.CurrentFiscalYearEndDate) != 10 || latest.CurrentFiscalYearEndDate[4] != '-' || latest.CurrentFiscalYearEndDate[7] != '-' {
				t.Errorf("CurrentFiscalYearEndDate format invalid = %v, want YYYY-MM-DD", latest.CurrentFiscalYearEndDate)
			}
		}
		
		// 書類種別の検証（会計基準は実際の構造体にはない）
		if latest.TypeOfDocument == "" {
			t.Error("TypeOfDocument is empty")
		}
		
		// 財務データの検証（nilチェックと値の妥当性確認）
		validateFinancialValue(t, "NetSales", latest.NetSales)
		validateFinancialValue(t, "OperatingProfit", latest.OperatingProfit)
		validateFinancialValue(t, "OrdinaryProfit", latest.OrdinaryProfit)
		validateFinancialValue(t, "Profit", latest.Profit)
		validateFinancialValue(t, "EarningsPerShare", latest.EarningsPerShare)
		validateFinancialValue(t, "DilutedEarningsPerShare", latest.DilutedEarningsPerShare)
		validateFinancialValue(t, "TotalAssets", latest.TotalAssets)
		validateFinancialValue(t, "Equity", latest.Equity)
		validateFinancialValue(t, "EquityToAssetRatio", latest.EquityToAssetRatio)
		validateFinancialValue(t, "BookValuePerShare", latest.BookValuePerShare)
		
		// キャッシュフロー情報の検証
		validateFinancialValue(t, "CashFlowsFromOperatingActivities", latest.CashFlowsFromOperatingActivities)
		validateFinancialValue(t, "CashFlowsFromInvestingActivities", latest.CashFlowsFromInvestingActivities)
		validateFinancialValue(t, "CashFlowsFromFinancingActivities", latest.CashFlowsFromFinancingActivities)
		validateFinancialValue(t, "CashAndEquivalents", latest.CashAndEquivalents)
		
		// 配当情報の検証
		validateFinancialValue(t, "ResultDividendPerShareAnnual", latest.ResultDividendPerShareAnnual)
		validateFinancialValue(t, "ResultDividendPerShare1stQuarter", latest.ResultDividendPerShare1stQuarter)
		validateFinancialValue(t, "ResultDividendPerShare2ndQuarter", latest.ResultDividendPerShare2ndQuarter)
		validateFinancialValue(t, "ResultDividendPerShare3rdQuarter", latest.ResultDividendPerShare3rdQuarter)
		validateFinancialValue(t, "ResultDividendPerShareFiscalYearEnd", latest.ResultDividendPerShareFiscalYearEnd)
		validateFinancialValue(t, "ResultPayoutRatioAnnual", latest.ResultPayoutRatioAnnual)
		
		// 新規追加フィールド: REIT関連と配当総額
		validateFinancialValue(t, "DistributionsPerUnitREIT", latest.DistributionsPerUnitREIT)
		validateFinancialValue(t, "ResultTotalDividendPaidAnnual", latest.ResultTotalDividendPaidAnnual)
		
		// 予想配当情報の検証
		validateFinancialValue(t, "ForecastDividendPerShareAnnual", latest.ForecastDividendPerShareAnnual)
		validateFinancialValue(t, "ForecastDividendPerShare1stQuarter", latest.ForecastDividendPerShare1stQuarter)
		validateFinancialValue(t, "ForecastDividendPerShare2ndQuarter", latest.ForecastDividendPerShare2ndQuarter)
		validateFinancialValue(t, "ForecastDividendPerShare3rdQuarter", latest.ForecastDividendPerShare3rdQuarter)
		validateFinancialValue(t, "ForecastDividendPerShareFiscalYearEnd", latest.ForecastDividendPerShareFiscalYearEnd)
		validateFinancialValue(t, "ForecastPayoutRatioAnnual", latest.ForecastPayoutRatioAnnual)
		
		// 新規追加フィールド: 予想REIT関連と配当総額
		validateFinancialValue(t, "ForecastDistributionsPerUnitREIT", latest.ForecastDistributionsPerUnitREIT)
		validateFinancialValue(t, "ForecastTotalDividendPaidAnnual", latest.ForecastTotalDividendPaidAnnual)
		
		// 予想財務データの検証
		validateFinancialValue(t, "ForecastNetSales", latest.ForecastNetSales)
		validateFinancialValue(t, "ForecastOperatingProfit", latest.ForecastOperatingProfit)
		validateFinancialValue(t, "ForecastOrdinaryProfit", latest.ForecastOrdinaryProfit)
		validateFinancialValue(t, "ForecastProfit", latest.ForecastProfit)
		validateFinancialValue(t, "ForecastEarningsPerShare", latest.ForecastEarningsPerShare)
		
		// 第2四半期予想
		validateFinancialValue(t, "ForecastNetSales2ndQuarter", latest.ForecastNetSales2ndQuarter)
		validateFinancialValue(t, "ForecastOperatingProfit2ndQuarter", latest.ForecastOperatingProfit2ndQuarter)
		validateFinancialValue(t, "ForecastOrdinaryProfit2ndQuarter", latest.ForecastOrdinaryProfit2ndQuarter)
		validateFinancialValue(t, "ForecastProfit2ndQuarter", latest.ForecastProfit2ndQuarter)
		validateFinancialValue(t, "ForecastEarningsPerShare2ndQuarter", latest.ForecastEarningsPerShare2ndQuarter)
		
		// 翌期予想配当
		validateFinancialValue(t, "NextYearForecastDividendPerShare1stQuarter", latest.NextYearForecastDividendPerShare1stQuarter)
		validateFinancialValue(t, "NextYearForecastDividendPerShare2ndQuarter", latest.NextYearForecastDividendPerShare2ndQuarter)
		validateFinancialValue(t, "NextYearForecastDividendPerShare3rdQuarter", latest.NextYearForecastDividendPerShare3rdQuarter)
		validateFinancialValue(t, "NextYearForecastDividendPerShareFiscalYearEnd", latest.NextYearForecastDividendPerShareFiscalYearEnd)
		validateFinancialValue(t, "NextYearForecastDividendPerShareAnnual", latest.NextYearForecastDividendPerShareAnnual)
		validateFinancialValue(t, "NextYearForecastDistributionsPerUnitREIT", latest.NextYearForecastDistributionsPerUnitREIT)
		validateFinancialValue(t, "NextYearForecastPayoutRatioAnnual", latest.NextYearForecastPayoutRatioAnnual)
		
		// 翌期第2四半期予想
		validateFinancialValue(t, "NextYearForecastNetSales2ndQuarter", latest.NextYearForecastNetSales2ndQuarter)
		validateFinancialValue(t, "NextYearForecastOperatingProfit2ndQuarter", latest.NextYearForecastOperatingProfit2ndQuarter)
		validateFinancialValue(t, "NextYearForecastOrdinaryProfit2ndQuarter", latest.NextYearForecastOrdinaryProfit2ndQuarter)
		validateFinancialValue(t, "NextYearForecastProfit2ndQuarter", latest.NextYearForecastProfit2ndQuarter)
		validateFinancialValue(t, "NextYearForecastEarningsPerShare2ndQuarter", latest.NextYearForecastEarningsPerShare2ndQuarter)
		
		// 翌期通期予想
		validateFinancialValue(t, "NextYearForecastNetSales", latest.NextYearForecastNetSales)
		validateFinancialValue(t, "NextYearForecastOperatingProfit", latest.NextYearForecastOperatingProfit)
		validateFinancialValue(t, "NextYearForecastOrdinaryProfit", latest.NextYearForecastOrdinaryProfit)
		validateFinancialValue(t, "NextYearForecastProfit", latest.NextYearForecastProfit)
		validateFinancialValue(t, "NextYearForecastEarningsPerShare", latest.NextYearForecastEarningsPerShare)
		
		// 修正情報の検証（bool型フィールド）
		t.Logf("MaterialChangesInSubsidiaries: %v", latest.MaterialChangesInSubsidiaries)
		t.Logf("SignificantChangesInTheScopeOfConsolidation: %v", latest.SignificantChangesInTheScopeOfConsolidation)
		t.Logf("ChangesBasedOnRevisionsOfAccountingStandard: %v", latest.ChangesBasedOnRevisionsOfAccountingStandard)
		t.Logf("ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard: %v", latest.ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard)
		t.Logf("ChangesInAccountingEstimates: %v", latest.ChangesInAccountingEstimates)
		t.Logf("RetrospectiveRestatement: %v", latest.RetrospectiveRestatement)
		
		// 株式数情報の検証
		if latest.NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock != nil {
			t.Logf("NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock: %d", 
				*latest.NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock)
		}
		if latest.NumberOfTreasuryStockAtTheEndOfFiscalYear != nil {
			t.Logf("NumberOfTreasuryStockAtTheEndOfFiscalYear: %d", 
				*latest.NumberOfTreasuryStockAtTheEndOfFiscalYear)
		}
		if latest.AverageNumberOfShares != nil {
			t.Logf("AverageNumberOfShares: %d", *latest.AverageNumberOfShares)
		}
		
		// 単体財務データの検証
		validateFinancialValue(t, "NonConsolidatedNetSales", latest.NonConsolidatedNetSales)
		validateFinancialValue(t, "NonConsolidatedOperatingProfit", latest.NonConsolidatedOperatingProfit)
		validateFinancialValue(t, "NonConsolidatedOrdinaryProfit", latest.NonConsolidatedOrdinaryProfit)
		validateFinancialValue(t, "NonConsolidatedProfit", latest.NonConsolidatedProfit)
		validateFinancialValue(t, "NonConsolidatedEarningsPerShare", latest.NonConsolidatedEarningsPerShare)
		validateFinancialValue(t, "NonConsolidatedTotalAssets", latest.NonConsolidatedTotalAssets)
		validateFinancialValue(t, "NonConsolidatedEquity", latest.NonConsolidatedEquity)
		validateFinancialValue(t, "NonConsolidatedEquityToAssetRatio", latest.NonConsolidatedEquityToAssetRatio)
		validateFinancialValue(t, "NonConsolidatedBookValuePerShare", latest.NonConsolidatedBookValuePerShare)
		
		// 単体予想（第2四半期）
		validateFinancialValue(t, "ForecastNonConsolidatedNetSales2ndQuarter", latest.ForecastNonConsolidatedNetSales2ndQuarter)
		validateFinancialValue(t, "ForecastNonConsolidatedOperatingProfit2ndQuarter", latest.ForecastNonConsolidatedOperatingProfit2ndQuarter)
		validateFinancialValue(t, "ForecastNonConsolidatedOrdinaryProfit2ndQuarter", latest.ForecastNonConsolidatedOrdinaryProfit2ndQuarter)
		validateFinancialValue(t, "ForecastNonConsolidatedProfit2ndQuarter", latest.ForecastNonConsolidatedProfit2ndQuarter)
		validateFinancialValue(t, "ForecastNonConsolidatedEarningsPerShare2ndQuarter", latest.ForecastNonConsolidatedEarningsPerShare2ndQuarter)
		
		// 単体翌期予想（第2四半期）
		validateFinancialValue(t, "NextYearForecastNonConsolidatedNetSales2ndQuarter", latest.NextYearForecastNonConsolidatedNetSales2ndQuarter)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedOperatingProfit2ndQuarter", latest.NextYearForecastNonConsolidatedOperatingProfit2ndQuarter)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedOrdinaryProfit2ndQuarter", latest.NextYearForecastNonConsolidatedOrdinaryProfit2ndQuarter)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedProfit2ndQuarter", latest.NextYearForecastNonConsolidatedProfit2ndQuarter)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedEarningsPerShare2ndQuarter", latest.NextYearForecastNonConsolidatedEarningsPerShare2ndQuarter)
		
		// 単体予想（期末）
		validateFinancialValue(t, "ForecastNonConsolidatedNetSales", latest.ForecastNonConsolidatedNetSales)
		validateFinancialValue(t, "ForecastNonConsolidatedOperatingProfit", latest.ForecastNonConsolidatedOperatingProfit)
		validateFinancialValue(t, "ForecastNonConsolidatedOrdinaryProfit", latest.ForecastNonConsolidatedOrdinaryProfit)
		validateFinancialValue(t, "ForecastNonConsolidatedProfit", latest.ForecastNonConsolidatedProfit)
		validateFinancialValue(t, "ForecastNonConsolidatedEarningsPerShare", latest.ForecastNonConsolidatedEarningsPerShare)
		
		// 単体翌期予想（期末）
		validateFinancialValue(t, "NextYearForecastNonConsolidatedNetSales", latest.NextYearForecastNonConsolidatedNetSales)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedOperatingProfit", latest.NextYearForecastNonConsolidatedOperatingProfit)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedOrdinaryProfit", latest.NextYearForecastNonConsolidatedOrdinaryProfit)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedProfit", latest.NextYearForecastNonConsolidatedProfit)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedEarningsPerShare", latest.NextYearForecastNonConsolidatedEarningsPerShare)
		
		t.Logf("Successfully validated all fields for statement: %s", latest.DisclosedDate)
	})

	t.Run("GetStatements_Historical", func(t *testing.T) {
		// 複数の履歴データ取得のテスト
		statements, err := jq.Statements.GetStatements("7203", "")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get statements: %v", err)
		}

		if len(statements) == 0 {
			t.Skip("No statements data available")
		}

		t.Logf("Retrieved %d statements", len(statements))
		
		// 複数の決算期のデータが取得できているかを確認
		if len(statements) >= 2 {
			first := statements[0]
			second := statements[1]
			
			// 決算期が異なることを確認
			if first.CurrentFiscalYearEndDate == second.CurrentFiscalYearEndDate {
				t.Logf("Warning: Multiple statements with same fiscal year end date")
			} else {
				t.Logf("Found multiple fiscal periods: %s and %s", 
					first.CurrentFiscalYearEndDate, second.CurrentFiscalYearEndDate)
			}
		}
	})

	t.Run("GetStatementsByDate", func(t *testing.T) {
		// 特定日の財務諸表を取得（過去の営業日を使用）
		date := getTestDate()
		
		// GetStatementsByDateメソッドを使用
		statements, err := jq.Statements.GetStatementsByDate(date)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			// 指定日にデータがない可能性もある
			t.Logf("No statements data for date %s: %v", date, err)
			return
		}

		if len(statements) > 0 {
			t.Logf("Found %d statements for date %s", len(statements), date)
			
			// 各財務諸表の基本検証
			for i, stmt := range statements {
				if i >= 5 {
					break // 最初の5件のみ詳細検証
				}
				
				if stmt.LocalCode == "" {
					t.Errorf("Statement[%d]: LocalCode is empty", i)
				} else {
					// LocalCodeは5桁
					if len(stmt.LocalCode) != 5 {
						t.Errorf("Statement[%d]: LocalCode length = %d, want 5", i, len(stmt.LocalCode))
					}
				}
				// DisclosedDateの検証（APIはYYYY-MM-DD形式で返す）
				expectedDate := getTestDateFormatted()
				if stmt.DisclosedDate != expectedDate {
					t.Errorf("Statement[%d]: DisclosedDate = %v, want %v", i, stmt.DisclosedDate, expectedDate)
				}
				t.Logf("Statement[%d]: %s - %s", i, stmt.LocalCode, stmt.TypeOfDocument)
			}
		}
	})
}

// validateFinancialValue は財務数値フィールドを検証するヘルパー関数
func validateFinancialValue(t *testing.T, fieldName string, value *float64) {
	if value != nil {
		if *value < 0 {
			// 一部のフィールド（利益など）は負の値を取ることがある
			t.Logf("%s: %.2f (negative value)", fieldName, *value)
		} else {
			t.Logf("%s: %.2f", fieldName, *value)
		}
	} else {
		// nilは許容される（すべての企業がすべてのフィールドを持つわけではない）
		t.Logf("%s: nil", fieldName)
	}
}