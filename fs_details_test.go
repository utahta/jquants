package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestFSDetailsService_GetFSDetails(t *testing.T) {
	tests := []struct {
		name     string
		params   FSDetailsParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with code and date",
			params: FSDetailsParams{
				Code: "86970",
				Date: "20230130",
			},
			wantPath: "/fins/details?code=86970&date=20230130",
		},
		{
			name: "with code only",
			params: FSDetailsParams{
				Code: "86970",
			},
			wantPath: "/fins/details?code=86970",
		},
		{
			name: "with date only",
			params: FSDetailsParams{
				Date: "2023-01-30",
			},
			wantPath: "/fins/details?date=2023-01-30",
		},
		{
			name: "with pagination key",
			params: FSDetailsParams{
				Code:          "86970",
				PaginationKey: "key123",
			},
			wantPath: "/fins/details?code=86970&pagination_key=key123",
		},
		{
			name:     "with no required parameters",
			params:   FSDetailsParams{},
			wantPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewFSDetailsService(mockClient)

			// Mock response based on documentation sample
			mockResponse := FSDetailsResponse{
				Data: []FSDetail{
					{
						DiscDate: "2023-01-30",
						DiscTime: "12:00:00",
						Code:     "86970",
						DiscNo:   "20230127594871",
						DocType:  "3QFinancialStatements_Consolidated_IFRS",
						FS: map[string]string{
							"Goodwill (IFRS)":                                       "67374000000",
							"Retained earnings (IFRS)":                              "263894000000",
							"Operating profit (loss) (IFRS)":                        "51765000000.0",
							"Previous fiscal year end date, DEI":                    "2022-03-31",
							"Basic earnings (loss) per share (IFRS)":                "66.76",
							"Document type, DEI":                                    "四半期第３号参考様式　[IFRS]（連結）",
							"Current period end date, DEI":                          "2022-12-31",
							"Revenue - 2 (IFRS)":                                    "100987000000.0",
							"Profit (loss) attributable to owners of parent (IFRS)": "35175000000.0",
							"Current liabilities (IFRS)":                            "78852363000000",
							"Equity attributable to owners of parent (IFRS)":        "311103000000",
							"Non-current liabilities (IFRS)":                        "33476000000",
							"Property, plant and equipment (IFRS)":                  "11277000000",
							"Cash and cash equivalents (IFRS)":                      "91135000000",
							"Share capital (IFRS)":                                  "11500000000",
							"Assets (IFRS)":                                         "79205861000000",
							"Equity (IFRS)":                                         "320021000000",
							"Liabilities (IFRS)":                                    "78885839000000",
							"Accounting standards, DEI":                             "IFRS",
							"Current assets (IFRS)":                                 "100000000000", // For ratio test
						},
					},
				},
				PaginationKey: "",
			}
			if !tt.wantErr {
				mockClient.SetResponse("GET", tt.wantPath, mockResponse)
			}

			// Execute
			resp, err := service.GetFSDetails(tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetFSDetails() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetFSDetails() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetFSDetails() returned nil response")
				return
			}
			if len(resp.Data) == 0 {
				t.Error("GetFSDetails() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetFSDetails() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestFSDetailsService_GetFSDetails_RequiresParameter(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewFSDetailsService(mockClient)

	// Execute with empty parameters
	_, err := service.GetFSDetails(FSDetailsParams{})

	// Verify
	if err == nil {
		t.Error("GetFSDetails() expected error for missing parameters but got nil")
	}
	if err.Error() != "either code or date parameter is required" {
		t.Errorf("GetFSDetails() error = %v, want 'either code or date parameter is required'", err)
	}
}

func TestFSDetailsService_GetFSDetailsByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewFSDetailsService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := FSDetailsResponse{
		Data: []FSDetail{
			{
				DiscDate: "2023-01-30",
				Code:     "86970",
				DocType:  "3QFinancialStatements_Consolidated_IFRS",
				FS: map[string]string{
					"Accounting standards, DEI": "IFRS",
				},
			},
			{
				DiscDate: "2022-10-31",
				Code:     "86970",
				DocType:  "2QFinancialStatements_Consolidated_IFRS",
				FS: map[string]string{
					"Accounting standards, DEI": "IFRS",
				},
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := FSDetailsResponse{
		Data: []FSDetail{
			{
				DiscDate: "2022-07-29",
				Code:     "86970",
				DocType:  "1QFinancialStatements_Consolidated_IFRS",
				FS: map[string]string{
					"Accounting standards, DEI": "IFRS",
				},
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/fins/details?code=86970", mockResponse1)
	mockClient.SetResponse("GET", "/fins/details?code=86970&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetFSDetailsByCode("86970")

	// Verify
	if err != nil {
		t.Fatalf("GetFSDetailsByCode() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetFSDetailsByCode() returned %d items, want 3", len(data))
	}
	for _, item := range data {
		if item.Code != "86970" {
			t.Errorf("GetFSDetailsByCode() returned code %v, want 86970", item.Code)
		}
	}
}

func TestFSDetailsService_GetFSDetailsByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewFSDetailsService(mockClient)

	// Mock response
	mockResponse := FSDetailsResponse{
		Data: []FSDetail{
			{
				DiscDate: "2023-01-30",
				Code:     "86970",
			},
			{
				DiscDate: "2023-01-30",
				Code:     "13010",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/fins/details?date=20230130", mockResponse)

	// Execute
	data, err := service.GetFSDetailsByDate("20230130")

	// Verify
	if err != nil {
		t.Fatalf("GetFSDetailsByDate() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetFSDetailsByDate() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.DiscDate != "2023-01-30" {
			t.Errorf("GetFSDetailsByDate() returned date %v, want 2023-01-30", item.DiscDate)
		}
	}
}

func TestFSDetail_AccountingStandardsMethods(t *testing.T) {
	tests := []struct {
		name           string
		fsDetail       FSDetail
		isIFRS         bool
		isJapaneseGAAP bool
	}{
		{
			name: "IFRS",
			fsDetail: FSDetail{
				FS: map[string]string{
					"Accounting standards, DEI": "IFRS",
				},
			},
			isIFRS:         true,
			isJapaneseGAAP: false,
		},
		{
			name: "Japanese GAAP",
			fsDetail: FSDetail{
				FS: map[string]string{
					"Accounting standards, DEI": "JapaneseGAAP",
				},
			},
			isIFRS:         false,
			isJapaneseGAAP: true,
		},
		{
			name: "No standard specified",
			fsDetail: FSDetail{
				FS: map[string]string{},
			},
			isIFRS:         false,
			isJapaneseGAAP: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fsDetail.IsIFRS(); got != tt.isIFRS {
				t.Errorf("IsIFRS() = %v, want %v", got, tt.isIFRS)
			}
			if got := tt.fsDetail.IsJapaneseGAAP(); got != tt.isJapaneseGAAP {
				t.Errorf("IsJapaneseGAAP() = %v, want %v", got, tt.isJapaneseGAAP)
			}
		})
	}
}

func TestFSDetail_DocumentTypeMethods(t *testing.T) {
	tests := []struct {
		name           string
		fsDetail       FSDetail
		isQuarterly    bool
		isAnnual       bool
		isConsolidated bool
		quarter        int
	}{
		{
			name: "1Q Consolidated IFRS",
			fsDetail: FSDetail{
				DocType: "1QFinancialStatements_Consolidated_IFRS",
			},
			isQuarterly:    true,
			isAnnual:       false,
			isConsolidated: true,
			quarter:        1,
		},
		{
			name: "2Q Consolidated IFRS",
			fsDetail: FSDetail{
				DocType: "2QFinancialStatements_Consolidated_IFRS",
			},
			isQuarterly:    true,
			isAnnual:       false,
			isConsolidated: true,
			quarter:        2,
		},
		{
			name: "3Q Consolidated IFRS",
			fsDetail: FSDetail{
				DocType: "3QFinancialStatements_Consolidated_IFRS",
			},
			isQuarterly:    true,
			isAnnual:       false,
			isConsolidated: true,
			quarter:        3,
		},
		{
			name: "FY Consolidated IFRS",
			fsDetail: FSDetail{
				DocType: "FYFinancialStatements_Consolidated_IFRS",
			},
			isQuarterly:    false,
			isAnnual:       true,
			isConsolidated: true,
			quarter:        0,
		},
		{
			name: "Other document type",
			fsDetail: FSDetail{
				DocType: "OtherDocument",
			},
			isQuarterly:    false,
			isAnnual:       false,
			isConsolidated: false,
			quarter:        -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fsDetail.IsQuarterly(); got != tt.isQuarterly {
				t.Errorf("IsQuarterly() = %v, want %v", got, tt.isQuarterly)
			}
			if got := tt.fsDetail.IsAnnual(); got != tt.isAnnual {
				t.Errorf("IsAnnual() = %v, want %v", got, tt.isAnnual)
			}
			if got := tt.fsDetail.IsConsolidated(); got != tt.isConsolidated {
				t.Errorf("IsConsolidated() = %v, want %v", got, tt.isConsolidated)
			}
			if got := tt.fsDetail.GetQuarter(); got != tt.quarter {
				t.Errorf("GetQuarter() = %v, want %v", got, tt.quarter)
			}
		})
	}
}

func TestFSDetail_GetValue(t *testing.T) {
	// Setup
	fsDetail := FSDetail{
		FS: map[string]string{
			"Assets (IFRS)": "79205861000000",
			"Equity (IFRS)": "320021000000",
		},
	}

	// Test existing key
	value, ok := fsDetail.GetValue("Assets (IFRS)")
	if !ok {
		t.Error("GetValue() expected to find key but got false")
	}
	if value != "79205861000000" {
		t.Errorf("GetValue() = %v, want %v", value, "79205861000000")
	}

	// Test non-existing key
	_, ok = fsDetail.GetValue("NonExistingKey")
	if ok {
		t.Error("GetValue() expected not to find key but got true")
	}
}

func TestFSDetail_GetFloatValue(t *testing.T) {
	// Setup
	fsDetail := FSDetail{
		FS: map[string]string{
			"Assets (IFRS)":                          "79205861000000",
			"Basic earnings (loss) per share (IFRS)": "66.76",
			"Invalid value":                          "not_a_number",
		},
	}

	// Test valid integer value
	value, err := fsDetail.GetFloatValue("Assets (IFRS)")
	if err != nil {
		t.Fatalf("GetFloatValue() error = %v", err)
	}
	if value != 79205861000000 {
		t.Errorf("GetFloatValue() = %v, want %v", value, 79205861000000)
	}

	// Test valid float value
	value, err = fsDetail.GetFloatValue("Basic earnings (loss) per share (IFRS)")
	if err != nil {
		t.Fatalf("GetFloatValue() error = %v", err)
	}
	if value != 66.76 {
		t.Errorf("GetFloatValue() = %v, want %v", value, 66.76)
	}

	// Test non-existing key
	_, err = fsDetail.GetFloatValue("NonExistingKey")
	if err == nil {
		t.Error("GetFloatValue() expected error for non-existing key but got nil")
	}

	// Test invalid value
	_, err = fsDetail.GetFloatValue("Invalid value")
	if err == nil {
		t.Error("GetFloatValue() expected error for invalid value but got nil")
	}
}

func TestFSDetail_FinancialRatios(t *testing.T) {
	// Setup - IFRS financial statement
	fsDetailIFRS := FSDetail{
		FS: map[string]string{
			"Accounting standards, DEI":                             "IFRS",
			"Profit (loss) attributable to owners of parent (IFRS)": "35175000000",
			"Equity attributable to owners of parent (IFRS)":        "311103000000",
			"Current assets (IFRS)":                                 "100000000000",
			"Current liabilities (IFRS)":                            "78852363000000",
			"Equity (IFRS)":                                         "320021000000",
			"Assets (IFRS)":                                         "79205861000000",
			"Basic earnings (loss) per share (IFRS)":                "66.76",
		},
	}

	// Test GetROE
	roe, err := fsDetailIFRS.GetROE()
	if err != nil {
		t.Fatalf("GetROE() error = %v", err)
	}
	if roe == nil {
		t.Fatal("GetROE() returned nil")
		return
	}
	expectedROE := (35175000000.0 / 311103000000.0) * 100
	if *roe != expectedROE {
		t.Errorf("GetROE() = %v, want %v", *roe, expectedROE)
	}

	// Test GetCurrentRatio
	currentRatio, err := fsDetailIFRS.GetCurrentRatio()
	if err != nil {
		t.Fatalf("GetCurrentRatio() error = %v", err)
	}
	if currentRatio == nil {
		t.Fatal("GetCurrentRatio() returned nil")
		return
	}
	expectedRatio := 100000000000.0 / 78852363000000.0
	if *currentRatio != expectedRatio {
		t.Errorf("GetCurrentRatio() = %v, want %v", *currentRatio, expectedRatio)
	}

	// Test GetEquityRatio
	equityRatio, err := fsDetailIFRS.GetEquityRatio()
	if err != nil {
		t.Fatalf("GetEquityRatio() error = %v", err)
	}
	if equityRatio == nil {
		t.Fatal("GetEquityRatio() returned nil")
		return
	}
	expectedEquityRatio := (320021000000.0 / 79205861000000.0) * 100
	if diff := *equityRatio - expectedEquityRatio; diff > 0.0000001 || diff < -0.0000001 {
		t.Errorf("GetEquityRatio() = %v, want %v", *equityRatio, expectedEquityRatio)
	}

	// Test GetBasicEPS
	eps, err := fsDetailIFRS.GetBasicEPS()
	if err != nil {
		t.Fatalf("GetBasicEPS() error = %v", err)
	}
	if eps == nil {
		t.Fatal("GetBasicEPS() returned nil")
		return
	}
	if *eps != 66.76 {
		t.Errorf("GetBasicEPS() = %v, want %v", *eps, 66.76)
	}

	// Test non-IFRS error
	fsDetailJGAAP := FSDetail{
		FS: map[string]string{
			"Accounting standards, DEI": "JapaneseGAAP",
		},
	}
	_, err = fsDetailJGAAP.GetROE()
	if err == nil {
		t.Error("GetROE() expected error for non-IFRS but got nil")
	}
}

func TestFSDetail_FinancialRatios_ErrorCases(t *testing.T) {
	// Setup - IFRS with zero values
	fsDetailZeroEquity := FSDetail{
		FS: map[string]string{
			"Accounting standards, DEI":                             "IFRS",
			"Profit (loss) attributable to owners of parent (IFRS)": "35175000000",
			"Equity attributable to owners of parent (IFRS)":        "0",
		},
	}

	// Test GetROE with zero equity
	_, err := fsDetailZeroEquity.GetROE()
	if err == nil {
		t.Error("GetROE() expected error for zero equity but got nil")
	}

	// Setup - IFRS with zero liabilities
	fsDetailZeroLiabilities := FSDetail{
		FS: map[string]string{
			"Accounting standards, DEI":  "IFRS",
			"Current assets (IFRS)":      "100000000000",
			"Current liabilities (IFRS)": "0",
		},
	}

	// Test GetCurrentRatio with zero liabilities
	_, err = fsDetailZeroLiabilities.GetCurrentRatio()
	if err == nil {
		t.Error("GetCurrentRatio() expected error for zero liabilities but got nil")
	}

	// Setup - IFRS with zero assets
	fsDetailZeroAssets := FSDetail{
		FS: map[string]string{
			"Accounting standards, DEI": "IFRS",
			"Equity (IFRS)":             "320021000000",
			"Assets (IFRS)":             "0",
		},
	}

	// Test GetEquityRatio with zero assets
	_, err = fsDetailZeroAssets.GetEquityRatio()
	if err == nil {
		t.Error("GetEquityRatio() expected error for zero assets but got nil")
	}
}

func TestFSDetailsService_GetFSDetails_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewFSDetailsService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/fins/details?code=86970", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetFSDetails(FSDetailsParams{Code: "86970"})

	// Verify
	if err == nil {
		t.Error("GetFSDetails() expected error but got nil")
	}
}
