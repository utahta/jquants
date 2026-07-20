package jquants

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// EdinetCrossShareholdingsService は政策保有株式（EDINET）を取得するサービスです。
// 有価証券報告書「株式の保有状況」に基づく政策保有株式データを取得します。
type EdinetCrossShareholdingsService struct {
	client client.HTTPClient
}

// NewEdinetCrossShareholdingsService は新しいEdinetCrossShareholdingsServiceを作成します。
func NewEdinetCrossShareholdingsService(c client.HTTPClient) *EdinetCrossShareholdingsService {
	return &EdinetCrossShareholdingsService{client: c}
}

// EdinetCrossShareholdingsParams は政策保有株式のリクエストパラメータです。
// すべて任意ですが、edinet_codeとcodeの同時指定はできません。
// すべて省略した場合はAPI実行日に提出された全有報のデータ一覧が返ります。
type EdinetCrossShareholdingsParams struct {
	EdinetCode    string // EDINETコード（例: E03814）（codeとの同時指定は不可）
	Code          string // 4桁もしくは5桁の銘柄コード（edinet_codeとの同時指定は不可）
	Date          string // 提出日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// EdinetCrossShareholdingsResponse は政策保有株式のレスポンスです。
type EdinetCrossShareholdingsResponse struct {
	Data          []EdinetCrossShareholdingDoc `json:"data"`
	PaginationKey string                       `json:"pagination_key"` // ページネーションキー
}

// EdinetCrossShareholdingDoc は有価証券報告書ごとの政策保有株式データを表します。
// J-Quants API /edinet/cross-shareholdings エンドポイントのレスポンスデータ。
// データは2020年3月31日以降、対象書類は有価証券報告書。Standardプラン以上で利用可能。
type EdinetCrossShareholdingDoc struct {
	// 書類メタ情報
	DocId       string `json:"DocId"`       // EDINET書類管理番号（S + 7桁英数字）
	Code        string `json:"Code"`        // 提出会社の銘柄コード（5桁）
	EdinetCode  string `json:"EdinetCode"`  // 提出会社のEDINETコード
	FilerName   string `json:"FilerName"`   // 提出者名（会社名）
	FilerNameEn string `json:"FilerNameEn"` // 提出者名（英語）
	DocTypeCode string `json:"DocTypeCode"` // 書類種別コード（120=有価証券報告書）
	SubDate     string `json:"SubDate"`     // 提出日（YYYY-MM-DD形式）
	SubTime     string `json:"SubTime"`     // 提出時刻（HH:MM:SS形式）
	PerSt       string `json:"PerSt"`       // 対象事業年度の開始日（YYYY-MM-DD形式）
	PerEn       string `json:"PerEn"`       // 対象事業年度の終了日（YYYY-MM-DD形式）

	// 保有主体ブロック（記載がない場合はnil）
	Report        *EdinetCrossShareholdingBlock `json:"Report"`        // 提出会社
	Largest       *EdinetCrossShareholdingBlock `json:"Largest"`       // 連結最大保有会社
	SecondLargest *EdinetCrossShareholdingBlock `json:"SecondLargest"` // 連結第二最大保有会社
}

// EdinetCrossShareholdingBlock は保有主体（提出会社/連結最大保有会社/連結第二最大保有会社）ごとの保有状況を表します。
type EdinetCrossShareholdingBlock struct {
	// 保有主体情報
	HldrName       string `json:"HldrName"`       // 当該保有主体の会社名
	HldrCode       string `json:"HldrCode"`       // 当該保有主体の証券コード（5桁）
	HldrEdinetCode string `json:"HldrEdinetCode"` // 当該保有主体のEDINETコード

	// 上場株式
	ListedIss        float64 `json:"ListedIss"`        // 銘柄数
	ListedBookVal    float64 `json:"ListedBookVal"`    // 貸借対照表計上額合計（円）
	ListedIncIss     float64 `json:"ListedIncIss"`     // 株式数が増加した銘柄数
	ListedIncAcqCost float64 `json:"ListedIncAcqCost"` // 増加に係る取得価額合計（円）
	ListedDecIss     float64 `json:"ListedDecIss"`     // 株式数が減少した銘柄数
	ListedDecSaleAmt float64 `json:"ListedDecSaleAmt"` // 減少に係る売却価額合計（円）
	ListedIncRsn     string  `json:"ListedIncRsn"`     // 株式数が増加した理由

	// 非上場株式
	NonListedIss        float64 `json:"NonListedIss"`        // 銘柄数
	NonListedBookVal    float64 `json:"NonListedBookVal"`    // 貸借対照表計上額合計（円）
	NonListedIncIss     float64 `json:"NonListedIncIss"`     // 株式数が増加した銘柄数
	NonListedIncAcqCost float64 `json:"NonListedIncAcqCost"` // 増加に係る取得価額合計（円）
	NonListedDecIss     float64 `json:"NonListedDecIss"`     // 株式数が減少した銘柄数
	NonListedDecSaleAmt float64 `json:"NonListedDecSaleAmt"` // 減少に係る売却価額合計（円）
	NonListedIncRsn     string  `json:"NonListedIncRsn"`     // 株式数が増加した理由

	// 銘柄レコード
	Spec []EdinetCrossShareholdingIssue `json:"Spec"` // 特定投資株式
	Deem []EdinetCrossShareholdingIssue `json:"Deem"` // みなし保有株式

	// 注釈（HTML断片）
	SpecFn string `json:"SpecFn"` // 特定投資株式の注釈
	DeemFn string `json:"DeemFn"` // みなし保有株式の注釈
}

// EdinetCrossShareholdingIssue は特定投資株式・みなし保有株式の1銘柄分のデータを表します。
// 株式数・貸借対照表計上額が非開示の場合は数値がnilとなり、対応する*NotDiscに
// 非開示マーカー生値（* / ＊ / ※ / （注N））が入ります。
type EdinetCrossShareholdingIssue struct {
	// 保有先銘柄情報
	IsrName       string `json:"IsrName"`       // 保有先銘柄名
	IsrCode       string `json:"IsrCode"`       // 保有先の銘柄コード（5桁、IsrNameから名寄せ）
	IsrEdinetCode string `json:"IsrEdinetCode"` // 保有先のEDINETコード（IsrNameから名寄せ）

	// 株式数・貸借対照表計上額（非開示の場合はnil）
	CurShs     *float64 `json:"CurShs"`     // 当事業年度の株式数（株）
	PriShs     *float64 `json:"PriShs"`     // 前事業年度の株式数（株）
	CurBookVal *float64 `json:"CurBookVal"` // 当事業年度の貸借対照表計上額（円）
	PriBookVal *float64 `json:"PriBookVal"` // 前事業年度の貸借対照表計上額（円）

	// 非開示マーカー生値（* / ＊ / ※ / （注N）。開示されている場合は空文字）
	CurShsNotDisc     string `json:"CurShsNotDisc"`     // 当期株式数の非開示マーカー
	PriShsNotDisc     string `json:"PriShsNotDisc"`     // 前期株式数の非開示マーカー
	CurBookValNotDisc string `json:"CurBookValNotDisc"` // 当期貸借対照表計上額の非開示マーカー
	PriBookValNotDisc string `json:"PriBookValNotDisc"` // 前期貸借対照表計上額の非開示マーカー

	// その他
	HoldRat      string `json:"HoldRat"`      // 保有目的・業務提携の概要・定量効果・増加理由（複合テキスト）
	IsrHolds     string `json:"IsrHolds"`     // 当社の株式の保有の有無（生データ。例: 有 / 無 / 無(注)３）
	IsrHoldsCode string `json:"IsrHoldsCode"` // 当社の株式の保有の有無の正規化値（1=有、0=無、2=判定不能）
}

// RawEdinetCrossShareholdingDoc is used for unmarshaling JSON response with mixed types
type RawEdinetCrossShareholdingDoc struct {
	// 書類メタ情報
	DocId       string `json:"DocId"`
	Code        string `json:"Code"`
	EdinetCode  string `json:"EdinetCode"`
	FilerName   string `json:"FilerName"`
	FilerNameEn string `json:"FilerNameEn"`
	DocTypeCode string `json:"DocTypeCode"`
	SubDate     string `json:"SubDate"`
	SubTime     string `json:"SubTime"`
	PerSt       string `json:"PerSt"`
	PerEn       string `json:"PerEn"`

	// 保有主体ブロック
	Report        *RawEdinetCrossShareholdingBlock `json:"Report"`
	Largest       *RawEdinetCrossShareholdingBlock `json:"Largest"`
	SecondLargest *RawEdinetCrossShareholdingBlock `json:"SecondLargest"`
}

// RawEdinetCrossShareholdingBlock is used for unmarshaling JSON response with mixed types
type RawEdinetCrossShareholdingBlock struct {
	HldrName       string `json:"HldrName"`
	HldrCode       string `json:"HldrCode"`
	HldrEdinetCode string `json:"HldrEdinetCode"`

	ListedIss        types.NullableFloat64 `json:"ListedIss"`
	ListedBookVal    types.NullableFloat64 `json:"ListedBookVal"`
	ListedIncIss     types.NullableFloat64 `json:"ListedIncIss"`
	ListedIncAcqCost types.NullableFloat64 `json:"ListedIncAcqCost"`
	ListedDecIss     types.NullableFloat64 `json:"ListedDecIss"`
	ListedDecSaleAmt types.NullableFloat64 `json:"ListedDecSaleAmt"`
	ListedIncRsn     string                `json:"ListedIncRsn"`

	NonListedIss        types.NullableFloat64 `json:"NonListedIss"`
	NonListedBookVal    types.NullableFloat64 `json:"NonListedBookVal"`
	NonListedIncIss     types.NullableFloat64 `json:"NonListedIncIss"`
	NonListedIncAcqCost types.NullableFloat64 `json:"NonListedIncAcqCost"`
	NonListedDecIss     types.NullableFloat64 `json:"NonListedDecIss"`
	NonListedDecSaleAmt types.NullableFloat64 `json:"NonListedDecSaleAmt"`
	NonListedIncRsn     string                `json:"NonListedIncRsn"`

	Spec []RawEdinetCrossShareholdingIssue `json:"Spec"`
	Deem []RawEdinetCrossShareholdingIssue `json:"Deem"`

	SpecFn string `json:"SpecFn"`
	DeemFn string `json:"DeemFn"`
}

// RawEdinetCrossShareholdingIssue is used for unmarshaling JSON response with mixed types
type RawEdinetCrossShareholdingIssue struct {
	IsrName       string `json:"IsrName"`
	IsrCode       string `json:"IsrCode"`
	IsrEdinetCode string `json:"IsrEdinetCode"`

	CurShs     types.NullableFloat64 `json:"CurShs"`
	PriShs     types.NullableFloat64 `json:"PriShs"`
	CurBookVal types.NullableFloat64 `json:"CurBookVal"`
	PriBookVal types.NullableFloat64 `json:"PriBookVal"`

	CurShsNotDisc     string `json:"CurShsNotDisc"`
	PriShsNotDisc     string `json:"PriShsNotDisc"`
	CurBookValNotDisc string `json:"CurBookValNotDisc"`
	PriBookValNotDisc string `json:"PriBookValNotDisc"`

	HoldRat      string `json:"HoldRat"`
	IsrHolds     string `json:"IsrHolds"`
	IsrHoldsCode string `json:"IsrHoldsCode"`
}

// UnmarshalJSON implements custom JSON unmarshaling for EdinetCrossShareholdingsResponse
func (r *EdinetCrossShareholdingsResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawEdinetCrossShareholdingDoc
	type rawResponse struct {
		Data          []RawEdinetCrossShareholdingDoc `json:"data"`
		PaginationKey string                          `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawEdinetCrossShareholdingDoc to EdinetCrossShareholdingDoc
	r.Data = make([]EdinetCrossShareholdingDoc, len(raw.Data))
	for idx, rd := range raw.Data {
		r.Data[idx] = EdinetCrossShareholdingDoc{
			// 書類メタ情報
			DocId:       rd.DocId,
			Code:        rd.Code,
			EdinetCode:  rd.EdinetCode,
			FilerName:   rd.FilerName,
			FilerNameEn: rd.FilerNameEn,
			DocTypeCode: rd.DocTypeCode,
			SubDate:     rd.SubDate,
			SubTime:     rd.SubTime,
			PerSt:       rd.PerSt,
			PerEn:       rd.PerEn,

			// 保有主体ブロック
			Report:        convertEdinetCrossShareholdingBlock(rd.Report),
			Largest:       convertEdinetCrossShareholdingBlock(rd.Largest),
			SecondLargest: convertEdinetCrossShareholdingBlock(rd.SecondLargest),
		}
	}

	return nil
}

// convertEdinetCrossShareholdingBlock converts RawEdinetCrossShareholdingBlock to EdinetCrossShareholdingBlock
func convertEdinetCrossShareholdingBlock(raw *RawEdinetCrossShareholdingBlock) *EdinetCrossShareholdingBlock {
	if raw == nil {
		return nil
	}

	return &EdinetCrossShareholdingBlock{
		HldrName:       raw.HldrName,
		HldrCode:       raw.HldrCode,
		HldrEdinetCode: raw.HldrEdinetCode,

		ListedIss:        raw.ListedIss.Or(0),
		ListedBookVal:    raw.ListedBookVal.Or(0),
		ListedIncIss:     raw.ListedIncIss.Or(0),
		ListedIncAcqCost: raw.ListedIncAcqCost.Or(0),
		ListedDecIss:     raw.ListedDecIss.Or(0),
		ListedDecSaleAmt: raw.ListedDecSaleAmt.Or(0),
		ListedIncRsn:     raw.ListedIncRsn,

		NonListedIss:        raw.NonListedIss.Or(0),
		NonListedBookVal:    raw.NonListedBookVal.Or(0),
		NonListedIncIss:     raw.NonListedIncIss.Or(0),
		NonListedIncAcqCost: raw.NonListedIncAcqCost.Or(0),
		NonListedDecIss:     raw.NonListedDecIss.Or(0),
		NonListedDecSaleAmt: raw.NonListedDecSaleAmt.Or(0),
		NonListedIncRsn:     raw.NonListedIncRsn,

		Spec: convertEdinetCrossShareholdingIssues(raw.Spec),
		Deem: convertEdinetCrossShareholdingIssues(raw.Deem),

		SpecFn: raw.SpecFn,
		DeemFn: raw.DeemFn,
	}
}

// convertEdinetCrossShareholdingIssues converts RawEdinetCrossShareholdingIssue slice to EdinetCrossShareholdingIssue slice
func convertEdinetCrossShareholdingIssues(raws []RawEdinetCrossShareholdingIssue) []EdinetCrossShareholdingIssue {
	if raws == nil {
		return nil
	}

	issues := make([]EdinetCrossShareholdingIssue, len(raws))
	for idx, ri := range raws {
		issues[idx] = EdinetCrossShareholdingIssue{
			IsrName:       ri.IsrName,
			IsrCode:       ri.IsrCode,
			IsrEdinetCode: ri.IsrEdinetCode,

			CurShs:     ri.CurShs.Ptr(),
			PriShs:     ri.PriShs.Ptr(),
			CurBookVal: ri.CurBookVal.Ptr(),
			PriBookVal: ri.PriBookVal.Ptr(),

			CurShsNotDisc:     ri.CurShsNotDisc,
			PriShsNotDisc:     ri.PriShsNotDisc,
			CurBookValNotDisc: ri.CurBookValNotDisc,
			PriBookValNotDisc: ri.PriBookValNotDisc,

			HoldRat:      ri.HoldRat,
			IsrHolds:     ri.IsrHolds,
			IsrHoldsCode: ri.IsrHoldsCode,
		}
	}

	return issues
}

// GetCrossShareholdings は政策保有株式を取得します。
// パラメータをすべて省略した場合はAPI実行日に提出された全有報のデータ一覧を取得します。
func (s *EdinetCrossShareholdingsService) GetCrossShareholdings(ctx context.Context, params EdinetCrossShareholdingsParams) (*EdinetCrossShareholdingsResponse, error) {
	// edinet_codeとcodeの同時指定は不可
	if params.EdinetCode != "" && params.Code != "" {
		return nil, fmt.Errorf("edinet_code and code cannot be specified together")
	}

	path := "/edinet/cross-shareholdings"

	query := "?"
	if params.EdinetCode != "" {
		query += fmt.Sprintf("edinet_code=%s&", params.EdinetCode)
	}
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

	var resp EdinetCrossShareholdingsResponse
	if err := s.client.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get cross shareholdings: %w", err)
	}

	return &resp, nil
}

// GetCrossShareholdingsByCode は指定銘柄の政策保有株式を取得します。
// ページネーションを使用して全データを取得します。
func (s *EdinetCrossShareholdingsService) GetCrossShareholdingsByCode(ctx context.Context, code string) ([]EdinetCrossShareholdingDoc, error) {
	var allData []EdinetCrossShareholdingDoc
	paginationKey := ""

	for {
		params := EdinetCrossShareholdingsParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetCrossShareholdings(ctx, params)
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

// GetCrossShareholdingsByEdinetCode は指定EDINETコードの政策保有株式を取得します。
// ページネーションを使用して全データを取得します。
func (s *EdinetCrossShareholdingsService) GetCrossShareholdingsByEdinetCode(ctx context.Context, edinetCode string) ([]EdinetCrossShareholdingDoc, error) {
	var allData []EdinetCrossShareholdingDoc
	paginationKey := ""

	for {
		params := EdinetCrossShareholdingsParams{
			EdinetCode:    edinetCode,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetCrossShareholdings(ctx, params)
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

// GetCrossShareholdingsByDate は指定提出日の全有報の政策保有株式を取得します。
// ページネーションを使用して全データを取得します。
func (s *EdinetCrossShareholdingsService) GetCrossShareholdingsByDate(ctx context.Context, date string) ([]EdinetCrossShareholdingDoc, error) {
	var allData []EdinetCrossShareholdingDoc
	paginationKey := ""

	for {
		params := EdinetCrossShareholdingsParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetCrossShareholdings(ctx, params)
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

// HasReport は提出会社の保有ブロックが存在するかを判定します。
func (d *EdinetCrossShareholdingDoc) HasReport() bool {
	return d.Report != nil
}

// HasLargest は連結最大保有会社の保有ブロックが存在するかを判定します。
func (d *EdinetCrossShareholdingDoc) HasLargest() bool {
	return d.Largest != nil
}

// HasSecondLargest は連結第二最大保有会社の保有ブロックが存在するかを判定します。
func (d *EdinetCrossShareholdingDoc) HasSecondLargest() bool {
	return d.SecondLargest != nil
}

// IsDisclosed は当事業年度の株式数が開示されているかを判定します。
func (i *EdinetCrossShareholdingIssue) IsDisclosed() bool {
	return i.CurShs != nil
}
