package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestListedService_GetListedInfo(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		date     string
		wantPath string
	}{
		{
			name:     "with code and date",
			code:     "7203",
			date:     "20240101",
			wantPath: "/equities/master?code=7203&date=20240101",
		},
		{
			name:     "with code only",
			code:     "7203",
			date:     "",
			wantPath: "/equities/master?code=7203",
		},
		{
			name:     "with date only",
			code:     "",
			date:     "20240101",
			wantPath: "/equities/master?date=20240101",
		},
		{
			name:     "with no parameters",
			code:     "",
			date:     "",
			wantPath: "/equities/master",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewListedService(mockClient)

			// Mock response
			mockResponse := ListedInfoResponse{
				Data: []ListedInfo{
					{
						Date:     "20240101",
						Code:     "7203",
						CoName:   "トヨタ自動車",
						CoNameEn: "TOYOTA MOTOR CORPORATION",
						S17:      "6",
						S17Nm:    "自動車・輸送機",
						S33:      "3700",
						S33Nm:    "輸送用機器",
						ScaleCat: "TOPIX Core30",
						Mkt:      "0111",
						MktNm:    "プライム",
					},
				},
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Test
			infos, err := service.GetListedInfo(tt.code, tt.date)
			if err != nil {
				t.Errorf("GetListedInfo failed: %v", err)
			}

			// Verify
			if len(infos) != 1 {
				t.Errorf("Expected 1 info, got %d", len(infos))
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

func TestListedService_GetCompanyInfo(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewListedService(mockClient)

	// Mock response
	mockResponse := ListedInfoResponse{
		Data: []ListedInfo{
			{
				Date:     "20240101",
				Code:     "7203",
				CoName:   "トヨタ自動車",
				CoNameEn: "TOYOTA MOTOR CORPORATION",
				S17:      "6",
				S17Nm:    "自動車・輸送機",
				S33:      "3700",
				S33Nm:    "輸送用機器",
				ScaleCat: "TOPIX Core30",
				Mkt:      "0111",
				MktNm:    "プライム",
			},
		},
	}
	mockClient.SetResponse("GET", "/equities/master?code=7203", mockResponse)

	// Test
	info, err := service.GetCompanyInfo("7203")
	if err != nil {
		t.Errorf("GetCompanyInfo failed: %v", err)
	}

	// Verify
	if info == nil {
		t.Errorf("Expected info, got nil")
		return
	}

	if info.Code != "7203" {
		t.Errorf("Expected code 7203, got %s", info.Code)
	}

	if info.CoName != "トヨタ自動車" {
		t.Errorf("Expected company name トヨタ自動車, got %s", info.CoName)
	}

	if info.MktNm != "プライム" {
		t.Errorf("Expected market code name プライム, got %s", info.MktNm)
	}
}

func TestListedService_GetCompanyInfo_NotFound(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewListedService(mockClient)

	// Mock empty response
	mockResponse := ListedInfoResponse{
		Data: []ListedInfo{},
	}
	mockClient.SetResponse("GET", "/equities/master?code=9999", mockResponse)

	// Test
	_, err := service.GetCompanyInfo("9999")
	if err == nil {
		t.Errorf("Expected error for non-existent company, got nil")
	}
}

func TestListedService_GetListedInfo_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewListedService(mockClient)

	// Mock error
	mockClient.SetError("GET", "/equities/master?code=7203", fmt.Errorf("API error"))

	// Test
	_, err := service.GetListedInfo("7203", "")
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestListedService_GetListedBySector17(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewListedService(mockClient)

	// Mock response with multiple companies
	mockResponse := ListedInfoResponse{
		Data: []ListedInfo{
			{
				Code:   "7203",
				CoName: "トヨタ自動車",
				S17:    "6",
				S17Nm:  "自動車・輸送機",
				MktNm:  "プライム",
			},
			{
				Code:   "7267",
				CoName: "本田技研工業",
				S17:    "6",
				S17Nm:  "自動車・輸送機",
				MktNm:  "プライム",
			},
			{
				Code:   "9984",
				CoName: "ソフトバンクグループ",
				S17:    "10",
				S17Nm:  "情報通信・サービスその他",
				MktNm:  "プライム",
			},
		},
	}
	mockClient.SetResponse("GET", "/equities/master", mockResponse)

	// Test - 自動車・輸送機セクターの銘柄を取得
	infos, err := service.GetListedBySector17(Sector17Auto, "")
	if err != nil {
		t.Errorf("GetListedBySector17 failed: %v", err)
	}

	// Verify - 2銘柄が返されるはず
	if len(infos) != 2 {
		t.Errorf("Expected 2 companies, got %d", len(infos))
	}

	for _, info := range infos {
		if info.S17 != "6" {
			t.Errorf("Expected S17 6, got %s for %s", info.S17, info.CoName)
		}
	}
}

func TestListedService_GetListedBySector33(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewListedService(mockClient)

	// Mock response
	mockResponse := ListedInfoResponse{
		Data: []ListedInfo{
			{
				Code:   "4755",
				CoName: "楽天グループ",
				S33:    "5250",
				S33Nm:  "情報・通信業",
				MktNm:  "プライム",
			},
			{
				Code:   "9984",
				CoName: "ソフトバンクグループ",
				S33:    "5250",
				S33Nm:  "情報・通信業",
				MktNm:  "プライム",
			},
			{
				Code:   "7203",
				CoName: "トヨタ自動車",
				S33:    "3700",
				S33Nm:  "輸送用機器",
				MktNm:  "プライム",
			},
		},
	}
	mockClient.SetResponse("GET", "/equities/master", mockResponse)

	// Test - 情報・通信業の銘柄を取得
	infos, err := service.GetListedBySector33(Sector33IT, "")
	if err != nil {
		t.Errorf("GetListedBySector33 failed: %v", err)
	}

	// Verify
	if len(infos) != 2 {
		t.Errorf("Expected 2 companies, got %d", len(infos))
	}
}

func TestListedService_GetListedByMarket(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewListedService(mockClient)

	// Mock response
	mockResponse := ListedInfoResponse{
		Data: []ListedInfo{
			{
				Code:   "7203",
				CoName: "トヨタ自動車",
				Mkt:    MarketPrime,
				MktNm:  "プライム",
			},
			{
				Code:   "9984",
				CoName: "ソフトバンクグループ",
				Mkt:    MarketPrime,
				MktNm:  "プライム",
			},
			{
				Code:   "4755",
				CoName: "楽天グループ",
				Mkt:    MarketPrime,
				MktNm:  "プライム",
			},
			{
				Code:   "3994",
				CoName: "マネーフォワード",
				Mkt:    MarketGrowth,
				MktNm:  "グロース",
			},
		},
	}
	mockClient.SetResponse("GET", "/equities/master", mockResponse)

	// Test - プライム市場の銘柄を取得
	infos, err := service.GetListedByMarket(MarketPrime, "")
	if err != nil {
		t.Errorf("GetListedByMarket failed: %v", err)
	}

	// Verify
	if len(infos) != 3 {
		t.Errorf("Expected 3 companies, got %d", len(infos))
	}

	for _, info := range infos {
		if info.Mkt != MarketPrime {
			t.Errorf("Expected market code %s, got %s", MarketPrime, info.Mkt)
		}
	}
}
