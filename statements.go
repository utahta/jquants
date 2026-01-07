package jquants

import (
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// StatementsService は財務諸表データを取得するサービスです。
// 売上高、利益、資産、ROE/ROAなどの財務指標を提供します。
// 財務情報API(/fins/summary)から四半期毎の決算短信サマリーデータを取得します。
type StatementsService struct {
	client client.HTTPClient
}

func NewStatementsService(c client.HTTPClient) *StatementsService {
	return &StatementsService{
		client: c,
	}
}

// Statement は財務諸表データを表します。
// J-Quants API /fins/summary エンドポイントのレスポンスデータ。
// すべてのフィールドはAPIでは文字列型で返されますが、このstructでは適切な型に変換されています。
type Statement struct {
	// 基本情報
	DiscDate   string         `json:"DiscDate"`   // 開示日
	DiscTime   string         `json:"DiscTime"`   // 開示時刻
	Code       string         `json:"Code"`       // 銘柄コード（5桁）
	DiscNo     string         `json:"DiscNo"`     // 開示番号
	DocType    TypeOfDocument `json:"DocType"`    // 開示書類種別
	CurPerType string         `json:"CurPerType"` // 当会計期間の種類 [1Q, 2Q, 3Q, 4Q, 5Q, FY]
	CurPerSt   string         `json:"CurPerSt"`   // 当会計期間開始日
	CurPerEn   string         `json:"CurPerEn"`   // 当会計期間終了日
	CurFYSt    string         `json:"CurFYSt"`    // 当事業年度開始日
	CurFYEn    string         `json:"CurFYEn"`    // 当事業年度終了日
	NxtFYSt    string         `json:"NxtFYSt"`    // 翌事業年度開始日
	NxtFYEn    string         `json:"NxtFYEn"`    // 翌事業年度終了日

	// 連結財務数値
	Sales  *float64 `json:"Sales"`  // 売上高
	OP     *float64 `json:"OP"`     // 営業利益
	OdP    *float64 `json:"OdP"`    // 経常利益
	NP     *float64 `json:"NP"`     // 当期純利益
	TA     *float64 `json:"TA"`     // 総資産
	Eq     *float64 `json:"Eq"`     // 純資産
	CashEq *float64 `json:"CashEq"` // 現金及び現金同等物期末残高
	CFO    *float64 `json:"CFO"`    // 営業活動によるキャッシュ・フロー
	CFI    *float64 `json:"CFI"`    // 投資活動によるキャッシュ・フロー
	CFF    *float64 `json:"CFF"`    // 財務活動によるキャッシュ・フロー

	// 財務指標
	EPS  *float64 `json:"EPS"`  // 一株あたり当期純利益
	DEPS *float64 `json:"DEPS"` // 潜在株式調整後一株あたり当期純利益
	BPS  *float64 `json:"BPS"`  // 一株あたり純資産
	EqAR *float64 `json:"EqAR"` // 自己資本比率

	// 配当実績
	Div1Q         *float64 `json:"Div1Q"`         // 一株あたり配当実績_第1四半期末
	Div2Q         *float64 `json:"Div2Q"`         // 一株あたり配当実績_第2四半期末
	Div3Q         *float64 `json:"Div3Q"`         // 一株あたり配当実績_第3四半期末
	DivFY         *float64 `json:"DivFY"`         // 一株あたり配当実績_期末
	DivAnn        *float64 `json:"DivAnn"`        // 一株あたり配当実績_合計
	DivUnit       *float64 `json:"DivUnit"`       // 1口当たり分配金
	DivTotalAnn   *float64 `json:"DivTotalAnn"`   // 配当金総額
	PayoutRatioAn *float64 `json:"PayoutRatioAn"` // 配当性向

	// 配当予想
	FDiv1Q         *float64 `json:"FDiv1Q"`         // 一株あたり配当予想_第1四半期末
	FDiv2Q         *float64 `json:"FDiv2Q"`         // 一株あたり配当予想_第2四半期末
	FDiv3Q         *float64 `json:"FDiv3Q"`         // 一株あたり配当予想_第3四半期末
	FDivFY         *float64 `json:"FDivFY"`         // 一株あたり配当予想_期末
	FDivAnn        *float64 `json:"FDivAnn"`        // 一株あたり配当予想_合計
	FDivUnit       *float64 `json:"FDivUnit"`       // 1口当たり予想分配金
	FDivTotalAnn   *float64 `json:"FDivTotalAnn"`   // 予想配当金総額
	FPayoutRatioAn *float64 `json:"FPayoutRatioAn"` // 予想配当性向

	// 翌期配当予想
	NxFDiv1Q         *float64 `json:"NxFDiv1Q"`         // 一株あたり配当予想_翌事業年度第1四半期末
	NxFDiv2Q         *float64 `json:"NxFDiv2Q"`         // 一株あたり配当予想_翌事業年度第2四半期末
	NxFDiv3Q         *float64 `json:"NxFDiv3Q"`         // 一株あたり配当予想_翌事業年度第3四半期末
	NxFDivFY         *float64 `json:"NxFDivFY"`         // 一株あたり配当予想_翌事業年度期末
	NxFDivAnn        *float64 `json:"NxFDivAnn"`        // 一株あたり配当予想_翌事業年度合計
	NxFDivUnit       *float64 `json:"NxFDivUnit"`       // 1口当たり翌事業年度予想分配金
	NxFPayoutRatioAn *float64 `json:"NxFPayoutRatioAn"` // 翌事業年度予想配当性向

	// 第2四半期業績予想
	FSales2Q *float64 `json:"FSales2Q"` // 売上高_予想_第2四半期末
	FOP2Q    *float64 `json:"FOP2Q"`    // 営業利益_予想_第2四半期末
	FOdP2Q   *float64 `json:"FOdP2Q"`   // 経常利益_予想_第2四半期末
	FNP2Q    *float64 `json:"FNP2Q"`    // 当期純利益_予想_第2四半期末
	FEPS2Q   *float64 `json:"FEPS2Q"`   // 一株あたり当期純利益_予想_第2四半期末

	// 翌期第2四半期業績予想
	NxFSales2Q *float64 `json:"NxFSales2Q"` // 売上高_予想_翌事業年度第2四半期末
	NxFOP2Q    *float64 `json:"NxFOP2Q"`    // 営業利益_予想_翌事業年度第2四半期末
	NxFOdP2Q   *float64 `json:"NxFOdP2Q"`   // 経常利益_予想_翌事業年度第2四半期末
	NxFNp2Q    *float64 `json:"NxFNp2Q"`    // 当期純利益_予想_翌事業年度第2四半期末
	NxFEPS2Q   *float64 `json:"NxFEPS2Q"`   // 一株あたり当期純利益_予想_翌事業年度第2四半期末

	// 期末業績予想
	FSales *float64 `json:"FSales"` // 売上高_予想_期末
	FOP    *float64 `json:"FOP"`    // 営業利益_予想_期末
	FOdP   *float64 `json:"FOdP"`   // 経常利益_予想_期末
	FNP    *float64 `json:"FNP"`    // 当期純利益_予想_期末
	FEPS   *float64 `json:"FEPS"`   // 一株あたり当期純利益_予想_期末

	// 翌期末業績予想
	NxFSales *float64 `json:"NxFSales"` // 売上高_予想_翌事業年度期末
	NxFOP    *float64 `json:"NxFOP"`    // 営業利益_予想_翌事業年度期末
	NxFOdP   *float64 `json:"NxFOdP"`   // 経常利益_予想_翌事業年度期末
	NxFNp    *float64 `json:"NxFNp"`    // 当期純利益_予想_翌事業年度期末
	NxFEPS   *float64 `json:"NxFEPS"`   // 一株あたり当期純利益_予想_翌事業年度期末

	// その他
	MatChgSub  string `json:"MatChgSub"`  // 期中における重要な子会社の異動
	SigChgInC  string `json:"SigChgInC"`  // 期中における連結範囲の重要な変更
	ChgByASRev string `json:"ChgByASRev"` // 会計基準等の改正に伴う会計方針の変更
	ChgNoASRev string `json:"ChgNoASRev"` // 会計基準等の改正に伴う変更以外の会計方針の変更
	ChgAcEst   string `json:"ChgAcEst"`   // 会計上の見積りの変更
	RetroRst   string `json:"RetroRst"`   // 修正再表示
	ShOutFY    *int64 `json:"ShOutFY"`    // 期末発行済株式数
	TrShFY     *int64 `json:"TrShFY"`     // 期末自己株式数
	AvgSh      *int64 `json:"AvgSh"`      // 期中平均株式数

	// 非連結財務数値
	NCSales *float64 `json:"NCSales"` // 売上高_非連結
	NCOP    *float64 `json:"NCOP"`    // 営業利益_非連結
	NCOdP   *float64 `json:"NCOdP"`   // 経常利益_非連結
	NCNP    *float64 `json:"NCNP"`    // 当期純利益_非連結
	NCEPS   *float64 `json:"NCEPS"`   // 一株あたり当期純利益_非連結
	NCTA    *float64 `json:"NCTA"`    // 総資産_非連結
	NCEq    *float64 `json:"NCEq"`    // 純資産_非連結
	NCEqAR  *float64 `json:"NCEqAR"`  // 自己資本比率_非連結
	NCBPS   *float64 `json:"NCBPS"`   // 一株あたり純資産_非連結

	// 非連結第2四半期予想
	FNCSales2Q *float64 `json:"FNCSales2Q"` // 売上高_予想_第2四半期末_非連結
	FNCOP2Q    *float64 `json:"FNCOP2Q"`    // 営業利益_予想_第2四半期末_非連結
	FNCOdP2Q   *float64 `json:"FNCOdP2Q"`   // 経常利益_予想_第2四半期末_非連結
	FNCNP2Q    *float64 `json:"FNCNP2Q"`    // 当期純利益_予想_第2四半期末_非連結
	FNCEPS2Q   *float64 `json:"FNCEPS2Q"`   // 一株あたり当期純利益_予想_第2四半期末_非連結

	// 翌期非連結第2四半期予想
	NxFNCSales2Q *float64 `json:"NxFNCSales2Q"` // 売上高_予想_翌事業年度第2四半期末_非連結
	NxFNCOP2Q    *float64 `json:"NxFNCOP2Q"`    // 営業利益_予想_翌事業年度第2四半期末_非連結
	NxFNCOdP2Q   *float64 `json:"NxFNCOdP2Q"`   // 経常利益_予想_翌事業年度第2四半期末_非連結
	NxFNCNP2Q    *float64 `json:"NxFNCNP2Q"`    // 当期純利益_予想_翌事業年度第2四半期末_非連結
	NxFNCEPS2Q   *float64 `json:"NxFNCEPS2Q"`   // 一株あたり当期純利益_予想_翌事業年度第2四半期末_非連結

	// 非連結期末予想
	FNCSales *float64 `json:"FNCSales"` // 売上高_予想_期末_非連結
	FNCOP    *float64 `json:"FNCOP"`    // 営業利益_予想_期末_非連結
	FNCOdP   *float64 `json:"FNCOdP"`   // 経常利益_予想_期末_非連結
	FNCNP    *float64 `json:"FNCNP"`    // 当期純利益_予想_期末_非連結
	FNCEPS   *float64 `json:"FNCEPS"`   // 一株あたり当期純利益_予想_期末_非連結

	// 翌期非連結期末予想
	NxFNCSales *float64 `json:"NxFNCSales"` // 売上高_予想_翌事業年度期末_非連結
	NxFNCOP    *float64 `json:"NxFNCOP"`    // 営業利益_予想_翌事業年度期末_非連結
	NxFNCOdP   *float64 `json:"NxFNCOdP"`   // 経常利益_予想_翌事業年度期末_非連結
	NxFNCNP    *float64 `json:"NxFNCNP"`    // 当期純利益_予想_翌事業年度期末_非連結
	NxFNCEPS   *float64 `json:"NxFNCEPS"`   // 一株あたり当期純利益_予想_翌事業年度期末_非連結
}

// RawStatement is used for unmarshaling JSON response with mixed types
type RawStatement struct {
	// 基本情報
	DiscDate   string `json:"DiscDate"`
	DiscTime   string `json:"DiscTime"`
	Code       string `json:"Code"`
	DiscNo     string `json:"DiscNo"`
	DocType    string `json:"DocType"`
	CurPerType string `json:"CurPerType"`
	CurPerSt   string `json:"CurPerSt"`
	CurPerEn   string `json:"CurPerEn"`
	CurFYSt    string `json:"CurFYSt"`
	CurFYEn    string `json:"CurFYEn"`
	NxtFYSt    string `json:"NxtFYSt"`
	NxtFYEn    string `json:"NxtFYEn"`

	// 連結財務数値
	Sales  *types.Float64String `json:"Sales"`
	OP     *types.Float64String `json:"OP"`
	OdP    *types.Float64String `json:"OdP"`
	NP     *types.Float64String `json:"NP"`
	TA     *types.Float64String `json:"TA"`
	Eq     *types.Float64String `json:"Eq"`
	CashEq *types.Float64String `json:"CashEq"`
	CFO    *types.Float64String `json:"CFO"`
	CFI    *types.Float64String `json:"CFI"`
	CFF    *types.Float64String `json:"CFF"`

	// 財務指標
	EPS  *types.Float64String `json:"EPS"`
	DEPS *types.Float64String `json:"DEPS"`
	BPS  *types.Float64String `json:"BPS"`
	EqAR *types.Float64String `json:"EqAR"`

	// 配当実績
	Div1Q         *types.Float64String `json:"Div1Q"`
	Div2Q         *types.Float64String `json:"Div2Q"`
	Div3Q         *types.Float64String `json:"Div3Q"`
	DivFY         *types.Float64String `json:"DivFY"`
	DivAnn        *types.Float64String `json:"DivAnn"`
	DivUnit       *types.Float64String `json:"DivUnit"`
	DivTotalAnn   *types.Float64String `json:"DivTotalAnn"`
	PayoutRatioAn *types.Float64String `json:"PayoutRatioAnn"`

	// 配当予想
	FDiv1Q         *types.Float64String `json:"FDiv1Q"`
	FDiv2Q         *types.Float64String `json:"FDiv2Q"`
	FDiv3Q         *types.Float64String `json:"FDiv3Q"`
	FDivFY         *types.Float64String `json:"FDivFY"`
	FDivAnn        *types.Float64String `json:"FDivAnn"`
	FDivUnit       *types.Float64String `json:"FDivUnit"`
	FDivTotalAnn   *types.Float64String `json:"FDivTotalAnn"`
	FPayoutRatioAn *types.Float64String `json:"FPayoutRatioAnn"`

	// 翌期配当予想
	NxFDiv1Q         *types.Float64String `json:"NxFDiv1Q"`
	NxFDiv2Q         *types.Float64String `json:"NxFDiv2Q"`
	NxFDiv3Q         *types.Float64String `json:"NxFDiv3Q"`
	NxFDivFY         *types.Float64String `json:"NxFDivFY"`
	NxFDivAnn        *types.Float64String `json:"NxFDivAnn"`
	NxFDivUnit       *types.Float64String `json:"NxFDivUnit"`
	NxFPayoutRatioAn *types.Float64String `json:"NxFPayoutRatioAnn"`

	// 第2四半期業績予想
	FSales2Q *types.Float64String `json:"FSales2Q"`
	FOP2Q    *types.Float64String `json:"FOP2Q"`
	FOdP2Q   *types.Float64String `json:"FOdP2Q"`
	FNP2Q    *types.Float64String `json:"FNP2Q"`
	FEPS2Q   *types.Float64String `json:"FEPS2Q"`

	// 翌期第2四半期業績予想
	NxFSales2Q *types.Float64String `json:"NxFSales2Q"`
	NxFOP2Q    *types.Float64String `json:"NxFOP2Q"`
	NxFOdP2Q   *types.Float64String `json:"NxFOdP2Q"`
	NxFNp2Q    *types.Float64String `json:"NxFNp2Q"`
	NxFEPS2Q   *types.Float64String `json:"NxFEPS2Q"`

	// 期末業績予想
	FSales *types.Float64String `json:"FSales"`
	FOP    *types.Float64String `json:"FOP"`
	FOdP   *types.Float64String `json:"FOdP"`
	FNP    *types.Float64String `json:"FNP"`
	FEPS   *types.Float64String `json:"FEPS"`

	// 翌期末業績予想
	NxFSales *types.Float64String `json:"NxFSales"`
	NxFOP    *types.Float64String `json:"NxFOP"`
	NxFOdP   *types.Float64String `json:"NxFOdP"`
	NxFNp    *types.Float64String `json:"NxFNp"`
	NxFEPS   *types.Float64String `json:"NxFEPS"`

	// その他
	MatChgSub  string               `json:"MatChgSub"`
	SigChgInC  string               `json:"SigChgInC"`
	ChgByASRev string               `json:"ChgByASRev"`
	ChgNoASRev string               `json:"ChgNoASRev"`
	ChgAcEst   string               `json:"ChgAcEst"`
	RetroRst   string               `json:"RetroRst"`
	ShOutFY    *types.Float64String `json:"ShOutFY"`
	TrShFY     *types.Float64String `json:"TrShFY"`
	AvgSh      *types.Float64String `json:"AvgSh"`

	// 非連結財務数値
	NCSales *types.Float64String `json:"NCSales"`
	NCOP    *types.Float64String `json:"NCOP"`
	NCOdP   *types.Float64String `json:"NCOdP"`
	NCNP    *types.Float64String `json:"NCNP"`
	NCEPS   *types.Float64String `json:"NCEPS"`
	NCTA    *types.Float64String `json:"NCTA"`
	NCEq    *types.Float64String `json:"NCEq"`
	NCEqAR  *types.Float64String `json:"NCEqAR"`
	NCBPS   *types.Float64String `json:"NCBPS"`

	// 非連結第2四半期予想
	FNCSales2Q *types.Float64String `json:"FNCSales2Q"`
	FNCOP2Q    *types.Float64String `json:"FNCOP2Q"`
	FNCOdP2Q   *types.Float64String `json:"FNCOdP2Q"`
	FNCNP2Q    *types.Float64String `json:"FNCNP2Q"`
	FNCEPS2Q   *types.Float64String `json:"FNCEPS2Q"`

	// 翌期非連結第2四半期予想
	NxFNCSales2Q *types.Float64String `json:"NxFNCSales2Q"`
	NxFNCOP2Q    *types.Float64String `json:"NxFNCOP2Q"`
	NxFNCOdP2Q   *types.Float64String `json:"NxFNCOdP2Q"`
	NxFNCNP2Q    *types.Float64String `json:"NxFNCNP2Q"`
	NxFNCEPS2Q   *types.Float64String `json:"NxFNCEPS2Q"`

	// 非連結期末予想
	FNCSales *types.Float64String `json:"FNCSales"`
	FNCOP    *types.Float64String `json:"FNCOP"`
	FNCOdP   *types.Float64String `json:"FNCOdP"`
	FNCNP    *types.Float64String `json:"FNCNP"`
	FNCEPS   *types.Float64String `json:"FNCEPS"`

	// 翌期非連結期末予想
	NxFNCSales *types.Float64String `json:"NxFNCSales"`
	NxFNCOP    *types.Float64String `json:"NxFNCOP"`
	NxFNCOdP   *types.Float64String `json:"NxFNCOdP"`
	NxFNCNP    *types.Float64String `json:"NxFNCNP"`
	NxFNCEPS   *types.Float64String `json:"NxFNCEPS"`
}

type StatementsResponse struct {
	Data          []Statement `json:"data"`
	PaginationKey string      `json:"pagination_key"`
}

// UnmarshalJSON implements custom JSON unmarshaling
func (s *StatementsResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawStatement
	type rawResponse struct {
		Data          []RawStatement `json:"data"`
		PaginationKey string         `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	s.PaginationKey = raw.PaginationKey

	// Convert RawStatement to Statement
	s.Data = make([]Statement, len(raw.Data))
	for i, rs := range raw.Data {
		s.Data[i] = Statement{
			// 基本情報
			DiscDate:   rs.DiscDate,
			DiscTime:   rs.DiscTime,
			Code:       rs.Code,
			DiscNo:     rs.DiscNo,
			DocType:    ParseTypeOfDocument(rs.DocType),
			CurPerType: rs.CurPerType,
			CurPerSt:   rs.CurPerSt,
			CurPerEn:   rs.CurPerEn,
			CurFYSt:    rs.CurFYSt,
			CurFYEn:    rs.CurFYEn,
			NxtFYSt:    rs.NxtFYSt,
			NxtFYEn:    rs.NxtFYEn,

			// 連結財務数値
			Sales:  types.ToFloat64Ptr(rs.Sales),
			OP:     types.ToFloat64Ptr(rs.OP),
			OdP:    types.ToFloat64Ptr(rs.OdP),
			NP:     types.ToFloat64Ptr(rs.NP),
			TA:     types.ToFloat64Ptr(rs.TA),
			Eq:     types.ToFloat64Ptr(rs.Eq),
			CashEq: types.ToFloat64Ptr(rs.CashEq),
			CFO:    types.ToFloat64Ptr(rs.CFO),
			CFI:    types.ToFloat64Ptr(rs.CFI),
			CFF:    types.ToFloat64Ptr(rs.CFF),

			// 財務指標
			EPS:  types.ToFloat64Ptr(rs.EPS),
			DEPS: types.ToFloat64Ptr(rs.DEPS),
			BPS:  types.ToFloat64Ptr(rs.BPS),
			EqAR: types.ToFloat64Ptr(rs.EqAR),

			// 配当実績
			Div1Q:         types.ToFloat64Ptr(rs.Div1Q),
			Div2Q:         types.ToFloat64Ptr(rs.Div2Q),
			Div3Q:         types.ToFloat64Ptr(rs.Div3Q),
			DivFY:         types.ToFloat64Ptr(rs.DivFY),
			DivAnn:        types.ToFloat64Ptr(rs.DivAnn),
			DivUnit:       types.ToFloat64Ptr(rs.DivUnit),
			DivTotalAnn:   types.ToFloat64Ptr(rs.DivTotalAnn),
			PayoutRatioAn: types.ToFloat64Ptr(rs.PayoutRatioAn),

			// 配当予想
			FDiv1Q:         types.ToFloat64Ptr(rs.FDiv1Q),
			FDiv2Q:         types.ToFloat64Ptr(rs.FDiv2Q),
			FDiv3Q:         types.ToFloat64Ptr(rs.FDiv3Q),
			FDivFY:         types.ToFloat64Ptr(rs.FDivFY),
			FDivAnn:        types.ToFloat64Ptr(rs.FDivAnn),
			FDivUnit:       types.ToFloat64Ptr(rs.FDivUnit),
			FDivTotalAnn:   types.ToFloat64Ptr(rs.FDivTotalAnn),
			FPayoutRatioAn: types.ToFloat64Ptr(rs.FPayoutRatioAn),

			// 翌期配当予想
			NxFDiv1Q:         types.ToFloat64Ptr(rs.NxFDiv1Q),
			NxFDiv2Q:         types.ToFloat64Ptr(rs.NxFDiv2Q),
			NxFDiv3Q:         types.ToFloat64Ptr(rs.NxFDiv3Q),
			NxFDivFY:         types.ToFloat64Ptr(rs.NxFDivFY),
			NxFDivAnn:        types.ToFloat64Ptr(rs.NxFDivAnn),
			NxFDivUnit:       types.ToFloat64Ptr(rs.NxFDivUnit),
			NxFPayoutRatioAn: types.ToFloat64Ptr(rs.NxFPayoutRatioAn),

			// 第2四半期業績予想
			FSales2Q: types.ToFloat64Ptr(rs.FSales2Q),
			FOP2Q:    types.ToFloat64Ptr(rs.FOP2Q),
			FOdP2Q:   types.ToFloat64Ptr(rs.FOdP2Q),
			FNP2Q:    types.ToFloat64Ptr(rs.FNP2Q),
			FEPS2Q:   types.ToFloat64Ptr(rs.FEPS2Q),

			// 翌期第2四半期業績予想
			NxFSales2Q: types.ToFloat64Ptr(rs.NxFSales2Q),
			NxFOP2Q:    types.ToFloat64Ptr(rs.NxFOP2Q),
			NxFOdP2Q:   types.ToFloat64Ptr(rs.NxFOdP2Q),
			NxFNp2Q:    types.ToFloat64Ptr(rs.NxFNp2Q),
			NxFEPS2Q:   types.ToFloat64Ptr(rs.NxFEPS2Q),

			// 期末業績予想
			FSales: types.ToFloat64Ptr(rs.FSales),
			FOP:    types.ToFloat64Ptr(rs.FOP),
			FOdP:   types.ToFloat64Ptr(rs.FOdP),
			FNP:    types.ToFloat64Ptr(rs.FNP),
			FEPS:   types.ToFloat64Ptr(rs.FEPS),

			// 翌期末業績予想
			NxFSales: types.ToFloat64Ptr(rs.NxFSales),
			NxFOP:    types.ToFloat64Ptr(rs.NxFOP),
			NxFOdP:   types.ToFloat64Ptr(rs.NxFOdP),
			NxFNp:    types.ToFloat64Ptr(rs.NxFNp),
			NxFEPS:   types.ToFloat64Ptr(rs.NxFEPS),

			// その他
			MatChgSub:  rs.MatChgSub,
			SigChgInC:  rs.SigChgInC,
			ChgByASRev: rs.ChgByASRev,
			ChgNoASRev: rs.ChgNoASRev,
			ChgAcEst:   rs.ChgAcEst,
			RetroRst:   rs.RetroRst,
			ShOutFY:    types.ToInt64Ptr(rs.ShOutFY),
			TrShFY:     types.ToInt64Ptr(rs.TrShFY),
			AvgSh:      types.ToInt64Ptr(rs.AvgSh),

			// 非連結財務数値
			NCSales: types.ToFloat64Ptr(rs.NCSales),
			NCOP:    types.ToFloat64Ptr(rs.NCOP),
			NCOdP:   types.ToFloat64Ptr(rs.NCOdP),
			NCNP:    types.ToFloat64Ptr(rs.NCNP),
			NCEPS:   types.ToFloat64Ptr(rs.NCEPS),
			NCTA:    types.ToFloat64Ptr(rs.NCTA),
			NCEq:    types.ToFloat64Ptr(rs.NCEq),
			NCEqAR:  types.ToFloat64Ptr(rs.NCEqAR),
			NCBPS:   types.ToFloat64Ptr(rs.NCBPS),

			// 非連結第2四半期予想
			FNCSales2Q: types.ToFloat64Ptr(rs.FNCSales2Q),
			FNCOP2Q:    types.ToFloat64Ptr(rs.FNCOP2Q),
			FNCOdP2Q:   types.ToFloat64Ptr(rs.FNCOdP2Q),
			FNCNP2Q:    types.ToFloat64Ptr(rs.FNCNP2Q),
			FNCEPS2Q:   types.ToFloat64Ptr(rs.FNCEPS2Q),

			// 翌期非連結第2四半期予想
			NxFNCSales2Q: types.ToFloat64Ptr(rs.NxFNCSales2Q),
			NxFNCOP2Q:    types.ToFloat64Ptr(rs.NxFNCOP2Q),
			NxFNCOdP2Q:   types.ToFloat64Ptr(rs.NxFNCOdP2Q),
			NxFNCNP2Q:    types.ToFloat64Ptr(rs.NxFNCNP2Q),
			NxFNCEPS2Q:   types.ToFloat64Ptr(rs.NxFNCEPS2Q),

			// 非連結期末予想
			FNCSales: types.ToFloat64Ptr(rs.FNCSales),
			FNCOP:    types.ToFloat64Ptr(rs.FNCOP),
			FNCOdP:   types.ToFloat64Ptr(rs.FNCOdP),
			FNCNP:    types.ToFloat64Ptr(rs.FNCNP),
			FNCEPS:   types.ToFloat64Ptr(rs.FNCEPS),

			// 翌期非連結期末予想
			NxFNCSales: types.ToFloat64Ptr(rs.NxFNCSales),
			NxFNCOP:    types.ToFloat64Ptr(rs.NxFNCOP),
			NxFNCOdP:   types.ToFloat64Ptr(rs.NxFNCOdP),
			NxFNCNP:    types.ToFloat64Ptr(rs.NxFNCNP),
			NxFNCEPS:   types.ToFloat64Ptr(rs.NxFNCEPS),
		}
	}

	return nil
}

type StatementsParams struct {
	Code          string // 銘柄コード（4桁または5桁）
	Date          string // 開示日付（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// GetStatements は財務諸表データを取得します。
// パラメータ:
// - code: 銘柄コード（空の場合は全銘柄）
// - date: 開示日付（空の場合は全期間）
func (s *StatementsService) GetStatements(code string, date string) ([]Statement, error) {
	params := StatementsParams{
		Code: code,
		Date: date,
	}
	resp, err := s.GetStatementsWithParams(params)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetStatementsWithParams は詳細パラメータを指定して財務諸表データを取得します。
func (s *StatementsService) GetStatementsWithParams(params StatementsParams) (*StatementsResponse, error) {
	path := "/fins/summary"

	query := "?"
	if params.Code != "" {
		query += fmt.Sprintf("code=%s&", params.Code)
	}
	if params.Date != "" {
		query += fmt.Sprintf("date=%s&", params.Date)
	}
	if params.PaginationKey != "" {
		query += fmt.Sprintf("pagination_key=%s&", params.PaginationKey)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp StatementsResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get statements: %w", err)
	}

	return &resp, nil
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

	// 最新のものを探す（DiscDateでソート）
	latestStmt := statements[0]
	for _, stmt := range statements {
		if stmt.DiscDate > latestStmt.DiscDate {
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

// GetAllStatementsByCode は指定銘柄の全財務諸表データを取得します（ページネーション対応）。
func (s *StatementsService) GetAllStatementsByCode(code string) ([]Statement, error) {
	var allStatements []Statement
	paginationKey := ""

	for {
		params := StatementsParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetStatementsWithParams(params)
		if err != nil {
			return nil, err
		}

		allStatements = append(allStatements, resp.Data...)

		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allStatements, nil
}
