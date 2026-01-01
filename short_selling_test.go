package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestShortSellingService_GetShortSelling(t *testing.T) {
	tests := []struct {
		name     string
		params   ShortSellingParams
		wantPath string
	}{
		{
			name: "with sector and date range",
			params: ShortSellingParams{
				Sector33Code: "0050",
				From:         "20220101",
				To:           "20221231",
			},
			wantPath: "/markets/short_selling?sector33code=0050&from=20220101&to=20221231",
		},
		{
			name: "with sector only",
			params: ShortSellingParams{
				Sector33Code: "0050",
			},
			wantPath: "/markets/short_selling?sector33code=0050",
		},
		{
			name: "with date only",
			params: ShortSellingParams{
				Date: "20221025",
			},
			wantPath: "/markets/short_selling?date=20221025",
		},
		{
			name: "with date and pagination key",
			params: ShortSellingParams{
				Date:          "20221025",
				PaginationKey: "key123",
			},
			wantPath: "/markets/short_selling?date=20221025&pagination_key=key123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewShortSellingService(mockClient)

			// Mock response based on documentation sample
			mockResponse := ShortSellingResponse{
				ShortSelling: []ShortSelling{
					{
						Date:         "2022-10-25",
						Sector33Code: "0050",
						SellingExcludingShortSellingTurnoverValue:    1333126400.0,
						ShortSellingWithRestrictionsTurnoverValue:    787355200.0,
						ShortSellingWithoutRestrictionsTurnoverValue: 149084300.0,
					},
				},
				PaginationKey: "",
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Execute
			resp, err := service.GetShortSelling(tt.params)

			// Verify
			if err != nil {
				t.Fatalf("GetShortSelling() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetShortSelling() returned nil response")
				return
			}
			if len(resp.ShortSelling) == 0 {
				t.Error("GetShortSelling() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetShortSelling() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestShortSellingService_GetShortSelling_RequiresSectorOrDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewShortSellingService(mockClient)

	// Execute with empty sector and date
	_, err := service.GetShortSelling(ShortSellingParams{})

	// Verify
	if err == nil {
		t.Error("GetShortSelling() expected error for missing sector and date but got nil")
	}
	if err.Error() != "either sector33code or date parameter is required" {
		t.Errorf("GetShortSelling() error = %v, want 'either sector33code or date parameter is required'", err)
	}
}

func TestShortSellingService_GetShortSellingBySector(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewShortSellingService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := ShortSellingResponse{
		ShortSelling: []ShortSelling{
			{
				Date:         "2022-10-25",
				Sector33Code: "0050",
				SellingExcludingShortSellingTurnoverValue:    1333126400.0,
				ShortSellingWithRestrictionsTurnoverValue:    787355200.0,
				ShortSellingWithoutRestrictionsTurnoverValue: 149084300.0,
			},
			{
				Date:         "2022-10-24",
				Sector33Code: "0050",
				SellingExcludingShortSellingTurnoverValue:    1200000000.0,
				ShortSellingWithRestrictionsTurnoverValue:    750000000.0,
				ShortSellingWithoutRestrictionsTurnoverValue: 140000000.0,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := ShortSellingResponse{
		ShortSelling: []ShortSelling{
			{
				Date:         "2022-10-21",
				Sector33Code: "0050",
				SellingExcludingShortSellingTurnoverValue:    1100000000.0,
				ShortSellingWithRestrictionsTurnoverValue:    700000000.0,
				ShortSellingWithoutRestrictionsTurnoverValue: 130000000.0,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/markets/short_selling?sector33code=0050", mockResponse1)
	mockClient.SetResponse("GET", "/markets/short_selling?sector33code=0050&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetShortSellingBySector("0050")

	// Verify
	if err != nil {
		t.Fatalf("GetShortSellingBySector() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetShortSellingBySector() returned %d items, want 3", len(data))
	}
	for _, item := range data {
		if item.Sector33Code != "0050" {
			t.Errorf("GetShortSellingBySector() returned sector33code %v, want 0050", item.Sector33Code)
		}
	}
}

func TestShortSellingService_GetShortSellingByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewShortSellingService(mockClient)

	// Mock response
	mockResponse := ShortSellingResponse{
		ShortSelling: []ShortSelling{
			{
				Date:         "2022-10-25",
				Sector33Code: "0050",
				SellingExcludingShortSellingTurnoverValue:    1333126400.0,
				ShortSellingWithRestrictionsTurnoverValue:    787355200.0,
				ShortSellingWithoutRestrictionsTurnoverValue: 149084300.0,
			},
			{
				Date:         "2022-10-25",
				Sector33Code: "1050",
				SellingExcludingShortSellingTurnoverValue:    500000000.0,
				ShortSellingWithRestrictionsTurnoverValue:    300000000.0,
				ShortSellingWithoutRestrictionsTurnoverValue: 50000000.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/markets/short_selling?date=20221025", mockResponse)

	// Execute
	data, err := service.GetShortSellingByDate("20221025")

	// Verify
	if err != nil {
		t.Fatalf("GetShortSellingByDate() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetShortSellingByDate() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.Date != "2022-10-25" {
			t.Errorf("GetShortSellingByDate() returned date %v, want 2022-10-25", item.Date)
		}
	}
}

func TestShortSellingService_GetShortSellingBySectorAndDateRange(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewShortSellingService(mockClient)

	// Mock response
	mockResponse := ShortSellingResponse{
		ShortSelling: []ShortSelling{
			{
				Date:         "2022-10-25",
				Sector33Code: "0050",
				SellingExcludingShortSellingTurnoverValue:    1333126400.0,
				ShortSellingWithRestrictionsTurnoverValue:    787355200.0,
				ShortSellingWithoutRestrictionsTurnoverValue: 149084300.0,
			},
			{
				Date:         "2022-10-24",
				Sector33Code: "0050",
				SellingExcludingShortSellingTurnoverValue:    1200000000.0,
				ShortSellingWithRestrictionsTurnoverValue:    750000000.0,
				ShortSellingWithoutRestrictionsTurnoverValue: 140000000.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/markets/short_selling?sector33code=0050&from=20220101&to=20221231", mockResponse)

	// Execute
	data, err := service.GetShortSellingBySectorAndDateRange("0050", "20220101", "20221231")

	// Verify
	if err != nil {
		t.Fatalf("GetShortSellingBySectorAndDateRange() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetShortSellingBySectorAndDateRange() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.Sector33Code != "0050" {
			t.Errorf("GetShortSellingBySectorAndDateRange() returned sector33code %v, want 0050", item.Sector33Code)
		}
	}
}

func TestShortSelling_GetTotalShortSellingValue(t *testing.T) {
	data := ShortSelling{
		ShortSellingWithRestrictionsTurnoverValue:    787355200.0,
		ShortSellingWithoutRestrictionsTurnoverValue: 149084300.0,
	}

	expected := 787355200.0 + 149084300.0
	got := data.GetTotalShortSellingValue()

	if got != expected {
		t.Errorf("ShortSelling.GetTotalShortSellingValue() = %v, want %v", got, expected)
	}
}

func TestShortSelling_GetTotalTurnoverValue(t *testing.T) {
	data := ShortSelling{
		SellingExcludingShortSellingTurnoverValue:    1333126400.0,
		ShortSellingWithRestrictionsTurnoverValue:    787355200.0,
		ShortSellingWithoutRestrictionsTurnoverValue: 149084300.0,
	}

	expected := 1333126400.0 + 787355200.0 + 149084300.0
	got := data.GetTotalTurnoverValue()

	if got != expected {
		t.Errorf("ShortSelling.GetTotalTurnoverValue() = %v, want %v", got, expected)
	}
}

func TestShortSelling_GetShortSellingRatio(t *testing.T) {
	data := ShortSelling{
		SellingExcludingShortSellingTurnoverValue:    1333126400.0,
		ShortSellingWithRestrictionsTurnoverValue:    787355200.0,
		ShortSellingWithoutRestrictionsTurnoverValue: 149084300.0,
	}

	totalShortSelling := 787355200.0 + 149084300.0
	totalTurnover := 1333126400.0 + totalShortSelling
	expected := (totalShortSelling / totalTurnover) * 100

	got := data.GetShortSellingRatio()

	if got != expected {
		t.Errorf("ShortSelling.GetShortSellingRatio() = %v, want %v", got, expected)
	}
}

func TestShortSelling_GetShortSellingRatio_ZeroTurnover(t *testing.T) {
	data := ShortSelling{
		SellingExcludingShortSellingTurnoverValue:    0.0,
		ShortSellingWithRestrictionsTurnoverValue:    0.0,
		ShortSellingWithoutRestrictionsTurnoverValue: 0.0,
	}

	got := data.GetShortSellingRatio()

	if got != 0.0 {
		t.Errorf("ShortSelling.GetShortSellingRatio() = %v, want 0.0", got)
	}
}

func TestShortSelling_GetRestrictedShortSellingRatio(t *testing.T) {
	data := ShortSelling{
		ShortSellingWithRestrictionsTurnoverValue:    787355200.0,
		ShortSellingWithoutRestrictionsTurnoverValue: 149084300.0,
	}

	totalShortSelling := 787355200.0 + 149084300.0
	expected := (787355200.0 / totalShortSelling) * 100

	got := data.GetRestrictedShortSellingRatio()

	if got != expected {
		t.Errorf("ShortSelling.GetRestrictedShortSellingRatio() = %v, want %v", got, expected)
	}
}

func TestShortSelling_GetUnrestrictedShortSellingRatio(t *testing.T) {
	data := ShortSelling{
		ShortSellingWithRestrictionsTurnoverValue:    787355200.0,
		ShortSellingWithoutRestrictionsTurnoverValue: 149084300.0,
	}

	totalShortSelling := 787355200.0 + 149084300.0
	expected := (149084300.0 / totalShortSelling) * 100

	got := data.GetUnrestrictedShortSellingRatio()

	if got != expected {
		t.Errorf("ShortSelling.GetUnrestrictedShortSellingRatio() = %v, want %v", got, expected)
	}
}

func TestShortSelling_GetRestrictedShortSellingRatio_ZeroShortSelling(t *testing.T) {
	data := ShortSelling{
		ShortSellingWithRestrictionsTurnoverValue:    0.0,
		ShortSellingWithoutRestrictionsTurnoverValue: 0.0,
	}

	got := data.GetRestrictedShortSellingRatio()

	if got != 0.0 {
		t.Errorf("ShortSelling.GetRestrictedShortSellingRatio() = %v, want 0.0", got)
	}
}

func TestShortSelling_GetUnrestrictedShortSellingRatio_ZeroShortSelling(t *testing.T) {
	data := ShortSelling{
		ShortSellingWithRestrictionsTurnoverValue:    0.0,
		ShortSellingWithoutRestrictionsTurnoverValue: 0.0,
	}

	got := data.GetUnrestrictedShortSellingRatio()

	if got != 0.0 {
		t.Errorf("ShortSelling.GetUnrestrictedShortSellingRatio() = %v, want 0.0", got)
	}
}

func TestShortSellingService_GetShortSelling_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewShortSellingService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/markets/short_selling?sector33code=0050", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetShortSelling(ShortSellingParams{Sector33Code: "0050"})

	// Verify
	if err == nil {
		t.Error("GetShortSelling() expected error but got nil")
	}
}
