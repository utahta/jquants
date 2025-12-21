# J-Quants Go Client

J-Quants APIのGo言語クライアントライブラリです。日本の株式市場データに簡単にアクセスできます。

## 特徴

- 📊 包括的な市場データアクセス（株価、財務情報、指数など）
- 🔐 自動的な認証とトークン管理
- 📄 ページネーション対応
- 🧪 充実したテストカバレッジ
- 📝 詳細なドキュメント

## インストール

```bash
go get github.com/utahta/jquants
```

## クイックスタート

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/utahta/jquants"
    "github.com/utahta/jquants/auth"
    "github.com/utahta/jquants/client"
)

func main() {
    // HTTPクライアントを作成
    httpClient := client.NewClient()
    
    // 認証を初期化（環境変数から）
    authClient := auth.NewAuth(httpClient)
    if err := authClient.InitFromEnv(); err != nil {
        log.Fatal(err)
    }
    
    // J-Quants APIクライアントを作成
    jq := jquants.NewJQuantsAPI(httpClient)
    
    // 株価データを取得
    quotes, err := jq.Quotes.GetDailyQuotesByCode("7203") // トヨタ自動車
    if err != nil {
        log.Fatal(err)
    }
    
    for _, quote := range quotes {
        fmt.Printf("%s: 終値 %.2f円\n", quote.Date, *quote.Close)
    }
}
```

## 認証設定

`InitFromEnv()` は以下の優先順位で認証情報を取得します：

1. 環境変数 `JQUANTS_REFRESH_TOKEN`
2. 設定ファイル `~/.jquants/refresh_token`
3. 環境変数 `JQUANTS_EMAIL` と `JQUANTS_PASSWORD`（自動ログイン）

### 推奨: Email/Passwordによる自動認証

```bash
export JQUANTS_EMAIL="your-email@example.com"
export JQUANTS_PASSWORD="your-password"
```

この設定により、リフレッシュトークンがない場合でも自動的にログインします。
ログイン成功時、リフレッシュトークンは `~/.jquants/refresh_token` に自動保存されます。

### リフレッシュトークンによる認証

```bash
export JQUANTS_REFRESH_TOKEN="your-refresh-token"
```

### 設定ファイルによる認証

リフレッシュトークンを `~/.jquants/refresh_token` に直接保存することもできます。

## 利用可能なAPI

このライブラリでアクセスできる主なAPIエンドポイント：

- **株価情報**: 日次株価四本値
- **上場銘柄一覧**: 上場している全銘柄の情報
- **財務情報**: 企業の財務諸表データ
- **決算発表予定**: 決算発表の予定日
- **指数情報**: TOPIX等の指数データ
- **売買内訳**: 投資部門別売買状況
- **空売り情報**: 空売り残高データ
- **信用取引残高**: 週次信用取引残高
- **営業日カレンダー**: 東証の営業日情報
- **財務詳細情報**: より詳細な財務データ
- **先物取引**: 先物取引データ
- **オプション取引**: オプション取引データ
- **午前終値**: 前場終値データ
- その他

※各APIの利用可能なプランについては、[J-Quants公式サイト](https://jpx-jquants.com/)でご確認ください。

## 使用例

### 日次株価を取得

```go
// 特定銘柄の直近データ
quotes, err := jq.Quotes.GetDailyQuotesByCode("7203")

// 日付範囲指定
params := jquants.DailyQuotesParams{
    Code: "7203",
    From: "2024-01-01",
    To:   "2024-01-31",
}
response, err := jq.Quotes.GetDailyQuotes(params)
```

### 上場銘柄一覧を取得

```go
// 全銘柄を取得
companies, err := jq.Listed.GetInfo()

// 特定銘柄の情報を取得
company, err := jq.Listed.GetInfoByCode("7203")
fmt.Printf("企業名: %s\n", company.CompanyName)

// 市場区分で絞り込み（定数を使用）
primeCompanies, err := jq.Listed.GetListedByMarket(jquants.MarketPrime, "")
for _, company := range primeCompanies {
    fmt.Printf("%s (%s) - %s\n", company.CompanyName, company.Code, company.MarketCodeName)
}

// 他の市場区分の例
standardCompanies, err := jq.Listed.GetListedByMarket(jquants.MarketStandard, "")
growthCompanies, err := jq.Listed.GetListedByMarket(jquants.MarketGrowth, "")
```

### 財務情報を取得

```go
// 最新の財務情報
statement, err := jq.Statements.GetLatestStatements("7203")
if statement.NetSales != nil {
    fmt.Printf("売上高: %.0f円\n", *statement.NetSales)
}

// 特定日の財務情報
statements, err := jq.Statements.GetStatementsByDate("2024-01-15")

// 開示書類種別での絞り込み
for _, stmt := range statements {
    // 連結決算のみ処理
    if stmt.TypeOfDocument.IsConsolidated() {
        // IFRS採用企業のみ
        if stmt.TypeOfDocument.GetAccountingStandard() == "IFRS" {
            fmt.Printf("%s: IFRS採用企業\n", stmt.LocalCode)
        }
    }
}
```

### ページネーション対応

大量のデータを扱うAPIではページネーションがサポートされています：

```go
params := jquants.DailyQuotesParams{
    Date: "2024-01-15",
}

// 最初のページ
response, err := jq.Quotes.GetDailyQuotes(params)
quotes := response.DailyQuotes

// 次のページがある場合
if response.PaginationKey != "" {
    params.PaginationKey = response.PaginationKey
    nextResponse, err := jq.Quotes.GetDailyQuotes(params)
    quotes = append(quotes, nextResponse.DailyQuotes...)
}
```

## 開発

### 必要な環境

- Go 1.24.0以上
- Make（オプション）

### ビルドとテスト

```bash
# 依存関係の取得
go mod download

# コンパイルチェック
make check

# テストの実行
make test

# カバレッジ付きテスト
make test-cover

# リントチェック
make lint

# E2Eテスト（認証情報が必要）
make test-e2e
```

### プロジェクト構造

```
jquants/
├── auth/          # 認証処理
├── client/        # HTTPクライアント
├── types/         # カスタム型定義
├── docs/          # APIドキュメント
├── test/e2e/      # E2Eテスト
├── cmd/           # コマンドラインツール
│   └── gitbook2md/  # GitBook→Markdown変換ツール
├── *.go           # 各APIサービス実装
└── Makefile       # ビルドタスク
```

## ツール

### gitbook2md

GitBookのドキュメントをMarkdown形式に変換するツールです。J-Quants APIの公式ドキュメントをローカルで参照する際に便利です。

```bash
# ビルド
cd cmd/gitbook2md && go build

# 使用方法1: HTMLファイルから変換
./gitbook2md input.html output.md

# 使用方法2: URLから直接変換（推奨）
./gitbook2md --url https://jpx.gitbook.io/j-quants-ja/api-reference/statements --output statements.md
```

#### 使用例

```bash
# 財務情報APIのドキュメントを取得
./gitbook2md --url https://jpx.gitbook.io/j-quants-ja/api-reference/statements --output docs/api/statements.md

# 上場銘柄一覧APIのドキュメントを取得
./gitbook2md --url https://jpx.gitbook.io/j-quants-ja/api-reference/listed_info --output docs/api/listed_info.md
```

## エラーハンドリング

```go
quotes, err := jq.Quotes.GetDailyQuotesByCode("9999")
if err != nil {
    // APIエラーの詳細を取得
    if apiErr, ok := err.(*client.APIError); ok {
        fmt.Printf("APIエラー: %d - %s\n", apiErr.StatusCode, apiErr.Message)
    }
}
```

## 注意事項

- J-Quants APIの利用には適切なサブスクリプションが必要です
- 各APIの利用可能なプランについては、[J-Quants公式サイト](https://jpx-jquants.com/)でご確認ください
- APIレート制限に注意してください
- 営業日以外はデータが取得できない場合があります
- 詳細なAPI仕様は[公式ドキュメント](https://jpx.gitbook.io/j-quants-ja/api-reference)を参照してください

## ライセンス

MITライセンス

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを作成して変更内容を議論してください。

## サポート

- [Issue Tracker](https://github.com/utahta/jquants/issues)
- [J-Quants公式サイト](https://jpx-jquants.com/)