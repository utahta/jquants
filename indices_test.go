package jquants

import (
	"fmt"
	"testing"

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
			wantPath: "/indices/bars/daily?code=0000&date=20240101&from=20240101&to=20240131&pagination_key=key123",
		},
		{
			name: "with code and date range",
			params: IndicesParams{
				Code: "0000",
				From: "20240101",
				To:   "20240131",
			},
			wantPath: "/indices/bars/daily?code=0000&from=20240101&to=20240131",
		},
		{
			name: "with code only",
			params: IndicesParams{
				Code: "0028",
			},
			wantPath: "/indices/bars/daily?code=0028",
		},
		{
			name: "with date only",
			params: IndicesParams{
				Date: "20240101",
			},
			wantPath: "/indices/bars/daily?date=20240101",
		},
		{
			name:     "with no parameters",
			params:   IndicesParams{},
			wantPath: "/indices/bars/daily",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewIndicesService(mockClient)

			// Mock response
			mockResponse := IndicesResponse{
				Data: []Index{
					{
						Date: "2024-01-31",
						Code: "0000",
						O:    2400.0,
						H:    2420.0,
						L:    2390.0,
						C:    2410.0,
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
				return
			}
			if len(resp.Data) == 0 {
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

	// Mock response - 最初のページ
	mockResponse1 := IndicesResponse{
		Data: []Index{
			{Date: "2024-01-01", Code: "0000", O: 2400.0, H: 2420.0, L: 2390.0, C: 2410.0},
			{Date: "2024-01-02", Code: "0000", C: 2420.0},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := IndicesResponse{
		Data: []Index{
			{Date: "2024-01-03", Code: "0000", C: 2430.0},
		},
		PaginationKey: "",
	}

	mockClient.SetResponse("GET", "/indices/bars/daily?code=0000", mockResponse1)
	mockClient.SetResponse("GET", "/indices/bars/daily?code=0000&pagination_key=next_page_key", mockResponse2)

	// Execute
	indices, err := service.GetIndicesByCode("0000")

	// Verify
	if err != nil {
		t.Fatalf("GetIndicesByCode() error = %v", err)
	}
	if len(indices) != 3 {
		t.Errorf("GetIndicesByCode() returned %d items, want 3", len(indices))
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
		Data: []Index{
			{
				Date: "2024-01-01",
				Code: "0000",
				C:    2400.0,
			},
			{
				Date: "2024-01-01",
				Code: "0028",
				C:    1200.0,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := IndicesResponse{
		Data: []Index{
			{
				Date: "2024-01-01",
				Code: "0050",
				C:    1600.0,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/indices/bars/daily?date=20240101", mockResponse1)
	mockClient.SetResponse("GET", "/indices/bars/daily?date=20240101&pagination_key=next_page_key", mockResponse2)

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

	// Mock response
	mockResponse := IndicesResponse{
		Data: []Index{
			{Date: "2024-02-01", Code: "0000", C: 2400.0},
		},
	}
	mockClient.SetResponse("GET", "/indices/bars/daily?code=0000", mockResponse)

	// Execute
	indices, err := service.GetTOPIX()

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

	// Mock response - 情報・通信業
	mockResponse := IndicesResponse{
		Data: []Index{
			{Date: "2024-02-01", Code: "0058", O: 3000.0, H: 3050.0, L: 2980.0, C: 3020.0},
		},
	}
	mockClient.SetResponse("GET", "/indices/bars/daily?code=0058", mockResponse)

	// Execute
	indices, err := service.GetSectorIndex(IndexSectorInfoComm)

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

	// Mock response
	mockResponse := IndicesResponse{
		Data: []Index{
			{Date: "2024-02-01", Code: "0500", C: 1800.0},
		},
	}
	mockClient.SetResponse("GET", "/indices/bars/daily?code=0500", mockResponse)

	// Execute
	indices, err := service.GetPrimeMarketIndex()

	// Verify
	if err != nil {
		t.Fatalf("GetPrimeMarketIndex() error = %v", err)
	}
	if indices[0].Code != "0500" {
		t.Errorf("GetPrimeMarketIndex() returned code %v, want 0500", indices[0].Code)
	}
}

func TestIndicesService_GetIndicesByCodeAndDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndicesService(mockClient)

	// Mock response
	mockResponse := IndicesResponse{
		Data: []Index{
			{
				Date: "2024-01-01",
				Code: "0000",
				O:    2400.0,
				H:    2420.0,
				L:    2390.0,
				C:    2410.0,
			},
		},
	}
	mockClient.SetResponse("GET", "/indices/bars/daily?code=0000&date=20240101", mockResponse)

	// Execute
	indices, err := service.GetIndicesByCodeAndDate("0000", "20240101")

	// Verify
	if err != nil {
		t.Fatalf("GetIndicesByCodeAndDate() error = %v", err)
	}
	if len(indices) != 1 {
		t.Errorf("GetIndicesByCodeAndDate() returned %d items, want 1", len(indices))
	}
	if indices[0].Code != "0000" || indices[0].Date != "2024-01-01" {
		t.Errorf("Index data mismatch: code=%s, date=%s", indices[0].Code, indices[0].Date)
	}
	if mockClient.LastPath != "/indices/bars/daily?code=0000&date=20240101" {
		t.Errorf("Expected path /indices/bars/daily?code=0000&date=20240101, got %s", mockClient.LastPath)
	}
}

func TestIndicesService_GetIndicesByCodeAndDateRange(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndicesService(mockClient)

	basePath := "/indices/bars/daily?code=0000&from=20240101&to=20240131"

	// Mock response - 最初のページ
	mockResponse1 := IndicesResponse{
		Data: []Index{
			{Date: "2024-01-01", Code: "0000", C: 2400.0},
			{Date: "2024-01-02", Code: "0000", C: 2410.0},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := IndicesResponse{
		Data: []Index{
			{Date: "2024-01-03", Code: "0000", C: 2420.0},
		},
		PaginationKey: "",
	}

	mockClient.SetResponse("GET", basePath, mockResponse1)
	mockClient.SetResponse("GET", basePath+"&pagination_key=next_page_key", mockResponse2)

	// Execute
	indices, err := service.GetIndicesByCodeAndDateRange("0000", "20240101", "20240131")

	// Verify
	if err != nil {
		t.Fatalf("GetIndicesByCodeAndDateRange() error = %v", err)
	}
	if len(indices) != 3 {
		t.Errorf("GetIndicesByCodeAndDateRange() returned %d items, want 3", len(indices))
	}
	if indices[0].Date != "2024-01-01" || indices[0].C != 2400.0 {
		t.Errorf("First index data mismatch")
	}
	if indices[2].Date != "2024-01-03" || indices[2].C != 2420.0 {
		t.Errorf("Last index data mismatch")
	}
}

func TestIndicesService_GetIndices_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewIndicesService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/indices/bars/daily?code=0000", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetIndices(IndicesParams{Code: "0000"})

	// Verify
	if err == nil {
		t.Error("GetIndices() expected error but got nil")
	}
}

