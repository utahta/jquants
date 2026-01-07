package jquants

import (
	"fmt"
	"testing"
	"time"

	"github.com/utahta/jquants/client"
)

func TestBreakdownService_GetBreakdown(t *testing.T) {
	tests := []struct {
		name     string
		params   BreakdownParams
		wantPath string
	}{
		{
			name: "with all parameters",
			params: BreakdownParams{
				Code:          "7203",
				Date:          "20240101",
				From:          "20240101",
				To:            "20240131",
				PaginationKey: "key123",
			},
			wantPath: "/markets/breakdown?code=7203&date=20240101&from=20240101&to=20240131&pagination_key=key123",
		},
		{
			name: "with code and date range",
			params: BreakdownParams{
				Code: "7203",
				From: "20240101",
				To:   "20240131",
			},
			wantPath: "/markets/breakdown?code=7203&from=20240101&to=20240131",
		},
		{
			name: "with code only",
			params: BreakdownParams{
				Code: "7203",
			},
			wantPath: "/markets/breakdown?code=7203",
		},
		{
			name: "with date only",
			params: BreakdownParams{
				Date: "20240101",
			},
			wantPath: "/markets/breakdown?date=20240101",
		},
		{
			name:     "with no parameters",
			params:   BreakdownParams{},
			wantPath: "/markets/breakdown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewBreakdownService(mockClient)

			// Mock response
			mockResponse := BreakdownResponse{
				Data: []Breakdown{
					{
						Date:            "2024-01-31",
						Code:            "72030",
						LongSellVa:      115164000.0,
						ShrtNoMrgnVa:    93561000.0,
						MrgnSellNewVa:   6412000.0,
						MrgnSellCloseVa: 23009000.0,
						LongBuyVa:       185114000.0,
						MrgnBuyNewVa:    35568000.0,
						MrgnBuyCloseVa:  17464000.0,
						LongSellVo:      415000.0,
						ShrtNoMrgnVo:    337000.0,
						MrgnSellNewVo:   23000.0,
						MrgnSellCloseVo: 83000.0,
						LongBuyVo:       667000.0,
						MrgnBuyNewVo:    128000.0,
						MrgnBuyCloseVo:  63000.0,
					},
				},
				PaginationKey: "",
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Execute
			resp, err := service.GetBreakdown(tt.params)

			// Verify
			if err != nil {
				t.Fatalf("GetBreakdown() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetBreakdown() returned nil response")
				return
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetBreakdown() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
			if len(resp.Data) != 1 {
				t.Errorf("GetBreakdown() returned %d items, want 1", len(resp.Data))
			}
		})
	}
}

func TestBreakdownService_GetBreakdownByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewBreakdownService(mockClient)

	// Calculate expected dates
	to := time.Now()
	from := to.AddDate(0, 0, -30)
	expectedPath := fmt.Sprintf("/markets/breakdown?code=7203&from=%s&to=%s",
		from.Format("20060102"), to.Format("20060102"))

	// Mock response
	mockResponse := BreakdownResponse{
		Data: []Breakdown{
			{
				Date:       "2024-02-01",
				Code:       "72030",
				LongSellVa: 100000000.0,
				LongBuyVa:  120000000.0,
				LongSellVo: 400000.0,
				LongBuyVo:  500000.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", expectedPath, mockResponse)

	// Execute
	breakdown, err := service.GetBreakdownByCode("7203", 30)

	// Verify
	if err != nil {
		t.Fatalf("GetBreakdownByCode() error = %v", err)
	}
	if len(breakdown) != 1 {
		t.Errorf("GetBreakdownByCode() returned %d items, want 1", len(breakdown))
	}
	if breakdown[0].Code != "72030" {
		t.Errorf("GetBreakdownByCode() returned code %v, want 72030", breakdown[0].Code)
	}
}

func TestBreakdownService_GetBreakdownByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewBreakdownService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := BreakdownResponse{
		Data: []Breakdown{
			{
				Date:       "2024-01-01",
				Code:       "13010",
				LongSellVa: 100000000.0,
			},
			{
				Date:       "2024-01-01",
				Code:       "13020",
				LongSellVa: 200000000.0,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := BreakdownResponse{
		Data: []Breakdown{
			{
				Date:       "2024-01-01",
				Code:       "72030",
				LongSellVa: 300000000.0,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/markets/breakdown?date=20240101", mockResponse1)
	mockClient.SetResponse("GET", "/markets/breakdown?date=20240101&pagination_key=next_page_key", mockResponse2)

	// Execute
	breakdown, err := service.GetBreakdownByDate("20240101")

	// Verify
	if err != nil {
		t.Fatalf("GetBreakdownByDate() error = %v", err)
	}
	if len(breakdown) != 3 {
		t.Errorf("GetBreakdownByDate() returned %d items, want 3", len(breakdown))
	}
}

func TestBreakdownService_GetBreakdown_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewBreakdownService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/markets/breakdown?code=7203", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetBreakdown(BreakdownParams{Code: "7203"})

	// Verify
	if err == nil {
		t.Error("GetBreakdown() expected error but got nil")
	}
}

func TestBreakdown_HelperMethods(t *testing.T) {
	b := &Breakdown{
		// 売りの約定代金
		LongSellVa:      100.0,
		ShrtNoMrgnVa:    50.0,
		MrgnSellNewVa:   30.0,
		MrgnSellCloseVa: 20.0,
		// 買いの約定代金
		LongBuyVa:      150.0,
		MrgnBuyNewVa:   40.0,
		MrgnBuyCloseVa: 10.0,
		// 売りの約定株数
		LongSellVo:      1000.0,
		ShrtNoMrgnVo:    500.0,
		MrgnSellNewVo:   300.0,
		MrgnSellCloseVo: 200.0,
		// 買いの約定株数
		LongBuyVo:      1500.0,
		MrgnBuyNewVo:   400.0,
		MrgnBuyCloseVo: 100.0,
	}

	// 売り合計テスト
	if totalSellValue := b.GetTotalSellValue(); totalSellValue != 200.0 {
		t.Errorf("GetTotalSellValue() = %v, want 200.0", totalSellValue)
	}

	if totalSellVolume := b.GetTotalSellVolume(); totalSellVolume != 2000.0 {
		t.Errorf("GetTotalSellVolume() = %v, want 2000.0", totalSellVolume)
	}

	// 買い合計テスト
	if totalBuyValue := b.GetTotalBuyValue(); totalBuyValue != 200.0 {
		t.Errorf("GetTotalBuyValue() = %v, want 200.0", totalBuyValue)
	}

	if totalBuyVolume := b.GetTotalBuyVolume(); totalBuyVolume != 2000.0 {
		t.Errorf("GetTotalBuyVolume() = %v, want 2000.0", totalBuyVolume)
	}

	// 信用取引合計テスト
	if marginNewValue := b.GetMarginNewValue(); marginNewValue != 70.0 {
		t.Errorf("GetMarginNewValue() = %v, want 70.0", marginNewValue)
	}

	if marginCloseValue := b.GetMarginCloseValue(); marginCloseValue != 30.0 {
		t.Errorf("GetMarginCloseValue() = %v, want 30.0", marginCloseValue)
	}

	// 空売り比率テスト
	// 空売り合計 = 50 + 30 = 80, 売り合計 = 200, 比率 = 80/200 = 0.4
	if shortSellRatio := b.GetShortSellRatio(); shortSellRatio != 0.4 {
		t.Errorf("GetShortSellRatio() = %v, want 0.4", shortSellRatio)
	}
}

func TestBreakdown_GetShortSellRatio_ZeroDivision(t *testing.T) {
	b := &Breakdown{
		// 全て0
		LongSellVa:      0,
		ShrtNoMrgnVa:    0,
		MrgnSellNewVa:   0,
		MrgnSellCloseVa: 0,
	}

	// ゼロ除算のテスト
	if ratio := b.GetShortSellRatio(); ratio != 0 {
		t.Errorf("GetShortSellRatio() with zero total sell = %v, want 0", ratio)
	}
}
