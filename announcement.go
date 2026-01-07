package jquants

import (
	"fmt"

	"github.com/utahta/jquants/client"
)

// AnnouncementService は決算発表予定データを取得するサービスです。
// 翌営業日に決算発表が行われる銘柄の情報を提供します。
// 現在は3月期・9月期決算会社のみ対応しています。
type AnnouncementService struct {
	client client.HTTPClient
}

func NewAnnouncementService(c client.HTTPClient) *AnnouncementService {
	return &AnnouncementService{
		client: c,
	}
}

// Announcement は決算発表予定を表します。
// J-Quants API /equities/earnings-calendar エンドポイントのレスポンスデータ。
// 注意: このAPIは翌営業日の決算発表予定のみを返します。
type Announcement struct {
	Date     string `json:"Date"`     // 日付（YYYY-MM-DD形式、決算発表予定日が未定の場合は空文字）
	Code     string `json:"Code"`     // 銘柄コード
	CoName   string `json:"CoName"`   // 会社名
	FY       string `json:"FY"`       // 決算期末（例: "9月30日"）
	SectorNm string `json:"SectorNm"` // 業種名
	FQ       string `json:"FQ"`       // 決算種別（例: "第１四半期"、"第２四半期"、"第３四半期"、"通期"）
	Section  string `json:"Section"`  // 市場区分（例: "マザーズ"）
}

// AnnouncementResponse は決算発表予定のレスポンスです。
type AnnouncementResponse struct {
	Data          []Announcement `json:"data"`
	PaginationKey string         `json:"pagination_key"` // ページネーションキー
}

// AnnouncementParams は決算発表予定のリクエストパラメータです。
type AnnouncementParams struct {
	PaginationKey string // ページネーションキー
}

// GetAnnouncement は翌営業日の決算発表予定を取得します。
// このAPIは翌営業日に決算発表が行われる銘柄の情報のみを返します。
// 3月期・9月期決算会社のみ対応しています。
// パラメータ:
// - params: ページネーション用のパラメータ（オプション）
func (s *AnnouncementService) GetAnnouncement(params AnnouncementParams) (*AnnouncementResponse, error) {
	path := "/equities/earnings-calendar"

	if params.PaginationKey != "" {
		path += fmt.Sprintf("?pagination_key=%s", params.PaginationKey)
	}

	var resp AnnouncementResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get announcement: %w", err)
	}

	return &resp, nil
}

// GetAllAnnouncements は翌営業日の全決算発表予定を取得します。
// ページネーションを使用して全データを取得します。
func (s *AnnouncementService) GetAllAnnouncements() ([]Announcement, error) {
	var allAnnouncements []Announcement
	paginationKey := ""

	for {
		params := AnnouncementParams{
			PaginationKey: paginationKey,
		}

		resp, err := s.GetAnnouncement(params)
		if err != nil {
			return nil, err
		}

		allAnnouncements = append(allAnnouncements, resp.Data...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allAnnouncements, nil
}

// GetAnnouncementByCode は指定銘柄の決算発表予定を取得します。
// 翌営業日に決算発表予定がある場合のみ取得できます。
func (s *AnnouncementService) GetAnnouncementByCode(code string) (*Announcement, error) {
	announcements, err := s.GetAllAnnouncements()
	if err != nil {
		return nil, err
	}

	for _, announcement := range announcements {
		if announcement.Code == code {
			return &announcement, nil
		}
	}

	return nil, fmt.Errorf("no announcement found for code: %s", code)
}
