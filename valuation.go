package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// ValuationService は日次のバリュエーション指標を取得するサービスです。
type ValuationService struct {
	client client.HTTPClient
}

func NewValuationService(c client.HTTPClient) *ValuationService {
	return &ValuationService{client: c}
}

// Valuation は /equities/valuation の指標を表します。
// 実績は直近12ヶ月、予想は進行期が対象です。算出できない項目はnilです。
type Valuation struct {
	Date   string   `json:"Date"`
	Code   string   `json:"Code"`
	EPS    *float64 `json:"EPS"`    // 1株当たり利益（実績・円）
	FwdEPS *float64 `json:"FwdEPS"` // 1株当たり利益（予想・円）
	BPS    *float64 `json:"BPS"`    // 1株当たり純資産（円）
	ROE    *float64 `json:"ROE"`    // 自己資本利益率（実績・小数。0.1 = 10%）
	FwdROE *float64 `json:"FwdROE"` // 自己資本利益率（予想・小数）
	PER    *float64 `json:"PER"`    // 株価収益率（実績・倍）
	FwdPER *float64 `json:"FwdPER"` // 株価収益率（予想・倍）
	PBR    *float64 `json:"PBR"`    // 株価純資産倍率（倍）
	MktCap *float64 `json:"MktCap"` // 時価総額（百万円）。自己株式を控除した株式数で算出
}

// RawValuation は数値と数値文字列が混在するレスポンスの読み取りに使用します。
type RawValuation struct {
	Date   string                `json:"Date"`
	Code   string                `json:"Code"`
	EPS    types.NullableFloat64 `json:"EPS"`
	FwdEPS types.NullableFloat64 `json:"FwdEPS"`
	BPS    types.NullableFloat64 `json:"BPS"`
	ROE    types.NullableFloat64 `json:"ROE"`
	FwdROE types.NullableFloat64 `json:"FwdROE"`
	PER    types.NullableFloat64 `json:"PER"`
	FwdPER types.NullableFloat64 `json:"FwdPER"`
	PBR    types.NullableFloat64 `json:"PBR"`
	MktCap types.NullableFloat64 `json:"MktCap"`
}

type ValuationResponse struct {
	Data          []Valuation `json:"data"`
	PaginationKey string      `json:"pagination_key"`
}

func (r *ValuationResponse) UnmarshalJSON(data []byte) error {
	var raw struct {
		Data          []RawValuation `json:"data"`
		PaginationKey string         `json:"pagination_key"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.PaginationKey = raw.PaginationKey
	r.Data = make([]Valuation, len(raw.Data))
	for i, v := range raw.Data {
		r.Data[i] = Valuation{
			Date:   v.Date,
			Code:   v.Code,
			EPS:    v.EPS.Ptr(),
			FwdEPS: v.FwdEPS.Ptr(),
			BPS:    v.BPS.Ptr(),
			ROE:    v.ROE.Ptr(),
			FwdROE: v.FwdROE.Ptr(),
			PER:    v.PER.Ptr(),
			FwdPER: v.FwdPER.Ptr(),
			PBR:    v.PBR.Ptr(),
			MktCap: v.MktCap.Ptr(),
		}
	}
	return nil
}

// ValuationParams はCodeまたはDateの指定が必要です。
type ValuationParams struct {
	Code          string // 銘柄コード（4桁または5桁）
	Date          string // 日付（YYYYMMDD または YYYY-MM-DD）
	From          string // 開始日付（Code指定時）
	To            string // 終了日付（Code指定時）
	PaginationKey string
}

func (s *ValuationService) GetValuations(ctx context.Context, params ValuationParams) (*ValuationResponse, error) {
	if params.Code == "" && params.Date == "" {
		return nil, fmt.Errorf("either code or date parameter is required")
	}

	query := url.Values{}
	if params.Code != "" {
		query.Set("code", params.Code)
	}
	if params.Date != "" {
		query.Set("date", params.Date)
	}
	if params.From != "" {
		query.Set("from", params.From)
	}
	if params.To != "" {
		query.Set("to", params.To)
	}
	if params.PaginationKey != "" {
		query.Set("pagination_key", params.PaginationKey)
	}

	var resp ValuationResponse
	if err := s.client.DoRequest(ctx, "GET", "/equities/valuation?"+query.Encode(), nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get valuations: %w", err)
	}
	return &resp, nil
}

// GetValuationsByCode は指定銘柄の全期間のデータを取得します。
func (s *ValuationService) GetValuationsByCode(ctx context.Context, code string) ([]Valuation, error) {
	return s.getAllValuations(ctx, ValuationParams{Code: code})
}

// GetValuationsByDate は指定日の全銘柄のデータを取得します。
func (s *ValuationService) GetValuationsByDate(ctx context.Context, date string) ([]Valuation, error) {
	return s.getAllValuations(ctx, ValuationParams{Date: date})
}

func (s *ValuationService) GetValuationsByCodeAndDate(ctx context.Context, code, date string) ([]Valuation, error) {
	return s.getAllValuations(ctx, ValuationParams{Code: code, Date: date})
}

func (s *ValuationService) GetValuationsByCodeAndDateRange(ctx context.Context, code, from, to string) ([]Valuation, error) {
	return s.getAllValuations(ctx, ValuationParams{Code: code, From: from, To: to})
}

func (s *ValuationService) getAllValuations(ctx context.Context, params ValuationParams) ([]Valuation, error) {
	var allData []Valuation
	for {
		resp, err := s.GetValuations(ctx, params)
		if err != nil {
			return nil, err
		}
		allData = append(allData, resp.Data...)
		if resp.PaginationKey == "" {
			return allData, nil
		}
		params.PaginationKey = resp.PaginationKey
	}
}
