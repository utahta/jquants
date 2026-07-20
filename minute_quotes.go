package jquants

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// MinuteQuotesService は株価分足データを取得するサービスです。
// 歩み値を1分単位に集計した四本値（始値・高値・安値・終値）、出来高、売買代金を提供します。
// 利用にはアドオン契約が必要です。データ取得可能期間は過去2年間です。
// 約定のない時間帯の分足は返却されません。地方取引所単独上場銘柄は収録対象外です。
type MinuteQuotesService struct {
	client client.HTTPClient
}

// NewMinuteQuotesService は新しいMinuteQuotesServiceを作成します。
func NewMinuteQuotesService(c client.HTTPClient) *MinuteQuotesService {
	return &MinuteQuotesService{client: c}
}

// MinuteQuotesParams は株価分足のリクエストパラメータです。
type MinuteQuotesParams struct {
	Code          string // 銘柄コード（4桁または5桁）（code、dateのいずれかが必須）
	Date          string // 日付（YYYYMMDD または YYYY-MM-DD）（code、dateのいずれかが必須）
	From          string // fromの指定（YYYYMMDD または YYYY-MM-DD）
	To            string // toの指定（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// MinuteQuotesResponse は株価分足のレスポンスです。
type MinuteQuotesResponse struct {
	Data          []MinuteQuote `json:"data"`
	PaginationKey string        `json:"pagination_key"` // ページネーションキー
}

// MinuteQuote は株価分足データを表します。
// J-Quants API /equities/bars/minute エンドポイントのレスポンスデータ。
// 約定のない時間帯の分足は返却されないため、四本値・出来高・売買代金は常に値を持ちます。
type MinuteQuote struct {
	// 基本情報
	Date string `json:"Date"` // 日付（YYYY-MM-DD形式）
	Time string `json:"Time"` // 時刻（HH:mm形式）
	Code string `json:"Code"` // 銘柄コード

	// 四本値
	O float64 `json:"O"` // 始値
	H float64 `json:"H"` // 高値
	L float64 `json:"L"` // 安値
	C float64 `json:"C"` // 終値

	// 出来高・売買代金
	Vo float64 `json:"Vo"` // 出来高
	Va float64 `json:"Va"` // 売買代金
}

// RawMinuteQuote is used for unmarshaling JSON response with mixed types
type RawMinuteQuote struct {
	// 基本情報
	Date string `json:"Date"`
	Time string `json:"Time"`
	Code string `json:"Code"`

	// 四本値
	O types.NullableFloat64 `json:"O"`
	H types.NullableFloat64 `json:"H"`
	L types.NullableFloat64 `json:"L"`
	C types.NullableFloat64 `json:"C"`

	// 出来高・売買代金
	Vo types.NullableFloat64 `json:"Vo"`
	Va types.NullableFloat64 `json:"Va"`
}

// UnmarshalJSON implements custom JSON unmarshaling for MinuteQuotesResponse
func (r *MinuteQuotesResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawMinuteQuote
	type rawResponse struct {
		Data          []RawMinuteQuote `json:"data"`
		PaginationKey string           `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawMinuteQuote to MinuteQuote
	r.Data = make([]MinuteQuote, len(raw.Data))
	for idx, rmq := range raw.Data {
		r.Data[idx] = MinuteQuote{
			// 基本情報
			Date: rmq.Date,
			Time: rmq.Time,
			Code: rmq.Code,

			// 四本値
			O: rmq.O.Or(0),
			H: rmq.H.Or(0),
			L: rmq.L.Or(0),
			C: rmq.C.Or(0),

			// 出来高・売買代金
			Vo: rmq.Vo.Or(0),
			Va: rmq.Va.Or(0),
		}
	}

	return nil
}

// GetMinuteQuotes は株価分足データを取得します。
func (s *MinuteQuotesService) GetMinuteQuotes(ctx context.Context, params MinuteQuotesParams) (*MinuteQuotesResponse, error) {
	// code、dateのいずれかが必須
	if params.Code == "" && params.Date == "" {
		return nil, fmt.Errorf("either code or date parameter is required")
	}

	path := "/equities/bars/minute"

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

	var resp MinuteQuotesResponse
	if err := s.client.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get minute quotes: %w", err)
	}

	return &resp, nil
}

// GetMinuteQuotesByCode は指定銘柄の株価分足データを取得します。
// ページネーションを使用して全データを取得します。
func (s *MinuteQuotesService) GetMinuteQuotesByCode(ctx context.Context, code string) ([]MinuteQuote, error) {
	var allData []MinuteQuote
	paginationKey := ""

	for {
		params := MinuteQuotesParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetMinuteQuotes(ctx, params)
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

// GetMinuteQuotesByCodeAndDate は指定銘柄の指定日の株価分足データを取得します。
// ページネーションを使用して全データを取得します。
func (s *MinuteQuotesService) GetMinuteQuotesByCodeAndDate(ctx context.Context, code, date string) ([]MinuteQuote, error) {
	var allData []MinuteQuote
	paginationKey := ""

	for {
		params := MinuteQuotesParams{
			Code:          code,
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetMinuteQuotes(ctx, params)
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

// GetMinuteQuotesByDate は指定日の全上場銘柄の株価分足データを取得します。
// ページネーションを使用して全データを取得します。
func (s *MinuteQuotesService) GetMinuteQuotesByDate(ctx context.Context, date string) ([]MinuteQuote, error) {
	var allData []MinuteQuote
	paginationKey := ""

	for {
		params := MinuteQuotesParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetMinuteQuotes(ctx, params)
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
