package jquants

import (
	"fmt"

	"github.com/utahta/jquants/client"
)

// FSDetailsService は財務諸表(BS/PL/CF)詳細情報を取得するサービスです。
// 四半期毎の詳細な財務諸表データを提供します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では利用できません。
type FSDetailsService struct {
	client client.HTTPClient
}

// NewFSDetailsService は新しいFSDetailsServiceを作成します。
func NewFSDetailsService(c client.HTTPClient) *FSDetailsService {
	return &FSDetailsService{client: c}
}

// FSDetailsParams は財務諸表詳細情報のリクエストパラメータです。
type FSDetailsParams struct {
	Code          string // 銘柄コード（codeまたはdateのいずれかが必須）
	Date          string // 開示日（YYYYMMDD または YYYY-MM-DD）（codeまたはdateのいずれかが必須）
	PaginationKey string // ページネーションキー
}

// FSDetailsResponse は財務諸表詳細情報のレスポンスです。
type FSDetailsResponse struct {
	Data          []FSDetail `json:"data"`
	PaginationKey string     `json:"pagination_key"` // ページネーションキー
}

// FSDetail は財務諸表詳細情報を表します。
// J-Quants API /fins/details エンドポイントのレスポンスデータ。
// 上場企業の四半期毎の財務情報における、貸借対照表、損益計算書、キャッシュ・フロー計算書に記載の項目。
//
// 注意: このデータはプレミアムプラン専用APIで取得されます。
type FSDetail struct {
	// 基本情報
	DiscDate string `json:"DiscDate"` // 開示日（YYYY-MM-DD形式）
	DiscTime string `json:"DiscTime"` // 開示時刻（HH:MM:SS形式）
	Code     string `json:"Code"`     // 銘柄コード（5桁）
	DiscNo   string `json:"DiscNo"`   // 開示番号（昇順ソートのキー）
	DocType  string `json:"DocType"`  // 開示書類種別

	// 財務諸表データ
	FS map[string]string `json:"FS"` // 財務諸表の各種項目（冗長ラベル（英語）とその値のマップ）
}

// GetFSDetails は財務諸表詳細情報を取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *FSDetailsService) GetFSDetails(params FSDetailsParams) (*FSDetailsResponse, error) {
	// codeまたはdateのいずれかが必須
	if params.Code == "" && params.Date == "" {
		return nil, fmt.Errorf("either code or date parameter is required")
	}

	path := "/fins/details"

	query := "?"
	if params.Code != "" {
		query += fmt.Sprintf("code=%s&", params.Code)
	}
	if params.Date != "" {
		query += fmt.Sprintf("date=%s&", params.Date)
	}
	if params.PaginationKey != "" {
		query += fmt.Sprintf("pagination_key=%s&", params.PaginationKey)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp FSDetailsResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get fs details: %w", err)
	}

	return &resp, nil
}

// GetFSDetailsByCode は指定銘柄の財務諸表詳細情報を取得します。
// ページネーションを使用して全データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *FSDetailsService) GetFSDetailsByCode(code string) ([]FSDetail, error) {
	var allData []FSDetail
	paginationKey := ""

	for {
		params := FSDetailsParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetFSDetails(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.Data...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetFSDetailsByDate は指定日の全銘柄財務諸表詳細情報を取得します。
// ページネーションを使用して全データを取得します。
func (s *FSDetailsService) GetFSDetailsByDate(date string) ([]FSDetail, error) {
	var allData []FSDetail
	paginationKey := ""

	for {
		params := FSDetailsParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetFSDetails(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.Data...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// IsIFRS は財務諸表がIFRS基準かどうかを判定します。
func (d *FSDetail) IsIFRS() bool {
	if standards, ok := d.FS["Accounting standards, DEI"]; ok {
		return standards == "IFRS"
	}
	return false
}

// IsJapaneseGAAP は財務諸表が日本基準かどうかを判定します。
func (d *FSDetail) IsJapaneseGAAP() bool {
	if standards, ok := d.FS["Accounting standards, DEI"]; ok {
		return standards == "JapaneseGAAP"
	}
	return false
}

// IsQuarterly は四半期財務諸表かどうかを判定します。
func (d *FSDetail) IsQuarterly() bool {
	return contains(d.DocType, "Q") && contains(d.DocType, "Financial")
}

// IsAnnual は年次財務諸表かどうかを判定します。
func (d *FSDetail) IsAnnual() bool {
	return contains(d.DocType, "FY") && contains(d.DocType, "Financial")
}

// IsConsolidated は連結財務諸表かどうかを判定します。
func (d *FSDetail) IsConsolidated() bool {
	return contains(d.DocType, "Consolidated")
}

// GetQuarter は四半期を取得します（1, 2, 3, 0（年次））。
func (d *FSDetail) GetQuarter() int {
	if contains(d.DocType, "1Q") {
		return 1
	} else if contains(d.DocType, "2Q") {
		return 2
	} else if contains(d.DocType, "3Q") {
		return 3
	} else if contains(d.DocType, "FY") {
		return 0
	}
	return -1
}

// GetValue は指定されたキーの値を取得します。
func (d *FSDetail) GetValue(key string) (string, bool) {
	value, ok := d.FS[key]
	return value, ok
}

// GetFloatValue は指定されたキーの値をfloat64として取得します。
func (d *FSDetail) GetFloatValue(key string) (float64, error) {
	value, ok := d.FS[key]
	if !ok {
		return 0, fmt.Errorf("key %s not found", key)
	}

	var floatVal float64
	_, err := fmt.Sscanf(value, "%f", &floatVal)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float value for key %s: %w", key, err)
	}

	return floatVal, nil
}

// Financial ratios and analysis methods

// GetROE は自己資本利益率（ROE）を計算します（IFRS）。
func (d *FSDetail) GetROE() (*float64, error) {
	if !d.IsIFRS() {
		return nil, fmt.Errorf("ROE calculation is only supported for IFRS")
	}

	netIncome, err := d.GetFloatValue("Profit (loss) attributable to owners of parent (IFRS)")
	if err != nil {
		return nil, err
	}

	equity, err := d.GetFloatValue("Equity attributable to owners of parent (IFRS)")
	if err != nil {
		return nil, err
	}

	if equity == 0 {
		return nil, fmt.Errorf("equity is zero")
	}

	roe := (netIncome / equity) * 100
	return &roe, nil
}

// GetCurrentRatio は流動比率を計算します（IFRS）。
func (d *FSDetail) GetCurrentRatio() (*float64, error) {
	if !d.IsIFRS() {
		return nil, fmt.Errorf("current ratio calculation is only supported for IFRS")
	}

	currentAssets, err := d.GetFloatValue("Current assets (IFRS)")
	if err != nil {
		return nil, err
	}

	currentLiabilities, err := d.GetFloatValue("Current liabilities (IFRS)")
	if err != nil {
		return nil, err
	}

	if currentLiabilities == 0 {
		return nil, fmt.Errorf("current liabilities is zero")
	}

	ratio := currentAssets / currentLiabilities
	return &ratio, nil
}

// GetEquityRatio は自己資本比率を計算します（IFRS）。
func (d *FSDetail) GetEquityRatio() (*float64, error) {
	if !d.IsIFRS() {
		return nil, fmt.Errorf("equity ratio calculation is only supported for IFRS")
	}

	equity, err := d.GetFloatValue("Equity (IFRS)")
	if err != nil {
		return nil, err
	}

	assets, err := d.GetFloatValue("Assets (IFRS)")
	if err != nil {
		return nil, err
	}

	if assets == 0 {
		return nil, fmt.Errorf("assets is zero")
	}

	ratio := (equity / assets) * 100
	return &ratio, nil
}

// GetBasicEPS は基本的1株当たり利益を取得します（IFRS）。
func (d *FSDetail) GetBasicEPS() (*float64, error) {
	if !d.IsIFRS() {
		return nil, fmt.Errorf("basic EPS is only available for IFRS")
	}

	eps, err := d.GetFloatValue("Basic earnings (loss) per share (IFRS)")
	if err != nil {
		return nil, err
	}

	return &eps, nil
}

// contains は文字列に部分文字列が含まれるかをチェックするヘルパー関数
func contains(s, substr string) bool {
	return len(s) >= len(substr) && len(substr) > 0 &&
		(s == substr || len(s) > len(substr) &&
			(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
				func() bool {
					for i := 1; i < len(s)-len(substr); i++ {
						if s[i:i+len(substr)] == substr {
							return true
						}
					}
					return false
				}()))
}
