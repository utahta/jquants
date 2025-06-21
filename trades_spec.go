package jquants

import (
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// TradesSpecService は投資部門別情報を取得するサービスです。
// 投資部門別売買状況（株式・金額）のデータを提供します。
type TradesSpecService struct {
	client client.HTTPClient
}

// NewTradesSpecService は新しいTradesSpecServiceを作成します。
func NewTradesSpecService(c client.HTTPClient) *TradesSpecService {
	return &TradesSpecService{client: c}
}

// TradesSpec は投資部門別情報データを表します。
// J-Quants API /markets/trades_spec エンドポイントのレスポンスデータ。
type TradesSpec struct {
	// 基本情報
	PublishedDate string `json:"PublishedDate"` // 公表日（YYYY-MM-DD形式）
	StartDate     string `json:"StartDate"`     // 開始日（YYYY-MM-DD形式）
	EndDate       string `json:"EndDate"`       // 終了日（YYYY-MM-DD形式）
	Section       string `json:"Section"`       // 市場名

	// 自己取引
	ProprietarySales     float64 `json:"ProprietarySales"`     // 自己計_売（千円）
	ProprietaryPurchases float64 `json:"ProprietaryPurchases"` // 自己計_買（千円）
	ProprietaryTotal     float64 `json:"ProprietaryTotal"`     // 自己計_合計（千円）
	ProprietaryBalance   float64 `json:"ProprietaryBalance"`   // 自己計_差引（千円）

	// 委託取引
	BrokerageSales     float64 `json:"BrokerageSales"`     // 委託計_売（千円）
	BrokeragePurchases float64 `json:"BrokeragePurchases"` // 委託計_買（千円）
	BrokerageTotal     float64 `json:"BrokerageTotal"`     // 委託計_合計（千円）
	BrokerageBalance   float64 `json:"BrokerageBalance"`   // 委託計_差引（千円）

	// 総計
	TotalSales     float64 `json:"TotalSales"`     // 総計_売（千円）
	TotalPurchases float64 `json:"TotalPurchases"` // 総計_買（千円）
	TotalTotal     float64 `json:"TotalTotal"`     // 総計_合計（千円）
	TotalBalance   float64 `json:"TotalBalance"`   // 総計_差引（千円）

	// 投資部門別内訳
	IndividualsSales                    float64 `json:"IndividualsSales"`                    // 個人_売（千円）
	IndividualsPurchases                float64 `json:"IndividualsPurchases"`                // 個人_買（千円）
	IndividualsTotal                    float64 `json:"IndividualsTotal"`                    // 個人_合計（千円）
	IndividualsBalance                  float64 `json:"IndividualsBalance"`                  // 個人_差引（千円）
	ForeignersSales                     float64 `json:"ForeignersSales"`                     // 海外投資家_売（千円）
	ForeignersPurchases                 float64 `json:"ForeignersPurchases"`                 // 海外投資家_買（千円）
	ForeignersTotal                     float64 `json:"ForeignersTotal"`                     // 海外投資家_合計（千円）
	ForeignersBalance                   float64 `json:"ForeignersBalance"`                   // 海外投資家_差引（千円）
	SecuritiesCosSales                  float64 `json:"SecuritiesCosSales"`                  // 証券会社_売（千円）
	SecuritiesCosPurchases              float64 `json:"SecuritiesCosPurchases"`              // 証券会社_買（千円）
	SecuritiesCosTotal                  float64 `json:"SecuritiesCosTotal"`                  // 証券会社_合計（千円）
	SecuritiesCosBalance                float64 `json:"SecuritiesCosBalance"`                // 証券会社_差引（千円）
	InvestmentTrustsSales               float64 `json:"InvestmentTrustsSales"`               // 投資信託_売（千円）
	InvestmentTrustsPurchases           float64 `json:"InvestmentTrustsPurchases"`           // 投資信託_買（千円）
	InvestmentTrustsTotal               float64 `json:"InvestmentTrustsTotal"`               // 投資信託_合計（千円）
	InvestmentTrustsBalance             float64 `json:"InvestmentTrustsBalance"`             // 投資信託_差引（千円）
	BusinessCosSales                    float64 `json:"BusinessCosSales"`                    // 事業法人_売（千円）
	BusinessCosPurchases                float64 `json:"BusinessCosPurchases"`                // 事業法人_買（千円）
	BusinessCosTotal                    float64 `json:"BusinessCosTotal"`                    // 事業法人_合計（千円）
	BusinessCosBalance                  float64 `json:"BusinessCosBalance"`                  // 事業法人_差引（千円）
	OtherCosSales                       float64 `json:"OtherCosSales"`                       // その他法人_売（千円）
	OtherCosPurchases                   float64 `json:"OtherCosPurchases"`                   // その他法人_買（千円）
	OtherCosTotal                       float64 `json:"OtherCosTotal"`                       // その他法人_合計（千円）
	OtherCosBalance                     float64 `json:"OtherCosBalance"`                     // その他法人_差引（千円）
	InsuranceCosSales                   float64 `json:"InsuranceCosSales"`                   // 生保・損保_売（千円）
	InsuranceCosPurchases               float64 `json:"InsuranceCosPurchases"`               // 生保・損保_買（千円）
	InsuranceCosTotal                   float64 `json:"InsuranceCosTotal"`                   // 生保・損保_合計（千円）
	InsuranceCosBalance                 float64 `json:"InsuranceCosBalance"`                 // 生保・損保_差引（千円）
	CityBKsRegionalBKsEtcSales          float64 `json:"CityBKsRegionalBKsEtcSales"`          // 都銀・地銀等_売（千円）
	CityBKsRegionalBKsEtcPurchases      float64 `json:"CityBKsRegionalBKsEtcPurchases"`      // 都銀・地銀等_買（千円）
	CityBKsRegionalBKsEtcTotal          float64 `json:"CityBKsRegionalBKsEtcTotal"`          // 都銀・地銀等_合計（千円）
	CityBKsRegionalBKsEtcBalance        float64 `json:"CityBKsRegionalBKsEtcBalance"`        // 都銀・地銀等_差引（千円）
	TrustBanksSales                     float64 `json:"TrustBanksSales"`                     // 信託銀行_売（千円）
	TrustBanksPurchases                 float64 `json:"TrustBanksPurchases"`                 // 信託銀行_買（千円）
	TrustBanksTotal                     float64 `json:"TrustBanksTotal"`                     // 信託銀行_合計（千円）
	TrustBanksBalance                   float64 `json:"TrustBanksBalance"`                   // 信託銀行_差引（千円）
	OtherFinancialInstitutionsSales     float64 `json:"OtherFinancialInstitutionsSales"`     // その他金融機関_売（千円）
	OtherFinancialInstitutionsPurchases float64 `json:"OtherFinancialInstitutionsPurchases"` // その他金融機関_買（千円）
	OtherFinancialInstitutionsTotal     float64 `json:"OtherFinancialInstitutionsTotal"`     // その他金融機関_合計（千円）
	OtherFinancialInstitutionsBalance   float64 `json:"OtherFinancialInstitutionsBalance"`   // その他金融機関_差引（千円）
}

// RawTradesSpec is used for unmarshaling JSON response with mixed types
type RawTradesSpec struct {
	// 基本情報
	PublishedDate string `json:"PublishedDate"`
	StartDate     string `json:"StartDate"`
	EndDate       string `json:"EndDate"`
	Section       string `json:"Section"`

	// 自己取引
	ProprietarySales     types.Float64String `json:"ProprietarySales"`
	ProprietaryPurchases types.Float64String `json:"ProprietaryPurchases"`
	ProprietaryTotal     types.Float64String `json:"ProprietaryTotal"`
	ProprietaryBalance   types.Float64String `json:"ProprietaryBalance"`

	// 委託取引
	BrokerageSales     types.Float64String `json:"BrokerageSales"`
	BrokeragePurchases types.Float64String `json:"BrokeragePurchases"`
	BrokerageTotal     types.Float64String `json:"BrokerageTotal"`
	BrokerageBalance   types.Float64String `json:"BrokerageBalance"`

	// 総計
	TotalSales     types.Float64String `json:"TotalSales"`
	TotalPurchases types.Float64String `json:"TotalPurchases"`
	TotalTotal     types.Float64String `json:"TotalTotal"`
	TotalBalance   types.Float64String `json:"TotalBalance"`

	// 投資部門別内訳
	IndividualsSales                    types.Float64String `json:"IndividualsSales"`
	IndividualsPurchases                types.Float64String `json:"IndividualsPurchases"`
	IndividualsTotal                    types.Float64String `json:"IndividualsTotal"`
	IndividualsBalance                  types.Float64String `json:"IndividualsBalance"`
	ForeignersSales                     types.Float64String `json:"ForeignersSales"`
	ForeignersPurchases                 types.Float64String `json:"ForeignersPurchases"`
	ForeignersTotal                     types.Float64String `json:"ForeignersTotal"`
	ForeignersBalance                   types.Float64String `json:"ForeignersBalance"`
	SecuritiesCosSales                  types.Float64String `json:"SecuritiesCosSales"`
	SecuritiesCosPurchases              types.Float64String `json:"SecuritiesCosPurchases"`
	SecuritiesCosTotal                  types.Float64String `json:"SecuritiesCosTotal"`
	SecuritiesCosBalance                types.Float64String `json:"SecuritiesCosBalance"`
	InvestmentTrustsSales               types.Float64String `json:"InvestmentTrustsSales"`
	InvestmentTrustsPurchases           types.Float64String `json:"InvestmentTrustsPurchases"`
	InvestmentTrustsTotal               types.Float64String `json:"InvestmentTrustsTotal"`
	InvestmentTrustsBalance             types.Float64String `json:"InvestmentTrustsBalance"`
	BusinessCosSales                    types.Float64String `json:"BusinessCosSales"`
	BusinessCosPurchases                types.Float64String `json:"BusinessCosPurchases"`
	BusinessCosTotal                    types.Float64String `json:"BusinessCosTotal"`
	BusinessCosBalance                  types.Float64String `json:"BusinessCosBalance"`
	OtherCosSales                       types.Float64String `json:"OtherCosSales"`
	OtherCosPurchases                   types.Float64String `json:"OtherCosPurchases"`
	OtherCosTotal                       types.Float64String `json:"OtherCosTotal"`
	OtherCosBalance                     types.Float64String `json:"OtherCosBalance"`
	InsuranceCosSales                   types.Float64String `json:"InsuranceCosSales"`
	InsuranceCosPurchases               types.Float64String `json:"InsuranceCosPurchases"`
	InsuranceCosTotal                   types.Float64String `json:"InsuranceCosTotal"`
	InsuranceCosBalance                 types.Float64String `json:"InsuranceCosBalance"`
	CityBKsRegionalBKsEtcSales          types.Float64String `json:"CityBKsRegionalBKsEtcSales"`
	CityBKsRegionalBKsEtcPurchases      types.Float64String `json:"CityBKsRegionalBKsEtcPurchases"`
	CityBKsRegionalBKsEtcTotal          types.Float64String `json:"CityBKsRegionalBKsEtcTotal"`
	CityBKsRegionalBKsEtcBalance        types.Float64String `json:"CityBKsRegionalBKsEtcBalance"`
	TrustBanksSales                     types.Float64String `json:"TrustBanksSales"`
	TrustBanksPurchases                 types.Float64String `json:"TrustBanksPurchases"`
	TrustBanksTotal                     types.Float64String `json:"TrustBanksTotal"`
	TrustBanksBalance                   types.Float64String `json:"TrustBanksBalance"`
	OtherFinancialInstitutionsSales     types.Float64String `json:"OtherFinancialInstitutionsSales"`
	OtherFinancialInstitutionsPurchases types.Float64String `json:"OtherFinancialInstitutionsPurchases"`
	OtherFinancialInstitutionsTotal     types.Float64String `json:"OtherFinancialInstitutionsTotal"`
	OtherFinancialInstitutionsBalance   types.Float64String `json:"OtherFinancialInstitutionsBalance"`
}

// TradesSpecResponse は投資部門別情報のレスポンスです。
type TradesSpecResponse struct {
	TradesSpec    []TradesSpec `json:"trades_spec"`
	PaginationKey string       `json:"pagination_key"` // ページネーションキー
}

// UnmarshalJSON implements custom JSON unmarshaling
func (t *TradesSpecResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawTradesSpec
	type rawResponse struct {
		TradesSpec    []RawTradesSpec `json:"trades_spec"`
		PaginationKey string          `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	t.PaginationKey = raw.PaginationKey

	// Convert RawTradesSpec to TradesSpec
	t.TradesSpec = make([]TradesSpec, len(raw.TradesSpec))
	for idx, rt := range raw.TradesSpec {
		t.TradesSpec[idx] = TradesSpec{
			// 基本情報
			PublishedDate: rt.PublishedDate,
			StartDate:     rt.StartDate,
			EndDate:       rt.EndDate,
			Section:       rt.Section,

			// 自己取引
			ProprietarySales:     float64(rt.ProprietarySales),
			ProprietaryPurchases: float64(rt.ProprietaryPurchases),
			ProprietaryTotal:     float64(rt.ProprietaryTotal),
			ProprietaryBalance:   float64(rt.ProprietaryBalance),

			// 委託取引
			BrokerageSales:     float64(rt.BrokerageSales),
			BrokeragePurchases: float64(rt.BrokeragePurchases),
			BrokerageTotal:     float64(rt.BrokerageTotal),
			BrokerageBalance:   float64(rt.BrokerageBalance),

			// 総計
			TotalSales:     float64(rt.TotalSales),
			TotalPurchases: float64(rt.TotalPurchases),
			TotalTotal:     float64(rt.TotalTotal),
			TotalBalance:   float64(rt.TotalBalance),

			// 投資部門別内訳
			IndividualsSales:                    float64(rt.IndividualsSales),
			IndividualsPurchases:                float64(rt.IndividualsPurchases),
			IndividualsTotal:                    float64(rt.IndividualsTotal),
			IndividualsBalance:                  float64(rt.IndividualsBalance),
			ForeignersSales:                     float64(rt.ForeignersSales),
			ForeignersPurchases:                 float64(rt.ForeignersPurchases),
			ForeignersTotal:                     float64(rt.ForeignersTotal),
			ForeignersBalance:                   float64(rt.ForeignersBalance),
			SecuritiesCosSales:                  float64(rt.SecuritiesCosSales),
			SecuritiesCosPurchases:              float64(rt.SecuritiesCosPurchases),
			SecuritiesCosTotal:                  float64(rt.SecuritiesCosTotal),
			SecuritiesCosBalance:                float64(rt.SecuritiesCosBalance),
			InvestmentTrustsSales:               float64(rt.InvestmentTrustsSales),
			InvestmentTrustsPurchases:           float64(rt.InvestmentTrustsPurchases),
			InvestmentTrustsTotal:               float64(rt.InvestmentTrustsTotal),
			InvestmentTrustsBalance:             float64(rt.InvestmentTrustsBalance),
			BusinessCosSales:                    float64(rt.BusinessCosSales),
			BusinessCosPurchases:                float64(rt.BusinessCosPurchases),
			BusinessCosTotal:                    float64(rt.BusinessCosTotal),
			BusinessCosBalance:                  float64(rt.BusinessCosBalance),
			OtherCosSales:                       float64(rt.OtherCosSales),
			OtherCosPurchases:                   float64(rt.OtherCosPurchases),
			OtherCosTotal:                       float64(rt.OtherCosTotal),
			OtherCosBalance:                     float64(rt.OtherCosBalance),
			InsuranceCosSales:                   float64(rt.InsuranceCosSales),
			InsuranceCosPurchases:               float64(rt.InsuranceCosPurchases),
			InsuranceCosTotal:                   float64(rt.InsuranceCosTotal),
			InsuranceCosBalance:                 float64(rt.InsuranceCosBalance),
			CityBKsRegionalBKsEtcSales:          float64(rt.CityBKsRegionalBKsEtcSales),
			CityBKsRegionalBKsEtcPurchases:      float64(rt.CityBKsRegionalBKsEtcPurchases),
			CityBKsRegionalBKsEtcTotal:          float64(rt.CityBKsRegionalBKsEtcTotal),
			CityBKsRegionalBKsEtcBalance:        float64(rt.CityBKsRegionalBKsEtcBalance),
			TrustBanksSales:                     float64(rt.TrustBanksSales),
			TrustBanksPurchases:                 float64(rt.TrustBanksPurchases),
			TrustBanksTotal:                     float64(rt.TrustBanksTotal),
			TrustBanksBalance:                   float64(rt.TrustBanksBalance),
			OtherFinancialInstitutionsSales:     float64(rt.OtherFinancialInstitutionsSales),
			OtherFinancialInstitutionsPurchases: float64(rt.OtherFinancialInstitutionsPurchases),
			OtherFinancialInstitutionsTotal:     float64(rt.OtherFinancialInstitutionsTotal),
			OtherFinancialInstitutionsBalance:   float64(rt.OtherFinancialInstitutionsBalance),
		}
	}

	return nil
}

// TradesSpecParams は投資部門別情報のリクエストパラメータです。
type TradesSpecParams struct {
	Section       string // セクション（市場）（例: TSEPrime, TSEStandard, TSEGrowth）
	From          string // 開始日（YYYYMMDD または YYYY-MM-DD）
	To            string // 終了日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// 市場コード定数
const (
	SectionTSE1st      = "TSE1st"      // 市場一部
	SectionTSE2nd      = "TSE2nd"      // 市場二部
	SectionTSEMothers  = "TSEMothers"  // マザーズ
	SectionTSEJASDAQ   = "TSEJASDAQ"   // JASDAQ
	SectionTSEPrime    = "TSEPrime"    // プライム
	SectionTSEStandard = "TSEStandard" // スタンダード
	SectionTSEGrowth   = "TSEGrowth"   // グロース
	SectionTokyoNagoya = "TokyoNagoya" // 東証および名証
)

// GetTradesSpec は投資部門別情報を取得します。
// パラメータの組み合わせ:
// - section指定あり、from/to指定あり: 指定したセクションの指定した期間のデータ
// - section指定あり、from/to指定なし: 指定したセクションの全期間のデータ
// - section指定なし、from/to指定あり: すべてのセクションの指定した期間のデータ
// - section指定なし、from/to指定なし: すべてのセクションの全期間のデータ
func (s *TradesSpecService) GetTradesSpec(params TradesSpecParams) (*TradesSpecResponse, error) {
	path := "/markets/trades_spec"

	query := "?"
	if params.Section != "" {
		query += fmt.Sprintf("section=%s&", params.Section)
	}
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

	var resp TradesSpecResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get trades spec: %w", err)
	}

	return &resp, nil
}

// GetTradesSpecByDateRange は指定期間の投資部門別情報を取得します。
// ページネーションを使用して全データを取得します。
func (s *TradesSpecService) GetTradesSpecByDateRange(from, to string) ([]TradesSpec, error) {
	var allTradesSpec []TradesSpec
	paginationKey := ""

	for {
		params := TradesSpecParams{
			From:          from,
			To:            to,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetTradesSpec(params)
		if err != nil {
			return nil, err
		}

		allTradesSpec = append(allTradesSpec, resp.TradesSpec...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allTradesSpec, nil
}

// GetTradesSpecBySection は指定セクションの投資部門別情報を取得します。
// ページネーションを使用して全データを取得します。
func (s *TradesSpecService) GetTradesSpecBySection(section string) ([]TradesSpec, error) {
	var allTradesSpec []TradesSpec
	paginationKey := ""

	for {
		params := TradesSpecParams{
			Section:       section,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetTradesSpec(params)
		if err != nil {
			return nil, err
		}

		allTradesSpec = append(allTradesSpec, resp.TradesSpec...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allTradesSpec, nil
}

// GetAllTradesSpec は全セクション・全期間の投資部門別情報を取得します。
// ページネーションを使用して大量データを分割取得します。
func (s *TradesSpecService) GetAllTradesSpec() ([]TradesSpec, error) {
	var allTradesSpec []TradesSpec
	paginationKey := ""

	for {
		params := TradesSpecParams{
			PaginationKey: paginationKey,
		}

		resp, err := s.GetTradesSpec(params)
		if err != nil {
			return nil, err
		}

		allTradesSpec = append(allTradesSpec, resp.TradesSpec...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allTradesSpec, nil
}

// IsBuyerDominant は買い手が優勢かどうかを判定します（差引がプラス）。
func (ts *TradesSpec) IsBuyerDominant(investorType string) bool {
	switch investorType {
	case "individuals":
		return ts.IndividualsBalance > 0
	case "foreigners":
		return ts.ForeignersBalance > 0
	case "securities":
		return ts.SecuritiesCosBalance > 0
	case "investment_trusts":
		return ts.InvestmentTrustsBalance > 0
	case "business":
		return ts.BusinessCosBalance > 0
	case "insurance":
		return ts.InsuranceCosBalance > 0
	case "trust_banks":
		return ts.TrustBanksBalance > 0
	case "total":
		return ts.TotalBalance > 0
	default:
		return false
	}
}

// GetNetFlow は指定した投資家タイプの純流入額を取得します（差引額）。
func (ts *TradesSpec) GetNetFlow(investorType string) float64 {
	switch investorType {
	case "individuals":
		return ts.IndividualsBalance
	case "foreigners":
		return ts.ForeignersBalance
	case "securities":
		return ts.SecuritiesCosBalance
	case "investment_trusts":
		return ts.InvestmentTrustsBalance
	case "business":
		return ts.BusinessCosBalance
	case "insurance":
		return ts.InsuranceCosBalance
	case "trust_banks":
		return ts.TrustBanksBalance
	case "total":
		return ts.TotalBalance
	default:
		return 0
	}
}
