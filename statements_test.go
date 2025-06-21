package jquants

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestStatementsService_GetStatements(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		date     string
		wantPath string
	}{
		{
			name:     "with code and date",
			code:     "7203",
			date:     "20240101",
			wantPath: "/fins/statements?code=7203&date=20240101",
		},
		{
			name:     "with code only",
			code:     "7203",
			date:     "",
			wantPath: "/fins/statements?code=7203",
		},
		{
			name:     "with no parameters",
			code:     "",
			date:     "",
			wantPath: "/fins/statements",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewStatementsService(mockClient)

			// Mock response
			mockResponse := StatementsResponse{
				Statements: []Statement{
					{
						// 基本情報
						DisclosedDate:              "2024-01-15",
						DisclosedTime:              "14:30:00",
						LocalCode:                  "72030",
						DisclosureNumber:           "20240115123456",
						TypeOfDocument:             "3QFinancialStatements_Consolidated_IFRS",
						TypeOfCurrentPeriod:        "3Q",
						CurrentPeriodStartDate:     "2023-10-01",
						CurrentPeriodEndDate:       "2023-12-31",
						CurrentFiscalYearStartDate: "2023-04-01",
						CurrentFiscalYearEndDate:   "2024-03-31",
						// 連結財務数値
						NetSales:                         floatPtr(10000000000),
						OperatingProfit:                  floatPtr(2000000000),
						OrdinaryProfit:                   floatPtr(2100000000),
						Profit:                           floatPtr(1500000000),
						TotalAssets:                      floatPtr(50000000000),
						Equity:                           floatPtr(30000000000),
						CashAndEquivalents:               floatPtr(5000000000),
						CashFlowsFromOperatingActivities: floatPtr(3000000000),
						CashFlowsFromInvestingActivities: floatPtr(-1500000000),
						CashFlowsFromFinancingActivities: floatPtr(-500000000),
						// 財務指標
						EarningsPerShare:        floatPtr(150.5),
						DilutedEarningsPerShare: floatPtr(149.8),
						BookValuePerShare:       floatPtr(3000.0),
						EquityToAssetRatio:      floatPtr(60.0),
						// 配当情報
						ResultDividendPerShareAnnual:        floatPtr(60.0),
						ResultPayoutRatioAnnual:             floatPtr(39.9),
						ForecastDividendPerShareAnnual:      floatPtr(65.0),
						ForecastPayoutRatioAnnual:           floatPtr(40.0),
						ResultDividendPerShare2ndQuarter:    floatPtr(25.0),
						ResultDividendPerShareFiscalYearEnd: floatPtr(35.0),
						// 業績予想
						ForecastNetSales:         floatPtr(15000000000),
						ForecastOperatingProfit:  floatPtr(3000000000),
						ForecastOrdinaryProfit:   floatPtr(3100000000),
						ForecastProfit:           floatPtr(2200000000),
						ForecastEarningsPerShare: floatPtr(220.0),
					},
				},
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Test
			statements, err := service.GetStatements(tt.code, tt.date)
			if err != nil {
				t.Errorf("GetStatements failed: %v", err)
			}

			// Verify
			if len(statements) != 1 {
				t.Errorf("Expected 1 statement, got %d", len(statements))
			}

			if mockClient.LastMethod != "GET" {
				t.Errorf("Expected GET method, got %s", mockClient.LastMethod)
			}

			if mockClient.LastPath != tt.wantPath {
				t.Errorf("Expected path %s, got %s", tt.wantPath, mockClient.LastPath)
			}
		})
	}
}

func TestStatementsService_GetLatestStatements(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewStatementsService(mockClient)

	// Mock response
	mockResponse := StatementsResponse{
		Statements: []Statement{
			{
				DisclosedDate:                  "2024-01-15",
				DisclosedTime:                  "14:30:00",
				LocalCode:                      "72030",
				DisclosureNumber:               "20240115123456",
				TypeOfDocument:                 "3QFinancialStatements_Consolidated_IFRS",
				TypeOfCurrentPeriod:            "3Q",
				CurrentFiscalYearEndDate:       "2024-03-31",
				NetSales:                       floatPtr(10000000000),
				OperatingProfit:                floatPtr(2000000000),
				OrdinaryProfit:                 floatPtr(2100000000),
				Profit:                         floatPtr(1500000000),
				ResultDividendPerShareAnnual:   floatPtr(50.0),
				ForecastDividendPerShareAnnual: floatPtr(55.0),
				EarningsPerShare:               floatPtr(150.5),
				DilutedEarningsPerShare:        floatPtr(149.8),
			},
		},
	}
	mockClient.SetResponse("GET", "/fins/statements?code=7203", mockResponse)

	// Test
	statement, err := service.GetLatestStatements("7203")
	if err != nil {
		t.Errorf("GetLatestStatements failed: %v", err)
	}

	// Verify
	if statement == nil {
		t.Errorf("Expected statement, got nil")
		return
	}

	if statement.NetSales == nil || *statement.NetSales != 10000000000 {
		t.Errorf("Expected NetSales 10000000000, got %v", statement.NetSales)
	}

	if statement.EarningsPerShare == nil || *statement.EarningsPerShare != 150.5 {
		t.Errorf("Expected EarningsPerShare 150.5, got %v", statement.EarningsPerShare)
	}
}

func TestStatementsService_GetLatestStatements_NotFound(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewStatementsService(mockClient)

	// Mock empty response
	mockResponse := StatementsResponse{
		Statements: []Statement{},
	}
	mockClient.SetResponse("GET", "/fins/statements?code=9999", mockResponse)

	// Test
	_, err := service.GetLatestStatements("9999")
	if err == nil {
		t.Errorf("Expected error for non-existent company, got nil")
	}
}

func TestStatementsService_GetStatements_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewStatementsService(mockClient)

	// Mock error
	mockClient.SetError("GET", "/fins/statements?code=7203", fmt.Errorf("API error"))

	// Test
	_, err := service.GetStatements("7203", "")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestStatementsResponse_UnmarshalJSON(t *testing.T) {
	// Test JSON with mixed types (strings that should be converted to float64/int64/bool)
	jsonData := `{
		"statements": [
			{
				"DisclosedDate": "2024-01-15",
				"DisclosedTime": "14:30:00",
				"LocalCode": "72030",
				"DisclosureNumber": "20240115123456",
				"TypeOfDocument": "3QFinancialStatements_Consolidated_IFRS",
				"TypeOfCurrentPeriod": "3Q",
				"CurrentPeriodStartDate": "2023-10-01",
				"CurrentPeriodEndDate": "2023-12-31",
				"CurrentFiscalYearStartDate": "2023-04-01",
				"CurrentFiscalYearEndDate": "2024-03-31",
				"NetSales": "10000000000",
				"OperatingProfit": "2000000000",
				"OrdinaryProfit": "2100000000",
				"Profit": "1500000000",
				"EarningsPerShare": "150.5",
				"DilutedEarningsPerShare": "149.8",
				"TotalAssets": "50000000000",
				"Equity": "30000000000",
				"EquityToAssetRatio": "60.0",
				"BookValuePerShare": "3000.0",
				"CashAndEquivalents": "5000000000",
				"CashFlowsFromOperatingActivities": "3000000000",
				"CashFlowsFromInvestingActivities": "-1500000000",
				"CashFlowsFromFinancingActivities": "-500000000",
				"MaterialChangesInSubsidiaries": "false",
				"SignificantChangesInTheScopeOfConsolidation": "true",
				"ChangesInAccountingEstimates": "false",
				"ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard": "false",
				"RetrospectiveRestatement": "false",
				"NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock": "1000000000",
				"NumberOfTreasuryStockAtTheEndOfFiscalYear": "50000000",
				"ResultDividendPerShareAnnual": "60.0",
				"ResultPayoutRatioAnnual": "39.9",
				"ForecastDividendPerShareAnnual": "65.0",
				"ForecastPayoutRatioAnnual": "40.0",
				"ForecastNetSales": "15000000000",
				"ForecastOperatingProfit": "3000000000",
				"ForecastEarningsPerShare": "220.0",
				"NextYearForecastDividendPerShareAnnual": "70.0",
				"NextYearForecastPayoutRatioAnnual": "35.0",
				"NextYearForecastNetSales": "16000000000",
				"NextYearForecastEarningsPerShare": "250.0",
				"NonConsolidatedNetSales": "8000000000",
				"NonConsolidatedOperatingProfit": "1600000000",
				"NonConsolidatedProfit": "1200000000",
				"NonConsolidatedEarningsPerShare": "120.0"
			}
		]
	}`

	var resp StatementsResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Verify the conversion
	if len(resp.Statements) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(resp.Statements))
	}

	s := resp.Statements[0]

	// Check basic info
	if s.DisclosedDate != "2024-01-15" {
		t.Errorf("Expected DisclosedDate 2024-01-15, got %s", s.DisclosedDate)
	}
	if s.LocalCode != "72030" {
		t.Errorf("Expected LocalCode 72030, got %s", s.LocalCode)
	}
	if s.TypeOfCurrentPeriod != "3Q" {
		t.Errorf("Expected TypeOfCurrentPeriod 3Q, got %s", s.TypeOfCurrentPeriod)
	}

	// Check float64 conversions
	if s.NetSales == nil || *s.NetSales != 10000000000 {
		t.Errorf("Expected NetSales 10000000000, got %v", s.NetSales)
	}
	if s.EarningsPerShare == nil || *s.EarningsPerShare != 150.5 {
		t.Errorf("Expected EarningsPerShare 150.5, got %v", s.EarningsPerShare)
	}
	if s.CashFlowsFromInvestingActivities == nil || *s.CashFlowsFromInvestingActivities != -1500000000 {
		t.Errorf("Expected CashFlowsFromInvestingActivities -1500000000, got %v", s.CashFlowsFromInvestingActivities)
	}

	// Check bool conversions
	if s.MaterialChangesInSubsidiaries != false {
		t.Errorf("Expected MaterialChangesInSubsidiaries false, got %v", s.MaterialChangesInSubsidiaries)
	}
	if s.SignificantChangesInTheScopeOfConsolidation != true {
		t.Errorf("Expected SignificantChangesInTheScopeOfConsolidation true, got %v", s.SignificantChangesInTheScopeOfConsolidation)
	}

	// Check int64 conversions
	if s.NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock == nil || *s.NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock != 1000000000 {
		t.Errorf("Expected NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock 1000000000, got %v", s.NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock)
	}

	// Check empty string handling
	// Note: Due to UnmarshalJSON implementation, empty strings are converted to 0 for numeric types
	jsonDataEmpty := `{
		"statements": [
			{
				"DisclosedDate": "2024-01-15",
				"LocalCode": "72030",
				"NetSales": "",
				"MaterialChangesInSubsidiaries": "",
				"NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock": ""
			}
		]
	}`

	var respEmpty StatementsResponse
	err = json.Unmarshal([]byte(jsonDataEmpty), &respEmpty)
	if err != nil {
		t.Fatalf("UnmarshalJSON with empty strings failed: %v", err)
	}

	sEmpty := respEmpty.Statements[0]
	// Empty string is converted to 0 for numeric types in current implementation
	if sEmpty.NetSales == nil || *sEmpty.NetSales != 0 {
		t.Errorf("Expected NetSales 0 for empty string, got %v", sEmpty.NetSales)
	}
	if sEmpty.MaterialChangesInSubsidiaries != false {
		t.Errorf("Expected MaterialChangesInSubsidiaries false for empty string, got %v", sEmpty.MaterialChangesInSubsidiaries)
	}
	// Empty string is converted to 0 for numeric types in current implementation
	if sEmpty.NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock == nil || *sEmpty.NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock != 0 {
		t.Errorf("Expected NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock 0 for empty string, got %v", sEmpty.NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock)
	}
}
