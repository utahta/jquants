package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestTOPIXService_GetTOPIXData(t *testing.T) {
	tests := []struct {
		name     string
		params   TOPIXParams
		wantPath string
	}{
		{
			name: "with all parameters",
			params: TOPIXParams{
				From:          "20240101",
				To:            "20240131",
				PaginationKey: "key123",
			},
			wantPath: "/indices/topix?from=20240101&to=20240131&pagination_key=key123",
		},
		{
			name: "with date range only",
			params: TOPIXParams{
				From: "20240101",
				To:   "20240131",
			},
			wantPath: "/indices/topix?from=20240101&to=20240131",
		},
		{
			name: "with from date only",
			params: TOPIXParams{
				From: "20240101",
			},
			wantPath: "/indices/topix?from=20240101",
		},
		{
			name:     "with no parameters",
			params:   TOPIXParams{},
			wantPath: "/indices/topix",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewTOPIXService(mockClient)

			// Mock response
			mockResponse := TOPIXResponse{
				TOPIX: []TOPIXData{
					{
						Date:  "2022-06-28",
						Open:  1885.52,
						High:  1907.38,
						Low:   1885.32,
						Close: 1907.38,
					},
				},
				PaginationKey: "",
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Execute
			resp, err := service.GetTOPIXData(tt.params)

			// Verify
			if err != nil {
				t.Fatalf("GetTOPIXData() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetTOPIXData() returned nil response")
				return
			}
			if len(resp.TOPIX) == 0 {
				t.Error("GetTOPIXData() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetTOPIXData() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestTOPIXService_GetTOPIXByDateRange(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTOPIXService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := TOPIXResponse{
		TOPIX: []TOPIXData{
			{
				Date:  "2024-01-01",
				Open:  2400.0,
				Close: 2420.0,
			},
			{
				Date:  "2024-01-02",
				Open:  2420.0,
				Close: 2435.0,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := TOPIXResponse{
		TOPIX: []TOPIXData{
			{
				Date:  "2024-01-03",
				Open:  2435.0,
				Close: 2450.0,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/indices/topix?from=20240101&to=20240103", mockResponse1)
	mockClient.SetResponse("GET", "/indices/topix?from=20240101&to=20240103&pagination_key=next_page_key", mockResponse2)

	// Execute
	topixData, err := service.GetTOPIXByDateRange("20240101", "20240103")

	// Verify
	if err != nil {
		t.Fatalf("GetTOPIXByDateRange() error = %v", err)
	}
	if len(topixData) != 3 {
		t.Errorf("GetTOPIXByDateRange() returned %d items, want 3", len(topixData))
	}
}

func TestTOPIXService_GetLatestTOPIX(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTOPIXService(mockClient)

	// Mock response with multiple data points (最新が最初にくる想定)
	mockResponse := TOPIXResponse{
		TOPIX: []TOPIXData{
			{
				Date:  "2024-02-01",
				Open:  2480.0,
				High:  2490.0,
				Low:   2475.0,
				Close: 2485.0,
			},
			{
				Date:  "2024-01-31",
				Close: 2470.0,
			},
			{
				Date:  "2024-01-30",
				Close: 2450.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/indices/topix", mockResponse)

	// Execute
	latest, err := service.GetLatestTOPIX()

	// Verify
	if err != nil {
		t.Fatalf("GetLatestTOPIX() error = %v", err)
	}
	if latest == nil {
		t.Fatal("GetLatestTOPIX() returned nil")
		return
	}
	if latest.Date != "2024-02-01" {
		t.Errorf("GetLatestTOPIX() returned date %v, want 2024-02-01", latest.Date)
	}
	if latest.Close != 2485.0 {
		t.Errorf("GetLatestTOPIX() returned close %v, want 2485.0", latest.Close)
	}
}

func TestTOPIXService_GetLatestTOPIX_NoData(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTOPIXService(mockClient)

	// Mock empty response
	mockResponse := TOPIXResponse{
		TOPIX:         []TOPIXData{},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/indices/topix", mockResponse)

	// Execute
	_, err := service.GetLatestTOPIX()

	// Verify
	if err == nil {
		t.Error("GetLatestTOPIX() expected error for empty data but got nil")
	}
}

func TestTOPIXService_GetTOPIXData_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTOPIXService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/indices/topix", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetTOPIXData(TOPIXParams{})

	// Verify
	if err == nil {
		t.Error("GetTOPIXData() expected error but got nil")
	}
}
