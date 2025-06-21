package jquants

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// QuotesService は株価データを取得するサービスです。
// 日次の株価情報（始値、高値、安値、終値、出来高）や調整後株価を提供します。
type QuotesService struct {
	client client.HTTPClient
}

func NewQuotesService(c client.HTTPClient) *QuotesService {
	return &QuotesService{
		client: c,
	}
}

// DailyQuote は日次株価データを表します。
// J-Quants API /prices/daily_quotes エンドポイントのレスポンスデータ。
// 売買が成立しなかった日は四本値、取引高、売買代金がnullになります。
type DailyQuote struct {
	// 基本情報
	Date string `json:"Date"` // 日付 (YYYY-MM-DD形式)
	Code string `json:"Code"` // 銘柄コード

	// 日通しデータ（調整前）
	Open          *float64 `json:"Open"`          // 始値（調整前）
	High          *float64 `json:"High"`          // 高値（調整前）
	Low           *float64 `json:"Low"`           // 安値（調整前）
	Close         *float64 `json:"Close"`         // 終値（調整前）
	UpperLimit    string   `json:"UpperLimit"`    // 日通ストップ高フラグ (0:通常, 1:ストップ高)
	LowerLimit    string   `json:"LowerLimit"`    // 日通ストップ安フラグ (0:通常, 1:ストップ安)
	Volume        *float64 `json:"Volume"`        // 取引高（調整前）
	TurnoverValue *float64 `json:"TurnoverValue"` // 取引代金（円）

	// 調整係数と調整後データ
	AdjustmentFactor float64  `json:"AdjustmentFactor"` // 調整係数（分割等）
	AdjustmentOpen   *float64 `json:"AdjustmentOpen"`   // 調整済み始値
	AdjustmentHigh   *float64 `json:"AdjustmentHigh"`   // 調整済み高値
	AdjustmentLow    *float64 `json:"AdjustmentLow"`    // 調整済み安値
	AdjustmentClose  *float64 `json:"AdjustmentClose"`  // 調整済み終値
	AdjustmentVolume *float64 `json:"AdjustmentVolume"` // 調整済み取引高

	// 前場データ（Premiumプランのみ）
	MorningOpen             *float64 `json:"MorningOpen"`             // 前場始値
	MorningHigh             *float64 `json:"MorningHigh"`             // 前場高値
	MorningLow              *float64 `json:"MorningLow"`              // 前場安値
	MorningClose            *float64 `json:"MorningClose"`            // 前場終値
	MorningUpperLimit       string   `json:"MorningUpperLimit"`       // 前場ストップ高フラグ
	MorningLowerLimit       string   `json:"MorningLowerLimit"`       // 前場ストップ安フラグ
	MorningVolume           *float64 `json:"MorningVolume"`           // 前場売買高
	MorningTurnoverValue    *float64 `json:"MorningTurnoverValue"`    // 前場取引代金
	MorningAdjustmentOpen   *float64 `json:"MorningAdjustmentOpen"`   // 調整済み前場始値
	MorningAdjustmentHigh   *float64 `json:"MorningAdjustmentHigh"`   // 調整済み前場高値
	MorningAdjustmentLow    *float64 `json:"MorningAdjustmentLow"`    // 調整済み前場安値
	MorningAdjustmentClose  *float64 `json:"MorningAdjustmentClose"`  // 調整済み前場終値
	MorningAdjustmentVolume *float64 `json:"MorningAdjustmentVolume"` // 調整済み前場売買高

	// 後場データ（Premiumプランのみ）
	AfternoonOpen             *float64 `json:"AfternoonOpen"`             // 後場始値
	AfternoonHigh             *float64 `json:"AfternoonHigh"`             // 後場高値
	AfternoonLow              *float64 `json:"AfternoonLow"`              // 後場安値
	AfternoonClose            *float64 `json:"AfternoonClose"`            // 後場終値
	AfternoonUpperLimit       string   `json:"AfternoonUpperLimit"`       // 後場ストップ高フラグ
	AfternoonLowerLimit       string   `json:"AfternoonLowerLimit"`       // 後場ストップ安フラグ
	AfternoonVolume           *float64 `json:"AfternoonVolume"`           // 後場売買高
	AfternoonTurnoverValue    *float64 `json:"AfternoonTurnoverValue"`    // 後場取引代金
	AfternoonAdjustmentOpen   *float64 `json:"AfternoonAdjustmentOpen"`   // 調整済み後場始値
	AfternoonAdjustmentHigh   *float64 `json:"AfternoonAdjustmentHigh"`   // 調整済み後場高値
	AfternoonAdjustmentLow    *float64 `json:"AfternoonAdjustmentLow"`    // 調整済み後場安値
	AfternoonAdjustmentClose  *float64 `json:"AfternoonAdjustmentClose"`  // 調整済み後場終値
	AfternoonAdjustmentVolume *float64 `json:"AfternoonAdjustmentVolume"` // 調整済み後場売買高
}

// RawDailyQuote is used for unmarshaling JSON response with mixed types
type RawDailyQuote struct {
	Date string `json:"Date"`
	Code string `json:"Code"`
	// 日通しデータ
	Open             *types.Float64String `json:"Open"`
	High             *types.Float64String `json:"High"`
	Low              *types.Float64String `json:"Low"`
	Close            *types.Float64String `json:"Close"`
	UpperLimit       string               `json:"UpperLimit"`
	LowerLimit       string               `json:"LowerLimit"`
	Volume           *types.Float64String `json:"Volume"`
	TurnoverValue    *types.Float64String `json:"TurnoverValue"`
	AdjustmentFactor types.Float64String  `json:"AdjustmentFactor"`
	AdjustmentOpen   *types.Float64String `json:"AdjustmentOpen"`
	AdjustmentHigh   *types.Float64String `json:"AdjustmentHigh"`
	AdjustmentLow    *types.Float64String `json:"AdjustmentLow"`
	AdjustmentClose  *types.Float64String `json:"AdjustmentClose"`
	AdjustmentVolume *types.Float64String `json:"AdjustmentVolume"`
	// 前場データ
	MorningOpen             *types.Float64String `json:"MorningOpen"`
	MorningHigh             *types.Float64String `json:"MorningHigh"`
	MorningLow              *types.Float64String `json:"MorningLow"`
	MorningClose            *types.Float64String `json:"MorningClose"`
	MorningUpperLimit       string               `json:"MorningUpperLimit"`
	MorningLowerLimit       string               `json:"MorningLowerLimit"`
	MorningVolume           *types.Float64String `json:"MorningVolume"`
	MorningTurnoverValue    *types.Float64String `json:"MorningTurnoverValue"`
	MorningAdjustmentOpen   *types.Float64String `json:"MorningAdjustmentOpen"`
	MorningAdjustmentHigh   *types.Float64String `json:"MorningAdjustmentHigh"`
	MorningAdjustmentLow    *types.Float64String `json:"MorningAdjustmentLow"`
	MorningAdjustmentClose  *types.Float64String `json:"MorningAdjustmentClose"`
	MorningAdjustmentVolume *types.Float64String `json:"MorningAdjustmentVolume"`
	// 後場データ
	AfternoonOpen             *types.Float64String `json:"AfternoonOpen"`
	AfternoonHigh             *types.Float64String `json:"AfternoonHigh"`
	AfternoonLow              *types.Float64String `json:"AfternoonLow"`
	AfternoonClose            *types.Float64String `json:"AfternoonClose"`
	AfternoonUpperLimit       string               `json:"AfternoonUpperLimit"`
	AfternoonLowerLimit       string               `json:"AfternoonLowerLimit"`
	AfternoonVolume           *types.Float64String `json:"AfternoonVolume"`
	AfternoonTurnoverValue    *types.Float64String `json:"AfternoonTurnoverValue"`
	AfternoonAdjustmentOpen   *types.Float64String `json:"AfternoonAdjustmentOpen"`
	AfternoonAdjustmentHigh   *types.Float64String `json:"AfternoonAdjustmentHigh"`
	AfternoonAdjustmentLow    *types.Float64String `json:"AfternoonAdjustmentLow"`
	AfternoonAdjustmentClose  *types.Float64String `json:"AfternoonAdjustmentClose"`
	AfternoonAdjustmentVolume *types.Float64String `json:"AfternoonAdjustmentVolume"`
}

type DailyQuotesResponse struct {
	DailyQuotes   []DailyQuote `json:"daily_quotes"`
	PaginationKey string       `json:"pagination_key"` // ページネーションキー
}

// UnmarshalJSON implements custom JSON unmarshaling
func (d *DailyQuotesResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawDailyQuote
	type rawResponse struct {
		DailyQuotes   []RawDailyQuote `json:"daily_quotes"`
		PaginationKey string          `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	d.PaginationKey = raw.PaginationKey

	// Convert RawDailyQuote to DailyQuote
	d.DailyQuotes = make([]DailyQuote, len(raw.DailyQuotes))
	for i, rdq := range raw.DailyQuotes {
		d.DailyQuotes[i] = DailyQuote{
			Date: rdq.Date,
			Code: rdq.Code,
			// 日通しデータ
			Open:             types.ToFloat64Ptr(rdq.Open),
			High:             types.ToFloat64Ptr(rdq.High),
			Low:              types.ToFloat64Ptr(rdq.Low),
			Close:            types.ToFloat64Ptr(rdq.Close),
			UpperLimit:       rdq.UpperLimit,
			LowerLimit:       rdq.LowerLimit,
			Volume:           types.ToFloat64Ptr(rdq.Volume),
			TurnoverValue:    types.ToFloat64Ptr(rdq.TurnoverValue),
			AdjustmentFactor: float64(rdq.AdjustmentFactor),
			AdjustmentOpen:   types.ToFloat64Ptr(rdq.AdjustmentOpen),
			AdjustmentHigh:   types.ToFloat64Ptr(rdq.AdjustmentHigh),
			AdjustmentLow:    types.ToFloat64Ptr(rdq.AdjustmentLow),
			AdjustmentClose:  types.ToFloat64Ptr(rdq.AdjustmentClose),
			AdjustmentVolume: types.ToFloat64Ptr(rdq.AdjustmentVolume),
			// 前場データ
			MorningOpen:             types.ToFloat64Ptr(rdq.MorningOpen),
			MorningHigh:             types.ToFloat64Ptr(rdq.MorningHigh),
			MorningLow:              types.ToFloat64Ptr(rdq.MorningLow),
			MorningClose:            types.ToFloat64Ptr(rdq.MorningClose),
			MorningUpperLimit:       rdq.MorningUpperLimit,
			MorningLowerLimit:       rdq.MorningLowerLimit,
			MorningVolume:           types.ToFloat64Ptr(rdq.MorningVolume),
			MorningTurnoverValue:    types.ToFloat64Ptr(rdq.MorningTurnoverValue),
			MorningAdjustmentOpen:   types.ToFloat64Ptr(rdq.MorningAdjustmentOpen),
			MorningAdjustmentHigh:   types.ToFloat64Ptr(rdq.MorningAdjustmentHigh),
			MorningAdjustmentLow:    types.ToFloat64Ptr(rdq.MorningAdjustmentLow),
			MorningAdjustmentClose:  types.ToFloat64Ptr(rdq.MorningAdjustmentClose),
			MorningAdjustmentVolume: types.ToFloat64Ptr(rdq.MorningAdjustmentVolume),
			// 後場データ
			AfternoonOpen:             types.ToFloat64Ptr(rdq.AfternoonOpen),
			AfternoonHigh:             types.ToFloat64Ptr(rdq.AfternoonHigh),
			AfternoonLow:              types.ToFloat64Ptr(rdq.AfternoonLow),
			AfternoonClose:            types.ToFloat64Ptr(rdq.AfternoonClose),
			AfternoonUpperLimit:       rdq.AfternoonUpperLimit,
			AfternoonLowerLimit:       rdq.AfternoonLowerLimit,
			AfternoonVolume:           types.ToFloat64Ptr(rdq.AfternoonVolume),
			AfternoonTurnoverValue:    types.ToFloat64Ptr(rdq.AfternoonTurnoverValue),
			AfternoonAdjustmentOpen:   types.ToFloat64Ptr(rdq.AfternoonAdjustmentOpen),
			AfternoonAdjustmentHigh:   types.ToFloat64Ptr(rdq.AfternoonAdjustmentHigh),
			AfternoonAdjustmentLow:    types.ToFloat64Ptr(rdq.AfternoonAdjustmentLow),
			AfternoonAdjustmentClose:  types.ToFloat64Ptr(rdq.AfternoonAdjustmentClose),
			AfternoonAdjustmentVolume: types.ToFloat64Ptr(rdq.AfternoonAdjustmentVolume),
		}
	}

	return nil
}

type DailyQuotesParams struct {
	Code          string // 銘柄コード（4桁または5桁）
	Date          string // 基準日付（YYYYMMDD または YYYY-MM-DD）
	From          string // 開始日付（YYYYMMDD または YYYY-MM-DD）
	To            string // 終了日付（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// GetDailyQuotes は指定された条件で日次株価データを取得します。
// codeまたはdateのいずれかが必須です。
// パラメータ:
// - Code: 銘柄コード（例: "7203" または "72030"）
// - Date: 基準日付（例: "20240101" または "2024-01-01"）
// - From/To: 期間指定（例: "20240101" または "2024-01-01"）
// - PaginationKey: ページネーション用キー
func (s *QuotesService) GetDailyQuotes(params DailyQuotesParams) (*DailyQuotesResponse, error) {
	path := "/prices/daily_quotes"

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

	var resp DailyQuotesResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get daily quotes: %w", err)
	}

	return &resp, nil
}

// GetDailyQuotesByCode は指定銘柄の過去N日間の株価データを取得します。
// 例: GetDailyQuotesByCode("7203", 30) でトヨタ自動車の過去30日間のデータを取得
// ページネーションを使用して全データを取得します。
func (s *QuotesService) GetDailyQuotesByCode(code string, days int) ([]DailyQuote, error) {
	to := time.Now()
	from := to.AddDate(0, 0, -days)

	var allQuotes []DailyQuote
	paginationKey := ""

	for {
		params := DailyQuotesParams{
			Code:          code,
			From:          from.Format("20060102"),
			To:            to.Format("20060102"),
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDailyQuotes(params)
		if err != nil {
			return nil, err
		}

		allQuotes = append(allQuotes, resp.DailyQuotes...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allQuotes, nil
}

// GetDailyQuotesByDate は指定日の全銘柄の株価データを取得します。
// ページネーションを使用して大量データを分割取得します。
func (s *QuotesService) GetDailyQuotesByDate(date string) ([]DailyQuote, error) {
	var allQuotes []DailyQuote
	paginationKey := ""

	for {
		params := DailyQuotesParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDailyQuotes(params)
		if err != nil {
			return nil, err
		}

		allQuotes = append(allQuotes, resp.DailyQuotes...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allQuotes, nil
}

// IsStopHigh はストップ高かどうかを判定します
func (q *DailyQuote) IsStopHigh() bool {
	return q.UpperLimit == "1"
}

// IsStopLow はストップ安かどうかを判定します
func (q *DailyQuote) IsStopLow() bool {
	return q.LowerLimit == "1"
}

// IsMorningStopHigh は前場でストップ高かどうかを判定します
func (q *DailyQuote) IsMorningStopHigh() bool {
	return q.MorningUpperLimit == "1"
}

// IsMorningStopLow は前場でストップ安かどうかを判定します
func (q *DailyQuote) IsMorningStopLow() bool {
	return q.MorningLowerLimit == "1"
}

// IsAfternoonStopHigh は後場でストップ高かどうかを判定します
func (q *DailyQuote) IsAfternoonStopHigh() bool {
	return q.AfternoonUpperLimit == "1"
}

// IsAfternoonStopLow は後場でストップ安かどうかを判定します
func (q *DailyQuote) IsAfternoonStopLow() bool {
	return q.AfternoonLowerLimit == "1"
}

// HasMorningData は前場データが存在するかを判定します（Premiumプランのみ）
func (q *DailyQuote) HasMorningData() bool {
	return q.MorningOpen != nil || q.MorningHigh != nil ||
		q.MorningLow != nil || q.MorningClose != nil
}

// HasAfternoonData は後場データが存在するかを判定します（Premiumプランのみ）
func (q *DailyQuote) HasAfternoonData() bool {
	return q.AfternoonOpen != nil || q.AfternoonHigh != nil ||
		q.AfternoonLow != nil || q.AfternoonClose != nil
}
