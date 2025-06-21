package jquants

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/utahta/jquants/client"
)

// OptionsService はオプション四本値データを取得するサービスです。
// オプションに関する、四本値や清算値段、理論価格に関する情報を提供します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では利用できません。
type OptionsService struct {
	client client.HTTPClient
}

// NewOptionsService は新しいOptionsServiceを作成します。
func NewOptionsService(c client.HTTPClient) *OptionsService {
	return &OptionsService{client: c}
}

// OptionsParams はオプション四本値のリクエストパラメータです。
type OptionsParams struct {
	Date          string // 取引日（必須）（YYYYMMDD または YYYY-MM-DD）
	Category      string // 商品区分の指定（TOPIXE, NK225E等）
	Code          string // 対象有価証券コード（categoryでEQOPを指定した場合に設定）
	ContractFlag  string // 中心限月フラグの指定（1: 中心限月のみ）
	PaginationKey string // ページネーションキー
}

// OptionsResponse はオプション四本値のレスポンスです。
type OptionsResponse struct {
	Options       []Option `json:"options"`
	PaginationKey string   `json:"pagination_key"` // ページネーションキー
}

// Option はオプション四本値データを表します。
// J-Quants API /derivatives/options エンドポイントのレスポンスデータ。
//
// 注意: このデータはプレミアムプラン専用APIで取得されます。
type Option struct {
	// 基本情報
	Code                           string  `json:"Code"`                           // 銘柄コード
	DerivativesProductCategory     string  `json:"DerivativesProductCategory"`     // オプション商品区分
	UnderlyingSSO                  string  `json:"UnderlyingSSO"`                  // 有価証券オプション対象銘柄（有価証券オプション以外は"-"）
	Date                           string  `json:"Date"`                           // 取引日（YYYY-MM-DD形式）
	ContractMonth                  string  `json:"ContractMonth"`                  // 限月（YYYY-MM形式、日経225miniオプションは週表記）
	StrikePrice                    float64 `json:"StrikePrice"`                    // 権利行使価格
	PutCallDivision                string  `json:"PutCallDivision"`                // プットコール区分（1: プット、2: コール）
	EmergencyMarginTriggerDivision string  `json:"EmergencyMarginTriggerDivision"` // 緊急取引証拠金発動区分（001: 発動時、002: 清算価格算出時）

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
	TheoreticalPrice         *float64 `json:"TheoreticalPrice"`         // 理論価格
	BaseVolatility           *float64 `json:"BaseVolatility"`           // 基準ボラティリティ
	UnderlyingPrice          *float64 `json:"UnderlyingPrice"`          // 原証券価格
	ImpliedVolatility        *float64 `json:"ImpliedVolatility"`        // インプライドボラティリティ
	InterestRate             *float64 `json:"InterestRate"`             // 理論価格計算用金利
	LastTradingDay           *string  `json:"LastTradingDay"`           // 取引最終年月日（YYYY-MM-DD形式）
	SpecialQuotationDay      *string  `json:"SpecialQuotationDay"`      // SQ日（YYYY-MM-DD形式）
	CentralContractMonthFlag *string  `json:"CentralContractMonthFlag"` // 中心限月フラグ（1:中心限月、0:その他）
}

// RawOption is used for unmarshaling JSON response with mixed types
type RawOption struct {
	Code                           string      `json:"Code"`
	DerivativesProductCategory     string      `json:"DerivativesProductCategory"`
	UnderlyingSSO                  string      `json:"UnderlyingSSO"`
	Date                           string      `json:"Date"`
	ContractMonth                  string      `json:"ContractMonth"`
	StrikePrice                    float64     `json:"StrikePrice"`
	PutCallDivision                string      `json:"PutCallDivision"`
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
	TheoreticalPrice               interface{} `json:"TheoreticalPrice"`
	BaseVolatility                 interface{} `json:"BaseVolatility"`
	UnderlyingPrice                interface{} `json:"UnderlyingPrice"`
	ImpliedVolatility              interface{} `json:"ImpliedVolatility"`
	InterestRate                   interface{} `json:"InterestRate"`
	LastTradingDay                 interface{} `json:"LastTradingDay"`
	SpecialQuotationDay            interface{} `json:"SpecialQuotationDay"`
	CentralContractMonthFlag       interface{} `json:"CentralContractMonthFlag"`
}

// UnmarshalJSON implements custom JSON unmarshaling for OptionsResponse
func (r *OptionsResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawOption
	type rawResponse struct {
		Options       []RawOption `json:"options"`
		PaginationKey string      `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawOption to Option
	r.Options = make([]Option, len(raw.Options))
	for idx, ro := range raw.Options {
		o := Option{
			Code:                           ro.Code,
			DerivativesProductCategory:     ro.DerivativesProductCategory,
			UnderlyingSSO:                  ro.UnderlyingSSO,
			Date:                           ro.Date,
			ContractMonth:                  ro.ContractMonth,
			StrikePrice:                    ro.StrikePrice,
			PutCallDivision:                ro.PutCallDivision,
			EmergencyMarginTriggerDivision: ro.EmergencyMarginTriggerDivision,
			WholeDayOpen:                   ro.WholeDayOpen,
			WholeDayHigh:                   ro.WholeDayHigh,
			WholeDayLow:                    ro.WholeDayLow,
			WholeDayClose:                  ro.WholeDayClose,
			NightSessionOpen:               ro.NightSessionOpen,
			NightSessionHigh:               ro.NightSessionHigh,
			NightSessionLow:                ro.NightSessionLow,
			NightSessionClose:              ro.NightSessionClose,
			DaySessionOpen:                 ro.DaySessionOpen,
			DaySessionHigh:                 ro.DaySessionHigh,
			DaySessionLow:                  ro.DaySessionLow,
			DaySessionClose:                ro.DaySessionClose,
			MorningSessionOpen:             ro.MorningSessionOpen,
			MorningSessionHigh:             ro.MorningSessionHigh,
			MorningSessionLow:              ro.MorningSessionLow,
			MorningSessionClose:            ro.MorningSessionClose,
			Volume:                         ro.Volume,
			OpenInterest:                   ro.OpenInterest,
			TurnoverValue:                  ro.TurnoverValue,
		}

		// Convert optional fields
		if v, ok := parseOptionalFloatOpt(ro.VolumeOnlyAuction); ok {
			o.VolumeOnlyAuction = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.SettlementPrice); ok {
			o.SettlementPrice = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.TheoreticalPrice); ok {
			o.TheoreticalPrice = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.BaseVolatility); ok {
			o.BaseVolatility = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.UnderlyingPrice); ok {
			o.UnderlyingPrice = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.ImpliedVolatility); ok {
			o.ImpliedVolatility = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.InterestRate); ok {
			o.InterestRate = &v
		}
		if v, ok := parseOptionalStringOpt(ro.LastTradingDay); ok {
			o.LastTradingDay = &v
		}
		if v, ok := parseOptionalStringOpt(ro.SpecialQuotationDay); ok {
			o.SpecialQuotationDay = &v
		}
		if v, ok := parseOptionalStringOpt(ro.CentralContractMonthFlag); ok {
			o.CentralContractMonthFlag = &v
		}

		r.Options[idx] = o
	}

	return nil
}

// GetOptions はオプション四本値データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *OptionsService) GetOptions(params OptionsParams) (*OptionsResponse, error) {
	// dateは必須パラメータ
	if params.Date == "" {
		return nil, fmt.Errorf("date parameter is required")
	}

	path := "/derivatives/options"

	query := fmt.Sprintf("?date=%s", params.Date)
	if params.Category != "" {
		query += fmt.Sprintf("&category=%s", params.Category)
	}
	if params.Code != "" {
		query += fmt.Sprintf("&code=%s", params.Code)
	}
	if params.ContractFlag != "" {
		query += fmt.Sprintf("&contract_flag=%s", params.ContractFlag)
	}
	if params.PaginationKey != "" {
		query += fmt.Sprintf("&pagination_key=%s", params.PaginationKey)
	}

	path += query

	var resp OptionsResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get options: %w", err)
	}

	return &resp, nil
}

// GetOptionsByDate は指定日の全オプションデータを取得します。
// ページネーションを使用して全データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *OptionsService) GetOptionsByDate(date string) ([]Option, error) {
	var allData []Option
	paginationKey := ""

	for {
		params := OptionsParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetOptions(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.Options...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetOptionsByCategory は指定日・商品カテゴリのオプションデータを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *OptionsService) GetOptionsByCategory(date, category string) ([]Option, error) {
	var allData []Option
	paginationKey := ""

	for {
		params := OptionsParams{
			Date:          date,
			Category:      category,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetOptions(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.Options...)

		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetSecurityOptionsByCode は指定日・銘柄の有価証券オプションデータを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *OptionsService) GetSecurityOptionsByCode(date, code string) ([]Option, error) {
	var allData []Option
	paginationKey := ""

	for {
		params := OptionsParams{
			Date:          date,
			Category:      "EQOP",
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetOptions(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.Options...)

		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetCentralContractMonthOptions は中心限月のオプションデータのみを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *OptionsService) GetCentralContractMonthOptions(date string) ([]Option, error) {
	var allData []Option
	paginationKey := ""

	for {
		params := OptionsParams{
			Date:          date,
			ContractFlag:  "1",
			PaginationKey: paginationKey,
		}

		resp, err := s.GetOptions(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.Options...)

		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// Helper methods for Option

// IsCall はコールオプションかどうかを判定します。
func (o *Option) IsCall() bool {
	return o.PutCallDivision == "2"
}

// IsPut はプットオプションかどうかを判定します。
func (o *Option) IsPut() bool {
	return o.PutCallDivision == "1"
}

// IsEmergencyMarginTriggered は緊急取引証拠金が発動されたかどうかを判定します。
func (o *Option) IsEmergencyMarginTriggered() bool {
	return o.EmergencyMarginTriggerDivision == "001"
}

// IsCentralContractMonth は中心限月かどうかを判定します。
func (o *Option) IsCentralContractMonth() bool {
	return o.CentralContractMonthFlag != nil && *o.CentralContractMonthFlag == "1"
}

// IsSecurityOption は有価証券オプションかどうかを判定します。
func (o *Option) IsSecurityOption() bool {
	return o.UnderlyingSSO != "-"
}

// HasNightSession はナイトセッションデータがあるかを判定します。
func (o *Option) HasNightSession() bool {
	// interface{}型のフィールドが空文字列でないかチェック
	if str, ok := o.NightSessionOpen.(string); ok && str == "" {
		return false
	}
	// 0の場合も取引なしと判定
	if val, ok := o.NightSessionOpen.(float64); ok && val == 0 {
		return false
	}
	return true
}

// HasMorningSession は前場データがあるかを判定します。
func (o *Option) HasMorningSession() bool {
	// interface{}型のフィールドが空文字列でないかチェック
	if str, ok := o.MorningSessionOpen.(string); ok && str == "" {
		return false
	}
	return true
}

// GetNightSessionOpen はナイトセッション始値を取得します。
func (o *Option) GetNightSessionOpen() *float64 {
	return parseInterfaceToFloat64Opt(o.NightSessionOpen)
}

// GetNightSessionHigh はナイトセッション高値を取得します。
func (o *Option) GetNightSessionHigh() *float64 {
	return parseInterfaceToFloat64Opt(o.NightSessionHigh)
}

// GetNightSessionLow はナイトセッション安値を取得します。
func (o *Option) GetNightSessionLow() *float64 {
	return parseInterfaceToFloat64Opt(o.NightSessionLow)
}

// GetNightSessionClose はナイトセッション終値を取得します。
func (o *Option) GetNightSessionClose() *float64 {
	return parseInterfaceToFloat64Opt(o.NightSessionClose)
}

// GetMorningSessionOpen は前場始値を取得します。
func (o *Option) GetMorningSessionOpen() *float64 {
	return parseInterfaceToFloat64Opt(o.MorningSessionOpen)
}

// GetMorningSessionHigh は前場高値を取得します。
func (o *Option) GetMorningSessionHigh() *float64 {
	return parseInterfaceToFloat64Opt(o.MorningSessionHigh)
}

// GetMorningSessionLow は前場安値を取得します。
func (o *Option) GetMorningSessionLow() *float64 {
	return parseInterfaceToFloat64Opt(o.MorningSessionLow)
}

// GetMorningSessionClose は前場終値を取得します。
func (o *Option) GetMorningSessionClose() *float64 {
	return parseInterfaceToFloat64Opt(o.MorningSessionClose)
}

// IsITM はイン・ザ・マネーかどうかを判定します。
func (o *Option) IsITM() bool {
	if o.UnderlyingPrice == nil {
		return false
	}
	if o.IsCall() {
		return *o.UnderlyingPrice > o.StrikePrice
	}
	return *o.UnderlyingPrice < o.StrikePrice
}

// IsOTM はアウト・オブ・ザ・マネーかどうかを判定します。
func (o *Option) IsOTM() bool {
	if o.UnderlyingPrice == nil {
		return false
	}
	if o.IsCall() {
		return *o.UnderlyingPrice < o.StrikePrice
	}
	return *o.UnderlyingPrice > o.StrikePrice
}

// IsATM はアット・ザ・マネーかどうかを判定します（許容誤差0.1%）。
func (o *Option) IsATM() bool {
	if o.UnderlyingPrice == nil {
		return false
	}
	diff := (*o.UnderlyingPrice - o.StrikePrice) / o.StrikePrice
	return diff > -0.001 && diff < 0.001
}

// GetMoneyness はマネーネスを計算します（原資産価格/権利行使価格）。
func (o *Option) GetMoneyness() *float64 {
	if o.UnderlyingPrice == nil || o.StrikePrice == 0 {
		return nil
	}
	moneyness := *o.UnderlyingPrice / o.StrikePrice
	return &moneyness
}

// GetIntrinsicValue は本質的価値を計算します。
func (o *Option) GetIntrinsicValue() float64 {
	if o.UnderlyingPrice == nil {
		return 0
	}
	if o.IsCall() {
		diff := *o.UnderlyingPrice - o.StrikePrice
		if diff > 0 {
			return diff
		}
		return 0
	}
	// Put option
	diff := o.StrikePrice - *o.UnderlyingPrice
	if diff > 0 {
		return diff
	}
	return 0
}

// GetTimeValue は時間価値を計算します。
func (o *Option) GetTimeValue() *float64 {
	if o.TheoreticalPrice == nil {
		return nil
	}
	intrinsic := o.GetIntrinsicValue()
	timeValue := *o.TheoreticalPrice - intrinsic
	return &timeValue
}

// Helper functions

func parseInterfaceToFloat64Opt(v interface{}) *float64 {
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

func parseOptionalFloatOpt(v interface{}) (float64, bool) {
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

func parseOptionalStringOpt(v interface{}) (string, bool) {
	if str, ok := v.(string); ok && str != "" {
		return str, true
	}
	return "", false
}
