package jquants

import (
	"context"
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
	Sales  types.NullableFloat64 `json:"Sales"`
	OP     types.NullableFloat64 `json:"OP"`
	OdP    types.NullableFloat64 `json:"OdP"`
	NP     types.NullableFloat64 `json:"NP"`
	TA     types.NullableFloat64 `json:"TA"`
	Eq     types.NullableFloat64 `json:"Eq"`
	CashEq types.NullableFloat64 `json:"CashEq"`
	CFO    types.NullableFloat64 `json:"CFO"`
	CFI    types.NullableFloat64 `json:"CFI"`
	CFF    types.NullableFloat64 `json:"CFF"`

	// 財務指標
	EPS  types.NullableFloat64 `json:"EPS"`
	DEPS types.NullableFloat64 `json:"DEPS"`
	BPS  types.NullableFloat64 `json:"BPS"`
	EqAR types.NullableFloat64 `json:"EqAR"`

	// 配当実績
	Div1Q         types.NullableFloat64 `json:"Div1Q"`
	Div2Q         types.NullableFloat64 `json:"Div2Q"`
	Div3Q         types.NullableFloat64 `json:"Div3Q"`
	DivFY         types.NullableFloat64 `json:"DivFY"`
	DivAnn        types.NullableFloat64 `json:"DivAnn"`
	DivUnit       types.NullableFloat64 `json:"DivUnit"`
	DivTotalAnn   types.NullableFloat64 `json:"DivTotalAnn"`
	PayoutRatioAn types.NullableFloat64 `json:"PayoutRatioAnn"`

	// 配当予想
	FDiv1Q         types.NullableFloat64 `json:"FDiv1Q"`
	FDiv2Q         types.NullableFloat64 `json:"FDiv2Q"`
	FDiv3Q         types.NullableFloat64 `json:"FDiv3Q"`
	FDivFY         types.NullableFloat64 `json:"FDivFY"`
	FDivAnn        types.NullableFloat64 `json:"FDivAnn"`
	FDivUnit       types.NullableFloat64 `json:"FDivUnit"`
	FDivTotalAnn   types.NullableFloat64 `json:"FDivTotalAnn"`
	FPayoutRatioAn types.NullableFloat64 `json:"FPayoutRatioAnn"`

	// 翌期配当予想
	NxFDiv1Q         types.NullableFloat64 `json:"NxFDiv1Q"`
	NxFDiv2Q         types.NullableFloat64 `json:"NxFDiv2Q"`
	NxFDiv3Q         types.NullableFloat64 `json:"NxFDiv3Q"`
	NxFDivFY         types.NullableFloat64 `json:"NxFDivFY"`
	NxFDivAnn        types.NullableFloat64 `json:"NxFDivAnn"`
	NxFDivUnit       types.NullableFloat64 `json:"NxFDivUnit"`
	NxFPayoutRatioAn types.NullableFloat64 `json:"NxFPayoutRatioAnn"`

	// 第2四半期業績予想
	FSales2Q types.NullableFloat64 `json:"FSales2Q"`
	FOP2Q    types.NullableFloat64 `json:"FOP2Q"`
	FOdP2Q   types.NullableFloat64 `json:"FOdP2Q"`
	FNP2Q    types.NullableFloat64 `json:"FNP2Q"`
	FEPS2Q   types.NullableFloat64 `json:"FEPS2Q"`

	// 翌期第2四半期業績予想
	NxFSales2Q types.NullableFloat64 `json:"NxFSales2Q"`
	NxFOP2Q    types.NullableFloat64 `json:"NxFOP2Q"`
	NxFOdP2Q   types.NullableFloat64 `json:"NxFOdP2Q"`
	NxFNp2Q    types.NullableFloat64 `json:"NxFNp2Q"`
	NxFEPS2Q   types.NullableFloat64 `json:"NxFEPS2Q"`

	// 期末業績予想
	FSales types.NullableFloat64 `json:"FSales"`
	FOP    types.NullableFloat64 `json:"FOP"`
	FOdP   types.NullableFloat64 `json:"FOdP"`
	FNP    types.NullableFloat64 `json:"FNP"`
	FEPS   types.NullableFloat64 `json:"FEPS"`

	// 翌期末業績予想
	NxFSales types.NullableFloat64 `json:"NxFSales"`
	NxFOP    types.NullableFloat64 `json:"NxFOP"`
	NxFOdP   types.NullableFloat64 `json:"NxFOdP"`
	NxFNp    types.NullableFloat64 `json:"NxFNp"`
	NxFEPS   types.NullableFloat64 `json:"NxFEPS"`

	// その他
	MatChgSub  string              `json:"MatChgSub"`
	SigChgInC  string              `json:"SigChgInC"`
	ChgByASRev string              `json:"ChgByASRev"`
	ChgNoASRev string              `json:"ChgNoASRev"`
	ChgAcEst   string              `json:"ChgAcEst"`
	RetroRst   string              `json:"RetroRst"`
	ShOutFY    types.NullableInt64 `json:"ShOutFY"`
	TrShFY     types.NullableInt64 `json:"TrShFY"`
	AvgSh      types.NullableInt64 `json:"AvgSh"`

	// 非連結財務数値
	NCSales types.NullableFloat64 `json:"NCSales"`
	NCOP    types.NullableFloat64 `json:"NCOP"`
	NCOdP   types.NullableFloat64 `json:"NCOdP"`
	NCNP    types.NullableFloat64 `json:"NCNP"`
	NCEPS   types.NullableFloat64 `json:"NCEPS"`
	NCTA    types.NullableFloat64 `json:"NCTA"`
	NCEq    types.NullableFloat64 `json:"NCEq"`
	NCEqAR  types.NullableFloat64 `json:"NCEqAR"`
	NCBPS   types.NullableFloat64 `json:"NCBPS"`

	// 非連結第2四半期予想
	FNCSales2Q types.NullableFloat64 `json:"FNCSales2Q"`
	FNCOP2Q    types.NullableFloat64 `json:"FNCOP2Q"`
	FNCOdP2Q   types.NullableFloat64 `json:"FNCOdP2Q"`
	FNCNP2Q    types.NullableFloat64 `json:"FNCNP2Q"`
	FNCEPS2Q   types.NullableFloat64 `json:"FNCEPS2Q"`

	// 翌期非連結第2四半期予想
	NxFNCSales2Q types.NullableFloat64 `json:"NxFNCSales2Q"`
	NxFNCOP2Q    types.NullableFloat64 `json:"NxFNCOP2Q"`
	NxFNCOdP2Q   types.NullableFloat64 `json:"NxFNCOdP2Q"`
	NxFNCNP2Q    types.NullableFloat64 `json:"NxFNCNP2Q"`
	NxFNCEPS2Q   types.NullableFloat64 `json:"NxFNCEPS2Q"`

	// 非連結期末予想
	FNCSales types.NullableFloat64 `json:"FNCSales"`
	FNCOP    types.NullableFloat64 `json:"FNCOP"`
	FNCOdP   types.NullableFloat64 `json:"FNCOdP"`
	FNCNP    types.NullableFloat64 `json:"FNCNP"`
	FNCEPS   types.NullableFloat64 `json:"FNCEPS"`

	// 翌期非連結期末予想
	NxFNCSales types.NullableFloat64 `json:"NxFNCSales"`
	NxFNCOP    types.NullableFloat64 `json:"NxFNCOP"`
	NxFNCOdP   types.NullableFloat64 `json:"NxFNCOdP"`
	NxFNCNP    types.NullableFloat64 `json:"NxFNCNP"`
	NxFNCEPS   types.NullableFloat64 `json:"NxFNCEPS"`
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
			Sales:  rs.Sales.Ptr(),
			OP:     rs.OP.Ptr(),
			OdP:    rs.OdP.Ptr(),
			NP:     rs.NP.Ptr(),
			TA:     rs.TA.Ptr(),
			Eq:     rs.Eq.Ptr(),
			CashEq: rs.CashEq.Ptr(),
			CFO:    rs.CFO.Ptr(),
			CFI:    rs.CFI.Ptr(),
			CFF:    rs.CFF.Ptr(),

			// 財務指標
			EPS:  rs.EPS.Ptr(),
			DEPS: rs.DEPS.Ptr(),
			BPS:  rs.BPS.Ptr(),
			EqAR: rs.EqAR.Ptr(),

			// 配当実績
			Div1Q:         rs.Div1Q.Ptr(),
			Div2Q:         rs.Div2Q.Ptr(),
			Div3Q:         rs.Div3Q.Ptr(),
			DivFY:         rs.DivFY.Ptr(),
			DivAnn:        rs.DivAnn.Ptr(),
			DivUnit:       rs.DivUnit.Ptr(),
			DivTotalAnn:   rs.DivTotalAnn.Ptr(),
			PayoutRatioAn: rs.PayoutRatioAn.Ptr(),

			// 配当予想
			FDiv1Q:         rs.FDiv1Q.Ptr(),
			FDiv2Q:         rs.FDiv2Q.Ptr(),
			FDiv3Q:         rs.FDiv3Q.Ptr(),
			FDivFY:         rs.FDivFY.Ptr(),
			FDivAnn:        rs.FDivAnn.Ptr(),
			FDivUnit:       rs.FDivUnit.Ptr(),
			FDivTotalAnn:   rs.FDivTotalAnn.Ptr(),
			FPayoutRatioAn: rs.FPayoutRatioAn.Ptr(),

			// 翌期配当予想
			NxFDiv1Q:         rs.NxFDiv1Q.Ptr(),
			NxFDiv2Q:         rs.NxFDiv2Q.Ptr(),
			NxFDiv3Q:         rs.NxFDiv3Q.Ptr(),
			NxFDivFY:         rs.NxFDivFY.Ptr(),
			NxFDivAnn:        rs.NxFDivAnn.Ptr(),
			NxFDivUnit:       rs.NxFDivUnit.Ptr(),
			NxFPayoutRatioAn: rs.NxFPayoutRatioAn.Ptr(),

			// 第2四半期業績予想
			FSales2Q: rs.FSales2Q.Ptr(),
			FOP2Q:    rs.FOP2Q.Ptr(),
			FOdP2Q:   rs.FOdP2Q.Ptr(),
			FNP2Q:    rs.FNP2Q.Ptr(),
			FEPS2Q:   rs.FEPS2Q.Ptr(),

			// 翌期第2四半期業績予想
			NxFSales2Q: rs.NxFSales2Q.Ptr(),
			NxFOP2Q:    rs.NxFOP2Q.Ptr(),
			NxFOdP2Q:   rs.NxFOdP2Q.Ptr(),
			NxFNp2Q:    rs.NxFNp2Q.Ptr(),
			NxFEPS2Q:   rs.NxFEPS2Q.Ptr(),

			// 期末業績予想
			FSales: rs.FSales.Ptr(),
			FOP:    rs.FOP.Ptr(),
			FOdP:   rs.FOdP.Ptr(),
			FNP:    rs.FNP.Ptr(),
			FEPS:   rs.FEPS.Ptr(),

			// 翌期末業績予想
			NxFSales: rs.NxFSales.Ptr(),
			NxFOP:    rs.NxFOP.Ptr(),
			NxFOdP:   rs.NxFOdP.Ptr(),
			NxFNp:    rs.NxFNp.Ptr(),
			NxFEPS:   rs.NxFEPS.Ptr(),

			// その他
			MatChgSub:  rs.MatChgSub,
			SigChgInC:  rs.SigChgInC,
			ChgByASRev: rs.ChgByASRev,
			ChgNoASRev: rs.ChgNoASRev,
			ChgAcEst:   rs.ChgAcEst,
			RetroRst:   rs.RetroRst,
			ShOutFY:    rs.ShOutFY.Ptr(),
			TrShFY:     rs.TrShFY.Ptr(),
			AvgSh:      rs.AvgSh.Ptr(),

			// 非連結財務数値
			NCSales: rs.NCSales.Ptr(),
			NCOP:    rs.NCOP.Ptr(),
			NCOdP:   rs.NCOdP.Ptr(),
			NCNP:    rs.NCNP.Ptr(),
			NCEPS:   rs.NCEPS.Ptr(),
			NCTA:    rs.NCTA.Ptr(),
			NCEq:    rs.NCEq.Ptr(),
			NCEqAR:  rs.NCEqAR.Ptr(),
			NCBPS:   rs.NCBPS.Ptr(),

			// 非連結第2四半期予想
			FNCSales2Q: rs.FNCSales2Q.Ptr(),
			FNCOP2Q:    rs.FNCOP2Q.Ptr(),
			FNCOdP2Q:   rs.FNCOdP2Q.Ptr(),
			FNCNP2Q:    rs.FNCNP2Q.Ptr(),
			FNCEPS2Q:   rs.FNCEPS2Q.Ptr(),

			// 翌期非連結第2四半期予想
			NxFNCSales2Q: rs.NxFNCSales2Q.Ptr(),
			NxFNCOP2Q:    rs.NxFNCOP2Q.Ptr(),
			NxFNCOdP2Q:   rs.NxFNCOdP2Q.Ptr(),
			NxFNCNP2Q:    rs.NxFNCNP2Q.Ptr(),
			NxFNCEPS2Q:   rs.NxFNCEPS2Q.Ptr(),

			// 非連結期末予想
			FNCSales: rs.FNCSales.Ptr(),
			FNCOP:    rs.FNCOP.Ptr(),
			FNCOdP:   rs.FNCOdP.Ptr(),
			FNCNP:    rs.FNCNP.Ptr(),
			FNCEPS:   rs.FNCEPS.Ptr(),

			// 翌期非連結期末予想
			NxFNCSales: rs.NxFNCSales.Ptr(),
			NxFNCOP:    rs.NxFNCOP.Ptr(),
			NxFNCOdP:   rs.NxFNCOdP.Ptr(),
			NxFNCNP:    rs.NxFNCNP.Ptr(),
			NxFNCEPS:   rs.NxFNCEPS.Ptr(),
		}
	}

	return nil
}

type StatementsParams struct {
	Code          string // 銘柄コード（4桁または5桁）
	Date          string // 開示日付（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// GetStatements は指定された条件で財務諸表データを取得します。
// codeまたはdateのいずれかが必須です。
// パラメータ:
// - Code: 銘柄コード（例: "7203" または "72030"）
// - Date: 開示日付（例: "20240101" または "2024-01-01"）
// - PaginationKey: ページネーション用キー
func (s *StatementsService) GetStatements(ctx context.Context, params StatementsParams) (*StatementsResponse, error) {
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
	if err := s.client.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get statements: %w", err)
	}

	return &resp, nil
}

// GetAllStatementsByCode は指定銘柄の全期間の財務諸表データを取得します。
// ページネーションを使用して全データを取得します。
func (s *StatementsService) GetAllStatementsByCode(ctx context.Context, code string) ([]Statement, error) {
	var allStatements []Statement
	paginationKey := ""

	for {
		resp, err := s.GetStatements(ctx, StatementsParams{
			Code:          code,
			PaginationKey: paginationKey,
		})
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

// GetStatementsByCodeAndDate は指定銘柄の指定日の財務諸表データを取得します。
func (s *StatementsService) GetStatementsByCodeAndDate(ctx context.Context, code, date string) ([]Statement, error) {
	resp, err := s.GetStatements(ctx, StatementsParams{
		Code: code,
		Date: date,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetStatementsByDate は指定日の全銘柄の財務諸表データを取得します。
// ページネーションを使用して全データを取得します。
func (s *StatementsService) GetStatementsByDate(ctx context.Context, date string) ([]Statement, error) {
	var allStatements []Statement
	paginationKey := ""

	for {
		resp, err := s.GetStatements(ctx, StatementsParams{
			Date:          date,
			PaginationKey: paginationKey,
		})
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

// GetLatestStatements は指定銘柄の最新財務諸表を取得します。
// 例: GetLatestStatements("7203") でトヨタ自動車の最新決算データを取得
func (s *StatementsService) GetLatestStatements(ctx context.Context, code string) (*Statement, error) {
	statements, err := s.GetAllStatementsByCode(ctx, code)
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
