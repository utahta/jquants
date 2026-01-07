package jquants

import (
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// ShortSellingPositionsService は空売り残高報告を取得するサービスです。
type ShortSellingPositionsService struct {
	client client.HTTPClient
}

// NewShortSellingPositionsService は新しいShortSellingPositionsServiceを作成します。
func NewShortSellingPositionsService(c client.HTTPClient) *ShortSellingPositionsService {
	return &ShortSellingPositionsService{client: c}
}

// ShortSellingPositionsParams は空売り残高報告のリクエストパラメータです。
type ShortSellingPositionsParams struct {
	Code              string // 4桁もしくは5桁の銘柄コード（code、disclosed_date、calculated_dateのいずれかが必須）
	DisclosedDate     string // 公表日（YYYYMMDD または YYYY-MM-DD）（code、disclosed_date、calculated_dateのいずれかが必須）
	DisclosedDateFrom string // 公表日のfrom指定（YYYYMMDD または YYYY-MM-DD）
	DisclosedDateTo   string // 公表日のto指定（YYYYMMDD または YYYY-MM-DD）
	CalculatedDate    string // 計算日（YYYYMMDD または YYYY-MM-DD）（code、disclosed_date、calculated_dateのいずれかが必須）
	PaginationKey     string // ページネーションキー
}

// ShortSellingPositionsResponse は空売り残高報告のレスポンスです。
type ShortSellingPositionsResponse struct {
	Data          []ShortSellingPosition `json:"data"`
	PaginationKey string                 `json:"pagination_key"` // ページネーションキー
}

// ShortSellingPosition は空売り残高報告のデータを表します。
// J-Quants API /markets/short-sale-report エンドポイントのレスポンスデータ。
// 有価証券の取引等の規制に関する内閣府令に基づく大口空売り残高（0.5％以上）の報告データ。
type ShortSellingPosition struct {
	// 基本情報
	DiscDate string `json:"DiscDate"` // 公表日（YYYY-MM-DD形式）
	CalcDate string `json:"CalcDate"` // 計算日（YYYY-MM-DD形式）
	Code     string `json:"Code"`     // 銘柄コード（5桁）

	// 空売り者情報
	SSName string `json:"SSName"` // 商号・名称・氏名（日本語名称または英語名称が混在）
	SSAddr string `json:"SSAddr"` // 住所・所在地

	// 投資一任契約関連情報
	DICName  string `json:"DICName"`  // 委託者・投資一任契約の相手方の商号・名称・氏名
	DICAddr  string `json:"DICAddr"`  // 委託者・投資一任契約の相手方の住所・所在地
	FundName string `json:"FundName"` // 信託財産・運用財産の名称

	// 空売り残高情報
	ShrtPosToSO    float64 `json:"ShrtPosToSO"`    // 空売り残高割合（発行済株式総数に対する比率）
	ShrtPosShares  float64 `json:"ShrtPosShares"`  // 空売り残高数量（株数）
	ShrtPosUnits   float64 `json:"ShrtPosUnits"`   // 空売り残高売買単位数

	// 前回報告データ
	PrevRptDate  string  `json:"PrevRptDate"`  // 直近計算年月日（YYYY-MM-DD形式）
	PrevRptRatio float64 `json:"PrevRptRatio"` // 直近空売り残高割合

	// その他
	Notes string `json:"Notes"` // 備考
}

// RawShortSellingPosition is used for unmarshaling JSON response with mixed types
type RawShortSellingPosition struct {
	// 基本情報
	DiscDate string `json:"DiscDate"`
	CalcDate string `json:"CalcDate"`
	Code     string `json:"Code"`

	// 空売り者情報
	SSName string `json:"SSName"`
	SSAddr string `json:"SSAddr"`

	// 投資一任契約関連情報
	DICName  string `json:"DICName"`
	DICAddr  string `json:"DICAddr"`
	FundName string `json:"FundName"`

	// 空売り残高情報
	ShrtPosToSO   types.Float64String `json:"ShrtPosToSO"`
	ShrtPosShares types.Float64String `json:"ShrtPosShares"`
	ShrtPosUnits  types.Float64String `json:"ShrtPosUnits"`

	// 前回報告データ
	PrevRptDate  string              `json:"PrevRptDate"`
	PrevRptRatio types.Float64String `json:"PrevRptRatio"`

	// その他
	Notes string `json:"Notes"`
}

// UnmarshalJSON implements custom JSON unmarshaling for ShortSellingPositionsResponse
func (r *ShortSellingPositionsResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawShortSellingPosition
	type rawResponse struct {
		Data          []RawShortSellingPosition `json:"data"`
		PaginationKey string                    `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawShortSellingPosition to ShortSellingPosition
	r.Data = make([]ShortSellingPosition, len(raw.Data))
	for idx, rsp := range raw.Data {
		r.Data[idx] = ShortSellingPosition{
			// 基本情報
			DiscDate: rsp.DiscDate,
			CalcDate: rsp.CalcDate,
			Code:     rsp.Code,

			// 空売り者情報
			SSName: rsp.SSName,
			SSAddr: rsp.SSAddr,

			// 投資一任契約関連情報
			DICName:  rsp.DICName,
			DICAddr:  rsp.DICAddr,
			FundName: rsp.FundName,

			// 空売り残高情報
			ShrtPosToSO:   float64(rsp.ShrtPosToSO),
			ShrtPosShares: float64(rsp.ShrtPosShares),
			ShrtPosUnits:  float64(rsp.ShrtPosUnits),

			// 前回報告データ
			PrevRptDate:  rsp.PrevRptDate,
			PrevRptRatio: float64(rsp.PrevRptRatio),

			// その他
			Notes: rsp.Notes,
		}
	}

	return nil
}

// GetShortSellingPositions は空売り残高報告を取得します。
func (s *ShortSellingPositionsService) GetShortSellingPositions(params ShortSellingPositionsParams) (*ShortSellingPositionsResponse, error) {
	// code、disc_date、calc_dateのいずれかが必須
	if params.Code == "" && params.DisclosedDate == "" && params.CalculatedDate == "" {
		return nil, fmt.Errorf("either code, disc_date, or calc_date parameter is required")
	}

	path := "/markets/short-sale-report"

	query := "?"
	if params.Code != "" {
		query += fmt.Sprintf("code=%s&", params.Code)
	}
	if params.DisclosedDate != "" {
		query += fmt.Sprintf("disc_date=%s&", params.DisclosedDate)
	}
	if params.DisclosedDateFrom != "" {
		query += fmt.Sprintf("disc_date_from=%s&", params.DisclosedDateFrom)
	}
	if params.DisclosedDateTo != "" {
		query += fmt.Sprintf("disc_date_to=%s&", params.DisclosedDateTo)
	}
	if params.CalculatedDate != "" {
		query += fmt.Sprintf("calc_date=%s&", params.CalculatedDate)
	}
	if params.PaginationKey != "" {
		query += fmt.Sprintf("pagination_key=%s&", params.PaginationKey)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp ShortSellingPositionsResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get short selling positions: %w", err)
	}

	return &resp, nil
}

// GetShortSellingPositionsByCode は指定銘柄の空売り残高報告を取得します。
// ページネーションを使用して全データを取得します。
func (s *ShortSellingPositionsService) GetShortSellingPositionsByCode(code string) ([]ShortSellingPosition, error) {
	var allData []ShortSellingPosition
	paginationKey := ""

	for {
		params := ShortSellingPositionsParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetShortSellingPositions(params)
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

// GetShortSellingPositionsByDisclosedDate は指定公表日の全銘柄空売り残高報告を取得します。
// ページネーションを使用して全データを取得します。
func (s *ShortSellingPositionsService) GetShortSellingPositionsByDisclosedDate(disclosedDate string) ([]ShortSellingPosition, error) {
	var allData []ShortSellingPosition
	paginationKey := ""

	for {
		params := ShortSellingPositionsParams{
			DisclosedDate: disclosedDate,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetShortSellingPositions(params)
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

// GetShortSellingPositionsByCalculatedDate は指定計算日の全銘柄空売り残高報告を取得します。
// ページネーションを使用して全データを取得します。
func (s *ShortSellingPositionsService) GetShortSellingPositionsByCalculatedDate(calculatedDate string) ([]ShortSellingPosition, error) {
	var allData []ShortSellingPosition
	paginationKey := ""

	for {
		params := ShortSellingPositionsParams{
			CalculatedDate: calculatedDate,
			PaginationKey:  paginationKey,
		}

		resp, err := s.GetShortSellingPositions(params)
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

// GetShortSellingPositionsByCodeAndDateRange は指定銘柄・期間の空売り残高報告を取得します。
func (s *ShortSellingPositionsService) GetShortSellingPositionsByCodeAndDateRange(code, fromDate, toDate string) ([]ShortSellingPosition, error) {
	var allData []ShortSellingPosition
	paginationKey := ""

	for {
		params := ShortSellingPositionsParams{
			Code:              code,
			DisclosedDateFrom: fromDate,
			DisclosedDateTo:   toDate,
			PaginationKey:     paginationKey,
		}

		resp, err := s.GetShortSellingPositions(params)
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

// GetPositionChange は前回報告からの残高変化を計算します（株数）。
func (ssp *ShortSellingPosition) GetPositionChange() float64 {
	// 前回報告時の残高割合から株数を逆算
	// 注：発行済株式総数が変わっていない前提での概算
	if ssp.ShrtPosToSO == 0 {
		return 0
	}

	// 現在の残高株数 / 現在の残高割合 = 発行済株式総数（概算）
	estimatedSharesOutstanding := ssp.ShrtPosShares / ssp.ShrtPosToSO

	// 前回の残高株数（概算）
	previousShares := estimatedSharesOutstanding * ssp.PrevRptRatio

	return ssp.ShrtPosShares - previousShares
}

// GetPositionChangeRatio は前回報告からの残高変化率を計算します（パーセント）。
func (ssp *ShortSellingPosition) GetPositionChangeRatio() float64 {
	if ssp.PrevRptRatio == 0 {
		return 0
	}

	changeRatio := (ssp.ShrtPosToSO - ssp.PrevRptRatio) / ssp.PrevRptRatio

	return changeRatio * 100
}

// IsIncrease は前回報告から残高が増加したかを判定します。
func (ssp *ShortSellingPosition) IsIncrease() bool {
	return ssp.ShrtPosToSO > ssp.PrevRptRatio
}

// IsDecrease は前回報告から残高が減少したかを判定します。
func (ssp *ShortSellingPosition) IsDecrease() bool {
	return ssp.ShrtPosToSO < ssp.PrevRptRatio
}

// IsNoChange は前回報告から残高が変化していないかを判定します。
func (ssp *ShortSellingPosition) IsNoChange() bool {
	return ssp.ShrtPosToSO == ssp.PrevRptRatio
}

// HasDiscretionaryInvestment は投資一任契約があるかを判定します。
func (ssp *ShortSellingPosition) HasDiscretionaryInvestment() bool {
	return ssp.DICName != "" || ssp.DICAddr != "" || ssp.FundName != ""
}

// IsIndividual は個人投資家かを判定します。
func (ssp *ShortSellingPosition) IsIndividual() bool {
	return ssp.SSName == "個人"
}

// GetPositionPercentage は残高割合をパーセント表記で取得します。
func (ssp *ShortSellingPosition) GetPositionPercentage() float64 {
	return ssp.ShrtPosToSO * 100
}

// GetPreviousPositionPercentage は前回報告の残高割合をパーセント表記で取得します。
func (ssp *ShortSellingPosition) GetPreviousPositionPercentage() float64 {
	return ssp.PrevRptRatio * 100
}
