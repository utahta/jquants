package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestTimelyDisclosureService_GetDisclosures(t *testing.T) {
	tests := []struct {
		name     string
		params   TimelyDisclosureParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with date only",
			params: TimelyDisclosureParams{
				Date: "20250401",
			},
			wantPath: "/td/list?date=20250401",
		},
		{
			name: "with hyphenated date",
			params: TimelyDisclosureParams{
				Date: "2025-04-01",
			},
			wantPath: "/td/list?date=2025-04-01",
		},
		{
			name: "with code only",
			params: TimelyDisclosureParams{
				Code: "86970",
			},
			wantPath: "/td/list?code=86970",
		},
		{
			name: "with code and date range",
			params: TimelyDisclosureParams{
				Code: "86970",
				From: "20250301",
				To:   "20250401",
			},
			wantPath: "/td/list?code=86970&from=20250301&to=20250401",
		},
		{
			name: "with date and discItems",
			params: TimelyDisclosureParams{
				Date:      "20250401",
				DiscItems: "11101,11102",
			},
			wantPath: "/td/list?date=20250401&discItems=11101,11102",
		},
		{
			name: "with date and cursor",
			params: TimelyDisclosureParams{
				Date:   "20250401",
				Cursor: "cursor123",
			},
			wantPath: "/td/list?date=20250401&cursor=cursor123",
		},
		{
			name: "with date and pagination key",
			params: TimelyDisclosureParams{
				Date:          "20250401",
				PaginationKey: "key123",
			},
			wantPath: "/td/list?date=20250401&pagination_key=key123",
		},
		{
			name:     "with no required parameters",
			params:   TimelyDisclosureParams{},
			wantPath: "",
			wantErr:  true,
		},
		{
			name: "with cursor and pagination key",
			params: TimelyDisclosureParams{
				Date:          "20250401",
				Cursor:        "cursor123",
				PaginationKey: "key123",
			},
			wantPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewTimelyDisclosureService(mockClient)

			// Mock response based on documentation sample
			mockResponse := TimelyDisclosureResponse{
				Data: []TimelyDisclosure{
					{
						DiscNo:     "20250401130100",
						Code:       "86970",
						Name:       "日本取引所グループ",
						DiscDate:   "2025-04-01",
						DiscTime:   "08:00",
						Title:      "2025年3月期 決算短信〔日本基準〕（連結）",
						DiscStatus: "",
						RevNo:      1,
						DiscItems:  []string{"11101"},
						Docs:       []string{"g", "s", "x"},
					},
				},
				Cursor:        "",
				PaginationKey: "",
			}
			if !tt.wantErr {
				mockClient.SetResponse("GET", tt.wantPath, mockResponse)
			}

			// Execute
			resp, err := service.GetDisclosures(context.Background(), tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetDisclosures() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetDisclosures() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetDisclosures() returned nil response")
				return
			}
			if len(resp.Data) == 0 {
				t.Error("GetDisclosures() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetDisclosures() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestTimelyDisclosureService_GetDisclosures_RequiresParameter(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTimelyDisclosureService(mockClient)

	// Execute with empty parameters
	_, err := service.GetDisclosures(context.Background(), TimelyDisclosureParams{})

	// Verify
	if err == nil {
		t.Error("GetDisclosures() expected error for missing parameters but got nil")
	}
	if err.Error() != "either date or code parameter is required" {
		t.Errorf("GetDisclosures() error = %v, want 'either date or code parameter is required'", err)
	}
}

func TestTimelyDisclosureService_GetDisclosures_CursorConflictsWithPaginationKey(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTimelyDisclosureService(mockClient)

	// Execute with cursor and pagination_key
	_, err := service.GetDisclosures(context.Background(), TimelyDisclosureParams{
		Date:          "20250401",
		Cursor:        "cursor123",
		PaginationKey: "key123",
	})

	// Verify
	if err == nil {
		t.Error("GetDisclosures() expected error for cursor and pagination_key but got nil")
	}
	if err.Error() != "cursor and pagination_key parameters cannot be specified together" {
		t.Errorf("GetDisclosures() error = %v, want 'cursor and pagination_key parameters cannot be specified together'", err)
	}
}

func TestTimelyDisclosureResponse_UnmarshalJSON(t *testing.T) {
	// DiscStatusのnullが空文字として、配列フィールドとcursorがそのまま取得できること
	jsonData := `{
		"data": [
			{
				"DiscNo": "20250401130100",
				"Code": "86970",
				"Name": "日本取引所グループ",
				"DiscDate": "2025-04-01",
				"DiscTime": "08:00",
				"Title": "2025年3月期 決算短信〔日本基準〕（連結）",
				"DiscStatus": null,
				"RevNo": 1,
				"DiscItems": ["11101", "11102"],
				"Docs": ["g", "s", "x"]
			},
			{
				"DiscNo": "20250401130200",
				"Code": "13010",
				"Name": "極洋",
				"DiscDate": "2025-04-01",
				"DiscTime": "09:00",
				"Title": "訂正開示",
				"DiscStatus": "revision",
				"RevNo": 2,
				"DiscItems": ["11102"],
				"Docs": ["g"]
			}
		],
		"cursor": "cursor123",
		"pagination_key": ""
	}`

	var resp TimelyDisclosureResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Fatalf("data count = %d, want 2", len(resp.Data))
	}
	if resp.Cursor != "cursor123" {
		t.Errorf("Cursor = %v, want cursor123", resp.Cursor)
	}
	if resp.PaginationKey != "" {
		t.Errorf("PaginationKey = %v, want empty", resp.PaginationKey)
	}

	// 1件目: DiscStatusはnull→空文字、配列フィールドあり
	d := resp.Data[0]
	if d.DiscNo != "20250401130100" {
		t.Errorf("DiscNo = %v, want 20250401130100", d.DiscNo)
	}
	if d.DiscStatus != "" {
		t.Errorf("DiscStatus = %v, want empty for null", d.DiscStatus)
	}
	if d.RevNo != 1 {
		t.Errorf("RevNo = %v, want 1", d.RevNo)
	}
	if len(d.DiscItems) != 2 || d.DiscItems[0] != "11101" || d.DiscItems[1] != "11102" {
		t.Errorf("DiscItems = %v, want [11101 11102]", d.DiscItems)
	}
	if len(d.Docs) != 3 || d.Docs[0] != "g" || d.Docs[1] != "s" || d.Docs[2] != "x" {
		t.Errorf("Docs = %v, want [g s x]", d.Docs)
	}
	if d.IsRevision() || d.IsDeleted() {
		t.Error("IsRevision()/IsDeleted() should be false for null DiscStatus")
	}
	if !d.HasXBRL() {
		t.Error("HasXBRL() should be true when Docs contains x")
	}

	// 2件目: 訂正開示、XBRLなし
	d = resp.Data[1]
	if !d.IsRevision() {
		t.Error("IsRevision() should be true for revision")
	}
	if d.RevNo != 2 {
		t.Errorf("RevNo = %v, want 2", d.RevNo)
	}
	if d.HasXBRL() {
		t.Error("HasXBRL() should be false when Docs does not contain x")
	}
}

func TestTimelyDisclosureService_GetDisclosuresByDate(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTimelyDisclosureService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := TimelyDisclosureResponse{
		Data: []TimelyDisclosure{
			{
				DiscNo:   "20250401130100",
				Code:     "86970",
				DiscDate: "2025-04-01",
				RevNo:    1,
			},
			{
				DiscNo:   "20250401130200",
				Code:     "13010",
				DiscDate: "2025-04-01",
				RevNo:    1,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := TimelyDisclosureResponse{
		Data: []TimelyDisclosure{
			{
				DiscNo:   "20250401130300",
				Code:     "72030",
				DiscDate: "2025-04-01",
				RevNo:    1,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/td/list?date=20250401", mockResponse1)
	mockClient.SetResponse("GET", "/td/list?date=20250401&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetDisclosuresByDate(context.Background(), "20250401")

	// Verify
	if err != nil {
		t.Fatalf("GetDisclosuresByDate() error = %v", err)
	}
	if len(data) != 3 {
		t.Errorf("GetDisclosuresByDate() returned %d items, want 3", len(data))
	}
	for _, item := range data {
		if item.DiscDate != "2025-04-01" {
			t.Errorf("GetDisclosuresByDate() returned date %v, want 2025-04-01", item.DiscDate)
		}
	}
}

func TestTimelyDisclosureService_GetDisclosuresByCode(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTimelyDisclosureService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := TimelyDisclosureResponse{
		Data: []TimelyDisclosure{
			{
				DiscNo:   "20250401130100",
				Code:     "86970",
				DiscDate: "2025-04-01",
				RevNo:    1,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := TimelyDisclosureResponse{
		Data: []TimelyDisclosure{
			{
				DiscNo:   "20250301120100",
				Code:     "86970",
				DiscDate: "2025-03-01",
				RevNo:    1,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/td/list?code=86970", mockResponse1)
	mockClient.SetResponse("GET", "/td/list?code=86970&pagination_key=next_page_key", mockResponse2)

	// Execute
	data, err := service.GetDisclosuresByCode(context.Background(), "86970")

	// Verify
	if err != nil {
		t.Fatalf("GetDisclosuresByCode() error = %v", err)
	}
	if len(data) != 2 {
		t.Errorf("GetDisclosuresByCode() returned %d items, want 2", len(data))
	}
	for _, item := range data {
		if item.Code != "86970" {
			t.Errorf("GetDisclosuresByCode() returned code %v, want 86970", item.Code)
		}
	}
}

func TestTimelyDisclosureService_GetDisclosureFiles(t *testing.T) {
	tests := []struct {
		name     string
		params   TimelyDisclosureFilesParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with discNo only",
			params: TimelyDisclosureFilesParams{
				DiscNo: "20250401130100",
			},
			wantPath: "/td/files?discNo=20250401130100",
		},
		{
			name: "with discNo and docs",
			params: TimelyDisclosureFilesParams{
				DiscNo: "20250401130100",
				Docs:   "g,s",
			},
			wantPath: "/td/files?discNo=20250401130100&docs=g,s",
		},
		{
			name: "with docs constant",
			params: TimelyDisclosureFilesParams{
				DiscNo: "20250401130100",
				Docs:   TimelyDisclosureDocXBRL,
			},
			wantPath: "/td/files?discNo=20250401130100&docs=x",
		},
		{
			name:     "without discNo",
			params:   TimelyDisclosureFilesParams{},
			wantPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewTimelyDisclosureService(mockClient)

			// Mock response based on documentation sample
			mockResponse := TimelyDisclosureFiles{
				DiscNo: "20250401130100",
				Files: TimelyDisclosureFileURLs{
					PDF:        "https://example.com/download-url-pdf",
					SummaryPDF: "https://example.com/download-url-summary",
					XBRL:       "https://example.com/download-url-xbrl",
				},
			}
			if !tt.wantErr {
				mockClient.SetResponse("GET", tt.wantPath, mockResponse)
			}

			// Execute
			resp, err := service.GetDisclosureFiles(context.Background(), tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetDisclosureFiles() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetDisclosureFiles() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetDisclosureFiles() returned nil response")
				return
			}
			if resp.DiscNo != "20250401130100" {
				t.Errorf("GetDisclosureFiles() DiscNo = %v, want 20250401130100", resp.DiscNo)
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetDisclosureFiles() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestTimelyDisclosureService_GetDisclosureFiles_RequiresDiscNo(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTimelyDisclosureService(mockClient)

	// Execute with empty parameters
	_, err := service.GetDisclosureFiles(context.Background(), TimelyDisclosureFilesParams{})

	// Verify
	if err == nil {
		t.Error("GetDisclosureFiles() expected error for missing discNo but got nil")
	}
	if err.Error() != "discNo parameter is required" {
		t.Errorf("GetDisclosureFiles() error = %v, want 'discNo parameter is required'", err)
	}
}

func TestTimelyDisclosureFiles_UnmarshalJSON_AllFiles(t *testing.T) {
	// 全種別のURLが取得できること
	jsonData := `{
		"discNo": "20250401130100",
		"files": {
			"pdf": "https://example.com/download-url-pdf",
			"summaryPdf": "https://example.com/download-url-summary",
			"xbrl": "https://example.com/download-url-xbrl"
		}
	}`

	var files TimelyDisclosureFiles
	if err := json.Unmarshal([]byte(jsonData), &files); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	if files.DiscNo != "20250401130100" {
		t.Errorf("DiscNo = %v, want 20250401130100", files.DiscNo)
	}
	if files.Files.PDF != "https://example.com/download-url-pdf" {
		t.Errorf("Files.PDF = %v, want https://example.com/download-url-pdf", files.Files.PDF)
	}
	if files.Files.SummaryPDF != "https://example.com/download-url-summary" {
		t.Errorf("Files.SummaryPDF = %v, want https://example.com/download-url-summary", files.Files.SummaryPDF)
	}
	if files.Files.XBRL != "https://example.com/download-url-xbrl" {
		t.Errorf("Files.XBRL = %v, want https://example.com/download-url-xbrl", files.Files.XBRL)
	}
}

func TestTimelyDisclosureFiles_UnmarshalJSON_PDFOnly(t *testing.T) {
	// docsで絞った場合など、存在しないキーは空文字になること
	jsonData := `{
		"discNo": "20250401130100",
		"files": {
			"pdf": "https://example.com/download-url-pdf"
		}
	}`

	var files TimelyDisclosureFiles
	if err := json.Unmarshal([]byte(jsonData), &files); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	if files.Files.PDF != "https://example.com/download-url-pdf" {
		t.Errorf("Files.PDF = %v, want https://example.com/download-url-pdf", files.Files.PDF)
	}
	if files.Files.SummaryPDF != "" {
		t.Errorf("Files.SummaryPDF = %v, want empty", files.Files.SummaryPDF)
	}
	if files.Files.XBRL != "" {
		t.Errorf("Files.XBRL = %v, want empty", files.Files.XBRL)
	}
}

func TestTimelyDisclosureService_GetBulkFile(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTimelyDisclosureService(mockClient)

	// Mock response based on documentation sample
	mockResponse := TimelyDisclosureBulk{
		LastUpdated: "2025-04-01T08:00:00Z",
		URL:         "https://example.com/download-url-bulk-csv",
	}
	mockClient.SetResponse("GET", "/td/bulk", mockResponse)

	// Execute
	resp, err := service.GetBulkFile(context.Background())

	// Verify
	if err != nil {
		t.Fatalf("GetBulkFile() error = %v", err)
	}
	if resp == nil {
		t.Fatal("GetBulkFile() returned nil response")
		return
	}
	if resp.LastUpdated != "2025-04-01T08:00:00Z" {
		t.Errorf("GetBulkFile() LastUpdated = %v, want 2025-04-01T08:00:00Z", resp.LastUpdated)
	}
	if resp.URL != "https://example.com/download-url-bulk-csv" {
		t.Errorf("GetBulkFile() URL = %v, want https://example.com/download-url-bulk-csv", resp.URL)
	}
	if mockClient.LastPath != "/td/bulk" {
		t.Errorf("GetBulkFile() path = %v, want /td/bulk", mockClient.LastPath)
	}
}

func TestTimelyDisclosureBulk_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"lastUpdated": "2025-04-01T08:00:00Z",
		"url": "https://example.com/download-url-bulk-csv"
	}`

	var bulk TimelyDisclosureBulk
	if err := json.Unmarshal([]byte(jsonData), &bulk); err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}
	if bulk.LastUpdated != "2025-04-01T08:00:00Z" {
		t.Errorf("LastUpdated = %v, want 2025-04-01T08:00:00Z", bulk.LastUpdated)
	}
	if bulk.URL != "https://example.com/download-url-bulk-csv" {
		t.Errorf("URL = %v, want https://example.com/download-url-bulk-csv", bulk.URL)
	}
}

func TestTimelyDisclosure_IsRevision(t *testing.T) {
	tests := []struct {
		name       string
		disclosure TimelyDisclosure
		want       bool
	}{
		{
			name:       "revision",
			disclosure: TimelyDisclosure{DiscStatus: "revision"},
			want:       true,
		},
		{
			name:       "new disclosure",
			disclosure: TimelyDisclosure{DiscStatus: ""},
			want:       false,
		},
		{
			name:       "delete",
			disclosure: TimelyDisclosure{DiscStatus: "delete"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.disclosure.IsRevision(); got != tt.want {
				t.Errorf("TimelyDisclosure.IsRevision() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimelyDisclosure_IsDeleted(t *testing.T) {
	tests := []struct {
		name       string
		disclosure TimelyDisclosure
		want       bool
	}{
		{
			name:       "delete",
			disclosure: TimelyDisclosure{DiscStatus: "delete"},
			want:       true,
		},
		{
			name:       "new disclosure",
			disclosure: TimelyDisclosure{DiscStatus: ""},
			want:       false,
		},
		{
			name:       "revision",
			disclosure: TimelyDisclosure{DiscStatus: "revision"},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.disclosure.IsDeleted(); got != tt.want {
				t.Errorf("TimelyDisclosure.IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimelyDisclosure_HasXBRL(t *testing.T) {
	tests := []struct {
		name       string
		disclosure TimelyDisclosure
		want       bool
	}{
		{
			name:       "has xbrl",
			disclosure: TimelyDisclosure{Docs: []string{"g", "s", "x"}},
			want:       true,
		},
		{
			name:       "no xbrl",
			disclosure: TimelyDisclosure{Docs: []string{"g", "s"}},
			want:       false,
		},
		{
			name:       "empty docs",
			disclosure: TimelyDisclosure{Docs: nil},
			want:       false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.disclosure.HasXBRL(); got != tt.want {
				t.Errorf("TimelyDisclosure.HasXBRL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTimelyDisclosureService_GetDisclosures_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTimelyDisclosureService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/td/list?code=86970", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetDisclosures(context.Background(), TimelyDisclosureParams{Code: "86970"})

	// Verify
	if err == nil {
		t.Error("GetDisclosures() expected error but got nil")
	}
}

func TestTimelyDisclosureService_SignedURLMethods_SkipCache(t *testing.T) {
	mockClient := client.NewMockClient()
	service := NewTimelyDisclosureService(mockClient)

	mockClient.SetResponse("GET", "/td/files?discNo=20250401130100", TimelyDisclosureFiles{DiscNo: "20250401130100"})
	if _, err := service.GetDisclosureFiles(context.Background(), TimelyDisclosureFilesParams{DiscNo: "20250401130100"}); err != nil {
		t.Fatalf("GetDisclosureFiles() error = %v", err)
	}
	// 署名付きURLは失効するため、キャッシュバイパスを要求していること
	if !mockClient.LastSkipCache {
		t.Error("GetDisclosureFiles() should request cache bypass for expiring signed URLs")
	}

	mockClient.SetResponse("GET", "/td/bulk", TimelyDisclosureBulk{URL: "https://example.com/bulk.csv.gz"})
	if _, err := service.GetBulkFile(context.Background()); err != nil {
		t.Fatalf("GetBulkFile() error = %v", err)
	}
	if !mockClient.LastSkipCache {
		t.Error("GetBulkFile() should request cache bypass for expiring signed URLs")
	}
}

func TestTimelyDisclosureService_GetDisclosures_CursorSkipsCache(t *testing.T) {
	mockClient := client.NewMockClient()
	service := NewTimelyDisclosureService(mockClient)
	mockClient.SetResponse("GET", "/td/list?date=2025-04-01&cursor=cur123", TimelyDisclosureResponse{})
	mockClient.SetResponse("GET", "/td/list?date=2025-04-01", TimelyDisclosureResponse{})

	// cursorによるポーリングはキャッシュをバイパスすること
	if _, err := service.GetDisclosures(context.Background(), TimelyDisclosureParams{Date: "2025-04-01", Cursor: "cur123"}); err != nil {
		t.Fatalf("GetDisclosures() error = %v", err)
	}
	if !mockClient.LastSkipCache {
		t.Error("GetDisclosures() with cursor should bypass the session cache")
	}

	// cursorなしは通常どおりキャッシュ対象
	if _, err := service.GetDisclosures(context.Background(), TimelyDisclosureParams{Date: "2025-04-01"}); err != nil {
		t.Fatalf("GetDisclosures() error = %v", err)
	}
	if mockClient.LastSkipCache {
		t.Error("GetDisclosures() without cursor should use the session cache")
	}
}
