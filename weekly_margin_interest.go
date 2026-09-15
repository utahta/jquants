package jquants

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// WeeklyMarginInterestService は信用取引残高を取得するサービスです。
type WeeklyMarginInterestService struct {
	client client.HTTPClient
}

// NewWeeklyMarginInterestService は新しいWeeklyMarginInterestServiceを作成します。
func NewWeeklyMarginInterestService(c client.HTTPClient) *WeeklyMarginInterestService {
	return &WeeklyMarginInterestService{client: c}
}

// WeeklyMarginInterestParams はCode、Date、PublishedDateのいずれかが必須です。
type WeeklyMarginInterestParams struct {
	Code          string // 銘柄コード
	Date          string // 申込日付（YYYYMMDD または YYYY-MM-DD）
	PublishedDate string // 公表日（Date/From/Toとの同時指定不可）
	From          string // 期間の開始日（YYYYMMDD または YYYY-MM-DD）
	To            string // 期間の終了日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// WeeklyMarginInterestResponse は信用取引残高のレスポンスです。
type WeeklyMarginInterestResponse struct {
	Data          []WeeklyMarginInterest `json:"data"`
	PaginationKey string                 `json:"pagination_key"` // ページネーションキー
}

// 銘柄区分定数
const (
	IssueTypeCredit   = "1" // 信用銘柄（制度信用取引のみ可能）
	IssueTypeLendable = "2" // 貸借銘柄（制度信用取引・貸借取引ともに可能）
	IssueTypeOther    = "3" // その他（一般信用取引のみまたは取引不可）
)

// WeeklyMarginInterest は信用取引残高のデータを表します。
// J-Quants API /markets/margin-interest エンドポイントのレスポンスデータ。
// 公表日と金額は2026年9月25日申込分以降が対象で、それ以前はnilです。
type WeeklyMarginInterest struct {
	// 基本情報
	PubDate *string `json:"PubDate"` // 公表日（YYYY-MM-DD形式）
	Date    string  `json:"Date"`    // 申込日付（YYYY-MM-DD形式）
	Code    string  `json:"Code"`    // 銘柄コード
	IssType string  `json:"IssType"` // 銘柄区分（1: 信用銘柄、2: 貸借銘柄、3: その他）

	// 信用取引残高（売建）
	ShrtVol    float64 `json:"ShrtVol"`    // 売合計信用取引残高
	ShrtNegVol float64 `json:"ShrtNegVol"` // 売一般信用取引残高
	ShrtStdVol float64 `json:"ShrtStdVol"` // 売制度信用取引残高

	// 信用取引残高（買建）
	LongVol    float64 `json:"LongVol"`    // 買合計信用取引残高
	LongNegVol float64 `json:"LongNegVol"` // 買一般信用取引残高
	LongStdVol float64 `json:"LongStdVol"` // 買制度信用取引残高

	// 信用取引残高（金額）
	ShrtVal    *float64 `json:"ShrtVal"`    // 売合計
	LongVal    *float64 `json:"LongVal"`    // 買合計
	ShrtNegVal *float64 `json:"ShrtNegVal"` // 売一般信用
	LongNegVal *float64 `json:"LongNegVal"` // 買一般信用
	ShrtStdVal *float64 `json:"ShrtStdVal"` // 売制度信用
	LongStdVal *float64 `json:"LongStdVal"` // 買制度信用
}

// RawWeeklyMarginInterest is used for unmarshaling JSON response with mixed types
type RawWeeklyMarginInterest struct {
	// 基本情報
	PubDate types.NullableString `json:"PubDate"`
	Date    string               `json:"Date"`
	Code    string               `json:"Code"`
	IssType string               `json:"IssType"`

	// 信用取引残高（売建）
	ShrtVol    types.NullableFloat64 `json:"ShrtVol"`
	ShrtNegVol types.NullableFloat64 `json:"ShrtNegVol"`
	ShrtStdVol types.NullableFloat64 `json:"ShrtStdVol"`

	// 信用取引残高（買建）
	LongVol    types.NullableFloat64 `json:"LongVol"`
	LongNegVol types.NullableFloat64 `json:"LongNegVol"`
	LongStdVol types.NullableFloat64 `json:"LongStdVol"`

	ShrtVal    types.NullableFloat64 `json:"ShrtVal"`
	LongVal    types.NullableFloat64 `json:"LongVal"`
	ShrtNegVal types.NullableFloat64 `json:"ShrtNegVal"`
	LongNegVal types.NullableFloat64 `json:"LongNegVal"`
	ShrtStdVal types.NullableFloat64 `json:"ShrtStdVal"`
	LongStdVal types.NullableFloat64 `json:"LongStdVal"`
}

// UnmarshalJSON implements custom JSON unmarshaling for WeeklyMarginInterestResponse
func (r *WeeklyMarginInterestResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawWeeklyMarginInterest
	type rawResponse struct {
		Data          []RawWeeklyMarginInterest `json:"data"`
		PaginationKey string                    `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	r.PaginationKey = raw.PaginationKey

	// Convert RawWeeklyMarginInterest to WeeklyMarginInterest
	r.Data = make([]WeeklyMarginInterest, len(raw.Data))
	for idx, rm := range raw.Data {
		r.Data[idx] = WeeklyMarginInterest{
			// 基本情報
			PubDate: rm.PubDate.Ptr(),
			Date:    rm.Date,
			Code:    rm.Code,
			IssType: rm.IssType,

			// 信用取引残高（売建）
			ShrtVol:    rm.ShrtVol.Or(0),
			ShrtNegVol: rm.ShrtNegVol.Or(0),
			ShrtStdVol: rm.ShrtStdVol.Or(0),

			// 信用取引残高（買建）
			LongVol:    rm.LongVol.Or(0),
			LongNegVol: rm.LongNegVol.Or(0),
			LongStdVol: rm.LongStdVol.Or(0),
			ShrtVal:    rm.ShrtVal.Ptr(),
			LongVal:    rm.LongVal.Ptr(),
			ShrtNegVal: rm.ShrtNegVal.Ptr(),
			LongNegVal: rm.LongNegVal.Ptr(),
			ShrtStdVal: rm.ShrtStdVal.Ptr(),
			LongStdVal: rm.LongStdVal.Ptr(),
		}
	}

	return nil
}

// GetWeeklyMarginInterest は信用取引残高を取得します。
func (s *WeeklyMarginInterestService) GetWeeklyMarginInterest(ctx context.Context, params WeeklyMarginInterestParams) (*WeeklyMarginInterestResponse, error) {
	if params.Code == "" && params.Date == "" && params.PublishedDate == "" {
		return nil, fmt.Errorf("either code, date or published_date parameter is required")
	}
	if params.PublishedDate != "" && (params.Date != "" || params.From != "" || params.To != "") {
		return nil, fmt.Errorf("published_date cannot be combined with date, from or to")
	}

	path := "/markets/margin-interest"

	query := "?"
	if params.Code != "" {
		query += fmt.Sprintf("code=%s&", params.Code)
	}
	if params.Date != "" {
		query += fmt.Sprintf("date=%s&", params.Date)
	}
	if params.PublishedDate != "" {
		query += fmt.Sprintf("published_date=%s&", params.PublishedDate)
	}
	if params.From != "" {
		query += fmt.Sprintf("from=%s&", params.From)
	}
	if params.To != "" {
		query += fmt.Sprintf("to=%s&", params.To)
	}
	if params.PaginationKey != "" {
		query += fmt.Sprintf("pagination_key=%s&", params.PaginationKey)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp WeeklyMarginInterestResponse
	if err := s.client.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get weekly margin interest: %w", err)
	}

	return &resp, nil
}

// GetWeeklyMarginInterestByCode は指定銘柄の信用取引残高を取得します。
// ページネーションを使用して全データを取得します。
func (s *WeeklyMarginInterestService) GetWeeklyMarginInterestByCode(ctx context.Context, code string) ([]WeeklyMarginInterest, error) {
	var allData []WeeklyMarginInterest
	paginationKey := ""

	for {
		params := WeeklyMarginInterestParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetWeeklyMarginInterest(ctx, params)
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

// GetWeeklyMarginInterestByDate は指定日の全銘柄信用取引残高を取得します。
// ページネーションを使用して全データを取得します。
func (s *WeeklyMarginInterestService) GetWeeklyMarginInterestByDate(ctx context.Context, date string) ([]WeeklyMarginInterest, error) {
	var allData []WeeklyMarginInterest
	paginationKey := ""

	for {
		params := WeeklyMarginInterestParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetWeeklyMarginInterest(ctx, params)
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

// GetWeeklyMarginInterestByPublishedDate は指定公表日の全銘柄の残高を取得します。
// 公表日が収録されていない過去データは取得できません。
func (s *WeeklyMarginInterestService) GetWeeklyMarginInterestByPublishedDate(ctx context.Context, publishedDate string) ([]WeeklyMarginInterest, error) {
	params := WeeklyMarginInterestParams{PublishedDate: publishedDate}
	var allData []WeeklyMarginInterest
	for {
		resp, err := s.GetWeeklyMarginInterest(ctx, params)
		if err != nil {
			return nil, err
		}
		allData = append(allData, resp.Data...)
		if resp.PaginationKey == "" {
			return allData, nil
		}
		params.PaginationKey = resp.PaginationKey
	}
}

// GetWeeklyMarginInterestByCodeAndDateRange は指定銘柄・期間の信用取引残高を取得します。
func (s *WeeklyMarginInterestService) GetWeeklyMarginInterestByCodeAndDateRange(ctx context.Context, code, from, to string) ([]WeeklyMarginInterest, error) {
	var allData []WeeklyMarginInterest
	paginationKey := ""

	for {
		params := WeeklyMarginInterestParams{
			Code:          code,
			From:          from,
			To:            to,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetWeeklyMarginInterest(ctx, params)
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

// GetWeeklyMarginInterestByCodeAndDate は指定銘柄の指定申込日の信用取引残高を取得します。
func (s *WeeklyMarginInterestService) GetWeeklyMarginInterestByCodeAndDate(ctx context.Context, code, date string) ([]WeeklyMarginInterest, error) {
	resp, err := s.GetWeeklyMarginInterest(ctx, WeeklyMarginInterestParams{
		Code: code,
		Date: date,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// IsCredit は信用銘柄かどうかを判定します。
func (wmi *WeeklyMarginInterest) IsCredit() bool {
	return wmi.IssType == IssueTypeCredit
}

// IsLendable は貸借銘柄かどうかを判定します。
func (wmi *WeeklyMarginInterest) IsLendable() bool {
	return wmi.IssType == IssueTypeLendable
}

// GetShortLongRatio は売建残高と買建残高の比率を計算します（売建/買建）。
func (wmi *WeeklyMarginInterest) GetShortLongRatio() float64 {
	if wmi.LongVol == 0 {
		return 0
	}
	return wmi.ShrtVol / wmi.LongVol
}

// GetStandardizedRatio は制度信用の割合を計算します（制度信用/合計）。
func (wmi *WeeklyMarginInterest) GetStandardizedRatio() (float64, float64) {
	// 売建の制度信用比率
	shortStandardizedRatio := float64(0)
	if wmi.ShrtVol > 0 {
		shortStandardizedRatio = wmi.ShrtStdVol / wmi.ShrtVol
	}

	// 買建の制度信用比率
	longStandardizedRatio := float64(0)
	if wmi.LongVol > 0 {
		longStandardizedRatio = wmi.LongStdVol / wmi.LongVol
	}

	return shortStandardizedRatio, longStandardizedRatio
}
