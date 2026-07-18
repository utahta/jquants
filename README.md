# J-Quants Go Client

J-Quants API v2のGo言語クライアントライブラリです。日本の株式市場データに簡単にアクセスできます。

## 特徴

- 📊 包括的な市場データアクセス（株価、財務情報、指数など）
- 🔐 APIキーによるシンプルな認証
- 📄 ページネーション対応
- 🚀 セッション単位のキャッシュ機能（オプション）
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
    "context"
    "fmt"
    "log"

    "github.com/utahta/jquants"
    "github.com/utahta/jquants/client"
)

func main() {
    // HTTPクライアントを作成（環境変数 JQUANTS_API_KEY から自動取得）
    httpClient, err := client.NewClientFromEnv()
    if err != nil {
        log.Fatal(err)
    }

    // J-Quants APIクライアントを作成
    jq := jquants.NewJQuantsAPI(httpClient)

    // 株価データを取得
    ctx := context.Background()
    quotes, err := jq.Quotes.GetDailyQuotesByCode(ctx, "7203") // トヨタ自動車
    if err != nil {
        log.Fatal(err)
    }

    for _, quote := range quotes {
        fmt.Printf("%s: 終値 %.2f円\n", quote.Date, *quote.C)
    }
}
```

## 認証設定

J-Quants API v2ではAPIキー方式を使用します。

### 環境変数による設定（推奨）

```bash
export JQUANTS_API_KEY="your-api-key"
```

APIキーは[J-Quantsダッシュボード](https://jpx-jquants.com/)から取得できます。

### 直接指定する場合

```go
httpClient := client.NewClient("your-api-key")
```

## キャッシュ機能

セッション単位のキャッシュを有効にすると、同じリクエストに対するAPI呼び出しを削減できます。

```go
// キャッシュを有効化
httpClient := client.NewClient("your-api-key", client.WithCache())

// または環境変数から
httpClient, err := client.NewClientFromEnv(client.WithCache())

// キャッシュをクリア
httpClient.ClearCache()

// キャッシュエントリ数を取得
size := httpClient.CacheSize()
```

- キャッシュはGETリクエストのみに適用されます
- キャッシュキーはリクエストパス（クエリパラメータ含む）で区別されます
- 同時リクエストの重複排除（singleflight）により効率的に動作します
- 重複排除の待機中にコンテキストをキャンセルした場合、その呼び出し元だけが即座にエラーで戻り、進行中のリクエスト自体は完了してキャッシュに保存されます
- キャッシュはクライアントインスタンスの生存期間のみ有効です

## 利用可能なAPI

このライブラリでアクセスできるAPIエンドポイント：

| カテゴリ | サービス | 説明 |
|---------|---------|------|
| **株価** | Quotes | 日次株価四本値 |
| **株価** | PricesAM | 前場四本値 |
| **銘柄情報** | Listed | 上場銘柄一覧 |
| **財務** | Statements | 財務情報 |
| **財務** | FSDetails | 財務諸表詳細（BS/PL/CF） |
| **財務** | Dividend | 配当金情報 |
| **財務** | Announcement | 決算発表予定 |
| **指数** | Indices | 指数四本値 |
| **指数** | TOPIX | TOPIX指数四本値 |
| **デリバティブ** | Futures | 先物四本値 |
| **デリバティブ** | Options | オプション四本値 |
| **デリバティブ** | IndexOption | 日経225オプション |
| **市場統計** | TradesSpec | 投資部門別売買状況 |
| **市場統計** | Breakdown | 売買内訳データ |
| **市場統計** | TradingCalendar | 取引カレンダー |
| **信用取引** | WeeklyMarginInterest | 信用取引週末残高 |
| **信用取引** | DailyMarginInterest | 日々公表信用取引残高 |
| **空売り** | ShortSelling | 業種別空売り比率 |
| **空売り** | ShortSellingPositions | 空売り残高報告 |

※各APIの利用可能なプランについては、[J-Quants公式サイト](https://jpx-jquants.com/)でご確認ください。

## 使用例

### 日次株価を取得

```go
// 特定銘柄の直近データ
quotes, err := jq.Quotes.GetDailyQuotesByCode(ctx, "7203")

// 日付範囲指定
params := jquants.DailyQuotesParams{
    Code: "7203",
    From: "2024-01-01",
    To:   "2024-01-31",
}
response, err := jq.Quotes.GetDailyQuotes(ctx, params)

// v2 APIでは短縮フィールド名を使用
for _, q := range response.Data {
    fmt.Printf("日付: %s, 始値: %.2f, 高値: %.2f, 安値: %.2f, 終値: %.2f\n",
        q.Date, *q.O, *q.H, *q.L, *q.C)
}
```

### 上場銘柄一覧を取得

```go
// 全銘柄を取得
companies, err := jq.Listed.GetAllListedInfo(ctx)

// 特定銘柄の情報を取得
companies, err := jq.Listed.GetListedInfoByCode(ctx, "7203")
fmt.Printf("企業名: %s\n", companies[0].Name)

// 市場区分で絞り込み（定数を使用）
primeCompanies, err := jq.Listed.GetListedByMarket(ctx, jquants.MarketPrime, "")
for _, company := range primeCompanies {
    fmt.Printf("%s (%s) - %s\n", company.Name, company.Code, company.MktName)
}
```

### 財務情報を取得

```go
// 最新の財務情報
statement, err := jq.Statements.GetLatestStatements(ctx, "7203")
if statement.NetSales != nil {
    fmt.Printf("売上高: %.0f円\n", *statement.NetSales)
}

// 特定日の財務情報
statements, err := jq.Statements.GetStatementsByDate(ctx, "2024-01-15")

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

### 日々公表信用取引残高を取得

```go
// 銘柄コードで取得
data, err := jq.DailyMarginInterest.GetDailyMarginInterestByCode(ctx, "13260")

// 公表日で取得
data, err := jq.DailyMarginInterest.GetDailyMarginInterestByDate(ctx, "20240208")

// 公表理由の確認
for _, d := range data {
    if d.PubReason.IsPrecautionByJSF() {
        fmt.Printf("%s: 日証金注意喚起銘柄\n", d.Code)
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
response, err := jq.Quotes.GetDailyQuotes(ctx, params)
quotes := response.Data

// 次のページがある場合
if response.PaginationKey != "" {
    params.PaginationKey = response.PaginationKey
    nextResponse, err := jq.Quotes.GetDailyQuotes(ctx, params)
    quotes = append(quotes, nextResponse.Data...)
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

# E2Eテスト（API Keyが必要）
make test-e2e
```

### プロジェクト構造

```
jquants/
├── client/        # HTTPクライアント（認証含む）
├── types/         # カスタム型定義
├── docs/          # APIドキュメント
│   └── v2/        # v2 APIドキュメント
├── test/e2e/      # E2Eテスト
├── cmd/           # コマンドラインツール
│   └── gitbook2md/  # GitBook→Markdown変換ツール
├── *.go           # 各APIサービス実装
└── Makefile       # ビルドタスク
```

## V1からV2への移行

J-Quants API v2では以下の変更があります：

### 認証方式の変更

| 項目 | V1 | V2 |
|------|-----|-----|
| 認証方式 | トークン方式（ID Token/Refresh Token） | APIキー方式（x-api-key） |
| 認証パッケージ | `auth/` パッケージを使用 | `client` パッケージに統合 |
| 環境変数 | `JQUANTS_EMAIL`, `JQUANTS_PASSWORD` | `JQUANTS_API_KEY` |

### レスポンス形式の変更

- レスポンスキー: 各API固有のキー → 統一された `data` キー
- フィールド名: 短縮形に変更（例: `Open` → `O`, `Close` → `C`, `Volume` → `Vo`）

### コード例

```go
// V1
quote.Open, quote.High, quote.Low, quote.Close

// V2
quote.O, quote.H, quote.L, quote.C
```

詳細は `docs/v2/migration-v1-v2.md` を参照してください。

## ツール

### gitbook2md

GitBookのドキュメントをMarkdown形式に変換するツールです。

```bash
# ビルド
cd cmd/gitbook2md && go build

# URLから直接変換
./gitbook2md --url https://jpx.gitbook.io/j-quants-ja/api-reference/statements --output statements.md
```

## エラーハンドリング

```go
quotes, err := jq.Quotes.GetDailyQuotesByCode(ctx, "9999")
if err != nil {
    // APIエラーの詳細を取得
    if apiErr, ok := err.(*client.APIError); ok {
        fmt.Printf("APIエラー: %d - %s\n", apiErr.StatusCode, apiErr.Message)
    }
}
```

## 注意事項

- J-Quants APIの利用には適切なサブスクリプションが必要です
- プランごとにレートリミットが設定されています（Free: 5/分, Light: 60/分, Standard: 120/分, Premium: 500/分）
- 営業日以外はデータが取得できない場合があります
- 詳細なAPI仕様は[公式ドキュメント](https://jpx.gitbook.io/j-quants-ja/api-reference)を参照してください

## ライセンス

MITライセンス

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを作成して変更内容を議論してください。

## サポート

- [Issue Tracker](https://github.com/utahta/jquants/issues)
- [J-Quants公式サイト](https://jpx-jquants.com/)
