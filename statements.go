package jquants

import (
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// StatementsService は財務諸表データを取得するサービスです。
// 売上高、利益、資産、ROE/ROAなどの財務指標を提供します。
// 財務諸表API(/fins/statements)から四半期毎の決算短信サマリーデータを取得します。
type StatementsService struct {
	client client.HTTPClient
}

func NewStatementsService(c client.HTTPClient) *StatementsService {
	return &StatementsService{
		client: c,
	}
}

// Statement は財務諸表データを表します。
// J-Quants API /fins/statements エンドポイントのレスポンスデータ。
// すべてのフィールドはAPIでは文字列型で返されますが、このstructでは適切な型に変換されています。
type Statement struct {
	// 基本情報
	DisclosedDate              string `json:"DisclosedDate"`
	DisclosedTime              string `json:"DisclosedTime"`
	LocalCode                  string `json:"LocalCode"`
	DisclosureNumber           string `json:"DisclosureNumber"`
	TypeOfDocument             TypeOfDocument `json:"TypeOfDocument"`
	TypeOfCurrentPeriod        string `json:"TypeOfCurrentPeriod"`
	CurrentPeriodStartDate     string `json:"CurrentPeriodStartDate"`
	CurrentPeriodEndDate       string `json:"CurrentPeriodEndDate"`
	CurrentFiscalYearStartDate string `json:"CurrentFiscalYearStartDate"`
	CurrentFiscalYearEndDate   string `json:"CurrentFiscalYearEndDate"`
	NextFiscalYearStartDate    string `json:"NextFiscalYearStartDate"`
	NextFiscalYearEndDate      string `json:"NextFiscalYearEndDate"`

	// 連結財務数値
	NetSales                         *float64 `json:"NetSales"`
	OperatingProfit                  *float64 `json:"OperatingProfit"`
	OrdinaryProfit                   *float64 `json:"OrdinaryProfit"`
	Profit                           *float64 `json:"Profit"`
	TotalAssets                      *float64 `json:"TotalAssets"`
	Equity                           *float64 `json:"Equity"`
	CashAndEquivalents               *float64 `json:"CashAndEquivalents"`
	CashFlowsFromOperatingActivities *float64 `json:"CashFlowsFromOperatingActivities"`
	CashFlowsFromInvestingActivities *float64 `json:"CashFlowsFromInvestingActivities"`
	CashFlowsFromFinancingActivities *float64 `json:"CashFlowsFromFinancingActivities"`

	// 財務指標
	EarningsPerShare        *float64 `json:"EarningsPerShare"`
	DilutedEarningsPerShare *float64 `json:"DilutedEarningsPerShare"`
	BookValuePerShare       *float64 `json:"BookValuePerShare"`
	EquityToAssetRatio      *float64 `json:"EquityToAssetRatio"`

	// 配当実績
	ResultDividendPerShare1stQuarter    *float64 `json:"ResultDividendPerShare1stQuarter"`
	ResultDividendPerShare2ndQuarter    *float64 `json:"ResultDividendPerShare2ndQuarter"`
	ResultDividendPerShare3rdQuarter    *float64 `json:"ResultDividendPerShare3rdQuarter"`
	ResultDividendPerShareFiscalYearEnd *float64 `json:"ResultDividendPerShareFiscalYearEnd"`
	ResultDividendPerShareAnnual        *float64 `json:"ResultDividendPerShareAnnual"`
	DistributionsPerUnitREIT            *float64 `json:"DistributionsPerUnit(REIT)"`
	ResultTotalDividendPaidAnnual       *float64 `json:"ResultTotalDividendPaidAnnual"`
	ResultPayoutRatioAnnual             *float64 `json:"ResultPayoutRatioAnnual"`

	// 配当予想
	ForecastDividendPerShare1stQuarter    *float64 `json:"ForecastDividendPerShare1stQuarter"`
	ForecastDividendPerShare2ndQuarter    *float64 `json:"ForecastDividendPerShare2ndQuarter"`
	ForecastDividendPerShare3rdQuarter    *float64 `json:"ForecastDividendPerShare3rdQuarter"`
	ForecastDividendPerShareFiscalYearEnd *float64 `json:"ForecastDividendPerShareFiscalYearEnd"`
	ForecastDividendPerShareAnnual        *float64 `json:"ForecastDividendPerShareAnnual"`
	ForecastDistributionsPerUnitREIT      *float64 `json:"ForecastDistributionsPerUnit(REIT)"`
	ForecastTotalDividendPaidAnnual       *float64 `json:"ForecastTotalDividendPaidAnnual"`
	ForecastPayoutRatioAnnual             *float64 `json:"ForecastPayoutRatioAnnual"`

	// 当期業績予想
	ForecastNetSales         *float64 `json:"ForecastNetSales"`
	ForecastOperatingProfit  *float64 `json:"ForecastOperatingProfit"`
	ForecastOrdinaryProfit   *float64 `json:"ForecastOrdinaryProfit"`
	ForecastProfit           *float64 `json:"ForecastProfit"`
	ForecastEarningsPerShare *float64 `json:"ForecastEarningsPerShare"`

	// 第2四半期業績予想
	ForecastNetSales2ndQuarter         *float64 `json:"ForecastNetSales2ndQuarter"`
	ForecastOperatingProfit2ndQuarter  *float64 `json:"ForecastOperatingProfit2ndQuarter"`
	ForecastOrdinaryProfit2ndQuarter   *float64 `json:"ForecastOrdinaryProfit2ndQuarter"`
	ForecastProfit2ndQuarter           *float64 `json:"ForecastProfit2ndQuarter"`
	ForecastEarningsPerShare2ndQuarter *float64 `json:"ForecastEarningsPerShare2ndQuarter"`

	// 翌期業績予想
	NextYearForecastDividendPerShare1stQuarter    *float64 `json:"NextYearForecastDividendPerShare1stQuarter"`
	NextYearForecastDividendPerShare2ndQuarter    *float64 `json:"NextYearForecastDividendPerShare2ndQuarter"`
	NextYearForecastDividendPerShare3rdQuarter    *float64 `json:"NextYearForecastDividendPerShare3rdQuarter"`
	NextYearForecastDividendPerShareFiscalYearEnd *float64 `json:"NextYearForecastDividendPerShareFiscalYearEnd"`
	NextYearForecastDividendPerShareAnnual        *float64 `json:"NextYearForecastDividendPerShareAnnual"`
	NextYearForecastDistributionsPerUnitREIT       *float64 `json:"NextYearForecastDistributionsPerUnit(REIT)"`
	NextYearForecastPayoutRatioAnnual              *float64 `json:"NextYearForecastPayoutRatioAnnual"`
	NextYearForecastNetSales2ndQuarter            *float64 `json:"NextYearForecastNetSales2ndQuarter"`
	NextYearForecastOperatingProfit2ndQuarter     *float64 `json:"NextYearForecastOperatingProfit2ndQuarter"`
	NextYearForecastOrdinaryProfit2ndQuarter      *float64 `json:"NextYearForecastOrdinaryProfit2ndQuarter"`
	NextYearForecastProfit2ndQuarter              *float64 `json:"NextYearForecastProfit2ndQuarter"`
	NextYearForecastEarningsPerShare2ndQuarter    *float64 `json:"NextYearForecastEarningsPerShare2ndQuarter"`
	NextYearForecastNetSales                       *float64 `json:"NextYearForecastNetSales"`
	NextYearForecastOperatingProfit                *float64 `json:"NextYearForecastOperatingProfit"`
	NextYearForecastOrdinaryProfit                 *float64 `json:"NextYearForecastOrdinaryProfit"`
	NextYearForecastProfit                         *float64 `json:"NextYearForecastProfit"`
	NextYearForecastEarningsPerShare               *float64 `json:"NextYearForecastEarningsPerShare"`

	// その他
	MaterialChangesInSubsidiaries                                                bool   `json:"MaterialChangesInSubsidiaries"`
	SignificantChangesInTheScopeOfConsolidation                                  bool   `json:"SignificantChangesInTheScopeOfConsolidation"`
	ChangesBasedOnRevisionsOfAccountingStandard                                  bool   `json:"ChangesBasedOnRevisionsOfAccountingStandard"`
	ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard                     bool   `json:"ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard"`
	ChangesInAccountingEstimates                                                 bool   `json:"ChangesInAccountingEstimates"`
	RetrospectiveRestatement                                                     bool   `json:"RetrospectiveRestatement"`
	NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock *int64 `json:"NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock"`
	NumberOfTreasuryStockAtTheEndOfFiscalYear                                    *int64 `json:"NumberOfTreasuryStockAtTheEndOfFiscalYear"`
	AverageNumberOfShares                                                        *int64 `json:"AverageNumberOfShares"`

	// 単体財務数値
	NonConsolidatedNetSales           *float64 `json:"NonConsolidatedNetSales"`
	NonConsolidatedOperatingProfit    *float64 `json:"NonConsolidatedOperatingProfit"`
	NonConsolidatedOrdinaryProfit     *float64 `json:"NonConsolidatedOrdinaryProfit"`
	NonConsolidatedProfit             *float64 `json:"NonConsolidatedProfit"`
	NonConsolidatedEarningsPerShare   *float64 `json:"NonConsolidatedEarningsPerShare"`
	NonConsolidatedTotalAssets        *float64 `json:"NonConsolidatedTotalAssets"`
	NonConsolidatedEquity             *float64 `json:"NonConsolidatedEquity"`
	NonConsolidatedEquityToAssetRatio *float64 `json:"NonConsolidatedEquityToAssetRatio"`
	NonConsolidatedBookValuePerShare  *float64 `json:"NonConsolidatedBookValuePerShare"`

	// 単体予想（第2四半期）
	ForecastNonConsolidatedNetSales2ndQuarter         *float64 `json:"ForecastNonConsolidatedNetSales2ndQuarter"`
	ForecastNonConsolidatedOperatingProfit2ndQuarter  *float64 `json:"ForecastNonConsolidatedOperatingProfit2ndQuarter"`
	ForecastNonConsolidatedOrdinaryProfit2ndQuarter   *float64 `json:"ForecastNonConsolidatedOrdinaryProfit2ndQuarter"`
	ForecastNonConsolidatedProfit2ndQuarter           *float64 `json:"ForecastNonConsolidatedProfit2ndQuarter"`
	ForecastNonConsolidatedEarningsPerShare2ndQuarter *float64 `json:"ForecastNonConsolidatedEarningsPerShare2ndQuarter"`

	// 単体翌期予想（第2四半期）
	NextYearForecastNonConsolidatedNetSales2ndQuarter         *float64 `json:"NextYearForecastNonConsolidatedNetSales2ndQuarter"`
	NextYearForecastNonConsolidatedOperatingProfit2ndQuarter  *float64 `json:"NextYearForecastNonConsolidatedOperatingProfit2ndQuarter"`
	NextYearForecastNonConsolidatedOrdinaryProfit2ndQuarter   *float64 `json:"NextYearForecastNonConsolidatedOrdinaryProfit2ndQuarter"`
	NextYearForecastNonConsolidatedProfit2ndQuarter           *float64 `json:"NextYearForecastNonConsolidatedProfit2ndQuarter"`
	NextYearForecastNonConsolidatedEarningsPerShare2ndQuarter *float64 `json:"NextYearForecastNonConsolidatedEarningsPerShare2ndQuarter"`

	// 単体予想（期末）
	ForecastNonConsolidatedNetSales         *float64 `json:"ForecastNonConsolidatedNetSales"`
	ForecastNonConsolidatedOperatingProfit  *float64 `json:"ForecastNonConsolidatedOperatingProfit"`
	ForecastNonConsolidatedOrdinaryProfit   *float64 `json:"ForecastNonConsolidatedOrdinaryProfit"`
	ForecastNonConsolidatedProfit           *float64 `json:"ForecastNonConsolidatedProfit"`
	ForecastNonConsolidatedEarningsPerShare *float64 `json:"ForecastNonConsolidatedEarningsPerShare"`

	// 単体翌期予想（期末）
	NextYearForecastNonConsolidatedNetSales         *float64 `json:"NextYearForecastNonConsolidatedNetSales"`
	NextYearForecastNonConsolidatedOperatingProfit  *float64 `json:"NextYearForecastNonConsolidatedOperatingProfit"`
	NextYearForecastNonConsolidatedOrdinaryProfit   *float64 `json:"NextYearForecastNonConsolidatedOrdinaryProfit"`
	NextYearForecastNonConsolidatedProfit           *float64 `json:"NextYearForecastNonConsolidatedProfit"`
	NextYearForecastNonConsolidatedEarningsPerShare *float64 `json:"NextYearForecastNonConsolidatedEarningsPerShare"`
}


// RawStatement is used for unmarshaling JSON response with mixed types
type RawStatement struct {
	// 基本情報
	DisclosedDate              string `json:"DisclosedDate"`
	DisclosedTime              string `json:"DisclosedTime"`
	LocalCode                  string `json:"LocalCode"`
	DisclosureNumber           string `json:"DisclosureNumber"`
	TypeOfDocument             string `json:"TypeOfDocument"`
	TypeOfCurrentPeriod        string `json:"TypeOfCurrentPeriod"`
	CurrentPeriodStartDate     string `json:"CurrentPeriodStartDate"`
	CurrentPeriodEndDate       string `json:"CurrentPeriodEndDate"`
	CurrentFiscalYearStartDate string `json:"CurrentFiscalYearStartDate"`
	CurrentFiscalYearEndDate   string `json:"CurrentFiscalYearEndDate"`
	NextFiscalYearStartDate    string `json:"NextFiscalYearStartDate"`
	NextFiscalYearEndDate      string `json:"NextFiscalYearEndDate"`

	// 連結財務数値
	NetSales                         *types.Float64String `json:"NetSales"`
	OperatingProfit                  *types.Float64String `json:"OperatingProfit"`
	OrdinaryProfit                   *types.Float64String `json:"OrdinaryProfit"`
	Profit                           *types.Float64String `json:"Profit"`
	TotalAssets                      *types.Float64String `json:"TotalAssets"`
	Equity                           *types.Float64String `json:"Equity"`
	CashAndEquivalents               *types.Float64String `json:"CashAndEquivalents"`
	CashFlowsFromOperatingActivities *types.Float64String `json:"CashFlowsFromOperatingActivities"`
	CashFlowsFromInvestingActivities *types.Float64String `json:"CashFlowsFromInvestingActivities"`
	CashFlowsFromFinancingActivities *types.Float64String `json:"CashFlowsFromFinancingActivities"`

	// 財務指標
	EarningsPerShare        *types.Float64String `json:"EarningsPerShare"`
	DilutedEarningsPerShare *types.Float64String `json:"DilutedEarningsPerShare"`
	BookValuePerShare       *types.Float64String `json:"BookValuePerShare"`
	EquityToAssetRatio      *types.Float64String `json:"EquityToAssetRatio"`

	// 配当実績
	ResultDividendPerShare1stQuarter    *types.Float64String `json:"ResultDividendPerShare1stQuarter"`
	ResultDividendPerShare2ndQuarter    *types.Float64String `json:"ResultDividendPerShare2ndQuarter"`
	ResultDividendPerShare3rdQuarter    *types.Float64String `json:"ResultDividendPerShare3rdQuarter"`
	ResultDividendPerShareFiscalYearEnd *types.Float64String `json:"ResultDividendPerShareFiscalYearEnd"`
	ResultDividendPerShareAnnual        *types.Float64String `json:"ResultDividendPerShareAnnual"`
	DistributionsPerUnitREIT            *types.Float64String `json:"DistributionsPerUnit(REIT)"`
	ResultTotalDividendPaidAnnual       *types.Float64String `json:"ResultTotalDividendPaidAnnual"`
	ResultPayoutRatioAnnual             *types.Float64String `json:"ResultPayoutRatioAnnual"`

	// 配当予想
	ForecastDividendPerShare1stQuarter    *types.Float64String `json:"ForecastDividendPerShare1stQuarter"`
	ForecastDividendPerShare2ndQuarter    *types.Float64String `json:"ForecastDividendPerShare2ndQuarter"`
	ForecastDividendPerShare3rdQuarter    *types.Float64String `json:"ForecastDividendPerShare3rdQuarter"`
	ForecastDividendPerShareFiscalYearEnd *types.Float64String `json:"ForecastDividendPerShareFiscalYearEnd"`
	ForecastDividendPerShareAnnual        *types.Float64String `json:"ForecastDividendPerShareAnnual"`
	ForecastDistributionsPerUnitREIT      *types.Float64String `json:"ForecastDistributionsPerUnit(REIT)"`
	ForecastTotalDividendPaidAnnual       *types.Float64String `json:"ForecastTotalDividendPaidAnnual"`
	ForecastPayoutRatioAnnual             *types.Float64String `json:"ForecastPayoutRatioAnnual"`

	// 当期業績予想
	ForecastNetSales         *types.Float64String `json:"ForecastNetSales"`
	ForecastOperatingProfit  *types.Float64String `json:"ForecastOperatingProfit"`
	ForecastOrdinaryProfit   *types.Float64String `json:"ForecastOrdinaryProfit"`
	ForecastProfit           *types.Float64String `json:"ForecastProfit"`
	ForecastEarningsPerShare *types.Float64String `json:"ForecastEarningsPerShare"`

	// 第2四半期業績予想
	ForecastNetSales2ndQuarter         *types.Float64String `json:"ForecastNetSales2ndQuarter"`
	ForecastOperatingProfit2ndQuarter  *types.Float64String `json:"ForecastOperatingProfit2ndQuarter"`
	ForecastOrdinaryProfit2ndQuarter   *types.Float64String `json:"ForecastOrdinaryProfit2ndQuarter"`
	ForecastProfit2ndQuarter           *types.Float64String `json:"ForecastProfit2ndQuarter"`
	ForecastEarningsPerShare2ndQuarter *types.Float64String `json:"ForecastEarningsPerShare2ndQuarter"`

	// 翌期業績予想
	NextYearForecastDividendPerShare1stQuarter    *types.Float64String `json:"NextYearForecastDividendPerShare1stQuarter"`
	NextYearForecastDividendPerShare2ndQuarter    *types.Float64String `json:"NextYearForecastDividendPerShare2ndQuarter"`
	NextYearForecastDividendPerShare3rdQuarter    *types.Float64String `json:"NextYearForecastDividendPerShare3rdQuarter"`
	NextYearForecastDividendPerShareFiscalYearEnd *types.Float64String `json:"NextYearForecastDividendPerShareFiscalYearEnd"`
	NextYearForecastDividendPerShareAnnual        *types.Float64String `json:"NextYearForecastDividendPerShareAnnual"`
	NextYearForecastDistributionsPerUnitREIT       *types.Float64String `json:"NextYearForecastDistributionsPerUnit(REIT)"`
	NextYearForecastPayoutRatioAnnual              *types.Float64String `json:"NextYearForecastPayoutRatioAnnual"`
	NextYearForecastNetSales2ndQuarter            *types.Float64String `json:"NextYearForecastNetSales2ndQuarter"`
	NextYearForecastOperatingProfit2ndQuarter     *types.Float64String `json:"NextYearForecastOperatingProfit2ndQuarter"`
	NextYearForecastOrdinaryProfit2ndQuarter      *types.Float64String `json:"NextYearForecastOrdinaryProfit2ndQuarter"`
	NextYearForecastProfit2ndQuarter              *types.Float64String `json:"NextYearForecastProfit2ndQuarter"`
	NextYearForecastEarningsPerShare2ndQuarter    *types.Float64String `json:"NextYearForecastEarningsPerShare2ndQuarter"`
	NextYearForecastNetSales                       *types.Float64String `json:"NextYearForecastNetSales"`
	NextYearForecastOperatingProfit                *types.Float64String `json:"NextYearForecastOperatingProfit"`
	NextYearForecastOrdinaryProfit                 *types.Float64String `json:"NextYearForecastOrdinaryProfit"`
	NextYearForecastProfit                         *types.Float64String `json:"NextYearForecastProfit"`
	NextYearForecastEarningsPerShare               *types.Float64String `json:"NextYearForecastEarningsPerShare"`

	// その他
	MaterialChangesInSubsidiaries                                                types.BoolString     `json:"MaterialChangesInSubsidiaries"`
	SignificantChangesInTheScopeOfConsolidation                                  types.BoolString     `json:"SignificantChangesInTheScopeOfConsolidation"`
	ChangesBasedOnRevisionsOfAccountingStandard                                  types.BoolString     `json:"ChangesBasedOnRevisionsOfAccountingStandard"`
	ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard                     types.BoolString     `json:"ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard"`
	ChangesInAccountingEstimates                                                 types.BoolString     `json:"ChangesInAccountingEstimates"`
	RetrospectiveRestatement                                                     types.BoolString     `json:"RetrospectiveRestatement"`
	NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock *types.Float64String `json:"NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock"`
	NumberOfTreasuryStockAtTheEndOfFiscalYear                                    *types.Float64String `json:"NumberOfTreasuryStockAtTheEndOfFiscalYear"`
	AverageNumberOfShares                                                        *types.Float64String `json:"AverageNumberOfShares"`

	// 単体財務数値
	NonConsolidatedNetSales           *types.Float64String `json:"NonConsolidatedNetSales"`
	NonConsolidatedOperatingProfit    *types.Float64String `json:"NonConsolidatedOperatingProfit"`
	NonConsolidatedOrdinaryProfit     *types.Float64String `json:"NonConsolidatedOrdinaryProfit"`
	NonConsolidatedProfit             *types.Float64String `json:"NonConsolidatedProfit"`
	NonConsolidatedEarningsPerShare   *types.Float64String `json:"NonConsolidatedEarningsPerShare"`
	NonConsolidatedTotalAssets        *types.Float64String `json:"NonConsolidatedTotalAssets"`
	NonConsolidatedEquity             *types.Float64String `json:"NonConsolidatedEquity"`
	NonConsolidatedEquityToAssetRatio *types.Float64String `json:"NonConsolidatedEquityToAssetRatio"`
	NonConsolidatedBookValuePerShare  *types.Float64String `json:"NonConsolidatedBookValuePerShare"`

	// 単体予想（第2四半期）
	ForecastNonConsolidatedNetSales2ndQuarter         *types.Float64String `json:"ForecastNonConsolidatedNetSales2ndQuarter"`
	ForecastNonConsolidatedOperatingProfit2ndQuarter  *types.Float64String `json:"ForecastNonConsolidatedOperatingProfit2ndQuarter"`
	ForecastNonConsolidatedOrdinaryProfit2ndQuarter   *types.Float64String `json:"ForecastNonConsolidatedOrdinaryProfit2ndQuarter"`
	ForecastNonConsolidatedProfit2ndQuarter           *types.Float64String `json:"ForecastNonConsolidatedProfit2ndQuarter"`
	ForecastNonConsolidatedEarningsPerShare2ndQuarter *types.Float64String `json:"ForecastNonConsolidatedEarningsPerShare2ndQuarter"`

	// 単体翌期予想（第2四半期）
	NextYearForecastNonConsolidatedNetSales2ndQuarter         *types.Float64String `json:"NextYearForecastNonConsolidatedNetSales2ndQuarter"`
	NextYearForecastNonConsolidatedOperatingProfit2ndQuarter  *types.Float64String `json:"NextYearForecastNonConsolidatedOperatingProfit2ndQuarter"`
	NextYearForecastNonConsolidatedOrdinaryProfit2ndQuarter   *types.Float64String `json:"NextYearForecastNonConsolidatedOrdinaryProfit2ndQuarter"`
	NextYearForecastNonConsolidatedProfit2ndQuarter           *types.Float64String `json:"NextYearForecastNonConsolidatedProfit2ndQuarter"`
	NextYearForecastNonConsolidatedEarningsPerShare2ndQuarter *types.Float64String `json:"NextYearForecastNonConsolidatedEarningsPerShare2ndQuarter"`

	// 単体予想（期末）
	ForecastNonConsolidatedNetSales         *types.Float64String `json:"ForecastNonConsolidatedNetSales"`
	ForecastNonConsolidatedOperatingProfit  *types.Float64String `json:"ForecastNonConsolidatedOperatingProfit"`
	ForecastNonConsolidatedOrdinaryProfit   *types.Float64String `json:"ForecastNonConsolidatedOrdinaryProfit"`
	ForecastNonConsolidatedProfit           *types.Float64String `json:"ForecastNonConsolidatedProfit"`
	ForecastNonConsolidatedEarningsPerShare *types.Float64String `json:"ForecastNonConsolidatedEarningsPerShare"`

	// 単体翌期予想（期末）
	NextYearForecastNonConsolidatedNetSales         *types.Float64String `json:"NextYearForecastNonConsolidatedNetSales"`
	NextYearForecastNonConsolidatedOperatingProfit  *types.Float64String `json:"NextYearForecastNonConsolidatedOperatingProfit"`
	NextYearForecastNonConsolidatedOrdinaryProfit   *types.Float64String `json:"NextYearForecastNonConsolidatedOrdinaryProfit"`
	NextYearForecastNonConsolidatedProfit           *types.Float64String `json:"NextYearForecastNonConsolidatedProfit"`
	NextYearForecastNonConsolidatedEarningsPerShare *types.Float64String `json:"NextYearForecastNonConsolidatedEarningsPerShare"`
}

type StatementsResponse struct {
	Statements []Statement `json:"statements"`
}

// UnmarshalJSON implements custom JSON unmarshaling
func (s *StatementsResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawStatement
	type rawResponse struct {
		Statements []RawStatement `json:"statements"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Convert RawStatement to Statement
	s.Statements = make([]Statement, len(raw.Statements))
	for i, rs := range raw.Statements {
		s.Statements[i] = Statement{
			// 基本情報
			DisclosedDate:              rs.DisclosedDate,
			DisclosedTime:              rs.DisclosedTime,
			LocalCode:                  rs.LocalCode,
			DisclosureNumber:           rs.DisclosureNumber,
			TypeOfDocument:             ParseTypeOfDocument(rs.TypeOfDocument),
			TypeOfCurrentPeriod:        rs.TypeOfCurrentPeriod,
			CurrentPeriodStartDate:     rs.CurrentPeriodStartDate,
			CurrentPeriodEndDate:       rs.CurrentPeriodEndDate,
			CurrentFiscalYearStartDate: rs.CurrentFiscalYearStartDate,
			CurrentFiscalYearEndDate:   rs.CurrentFiscalYearEndDate,
			NextFiscalYearStartDate:    rs.NextFiscalYearStartDate,
			NextFiscalYearEndDate:      rs.NextFiscalYearEndDate,

			// 連結財務数値
			NetSales:                         types.ToFloat64Ptr(rs.NetSales),
			OperatingProfit:                  types.ToFloat64Ptr(rs.OperatingProfit),
			OrdinaryProfit:                   types.ToFloat64Ptr(rs.OrdinaryProfit),
			Profit:                           types.ToFloat64Ptr(rs.Profit),
			TotalAssets:                      types.ToFloat64Ptr(rs.TotalAssets),
			Equity:                           types.ToFloat64Ptr(rs.Equity),
			CashAndEquivalents:               types.ToFloat64Ptr(rs.CashAndEquivalents),
			CashFlowsFromOperatingActivities: types.ToFloat64Ptr(rs.CashFlowsFromOperatingActivities),
			CashFlowsFromInvestingActivities: types.ToFloat64Ptr(rs.CashFlowsFromInvestingActivities),
			CashFlowsFromFinancingActivities: types.ToFloat64Ptr(rs.CashFlowsFromFinancingActivities),

			// 財務指標
			EarningsPerShare:        types.ToFloat64Ptr(rs.EarningsPerShare),
			DilutedEarningsPerShare: types.ToFloat64Ptr(rs.DilutedEarningsPerShare),
			BookValuePerShare:       types.ToFloat64Ptr(rs.BookValuePerShare),
			EquityToAssetRatio:      types.ToFloat64Ptr(rs.EquityToAssetRatio),

			// 配当実績
			ResultDividendPerShare1stQuarter:    types.ToFloat64Ptr(rs.ResultDividendPerShare1stQuarter),
			ResultDividendPerShare2ndQuarter:    types.ToFloat64Ptr(rs.ResultDividendPerShare2ndQuarter),
			ResultDividendPerShare3rdQuarter:    types.ToFloat64Ptr(rs.ResultDividendPerShare3rdQuarter),
			ResultDividendPerShareFiscalYearEnd: types.ToFloat64Ptr(rs.ResultDividendPerShareFiscalYearEnd),
			ResultDividendPerShareAnnual:        types.ToFloat64Ptr(rs.ResultDividendPerShareAnnual),
			DistributionsPerUnitREIT:            types.ToFloat64Ptr(rs.DistributionsPerUnitREIT),
			ResultTotalDividendPaidAnnual:       types.ToFloat64Ptr(rs.ResultTotalDividendPaidAnnual),
			ResultPayoutRatioAnnual:             types.ToFloat64Ptr(rs.ResultPayoutRatioAnnual),

			// 配当予想
			ForecastDividendPerShare1stQuarter:    types.ToFloat64Ptr(rs.ForecastDividendPerShare1stQuarter),
			ForecastDividendPerShare2ndQuarter:    types.ToFloat64Ptr(rs.ForecastDividendPerShare2ndQuarter),
			ForecastDividendPerShare3rdQuarter:    types.ToFloat64Ptr(rs.ForecastDividendPerShare3rdQuarter),
			ForecastDividendPerShareFiscalYearEnd: types.ToFloat64Ptr(rs.ForecastDividendPerShareFiscalYearEnd),
			ForecastDividendPerShareAnnual:        types.ToFloat64Ptr(rs.ForecastDividendPerShareAnnual),
			ForecastDistributionsPerUnitREIT:      types.ToFloat64Ptr(rs.ForecastDistributionsPerUnitREIT),
			ForecastTotalDividendPaidAnnual:       types.ToFloat64Ptr(rs.ForecastTotalDividendPaidAnnual),
			ForecastPayoutRatioAnnual:             types.ToFloat64Ptr(rs.ForecastPayoutRatioAnnual),

			// 当期業績予想
			ForecastNetSales:         types.ToFloat64Ptr(rs.ForecastNetSales),
			ForecastOperatingProfit:  types.ToFloat64Ptr(rs.ForecastOperatingProfit),
			ForecastOrdinaryProfit:   types.ToFloat64Ptr(rs.ForecastOrdinaryProfit),
			ForecastProfit:           types.ToFloat64Ptr(rs.ForecastProfit),
			ForecastEarningsPerShare: types.ToFloat64Ptr(rs.ForecastEarningsPerShare),

			// 第2四半期業績予想
			ForecastNetSales2ndQuarter:         types.ToFloat64Ptr(rs.ForecastNetSales2ndQuarter),
			ForecastOperatingProfit2ndQuarter:  types.ToFloat64Ptr(rs.ForecastOperatingProfit2ndQuarter),
			ForecastOrdinaryProfit2ndQuarter:   types.ToFloat64Ptr(rs.ForecastOrdinaryProfit2ndQuarter),
			ForecastProfit2ndQuarter:           types.ToFloat64Ptr(rs.ForecastProfit2ndQuarter),
			ForecastEarningsPerShare2ndQuarter: types.ToFloat64Ptr(rs.ForecastEarningsPerShare2ndQuarter),

			// 翌期業績予想
			NextYearForecastDividendPerShare1stQuarter:    types.ToFloat64Ptr(rs.NextYearForecastDividendPerShare1stQuarter),
			NextYearForecastDividendPerShare2ndQuarter:    types.ToFloat64Ptr(rs.NextYearForecastDividendPerShare2ndQuarter),
			NextYearForecastDividendPerShare3rdQuarter:    types.ToFloat64Ptr(rs.NextYearForecastDividendPerShare3rdQuarter),
			NextYearForecastDividendPerShareFiscalYearEnd: types.ToFloat64Ptr(rs.NextYearForecastDividendPerShareFiscalYearEnd),
			NextYearForecastDividendPerShareAnnual:        types.ToFloat64Ptr(rs.NextYearForecastDividendPerShareAnnual),
			NextYearForecastDistributionsPerUnitREIT:       types.ToFloat64Ptr(rs.NextYearForecastDistributionsPerUnitREIT),
			NextYearForecastPayoutRatioAnnual:              types.ToFloat64Ptr(rs.NextYearForecastPayoutRatioAnnual),
			NextYearForecastNetSales2ndQuarter:            types.ToFloat64Ptr(rs.NextYearForecastNetSales2ndQuarter),
			NextYearForecastOperatingProfit2ndQuarter:     types.ToFloat64Ptr(rs.NextYearForecastOperatingProfit2ndQuarter),
			NextYearForecastOrdinaryProfit2ndQuarter:      types.ToFloat64Ptr(rs.NextYearForecastOrdinaryProfit2ndQuarter),
			NextYearForecastProfit2ndQuarter:              types.ToFloat64Ptr(rs.NextYearForecastProfit2ndQuarter),
			NextYearForecastEarningsPerShare2ndQuarter:    types.ToFloat64Ptr(rs.NextYearForecastEarningsPerShare2ndQuarter),
			NextYearForecastNetSales:                       types.ToFloat64Ptr(rs.NextYearForecastNetSales),
			NextYearForecastOperatingProfit:                types.ToFloat64Ptr(rs.NextYearForecastOperatingProfit),
			NextYearForecastOrdinaryProfit:                 types.ToFloat64Ptr(rs.NextYearForecastOrdinaryProfit),
			NextYearForecastProfit:                         types.ToFloat64Ptr(rs.NextYearForecastProfit),
			NextYearForecastEarningsPerShare:               types.ToFloat64Ptr(rs.NextYearForecastEarningsPerShare),

			// その他
			MaterialChangesInSubsidiaries:                                                types.ToBool(rs.MaterialChangesInSubsidiaries),
			SignificantChangesInTheScopeOfConsolidation:                                  types.ToBool(rs.SignificantChangesInTheScopeOfConsolidation),
			ChangesBasedOnRevisionsOfAccountingStandard:                                  types.ToBool(rs.ChangesBasedOnRevisionsOfAccountingStandard),
			ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard:                     types.ToBool(rs.ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard),
			ChangesInAccountingEstimates:                                                 types.ToBool(rs.ChangesInAccountingEstimates),
			RetrospectiveRestatement:                                                     types.ToBool(rs.RetrospectiveRestatement),
			NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock: types.ToInt64Ptr(rs.NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock),
			NumberOfTreasuryStockAtTheEndOfFiscalYear:                                    types.ToInt64Ptr(rs.NumberOfTreasuryStockAtTheEndOfFiscalYear),
			AverageNumberOfShares:                                                        types.ToInt64Ptr(rs.AverageNumberOfShares),

			// 単体財務数値
			NonConsolidatedNetSales:           types.ToFloat64Ptr(rs.NonConsolidatedNetSales),
			NonConsolidatedOperatingProfit:    types.ToFloat64Ptr(rs.NonConsolidatedOperatingProfit),
			NonConsolidatedOrdinaryProfit:     types.ToFloat64Ptr(rs.NonConsolidatedOrdinaryProfit),
			NonConsolidatedProfit:             types.ToFloat64Ptr(rs.NonConsolidatedProfit),
			NonConsolidatedEarningsPerShare:   types.ToFloat64Ptr(rs.NonConsolidatedEarningsPerShare),
			NonConsolidatedTotalAssets:        types.ToFloat64Ptr(rs.NonConsolidatedTotalAssets),
			NonConsolidatedEquity:             types.ToFloat64Ptr(rs.NonConsolidatedEquity),
			NonConsolidatedEquityToAssetRatio: types.ToFloat64Ptr(rs.NonConsolidatedEquityToAssetRatio),
			NonConsolidatedBookValuePerShare:  types.ToFloat64Ptr(rs.NonConsolidatedBookValuePerShare),

			// 単体予想（第2四半期）
			ForecastNonConsolidatedNetSales2ndQuarter:         types.ToFloat64Ptr(rs.ForecastNonConsolidatedNetSales2ndQuarter),
			ForecastNonConsolidatedOperatingProfit2ndQuarter:  types.ToFloat64Ptr(rs.ForecastNonConsolidatedOperatingProfit2ndQuarter),
			ForecastNonConsolidatedOrdinaryProfit2ndQuarter:   types.ToFloat64Ptr(rs.ForecastNonConsolidatedOrdinaryProfit2ndQuarter),
			ForecastNonConsolidatedProfit2ndQuarter:           types.ToFloat64Ptr(rs.ForecastNonConsolidatedProfit2ndQuarter),
			ForecastNonConsolidatedEarningsPerShare2ndQuarter: types.ToFloat64Ptr(rs.ForecastNonConsolidatedEarningsPerShare2ndQuarter),

			// 単体翌期予想（第2四半期）
			NextYearForecastNonConsolidatedNetSales2ndQuarter:         types.ToFloat64Ptr(rs.NextYearForecastNonConsolidatedNetSales2ndQuarter),
			NextYearForecastNonConsolidatedOperatingProfit2ndQuarter:  types.ToFloat64Ptr(rs.NextYearForecastNonConsolidatedOperatingProfit2ndQuarter),
			NextYearForecastNonConsolidatedOrdinaryProfit2ndQuarter:   types.ToFloat64Ptr(rs.NextYearForecastNonConsolidatedOrdinaryProfit2ndQuarter),
			NextYearForecastNonConsolidatedProfit2ndQuarter:           types.ToFloat64Ptr(rs.NextYearForecastNonConsolidatedProfit2ndQuarter),
			NextYearForecastNonConsolidatedEarningsPerShare2ndQuarter: types.ToFloat64Ptr(rs.NextYearForecastNonConsolidatedEarningsPerShare2ndQuarter),

			// 単体予想（期末）
			ForecastNonConsolidatedNetSales:         types.ToFloat64Ptr(rs.ForecastNonConsolidatedNetSales),
			ForecastNonConsolidatedOperatingProfit:  types.ToFloat64Ptr(rs.ForecastNonConsolidatedOperatingProfit),
			ForecastNonConsolidatedOrdinaryProfit:   types.ToFloat64Ptr(rs.ForecastNonConsolidatedOrdinaryProfit),
			ForecastNonConsolidatedProfit:           types.ToFloat64Ptr(rs.ForecastNonConsolidatedProfit),
			ForecastNonConsolidatedEarningsPerShare: types.ToFloat64Ptr(rs.ForecastNonConsolidatedEarningsPerShare),

			// 単体翌期予想（期末）
			NextYearForecastNonConsolidatedNetSales:         types.ToFloat64Ptr(rs.NextYearForecastNonConsolidatedNetSales),
			NextYearForecastNonConsolidatedOperatingProfit:  types.ToFloat64Ptr(rs.NextYearForecastNonConsolidatedOperatingProfit),
			NextYearForecastNonConsolidatedOrdinaryProfit:   types.ToFloat64Ptr(rs.NextYearForecastNonConsolidatedOrdinaryProfit),
			NextYearForecastNonConsolidatedProfit:           types.ToFloat64Ptr(rs.NextYearForecastNonConsolidatedProfit),
			NextYearForecastNonConsolidatedEarningsPerShare: types.ToFloat64Ptr(rs.NextYearForecastNonConsolidatedEarningsPerShare),
		}
	}

	return nil
}

// GetStatements は財務諸表データを取得します。
// パラメータ:
// - code: 銘柄コード
// - date: 基準日（YYYYMMDD形式、空の場合は最新）
func (s *StatementsService) GetStatements(code string, date string) ([]Statement, error) {
	path := "/fins/statements"

	query := "?"
	if code != "" {
		query += fmt.Sprintf("code=%s&", code)
	}
	if date != "" {
		query += fmt.Sprintf("date=%s&", date)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp StatementsResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get statements: %w", err)
	}

	return resp.Statements, nil
}

// GetLatestStatements は指定銘柄の最新財務諸表を取得します。
// 例: GetLatestStatements("7203") でトヨタ自動車の最新決算データを取得
func (s *StatementsService) GetLatestStatements(code string) (*Statement, error) {
	// 最新日付を指定して取得
	statements, err := s.GetStatements(code, "")
	if err != nil {
		return nil, err
	}

	if len(statements) == 0 {
		return nil, fmt.Errorf("no statements found for code: %s", code)
	}

	// 最新のものを探す（DisclosedDateでソート ）
	latestStmt := statements[0]
	for _, stmt := range statements {
		if stmt.DisclosedDate > latestStmt.DisclosedDate {
			latestStmt = stmt
		}
	}

	return &latestStmt, nil
}

// GetStatementsByDate は指定日の財務諸表データを取得します。
// 例: GetStatementsByDate("2024-01-15") で特定日に開示された全銘柄の決算データを取得
func (s *StatementsService) GetStatementsByDate(date string) ([]Statement, error) {
	return s.GetStatements("", date)
}
