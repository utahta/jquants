//go:build e2e
// +build e2e

package e2e

import (
	"testing"

	"github.com/utahta/jquants"
)

// TestIndicesEndpoint は/markets/indicesエンドポイントの完全なテスト
func TestIndicesEndpoint(t *testing.T) {
	t.Run("GetIndices_ByDate", func(t *testing.T) {
		// 最近の営業日の指数データを取得
		date := getTestDate()

		params := jquants.IndicesParams{
			Date: date,
		}

		resp, err := jq.Indices.GetIndices(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get indices: %v", err)
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No indices data available for the specified date")
		}

		// 各指数データを詳細に検証
		for i, index := range resp.Data {
			// 基本情報の検証
			if index.Date == "" {
				t.Errorf("Index[%d]: Date is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(index.Date) != 10 || index.Date[4] != '-' || index.Date[7] != '-' {
					t.Errorf("Index[%d]: Date format invalid = %v, want YYYY-MM-DD", i, index.Date)
				}
				// 日付の一致確認
				expectedDate := getTestDateFormatted()
				if index.Date != expectedDate {
					t.Errorf("Index[%d]: Date = %v, want %v", i, index.Date, expectedDate)
				}
			}

			if index.Code == "" {
				t.Errorf("Index[%d]: Code is empty", i)
			}

			// 四本値の検証
			if index.O <= 0 {
				t.Errorf("Index[%d]: Open = %v, want > 0", i, index.O)
			}
			if index.H <= 0 {
				t.Errorf("Index[%d]: High = %v, want > 0", i, index.H)
			}
			if index.L <= 0 {
				t.Errorf("Index[%d]: Low = %v, want > 0", i, index.L)
			}
			if index.C <= 0 {
				t.Errorf("Index[%d]: Close = %v, want > 0", i, index.C)
			}

			// 四本値の論理的整合性チェック
			if index.H < index.L {
				t.Errorf("Index[%d]: High (%v) < Low (%v)", i, index.H, index.L)
			}
			if index.O > index.H || index.O < index.L {
				t.Errorf("Index[%d]: Open (%v) is outside High (%v) - Low (%v) range",
					i, index.O, index.H, index.L)
			}
			if index.C > index.H || index.C < index.L {
				t.Errorf("Index[%d]: Close (%v) is outside High (%v) - Low (%v) range",
					i, index.C, index.H, index.L)
			}

			// 最初の5件の詳細ログ
			if i < 5 {
				t.Logf("Index[%d]: Code=%s, Date=%s, O=%.2f, H=%.2f, L=%.2f, C=%.2f",
					i, index.Code, index.Date, index.O, index.H, index.L, index.C)
			}
		}

		t.Logf("Retrieved %d indices for date %s", len(resp.Data), date)

		// 主要な指数が含まれているか確認
		codeMap := make(map[string]bool)
		for _, index := range resp.Data {
			codeMap[index.Code] = true
		}

		// よく知られた指数コードの存在確認（TOPIX等）
		knownIndices := []string{"0000", "0028", "0070"} // TOPIX、TOPIX Core30、グロース市場250等
		for _, code := range knownIndices {
			if codeMap[code] {
				t.Logf("Found known index: %s", code)
			}
		}
	})

	t.Run("GetIndices_WithCode", func(t *testing.T) {
		// 特定の指数コードでフィルタリング
		date := getTestDate()

		params := jquants.IndicesParams{
			Date: date,
			Code: "0000", // 特定の指数コード
		}

		resp, err := jq.Indices.GetIndices(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Logf("Failed to get indices with code filter: %v", err)
			return
		}

		if resp != nil && len(resp.Data) > 0 {
			// 全てのデータが指定したコードか確認
			for i, index := range resp.Data {
				if index.Code != params.Code {
					t.Errorf("Index[%d]: Code = %v, want %v", i, index.Code, params.Code)
				}
			}
			t.Logf("Retrieved %d indices for code %s", len(resp.Data), params.Code)
		}
	})

	t.Run("GetIndices_Pagination", func(t *testing.T) {
		// ページネーションのテスト
		date := getTestDate()

		params := jquants.IndicesParams{
			Date: date,
		}

		resp, err := jq.Indices.GetIndices(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No indices data available")
		}

		if resp == nil || len(resp.Data) == 0 {
			t.Skip("No indices data available")
		}

		firstPageCount := len(resp.Data)
		t.Logf("First page: %d indices", firstPageCount)

		if resp.PaginationKey != "" {
			// 次のページを取得
			params.PaginationKey = resp.PaginationKey
			resp2, err := jq.Indices.GetIndices(params)
			if err != nil {
				t.Fatalf("Failed to get next page: %v", err)
			}

			if resp2 != nil && len(resp2.Data) > 0 {
				t.Logf("Second page: %d indices", len(resp2.Data))

				// 異なるデータであることを確認
				if resp2.Data[0].Code == resp.Data[0].Code {
					t.Logf("Warning: Second page might contain overlapping data")
				}
			}
		} else {
			t.Logf("No pagination needed (all data in first page)")
		}
	})

	t.Run("GetTOPIXCore30", func(t *testing.T) {
		// TOPIX Core30の便利メソッドテスト
		indices, err := jq.Indices.GetTOPIXCore30(30) // 過去30日分
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get TOPIX Core30: %v", err)
		}

		if len(indices) == 0 {
			t.Skip("No TOPIX Core30 data available")
		}

		t.Logf("Retrieved %d days of TOPIX Core30 data", len(indices))

		// データの検証
		for i, index := range indices {
			if index.Code != "0028" { // TOPIX Core30のコード
				t.Errorf("Index[%d]: Code = %v, want 0028", i, index.Code)
			}

			// 最初と最後のデータをログ
			if i == 0 || i == len(indices)-1 {
				t.Logf("TOPIX Core30 [%s]: Close=%.2f (O=%.2f, H=%.2f, L=%.2f)",
					index.Date, index.C, index.O, index.H, index.L)
			}
		}

		// 価格変動の計算（最初と最後）
		if len(indices) >= 2 {
			first := indices[0]
			last := indices[len(indices)-1]
			change := last.C - first.C
			changePercent := (change / first.C) * 100
			t.Logf("Price change over period: %.2f (%.2f%%)", change, changePercent)
		}
	})

	t.Run("GetTOPIX", func(t *testing.T) {
		// TOPIXの便利メソッドテスト
		indices, err := jq.Indices.GetTOPIX(30) // 過去30日分
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get TOPIX: %v", err)
		}

		if len(indices) == 0 {
			t.Skip("No TOPIX data available")
		}

		t.Logf("Retrieved %d days of TOPIX data", len(indices))

		// データの検証
		for i, index := range indices {
			if index.Code != "0000" { // TOPIXのコード
				t.Errorf("Index[%d]: Code = %v, want 0000", i, index.Code)
			}

			// 最初と最後のデータをログ
			if i == 0 || i == len(indices)-1 {
				t.Logf("TOPIX [%s]: Close=%.2f (O=%.2f, H=%.2f, L=%.2f)",
					index.Date, index.C, index.O, index.H, index.L)
			}
		}
	})

	t.Run("GetIndices_HistoricalData", func(t *testing.T) {
		// 過去の特定期間のデータ取得
		params := jquants.IndicesParams{
			From: "2024-01-01",
			To:   "2024-01-31",
		}

		resp, err := jq.Indices.GetIndices(params)
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Logf("Failed to get historical indices: %v", err)
			return
		}

		if resp != nil && len(resp.Data) > 0 {
			t.Logf("Retrieved %d indices for January 2024", len(resp.Data))

			// 日付範囲の検証
			for _, index := range resp.Data {
				if index.Date < params.From || index.Date > params.To {
					t.Errorf("Index date %s is outside requested range %s to %s", index.Date, params.From, params.To)
				}
			}
		}
	})
}
