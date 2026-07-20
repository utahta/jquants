package jquants

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// EdinetLargeVolumeShareholdersService は大量保有報告書（EDINET）を取得するサービスです。
// 大量保有報告書・変更報告書に記載されている発行者、提出者情報を取得します。
type EdinetLargeVolumeShareholdersService struct {
	client client.HTTPClient
}

// NewEdinetLargeVolumeShareholdersService は新しいEdinetLargeVolumeShareholdersServiceを作成します。
func NewEdinetLargeVolumeShareholdersService(c client.HTTPClient) *EdinetLargeVolumeShareholdersService {
	return &EdinetLargeVolumeShareholdersService{client: c}
}

// 大量保有書類種別コード（LargeHldgTypeCode）の定義
const (
	LargeHldgTypeCodeUnknown                       = "0" // 不明
	LargeHldgTypeCodeReport                        = "1" // 大量保有報告書
	LargeHldgTypeCodeChangeReport                  = "2" // 変更報告書
	LargeHldgTypeCodeChangeReportShortTermTransfer = "3" // 変更報告書（短期大量譲渡）
	LargeHldgTypeCodeSpecialReport                 = "4" // 大量保有報告書（特例対象株券等）
	LargeHldgTypeCodeSpecialChangeReport           = "5" // 変更報告書（特例対象株券等）
)

// EdinetLargeVolumeShareholdersParams は大量保有報告書のリクエストパラメータです。
// すべて任意ですが、edinet_codeとcodeの同時指定はできません。
// すべて省略した場合はAPI実行日に提出された全書類のデータ一覧が返ります。
type EdinetLargeVolumeShareholdersParams struct {
	EdinetCode    string // 発行者のEDINETコード（例: E03814）（codeとの同時指定は不可）
	Code          string // 発行者の4桁もしくは5桁の銘柄コード（edinet_codeとの同時指定は不可）
	Date          string // 提出日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// EdinetLargeVolumeShareholdersResponse は大量保有報告書のレスポンスです。
type EdinetLargeVolumeShareholdersResponse struct {
	Data          []EdinetLargeVolumeShareholderDoc `json:"data"`
	PaginationKey string                            `json:"pagination_key"` // ページネーションキー
}

// EdinetLargeVolumeShareholderDoc は大量保有報告書・変更報告書1書類分のデータを表します。
// J-Quants API /edinet/large-volume-shareholders エンドポイントのレスポンスデータ。
// データは提出日2021年7月1日以降、対象書類は大量保有報告書・変更報告書等（書類種別コード350）。
// Standardプラン以上で利用可能。
type EdinetLargeVolumeShareholderDoc struct {
	// 書類メタ情報
	DocId             string `json:"DocId"`             // EDINET書類管理番号（S + 7桁英数字）
	Code              string `json:"Code"`              // 発行者（保有対象銘柄）の銘柄コード（5桁）
	EdinetCode        string `json:"EdinetCode"`        // 発行者のEDINETコード
	IsrName           string `json:"IsrName"`           // 発行者名
	DocTypeCode       string `json:"DocTypeCode"`       // 書類種別コード（350=大量保有報告書関連）
	SubDate           string `json:"SubDate"`           // 提出日（YYYY-MM-DD形式）
	SubTime           string `json:"SubTime"`           // 提出時刻（HH:MM:SS形式）
	LargeHldgTypeCode string `json:"LargeHldgTypeCode"` // 大量保有書類種別コード（LargeHldgTypeCode定数を参照）
	DocTitle          string `json:"DocTitle"`          // 書類表題（例: 大量保有報告書）
	ChgRsn            string `json:"ChgRsn"`            // 報告義務発生日における変更事由（変更報告書のみ）

	// 保有状況の合計
	TotalShsHeld      *float64 `json:"TotalShsHeld"`      // 保有株券等の数の合計（株）（合計欄の記載がない書類ではnull）
	TotalShsRatio     *float64 `json:"TotalShsRatio"`     // 株券等保有割合の合計（小数表現。0.1343=13.43%）（合計欄の記載がない書類ではnull）
	TotalShsRatioLast *float64 `json:"TotalShsRatioLast"` // 直前の報告書に係る株券等保有割合の合計（変更報告書のみ）
	TotalOutStks      float64  `json:"TotalOutStks"`      // 発行済株式等総数（株）

	// 提出者及び共同保有者のレコード
	Hldrs []EdinetLargeVolumeShareholderHolder `json:"Hldrs"`
}

// EdinetLargeVolumeShareholderHolder は提出者・共同保有者1名分のデータを表します。
type EdinetLargeVolumeShareholderHolder struct {
	// 保有者情報
	HldrName          string `json:"HldrName"`          // 保有者の氏名又は名称
	HldrNameEn        string `json:"HldrNameEn"`        // 保有者の名称（英語）
	HldrEdinetCode    string `json:"HldrEdinetCode"`    // 保有者のEDINETコード
	HldrCode          string `json:"HldrCode"`          // 保有者の銘柄コード（保有者が上場会社の場合等）
	LargeHldrTypeCode string `json:"LargeHldrTypeCode"` // 保有者区分コード（1=個人、2=法人、0=不明）
	LargeHldrTypeRaw  string `json:"LargeHldrTypeRaw"`  // 保有者区分（個人法人の別）の書類記載生値

	// 保有状況
	HldgPurp     string   `json:"HldgPurp"`     // 保有目的
	ImpProp      string   `json:"ImpProp"`      // 重要提案行為等
	ColAgr       string   `json:"ColAgr"`       // 担保契約等重要な契約
	ShsHeld      float64  `json:"ShsHeld"`      // 保有株券等の数（株）
	ShsRatio     float64  `json:"ShsRatio"`     // 株券等保有割合（小数表現。0.0572=5.72%）
	ShsRatioLast *float64 `json:"ShsRatioLast"` // 直前の報告書に係る株券等保有割合（変更報告書のみ）

	// 取得資金
	OwnFund    *float64 `json:"OwnFund"`    // 取得資金のうち自己資金額（円）
	TotalBrw   *float64 `json:"TotalBrw"`   // 取得資金のうち借入金額計（円）
	TotalOther *float64 `json:"TotalOther"` // 取得資金のうちその他金額計（円）
	OtherBrk   string   `json:"OtherBrk"`   // その他金額計の内訳（株式分割による取得等、記載がある場合）
	TotalFund  *float64 `json:"TotalFund"`  // 取得資金合計（円）

	// 明細
	AcqDisp  []EdinetLargeVolumeShareholderAcqDisp   `json:"AcqDisp"`  // 最近60日間の取得又は処分の状況
	BrwList  []EdinetLargeVolumeShareholderBorrowing `json:"BrwList"`  // 借入金の内訳
	CredList []EdinetLargeVolumeShareholderCreditor  `json:"CredList"` // 借入先の名称等
}

// EdinetLargeVolumeShareholderAcqDisp は最近60日間の取得又は処分の状況1件分を表します。
type EdinetLargeVolumeShareholderAcqDisp struct {
	Date        string   `json:"Date"`        // 年月日（YYYY-MM-DD形式）
	SecType     string   `json:"SecType"`     // 株券等の種類（例: 普通株式）
	Shs         float64  `json:"Shs"`         // 数量（株）
	Ratio       *float64 `json:"Ratio"`       // 割合（%表記。他の保有割合と異なり小数表現ではない）
	Mkt         string   `json:"Mkt"`         // 市場内外取引の別（書類記載の生値）
	MktCode     string   `json:"MktCode"`     // 市場内外取引コード（1=市場内、2=市場外）
	TxnType     string   `json:"TxnType"`     // 取得又は処分の別（書類記載の生値）
	TxnTypeCode string   `json:"TxnTypeCode"` // 取得又は処分コード（1=取得、2=処分）
	Cptty       string   `json:"Cptty"`       // 譲渡の相手方（短期大量譲渡変更の書類でのみ記載）
	Price       *float64 `json:"Price"`       // 単価（円）
	PriceRaw    string   `json:"PriceRaw"`    // 単価を数値化できない場合の生値
}

// EdinetLargeVolumeShareholderBorrowing は借入金の内訳1件分を表します。
type EdinetLargeVolumeShareholderBorrowing struct {
	Name        string   `json:"Name"`        // 名称（支店名を含む）
	Ind         string   `json:"Ind"`         // 業種
	Rep         string   `json:"Rep"`         // 代表者氏名
	Addr        string   `json:"Addr"`        // 所在地
	DiscBrwPurp string   `json:"DiscBrwPurp"` // 借入目的の開示区分（1=銀行等に開示せず、2=銀行等に開示及び銀行等以外の借入）
	Amt         *float64 `json:"Amt"`         // 金額（円）
}

// EdinetLargeVolumeShareholderCreditor は借入先の名称等1件分を表します。
type EdinetLargeVolumeShareholderCreditor struct {
	Name string `json:"Name"` // 名称（支店名を含む）
	Rep  string `json:"Rep"`  // 代表者氏名
	Addr string `json:"Addr"` // 所在地
}

// RawEdinetLargeVolumeShareholderDoc is used for unmarshaling JSON response with mixed types
type RawEdinetLargeVolumeShareholderDoc struct {
	// 書類メタ情報
	DocId             string `json:"DocId"`
	Code              string `json:"Code"`
	EdinetCode        string `json:"EdinetCode"`
	IsrName           string `json:"IsrName"`
	DocTypeCode       string `json:"DocTypeCode"`
	SubDate           string `json:"SubDate"`
	SubTime           string `json:"SubTime"`
	LargeHldgTypeCode string `json:"LargeHldgTypeCode"`
	DocTitle          string `json:"DocTitle"`
	ChgRsn            string `json:"ChgRsn"`

	// 保有状況の合計
	TotalShsHeld      types.NullableFloat64 `json:"TotalShsHeld"`
	TotalShsRatio     types.NullableFloat64 `json:"TotalShsRatio"`
	TotalShsRatioLast types.NullableFloat64 `json:"TotalShsRatioLast"`
	TotalOutStks      types.NullableFloat64 `json:"TotalOutStks"`

	// 提出者及び共同保有者のレコード
	Hldrs []RawEdinetLargeVolumeShareholderHolder `json:"Hldrs"`
}

// RawEdinetLargeVolumeShareholderHolder is used for unmarshaling JSON response with mixed types
type RawEdinetLargeVolumeShareholderHolder struct {
	HldrName          string `json:"HldrName"`
	HldrNameEn        string `json:"HldrNameEn"`
	HldrEdinetCode    string `json:"HldrEdinetCode"`
	HldrCode          string `json:"HldrCode"`
	LargeHldrTypeCode string `json:"LargeHldrTypeCode"`
	LargeHldrTypeRaw  string `json:"LargeHldrTypeRaw"`

	HldgPurp     string                `json:"HldgPurp"`
	ImpProp      string                `json:"ImpProp"`
	ColAgr       string                `json:"ColAgr"`
	ShsHeld      types.NullableFloat64 `json:"ShsHeld"`
	ShsRatio     types.NullableFloat64 `json:"ShsRatio"`
	ShsRatioLast types.NullableFloat64 `json:"ShsRatioLast"`

	OwnFund    types.NullableFloat64 `json:"OwnFund"`
	TotalBrw   types.NullableFloat64 `json:"TotalBrw"`
	TotalOther types.NullableFloat64 `json:"TotalOther"`
	OtherBrk   string                `json:"OtherBrk"`
	TotalFund  types.NullableFloat64 `json:"TotalFund"`

	AcqDisp  []RawEdinetLargeVolumeShareholderAcqDisp   `json:"AcqDisp"`
	BrwList  []RawEdinetLargeVolumeShareholderBorrowing `json:"BrwList"`
	CredList []EdinetLargeVolumeShareholderCreditor     `json:"CredList"`
}

// RawEdinetLargeVolumeShareholderAcqDisp is used for unmarshaling JSON response with mixed types
type RawEdinetLargeVolumeShareholderAcqDisp struct {
	Date        string                `json:"Date"`
	SecType     string                `json:"SecType"`
	Shs         types.NullableFloat64 `json:"Shs"`
	Ratio       types.NullableFloat64 `json:"Ratio"`
	Mkt         string                `json:"Mkt"`
	MktCode     string                `json:"MktCode"`
	TxnType     string                `json:"TxnType"`
	TxnTypeCode string                `json:"TxnTypeCode"`
	Cptty       string                `json:"Cptty"`
	Price       types.NullableFloat64 `json:"Price"`
	PriceRaw    string                `json:"PriceRaw"`
}

// RawEdinetLargeVolumeShareholderBorrowing is used for unmarshaling JSON response with mixed types
type RawEdinetLargeVolumeShareholderBorrowing struct {
	Name        string                `json:"Name"`
	Ind         string                `json:"Ind"`
	Rep         string                `json:"Rep"`
	Addr        string                `json:"Addr"`
	DiscBrwPurp string                `json:"DiscBrwPurp"`
	Amt         types.NullableFloat64 `json:"Amt"`
}

// UnmarshalJSON implements custom JSON unmarshaling for EdinetLargeVolumeShareholdersResponse
func (r *EdinetLargeVolumeShareholdersResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawEdinetLargeVolumeShareholderDoc
	type rawResponse struct {
		Data          []RawEdinetLargeVolumeShareholderDoc `json:"data"`
		PaginationKey string                               `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawEdinetLargeVolumeShareholderDoc to EdinetLargeVolumeShareholderDoc
	r.Data = make([]EdinetLargeVolumeShareholderDoc, len(raw.Data))
	for idx, rd := range raw.Data {
		r.Data[idx] = rd.toDoc()
	}

	return nil
}

// toDoc converts RawEdinetLargeVolumeShareholderDoc to EdinetLargeVolumeShareholderDoc
func (rd RawEdinetLargeVolumeShareholderDoc) toDoc() EdinetLargeVolumeShareholderDoc {
	holders := make([]EdinetLargeVolumeShareholderHolder, len(rd.Hldrs))
	for idx, rh := range rd.Hldrs {
		holders[idx] = rh.toHolder()
	}

	return EdinetLargeVolumeShareholderDoc{
		// 書類メタ情報
		DocId:             rd.DocId,
		Code:              rd.Code,
		EdinetCode:        rd.EdinetCode,
		IsrName:           rd.IsrName,
		DocTypeCode:       rd.DocTypeCode,
		SubDate:           rd.SubDate,
		SubTime:           rd.SubTime,
		LargeHldgTypeCode: rd.LargeHldgTypeCode,
		DocTitle:          rd.DocTitle,
		ChgRsn:            rd.ChgRsn,

		// 保有状況の合計
		TotalShsHeld:      rd.TotalShsHeld.Ptr(),
		TotalShsRatio:     rd.TotalShsRatio.Ptr(),
		TotalShsRatioLast: rd.TotalShsRatioLast.Ptr(),
		TotalOutStks:      rd.TotalOutStks.Or(0),

		// 提出者及び共同保有者のレコード
		Hldrs: holders,
	}
}

// toHolder converts RawEdinetLargeVolumeShareholderHolder to EdinetLargeVolumeShareholderHolder
func (rh RawEdinetLargeVolumeShareholderHolder) toHolder() EdinetLargeVolumeShareholderHolder {
	acqDisps := make([]EdinetLargeVolumeShareholderAcqDisp, len(rh.AcqDisp))
	for idx, ra := range rh.AcqDisp {
		acqDisps[idx] = EdinetLargeVolumeShareholderAcqDisp{
			Date:        ra.Date,
			SecType:     ra.SecType,
			Shs:         ra.Shs.Or(0),
			Ratio:       ra.Ratio.Ptr(),
			Mkt:         ra.Mkt,
			MktCode:     ra.MktCode,
			TxnType:     ra.TxnType,
			TxnTypeCode: ra.TxnTypeCode,
			Cptty:       ra.Cptty,
			Price:       ra.Price.Ptr(),
			PriceRaw:    ra.PriceRaw,
		}
	}

	borrowings := make([]EdinetLargeVolumeShareholderBorrowing, len(rh.BrwList))
	for idx, rb := range rh.BrwList {
		borrowings[idx] = EdinetLargeVolumeShareholderBorrowing{
			Name:        rb.Name,
			Ind:         rb.Ind,
			Rep:         rb.Rep,
			Addr:        rb.Addr,
			DiscBrwPurp: rb.DiscBrwPurp,
			Amt:         rb.Amt.Ptr(),
		}
	}

	return EdinetLargeVolumeShareholderHolder{
		// 保有者情報
		HldrName:          rh.HldrName,
		HldrNameEn:        rh.HldrNameEn,
		HldrEdinetCode:    rh.HldrEdinetCode,
		HldrCode:          rh.HldrCode,
		LargeHldrTypeCode: rh.LargeHldrTypeCode,
		LargeHldrTypeRaw:  rh.LargeHldrTypeRaw,

		// 保有状況
		HldgPurp:     rh.HldgPurp,
		ImpProp:      rh.ImpProp,
		ColAgr:       rh.ColAgr,
		ShsHeld:      rh.ShsHeld.Or(0),
		ShsRatio:     rh.ShsRatio.Or(0),
		ShsRatioLast: rh.ShsRatioLast.Ptr(),

		// 取得資金
		OwnFund:    rh.OwnFund.Ptr(),
		TotalBrw:   rh.TotalBrw.Ptr(),
		TotalOther: rh.TotalOther.Ptr(),
		OtherBrk:   rh.OtherBrk,
		TotalFund:  rh.TotalFund.Ptr(),

		// 明細
		AcqDisp:  acqDisps,
		BrwList:  borrowings,
		CredList: rh.CredList,
	}
}

// GetLargeVolumeShareholders は大量保有報告書を取得します。
// パラメータをすべて省略した場合はAPI実行日に提出された全書類のデータ一覧を取得します。
func (s *EdinetLargeVolumeShareholdersService) GetLargeVolumeShareholders(ctx context.Context, params EdinetLargeVolumeShareholdersParams) (*EdinetLargeVolumeShareholdersResponse, error) {
	// edinet_codeとcodeの同時指定は不可
	if params.EdinetCode != "" && params.Code != "" {
		return nil, fmt.Errorf("edinet_code and code cannot be specified together")
	}

	path := "/edinet/large-volume-shareholders"

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

	var resp EdinetLargeVolumeShareholdersResponse
	if err := s.client.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get large volume shareholders: %w", err)
	}

	return &resp, nil
}

// GetLargeVolumeShareholdersByCode は指定銘柄の大量保有報告書を取得します。
// ページネーションを使用して全データを取得します。
func (s *EdinetLargeVolumeShareholdersService) GetLargeVolumeShareholdersByCode(ctx context.Context, code string) ([]EdinetLargeVolumeShareholderDoc, error) {
	var allData []EdinetLargeVolumeShareholderDoc
	paginationKey := ""

	for {
		params := EdinetLargeVolumeShareholdersParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetLargeVolumeShareholders(ctx, params)
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

// GetLargeVolumeShareholdersByEdinetCode は指定EDINETコードの発行者に係る大量保有報告書を取得します。
// ページネーションを使用して全データを取得します。
func (s *EdinetLargeVolumeShareholdersService) GetLargeVolumeShareholdersByEdinetCode(ctx context.Context, edinetCode string) ([]EdinetLargeVolumeShareholderDoc, error) {
	var allData []EdinetLargeVolumeShareholderDoc
	paginationKey := ""

	for {
		params := EdinetLargeVolumeShareholdersParams{
			EdinetCode:    edinetCode,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetLargeVolumeShareholders(ctx, params)
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

// GetLargeVolumeShareholdersByDate は指定提出日の全書類の大量保有報告書を取得します。
// ページネーションを使用して全データを取得します。
func (s *EdinetLargeVolumeShareholdersService) GetLargeVolumeShareholdersByDate(ctx context.Context, date string) ([]EdinetLargeVolumeShareholderDoc, error) {
	var allData []EdinetLargeVolumeShareholderDoc
	paginationKey := ""

	for {
		params := EdinetLargeVolumeShareholdersParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetLargeVolumeShareholders(ctx, params)
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

// IsChangeReport は変更報告書（変更報告書、短期大量譲渡、特例対象株券等の変更報告書）かを判定します。
func (d *EdinetLargeVolumeShareholderDoc) IsChangeReport() bool {
	switch d.LargeHldgTypeCode {
	case LargeHldgTypeCodeChangeReport, LargeHldgTypeCodeChangeReportShortTermTransfer, LargeHldgTypeCodeSpecialChangeReport:
		return true
	}
	return false
}

// GetTotalShareholdingPercentage は株券等保有割合の合計をパーセント表記で取得します。
// 合計欄の記載がない書類（TotalShsRatioがnull）の場合は0を返します。
func (d *EdinetLargeVolumeShareholderDoc) GetTotalShareholdingPercentage() float64 {
	if d.TotalShsRatio == nil {
		return 0
	}
	return *d.TotalShsRatio * 100
}

// GetShareholdingPercentage は株券等保有割合をパーセント表記で取得します。
func (h *EdinetLargeVolumeShareholderHolder) GetShareholdingPercentage() float64 {
	return h.ShsRatio * 100
}
