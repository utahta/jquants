package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestEdinetMajorShareholdersService_GetMajorShareholders(t *testing.T) {
	tests := []struct {
		name     string
		params   EdinetMajorShareholdersParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with edinet code",
			params: EdinetMajorShareholdersParams{
				EdinetCode: "E03814",
			},
			wantPath: "/edinet/major-shareholders?edinet_code=E03814",
		},
		{
			name: "with code",
			params: EdinetMajorShareholdersParams{
				Code: "86970",
			},
			wantPath: "/edinet/major-shareholders?code=86970",
		},
		{
			name: "with date",
			params: EdinetMajorShareholdersParams{
				Date: "20250620",
			},
			wantPath: "/edinet/major-shareholders?date=20250620",
		},
		{
			name: "with date in hyphen format",
			params: EdinetMajorShareholdersParams{
				Date: "2025-06-20",
			},
			wantPath: "/edinet/major-shareholders?date=2025-06-20",
		},
		{
			name: "with code and date",
			params: EdinetMajorShareholdersParams{
				Code: "86970",
				Date: "20250620",
			},
			wantPath: "/edinet/major-shareholders?code=86970&date=20250620",
		},
		{
			name: "with edinet code and date",
			params: EdinetMajorShareholdersParams{
				EdinetCode: "E03814",
				Date:       "20250620",
			},
			wantPath: "/edinet/major-shareholders?edinet_code=E03814&date=20250620",
		},
		{
			name: "with pagination key",
			params: EdinetMajorShareholdersParams{
				Code:          "86970",
				PaginationKey: "key123",
			},
			wantPath: "/edinet/major-shareholders?code=86970&pagination_key=key123",
		},
		{
			name:     "with no parameters",
			params:   EdinetMajorShareholdersParams{},
			wantPath: "/edinet/major-shareholders",
		},
		{
			name: "with edinet code and code",
			params: EdinetMajorShareholdersParams{
				EdinetCode: "E03814",
				Code:       "86970",
			},
			wantPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewEdinetMajorShareholdersService(mockClient)

			// Mock response based on documentation sample
			mockResponse := EdinetMajorShareholdersResponse{
				Data: []EdinetMajorShareholderDoc{
					{
						DocId:       "S100YA84",
						Code:        "86970",
						EdinetCode:  "E03814",
						FilerName:   "株式会社日本取引所グループ",
						FilerNameEn: "Japan Exchange Group, Inc.",
						DocTypeCode: "120",
						SubDate:     "2026-06-11",
						SubTime:     "15:00:00",
						PerSt:       "2025-04-01",
						PerEn:       "2026-03-31",
						Hldrs: []EdinetMajorShareholderHolder{
							{
								Rank:     1,
								HldrName: "日本マスタートラスト信託銀行株式会社（信託口）",
								HldrAddr: "東京都港区赤坂１丁目８番１号　赤坂インターシティＡＩＲ",
								ShsHeld:  175830000,
								ShsRatio: 0.1704,
							},
						},
					},
				},
				PaginationKey: "",
			}
			if !tt.wantErr {
				mockClient.SetResponse("GET", tt.wantPath, mockResponse)
			}

			// Execute
			resp, err := service.GetMajorShareholders(context.Background(), tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetMajorShareholders() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetMajorShareholders() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetMajorShareholders() returned nil response")
				return
			}
			if len(resp.Data) == 0 {
				t.Error("GetMajorShareholders() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetMajorShareholders() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestEdinetMajorShareholdersService_GetMajorShareholders_RejectsEdinetCodeAndCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetMajorShareholdersService(mockClient)

	// Execute with both edinet_code and code
	_, err := service.GetMajorShareholders(context.Background(), EdinetMajorShareholdersParams{
		EdinetCode: "E03814",
		Code:       "86970",
	})

	// Verify
	if err == nil {
		t.Error("GetMajorShareholders() expected error for edinet_code and code but got nil")
	}
	if err.Error() != "edinet_code and code cannot be specified together" {
		t.Errorf("GetMajorShareholders() error = %v, want 'edinet_code and code cannot be specified together'", err)
	}
	if mockClient.RequestCount != 0 {
		t.Errorf("GetMajorShareholders() request count = %v, want 0", mockClient.RequestCount)
	}
}

func TestEdinetMajorShareholdersService_GetMajorShareholders_DecodesNestedHolders(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetMajorShareholdersService(mockClient)

	// JSON-shaped mock response（数値が文字列で返るケースを含む）
	mockJSON := json.RawMessage(`{
		"data": [
			{
				"DocId": "S100YA84",
				"Code": "86970",
				"EdinetCode": "E03814",
				"FilerName": "株式会社日本取引所グループ",
				"FilerNameEn": "Japan Exchange Group, Inc.",
				"DocTypeCode": "120",
				"SubDate": "2026-06-11",
				"SubTime": "15:00:00",
				"PerSt": "2025-04-01",
				"PerEn": "2026-03-31",
				"Hldrs": [
					{
						"Rank": 1,
						"HldrName": "日本マスタートラスト信託銀行株式会社（信託口）",
						"HldrAddr": "東京都港区赤坂１丁目８番１号　赤坂インターシティＡＩＲ",
						"ShsHeld": 175830000,
						"ShsRatio": 0.1704
					},
					{
						"Rank": 2,
						"HldrName": "株式会社日本カストディ銀行（信託口）",
						"HldrAddr": "東京都中央区晴海１丁目８－１２",
						"ShsHeld": "56970000",
						"ShsRatio": "0.0552"
					}
				]
			}
		],
		"pagination_key": ""
	}`)
	mockClient.SetResponse("GET", "/edinet/major-shareholders?edinet_code=E03814", mockJSON)

	// Execute
	resp, err := service.GetMajorShareholders(context.Background(), EdinetMajorShareholdersParams{EdinetCode: "E03814"})

	// Verify
	if err != nil {
		t.Fatalf("GetMajorShareholders() error = %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("GetMajorShareholders() returned %d items, want 1", len(resp.Data))
	}

	doc := resp.Data[0]
	if doc.DocId != "S100YA84" {
		t.Errorf("GetMajorShareholders() DocId = %v, want S100YA84", doc.DocId)
	}
	if doc.Code != "86970" {
		t.Errorf("GetMajorShareholders() Code = %v, want 86970", doc.Code)
	}
	if doc.DocTypeCode != "120" {
		t.Errorf("GetMajorShareholders() DocTypeCode = %v, want 120", doc.DocTypeCode)
	}
	if len(doc.Hldrs) != 2 {
		t.Fatalf("GetMajorShareholders() returned %d holders, want 2", len(doc.Hldrs))
	}

	holder1 := doc.Hldrs[0]
	if holder1.Rank != 1 {
		t.Errorf("GetMajorShareholders() Hldrs[0].Rank = %v, want 1", holder1.Rank)
	}
	if holder1.HldrName != "日本マスタートラスト信託銀行株式会社（信託口）" {
		t.Errorf("GetMajorShareholders() Hldrs[0].HldrName = %v, want 日本マスタートラスト信託銀行株式会社（信託口）", holder1.HldrName)
	}
	if holder1.ShsHeld != 175830000 {
		t.Errorf("GetMajorShareholders() Hldrs[0].ShsHeld = %v, want 175830000", holder1.ShsHeld)
	}
	if holder1.ShsRatio != 0.1704 {
		t.Errorf("GetMajorShareholders() Hldrs[0].ShsRatio = %v, want 0.1704", holder1.ShsRatio)
	}

	// 数値が文字列で返るケース
	holder2 := doc.Hldrs[1]
	if holder2.Rank != 2 {
		t.Errorf("GetMajorShareholders() Hldrs[1].Rank = %v, want 2", holder2.Rank)
	}
	if holder2.ShsHeld != 56970000 {
		t.Errorf("GetMajorShareholders() Hldrs[1].ShsHeld = %v, want 56970000", holder2.ShsHeld)
	}
	if holder2.ShsRatio != 0.0552 {
		t.Errorf("GetMajorShareholders() Hldrs[1].ShsRatio = %v, want 0.0552", holder2.ShsRatio)
	}
}

func TestEdinetMajorShareholdersService_GetMajorShareholdersByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetMajorShareholdersService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := EdinetMajorShareholdersResponse{
		Data: []EdinetMajorShareholderDoc{
			{
				DocId:   "S100YA84",
				Code:    "86970",
				SubDate: "2026-06-11",
				Hldrs: []EdinetMajorShareholderHolder{
					{Rank: 1, HldrName: "日本マスタートラスト信託銀行株式会社（信託口）", ShsHeld: 175830000, ShsRatio: 0.1704},
				},
			},
			{
				DocId:   "S100XB12",
				Code:    "86970",
				SubDate: "2025-06-12",
				Hldrs: []EdinetMajorShareholderHolder{
					{Rank: 1, HldrName: "日本マスタートラスト信託銀行株式会社（信託口）", ShsHeld: 170000000, ShsRatio: 0.1650},
				},
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := EdinetMajorShareholdersResponse{
		Data: []EdinetMajorShareholderDoc{
			{
				DocId:   "S100WC34",
				Code:    "86970",
				SubDate: "2024-06-13",
				Hldrs: []EdinetMajorShareholderHolder{
					{Rank: 1, HldrName: "日本マスタートラスト信託銀行株式会社（信託口）", ShsHeld: 165000000, ShsRatio: 0.1600},
				},
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/edinet/major-shareholders?code=86970", mockResponse1)
	mockClient.SetResponse("GET", "/edinet/major-shareholders?code=86970&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetMajorShareholdersByCode(context.Background(), "86970")

	// Verify
	if err != nil {
		t.Fatalf("GetMajorShareholdersByCode() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetMajorShareholdersByCode() returned %d items, want 3", len(data))
	}
	for _, item := range data {
		if item.Code != "86970" {
			t.Errorf("GetMajorShareholdersByCode() returned code %v, want 86970", item.Code)
		}
	}
}

func TestEdinetMajorShareholdersService_GetMajorShareholdersByEdinetCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetMajorShareholdersService(mockClient)

	// Mock response
	mockResponse := EdinetMajorShareholdersResponse{
		Data: []EdinetMajorShareholderDoc{
			{
				DocId:      "S100YA84",
				Code:       "86970",
				EdinetCode: "E03814",
				SubDate:    "2026-06-11",
				Hldrs: []EdinetMajorShareholderHolder{
					{Rank: 1, HldrName: "日本マスタートラスト信託銀行株式会社（信託口）", ShsHeld: 175830000, ShsRatio: 0.1704},
				},
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/edinet/major-shareholders?edinet_code=E03814", mockResponse)

	// Execute
	data, err := service.GetMajorShareholdersByEdinetCode(context.Background(), "E03814")

	// Verify
	if err != nil {
		t.Fatalf("GetMajorShareholdersByEdinetCode() error = %v", err)
	}
	if len(data) != 1 {
		t.Errorf("GetMajorShareholdersByEdinetCode() returned %d items, want 1", len(data))
	}
	if data[0].EdinetCode != "E03814" {
		t.Errorf("GetMajorShareholdersByEdinetCode() returned edinet code %v, want E03814", data[0].EdinetCode)
	}
}

func TestEdinetMajorShareholdersService_GetMajorShareholdersByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetMajorShareholdersService(mockClient)

	// Mock response
	mockResponse := EdinetMajorShareholdersResponse{
		Data: []EdinetMajorShareholderDoc{
			{
				DocId:   "S100YA84",
				Code:    "86970",
				SubDate: "2025-06-20",
				Hldrs: []EdinetMajorShareholderHolder{
					{Rank: 1, HldrName: "日本マスタートラスト信託銀行株式会社（信託口）", ShsHeld: 175830000, ShsRatio: 0.1704},
				},
			},
			{
				DocId:   "S100YB56",
				Code:    "13660",
				SubDate: "2025-06-20",
				Hldrs: []EdinetMajorShareholderHolder{
					{Rank: 1, HldrName: "株式会社日本カストディ銀行（信託口）", ShsHeld: 56970000, ShsRatio: 0.0552},
				},
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/edinet/major-shareholders?date=20250620", mockResponse)

	// Execute
	data, err := service.GetMajorShareholdersByDate(context.Background(), "20250620")

	// Verify
	if err != nil {
		t.Fatalf("GetMajorShareholdersByDate() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetMajorShareholdersByDate() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.SubDate != "2025-06-20" {
			t.Errorf("GetMajorShareholdersByDate() returned sub date %v, want 2025-06-20", item.SubDate)
		}
	}
}

func TestEdinetMajorShareholderHolder_GetShareholdingPercentage(t *testing.T) {
	holder := EdinetMajorShareholderHolder{
		ShsRatio: 0.1704,
	}

	expected := 17.04
	got := holder.GetShareholdingPercentage()

	if got != expected {
		t.Errorf("EdinetMajorShareholderHolder.GetShareholdingPercentage() = %v, want %v", got, expected)
	}
}

func TestEdinetMajorShareholdersService_GetMajorShareholders_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetMajorShareholdersService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/edinet/major-shareholders?code=86970", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetMajorShareholders(context.Background(), EdinetMajorShareholdersParams{Code: "86970"})

	// Verify
	if err == nil {
		t.Error("GetMajorShareholders() expected error but got nil")
	}
}
