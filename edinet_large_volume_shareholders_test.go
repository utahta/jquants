package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestEdinetLargeVolumeShareholdersService_GetLargeVolumeShareholders(t *testing.T) {
	tests := []struct {
		name     string
		params   EdinetLargeVolumeShareholdersParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with edinet code",
			params: EdinetLargeVolumeShareholdersParams{
				EdinetCode: "E03814",
			},
			wantPath: "/edinet/large-volume-shareholders?edinet_code=E03814",
		},
		{
			name: "with code",
			params: EdinetLargeVolumeShareholdersParams{
				Code: "86970",
			},
			wantPath: "/edinet/large-volume-shareholders?code=86970",
		},
		{
			name: "with date",
			params: EdinetLargeVolumeShareholdersParams{
				Date: "20250707",
			},
			wantPath: "/edinet/large-volume-shareholders?date=20250707",
		},
		{
			name: "with date in hyphen format",
			params: EdinetLargeVolumeShareholdersParams{
				Date: "2025-07-07",
			},
			wantPath: "/edinet/large-volume-shareholders?date=2025-07-07",
		},
		{
			name: "with code and date",
			params: EdinetLargeVolumeShareholdersParams{
				Code: "86970",
				Date: "20250707",
			},
			wantPath: "/edinet/large-volume-shareholders?code=86970&date=20250707",
		},
		{
			name: "with edinet code and date",
			params: EdinetLargeVolumeShareholdersParams{
				EdinetCode: "E03814",
				Date:       "20250707",
			},
			wantPath: "/edinet/large-volume-shareholders?edinet_code=E03814&date=20250707",
		},
		{
			name: "with pagination key",
			params: EdinetLargeVolumeShareholdersParams{
				Code:          "86970",
				PaginationKey: "key123",
			},
			wantPath: "/edinet/large-volume-shareholders?code=86970&pagination_key=key123",
		},
		{
			name:     "with no parameters",
			params:   EdinetLargeVolumeShareholdersParams{},
			wantPath: "/edinet/large-volume-shareholders",
		},
		{
			name: "with edinet code and code",
			params: EdinetLargeVolumeShareholdersParams{
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
			service := NewEdinetLargeVolumeShareholdersService(mockClient)

			// Mock response based on documentation sample
			ratioLast := 0.0614
			mockResponse := EdinetLargeVolumeShareholdersResponse{
				Data: []EdinetLargeVolumeShareholderDoc{
					{
						DocId:             "S100WBIV",
						Code:              "86970",
						EdinetCode:        "E03814",
						IsrName:           "株式会社日本取引所グループ",
						DocTypeCode:       "350",
						SubDate:           "2025-07-07",
						SubTime:           "12:09:00",
						LargeHldgTypeCode: "5",
						DocTitle:          "変更報告書ＮＯ．9",
						ChgRsn:            "・株券等保有割合の１％以上の増加",
						TotalShsHeld:      floatPtr(76018630),
						TotalShsRatio:     floatPtr(0.0728),
						TotalShsRatioLast: &ratioLast,
						TotalOutStks:      1044578366,
						Hldrs: []EdinetLargeVolumeShareholderHolder{
							{
								HldrName:          "サンプル・アセットマネジメント株式会社",
								HldrEdinetCode:    "E99990",
								LargeHldrTypeCode: "2",
								ShsHeld:           59555500,
								ShsRatio:          0.057,
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
			resp, err := service.GetLargeVolumeShareholders(context.Background(), tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetLargeVolumeShareholders() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetLargeVolumeShareholders() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetLargeVolumeShareholders() returned nil response")
				return
			}
			if len(resp.Data) == 0 {
				t.Error("GetLargeVolumeShareholders() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetLargeVolumeShareholders() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestEdinetLargeVolumeShareholdersService_GetLargeVolumeShareholders_RejectsEdinetCodeAndCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetLargeVolumeShareholdersService(mockClient)

	// Execute with both edinet_code and code
	_, err := service.GetLargeVolumeShareholders(context.Background(), EdinetLargeVolumeShareholdersParams{
		EdinetCode: "E03814",
		Code:       "86970",
	})

	// Verify
	if err == nil {
		t.Error("GetLargeVolumeShareholders() expected error for edinet_code and code but got nil")
	}
	if err.Error() != "edinet_code and code cannot be specified together" {
		t.Errorf("GetLargeVolumeShareholders() error = %v, want 'edinet_code and code cannot be specified together'", err)
	}
	if mockClient.RequestCount != 0 {
		t.Errorf("GetLargeVolumeShareholders() request count = %v, want 0", mockClient.RequestCount)
	}
}

func TestEdinetLargeVolumeShareholdersService_GetLargeVolumeShareholders_DecodesNestedHolders(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetLargeVolumeShareholdersService(mockClient)

	// JSON-shaped mock response（null値・数値が文字列で返るケースを含む）
	mockJSON := json.RawMessage(`{
		"data": [
			{
				"DocId": "S100WBIV",
				"Code": "86970",
				"EdinetCode": "E03814",
				"IsrName": "株式会社日本取引所グループ",
				"DocTypeCode": "350",
				"SubDate": "2025-07-07",
				"SubTime": "12:09:00",
				"LargeHldgTypeCode": "5",
				"DocTitle": "変更報告書ＮＯ．9",
				"ChgRsn": "・株券等保有割合の１％以上の増加",
				"TotalShsHeld": "76018630",
				"TotalShsRatio": 0.0728,
				"TotalShsRatioLast": 0.0614,
				"TotalOutStks": 1044578366,
				"Hldrs": [
					{
						"HldrName": "サンプル・アセットマネジメント株式会社",
						"HldrNameEn": "Sample Asset Management Co., Ltd.",
						"HldrEdinetCode": "E99990",
						"HldrCode": null,
						"LargeHldrTypeCode": "2",
						"LargeHldrTypeRaw": "法人(株式会社)",
						"HldgPurp": "信託財産の運用として保有している。",
						"ImpProp": null,
						"ColAgr": null,
						"ShsHeld": 59555500,
						"ShsRatio": 0.057,
						"ShsRatioLast": 0.0506,
						"OwnFund": null,
						"TotalBrw": null,
						"TotalOther": null,
						"OtherBrk": null,
						"TotalFund": null,
						"AcqDisp": [],
						"BrwList": [],
						"CredList": []
					},
					{
						"HldrName": "サンプル証券株式会社",
						"HldrNameEn": "Sample Securities Co., Ltd.",
						"HldrEdinetCode": "E99991",
						"HldrCode": null,
						"LargeHldrTypeCode": "2",
						"LargeHldrTypeRaw": "法人(株式会社)",
						"HldgPurp": "証券業務に係る商品在庫として保有している。",
						"ImpProp": null,
						"ColAgr": "消費貸借契約により、サンプル信託銀行株式会社から1,000,000株借入れている。",
						"ShsHeld": 8893542,
						"ShsRatio": 0.0085,
						"ShsRatioLast": 0.0095,
						"OwnFund": 300000000,
						"TotalBrw": 500000000,
						"TotalOther": null,
						"OtherBrk": null,
						"TotalFund": 800000000,
						"AcqDisp": [
							{
								"Date": "2025-06-20",
								"SecType": "普通株式",
								"Shs": 100000,
								"Ratio": 0.01,
								"Mkt": "市場内",
								"MktCode": "1",
								"TxnType": "取得",
								"TxnTypeCode": "1",
								"Cptty": null,
								"Price": 3800,
								"PriceRaw": null
							},
							{
								"Date": "2025-06-23",
								"SecType": "普通株式",
								"Shs": "50000",
								"Ratio": null,
								"Mkt": "市場外",
								"MktCode": "2",
								"TxnType": "処分",
								"TxnTypeCode": "2",
								"Cptty": null,
								"Price": null,
								"PriceRaw": "3,800円"
							}
						],
						"BrwList": [
							{
								"Name": "サンプル銀行株式会社",
								"Ind": "銀行",
								"Rep": "代表取締役 見本 太郎",
								"Addr": "東京都千代田区丸の内一丁目1番1号",
								"DiscBrwPurp": "2",
								"Amt": 500000000
							}
						],
						"CredList": [
							{
								"Name": "サンプル信託銀行株式会社",
								"Rep": "代表取締役 例示 花子",
								"Addr": "東京都千代田区大手町一丁目1番1号"
							}
						]
					}
				]
			}
		],
		"pagination_key": ""
	}`)
	mockClient.SetResponse("GET", "/edinet/large-volume-shareholders?edinet_code=E03814", mockJSON)

	// Execute
	resp, err := service.GetLargeVolumeShareholders(context.Background(), EdinetLargeVolumeShareholdersParams{EdinetCode: "E03814"})

	// Verify
	if err != nil {
		t.Fatalf("GetLargeVolumeShareholders() error = %v", err)
	}
	if len(resp.Data) != 1 {
		t.Fatalf("GetLargeVolumeShareholders() returned %d items, want 1", len(resp.Data))
	}

	doc := resp.Data[0]
	if doc.DocId != "S100WBIV" {
		t.Errorf("GetLargeVolumeShareholders() DocId = %v, want S100WBIV", doc.DocId)
	}
	if doc.LargeHldgTypeCode != LargeHldgTypeCodeSpecialChangeReport {
		t.Errorf("GetLargeVolumeShareholders() LargeHldgTypeCode = %v, want 5", doc.LargeHldgTypeCode)
	}
	if !doc.IsChangeReport() {
		t.Error("GetLargeVolumeShareholders() IsChangeReport() = false, want true")
	}
	// 数値が文字列で返るケース
	if doc.TotalShsHeld == nil || *doc.TotalShsHeld != 76018630 {
		t.Errorf("GetLargeVolumeShareholders() TotalShsHeld = %v, want 76018630", ptrToStr(doc.TotalShsHeld))
	}
	if doc.TotalShsRatio == nil || *doc.TotalShsRatio != 0.0728 {
		t.Errorf("GetLargeVolumeShareholders() TotalShsRatio = %v, want 0.0728", ptrToStr(doc.TotalShsRatio))
	}
	if doc.TotalShsRatioLast == nil || *doc.TotalShsRatioLast != 0.0614 {
		t.Errorf("GetLargeVolumeShareholders() TotalShsRatioLast = %v, want 0.0614", doc.TotalShsRatioLast)
	}
	if doc.TotalOutStks != 1044578366 {
		t.Errorf("GetLargeVolumeShareholders() TotalOutStks = %v, want 1044578366", doc.TotalOutStks)
	}
	if len(doc.Hldrs) != 2 {
		t.Fatalf("GetLargeVolumeShareholders() returned %d holders, want 2", len(doc.Hldrs))
	}

	// null資金の保有者（nilポインタになること）
	holder1 := doc.Hldrs[0]
	if holder1.HldrName != "サンプル・アセットマネジメント株式会社" {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[0].HldrName = %v, want サンプル・アセットマネジメント株式会社", holder1.HldrName)
	}
	if holder1.HldrCode != "" {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[0].HldrCode = %v, want empty", holder1.HldrCode)
	}
	if holder1.ShsHeld != 59555500 {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[0].ShsHeld = %v, want 59555500", holder1.ShsHeld)
	}
	if holder1.ShsRatioLast == nil || *holder1.ShsRatioLast != 0.0506 {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[0].ShsRatioLast = %v, want 0.0506", holder1.ShsRatioLast)
	}
	if holder1.OwnFund != nil {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[0].OwnFund = %v, want nil", holder1.OwnFund)
	}
	if holder1.TotalBrw != nil {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[0].TotalBrw = %v, want nil", holder1.TotalBrw)
	}
	if holder1.TotalOther != nil {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[0].TotalOther = %v, want nil", holder1.TotalOther)
	}
	if holder1.TotalFund != nil {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[0].TotalFund = %v, want nil", holder1.TotalFund)
	}
	if len(holder1.AcqDisp) != 0 {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[0].AcqDisp has %d items, want 0", len(holder1.AcqDisp))
	}

	// 資金・明細ありの保有者
	holder2 := doc.Hldrs[1]
	if holder2.OwnFund == nil || *holder2.OwnFund != 300000000 {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[1].OwnFund = %v, want 300000000", holder2.OwnFund)
	}
	if holder2.TotalBrw == nil || *holder2.TotalBrw != 500000000 {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[1].TotalBrw = %v, want 500000000", holder2.TotalBrw)
	}
	if holder2.TotalOther != nil {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[1].TotalOther = %v, want nil", holder2.TotalOther)
	}
	if holder2.TotalFund == nil || *holder2.TotalFund != 800000000 {
		t.Errorf("GetLargeVolumeShareholders() Hldrs[1].TotalFund = %v, want 800000000", holder2.TotalFund)
	}
	if len(holder2.AcqDisp) != 2 {
		t.Fatalf("GetLargeVolumeShareholders() Hldrs[1].AcqDisp has %d items, want 2", len(holder2.AcqDisp))
	}

	// 単価が数値で返るケース
	acqDisp1 := holder2.AcqDisp[0]
	if acqDisp1.Shs != 100000 {
		t.Errorf("GetLargeVolumeShareholders() AcqDisp[0].Shs = %v, want 100000", acqDisp1.Shs)
	}
	if acqDisp1.Ratio == nil || *acqDisp1.Ratio != 0.01 {
		t.Errorf("GetLargeVolumeShareholders() AcqDisp[0].Ratio = %v, want 0.01", acqDisp1.Ratio)
	}
	if acqDisp1.Price == nil || *acqDisp1.Price != 3800 {
		t.Errorf("GetLargeVolumeShareholders() AcqDisp[0].Price = %v, want 3800", acqDisp1.Price)
	}
	if acqDisp1.PriceRaw != "" {
		t.Errorf("GetLargeVolumeShareholders() AcqDisp[0].PriceRaw = %v, want empty", acqDisp1.PriceRaw)
	}

	// 単価を数値化できず生値のみ返るケース（数量が文字列で返るケースを含む）
	acqDisp2 := holder2.AcqDisp[1]
	if acqDisp2.Shs != 50000 {
		t.Errorf("GetLargeVolumeShareholders() AcqDisp[1].Shs = %v, want 50000", acqDisp2.Shs)
	}
	if acqDisp2.Ratio != nil {
		t.Errorf("GetLargeVolumeShareholders() AcqDisp[1].Ratio = %v, want nil", acqDisp2.Ratio)
	}
	if acqDisp2.Price != nil {
		t.Errorf("GetLargeVolumeShareholders() AcqDisp[1].Price = %v, want nil", acqDisp2.Price)
	}
	if acqDisp2.PriceRaw != "3,800円" {
		t.Errorf("GetLargeVolumeShareholders() AcqDisp[1].PriceRaw = %v, want 3,800円", acqDisp2.PriceRaw)
	}

	// 借入金の内訳・借入先
	if len(holder2.BrwList) != 1 {
		t.Fatalf("GetLargeVolumeShareholders() Hldrs[1].BrwList has %d items, want 1", len(holder2.BrwList))
	}
	borrowing := holder2.BrwList[0]
	if borrowing.Name != "サンプル銀行株式会社" {
		t.Errorf("GetLargeVolumeShareholders() BrwList[0].Name = %v, want サンプル銀行株式会社", borrowing.Name)
	}
	if borrowing.Amt == nil || *borrowing.Amt != 500000000 {
		t.Errorf("GetLargeVolumeShareholders() BrwList[0].Amt = %v, want 500000000", borrowing.Amt)
	}
	if len(holder2.CredList) != 1 {
		t.Fatalf("GetLargeVolumeShareholders() Hldrs[1].CredList has %d items, want 1", len(holder2.CredList))
	}
	if holder2.CredList[0].Name != "サンプル信託銀行株式会社" {
		t.Errorf("GetLargeVolumeShareholders() CredList[0].Name = %v, want サンプル信託銀行株式会社", holder2.CredList[0].Name)
	}
}

func TestEdinetLargeVolumeShareholdersService_GetLargeVolumeShareholdersByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetLargeVolumeShareholdersService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := EdinetLargeVolumeShareholdersResponse{
		Data: []EdinetLargeVolumeShareholderDoc{
			{
				DocId:             "S100WBIV",
				Code:              "86970",
				SubDate:           "2025-07-07",
				LargeHldgTypeCode: "5",
				TotalShsRatio:     floatPtr(0.0728),
			},
			{
				DocId:             "S100WA01",
				Code:              "86970",
				SubDate:           "2025-06-20",
				LargeHldgTypeCode: "2",
				TotalShsRatio:     floatPtr(0.0614),
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := EdinetLargeVolumeShareholdersResponse{
		Data: []EdinetLargeVolumeShareholderDoc{
			{
				DocId:             "S100W901",
				Code:              "86970",
				SubDate:           "2025-05-15",
				LargeHldgTypeCode: "1",
				TotalShsRatio:     floatPtr(0.0512),
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/edinet/large-volume-shareholders?code=86970", mockResponse1)
	mockClient.SetResponse("GET", "/edinet/large-volume-shareholders?code=86970&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetLargeVolumeShareholdersByCode(context.Background(), "86970")

	// Verify
	if err != nil {
		t.Fatalf("GetLargeVolumeShareholdersByCode() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetLargeVolumeShareholdersByCode() returned %d items, want 3", len(data))
	}
	for _, item := range data {
		if item.Code != "86970" {
			t.Errorf("GetLargeVolumeShareholdersByCode() returned code %v, want 86970", item.Code)
		}
	}
}

func TestEdinetLargeVolumeShareholdersService_GetLargeVolumeShareholdersByEdinetCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetLargeVolumeShareholdersService(mockClient)

	// Mock response
	mockResponse := EdinetLargeVolumeShareholdersResponse{
		Data: []EdinetLargeVolumeShareholderDoc{
			{
				DocId:             "S100WBIV",
				Code:              "86970",
				EdinetCode:        "E03814",
				SubDate:           "2025-07-07",
				LargeHldgTypeCode: "5",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/edinet/large-volume-shareholders?edinet_code=E03814", mockResponse)

	// Execute
	data, err := service.GetLargeVolumeShareholdersByEdinetCode(context.Background(), "E03814")

	// Verify
	if err != nil {
		t.Fatalf("GetLargeVolumeShareholdersByEdinetCode() error = %v", err)
	}
	if len(data) != 1 {
		t.Errorf("GetLargeVolumeShareholdersByEdinetCode() returned %d items, want 1", len(data))
	}
	if data[0].EdinetCode != "E03814" {
		t.Errorf("GetLargeVolumeShareholdersByEdinetCode() returned edinet code %v, want E03814", data[0].EdinetCode)
	}
}

func TestEdinetLargeVolumeShareholdersService_GetLargeVolumeShareholdersByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetLargeVolumeShareholdersService(mockClient)

	// Mock response
	mockResponse := EdinetLargeVolumeShareholdersResponse{
		Data: []EdinetLargeVolumeShareholderDoc{
			{
				DocId:             "S100WBIV",
				Code:              "86970",
				SubDate:           "2025-07-07",
				LargeHldgTypeCode: "5",
			},
			{
				DocId:             "S100WBJZ",
				Code:              "13660",
				SubDate:           "2025-07-07",
				LargeHldgTypeCode: "1",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/edinet/large-volume-shareholders?date=20250707", mockResponse)

	// Execute
	data, err := service.GetLargeVolumeShareholdersByDate(context.Background(), "20250707")

	// Verify
	if err != nil {
		t.Fatalf("GetLargeVolumeShareholdersByDate() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetLargeVolumeShareholdersByDate() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.SubDate != "2025-07-07" {
			t.Errorf("GetLargeVolumeShareholdersByDate() returned sub date %v, want 2025-07-07", item.SubDate)
		}
	}
}

func TestEdinetLargeVolumeShareholderDoc_IsChangeReport(t *testing.T) {
	tests := []struct {
		name string
		code string
		want bool
	}{
		{name: "large volume report", code: LargeHldgTypeCodeReport, want: false},
		{name: "change report", code: LargeHldgTypeCodeChangeReport, want: true},
		{name: "change report with short term transfer", code: LargeHldgTypeCodeChangeReportShortTermTransfer, want: true},
		{name: "special report", code: LargeHldgTypeCodeSpecialReport, want: false},
		{name: "special change report", code: LargeHldgTypeCodeSpecialChangeReport, want: true},
		{name: "unknown", code: LargeHldgTypeCodeUnknown, want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := EdinetLargeVolumeShareholderDoc{LargeHldgTypeCode: tt.code}
			if got := doc.IsChangeReport(); got != tt.want {
				t.Errorf("EdinetLargeVolumeShareholderDoc.IsChangeReport() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEdinetLargeVolumeShareholderDoc_GetTotalShareholdingPercentage(t *testing.T) {
	doc := EdinetLargeVolumeShareholderDoc{
		TotalShsRatio: floatPtr(0.1343),
	}

	expected := 13.43
	got := doc.GetTotalShareholdingPercentage()

	if got != expected {
		t.Errorf("EdinetLargeVolumeShareholderDoc.GetTotalShareholdingPercentage() = %v, want %v", got, expected)
	}

	// 合計欄の記載がない書類（null）の場合は0
	empty := EdinetLargeVolumeShareholderDoc{}
	if got := empty.GetTotalShareholdingPercentage(); got != 0 {
		t.Errorf("GetTotalShareholdingPercentage() with nil TotalShsRatio = %v, want 0", got)
	}
}

func TestEdinetLargeVolumeShareholderHolder_GetShareholdingPercentage(t *testing.T) {
	holder := EdinetLargeVolumeShareholderHolder{
		ShsRatio: 0.0572,
	}

	expected := 5.72
	got := holder.GetShareholdingPercentage()

	if got != expected {
		t.Errorf("EdinetLargeVolumeShareholderHolder.GetShareholdingPercentage() = %v, want %v", got, expected)
	}
}

func TestEdinetLargeVolumeShareholdersService_GetLargeVolumeShareholders_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewEdinetLargeVolumeShareholdersService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/edinet/large-volume-shareholders?code=86970", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetLargeVolumeShareholders(context.Background(), EdinetLargeVolumeShareholdersParams{Code: "86970"})

	// Verify
	if err == nil {
		t.Error("GetLargeVolumeShareholders() expected error but got nil")
	}
}
