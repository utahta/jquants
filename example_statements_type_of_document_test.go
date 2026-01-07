package jquants_test

import (
	"fmt"

	"github.com/utahta/jquants"
)

func ExampleTypeOfDocument_usage() {
	// 定義済みの定数を使用
	docType := jquants.TypeOfDocumentFYConsolidatedIFRS
	fmt.Printf("Document type: %s\n", docType)
	fmt.Printf("Is consolidated: %v\n", docType.IsConsolidated())
	fmt.Printf("Is annual: %v\n", docType.IsAnnual())
	fmt.Printf("Accounting standard: %s\n", docType.GetAccountingStandard())

	// Output:
	// Document type: FYFinancialStatements_Consolidated_IFRS
	// Is consolidated: true
	// Is annual: true
	// Accounting standard: IFRS
}

func ExampleStatement_DocType() {
	// Statementからの直接使用例
	stmt := &jquants.Statement{
		DocType: jquants.ParseTypeOfDocument("2QFinancialStatements_Consolidated_JP"),
	}

	fmt.Printf("Is quarterly: %v\n", stmt.DocType.IsQuarterly())
	fmt.Printf("Period: %s\n", stmt.DocType.GetPeriod())
	fmt.Printf("Accounting standard: %s\n", stmt.DocType.GetAccountingStandard())

	// Output:
	// Is quarterly: true
	// Period: 2Q
	// Accounting standard: JP
}

func ExampleTypeOfDocument_filtering() {
	// 財務諸表のフィルタリング例
	statements := []jquants.Statement{
		{DocType: jquants.TypeOfDocumentFYConsolidatedIFRS},
		{DocType: jquants.TypeOfDocumentFYNonConsolidatedJP},
		{DocType: jquants.TypeOfDocumentDividendRevision},
		{DocType: jquants.TypeOfDocument2QConsolidatedJP},
	}

	// IFRS採用企業の通期決算のみ抽出
	for _, stmt := range statements {
		if stmt.DocType.IsAnnual() && stmt.DocType.GetAccountingStandard() == "IFRS" {
			fmt.Printf("Found IFRS annual report: %s\n", stmt.DocType)
		}
	}

	// Output:
	// Found IFRS annual report: FYFinancialStatements_Consolidated_IFRS
}
