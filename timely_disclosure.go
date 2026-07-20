package jquants

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// 適時開示書類タイプ（Docs・docsパラメータで使用する値）
const (
	TimelyDisclosureDocFullPDF    = "g" // 全文情報PDF
	TimelyDisclosureDocSummaryPDF = "s" // サマリ情報PDF
	TimelyDisclosureDocXBRL       = "x" // XBRL関連ファイル
)

// TimelyDisclosureService はTDnet適時開示情報を取得するサービスです。
// 利用にはTDnet/適時開示情報アドオン契約が必要です。データ取得可能期間は過去5年間です。
type TimelyDisclosureService struct {
	client client.HTTPClient
}

// NewTimelyDisclosureService は新しいTimelyDisclosureServiceを作成します。
func NewTimelyDisclosureService(c client.HTTPClient) *TimelyDisclosureService {
	return &TimelyDisclosureService{client: c}
}

// TimelyDisclosureParams は適時開示インデックス一覧のリクエストパラメータです。
type TimelyDisclosureParams struct {
	Date          string // 開示日（YYYYMMDD または YYYY-MM-DD）（date、codeのいずれかが必須）
	Code          string // 4桁もしくは5桁の銘柄コード（date、codeのいずれかが必須）
	From          string // 取得開始日（YYYYMMDD または YYYY-MM-DD）。codeと組み合わせ、toと必ずペアで指定
	To            string // 取得終了日（YYYYMMDD または YYYY-MM-DD）。codeと組み合わせ、fromと必ずペアで指定
	DiscItems     string // 公開項目コードで絞り込む（カンマ区切りで複数指定可能、AND条件）
	Cursor        string // 差分取得用カーソル（前回レスポンスのcursorを指定）。pagination_keyと同時指定不可
	PaginationKey string // ページネーションキー。cursorと同時指定不可
}

// TimelyDisclosureResponse は適時開示インデックス一覧のレスポンスです。
type TimelyDisclosureResponse struct {
	Data          []TimelyDisclosure `json:"data"`
	Cursor        string             `json:"cursor"`         // 差分取得用カーソル（dateに当日を指定し、ページネーションなしで全件取得できた場合のみ返却）
	PaginationKey string             `json:"pagination_key"` // ページネーションキー
}

// TimelyDisclosure は適時開示インデックスのデータを表します。
// J-Quants API /td/list エンドポイントのレスポンスデータ。
type TimelyDisclosure struct {
	DiscNo     string   `json:"DiscNo"`     // 開示番号（14桁）
	Code       string   `json:"Code"`       // 銘柄コード
	Name       string   `json:"Name"`       // 会社名
	DiscDate   string   `json:"DiscDate"`   // 開示日（YYYY-MM-DD形式）
	DiscTime   string   `json:"DiscTime"`   // 開示時刻（HH:MM形式）
	Title      string   `json:"Title"`      // 開示タイトル
	DiscStatus string   `json:"DiscStatus"` // 取扱属性（""=新規、"revision"=訂正、"delete"=削除。現仕様では常に新規）
	RevNo      int      `json:"RevNo"`      // 開示履歴番号（1～99。現仕様では常に1）
	DiscItems  []string `json:"DiscItems"`  // 公開項目コードのリスト
	Docs       []string `json:"Docs"`       // 書類タイプのリスト（g=全文情報PDF、s=サマリ情報PDF、x=XBRL関連ファイル）
}

// RawTimelyDisclosure is used for unmarshaling JSON response with mixed types
type RawTimelyDisclosure struct {
	DiscNo     string               `json:"DiscNo"`
	Code       string               `json:"Code"`
	Name       string               `json:"Name"`
	DiscDate   string               `json:"DiscDate"`
	DiscTime   string               `json:"DiscTime"`
	Title      string               `json:"Title"`
	DiscStatus types.NullableString `json:"DiscStatus"`
	RevNo      types.NullableInt64  `json:"RevNo"`
	DiscItems  []string             `json:"DiscItems"`
	Docs       []string             `json:"Docs"`
}

// UnmarshalJSON implements custom JSON unmarshaling for TimelyDisclosureResponse
func (r *TimelyDisclosureResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawTimelyDisclosure
	type rawResponse struct {
		Data          []RawTimelyDisclosure `json:"data"`
		Cursor        string                `json:"cursor"`
		PaginationKey string                `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.Cursor = raw.Cursor
	r.PaginationKey = raw.PaginationKey

	// Convert RawTimelyDisclosure to TimelyDisclosure
	r.Data = make([]TimelyDisclosure, len(raw.Data))
	for idx, rtd := range raw.Data {
		r.Data[idx] = TimelyDisclosure{
			DiscNo:     rtd.DiscNo,
			Code:       rtd.Code,
			Name:       rtd.Name,
			DiscDate:   rtd.DiscDate,
			DiscTime:   rtd.DiscTime,
			Title:      rtd.Title,
			DiscStatus: rtd.DiscStatus.Or(""),
			RevNo:      int(rtd.RevNo.Or(0)),
			DiscItems:  rtd.DiscItems,
			Docs:       rtd.Docs,
		}
	}

	return nil
}

// TimelyDisclosureFilesParams は適時開示ファイル取得のリクエストパラメータです。
type TimelyDisclosureFilesParams struct {
	DiscNo string // 開示番号（14桁）（必須）
	Docs   string // 取得するファイル種別（カンマ区切りで複数指定可能。g=全文情報PDF、s=サマリ情報PDF、x=XBRL関連ファイル。省略時は全種別）
}

// TimelyDisclosureFiles は適時開示ファイルのダウンロードURL情報を表します。
// J-Quants API /td/files エンドポイントのレスポンスデータ。
type TimelyDisclosureFiles struct {
	DiscNo string                   `json:"discNo"` // 開示番号（14桁）
	Files  TimelyDisclosureFileURLs `json:"files"`  // ファイルのダウンロードURL一覧
}

// TimelyDisclosureFileURLs は書類タイプごとの署名付きダウンロードURLです。
// URLの有効期限は15分です。docsで絞った場合や書類が存在しない場合、該当フィールドは空文字になります。
type TimelyDisclosureFileURLs struct {
	PDF        string `json:"pdf"`        // 全文情報PDFのダウンロードURL
	SummaryPDF string `json:"summaryPdf"` // サマリ情報PDFのダウンロードURL
	XBRL       string `json:"xbrl"`       // XBRL関連ファイルのダウンロードURL
}

// TimelyDisclosureBulk は適時開示インデックス一括ダウンロード情報を表します。
// J-Quants API /td/bulk エンドポイントのレスポンスデータ。
// URLはgzip圧縮CSVの署名付きダウンロードURLで、有効期限は15分です。
type TimelyDisclosureBulk struct {
	LastUpdated string `json:"lastUpdated"` // CSVファイルの最終更新日時（ISO 8601形式）
	URL         string `json:"url"`         // CSVファイル（gzip圧縮）のダウンロードURL
}

// GetDisclosures は適時開示インデックス一覧を取得します。
func (s *TimelyDisclosureService) GetDisclosures(ctx context.Context, params TimelyDisclosureParams) (*TimelyDisclosureResponse, error) {
	// date、codeのいずれかが必須
	if params.Date == "" && params.Code == "" {
		return nil, fmt.Errorf("either date or code parameter is required")
	}
	// cursorとpagination_keyは同時指定不可
	if params.Cursor != "" && params.PaginationKey != "" {
		return nil, fmt.Errorf("cursor and pagination_key parameters cannot be specified together")
	}

	path := "/td/list"

	query := "?"
	if params.Date != "" {
		query += fmt.Sprintf("date=%s&", params.Date)
	}
	if params.Code != "" {
		query += fmt.Sprintf("code=%s&", params.Code)
	}
	if params.From != "" {
		query += fmt.Sprintf("from=%s&", params.From)
	}
	if params.To != "" {
		query += fmt.Sprintf("to=%s&", params.To)
	}
	if params.DiscItems != "" {
		query += fmt.Sprintf("discItems=%s&", params.DiscItems)
	}
	if params.Cursor != "" {
		query += fmt.Sprintf("cursor=%s&", params.Cursor)
	}
	if params.PaginationKey != "" {
		query += fmt.Sprintf("pagination_key=%s&", params.PaginationKey)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp TimelyDisclosureResponse
	var err error
	if params.Cursor != "" {
		// cursorによる差分取得はポーリング用途のため、キャッシュを経由しない
		err = client.DoRequestNoCache(ctx, s.client, "GET", path, nil, &resp)
	} else {
		err = s.client.DoRequest(ctx, "GET", path, nil, &resp)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get timely disclosures: %w", err)
	}

	return &resp, nil
}

// GetDisclosuresByDate は指定開示日の適時開示インデックス一覧を取得します。
// ページネーションを使用して全データを取得します。
func (s *TimelyDisclosureService) GetDisclosuresByDate(ctx context.Context, date string) ([]TimelyDisclosure, error) {
	var allData []TimelyDisclosure
	paginationKey := ""

	for {
		params := TimelyDisclosureParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDisclosures(ctx, params)
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

// GetDisclosuresByCode は指定銘柄の適時開示インデックス一覧を取得します（直近5年間）。
// ページネーションを使用して全データを取得します。
func (s *TimelyDisclosureService) GetDisclosuresByCode(ctx context.Context, code string) ([]TimelyDisclosure, error) {
	var allData []TimelyDisclosure
	paginationKey := ""

	for {
		params := TimelyDisclosureParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDisclosures(ctx, params)
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

// GetDisclosureFiles は開示番号に対応する適時開示ファイルのダウンロードURLを取得します。
// URLの有効期限は15分です。
func (s *TimelyDisclosureService) GetDisclosureFiles(ctx context.Context, params TimelyDisclosureFilesParams) (*TimelyDisclosureFiles, error) {
	// discNoは必須
	if params.DiscNo == "" {
		return nil, fmt.Errorf("discNo parameter is required")
	}

	path := "/td/files"

	query := "?"
	query += fmt.Sprintf("discNo=%s&", params.DiscNo)
	if params.Docs != "" {
		query += fmt.Sprintf("docs=%s&", params.Docs)
	}

	path += query[:len(query)-1] // Remove trailing &

	var resp TimelyDisclosureFiles
	// 署名付きURLは15分で失効するため、キャッシュを経由しない
	if err := client.DoRequestNoCache(ctx, s.client, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get timely disclosure files: %w", err)
	}

	return &resp, nil
}

// GetBulkFile は全開示インデックスCSV（gzip圧縮）のダウンロードURLを取得します。
// URLの有効期限は15分です。
func (s *TimelyDisclosureService) GetBulkFile(ctx context.Context) (*TimelyDisclosureBulk, error) {
	var resp TimelyDisclosureBulk
	// 署名付きURLは15分で失効するため、キャッシュを経由しない
	if err := client.DoRequestNoCache(ctx, s.client, "GET", "/td/bulk", nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get timely disclosure bulk file: %w", err)
	}

	return &resp, nil
}

// IsRevision は訂正開示情報かを判定します。
func (td *TimelyDisclosure) IsRevision() bool {
	return td.DiscStatus == "revision"
}

// IsDeleted は削除開示情報かを判定します。
func (td *TimelyDisclosure) IsDeleted() bool {
	return td.DiscStatus == "delete"
}

// HasXBRL はXBRL関連ファイルが存在するかを判定します。
func (td *TimelyDisclosure) HasXBRL() bool {
	return slices.Contains(td.Docs, TimelyDisclosureDocXBRL)
}
