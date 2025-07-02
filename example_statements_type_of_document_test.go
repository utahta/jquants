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

func ExampleStatement_TypeOfDocument() {
	// Statementからの直接使用例
	stmt := &jquants.Statement{
		TypeOfDocument: jquants.ParseTypeOfDocument("2QFinancialStatements_Consolidated_JP"),
	}

	fmt.Printf("Is quarterly: %v\n", stmt.TypeOfDocument.IsQuarterly())
	fmt.Printf("Period: %s\n", stmt.TypeOfDocument.GetPeriod())
	fmt.Printf("Accounting standard: %s\n", stmt.TypeOfDocument.GetAccountingStandard())

	// Output:
	// Is quarterly: true
	// Period: 2Q
	// Accounting standard: JP
}

func ExampleTypeOfDocument_filtering() {
	// 財務諸表のフィルタリング例
	statements := []jquants.Statement{
		{TypeOfDocument: jquants.TypeOfDocumentFYConsolidatedIFRS},
		{TypeOfDocument: jquants.TypeOfDocumentFYNonConsolidatedJP},
		{TypeOfDocument: jquants.TypeOfDocumentDividendRevision},
		{TypeOfDocument: jquants.TypeOfDocument2QConsolidatedJP},
	}

	// IFRS採用企業の通期決算のみ抽出
	for _, stmt := range statements {
		if stmt.TypeOfDocument.IsAnnual() && stmt.TypeOfDocument.GetAccountingStandard() == "IFRS" {
			fmt.Printf("Found IFRS annual report: %s\n", stmt.TypeOfDocument)
		}
	}

	// Output:
	// Found IFRS annual report: FYFinancialStatements_Consolidated_IFRS
}