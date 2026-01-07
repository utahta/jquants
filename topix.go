package jquants

import (
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// TOPIXService はTOPIX指数四本値データを取得するサービスです。
// TOPIX専用のエンドポイントを使用してTOPIXデータのみを提供します。
type TOPIXService struct {
	client client.HTTPClient
}

// NewTOPIXService は新しいTOPIXServiceを作成します。
func NewTOPIXService(c client.HTTPClient) *TOPIXService {
	return &TOPIXService{client: c}
}

// TOPIXData はTOPIX指数の四本値データを表します。
// J-Quants API /indices/bars/daily/topix エンドポイントのレスポンスデータ。
type TOPIXData struct {
	Date string  `json:"Date"` // 日付（YYYY-MM-DD形式）
	O    float64 `json:"O"`    // 始値
	H    float64 `json:"H"`    // 高値
	L    float64 `json:"L"`    // 安値
	C    float64 `json:"C"`    // 終値
}

// RawTOPIXData is used for unmarshaling JSON response with mixed types
type RawTOPIXData struct {
	Date string              `json:"Date"`
	O    types.Float64String `json:"O"`
	H    types.Float64String `json:"H"`
	L    types.Float64String `json:"L"`
	C    types.Float64String `json:"C"`
}

// TOPIXResponse はTOPIX指数のレスポンスです。
type TOPIXResponse struct {
	Data          []TOPIXData `json:"data"`
	PaginationKey string      `json:"pagination_key"` // ページネーションキー
}

// UnmarshalJSON implements custom JSON unmarshaling
func (t *TOPIXResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawTOPIXData
	type rawResponse struct {
		Data          []RawTOPIXData `json:"data"`
		PaginationKey string         `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	t.PaginationKey = raw.PaginationKey

	// Convert RawTOPIXData to TOPIXData
	t.Data = make([]TOPIXData, len(raw.Data))
	for idx, rt := range raw.Data {
		t.Data[idx] = TOPIXData{
			Date: rt.Date,
			O:    float64(rt.O),
			H:    float64(rt.H),
			L:    float64(rt.L),
			C:    float64(rt.C),
		}
	}

	return nil
}

// TOPIXParams はTOPIX指数のリクエストパラメータです。
type TOPIXParams struct {
	From          string // 開始日（YYYYMMDD または YYYY-MM-DD）
	To            string // 終了日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// GetTOPIXData は指定された条件でTOPIX指数四本値データを取得します。
// パラメータを指定しない場合は全期間のデータが返されます。
// パラメータ:
// - From/To: 期間指定（例: "20240101" または "2024-01-01"）
// - PaginationKey: ページネーション用キー
func (s *TOPIXService) GetTOPIXData(params TOPIXParams) (*TOPIXResponse, error) {
	path := "/indices/bars/daily/topix"

	query := "?"
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

	var resp TOPIXResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get TOPIX data: %w", err)
	}

	return &resp, nil
}

// GetTOPIXByDateRange は指定した期間のTOPIX指数データを取得します。
// ページネーションを使用して全データを取得します。
func (s *TOPIXService) GetTOPIXByDateRange(from, to string) ([]TOPIXData, error) {
	var allTOPIX []TOPIXData
	paginationKey := ""

	for {
		params := TOPIXParams{
			From:          from,
			To:            to,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetTOPIXData(params)
		if err != nil {
			return nil, err
		}

		allTOPIX = append(allTOPIX, resp.Data...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allTOPIX, nil
}

// GetAllTOPIXData は全期間のTOPIX指数データを取得します。
// ページネーションを使用して大量データを分割取得します。
func (s *TOPIXService) GetAllTOPIXData() ([]TOPIXData, error) {
	var allTOPIX []TOPIXData
	paginationKey := ""

	for {
		params := TOPIXParams{
			PaginationKey: paginationKey,
		}

		resp, err := s.GetTOPIXData(params)
		if err != nil {
			return nil, err
		}

		allTOPIX = append(allTOPIX, resp.Data...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allTOPIX, nil
}

// GetLatestTOPIX は最新のTOPIX指数データを取得します。
func (s *TOPIXService) GetLatestTOPIX() (*TOPIXData, error) {
	resp, err := s.GetTOPIXData(TOPIXParams{})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no TOPIX data found")
	}

	// 最新のデータを返す（通常はレスポンスの最初の要素が最新）
	return &resp.Data[0], nil
}
