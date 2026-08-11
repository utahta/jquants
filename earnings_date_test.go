package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestEarningsDateService_GetEarningsDates(t *testing.T) {
	tests := []struct {
		name     string
		params   EarningsDateParams
		wantPath string
	}{
		{
			name:     "with code",
			params:   EarningsDateParams{Code: "86970"},
			wantPath: "/fins/earnings-date?code=86970",
		},
		{
			name:     "with date",
			params:   EarningsDateParams{Date: "20250620"},
			wantPath: "/fins/earnings-date?date=20250620",
		},
		{
			name:     "with date in hyphen format",
			params:   EarningsDateParams{Date: "2025-06-20"},
			wantPath: "/fins/earnings-date?date=2025-06-20",
		},
		{
			name:     "with scheduled date",
			params:   EarningsDateParams{ScheduledDate: "20250805"},
			wantPath: "/fins/earnings-date?scheduled_date=20250805",
		},
		{
			name:     "with pagination key",
			params:   EarningsDateParams{Code: "86970", PaginationKey: "key123"},
			wantPath: "/fins/earnings-date?code=86970&pagination_key=key123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := client.NewMockClient()
			service := NewEarningsDateService(mockClient)

			mockResponse := EarningsDateResponse{
				Data: []EarningsDate{
					{
						PubDate:  "2025-06-03",
						SchDate:  "2025-07-30",
						FQName:   "1Q",
						FYE:      "0331",
						Code:     "86970",
						CoName:   "日本取引所グループ",
						CoNameEn: "Japan Exchange Group,Inc.",
					},
				},
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			resp, err := service.GetEarningsDates(context.Background(), tt.params)
			if err != nil {
				t.Fatalf("GetEarningsDates() error = %v", err)
			}

			if len(resp.Data) != 1 {
				t.Errorf("GetEarningsDates() returned %d items, want 1", len(resp.Data))
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetEarningsDates() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestEarningsDateService_GetEarningsDates_ParamConflict(t *testing.T) {
	tests := []struct {
		name   string
		params EarningsDateParams
	}{
		{
			name:   "no parameter",
			params: EarningsDateParams{},
		},
		{
			name:   "pagination key only",
			params: EarningsDateParams{PaginationKey: "key123"},
		},
		{
			name:   "code and date",
			params: EarningsDateParams{Code: "86970", Date: "20250620"},
		},
		{
			name:   "date and scheduled date",
			params: EarningsDateParams{Date: "20250620", ScheduledDate: "20250805"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := client.NewMockClient()
			service := NewEarningsDateService(mockClient)

			if _, err := service.GetEarningsDates(context.Background(), tt.params); err == nil {
				t.Error("GetEarningsDates() expected error, got nil")
			}
			if mockClient.RequestCount != 0 {
				t.Errorf("GetEarningsDates() sent %d requests, want 0", mockClient.RequestCount)
			}
		})
	}
}

func TestEarningsDateService_GetEarningsDatesByCode(t *testing.T) {
	mockClient := client.NewMockClient()
	service := NewEarningsDateService(mockClient)

	// 1ページ目: pagination_keyあり
	mockClient.SetResponse("GET", "/fins/earnings-date?code=86970", EarningsDateResponse{
		Data: []EarningsDate{
			{PubDate: "2025-06-03", SchDate: "2025-07-30", FQName: "1Q", Code: "86970"},
		},
		PaginationKey: "key123",
	})
	// 2ページ目: pagination_keyなし
	mockClient.SetResponse("GET", "/fins/earnings-date?code=86970&pagination_key=key123", EarningsDateResponse{
		Data: []EarningsDate{
			{PubDate: "2025-07-01", SchDate: "2025-08-05", FQName: "1Q", Code: "86970"},
		},
	})

	dates, err := service.GetEarningsDatesByCode(context.Background(), "86970")
	if err != nil {
		t.Fatalf("GetEarningsDatesByCode() error = %v", err)
	}

	if len(dates) != 2 {
		t.Fatalf("GetEarningsDatesByCode() returned %d items, want 2", len(dates))
	}
	if dates[1].SchDate != "2025-08-05" {
		t.Errorf("GetEarningsDatesByCode() SchDate = %v, want 2025-08-05", dates[1].SchDate)
	}
	if mockClient.RequestCount != 2 {
		t.Errorf("GetEarningsDatesByCode() sent %d requests, want 2", mockClient.RequestCount)
	}
}

func TestEarningsDateService_GetEarningsDatesByScheduledDate(t *testing.T) {
	mockClient := client.NewMockClient()
	service := NewEarningsDateService(mockClient)

	mockClient.SetResponse("GET", "/fins/earnings-date?scheduled_date=2025-08-05", EarningsDateResponse{
		Data: []EarningsDate{
			{PubDate: "2025-07-01", SchDate: "2025-08-05", FQName: "1Q", Code: "86970"},
		},
	})

	dates, err := service.GetEarningsDatesByScheduledDate(context.Background(), "2025-08-05")
	if err != nil {
		t.Fatalf("GetEarningsDatesByScheduledDate() error = %v", err)
	}

	if len(dates) != 1 {
		t.Fatalf("GetEarningsDatesByScheduledDate() returned %d items, want 1", len(dates))
	}
}

func TestEarningsDateService_GetEarningsDates_Error(t *testing.T) {
	mockClient := client.NewMockClient()
	service := NewEarningsDateService(mockClient)

	mockClient.SetError("GET", "/fins/earnings-date?code=86970", fmt.Errorf("API error"))

	if _, err := service.GetEarningsDates(context.Background(), EarningsDateParams{Code: "86970"}); err == nil {
		t.Error("GetEarningsDates() expected error, got nil")
	}
}

func TestEarningsDate_IsUndetermined(t *testing.T) {
	// 予定日が未定に変更された場合、SchDateは空文字で届く
	jsonData := `{
		"data": [
			{
				"PubDate": "2025-06-03",
				"SchDate": "2025-07-30",
				"FQName": "1Q",
				"FYE": "0331",
				"Code": "86970",
				"CoName": "日本取引所グループ",
				"CoNameEn": "Japan Exchange Group,Inc."
			},
			{
				"PubDate": "2025-07-10",
				"SchDate": "",
				"FQName": "1Q",
				"FYE": "0331",
				"Code": "86970",
				"CoName": "日本取引所グループ",
				"CoNameEn": "Japan Exchange Group,Inc."
			}
		]
	}`

	var resp EarningsDateResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if resp.Data[0].IsUndetermined() {
		t.Error("IsUndetermined() = true, want false")
	}
	if !resp.Data[1].IsUndetermined() {
		t.Error("IsUndetermined() = false, want true")
	}
}
