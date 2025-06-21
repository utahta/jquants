package jquants

import (
	"fmt"
	"testing"

	"github.com/utahta/jquants/client"
)

func TestTradesSpecService_GetTradesSpec(t *testing.T) {
	tests := []struct {
		name     string
		params   TradesSpecParams
		wantPath string
	}{
		{
			name: "with all parameters",
			params: TradesSpecParams{
				Section:       "TSEPrime",
				From:          "20230324",
				To:            "20230403",
				PaginationKey: "key123",
			},
			wantPath: "/markets/trades_spec?section=TSEPrime&from=20230324&to=20230403&pagination_key=key123",
		},
		{
			name: "with section and date range",
			params: TradesSpecParams{
				Section: "TSEPrime",
				From:    "20230324",
				To:      "20230403",
			},
			wantPath: "/markets/trades_spec?section=TSEPrime&from=20230324&to=20230403",
		},
		{
			name: "with section only",
			params: TradesSpecParams{
				Section: "TSEStandard",
			},
			wantPath: "/markets/trades_spec?section=TSEStandard",
		},
		{
			name: "with date range only",
			params: TradesSpecParams{
				From: "20230324",
				To:   "20230403",
			},
			wantPath: "/markets/trades_spec?from=20230324&to=20230403",
		},
		{
			name:     "with no parameters",
			params:   TradesSpecParams{},
			wantPath: "/markets/trades_spec",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockClient := client.NewMockClient()
			service := NewTradesSpecService(mockClient)

			// Mock response
			mockResponse := TradesSpecResponse{
				TradesSpec: []TradesSpec{
					{
						PublishedDate:        "2017-01-13",
						StartDate:            "2017-01-04",
						EndDate:              "2017-01-06",
						Section:              "TSE1st",
						ProprietarySales:     1311271004,
						ProprietaryPurchases: 1453326508,
						ProprietaryTotal:     2764597512,
						ProprietaryBalance:   142055504,
						BrokerageSales:       7165529005,
						BrokeragePurchases:   7030019854,
						BrokerageTotal:       14195548859,
						BrokerageBalance:     -135509151,
						TotalSales:           8476800009,
						TotalPurchases:       8483346362,
						TotalTotal:           16960146371,
						TotalBalance:         6546353,
						IndividualsSales:     1401711615,
						IndividualsPurchases: 1161801155,
						IndividualsTotal:     2563512770,
						IndividualsBalance:   -239910460,
						ForeignersSales:      5094891735,
						ForeignersPurchases:  5317151774,
						ForeignersTotal:      10412043509,
						ForeignersBalance:    222260039,
					},
				},
				PaginationKey: "",
			}
			mockClient.SetResponse("GET", tt.wantPath, mockResponse)

			// Execute
			resp, err := service.GetTradesSpec(tt.params)

			// Verify
			if err != nil {
				t.Fatalf("GetTradesSpec() error = %v", err)
			}
			if resp == nil {
				t.Fatal("GetTradesSpec() returned nil response")
			}
			if len(resp.TradesSpec) == 0 {
				t.Error("GetTradesSpec() returned empty data")
			}
			if mockClient.LastPath != tt.wantPath {
				t.Errorf("GetTradesSpec() path = %v, want %v", mockClient.LastPath, tt.wantPath)
			}
		})
	}
}

func TestTradesSpecService_GetTradesSpecByDateRange(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTradesSpecService(mockClient)

	// Mock response - 最初のページ
	mockResponse1 := TradesSpecResponse{
		TradesSpec: []TradesSpec{
			{
				PublishedDate:      "2017-01-13",
				StartDate:          "2017-01-04",
				EndDate:            "2017-01-06",
				Section:            "TSEPrime",
				TotalBalance:       1000000,
				ForeignersBalance:  500000,
				IndividualsBalance: -200000,
			},
			{
				PublishedDate:      "2017-01-20",
				StartDate:          "2017-01-11",
				EndDate:            "2017-01-13",
				Section:            "TSEPrime",
				TotalBalance:       800000,
				ForeignersBalance:  300000,
				IndividualsBalance: -100000,
			},
		},
		PaginationKey: "next_page_key",
	}

	// Mock response - 2ページ目
	mockResponse2 := TradesSpecResponse{
		TradesSpec: []TradesSpec{
			{
				PublishedDate:      "2017-01-27",
				StartDate:          "2017-01-18",
				EndDate:            "2017-01-20",
				Section:            "TSEStandard",
				TotalBalance:       600000,
				ForeignersBalance:  200000,
				IndividualsBalance: -50000,
			},
		},
		PaginationKey: "", // 最後のページ
	}

	mockClient.SetResponse("GET", "/markets/trades_spec?from=20170104&to=20170120", mockResponse1)
	mockClient.SetResponse("GET", "/markets/trades_spec?from=20170104&to=20170120&pagination_key=next_page_key", mockResponse2)

	// Execute
	tradesSpec, err := service.GetTradesSpecByDateRange("20170104", "20170120")

	// Verify
	if err != nil {
		t.Fatalf("GetTradesSpecByDateRange() error = %v", err)
	}
	if len(tradesSpec) != 3 {
		t.Errorf("GetTradesSpecByDateRange() returned %d items, want 3", len(tradesSpec))
	}
}

func TestTradesSpecService_GetTradesSpecBySection(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTradesSpecService(mockClient)

	// Mock response
	mockResponse := TradesSpecResponse{
		TradesSpec: []TradesSpec{
			{
				Section:            "TSEPrime",
				TotalBalance:       1000000,
				ForeignersBalance:  500000,
				IndividualsBalance: -200000,
			},
			{
				Section:            "TSEPrime",
				TotalBalance:       800000,
				ForeignersBalance:  300000,
				IndividualsBalance: -100000,
			},
		},
		PaginationKey: "",
	}
	mockClient.SetResponse("GET", "/markets/trades_spec?section=TSEPrime", mockResponse)

	// Execute
	tradesSpec, err := service.GetTradesSpecBySection("TSEPrime")

	// Verify
	if err != nil {
		t.Fatalf("GetTradesSpecBySection() error = %v", err)
	}
	if len(tradesSpec) != 2 {
		t.Errorf("GetTradesSpecBySection() returned %d items, want 2", len(tradesSpec))
	}
	for _, ts := range tradesSpec {
		if ts.Section != "TSEPrime" {
			t.Errorf("GetTradesSpecBySection() returned section %v, want TSEPrime", ts.Section)
		}
	}
}

func TestTradesSpec_IsBuyerDominant(t *testing.T) {
	tests := []struct {
		name         string
		tradesSpec   TradesSpec
		investorType string
		want         bool
	}{
		{
			name: "individuals buyer dominant",
			tradesSpec: TradesSpec{
				IndividualsBalance: 100000,
			},
			investorType: "individuals",
			want:         true,
		},
		{
			name: "individuals seller dominant",
			tradesSpec: TradesSpec{
				IndividualsBalance: -100000,
			},
			investorType: "individuals",
			want:         false,
		},
		{
			name: "foreigners buyer dominant",
			tradesSpec: TradesSpec{
				ForeignersBalance: 500000,
			},
			investorType: "foreigners",
			want:         true,
		},
		{
			name: "foreigners seller dominant",
			tradesSpec: TradesSpec{
				ForeignersBalance: -300000,
			},
			investorType: "foreigners",
			want:         false,
		},
		{
			name: "total buyer dominant",
			tradesSpec: TradesSpec{
				TotalBalance: 1000000,
			},
			investorType: "total",
			want:         true,
		},
		{
			name: "unknown investor type",
			tradesSpec: TradesSpec{
				TotalBalance: 1000000,
			},
			investorType: "unknown",
			want:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.tradesSpec.IsBuyerDominant(tt.investorType)
			if got != tt.want {
				t.Errorf("TradesSpec.IsBuyerDominant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTradesSpec_GetNetFlow(t *testing.T) {
	tradesSpec := TradesSpec{
		IndividualsBalance:      -100000,
		ForeignersBalance:       500000,
		SecuritiesCosBalance:    50000,
		InvestmentTrustsBalance: 25000,
		BusinessCosBalance:      -20000,
		InsuranceCosBalance:     75000,
		TrustBanksBalance:       30000,
		TotalBalance:            1000000,
	}

	tests := []struct {
		investorType string
		want         float64
	}{
		{"individuals", -100000},
		{"foreigners", 500000},
		{"securities", 50000},
		{"investment_trusts", 25000},
		{"business", -20000},
		{"insurance", 75000},
		{"trust_banks", 30000},
		{"total", 1000000},
		{"unknown", 0},
	}

	for _, tt := range tests {
		t.Run(tt.investorType, func(t *testing.T) {
			got := tradesSpec.GetNetFlow(tt.investorType)
			if got != tt.want {
				t.Errorf("TradesSpec.GetNetFlow() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTradesSpecService_GetTradesSpec_Error(t *testing.T) {
	// Setup
	mockClient := client.NewMockClient()
	service := NewTradesSpecService(mockClient)

	// Set error response
	mockClient.SetError("GET", "/markets/trades_spec", fmt.Errorf("unauthorized"))

	// Execute
	_, err := service.GetTradesSpec(TradesSpecParams{})

	// Verify
	if err == nil {
		t.Error("GetTradesSpec() expected error but got nil")
	}
}

func TestMarketSectionConstants(t *testing.T) {
	// 市場コード定数の値を確認
	tests := []struct {
		constant string
		expected string
	}{
		{SectionTSE1st, "TSE1st"},
		{SectionTSE2nd, "TSE2nd"},
		{SectionTSEMothers, "TSEMothers"},
		{SectionTSEJASDAQ, "TSEJASDAQ"},
		{SectionTSEPrime, "TSEPrime"},
		{SectionTSEStandard, "TSEStandard"},
		{SectionTSEGrowth, "TSEGrowth"},
		{SectionTokyoNagoya, "TokyoNagoya"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("Market section constant = %v, want %v", tt.constant, tt.expected)
			}
		})
	}
}
