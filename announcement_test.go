package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestAnnouncementService_GetAnnouncement(t *testing.T) {
	tests := []struct {
		name     string
		params   AnnouncementParams
		wantPath string
	}{
		{
			name:     "without pagination",
			params:   AnnouncementParams{},
			wantPath: "/fins/announcement",
		},
		{
			name: "with pagination",
			params: AnnouncementParams{
				PaginationKey: "key123",
			},
			wantPath: "/fins/announcement?pagination_key=key123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewAnnouncementService(mockClient)

			// Mock response
			mockResponse := AnnouncementResponse{
				Announcement: []Announcement{
					{
						Date:          "2024-02-14",
						Code:          "43760",
						CompanyName:   "くふうカンパニー",
						FiscalYear:    "9月30日",
						SectorName:    "情報・通信業",
						FiscalQuarter: "第１四半期",
						Section:       "マザーズ",
					},
					{
						Date:          "2024-02-14",
						Code:          "7203",
						CompanyName:   "トヨタ自動車",
						FiscalYear:    "3月31日",
						SectorName:    "輸送用機器",
						FiscalQuarter: "第３四半期",
						Section:       "プライム",
					},
				},
				PaginationKey: "",
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Test
			resp, err := service.GetAnnouncement(tt.params)
			if err != nil {
				t.Errorf("GetAnnouncement failed: %v", err)
			}

			// Verify
			if resp == nil {
				t.Fatal("Expected response, got nil")
			}
			if len(resp.Announcement) != 2 {
				t.Errorf("Expected 2 announcements, got %d", len(resp.Announcement))
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

func TestAnnouncementService_GetAllAnnouncements(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewAnnouncementService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := AnnouncementResponse{
		Announcement: []Announcement{
			{
				Date:          "2024-02-14",
				Code:          "7203",
				CompanyName:   "トヨタ自動車",
				FiscalYear:    "3月31日",
				SectorName:    "輸送用機器",
				FiscalQuarter: "第３四半期",
				Section:       "プライム",
			},
			{
				Date:          "2024-02-14",
				Code:          "9984",
				CompanyName:   "ソフトバンクグループ",
				FiscalYear:    "3月31日",
				SectorName:    "情報・通信業",
				FiscalQuarter: "第３四半期",
				Section:       "プライム",
			},
		},
		PaginationKey: "next_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := AnnouncementResponse{
		Announcement: []Announcement{
			{
				Date:          "2024-02-14",
				Code:          "43760",
				CompanyName:   "くふうカンパニー",
				FiscalYear:    "9月30日",
				SectorName:    "情報・通信業",
				FiscalQuarter: "第１四半期",
				Section:       "グロース",
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/fins/announcement", mockResponse1)
	mockClient.SetResponse("GET", "/fins/announcement?pagination_key=next_key", mockResponse2)

	// Test
	announcements, err := service.GetAllAnnouncements()
	if err != nil {
		t.Errorf("GetAllAnnouncements failed: %v", err)
	}

	// Verify
	if len(announcements) != 3 {
		t.Errorf("Expected 3 announcements, got %d", len(announcements))
	}

	// 最初の銘柄を確認
	if announcements[0].Code != "7203" {
		t.Errorf("Expected first code to be 7203, got %s", announcements[0].Code)
	}

	// 最後の銘柄を確認
	if announcements[2].Code != "43760" {
		t.Errorf("Expected last code to be 43760, got %s", announcements[2].Code)
	}
}

func TestAnnouncementService_GetAnnouncementByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewAnnouncementService(mockClient)

	// Mock response
	mockResponse := AnnouncementResponse{
		Announcement: []Announcement{
			{
				Date:          "2024-02-14",
				Code:          "7203",
				CompanyName:   "トヨタ自動車",
				FiscalYear:    "3月31日",
				SectorName:    "輸送用機器",
				FiscalQuarter: "第３四半期",
				Section:       "プライム",
			},
			{
				Date:          "2024-02-14",
				Code:          "9984",
				CompanyName:   "ソフトバンクグループ",
				FiscalYear:    "3月31日",
				SectorName:    "情報・通信業",
				FiscalQuarter: "第３四半期",
				Section:       "プライム",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/fins/announcement", mockResponse)

	// Test
	announcement, err := service.GetAnnouncementByCode("7203")
	if err != nil {
		t.Errorf("GetAnnouncementByCode failed: %v", err)
	}

	// Verify
	if announcement == nil {
		t.Errorf("Expected announcement, got nil")
		return
	}

	if announcement.Code != "7203" {
		t.Errorf("Expected code 7203, got %s", announcement.Code)
	}

	if announcement.Section != "プライム" {
		t.Errorf("Expected section プライム, got %s", announcement.Section)
	}

	if announcement.FiscalQuarter != "第３四半期" {
		t.Errorf("Expected fiscal quarter 第３四半期, got %s", announcement.FiscalQuarter)
	}
}

func TestAnnouncementService_GetAnnouncementByCode_NotFound(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewAnnouncementService(mockClient)

	// Mock response
	mockResponse := AnnouncementResponse{
		Announcement: []Announcement{
			{
				Code:        "9984",
				CompanyName: "ソフトバンクグループ",
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/fins/announcement", mockResponse)

	// Test
	_, err := service.GetAnnouncementByCode("7203")
	if err == nil {
		t.Errorf("Expected error for non-existent code, got nil")
	}
}

func TestAnnouncementService_GetAnnouncement_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewAnnouncementService(mockClient)

	// Mock error
	mockClient.SetError("GET", "/fins/announcement", fmt.Errorf("API error"))

	// Test
	_, err := service.GetAnnouncement(AnnouncementParams{})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
