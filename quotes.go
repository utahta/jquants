package jquants

import (
	"context"
	"encoding/json"
	"fmt"

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
// J-Quants API /equities/bars/daily エンドポイントのレスポンスデータ。
// 売買が成立しなかった日は四本値、取引高、売買代金がnullになります。
type DailyQuote struct {
	// 基本情報
	Date string `json:"Date"` // 日付 (YYYY-MM-DD形式)
	Code string `json:"Code"` // 銘柄コード

	// 日通しデータ（調整前）
	O  *float64 `json:"O"`  // 始値（調整前）
	H  *float64 `json:"H"`  // 高値（調整前）
	L  *float64 `json:"L"`  // 安値（調整前）
	C  *float64 `json:"C"`  // 終値（調整前）
	UL string   `json:"UL"` // 日通ストップ高フラグ (0:通常, 1:ストップ高)
	LL string   `json:"LL"` // 日通ストップ安フラグ (0:通常, 1:ストップ安)
	Vo *float64 `json:"Vo"` // 取引高（調整前）
	Va *float64 `json:"Va"` // 取引代金（円）

	// 調整係数と調整後データ
	AdjFactor float64  `json:"AdjFactor"` // 調整係数（分割等）
	AdjO      *float64 `json:"AdjO"`      // 調整済み始値
	AdjH      *float64 `json:"AdjH"`      // 調整済み高値
	AdjL      *float64 `json:"AdjL"`      // 調整済み安値
	AdjC      *float64 `json:"AdjC"`      // 調整済み終値
	AdjVo     *float64 `json:"AdjVo"`     // 調整済み取引高

	// 前場データ（Premiumプランのみ）
	MO     *float64 `json:"MO"`     // 前場始値
	MH     *float64 `json:"MH"`     // 前場高値
	ML     *float64 `json:"ML"`     // 前場安値
	MC     *float64 `json:"MC"`     // 前場終値
	MUL    string   `json:"MUL"`    // 前場ストップ高フラグ
	MLL    string   `json:"MLL"`    // 前場ストップ安フラグ
	MVo    *float64 `json:"MVo"`    // 前場売買高
	MVa    *float64 `json:"MVa"`    // 前場取引代金
	MAdjO  *float64 `json:"MAdjO"`  // 調整済み前場始値
	MAdjH  *float64 `json:"MAdjH"`  // 調整済み前場高値
	MAdjL  *float64 `json:"MAdjL"`  // 調整済み前場安値
	MAdjC  *float64 `json:"MAdjC"`  // 調整済み前場終値
	MAdjVo *float64 `json:"MAdjVo"` // 調整済み前場売買高

	// 後場データ（Premiumプランのみ）
	AO     *float64 `json:"AO"`     // 後場始値
	AH     *float64 `json:"AH"`     // 後場高値
	AL     *float64 `json:"AL"`     // 後場安値
	AC     *float64 `json:"AC"`     // 後場終値
	AUL    string   `json:"AUL"`    // 後場ストップ高フラグ
	ALL    string   `json:"ALL"`    // 後場ストップ安フラグ
	AVo    *float64 `json:"AVo"`    // 後場売買高
	AVa    *float64 `json:"AVa"`    // 後場取引代金
	AAdjO  *float64 `json:"AAdjO"`  // 調整済み後場始値
	AAdjH  *float64 `json:"AAdjH"`  // 調整済み後場高値
	AAdjL  *float64 `json:"AAdjL"`  // 調整済み後場安値
	AAdjC  *float64 `json:"AAdjC"`  // 調整済み後場終値
	AAdjVo *float64 `json:"AAdjVo"` // 調整済み後場売買高
}

// RawDailyQuote is used for unmarshaling JSON response with mixed types
type RawDailyQuote struct {
	Date string `json:"Date"`
	Code string `json:"Code"`
	// 日通しデータ
	O         types.NullableFloat64 `json:"O"`
	H         types.NullableFloat64 `json:"H"`
	L         types.NullableFloat64 `json:"L"`
	C         types.NullableFloat64 `json:"C"`
	UL        string                `json:"UL"`
	LL        string                `json:"LL"`
	Vo        types.NullableFloat64 `json:"Vo"`
	Va        types.NullableFloat64 `json:"Va"`
	AdjFactor types.NullableFloat64 `json:"AdjFactor"`
	AdjO      types.NullableFloat64 `json:"AdjO"`
	AdjH      types.NullableFloat64 `json:"AdjH"`
	AdjL      types.NullableFloat64 `json:"AdjL"`
	AdjC      types.NullableFloat64 `json:"AdjC"`
	AdjVo     types.NullableFloat64 `json:"AdjVo"`
	// 前場データ
	MO     types.NullableFloat64 `json:"MO"`
	MH     types.NullableFloat64 `json:"MH"`
	ML     types.NullableFloat64 `json:"ML"`
	MC     types.NullableFloat64 `json:"MC"`
	MUL    string                `json:"MUL"`
	MLL    string                `json:"MLL"`
	MVo    types.NullableFloat64 `json:"MVo"`
	MVa    types.NullableFloat64 `json:"MVa"`
	MAdjO  types.NullableFloat64 `json:"MAdjO"`
	MAdjH  types.NullableFloat64 `json:"MAdjH"`
	MAdjL  types.NullableFloat64 `json:"MAdjL"`
	MAdjC  types.NullableFloat64 `json:"MAdjC"`
	MAdjVo types.NullableFloat64 `json:"MAdjVo"`
	// 後場データ
	AO     types.NullableFloat64 `json:"AO"`
	AH     types.NullableFloat64 `json:"AH"`
	AL     types.NullableFloat64 `json:"AL"`
	AC     types.NullableFloat64 `json:"AC"`
	AUL    string                `json:"AUL"`
	ALL    string                `json:"ALL"`
	AVo    types.NullableFloat64 `json:"AVo"`
	AVa    types.NullableFloat64 `json:"AVa"`
	AAdjO  types.NullableFloat64 `json:"AAdjO"`
	AAdjH  types.NullableFloat64 `json:"AAdjH"`
	AAdjL  types.NullableFloat64 `json:"AAdjL"`
	AAdjC  types.NullableFloat64 `json:"AAdjC"`
	AAdjVo types.NullableFloat64 `json:"AAdjVo"`
}

type DailyQuotesResponse struct {
	Data          []DailyQuote `json:"data"`
	PaginationKey string       `json:"pagination_key"` // ページネーションキー
}

// UnmarshalJSON implements custom JSON unmarshaling
func (d *DailyQuotesResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawDailyQuote
	type rawResponse struct {
		Data          []RawDailyQuote `json:"data"`
		PaginationKey string          `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	d.PaginationKey = raw.PaginationKey

	// Convert RawDailyQuote to DailyQuote
	d.Data = make([]DailyQuote, len(raw.Data))
	for i, rdq := range raw.Data {
		d.Data[i] = DailyQuote{
			Date: rdq.Date,
			Code: rdq.Code,
			// 日通しデータ
			O:         rdq.O.Ptr(),
			H:         rdq.H.Ptr(),
			L:         rdq.L.Ptr(),
			C:         rdq.C.Ptr(),
			UL:        rdq.UL,
			LL:        rdq.LL,
			Vo:        rdq.Vo.Ptr(),
			Va:        rdq.Va.Ptr(),
			AdjFactor: rdq.AdjFactor.Or(0),
			AdjO:      rdq.AdjO.Ptr(),
			AdjH:      rdq.AdjH.Ptr(),
			AdjL:      rdq.AdjL.Ptr(),
			AdjC:      rdq.AdjC.Ptr(),
			AdjVo:     rdq.AdjVo.Ptr(),
			// 前場データ
			MO:     rdq.MO.Ptr(),
			MH:     rdq.MH.Ptr(),
			ML:     rdq.ML.Ptr(),
			MC:     rdq.MC.Ptr(),
			MUL:    rdq.MUL,
			MLL:    rdq.MLL,
			MVo:    rdq.MVo.Ptr(),
			MVa:    rdq.MVa.Ptr(),
			MAdjO:  rdq.MAdjO.Ptr(),
			MAdjH:  rdq.MAdjH.Ptr(),
			MAdjL:  rdq.MAdjL.Ptr(),
			MAdjC:  rdq.MAdjC.Ptr(),
			MAdjVo: rdq.MAdjVo.Ptr(),
			// 後場データ
			AO:     rdq.AO.Ptr(),
			AH:     rdq.AH.Ptr(),
			AL:     rdq.AL.Ptr(),
			AC:     rdq.AC.Ptr(),
			AUL:    rdq.AUL,
			ALL:    rdq.ALL,
			AVo:    rdq.AVo.Ptr(),
			AVa:    rdq.AVa.Ptr(),
			AAdjO:  rdq.AAdjO.Ptr(),
			AAdjH:  rdq.AAdjH.Ptr(),
			AAdjL:  rdq.AAdjL.Ptr(),
			AAdjC:  rdq.AAdjC.Ptr(),
			AAdjVo: rdq.AAdjVo.Ptr(),
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
func (s *QuotesService) GetDailyQuotes(ctx context.Context, params DailyQuotesParams) (*DailyQuotesResponse, error) {
	path := "/equities/bars/daily"

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
	if err := s.client.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get daily quotes: %w", err)
	}

	return &resp, nil
}

// GetDailyQuotesByCode は指定銘柄の全期間の株価データを取得します。
// ページネーションを使用して全データを取得します。
func (s *QuotesService) GetDailyQuotesByCode(ctx context.Context, code string) ([]DailyQuote, error) {
	var allQuotes []DailyQuote
	paginationKey := ""

	for {
		params := DailyQuotesParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDailyQuotes(ctx, params)
		if err != nil {
			return nil, err
		}

		allQuotes = append(allQuotes, resp.Data...)

		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allQuotes, nil
}

// GetDailyQuotesByCodeAndDate は指定銘柄の指定日の株価データを取得します。
func (s *QuotesService) GetDailyQuotesByCodeAndDate(ctx context.Context, code, date string) ([]DailyQuote, error) {
	resp, err := s.GetDailyQuotes(ctx, DailyQuotesParams{
		Code: code,
		Date: date,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetDailyQuotesByCodeAndDateRange は指定銘柄の指定期間の株価データを取得します。
// ページネーションを使用して全データを取得します。
func (s *QuotesService) GetDailyQuotesByCodeAndDateRange(ctx context.Context, code, from, to string) ([]DailyQuote, error) {
	var allQuotes []DailyQuote
	paginationKey := ""

	for {
		params := DailyQuotesParams{
			Code:          code,
			From:          from,
			To:            to,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDailyQuotes(ctx, params)
		if err != nil {
			return nil, err
		}

		allQuotes = append(allQuotes, resp.Data...)

		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allQuotes, nil
}

// GetDailyQuotesByDate は指定日の全銘柄の株価データを取得します。
// ページネーションを使用して大量データを分割取得します。
func (s *QuotesService) GetDailyQuotesByDate(ctx context.Context, date string) ([]DailyQuote, error) {
	var allQuotes []DailyQuote
	paginationKey := ""

	for {
		params := DailyQuotesParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDailyQuotes(ctx, params)
		if err != nil {
			return nil, err
		}

		allQuotes = append(allQuotes, resp.Data...)

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
	return q.UL == "1"
}

// IsStopLow はストップ安かどうかを判定します
func (q *DailyQuote) IsStopLow() bool {
	return q.LL == "1"
}

// IsMorningStopHigh は前場でストップ高かどうかを判定します
func (q *DailyQuote) IsMorningStopHigh() bool {
	return q.MUL == "1"
}

// IsMorningStopLow は前場でストップ安かどうかを判定します
func (q *DailyQuote) IsMorningStopLow() bool {
	return q.MLL == "1"
}

// IsAfternoonStopHigh は後場でストップ高かどうかを判定します
func (q *DailyQuote) IsAfternoonStopHigh() bool {
	return q.AUL == "1"
}

// IsAfternoonStopLow は後場でストップ安かどうかを判定します
func (q *DailyQuote) IsAfternoonStopLow() bool {
	return q.ALL == "1"
}

// HasMorningData は前場データが存在するかを判定します（Premiumプランのみ）
func (q *DailyQuote) HasMorningData() bool {
	return q.MO != nil || q.MH != nil || q.ML != nil || q.MC != nil
}

// HasAfternoonData は後場データが存在するかを判定します（Premiumプランのみ）
func (q *DailyQuote) HasAfternoonData() bool {
	return q.AO != nil || q.AH != nil || q.AL != nil || q.AC != nil
}
