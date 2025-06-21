package jquants

import "github.com/utahta/jquants/client"

// JQuantsAPI はJ-Quants APIの各サービスを統合したインターフェースです。
// 各サービスの説明:
// - Quotes: 株価データ（日次の始値、高値、安値、終値、出来高）
// - PricesAM: 前場四本値データ（午前中の取引データ）
// - Listed: 企業情報（企業名、業種、市場区分）
// - Statements: 財務諸表（売上高、利益、ROE/ROA）
// - Dividend: 配当情報（配当金、権利確定日、支払日）
// - Announcement: 決算発表予定（発表日、発表時刻）
// - TradesSpec: 投資部門別売買状況（機関投資家、個人投資家等の売買動向）
// - WeeklyMarginInterest: 信用取引週末残高（信用買い/売り残高）
// - ShortSelling: 業種別空売り比率（業種ごとの空売り状況）
// - ShortSellingPositions: 空売り残高報告（大口の空売りポジション）
// - Breakdown: 売買内訳データ（売買の詳細な内訳）
// - TradingCalendar: 取引カレンダー（営業日/休日情報）
// - Indices: 指数四本値（日経平均、TOPIX等）
// - TOPIX: TOPIX専用詳細データ
// - IndexOption: 日経225オプション
// - Futures: 先物取引データ
// - Options: 個別株オプション
// - FSDetails: 財務諸表詳細（BS/PL詳細データ）
type JQuantsAPI struct {
	client                client.HTTPClient
	Quotes                *QuotesService
	PricesAM              *PricesAMService
	Listed                *ListedService
	Statements            *StatementsService
	Dividend              *DividendService
	Announcement          *AnnouncementService
	TradesSpec            *TradesSpecService
	WeeklyMarginInterest  *WeeklyMarginInterestService
	ShortSelling          *ShortSellingService
	ShortSellingPositions *ShortSellingPositionsService
	Breakdown             *BreakdownService
	TradingCalendar       *TradingCalendarService
	Indices               *IndicesService
	TOPIX                 *TOPIXService
	IndexOption           *IndexOptionService
	Futures               *FuturesService
	Options               *OptionsService
	FSDetails             *FSDetailsService
}

func NewJQuantsAPI(c client.HTTPClient) *JQuantsAPI {
	return &JQuantsAPI{
		client:                c,
		Quotes:                NewQuotesService(c),
		PricesAM:              NewPricesAMService(c),
		Listed:                NewListedService(c),
		Statements:            NewStatementsService(c),
		Dividend:              NewDividendService(c),
		Announcement:          NewAnnouncementService(c),
		TradesSpec:            NewTradesSpecService(c),
		WeeklyMarginInterest:  NewWeeklyMarginInterestService(c),
		ShortSelling:          NewShortSellingService(c),
		ShortSellingPositions: NewShortSellingPositionsService(c),
		Breakdown:             NewBreakdownService(c),
		TradingCalendar:       NewTradingCalendarService(c),
		Indices:               NewIndicesService(c),
		TOPIX:                 NewTOPIXService(c),
		IndexOption:           NewIndexOptionService(c),
		Futures:               NewFuturesService(c),
		Options:               NewOptionsService(c),
		FSDetails:             NewFSDetailsService(c),
	}
}
