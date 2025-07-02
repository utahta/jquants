package jquants

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestTypeOfDocument_IsConsolidated(t *testing.T) {
	tests := []struct {
		name     string
		doc      TypeOfDocument
		expected bool
	}{
		{
			name:     "連結・日本基準",
			doc:      TypeOfDocumentFYConsolidatedJP,
			expected: true,
		},
		{
			name:     "連結・IFRS",
			doc:      TypeOfDocument2QConsolidatedIFRS,
			expected: true,
		},
		{
			name:     "REIT",
			doc:      TypeOfDocumentFYConsolidatedREIT,
			expected: true,
		},
		{
			name:     "非連結",
			doc:      TypeOfDocumentFYNonConsolidatedJP,
			expected: false,
		},
		{
			name:     "予想修正",
			doc:      TypeOfDocumentDividendRevision,
			expected: false,
		},
		{
			name:     "カスタム連結",
			doc:      ParseTypeOfDocument("OtherPeriodFinancialStatements_Consolidated_JMIS"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.doc.IsConsolidated(); got != tt.expected {
				t.Errorf("IsConsolidated() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTypeOfDocument_IsQuarterly(t *testing.T) {
	tests := []struct {
		name     string
		doc      TypeOfDocument
		expected bool
	}{
		{
			name:     "第1四半期",
			doc:      TypeOfDocument1QConsolidatedJP,
			expected: true,
		},
		{
			name:     "第2四半期",
			doc:      TypeOfDocument2QConsolidatedIFRS,
			expected: true,
		},
		{
			name:     "第3四半期",
			doc:      TypeOfDocument3QConsolidatedUS,
			expected: true,
		},
		{
			name:     "通期",
			doc:      TypeOfDocumentFYConsolidatedJP,
			expected: false,
		},
		{
			name:     "その他四半期",
			doc:      ParseTypeOfDocument("OtherPeriodFinancialStatements_Consolidated_JP"),
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.doc.IsQuarterly(); got != tt.expected {
				t.Errorf("IsQuarterly() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTypeOfDocument_GetAccountingStandard(t *testing.T) {
	tests := []struct {
		name     string
		doc      TypeOfDocument
		expected string
	}{
		{
			name:     "日本基準",
			doc:      TypeOfDocumentFYConsolidatedJP,
			expected: "JP",
		},
		{
			name:     "米国基準",
			doc:      TypeOfDocument1QConsolidatedUS,
			expected: "US",
		},
		{
			name:     "IFRS",
			doc:      TypeOfDocument2QConsolidatedIFRS,
			expected: "IFRS",
		},
		{
			name:     "JMIS",
			doc:      ParseTypeOfDocument("3QFinancialStatements_Consolidated_JMIS"),
			expected: "JMIS",
		},
		{
			name:     "REIT",
			doc:      TypeOfDocumentFYConsolidatedREIT,
			expected: "REIT",
		},
		{
			name:     "外国株",
			doc:      TypeOfDocumentFYConsolidatedForeign,
			expected: "Foreign",
		},
		{
			name:     "予想修正（会計基準なし）",
			doc:      TypeOfDocumentDividendRevision,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.doc.GetAccountingStandard(); got != tt.expected {
				t.Errorf("GetAccountingStandard() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTypeOfDocument_GetPeriod(t *testing.T) {
	tests := []struct {
		name     string
		doc      TypeOfDocument
		expected string
	}{
		{
			name:     "通期",
			doc:      TypeOfDocumentFYConsolidatedJP,
			expected: "FY",
		},
		{
			name:     "第1四半期",
			doc:      TypeOfDocument1QConsolidatedJP,
			expected: "1Q",
		},
		{
			name:     "第2四半期",
			doc:      TypeOfDocument2QConsolidatedIFRS,
			expected: "2Q",
		},
		{
			name:     "第3四半期",
			doc:      TypeOfDocument3QConsolidatedUS,
			expected: "3Q",
		},
		{
			name:     "その他四半期",
			doc:      ParseTypeOfDocument("OtherPeriodFinancialStatements_Consolidated_JP"),
			expected: "OtherPeriod",
		},
		{
			name:     "予想修正（期間なし）",
			doc:      TypeOfDocumentEarningsRevision,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.doc.GetPeriod(); got != tt.expected {
				t.Errorf("GetPeriod() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTypeOfDocument_IsForecastRevision(t *testing.T) {
	tests := []struct {
		name     string
		doc      TypeOfDocument
		expected bool
	}{
		{
			name:     "配当予想修正",
			doc:      TypeOfDocumentDividendRevision,
			expected: true,
		},
		{
			name:     "業績予想修正",
			doc:      TypeOfDocumentEarningsRevision,
			expected: true,
		},
		{
			name:     "REIT分配予想修正",
			doc:      TypeOfDocumentREITDividendRevision,
			expected: true,
		},
		{
			name:     "通常の決算短信",
			doc:      TypeOfDocumentFYConsolidatedJP,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.doc.IsForecastRevision(); got != tt.expected {
				t.Errorf("IsForecastRevision() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTypeOfDocument_IsREIT(t *testing.T) {
	tests := []struct {
		name     string
		doc      TypeOfDocument
		expected bool
	}{
		{
			name:     "REIT決算短信",
			doc:      TypeOfDocumentFYConsolidatedREIT,
			expected: true,
		},
		{
			name:     "REIT分配予想修正",
			doc:      TypeOfDocumentREITDividendRevision,
			expected: true,
		},
		{
			name:     "REIT利益予想修正",
			doc:      TypeOfDocumentREITEarningsRevision,
			expected: true,
		},
		{
			name:     "通常の決算短信",
			doc:      TypeOfDocumentFYConsolidatedJP,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.doc.IsREIT(); got != tt.expected {
				t.Errorf("IsREIT() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParseTypeOfDocument(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  TypeOfDocument
	}{
		{
			name:  "標準的な値",
			input: "FYFinancialStatements_Consolidated_JP",
			want:  TypeOfDocumentFYConsolidatedJP,
		},
		{
			name:  "カスタム値",
			input: "CustomFinancialStatements_Special_Type",
			want:  TypeOfDocument("CustomFinancialStatements_Special_Type"),
		},
		{
			name:  "空文字",
			input: "",
			want:  TypeOfDocument(""),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseTypeOfDocument(tt.input); got != tt.want {
				t.Errorf("ParseTypeOfDocument() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTypeOfDocument_JSON(t *testing.T) {
	tests := []struct {
		name string
		doc  TypeOfDocument
		json string
	}{
		{
			name: "標準的な値",
			doc:  TypeOfDocumentFYConsolidatedJP,
			json: `"FYFinancialStatements_Consolidated_JP"`,
		},
		{
			name: "カスタム値",
			doc:  TypeOfDocument("CustomType"),
			json: `"CustomType"`,
		},
		{
			name: "空文字",
			doc:  TypeOfDocument(""),
			json: `""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal
			got, err := json.Marshal(tt.doc)
			if err != nil {
				t.Fatalf("Marshal failed: %v", err)
			}
			if string(got) != tt.json {
				t.Errorf("Marshal() = %s, want %s", got, tt.json)
			}

			// Unmarshal
			var doc TypeOfDocument
			if err := json.Unmarshal([]byte(tt.json), &doc); err != nil {
				t.Fatalf("Unmarshal failed: %v", err)
			}
			if doc != tt.doc {
				t.Errorf("Unmarshal() = %v, want %v", doc, tt.doc)
			}
		})
	}
}

func TestStatement_TypeOfDocument_JSON(t *testing.T) {
	// Statement構造体でのJSON変換テスト
	stmt := Statement{
		TypeOfDocument: TypeOfDocumentFYConsolidatedIFRS,
		LocalCode:      "12345",
	}

	// Marshal
	data, err := json.Marshal(stmt)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// JSONに正しい文字列が含まれているか確認
	if !strings.Contains(string(data), `"TypeOfDocument":"FYFinancialStatements_Consolidated_IFRS"`) {
		t.Errorf("JSON does not contain expected TypeOfDocument string: %s", data)
	}

	// Unmarshal
	var decoded Statement
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if decoded.TypeOfDocument != TypeOfDocumentFYConsolidatedIFRS {
		t.Errorf("Unmarshal TypeOfDocument = %v, want %v", decoded.TypeOfDocument, TypeOfDocumentFYConsolidatedIFRS)
	}
}