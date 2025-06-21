package types

import (
	"encoding/json"
	"testing"
)

func TestFloat64String_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		want     *float64
		wantErr  bool
	}{
		{
			name: "float64値",
			json: `{"value": 123.45}`,
			want: floatPtr(123.45),
		},
		{
			name: "文字列値",
			json: `{"value": "678.90"}`,
			want: floatPtr(678.90),
		},
		{
			name: "空文字列",
			json: `{"value": ""}`,
			want: floatPtr(0), // 空文字列は0として扱われる
		},
		{
			name: "ハイフン",
			json: `{"value": "-"}`,
			want: floatPtr(0), // ハイフンは0として扱われる
		},
		{
			name: "不正な文字列",
			json: `{"value": "abc"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data struct {
				Value *Float64String `json:"value"`
			}
			
			err := json.Unmarshal([]byte(tt.json), &data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err == nil {
				if tt.want == nil && data.Value != nil {
					t.Errorf("Expected nil, got %v", *data.Value)
				} else if tt.want != nil && data.Value == nil {
					t.Errorf("Expected %v, got nil", *tt.want)
				} else if tt.want != nil && data.Value != nil {
					if float64(*data.Value) != *tt.want {
						t.Errorf("Expected %v, got %v", *tt.want, float64(*data.Value))
					}
				}
			}
		})
	}
}

func TestBoolString_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    bool
		wantErr bool
	}{
		{
			name: "bool true",
			json: `{"value": true}`,
			want: true,
		},
		{
			name: "bool false",
			json: `{"value": false}`,
			want: false,
		},
		{
			name: "文字列 true",
			json: `{"value": "true"}`,
			want: true,
		},
		{
			name: "文字列 false",
			json: `{"value": "false"}`,
			want: false,
		},
		{
			name: "文字列 1",
			json: `{"value": "1"}`,
			want: true,
		},
		{
			name: "文字列 0",
			json: `{"value": "0"}`,
			want: false,
		},
		{
			name: "空文字列",
			json: `{"value": ""}`,
			want: false,
		},
		{
			name: "不正な文字列",
			json: `{"value": "invalid"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data struct {
				Value BoolString `json:"value"`
			}
			
			err := json.Unmarshal([]byte(tt.json), &data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err == nil && bool(data.Value) != tt.want {
				t.Errorf("Expected %v, got %v", tt.want, bool(data.Value))
			}
		})
	}
}

func TestInt64String_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		json    string
		want    *int64
		wantErr bool
	}{
		{
			name: "int64値",
			json: `{"value": 12345}`,
			want: int64Ptr(12345),
		},
		{
			name: "文字列値",
			json: `{"value": "67890"}`,
			want: int64Ptr(67890),
		},
		{
			name: "空文字列",
			json: `{"value": ""}`,
			want: int64Ptr(0), // 空文字列は0として扱われる
		},
		{
			name: "ハイフン",
			json: `{"value": "-"}`,
			want: int64Ptr(0), // ハイフンは0として扱われる
		},
		{
			name: "不正な文字列",
			json: `{"value": "abc"}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var data struct {
				Value *Int64String `json:"value"`
			}
			
			err := json.Unmarshal([]byte(tt.json), &data)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if err == nil {
				if tt.want == nil && data.Value != nil {
					t.Errorf("Expected nil, got %v", *data.Value)
				} else if tt.want != nil && data.Value == nil {
					t.Errorf("Expected %v, got nil", *tt.want)
				} else if tt.want != nil && data.Value != nil {
					if int64(*data.Value) != *tt.want {
						t.Errorf("Expected %v, got %v", *tt.want, int64(*data.Value))
					}
				}
			}
		})
	}
}

func floatPtr(f float64) *float64 {
	return &f
}

func int64Ptr(i int64) *int64 {
	return &i
}