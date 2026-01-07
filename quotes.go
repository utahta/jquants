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
	O         *types.Float64String `json:"O"`
	H         *types.Float64String `json:"H"`
	L         *types.Float64String `json:"L"`
	C         *types.Float64String `json:"C"`
	UL        string               `json:"UL"`
	LL        string               `json:"LL"`
	Vo        *types.Float64String `json:"Vo"`
	Va        *types.Float64String `json:"Va"`
	AdjFactor types.Float64String  `json:"AdjFactor"`
	AdjO      *types.Float64String `json:"AdjO"`
	AdjH      *types.Float64String `json:"AdjH"`
	AdjL      *types.Float64String `json:"AdjL"`
	AdjC      *types.Float64String `json:"AdjC"`
	AdjVo     *types.Float64String `json:"AdjVo"`
	// 前場データ
	MO     *types.Float64String `json:"MO"`
	MH     *types.Float64String `json:"MH"`
	ML     *types.Float64String `json:"ML"`
	MC     *types.Float64String `json:"MC"`
	MUL    string               `json:"MUL"`
	MLL    string               `json:"MLL"`
	MVo    *types.Float64String `json:"MVo"`
	MVa    *types.Float64String `json:"MVa"`
	MAdjO  *types.Float64String `json:"MAdjO"`
	MAdjH  *types.Float64String `json:"MAdjH"`
	MAdjL  *types.Float64String `json:"MAdjL"`
	MAdjC  *types.Float64String `json:"MAdjC"`
	MAdjVo *types.Float64String `json:"MAdjVo"`
	// 後場データ
	AO     *types.Float64String `json:"AO"`
	AH     *types.Float64String `json:"AH"`
	AL     *types.Float64String `json:"AL"`
	AC     *types.Float64String `json:"AC"`
	AUL    string               `json:"AUL"`
	ALL    string               `json:"ALL"`
	AVo    *types.Float64String `json:"AVo"`
	AVa    *types.Float64String `json:"AVa"`
	AAdjO  *types.Float64String `json:"AAdjO"`
	AAdjH  *types.Float64String `json:"AAdjH"`
	AAdjL  *types.Float64String `json:"AAdjL"`
	AAdjC  *types.Float64String `json:"AAdjC"`
	AAdjVo *types.Float64String `json:"AAdjVo"`
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
			O:         types.ToFloat64Ptr(rdq.O),
			H:         types.ToFloat64Ptr(rdq.H),
			L:         types.ToFloat64Ptr(rdq.L),
			C:         types.ToFloat64Ptr(rdq.C),
			UL:        rdq.UL,
			LL:        rdq.LL,
			Vo:        types.ToFloat64Ptr(rdq.Vo),
			Va:        types.ToFloat64Ptr(rdq.Va),
			AdjFactor: float64(rdq.AdjFactor),
			AdjO:      types.ToFloat64Ptr(rdq.AdjO),
			AdjH:      types.ToFloat64Ptr(rdq.AdjH),
			AdjL:      types.ToFloat64Ptr(rdq.AdjL),
			AdjC:      types.ToFloat64Ptr(rdq.AdjC),
			AdjVo:     types.ToFloat64Ptr(rdq.AdjVo),
			// 前場データ
			MO:     types.ToFloat64Ptr(rdq.MO),
			MH:     types.ToFloat64Ptr(rdq.MH),
			ML:     types.ToFloat64Ptr(rdq.ML),
			MC:     types.ToFloat64Ptr(rdq.MC),
			MUL:    rdq.MUL,
			MLL:    rdq.MLL,
			MVo:    types.ToFloat64Ptr(rdq.MVo),
			MVa:    types.ToFloat64Ptr(rdq.MVa),
			MAdjO:  types.ToFloat64Ptr(rdq.MAdjO),
			MAdjH:  types.ToFloat64Ptr(rdq.MAdjH),
			MAdjL:  types.ToFloat64Ptr(rdq.MAdjL),
			MAdjC:  types.ToFloat64Ptr(rdq.MAdjC),
			MAdjVo: types.ToFloat64Ptr(rdq.MAdjVo),
			// 後場データ
			AO:     types.ToFloat64Ptr(rdq.AO),
			AH:     types.ToFloat64Ptr(rdq.AH),
			AL:     types.ToFloat64Ptr(rdq.AL),
			AC:     types.ToFloat64Ptr(rdq.AC),
			AUL:    rdq.AUL,
			ALL:    rdq.ALL,
			AVo:    types.ToFloat64Ptr(rdq.AVo),
			AVa:    types.ToFloat64Ptr(rdq.AVa),
			AAdjO:  types.ToFloat64Ptr(rdq.AAdjO),
			AAdjH:  types.ToFloat64Ptr(rdq.AAdjH),
			AAdjL:  types.ToFloat64Ptr(rdq.AAdjL),
			AAdjC:  types.ToFloat64Ptr(rdq.AAdjC),
			AAdjVo: types.ToFloat64Ptr(rdq.AAdjVo),
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

		allQuotes = append(allQuotes, resp.Data...)

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
