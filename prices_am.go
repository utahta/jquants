package jquants

import (
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// PricesAMService は前場四本値データを取得するサービスです。
// 午前中（前場）の取引における始値、高値、安値、終値を提供します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では利用できません。
type PricesAMService struct {
	client client.HTTPClient
}

// NewPricesAMService は新しいPricesAMServiceを作成します。
func NewPricesAMService(c client.HTTPClient) *PricesAMService {
	return &PricesAMService{client: c}
}

// PricesAMParams は前場四本値のリクエストパラメータです。
type PricesAMParams struct {
	Code          string // 銘柄コード（4桁または5桁）
	PaginationKey string // ページネーションキー
}

// PricesAMResponse は前場四本値のレスポンスです。
type PricesAMResponse struct {
	Data          []PriceAM `json:"data"`
	PaginationKey string    `json:"pagination_key"` // ページネーションキー
}

// PriceAM は前場四本値データを表します。
// J-Quants API /equities/bars/daily/am エンドポイントのレスポンスデータ。
// 前場終了後に当日の前場データを取得可能（翌日6:00頃まで）。
//
// 注意: このデータはプレミアムプラン専用APIで取得されます。
type PriceAM struct {
	// 基本情報
	Date string `json:"Date"` // 日付（YYYY-MM-DD形式）
	Code string `json:"Code"` // 銘柄コード

	// 前場四本値データ
	MO *float64 `json:"MO"` // 前場始値（前場最初の約定価格）
	MH *float64 `json:"MH"` // 前場高値（前場中の最高約定価格）
	ML *float64 `json:"ML"` // 前場安値（前場中の最低約定価格）
	MC *float64 `json:"MC"` // 前場終値（前場最後の約定価格、前引け）

	// 前場取引情報
	MVo *float64 `json:"MVo"` // 前場売買高（前場中の総売買株数）
	MVa *float64 `json:"MVa"` // 前場取引代金（前場中の総取引代金）
}

// RawPriceAM is used for unmarshaling JSON response with mixed types
type RawPriceAM struct {
	// 基本情報
	Date string `json:"Date"`
	Code string `json:"Code"`

	// 前場四本値データ
	MO *types.Float64String `json:"MO"`
	MH *types.Float64String `json:"MH"`
	ML *types.Float64String `json:"ML"`
	MC *types.Float64String `json:"MC"`

	// 前場取引情報
	MVo *types.Float64String `json:"MVo"`
	MVa *types.Float64String `json:"MVa"`
}

// UnmarshalJSON implements custom JSON unmarshaling for PricesAMResponse
func (r *PricesAMResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawPriceAM
	type rawResponse struct {
		Data          []RawPriceAM `json:"data"`
		PaginationKey string       `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawPriceAM to PriceAM
	r.Data = make([]PriceAM, len(raw.Data))
	for idx, rpa := range raw.Data {
		r.Data[idx] = PriceAM{
			// 基本情報
			Date: rpa.Date,
			Code: rpa.Code,

			// 前場四本値データ
			MO: types.ToFloat64Ptr(rpa.MO),
			MH: types.ToFloat64Ptr(rpa.MH),
			ML: types.ToFloat64Ptr(rpa.ML),
			MC: types.ToFloat64Ptr(rpa.MC),

			// 前場取引情報
			MVo: types.ToFloat64Ptr(rpa.MVo),
			MVa: types.ToFloat64Ptr(rpa.MVa),
		}
	}

	return nil
}

// GetPricesAM は前場四本値データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *PricesAMService) GetPricesAM(params PricesAMParams) (*PricesAMResponse, error) {
	path := "/equities/bars/daily/am"

	query := "?"
	if params.Code != "" {
		query += fmt.Sprintf("code=%s&", params.Code)
	}
	if params.PaginationKey != "" {
		query += fmt.Sprintf("pagination_key=%s&", params.PaginationKey)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp PricesAMResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get prices AM: %w", err)
	}

	return &resp, nil
}

// GetPricesAMByCode は指定銘柄の前場四本値データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *PricesAMService) GetPricesAMByCode(code string) (*PricesAMResponse, error) {
	return s.GetPricesAM(PricesAMParams{Code: code})
}

// GetAllPricesAM は全銘柄の前場四本値データを取得します。
// ページネーションを使用して全データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *PricesAMService) GetAllPricesAM() ([]PriceAM, error) {
	var allData []PriceAM
	paginationKey := ""

	for {
		params := PricesAMParams{
			PaginationKey: paginationKey,
		}

		resp, err := s.GetPricesAM(params)
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

// GetMorningRange は前場の値幅を計算します。
func (pa *PriceAM) GetMorningRange() *float64 {
	if pa.MH == nil || pa.ML == nil {
		return nil
	}

	result := *pa.MH - *pa.ML
	return &result
}

// GetMorningChangeFromOpen は前場の始値からの変動幅を計算します。
func (pa *PriceAM) GetMorningChangeFromOpen() *float64 {
	if pa.MC == nil || pa.MO == nil {
		return nil
	}

	result := *pa.MC - *pa.MO
	return &result
}

// GetMorningChangeRate は前場の始値からの変動率を計算します（パーセント）。
func (pa *PriceAM) GetMorningChangeRate() *float64 {
	if pa.MC == nil || pa.MO == nil || *pa.MO == 0 {
		return nil
	}

	result := ((*pa.MC - *pa.MO) / *pa.MO) * 100
	return &result
}

// HasMorningTrade は前場に取引があったかを判定します。
func (pa *PriceAM) HasMorningTrade() bool {
	return pa.MVo != nil && *pa.MVo > 0
}

// IsActiveTrading は前場に活発な取引があったかを判定します（売買代金1億円以上）。
func (pa *PriceAM) IsActiveTrading() bool {
	return pa.MVa != nil && *pa.MVa >= 100000000
}

// GetAveragePrice は前場の平均約定価格を計算します（概算）。
func (pa *PriceAM) GetAveragePrice() *float64 {
	if pa.MVa == nil || pa.MVo == nil || *pa.MVo == 0 {
		return nil
	}

	result := *pa.MVa / *pa.MVo
	return &result
}

// IsUpperLimit は前場にストップ高だったかを判定します（始値=高値=安値=終値）。
func (pa *PriceAM) IsUpperLimit() bool {
	if pa.MO == nil || pa.MH == nil || pa.ML == nil || pa.MC == nil {
		return false
	}

	return *pa.MO == *pa.MH &&
		*pa.MH == *pa.ML &&
		*pa.ML == *pa.MC
}

// IsLowerLimit は前場にストップ安だったかを判定します（始値=高値=安値=終値）。
func (pa *PriceAM) IsLowerLimit() bool {
	// ストップ高と同じ条件だが、価格水準で区別する必要がある場合は前日終値との比較が必要
	return pa.IsUpperLimit()
}
