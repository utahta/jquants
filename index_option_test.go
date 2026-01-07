package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestIndexOptionService_GetIndexOptions(t *testing.T) {
	tests := []struct {
		name     string
		params   IndexOptionParams
		wantPath string
	}{
		{
			name: "with date and pagination key",
			params: IndexOptionParams{
				Date:          "20230322",
				PaginationKey: "key123",
			},
			wantPath: "/derivatives/bars/daily/options/225?date=20230322&pagination_key=key123",
		},
		{
			name: "with date only",
			params: IndexOptionParams{
				Date: "20230322",
			},
			wantPath: "/derivatives/bars/daily/options/225?date=20230322",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewIndexOptionService(mockClient)

			// Mock response based on documentation sample
			mockResponse := IndexOptionResponse{
				Data: []IndexOption{
					{
						Date:         "2023-03-22",
						Code:         "130060018",
						CM:           "2025-06",
						Strike:       20000.0,
						PCDiv:        "1",
						LTD:          "2025-06-12",
						SQD:          "2025-06-13",
						EmMrgnTrgDiv: "002",
						O:            0.0,
						H:            0.0,
						L:            0.0,
						C:            0.0,
						EO:           floatPtr(0.0),
						EH:           floatPtr(0.0),
						EL:           floatPtr(0.0),
						EC:           floatPtr(0.0),
						AO:           0.0,
						AH:           0.0,
						AL:           0.0,
						AC:           0.0,
						Vo:           0.0,
						VoOA:         floatPtr(0.0),
						OI:           330.0,
						Va:           0.0,
						Settle:       floatPtr(980.0),
						Theo:         floatPtr(974.641),
						BaseVol:      floatPtr(17.93025),
						UnderPx:      floatPtr(27466.61),
						IV:           floatPtr(23.1816),
						IR:           floatPtr(0.2336),
					},
				},
				PaginationKey: "",
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Execute
			resp, err := service.GetIndexOptions(tt.params)

			// Verify
			if err != nil {
				t.Fatalf("GetIndexOptions() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetIndexOptions() returned nil response")
				return
			}
			if len(resp.Data) == 0 {
				t.Error("GetIndexOptions() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetIndexOptions() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestIndexOptionService_GetIndexOptions_RequiresDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndexOptionService(mockClient)

	// Execute with empty date
	_, err := service.GetIndexOptions(IndexOptionParams{})

	// Verify
	if err == nil {
		t.Error("GetIndexOptions() expected error for missing date but got nil")
	}
	if err.Error() != "date parameter is required" {
		t.Errorf("GetIndexOptions() error = %v, want 'date parameter is required'", err)
	}
}

func TestIndexOptionService_GetIndexOptionsByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndexOptionService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := IndexOptionResponse{
		Data: []IndexOption{
			{
				Date:         "2023-03-22",
				Code:         "130060018",
				PCDiv:        "1",
				Strike:       20000.0,
				EmMrgnTrgDiv: "002",
			},
			{
				Date:         "2023-03-22",
				Code:         "130060019",
				PCDiv:        "2",
				Strike:       21000.0,
				EmMrgnTrgDiv: "002",
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := IndexOptionResponse{
		Data: []IndexOption{
			{
				Date:         "2023-03-22",
				Code:         "130060020",
				PCDiv:        "1",
				Strike:       22000.0,
				EmMrgnTrgDiv: "002",
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/derivatives/bars/daily/options/225?date=20230322", mockResponse1)
	mockClient.SetResponse("GET", "/derivatives/bars/daily/options/225?date=20230322&pagination_key=next_page_key", mockResponse2)

	// Execute
	options, err := service.GetIndexOptionsByDate("20230322")

	// Verify
	if err != nil {
		t.Fatalf("GetIndexOptionsByDate() error = %v", err)
	}
	if len(options) != 3 {
		t.Errorf("GetIndexOptionsByDate() returned %d options, want 3", len(options))
	}
}

func TestIndexOptionService_GetCallOptions(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndexOptionService(mockClient)

	// Mock response with both call and put options
	mockResponse := IndexOptionResponse{
		Data: []IndexOption{
			{
				Date:         "2023-03-22",
				Code:         "130060018",
				PCDiv:        "1", // Put
				Strike:       20000.0,
				EmMrgnTrgDiv: "002",
			},
			{
				Date:         "2023-03-22",
				Code:         "130060019",
				PCDiv:        "2", // Call
				Strike:       21000.0,
				EmMrgnTrgDiv: "002",
			},
			{
				Date:         "2023-03-22",
				Code:         "130060020",
				PCDiv:        "2", // Call
				Strike:       22000.0,
				EmMrgnTrgDiv: "002",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/derivatives/bars/daily/options/225?date=20230322", mockResponse)

	// Execute
	callOptions, err := service.GetCallOptions("20230322")

	// Verify
	if err != nil {
		t.Fatalf("GetCallOptions() error = %v", err)
	}
	if len(callOptions) != 2 {
		t.Errorf("GetCallOptions() returned %d options, want 2", len(callOptions))
	}
	for _, option := range callOptions {
		if option.PCDiv != PutCallDivisionCall {
			t.Errorf("GetCallOptions() returned option with PCDiv %v, want %v",
				option.PCDiv, PutCallDivisionCall)
		}
	}
}

func TestIndexOptionService_GetPutOptions(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndexOptionService(mockClient)

	// Mock response with both call and put options
	mockResponse := IndexOptionResponse{
		Data: []IndexOption{
			{
				Date:         "2023-03-22",
				Code:         "130060018",
				PCDiv:        "1", // Put
				Strike:       20000.0,
				EmMrgnTrgDiv: "002",
			},
			{
				Date:         "2023-03-22",
				Code:         "130060019",
				PCDiv:        "2", // Call
				Strike:       21000.0,
				EmMrgnTrgDiv: "002",
			},
			{
				Date:         "2023-03-22",
				Code:         "130060020",
				PCDiv:        "1", // Put
				Strike:       22000.0,
				EmMrgnTrgDiv: "002",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/derivatives/bars/daily/options/225?date=20230322", mockResponse)

	// Execute
	putOptions, err := service.GetPutOptions("20230322")

	// Verify
	if err != nil {
		t.Fatalf("GetPutOptions() error = %v", err)
	}
	if len(putOptions) != 2 {
		t.Errorf("GetPutOptions() returned %d options, want 2", len(putOptions))
	}
	for _, option := range putOptions {
		if option.PCDiv != PutCallDivisionPut {
			t.Errorf("GetPutOptions() returned option with PCDiv %v, want %v",
				option.PCDiv, PutCallDivisionPut)
		}
	}
}

func TestIndexOptionService_GetOptionChain(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndexOptionService(mockClient)

	// Mock response
	mockResponse := IndexOptionResponse{
		Data: []IndexOption{
			{
				Date:         "2023-03-22",
				Code:         "130060018",
				PCDiv:        "1",
				Strike:       20000.0,
				EmMrgnTrgDiv: "002",
			},
			{
				Date:         "2023-03-22",
				Code:         "130060019",
				PCDiv:        "2",
				Strike:       20000.0,
				EmMrgnTrgDiv: "002",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/derivatives/bars/daily/options/225?date=20230322", mockResponse)

	// Execute
	options, err := service.GetOptionChain("20230322")

	// Verify
	if err != nil {
		t.Fatalf("GetOptionChain() error = %v", err)
	}
	if len(options) != 2 {
		t.Errorf("GetOptionChain() returned %d options, want 2", len(options))
	}
}

func TestIndexOption_IsCall(t *testing.T) {
	tests := []struct {
		name   string
		option IndexOption
		want   bool
	}{
		{
			name: "call option",
			option: IndexOption{
				PCDiv: PutCallDivisionCall,
			},
			want: true,
		},
		{
			name: "put option",
			option: IndexOption{
				PCDiv: PutCallDivisionPut,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.option.IsCall()
			if got != tt.want {
				t.Errorf("IndexOption.IsCall() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOption_IsPut(t *testing.T) {
	tests := []struct {
		name   string
		option IndexOption
		want   bool
	}{
		{
			name: "put option",
			option: IndexOption{
				PCDiv: PutCallDivisionPut,
			},
			want: true,
		},
		{
			name: "call option",
			option: IndexOption{
				PCDiv: PutCallDivisionCall,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.option.IsPut()
			if got != tt.want {
				t.Errorf("IndexOption.IsPut() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOption_IsEmergencyMarginTriggered(t *testing.T) {
	tests := []struct {
		name   string
		option IndexOption
		want   bool
	}{
		{
			name: "emergency margin triggered",
			option: IndexOption{
				EmMrgnTrgDiv: EmergencyMarginTriggerDivisionEmergency,
			},
			want: true,
		},
		{
			name: "normal settlement",
			option: IndexOption{
				EmMrgnTrgDiv: EmergencyMarginTriggerDivisionNormal,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.option.IsEmergencyMarginTriggered()
			if got != tt.want {
				t.Errorf("IndexOption.IsEmergencyMarginTriggered() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOption_HasNightSession(t *testing.T) {
	tests := []struct {
		name   string
		option IndexOption
		want   bool
	}{
		{
			name: "has night session data",
			option: IndexOption{
				EO: floatPtr(100.0),
			},
			want: true,
		},
		{
			name: "no night session data",
			option: IndexOption{
				EO: nil,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.option.HasNightSession()
			if got != tt.want {
				t.Errorf("IndexOption.HasNightSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIndexOptionService_GetIndexOptions_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndexOptionService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/derivatives/bars/daily/options/225?date=20230322", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetIndexOptions(IndexOptionParams{Date: "20230322"})

	// Verify
	if err == nil {
		t.Error("GetIndexOptions() expected error but got nil")
	}
}

func TestPutCallDivisionConstants(t *testing.T) {
	// プットコール区分定数の値を確認
	tests := []struct {
		constant string
		expected string
	}{
		{PutCallDivisionPut, "1"},
		{PutCallDivisionCall, "2"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("PutCall division constant = %v, want %v", tt.constant, tt.expected)
			}
		})
	}
}

func TestEmergencyMarginTriggerDivisionConstants(t *testing.T) {
	// 緊急取引証拠金発動区分定数の値を確認
	tests := []struct {
		constant string
		expected string
	}{
		{EmergencyMarginTriggerDivisionEmergency, "001"},
		{EmergencyMarginTriggerDivisionNormal, "002"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Emergency margin trigger division constant = %v, want %v", tt.constant, tt.expected)
			}
		})
	}
}
