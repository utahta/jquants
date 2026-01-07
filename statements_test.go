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
			wantPath: "/fins/summary?code=7203&date=20240101",
		},
		{
			name:     "with code only",
			code:     "7203",
			date:     "",
			wantPath: "/fins/summary?code=7203",
		},
		{
			name:     "with no parameters",
			code:     "",
			date:     "",
			wantPath: "/fins/summary",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewStatementsService(mockClient)

			// Mock response (v2 format)
			mockResponse := StatementsResponse{
				Data: []Statement{
					{
						// 基本情報
						DiscDate:   "2024-01-15",
						DiscTime:   "14:30:00",
						Code:       "72030",
						DiscNo:     "20240115123456",
						DocType:    "3QFinancialStatements_Consolidated_IFRS",
						CurPerType: "3Q",
						CurPerSt:   "2023-10-01",
						CurPerEn:   "2023-12-31",
						CurFYSt:    "2023-04-01",
						CurFYEn:    "2024-03-31",
						// 連結財務数値
						Sales:  floatPtr(10000000000),
						OP:     floatPtr(2000000000),
						OdP:    floatPtr(2100000000),
						NP:     floatPtr(1500000000),
						TA:     floatPtr(50000000000),
						Eq:     floatPtr(30000000000),
						CashEq: floatPtr(5000000000),
						CFO:    floatPtr(3000000000),
						CFI:    floatPtr(-1500000000),
						CFF:    floatPtr(-500000000),
						// 財務指標
						EPS:  floatPtr(150.5),
						DEPS: floatPtr(149.8),
						BPS:  floatPtr(3000.0),
						EqAR: floatPtr(60.0),
						// 配当情報
						DivAnn:        floatPtr(60.0),
						PayoutRatioAn: floatPtr(39.9),
						FDivAnn:       floatPtr(65.0),
						FPayoutRatioAn: floatPtr(40.0),
						Div2Q:         floatPtr(25.0),
						DivFY:         floatPtr(35.0),
						// 業績予想
						FSales: floatPtr(15000000000),
						FOP:    floatPtr(3000000000),
						FOdP:   floatPtr(3100000000),
						FNP:    floatPtr(2200000000),
						FEPS:   floatPtr(220.0),
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
		Data: []Statement{
			{
				DiscDate:   "2024-01-15",
				DiscTime:   "14:30:00",
				Code:       "72030",
				DiscNo:     "20240115123456",
				DocType:    "3QFinancialStatements_Consolidated_IFRS",
				CurPerType: "3Q",
				CurFYEn:    "2024-03-31",
				Sales:      floatPtr(10000000000),
				OP:         floatPtr(2000000000),
				OdP:        floatPtr(2100000000),
				NP:         floatPtr(1500000000),
				DivAnn:     floatPtr(50.0),
				FDivAnn:    floatPtr(55.0),
				EPS:        floatPtr(150.5),
				DEPS:       floatPtr(149.8),
			},
		},
	}
	mockClient.SetResponse("GET", "/fins/summary?code=7203", mockResponse)

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

	if statement.Sales == nil || *statement.Sales != 10000000000 {
		t.Errorf("Expected Sales 10000000000, got %v", statement.Sales)
	}

	if statement.EPS == nil || *statement.EPS != 150.5 {
		t.Errorf("Expected EPS 150.5, got %v", statement.EPS)
	}
}

func TestStatementsService_GetLatestStatements_NotFound(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewStatementsService(mockClient)

	// Mock empty response
	mockResponse := StatementsResponse{
		Data: []Statement{},
	}
	mockClient.SetResponse("GET", "/fins/summary?code=9999", mockResponse)

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
	mockClient.SetError("GET", "/fins/summary?code=7203", fmt.Errorf("API error"))

	// Test
	_, err := service.GetStatements("7203", "")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestStatementsService_GetStatementsByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewStatementsService(mockClient)

	// Mock response
	mockResponse := StatementsResponse{
		Data: []Statement{
			{
				DiscDate: "2024-01-15",
				Code:     "72030",
				Sales:    floatPtr(10000000000),
				OP:       floatPtr(2000000000),
			},
			{
				DiscDate: "2024-01-15",
				Code:     "86970",
				Sales:    floatPtr(5000000000),
				OP:       floatPtr(1000000000),
			},
		},
	}
	mockClient.SetResponse("GET", "/fins/summary?date=2024-01-15", mockResponse)

	// Test
	statements, err := service.GetStatementsByDate("2024-01-15")
	if err != nil {
		t.Errorf("GetStatementsByDate failed: %v", err)
	}

	// Verify
	if len(statements) != 2 {
		t.Errorf("Expected 2 statements, got %d", len(statements))
	}

	if mockClient.LastPath != "/fins/summary?date=2024-01-15" {
		t.Errorf("Expected path /fins/summary?date=2024-01-15, got %s", mockClient.LastPath)
	}
}

func TestStatementsResponse_UnmarshalJSON(t *testing.T) {
	// Test JSON with mixed types (strings that should be converted to float64/int64)
	// Using v2 API field names
	jsonData := `{
		"data": [
			{
				"DiscDate": "2024-01-15",
				"DiscTime": "14:30:00",
				"Code": "72030",
				"DiscNo": "20240115123456",
				"DocType": "3QFinancialStatements_Consolidated_IFRS",
				"CurPerType": "3Q",
				"CurPerSt": "2023-10-01",
				"CurPerEn": "2023-12-31",
				"CurFYSt": "2023-04-01",
				"CurFYEn": "2024-03-31",
				"Sales": "10000000000",
				"OP": "2000000000",
				"OdP": "2100000000",
				"NP": "1500000000",
				"EPS": "150.5",
				"DEPS": "149.8",
				"TA": "50000000000",
				"Eq": "30000000000",
				"EqAR": "60.0",
				"BPS": "3000.0",
				"CashEq": "5000000000",
				"CFO": "3000000000",
				"CFI": "-1500000000",
				"CFF": "-500000000",
				"MatChgSub": "false",
				"SigChgInC": "true",
				"ChgAcEst": "false",
				"ChgNoASRev": "false",
				"RetroRst": "false",
				"ShOutFY": "1000000000",
				"TrShFY": "50000000",
				"AvgSh": "975000000",
				"DivAnn": "60.0",
				"PayoutRatioAnn": "39.9",
				"DivUnit": "1000.0",
				"DivTotalAnn": "58500000000",
				"FDivAnn": "65.0",
				"FPayoutRatioAnn": "40.0",
				"FDivUnit": "1100.0",
				"FDivTotalAnn": "63375000000",
				"FSales": "15000000000",
				"FOP": "3000000000",
				"FEPS": "220.0",
				"NxFDiv1Q": "15.0",
				"NxFDiv2Q": "17.5",
				"NxFDiv3Q": "17.5",
				"NxFDivFY": "20.0",
				"NxFDivAnn": "70.0",
				"NxFDivUnit": "1200.0",
				"NxFPayoutRatioAnn": "35.0",
				"NxFSales2Q": "8000000000",
				"NxFOP2Q": "1600000000",
				"NxFOdP2Q": "1650000000",
				"NxFNp2Q": "1200000000",
				"NxFEPS2Q": "125.0",
				"NxFSales": "16000000000",
				"NxFEPS": "250.0",
				"ChgByASRev": "true",
				"NCSales": "8000000000",
				"NCOP": "1600000000",
				"NCNP": "1200000000",
				"NCEPS": "120.0",
				"FNCSales2Q": "4000000000",
				"FNCOP2Q": "800000000",
				"FNCOdP2Q": "820000000",
				"FNCNP2Q": "600000000",
				"FNCEPS2Q": "60.0",
				"NxFNCSales2Q": "4200000000",
				"NxFNCOP2Q": "840000000",
				"NxFNCOdP2Q": "860000000",
				"NxFNCNP2Q": "630000000",
				"NxFNCEPS2Q": "63.0",
				"FNCSales": "8500000000",
				"FNCOP": "1700000000",
				"FNCOdP": "1750000000",
				"FNCNP": "1300000000",
				"FNCEPS": "130.0",
				"NxFNCSales": "9000000000",
				"NxFNCOP": "1800000000",
				"NxFNCOdP": "1850000000",
				"NxFNCNP": "1400000000",
				"NxFNCEPS": "140.0"
			}
		]
	}`

	var resp StatementsResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	// Verify the conversion
	if len(resp.Data) != 1 {
		t.Fatalf("Expected 1 statement, got %d", len(resp.Data))
	}

	s := resp.Data[0]

	// Check basic info
	if s.DiscDate != "2024-01-15" {
		t.Errorf("Expected DiscDate 2024-01-15, got %s", s.DiscDate)
	}
	if s.Code != "72030" {
		t.Errorf("Expected Code 72030, got %s", s.Code)
	}
	if s.CurPerType != "3Q" {
		t.Errorf("Expected CurPerType 3Q, got %s", s.CurPerType)
	}

	// Check float64 conversions
	if s.Sales == nil || *s.Sales != 10000000000 {
		t.Errorf("Expected Sales 10000000000, got %v", s.Sales)
	}
	if s.EPS == nil || *s.EPS != 150.5 {
		t.Errorf("Expected EPS 150.5, got %v", s.EPS)
	}
	if s.CFI == nil || *s.CFI != -1500000000 {
		t.Errorf("Expected CFI -1500000000, got %v", s.CFI)
	}

	// Check string fields (bool values are now strings in v2)
	if s.MatChgSub != "false" {
		t.Errorf("Expected MatChgSub 'false', got %s", s.MatChgSub)
	}
	if s.SigChgInC != "true" {
		t.Errorf("Expected SigChgInC 'true', got %s", s.SigChgInC)
	}

	// Check int64 conversions
	if s.ShOutFY == nil || *s.ShOutFY != 1000000000 {
		t.Errorf("Expected ShOutFY 1000000000, got %v", s.ShOutFY)
	}
	if s.AvgSh == nil || *s.AvgSh != 975000000 {
		t.Errorf("Expected AvgSh 975000000, got %v", s.AvgSh)
	}

	// Check dividend fields
	if s.DivUnit == nil || *s.DivUnit != 1000.0 {
		t.Errorf("Expected DivUnit 1000.0, got %v", s.DivUnit)
	}
	if s.DivTotalAnn == nil || *s.DivTotalAnn != 58500000000 {
		t.Errorf("Expected DivTotalAnn 58500000000, got %v", s.DivTotalAnn)
	}
	if s.FDivUnit == nil || *s.FDivUnit != 1100.0 {
		t.Errorf("Expected FDivUnit 1100.0, got %v", s.FDivUnit)
	}
	if s.FDivTotalAnn == nil || *s.FDivTotalAnn != 63375000000 {
		t.Errorf("Expected FDivTotalAnn 63375000000, got %v", s.FDivTotalAnn)
	}

	// Check NextYear quarterly dividend fields
	if s.NxFDiv1Q == nil || *s.NxFDiv1Q != 15.0 {
		t.Errorf("Expected NxFDiv1Q 15.0, got %v", s.NxFDiv1Q)
	}
	if s.NxFDiv2Q == nil || *s.NxFDiv2Q != 17.5 {
		t.Errorf("Expected NxFDiv2Q 17.5, got %v", s.NxFDiv2Q)
	}
	if s.NxFDiv3Q == nil || *s.NxFDiv3Q != 17.5 {
		t.Errorf("Expected NxFDiv3Q 17.5, got %v", s.NxFDiv3Q)
	}
	if s.NxFDivFY == nil || *s.NxFDivFY != 20.0 {
		t.Errorf("Expected NxFDivFY 20.0, got %v", s.NxFDivFY)
	}
	if s.NxFDivUnit == nil || *s.NxFDivUnit != 1200.0 {
		t.Errorf("Expected NxFDivUnit 1200.0, got %v", s.NxFDivUnit)
	}

	// Check NextYear quarterly forecast fields
	if s.NxFSales2Q == nil || *s.NxFSales2Q != 8000000000 {
		t.Errorf("Expected NxFSales2Q 8000000000, got %v", s.NxFSales2Q)
	}
	if s.NxFOP2Q == nil || *s.NxFOP2Q != 1600000000 {
		t.Errorf("Expected NxFOP2Q 1600000000, got %v", s.NxFOP2Q)
	}
	if s.NxFOdP2Q == nil || *s.NxFOdP2Q != 1650000000 {
		t.Errorf("Expected NxFOdP2Q 1650000000, got %v", s.NxFOdP2Q)
	}
	if s.NxFNp2Q == nil || *s.NxFNp2Q != 1200000000 {
		t.Errorf("Expected NxFNp2Q 1200000000, got %v", s.NxFNp2Q)
	}
	if s.NxFEPS2Q == nil || *s.NxFEPS2Q != 125.0 {
		t.Errorf("Expected NxFEPS2Q 125.0, got %v", s.NxFEPS2Q)
	}

	// Check accounting standard field
	if s.ChgByASRev != "true" {
		t.Errorf("Expected ChgByASRev 'true', got %s", s.ChgByASRev)
	}

	// Check NonConsolidated forecast fields
	if s.FNCSales2Q == nil || *s.FNCSales2Q != 4000000000 {
		t.Errorf("Expected FNCSales2Q 4000000000, got %v", s.FNCSales2Q)
	}
	if s.FNCOP2Q == nil || *s.FNCOP2Q != 800000000 {
		t.Errorf("Expected FNCOP2Q 800000000, got %v", s.FNCOP2Q)
	}
	if s.FNCOdP2Q == nil || *s.FNCOdP2Q != 820000000 {
		t.Errorf("Expected FNCOdP2Q 820000000, got %v", s.FNCOdP2Q)
	}
	if s.FNCNP2Q == nil || *s.FNCNP2Q != 600000000 {
		t.Errorf("Expected FNCNP2Q 600000000, got %v", s.FNCNP2Q)
	}
	if s.FNCEPS2Q == nil || *s.FNCEPS2Q != 60.0 {
		t.Errorf("Expected FNCEPS2Q 60.0, got %v", s.FNCEPS2Q)
	}

	// Check NextYear NonConsolidated forecast fields
	if s.NxFNCSales2Q == nil || *s.NxFNCSales2Q != 4200000000 {
		t.Errorf("Expected NxFNCSales2Q 4200000000, got %v", s.NxFNCSales2Q)
	}
	if s.NxFNCEPS2Q == nil || *s.NxFNCEPS2Q != 63.0 {
		t.Errorf("Expected NxFNCEPS2Q 63.0, got %v", s.NxFNCEPS2Q)
	}

	// Check NonConsolidated annual forecast fields
	if s.FNCSales == nil || *s.FNCSales != 8500000000 {
		t.Errorf("Expected FNCSales 8500000000, got %v", s.FNCSales)
	}
	if s.FNCOP == nil || *s.FNCOP != 1700000000 {
		t.Errorf("Expected FNCOP 1700000000, got %v", s.FNCOP)
	}
	if s.FNCOdP == nil || *s.FNCOdP != 1750000000 {
		t.Errorf("Expected FNCOdP 1750000000, got %v", s.FNCOdP)
	}
	if s.FNCNP == nil || *s.FNCNP != 1300000000 {
		t.Errorf("Expected FNCNP 1300000000, got %v", s.FNCNP)
	}
	if s.FNCEPS == nil || *s.FNCEPS != 130.0 {
		t.Errorf("Expected FNCEPS 130.0, got %v", s.FNCEPS)
	}

	// Check NextYear NonConsolidated annual forecast fields
	if s.NxFNCSales == nil || *s.NxFNCSales != 9000000000 {
		t.Errorf("Expected NxFNCSales 9000000000, got %v", s.NxFNCSales)
	}
	if s.NxFNCEPS == nil || *s.NxFNCEPS != 140.0 {
		t.Errorf("Expected NxFNCEPS 140.0, got %v", s.NxFNCEPS)
	}

	// Check empty string handling
	jsonDataEmpty := `{
		"data": [
			{
				"DiscDate": "2024-01-15",
				"Code": "72030",
				"Sales": "",
				"MatChgSub": "",
				"ShOutFY": ""
			}
		]
	}`

	var respEmpty StatementsResponse
	err = json.Unmarshal([]byte(jsonDataEmpty), &respEmpty)
	if err != nil {
		t.Fatalf("UnmarshalJSON with empty strings failed: %v", err)
	}

	sEmpty := respEmpty.Data[0]
	// Empty string is converted to 0 for numeric types in current implementation
	if sEmpty.Sales == nil || *sEmpty.Sales != 0 {
		t.Errorf("Expected Sales 0 for empty string, got %v", sEmpty.Sales)
	}
	if sEmpty.MatChgSub != "" {
		t.Errorf("Expected MatChgSub '' for empty string, got %s", sEmpty.MatChgSub)
	}
	// Empty string is converted to 0 for numeric types in current implementation
	if sEmpty.ShOutFY == nil || *sEmpty.ShOutFY != 0 {
		t.Errorf("Expected ShOutFY 0 for empty string, got %v", sEmpty.ShOutFY)
	}
}
