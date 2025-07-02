package jquants

import (
	"encoding/json"
	"strings"
)

// TypeOfDocument は財務諸表の開示書類種別を表す型です
type TypeOfDocument string

// 開示書類種別の定義
// 財務諸表API (/fins/statements) のTypeOfDocumentフィールドで使用される値です
const (
	// 決算短信（通期・連結）
	TypeOfDocumentFYConsolidatedJP      TypeOfDocument = "FYFinancialStatements_Consolidated_JP"      // 決算短信（連結・日本基準）
	TypeOfDocumentFYConsolidatedUS      TypeOfDocument = "FYFinancialStatements_Consolidated_US"      // 決算短信（連結・米国基準）
	TypeOfDocumentFYConsolidatedIFRS    TypeOfDocument = "FYFinancialStatements_Consolidated_IFRS"    // 決算短信（連結・IFRS）
	TypeOfDocumentFYConsolidatedJMIS    TypeOfDocument = "FYFinancialStatements_Consolidated_JMIS"    // 決算短信（連結・JMIS）
	TypeOfDocumentFYConsolidatedREIT    TypeOfDocument = "FYFinancialStatements_Consolidated_REIT"    // 決算短信（REIT）
	TypeOfDocumentFYConsolidatedForeign TypeOfDocument = "FYFinancialStatements_Consolidated_Foreign" // 決算短信（連結・外国株）

	// 決算短信（通期・非連結）
	TypeOfDocumentFYNonConsolidatedJP      TypeOfDocument = "FYFinancialStatements_NonConsolidated_JP"      // 決算短信（非連結・日本基準）
	TypeOfDocumentFYNonConsolidatedIFRS    TypeOfDocument = "FYFinancialStatements_NonConsolidated_IFRS"    // 決算短信（非連結・IFRS）
	TypeOfDocumentFYNonConsolidatedForeign TypeOfDocument = "FYFinancialStatements_NonConsolidated_Foreign" // 決算短信（非連結・外国株）

	// 第1四半期決算短信（連結）
	TypeOfDocument1QConsolidatedJP   TypeOfDocument = "1QFinancialStatements_Consolidated_JP"   // 第1四半期決算短信（連結・日本基準）
	TypeOfDocument1QConsolidatedUS   TypeOfDocument = "1QFinancialStatements_Consolidated_US"   // 第1四半期決算短信（連結・米国基準）
	TypeOfDocument1QConsolidatedIFRS TypeOfDocument = "1QFinancialStatements_Consolidated_IFRS" // 第1四半期決算短信（連結・IFRS）

	// 第2四半期決算短信（連結）
	TypeOfDocument2QConsolidatedJP   TypeOfDocument = "2QFinancialStatements_Consolidated_JP"   // 第2四半期決算短信（連結・日本基準）
	TypeOfDocument2QConsolidatedUS   TypeOfDocument = "2QFinancialStatements_Consolidated_US"   // 第2四半期決算短信（連結・米国基準）
	TypeOfDocument2QConsolidatedIFRS TypeOfDocument = "2QFinancialStatements_Consolidated_IFRS" // 第2四半期決算短信（連結・IFRS）

	// 第3四半期決算短信（連結）
	TypeOfDocument3QConsolidatedJP   TypeOfDocument = "3QFinancialStatements_Consolidated_JP"   // 第3四半期決算短信（連結・日本基準）
	TypeOfDocument3QConsolidatedUS   TypeOfDocument = "3QFinancialStatements_Consolidated_US"   // 第3四半期決算短信（連結・米国基準）
	TypeOfDocument3QConsolidatedIFRS TypeOfDocument = "3QFinancialStatements_Consolidated_IFRS" // 第3四半期決算短信（連結・IFRS）

	// 予想修正
	TypeOfDocumentDividendRevision     TypeOfDocument = "DividendForecastRevision"     // 配当予想の修正
	TypeOfDocumentEarningsRevision     TypeOfDocument = "EarnForecastRevision"         // 業績予想の修正
	TypeOfDocumentREITDividendRevision TypeOfDocument = "REITDividendForecastRevision" // 分配予想の修正（REIT）
	TypeOfDocumentREITEarningsRevision TypeOfDocument = "REITEarnForecastRevision"     // 利益予想の修正（REIT）
)

// String はTypeOfDocumentを文字列として返します
func (t TypeOfDocument) String() string {
	return string(t)
}

// IsConsolidated は連結財務諸表かどうかを判定します
func (t TypeOfDocument) IsConsolidated() bool {
	s := string(t)
	return strings.Contains(s, "_Consolidated_") || strings.Contains(s, "_REIT")
}

// IsNonConsolidated は非連結（単体）財務諸表かどうかを判定します
func (t TypeOfDocument) IsNonConsolidated() bool {
	return strings.Contains(string(t), "_NonConsolidated_")
}

// IsQuarterly は四半期決算短信かどうかを判定します
func (t TypeOfDocument) IsQuarterly() bool {
	s := string(t)
	return strings.Contains(s, "1QFinancial") ||
		strings.Contains(s, "2QFinancial") ||
		strings.Contains(s, "3QFinancial") ||
		strings.Contains(s, "OtherPeriodFinancial")
}

// IsAnnual は通期決算短信かどうかを判定します
func (t TypeOfDocument) IsAnnual() bool {
	return strings.HasPrefix(string(t), "FYFinancial")
}

// IsForecastRevision は予想修正かどうかを判定します
func (t TypeOfDocument) IsForecastRevision() bool {
	s := string(t)
	return strings.Contains(s, "ForecastRevision")
}

// IsREIT はREIT関連の書類かどうかを判定します
func (t TypeOfDocument) IsREIT() bool {
	return strings.Contains(string(t), "REIT")
}

// GetAccountingStandard は会計基準を返します
// 例: "JP" (日本基準), "US" (米国基準), "IFRS" (国際財務報告基準), "JMIS" (修正国際基準), "Foreign" (外国株), "REIT"
func (t TypeOfDocument) GetAccountingStandard() string {
	s := string(t)
	
	// 予想修正の場合は空文字を返す
	if t.IsForecastRevision() {
		return ""
	}
	
	// 最後のアンダースコアの後ろが会計基準
	parts := strings.Split(s, "_")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	
	return ""
}

// GetPeriod は期間を返します
// 例: "FY" (通期), "1Q" (第1四半期), "2Q" (第2四半期), "3Q" (第3四半期), "OtherPeriod" (その他)
func (t TypeOfDocument) GetPeriod() string {
	s := string(t)
	
	// 予想修正の場合は空文字を返す
	if t.IsForecastRevision() {
		return ""
	}
	
	if strings.HasPrefix(s, "FY") {
		return "FY"
	} else if strings.HasPrefix(s, "1Q") {
		return "1Q"
	} else if strings.HasPrefix(s, "2Q") {
		return "2Q"
	} else if strings.HasPrefix(s, "3Q") {
		return "3Q"
	} else if strings.HasPrefix(s, "OtherPeriod") {
		return "OtherPeriod"
	}
	
	return ""
}

// ParseTypeOfDocument は文字列をTypeOfDocument型に変換します
// 無効な値の場合でもそのまま返します（後方互換性のため）
func ParseTypeOfDocument(s string) TypeOfDocument {
	return TypeOfDocument(s)
}

// UnmarshalJSON はJSONからTypeOfDocumentをデコードします
func (t *TypeOfDocument) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*t = TypeOfDocument(s)
	return nil
}

// MarshalJSON はTypeOfDocumentをJSONにエンコードします
func (t TypeOfDocument) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(t))
}