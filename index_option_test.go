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
			wantPath: "/option/index_option?date=20230322&pagination_key=key123",
		},
		{
			name: "with date only",
			params: IndexOptionParams{
				Date: "20230322",
			},
			wantPath: "/option/index_option?date=20230322",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewIndexOptionService(mockClient)

			// Mock response based on documentation sample
			mockResponse := IndexOptionResponse{
				IndexOptions: []IndexOption{
					{
						Date:                           "2023-03-22",
						Code:                           "130060018",
						ContractMonth:                  "2025-06",
						StrikePrice:                    20000.0,
						PutCallDivision:                "1",
						LastTradingDay:                 "2025-06-12",
						SpecialQuotationDay:            "2025-06-13",
						EmergencyMarginTriggerDivision: "002",
						WholeDayOpen:                   0.0,
						WholeDayHigh:                   0.0,
						WholeDayLow:                    0.0,
						WholeDayClose:                  0.0,
						NightSessionOpen:               floatPtr(0.0),
						NightSessionHigh:               floatPtr(0.0),
						NightSessionLow:                floatPtr(0.0),
						NightSessionClose:              floatPtr(0.0),
						DaySessionOpen:                 0.0,
						DaySessionHigh:                 0.0,
						DaySessionLow:                  0.0,
						DaySessionClose:                0.0,
						Volume:                         0.0,
						VolumeOnlyAuction:              floatPtr(0.0),
						OpenInterest:                   330.0,
						TurnoverValue:                  0.0,
						SettlementPrice:                floatPtr(980.0),
						TheoreticalPrice:               floatPtr(974.641),
						BaseVolatility:                 floatPtr(17.93025),
						UnderlyingPrice:                floatPtr(27466.61),
						ImpliedVolatility:              floatPtr(23.1816),
						InterestRate:                   floatPtr(0.2336),
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
			if len(resp.IndexOptions) == 0 {
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
		IndexOptions: []IndexOption{
			{
				Date:                           "2023-03-22",
				Code:                           "130060018",
				PutCallDivision:                "1",
				StrikePrice:                    20000.0,
				EmergencyMarginTriggerDivision: "002",
			},
			{
				Date:                           "2023-03-22",
				Code:                           "130060019",
				PutCallDivision:                "2",
				StrikePrice:                    21000.0,
				EmergencyMarginTriggerDivision: "002",
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := IndexOptionResponse{
		IndexOptions: []IndexOption{
			{
				Date:                           "2023-03-22",
				Code:                           "130060020",
				PutCallDivision:                "1",
				StrikePrice:                    22000.0,
				EmergencyMarginTriggerDivision: "002",
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/option/index_option?date=20230322", mockResponse1)
	mockClient.SetResponse("GET", "/option/index_option?date=20230322&pagination_key=next_page_key", mockResponse2)

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
		IndexOptions: []IndexOption{
			{
				Date:                           "2023-03-22",
				Code:                           "130060018",
				PutCallDivision:                "1", // Put
				StrikePrice:                    20000.0,
				EmergencyMarginTriggerDivision: "002",
			},
			{
				Date:                           "2023-03-22",
				Code:                           "130060019",
				PutCallDivision:                "2", // Call
				StrikePrice:                    21000.0,
				EmergencyMarginTriggerDivision: "002",
			},
			{
				Date:                           "2023-03-22",
				Code:                           "130060020",
				PutCallDivision:                "2", // Call
				StrikePrice:                    22000.0,
				EmergencyMarginTriggerDivision: "002",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/option/index_option?date=20230322", mockResponse)

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
		if option.PutCallDivision != PutCallDivisionCall {
			t.Errorf("GetCallOptions() returned option with PutCallDivision %v, want %v",
				option.PutCallDivision, PutCallDivisionCall)
		}
	}
}

func TestIndexOptionService_GetPutOptions(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndexOptionService(mockClient)

	// Mock response with both call and put options
	mockResponse := IndexOptionResponse{
		IndexOptions: []IndexOption{
			{
				Date:                           "2023-03-22",
				Code:                           "130060018",
				PutCallDivision:                "1", // Put
				StrikePrice:                    20000.0,
				EmergencyMarginTriggerDivision: "002",
			},
			{
				Date:                           "2023-03-22",
				Code:                           "130060019",
				PutCallDivision:                "2", // Call
				StrikePrice:                    21000.0,
				EmergencyMarginTriggerDivision: "002",
			},
			{
				Date:                           "2023-03-22",
				Code:                           "130060020",
				PutCallDivision:                "1", // Put
				StrikePrice:                    22000.0,
				EmergencyMarginTriggerDivision: "002",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/option/index_option?date=20230322", mockResponse)

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
		if option.PutCallDivision != PutCallDivisionPut {
			t.Errorf("GetPutOptions() returned option with PutCallDivision %v, want %v",
				option.PutCallDivision, PutCallDivisionPut)
		}
	}
}

func TestIndexOptionService_GetOptionChain(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndexOptionService(mockClient)

	// Mock response
	mockResponse := IndexOptionResponse{
		IndexOptions: []IndexOption{
			{
				Date:                           "2023-03-22",
				Code:                           "130060018",
				PutCallDivision:                "1",
				StrikePrice:                    20000.0,
				EmergencyMarginTriggerDivision: "002",
			},
			{
				Date:                           "2023-03-22",
				Code:                           "130060019",
				PutCallDivision:                "2",
				StrikePrice:                    20000.0,
				EmergencyMarginTriggerDivision: "002",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/option/index_option?date=20230322", mockResponse)

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
				PutCallDivision: PutCallDivisionCall,
			},
			want: true,
		},
		{
			name: "put option",
			option: IndexOption{
				PutCallDivision: PutCallDivisionPut,
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
				PutCallDivision: PutCallDivisionPut,
			},
			want: true,
		},
		{
			name: "call option",
			option: IndexOption{
				PutCallDivision: PutCallDivisionCall,
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
				EmergencyMarginTriggerDivision: EmergencyMarginTriggerDivisionEmergency,
			},
			want: true,
		},
		{
			name: "normal settlement",
			option: IndexOption{
				EmergencyMarginTriggerDivision: EmergencyMarginTriggerDivisionNormal,
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
				NightSessionOpen: floatPtr(100.0),
			},
			want: true,
		},
		{
			name: "no night session data",
			option: IndexOption{
				NightSessionOpen: nil,
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
	mockClient.SetError("GET", "/option/index_option?date=20230322", fmt.Errorf("unauthorized"))

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
