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
			wantPath: "/listed/info?code=7203&date=20240101",
		},
		{
			name:     "with code only",
			code:     "7203",
			date:     "",
			wantPath: "/listed/info?code=7203",
		},
		{
			name:     "with date only",
			code:     "",
			date:     "20240101",
			wantPath: "/listed/info?date=20240101",
		},
		{
			name:     "with no parameters",
			code:     "",
			date:     "",
			wantPath: "/listed/info",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewListedService(mockClient)

			// Mock response
			mockResponse := ListedInfoResponse{
				Info: []ListedInfo{
					{
						Date:               "20240101",
						Code:               "7203",
						CompanyName:        "トヨタ自動車",
						CompanyNameEnglish: "TOYOTA MOTOR CORPORATION",
						LocalCode:          "72030",
						Sector17Code:       "6",
						Sector17CodeName:   "自動車・輸送機",
						Sector33Code:       "3700",
						Sector33CodeName:   "輸送用機器",
						ScaleCategory:      "TOPIX Core30",
						MarketCode:         "0111",
						MarketCodeName:     "プライム",
						IsDelisted:         "false",
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
		Info: []ListedInfo{
			{
				Date:               "20240101",
				Code:               "7203",
				CompanyName:        "トヨタ自動車",
				CompanyNameEnglish: "TOYOTA MOTOR CORPORATION",
				LocalCode:          "72030",
				Sector17Code:       "6",
				Sector17CodeName:   "自動車・輸送機",
				Sector33Code:       "3700",
				Sector33CodeName:   "輸送用機器",
				ScaleCategory:      "TOPIX Core30",
				MarketCode:         "0111",
				MarketCodeName:     "プライム",
				IsDelisted:         "false",
			},
		},
	}
	mockClient.SetResponse("GET", "/listed/info?code=7203", mockResponse)

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

	if info.CompanyName != "トヨタ自動車" {
		t.Errorf("Expected company name トヨタ自動車, got %s", info.CompanyName)
	}

	if info.MarketCodeName != "プライム" {
		t.Errorf("Expected market code name プライム, got %s", info.MarketCodeName)
	}
}

func TestListedService_GetCompanyInfo_NotFound(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewListedService(mockClient)

	// Mock empty response
	mockResponse := ListedInfoResponse{
		Info: []ListedInfo{},
	}
	mockClient.SetResponse("GET", "/listed/info?code=9999", mockResponse)

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
	mockClient.SetError("GET", "/listed/info?code=7203", fmt.Errorf("API error"))

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
		Info: []ListedInfo{
			{
				Code:             "7203",
				CompanyName:      "トヨタ自動車",
				LocalCode:        "72030",
				Sector17Code:     "6",
				Sector17CodeName: "自動車・輸送機",
				MarketCodeName:   "プライム",
				IsDelisted:       "false",
			},
			{
				Code:             "7267",
				CompanyName:      "本田技研工業",
				LocalCode:        "72670",
				Sector17Code:     "6",
				Sector17CodeName: "自動車・輸送機",
				MarketCodeName:   "プライム",
				IsDelisted:       "false",
			},
			{
				Code:             "9984",
				CompanyName:      "ソフトバンクグループ",
				LocalCode:        "99840",
				Sector17Code:     "10",
				Sector17CodeName: "情報通信・サービスその他",
				MarketCodeName:   "プライム",
				IsDelisted:       "false",
			},
		},
	}
	mockClient.SetResponse("GET", "/listed/info", mockResponse)

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
		if info.Sector17Code != "6" {
			t.Errorf("Expected sector17Code 6, got %s for %s", info.Sector17Code, info.CompanyName)
		}
	}
}

func TestListedService_GetListedBySector33(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewListedService(mockClient)

	// Mock response
	mockResponse := ListedInfoResponse{
		Info: []ListedInfo{
			{
				Code:             "4755",
				CompanyName:      "楽天グループ",
				LocalCode:        "47550",
				Sector33Code:     "5250",
				Sector33CodeName: "情報・通信業",
				MarketCodeName:   "プライム",
			},
			{
				Code:             "9984",
				CompanyName:      "ソフトバンクグループ",
				LocalCode:        "99840",
				Sector33Code:     "5250",
				Sector33CodeName: "情報・通信業",
				MarketCodeName:   "プライム",
			},
			{
				Code:             "7203",
				CompanyName:      "トヨタ自動車",
				LocalCode:        "72030",
				Sector33Code:     "3700",
				Sector33CodeName: "輸送用機器",
				MarketCodeName:   "プライム",
			},
		},
	}
	mockClient.SetResponse("GET", "/listed/info", mockResponse)

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
		Info: []ListedInfo{
			{
				Code:           "7203",
				CompanyName:    "トヨタ自動車",
				MarketCodeName: "プライム",
			},
			{
				Code:           "9984",
				CompanyName:    "ソフトバンクグループ",
				MarketCodeName: "プライム",
			},
			{
				Code:           "4755",
				CompanyName:    "楽天グループ",
				MarketCodeName: "プライム",
			},
			{
				Code:           "3994",
				CompanyName:    "マネーフォワード",
				MarketCodeName: "グロース",
			},
		},
	}
	mockClient.SetResponse("GET", "/listed/info", mockResponse)

	// Test - プライム市場の銘柄を取得
	infos, err := service.GetListedByMarket("プライム", "")
	if err != nil {
		t.Errorf("GetListedByMarket failed: %v", err)
	}

	// Verify
	if len(infos) != 3 {
		t.Errorf("Expected 3 companies, got %d", len(infos))
	}

	for _, info := range infos {
		if info.MarketCodeName != "プライム" {
			t.Errorf("Expected market プライム, got %s", info.MarketCodeName)
		}
	}
}

func TestListedInfo_IsDelistedBool(t *testing.T) {
	tests := []struct {
		name       string
		isDelisted string
		want       bool
	}{
		{"上場中", "false", false},
		{"上場廃止", "true", true},
		{"空文字", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := &ListedInfo{
				IsDelisted: tt.isDelisted,
			}
			if got := info.IsDelistedBool(); got != tt.want {
				t.Errorf("IsDelistedBool() = %v, want %v", got, tt.want)
			}
		})
	}
}
