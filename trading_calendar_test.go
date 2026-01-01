package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestTradingCalendarService_GetTradingCalendar(t *testing.T) {
	tests := []struct {
		name     string
		params   TradingCalendarParams
		wantPath string
	}{
		{
			name: "with all parameters",
			params: TradingCalendarParams{
				HolidayDivision: "1",
				From:            "20240101",
				To:              "20240131",
			},
			wantPath: "/markets/trading_calendar?holidaydivision=1&from=20240101&to=20240131",
		},
		{
			name: "with date range only",
			params: TradingCalendarParams{
				From: "20240101",
				To:   "20240107",
			},
			wantPath: "/markets/trading_calendar?from=20240101&to=20240107",
		},
		{
			name: "with holiday division only",
			params: TradingCalendarParams{
				HolidayDivision: "1",
			},
			wantPath: "/markets/trading_calendar?holidaydivision=1",
		},
		{
			name:     "with no parameters",
			params:   TradingCalendarParams{},
			wantPath: "/markets/trading_calendar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewTradingCalendarService(mockClient)

			// Mock response
			mockResponse := TradingCalendarResponse{
				TradingCalendar: []TradingCalendar{
					{
						Date:            "2015-04-01",
						HolidayDivision: "1",
					},
					{
						Date:            "2015-04-02",
						HolidayDivision: "1",
					},
				},
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Execute
			resp, err := service.GetTradingCalendar(tt.params)

			// Verify
			if err != nil {
				t.Fatalf("GetTradingCalendar() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetTradingCalendar() returned nil response")
				return
			}
			if len(resp.TradingCalendar) != 2 {
				t.Errorf("GetTradingCalendar() returned %d items, want 2", len(resp.TradingCalendar))
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetTradingCalendar() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestTradingCalendarService_GetTradingCalendarByDateRange(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTradingCalendarService(mockClient)

	// Mock response
	mockResponse := TradingCalendarResponse{
		TradingCalendar: []TradingCalendar{
			{
				Date:            "2024-01-01",
				HolidayDivision: "0",
			},
			{
				Date:            "2024-01-02",
				HolidayDivision: "1",
			},
		},
	}
	mockClient.SetResponse("GET", "/markets/trading_calendar?from=20240101&to=20240107", mockResponse)

	// Execute
	resp, err := service.GetTradingCalendarByDateRange("20240101", "20240107")

	// Verify
	if err != nil {
		t.Fatalf("GetTradingCalendarByDateRange() error = %v", err)
	}
	if len(resp) != 2 {
		t.Errorf("GetTradingCalendarByDateRange() returned %d items, want 2", len(resp))
	}
	if resp[0].Date != "2024-01-01" {
		t.Errorf("GetTradingCalendarByDateRange() returned date %v, want 2024-01-01", resp[0].Date)
	}
	if resp[0].HolidayDivision != "0" {
		t.Errorf("GetTradingCalendarByDateRange() returned holiday division %v, want 0", resp[0].HolidayDivision)
	}
}

func TestTradingCalendarService_GetTradingDays(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTradingCalendarService(mockClient)

	// Mock response
	mockResponse := TradingCalendarResponse{
		TradingCalendar: []TradingCalendar{
			{
				Date:            "2024-01-02",
				HolidayDivision: "1",
			},
			{
				Date:            "2024-01-03",
				HolidayDivision: "1",
			},
		},
	}
	mockClient.SetResponse("GET", "/markets/trading_calendar?holidaydivision=1&from=20240101&to=20240107", mockResponse)

	// Execute
	tradingDays, err := service.GetTradingDays("20240101", "20240107")

	// Verify
	if err != nil {
		t.Fatalf("GetTradingDays() error = %v", err)
	}
	if len(tradingDays) != 2 {
		t.Errorf("GetTradingDays() returned %d items, want 2", len(tradingDays))
	}
	for _, day := range tradingDays {
		if !day.IsTradingDay() {
			t.Errorf("GetTradingDays() returned non-trading day %v", day.Date)
		}
	}
}

func TestTradingCalendar_HelperMethods(t *testing.T) {
	tests := []struct {
		name            string
		holidayDivision string
		isTradingDay    bool
		isNonTradingDay bool
		isHalfDay       bool
		hasOSETrading   bool
	}{
		{
			name:            "trading day",
			holidayDivision: "1",
			isTradingDay:    true,
			isNonTradingDay: false,
			isHalfDay:       false,
			hasOSETrading:   false,
		},
		{
			name:            "non-trading day",
			holidayDivision: "0",
			isTradingDay:    false,
			isNonTradingDay: true,
			isHalfDay:       false,
			hasOSETrading:   false,
		},
		{
			name:            "half trading day",
			holidayDivision: "2",
			isTradingDay:    false,
			isNonTradingDay: false,
			isHalfDay:       true,
			hasOSETrading:   false,
		},
		{
			name:            "OSE holiday trading",
			holidayDivision: "3",
			isTradingDay:    false,
			isNonTradingDay: false,
			isHalfDay:       false,
			hasOSETrading:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := TradingCalendar{HolidayDivision: tt.holidayDivision}

			if tc.IsTradingDay() != tt.isTradingDay {
				t.Errorf("IsTradingDay() = %v, want %v", tc.IsTradingDay(), tt.isTradingDay)
			}
			if tc.IsNonTradingDay() != tt.isNonTradingDay {
				t.Errorf("IsNonTradingDay() = %v, want %v", tc.IsNonTradingDay(), tt.isNonTradingDay)
			}
			if tc.IsHalfTradingDay() != tt.isHalfDay {
				t.Errorf("IsHalfTradingDay() = %v, want %v", tc.IsHalfTradingDay(), tt.isHalfDay)
			}
			if tc.HasOSEHolidayTrading() != tt.hasOSETrading {
				t.Errorf("HasOSEHolidayTrading() = %v, want %v", tc.HasOSEHolidayTrading(), tt.hasOSETrading)
			}
		})
	}
}

func TestTradingCalendarService_GetTradingCalendar_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTradingCalendarService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/markets/trading_calendar", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetTradingCalendar(TradingCalendarParams{})

	// Verify
	if err == nil {
		t.Error("GetTradingCalendar() expected error but got nil")
	}
}
