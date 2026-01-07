package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestDailyMarginInterestService_GetDailyMarginInterest(t *testing.T) {
	tests := []struct {
		name     string
		params   DailyMarginInterestParams
		wantPath string
	}{
		{
			name: "with code and date range",
			params: DailyMarginInterestParams{
				Code: "13260",
				From: "20240101",
				To:   "20240331",
			},
			wantPath: "/markets/margin-alert?code=13260&from=20240101&to=20240331",
		},
		{
			name: "with code only",
			params: DailyMarginInterestParams{
				Code: "13260",
			},
			wantPath: "/markets/margin-alert?code=13260",
		},
		{
			name: "with date only",
			params: DailyMarginInterestParams{
				Date: "20240208",
			},
			wantPath: "/markets/margin-alert?date=20240208",
		},
		{
			name: "with code and date",
			params: DailyMarginInterestParams{
				Code: "13260",
				Date: "20240208",
			},
			wantPath: "/markets/margin-alert?code=13260&date=20240208",
		},
		{
			name: "with pagination key",
			params: DailyMarginInterestParams{
				Date:          "20240208",
				PaginationKey: "key123",
			},
			wantPath: "/markets/margin-alert?date=20240208&pagination_key=key123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewDailyMarginInterestService(mockClient)

			// Mock response based on documentation sample
			mockResponse := DailyMarginInterestResponse{
				Data: []DailyMarginInterest{
					{
						PubDate: "2024-02-08",
						Code:    "13260",
						AppDate: "2024-02-07",
						PubReason: PublishReason{
							Restricted:          "0",
							DailyPublication:    "0",
							Monitoring:          "0",
							RestrictedByJSF:     "0",
							PrecautionByJSF:     "1",
							UnclearOrSecOnAlert: "0",
						},
						ShrtOut:       11.0,
						ShrtOutChg:    0.0,
						ShrtOutRatio:  "*",
						LongOut:       676.0,
						LongOutChg:    -20.0,
						LongOutRatio:  "*",
						SLRatio:       1.6,
						ShrtNegOut:    0.0,
						ShrtNegOutChg: 0.0,
						ShrtStdOut:    11.0,
						ShrtStdOutChg: 0.0,
						LongNegOut:    192.0,
						LongNegOutChg: -20.0,
						LongStdOut:    484.0,
						LongStdOutChg: 0.0,
						TSEMrgnRegCls: "001",
					},
				},
				PaginationKey: "",
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Execute
			resp, err := service.GetDailyMarginInterest(tt.params)

			// Verify
			if err != nil {
				t.Fatalf("GetDailyMarginInterest() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetDailyMarginInterest() returned nil response")
				return
			}
			if len(resp.Data) == 0 {
				t.Error("GetDailyMarginInterest() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetDailyMarginInterest() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestDailyMarginInterestService_GetDailyMarginInterest_RequiresCodeOrDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewDailyMarginInterestService(mockClient)

	// Execute with empty code and date
	_, err := service.GetDailyMarginInterest(DailyMarginInterestParams{})

	// Verify
	if err == nil {
		t.Error("GetDailyMarginInterest() expected error for missing code and date but got nil")
	}
	if err.Error() != "either code or date parameter is required" {
		t.Errorf("GetDailyMarginInterest() error = %v, want 'either code or date parameter is required'", err)
	}
}

func TestDailyMarginInterestService_GetDailyMarginInterestByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewDailyMarginInterestService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := DailyMarginInterestResponse{
		Data: []DailyMarginInterest{
			{
				PubDate: "2024-02-08",
				Code:    "13260",
				AppDate: "2024-02-07",
				PubReason: PublishReason{
					PrecautionByJSF: "1",
				},
				ShrtOut: 11.0,
				LongOut: 676.0,
			},
			{
				PubDate: "2024-02-07",
				Code:    "13260",
				AppDate: "2024-02-06",
				PubReason: PublishReason{
					PrecautionByJSF: "1",
				},
				ShrtOut: 11.0,
				LongOut: 696.0,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := DailyMarginInterestResponse{
		Data: []DailyMarginInterest{
			{
				PubDate: "2024-02-06",
				Code:    "13260",
				AppDate: "2024-02-05",
				PubReason: PublishReason{
					PrecautionByJSF: "1",
				},
				ShrtOut: 10.0,
				LongOut: 700.0,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/markets/margin-alert?code=13260", mockResponse1)
	mockClient.SetResponse("GET", "/markets/margin-alert?code=13260&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetDailyMarginInterestByCode("13260")

	// Verify
	if err != nil {
		t.Fatalf("GetDailyMarginInterestByCode() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetDailyMarginInterestByCode() returned %d items, want 3", len(data))
	}
	for _, item := range data {
		if item.Code != "13260" {
			t.Errorf("GetDailyMarginInterestByCode() returned code %v, want 13260", item.Code)
		}
	}
}

func TestDailyMarginInterestService_GetDailyMarginInterestByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewDailyMarginInterestService(mockClient)

	// Mock response
	mockResponse := DailyMarginInterestResponse{
		Data: []DailyMarginInterest{
			{
				PubDate: "2024-02-08",
				Code:    "13260",
				AppDate: "2024-02-07",
				ShrtOut: 11.0,
				LongOut: 676.0,
			},
			{
				PubDate: "2024-02-08",
				Code:    "72030",
				AppDate: "2024-02-07",
				ShrtOut: 1000.0,
				LongOut: 5000.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/markets/margin-alert?date=20240208", mockResponse)

	// Execute
	data, err := service.GetDailyMarginInterestByDate("20240208")

	// Verify
	if err != nil {
		t.Fatalf("GetDailyMarginInterestByDate() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetDailyMarginInterestByDate() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.PubDate != "2024-02-08" {
			t.Errorf("GetDailyMarginInterestByDate() returned date %v, want 2024-02-08", item.PubDate)
		}
	}
}

func TestDailyMarginInterestService_GetDailyMarginInterestByCodeAndDateRange(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewDailyMarginInterestService(mockClient)

	// Mock response
	mockResponse := DailyMarginInterestResponse{
		Data: []DailyMarginInterest{
			{
				PubDate: "2024-02-08",
				Code:    "13260",
				AppDate: "2024-02-07",
				ShrtOut: 11.0,
				LongOut: 676.0,
			},
			{
				PubDate: "2024-02-07",
				Code:    "13260",
				AppDate: "2024-02-06",
				ShrtOut: 11.0,
				LongOut: 696.0,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/markets/margin-alert?code=13260&from=20240201&to=20240208", mockResponse)

	// Execute
	data, err := service.GetDailyMarginInterestByCodeAndDateRange("13260", "20240201", "20240208")

	// Verify
	if err != nil {
		t.Fatalf("GetDailyMarginInterestByCodeAndDateRange() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetDailyMarginInterestByCodeAndDateRange() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.Code != "13260" {
			t.Errorf("GetDailyMarginInterestByCodeAndDateRange() returned code %v, want 13260", item.Code)
		}
	}
}

func TestDailyMarginInterest_GetOutChgValue(t *testing.T) {
	tests := []struct {
		name           string
		data           DailyMarginInterest
		wantShortChg   float64
		wantLongChg    float64
	}{
		{
			name: "numeric values",
			data: DailyMarginInterest{
				ShrtOutChg: 100.0,
				LongOutChg: -50.0,
			},
			wantShortChg: 100.0,
			wantLongChg:  -50.0,
		},
		{
			name: "string values (not published)",
			data: DailyMarginInterest{
				ShrtOutChg: "-",
				LongOutChg: "-",
			},
			wantShortChg: 0.0,
			wantLongChg:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotShort := tt.data.GetShortOutChgValue()
			if gotShort != tt.wantShortChg {
				t.Errorf("GetShortOutChgValue() = %v, want %v", gotShort, tt.wantShortChg)
			}

			gotLong := tt.data.GetLongOutChgValue()
			if gotLong != tt.wantLongChg {
				t.Errorf("GetLongOutChgValue() = %v, want %v", gotLong, tt.wantLongChg)
			}
		})
	}
}

func TestDailyMarginInterest_GetOutRatioValue(t *testing.T) {
	tests := []struct {
		name            string
		data            DailyMarginInterest
		wantShortRatio  float64
		wantLongRatio   float64
	}{
		{
			name: "numeric values",
			data: DailyMarginInterest{
				ShrtOutRatio: 0.5,
				LongOutRatio: 1.2,
			},
			wantShortRatio: 0.5,
			wantLongRatio:  1.2,
		},
		{
			name: "string values (ETF)",
			data: DailyMarginInterest{
				ShrtOutRatio: "*",
				LongOutRatio: "*",
			},
			wantShortRatio: 0.0,
			wantLongRatio:  0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotShort := tt.data.GetShortOutRatioValue()
			if gotShort != tt.wantShortRatio {
				t.Errorf("GetShortOutRatioValue() = %v, want %v", gotShort, tt.wantShortRatio)
			}

			gotLong := tt.data.GetLongOutRatioValue()
			if gotLong != tt.wantLongRatio {
				t.Errorf("GetLongOutRatioValue() = %v, want %v", gotLong, tt.wantLongRatio)
			}
		})
	}
}

func TestPublishReason_HelperMethods(t *testing.T) {
	tests := []struct {
		name       string
		reason     PublishReason
		isRestricted       bool
		isDailyPublication bool
		isMonitoring       bool
		isRestrictedByJSF  bool
		isPrecautionByJSF  bool
		isUnclearOrSecOnAlert bool
	}{
		{
			name: "precaution by JSF only",
			reason: PublishReason{
				Restricted:          "0",
				DailyPublication:    "0",
				Monitoring:          "0",
				RestrictedByJSF:     "0",
				PrecautionByJSF:     "1",
				UnclearOrSecOnAlert: "0",
			},
			isRestricted:       false,
			isDailyPublication: false,
			isMonitoring:       false,
			isRestrictedByJSF:  false,
			isPrecautionByJSF:  true,
			isUnclearOrSecOnAlert: false,
		},
		{
			name: "restricted and daily publication",
			reason: PublishReason{
				Restricted:          "1",
				DailyPublication:    "1",
				Monitoring:          "0",
				RestrictedByJSF:     "0",
				PrecautionByJSF:     "0",
				UnclearOrSecOnAlert: "0",
			},
			isRestricted:       true,
			isDailyPublication: true,
			isMonitoring:       false,
			isRestrictedByJSF:  false,
			isPrecautionByJSF:  false,
			isUnclearOrSecOnAlert: false,
		},
		{
			name: "all flags set",
			reason: PublishReason{
				Restricted:          "1",
				DailyPublication:    "1",
				Monitoring:          "1",
				RestrictedByJSF:     "1",
				PrecautionByJSF:     "1",
				UnclearOrSecOnAlert: "1",
			},
			isRestricted:       true,
			isDailyPublication: true,
			isMonitoring:       true,
			isRestrictedByJSF:  true,
			isPrecautionByJSF:  true,
			isUnclearOrSecOnAlert: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.reason.IsRestricted() != tt.isRestricted {
				t.Errorf("IsRestricted() = %v, want %v", tt.reason.IsRestricted(), tt.isRestricted)
			}
			if tt.reason.IsDailyPublication() != tt.isDailyPublication {
				t.Errorf("IsDailyPublication() = %v, want %v", tt.reason.IsDailyPublication(), tt.isDailyPublication)
			}
			if tt.reason.IsMonitoring() != tt.isMonitoring {
				t.Errorf("IsMonitoring() = %v, want %v", tt.reason.IsMonitoring(), tt.isMonitoring)
			}
			if tt.reason.IsRestrictedByJSF() != tt.isRestrictedByJSF {
				t.Errorf("IsRestrictedByJSF() = %v, want %v", tt.reason.IsRestrictedByJSF(), tt.isRestrictedByJSF)
			}
			if tt.reason.IsPrecautionByJSF() != tt.isPrecautionByJSF {
				t.Errorf("IsPrecautionByJSF() = %v, want %v", tt.reason.IsPrecautionByJSF(), tt.isPrecautionByJSF)
			}
			if tt.reason.IsUnclearOrSecOnAlert() != tt.isUnclearOrSecOnAlert {
				t.Errorf("IsUnclearOrSecOnAlert() = %v, want %v", tt.reason.IsUnclearOrSecOnAlert(), tt.isUnclearOrSecOnAlert)
			}
		})
	}
}

func TestDailyMarginInterestService_GetDailyMarginInterest_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewDailyMarginInterestService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/markets/margin-alert?code=13260", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetDailyMarginInterest(DailyMarginInterestParams{Code: "13260"})

	// Verify
	if err == nil {
		t.Error("GetDailyMarginInterest() expected error but got nil")
	}
}

func TestTSEMarginRegulationConstants(t *testing.T) {
	// 東証信用貸借規制区分定数の値を確認
	tests := []struct {
		constant string
		expected string
	}{
		{TSEMarginRegulationNone, "000"},
		{TSEMarginRegulationCautionForNew, "001"},
		{TSEMarginRegulationCautionForSelling, "002"},
		{TSEMarginRegulationCautionForBuying, "003"},
		{TSEMarginRegulationRestrictedNew, "011"},
		{TSEMarginRegulationRestrictedSelling, "012"},
		{TSEMarginRegulationRestrictedBuying, "013"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("TSE margin regulation constant = %v, want %v", tt.constant, tt.expected)
			}
		})
	}
}
