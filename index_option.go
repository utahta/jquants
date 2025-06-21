package jquants

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// IndexOptionService は日経225オプションデータを取得するサービスです。
type IndexOptionService struct {
	client client.HTTPClient
}

// NewIndexOptionService は新しいIndexOptionServiceを作成します。
func NewIndexOptionService(c client.HTTPClient) *IndexOptionService {
	return &IndexOptionService{client: c}
}

// IndexOption は日経225オプションの四本値データを表します。
// J-Quants API /option/index_option エンドポイントのレスポンスデータ。
type IndexOption struct {
	// 基本情報
	Date                           string  `json:"Date"`                           // 取引日（YYYY-MM-DD形式）
	Code                           string  `json:"Code"`                           // 銘柄コード
	ContractMonth                  string  `json:"ContractMonth"`                  // 限月（YYYY-MM形式）
	StrikePrice                    float64 `json:"StrikePrice"`                    // 権利行使価格
	PutCallDivision                string  `json:"PutCallDivision"`                // プットコール区分（1: プット、2: コール）
	LastTradingDay                 string  `json:"LastTradingDay"`                 // 取引最終年月日（YYYY-MM-DD形式）
	SpecialQuotationDay            string  `json:"SpecialQuotationDay"`            // SQ日（YYYY-MM-DD形式）
	EmergencyMarginTriggerDivision string  `json:"EmergencyMarginTriggerDivision"` // 緊急取引証拠金発動区分（001: 発動時、002: 通常時）

	// 日通し四本値
	WholeDayOpen  float64 `json:"WholeDayOpen"`  // 日通し始値
	WholeDayHigh  float64 `json:"WholeDayHigh"`  // 日通し高値
	WholeDayLow   float64 `json:"WholeDayLow"`   // 日通し安値
	WholeDayClose float64 `json:"WholeDayClose"` // 日通し終値

	// ナイトセッション四本値（取引開始日初日は空文字）
	NightSessionOpen  *float64 `json:"NightSessionOpen"`  // ナイト・セッション始値
	NightSessionHigh  *float64 `json:"NightSessionHigh"`  // ナイト・セッション高値
	NightSessionLow   *float64 `json:"NightSessionLow"`   // ナイト・セッション安値
	NightSessionClose *float64 `json:"NightSessionClose"` // ナイト・セッション終値

	// 日中セッション四本値
	DaySessionOpen  float64 `json:"DaySessionOpen"`  // 日中始値
	DaySessionHigh  float64 `json:"DaySessionHigh"`  // 日中高値
	DaySessionLow   float64 `json:"DaySessionLow"`   // 日中安値
	DaySessionClose float64 `json:"DaySessionClose"` // 日中終値

	// 取引情報
	Volume            float64  `json:"Volume"`              // 取引高
	VolumeOnlyAuction *float64 `json:"Volume(OnlyAuction)"` // 立会内取引高（2016年7月19日以降）
	OpenInterest      float64  `json:"OpenInterest"`        // 建玉
	TurnoverValue     float64  `json:"TurnoverValue"`       // 取引代金

	// 価格・リスク情報（2016年7月19日以降）
	SettlementPrice   *float64 `json:"SettlementPrice"`   // 清算値段
	TheoreticalPrice  *float64 `json:"TheoreticalPrice"`  // 理論価格
	BaseVolatility    *float64 `json:"BaseVolatility"`    // 基準ボラティリティ
	UnderlyingPrice   *float64 `json:"UnderlyingPrice"`   // 原証券価格
	ImpliedVolatility *float64 `json:"ImpliedVolatility"` // インプライドボラティリティ
	InterestRate      *float64 `json:"InterestRate"`      // 理論価格計算用金利
}

// RawIndexOption is used for unmarshaling JSON response with mixed types
type RawIndexOption struct {
	// 基本情報
	Date                           string              `json:"Date"`
	Code                           string              `json:"Code"`
	ContractMonth                  string              `json:"ContractMonth"`
	StrikePrice                    types.Float64String `json:"StrikePrice"`
	PutCallDivision                string              `json:"PutCallDivision"`
	LastTradingDay                 string              `json:"LastTradingDay"`
	SpecialQuotationDay            string              `json:"SpecialQuotationDay"`
	EmergencyMarginTriggerDivision string              `json:"EmergencyMarginTriggerDivision"`

	// 日通し四本値
	WholeDayOpen  types.Float64String `json:"WholeDayOpen"`
	WholeDayHigh  types.Float64String `json:"WholeDayHigh"`
	WholeDayLow   types.Float64String `json:"WholeDayLow"`
	WholeDayClose types.Float64String `json:"WholeDayClose"`

	// ナイトセッション四本値（空文字の可能性があるためString型として受け取る）
	NightSessionOpen  interface{} `json:"NightSessionOpen"`
	NightSessionHigh  interface{} `json:"NightSessionHigh"`
	NightSessionLow   interface{} `json:"NightSessionLow"`
	NightSessionClose interface{} `json:"NightSessionClose"`

	// 日中セッション四本値
	DaySessionOpen  types.Float64String `json:"DaySessionOpen"`
	DaySessionHigh  types.Float64String `json:"DaySessionHigh"`
	DaySessionLow   types.Float64String `json:"DaySessionLow"`
	DaySessionClose types.Float64String `json:"DaySessionClose"`

	// 取引情報
	Volume            types.Float64String `json:"Volume"`
	VolumeOnlyAuction interface{}         `json:"Volume(OnlyAuction)"`
	OpenInterest      types.Float64String `json:"OpenInterest"`
	TurnoverValue     types.Float64String `json:"TurnoverValue"`

	// 価格・リスク情報
	SettlementPrice   interface{} `json:"SettlementPrice"`
	TheoreticalPrice  interface{} `json:"TheoreticalPrice"`
	BaseVolatility    interface{} `json:"BaseVolatility"`
	UnderlyingPrice   interface{} `json:"UnderlyingPrice"`
	ImpliedVolatility interface{} `json:"ImpliedVolatility"`
	InterestRate      interface{} `json:"InterestRate"`
}

// IndexOptionResponse は日経225オプションのレスポンスです。
type IndexOptionResponse struct {
	IndexOptions  []IndexOption `json:"index_option"`
	PaginationKey string        `json:"pagination_key"` // ページネーションキー
}

// IndexOptionParams は日経225オプションのリクエストパラメータです。
type IndexOptionParams struct {
	Date          string // 取引日（YYYYMMDD または YYYY-MM-DD）（必須）
	PaginationKey string // ページネーションキー
}

// プットコール区分定数
const (
	PutCallDivisionPut  = "1" // プット
	PutCallDivisionCall = "2" // コール
)

// 緊急取引証拠金発動区分定数
const (
	EmergencyMarginTriggerDivisionEmergency = "001" // 緊急取引証拠金発動時
	EmergencyMarginTriggerDivisionNormal    = "002" // 清算価格算出時（通常時）
)

// UnmarshalJSON implements custom JSON unmarshaling for IndexOptionResponse
func (r *IndexOptionResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawIndexOption
	type rawResponse struct {
		IndexOptions  []RawIndexOption `json:"index_option"`
		PaginationKey string           `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawIndexOption to IndexOption
	r.IndexOptions = make([]IndexOption, len(raw.IndexOptions))
	for idx, ro := range raw.IndexOptions {
		r.IndexOptions[idx] = IndexOption{
			// 基本情報
			Date:                           ro.Date,
			Code:                           ro.Code,
			ContractMonth:                  ro.ContractMonth,
			StrikePrice:                    float64(ro.StrikePrice),
			PutCallDivision:                ro.PutCallDivision,
			LastTradingDay:                 ro.LastTradingDay,
			SpecialQuotationDay:            ro.SpecialQuotationDay,
			EmergencyMarginTriggerDivision: ro.EmergencyMarginTriggerDivision,

			// 日通し四本値
			WholeDayOpen:  float64(ro.WholeDayOpen),
			WholeDayHigh:  float64(ro.WholeDayHigh),
			WholeDayLow:   float64(ro.WholeDayLow),
			WholeDayClose: float64(ro.WholeDayClose),

			// ナイトセッション四本値（空文字の場合はnil）
			NightSessionOpen:  parseNullableFloat64(ro.NightSessionOpen),
			NightSessionHigh:  parseNullableFloat64(ro.NightSessionHigh),
			NightSessionLow:   parseNullableFloat64(ro.NightSessionLow),
			NightSessionClose: parseNullableFloat64(ro.NightSessionClose),

			// 日中セッション四本値
			DaySessionOpen:  float64(ro.DaySessionOpen),
			DaySessionHigh:  float64(ro.DaySessionHigh),
			DaySessionLow:   float64(ro.DaySessionLow),
			DaySessionClose: float64(ro.DaySessionClose),

			// 取引情報
			Volume:            float64(ro.Volume),
			VolumeOnlyAuction: parseNullableFloat64(ro.VolumeOnlyAuction),
			OpenInterest:      float64(ro.OpenInterest),
			TurnoverValue:     float64(ro.TurnoverValue),

			// 価格・リスク情報
			SettlementPrice:   parseNullableFloat64(ro.SettlementPrice),
			TheoreticalPrice:  parseNullableFloat64(ro.TheoreticalPrice),
			BaseVolatility:    parseNullableFloat64(ro.BaseVolatility),
			UnderlyingPrice:   parseNullableFloat64(ro.UnderlyingPrice),
			ImpliedVolatility: parseNullableFloat64(ro.ImpliedVolatility),
			InterestRate:      parseNullableFloat64(ro.InterestRate),
		}
	}

	return nil
}

// parseNullableFloat64 converts interface{} to *float64, handling empty strings and nil
func parseNullableFloat64(v interface{}) *float64 {
	switch val := v.(type) {
	case float64:
		return &val
	case string:
		if val == "" {
			return nil
		}
		// Try to convert string to float64 using strconv.ParseFloat
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil
		}
		return &f
	default:
		return nil
	}
}

// GetIndexOptions は指定日の日経225オプションデータを取得します。
func (s *IndexOptionService) GetIndexOptions(params IndexOptionParams) (*IndexOptionResponse, error) {
	if params.Date == "" {
		return nil, fmt.Errorf("date parameter is required")
	}

	path := "/option/index_option"

	query := fmt.Sprintf("?date=%s", params.Date)
	if params.PaginationKey != "" {
		query += fmt.Sprintf("&pagination_key=%s", params.PaginationKey)
	}

	path += query

	var resp IndexOptionResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get index options: %w", err)
	}

	return &resp, nil
}

// GetIndexOptionsByDate は指定日の全日経225オプションデータを取得します。
// ページネーションを使用して全データを取得します。
func (s *IndexOptionService) GetIndexOptionsByDate(date string) ([]IndexOption, error) {
	var allOptions []IndexOption
	paginationKey := ""

	for {
		params := IndexOptionParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetIndexOptions(params)
		if err != nil {
			return nil, err
		}

		allOptions = append(allOptions, resp.IndexOptions...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allOptions, nil
}

// GetCallOptions は指定日のコールオプションを取得します。
func (s *IndexOptionService) GetCallOptions(date string) ([]IndexOption, error) {
	options, err := s.GetIndexOptionsByDate(date)
	if err != nil {
		return nil, err
	}

	var callOptions []IndexOption
	for _, option := range options {
		if option.PutCallDivision == PutCallDivisionCall {
			callOptions = append(callOptions, option)
		}
	}

	return callOptions, nil
}

// GetPutOptions は指定日のプットオプションを取得します。
func (s *IndexOptionService) GetPutOptions(date string) ([]IndexOption, error) {
	options, err := s.GetIndexOptionsByDate(date)
	if err != nil {
		return nil, err
	}

	var putOptions []IndexOption
	for _, option := range options {
		if option.PutCallDivision == PutCallDivisionPut {
			putOptions = append(putOptions, option)
		}
	}

	return putOptions, nil
}

// GetOptionChain は指定日のオプションチェーン（全ての権利行使価格）を取得します。
func (s *IndexOptionService) GetOptionChain(date string) ([]IndexOption, error) {
	return s.GetIndexOptionsByDate(date)
}

// IsCall はコールオプションかどうかを判定します。
func (io *IndexOption) IsCall() bool {
	return io.PutCallDivision == PutCallDivisionCall
}

// IsPut はプットオプションかどうかを判定します。
func (io *IndexOption) IsPut() bool {
	return io.PutCallDivision == PutCallDivisionPut
}

// IsEmergencyMarginTriggered は緊急取引証拠金が発動しているかどうかを判定します。
func (io *IndexOption) IsEmergencyMarginTriggered() bool {
	return io.EmergencyMarginTriggerDivision == EmergencyMarginTriggerDivisionEmergency
}

// HasNightSession はナイトセッションデータがあるかどうかを判定します。
func (io *IndexOption) HasNightSession() bool {
	return io.NightSessionOpen != nil
}
