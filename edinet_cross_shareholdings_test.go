package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestEdinetCrossShareholdingsService_GetCrossShareholdings(t *testing.T) {
	tests := []struct {
		name     string
		params   EdinetCrossShareholdingsParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with edinet code",
			params: EdinetCrossShareholdingsParams{
				EdinetCode: "E03814",
			},
			wantPath: "/edinet/cross-shareholdings?edinet_code=E03814",
		},
		{
			name: "with code",
			params: EdinetCrossShareholdingsParams{
				Code: "86970",
			},
			wantPath: "/edinet/cross-shareholdings?code=86970",
		},
		{
			name: "with date",
			params: EdinetCrossShareholdingsParams{
				Date: "20250620",
			},
			wantPath: "/edinet/cross-shareholdings?date=20250620",
		},
		{
			name: "with date in hyphen format",
			params: EdinetCrossShareholdingsParams{
				Date: "2025-06-20",
			},
			wantPath: "/edinet/cross-shareholdings?date=2025-06-20",
		},
		{
			name: "with code and date",
			params: EdinetCrossShareholdingsParams{
				Code: "86970",
				Date: "20250620",
			},
			wantPath: "/edinet/cross-shareholdings?code=86970&date=20250620",
		},
		{
			name: "with edinet code and date",
			params: EdinetCrossShareholdingsParams{
				EdinetCode: "E03814",
				Date:       "20250620",
			},
			wantPath: "/edinet/cross-shareholdings?edinet_code=E03814&date=20250620",
		},
		{
			name: "with pagination key",
			params: EdinetCrossShareholdingsParams{
				Code:          "86970",
				PaginationKey: "key123",
			},
			wantPath: "/edinet/cross-shareholdings?code=86970&pagination_key=key123",
		},
		{
			name:     "with no parameters",
			params:   EdinetCrossShareholdingsParams{},
			wantPath: "/edinet/cross-shareholdings",
		},
		{
			name: "with edinet code and code",
			params: EdinetCrossShareholdingsParams{
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
			service := NewEdinetCrossShareholdingsService(mockClient)

			// Mock response based on documentation sample
			mockResponse := EdinetCrossShareholdingsResponse{
				Data: []EdinetCrossShareholdingDoc{
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
						Report: &EdinetCrossShareholdingBlock{
							HldrName:         "株式会社日本取引所グループ",
							HldrCode:         "86970",
							HldrEdinetCode:   "E03814",
							NonListedIss:     6,
							NonListedBookVal: 1035000000,
						},
					},
				},
				PaginationKey: "",
			}
			if !tt.wantErr {
				mockClient.SetResponse("GET", tt.wantPath, mockResponse)
			}

			// Execute
			resp, err := service.GetCrossShareholdings(context.Background(), tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetCrossShareholdings() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetCrossShareholdings() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetCrossShareholdings() returned nil response")
				return
			}
			if len(resp.Data) == 0 {
				t.Error("GetCrossShareholdings() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetCrossShareholdings() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestEdinetCrossShareholdingsService_GetCrossShareholdings_RejectsEdinetCodeAndCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetCrossShareholdingsService(mockClient)

	// Execute with both edinet_code and code
	_, err := service.GetCrossShareholdings(context.Background(), EdinetCrossShareholdingsParams{
		EdinetCode: "E03814",
		Code:       "86970",
	})

	// Verify
	if err == nil {
		t.Error("GetCrossShareholdings() expected error for edinet_code and code but got nil")
	}
	if err.Error() != "edinet_code and code cannot be specified together" {
		t.Errorf("GetCrossShareholdings() error = %v, want 'edinet_code and code cannot be specified together'", err)
	}
	if mockClient.RequestCount != 0 {
		t.Errorf("GetCrossShareholdings() request count = %v, want 0", mockClient.RequestCount)
	}
}

func TestEdinetCrossShareholdingsService_GetCrossShareholdings_DecodesNestedBlocks(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetCrossShareholdingsService(mockClient)

	// JSON-shaped mock response（ブロックのnull、非開示マーカー、数値が文字列で返るケースを含む）
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
				"Report": {
					"HldrName": "株式会社日本取引所グループ",
					"HldrCode": "86970",
					"HldrEdinetCode": "E03814",
					"ListedIss": 0,
					"ListedBookVal": 0,
					"ListedIncIss": 0,
					"ListedIncAcqCost": 0,
					"ListedDecIss": 0,
					"ListedDecSaleAmt": 0,
					"ListedIncRsn": null,
					"NonListedIss": 6,
					"NonListedBookVal": "1035000000",
					"NonListedIncIss": 0,
					"NonListedIncAcqCost": 0,
					"NonListedDecIss": 0,
					"NonListedDecSaleAmt": 0,
					"NonListedIncRsn": null,
					"Spec": [
						{
							"IsrName": "株式会社サンプル銀行",
							"IsrCode": "56780",
							"IsrEdinetCode": "E05678",
							"CurShs": "1200000",
							"PriShs": 1200000,
							"CurBookVal": 850000000,
							"PriBookVal": 820000000,
							"CurShsNotDisc": null,
							"PriShsNotDisc": null,
							"CurBookValNotDisc": null,
							"PriBookValNotDisc": null,
							"HoldRat": "取引関係の維持・強化のため",
							"IsrHolds": "有",
							"IsrHoldsCode": "1"
						},
						{
							"IsrName": "株式会社サンプル電機",
							"IsrCode": null,
							"IsrEdinetCode": null,
							"CurShs": 500000,
							"PriShs": null,
							"CurBookVal": 350000000,
							"PriBookVal": null,
							"CurShsNotDisc": null,
							"PriShsNotDisc": "（注3）",
							"CurBookValNotDisc": null,
							"PriBookValNotDisc": "（注3）",
							"HoldRat": "議決権行使指図権を持つため",
							"IsrHolds": "無(注)３",
							"IsrHoldsCode": "0"
						}
					],
					"Deem": [],
					"SpecFn": "<p>※ 特定投資株式は、事業関係の維持・強化を目的として保有しています。</p>",
					"DeemFn": null
				},
				"Largest": null,
				"SecondLargest": {
					"HldrName": "株式会社東京証券取引所",
					"HldrCode": null,
					"HldrEdinetCode": null,
					"ListedIss": 0,
					"ListedBookVal": 0,
					"ListedIncIss": 0,
					"ListedIncAcqCost": 0,
					"ListedDecIss": 0,
					"ListedDecSaleAmt": 0,
					"ListedIncRsn": null,
					"NonListedIss": 2,
					"NonListedBookVal": 953000000,
					"NonListedIncIss": 0,
					"NonListedIncAcqCost": 0,
					"NonListedDecIss": 0,
					"NonListedDecSaleAmt": 0,
					"NonListedIncRsn": null,
					"Spec": [],
					"Deem": [],
					"SpecFn": null,
					"DeemFn": null
				}
			}
		],
		"pagination_key": ""
	}`)
	mockClient.SetResponse("GET", "/edinet/cross-shareholdings?edinet_code=E03814", mockJSON)

	// Execute
	resp, err := service.GetCrossShareholdings(context.Background(), EdinetCrossShareholdingsParams{EdinetCode: "E03814"})

	// Verify
	if err != nil {
		t.Fatalf("GetCrossShareholdings() error = %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("GetCrossShareholdings() returned %d items, want 1", len(resp.Data))
	}

	doc := resp.Data[0]
	if doc.DocId != "S100YA84" {
		t.Errorf("GetCrossShareholdings() DocId = %v, want S100YA84", doc.DocId)
	}
	if doc.DocTypeCode != "120" {
		t.Errorf("GetCrossShareholdings() DocTypeCode = %v, want 120", doc.DocTypeCode)
	}

	// ブロックの有無
	if !doc.HasReport() {
		t.Fatal("GetCrossShareholdings() HasReport() = false, want true")
	}
	if doc.HasLargest() {
		t.Error("GetCrossShareholdings() HasLargest() = true, want false")
	}
	if doc.Largest != nil {
		t.Error("GetCrossShareholdings() Largest should be nil")
	}
	if !doc.HasSecondLargest() {
		t.Fatal("GetCrossShareholdings() HasSecondLargest() = false, want true")
	}

	// 提出会社ブロック
	report := doc.Report
	if report.HldrName != "株式会社日本取引所グループ" {
		t.Errorf("GetCrossShareholdings() Report.HldrName = %v, want 株式会社日本取引所グループ", report.HldrName)
	}
	if report.NonListedIss != 6 {
		t.Errorf("GetCrossShareholdings() Report.NonListedIss = %v, want 6", report.NonListedIss)
	}
	// 数値が文字列で返るケース
	if report.NonListedBookVal != 1035000000 {
		t.Errorf("GetCrossShareholdings() Report.NonListedBookVal = %v, want 1035000000", report.NonListedBookVal)
	}
	// nullの文字列は空文字になる
	if report.ListedIncRsn != "" {
		t.Errorf("GetCrossShareholdings() Report.ListedIncRsn = %v, want empty", report.ListedIncRsn)
	}
	if report.DeemFn != "" {
		t.Errorf("GetCrossShareholdings() Report.DeemFn = %v, want empty", report.DeemFn)
	}
	if report.SpecFn == "" {
		t.Error("GetCrossShareholdings() Report.SpecFn should not be empty")
	}
	if len(report.Deem) != 0 {
		t.Errorf("GetCrossShareholdings() Report.Deem returned %d items, want 0", len(report.Deem))
	}
	if len(report.Spec) != 2 {
		t.Fatalf("GetCrossShareholdings() Report.Spec returned %d items, want 2", len(report.Spec))
	}

	// 開示されている銘柄（数値が文字列で返るケースを含む）
	issue1 := report.Spec[0]
	if issue1.IsrName != "株式会社サンプル銀行" {
		t.Errorf("GetCrossShareholdings() Spec[0].IsrName = %v, want 株式会社サンプル銀行", issue1.IsrName)
	}
	if issue1.IsrCode != "56780" {
		t.Errorf("GetCrossShareholdings() Spec[0].IsrCode = %v, want 56780", issue1.IsrCode)
	}
	if issue1.CurShs == nil || *issue1.CurShs != 1200000 {
		t.Errorf("GetCrossShareholdings() Spec[0].CurShs = %v, want 1200000", ptrToStr(issue1.CurShs))
	}
	if issue1.PriShs == nil || *issue1.PriShs != 1200000 {
		t.Errorf("GetCrossShareholdings() Spec[0].PriShs = %v, want 1200000", ptrToStr(issue1.PriShs))
	}
	if issue1.CurBookVal == nil || *issue1.CurBookVal != 850000000 {
		t.Errorf("GetCrossShareholdings() Spec[0].CurBookVal = %v, want 850000000", ptrToStr(issue1.CurBookVal))
	}
	if issue1.PriShsNotDisc != "" {
		t.Errorf("GetCrossShareholdings() Spec[0].PriShsNotDisc = %v, want empty", issue1.PriShsNotDisc)
	}
	if !issue1.IsDisclosed() {
		t.Error("GetCrossShareholdings() Spec[0].IsDisclosed() = false, want true")
	}
	if issue1.IsrHoldsCode != "1" {
		t.Errorf("GetCrossShareholdings() Spec[0].IsrHoldsCode = %v, want 1", issue1.IsrHoldsCode)
	}

	// 前期が非開示の銘柄（数値はnil、NotDiscにマーカー生値）
	issue2 := report.Spec[1]
	if issue2.IsrCode != "" {
		t.Errorf("GetCrossShareholdings() Spec[1].IsrCode = %v, want empty", issue2.IsrCode)
	}
	if issue2.CurShs == nil || *issue2.CurShs != 500000 {
		t.Errorf("GetCrossShareholdings() Spec[1].CurShs = %v, want 500000", ptrToStr(issue2.CurShs))
	}
	if issue2.PriShs != nil {
		t.Errorf("GetCrossShareholdings() Spec[1].PriShs = %v, want nil", ptrToStr(issue2.PriShs))
	}
	if issue2.PriBookVal != nil {
		t.Errorf("GetCrossShareholdings() Spec[1].PriBookVal = %v, want nil", ptrToStr(issue2.PriBookVal))
	}
	if issue2.PriShsNotDisc != "（注3）" {
		t.Errorf("GetCrossShareholdings() Spec[1].PriShsNotDisc = %v, want （注3）", issue2.PriShsNotDisc)
	}
	if issue2.PriBookValNotDisc != "（注3）" {
		t.Errorf("GetCrossShareholdings() Spec[1].PriBookValNotDisc = %v, want （注3）", issue2.PriBookValNotDisc)
	}
	if issue2.IsrHolds != "無(注)３" {
		t.Errorf("GetCrossShareholdings() Spec[1].IsrHolds = %v, want 無(注)３", issue2.IsrHolds)
	}

	// 連結第二最大保有会社ブロック（nullの文字列は空文字になる）
	secondLargest := doc.SecondLargest
	if secondLargest.HldrName != "株式会社東京証券取引所" {
		t.Errorf("GetCrossShareholdings() SecondLargest.HldrName = %v, want 株式会社東京証券取引所", secondLargest.HldrName)
	}
	if secondLargest.HldrCode != "" {
		t.Errorf("GetCrossShareholdings() SecondLargest.HldrCode = %v, want empty", secondLargest.HldrCode)
	}
	if secondLargest.NonListedBookVal != 953000000 {
		t.Errorf("GetCrossShareholdings() SecondLargest.NonListedBookVal = %v, want 953000000", secondLargest.NonListedBookVal)
	}
}

func TestEdinetCrossShareholdingsService_GetCrossShareholdingsByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetCrossShareholdingsService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := EdinetCrossShareholdingsResponse{
		Data: []EdinetCrossShareholdingDoc{
			{
				DocId:   "S100YA84",
				Code:    "86970",
				SubDate: "2026-06-11",
				Report: &EdinetCrossShareholdingBlock{
					HldrName:     "株式会社日本取引所グループ",
					NonListedIss: 6,
				},
			},
			{
				DocId:   "S100XB12",
				Code:    "86970",
				SubDate: "2025-06-12",
				Report: &EdinetCrossShareholdingBlock{
					HldrName:     "株式会社日本取引所グループ",
					NonListedIss: 5,
				},
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := EdinetCrossShareholdingsResponse{
		Data: []EdinetCrossShareholdingDoc{
			{
				DocId:   "S100WC34",
				Code:    "86970",
				SubDate: "2024-06-13",
				Report: &EdinetCrossShareholdingBlock{
					HldrName:     "株式会社日本取引所グループ",
					NonListedIss: 5,
				},
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/edinet/cross-shareholdings?code=86970", mockResponse1)
	mockClient.SetResponse("GET", "/edinet/cross-shareholdings?code=86970&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetCrossShareholdingsByCode(context.Background(), "86970")

	// Verify
	if err != nil {
		t.Fatalf("GetCrossShareholdingsByCode() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetCrossShareholdingsByCode() returned %d items, want 3", len(data))
	}
	for _, item := range data {
		if item.Code != "86970" {
			t.Errorf("GetCrossShareholdingsByCode() returned code %v, want 86970", item.Code)
		}
	}
}

func TestEdinetCrossShareholdingsService_GetCrossShareholdingsByEdinetCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetCrossShareholdingsService(mockClient)

	// Mock response
	mockResponse := EdinetCrossShareholdingsResponse{
		Data: []EdinetCrossShareholdingDoc{
			{
				DocId:      "S100YA84",
				Code:       "86970",
				EdinetCode: "E03814",
				SubDate:    "2026-06-11",
				Report: &EdinetCrossShareholdingBlock{
					HldrName:     "株式会社日本取引所グループ",
					NonListedIss: 6,
				},
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/edinet/cross-shareholdings?edinet_code=E03814", mockResponse)

	// Execute
	data, err := service.GetCrossShareholdingsByEdinetCode(context.Background(), "E03814")

	// Verify
	if err != nil {
		t.Fatalf("GetCrossShareholdingsByEdinetCode() error = %v", err)
	}
	if len(data) != 1 {
		t.Errorf("GetCrossShareholdingsByEdinetCode() returned %d items, want 1", len(data))
	}
	if data[0].EdinetCode != "E03814" {
		t.Errorf("GetCrossShareholdingsByEdinetCode() returned edinet code %v, want E03814", data[0].EdinetCode)
	}
}

func TestEdinetCrossShareholdingsService_GetCrossShareholdingsByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetCrossShareholdingsService(mockClient)

	// Mock response
	mockResponse := EdinetCrossShareholdingsResponse{
		Data: []EdinetCrossShareholdingDoc{
			{
				DocId:   "S100YA84",
				Code:    "86970",
				SubDate: "2025-06-20",
				Report: &EdinetCrossShareholdingBlock{
					HldrName:     "株式会社日本取引所グループ",
					NonListedIss: 6,
				},
			},
			{
				DocId:   "S100YB56",
				Code:    "13660",
				SubDate: "2025-06-20",
				Report: &EdinetCrossShareholdingBlock{
					HldrName:  "株式会社サンプル",
					ListedIss: 3,
					Spec: []EdinetCrossShareholdingIssue{
						{
							IsrName: "株式会社サンプル銀行",
							IsrCode: "56780",
							CurShs:  floatPtr(1200000),
						},
					},
				},
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/edinet/cross-shareholdings?date=20250620", mockResponse)

	// Execute
	data, err := service.GetCrossShareholdingsByDate(context.Background(), "20250620")

	// Verify
	if err != nil {
		t.Fatalf("GetCrossShareholdingsByDate() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetCrossShareholdingsByDate() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.SubDate != "2025-06-20" {
			t.Errorf("GetCrossShareholdingsByDate() returned sub date %v, want 2025-06-20", item.SubDate)
		}
	}
}

func TestEdinetCrossShareholdingIssue_IsDisclosed(t *testing.T) {
	tests := []struct {
		name  string
		issue EdinetCrossShareholdingIssue
		want  bool
	}{
		{
			name: "disclosed",
			issue: EdinetCrossShareholdingIssue{
				CurShs: floatPtr(1200000),
			},
			want: true,
		},
		{
			name: "not disclosed",
			issue: EdinetCrossShareholdingIssue{
				CurShs:        nil,
				CurShsNotDisc: "※",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.issue.IsDisclosed(); got != tt.want {
				t.Errorf("EdinetCrossShareholdingIssue.IsDisclosed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEdinetCrossShareholdingsService_GetCrossShareholdings_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetCrossShareholdingsService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/edinet/cross-shareholdings?code=86970", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetCrossShareholdings(context.Background(), EdinetCrossShareholdingsParams{Code: "86970"})

	// Verify
	if err == nil {
		t.Error("GetCrossShareholdings() expected error but got nil")
	}
}
