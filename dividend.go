package jquants

import (
	"context"
	"fmt"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// DividendService は配当情報を取得するサービスです。
// 過去の配当実績や今後の配当予定を提供します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では利用できません。
type DividendService struct {
	client client.HTTPClient
}

// NewDividendService は新しいDividendServiceを作成します。
func NewDividendService(c client.HTTPClient) *DividendService {
	return &DividendService{client: c}
}

// DividendParams は配当情報のリクエストパラメータです。
type DividendParams struct {
	Code          string // 銘柄コード（codeまたはdateのいずれかが必須）
	Date          string // 通知日付（YYYYMMDD または YYYY-MM-DD）（codeまたはdateのいずれかが必須）
	From          string // 期間の開始日（YYYYMMDD または YYYY-MM-DD）
	To            string // 期間の終了日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// DividendResponse は配当情報のレスポンスです。
type DividendResponse struct {
	Data          []Dividend `json:"data"`
	PaginationKey string     `json:"pagination_key"` // ページネーションキー
}

// Dividend は配当情報を表します。
// J-Quants API /fins/dividend エンドポイントのレスポンスデータ。
// 上場会社の配当（決定・予想）に関する１株当たり配当金額、基準日、権利落日及び支払開始予定日等の情報。
//
// 注意: このデータはプレミアムプラン専用APIで取得されます。
type Dividend struct {
	// 基本情報
	PubDate string `json:"PubDate"` // 通知日時（年月日）（YYYY-MM-DD形式）
	PubTime string `json:"PubTime"` // 通知日時（時分）（HH:MI形式）
	Code    string `json:"Code"`    // 銘柄コード
	RefNo   string `json:"RefNo"`   // リファレンスナンバー（配当通知を一意に特定するための番号）
	CARefNo string `json:"CARefNo"` // CAリファレンスナンバー（訂正・削除の対象となっている配当通知のリファレンスナンバー）

	// 更新・配当区分情報
	StatCode     string `json:"StatCode"`     // 更新区分（コード）（1: 新規、2: 訂正、3: 削除）
	IFCode       string `json:"IFCode"`       // 配当種類（コード）（1: 中間配当、2: 期末配当）
	FRCode       string `json:"FRCode"`       // 予想／決定（コード）（1: 決定、2: 予想）
	CommSpecCode string `json:"CommSpecCode"` // 記念配当/特別配当コード（0: 通常、1: 記念配当、2: 特別配当、3: 記念・特別配当）

	// 日程情報
	// PayDate等のNullable型フィールドは、未定（-）をIsUndetermined()で判別できます。
	BoardDate  string               `json:"BoardDate"`  // 取締役会決議日（YYYY-MM-DD形式）
	IFTerm     string               `json:"IFTerm"`     // 配当基準日年月（YYYY-MM形式）
	RecDate    string               `json:"RecDate"`    // 基準日（YYYY-MM-DD形式）
	ExDate     string               `json:"ExDate"`     // 権利落日（YYYY-MM-DD形式）
	ActRecDate string               `json:"ActRecDate"` // 権利確定日（YYYY-MM-DD形式）
	PayDate    types.NullableString `json:"PayDate"`    // 支払開始予定日（YYYY-MM-DD形式）

	// 配当金額情報
	DivRate     types.NullableFloat64 `json:"DivRate"`     // １株当たり配当金額
	CommDivRate types.NullableFloat64 `json:"CommDivRate"` // １株当たり記念配当金額（2022年6月6日以降のみ）
	SpecDivRate types.NullableFloat64 `json:"SpecDivRate"` // １株当たり特別配当金額（2022年6月6日以降のみ）

	// 税務関連情報（2014年2月24日以降のみ提供）
	DistAmt          types.NullableFloat64 `json:"DistAmt"`          // 1株当たりの交付金銭等の額
	RetEarn          types.NullableFloat64 `json:"RetEarn"`          // 1株当たりの利益剰余金の額
	DeemDiv          types.NullableFloat64 `json:"DeemDiv"`          // 1株当たりのみなし配当の額
	DeemCapGains     types.NullableFloat64 `json:"DeemCapGains"`     // 1株当たりのみなし譲渡収入の額
	NetAssetDecRatio types.NullableFloat64 `json:"NetAssetDecRatio"` // 純資産減少割合
}

// GetDividend は配当情報を取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *DividendService) GetDividend(ctx context.Context, params DividendParams) (*DividendResponse, error) {
	// codeまたはdateのいずれかが必須
	if params.Code == "" && params.Date == "" {
		return nil, fmt.Errorf("either code or date parameter is required")
	}

	path := "/fins/dividend"

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

	var resp DividendResponse
	if err := s.client.DoRequest(ctx, "GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get dividend: %w", err)
	}

	return &resp, nil
}

// GetDividendByCode は指定銘柄の配当情報を取得します。
// ページネーションを使用して全データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *DividendService) GetDividendByCode(ctx context.Context, code string) ([]Dividend, error) {
	var allData []Dividend
	paginationKey := ""

	for {
		params := DividendParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDividend(ctx, params)
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

// GetDividendByDate は指定日の全銘柄配当情報を取得します。
// ページネーションを使用して全データを取得します。
//
// 注意: このAPIはプレミアムプラン専用です。
// スタンダードプラン以下では "This API is not available on your subscription" エラーが返されます。
func (s *DividendService) GetDividendByDate(ctx context.Context, date string) ([]Dividend, error) {
	var allData []Dividend
	paginationKey := ""

	for {
		params := DividendParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDividend(ctx, params)
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

// GetDividendByCodeAndDateRange は指定銘柄・期間の配当情報を取得します。
func (s *DividendService) GetDividendByCodeAndDateRange(ctx context.Context, code, from, to string) ([]Dividend, error) {
	var allData []Dividend
	paginationKey := ""

	for {
		params := DividendParams{
			Code:          code,
			From:          from,
			To:            to,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDividend(ctx, params)
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

// IsNew は新規通知かを判定します。
func (d *Dividend) IsNew() bool {
	return d.StatCode == "1"
}

// IsRevision は訂正通知かを判定します。
func (d *Dividend) IsRevision() bool {
	return d.StatCode == "2"
}

// IsDeleted は削除通知かを判定します。
func (d *Dividend) IsDeleted() bool {
	return d.StatCode == "3"
}

// IsInterim は中間配当かを判定します。
func (d *Dividend) IsInterim() bool {
	return d.IFCode == "1"
}

// IsFinal は期末配当かを判定します。
func (d *Dividend) IsFinal() bool {
	return d.IFCode == "2"
}

// IsForecast は予想配当かを判定します。
func (d *Dividend) IsForecast() bool {
	return d.FRCode == "2"
}

// IsResult は決定配当かを判定します。
func (d *Dividend) IsResult() bool {
	return d.FRCode == "1"
}

// IsCommemorative は記念配当を含むかを判定します。
func (d *Dividend) IsCommemorative() bool {
	return d.CommSpecCode == "1" || d.CommSpecCode == "3"
}

// IsSpecial は特別配当を含むかを判定します。
func (d *Dividend) IsSpecial() bool {
	return d.CommSpecCode == "2" || d.CommSpecCode == "3"
}

// IsOrdinary は通常配当のみかを判定します。
func (d *Dividend) IsOrdinary() bool {
	return d.CommSpecCode == "0"
}

// GetTotalDividendRate は配当金額の合計を計算します（通常＋記念＋特別）。
func (d *Dividend) GetTotalDividendRate() *float64 {
	rate, ok := d.DivRate.Get()
	if !ok {
		return nil
	}

	total := rate + d.CommDivRate.Or(0) + d.SpecDivRate.Or(0)
	return &total
}

// GetOrdinaryDividendRate は通常配当金額を計算します（総額から記念・特別を除く）。
func (d *Dividend) GetOrdinaryDividendRate() *float64 {
	rate, ok := d.DivRate.Get()
	if !ok {
		return nil
	}

	ordinary := rate - d.CommDivRate.Or(0) - d.SpecDivRate.Or(0)
	return &ordinary
}

// HasPayableDate は支払開始予定日が設定されているかを判定します。
func (d *Dividend) HasPayableDate() bool {
	_, ok := d.PayDate.Get()
	return ok
}

// IsPayableDateUndecided は支払開始予定日が未定（-）かを判定します。
func (d *Dividend) IsPayableDateUndecided() bool {
	return d.PayDate.IsUndetermined()
}

// IsDividendRateUndecided は配当金額が未定（-）かを判定します。
// 非設定（空文字）の場合はfalseを返します。
func (d *Dividend) IsDividendRateUndecided() bool {
	return d.DivRate.IsUndetermined()
}
