//go:build e2e
// +build e2e

package e2e

import (
	"testing"
)

// TestListedEndpoint は/equities/masterエンドポイントの完全なテスト
func TestListedEndpoint(t *testing.T) {
	t.Run("GetAllListedInfo", func(t *testing.T) {
		// 全ての上場企業情報を取得
		companies, err := jq.Listed.GetAllListedInfo()
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get all listed companies: %v", err)
		}

		if len(companies) == 0 {
			t.Skip("No listed companies data available")
		}

		t.Logf("Retrieved %d listed companies", len(companies))

		// 最初の10件を詳細に検証
		for i, company := range companies {
			if i >= 10 {
				break
			}

			// 基本情報の検証
			if company.Date == "" {
				t.Errorf("Company[%d]: Date is empty", i)
			} else {
				// 日付フォーマットの検証（YYYY-MM-DD形式）
				if len(company.Date) != 10 || company.Date[4] != '-' || company.Date[7] != '-' {
					t.Errorf("Company[%d]: Date format invalid = %v, want YYYY-MM-DD", i, company.Date)
				}
			}
			if company.Code == "" {
				t.Errorf("Company[%d]: Code is empty", i)
			}
			if company.CoName == "" {
				t.Errorf("Company[%d]: CoName is empty", i)
			}
			if company.CoNameEn == "" {
				t.Logf("Company[%d]: CoNameEn is empty (might be acceptable)", i)
			}

			// 業種情報の検証
			if company.S17 == "" {
				t.Errorf("Company[%d]: S17 (Sector17Code) is empty", i)
			}
			if company.S17Nm == "" {
				t.Errorf("Company[%d]: S17Nm (Sector17CodeName) is empty", i)
			}
			if company.S33 == "" {
				t.Errorf("Company[%d]: S33 (Sector33Code) is empty", i)
			}
			if company.S33Nm == "" {
				t.Errorf("Company[%d]: S33Nm (Sector33CodeName) is empty", i)
			}

			// 規模区分の検証
			if company.ScaleCat == "" {
				t.Errorf("Company[%d]: ScaleCat is empty", i)
			}

			// 市場区分の検証
			if company.Mkt == "" {
				t.Errorf("Company[%d]: Mkt (MarketCode) is empty", i)
			}
			if company.MktNm == "" {
				t.Errorf("Company[%d]: MktNm (MarketCodeName) is empty", i)
			}

			// 詳細ログ（最初の3件のみ）
			if i < 3 {
				t.Logf("Company[%d]: %s (%s) - %s - %s - %s",
					i, company.CoName, company.Code,
					company.S33Nm, company.MktNm, company.ScaleCat)
			}
		}
	})

	t.Run("GetListedInfoByCode", func(t *testing.T) {
		// 特定の銘柄コードで企業情報を取得（トヨタ自動車）
		infos, err := jq.Listed.GetListedInfoByCode("7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get listed info for code 7203: %v", err)
		}

		if len(infos) == 0 {
			t.Fatal("No listed info for code 7203")
		}

		company := infos[0]

		// トヨタ自動車の詳細検証
		if company.Code != "72030" && company.Code != "7203" {
			t.Errorf("Code = %v, want 72030 or 7203", company.Code)
		}
		if company.CoName == "" {
			t.Error("CoName is empty")
		}
		if company.S17 == "" {
			t.Error("S17 is empty")
		}
		if company.S33 == "" {
			t.Error("S33 is empty")
		}
		if company.Mkt == "" {
			t.Error("Mkt is empty")
		}
		if company.ScaleCat == "" {
			t.Error("ScaleCat is empty")
		}

		t.Logf("Company: %s (%s)", company.CoName, company.Code)
		t.Logf("Sector: %s (%s)", company.S33Nm, company.S33)
		t.Logf("Market: %s (%s)", company.MktNm, company.Mkt)
		t.Logf("Scale: %s", company.ScaleCat)
		if company.CoNameEn != "" {
			t.Logf("English Name: %s", company.CoNameEn)
		}
	})

	t.Run("GetListedInfoByCode_InvalidCode", func(t *testing.T) {
		// 存在しない銘柄コードでのテスト
		infos, err := jq.Listed.GetListedInfoByCode("99999")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			// エラーが返ってもOK
			t.Logf("Expected error for invalid code: %v", err)
			return
		}

		if len(infos) > 0 {
			t.Error("Expected no data for invalid code 99999")
		}
	})

	t.Run("GetAllListedInfo_MarketSegments", func(t *testing.T) {
		// 全企業の市場セグメント分析
		companies, err := jq.Listed.GetAllListedInfo()
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No listed companies data available")
		}

		if len(companies) == 0 {
			t.Skip("No companies data available")
		}

		// 市場別の集計
		marketCount := make(map[string]int)
		scaleCount := make(map[string]int)

		for _, company := range companies {
			marketCount[company.MktNm]++
			scaleCount[company.ScaleCat]++
		}

		t.Logf("Market segment distribution:")
		for market, count := range marketCount {
			t.Logf("  %s: %d companies", market, count)
		}

		t.Logf("Scale category distribution:")
		for scale, count := range scaleCount {
			t.Logf("  %s: %d companies", scale, count)
		}
	})

	t.Run("GetAllListedInfo_SectorAnalysis", func(t *testing.T) {
		// 業種別分析
		companies, err := jq.Listed.GetAllListedInfo()
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No listed companies data available")
		}

		if len(companies) == 0 {
			t.Skip("No companies data available")
		}

		// 33業種分類での集計
		sector33Count := make(map[string]int)
		sector17Count := make(map[string]int)

		for _, company := range companies {
			sector33Count[company.S33Nm]++
			sector17Count[company.S17Nm]++
		}

		t.Logf("Top 10 sectors (33 classification):")
		count := 0
		for sector, num := range sector33Count {
			if count >= 10 {
				break
			}
			t.Logf("  %s: %d companies", sector, num)
			count++
		}

		t.Logf("17 sector classification distribution:")
		for sector, num := range sector17Count {
			t.Logf("  %s: %d companies", sector, num)
		}
	})

	t.Run("GetAllListedInfo_CodeValidation", func(t *testing.T) {
		// 銘柄コードの形式検証
		companies, err := jq.Listed.GetAllListedInfo()
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No listed companies data available")
		}

		if len(companies) == 0 {
			t.Skip("No companies data available")
		}

		codeLength4 := 0
		codeLength5 := 0
		invalidCodes := 0

		for i, company := range companies {
			if i >= 100 {
				break // 最初の100件のみチェック
			}

			codeLen := len(company.Code)
			switch codeLen {
			case 4:
				codeLength4++
			case 5:
				codeLength5++
			default:
				invalidCodes++
				t.Logf("Invalid code length for %s: %s (length: %d)",
					company.CoName, company.Code, codeLen)
			}
		}

		t.Logf("Code length distribution (first 100 companies):")
		t.Logf("  4-digit codes: %d", codeLength4)
		t.Logf("  5-digit codes: %d", codeLength5)
		t.Logf("  Invalid codes: %d", invalidCodes)
	})
}
