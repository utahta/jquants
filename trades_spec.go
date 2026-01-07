package jquants

import (
	"fmt"

	"github.com/utahta/jquants/client"
)

// TradesSpecService は投資部門別情報を取得するサービスです。
// 投資部門別売買状況（株式・金額）のデータを提供します。
type TradesSpecService struct {
	client client.HTTPClient
}

// NewTradesSpecService は新しいTradesSpecServiceを作成します。
func NewTradesSpecService(c client.HTTPClient) *TradesSpecService {
	return &TradesSpecService{client: c}
}

// TradesSpec は投資部門別情報データを表します。
// J-Quants API /equities/investor-types エンドポイントのレスポンスデータ。
type TradesSpec struct {
	// 基本情報
	PubDate string `json:"PubDate"` // 公表日（YYYY-MM-DD形式）
	StDate  string `json:"StDate"`  // 開始日（YYYY-MM-DD形式）
	EnDate  string `json:"EnDate"`  // 終了日（YYYY-MM-DD形式）
	Section string `json:"Section"` // 市場名

	// 自己取引
	PropSell float64 `json:"PropSell"` // 自己計_売（千円）
	PropBuy  float64 `json:"PropBuy"`  // 自己計_買（千円）
	PropTot  float64 `json:"PropTot"`  // 自己計_合計（千円）
	PropBal  float64 `json:"PropBal"`  // 自己計_差引（千円）

	// 委託取引
	BrkSell float64 `json:"BrkSell"` // 委託計_売（千円）
	BrkBuy  float64 `json:"BrkBuy"`  // 委託計_買（千円）
	BrkTot  float64 `json:"BrkTot"`  // 委託計_合計（千円）
	BrkBal  float64 `json:"BrkBal"`  // 委託計_差引（千円）

	// 総計
	TotSell float64 `json:"TotSell"` // 総計_売（千円）
	TotBuy  float64 `json:"TotBuy"`  // 総計_買（千円）
	TotTot  float64 `json:"TotTot"`  // 総計_合計（千円）
	TotBal  float64 `json:"TotBal"`  // 総計_差引（千円）

	// 投資部門別内訳
	IndSell    float64 `json:"IndSell"`    // 個人_売（千円）
	IndBuy     float64 `json:"IndBuy"`     // 個人_買（千円）
	IndTot     float64 `json:"IndTot"`     // 個人_合計（千円）
	IndBal     float64 `json:"IndBal"`     // 個人_差引（千円）
	FrgnSell   float64 `json:"FrgnSell"`   // 海外投資家_売（千円）
	FrgnBuy    float64 `json:"FrgnBuy"`    // 海外投資家_買（千円）
	FrgnTot    float64 `json:"FrgnTot"`    // 海外投資家_合計（千円）
	FrgnBal    float64 `json:"FrgnBal"`    // 海外投資家_差引（千円）
	SecCoSell  float64 `json:"SecCoSell"`  // 証券会社_売（千円）
	SecCoBuy   float64 `json:"SecCoBuy"`   // 証券会社_買（千円）
	SecCoTot   float64 `json:"SecCoTot"`   // 証券会社_合計（千円）
	SecCoBal   float64 `json:"SecCoBal"`   // 証券会社_差引（千円）
	InvTrSell  float64 `json:"InvTrSell"`  // 投資信託_売（千円）
	InvTrBuy   float64 `json:"InvTrBuy"`   // 投資信託_買（千円）
	InvTrTot   float64 `json:"InvTrTot"`   // 投資信託_合計（千円）
	InvTrBal   float64 `json:"InvTrBal"`   // 投資信託_差引（千円）
	BusCoSell  float64 `json:"BusCoSell"`  // 事業法人_売（千円）
	BusCoBuy   float64 `json:"BusCoBuy"`   // 事業法人_買（千円）
	BusCoTot   float64 `json:"BusCoTot"`   // 事業法人_合計（千円）
	BusCoBal   float64 `json:"BusCoBal"`   // 事業法人_差引（千円）
	OthCoSell  float64 `json:"OthCoSell"`  // その他法人_売（千円）
	OthCoBuy   float64 `json:"OthCoBuy"`   // その他法人_買（千円）
	OthCoTot   float64 `json:"OthCoTot"`   // その他法人_合計（千円）
	OthCoBal   float64 `json:"OthCoBal"`   // その他法人_差引（千円）
	InsCoSell  float64 `json:"InsCoSell"`  // 生保・損保_売（千円）
	InsCoBuy   float64 `json:"InsCoBuy"`   // 生保・損保_買（千円）
	InsCoTot   float64 `json:"InsCoTot"`   // 生保・損保_合計（千円）
	InsCoBal   float64 `json:"InsCoBal"`   // 生保・損保_差引（千円）
	BankSell   float64 `json:"BankSell"`   // 都銀・地銀等_売（千円）
	BankBuy    float64 `json:"BankBuy"`    // 都銀・地銀等_買（千円）
	BankTot    float64 `json:"BankTot"`    // 都銀・地銀等_合計（千円）
	BankBal    float64 `json:"BankBal"`    // 都銀・地銀等_差引（千円）
	TrstBnkSel float64 `json:"TrstBnkSell"` // 信託銀行_売（千円）
	TrstBnkBuy float64 `json:"TrstBnkBuy"`  // 信託銀行_買（千円）
	TrstBnkTot float64 `json:"TrstBnkTot"`  // 信託銀行_合計（千円）
	TrstBnkBal float64 `json:"TrstBnkBal"`  // 信託銀行_差引（千円）
	OthFinSell float64 `json:"OthFinSell"`  // その他金融機関_売（千円）
	OthFinBuy  float64 `json:"OthFinBuy"`   // その他金融機関_買（千円）
	OthFinTot  float64 `json:"OthFinTot"`   // その他金融機関_合計（千円）
	OthFinBal  float64 `json:"OthFinBal"`   // その他金融機関_差引（千円）
}

// TradesSpecResponse は投資部門別情報のレスポンスです。
type TradesSpecResponse struct {
	Data          []TradesSpec `json:"data"`
	PaginationKey string       `json:"pagination_key"` // ページネーションキー
}

// TradesSpecParams は投資部門別情報のリクエストパラメータです。
type TradesSpecParams struct {
	Section       string // セクション（市場）（例: TSEPrime, TSEStandard, TSEGrowth）
	From          string // 開始日（YYYYMMDD または YYYY-MM-DD）
	To            string // 終了日（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// 市場コード定数
const (
	SectionTSE1st      = "TSE1st"      // 市場一部
	SectionTSE2nd      = "TSE2nd"      // 市場二部
	SectionTSEMothers  = "TSEMothers"  // マザーズ
	SectionTSEJASDAQ   = "TSEJASDAQ"   // JASDAQ
	SectionTSEPrime    = "TSEPrime"    // プライム
	SectionTSEStandard = "TSEStandard" // スタンダード
	SectionTSEGrowth   = "TSEGrowth"   // グロース
	SectionTokyoNagoya = "TokyoNagoya" // 東証および名証
)

// GetTradesSpec は投資部門別情報を取得します。
// パラメータの組み合わせ:
// - section指定あり、from/to指定あり: 指定したセクションの指定した期間のデータ
// - section指定あり、from/to指定なし: 指定したセクションの全期間のデータ
// - section指定なし、from/to指定あり: すべてのセクションの指定した期間のデータ
// - section指定なし、from/to指定なし: すべてのセクションの全期間のデータ
func (s *TradesSpecService) GetTradesSpec(params TradesSpecParams) (*TradesSpecResponse, error) {
	path := "/equities/investor-types"

	query := "?"
	if params.Section != "" {
		query += fmt.Sprintf("section=%s&", params.Section)
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

	var resp TradesSpecResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get trades spec: %w", err)
	}

	return &resp, nil
}

// GetTradesSpecByDateRange は指定期間の投資部門別情報を取得します。
// ページネーションを使用して全データを取得します。
func (s *TradesSpecService) GetTradesSpecByDateRange(from, to string) ([]TradesSpec, error) {
	var allTradesSpec []TradesSpec
	paginationKey := ""

	for {
		params := TradesSpecParams{
			From:          from,
			To:            to,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetTradesSpec(params)
		if err != nil {
			return nil, err
		}

		allTradesSpec = append(allTradesSpec, resp.Data...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allTradesSpec, nil
}

// GetTradesSpecBySection は指定セクションの投資部門別情報を取得します。
// ページネーションを使用して全データを取得します。
func (s *TradesSpecService) GetTradesSpecBySection(section string) ([]TradesSpec, error) {
	var allTradesSpec []TradesSpec
	paginationKey := ""

	for {
		params := TradesSpecParams{
			Section:       section,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetTradesSpec(params)
		if err != nil {
			return nil, err
		}

		allTradesSpec = append(allTradesSpec, resp.Data...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allTradesSpec, nil
}

// GetAllTradesSpec は全セクション・全期間の投資部門別情報を取得します。
// ページネーションを使用して大量データを分割取得します。
func (s *TradesSpecService) GetAllTradesSpec() ([]TradesSpec, error) {
	var allTradesSpec []TradesSpec
	paginationKey := ""

	for {
		params := TradesSpecParams{
			PaginationKey: paginationKey,
		}

		resp, err := s.GetTradesSpec(params)
		if err != nil {
			return nil, err
		}

		allTradesSpec = append(allTradesSpec, resp.Data...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allTradesSpec, nil
}

// IsBuyerDominant は買い手が優勢かどうかを判定します（差引がプラス）。
func (ts *TradesSpec) IsBuyerDominant(investorType string) bool {
	switch investorType {
	case "individuals":
		return ts.IndBal > 0
	case "foreigners":
		return ts.FrgnBal > 0
	case "securities":
		return ts.SecCoBal > 0
	case "investment_trusts":
		return ts.InvTrBal > 0
	case "business":
		return ts.BusCoBal > 0
	case "insurance":
		return ts.InsCoBal > 0
	case "trust_banks":
		return ts.TrstBnkBal > 0
	case "total":
		return ts.TotBal > 0
	default:
		return false
	}
}

// GetNetFlow は指定した投資家タイプの純流入額を取得します（差引額）。
func (ts *TradesSpec) GetNetFlow(investorType string) float64 {
	switch investorType {
	case "individuals":
		return ts.IndBal
	case "foreigners":
		return ts.FrgnBal
	case "securities":
		return ts.SecCoBal
	case "investment_trusts":
		return ts.InvTrBal
	case "business":
		return ts.BusCoBal
	case "insurance":
		return ts.InsCoBal
	case "trust_banks":
		return ts.TrstBnkBal
	case "total":
		return ts.TotBal
	default:
		return 0
	}
}
