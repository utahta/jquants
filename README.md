# J-Quants Go Client

J-Quants API v2 の Go クライアントライブラリです。株価・財務情報など、日本の株式市場データを取得できます。

利用には[J-Quants](https://jpx-jquants.com/)のAPIキーが必要です（取得できるデータは契約プランに応じます）。

## 特徴

- 📊 株価、財務情報、指数などのデータ取得
- 🔐 APIキーによるシンプルな認証
- 📄 ページネーション対応
- 🚀 セッション単位のキャッシュ機能（オプション）

## インストール

Go 1.25.0 以上が必要です。

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
        if quote.C != nil {
            fmt.Printf("%s: 終値 %.2f円\n", quote.Date, *quote.C)
        }
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

キャッシュを有効にすると、同じ HTTP クライアントで同じリクエストを行った際に、取得済みのレスポンスを再利用します。

```go
httpClient, err := client.NewClientFromEnv(client.WithCache())
if err != nil {
    log.Fatal(err)
}

// キャッシュをクリア
httpClient.ClearCache()

// キャッシュエントリ数を取得
fmt.Printf("キャッシュ件数: %d\n", httpClient.CacheSize())
```

API キーを直接指定する場合は、`client.NewClient("your-api-key", client.WithCache())` を使用します。

- キャッシュはGETリクエストのみに適用されます
- URL とクエリパラメータが同じリクエストでキャッシュを再利用します
- 同じリクエストを同時に実行した場合、API 呼び出しを1回にまとめます
- レスポンスの待機中にコンテキストをキャンセルすると、その呼び出しは `ctx.Err()` を返します。同じレスポンスを待つ他の呼び出しは継続します
- 署名付きダウンロードURLを返すAPI（Bulk・TDnetのファイルURL取得）はURLが失効するためキャッシュを経由しません
- キャッシュは自動更新されません。最新データを再取得する場合は `ClearCache()` を呼び出すか、新しい HTTP クライアントを作成してください

## エラーハンドリング

API がエラーステータスを返した場合は、`errors.As` で `*client.APIError` を取り出し、ステータスコードとレスポンスボディを確認できます。

```go
_, err := jq.Quotes.GetDailyQuotesByCode(ctx, "7203")

var apiErr *client.APIError
if errors.As(err, &apiErr) {
    log.Printf("status=%d body=%s", apiErr.StatusCode, apiErr.Body)
} else if err != nil {
    log.Printf("データ取得に失敗しました: %v", err)
}
```

利用頻度の高い判定にはヘルパーを用意しています。

| ヘルパー | 対象ステータス |
|---------|--------------|
| `client.IsRateLimitExceeded(err)` | 429 |
| `client.IsAuthError(err)` | 401, 403 |
| `client.IsServerError(err)` | 5xx |
| `client.StatusCode(err)` | 任意（`(code int, ok bool)` を返す） |

```go
for {
    quotes, err := jq.Quotes.GetDailyQuotesByCode(ctx, "7203")
    switch {
    case err == nil:
        return quotes, nil
    case client.IsRateLimitExceeded(err):
        time.Sleep(time.Minute) // 待機してから再試行
    case client.IsServerError(err):
        time.Sleep(10 * time.Second)
    default:
        return nil, err // 認証エラー・パラメータ不正などは再試行しても解決しない
    }
}
```

- 403はAPIキーが不正な場合のほか、契約プランに含まれないデータへアクセスした場合にも返ります
- データが存在しない場合に210を返すAPI（前場四本値など）があり、これも `*client.APIError`（`StatusCode` が210）になります
- 通信エラーやレスポンスのデコードエラーは `*client.APIError` にならないため、`client.StatusCode` の `ok` は false になります

## 利用可能なAPI

このライブラリでは以下のAPIエンドポイントにアクセスできます。

| カテゴリ | サービス | 説明 |
|---------|---------|------|
| **株価** | Quotes | 日次株価四本値 |
| **株価** | Valuation | 日次バリュエーション指標・時価総額 |
| **株価** | PricesAM | 前場四本値 |
| **株価** | MinuteQuotes | 株価分足 |
| **銘柄情報** | Listed | 上場銘柄一覧 |
| **財務** | Statements | 財務情報 |
| **財務** | FSDetails | 財務諸表詳細（BS/PL/CF） |
| **財務** | Dividend | 配当金情報 |
| **財務** | EarningsDate | 決算発表予定日 |
| **財務** | Announcement | 翌営業日の決算発表予定（3月期・9月期決算会社のみ） |
| **指数** | Indices | 指数四本値 |
| **指数** | TOPIX | TOPIX指数四本値 |
| **デリバティブ** | Futures | 先物四本値 |
| **デリバティブ** | Options | オプション四本値 |
| **デリバティブ** | IndexOption | 日経225オプション |
| **市場統計** | TradesSpec | 投資部門別売買状況 |
| **市場統計** | Breakdown | 売買内訳データ |
| **市場統計** | TradingCalendar | 取引カレンダー |
| **信用取引** | WeeklyMarginInterest | 信用取引残高 |
| **信用取引** | DailyMarginInterest | 日々公表信用取引残高 |
| **空売り** | ShortSelling | 業種別空売り比率 |
| **空売り** | ShortSellingPositions | 空売り残高報告 |
| **適時開示** | TimelyDisclosure | TDnet適時開示情報 |
| **EDINET** | EdinetMajorShareholders | 大株主状況（有価証券・半期・四半期報告書） |
| **EDINET** | EdinetCrossShareholdings | 政策保有株式 |
| **EDINET** | EdinetLargeVolumeShareholders | 大量保有報告書・変更報告書・訂正報告書 |
| **ダウンロード** | Bulk | CSV一括ダウンロード |

※各APIの利用可能なプランについては、[J-Quants公式サイト](https://jpx-jquants.com/)で確認してください。

## 使用例

以下の例では、クイックスタートと同様に作成した `jq` と `ctx` を使用します。日付は契約プランで取得できる期間に合わせて指定してください。ポインタ型のフィールドは、値がない場合に `nil` になります。

### 日次株価を取得

```go
// 特定銘柄の取得可能な全期間のデータ
quotes, err := jq.Quotes.GetDailyQuotesByCode(ctx, "7203")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("取得件数: %d\n", len(quotes))

// 日付範囲指定
params := jquants.DailyQuotesParams{
    Code: "7203",
    From: "2024-01-01",
    To:   "2024-01-31",
}
response, err := jq.Quotes.GetDailyQuotes(ctx, params)
if err != nil {
    log.Fatal(err)
}

for _, q := range response.Data {
    if q.O != nil && q.H != nil && q.L != nil && q.C != nil {
        fmt.Printf("日付: %s, 始値: %.2f, 高値: %.2f, 安値: %.2f, 終値: %.2f\n",
            q.Date, *q.O, *q.H, *q.L, *q.C)
    }
}
```

### バリュエーション指標を取得

```go
values, err := jq.Valuation.GetValuationsByCodeAndDateRange(ctx, "7203", "2026-08-01", "2026-08-31")
if err != nil {
    log.Fatal(err)
}
for _, v := range values {
    if v.PER != nil && v.PBR != nil {
        fmt.Printf("%s: PER %.2f倍、PBR %.2f倍\n", v.Date, *v.PER, *v.PBR)
    }
}
```

`EPS`・`FwdEPS`・`BPS`・`ROE`・`FwdROE`・`PER`・`FwdPER`・`PBR`・`MktCap` を取得できます。数値項目はポインタで、算出できない項目は `nil` です。ROE は小数表現（`0.1` = 10%）、時価総額は百万円単位です。[指標の公式仕様](https://jpx-jquants.com/ja/spec/eq-valuation)

取得可能な期間とデータの遅延は契約プランによって異なります。[公式のプラン別仕様](https://jpx-jquants.com/ja/spec/data-spec)を確認してください。

`Valuation.MktCap` は自己株式を除く株式数で算出されます。自己株式を含む `DailyQuote.MktCap` とは定義が異なります。株価四本値 API の `MktCap` は削除予定ですが、時期は未定です。[削除予告](https://jpx-jquants.com/ja/spec/release)

### 上場銘柄一覧を取得

```go
// 全銘柄を取得
companies, err := jq.Listed.GetAllListedInfo(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("銘柄数: %d\n", len(companies))

// 特定銘柄の情報を取得
companies, err = jq.Listed.GetListedInfoByCode(ctx, "7203")
if err != nil {
    log.Fatal(err)
}
for _, company := range companies {
    fmt.Printf("企業名: %s\n", company.CoName)
}

// 市場区分で絞り込み（定数を使用）
primeCompanies, err := jq.Listed.GetListedByMarket(ctx, jquants.MarketPrime, "")
if err != nil {
    log.Fatal(err)
}
for _, company := range primeCompanies {
    fmt.Printf("%s (%s) - %s\n", company.CoName, company.Code, company.MktNm)
}
```

### 財務情報を取得

```go
// 最新の財務情報
statement, err := jq.Statements.GetLatestStatements(ctx, "7203")
if err != nil {
    log.Fatal(err)
}
if statement.Sales != nil {
    fmt.Printf("売上高: %.0f円\n", *statement.Sales)
}

// 特定日の財務情報
statements, err := jq.Statements.GetStatementsByDate(ctx, "2024-01-15")
if err != nil {
    log.Fatal(err)
}

// 開示書類種別での絞り込み
for _, stmt := range statements {
    if stmt.DocType.IsConsolidated() && stmt.DocType.GetAccountingStandard() == "IFRS" {
        fmt.Printf("%s: IFRS採用企業\n", stmt.Code)
    }
}
```

### 信用取引残高を公表日で取得

```go
data, err := jq.WeeklyMarginInterest.GetWeeklyMarginInterestByPublishedDate(ctx, "2026-09-28")
if err != nil {
    log.Fatal(err)
}
for _, d := range data {
    if d.PubDate != nil && d.LongVal != nil {
        fmt.Printf("%s: 公表日 %s、買残高金額 %.0f\n", d.Code, *d.PubDate, *d.LongVal)
    }
}
```

信用取引残高は `WeeklyMarginInterest` で取得できます。`GetWeeklyMarginInterest` の `PublishedDate` は `Code` と併用できますが、`Date`・`From`・`To` とは併用できません。

日次データ、公表日 `PubDate`、金額6項目（`ShrtVal`・`LongVal`・`ShrtNegVal`・`LongNegVal`・`ShrtStdVal`・`LongStdVal`）の対象は2026年9月25日申込分以降です。それ以前のデータは週末残高のみで、公表日と金額は `nil` になります。公表日未収録の過去データは `PublishedDate` では検索できません。[公式仕様](https://jpx-jquants.com/ja/spec/mkt-margin-int-daily)

### EDINET の会計期間・訂正情報

`EdinetMajorShareholderDoc.CurPerSt`・`CurPerEn` は当会計期間の開始日・終了日です。事業年度を表す `PerSt`・`PerEn` と区別できます。`DocTypeCode` は有価証券報告書 `120`、四半期報告書 `140`、半期報告書 `160` に対応します。[大株主状況の仕様](https://jpx-jquants.com/ja/spec/edinet-major-shareholders)

大量保有の訂正報告書は `doc.LargeHldgTypeCode == jquants.LargeHldgTypeCodeCorrectionReport` で判定できます。`EdinetLargeVolumeShareholderDoc.ParDocId` は訂正元書類の管理番号（訂正報告書以外は `nil`）、`RptOblgDate` は報告義務発生日です。訂正報告書は訂正元を置き換えず、別レコードとして返されます。[大量保有報告書の仕様](https://jpx-jquants.com/ja/spec/edinet-large-volume-shareholders)

### 日々公表信用取引残高を取得

```go
// 銘柄コードで取得
data, err := jq.DailyMarginInterest.GetDailyMarginInterestByCode(ctx, "1326")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("取得件数: %d\n", len(data))

// 公表日で取得
data, err = jq.DailyMarginInterest.GetDailyMarginInterestByDate(ctx, "2024-02-08")
if err != nil {
    log.Fatal(err)
}

// 公表理由の確認
for _, d := range data {
    if d.PubReason.IsPrecautionByJSF() {
        fmt.Printf("%s: 日証金注意喚起銘柄\n", d.Code)
    }
}
```

### ページネーション対応

`GetDailyQuotesByCode` や `GetDailyQuotesByDate` は全ページを取得します。`GetDailyQuotes` で1ページずつ取得する場合は、レスポンスの `PaginationKey` を次のリクエストに指定します。

```go
params := jquants.DailyQuotesParams{
    Date: "2024-01-15",
}

var quotes []jquants.DailyQuote
for {
    response, err := jq.Quotes.GetDailyQuotes(ctx, params)
    if err != nil {
        log.Fatal(err)
    }
    quotes = append(quotes, response.Data...)
    if response.PaginationKey == "" {
        break
    }
    params.PaginationKey = response.PaginationKey
}
fmt.Printf("取得件数: %d\n", len(quotes))
```

## 開発

### 必要な環境

- Go 1.25.0 以上
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

# E2Eテスト（APIキーが必要）
make test-e2e
```

### プロジェクト構造

```
jquants/
├── client/        # HTTPクライアント（認証含む）
├── types/         # カスタム型定義
├── docs/v2/       # 公式APIドキュメントのローカルキャッシュ（make docs-sync で取得）
├── scripts/       # 開発用スクリプト
├── test/e2e/      # E2Eテスト
├── *.go           # 各APIサービス実装
└── Makefile       # ビルドタスク
```

## v1からv2への移行

J-Quants API v2では以下の変更があります。

### 認証方式の変更

| 項目 | v1 | v2 |
|------|-----|-----|
| 認証方式 | トークン方式（ID Token/Refresh Token） | APIキー方式（x-api-key） |
| 認証パッケージ | `auth/` パッケージを使用 | `client` パッケージに統合 |
| 環境変数 | `JQUANTS_EMAIL`, `JQUANTS_PASSWORD` | `JQUANTS_API_KEY` |

### レスポンス形式の変更

- レスポンスキー: 各API固有のキー → 統一された `data` キー
- フィールド名: 短縮形に変更（例: `Open` → `O`, `Close` → `C`, `Volume` → `Vo`）

### コード例

```go
// v1
quote.Open, quote.High, quote.Low, quote.Close

// v2
quote.O, quote.H, quote.L, quote.C
```

詳細は[公式の移行ガイド](https://jpx-jquants.com/ja/spec/migration-v1-v2)を参照してください。

## 注意事項

- J-Quants APIの利用には適切なサブスクリプションが必要です
- リクエスト数の上限はプラン・API・アドオンによって異なります。[公式のレートリミット](https://jpx-jquants.com/ja/spec/rate-limits)を確認してください
- データの更新日時や取得可能な期間は API と契約プランによって異なります
- 詳細なAPI仕様は[公式ドキュメント](https://jpx-jquants.com/ja/spec/)を参照してください

## ライセンス

MITライセンス

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを作成して変更内容を議論してください。

## サポート

- [Issue Tracker](https://github.com/utahta/jquants/issues)
- [J-Quants公式サイト](https://jpx-jquants.com/)
