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
	Sector33Code  string // 33業種コード（s33またはdateのいずれかが必須）
	Date          string // 日付（YYYYMMDD または YYYY-MM-DD）（s33またはdateのいずれかが必須）
	From          string // 期間の開始日（YYYYMMDD または YYYY-MM-DD）
	To            string // 期間の終了日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// ShortSellingResponse は業種別空売り比率のレスポンスです。
type ShortSellingResponse struct {
	Data          []ShortSelling `json:"data"`
	PaginationKey string         `json:"pagination_key"` // ページネーションキー
}

// ShortSelling は業種別空売り比率のデータを表します。
// J-Quants API /markets/short-ratio エンドポイントのレスポンスデータ。
type ShortSelling struct {
	// 基本情報
	Date string `json:"Date"` // 日付（YYYY-MM-DD形式）
	S33  string `json:"S33"`  // 33業種コード

	// 売買代金データ（単位：円）
	SellExShortVa float64 `json:"SellExShortVa"` // 実注文の売買代金（空売り以外の通常売り注文）
	ShrtWithResVa float64 `json:"ShrtWithResVa"` // 価格規制有りの空売り売買代金（アップティック・ルール等）
	ShrtNoResVa   float64 `json:"ShrtNoResVa"`   // 価格規制無しの空売り売買代金（ETF、REIT、裁定取引等）
}

// RawShortSelling is used for unmarshaling JSON response with mixed types
type RawShortSelling struct {
	// 基本情報
	Date string `json:"Date"`
	S33  string `json:"S33"`

	// 売買代金データ
	SellExShortVa types.Float64String `json:"SellExShortVa"`
	ShrtWithResVa types.Float64String `json:"ShrtWithResVa"`
	ShrtNoResVa   types.Float64String `json:"ShrtNoResVa"`
}

// UnmarshalJSON implements custom JSON unmarshaling for ShortSellingResponse
func (r *ShortSellingResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawShortSelling
	type rawResponse struct {
		Data          []RawShortSelling `json:"data"`
		PaginationKey string            `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawShortSelling to ShortSelling
	r.Data = make([]ShortSelling, len(raw.Data))
	for idx, rs := range raw.Data {
		r.Data[idx] = ShortSelling{
			// 基本情報
			Date: rs.Date,
			S33:  rs.S33,

			// 売買代金データ
			SellExShortVa: float64(rs.SellExShortVa),
			ShrtWithResVa: float64(rs.ShrtWithResVa),
			ShrtNoResVa:   float64(rs.ShrtNoResVa),
		}
	}

	return nil
}

// GetShortSelling は業種別空売り比率を取得します。
func (s *ShortSellingService) GetShortSelling(params ShortSellingParams) (*ShortSellingResponse, error) {
	// s33またはdateのいずれかが必須
	if params.Sector33Code == "" && params.Date == "" {
		return nil, fmt.Errorf("either s33 or date parameter is required")
	}

	path := "/markets/short-ratio"

	query := "?"
	if params.Sector33Code != "" {
		query += fmt.Sprintf("s33=%s&", params.Sector33Code)
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

		allData = append(allData, resp.Data...)

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

		allData = append(allData, resp.Data...)

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

		allData = append(allData, resp.Data...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetShortSellingBySectorAndDate は指定業種の指定日の空売り比率を取得します。
func (s *ShortSellingService) GetShortSellingBySectorAndDate(sector33Code, date string) ([]ShortSelling, error) {
	resp, err := s.GetShortSelling(ShortSellingParams{
		Sector33Code: sector33Code,
		Date:         date,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetTotalShortSellingValue は空売り合計金額を計算します。
func (ss *ShortSelling) GetTotalShortSellingValue() float64 {
	return ss.ShrtWithResVa + ss.ShrtNoResVa
}

// GetTotalTurnoverValue は総売買代金を計算します。
func (ss *ShortSelling) GetTotalTurnoverValue() float64 {
	return ss.SellExShortVa + ss.GetTotalShortSellingValue()
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
	return (ss.ShrtWithResVa / totalShortSelling) * 100
}

// GetUnrestrictedShortSellingRatio は価格規制なし空売りの割合を計算します（パーセント）。
func (ss *ShortSelling) GetUnrestrictedShortSellingRatio() float64 {
	totalShortSelling := ss.GetTotalShortSellingValue()
	if totalShortSelling == 0 {
		return 0
	}
	return (ss.ShrtNoResVa / totalShortSelling) * 100
}
