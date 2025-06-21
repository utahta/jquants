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
	ShortSellingPositions []ShortSellingPosition `json:"short_selling_positions"`
	PaginationKey         string                 `json:"pagination_key"` // ページネーションキー
}

// ShortSellingPosition は空売り残高報告のデータを表します。
// J-Quants API /markets/short_selling_positions エンドポイントのレスポンスデータ。
// 有価証券の取引等の規制に関する内閣府令に基づく大口空売り残高（0.5％以上）の報告データ。
type ShortSellingPosition struct {
	// 基本情報
	DisclosedDate  string `json:"DisclosedDate"`  // 公表日（YYYY-MM-DD形式）
	CalculatedDate string `json:"CalculatedDate"` // 計算日（YYYY-MM-DD形式）
	Code           string `json:"Code"`           // 銘柄コード（5桁）

	// 空売り者情報
	ShortSellerName    string `json:"ShortSellerName"`    // 商号・名称・氏名（日本語名称または英語名称が混在）
	ShortSellerAddress string `json:"ShortSellerAddress"` // 住所・所在地

	// 投資一任契約関連情報
	DiscretionaryInvestmentContractorName    string `json:"DiscretionaryInvestmentContractorName"`    // 委託者・投資一任契約の相手方の商号・名称・氏名
	DiscretionaryInvestmentContractorAddress string `json:"DiscretionaryInvestmentContractorAddress"` // 委託者・投資一任契約の相手方の住所・所在地
	InvestmentFundName                       string `json:"InvestmentFundName"`                       // 信託財産・運用財産の名称

	// 空売り残高情報
	ShortPositionsToSharesOutstandingRatio float64 `json:"ShortPositionsToSharesOutstandingRatio"` // 空売り残高割合（発行済株式総数に対する比率）
	ShortPositionsInSharesNumber           float64 `json:"ShortPositionsInSharesNumber"`           // 空売り残高数量（株数）
	ShortPositionsInTradingUnitsNumber     float64 `json:"ShortPositionsInTradingUnitsNumber"`     // 空売り残高売買単位数

	// 前回報告データ
	CalculationInPreviousReportingDate     string  `json:"CalculationInPreviousReportingDate"`     // 直近計算年月日（YYYY-MM-DD形式）
	ShortPositionsInPreviousReportingRatio float64 `json:"ShortPositionsInPreviousReportingRatio"` // 直近空売り残高割合

	// その他
	Notes string `json:"Notes"` // 備考
}

// RawShortSellingPosition is used for unmarshaling JSON response with mixed types
type RawShortSellingPosition struct {
	// 基本情報
	DisclosedDate  string `json:"DisclosedDate"`
	CalculatedDate string `json:"CalculatedDate"`
	Code           string `json:"Code"`

	// 空売り者情報
	ShortSellerName    string `json:"ShortSellerName"`
	ShortSellerAddress string `json:"ShortSellerAddress"`

	// 投資一任契約関連情報
	DiscretionaryInvestmentContractorName    string `json:"DiscretionaryInvestmentContractorName"`
	DiscretionaryInvestmentContractorAddress string `json:"DiscretionaryInvestmentContractorAddress"`
	InvestmentFundName                       string `json:"InvestmentFundName"`

	// 空売り残高情報
	ShortPositionsToSharesOutstandingRatio types.Float64String `json:"ShortPositionsToSharesOutstandingRatio"`
	ShortPositionsInSharesNumber           types.Float64String `json:"ShortPositionsInSharesNumber"`
	ShortPositionsInTradingUnitsNumber     types.Float64String `json:"ShortPositionsInTradingUnitsNumber"`

	// 前回報告データ
	CalculationInPreviousReportingDate     string              `json:"CalculationInPreviousReportingDate"`
	ShortPositionsInPreviousReportingRatio types.Float64String `json:"ShortPositionsInPreviousReportingRatio"`

	// その他
	Notes string `json:"Notes"`
}

// UnmarshalJSON implements custom JSON unmarshaling for ShortSellingPositionsResponse
func (r *ShortSellingPositionsResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawShortSellingPosition
	type rawResponse struct {
		ShortSellingPositions []RawShortSellingPosition `json:"short_selling_positions"`
		PaginationKey         string                    `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawShortSellingPosition to ShortSellingPosition
	r.ShortSellingPositions = make([]ShortSellingPosition, len(raw.ShortSellingPositions))
	for idx, rsp := range raw.ShortSellingPositions {
		r.ShortSellingPositions[idx] = ShortSellingPosition{
			// 基本情報
			DisclosedDate:  rsp.DisclosedDate,
			CalculatedDate: rsp.CalculatedDate,
			Code:           rsp.Code,

			// 空売り者情報
			ShortSellerName:    rsp.ShortSellerName,
			ShortSellerAddress: rsp.ShortSellerAddress,

			// 投資一任契約関連情報
			DiscretionaryInvestmentContractorName:    rsp.DiscretionaryInvestmentContractorName,
			DiscretionaryInvestmentContractorAddress: rsp.DiscretionaryInvestmentContractorAddress,
			InvestmentFundName:                       rsp.InvestmentFundName,

			// 空売り残高情報
			ShortPositionsToSharesOutstandingRatio: float64(rsp.ShortPositionsToSharesOutstandingRatio),
			ShortPositionsInSharesNumber:           float64(rsp.ShortPositionsInSharesNumber),
			ShortPositionsInTradingUnitsNumber:     float64(rsp.ShortPositionsInTradingUnitsNumber),

			// 前回報告データ
			CalculationInPreviousReportingDate:     rsp.CalculationInPreviousReportingDate,
			ShortPositionsInPreviousReportingRatio: float64(rsp.ShortPositionsInPreviousReportingRatio),

			// その他
			Notes: rsp.Notes,
		}
	}

	return nil
}

// GetShortSellingPositions は空売り残高報告を取得します。
func (s *ShortSellingPositionsService) GetShortSellingPositions(params ShortSellingPositionsParams) (*ShortSellingPositionsResponse, error) {
	// code、disclosed_date、calculated_dateのいずれかが必須
	if params.Code == "" && params.DisclosedDate == "" && params.CalculatedDate == "" {
		return nil, fmt.Errorf("either code, disclosed_date, or calculated_date parameter is required")
	}

	path := "/markets/short_selling_positions"

	query := "?"
	if params.Code != "" {
		query += fmt.Sprintf("code=%s&", params.Code)
	}
	if params.DisclosedDate != "" {
		query += fmt.Sprintf("disclosed_date=%s&", params.DisclosedDate)
	}
	if params.DisclosedDateFrom != "" {
		query += fmt.Sprintf("disclosed_date_from=%s&", params.DisclosedDateFrom)
	}
	if params.DisclosedDateTo != "" {
		query += fmt.Sprintf("disclosed_date_to=%s&", params.DisclosedDateTo)
	}
	if params.CalculatedDate != "" {
		query += fmt.Sprintf("calculated_date=%s&", params.CalculatedDate)
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

		allData = append(allData, resp.ShortSellingPositions...)

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

		allData = append(allData, resp.ShortSellingPositions...)

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

		allData = append(allData, resp.ShortSellingPositions...)

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

		allData = append(allData, resp.ShortSellingPositions...)

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
	if ssp.ShortPositionsToSharesOutstandingRatio == 0 {
		return 0
	}

	// 現在の残高株数 / 現在の残高割合 = 発行済株式総数（概算）
	estimatedSharesOutstanding := ssp.ShortPositionsInSharesNumber / ssp.ShortPositionsToSharesOutstandingRatio

	// 前回の残高株数（概算）
	previousShares := estimatedSharesOutstanding * ssp.ShortPositionsInPreviousReportingRatio

	return ssp.ShortPositionsInSharesNumber - previousShares
}

// GetPositionChangeRatio は前回報告からの残高変化率を計算します（パーセント）。
func (ssp *ShortSellingPosition) GetPositionChangeRatio() float64 {
	if ssp.ShortPositionsInPreviousReportingRatio == 0 {
		return 0
	}

	changeRatio := (ssp.ShortPositionsToSharesOutstandingRatio - ssp.ShortPositionsInPreviousReportingRatio) /
		ssp.ShortPositionsInPreviousReportingRatio

	return changeRatio * 100
}

// IsIncrease は前回報告から残高が増加したかを判定します。
func (ssp *ShortSellingPosition) IsIncrease() bool {
	return ssp.ShortPositionsToSharesOutstandingRatio > ssp.ShortPositionsInPreviousReportingRatio
}

// IsDecrease は前回報告から残高が減少したかを判定します。
func (ssp *ShortSellingPosition) IsDecrease() bool {
	return ssp.ShortPositionsToSharesOutstandingRatio < ssp.ShortPositionsInPreviousReportingRatio
}

// IsNoChange は前回報告から残高が変化していないかを判定します。
func (ssp *ShortSellingPosition) IsNoChange() bool {
	return ssp.ShortPositionsToSharesOutstandingRatio == ssp.ShortPositionsInPreviousReportingRatio
}

// HasDiscretionaryInvestment は投資一任契約があるかを判定します。
func (ssp *ShortSellingPosition) HasDiscretionaryInvestment() bool {
	return ssp.DiscretionaryInvestmentContractorName != "" ||
		ssp.DiscretionaryInvestmentContractorAddress != "" ||
		ssp.InvestmentFundName != ""
}

// IsIndividual は個人投資家かを判定します。
func (ssp *ShortSellingPosition) IsIndividual() bool {
	return ssp.ShortSellerName == "個人"
}

// GetPositionPercentage は残高割合をパーセント表記で取得します。
func (ssp *ShortSellingPosition) GetPositionPercentage() float64 {
	return ssp.ShortPositionsToSharesOutstandingRatio * 100
}

// GetPreviousPositionPercentage は前回報告の残高割合をパーセント表記で取得します。
func (ssp *ShortSellingPosition) GetPreviousPositionPercentage() float64 {
	return ssp.ShortPositionsInPreviousReportingRatio * 100
}
