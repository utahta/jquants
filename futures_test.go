package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestFuturesService_GetFutures(t *testing.T) {
	tests := []struct {
		name     string
		params   FuturesParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with all parameters",
			params: FuturesParams{
				Date:          "20240723",
				Category:      "TOPIXF",
				ContractFlag:  "1",
				PaginationKey: "key123",
			},
			wantPath: "/derivatives/bars/daily/futures?date=20240723&category=TOPIXF&contract_flag=1&pagination_key=key123",
		},
		{
			name: "with date and category",
			params: FuturesParams{
				Date:     "20240723",
				Category: "NK225F",
			},
			wantPath: "/derivatives/bars/daily/futures?date=20240723&category=NK225F",
		},
		{
			name: "with date only",
			params: FuturesParams{
				Date: "2024-07-23",
			},
			wantPath: "/derivatives/bars/daily/futures?date=2024-07-23",
		},
		{
			name: "with date and central contract flag",
			params: FuturesParams{
				Date:         "20240723",
				ContractFlag: "1",
			},
			wantPath: "/derivatives/bars/daily/futures?date=20240723&contract_flag=1",
		},
		{
			name:     "without date (should error)",
			params:   FuturesParams{},
			wantPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewFuturesService(mockClient)

			// Mock response based on documentation sample
			mockResponse := FuturesResponse{
				Data: []Futures{
					{
						Code:         "169090005",
						ProdCat:      "TOPIXF",
						Date:         "2024-07-23",
						CM:           "2024-09",
						EmMrgnTrgDiv: "002",
						O:            2825.5,
						H:            2853.0,
						L:            2825.5,
						C:            2829.0,
						EO:           2825.5,
						EH:           2850.0,
						EL:           2825.5,
						EC:           2845.0,
						AO:           2850.5,
						AH:           2853.0,
						AL:           2826.0,
						AC:           2829.0,
						MO:           "",
						MH:           "",
						ML:           "",
						MC:           "",
						Vo:           42910.0,
						OI:           479812.0,
						Va:           1217918971856.0,
						VoOA:         floatPtr(40405.0),
						Settle:       floatPtr(2829.0),
						LTD:          stringPtr("2024-09-12"),
						SQD:          stringPtr("2024-09-13"),
						CCMFlag:      stringPtr("1"),
					},
				},
				PaginationKey: "",
			}
			if !tt.wantErr {
				mockClient.SetResponse("GET", tt.wantPath, mockResponse)
			}

			// Execute
			resp, err := service.GetFutures(tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetFutures() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetFutures() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetFutures() returned nil response")
				return
			}
			if len(resp.Data) == 0 {
				t.Error("GetFutures() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetFutures() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestFuturesService_GetFutures_RequiresDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewFuturesService(mockClient)

	// Execute with empty date
	_, err := service.GetFutures(FuturesParams{})

	// Verify
	if err == nil {
		t.Error("GetFutures() expected error for missing date but got nil")
	}
	if err.Error() != "date parameter is required" {
		t.Errorf("GetFutures() error = %v, want 'date parameter is required'", err)
	}
}

func TestFuturesService_GetFuturesByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewFuturesService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := FuturesResponse{
		Data: []Futures{
			{
				Code:    "169090005",
				ProdCat: "TOPIXF",
				Date:    "2024-07-23",
				CM:      "2024-09",
			},
			{
				Code:    "169120005",
				ProdCat: "TOPIXF",
				Date:    "2024-07-23",
				CM:      "2024-12",
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := FuturesResponse{
		Data: []Futures{
			{
				Code:    "167090018",
				ProdCat: "NK225F",
				Date:    "2024-07-23",
				CM:      "2024-09",
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/derivatives/bars/daily/futures?date=20240723", mockResponse1)
	mockClient.SetResponse("GET", "/derivatives/bars/daily/futures?date=20240723&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetFuturesByDate("20240723")

	// Verify
	if err != nil {
		t.Fatalf("GetFuturesByDate() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetFuturesByDate() returned %d items, want 3", len(data))
	}
	for _, item := range data {
		if item.Date != "2024-07-23" {
			t.Errorf("GetFuturesByDate() returned date %v, want 2024-07-23", item.Date)
		}
	}
}

func TestFuturesService_GetFuturesByCategory(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewFuturesService(mockClient)

	// Mock response
	mockResponse := FuturesResponse{
		Data: []Futures{
			{
				Code:    "167090018",
				ProdCat: "NK225F",
				Date:    "2024-07-23",
				CM:      "2024-09",
			},
			{
				Code:    "167120018",
				ProdCat: "NK225F",
				Date:    "2024-07-23",
				CM:      "2024-12",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/derivatives/bars/daily/futures?date=20240723&category=NK225F", mockResponse)

	// Execute
	data, err := service.GetFuturesByCategory("20240723", "NK225F")

	// Verify
	if err != nil {
		t.Fatalf("GetFuturesByCategory() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetFuturesByCategory() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.ProdCat != "NK225F" {
			t.Errorf("GetFuturesByCategory() returned category %v, want NK225F", item.ProdCat)
		}
	}
}

func TestFuturesService_GetCentralContractMonthFutures(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewFuturesService(mockClient)

	// Mock response
	mockResponse := FuturesResponse{
		Data: []Futures{
			{
				Code:    "167090018",
				ProdCat: "NK225F",
				Date:    "2024-07-23",
				CM:      "2024-09",
				CCMFlag: stringPtr("1"),
			},
			{
				Code:    "169090005",
				ProdCat: "TOPIXF",
				Date:    "2024-07-23",
				CM:      "2024-09",
				CCMFlag: stringPtr("1"),
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/derivatives/bars/daily/futures?date=20240723&contract_flag=1", mockResponse)

	// Execute
	data, err := service.GetCentralContractMonthFutures("20240723")

	// Verify
	if err != nil {
		t.Fatalf("GetCentralContractMonthFutures() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetCentralContractMonthFutures() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if !item.IsCentralContractMonth() {
			t.Errorf("GetCentralContractMonthFutures() returned non-central contract month")
		}
	}
}

func TestFutures_HelperMethods(t *testing.T) {
	tests := []struct {
		name                       string
		futures                    Futures
		isEmergencyMarginTriggered bool
		isCentralContractMonth     bool
		hasNightSession            bool
		hasMorningSession          bool
	}{
		{
			name: "emergency margin triggered",
			futures: Futures{
				EmMrgnTrgDiv: "001",
			},
			isEmergencyMarginTriggered: true,
			isCentralContractMonth:     false,
			hasNightSession:            true,
			hasMorningSession:          true,
		},
		{
			name: "normal settlement price",
			futures: Futures{
				EmMrgnTrgDiv: "002",
				CCMFlag:      stringPtr("1"),
				EO:           2825.5,
				MO:           "",
			},
			isEmergencyMarginTriggered: false,
			isCentralContractMonth:     true,
			hasNightSession:            true,
			hasMorningSession:          false,
		},
		{
			name: "no night session on first day",
			futures: Futures{
				EmMrgnTrgDiv: "002",
				CCMFlag:      stringPtr("0"),
				EO:           "",
				MO:           2500.0,
			},
			isEmergencyMarginTriggered: false,
			isCentralContractMonth:     false,
			hasNightSession:            false,
			hasMorningSession:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.futures.IsEmergencyMarginTriggered(); got != tt.isEmergencyMarginTriggered {
				t.Errorf("IsEmergencyMarginTriggered() = %v, want %v", got, tt.isEmergencyMarginTriggered)
			}
			if got := tt.futures.IsCentralContractMonth(); got != tt.isCentralContractMonth {
				t.Errorf("IsCentralContractMonth() = %v, want %v", got, tt.isCentralContractMonth)
			}
			if got := tt.futures.HasNightSession(); got != tt.hasNightSession {
				t.Errorf("HasNightSession() = %v, want %v", got, tt.hasNightSession)
			}
			if got := tt.futures.HasMorningSession(); got != tt.hasMorningSession {
				t.Errorf("HasMorningSession() = %v, want %v", got, tt.hasMorningSession)
			}
		})
	}
}

func TestFutures_SessionGetters(t *testing.T) {
	// Setup
	futures := Futures{
		EO: 2825.5,
		EH: 2850.0,
		EL: 2825.5,
		EC: 2845.0,
		MO: "",
		MH: "",
		ML: "",
		MC: "",
	}

	// Test night session getters
	if open := futures.GetNightSessionOpen(); open == nil || *open != 2825.5 {
		t.Errorf("GetNightSessionOpen() = %v, want 2825.5", open)
	}
	if high := futures.GetNightSessionHigh(); high == nil || *high != 2850.0 {
		t.Errorf("GetNightSessionHigh() = %v, want 2850.0", high)
	}
	if low := futures.GetNightSessionLow(); low == nil || *low != 2825.5 {
		t.Errorf("GetNightSessionLow() = %v, want 2825.5", low)
	}
	if close := futures.GetNightSessionClose(); close == nil || *close != 2845.0 {
		t.Errorf("GetNightSessionClose() = %v, want 2845.0", close)
	}

	// Test morning session getters (should return nil for empty strings)
	if open := futures.GetMorningSessionOpen(); open != nil {
		t.Errorf("GetMorningSessionOpen() = %v, want nil", open)
	}
	if high := futures.GetMorningSessionHigh(); high != nil {
		t.Errorf("GetMorningSessionHigh() = %v, want nil", high)
	}
	if low := futures.GetMorningSessionLow(); low != nil {
		t.Errorf("GetMorningSessionLow() = %v, want nil", low)
	}
	if close := futures.GetMorningSessionClose(); close != nil {
		t.Errorf("GetMorningSessionClose() = %v, want nil", close)
	}
}

func TestFutures_GetDayNightGap(t *testing.T) {
	tests := []struct {
		name    string
		futures Futures
		want    *float64
	}{
		{
			name: "positive gap",
			futures: Futures{
				EC: 2845.0,
				AO: 2850.5,
			},
			want: floatPtr(5.5),
		},
		{
			name: "negative gap",
			futures: Futures{
				EC: 2850.0,
				AO: 2845.0,
			},
			want: floatPtr(-5.0),
		},
		{
			name: "no night session",
			futures: Futures{
				EC: "",
				AO: 2850.5,
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.futures.GetDayNightGap()
			if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) {
				t.Errorf("GetDayNightGap() = %v, want %v", got, tt.want)
			} else if got != nil && tt.want != nil && *got != *tt.want {
				t.Errorf("GetDayNightGap() = %v, want %v", *got, *tt.want)
			}
		})
	}
}

func TestFutures_GetWholeDayRange(t *testing.T) {
	// Setup
	futures := Futures{
		H: 2853.0,
		L: 2825.5,
	}

	// Execute
	got := futures.GetWholeDayRange()

	// Verify
	want := 27.5
	if got != want {
		t.Errorf("GetWholeDayRange() = %v, want %v", got, want)
	}
}

func TestFuturesResponse_UnmarshalJSON(t *testing.T) {
	// Test JSON with mixed types
	jsonData := `{
		"data": [
			{
				"Code": "169090005",
				"ProdCat": "TOPIXF",
				"Date": "2024-07-23",
				"CM": "2024-09",
				"EmMrgnTrgDiv": "002",
				"O": 2825.5,
				"H": 2853.0,
				"L": 2825.5,
				"C": 2829.0,
				"EO": 2825.5,
				"EH": 2850.0,
				"EL": 2825.5,
				"EC": 2845.0,
				"AO": 2850.5,
				"AH": 2853.0,
				"AL": 2826.0,
				"AC": 2829.0,
				"MO": "",
				"MH": "",
				"ML": "",
				"MC": "",
				"Vo": 42910.0,
				"OI": 479812.0,
				"Va": 1217918971856.0,
				"VoOA": 40405.0,
				"Settle": 2829.0,
				"LTD": "2024-09-12",
				"SQD": "2024-09-13",
				"CCMFlag": "1"
			}
		],
		"pagination_key": "next_key"
	}`

	// Unmarshal
	var resp FuturesResponse
	err := resp.UnmarshalJSON([]byte(jsonData))

	// Verify
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("UnmarshalJSON() returned %d items, want 1", len(resp.Data))
	}

	f := resp.Data[0]
	if f.Code != "169090005" {
		t.Errorf("Code = %v, want 169090005", f.Code)
	}
	if f.O != 2825.5 {
		t.Errorf("O = %v, want 2825.5", f.O)
	}
	if !f.HasNightSession() {
		t.Error("HasNightSession() = false, want true")
	}
	if f.HasMorningSession() {
		t.Error("HasMorningSession() = true, want false")
	}
	if f.VoOA == nil || *f.VoOA != 40405.0 {
		t.Errorf("VoOA = %v, want 40405.0", f.VoOA)
	}
	if f.CCMFlag == nil || *f.CCMFlag != "1" {
		t.Errorf("CCMFlag = %v, want 1", f.CCMFlag)
	}
}

func TestFuturesService_GetFutures_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewFuturesService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/derivatives/bars/daily/futures?date=20240723", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetFutures(FuturesParams{Date: "20240723"})

	// Verify
	if err == nil {
		t.Error("GetFutures() expected error but got nil")
	}
}
