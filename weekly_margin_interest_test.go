package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestWeeklyMarginInterestService_GetWeeklyMarginInterest(t *testing.T) {
	tests := []struct {
		name     string
		params   WeeklyMarginInterestParams
		wantPath string
	}{
		{
			name: "with code and date range",
			params: WeeklyMarginInterestParams{
				Code: "13010",
				From: "20230101",
				To:   "20230331",
			},
			wantPath: "/markets/weekly_margin_interest?code=13010&from=20230101&to=20230331",
		},
		{
			name: "with code only",
			params: WeeklyMarginInterestParams{
				Code: "13010",
			},
			wantPath: "/markets/weekly_margin_interest?code=13010",
		},
		{
			name: "with date only",
			params: WeeklyMarginInterestParams{
				Date: "20230217",
			},
			wantPath: "/markets/weekly_margin_interest?date=20230217",
		},
		{
			name: "with date and pagination key",
			params: WeeklyMarginInterestParams{
				Date:          "20230217",
				PaginationKey: "key123",
			},
			wantPath: "/markets/weekly_margin_interest?date=20230217&pagination_key=key123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewWeeklyMarginInterestService(mockClient)

			// Mock response based on documentation sample
			mockResponse := WeeklyMarginInterestResponse{
				WeeklyMarginInterest: []WeeklyMarginInterest{
					{
						Date:                               "2023-02-17",
						Code:                               "13010",
						IssueType:                          "2",
						ShortMarginTradeVolume:             4100.0,
						LongMarginTradeVolume:              27600.0,
						ShortNegotiableMarginTradeVolume:   1300.0,
						LongNegotiableMarginTradeVolume:    7600.0,
						ShortStandardizedMarginTradeVolume: 2800.0,
						LongStandardizedMarginTradeVolume:  20000.0,
					},
				},
				PaginationKey: "",
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Execute
			resp, err := service.GetWeeklyMarginInterest(tt.params)

			// Verify
			if err != nil {
				t.Fatalf("GetWeeklyMarginInterest() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetWeeklyMarginInterest() returned nil response")
			}
			if len(resp.WeeklyMarginInterest) == 0 {
				t.Error("GetWeeklyMarginInterest() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetWeeklyMarginInterest() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestWeeklyMarginInterestService_GetWeeklyMarginInterest_RequiresCodeOrDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewWeeklyMarginInterestService(mockClient)

	// Execute with empty code and date
	_, err := service.GetWeeklyMarginInterest(WeeklyMarginInterestParams{})

	// Verify
	if err == nil {
		t.Error("GetWeeklyMarginInterest() expected error for missing code and date but got nil")
	}
	if err.Error() != "either code or date parameter is required" {
		t.Errorf("GetWeeklyMarginInterest() error = %v, want 'either code or date parameter is required'", err)
	}
}

func TestWeeklyMarginInterestService_GetWeeklyMarginInterestByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewWeeklyMarginInterestService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := WeeklyMarginInterestResponse{
		WeeklyMarginInterest: []WeeklyMarginInterest{
			{
				Date:                   "2023-02-17",
				Code:                   "13010",
				IssueType:              "2",
				ShortMarginTradeVolume: 4100.0,
				LongMarginTradeVolume:  27600.0,
			},
			{
				Date:                   "2023-02-10",
				Code:                   "13010",
				IssueType:              "2",
				ShortMarginTradeVolume: 3900.0,
				LongMarginTradeVolume:  26800.0,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := WeeklyMarginInterestResponse{
		WeeklyMarginInterest: []WeeklyMarginInterest{
			{
				Date:                   "2023-02-03",
				Code:                   "13010",
				IssueType:              "2",
				ShortMarginTradeVolume: 3800.0,
				LongMarginTradeVolume:  25900.0,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/markets/weekly_margin_interest?code=13010", mockResponse1)
	mockClient.SetResponse("GET", "/markets/weekly_margin_interest?code=13010&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetWeeklyMarginInterestByCode("13010")

	// Verify
	if err != nil {
		t.Fatalf("GetWeeklyMarginInterestByCode() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetWeeklyMarginInterestByCode() returned %d items, want 3", len(data))
	}
	for _, item := range data {
		if item.Code != "13010" {
			t.Errorf("GetWeeklyMarginInterestByCode() returned code %v, want 13010", item.Code)
		}
	}
}

func TestWeeklyMarginInterestService_GetWeeklyMarginInterestByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewWeeklyMarginInterestService(mockClient)

	// Mock response
	mockResponse := WeeklyMarginInterestResponse{
		WeeklyMarginInterest: []WeeklyMarginInterest{
			{
				Date:                   "2023-02-17",
				Code:                   "13010",
				IssueType:              "2",
				ShortMarginTradeVolume: 4100.0,
				LongMarginTradeVolume:  27600.0,
			},
			{
				Date:                   "2023-02-17",
				Code:                   "86970",
				IssueType:              "1",
				ShortMarginTradeVolume: 2300.0,
				LongMarginTradeVolume:  15400.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/markets/weekly_margin_interest?date=20230217", mockResponse)

	// Execute
	data, err := service.GetWeeklyMarginInterestByDate("20230217")

	// Verify
	if err != nil {
		t.Fatalf("GetWeeklyMarginInterestByDate() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetWeeklyMarginInterestByDate() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.Date != "2023-02-17" {
			t.Errorf("GetWeeklyMarginInterestByDate() returned date %v, want 2023-02-17", item.Date)
		}
	}
}

func TestWeeklyMarginInterestService_GetWeeklyMarginInterestByCodeAndDateRange(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewWeeklyMarginInterestService(mockClient)

	// Mock response
	mockResponse := WeeklyMarginInterestResponse{
		WeeklyMarginInterest: []WeeklyMarginInterest{
			{
				Date:                   "2023-02-17",
				Code:                   "86970",
				IssueType:              "1",
				ShortMarginTradeVolume: 2300.0,
				LongMarginTradeVolume:  15400.0,
			},
			{
				Date:                   "2023-02-10",
				Code:                   "86970",
				IssueType:              "1",
				ShortMarginTradeVolume: 2200.0,
				LongMarginTradeVolume:  15100.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/markets/weekly_margin_interest?code=86970&from=20230101&to=20230331", mockResponse)

	// Execute
	data, err := service.GetWeeklyMarginInterestByCodeAndDateRange("86970", "20230101", "20230331")

	// Verify
	if err != nil {
		t.Fatalf("GetWeeklyMarginInterestByCodeAndDateRange() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetWeeklyMarginInterestByCodeAndDateRange() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.Code != "86970" {
			t.Errorf("GetWeeklyMarginInterestByCodeAndDateRange() returned code %v, want 86970", item.Code)
		}
	}
}

func TestWeeklyMarginInterest_IsCredit(t *testing.T) {
	tests := []struct {
		name string
		data WeeklyMarginInterest
		want bool
	}{
		{
			name: "credit issue",
			data: WeeklyMarginInterest{
				IssueType: IssueTypeCredit,
			},
			want: true,
		},
		{
			name: "lendable issue",
			data: WeeklyMarginInterest{
				IssueType: IssueTypeLendable,
			},
			want: false,
		},
		{
			name: "other issue",
			data: WeeklyMarginInterest{
				IssueType: IssueTypeOther,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.data.IsCredit()
			if got != tt.want {
				t.Errorf("WeeklyMarginInterest.IsCredit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeeklyMarginInterest_IsLendable(t *testing.T) {
	tests := []struct {
		name string
		data WeeklyMarginInterest
		want bool
	}{
		{
			name: "lendable issue",
			data: WeeklyMarginInterest{
				IssueType: IssueTypeLendable,
			},
			want: true,
		},
		{
			name: "credit issue",
			data: WeeklyMarginInterest{
				IssueType: IssueTypeCredit,
			},
			want: false,
		},
		{
			name: "other issue",
			data: WeeklyMarginInterest{
				IssueType: IssueTypeOther,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.data.IsLendable()
			if got != tt.want {
				t.Errorf("WeeklyMarginInterest.IsLendable() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeeklyMarginInterest_GetShortLongRatio(t *testing.T) {
	tests := []struct {
		name string
		data WeeklyMarginInterest
		want float64
	}{
		{
			name: "normal ratio",
			data: WeeklyMarginInterest{
				ShortMarginTradeVolume: 4100.0,
				LongMarginTradeVolume:  27600.0,
			},
			want: 4100.0 / 27600.0,
		},
		{
			name: "zero long margin",
			data: WeeklyMarginInterest{
				ShortMarginTradeVolume: 4100.0,
				LongMarginTradeVolume:  0.0,
			},
			want: 0.0,
		},
		{
			name: "zero short margin",
			data: WeeklyMarginInterest{
				ShortMarginTradeVolume: 0.0,
				LongMarginTradeVolume:  27600.0,
			},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.data.GetShortLongRatio()
			if got != tt.want {
				t.Errorf("WeeklyMarginInterest.GetShortLongRatio() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeeklyMarginInterest_GetStandardizedRatio(t *testing.T) {
	data := WeeklyMarginInterest{
		ShortMarginTradeVolume:             4100.0,
		LongMarginTradeVolume:              27600.0,
		ShortStandardizedMarginTradeVolume: 2800.0,
		LongStandardizedMarginTradeVolume:  20000.0,
	}

	shortRatio, longRatio := data.GetStandardizedRatio()

	expectedShortRatio := 2800.0 / 4100.0
	expectedLongRatio := 20000.0 / 27600.0

	if shortRatio != expectedShortRatio {
		t.Errorf("WeeklyMarginInterest.GetStandardizedRatio() short ratio = %v, want %v", shortRatio, expectedShortRatio)
	}
	if longRatio != expectedLongRatio {
		t.Errorf("WeeklyMarginInterest.GetStandardizedRatio() long ratio = %v, want %v", longRatio, expectedLongRatio)
	}
}

func TestWeeklyMarginInterest_GetStandardizedRatio_ZeroTotal(t *testing.T) {
	data := WeeklyMarginInterest{
		ShortMarginTradeVolume:             0.0,
		LongMarginTradeVolume:              0.0,
		ShortStandardizedMarginTradeVolume: 1000.0,
		LongStandardizedMarginTradeVolume:  2000.0,
	}

	shortRatio, longRatio := data.GetStandardizedRatio()

	if shortRatio != 0.0 {
		t.Errorf("WeeklyMarginInterest.GetStandardizedRatio() short ratio = %v, want 0.0", shortRatio)
	}
	if longRatio != 0.0 {
		t.Errorf("WeeklyMarginInterest.GetStandardizedRatio() long ratio = %v, want 0.0", longRatio)
	}
}

func TestWeeklyMarginInterestService_GetWeeklyMarginInterest_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewWeeklyMarginInterestService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/markets/weekly_margin_interest?code=13010", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetWeeklyMarginInterest(WeeklyMarginInterestParams{Code: "13010"})

	// Verify
	if err == nil {
		t.Error("GetWeeklyMarginInterest() expected error but got nil")
	}
}

func TestIssueTypeConstants(t *testing.T) {
	// 銘柄区分定数の値を確認
	tests := []struct {
		constant string
		expected string
	}{
		{IssueTypeCredit, "1"},
		{IssueTypeLendable, "2"},
		{IssueTypeOther, "3"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Issue type constant = %v, want %v", tt.constant, tt.expected)
			}
		})
	}
}
