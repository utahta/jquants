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
	Futures       []Futures `json:"futures"`
	PaginationKey string    `json:"pagination_key"` // ページネーションキー
}

// Futures は先物四本値データを表します。
// J-Quants API /derivatives/futures エンドポイントのレスポンスデータ。
//
// 注意: このデータはプレミアムプラン専用APIで取得されます。
type Futures struct {
	// 基本情報
	Code                           string `json:"Code"`                           // 銘柄コード
	DerivativesProductCategory     string `json:"DerivativesProductCategory"`     // 先物商品区分
	Date                           string `json:"Date"`                           // 取引日（YYYY-MM-DD形式）
	ContractMonth                  string `json:"ContractMonth"`                  // 限月（YYYY-MM形式）
	EmergencyMarginTriggerDivision string `json:"EmergencyMarginTriggerDivision"` // 緊急取引証拠金発動区分（001: 発動時、002: 清算価格算出時）

	// 日通し四本値
	WholeDayOpen  float64 `json:"WholeDayOpen"`  // 日通し始値
	WholeDayHigh  float64 `json:"WholeDayHigh"`  // 日通し高値
	WholeDayLow   float64 `json:"WholeDayLow"`   // 日通し安値
	WholeDayClose float64 `json:"WholeDayClose"` // 日通し終値

	// ナイト・セッション四本値（取引開始日初日は空文字）
	NightSessionOpen  interface{} `json:"NightSessionOpen"`  // ナイト・セッション始値
	NightSessionHigh  interface{} `json:"NightSessionHigh"`  // ナイト・セッション高値
	NightSessionLow   interface{} `json:"NightSessionLow"`   // ナイト・セッション安値
	NightSessionClose interface{} `json:"NightSessionClose"` // ナイト・セッション終値

	// 日中セッション四本値
	DaySessionOpen  float64 `json:"DaySessionOpen"`  // 日中始値
	DaySessionHigh  float64 `json:"DaySessionHigh"`  // 日中高値
	DaySessionLow   float64 `json:"DaySessionLow"`   // 日中安値
	DaySessionClose float64 `json:"DaySessionClose"` // 日中終値

	// 前場四本値（前後場取引対象銘柄でない場合、空文字）
	MorningSessionOpen  interface{} `json:"MorningSessionOpen"`  // 前場始値
	MorningSessionHigh  interface{} `json:"MorningSessionHigh"`  // 前場高値
	MorningSessionLow   interface{} `json:"MorningSessionLow"`   // 前場安値
	MorningSessionClose interface{} `json:"MorningSessionClose"` // 前場終値

	// 取引情報
	Volume        float64 `json:"Volume"`        // 取引高
	OpenInterest  float64 `json:"OpenInterest"`  // 建玉
	TurnoverValue float64 `json:"TurnoverValue"` // 取引代金

	// 2016年7月19日以降のみ提供されるフィールド
	VolumeOnlyAuction        *float64 `json:"Volume(OnlyAuction)"`      // 立会内取引高
	SettlementPrice          *float64 `json:"SettlementPrice"`          // 清算値段
	LastTradingDay           *string  `json:"LastTradingDay"`           // 取引最終年月日（YYYY-MM-DD形式）
	SpecialQuotationDay      *string  `json:"SpecialQuotationDay"`      // SQ日（YYYY-MM-DD形式）
	CentralContractMonthFlag *string  `json:"CentralContractMonthFlag"` // 中心限月フラグ（1:中心限月、0:その他）
}

// RawFutures is used for unmarshaling JSON response with mixed types
type RawFutures struct {
	Code                           string      `json:"Code"`
	DerivativesProductCategory     string      `json:"DerivativesProductCategory"`
	Date                           string      `json:"Date"`
	ContractMonth                  string      `json:"ContractMonth"`
	EmergencyMarginTriggerDivision string      `json:"EmergencyMarginTriggerDivision"`
	WholeDayOpen                   float64     `json:"WholeDayOpen"`
	WholeDayHigh                   float64     `json:"WholeDayHigh"`
	WholeDayLow                    float64     `json:"WholeDayLow"`
	WholeDayClose                  float64     `json:"WholeDayClose"`
	NightSessionOpen               interface{} `json:"NightSessionOpen"`
	NightSessionHigh               interface{} `json:"NightSessionHigh"`
	NightSessionLow                interface{} `json:"NightSessionLow"`
	NightSessionClose              interface{} `json:"NightSessionClose"`
	DaySessionOpen                 float64     `json:"DaySessionOpen"`
	DaySessionHigh                 float64     `json:"DaySessionHigh"`
	DaySessionLow                  float64     `json:"DaySessionLow"`
	DaySessionClose                float64     `json:"DaySessionClose"`
	MorningSessionOpen             interface{} `json:"MorningSessionOpen"`
	MorningSessionHigh             interface{} `json:"MorningSessionHigh"`
	MorningSessionLow              interface{} `json:"MorningSessionLow"`
	MorningSessionClose            interface{} `json:"MorningSessionClose"`
	Volume                         float64     `json:"Volume"`
	OpenInterest                   float64     `json:"OpenInterest"`
	TurnoverValue                  float64     `json:"TurnoverValue"`
	VolumeOnlyAuction              interface{} `json:"Volume(OnlyAuction)"`
	SettlementPrice                interface{} `json:"SettlementPrice"`
	LastTradingDay                 interface{} `json:"LastTradingDay"`
	SpecialQuotationDay            interface{} `json:"SpecialQuotationDay"`
	CentralContractMonthFlag       interface{} `json:"CentralContractMonthFlag"`
}

// UnmarshalJSON implements custom JSON unmarshaling for FuturesResponse
func (r *FuturesResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawFutures
	type rawResponse struct {
		Futures       []RawFutures `json:"futures"`
		PaginationKey string       `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawFutures to Futures
	r.Futures = make([]Futures, len(raw.Futures))
	for idx, rf := range raw.Futures {
		f := Futures{
			Code:                           rf.Code,
			DerivativesProductCategory:     rf.DerivativesProductCategory,
			Date:                           rf.Date,
			ContractMonth:                  rf.ContractMonth,
			EmergencyMarginTriggerDivision: rf.EmergencyMarginTriggerDivision,
			WholeDayOpen:                   rf.WholeDayOpen,
			WholeDayHigh:                   rf.WholeDayHigh,
			WholeDayLow:                    rf.WholeDayLow,
			WholeDayClose:                  rf.WholeDayClose,
			NightSessionOpen:               rf.NightSessionOpen,
			NightSessionHigh:               rf.NightSessionHigh,
			NightSessionLow:                rf.NightSessionLow,
			NightSessionClose:              rf.NightSessionClose,
			DaySessionOpen:                 rf.DaySessionOpen,
			DaySessionHigh:                 rf.DaySessionHigh,
			DaySessionLow:                  rf.DaySessionLow,
			DaySessionClose:                rf.DaySessionClose,
			MorningSessionOpen:             rf.MorningSessionOpen,
			MorningSessionHigh:             rf.MorningSessionHigh,
			MorningSessionLow:              rf.MorningSessionLow,
			MorningSessionClose:            rf.MorningSessionClose,
			Volume:                         rf.Volume,
			OpenInterest:                   rf.OpenInterest,
			TurnoverValue:                  rf.TurnoverValue,
		}

		// Convert optional fields
		if v, ok := parseOptionalFloat(rf.VolumeOnlyAuction); ok {
			f.VolumeOnlyAuction = &v
		}
		if v, ok := parseOptionalFloat(rf.SettlementPrice); ok {
			f.SettlementPrice = &v
		}
		if v, ok := parseOptionalString(rf.LastTradingDay); ok {
			f.LastTradingDay = &v
		}
		if v, ok := parseOptionalString(rf.SpecialQuotationDay); ok {
			f.SpecialQuotationDay = &v
		}
		if v, ok := parseOptionalString(rf.CentralContractMonthFlag); ok {
			f.CentralContractMonthFlag = &v
		}

		r.Futures[idx] = f
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

	path := "/derivatives/futures"

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

		allData = append(allData, resp.Futures...)

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

		allData = append(allData, resp.Futures...)

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

		allData = append(allData, resp.Futures...)

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
	return f.EmergencyMarginTriggerDivision == "001"
}

// IsCentralContractMonth は中心限月かどうかを判定します。
func (f *Futures) IsCentralContractMonth() bool {
	return f.CentralContractMonthFlag != nil && *f.CentralContractMonthFlag == "1"
}

// HasNightSession はナイトセッションデータがあるかを判定します。
func (f *Futures) HasNightSession() bool {
	// interface{}型のフィールドが空文字列でないかチェック
	if str, ok := f.NightSessionOpen.(string); ok && str == "" {
		return false
	}
	return true
}

// HasMorningSession は前場データがあるかを判定します。
func (f *Futures) HasMorningSession() bool {
	// interface{}型のフィールドが空文字列でないかチェック
	if str, ok := f.MorningSessionOpen.(string); ok && str == "" {
		return false
	}
	return true
}

// GetNightSessionOpen はナイトセッション始値を取得します。
func (f *Futures) GetNightSessionOpen() *float64 {
	return parseInterfaceToFloat64(f.NightSessionOpen)
}

// GetNightSessionHigh はナイトセッション高値を取得します。
func (f *Futures) GetNightSessionHigh() *float64 {
	return parseInterfaceToFloat64(f.NightSessionHigh)
}

// GetNightSessionLow はナイトセッション安値を取得します。
func (f *Futures) GetNightSessionLow() *float64 {
	return parseInterfaceToFloat64(f.NightSessionLow)
}

// GetNightSessionClose はナイトセッション終値を取得します。
func (f *Futures) GetNightSessionClose() *float64 {
	return parseInterfaceToFloat64(f.NightSessionClose)
}

// GetMorningSessionOpen は前場始値を取得します。
func (f *Futures) GetMorningSessionOpen() *float64 {
	return parseInterfaceToFloat64(f.MorningSessionOpen)
}

// GetMorningSessionHigh は前場高値を取得します。
func (f *Futures) GetMorningSessionHigh() *float64 {
	return parseInterfaceToFloat64(f.MorningSessionHigh)
}

// GetMorningSessionLow は前場安値を取得します。
func (f *Futures) GetMorningSessionLow() *float64 {
	return parseInterfaceToFloat64(f.MorningSessionLow)
}

// GetMorningSessionClose は前場終値を取得します。
func (f *Futures) GetMorningSessionClose() *float64 {
	return parseInterfaceToFloat64(f.MorningSessionClose)
}

// GetDayNightGap は日中始値とナイト終値のギャップを計算します。
func (f *Futures) GetDayNightGap() *float64 {
	nightClose := f.GetNightSessionClose()
	if nightClose == nil {
		return nil
	}
	gap := f.DaySessionOpen - *nightClose
	return &gap
}

// GetWholeDayRange は日通しの値幅を計算します。
func (f *Futures) GetWholeDayRange() float64 {
	return f.WholeDayHigh - f.WholeDayLow
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
