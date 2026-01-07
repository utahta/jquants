package jquants

import (
	"fmt"

	"github.com/utahta/jquants/client"
)

// 市場区分コード定義
const (
	MarketTSE1st         = "0101" // 東証一部
	MarketTSE2nd         = "0102" // 東証二部
	MarketMothers        = "0104" // マザーズ
	MarketTokyoProMarket = "0105" // TOKYO PRO MARKET
	MarketJASDAQStandard = "0106" // JASDAQ スタンダード
	MarketJASDAQGrowth   = "0107" // JASDAQ グロース
	MarketOther          = "0109" // その他
	MarketPrime          = "0111" // プライム
	MarketStandard       = "0112" // スタンダード
	MarketGrowth         = "0113" // グロース
)

// 17業種コード定義
const (
	Sector17Food         = "1"  // 食品
	Sector17Energy       = "2"  // エネルギー資源
	Sector17Construction = "3"  // 建設・資材
	Sector17Materials    = "4"  // 素材・化学
	Sector17Pharma       = "5"  // 医薬品
	Sector17Auto         = "6"  // 自動車・輸送機
	Sector17Steel        = "7"  // 鉄鋼・非鉄
	Sector17Machinery    = "8"  // 機械
	Sector17Electric     = "9"  // 電機・精密
	Sector17IT           = "10" // 情報通信・サービスその他
	Sector17Utilities    = "11" // 電力・ガス
	Sector17Transport    = "12" // 運輸・物流
	Sector17Trading      = "13" // 商社・卸売
	Sector17Retail       = "14" // 小売
	Sector17Banks        = "15" // 銀行
	Sector17Finance      = "16" // 金融（除く銀行）
	Sector17RealEstate   = "17" // 不動産
	Sector17Other        = "99" // その他
)

// 33業種コード定義（主要なもの）
const (
	Sector33Fishery       = "0050" // 水産・農林業
	Sector33Mining        = "1050" // 鉱業
	Sector33Construction  = "2050" // 建設業
	Sector33Foods         = "3050" // 食料品
	Sector33Textiles      = "3100" // 繊維製品
	Sector33Paper         = "3150" // パルプ・紙
	Sector33Chemicals     = "3200" // 化学
	Sector33Pharma        = "3250" // 医薬品
	Sector33Oil           = "3300" // 石油・石炭製品
	Sector33Rubber        = "3350" // ゴム製品
	Sector33Glass         = "3400" // ガラス・土石製品
	Sector33Steel         = "3450" // 鉄鋼
	Sector33NonFerrous    = "3500" // 非鉄金属
	Sector33Metal         = "3550" // 金属製品
	Sector33Machinery     = "3600" // 機械
	Sector33Electric      = "3650" // 電気機器
	Sector33Transport     = "3700" // 輸送用機器
	Sector33Precision     = "3750" // 精密機器
	Sector33OtherProducts = "3800" // その他製品
	Sector33Utilities     = "4050" // 電気・ガス業
	Sector33Land          = "5050" // 陸運業
	Sector33Marine        = "5100" // 海運業
	Sector33Air           = "5150" // 空運業
	Sector33Warehouse     = "5200" // 倉庫・運輸関連業
	Sector33IT            = "5250" // 情報・通信業
	Sector33Wholesale     = "6050" // 卸売業
	Sector33Retail        = "6100" // 小売業
	Sector33Banks         = "7050" // 銀行業
	Sector33Securities    = "7100" // 証券・商品先物取引業
	Sector33Insurance     = "7150" // 保険業
	Sector33OtherFinance  = "7200" // その他金融業
	Sector33RealEstate    = "8050" // 不動産業
	Sector33Services      = "9050" // サービス業
	Sector33Other         = "9999" // その他
)

// ListedService は上場企業情報を取得するサービスです。
// 企業名、業種分類、市場区分などの基本情報を提供します。
type ListedService struct {
	client client.HTTPClient
}

func NewListedService(c client.HTTPClient) *ListedService {
	return &ListedService{
		client: c,
	}
}

// ListedInfo は上場企業の基本情報を表します。
// J-Quants API /equities/master エンドポイントのレスポンスデータ。
// 過去時点、当日、翌営業日時点の銘柄情報が取得可能（翌営業日は17:30以降）。
type ListedInfo struct {
	Date     string `json:"Date"`     // 情報適用年月日（YYYY-MM-DD形式）
	Code     string `json:"Code"`     // 銘柄コード（4桁または5桁）
	CoName   string `json:"CoName"`   // 企業名（日本語）
	CoNameEn string `json:"CoNameEn"` // 企業名（英語）
	S17      string `json:"S17"`      // 17業種コード
	S17Nm    string `json:"S17Nm"`    // 17業種コード名
	S33      string `json:"S33"`      // 33業種コード
	S33Nm    string `json:"S33Nm"`    // 33業種コード名
	ScaleCat string `json:"ScaleCat"` // 規模コード（TOPIX Core30、Large70等）
	Mkt      string `json:"Mkt"`      // 市場区分コード
	MktNm    string `json:"MktNm"`    // 市場区分名（プライム、スタンダード、グロース等）
	Mrgn     string `json:"Mrgn"`     // 貸借信用区分（1: 信用 / 2: 貸借 / 3: その他）（Standard/Premiumのみ）
	MrgnNm   string `json:"MrgnNm"`   // 貸借信用区分名（Standard/Premiumのみ）
}

type ListedInfoResponse struct {
	Data []ListedInfo `json:"data"`
}

// GetListedInfo は上場企業情報を取得します。
// パラメータ:
// - code: 銘柄コード（空の場合は全銘柄）4桁または5桁形式をサポート
// - date: 基準日（YYYYMMDD または YYYY-MM-DD 形式、空の場合は最新）
//
// 使用例:
// - 特定銘柄: GetListedInfo("7203", "")
// - 全銘柄（最新）: GetListedInfo("", "")
// - 過去時点: GetListedInfo("", "20210907")
func (s *ListedService) GetListedInfo(code string, date string) ([]ListedInfo, error) {
	path := "/equities/master"

	query := "?"
	if code != "" {
		query += fmt.Sprintf("code=%s&", code)
	}
	if date != "" {
		query += fmt.Sprintf("date=%s&", date)
	}

	if len(query) > 1 {
		path += query[:len(query)-1] // Remove trailing &
	}

	var resp ListedInfoResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get listed info: %w", err)
	}

	return resp.Data, nil
}

// GetCompanyInfo は指定銘柄の最新企業情報を取得します。
// 例: GetCompanyInfo("7203") でトヨタ自動車の企業情報を取得
func (s *ListedService) GetCompanyInfo(code string) (*ListedInfo, error) {
	infos, err := s.GetListedInfo(code, "")
	if err != nil {
		return nil, err
	}

	if len(infos) == 0 {
		return nil, fmt.Errorf("no company info found for code: %s", code)
	}

	return &infos[0], nil
}

// GetListedBySector17 は指定した17業種コードの銘柄一覧を取得します。
// 例: GetListedBySector17(Sector17IT, "") でIT関連銘柄を取得
func (s *ListedService) GetListedBySector17(sector17Code string, date string) ([]ListedInfo, error) {
	allInfo, err := s.GetListedInfo("", date)
	if err != nil {
		return nil, err
	}

	var filtered []ListedInfo
	for _, info := range allInfo {
		if info.S17 == sector17Code {
			filtered = append(filtered, info)
		}
	}

	return filtered, nil
}

// GetListedBySector33 は指定した33業種コードの銘柄一覧を取得します。
// 例: GetListedBySector33(Sector33IT, "") で情報・通信業銘柄を取得
func (s *ListedService) GetListedBySector33(sector33Code string, date string) ([]ListedInfo, error) {
	allInfo, err := s.GetListedInfo("", date)
	if err != nil {
		return nil, err
	}

	var filtered []ListedInfo
	for _, info := range allInfo {
		if info.S33 == sector33Code {
			filtered = append(filtered, info)
		}
	}

	return filtered, nil
}

// GetListedByMarket は指定した市場区分の銘柄一覧を取得します。
// marketCode: MarketPrime, MarketStandard, MarketGrowth など
func (s *ListedService) GetListedByMarket(marketCode string, date string) ([]ListedInfo, error) {
	allInfo, err := s.GetListedInfo("", date)
	if err != nil {
		return nil, err
	}

	var filtered []ListedInfo
	for _, info := range allInfo {
		if info.Mkt == marketCode {
			filtered = append(filtered, info)
		}
	}

	return filtered, nil
}
