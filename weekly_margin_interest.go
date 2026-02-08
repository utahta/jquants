package jquants

import (
	"encoding/json"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// WeeklyMarginInterestService は信用取引週末残高を取得するサービスです。
type WeeklyMarginInterestService struct {
	client client.HTTPClient
}

// NewWeeklyMarginInterestService は新しいWeeklyMarginInterestServiceを作成します。
func NewWeeklyMarginInterestService(c client.HTTPClient) *WeeklyMarginInterestService {
	return &WeeklyMarginInterestService{client: c}
}

// WeeklyMarginInterestParams は信用取引週末残高のリクエストパラメータです。
type WeeklyMarginInterestParams struct {
	Code          string // 銘柄コード（codeまたはdateのいずれかが必須）
	Date          string // 申込日付（YYYYMMDD または YYYY-MM-DD）（codeまたはdateのいずれかが必須）
	From          string // 期間の開始日（YYYYMMDD または YYYY-MM-DD）
	To            string // 期間の終了日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// WeeklyMarginInterestResponse は信用取引週末残高のレスポンスです。
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

// WeeklyMarginInterest は信用取引週末残高のデータを表します。
// J-Quants API /markets/margin-interest エンドポイントのレスポンスデータ。
type WeeklyMarginInterest struct {
	// 基本情報
	Date    string `json:"Date"`    // 申込日付（YYYY-MM-DD形式）
	Code    string `json:"Code"`    // 銘柄コード
	IssType string `json:"IssType"` // 銘柄区分（1: 信用銘柄、2: 貸借銘柄、3: その他）

	// 信用取引残高（売建）
	ShrtVol    float64 `json:"ShrtVol"`    // 売合計信用取引週末残高
	ShrtNegVol float64 `json:"ShrtNegVol"` // 売一般信用取引週末残高
	ShrtStdVol float64 `json:"ShrtStdVol"` // 売制度信用取引週末残高

	// 信用取引残高（買建）
	LongVol    float64 `json:"LongVol"`    // 買合計信用取引週末残高
	LongNegVol float64 `json:"LongNegVol"` // 買一般信用取引週末残高
	LongStdVol float64 `json:"LongStdVol"` // 買制度信用取引週末残高
}

// RawWeeklyMarginInterest is used for unmarshaling JSON response with mixed types
type RawWeeklyMarginInterest struct {
	// 基本情報
	Date    string `json:"Date"`
	Code    string `json:"Code"`
	IssType string `json:"IssType"`

	// 信用取引残高（売建）
	ShrtVol    types.Float64String `json:"ShrtVol"`
	ShrtNegVol types.Float64String `json:"ShrtNegVol"`
	ShrtStdVol types.Float64String `json:"ShrtStdVol"`

	// 信用取引残高（買建）
	LongVol    types.Float64String `json:"LongVol"`
	LongNegVol types.Float64String `json:"LongNegVol"`
	LongStdVol types.Float64String `json:"LongStdVol"`
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
			Date:    rm.Date,
			Code:    rm.Code,
			IssType: rm.IssType,

			// 信用取引残高（売建）
			ShrtVol:    float64(rm.ShrtVol),
			ShrtNegVol: float64(rm.ShrtNegVol),
			ShrtStdVol: float64(rm.ShrtStdVol),

			// 信用取引残高（買建）
			LongVol:    float64(rm.LongVol),
			LongNegVol: float64(rm.LongNegVol),
			LongStdVol: float64(rm.LongStdVol),
		}
	}

	return nil
}

// GetWeeklyMarginInterest は信用取引週末残高を取得します。
func (s *WeeklyMarginInterestService) GetWeeklyMarginInterest(params WeeklyMarginInterestParams) (*WeeklyMarginInterestResponse, error) {
	// codeまたはdateのいずれかが必須
	if params.Code == "" && params.Date == "" {
		return nil, fmt.Errorf("either code or date parameter is required")
	}

	path := "/markets/margin-interest"

	query := "?"
	if params.Code != "" {
		query += fmt.Sprintf("code=%s&", params.Code)
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
	if params.PaginationKey != "" {
		query += fmt.Sprintf("pagination_key=%s&", params.PaginationKey)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp WeeklyMarginInterestResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get weekly margin interest: %w", err)
	}

	return &resp, nil
}

// GetWeeklyMarginInterestByCode は指定銘柄の信用取引週末残高を取得します。
// ページネーションを使用して全データを取得します。
func (s *WeeklyMarginInterestService) GetWeeklyMarginInterestByCode(code string) ([]WeeklyMarginInterest, error) {
	var allData []WeeklyMarginInterest
	paginationKey := ""

	for {
		params := WeeklyMarginInterestParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetWeeklyMarginInterest(params)
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

// GetWeeklyMarginInterestByDate は指定日の全銘柄信用取引週末残高を取得します。
// ページネーションを使用して全データを取得します。
func (s *WeeklyMarginInterestService) GetWeeklyMarginInterestByDate(date string) ([]WeeklyMarginInterest, error) {
	var allData []WeeklyMarginInterest
	paginationKey := ""

	for {
		params := WeeklyMarginInterestParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetWeeklyMarginInterest(params)
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

// GetWeeklyMarginInterestByCodeAndDateRange は指定銘柄・期間の信用取引週末残高を取得します。
func (s *WeeklyMarginInterestService) GetWeeklyMarginInterestByCodeAndDateRange(code, from, to string) ([]WeeklyMarginInterest, error) {
	var allData []WeeklyMarginInterest
	paginationKey := ""

	for {
		params := WeeklyMarginInterestParams{
			Code:          code,
			From:          from,
			To:            to,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetWeeklyMarginInterest(params)
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

// GetWeeklyMarginInterestByCodeAndDate は指定銘柄の指定公表日の信用取引週末残高を取得します。
func (s *WeeklyMarginInterestService) GetWeeklyMarginInterestByCodeAndDate(code, date string) ([]WeeklyMarginInterest, error) {
	resp, err := s.GetWeeklyMarginInterest(WeeklyMarginInterestParams{
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
