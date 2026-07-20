package jquants

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// BulkService はCSVファイルの一括ダウンロード（Bulk API）を提供するサービスです。
// ご契約プランでアクセス可能なエンドポイントのファイルが対象です。
// 株式ティックデータ（/equities/trades）はAPIでの提供がなく、このBulk API経由でのみ取得できます。
type BulkService struct {
	client client.HTTPClient
}

// NewBulkService は新しいBulkServiceを作成します。
func NewBulkService(c client.HTTPClient) *BulkService {
	return &BulkService{client: c}
}

// BulkListParams はダウンロード可能ファイル一覧のリクエストパラメータです。
type BulkListParams struct {
	Endpoint string // 対象データのエンドポイント名（e.g. /equities/bars/daily）（endpoint、dateのいずれかが必須）
	Date     string // 対象日付（YYYY-MM、YYYYMM、YYYY-MM-DD、YYYYMMDD）（endpoint、dateのいずれかが必須）
	From     string // 取得期間の開始日（YYYY-MM、YYYYMM、YYYY-MM-DD、YYYYMMDD）（endpoint指定時のみ使用可）
	To       string // 取得期間の終了日（YYYY-MM、YYYYMM、YYYY-MM-DD、YYYYMMDD）（endpoint指定時のみ使用可）
}

// BulkListResponse はダウンロード可能ファイル一覧のレスポンスです。
type BulkListResponse struct {
	Data []BulkFile `json:"data"`
}

// BulkFile はダウンロード可能ファイルの情報を表します。
// J-Quants API /bulk/list エンドポイントのレスポンスデータ。
type BulkFile struct {
	Key          string `json:"Key"`          // ファイルのキー（/bulk/get でのダウンロード時に使用。e.g. equities/bars/daily/historical/2025/equities_bars_daily_202501.csv.gz）
	LastModified string `json:"LastModified"` // 最終更新日時（ISO 8601形式）
	Size         int64  `json:"Size"`         // ファイルサイズ（バイト）
}

// RawBulkFile is used for unmarshaling JSON response with mixed types
type RawBulkFile struct {
	Key          string              `json:"Key"`
	LastModified string              `json:"LastModified"`
	Size         types.NullableInt64 `json:"Size"`
}

// UnmarshalJSON implements custom JSON unmarshaling for BulkListResponse
func (r *BulkListResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawBulkFile
	type rawResponse struct {
		Data []RawBulkFile `json:"data"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Convert RawBulkFile to BulkFile
	r.Data = make([]BulkFile, len(raw.Data))
	for idx, rbf := range raw.Data {
		r.Data[idx] = BulkFile{
			Key:          rbf.Key,
			LastModified: rbf.LastModified,
			Size:         rbf.Size.Or(0),
		}
	}

	return nil
}

// BulkGetParams はファイルダウンロード用URL取得のリクエストパラメータです。
// key単独、またはendpointとdateの組み合わせのどちらかを指定します（keyとendpoint/dateの併用は不可）。
type BulkGetParams struct {
	Key      string // ファイルのキー（/bulk/list で取得したKey）
	Endpoint string // 対象データのエンドポイント名（e.g. /equities/bars/daily。dateと組み合わせて使用）
	Date     string // 対象日付（YYYY-MM、YYYYMM、YYYY-MM-DD、YYYYMMDD。endpointと組み合わせて使用）
}

// GetFiles はダウンロード可能ファイル一覧を取得します。
// endpointとして取引カレンダー（/markets/calendar）を指定した場合、from/toの期間に関わらず最新の1ファイルのみが返されます。
func (s *BulkService) GetFiles(ctx context.Context, params BulkListParams) (*BulkListResponse, error) {
	// endpoint、dateのいずれかが必須
	if params.Endpoint == "" && params.Date == "" {
		return nil, fmt.Errorf("either endpoint or date parameter is required")
	}
	// from/toはendpoint指定時のみ使用可能
	if (params.From != "" || params.To != "") && params.Endpoint == "" {
		return nil, fmt.Errorf("from/to parameters can only be used with endpoint parameter")
	}

	path := "/bulk/list"

	query := "?"
	if params.Endpoint != "" {
		query += fmt.Sprintf("endpoint=%s&", params.Endpoint)
	}
	if params.Date != "" {
		query += fmt.Sprintf("date=%s&", params.Date)
	}
	if params.From != "" {
		query += fmt.Sprintf("from=%s&", params.From)
	}
	if params.To != "" {
		query += fmt.Sprintf("to=%s&", params.To)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp BulkListResponse
	if err := s.client.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get bulk file list: %w", err)
	}

	return &resp, nil
}

// GetDownloadURL はファイルダウンロード用の署名付きURLを取得します。
// 取得したURLの有効期限は5分で、再利用はできません。ファイルはgzip形式で圧縮されたCSVです。
func (s *BulkService) GetDownloadURL(ctx context.Context, params BulkGetParams) (string, error) {
	// key単独、またはendpointとdateの組み合わせのどちらかが必須
	if params.Key != "" && (params.Endpoint != "" || params.Date != "") {
		return "", fmt.Errorf("key parameter cannot be combined with endpoint or date parameter")
	}
	if params.Key == "" && (params.Endpoint == "" || params.Date == "") {
		return "", fmt.Errorf("either key or endpoint and date parameters are required")
	}

	path := "/bulk/get"

	query := "?"
	if params.Key != "" {
		query += fmt.Sprintf("key=%s&", params.Key)
	}
	if params.Endpoint != "" {
		query += fmt.Sprintf("endpoint=%s&", params.Endpoint)
	}
	if params.Date != "" {
		query += fmt.Sprintf("date=%s&", params.Date)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp struct {
		URL string `json:"url"`
	}
	// 署名付きURLは5分で失効し1回限りのため、キャッシュを経由しない
	if err := client.DoRequestNoCache(ctx, s.client, "GET", path, nil, &resp); err != nil {
		return "", fmt.Errorf("failed to get bulk download URL: %w", err)
	}

	return resp.URL, nil
}
