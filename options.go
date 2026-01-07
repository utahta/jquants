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
	Data          []Option `json:"data"`
	PaginationKey string   `json:"pagination_key"` // ページネーションキー
}

// Option はオプション四本値データを表します。
// J-Quants API /derivatives/bars/daily/options エンドポイントのレスポンスデータ。
//
// 注意: このデータはプレミアムプラン専用APIで取得されます。
type Option struct {
	// 基本情報
	Code         string  `json:"Code"`         // 銘柄コード
	ProdCat      string  `json:"ProdCat"`      // オプション商品区分
	UndSSO       string  `json:"UndSSO"`       // 有価証券オプション対象銘柄（有価証券オプション以外は"-"）
	Date         string  `json:"Date"`         // 取引日（YYYY-MM-DD形式）
	CM           string  `json:"CM"`           // 限月（YYYY-MM形式、日経225miniオプションは週表記）
	Strike       float64 `json:"Strike"`       // 権利行使価格
	PCDiv        string  `json:"PCDiv"`        // プットコール区分（1: プット、2: コール）
	EmMrgnTrgDiv string  `json:"EmMrgnTrgDiv"` // 緊急取引証拠金発動区分（001: 発動時、002: 清算価格算出時）

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
	Theo    *float64 `json:"Theo"`    // 理論価格
	BaseVol *float64 `json:"BaseVol"` // 基準ボラティリティ
	UnderPx *float64 `json:"UnderPx"` // 原証券価格
	IV      *float64 `json:"IV"`      // インプライドボラティリティ
	IR      *float64 `json:"IR"`      // 理論価格計算用金利
	LTD     *string  `json:"LTD"`     // 取引最終年月日（YYYY-MM-DD形式）
	SQD     *string  `json:"SQD"`     // SQ日（YYYY-MM-DD形式）
	CCMFlag *string  `json:"CCMFlag"` // 中心限月フラグ（1:中心限月、0:その他）
}

// RawOption is used for unmarshaling JSON response with mixed types
type RawOption struct {
	Code         string      `json:"Code"`
	ProdCat      string      `json:"ProdCat"`
	UndSSO       string      `json:"UndSSO"`
	Date         string      `json:"Date"`
	CM           string      `json:"CM"`
	Strike       float64     `json:"Strike"`
	PCDiv        string      `json:"PCDiv"`
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
	Theo         interface{} `json:"Theo"`
	BaseVol      interface{} `json:"BaseVol"`
	UnderPx      interface{} `json:"UnderPx"`
	IV           interface{} `json:"IV"`
	IR           interface{} `json:"IR"`
	LTD          interface{} `json:"LTD"`
	SQD          interface{} `json:"SQD"`
	CCMFlag      interface{} `json:"CCMFlag"`
}

// UnmarshalJSON implements custom JSON unmarshaling for OptionsResponse
func (r *OptionsResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawOption
	type rawResponse struct {
		Data          []RawOption `json:"data"`
		PaginationKey string      `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawOption to Option
	r.Data = make([]Option, len(raw.Data))
	for idx, ro := range raw.Data {
		o := Option{
			Code:         ro.Code,
			ProdCat:      ro.ProdCat,
			UndSSO:       ro.UndSSO,
			Date:         ro.Date,
			CM:           ro.CM,
			Strike:       ro.Strike,
			PCDiv:        ro.PCDiv,
			EmMrgnTrgDiv: ro.EmMrgnTrgDiv,
			O:            ro.O,
			H:            ro.H,
			L:            ro.L,
			C:            ro.C,
			EO:           ro.EO,
			EH:           ro.EH,
			EL:           ro.EL,
			EC:           ro.EC,
			AO:           ro.AO,
			AH:           ro.AH,
			AL:           ro.AL,
			AC:           ro.AC,
			MO:           ro.MO,
			MH:           ro.MH,
			ML:           ro.ML,
			MC:           ro.MC,
			Vo:           ro.Vo,
			OI:           ro.OI,
			Va:           ro.Va,
		}

		// Convert optional fields
		if v, ok := parseOptionalFloatOpt(ro.VoOA); ok {
			o.VoOA = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.Settle); ok {
			o.Settle = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.Theo); ok {
			o.Theo = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.BaseVol); ok {
			o.BaseVol = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.UnderPx); ok {
			o.UnderPx = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.IV); ok {
			o.IV = &v
		}
		if v, ok := parseOptionalFloatOpt(ro.IR); ok {
			o.IR = &v
		}
		if v, ok := parseOptionalStringOpt(ro.LTD); ok {
			o.LTD = &v
		}
		if v, ok := parseOptionalStringOpt(ro.SQD); ok {
			o.SQD = &v
		}
		if v, ok := parseOptionalStringOpt(ro.CCMFlag); ok {
			o.CCMFlag = &v
		}

		r.Data[idx] = o
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

	path := "/derivatives/bars/daily/options"

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

		allData = append(allData, resp.Data...)

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

		allData = append(allData, resp.Data...)

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

		allData = append(allData, resp.Data...)

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

		allData = append(allData, resp.Data...)

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
	return o.PCDiv == "2"
}

// IsPut はプットオプションかどうかを判定します。
func (o *Option) IsPut() bool {
	return o.PCDiv == "1"
}

// IsEmergencyMarginTriggered は緊急取引証拠金が発動されたかどうかを判定します。
func (o *Option) IsEmergencyMarginTriggered() bool {
	return o.EmMrgnTrgDiv == "001"
}

// IsCentralContractMonth は中心限月かどうかを判定します。
func (o *Option) IsCentralContractMonth() bool {
	return o.CCMFlag != nil && *o.CCMFlag == "1"
}

// IsSecurityOption は有価証券オプションかどうかを判定します。
func (o *Option) IsSecurityOption() bool {
	return o.UndSSO != "-"
}

// HasNightSession はナイトセッションデータがあるかを判定します。
func (o *Option) HasNightSession() bool {
	// interface{}型のフィールドが空文字列でないかチェック
	if str, ok := o.EO.(string); ok && str == "" {
		return false
	}
	// 0の場合も取引なしと判定
	if val, ok := o.EO.(float64); ok && val == 0 {
		return false
	}
	return true
}

// HasMorningSession は前場データがあるかを判定します。
func (o *Option) HasMorningSession() bool {
	// interface{}型のフィールドが空文字列でないかチェック
	if str, ok := o.MO.(string); ok && str == "" {
		return false
	}
	return true
}

// GetNightSessionOpen はナイトセッション始値を取得します。
func (o *Option) GetNightSessionOpen() *float64 {
	return parseInterfaceToFloat64Opt(o.EO)
}

// GetNightSessionHigh はナイトセッション高値を取得します。
func (o *Option) GetNightSessionHigh() *float64 {
	return parseInterfaceToFloat64Opt(o.EH)
}

// GetNightSessionLow はナイトセッション安値を取得します。
func (o *Option) GetNightSessionLow() *float64 {
	return parseInterfaceToFloat64Opt(o.EL)
}

// GetNightSessionClose はナイトセッション終値を取得します。
func (o *Option) GetNightSessionClose() *float64 {
	return parseInterfaceToFloat64Opt(o.EC)
}

// GetMorningSessionOpen は前場始値を取得します。
func (o *Option) GetMorningSessionOpen() *float64 {
	return parseInterfaceToFloat64Opt(o.MO)
}

// GetMorningSessionHigh は前場高値を取得します。
func (o *Option) GetMorningSessionHigh() *float64 {
	return parseInterfaceToFloat64Opt(o.MH)
}

// GetMorningSessionLow は前場安値を取得します。
func (o *Option) GetMorningSessionLow() *float64 {
	return parseInterfaceToFloat64Opt(o.ML)
}

// GetMorningSessionClose は前場終値を取得します。
func (o *Option) GetMorningSessionClose() *float64 {
	return parseInterfaceToFloat64Opt(o.MC)
}

// IsITM はイン・ザ・マネーかどうかを判定します。
func (o *Option) IsITM() bool {
	if o.UnderPx == nil {
		return false
	}
	if o.IsCall() {
		return *o.UnderPx > o.Strike
	}
	return *o.UnderPx < o.Strike
}

// IsOTM はアウト・オブ・ザ・マネーかどうかを判定します。
func (o *Option) IsOTM() bool {
	if o.UnderPx == nil {
		return false
	}
	if o.IsCall() {
		return *o.UnderPx < o.Strike
	}
	return *o.UnderPx > o.Strike
}

// IsATM はアット・ザ・マネーかどうかを判定します（許容誤差0.1%）。
func (o *Option) IsATM() bool {
	if o.UnderPx == nil {
		return false
	}
	diff := (*o.UnderPx - o.Strike) / o.Strike
	return diff > -0.001 && diff < 0.001
}

// GetMoneyness はマネーネスを計算します（原資産価格/権利行使価格）。
func (o *Option) GetMoneyness() *float64 {
	if o.UnderPx == nil || o.Strike == 0 {
		return nil
	}
	moneyness := *o.UnderPx / o.Strike
	return &moneyness
}

// GetIntrinsicValue は本質的価値を計算します。
func (o *Option) GetIntrinsicValue() float64 {
	if o.UnderPx == nil {
		return 0
	}
	if o.IsCall() {
		diff := *o.UnderPx - o.Strike
		if diff > 0 {
			return diff
		}
		return 0
	}
	// Put option
	diff := o.Strike - *o.UnderPx
	if diff > 0 {
		return diff
	}
	return 0
}

// GetTimeValue は時間価値を計算します。
func (o *Option) GetTimeValue() *float64 {
	if o.Theo == nil {
		return nil
	}
	intrinsic := o.GetIntrinsicValue()
	timeValue := *o.Theo - intrinsic
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
