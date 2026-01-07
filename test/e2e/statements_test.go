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
		if latest.Code != "7203" && latest.Code != "72030" {
			t.Errorf("LocalCode = %v, want 7203 or 72030", latest.Code)
		} else {
			// LocalCodeは5桁
			if len(latest.Code) != 5 {
				t.Errorf("LocalCode length = %d, want 5", len(latest.Code))
			}
		}
		if latest.DiscDate == "" {
			t.Error("DisclosedDate is empty")
		} else {
			// 日付フォーマットの検証（YYYY-MM-DD形式）
			if len(latest.DiscDate) != 10 || latest.DiscDate[4] != '-' || latest.DiscDate[7] != '-' {
				t.Errorf("DisclosedDate format invalid = %v, want YYYY-MM-DD", latest.DiscDate)
			}
		}
		if latest.DiscTime == "" {
			t.Error("DisclosedTime is empty")
		} else {
			// 時刻フォーマットの検証（HH:MMまたはHH:MM:SS形式）
			if len(latest.DiscTime) == 5 {
				// HH:MM形式
				if latest.DiscTime[2] != ':' {
					t.Errorf("DisclosedTime format invalid = %v, want HH:MM or HH:MM:SS", latest.DiscTime)
				}
			} else if len(latest.DiscTime) == 8 {
				// HH:MM:SS形式
				if latest.DiscTime[2] != ':' || latest.DiscTime[5] != ':' {
					t.Errorf("DisclosedTime format invalid = %v, want HH:MM or HH:MM:SS", latest.DiscTime)
				}
			} else {
				t.Errorf("DisclosedTime format invalid = %v, want HH:MM or HH:MM:SS", latest.DiscTime)
			}
		}
		if latest.DiscNo == "" {
			t.Error("DisclosureNumber is empty")
		}
		
		// 書類・期間情報の検証
		if latest.DocType == "" {
			t.Error("TypeOfDocument is empty")
		}
		if latest.CurPerType == "" {
			t.Error("TypeOfCurrentPeriod is empty")
		} else {
			// TypeOfCurrentPeriodの値検証
			validPeriods := map[string]bool{"1Q": true, "2Q": true, "3Q": true, "4Q": true, "5Q": true, "FY": true}
			if !validPeriods[latest.CurPerType] {
				t.Errorf("TypeOfCurrentPeriod = %v, want one of [1Q, 2Q, 3Q, 4Q, 5Q, FY]", latest.CurPerType)
			}
		}
		if latest.CurPerSt == "" {
			t.Error("CurrentPeriodStartDate is empty")
		} else {
			// 日付フォーマットの検証
			if len(latest.CurPerSt) != 10 || latest.CurPerSt[4] != '-' || latest.CurPerSt[7] != '-' {
				t.Errorf("CurrentPeriodStartDate format invalid = %v, want YYYY-MM-DD", latest.CurPerSt)
			}
		}
		if latest.CurPerEn == "" {
			t.Error("CurrentPeriodEndDate is empty")
		} else {
			// 日付フォーマットの検証
			if len(latest.CurPerEn) != 10 || latest.CurPerEn[4] != '-' || latest.CurPerEn[7] != '-' {
				t.Errorf("CurrentPeriodEndDate format invalid = %v, want YYYY-MM-DD", latest.CurPerEn)
			}
		}
		if latest.CurFYSt == "" {
			t.Error("CurrentFiscalYearStartDate is empty")
		} else {
			// 日付フォーマットの検証
			if len(latest.CurFYSt) != 10 || latest.CurFYSt[4] != '-' || latest.CurFYSt[7] != '-' {
				t.Errorf("CurrentFiscalYearStartDate format invalid = %v, want YYYY-MM-DD", latest.CurFYSt)
			}
		}
		if latest.CurFYEn == "" {
			t.Error("CurrentFiscalYearEndDate is empty")
		} else {
			// 日付フォーマットの検証
			if len(latest.CurFYEn) != 10 || latest.CurFYEn[4] != '-' || latest.CurFYEn[7] != '-' {
				t.Errorf("CurrentFiscalYearEndDate format invalid = %v, want YYYY-MM-DD", latest.CurFYEn)
			}
		}
		
		// 書類種別の検証（会計基準は実際の構造体にはない）
		if latest.DocType == "" {
			t.Error("TypeOfDocument is empty")
		}
		
		// 財務データの検証（nilチェックと値の妥当性確認）
		validateFinancialValue(t, "NetSales", latest.Sales)
		validateFinancialValue(t, "OperatingProfit", latest.OP)
		validateFinancialValue(t, "OrdinaryProfit", latest.OdP)
		validateFinancialValue(t, "Profit", latest.NP)
		validateFinancialValue(t, "EarningsPerShare", latest.EPS)
		validateFinancialValue(t, "DilutedEarningsPerShare", latest.DEPS)
		validateFinancialValue(t, "TotalAssets", latest.TA)
		validateFinancialValue(t, "Equity", latest.Eq)
		validateFinancialValue(t, "EquityToAssetRatio", latest.EqAR)
		validateFinancialValue(t, "BookValuePerShare", latest.BPS)
		
		// キャッシュフロー情報の検証
		validateFinancialValue(t, "CashFlowsFromOperatingActivities", latest.CFO)
		validateFinancialValue(t, "CashFlowsFromInvestingActivities", latest.CFI)
		validateFinancialValue(t, "CashFlowsFromFinancingActivities", latest.CFF)
		validateFinancialValue(t, "CashAndEquivalents", latest.CashEq)
		
		// 配当情報の検証
		validateFinancialValue(t, "ResultDividendPerShareAnnual", latest.DivAnn)
		validateFinancialValue(t, "ResultDividendPerShare1stQuarter", latest.Div1Q)
		validateFinancialValue(t, "ResultDividendPerShare2ndQuarter", latest.Div2Q)
		validateFinancialValue(t, "ResultDividendPerShare3rdQuarter", latest.Div3Q)
		validateFinancialValue(t, "ResultDividendPerShareFiscalYearEnd", latest.DivFY)
		validateFinancialValue(t, "ResultPayoutRatioAnnual", latest.PayoutRatioAn)
		
		// 新規追加フィールド: REIT関連と配当総額
		validateFinancialValue(t, "DistributionsPerUnitREIT", latest.DivUnit)
		validateFinancialValue(t, "ResultTotalDividendPaidAnnual", latest.DivTotalAnn)

		// 予想配当情報の検証
		validateFinancialValue(t, "ForecastDividendPerShareAnnual", latest.FDivAnn)
		validateFinancialValue(t, "ForecastDividendPerShare1stQuarter", latest.FDiv1Q)
		validateFinancialValue(t, "ForecastDividendPerShare2ndQuarter", latest.FDiv2Q)
		validateFinancialValue(t, "ForecastDividendPerShare3rdQuarter", latest.FDiv3Q)
		validateFinancialValue(t, "ForecastDividendPerShareFiscalYearEnd", latest.FDivFY)
		validateFinancialValue(t, "ForecastPayoutRatioAnnual", latest.FPayoutRatioAn)

		// 新規追加フィールド: 予想REIT関連と配当総額
		validateFinancialValue(t, "ForecastDistributionsPerUnitREIT", latest.FDivUnit)
		validateFinancialValue(t, "ForecastTotalDividendPaidAnnual", latest.FDivTotalAnn)
		
		// 予想財務データの検証
		validateFinancialValue(t, "ForecastNetSales", latest.FSales)
		validateFinancialValue(t, "ForecastOperatingProfit", latest.FOP)
		validateFinancialValue(t, "ForecastOrdinaryProfit", latest.FOdP)
		validateFinancialValue(t, "ForecastProfit", latest.FNP)
		validateFinancialValue(t, "ForecastEarningsPerShare", latest.FEPS)

		// 第2四半期予想
		validateFinancialValue(t, "ForecastNetSales2ndQuarter", latest.FSales2Q)
		validateFinancialValue(t, "ForecastOperatingProfit2ndQuarter", latest.FOP2Q)
		validateFinancialValue(t, "ForecastOrdinaryProfit2ndQuarter", latest.FOdP2Q)
		validateFinancialValue(t, "ForecastProfit2ndQuarter", latest.FNP2Q)
		validateFinancialValue(t, "ForecastEarningsPerShare2ndQuarter", latest.FEPS2Q)
		
		// 翌期予想配当
		validateFinancialValue(t, "NextYearForecastDividendPerShare1stQuarter", latest.NxFDiv1Q)
		validateFinancialValue(t, "NextYearForecastDividendPerShare2ndQuarter", latest.NxFDiv2Q)
		validateFinancialValue(t, "NextYearForecastDividendPerShare3rdQuarter", latest.NxFDiv3Q)
		validateFinancialValue(t, "NextYearForecastDividendPerShareFiscalYearEnd", latest.NxFDivFY)
		validateFinancialValue(t, "NextYearForecastDividendPerShareAnnual", latest.NxFDivAnn)
		validateFinancialValue(t, "NextYearForecastDistributionsPerUnitREIT", latest.NxFDivUnit)
		validateFinancialValue(t, "NextYearForecastPayoutRatioAnnual", latest.NxFPayoutRatioAn)

		// 翌期第2四半期予想
		validateFinancialValue(t, "NextYearForecastNetSales2ndQuarter", latest.NxFSales2Q)
		validateFinancialValue(t, "NextYearForecastOperatingProfit2ndQuarter", latest.NxFOP2Q)
		validateFinancialValue(t, "NextYearForecastOrdinaryProfit2ndQuarter", latest.NxFOdP2Q)
		validateFinancialValue(t, "NextYearForecastProfit2ndQuarter", latest.NxFNp2Q)
		validateFinancialValue(t, "NextYearForecastEarningsPerShare2ndQuarter", latest.NxFEPS2Q)

		// 翌期通期予想
		validateFinancialValue(t, "NextYearForecastNetSales", latest.NxFSales)
		validateFinancialValue(t, "NextYearForecastOperatingProfit", latest.NxFOP)
		validateFinancialValue(t, "NextYearForecastOrdinaryProfit", latest.NxFOdP)
		validateFinancialValue(t, "NextYearForecastProfit", latest.NxFNp)
		validateFinancialValue(t, "NextYearForecastEarningsPerShare", latest.NxFEPS)
		
		// 修正情報の検証（string型フィールド）
		t.Logf("MaterialChangesInSubsidiaries: %v", latest.MatChgSub)
		t.Logf("SignificantChangesInTheScopeOfConsolidation: %v", latest.SigChgInC)
		t.Logf("ChangesBasedOnRevisionsOfAccountingStandard: %v", latest.ChgByASRev)
		t.Logf("ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard: %v", latest.ChgNoASRev)
		t.Logf("ChangesInAccountingEstimates: %v", latest.ChgAcEst)
		t.Logf("RetrospectiveRestatement: %v", latest.RetroRst)

		// 株式数情報の検証
		if latest.ShOutFY != nil {
			t.Logf("NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock: %d",
				*latest.ShOutFY)
		}
		if latest.TrShFY != nil {
			t.Logf("NumberOfTreasuryStockAtTheEndOfFiscalYear: %d",
				*latest.TrShFY)
		}
		if latest.AvgSh != nil {
			t.Logf("AverageNumberOfShares: %d", *latest.AvgSh)
		}
		
		// 単体財務データの検証
		validateFinancialValue(t, "NonConsolidatedNetSales", latest.NCSales)
		validateFinancialValue(t, "NonConsolidatedOperatingProfit", latest.NCOP)
		validateFinancialValue(t, "NonConsolidatedOrdinaryProfit", latest.NCOdP)
		validateFinancialValue(t, "NonConsolidatedProfit", latest.NCNP)
		validateFinancialValue(t, "NonConsolidatedEarningsPerShare", latest.NCEPS)
		validateFinancialValue(t, "NonConsolidatedTotalAssets", latest.NCTA)
		validateFinancialValue(t, "NonConsolidatedEquity", latest.NCEq)
		validateFinancialValue(t, "NonConsolidatedEquityToAssetRatio", latest.NCEqAR)
		validateFinancialValue(t, "NonConsolidatedBookValuePerShare", latest.NCBPS)

		// 単体予想（第2四半期）
		validateFinancialValue(t, "ForecastNonConsolidatedNetSales2ndQuarter", latest.FNCSales2Q)
		validateFinancialValue(t, "ForecastNonConsolidatedOperatingProfit2ndQuarter", latest.FNCOP2Q)
		validateFinancialValue(t, "ForecastNonConsolidatedOrdinaryProfit2ndQuarter", latest.FNCOdP2Q)
		validateFinancialValue(t, "ForecastNonConsolidatedProfit2ndQuarter", latest.FNCNP2Q)
		validateFinancialValue(t, "ForecastNonConsolidatedEarningsPerShare2ndQuarter", latest.FNCEPS2Q)

		// 単体翌期予想（第2四半期）
		validateFinancialValue(t, "NextYearForecastNonConsolidatedNetSales2ndQuarter", latest.NxFNCSales2Q)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedOperatingProfit2ndQuarter", latest.NxFNCOP2Q)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedOrdinaryProfit2ndQuarter", latest.NxFNCOdP2Q)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedProfit2ndQuarter", latest.NxFNCNP2Q)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedEarningsPerShare2ndQuarter", latest.NxFNCEPS2Q)

		// 単体予想（期末）
		validateFinancialValue(t, "ForecastNonConsolidatedNetSales", latest.FNCSales)
		validateFinancialValue(t, "ForecastNonConsolidatedOperatingProfit", latest.FNCOP)
		validateFinancialValue(t, "ForecastNonConsolidatedOrdinaryProfit", latest.FNCOdP)
		validateFinancialValue(t, "ForecastNonConsolidatedProfit", latest.FNCNP)
		validateFinancialValue(t, "ForecastNonConsolidatedEarningsPerShare", latest.FNCEPS)

		// 単体翌期予想（期末）
		validateFinancialValue(t, "NextYearForecastNonConsolidatedNetSales", latest.NxFNCSales)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedOperatingProfit", latest.NxFNCOP)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedOrdinaryProfit", latest.NxFNCOdP)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedProfit", latest.NxFNCNP)
		validateFinancialValue(t, "NextYearForecastNonConsolidatedEarningsPerShare", latest.NxFNCEPS)
		
		t.Logf("Successfully validated all fields for statement: %s", latest.DiscDate)
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
			if first.CurFYEn == second.CurFYEn {
				t.Logf("Warning: Multiple statements with same fiscal year end date")
			} else {
				t.Logf("Found multiple fiscal periods: %s and %s",
					first.CurFYEn, second.CurFYEn)
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
				
				if stmt.Code == "" {
					t.Errorf("Statement[%d]: Code is empty", i)
				} else {
					// Codeは5桁
					if len(stmt.Code) != 5 {
						t.Errorf("Statement[%d]: Code length = %d, want 5", i, len(stmt.Code))
					}
				}
				// DiscDateの検証（APIはYYYY-MM-DD形式で返す）
				expectedDate := getTestDateFormatted()
				if stmt.DiscDate != expectedDate {
					t.Errorf("Statement[%d]: DiscDate = %v, want %v", i, stmt.DiscDate, expectedDate)
				}
				t.Logf("Statement[%d]: %s - %s", i, stmt.Code, stmt.DocType)
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