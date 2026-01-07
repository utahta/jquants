package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestShortSellingPositionsService_GetShortSellingPositions(t *testing.T) {
	tests := []struct {
		name     string
		params   ShortSellingPositionsParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with code and disclosed date",
			params: ShortSellingPositionsParams{
				Code:          "13660",
				DisclosedDate: "20240801",
			},
			wantPath: "/markets/short-sale-report?code=13660&disc_date=20240801",
		},
		{
			name: "with code only",
			params: ShortSellingPositionsParams{
				Code: "86970",
			},
			wantPath: "/markets/short-sale-report?code=86970",
		},
		{
			name: "with disclosed date only",
			params: ShortSellingPositionsParams{
				DisclosedDate: "2024-08-01",
			},
			wantPath: "/markets/short-sale-report?disc_date=2024-08-01",
		},
		{
			name: "with calculated date only",
			params: ShortSellingPositionsParams{
				CalculatedDate: "20240731",
			},
			wantPath: "/markets/short-sale-report?calc_date=20240731",
		},
		{
			name: "with code and date range",
			params: ShortSellingPositionsParams{
				Code:              "86970",
				DisclosedDateFrom: "20240101",
				DisclosedDateTo:   "20241231",
			},
			wantPath: "/markets/short-sale-report?code=86970&disc_date_from=20240101&disc_date_to=20241231",
		},
		{
			name: "with pagination key",
			params: ShortSellingPositionsParams{
				DisclosedDate: "20240801",
				PaginationKey: "key123",
			},
			wantPath: "/markets/short-sale-report?disc_date=20240801&pagination_key=key123",
		},
		{
			name:     "with no required parameters",
			params:   ShortSellingPositionsParams{},
			wantPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewShortSellingPositionsService(mockClient)

			// Mock response based on documentation sample
			mockResponse := ShortSellingPositionsResponse{
				Data: []ShortSellingPosition{
					{
						DiscDate:      "2024-08-01",
						CalcDate:      "2024-07-31",
						Code:          "13660",
						SSName:        "個人",
						SSAddr:        "",
						DICName:       "",
						DICAddr:       "",
						FundName:      "",
						ShrtPosToSO:   0.0053,
						ShrtPosShares: 140000,
						ShrtPosUnits:  140000,
						PrevRptDate:   "2024-07-22",
						PrevRptRatio:  0.0043,
						Notes:         "",
					},
				},
				PaginationKey: "",
			}
			if !tt.wantErr {
				mockClient.SetResponse("GET", tt.wantPath, mockResponse)
			}

			// Execute
			resp, err := service.GetShortSellingPositions(tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetShortSellingPositions() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetShortSellingPositions() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetShortSellingPositions() returned nil response")
				return
			}
			if len(resp.Data) == 0 {
				t.Error("GetShortSellingPositions() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetShortSellingPositions() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestShortSellingPositionsService_GetShortSellingPositions_RequiresParameter(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewShortSellingPositionsService(mockClient)

	// Execute with empty parameters
	_, err := service.GetShortSellingPositions(ShortSellingPositionsParams{})

	// Verify
	if err == nil {
		t.Error("GetShortSellingPositions() expected error for missing parameters but got nil")
	}
	if err.Error() != "either code, disc_date, or calc_date parameter is required" {
		t.Errorf("GetShortSellingPositions() error = %v, want 'either code, disc_date, or calc_date parameter is required'", err)
	}
}

func TestShortSellingPositionsService_GetShortSellingPositionsByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewShortSellingPositionsService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := ShortSellingPositionsResponse{
		Data: []ShortSellingPosition{
			{
				DiscDate:      "2024-08-01",
				CalcDate:      "2024-07-31",
				Code:          "86970",
				SSName:        "ABC Investment Management",
				SSAddr:        "123 Main St, New York",
				ShrtPosToSO:   0.0153,
				ShrtPosShares: 520000,
			},
			{
				DiscDate:      "2024-07-25",
				CalcDate:      "2024-07-24",
				Code:          "86970",
				SSName:        "ABC Investment Management",
				ShrtPosToSO:   0.0143,
				ShrtPosShares: 500000,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := ShortSellingPositionsResponse{
		Data: []ShortSellingPosition{
			{
				DiscDate:      "2024-07-18",
				CalcDate:      "2024-07-17",
				Code:          "86970",
				SSName:        "ABC Investment Management",
				ShrtPosToSO:   0.0133,
				ShrtPosShares: 480000,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/markets/short-sale-report?code=86970", mockResponse1)
	mockClient.SetResponse("GET", "/markets/short-sale-report?code=86970&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetShortSellingPositionsByCode("86970")

	// Verify
	if err != nil {
		t.Fatalf("GetShortSellingPositionsByCode() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetShortSellingPositionsByCode() returned %d items, want 3", len(data))
	}
	for _, item := range data {
		if item.Code != "86970" {
			t.Errorf("GetShortSellingPositionsByCode() returned code %v, want 86970", item.Code)
		}
	}
}

func TestShortSellingPositionsService_GetShortSellingPositionsByDisclosedDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewShortSellingPositionsService(mockClient)

	// Mock response
	mockResponse := ShortSellingPositionsResponse{
		Data: []ShortSellingPosition{
			{
				DiscDate:    "2024-08-01",
				CalcDate:    "2024-07-31",
				Code:        "13660",
				SSName:      "個人",
				ShrtPosToSO: 0.0053,
			},
			{
				DiscDate:    "2024-08-01",
				CalcDate:    "2024-07-31",
				Code:        "86970",
				SSName:      "XYZ Capital",
				ShrtPosToSO: 0.0087,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/markets/short-sale-report?disc_date=20240801", mockResponse)

	// Execute
	data, err := service.GetShortSellingPositionsByDisclosedDate("20240801")

	// Verify
	if err != nil {
		t.Fatalf("GetShortSellingPositionsByDisclosedDate() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetShortSellingPositionsByDisclosedDate() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.DiscDate != "2024-08-01" {
			t.Errorf("GetShortSellingPositionsByDisclosedDate() returned date %v, want 2024-08-01", item.DiscDate)
		}
	}
}

func TestShortSellingPositionsService_GetShortSellingPositionsByCalculatedDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewShortSellingPositionsService(mockClient)

	// Mock response
	mockResponse := ShortSellingPositionsResponse{
		Data: []ShortSellingPosition{
			{
				DiscDate:    "2024-08-01",
				CalcDate:    "2024-07-31",
				Code:        "13660",
				SSName:      "個人",
				ShrtPosToSO: 0.0053,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/markets/short-sale-report?calc_date=20240731", mockResponse)

	// Execute
	data, err := service.GetShortSellingPositionsByCalculatedDate("20240731")

	// Verify
	if err != nil {
		t.Fatalf("GetShortSellingPositionsByCalculatedDate() error = %v", err)
	}
	if len(data) != 1 {
		t.Errorf("GetShortSellingPositionsByCalculatedDate() returned %d items, want 1", len(data))
	}
	if data[0].CalcDate != "2024-07-31" {
		t.Errorf("GetShortSellingPositionsByCalculatedDate() returned date %v, want 2024-07-31", data[0].CalcDate)
	}
}

func TestShortSellingPositionsService_GetShortSellingPositionsByCodeAndDateRange(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewShortSellingPositionsService(mockClient)

	// Mock response
	mockResponse := ShortSellingPositionsResponse{
		Data: []ShortSellingPosition{
			{
				DiscDate:    "2024-03-01",
				Code:        "86970",
				ShrtPosToSO: 0.0087,
			},
			{
				DiscDate:    "2024-02-01",
				Code:        "86970",
				ShrtPosToSO: 0.0075,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/markets/short-sale-report?code=86970&disc_date_from=20240101&disc_date_to=20240331", mockResponse)

	// Execute
	data, err := service.GetShortSellingPositionsByCodeAndDateRange("86970", "20240101", "20240331")

	// Verify
	if err != nil {
		t.Fatalf("GetShortSellingPositionsByCodeAndDateRange() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetShortSellingPositionsByCodeAndDateRange() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.Code != "86970" {
			t.Errorf("GetShortSellingPositionsByCodeAndDateRange() returned code %v, want 86970", item.Code)
		}
	}
}

func TestShortSellingPosition_GetPositionChange(t *testing.T) {
	position := ShortSellingPosition{
		ShrtPosToSO:   0.0053,
		ShrtPosShares: 140000,
		PrevRptRatio:  0.0043,
	}

	change := position.GetPositionChange()

	// 概算値なので、おおよその値で確認
	if change < 25000 || change > 27000 {
		t.Errorf("ShortSellingPosition.GetPositionChange() = %v, want around 26400", change)
	}
}

func TestShortSellingPosition_GetPositionChangeRatio(t *testing.T) {
	position := ShortSellingPosition{
		ShrtPosToSO:  0.0053,
		PrevRptRatio: 0.0043,
	}

	expected := ((0.0053 - 0.0043) / 0.0043) * 100
	got := position.GetPositionChangeRatio()

	if got != expected {
		t.Errorf("ShortSellingPosition.GetPositionChangeRatio() = %v, want %v", got, expected)
	}
}

func TestShortSellingPosition_GetPositionChangeRatio_ZeroPrevious(t *testing.T) {
	position := ShortSellingPosition{
		ShrtPosToSO:  0.0053,
		PrevRptRatio: 0,
	}

	got := position.GetPositionChangeRatio()

	if got != 0 {
		t.Errorf("ShortSellingPosition.GetPositionChangeRatio() = %v, want 0", got)
	}
}

func TestShortSellingPosition_IsIncrease(t *testing.T) {
	tests := []struct {
		name     string
		position ShortSellingPosition
		want     bool
	}{
		{
			name: "increase",
			position: ShortSellingPosition{
				ShrtPosToSO:  0.0053,
				PrevRptRatio: 0.0043,
			},
			want: true,
		},
		{
			name: "decrease",
			position: ShortSellingPosition{
				ShrtPosToSO:  0.0043,
				PrevRptRatio: 0.0053,
			},
			want: false,
		},
		{
			name: "no change",
			position: ShortSellingPosition{
				ShrtPosToSO:  0.0053,
				PrevRptRatio: 0.0053,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.position.IsIncrease(); got != tt.want {
				t.Errorf("ShortSellingPosition.IsIncrease() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortSellingPosition_HasDiscretionaryInvestment(t *testing.T) {
	tests := []struct {
		name     string
		position ShortSellingPosition
		want     bool
	}{
		{
			name: "has discretionary investment",
			position: ShortSellingPosition{
				DICName: "ABC Management",
			},
			want: true,
		},
		{
			name: "has investment fund name",
			position: ShortSellingPosition{
				FundName: "XYZ Fund",
			},
			want: true,
		},
		{
			name: "no discretionary investment",
			position: ShortSellingPosition{
				DICName:  "",
				DICAddr:  "",
				FundName: "",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.position.HasDiscretionaryInvestment(); got != tt.want {
				t.Errorf("ShortSellingPosition.HasDiscretionaryInvestment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortSellingPosition_IsIndividual(t *testing.T) {
	tests := []struct {
		name     string
		position ShortSellingPosition
		want     bool
	}{
		{
			name: "individual",
			position: ShortSellingPosition{
				SSName: "個人",
			},
			want: true,
		},
		{
			name: "institution",
			position: ShortSellingPosition{
				SSName: "ABC Investment",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.position.IsIndividual(); got != tt.want {
				t.Errorf("ShortSellingPosition.IsIndividual() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortSellingPosition_GetPositionPercentage(t *testing.T) {
	position := ShortSellingPosition{
		ShrtPosToSO: 0.0053,
	}

	expected := 0.53
	got := position.GetPositionPercentage()

	if got != expected {
		t.Errorf("ShortSellingPosition.GetPositionPercentage() = %v, want %v", got, expected)
	}
}

func TestShortSellingPosition_GetPreviousPositionPercentage(t *testing.T) {
	position := ShortSellingPosition{
		PrevRptRatio: 0.0043,
	}

	expected := 0.43
	got := position.GetPreviousPositionPercentage()

	if got != expected {
		t.Errorf("ShortSellingPosition.GetPreviousPositionPercentage() = %v, want %v", got, expected)
	}
}

func TestShortSellingPositionsService_GetShortSellingPositions_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewShortSellingPositionsService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/markets/short-sale-report?code=86970", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetShortSellingPositions(ShortSellingPositionsParams{Code: "86970"})

	// Verify
	if err == nil {
		t.Error("GetShortSellingPositions() expected error but got nil")
	}
}
