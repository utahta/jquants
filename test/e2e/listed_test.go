//go:build e2e
// +build e2e

package e2e

import (
	"testing"
)

// TestListedEndpoint は/listed/infoエンドポイントの完全なテスト
func TestListedEndpoint(t *testing.T) {
	t.Run("GetInfo_All", func(t *testing.T) {
		// 全ての上場企業情報を取得
		companies, err := jq.Listed.GetListedInfo("", "")
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
			// LocalCodeの検証（オプショナル）
			if company.LocalCode != "" {
				// LocalCodeが存在する場合は5桁形式であることを確認
				if len(company.LocalCode) != 5 {
					t.Errorf("Company[%d]: LocalCode length = %d, want 5", i, len(company.LocalCode))
				}
			}
			if company.CompanyName == "" {
				t.Errorf("Company[%d]: CompanyName is empty", i)
			}
			if company.CompanyNameEnglish == "" {
				t.Logf("Company[%d]: CompanyNameEnglish is empty (might be acceptable)", i)
			}
			
			// 業種情報の検証
			if company.Sector17Code == "" {
				t.Errorf("Company[%d]: Sector17Code is empty", i)
			}
			if company.Sector17CodeName == "" {
				t.Errorf("Company[%d]: Sector17CodeName is empty", i)
			}
			if company.Sector33Code == "" {
				t.Errorf("Company[%d]: Sector33Code is empty", i)
			}
			if company.Sector33CodeName == "" {
				t.Errorf("Company[%d]: Sector33CodeName is empty", i)
			}
			
			// 規模区分の検証
			if company.ScaleCategory == "" {
				t.Errorf("Company[%d]: ScaleCategory is empty", i)
			}
			
			// 市場区分の検証
			if company.MarketCode == "" {
				t.Errorf("Company[%d]: MarketCode is empty", i)
			}
			if company.MarketCodeName == "" {
				t.Errorf("Company[%d]: MarketCodeName is empty", i)
			}
			
			// 廃止情報の検証
			if company.IsDelisted != "" {
				// IsDelistedはtrue/falseの値のみ許可
				if company.IsDelisted != "true" && company.IsDelisted != "false" {
					t.Errorf("Company[%d]: IsDelisted = %v, want true or false", i, company.IsDelisted)
				}
				t.Logf("Company[%d]: IsDelisted: %s", i, company.IsDelisted)
			}
			
			// 詳細ログ（最初の3件のみ）
			if i < 3 {
				t.Logf("Company[%d]: %s (%s/%s) - %s - %s - %s",
					i, company.CompanyName, company.Code, company.LocalCode,
					company.Sector33CodeName, company.MarketCodeName, company.ScaleCategory)
			}
		}
	})

	t.Run("GetInfoByCode_ValidCode", func(t *testing.T) {
		// 特定の銘柄コードで企業情報を取得（トヨタ自動車）
		company, err := jq.Listed.GetCompanyInfo("7203")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Fatalf("Failed to get company info for code 7203: %v", err)
		}

		if company == nil {
			t.Fatal("Company info is nil for code 7203")
		}

		// トヨタ自動車の詳細検証
		if company.Code != "72030" && company.Code != "7203" {
			t.Errorf("Code = %v, want 72030 or 7203", company.Code)
		}
		if company.CompanyName == "" {
			t.Error("CompanyName is empty")
		}
		if company.Sector17Code == "" {
			t.Error("Sector17Code is empty")
		}
		if company.Sector33Code == "" {
			t.Error("Sector33Code is empty")
		}
		if company.MarketCode == "" {
			t.Error("MarketCode is empty")
		}
		if company.ScaleCategory == "" {
			t.Error("ScaleCategory is empty")
		}

		t.Logf("Company: %s (%s)", company.CompanyName, company.Code)
		t.Logf("Sector: %s (%s)", company.Sector33CodeName, company.Sector33Code)
		t.Logf("Market: %s (%s)", company.MarketCodeName, company.MarketCode)
		t.Logf("Scale: %s", company.ScaleCategory)
		if company.CompanyNameEnglish != "" {
			t.Logf("English Name: %s", company.CompanyNameEnglish)
		}
		if company.IsDelisted != "" {
			t.Logf("Delisted Status: %s", company.IsDelisted)
		}
	})

	t.Run("GetInfoByCode_InvalidCode", func(t *testing.T) {
		// 存在しない銘柄コードでのテスト
		company, err := jq.Listed.GetCompanyInfo("99999")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			// エラーが返ってもOK
			t.Logf("Expected error for invalid code: %v", err)
			return
		}

		if company != nil {
			t.Error("Expected nil company for invalid code 99999")
		}
	})

	t.Run("GetInfo_MarketSegments", func(t *testing.T) {
		// 全企業の市場セグメント分析
		companies, err := jq.Listed.GetListedInfo("", "")
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
			marketCount[company.MarketCodeName]++
			scaleCount[company.ScaleCategory]++
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

	t.Run("GetInfo_SectorAnalysis", func(t *testing.T) {
		// 業種別分析
		companies, err := jq.Listed.GetListedInfo("", "")
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
			sector33Count[company.Sector33CodeName]++
			sector17Count[company.Sector17CodeName]++
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

	t.Run("GetInfo_ListingStatus", func(t *testing.T) {
		// 上場・廃止ステータス分析
		companies, err := jq.Listed.GetListedInfo("", "")
		if err != nil {
			if isSubscriptionLimited(err) {
				t.Skip("Skipping due to subscription limitation")
			}
			t.Skip("No listed companies data available")
		}

		if len(companies) == 0 {
			t.Skip("No companies data available")
		}

		listedCount := 0
		delistedCount := 0
		
		for _, company := range companies {
			if company.IsDelisted == "true" {
				delistedCount++
				// 廃止企業の最初の5件のみログ出力
				if delistedCount <= 5 {
					t.Logf("Delisted company: %s (%s) - IsDelisted: %s",
						company.CompanyName, company.Code, company.IsDelisted)
				}
			} else if company.IsDelisted == "false" || company.IsDelisted == "" {
				listedCount++
			} else {
				t.Errorf("Invalid IsDelisted value: %s for company %s (%s)",
					company.IsDelisted, company.CompanyName, company.Code)
			}
		}

		t.Logf("Listing status:")
		t.Logf("  Active listings: %d", listedCount)
		t.Logf("  Delisted companies: %d", delistedCount)
		t.Logf("  Total: %d", len(companies))
	})

	t.Run("GetInfo_CodeValidation", func(t *testing.T) {
		// 銘柄コードの形式検証
		companies, err := jq.Listed.GetListedInfo("", "")
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
					company.CompanyName, company.Code, codeLen)
			}
		}

		t.Logf("Code length distribution (first 100 companies):")
		t.Logf("  4-digit codes: %d", codeLength4)
		t.Logf("  5-digit codes: %d", codeLength5)
		t.Logf("  Invalid codes: %d", invalidCodes)
	})
}