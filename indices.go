package jquants

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/utahta/jquants/client"
	"github.com/utahta/jquants/types"
)

// IndicesService は指数四本値データを取得するサービスです。
// TOPIX、日経平均等の主要指数の四本値データを提供します。
type IndicesService struct {
	client client.HTTPClient
}

// NewIndicesService は新しいIndicesServiceを作成します。
func NewIndicesService(c client.HTTPClient) *IndicesService {
	return &IndicesService{client: c}
}

// Index は指数四本値データを表します。
// J-Quants API /indices エンドポイントのレスポンスデータ。
type Index struct {
	Date  string  `json:"Date"`  // 日付（YYYY-MM-DD形式）
	Code  string  `json:"Code"`  // 指数コード
	Open  float64 `json:"Open"`  // 始値
	High  float64 `json:"High"`  // 高値
	Low   float64 `json:"Low"`   // 安値
	Close float64 `json:"Close"` // 終値
}

// RawIndex is used for unmarshaling JSON response with mixed types
type RawIndex struct {
	Date  string              `json:"Date"`
	Code  string              `json:"Code"`
	Open  types.Float64String `json:"Open"`
	High  types.Float64String `json:"High"`
	Low   types.Float64String `json:"Low"`
	Close types.Float64String `json:"Close"`
}

// IndicesResponse は指数四本値のレスポンスです。
type IndicesResponse struct {
	Indices       []Index `json:"indices"`
	PaginationKey string  `json:"pagination_key"` // ページネーションキー
}

// UnmarshalJSON implements custom JSON unmarshaling
func (i *IndicesResponse) UnmarshalJSON(data []byte) error {
	// First unmarshal into RawIndex
	type rawResponse struct {
		Indices       []RawIndex `json:"indices"`
		PaginationKey string     `json:"pagination_key"`
	}

	var raw rawResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Set pagination key
	i.PaginationKey = raw.PaginationKey

	// Convert RawIndex to Index
	i.Indices = make([]Index, len(raw.Indices))
	for idx, ri := range raw.Indices {
		i.Indices[idx] = Index{
			Date:  ri.Date,
			Code:  ri.Code,
			Open:  float64(ri.Open),
			High:  float64(ri.High),
			Low:   float64(ri.Low),
			Close: float64(ri.Close),
		}
	}

	return nil
}

// IndicesParams は指数四本値のリクエストパラメータです。
type IndicesParams struct {
	Code          string // 指数コード（例: "0000" TOPIX、"0028" TOPIX Core30）
	Date          string // 基準日付（YYYYMMDD または YYYY-MM-DD）
	From          string // 開始日付（YYYYMMDD または YYYY-MM-DD）
	To            string // 終了日付（YYYYMMDD または YYYY-MM-DD）
	PaginationKey string // ページネーションキー
}

// GetIndices は指定された条件で指数四本値データを取得します。
// codeまたはdateのいずれかが必須です。
// パラメータ:
// - Code: 指数コード（例: "0000" または "0028"）
// - Date: 基準日付（例: "20240101" または "2024-01-01"）
// - From/To: 期間指定（例: "20240101" または "2024-01-01"）
// - PaginationKey: ページネーション用キー
func (s *IndicesService) GetIndices(params IndicesParams) (*IndicesResponse, error) {
	path := "/indices"

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

	var resp IndicesResponse
	if err := s.client.DoRequest("GET", path, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get indices: %w", err)
	}

	return &resp, nil
}

// GetIndicesByCode は指定指数の過去N日間のデータを取得します。
// ページネーションを使用して全データを取得します。
func (s *IndicesService) GetIndicesByCode(code string, days int) ([]Index, error) {
	to := time.Now()
	from := to.AddDate(0, 0, -days)

	var allIndices []Index
	paginationKey := ""

	for {
		params := IndicesParams{
			Code:          code,
			From:          from.Format("20060102"),
			To:            to.Format("20060102"),
			PaginationKey: paginationKey,
		}

		resp, err := s.GetIndices(params)
		if err != nil {
			return nil, err
		}

		allIndices = append(allIndices, resp.Indices...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allIndices, nil
}

// GetIndicesByDate は指定日の全指数データを取得します。
// ページネーションを使用して大量データを分割取得します。
func (s *IndicesService) GetIndicesByDate(date string) ([]Index, error) {
	var allIndices []Index
	paginationKey := ""

	for {
		params := IndicesParams{
			Date:          date,
			PaginationKey: paginationKey,
		}

		resp, err := s.GetIndices(params)
		if err != nil {
			return nil, err
		}

		allIndices = append(allIndices, resp.Indices...)

		// ページネーションキーがなければ終了
		if resp.PaginationKey == "" {
			break
		}
		paginationKey = resp.PaginationKey
	}

	return allIndices, nil
}

// 主要指数コード定数
const (
	// TOPIX関連
	IndexTOPIX         = "0000" // TOPIX
	IndexTOPIXCore30   = "0028" // TOPIX Core30
	IndexTOPIXLarge70  = "0029" // TOPIX Large 70
	IndexTOPIX100      = "002A" // TOPIX 100
	IndexTOPIXMid400   = "002B" // TOPIX Mid400
	IndexTOPIX500      = "002C" // TOPIX 500
	IndexTOPIXSmall    = "002D" // TOPIX Small
	IndexTOPIX1000     = "002E" // TOPIX 1000
	IndexTOPIXSmall500 = "002F" // TOPIX Small500

	// その他市場指数
	IndexMothers250 = "0070" // 東証グロース市場250指数（旧：東証マザーズ指数）
	IndexREIT       = "0075" // REIT

	// 新市場区分指数
	IndexPrime       = "0500" // 東証プライム市場指数
	IndexStandard    = "0501" // 東証スタンダード市場指数
	IndexGrowth      = "0502" // 東証グロース市場指数
	IndexJPXPrime150 = "0503" // JPXプライム150指数
)

// GetTOPIX はTOPIXのデータを取得します。
func (s *IndicesService) GetTOPIX(days int) ([]Index, error) {
	return s.GetIndicesByCode(IndexTOPIX, days)
}

// GetTOPIXCore30 はTOPIX Core30のデータを取得します。
func (s *IndicesService) GetTOPIXCore30(days int) ([]Index, error) {
	return s.GetIndicesByCode(IndexTOPIXCore30, days)
}

// GetPrimeMarketIndex は東証プライム市場指数のデータを取得します。
func (s *IndicesService) GetPrimeMarketIndex(days int) ([]Index, error) {
	return s.GetIndicesByCode(IndexPrime, days)
}

// GetREIT はREIT指数のデータを取得します。
func (s *IndicesService) GetREIT(days int) ([]Index, error) {
	return s.GetIndicesByCode(IndexREIT, days)
}

// 業種別指数コード（東証33業種）
const (
	IndexSectorFishery      = "0040" // 水産・農林業
	IndexSectorMining       = "0041" // 鉱業
	IndexSectorConstruction = "0042" // 建設業
	IndexSectorFoods        = "0043" // 食料品
	IndexSectorTextiles     = "0044" // 繊維製品
	IndexSectorPulpPaper    = "0045" // パルプ・紙
	IndexSectorChemicals    = "0046" // 化学
	IndexSectorPharmaceut   = "0047" // 医薬品
	IndexSectorOilCoal      = "0048" // 石油・石炭製品
	IndexSectorRubber       = "0049" // ゴム製品
	IndexSectorGlassCeram   = "004A" // ガラス・土石製品
	IndexSectorIronSteel    = "004B" // 鉄鋼
	IndexSectorNonferrous   = "004C" // 非鉄金属
	IndexSectorMetalProd    = "004D" // 金属製品
	IndexSectorMachinery    = "004E" // 機械
	IndexSectorElecAppl     = "004F" // 電気機器
	IndexSectorTransEquip   = "0050" // 輸送用機器
	IndexSectorPrecision    = "0051" // 精密機器
	IndexSectorOtherProd    = "0052" // その他製品
	IndexSectorElecGas      = "0053" // 電気・ガス業
	IndexSectorLandTrans    = "0054" // 陸運業
	IndexSectorMarine       = "0055" // 海運業
	IndexSectorAir          = "0056" // 空運業
	IndexSectorWarehouse    = "0057" // 倉庫・運輸関連業
	IndexSectorInfoComm     = "0058" // 情報・通信業
	IndexSectorWholesale    = "0059" // 卸売業
	IndexSectorRetail       = "005A" // 小売業
	IndexSectorBanks        = "005B" // 銀行業
	IndexSectorSecurities   = "005C" // 証券・商品先物取引業
	IndexSectorInsurance    = "005D" // 保険業
	IndexSectorOtherFin     = "005E" // その他金融業
	IndexSectorRealEstate   = "005F" // 不動産業
	IndexSectorServices     = "0060" // サービス業
)

// GetSectorIndex は指定した業種別指数のデータを取得します。
func (s *IndicesService) GetSectorIndex(sectorCode string, days int) ([]Index, error) {
	return s.GetIndicesByCode(sectorCode, days)
}
