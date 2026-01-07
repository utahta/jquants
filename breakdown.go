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
	Data          []Breakdown `json:"data"`
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
	LongSellVa     float64 `json:"LongSellVa"`     // 実売りの約定代金
	ShrtNoMrgnVa   float64 `json:"ShrtNoMrgnVa"`   // 空売り（信用新規売りを除く）の約定代金
	MrgnSellNewVa  float64 `json:"MrgnSellNewVa"`  // 信用新規売りの約定代金
	MrgnSellCloseVa float64 `json:"MrgnSellCloseVa"` // 信用返済売りの約定代金

	// 買いの約定代金内訳（単位：円）
	LongBuyVa      float64 `json:"LongBuyVa"`      // 現物買いの約定代金
	MrgnBuyNewVa   float64 `json:"MrgnBuyNewVa"`   // 信用新規買いの約定代金
	MrgnBuyCloseVa float64 `json:"MrgnBuyCloseVa"` // 信用返済買いの約定代金

	// 売りの約定株数内訳（単位：株）
	LongSellVo      float64 `json:"LongSellVo"`      // 実売りの約定株数
	ShrtNoMrgnVo    float64 `json:"ShrtNoMrgnVo"`    // 空売り（信用新規売りを除く）の約定株数
	MrgnSellNewVo   float64 `json:"MrgnSellNewVo"`   // 信用新規売りの約定株数
	MrgnSellCloseVo float64 `json:"MrgnSellCloseVo"` // 信用返済売りの約定株数

	// 買いの約定株数内訳（単位：株）
	LongBuyVo      float64 `json:"LongBuyVo"`      // 現物買いの約定株数
	MrgnBuyNewVo   float64 `json:"MrgnBuyNewVo"`   // 信用新規買いの約定株数
	MrgnBuyCloseVo float64 `json:"MrgnBuyCloseVo"` // 信用返済買いの約定株数
}

// RawBreakdown is used for unmarshaling JSON response with mixed types
type RawBreakdown struct {
	Date            string              `json:"Date"`
	Code            string              `json:"Code"`
	LongSellVa      types.Float64String `json:"LongSellVa"`
	ShrtNoMrgnVa    types.Float64String `json:"ShrtNoMrgnVa"`
	MrgnSellNewVa   types.Float64String `json:"MrgnSellNewVa"`
	MrgnSellCloseVa types.Float64String `json:"MrgnSellCloseVa"`
	LongBuyVa       types.Float64String `json:"LongBuyVa"`
	MrgnBuyNewVa    types.Float64String `json:"MrgnBuyNewVa"`
	MrgnBuyCloseVa  types.Float64String `json:"MrgnBuyCloseVa"`
	LongSellVo      types.Float64String `json:"LongSellVo"`
	ShrtNoMrgnVo    types.Float64String `json:"ShrtNoMrgnVo"`
	MrgnSellNewVo   types.Float64String `json:"MrgnSellNewVo"`
	MrgnSellCloseVo types.Float64String `json:"MrgnSellCloseVo"`
	LongBuyVo       types.Float64String `json:"LongBuyVo"`
	MrgnBuyNewVo    types.Float64String `json:"MrgnBuyNewVo"`
	MrgnBuyCloseVo  types.Float64String `json:"MrgnBuyCloseVo"`
}

// UnmarshalJSON implements custom JSON unmarshaling
func (r *BreakdownResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawBreakdown
	type rawResponse struct {
		Data          []RawBreakdown `json:"data"`
		PaginationKey string         `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawBreakdown to Breakdown
	r.Data = make([]Breakdown, len(raw.Data))
	for i, rb := range raw.Data {
		r.Data[i] = Breakdown{
			Date:            rb.Date,
			Code:            rb.Code,
			LongSellVa:      float64(rb.LongSellVa),
			ShrtNoMrgnVa:    float64(rb.ShrtNoMrgnVa),
			MrgnSellNewVa:   float64(rb.MrgnSellNewVa),
			MrgnSellCloseVa: float64(rb.MrgnSellCloseVa),
			LongBuyVa:       float64(rb.LongBuyVa),
			MrgnBuyNewVa:    float64(rb.MrgnBuyNewVa),
			MrgnBuyCloseVa:  float64(rb.MrgnBuyCloseVa),
			LongSellVo:      float64(rb.LongSellVo),
			ShrtNoMrgnVo:    float64(rb.ShrtNoMrgnVo),
			MrgnSellNewVo:   float64(rb.MrgnSellNewVo),
			MrgnSellCloseVo: float64(rb.MrgnSellCloseVo),
			LongBuyVo:       float64(rb.LongBuyVo),
			MrgnBuyNewVo:    float64(rb.MrgnBuyNewVo),
			MrgnBuyCloseVo:  float64(rb.MrgnBuyCloseVo),
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

		allBreakdown = append(allBreakdown, resp.Data...)

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

		allBreakdown = append(allBreakdown, resp.Data...)

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
	return b.LongSellVa + b.ShrtNoMrgnVa +
		b.MrgnSellNewVa + b.MrgnSellCloseVa
}

// GetTotalSellVolume は売りの約定株数合計を返します
func (b *Breakdown) GetTotalSellVolume() float64 {
	return b.LongSellVo + b.ShrtNoMrgnVo +
		b.MrgnSellNewVo + b.MrgnSellCloseVo
}

// 買い合計を計算するヘルパーメソッド

// GetTotalBuyValue は買いの約定代金合計を返します
func (b *Breakdown) GetTotalBuyValue() float64 {
	return b.LongBuyVa + b.MrgnBuyNewVa + b.MrgnBuyCloseVa
}

// GetTotalBuyVolume は買いの約定株数合計を返します
func (b *Breakdown) GetTotalBuyVolume() float64 {
	return b.LongBuyVo + b.MrgnBuyNewVo + b.MrgnBuyCloseVo
}

// 信用取引関連のヘルパーメソッド

// GetMarginNewValue は信用新規取引の約定代金合計を返します
func (b *Breakdown) GetMarginNewValue() float64 {
	return b.MrgnSellNewVa + b.MrgnBuyNewVa
}

// GetMarginCloseValue は信用返済取引の約定代金合計を返します
func (b *Breakdown) GetMarginCloseValue() float64 {
	return b.MrgnSellCloseVa + b.MrgnBuyCloseVa
}

// GetShortSellRatio は空売り比率を返します（空売り÷全売り）
func (b *Breakdown) GetShortSellRatio() float64 {
	totalSell := b.GetTotalSellValue()
	if totalSell == 0 {
		return 0
	}
	shortSell := b.ShrtNoMrgnVa + b.MrgnSellNewVa
	return shortSell / totalSell
}
