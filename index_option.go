package jquants

import (
	"context"
	"encoding/json"
	"fmt"

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
// J-Quants API /derivatives/bars/daily/options/225 エンドポイントのレスポンスデータ。
type IndexOption struct {
	// 基本情報
	Date         string  `json:"Date"`         // 取引日（YYYY-MM-DD形式）
	Code         string  `json:"Code"`         // 銘柄コード
	CM           string  `json:"CM"`           // 限月（YYYY-MM形式）
	Strike       float64 `json:"Strike"`       // 権利行使価格
	PCDiv        string  `json:"PCDiv"`        // プットコール区分（1: プット、2: コール）
	LTD          string  `json:"LTD"`          // 取引最終年月日（YYYY-MM-DD形式）
	SQD          string  `json:"SQD"`          // SQ日（YYYY-MM-DD形式）
	EmMrgnTrgDiv string  `json:"EmMrgnTrgDiv"` // 緊急取引証拠金発動区分（001: 発動時、002: 通常時）

	// 日通し四本値
	O float64 `json:"O"` // 日通し始値
	H float64 `json:"H"` // 日通し高値
	L float64 `json:"L"` // 日通し安値
	C float64 `json:"C"` // 日通し終値

	// ナイトセッション四本値（取引開始日初日は空文字）
	EO *float64 `json:"EO"` // ナイト・セッション始値
	EH *float64 `json:"EH"` // ナイト・セッション高値
	EL *float64 `json:"EL"` // ナイト・セッション安値
	EC *float64 `json:"EC"` // ナイト・セッション終値

	// 日中セッション四本値
	AO float64 `json:"AO"` // 日中始値
	AH float64 `json:"AH"` // 日中高値
	AL float64 `json:"AL"` // 日中安値
	AC float64 `json:"AC"` // 日中終値

	// 取引情報
	Vo   float64  `json:"Vo"`   // 取引高
	VoOA *float64 `json:"VoOA"` // 立会内取引高（2016年7月19日以降）
	OI   float64  `json:"OI"`   // 建玉
	Va   float64  `json:"Va"`   // 取引代金

	// 価格・リスク情報（2016年7月19日以降）
	Settle  *float64 `json:"Settle"`  // 清算値段
	Theo    *float64 `json:"Theo"`    // 理論価格
	BaseVol *float64 `json:"BaseVol"` // 基準ボラティリティ
	UnderPx *float64 `json:"UnderPx"` // 原証券価格
	IV      *float64 `json:"IV"`      // インプライドボラティリティ
	IR      *float64 `json:"IR"`      // 理論価格計算用金利
}

// RawIndexOption is used for unmarshaling JSON response with mixed types
type RawIndexOption struct {
	// 基本情報
	Date         string                `json:"Date"`
	Code         string                `json:"Code"`
	CM           string                `json:"CM"`
	Strike       types.NullableFloat64 `json:"Strike"`
	PCDiv        string                `json:"PCDiv"`
	LTD          string                `json:"LTD"`
	SQD          string                `json:"SQD"`
	EmMrgnTrgDiv string                `json:"EmMrgnTrgDiv"`

	// 日通し四本値
	O types.NullableFloat64 `json:"O"`
	H types.NullableFloat64 `json:"H"`
	L types.NullableFloat64 `json:"L"`
	C types.NullableFloat64 `json:"C"`

	// ナイトセッション四本値（空文字の可能性があるためinterface{}型として受け取る）
	EO types.NullableFloat64 `json:"EO"`
	EH types.NullableFloat64 `json:"EH"`
	EL types.NullableFloat64 `json:"EL"`
	EC types.NullableFloat64 `json:"EC"`

	// 日中セッション四本値
	AO types.NullableFloat64 `json:"AO"`
	AH types.NullableFloat64 `json:"AH"`
	AL types.NullableFloat64 `json:"AL"`
	AC types.NullableFloat64 `json:"AC"`

	// 取引情報
	Vo   types.NullableFloat64 `json:"Vo"`
	VoOA types.NullableFloat64 `json:"VoOA"`
	OI   types.NullableFloat64 `json:"OI"`
	Va   types.NullableFloat64 `json:"Va"`

	// 価格・リスク情報
	Settle  types.NullableFloat64 `json:"Settle"`
	Theo    types.NullableFloat64 `json:"Theo"`
	BaseVol types.NullableFloat64 `json:"BaseVol"`
	UnderPx types.NullableFloat64 `json:"UnderPx"`
	IV      types.NullableFloat64 `json:"IV"`
	IR      types.NullableFloat64 `json:"IR"`
}

// IndexOptionResponse は日経225オプションのレスポンスです。
type IndexOptionResponse struct {
	Data          []IndexOption `json:"data"`
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
		Data          []RawIndexOption `json:"data"`
		PaginationKey string           `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawIndexOption to IndexOption
	r.Data = make([]IndexOption, len(raw.Data))
	for idx, ro := range raw.Data {
		r.Data[idx] = IndexOption{
			// 基本情報
			Date:         ro.Date,
			Code:         ro.Code,
			CM:           ro.CM,
			Strike:       ro.Strike.Or(0),
			PCDiv:        ro.PCDiv,
			LTD:          ro.LTD,
			SQD:          ro.SQD,
			EmMrgnTrgDiv: ro.EmMrgnTrgDiv,

			// 日通し四本値
			O: ro.O.Or(0),
			H: ro.H.Or(0),
			L: ro.L.Or(0),
			C: ro.C.Or(0),

			// ナイトセッション四本値（空文字の場合はnil）
			EO: ro.EO.Ptr(),
			EH: ro.EH.Ptr(),
			EL: ro.EL.Ptr(),
			EC: ro.EC.Ptr(),

			// 日中セッション四本値
			AO: ro.AO.Or(0),
			AH: ro.AH.Or(0),
			AL: ro.AL.Or(0),
			AC: ro.AC.Or(0),

			// 取引情報
			Vo:   ro.Vo.Or(0),
			VoOA: ro.VoOA.Ptr(),
			OI:   ro.OI.Or(0),
			Va:   ro.Va.Or(0),

			// 価格・リスク情報
			Settle:  ro.Settle.Ptr(),
			Theo:    ro.Theo.Ptr(),
			BaseVol: ro.BaseVol.Ptr(),
			UnderPx: ro.UnderPx.Ptr(),
			IV:      ro.IV.Ptr(),
			IR:      ro.IR.Ptr(),
		}
	}

	return nil
}

// GetIndexOptions は指定日の日経225オプションデータを取得します。
func (s *IndexOptionService) GetIndexOptions(ctx context.Context, params IndexOptionParams) (*IndexOptionResponse, error) {
	if params.Date == "" {
		return nil, fmt.Errorf("date parameter is required")
	}

	path := "/derivatives/bars/daily/options/225"

	query := fmt.Sprintf("?date=%s", params.Date)
	if params.PaginationKey != "" {
		query += fmt.Sprintf("&pagination_key=%s", params.PaginationKey)
	}

	path += query

	var resp IndexOptionResponse
	if err := s.client.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get index options: %w", err)
	}

	return &resp, nil
}

// GetIndexOptionsByDate は指定日の全日経225オプションデータを取得します。
// ページネーションを使用して全データを取得します。
func (s *IndexOptionService) GetIndexOptionsByDate(ctx context.Context, date string) ([]IndexOption, error) {
	var allOptions []IndexOption
	paginationKey := ""

	for {
		params := IndexOptionParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetIndexOptions(ctx, params)
		if err != nil {
			return nil, err
		}

		allOptions = append(allOptions, resp.Data...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allOptions, nil
}

// GetCallOptions は指定日のコールオプションを取得します。
func (s *IndexOptionService) GetCallOptions(ctx context.Context, date string) ([]IndexOption, error) {
	options, err := s.GetIndexOptionsByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	var callOptions []IndexOption
	for _, option := range options {
		if option.PCDiv == PutCallDivisionCall {
			callOptions = append(callOptions, option)
		}
	}

	return callOptions, nil
}

// GetPutOptions は指定日のプットオプションを取得します。
func (s *IndexOptionService) GetPutOptions(ctx context.Context, date string) ([]IndexOption, error) {
	options, err := s.GetIndexOptionsByDate(ctx, date)
	if err != nil {
		return nil, err
	}

	var putOptions []IndexOption
	for _, option := range options {
		if option.PCDiv == PutCallDivisionPut {
			putOptions = append(putOptions, option)
		}
	}

	return putOptions, nil
}

// GetOptionChain は指定日のオプションチェーン（全ての権利行使価格）を取得します。
func (s *IndexOptionService) GetOptionChain(ctx context.Context, date string) ([]IndexOption, error) {
	return s.GetIndexOptionsByDate(ctx, date)
}

// IsCall はコールオプションかどうかを判定します。
func (io *IndexOption) IsCall() bool {
	return io.PCDiv == PutCallDivisionCall
}

// IsPut はプットオプションかどうかを判定します。
func (io *IndexOption) IsPut() bool {
	return io.PCDiv == PutCallDivisionPut
}

// IsEmergencyMarginTriggered は緊急取引証拠金が発動しているかどうかを判定します。
func (io *IndexOption) IsEmergencyMarginTriggered() bool {
	return io.EmMrgnTrgDiv == EmergencyMarginTriggerDivisionEmergency
}

// HasNightSession はナイトセッションデータがあるかどうかを判定します。
func (io *IndexOption) HasNightSession() bool {
	return io.EO != nil
}
