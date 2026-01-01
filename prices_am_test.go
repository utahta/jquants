package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestPricesAMService_GetPricesAM(t *testing.T) {
	tests := []struct {
		name     string
		params   PricesAMParams
		wantPath string
	}{
		{
			name: "with code",
			params: PricesAMParams{
				Code: "39400",
			},
			wantPath: "/prices/prices_am?code=39400",
		},
		{
			name: "with code and pagination key",
			params: PricesAMParams{
				Code:          "27800",
				PaginationKey: "key123",
			},
			wantPath: "/prices/prices_am?code=27800&pagination_key=key123",
		},
		{
			name:     "with no parameters",
			params:   PricesAMParams{},
			wantPath: "/prices/prices_am",
		},
		{
			name: "with pagination key only",
			params: PricesAMParams{
				PaginationKey: "key456",
			},
			wantPath: "/prices/prices_am?pagination_key=key456",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewPricesAMService(mockClient)

			// Mock response based on documentation sample
			mockResponse := PricesAMResponse{
				PricesAM: []PriceAM{
					{
						Date:                 "2023-03-20",
						Code:                 "39400",
						MorningOpen:          floatPtr(232.0),
						MorningHigh:          floatPtr(244.0),
						MorningLow:           floatPtr(232.0),
						MorningClose:         floatPtr(240.0),
						MorningVolume:        floatPtr(52600.0),
						MorningTurnoverValue: floatPtr(12518800.0),
					},
				},
				PaginationKey: "",
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Execute
			resp, err := service.GetPricesAM(tt.params)

			// Verify
			if err != nil {
				t.Fatalf("GetPricesAM() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetPricesAM() returned nil response")
				return
			}
			if len(resp.PricesAM) == 0 {
				t.Error("GetPricesAM() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetPricesAM() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestPricesAMService_GetPricesAMByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewPricesAMService(mockClient)

	// Mock response
	mockResponse := PricesAMResponse{
		PricesAM: []PriceAM{
			{
				Date:                 "2023-03-20",
				Code:                 "39400",
				MorningOpen:          floatPtr(232.0),
				MorningHigh:          floatPtr(244.0),
				MorningLow:           floatPtr(232.0),
				MorningClose:         floatPtr(240.0),
				MorningVolume:        floatPtr(52600.0),
				MorningTurnoverValue: floatPtr(12518800.0),
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/prices/prices_am?code=39400", mockResponse)

	// Execute
	resp, err := service.GetPricesAMByCode("39400")

	// Verify
	if err != nil {
		t.Fatalf("GetPricesAMByCode() error = %v", err)
	}
	if len(resp.PricesAM) != 1 {
		t.Errorf("GetPricesAMByCode() returned %d items, want 1", len(resp.PricesAM))
	}
	if resp.PricesAM[0].Code != "39400" {
		t.Errorf("GetPricesAMByCode() returned code %v, want 39400", resp.PricesAM[0].Code)
	}
}

func TestPricesAMService_GetAllPricesAM(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewPricesAMService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := PricesAMResponse{
		PricesAM: []PriceAM{
			{
				Date:         "2023-03-20",
				Code:         "13010",
				MorningClose: floatPtr(2000.0),
			},
			{
				Date:         "2023-03-20",
				Code:         "13020",
				MorningClose: floatPtr(1500.0),
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := PricesAMResponse{
		PricesAM: []PriceAM{
			{
				Date:         "2023-03-20",
				Code:         "13030",
				MorningClose: floatPtr(1800.0),
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/prices/prices_am", mockResponse1)
	mockClient.SetResponse("GET", "/prices/prices_am?pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetAllPricesAM()

	// Verify
	if err != nil {
		t.Fatalf("GetAllPricesAM() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetAllPricesAM() returned %d items, want 3", len(data))
	}
}

func TestPriceAM_GetMorningRange(t *testing.T) {
	tests := []struct {
		name  string
		price PriceAM
		want  *float64
	}{
		{
			name: "normal case",
			price: PriceAM{
				MorningHigh: floatPtr(244.0),
				MorningLow:  floatPtr(232.0),
			},
			want: floatPtr(12.0),
		},
		{
			name: "missing high",
			price: PriceAM{
				MorningLow: floatPtr(232.0),
			},
			want: nil,
		},
		{
			name: "missing low",
			price: PriceAM{
				MorningHigh: floatPtr(244.0),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.price.GetMorningRange()
			if !compareFloat64Ptr(got, tt.want) {
				t.Errorf("GetMorningRange() = %v, want %v", ptrToStr(got), ptrToStr(tt.want))
			}
		})
	}
}

func TestPriceAM_GetMorningChangeFromOpen(t *testing.T) {
	tests := []struct {
		name  string
		price PriceAM
		want  *float64
	}{
		{
			name: "positive change",
			price: PriceAM{
				MorningOpen:  floatPtr(232.0),
				MorningClose: floatPtr(240.0),
			},
			want: floatPtr(8.0),
		},
		{
			name: "negative change",
			price: PriceAM{
				MorningOpen:  floatPtr(240.0),
				MorningClose: floatPtr(232.0),
			},
			want: floatPtr(-8.0),
		},
		{
			name: "missing open",
			price: PriceAM{
				MorningClose: floatPtr(240.0),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.price.GetMorningChangeFromOpen()
			if !compareFloat64Ptr(got, tt.want) {
				t.Errorf("GetMorningChangeFromOpen() = %v, want %v", ptrToStr(got), ptrToStr(tt.want))
			}
		})
	}
}

func TestPriceAM_GetMorningChangeRate(t *testing.T) {
	tests := []struct {
		name  string
		price PriceAM
		want  *float64
	}{
		{
			name: "positive rate",
			price: PriceAM{
				MorningOpen:  floatPtr(232.0),
				MorningClose: floatPtr(240.0),
			},
			want: floatPtr((8.0 / 232.0) * 100),
		},
		{
			name: "negative rate",
			price: PriceAM{
				MorningOpen:  floatPtr(240.0),
				MorningClose: floatPtr(232.0),
			},
			want: floatPtr((-8.0 / 240.0) * 100),
		},
		{
			name: "zero open",
			price: PriceAM{
				MorningOpen:  floatPtr(0.0),
				MorningClose: floatPtr(240.0),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.price.GetMorningChangeRate()
			if !compareFloat64Ptr(got, tt.want) {
				t.Errorf("GetMorningChangeRate() = %v, want %v", ptrToStr(got), ptrToStr(tt.want))
			}
		})
	}
}

func TestPriceAM_HasMorningTrade(t *testing.T) {
	tests := []struct {
		name  string
		price PriceAM
		want  bool
	}{
		{
			name: "has trade",
			price: PriceAM{
				MorningVolume: floatPtr(52600.0),
			},
			want: true,
		},
		{
			name: "no trade",
			price: PriceAM{
				MorningVolume: floatPtr(0.0),
			},
			want: false,
		},
		{
			name: "null volume",
			price: PriceAM{
				MorningVolume: nil,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price.HasMorningTrade(); got != tt.want {
				t.Errorf("HasMorningTrade() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceAM_IsActiveTrading(t *testing.T) {
	tests := []struct {
		name  string
		price PriceAM
		want  bool
	}{
		{
			name: "active trading",
			price: PriceAM{
				MorningTurnoverValue: floatPtr(150000000.0), // 1.5億円
			},
			want: true,
		},
		{
			name: "exactly 100M",
			price: PriceAM{
				MorningTurnoverValue: floatPtr(100000000.0), // 1億円
			},
			want: true,
		},
		{
			name: "below threshold",
			price: PriceAM{
				MorningTurnoverValue: floatPtr(50000000.0), // 5千万円
			},
			want: false,
		},
		{
			name: "null turnover",
			price: PriceAM{
				MorningTurnoverValue: nil,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price.IsActiveTrading(); got != tt.want {
				t.Errorf("IsActiveTrading() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPriceAM_GetAveragePrice(t *testing.T) {
	tests := []struct {
		name  string
		price PriceAM
		want  *float64
	}{
		{
			name: "normal case",
			price: PriceAM{
				MorningTurnoverValue: floatPtr(12518800.0),
				MorningVolume:        floatPtr(52600.0),
			},
			want: floatPtr(238.0),
		},
		{
			name: "zero volume",
			price: PriceAM{
				MorningTurnoverValue: floatPtr(12518800.0),
				MorningVolume:        floatPtr(0.0),
			},
			want: nil,
		},
		{
			name: "null turnover",
			price: PriceAM{
				MorningTurnoverValue: nil,
				MorningVolume:        floatPtr(52600.0),
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.price.GetAveragePrice()
			if !compareFloat64Ptr(got, tt.want) {
				t.Errorf("GetAveragePrice() = %v, want %v", ptrToStr(got), ptrToStr(tt.want))
			}
		})
	}
}

func TestPriceAM_IsUpperLimit(t *testing.T) {
	tests := []struct {
		name  string
		price PriceAM
		want  bool
	}{
		{
			name: "upper limit",
			price: PriceAM{
				MorningOpen:  floatPtr(1000.0),
				MorningHigh:  floatPtr(1000.0),
				MorningLow:   floatPtr(1000.0),
				MorningClose: floatPtr(1000.0),
			},
			want: true,
		},
		{
			name: "not upper limit",
			price: PriceAM{
				MorningOpen:  floatPtr(232.0),
				MorningHigh:  floatPtr(244.0),
				MorningLow:   floatPtr(232.0),
				MorningClose: floatPtr(240.0),
			},
			want: false,
		},
		{
			name: "missing data",
			price: PriceAM{
				MorningOpen:  floatPtr(1000.0),
				MorningHigh:  floatPtr(1000.0),
				MorningLow:   floatPtr(1000.0),
				MorningClose: nil,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.price.IsUpperLimit(); got != tt.want {
				t.Errorf("IsUpperLimit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPricesAMService_GetPricesAM_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewPricesAMService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/prices/prices_am", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetPricesAM(PricesAMParams{})

	// Verify
	if err == nil {
		t.Error("GetPricesAM() expected error but got nil")
	}
}

// Note: Helper functions are now in test_helpers.go
