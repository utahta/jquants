package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestBulkService_GetFiles(t *testing.T) {
	tests := []struct {
		name     string
		params   BulkListParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with endpoint only",
			params: BulkListParams{
				Endpoint: "/equities/bars/daily",
			},
			wantPath: "/bulk/list?endpoint=/equities/bars/daily",
		},
		{
			name: "with date only",
			params: BulkListParams{
				Date: "2025-01",
			},
			wantPath: "/bulk/list?date=2025-01",
		},
		{
			name: "with endpoint and period",
			params: BulkListParams{
				Endpoint: "/equities/bars/daily",
				From:     "202401",
				To:       "202412",
			},
			wantPath: "/bulk/list?endpoint=/equities/bars/daily&from=202401&to=202412",
		},
		{
			name:     "with no required parameters",
			params:   BulkListParams{},
			wantPath: "",
			wantErr:  true,
		},
		{
			name: "with from/to but no endpoint",
			params: BulkListParams{
				Date: "20250101",
				From: "202401",
				To:   "202412",
			},
			wantPath: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewBulkService(mockClient)

			// Mock response based on documentation sample
			mockResponse := BulkListResponse{
				Data: []BulkFile{
					{
						Key:          "equities/bars/daily/historical/2025/equities_bars_daily_202501.csv.gz",
						LastModified: "2025-11-07T20:48:51.295000+00:00",
						Size:         6933528,
					},
				},
			}
			if !tt.wantErr {
				mockClient.SetResponse("GET", tt.wantPath, mockResponse)
			}

			// Execute
			resp, err := service.GetFiles(context.Background(), tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetFiles() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetFiles() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetFiles() returned nil response")
			}
			if len(resp.Data) == 0 {
				t.Error("GetFiles() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetFiles() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestBulkService_GetFiles_RequiresParameter(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewBulkService(mockClient)

	// Execute with empty parameters
	_, err := service.GetFiles(context.Background(), BulkListParams{})

	// Verify
	if err == nil {
		t.Error("GetFiles() expected error for missing parameters but got nil")
	}
	if err.Error() != "either endpoint or date parameter is required" {
		t.Errorf("GetFiles() error = %v, want 'either endpoint or date parameter is required'", err)
	}
}

func TestBulkService_GetFiles_FromToRequiresEndpoint(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewBulkService(mockClient)

	// Execute with from/to but no endpoint
	_, err := service.GetFiles(context.Background(), BulkListParams{Date: "20250101", From: "202401"})

	// Verify
	if err == nil {
		t.Error("GetFiles() expected error for from/to without endpoint but got nil")
	}
	if err.Error() != "from/to parameters can only be used with endpoint parameter" {
		t.Errorf("GetFiles() error = %v, want 'from/to parameters can only be used with endpoint parameter'", err)
	}
}

func TestBulkListResponse_UnmarshalJSON(t *testing.T) {
	// ドキュメントのレスポンスサンプルに基づくJSON
	// Sizeはfloat64の丸めが発生する2^53超の値で厳密なint64パースを確認する
	raw := `{
		"data": [
			{
				"Key": "equities/bars/daily/historical/2025/equities_bars_daily_202501.csv.gz",
				"LastModified": "2025-11-07T20:48:51.295000+00:00",
				"Size": 9007199254740993
			},
			{
				"Key": "equities/bars/daily/historical/2024/equities_bars_daily_202412.csv.gz",
				"LastModified": "2025-01-07T18:30:15.123000+00:00",
				"Size": 6845123
			}
		]
	}`

	var resp BulkListResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if len(resp.Data) != 2 {
		t.Fatalf("UnmarshalJSON() returned %d items, want 2", len(resp.Data))
	}

	file := resp.Data[0]
	if file.Key != "equities/bars/daily/historical/2025/equities_bars_daily_202501.csv.gz" {
		t.Errorf("UnmarshalJSON() Key = %v", file.Key)
	}
	if file.LastModified != "2025-11-07T20:48:51.295000+00:00" {
		t.Errorf("UnmarshalJSON() LastModified = %v", file.LastModified)
	}
	if file.Size != 9007199254740993 {
		t.Errorf("UnmarshalJSON() Size = %v, want 9007199254740993", file.Size)
	}
	if resp.Data[1].Size != 6845123 {
		t.Errorf("UnmarshalJSON() Size = %v, want 6845123", resp.Data[1].Size)
	}
}

func TestBulkService_GetDownloadURL(t *testing.T) {
	tests := []struct {
		name     string
		params   BulkGetParams
		wantPath string
		wantErr  bool
	}{
		{
			name: "with key only",
			params: BulkGetParams{
				Key: "equities/bars/daily/historical/2025/equities_bars_daily_202501.csv.gz",
			},
			wantPath: "/bulk/get?key=equities/bars/daily/historical/2025/equities_bars_daily_202501.csv.gz",
		},
		{
			name: "with endpoint and date",
			params: BulkGetParams{
				Endpoint: "/equities/bars/daily",
				Date:     "202501",
			},
			wantPath: "/bulk/get?endpoint=/equities/bars/daily&date=202501",
		},
		{
			name: "with key and endpoint",
			params: BulkGetParams{
				Key:      "equities/bars/daily/historical/2025/equities_bars_daily_202501.csv.gz",
				Endpoint: "/equities/bars/daily",
			},
			wantErr: true,
		},
		{
			name: "with key and date",
			params: BulkGetParams{
				Key:  "equities/bars/daily/historical/2025/equities_bars_daily_202501.csv.gz",
				Date: "202501",
			},
			wantErr: true,
		},
		{
			name: "with endpoint only",
			params: BulkGetParams{
				Endpoint: "/equities/bars/daily",
			},
			wantErr: true,
		},
		{
			name: "with date only",
			params: BulkGetParams{
				Date: "202501",
			},
			wantErr: true,
		},
		{
			name:    "with no parameters",
			params:  BulkGetParams{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewBulkService(mockClient)

			wantURL := "https://example.presigned-url.com/file.csv.gz"
			if !tt.wantErr {
				mockClient.SetResponse("GET", tt.wantPath, map[string]string{"url": wantURL})
			}

			// Execute
			url, err := service.GetDownloadURL(context.Background(), tt.params)

			// Verify
			if tt.wantErr {
				if err == nil {
					t.Error("GetDownloadURL() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("GetDownloadURL() error = %v", err)
			}
			if url != wantURL {
				t.Errorf("GetDownloadURL() url = %v, want %v", url, wantURL)
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetDownloadURL() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestBulkService_GetDownloadURL_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewBulkService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/bulk/get?key=some/file.csv.gz", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetDownloadURL(context.Background(), BulkGetParams{Key: "some/file.csv.gz"})

	// Verify
	if err == nil {
		t.Error("GetDownloadURL() expected error but got nil")
	}
}

func TestBulkService_GetDownloadURL_SkipsCache(t *testing.T) {
	mockClient := client.NewMockClient()
	service := NewBulkService(mockClient)
	mockClient.SetResponse("GET", "/bulk/get?key=abc", map[string]string{"url": "https://example.com/file.csv.gz"})

	if _, err := service.GetDownloadURL(context.Background(), BulkGetParams{Key: "abc"}); err != nil {
		t.Fatalf("GetDownloadURL() error = %v", err)
	}
	// 署名付きURLは失効するため、キャッシュバイパスを要求していること
	if !mockClient.LastSkipCache {
		t.Error("GetDownloadURL() should request cache bypass for expiring signed URLs")
	}
}
