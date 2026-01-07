package jquants

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestOptionsService_GetOptions(t *testing.T) {
	tests := []struct {
		name     string
		params   OptionsParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with all parameters",
			params: OptionsParams{
				Date:          "20240723",
				Category:      "NK225E",
				Code:          "7203",
				ContractFlag:  "1",
				PaginationKey: "test_key",
			},
			wantPath: "/derivatives/bars/daily/options?date=20240723&category=NK225E&code=7203&contract_flag=1&pagination_key=test_key",
		},
		{
			name: "with date only (required)",
			params: OptionsParams{
				Date: "20240723",
			},
			wantPath: "/derivatives/bars/daily/options?date=20240723",
		},
		{
			name: "with date and category",
			params: OptionsParams{
				Date:     "20240723",
				Category: "TOPIXE",
			},
			wantPath: "/derivatives/bars/daily/options?date=20240723&category=TOPIXE",
		},
		{
			name: "with EQOP category and code",
			params: OptionsParams{
				Date:     "20240723",
				Category: "EQOP",
				Code:     "7203",
			},
			wantPath: "/derivatives/bars/daily/options?date=20240723&category=EQOP&code=7203",
		},
		{
			name:    "without date (should error)",
			params:  OptionsParams{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewOptionsService(mockClient)

			if !tt.wantErr {
				// Mock response
				mockResponse := OptionsResponse{
					Data: []Option{
						{
							Code:         "140014505",
							ProdCat:      "TOPIXE",
							UndSSO:       "-",
							Date:         "2024-07-23",
							CM:           "2025-01",
							Strike:       2450.0,
							PCDiv:        "2",
							EmMrgnTrgDiv: "002",
							O:            0.0,
							H:            0.0,
							L:            0.0,
							C:            0.0,
							EO:           "",
							EH:           "",
							EL:           "",
							EC:           "",
							AO:           0.0,
							AH:           0.0,
							AL:           0.0,
							AC:           0.0,
							MO:           "",
							MH:           "",
							ML:           "",
							MC:           "",
							Vo:           0.0,
							OI:           0.0,
							Va:           0.0,
							VoOA:         floatPtr(0.0),
							Settle:       floatPtr(377.0),
							Theo:         floatPtr(380.3801),
							BaseVol:      floatPtr(18.115),
							UnderPx:      floatPtr(2833.39),
							IV:           floatPtr(17.2955),
							IR:           floatPtr(0.3527),
							LTD:          stringPtr("2025-01-09"),
							SQD:          stringPtr("2025-01-10"),
							CCMFlag:      stringPtr("0"),
						},
					},
					PaginationKey: "value1.value2.",
				}
				mockClient.SetResponse("GET", tt.wantPath, mockResponse)
			}

			// Execute
			resp, err := service.GetOptions(tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetOptions() expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("GetOptions() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetOptions() returned nil response")
				return
			}
			if len(resp.Data) == 0 {
				t.Error("GetOptions() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetOptions() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestOptionsService_GetOptionsByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewOptionsService(mockClient)

	// Mock response - first page
	mockResponse1 := OptionsResponse{
		Data: []Option{
			{
				Code:         "140014505",
				ProdCat:      "TOPIXE",
				UndSSO:       "-",
				Date:         "2024-07-23",
				CM:           "2025-01",
				Strike:       2450.0,
				PCDiv:        "2",
				EmMrgnTrgDiv: "002",
				C:            0.0,
				Vo:           0.0,
			},
		},
		PaginationKey: "next_page",
	}
	mockClient.SetResponse("GET", "/derivatives/bars/daily/options?date=20240723", mockResponse1)

	// Mock response - second page
	mockResponse2 := OptionsResponse{
		Data: []Option{
			{
				Code:         "140014506",
				ProdCat:      "NK225E",
				UndSSO:       "-",
				Date:         "2024-07-23",
				CM:           "2025-01",
				Strike:       40000.0,
				PCDiv:        "1",
				EmMrgnTrgDiv: "002",
				C:            50.0,
				Vo:           100.0,
			},
		},
		PaginationKey: "", // 最終ページ
	}
	mockClient.SetResponse("GET", "/derivatives/bars/daily/options?date=20240723&pagination_key=next_page", mockResponse2)

	// Execute
	options, err := service.GetOptionsByDate("20240723")

	// Verify
	if err != nil {
		t.Fatalf("GetOptionsByDate() error = %v", err)
	}
	if len(options) != 2 {
		t.Errorf("GetOptionsByDate() returned %d items, want 2", len(options))
	}
	if options[0].Code != "140014505" {
		t.Errorf("GetOptionsByDate() first item code = %v, want 140014505", options[0].Code)
	}
	if options[1].Code != "140014506" {
		t.Errorf("GetOptionsByDate() second item code = %v, want 140014506", options[1].Code)
	}
}

func TestOptionsService_GetOptionsByCategory(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewOptionsService(mockClient)

	// Mock response
	mockResponse := OptionsResponse{
		Data: []Option{
			{
				Code:         "140014505",
				ProdCat:      "NK225E",
				UndSSO:       "-",
				Date:         "2024-07-23",
				CM:           "2025-01",
				Strike:       40000.0,
				PCDiv:        "2",
				EmMrgnTrgDiv: "002",
			},
		},
	}
	mockClient.SetResponse("GET", "/derivatives/bars/daily/options?date=20240723&category=NK225E", mockResponse)

	// Execute
	options, err := service.GetOptionsByCategory("20240723", "NK225E")

	// Verify
	if err != nil {
		t.Fatalf("GetOptionsByCategory() error = %v", err)
	}
	if len(options) != 1 {
		t.Errorf("GetOptionsByCategory() returned %d items, want 1", len(options))
	}
	if options[0].ProdCat != "NK225E" {
		t.Errorf("GetOptionsByCategory() category = %v, want NK225E", options[0].ProdCat)
	}
}

func TestOptionsService_GetSecurityOptionsByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewOptionsService(mockClient)

	// Mock response
	mockResponse := OptionsResponse{
		Data: []Option{
			{
				Code:         "10014505",
				ProdCat:      "EQOP",
				UndSSO:       "7203",
				Date:         "2024-07-23",
				CM:           "2025-01",
				Strike:       2500.0,
				PCDiv:        "1",
				EmMrgnTrgDiv: "002",
			},
		},
	}
	mockClient.SetResponse("GET", "/derivatives/bars/daily/options?date=20240723&category=EQOP&code=7203", mockResponse)

	// Execute
	options, err := service.GetSecurityOptionsByCode("20240723", "7203")

	// Verify
	if err != nil {
		t.Fatalf("GetSecurityOptionsByCode() error = %v", err)
	}
	if len(options) != 1 {
		t.Errorf("GetSecurityOptionsByCode() returned %d items, want 1", len(options))
	}
	if options[0].UndSSO != "7203" {
		t.Errorf("GetSecurityOptionsByCode() UndSSO = %v, want 7203", options[0].UndSSO)
	}
	if options[0].ProdCat != "EQOP" {
		t.Errorf("GetSecurityOptionsByCode() category = %v, want EQOP", options[0].ProdCat)
	}
}

func TestOptionsService_GetCentralContractMonthOptions(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewOptionsService(mockClient)

	// Mock response
	mockResponse := OptionsResponse{
		Data: []Option{
			{
				Code:         "140014505",
				ProdCat:      "NK225E",
				UndSSO:       "-",
				Date:         "2024-07-23",
				CM:           "2025-01",
				Strike:       40000.0,
				PCDiv:        "2",
				EmMrgnTrgDiv: "002",
				CCMFlag:      stringPtr("1"),
			},
		},
	}
	mockClient.SetResponse("GET", "/derivatives/bars/daily/options?date=20240723&contract_flag=1", mockResponse)

	// Execute
	options, err := service.GetCentralContractMonthOptions("20240723")

	// Verify
	if err != nil {
		t.Fatalf("GetCentralContractMonthOptions() error = %v", err)
	}
	if len(options) != 1 {
		t.Errorf("GetCentralContractMonthOptions() returned %d items, want 1", len(options))
	}
	if !options[0].IsCentralContractMonth() {
		t.Error("GetCentralContractMonthOptions() returned non-central contract month option")
	}
}

func TestOptionsService_GetOptions_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewOptionsService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/derivatives/bars/daily/options?date=20240723", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetOptions(OptionsParams{Date: "20240723"})

	// Verify
	if err == nil {
		t.Error("GetOptions() expected error but got nil")
	}
}

func TestOption_HelperMethods(t *testing.T) {
	t.Run("IsCall and IsPut", func(t *testing.T) {
		callOption := Option{PCDiv: "2"}
		putOption := Option{PCDiv: "1"}

		if !callOption.IsCall() {
			t.Error("IsCall() returned false for call option")
		}
		if callOption.IsPut() {
			t.Error("IsPut() returned true for call option")
		}
		if putOption.IsCall() {
			t.Error("IsCall() returned true for put option")
		}
		if !putOption.IsPut() {
			t.Error("IsPut() returned false for put option")
		}
	})

	t.Run("IsEmergencyMarginTriggered", func(t *testing.T) {
		triggered := Option{EmMrgnTrgDiv: "001"}
		normal := Option{EmMrgnTrgDiv: "002"}

		if !triggered.IsEmergencyMarginTriggered() {
			t.Error("IsEmergencyMarginTriggered() returned false for triggered")
		}
		if normal.IsEmergencyMarginTriggered() {
			t.Error("IsEmergencyMarginTriggered() returned true for normal")
		}
	})

	t.Run("IsCentralContractMonth", func(t *testing.T) {
		central := Option{CCMFlag: stringPtr("1")}
		nonCentral := Option{CCMFlag: stringPtr("0")}
		nilFlag := Option{}

		if !central.IsCentralContractMonth() {
			t.Error("IsCentralContractMonth() returned false for central")
		}
		if nonCentral.IsCentralContractMonth() {
			t.Error("IsCentralContractMonth() returned true for non-central")
		}
		if nilFlag.IsCentralContractMonth() {
			t.Error("IsCentralContractMonth() returned true for nil flag")
		}
	})

	t.Run("IsSecurityOption", func(t *testing.T) {
		securityOption := Option{UndSSO: "7203"}
		indexOption := Option{UndSSO: "-"}

		if !securityOption.IsSecurityOption() {
			t.Error("IsSecurityOption() returned false for security option")
		}
		if indexOption.IsSecurityOption() {
			t.Error("IsSecurityOption() returned true for index option")
		}
	})

	t.Run("Session data helpers", func(t *testing.T) {
		optionWithNight := Option{
			EO: 100.0,
			EH: 110.0,
			EL: 90.0,
			EC: 105.0,
		}

		optionNoNight := Option{
			EO: "",
			EH: "",
			EL: "",
			EC: "",
		}

		if !optionWithNight.HasNightSession() {
			t.Error("HasNightSession() returned false for option with night session")
		}
		if optionNoNight.HasNightSession() {
			t.Error("HasNightSession() returned true for option without night session")
		}

		// Test session value getters
		if open := optionWithNight.GetNightSessionOpen(); open == nil || *open != 100.0 {
			t.Errorf("GetNightSessionOpen() = %v, want 100.0", open)
		}
		if open := optionNoNight.GetNightSessionOpen(); open != nil {
			t.Errorf("GetNightSessionOpen() = %v, want nil", open)
		}
	})

	t.Run("ITM/OTM/ATM", func(t *testing.T) {
		callITM := Option{
			PCDiv:   "2",
			Strike:  100.0,
			UnderPx: floatPtr(110.0),
		}
		callOTM := Option{
			PCDiv:   "2",
			Strike:  100.0,
			UnderPx: floatPtr(90.0),
		}
		callATM := Option{
			PCDiv:   "2",
			Strike:  100.0,
			UnderPx: floatPtr(100.0),
		}
		putITM := Option{
			PCDiv:   "1",
			Strike:  100.0,
			UnderPx: floatPtr(90.0),
		}
		putOTM := Option{
			PCDiv:   "1",
			Strike:  100.0,
			UnderPx: floatPtr(110.0),
		}

		if !callITM.IsITM() {
			t.Error("IsITM() returned false for ITM call")
		}
		if callITM.IsOTM() {
			t.Error("IsOTM() returned true for ITM call")
		}
		if !callOTM.IsOTM() {
			t.Error("IsOTM() returned false for OTM call")
		}
		if callOTM.IsITM() {
			t.Error("IsITM() returned true for OTM call")
		}
		if !callATM.IsATM() {
			t.Error("IsATM() returned false for ATM call")
		}
		if !putITM.IsITM() {
			t.Error("IsITM() returned false for ITM put")
		}
		if !putOTM.IsOTM() {
			t.Error("IsOTM() returned false for OTM put")
		}
	})

	t.Run("Value calculations", func(t *testing.T) {
		callOption := Option{
			PCDiv:   "2",
			Strike:  100.0,
			UnderPx: floatPtr(110.0),
			Theo:    floatPtr(15.0),
		}

		// Test moneyness
		if moneyness := callOption.GetMoneyness(); moneyness == nil || *moneyness != 1.1 {
			t.Errorf("GetMoneyness() = %v, want 1.1", moneyness)
		}

		// Test intrinsic value
		if intrinsic := callOption.GetIntrinsicValue(); intrinsic != 10.0 {
			t.Errorf("GetIntrinsicValue() = %v, want 10.0", intrinsic)
		}

		// Test time value
		if timeValue := callOption.GetTimeValue(); timeValue == nil || *timeValue != 5.0 {
			t.Errorf("GetTimeValue() = %v, want 5.0", timeValue)
		}
	})
}

func TestOptionsResponse_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"data": [
			{
				"Code": "140014505",
				"ProdCat": "TOPIXE",
				"UndSSO": "-",
				"Date": "2024-07-23",
				"CM": "2025-01",
				"Strike": 2450.0,
				"PCDiv": "2",
				"EmMrgnTrgDiv": "002",
				"O": 0.0,
				"H": 0.0,
				"L": 0.0,
				"C": 0.0,
				"EO": "",
				"EH": "",
				"EL": "",
				"EC": "",
				"AO": 0.0,
				"AH": 0.0,
				"AL": 0.0,
				"AC": 0.0,
				"MO": "",
				"MH": "",
				"ML": "",
				"MC": "",
				"Vo": 0.0,
				"OI": 0.0,
				"Va": 0.0,
				"VoOA": 0.0,
				"Settle": 377.0,
				"Theo": 380.3801,
				"BaseVol": 18.115,
				"UnderPx": 2833.39,
				"IV": 17.2955,
				"IR": 0.3527,
				"LTD": "2025-01-09",
				"SQD": "2025-01-10",
				"CCMFlag": "0"
			}
		],
		"pagination_key": "value1.value2."
	}`

	var resp OptionsResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if len(resp.Data) != 1 {
		t.Errorf("UnmarshalJSON() options count = %d, want 1", len(resp.Data))
	}

	opt := resp.Data[0]
	if opt.Code != "140014505" {
		t.Errorf("UnmarshalJSON() Code = %v, want 140014505", opt.Code)
	}
	if opt.Strike != 2450.0 {
		t.Errorf("UnmarshalJSON() Strike = %v, want 2450.0", opt.Strike)
	}

	// Check empty string fields are handled correctly
	if str, ok := opt.EO.(string); !ok || str != "" {
		t.Errorf("UnmarshalJSON() EO = %v, want empty string", opt.EO)
	}

	// Check optional fields
	if opt.Settle == nil || *opt.Settle != 377.0 {
		t.Errorf("UnmarshalJSON() Settle = %v, want 377.0", opt.Settle)
	}
	if opt.LTD == nil || *opt.LTD != "2025-01-09" {
		t.Errorf("UnmarshalJSON() LTD = %v, want 2025-01-09", opt.LTD)
	}

	if resp.PaginationKey != "value1.value2." {
		t.Errorf("UnmarshalJSON() PaginationKey = %v, want value1.value2.", resp.PaginationKey)
	}
}
