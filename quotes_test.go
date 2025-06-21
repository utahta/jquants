package jquants

import (
	"fmt"
	"testing"
	"time"

	"github.com/utahta/jquants/client"
)

func TestQuotesService_GetDailyQuotes(t *testing.T) {
	tests := []struct {
		name     string
		params   DailyQuotesParams
		wantPath string
	}{
		{
			name:     "with all parameters",
			params:   DailyQuotesParams{Code: "7203", From: "20240101", To: "20240131"},
			wantPath: "/prices/daily_quotes?code=7203&from=20240101&to=20240131",
		},
		{
			name:     "with code only",
			params:   DailyQuotesParams{Code: "7203"},
			wantPath: "/prices/daily_quotes?code=7203",
		},
		{
			name:     "with date only",
			params:   DailyQuotesParams{Date: "20240101"},
			wantPath: "/prices/daily_quotes?date=20240101",
		},
		{
			name:     "with pagination key",
			params:   DailyQuotesParams{Date: "20240101", PaginationKey: "key123"},
			wantPath: "/prices/daily_quotes?date=20240101&pagination_key=key123",
		},
		{
			name:     "with no parameters",
			params:   DailyQuotesParams{},
			wantPath: "/prices/daily_quotes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewQuotesService(mockClient)

			// Mock response
			mockResponse := DailyQuotesResponse{
				DailyQuotes: []DailyQuote{
					{
						Date:             "20240101",
						Code:             "7203",
						Open:             floatPtr(2490.0),
						High:             floatPtr(2510.0),
						Low:              floatPtr(2480.0),
						Close:            floatPtr(2500.0),
						UpperLimit:       "0",
						LowerLimit:       "0",
						Volume:           floatPtr(1000000),
						AdjustmentFactor: 1.0,
						AdjustmentClose:  floatPtr(2500.0),
					},
					{
						Date:       "20240102",
						Code:       "7203",
						Close:      floatPtr(2520.0),
						UpperLimit: "1", // ストップ高
						LowerLimit: "0",
					},
				},
				PaginationKey: "",
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Test
			resp, err := service.GetDailyQuotes(tt.params)
			if err != nil {
				t.Errorf("GetDailyQuotes failed: %v", err)
			}

			// Verify
			if len(resp.DailyQuotes) != 2 {
				t.Errorf("Expected 2 quotes, got %d", len(resp.DailyQuotes))
			}

			if mockClient.LastMethod != "GET" {
				t.Errorf("Expected GET method, got %s", mockClient.LastMethod)
			}

			if mockClient.LastPath != tt.wantPath {
				t.Errorf("Expected path %s, got %s", tt.wantPath, mockClient.LastPath)
			}
		})
	}
}

func TestQuotesService_GetDailyQuotesByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewQuotesService(mockClient)

	// Calculate expected dates
	to := time.Now()
	from := to.AddDate(0, 0, -30)
	expectedPath := fmt.Sprintf("/prices/daily_quotes?code=7203&from=%s&to=%s",
		from.Format("20060102"), to.Format("20060102"))

	// Mock response
	mockResponse := DailyQuotesResponse{
		DailyQuotes: []DailyQuote{
			{
				Date:            "20240101",
				Code:            "7203",
				Open:            floatPtr(2480.0),
				High:            floatPtr(2510.0),
				Low:             floatPtr(2470.0),
				Close:           floatPtr(2500.0),
				Volume:          floatPtr(1000000),
				AdjustmentClose: floatPtr(2500.0),
			},
		},
		PaginationKey: "", // ページネーションなし
	}

	mockClient.SetResponse("GET", expectedPath, mockResponse)

	// Test
	quotes, err := service.GetDailyQuotesByCode("7203", 30)
	if err != nil {
		t.Errorf("GetDailyQuotesByCode failed: %v", err)
	}

	// Verify
	if len(quotes) != 1 {
		t.Errorf("Expected 1 quote, got %d", len(quotes))
	}

	if quotes[0].Code != "7203" {
		t.Errorf("Expected code 7203, got %s", quotes[0].Code)
	}

	// Check that request was made with correct parameters
	if mockClient.LastMethod != "GET" {
		t.Errorf("Expected GET method, got %s", mockClient.LastMethod)
	}

	if mockClient.LastPath != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, mockClient.LastPath)
	}
}

func TestQuotesService_GetDailyQuotes_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewQuotesService(mockClient)

	// Mock error
	mockClient.SetError("GET", "/prices/daily_quotes?code=7203", fmt.Errorf("API error"))

	// Test
	_, err := service.GetDailyQuotes(DailyQuotesParams{Code: "7203"})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestQuotesService_GetDailyQuotesByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewQuotesService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := DailyQuotesResponse{
		DailyQuotes: []DailyQuote{
			{
				Date:  "20240101",
				Code:  "1301",
				Close: floatPtr(1000.0),
			},
			{
				Date:  "20240101",
				Code:  "1332",
				Close: floatPtr(2000.0),
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := DailyQuotesResponse{
		DailyQuotes: []DailyQuote{
			{
				Date:  "20240101",
				Code:  "7203",
				Close: floatPtr(2500.0),
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/prices/daily_quotes?date=20240101", mockResponse1)
	mockClient.SetResponse("GET", "/prices/daily_quotes?date=20240101&pagination_key=next_page_key", mockResponse2)

	// Test
	quotes, err := service.GetDailyQuotesByDate("20240101")
	if err != nil {
		t.Errorf("GetDailyQuotesByDate failed: %v", err)
	}

	// Verify
	if len(quotes) != 3 {
		t.Errorf("Expected 3 quotes total, got %d", len(quotes))
	}
}

func TestDailyQuote_StopLimitMethods(t *testing.T) {
	tests := []struct {
		name   string
		quote  DailyQuote
		method string
		want   bool
	}{
		{
			name:   "ストップ高",
			quote:  DailyQuote{UpperLimit: "1", LowerLimit: "0"},
			method: "IsStopHigh",
			want:   true,
		},
		{
			name:   "ストップ安",
			quote:  DailyQuote{UpperLimit: "0", LowerLimit: "1"},
			method: "IsStopLow",
			want:   true,
		},
		{
			name:   "通常",
			quote:  DailyQuote{UpperLimit: "0", LowerLimit: "0"},
			method: "IsStopHigh",
			want:   false,
		},
		{
			name:   "前場ストップ高",
			quote:  DailyQuote{MorningUpperLimit: "1", MorningLowerLimit: "0"},
			method: "IsMorningStopHigh",
			want:   true,
		},
		{
			name:   "後場ストップ安",
			quote:  DailyQuote{AfternoonUpperLimit: "0", AfternoonLowerLimit: "1"},
			method: "IsAfternoonStopLow",
			want:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bool
			switch tt.method {
			case "IsStopHigh":
				got = tt.quote.IsStopHigh()
			case "IsStopLow":
				got = tt.quote.IsStopLow()
			case "IsMorningStopHigh":
				got = tt.quote.IsMorningStopHigh()
			case "IsAfternoonStopLow":
				got = tt.quote.IsAfternoonStopLow()
			}

			if got != tt.want {
				t.Errorf("%s() = %v, want %v", tt.method, got, tt.want)
			}
		})
	}
}

func TestDailyQuote_HasData(t *testing.T) {
	tests := []struct {
		name  string
		quote DailyQuote
		want  bool
	}{
		{
			name: "前場データあり",
			quote: DailyQuote{
				MorningOpen:  floatPtr(1000),
				MorningHigh:  floatPtr(1010),
				MorningLow:   floatPtr(990),
				MorningClose: floatPtr(1005),
			},
			want: true,
		},
		{
			name: "前場データなし",
			quote: DailyQuote{
				Open:  floatPtr(1000),
				Close: floatPtr(1005),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.quote.HasMorningData(); got != tt.want {
				t.Errorf("HasMorningData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuotesService_DateFormatting(t *testing.T) {
	// This test verifies date formatting works correctly
	to := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	from := to.AddDate(0, 0, -30)

	expectedFrom := "20231216"
	expectedTo := "20240115"

	if from.Format("20060102") != expectedFrom {
		t.Errorf("Expected from date %s, got %s", expectedFrom, from.Format("20060102"))
	}

	if to.Format("20060102") != expectedTo {
		t.Errorf("Expected to date %s, got %s", expectedTo, to.Format("20060102"))
	}
}

func TestQuotesService_GetDailyQuotesByCode_WithPagination(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewQuotesService(mockClient)

	// Calculate expected dates
	to := time.Now()
	from := to.AddDate(0, 0, -30)
	basePath := fmt.Sprintf("/prices/daily_quotes?code=7203&from=%s&to=%s",
		from.Format("20060102"), to.Format("20060102"))

	// Mock response - 最初のページ
	mockResponse1 := DailyQuotesResponse{
		DailyQuotes: []DailyQuote{
			{
				Date:  "20240101",
				Code:  "7203",
				Close: floatPtr(2500.0),
			},
			{
				Date:  "20240102",
				Code:  "7203",
				Close: floatPtr(2510.0),
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := DailyQuotesResponse{
		DailyQuotes: []DailyQuote{
			{
				Date:  "20240103",
				Code:  "7203",
				Close: floatPtr(2520.0),
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", basePath, mockResponse1)
	mockClient.SetResponse("GET", basePath+"&pagination_key=next_page_key", mockResponse2)

	// Test
	quotes, err := service.GetDailyQuotesByCode("7203", 30)
	if err != nil {
		t.Errorf("GetDailyQuotesByCode failed: %v", err)
	}

	// Verify
	if len(quotes) != 3 {
		t.Errorf("Expected 3 quotes total, got %d", len(quotes))
	}

	// 各データの確認
	if quotes[0].Date != "20240101" || *quotes[0].Close != 2500.0 {
		t.Errorf("First quote data mismatch")
	}
	if quotes[2].Date != "20240103" || *quotes[2].Close != 2520.0 {
		t.Errorf("Last quote data mismatch")
	}
}

// Note: Helper functions are now in test_helpers.go
