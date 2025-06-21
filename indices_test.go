package jquants

import (
	"fmt"
	"testing"
	"time"

	"github.com/utahta/jquants/client"
)

func TestIndicesService_GetIndices(t *testing.T) {
	tests := []struct {
		name     string
		params   IndicesParams
		wantPath string
	}{
		{
			name: "with all parameters",
			params: IndicesParams{
				Code:          "0000",
				Date:          "20240101",
				From:          "20240101",
				To:            "20240131",
				PaginationKey: "key123",
			},
			wantPath: "/indices?code=0000&date=20240101&from=20240101&to=20240131&pagination_key=key123",
		},
		{
			name: "with code and date range",
			params: IndicesParams{
				Code: "0000",
				From: "20240101",
				To:   "20240131",
			},
			wantPath: "/indices?code=0000&from=20240101&to=20240131",
		},
		{
			name: "with code only",
			params: IndicesParams{
				Code: "0028",
			},
			wantPath: "/indices?code=0028",
		},
		{
			name: "with date only",
			params: IndicesParams{
				Date: "20240101",
			},
			wantPath: "/indices?date=20240101",
		},
		{
			name:     "with no parameters",
			params:   IndicesParams{},
			wantPath: "/indices",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewIndicesService(mockClient)

			// Mock response
			mockResponse := IndicesResponse{
				Indices: []Index{
					{
						Date:  "2024-01-31",
						Code:  "0000",
						Open:  2400.0,
						High:  2420.0,
						Low:   2390.0,
						Close: 2410.0,
					},
				},
				PaginationKey: "",
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Execute
			resp, err := service.GetIndices(tt.params)

			// Verify
			if err != nil {
				t.Fatalf("GetIndices() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetIndices() returned nil response")
			}
			if len(resp.Indices) == 0 {
				t.Error("GetIndices() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetIndices() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestIndicesService_GetIndicesByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndicesService(mockClient)

	// Calculate expected dates
	to := time.Now()
	from := to.AddDate(0, 0, -30)
	expectedPath := fmt.Sprintf("/indices?code=0000&from=%s&to=%s",
		from.Format("20060102"), to.Format("20060102"))

	// Mock response
	mockResponse := IndicesResponse{
		Indices: []Index{
			{
				Date:  "2024-02-01",
				Code:  "0000",
				Open:  2400.0,
				High:  2420.0,
				Low:   2390.0,
				Close: 2410.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", expectedPath, mockResponse)

	// Execute
	indices, err := service.GetIndicesByCode("0000", 30)

	// Verify
	if err != nil {
		t.Fatalf("GetIndicesByCode() error = %v", err)
	}
	if len(indices) != 1 {
		t.Errorf("GetIndicesByCode() returned %d items, want 1", len(indices))
	}
	if indices[0].Code != "0000" {
		t.Errorf("GetIndicesByCode() returned code %v, want 0000", indices[0].Code)
	}
}

func TestIndicesService_GetIndicesByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndicesService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := IndicesResponse{
		Indices: []Index{
			{
				Date:  "2024-01-01",
				Code:  "0000",
				Close: 2400.0,
			},
			{
				Date:  "2024-01-01",
				Code:  "0028",
				Close: 1200.0,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := IndicesResponse{
		Indices: []Index{
			{
				Date:  "2024-01-01",
				Code:  "0050",
				Close: 1600.0,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/indices?date=20240101", mockResponse1)
	mockClient.SetResponse("GET", "/indices?date=20240101&pagination_key=next_page_key", mockResponse2)

	// Execute
	indices, err := service.GetIndicesByDate("20240101")

	// Verify
	if err != nil {
		t.Fatalf("GetIndicesByDate() error = %v", err)
	}
	if len(indices) != 3 {
		t.Errorf("GetIndicesByDate() returned %d items, want 3", len(indices))
	}
}

func TestIndicesService_GetTOPIX(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndicesService(mockClient)

	// Calculate expected dates
	to := time.Now()
	from := to.AddDate(0, 0, -30)
	expectedPath := fmt.Sprintf("/indices?code=0000&from=%s&to=%s",
		from.Format("20060102"), to.Format("20060102"))

	// Mock response
	mockResponse := IndicesResponse{
		Indices: []Index{
			{
				Date:  "2024-02-01",
				Code:  "0000",
				Close: 2400.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", expectedPath, mockResponse)

	// Execute
	indices, err := service.GetTOPIX(30)

	// Verify
	if err != nil {
		t.Fatalf("GetTOPIX() error = %v", err)
	}
	if len(indices) != 1 {
		t.Errorf("GetTOPIX() returned %d items, want 1", len(indices))
	}
	if indices[0].Code != "0000" {
		t.Errorf("GetTOPIX() returned code %v, want 0000", indices[0].Code)
	}
}

func TestIndicesService_GetSectorIndex(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndicesService(mockClient)

	// Calculate expected dates
	to := time.Now()
	from := to.AddDate(0, 0, -30)
	expectedPath := fmt.Sprintf("/indices?code=0058&from=%s&to=%s",
		from.Format("20060102"), to.Format("20060102"))

	// Mock response - 情報・通信業
	mockResponse := IndicesResponse{
		Indices: []Index{
			{
				Date:  "2024-02-01",
				Code:  "0058",
				Open:  3000.0,
				High:  3050.0,
				Low:   2980.0,
				Close: 3020.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", expectedPath, mockResponse)

	// Execute
	indices, err := service.GetSectorIndex(IndexSectorInfoComm, 30)

	// Verify
	if err != nil {
		t.Fatalf("GetSectorIndex() error = %v", err)
	}
	if len(indices) != 1 {
		t.Errorf("GetSectorIndex() returned %d items, want 1", len(indices))
	}
	if indices[0].Code != "0058" {
		t.Errorf("GetSectorIndex() returned code %v, want 0058", indices[0].Code)
	}
}

func TestIndicesService_GetPrimeMarketIndex(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndicesService(mockClient)

	// Calculate expected dates
	to := time.Now()
	from := to.AddDate(0, 0, -30)
	expectedPath := fmt.Sprintf("/indices?code=0500&from=%s&to=%s",
		from.Format("20060102"), to.Format("20060102"))

	// Mock response
	mockResponse := IndicesResponse{
		Indices: []Index{
			{
				Date:  "2024-02-01",
				Code:  "0500",
				Close: 1800.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", expectedPath, mockResponse)

	// Execute
	indices, err := service.GetPrimeMarketIndex(30)

	// Verify
	if err != nil {
		t.Fatalf("GetPrimeMarketIndex() error = %v", err)
	}
	if indices[0].Code != "0500" {
		t.Errorf("GetPrimeMarketIndex() returned code %v, want 0500", indices[0].Code)
	}
}

func TestIndicesService_GetIndices_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndicesService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/indices?code=0000", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetIndices(IndicesParams{Code: "0000"})

	// Verify
	if err == nil {
		t.Error("GetIndices() expected error but got nil")
	}
}

func TestIndicesService_GetIndicesByCode_WithPagination(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndicesService(mockClient)

	// Calculate expected dates
	to := time.Now()
	from := to.AddDate(0, 0, -100)
	basePath := fmt.Sprintf("/indices?code=0000&from=%s&to=%s",
		from.Format("20060102"), to.Format("20060102"))

	// Mock response - 最初のページ
	mockResponse1 := IndicesResponse{
		Indices: []Index{
			{Date: "2024-01-01", Code: "0000", Close: 2400.0},
			{Date: "2024-01-02", Code: "0000", Close: 2410.0},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := IndicesResponse{
		Indices: []Index{
			{Date: "2024-01-03", Code: "0000", Close: 2420.0},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", basePath, mockResponse1)
	mockClient.SetResponse("GET", basePath+"&pagination_key=next_page_key", mockResponse2)

	// Execute
	indices, err := service.GetIndicesByCode("0000", 100)

	// Verify
	if err != nil {
		t.Fatalf("GetIndicesByCode() error = %v", err)
	}
	if len(indices) != 3 {
		t.Errorf("GetIndicesByCode() returned %d items, want 3", len(indices))
	}
}
