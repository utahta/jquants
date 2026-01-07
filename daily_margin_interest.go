package jquants

import (
	"fmt"

	"github.com/utahta/jquants/client"
)

// DailyMarginInterestService は日々公表信用取引残高を取得するサービスです。
// 日々公表銘柄に指定された個別銘柄の日々の信用取引残高を取得できます。
type DailyMarginInterestService struct {
	client client.HTTPClient
}

// NewDailyMarginInterestService は新しいDailyMarginInterestServiceを作成します。
func NewDailyMarginInterestService(c client.HTTPClient) *DailyMarginInterestService {
	return &DailyMarginInterestService{client: c}
}

// DailyMarginInterestParams は日々公表信用取引残高のリクエストパラメータです。
type DailyMarginInterestParams struct {
	Code          string // 銘柄コード（codeまたはdateのいずれかが必須）
	Date          string // 公表日（YYYYMMDD または YYYY-MM-DD）（codeまたはdateのいずれかが必須）
	From          string // 期間の開始日（YYYYMMDD または YYYY-MM-DD）
	To            string // 期間の終了日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// DailyMarginInterestResponse は日々公表信用取引残高のレスポンスです。
type DailyMarginInterestResponse struct {
	Data          []DailyMarginInterest `json:"data"`
	PaginationKey string                `json:"pagination_key"`
}

// PublishReason は公表の理由を表します。
type PublishReason struct {
	Restricted         string `json:"Restricted"`         // 規制措置（0: 対象外, 1: 対象）
	DailyPublication   string `json:"DailyPublication"`   // 日々公表（0: 対象外, 1: 対象）
	Monitoring         string `json:"Monitoring"`         // 監視（0: 対象外, 1: 対象）
	RestrictedByJSF    string `json:"RestrictedByJSF"`    // 日証金の規制（0: 対象外, 1: 対象）
	PrecautionByJSF    string `json:"PrecautionByJSF"`    // 日証金の注意喚起（0: 対象外, 1: 対象）
	UnclearOrSecOnAlert string `json:"UnclearOrSecOnAlert"` // 不明瞭または証券会社の注意喚起（0: 対象外, 1: 対象）
}

// DailyMarginInterest は日々公表信用取引残高のデータを表します。
// J-Quants API /markets/margin-alert エンドポイントのレスポンスデータ。
type DailyMarginInterest struct {
	// 基本情報
	PubDate   string        `json:"PubDate"`   // 公表日（YYYY-MM-DD形式）
	Code      string        `json:"Code"`      // 銘柄コード
	AppDate   string        `json:"AppDate"`   // 申込日（YYYY-MM-DD形式）
	PubReason PublishReason `json:"PubReason"` // 公表の理由

	// 売合計信用残高
	ShrtOut      float64     `json:"ShrtOut"`      // 売合計信用残高（株）
	ShrtOutChg   interface{} `json:"ShrtOutChg"`   // 前日比 売合計信用残高（株）。前日に公表されていない場合は「-」
	ShrtOutRatio interface{} `json:"ShrtOutRatio"` // 上場比 売合計信用残高（%）。ETFの場合は「*」

	// 買合計信用残高
	LongOut      float64     `json:"LongOut"`      // 買合計信用残高（株）
	LongOutChg   interface{} `json:"LongOutChg"`   // 前日比 買合計信用残高（株）。前日に公表されていない場合は「-」
	LongOutRatio interface{} `json:"LongOutRatio"` // 上場比 買合計信用残高（%）。ETFの場合は「*」

	// 取組比率
	SLRatio float64 `json:"SLRatio"` // 取組比率（%）= 売合計信用残高 / 買合計信用残高 × 100

	// 一般信用取引残高
	ShrtNegOut    float64     `json:"ShrtNegOut"`    // 一般信用取引売残高（株）
	ShrtNegOutChg interface{} `json:"ShrtNegOutChg"` // 前日比 一般信用取引売残高（株）
	LongNegOut    float64     `json:"LongNegOut"`    // 一般信用取引買残高（株）
	LongNegOutChg interface{} `json:"LongNegOutChg"` // 前日比 一般信用取引買残高（株）

	// 制度信用取引残高
	ShrtStdOut    float64     `json:"ShrtStdOut"`    // 制度信用取引売残高（株）
	ShrtStdOutChg interface{} `json:"ShrtStdOutChg"` // 前日比 制度信用取引売残高（株）
	LongStdOut    float64     `json:"LongStdOut"`    // 制度信用取引買残高（株）
	LongStdOutChg interface{} `json:"LongStdOutChg"` // 前日比 制度信用取引買残高（株）

	// 規制区分
	TSEMrgnRegCls string `json:"TSEMrgnRegCls"` // 東証信用貸借規制区分
}

// 東証信用貸借規制区分定数
const (
	TSEMarginRegulationNone              = "000" // 規制なし
	TSEMarginRegulationCautionForNew     = "001" // 新規建て注意喚起
	TSEMarginRegulationCautionForSelling = "002" // 売り注意喚起
	TSEMarginRegulationCautionForBuying  = "003" // 買い注意喚起
	TSEMarginRegulationRestrictedNew     = "011" // 新規建て規制
	TSEMarginRegulationRestrictedSelling = "012" // 売り規制
	TSEMarginRegulationRestrictedBuying  = "013" // 買い規制
)

// GetDailyMarginInterest は日々公表信用取引残高を取得します。
func (s *DailyMarginInterestService) GetDailyMarginInterest(params DailyMarginInterestParams) (*DailyMarginInterestResponse, error) {
	// codeまたはdateのいずれかが必須
	if params.Code == "" && params.Date == "" {
		return nil, fmt.Errorf("either code or date parameter is required")
	}

	path := "/markets/margin-alert"

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

	var resp DailyMarginInterestResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get daily margin interest: %w", err)
	}

	return &resp, nil
}

// GetDailyMarginInterestByCode は指定銘柄の日々公表信用取引残高を取得します。
// ページネーションを使用して全データを取得します。
func (s *DailyMarginInterestService) GetDailyMarginInterestByCode(code string) ([]DailyMarginInterest, error) {
	var allData []DailyMarginInterest
	paginationKey := ""

	for {
		params := DailyMarginInterestParams{
			Code:          code,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDailyMarginInterest(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.Data...)

		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetDailyMarginInterestByDate は指定日の全銘柄の日々公表信用取引残高を取得します。
// ページネーションを使用して全データを取得します。
func (s *DailyMarginInterestService) GetDailyMarginInterestByDate(date string) ([]DailyMarginInterest, error) {
	var allData []DailyMarginInterest
	paginationKey := ""

	for {
		params := DailyMarginInterestParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDailyMarginInterest(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.Data...)

		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetDailyMarginInterestByCodeAndDateRange は指定銘柄・期間の日々公表信用取引残高を取得します。
func (s *DailyMarginInterestService) GetDailyMarginInterestByCodeAndDateRange(code, from, to string) ([]DailyMarginInterest, error) {
	var allData []DailyMarginInterest
	paginationKey := ""

	for {
		params := DailyMarginInterestParams{
			Code:          code,
			From:          from,
			To:            to,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetDailyMarginInterest(params)
		if err != nil {
			return nil, err
		}

		allData = append(allData, resp.Data...)

		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allData, nil
}

// GetShortOutChgValue は前日比売合計信用残高を数値で取得します。
// 前日に公表されていない場合（「-」）は0を返します。
func (d *DailyMarginInterest) GetShortOutChgValue() float64 {
	if v, ok := d.ShrtOutChg.(float64); ok {
		return v
	}
	return 0
}

// GetLongOutChgValue は前日比買合計信用残高を数値で取得します。
// 前日に公表されていない場合（「-」）は0を返します。
func (d *DailyMarginInterest) GetLongOutChgValue() float64 {
	if v, ok := d.LongOutChg.(float64); ok {
		return v
	}
	return 0
}

// GetShortOutRatioValue は上場比売合計信用残高を数値で取得します。
// ETFの場合（「*」）は0を返します。
func (d *DailyMarginInterest) GetShortOutRatioValue() float64 {
	if v, ok := d.ShrtOutRatio.(float64); ok {
		return v
	}
	return 0
}

// GetLongOutRatioValue は上場比買合計信用残高を数値で取得します。
// ETFの場合（「*」）は0を返します。
func (d *DailyMarginInterest) GetLongOutRatioValue() float64 {
	if v, ok := d.LongOutRatio.(float64); ok {
		return v
	}
	return 0
}

// IsRestricted は規制措置対象かどうかを判定します。
func (p *PublishReason) IsRestricted() bool {
	return p.Restricted == "1"
}

// IsDailyPublication は日々公表対象かどうかを判定します。
func (p *PublishReason) IsDailyPublication() bool {
	return p.DailyPublication == "1"
}

// IsMonitoring は監視対象かどうかを判定します。
func (p *PublishReason) IsMonitoring() bool {
	return p.Monitoring == "1"
}

// IsRestrictedByJSF は日証金の規制対象かどうかを判定します。
func (p *PublishReason) IsRestrictedByJSF() bool {
	return p.RestrictedByJSF == "1"
}

// IsPrecautionByJSF は日証金の注意喚起対象かどうかを判定します。
func (p *PublishReason) IsPrecautionByJSF() bool {
	return p.PrecautionByJSF == "1"
}

// IsUnclearOrSecOnAlert は不明瞭または証券会社の注意喚起対象かどうかを判定します。
func (p *PublishReason) IsUnclearOrSecOnAlert() bool {
	return p.UnclearOrSecOnAlert == "1"
}
