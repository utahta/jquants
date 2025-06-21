package jquants

import (
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// StatementsService は財務諸表データを取得するサービスです。
// 売上高、利益、資産、ROE/ROAなどの財務指標を提供します。
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
	TypeOfDocument             string `json:"TypeOfDocument"`
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
	ResultPayoutRatioAnnual             *float64 `json:"ResultPayoutRatioAnnual"`

	// 配当予想
	ForecastDividendPerShare1stQuarter    *float64 `json:"ForecastDividendPerShare1stQuarter"`
	ForecastDividendPerShare2ndQuarter    *float64 `json:"ForecastDividendPerShare2ndQuarter"`
	ForecastDividendPerShare3rdQuarter    *float64 `json:"ForecastDividendPerShare3rdQuarter"`
	ForecastDividendPerShareFiscalYearEnd *float64 `json:"ForecastDividendPerShareFiscalYearEnd"`
	ForecastDividendPerShareAnnual        *float64 `json:"ForecastDividendPerShareAnnual"`
	ForecastPayoutRatioAnnual             *float64 `json:"ForecastPayoutRatioAnnual"`
	ForecastDistributionsPerUnitREIT      *float64 `json:"ForecastDistributionsPerUnit(REIT)"`

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
	NextYearForecastDividendPerShareAnnual *float64 `json:"NextYearForecastDividendPerShareAnnual"`
	NextYearForecastPayoutRatioAnnual      *float64 `json:"NextYearForecastPayoutRatioAnnual"`
	NextYearForecastNetSales               *float64 `json:"NextYearForecastNetSales"`
	NextYearForecastOperatingProfit        *float64 `json:"NextYearForecastOperatingProfit"`
	NextYearForecastOrdinaryProfit         *float64 `json:"NextYearForecastOrdinaryProfit"`
	NextYearForecastProfit                 *float64 `json:"NextYearForecastProfit"`
	NextYearForecastEarningsPerShare       *float64 `json:"NextYearForecastEarningsPerShare"`

	// その他
	MaterialChangesInSubsidiaries                                                bool   `json:"MaterialChangesInSubsidiaries"`
	SignificantChangesInTheScopeOfConsolidation                                  bool   `json:"SignificantChangesInTheScopeOfConsolidation"`
	ChangesInAccountingEstimates                                                 bool   `json:"ChangesInAccountingEstimates"`
	ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard                     bool   `json:"ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard"`
	RetrospectiveRestatement                                                     bool   `json:"RetrospectiveRestatement"`
	NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock *int64 `json:"NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock"`
	NumberOfTreasuryStockAtTheEndOfFiscalYear                                    *int64 `json:"NumberOfTreasuryStockAtTheEndOfFiscalYear"`

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
	ResultPayoutRatioAnnual             *types.Float64String `json:"ResultPayoutRatioAnnual"`

	// 配当予想
	ForecastDividendPerShare1stQuarter    *types.Float64String `json:"ForecastDividendPerShare1stQuarter"`
	ForecastDividendPerShare2ndQuarter    *types.Float64String `json:"ForecastDividendPerShare2ndQuarter"`
	ForecastDividendPerShare3rdQuarter    *types.Float64String `json:"ForecastDividendPerShare3rdQuarter"`
	ForecastDividendPerShareFiscalYearEnd *types.Float64String `json:"ForecastDividendPerShareFiscalYearEnd"`
	ForecastDividendPerShareAnnual        *types.Float64String `json:"ForecastDividendPerShareAnnual"`
	ForecastPayoutRatioAnnual             *types.Float64String `json:"ForecastPayoutRatioAnnual"`
	ForecastDistributionsPerUnitREIT      *types.Float64String `json:"ForecastDistributionsPerUnit(REIT)"`

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
	NextYearForecastDividendPerShareAnnual *types.Float64String `json:"NextYearForecastDividendPerShareAnnual"`
	NextYearForecastPayoutRatioAnnual      *types.Float64String `json:"NextYearForecastPayoutRatioAnnual"`
	NextYearForecastNetSales               *types.Float64String `json:"NextYearForecastNetSales"`
	NextYearForecastOperatingProfit        *types.Float64String `json:"NextYearForecastOperatingProfit"`
	NextYearForecastOrdinaryProfit         *types.Float64String `json:"NextYearForecastOrdinaryProfit"`
	NextYearForecastProfit                 *types.Float64String `json:"NextYearForecastProfit"`
	NextYearForecastEarningsPerShare       *types.Float64String `json:"NextYearForecastEarningsPerShare"`

	// その他
	MaterialChangesInSubsidiaries                                                types.BoolString     `json:"MaterialChangesInSubsidiaries"`
	SignificantChangesInTheScopeOfConsolidation                                  types.BoolString     `json:"SignificantChangesInTheScopeOfConsolidation"`
	ChangesInAccountingEstimates                                                 types.BoolString     `json:"ChangesInAccountingEstimates"`
	ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard                     types.BoolString     `json:"ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard"`
	RetrospectiveRestatement                                                     types.BoolString     `json:"RetrospectiveRestatement"`
	NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock *types.Float64String `json:"NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock"`
	NumberOfTreasuryStockAtTheEndOfFiscalYear                                    *types.Float64String `json:"NumberOfTreasuryStockAtTheEndOfFiscalYear"`

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
			TypeOfDocument:             rs.TypeOfDocument,
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
			ResultPayoutRatioAnnual:             types.ToFloat64Ptr(rs.ResultPayoutRatioAnnual),

			// 配当予想
			ForecastDividendPerShare1stQuarter:    types.ToFloat64Ptr(rs.ForecastDividendPerShare1stQuarter),
			ForecastDividendPerShare2ndQuarter:    types.ToFloat64Ptr(rs.ForecastDividendPerShare2ndQuarter),
			ForecastDividendPerShare3rdQuarter:    types.ToFloat64Ptr(rs.ForecastDividendPerShare3rdQuarter),
			ForecastDividendPerShareFiscalYearEnd: types.ToFloat64Ptr(rs.ForecastDividendPerShareFiscalYearEnd),
			ForecastDividendPerShareAnnual:        types.ToFloat64Ptr(rs.ForecastDividendPerShareAnnual),
			ForecastPayoutRatioAnnual:             types.ToFloat64Ptr(rs.ForecastPayoutRatioAnnual),
			ForecastDistributionsPerUnitREIT:      types.ToFloat64Ptr(rs.ForecastDistributionsPerUnitREIT),

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
			NextYearForecastDividendPerShareAnnual: types.ToFloat64Ptr(rs.NextYearForecastDividendPerShareAnnual),
			NextYearForecastPayoutRatioAnnual:      types.ToFloat64Ptr(rs.NextYearForecastPayoutRatioAnnual),
			NextYearForecastNetSales:               types.ToFloat64Ptr(rs.NextYearForecastNetSales),
			NextYearForecastOperatingProfit:        types.ToFloat64Ptr(rs.NextYearForecastOperatingProfit),
			NextYearForecastOrdinaryProfit:         types.ToFloat64Ptr(rs.NextYearForecastOrdinaryProfit),
			NextYearForecastProfit:                 types.ToFloat64Ptr(rs.NextYearForecastProfit),
			NextYearForecastEarningsPerShare:       types.ToFloat64Ptr(rs.NextYearForecastEarningsPerShare),

			// その他
			MaterialChangesInSubsidiaries:                                                types.ToBool(rs.MaterialChangesInSubsidiaries),
			SignificantChangesInTheScopeOfConsolidation:                                  types.ToBool(rs.SignificantChangesInTheScopeOfConsolidation),
			ChangesInAccountingEstimates:                                                 types.ToBool(rs.ChangesInAccountingEstimates),
			ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard:                     types.ToBool(rs.ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard),
			RetrospectiveRestatement:                                                     types.ToBool(rs.RetrospectiveRestatement),
			NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock: types.ToInt64Ptr(rs.NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock),
			NumberOfTreasuryStockAtTheEndOfFiscalYear:                                    types.ToInt64Ptr(rs.NumberOfTreasuryStockAtTheEndOfFiscalYear),

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
