package jquants

import (
	"context"
	"fmt"

	"github.com/utahta/jquants/client"
)

// EarningsDateService は決算発表予定日を取得するサービスです。
// 決算期によらず、東証に報告を行った全上場銘柄（REIT等を含む）が対象で、
// 予定日の変更・未定の履歴も公表日単位で提供します。
type EarningsDateService struct {
	client client.HTTPClient
}

// NewEarningsDateService は新しいEarningsDateServiceを作成します。
func NewEarningsDateService(c client.HTTPClient) *EarningsDateService {
	return &EarningsDateService{client: c}
}

// EarningsDateParams は決算発表予定日のリクエストパラメータです。
// Code、Date、ScheduledDateのいずれか1つの指定が必須で、同時指定はできません。
type EarningsDateParams struct {
	Code          string // 銘柄コード（4桁または5桁）
	Date          string // 公表日（YYYYMMDD または YYYY-MM-DD）
	ScheduledDate string // 決算発表予定日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// EarningsDateResponse は決算発表予定日のレスポンスです。
type EarningsDateResponse struct {
	Data          []EarningsDate `json:"data"`
	PaginationKey string         `json:"pagination_key"` // ページネーションキー
}

// EarningsDate は決算発表予定日を表します。
// J-Quants API /fins/earnings-date エンドポイントのレスポンスデータ。
// 予定日が変更された場合は以前のレコードが残ったまま新たなレコードが追加されるため、
// 銘柄コード指定時は変更履歴を含む全レコードが返されます。
type EarningsDate struct {
	PubDate  string `json:"PubDate"`  // 公表日（YYYY-MM-DD形式）。この予定日が公表・変更された日
	SchDate  string `json:"SchDate"`  // 決算発表予定日（YYYY-MM-DD形式）。未定の場合は空文字
	FQName   string `json:"FQName"`   // 決算区分（1Q / 2Q / 3Q / FY）
	FYE      string `json:"FYE"`      // 決算期末（MMDD形式）
	Code     string `json:"Code"`     // 銘柄コード（5桁）
	CoName   string `json:"CoName"`   // 会社名（PubDate時点）
	CoNameEn string `json:"CoNameEn"` // 会社名（英語、PubDate時点）
}

// GetEarningsDates は指定された条件で決算発表予定日を取得します。
// Code、Date、ScheduledDateのいずれか1つの指定が必須です。
// パラメータ:
// - Code: 銘柄コード（例: "8697" または "86970"）。予定日の公表履歴を取得
// - Date: 公表日（例: "20250620" または "2025-06-20"）。その日に公表・変更された全銘柄を取得
// - ScheduledDate: 決算発表予定日（例: "20250805" または "2025-08-05"）。その日を現在有効な予定日とする全銘柄を取得
// - PaginationKey: ページネーション用キー
func (s *EarningsDateService) GetEarningsDates(ctx context.Context, params EarningsDateParams) (*EarningsDateResponse, error) {
	specified := 0
	for _, v := range []string{params.Code, params.Date, params.ScheduledDate} {
		if v != "" {
			specified++
		}
	}
	if specified != 1 {
		return nil, fmt.Errorf("exactly one of code, date or scheduled_date is required")
	}

	path := "/fins/earnings-date"

	query := "?"
	if params.Code != "" {
		query += fmt.Sprintf("code=%s&", params.Code)
	}
	if params.Date != "" {
		query += fmt.Sprintf("date=%s&", params.Date)
	}
	if params.ScheduledDate != "" {
		query += fmt.Sprintf("scheduled_date=%s&", params.ScheduledDate)
	}
	if params.PaginationKey != "" {
		query += fmt.Sprintf("pagination_key=%s&", params.PaginationKey)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp EarningsDateResponse
	if err := s.client.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get earnings dates: %w", err)
	}

	return &resp, nil
}

// GetEarningsDatesByCode は指定銘柄の決算発表予定日の公表履歴を取得します。
// ページネーションを使用して全データを取得します。
func (s *EarningsDateService) GetEarningsDatesByCode(ctx context.Context, code string) ([]EarningsDate, error) {
	return s.getAllEarningsDates(ctx, EarningsDateParams{Code: code})
}

// GetEarningsDatesByDate は指定日に公表・変更された全銘柄の決算発表予定日を取得します。
// ページネーションを使用して全データを取得します。
func (s *EarningsDateService) GetEarningsDatesByDate(ctx context.Context, date string) ([]EarningsDate, error) {
	return s.getAllEarningsDates(ctx, EarningsDateParams{Date: date})
}

// GetEarningsDatesByScheduledDate は指定日を現在有効な決算発表予定日とする全銘柄を取得します。
// 予定日がその後変更された銘柄は、変更前の予定日ではヒットしません。
// ページネーションを使用して全データを取得します。
func (s *EarningsDateService) GetEarningsDatesByScheduledDate(ctx context.Context, scheduledDate string) ([]EarningsDate, error) {
	return s.getAllEarningsDates(ctx, EarningsDateParams{ScheduledDate: scheduledDate})
}

func (s *EarningsDateService) getAllEarningsDates(ctx context.Context, params EarningsDateParams) ([]EarningsDate, error) {
	var allEarningsDates []EarningsDate
	paginationKey := ""

	for {
		params.PaginationKey = paginationKey

		resp, err := s.GetEarningsDates(ctx, params)
		if err != nil {
			return nil, err
		}

		allEarningsDates = append(allEarningsDates, resp.Data...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allEarningsDates, nil
}

// IsUndetermined は決算発表予定日が未定かどうかを判定します。
func (e *EarningsDate) IsUndetermined() bool {
	return e.SchDate == ""
}
