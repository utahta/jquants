package jquants

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/utahta/jquants/client"
)

// FuturesService は先物四本値データを取得するサービスです。
// 先物に関する、四本値や清算値段、理論価格に関する情報を提供します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では利用できません。
type FuturesService struct {
	client client.HTTPClient
}

// NewFuturesService は新しいFuturesServiceを作成します。
func NewFuturesService(c client.HTTPClient) *FuturesService {
	return &FuturesService{client: c}
}

// FuturesParams は先物四本値のリクエストパラメータです。
type FuturesParams struct {
	Date          string // 取引日（必須）（YYYYMMDD または YYYY-MM-DD）
	Category      string // 商品区分の指定（TOPIXF, NK225F等）
	ContractFlag  string // 中心限月フラグの指定（1: 中心限月のみ）
	PaginationKey string // ページネーションキー
}

// FuturesResponse は先物四本値のレスポンスです。
type FuturesResponse struct {
	Data          []Futures `json:"data"`
	PaginationKey string    `json:"pagination_key"` // ページネーションキー
}

// Futures は先物四本値データを表します。
// J-Quants API /derivatives/bars/daily/futures エンドポイントのレスポンスデータ。
//
// 注意: このデータはプレミアムプラン専用APIで取得されます。
type Futures struct {
	// 基本情報
	Code         string `json:"Code"`         // 銘柄コード
	ProdCat      string `json:"ProdCat"`      // 先物商品区分
	Date         string `json:"Date"`         // 取引日（YYYY-MM-DD形式）
	CM           string `json:"CM"`           // 限月（YYYY-MM形式）
	EmMrgnTrgDiv string `json:"EmMrgnTrgDiv"` // 緊急取引証拠金発動区分（001: 発動時、002: 清算価格算出時）

	// 日通し四本値
	O float64 `json:"O"` // 日通し始値
	H float64 `json:"H"` // 日通し高値
	L float64 `json:"L"` // 日通し安値
	C float64 `json:"C"` // 日通し終値

	// ナイト・セッション四本値（取引開始日初日は空文字）
	EO interface{} `json:"EO"` // ナイト・セッション始値
	EH interface{} `json:"EH"` // ナイト・セッション高値
	EL interface{} `json:"EL"` // ナイト・セッション安値
	EC interface{} `json:"EC"` // ナイト・セッション終値

	// 日中セッション四本値
	AO float64 `json:"AO"` // 日中始値
	AH float64 `json:"AH"` // 日中高値
	AL float64 `json:"AL"` // 日中安値
	AC float64 `json:"AC"` // 日中終値

	// 前場四本値（前後場取引対象銘柄でない場合、空文字）
	MO interface{} `json:"MO"` // 前場始値
	MH interface{} `json:"MH"` // 前場高値
	ML interface{} `json:"ML"` // 前場安値
	MC interface{} `json:"MC"` // 前場終値

	// 取引情報
	Vo float64 `json:"Vo"` // 取引高
	OI float64 `json:"OI"` // 建玉
	Va float64 `json:"Va"` // 取引代金

	// 2016年7月19日以降のみ提供されるフィールド
	VoOA    *float64 `json:"VoOA"`    // 立会内取引高
	Settle  *float64 `json:"Settle"`  // 清算値段
	LTD     *string  `json:"LTD"`     // 取引最終年月日（YYYY-MM-DD形式）
	SQD     *string  `json:"SQD"`     // SQ日（YYYY-MM-DD形式）
	CCMFlag *string  `json:"CCMFlag"` // 中心限月フラグ（1:中心限月、0:その他）
}

// RawFutures is used for unmarshaling JSON response with mixed types
type RawFutures struct {
	Code         string      `json:"Code"`
	ProdCat      string      `json:"ProdCat"`
	Date         string      `json:"Date"`
	CM           string      `json:"CM"`
	EmMrgnTrgDiv string      `json:"EmMrgnTrgDiv"`
	O            float64     `json:"O"`
	H            float64     `json:"H"`
	L            float64     `json:"L"`
	C            float64     `json:"C"`
	EO           interface{} `json:"EO"`
	EH           interface{} `json:"EH"`
	EL           interface{} `json:"EL"`
	EC           interface{} `json:"EC"`
	AO           float64     `json:"AO"`
	AH           float64     `json:"AH"`
	AL           float64     `json:"AL"`
	AC           float64     `json:"AC"`
	MO           interface{} `json:"MO"`
	MH           interface{} `json:"MH"`
	ML           interface{} `json:"ML"`
	MC           interface{} `json:"MC"`
	Vo           float64     `json:"Vo"`
	OI           float64     `json:"OI"`
	Va           float64     `json:"Va"`
	VoOA         interface{} `json:"VoOA"`
	Settle       interface{} `json:"Settle"`
	LTD          interface{} `json:"LTD"`
	SQD          interface{} `json:"SQD"`
	CCMFlag      interface{} `json:"CCMFlag"`
}

// UnmarshalJSON implements custom JSON unmarshaling for FuturesResponse
func (r *FuturesResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawFutures
	type rawResponse struct {
		Data          []RawFutures `json:"data"`
		PaginationKey string       `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawFutures to Futures
	r.Data = make([]Futures, len(raw.Data))
	for idx, rf := range raw.Data {
		f := Futures{
			Code:         rf.Code,
			ProdCat:      rf.ProdCat,
			Date:         rf.Date,
			CM:           rf.CM,
			EmMrgnTrgDiv: rf.EmMrgnTrgDiv,
			O:            rf.O,
			H:            rf.H,
			L:            rf.L,
			C:            rf.C,
			EO:           rf.EO,
			EH:           rf.EH,
			EL:           rf.EL,
			EC:           rf.EC,
			AO:           rf.AO,
			AH:           rf.AH,
			AL:           rf.AL,
			AC:           rf.AC,
			MO:           rf.MO,
			MH:           rf.MH,
			ML:           rf.ML,
			MC:           rf.MC,
			Vo:           rf.Vo,
			OI:           rf.OI,
			Va:           rf.Va,
		}

		// Convert optional fields
		if v, ok := parseOptionalFloat(rf.VoOA); ok {
			f.VoOA = &v
		}
		if v, ok := parseOptionalFloat(rf.Settle); ok {
			f.Settle = &v
		}
		if v, ok := parseOptionalString(rf.LTD); ok {
			f.LTD = &v
		}
		if v, ok := parseOptionalString(rf.SQD); ok {
			f.SQD = &v
		}
		if v, ok := parseOptionalString(rf.CCMFlag); ok {
			f.CCMFlag = &v
		}

		r.Data[idx] = f
	}

	return nil
}

// GetFutures は先物四本値データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *FuturesService) GetFutures(params FuturesParams) (*FuturesResponse, error) {
	// dateは必須パラメータ
	if params.Date == "" {
		return nil, fmt.Errorf("date parameter is required")
	}

	path := "/derivatives/bars/daily/futures"

	query := fmt.Sprintf("?date=%s", params.Date)
	if params.Category != "" {
		query += fmt.Sprintf("&category=%s", params.Category)
	}
	if params.ContractFlag != "" {
		query += fmt.Sprintf("&contract_flag=%s", params.ContractFlag)
	}
	if params.PaginationKey != "" {
		query += fmt.Sprintf("&pagination_key=%s", params.PaginationKey)
	}

	path += query

	var resp FuturesResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get futures: %w", err)
	}

	return &resp, nil
}

// GetFuturesByDate は指定日の全先物データを取得します。
// ページネーションを使用して全データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *FuturesService) GetFuturesByDate(date string) ([]Futures, error) {
	var allData []Futures
	paginationKey := ""

	for {
		params := FuturesParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetFutures(params)
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

// GetFuturesByCategory は指定日・商品カテゴリの先物データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *FuturesService) GetFuturesByCategory(date, category string) ([]Futures, error) {
	var allData []Futures
	paginationKey := ""

	for {
		params := FuturesParams{
			Date:          date,
			Category:      category,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetFutures(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.Data...)

		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetCentralContractMonthFutures は中心限月の先物データのみを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *FuturesService) GetCentralContractMonthFutures(date string) ([]Futures, error) {
	var allData []Futures
	paginationKey := ""

	for {
		params := FuturesParams{
			Date:          date,
			ContractFlag:  "1",
			PaginationKey: paginationKey,
		}

		resp, err := s.GetFutures(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.Data...)

		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// Helper methods for Futures

// IsEmergencyMarginTriggered は緊急取引証拠金が発動されたかどうかを判定します。
func (f *Futures) IsEmergencyMarginTriggered() bool {
	return f.EmMrgnTrgDiv == "001"
}

// IsCentralContractMonth は中心限月かどうかを判定します。
func (f *Futures) IsCentralContractMonth() bool {
	return f.CCMFlag != nil && *f.CCMFlag == "1"
}

// HasNightSession はナイトセッションデータがあるかを判定します。
func (f *Futures) HasNightSession() bool {
	// interface{}型のフィールドが空文字列でないかチェック
	if str, ok := f.EO.(string); ok && str == "" {
		return false
	}
	return true
}

// HasMorningSession は前場データがあるかを判定します。
func (f *Futures) HasMorningSession() bool {
	// interface{}型のフィールドが空文字列でないかチェック
	if str, ok := f.MO.(string); ok && str == "" {
		return false
	}
	return true
}

// GetNightSessionOpen はナイトセッション始値を取得します。
func (f *Futures) GetNightSessionOpen() *float64 {
	return parseInterfaceToFloat64(f.EO)
}

// GetNightSessionHigh はナイトセッション高値を取得します。
func (f *Futures) GetNightSessionHigh() *float64 {
	return parseInterfaceToFloat64(f.EH)
}

// GetNightSessionLow はナイトセッション安値を取得します。
func (f *Futures) GetNightSessionLow() *float64 {
	return parseInterfaceToFloat64(f.EL)
}

// GetNightSessionClose はナイトセッション終値を取得します。
func (f *Futures) GetNightSessionClose() *float64 {
	return parseInterfaceToFloat64(f.EC)
}

// GetMorningSessionOpen は前場始値を取得します。
func (f *Futures) GetMorningSessionOpen() *float64 {
	return parseInterfaceToFloat64(f.MO)
}

// GetMorningSessionHigh は前場高値を取得します。
func (f *Futures) GetMorningSessionHigh() *float64 {
	return parseInterfaceToFloat64(f.MH)
}

// GetMorningSessionLow は前場安値を取得します。
func (f *Futures) GetMorningSessionLow() *float64 {
	return parseInterfaceToFloat64(f.ML)
}

// GetMorningSessionClose は前場終値を取得します。
func (f *Futures) GetMorningSessionClose() *float64 {
	return parseInterfaceToFloat64(f.MC)
}

// GetDayNightGap は日中始値とナイト終値のギャップを計算します。
func (f *Futures) GetDayNightGap() *float64 {
	nightClose := f.GetNightSessionClose()
	if nightClose == nil {
		return nil
	}
	gap := f.AO - *nightClose
	return &gap
}

// GetWholeDayRange は日通しの値幅を計算します。
func (f *Futures) GetWholeDayRange() float64 {
	return f.H - f.L
}

// Helper functions

func parseInterfaceToFloat64(v interface{}) *float64 {
	switch val := v.(type) {
	case float64:
		return &val
	case int:
		f := float64(val)
		return &f
	case string:
		if val == "" {
			return nil
		}
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return &f
		}
	}
	return nil
}

func parseOptionalFloat(v interface{}) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case int:
		return float64(val), true
	case string:
		if val == "" {
			return 0, false
		}
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

func parseOptionalString(v interface{}) (string, bool) {
	if str, ok := v.(string); ok && str != "" {
		return str, true
	}
	return "", false
}
