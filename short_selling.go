package jquants

import (
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// ShortSellingService は業種別空売り比率を取得するサービスです。
type ShortSellingService struct {
	client client.HTTPClient
}

// NewShortSellingService は新しいShortSellingServiceを作成します。
func NewShortSellingService(c client.HTTPClient) *ShortSellingService {
	return &ShortSellingService{client: c}
}

// ShortSellingParams は業種別空売り比率のリクエストパラメータです。
type ShortSellingParams struct {
	Sector33Code  string // 33業種コード（sector33codeまたはdateのいずれかが必須）
	Date          string // 日付（YYYYMMDD または YYYY-MM-DD）（sector33codeまたはdateのいずれかが必須）
	From          string // 期間の開始日（YYYYMMDD または YYYY-MM-DD）
	To            string // 期間の終了日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// ShortSellingResponse は業種別空売り比率のレスポンスです。
type ShortSellingResponse struct {
	ShortSelling  []ShortSelling `json:"short_selling"`
	PaginationKey string         `json:"pagination_key"` // ページネーションキー
}

// ShortSelling は業種別空売り比率のデータを表します。
// J-Quants API /markets/short_selling エンドポイントのレスポンスデータ。
type ShortSelling struct {
	// 基本情報
	Date         string `json:"Date"`         // 日付（YYYY-MM-DD形式）
	Sector33Code string `json:"Sector33Code"` // 33業種コード

	// 売買代金データ（単位：円）
	SellingExcludingShortSellingTurnoverValue    float64 `json:"SellingExcludingShortSellingTurnoverValue"`    // 実注文の売買代金（空売り以外の通常売り注文）
	ShortSellingWithRestrictionsTurnoverValue    float64 `json:"ShortSellingWithRestrictionsTurnoverValue"`    // 価格規制有りの空売り売買代金（アップティック・ルール等）
	ShortSellingWithoutRestrictionsTurnoverValue float64 `json:"ShortSellingWithoutRestrictionsTurnoverValue"` // 価格規制無しの空売り売買代金（ETF、REIT、裁定取引等）
}

// RawShortSelling is used for unmarshaling JSON response with mixed types
type RawShortSelling struct {
	// 基本情報
	Date         string `json:"Date"`
	Sector33Code string `json:"Sector33Code"`

	// 売買代金データ
	SellingExcludingShortSellingTurnoverValue    types.Float64String `json:"SellingExcludingShortSellingTurnoverValue"`
	ShortSellingWithRestrictionsTurnoverValue    types.Float64String `json:"ShortSellingWithRestrictionsTurnoverValue"`
	ShortSellingWithoutRestrictionsTurnoverValue types.Float64String `json:"ShortSellingWithoutRestrictionsTurnoverValue"`
}

// UnmarshalJSON implements custom JSON unmarshaling for ShortSellingResponse
func (r *ShortSellingResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawShortSelling
	type rawResponse struct {
		ShortSelling  []RawShortSelling `json:"short_selling"`
		PaginationKey string            `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawShortSelling to ShortSelling
	r.ShortSelling = make([]ShortSelling, len(raw.ShortSelling))
	for idx, rs := range raw.ShortSelling {
		r.ShortSelling[idx] = ShortSelling{
			// 基本情報
			Date:         rs.Date,
			Sector33Code: rs.Sector33Code,

			// 売買代金データ
			SellingExcludingShortSellingTurnoverValue:    float64(rs.SellingExcludingShortSellingTurnoverValue),
			ShortSellingWithRestrictionsTurnoverValue:    float64(rs.ShortSellingWithRestrictionsTurnoverValue),
			ShortSellingWithoutRestrictionsTurnoverValue: float64(rs.ShortSellingWithoutRestrictionsTurnoverValue),
		}
	}

	return nil
}

// GetShortSelling は業種別空売り比率を取得します。
func (s *ShortSellingService) GetShortSelling(params ShortSellingParams) (*ShortSellingResponse, error) {
	// sector33codeまたはdateのいずれかが必須
	if params.Sector33Code == "" && params.Date == "" {
		return nil, fmt.Errorf("either sector33code or date parameter is required")
	}

	path := "/markets/short_selling"

	query := "?"
	if params.Sector33Code != "" {
		query += fmt.Sprintf("sector33code=%s&", params.Sector33Code)
	}
	if params.Date != "" {
		query += fmt.Sprintf("date=%s&", params.Date)
	}
	if params.From != "" {
		query += fmt.Sprintf("from=%s&", params.From)
	}
	if params.To != "" {
		query += fmt.Sprintf("to=%s&", params.To)
	}
	if params.PaginationKey != "" {
		query += fmt.Sprintf("pagination_key=%s&", params.PaginationKey)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp ShortSellingResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get short selling: %w", err)
	}

	return &resp, nil
}

// GetShortSellingBySector は指定業種の空売り比率を取得します。
// ページネーションを使用して全データを取得します。
func (s *ShortSellingService) GetShortSellingBySector(sector33Code string) ([]ShortSelling, error) {
	var allData []ShortSelling
	paginationKey := ""

	for {
		params := ShortSellingParams{
			Sector33Code:  sector33Code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetShortSelling(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.ShortSelling...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetShortSellingByDate は指定日の全業種空売り比率を取得します。
// ページネーションを使用して全データを取得します。
func (s *ShortSellingService) GetShortSellingByDate(date string) ([]ShortSelling, error) {
	var allData []ShortSelling
	paginationKey := ""

	for {
		params := ShortSellingParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetShortSelling(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.ShortSelling...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetShortSellingBySectorAndDateRange は指定業種・期間の空売り比率を取得します。
func (s *ShortSellingService) GetShortSellingBySectorAndDateRange(sector33Code, from, to string) ([]ShortSelling, error) {
	var allData []ShortSelling
	paginationKey := ""

	for {
		params := ShortSellingParams{
			Sector33Code:  sector33Code,
			From:          from,
			To:            to,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetShortSelling(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.ShortSelling...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetTotalShortSellingValue は空売り合計金額を計算します。
func (ss *ShortSelling) GetTotalShortSellingValue() float64 {
	return ss.ShortSellingWithRestrictionsTurnoverValue + ss.ShortSellingWithoutRestrictionsTurnoverValue
}

// GetTotalTurnoverValue は総売買代金を計算します。
func (ss *ShortSelling) GetTotalTurnoverValue() float64 {
	return ss.SellingExcludingShortSellingTurnoverValue + ss.GetTotalShortSellingValue()
}

// GetShortSellingRatio は空売り比率を計算します（パーセント）。
func (ss *ShortSelling) GetShortSellingRatio() float64 {
	totalTurnover := ss.GetTotalTurnoverValue()
	if totalTurnover == 0 {
		return 0
	}
	return (ss.GetTotalShortSellingValue() / totalTurnover) * 100
}

// GetRestrictedShortSellingRatio は価格規制付き空売りの割合を計算します（パーセント）。
func (ss *ShortSelling) GetRestrictedShortSellingRatio() float64 {
	totalShortSelling := ss.GetTotalShortSellingValue()
	if totalShortSelling == 0 {
		return 0
	}
	return (ss.ShortSellingWithRestrictionsTurnoverValue / totalShortSelling) * 100
}

// GetUnrestrictedShortSellingRatio は価格規制なし空売りの割合を計算します（パーセント）。
func (ss *ShortSelling) GetUnrestrictedShortSellingRatio() float64 {
	totalShortSelling := ss.GetTotalShortSellingValue()
	if totalShortSelling == 0 {
		return 0
	}
	return (ss.ShortSellingWithoutRestrictionsTurnoverValue / totalShortSelling) * 100
}
