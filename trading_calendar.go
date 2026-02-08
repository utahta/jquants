package jquants

import (
	"fmt"

	"github.com/utahta/jquants/client"
)

// TradingCalendarService は取引カレンダーを取得するサービスです。
// 東証およびOSEにおける営業日、休業日、祝日取引の有無の情報を提供します。
type TradingCalendarService struct {
	client client.HTTPClient
}

// NewTradingCalendarService は新しいTradingCalendarServiceを作成します。
func NewTradingCalendarService(c client.HTTPClient) *TradingCalendarService {
	return &TradingCalendarService{client: c}
}

// TradingCalendarParams は取引カレンダーのリクエストパラメータです。
type TradingCalendarParams struct {
	HolidayDivision string // 休日区分（0: 非営業日, 1: 営業日, 2: 東証半日立会日, 3: 非営業日(祝日取引あり)）
	From            string // 開始日（YYYYMMDD または YYYY-MM-DD）
	To              string // 終了日（YYYYMMDD または YYYY-MM-DD）
}

// TradingCalendarResponse は取引カレンダーのレスポンスです。
type TradingCalendarResponse struct {
	Data []TradingCalendar `json:"data"`
}

// TradingCalendar は取引カレンダーのデータです。
// J-Quants API /markets/calendar エンドポイントのレスポンスデータ。
type TradingCalendar struct {
	Date   string `json:"Date"`   // 日付（YYYY-MM-DD形式）
	HolDiv string `json:"HolDiv"` // 休日区分（0: 非営業日, 1: 営業日, 2: 東証半日立会日, 3: 非営業日(祝日取引あり)）
}

// 休日区分定数
const (
	HolidayDivisionNonTradingDay     = "0" // 非営業日（東証・OSEともに休業）
	HolidayDivisionTradingDay        = "1" // 営業日（通常の営業日）
	HolidayDivisionTSEHalfDay        = "2" // 東証半日立会日（東証が半日のみ取引）
	HolidayDivisionOSEHolidayTrading = "3" // 非営業日(祝日取引あり)（東証は休業だがOSEで祝日取引を実施）
)

// GetTradingCalendar は取引カレンダーを取得します。
// パラメータの組み合わせ:
// - hol_div指定あり、from/to指定なし: 指定された休日区分について全期間分のデータ
// - hol_div指定あり、from/to指定あり: 指定された休日区分について指定された期間分のデータ
// - hol_div指定なし、from/to指定あり: 指定された期間分のデータ
// - hol_div指定なし、from/to指定なし: 全期間分のデータ
func (s *TradingCalendarService) GetTradingCalendar(params TradingCalendarParams) (*TradingCalendarResponse, error) {
	path := "/markets/calendar"

	query := "?"
	if params.HolidayDivision != "" {
		query += fmt.Sprintf("hol_div=%s&", params.HolidayDivision)
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

	var resp TradingCalendarResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get trading calendar: %w", err)
	}

	return &resp, nil
}

// GetAllTradingCalendar は全期間・全区分の取引カレンダーを取得します。
func (s *TradingCalendarService) GetAllTradingCalendar() ([]TradingCalendar, error) {
	resp, err := s.GetTradingCalendar(TradingCalendarParams{})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetTradingCalendarByHolidayDivision は指定した休日区分の全期間の取引カレンダーを取得します。
func (s *TradingCalendarService) GetTradingCalendarByHolidayDivision(holDiv string) ([]TradingCalendar, error) {
	resp, err := s.GetTradingCalendar(TradingCalendarParams{
		HolidayDivision: holDiv,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetTradingCalendarByHolidayDivisionAndDateRange は指定した休日区分・期間の取引カレンダーを取得します。
func (s *TradingCalendarService) GetTradingCalendarByHolidayDivisionAndDateRange(holDiv, from, to string) ([]TradingCalendar, error) {
	resp, err := s.GetTradingCalendar(TradingCalendarParams{
		HolidayDivision: holDiv,
		From:            from,
		To:              to,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetTradingCalendarByDateRange は指定した期間の取引カレンダーを取得します。
func (s *TradingCalendarService) GetTradingCalendarByDateRange(from, to string) ([]TradingCalendar, error) {
	resp, err := s.GetTradingCalendar(TradingCalendarParams{
		From: from,
		To:   to,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetTradingDays は指定した期間の営業日のみを取得します。
func (s *TradingCalendarService) GetTradingDays(from, to string) ([]TradingCalendar, error) {
	resp, err := s.GetTradingCalendar(TradingCalendarParams{
		HolidayDivision: HolidayDivisionTradingDay,
		From:            from,
		To:              to,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetNonTradingDays は指定した期間の非営業日のみを取得します。
func (s *TradingCalendarService) GetNonTradingDays(from, to string) ([]TradingCalendar, error) {
	resp, err := s.GetTradingCalendar(TradingCalendarParams{
		HolidayDivision: HolidayDivisionNonTradingDay,
		From:            from,
		To:              to,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetOSEHolidayTradingDays はOSEで祝日取引を実施する日を取得します。
func (s *TradingCalendarService) GetOSEHolidayTradingDays(from, to string) ([]TradingCalendar, error) {
	resp, err := s.GetTradingCalendar(TradingCalendarParams{
		HolidayDivision: HolidayDivisionOSEHolidayTrading,
		From:            from,
		To:              to,
	})
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// IsTradingDay は指定した日付が営業日かどうかを判定します。
func (tc *TradingCalendar) IsTradingDay() bool {
	return tc.HolDiv == HolidayDivisionTradingDay
}

// IsNonTradingDay は指定した日付が非営業日かどうかを判定します。
func (tc *TradingCalendar) IsNonTradingDay() bool {
	return tc.HolDiv == HolidayDivisionNonTradingDay
}

// IsHalfTradingDay は指定した日付が半日立会日かどうかを判定します。
func (tc *TradingCalendar) IsHalfTradingDay() bool {
	return tc.HolDiv == HolidayDivisionTSEHalfDay
}

// HasOSEHolidayTrading は指定した日付にOSEで祝日取引があるかどうかを判定します。
func (tc *TradingCalendar) HasOSEHolidayTrading() bool {
	return tc.HolDiv == HolidayDivisionOSEHolidayTrading
}
