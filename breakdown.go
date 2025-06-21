package jquants

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// BreakdownService は売買内訳データを取得するサービスです。
// 東証上場銘柄の東証市場における銘柄別の日次売買代金・売買高（立会内取引に限る）について、
// 信用取引や空売りの利用に関する発注時のフラグ情報を用いて細分化したデータを提供します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では利用できません。
type BreakdownService struct {
	client client.HTTPClient
}

// NewBreakdownService は新しいBreakdownServiceを作成します。
func NewBreakdownService(c client.HTTPClient) *BreakdownService {
	return &BreakdownService{client: c}
}

// BreakdownParams は売買内訳データのリクエストパラメータです。
type BreakdownParams struct {
	Code          string // 銘柄コード（4桁または5桁）
	Date          string // 基準日付（YYYYMMDD または YYYY-MM-DD）
	From          string // 開始日付（YYYYMMDD または YYYY-MM-DD）
	To            string // 終了日付（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// BreakdownResponse は売買内訳データのレスポンスです。
type BreakdownResponse struct {
	Breakdown     []Breakdown `json:"breakdown"`
	PaginationKey string      `json:"pagination_key"` // ページネーションキー
}

// Breakdown は売買内訳データです。
// 売買代金・売買高を取引種別ごとに細分化したデータを表します。
//
// 注意: このデータはプレミアムプラン専用APIで取得されます。
type Breakdown struct {
	// 基本情報
	Date string `json:"Date"` // 売買日（YYYY-MM-DD形式）
	Code string `json:"Code"` // 銘柄コード

	// 売りの約定代金内訳（単位：円）
	LongSellValue               float64 `json:"LongSellValue"`               // 実売りの約定代金
	ShortSellWithoutMarginValue float64 `json:"ShortSellWithoutMarginValue"` // 空売り（信用新規売りを除く）の約定代金
	MarginSellNewValue          float64 `json:"MarginSellNewValue"`          // 信用新規売りの約定代金
	MarginSellCloseValue        float64 `json:"MarginSellCloseValue"`        // 信用返済売りの約定代金

	// 買いの約定代金内訳（単位：円）
	LongBuyValue        float64 `json:"LongBuyValue"`        // 現物買いの約定代金
	MarginBuyNewValue   float64 `json:"MarginBuyNewValue"`   // 信用新規買いの約定代金
	MarginBuyCloseValue float64 `json:"MarginBuyCloseValue"` // 信用返済買いの約定代金

	// 売りの約定株数内訳（単位：株）
	LongSellVolume               float64 `json:"LongSellVolume"`               // 実売りの約定株数
	ShortSellWithoutMarginVolume float64 `json:"ShortSellWithoutMarginVolume"` // 空売り（信用新規売りを除く）の約定株数
	MarginSellNewVolume          float64 `json:"MarginSellNewVolume"`          // 信用新規売りの約定株数
	MarginSellCloseVolume        float64 `json:"MarginSellCloseVolume"`        // 信用返済売りの約定株数

	// 買いの約定株数内訳（単位：株）
	LongBuyVolume        float64 `json:"LongBuyVolume"`        // 現物買いの約定株数
	MarginBuyNewVolume   float64 `json:"MarginBuyNewVolume"`   // 信用新規買いの約定株数
	MarginBuyCloseVolume float64 `json:"MarginBuyCloseVolume"` // 信用返済買いの約定株数
}

// RawBreakdown is used for unmarshaling JSON response with mixed types
type RawBreakdown struct {
	Date                         string              `json:"Date"`
	Code                         string              `json:"Code"`
	LongSellValue                types.Float64String `json:"LongSellValue"`
	ShortSellWithoutMarginValue  types.Float64String `json:"ShortSellWithoutMarginValue"`
	MarginSellNewValue           types.Float64String `json:"MarginSellNewValue"`
	MarginSellCloseValue         types.Float64String `json:"MarginSellCloseValue"`
	LongBuyValue                 types.Float64String `json:"LongBuyValue"`
	MarginBuyNewValue            types.Float64String `json:"MarginBuyNewValue"`
	MarginBuyCloseValue          types.Float64String `json:"MarginBuyCloseValue"`
	LongSellVolume               types.Float64String `json:"LongSellVolume"`
	ShortSellWithoutMarginVolume types.Float64String `json:"ShortSellWithoutMarginVolume"`
	MarginSellNewVolume          types.Float64String `json:"MarginSellNewVolume"`
	MarginSellCloseVolume        types.Float64String `json:"MarginSellCloseVolume"`
	LongBuyVolume                types.Float64String `json:"LongBuyVolume"`
	MarginBuyNewVolume           types.Float64String `json:"MarginBuyNewVolume"`
	MarginBuyCloseVolume         types.Float64String `json:"MarginBuyCloseVolume"`
}

// UnmarshalJSON implements custom JSON unmarshaling
func (r *BreakdownResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawBreakdown
	type rawResponse struct {
		Breakdown     []RawBreakdown `json:"breakdown"`
		PaginationKey string         `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawBreakdown to Breakdown
	r.Breakdown = make([]Breakdown, len(raw.Breakdown))
	for i, rb := range raw.Breakdown {
		r.Breakdown[i] = Breakdown{
			Date:                         rb.Date,
			Code:                         rb.Code,
			LongSellValue:                float64(rb.LongSellValue),
			ShortSellWithoutMarginValue:  float64(rb.ShortSellWithoutMarginValue),
			MarginSellNewValue:           float64(rb.MarginSellNewValue),
			MarginSellCloseValue:         float64(rb.MarginSellCloseValue),
			LongBuyValue:                 float64(rb.LongBuyValue),
			MarginBuyNewValue:            float64(rb.MarginBuyNewValue),
			MarginBuyCloseValue:          float64(rb.MarginBuyCloseValue),
			LongSellVolume:               float64(rb.LongSellVolume),
			ShortSellWithoutMarginVolume: float64(rb.ShortSellWithoutMarginVolume),
			MarginSellNewVolume:          float64(rb.MarginSellNewVolume),
			MarginSellCloseVolume:        float64(rb.MarginSellCloseVolume),
			LongBuyVolume:                float64(rb.LongBuyVolume),
			MarginBuyNewVolume:           float64(rb.MarginBuyNewVolume),
			MarginBuyCloseVolume:         float64(rb.MarginBuyCloseVolume),
		}
	}

	return nil
}

// GetBreakdown は売買内訳データを取得します。
// codeまたはdateのいずれかが必須です。
// パラメータ:
// - Code: 銘柄コード（例: "7203" または "72030"）
// - Date: 基準日付（例: "20240101" または "2024-01-01"）
// - From/To: 期間指定（例: "20240101" または "2024-01-01"）
// - PaginationKey: ページネーション用キー
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *BreakdownService) GetBreakdown(params BreakdownParams) (*BreakdownResponse, error) {
	path := "/markets/breakdown"

	query := "?"
	if params.Code != "" {
		query += fmt.Sprintf("code=%s&", params.Code)
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

	var resp BreakdownResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get breakdown: %w", err)
	}

	return &resp, nil
}

// GetBreakdownByCode は指定銘柄の過去N日間の売買内訳データを取得します。
// ページネーションを使用して全データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *BreakdownService) GetBreakdownByCode(code string, days int) ([]Breakdown, error) {
	to := time.Now()
	from := to.AddDate(0, 0, -days)

	var allBreakdown []Breakdown
	paginationKey := ""

	for {
		params := BreakdownParams{
			Code:          code,
			From:          from.Format("20060102"),
			To:            to.Format("20060102"),
			PaginationKey: paginationKey,
		}

		resp, err := s.GetBreakdown(params)
		if err != nil {
			return nil, err
		}

		allBreakdown = append(allBreakdown, resp.Breakdown...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allBreakdown, nil
}

// GetBreakdownByDate は指定日の全銘柄の売買内訳データを取得します。
// ページネーションを使用して大量データを分割取得します。
func (s *BreakdownService) GetBreakdownByDate(date string) ([]Breakdown, error) {
	var allBreakdown []Breakdown
	paginationKey := ""

	for {
		params := BreakdownParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetBreakdown(params)
		if err != nil {
			return nil, err
		}

		allBreakdown = append(allBreakdown, resp.Breakdown...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allBreakdown, nil
}

// 売り合計を計算するヘルパーメソッド

// GetTotalSellValue は売りの約定代金合計を返します
func (b *Breakdown) GetTotalSellValue() float64 {
	return b.LongSellValue + b.ShortSellWithoutMarginValue +
		b.MarginSellNewValue + b.MarginSellCloseValue
}

// GetTotalSellVolume は売りの約定株数合計を返します
func (b *Breakdown) GetTotalSellVolume() float64 {
	return b.LongSellVolume + b.ShortSellWithoutMarginVolume +
		b.MarginSellNewVolume + b.MarginSellCloseVolume
}

// 買い合計を計算するヘルパーメソッド

// GetTotalBuyValue は買いの約定代金合計を返します
func (b *Breakdown) GetTotalBuyValue() float64 {
	return b.LongBuyValue + b.MarginBuyNewValue + b.MarginBuyCloseValue
}

// GetTotalBuyVolume は買いの約定株数合計を返します
func (b *Breakdown) GetTotalBuyVolume() float64 {
	return b.LongBuyVolume + b.MarginBuyNewVolume + b.MarginBuyCloseVolume
}

// 信用取引関連のヘルパーメソッド

// GetMarginNewValue は信用新規取引の約定代金合計を返します
func (b *Breakdown) GetMarginNewValue() float64 {
	return b.MarginSellNewValue + b.MarginBuyNewValue
}

// GetMarginCloseValue は信用返済取引の約定代金合計を返します
func (b *Breakdown) GetMarginCloseValue() float64 {
	return b.MarginSellCloseValue + b.MarginBuyCloseValue
}

// GetShortSellRatio は空売り比率を返します（空売り÷全売り）
func (b *Breakdown) GetShortSellRatio() float64 {
	totalSell := b.GetTotalSellValue()
	if totalSell == 0 {
		return 0
	}
	shortSell := b.ShortSellWithoutMarginValue + b.MarginSellNewValue
	return shortSell / totalSell
}
