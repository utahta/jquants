package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestMinuteQuotesService_GetMinuteQuotes(t *testing.T) {
	tests := []struct {
		name     string
		params   MinuteQuotesParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with code only",
			params: MinuteQuotesParams{
				Code: "86970",
			},
			wantPath: "/equities/bars/minute?code=86970",
		},
		{
			name: "with code and date",
			params: MinuteQuotesParams{
				Code: "86970",
				Date: "20230324",
			},
			wantPath: "/equities/bars/minute?code=86970&date=20230324",
		},
		{
			name: "with date only",
			params: MinuteQuotesParams{
				Date: "2023-03-24",
			},
			wantPath: "/equities/bars/minute?date=2023-03-24",
		},
		{
			name: "with code and date range",
			params: MinuteQuotesParams{
				Code: "86970",
				From: "20230301",
				To:   "20230331",
			},
			wantPath: "/equities/bars/minute?code=86970&from=20230301&to=20230331",
		},
		{
			name: "with pagination key",
			params: MinuteQuotesParams{
				Code:          "86970",
				PaginationKey: "key123",
			},
			wantPath: "/equities/bars/minute?code=86970&pagination_key=key123",
		},
		{
			name:     "with no required parameters",
			params:   MinuteQuotesParams{},
			wantPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewMinuteQuotesService(mockClient)

			// Mock response based on documentation sample
			mockResponse := MinuteQuotesResponse{
				Data: []MinuteQuote{
					{
						Date: "2023-03-24",
						Time: "09:00",
						Code: "86970",
						O:    2047.0,
						H:    2055.0,
						L:    2045.0,
						C:    2050.0,
						Vo:   12500.0,
						Va:   25625000.0,
					},
				},
				PaginationKey: "",
			}
			if !tt.wantErr {
				mockClient.SetResponse("GET", tt.wantPath, mockResponse)
			}

			// Execute
			resp, err := service.GetMinuteQuotes(context.Background(), tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetMinuteQuotes() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetMinuteQuotes() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetMinuteQuotes() returned nil response")
			}
			if len(resp.Data) == 0 {
				t.Error("GetMinuteQuotes() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetMinuteQuotes() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestMinuteQuotesService_GetMinuteQuotes_RequiresParameter(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewMinuteQuotesService(mockClient)

	// Execute with empty parameters
	_, err := service.GetMinuteQuotes(context.Background(), MinuteQuotesParams{})

	// Verify
	if err == nil {
		t.Error("GetMinuteQuotes() expected error for missing parameters but got nil")
	}
	if err.Error() != "either code or date parameter is required" {
		t.Errorf("GetMinuteQuotes() error = %v, want 'either code or date parameter is required'", err)
	}
}

func TestMinuteQuotesResponse_UnmarshalJSON(t *testing.T) {
	// ドキュメントのレスポンスサンプルに基づくJSON
	// 2件目は数値の文字列・nullが混在するケース（防御的なデコードの確認）
	raw := `{
		"data": [
			{
				"Date": "2023-03-24",
				"Time": "09:00",
				"Code": "86970",
				"O": 2047.0,
				"H": 2055.0,
				"L": 2045.0,
				"C": 2050.0,
				"Vo": 12500.0,
				"Va": 25625000.0
			},
			{
				"Date": "2023-03-24",
				"Time": "09:01",
				"Code": "86970",
				"O": "2050.5",
				"H": 2052.0,
				"L": 2049.0,
				"C": 2051.0,
				"Vo": 3200.0,
				"Va": null
			}
		],
		"pagination_key": "value1.value2."
	}`

	var resp MinuteQuotesResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if len(resp.Data) != 2 {
		t.Fatalf("UnmarshalJSON() returned %d items, want 2", len(resp.Data))
	}
	if resp.PaginationKey != "value1.value2." {
		t.Errorf("UnmarshalJSON() PaginationKey = %v, want value1.value2.", resp.PaginationKey)
	}

	q := resp.Data[0]
	if q.Date != "2023-03-24" {
		t.Errorf("UnmarshalJSON() Date = %v, want 2023-03-24", q.Date)
	}
	if q.Time != "09:00" {
		t.Errorf("UnmarshalJSON() Time = %v, want 09:00", q.Time)
	}
	if q.Code != "86970" {
		t.Errorf("UnmarshalJSON() Code = %v, want 86970", q.Code)
	}
	if q.O != 2047.0 {
		t.Errorf("UnmarshalJSON() O = %v, want 2047.0", q.O)
	}
	if q.H != 2055.0 {
		t.Errorf("UnmarshalJSON() H = %v, want 2055.0", q.H)
	}
	if q.L != 2045.0 {
		t.Errorf("UnmarshalJSON() L = %v, want 2045.0", q.L)
	}
	if q.C != 2050.0 {
		t.Errorf("UnmarshalJSON() C = %v, want 2050.0", q.C)
	}
	if q.Vo != 12500.0 {
		t.Errorf("UnmarshalJSON() Vo = %v, want 12500.0", q.Vo)
	}
	if q.Va != 25625000.0 {
		t.Errorf("UnmarshalJSON() Va = %v, want 25625000.0", q.Va)
	}

	// 数値文字列は数値としてパースされ、nullはゼロ値になる
	if resp.Data[1].O != 2050.5 {
		t.Errorf("UnmarshalJSON() O = %v, want 2050.5", resp.Data[1].O)
	}
	if resp.Data[1].Va != 0 {
		t.Errorf("UnmarshalJSON() Va = %v, want 0", resp.Data[1].Va)
	}
}

func TestMinuteQuotesService_GetMinuteQuotesByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewMinuteQuotesService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := MinuteQuotesResponse{
		Data: []MinuteQuote{
			{
				Date: "2023-03-24",
				Time: "09:00",
				Code: "86970",
				O:    2047.0,
				H:    2055.0,
				L:    2045.0,
				C:    2050.0,
				Vo:   12500.0,
				Va:   25625000.0,
			},
			{
				Date: "2023-03-24",
				Time: "09:01",
				Code: "86970",
				O:    2050.0,
				H:    2052.0,
				L:    2049.0,
				C:    2051.0,
				Vo:   3200.0,
				Va:   6563200.0,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := MinuteQuotesResponse{
		Data: []MinuteQuote{
			{
				Date: "2023-03-24",
				Time: "09:02",
				Code: "86970",
				O:    2051.0,
				H:    2053.0,
				L:    2050.0,
				C:    2052.0,
				Vo:   1800.0,
				Va:   3693600.0,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/equities/bars/minute?code=86970", mockResponse1)
	mockClient.SetResponse("GET", "/equities/bars/minute?code=86970&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetMinuteQuotesByCode(context.Background(), "86970")

	// Verify
	if err != nil {
		t.Fatalf("GetMinuteQuotesByCode() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetMinuteQuotesByCode() returned %d items, want 3", len(data))
	}
	for _, item := range data {
		if item.Code != "86970" {
			t.Errorf("GetMinuteQuotesByCode() returned code %v, want 86970", item.Code)
		}
	}
}

func TestMinuteQuotesService_GetMinuteQuotesByCodeAndDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewMinuteQuotesService(mockClient)

	// Mock response
	mockResponse := MinuteQuotesResponse{
		Data: []MinuteQuote{
			{
				Date: "2023-03-24",
				Time: "09:00",
				Code: "86970",
				O:    2047.0,
				C:    2050.0,
				Vo:   12500.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/equities/bars/minute?code=86970&date=20230324", mockResponse)

	// Execute
	data, err := service.GetMinuteQuotesByCodeAndDate(context.Background(), "86970", "20230324")

	// Verify
	if err != nil {
		t.Fatalf("GetMinuteQuotesByCodeAndDate() error = %v", err)
	}
	if len(data) != 1 {
		t.Errorf("GetMinuteQuotesByCodeAndDate() returned %d items, want 1", len(data))
	}
	if data[0].Code != "86970" {
		t.Errorf("GetMinuteQuotesByCodeAndDate() returned code %v, want 86970", data[0].Code)
	}
	if data[0].Date != "2023-03-24" {
		t.Errorf("GetMinuteQuotesByCodeAndDate() returned date %v, want 2023-03-24", data[0].Date)
	}
}

func TestMinuteQuotesService_GetMinuteQuotesByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewMinuteQuotesService(mockClient)

	// Mock response
	mockResponse := MinuteQuotesResponse{
		Data: []MinuteQuote{
			{
				Date: "2023-03-24",
				Time: "09:00",
				Code: "86970",
				O:    2047.0,
			},
			{
				Date: "2023-03-24",
				Time: "09:00",
				Code: "13660",
				O:    1500.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/equities/bars/minute?date=20230324", mockResponse)

	// Execute
	data, err := service.GetMinuteQuotesByDate(context.Background(), "20230324")

	// Verify
	if err != nil {
		t.Fatalf("GetMinuteQuotesByDate() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetMinuteQuotesByDate() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.Date != "2023-03-24" {
			t.Errorf("GetMinuteQuotesByDate() returned date %v, want 2023-03-24", item.Date)
		}
	}
}

func TestMinuteQuotesService_GetMinuteQuotes_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewMinuteQuotesService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/equities/bars/minute?code=86970", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetMinuteQuotes(context.Background(), MinuteQuotesParams{Code: "86970"})

	// Verify
	if err == nil {
		t.Error("GetMinuteQuotes() expected error but got nil")
	}
}
