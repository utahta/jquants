package jquants

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// EdinetMajorShareholdersService は大株主状況（EDINET）を取得するサービスです。
// 有価証券報告書に記載されている大株主の状況を取得します。
type EdinetMajorShareholdersService struct {
	client client.HTTPClient
}

// NewEdinetMajorShareholdersService は新しいEdinetMajorShareholdersServiceを作成します。
func NewEdinetMajorShareholdersService(c client.HTTPClient) *EdinetMajorShareholdersService {
	return &EdinetMajorShareholdersService{client: c}
}

// EdinetMajorShareholdersParams は大株主状況のリクエストパラメータです。
// すべて任意ですが、edinet_codeとcodeの同時指定はできません。
// すべて省略した場合はAPI実行日に提出された全有報のデータ一覧が返ります。
type EdinetMajorShareholdersParams struct {
	EdinetCode    string // EDINETコード（例: E03814）（codeとの同時指定は不可）
	Code          string // 4桁もしくは5桁の銘柄コード（edinet_codeとの同時指定は不可）
	Date          string // 提出日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// EdinetMajorShareholdersResponse は大株主状況のレスポンスです。
type EdinetMajorShareholdersResponse struct {
	Data          []EdinetMajorShareholderDoc `json:"data"`
	PaginationKey string                      `json:"pagination_key"` // ページネーションキー
}

// EdinetMajorShareholderDoc は有価証券報告書ごとの大株主状況データを表します。
// J-Quants API /edinet/major-shareholders エンドポイントのレスポンスデータ。
// データは2016年6月1日以降、対象書類は有価証券報告書 第三号様式。Standardプラン以上で利用可能。
type EdinetMajorShareholderDoc struct {
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

	// 大株主レコード（順位順、Rank昇順）
	Hldrs []EdinetMajorShareholderHolder `json:"Hldrs"`
}

// EdinetMajorShareholderHolder は大株主1名分のデータを表します。
// 通常は上位10名ですが、同順位タイで11位以降が記載される場合や1名のみの場合もあります。
type EdinetMajorShareholderHolder struct {
	Rank     int     `json:"Rank"`     // 順位（1〜10、同順位タイで11以降あり）
	HldrName string  `json:"HldrName"` // 株主氏名又は名称
	HldrAddr string  `json:"HldrAddr"` // 株主住所
	ShsHeld  float64 `json:"ShsHeld"`  // 所有株式数（株）
	ShsRatio float64 `json:"ShsRatio"` // 発行済株式（自己株式を除く）に対する所有割合（小数表現。0.1881=18.81%）
}

// RawEdinetMajorShareholderDoc is used for unmarshaling JSON response with mixed types
type RawEdinetMajorShareholderDoc struct {
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

	// 大株主レコード
	Hldrs []RawEdinetMajorShareholderHolder `json:"Hldrs"`
}

// RawEdinetMajorShareholderHolder is used for unmarshaling JSON response with mixed types
type RawEdinetMajorShareholderHolder struct {
	Rank     types.NullableInt64   `json:"Rank"`
	HldrName string                `json:"HldrName"`
	HldrAddr string                `json:"HldrAddr"`
	ShsHeld  types.NullableFloat64 `json:"ShsHeld"`
	ShsRatio types.NullableFloat64 `json:"ShsRatio"`
}

// UnmarshalJSON implements custom JSON unmarshaling for EdinetMajorShareholdersResponse
func (r *EdinetMajorShareholdersResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawEdinetMajorShareholderDoc
	type rawResponse struct {
		Data          []RawEdinetMajorShareholderDoc `json:"data"`
		PaginationKey string                         `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawEdinetMajorShareholderDoc to EdinetMajorShareholderDoc
	r.Data = make([]EdinetMajorShareholderDoc, len(raw.Data))
	for idx, rd := range raw.Data {
		holders := make([]EdinetMajorShareholderHolder, len(rd.Hldrs))
		for hIdx, rh := range rd.Hldrs {
			holders[hIdx] = EdinetMajorShareholderHolder{
				Rank:     int(rh.Rank.Or(0)),
				HldrName: rh.HldrName,
				HldrAddr: rh.HldrAddr,
				ShsHeld:  rh.ShsHeld.Or(0),
				ShsRatio: rh.ShsRatio.Or(0),
			}
		}

		r.Data[idx] = EdinetMajorShareholderDoc{
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

			// 大株主レコード
			Hldrs: holders,
		}
	}

	return nil
}

// GetMajorShareholders は大株主状況を取得します。
// パラメータをすべて省略した場合はAPI実行日に提出された全有報のデータ一覧を取得します。
func (s *EdinetMajorShareholdersService) GetMajorShareholders(ctx context.Context, params EdinetMajorShareholdersParams) (*EdinetMajorShareholdersResponse, error) {
	// edinet_codeとcodeの同時指定は不可
	if params.EdinetCode != "" && params.Code != "" {
		return nil, fmt.Errorf("edinet_code and code cannot be specified together")
	}

	path := "/edinet/major-shareholders"

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

	var resp EdinetMajorShareholdersResponse
	if err := s.client.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get major shareholders: %w", err)
	}

	return &resp, nil
}

// GetMajorShareholdersByCode は指定銘柄の大株主状況を取得します。
// ページネーションを使用して全データを取得します。
func (s *EdinetMajorShareholdersService) GetMajorShareholdersByCode(ctx context.Context, code string) ([]EdinetMajorShareholderDoc, error) {
	var allData []EdinetMajorShareholderDoc
	paginationKey := ""

	for {
		params := EdinetMajorShareholdersParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetMajorShareholders(ctx, params)
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

// GetMajorShareholdersByEdinetCode は指定EDINETコードの大株主状況を取得します。
// ページネーションを使用して全データを取得します。
func (s *EdinetMajorShareholdersService) GetMajorShareholdersByEdinetCode(ctx context.Context, edinetCode string) ([]EdinetMajorShareholderDoc, error) {
	var allData []EdinetMajorShareholderDoc
	paginationKey := ""

	for {
		params := EdinetMajorShareholdersParams{
			EdinetCode:    edinetCode,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetMajorShareholders(ctx, params)
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

// GetMajorShareholdersByDate は指定提出日の全有報の大株主状況を取得します。
// ページネーションを使用して全データを取得します。
func (s *EdinetMajorShareholdersService) GetMajorShareholdersByDate(ctx context.Context, date string) ([]EdinetMajorShareholderDoc, error) {
	var allData []EdinetMajorShareholderDoc
	paginationKey := ""

	for {
		params := EdinetMajorShareholdersParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetMajorShareholders(ctx, params)
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

// GetShareholdingPercentage は所有割合をパーセント表記で取得します。
func (h *EdinetMajorShareholderHolder) GetShareholdingPercentage() float64 {
	return h.ShsRatio * 100
}
