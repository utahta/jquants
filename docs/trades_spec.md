# J-Quants API 投資部門別情報データ項目概要

このドキュメントは、J-Quants APIの投資部門別情報エンドポイント (`/markets/trades_spec`) のレスポンスデータ項目について説明します。

## データ型について
- 日付フィールドは `String` 型で返されます（YYYY-MM-DD形式）
- 金額データは `Number` 型で返されます
- 金額の単位は千円です

## エンドポイント概要

### 投資部門別情報 (/markets/trades_spec)

#### 概要
- **エンドポイント**: `https://api.jquants.com/v1/markets/trades_spec`
- **HTTPメソッド**: GET
- **説明**: 投資部門別売買状況（株式・金額）のデータを取得することができます。

#### 公式情報源
配信データは下記のページで公表している内容と同一です：
- https://www.jpx.co.jp/markets/statistics-equities/investor-type/index.html

#### 本APIの留意点
- 2022年4月4日に行われた市場区分見直しに伴い、市場区分に応じた内容となっている統計資料は、見直し後の市場区分に変更して掲載しています
- 過誤訂正により過去の投資部門別売買状況データが訂正された場合の対応：
  - **2023年4月3日以前の訂正**：訂正前のデータは提供せず、訂正後のデータのみ提供
  - **2023年4月3日以降の訂正**：訂正前と訂正後のデータの両方を提供（公表日で識別可能）
- 過誤訂正により過去のデータが訂正された場合は、過誤訂正が公表された翌営業日にデータが更新されます

#### リクエストパラメータ

データの取得では、セクション（section）または日付（from/to）の指定が可能です。

| パラメータ名 | 型 | 必須 | 説明 | 例 |
|------------|-----|------|------|-----|
| section | String | いいえ | セクション（市場） | TSEPrime |
| from | String | いいえ | 期間の開始日 | 20210901 or 2021-09-01 |
| to | String | いいえ | 期間の終了日 | 20210907 or 2021-09-07 |
| pagination_key | String | いいえ | 検索の先頭を指定する文字列 | 過去の検索で返却されたpagination_key |

#### パラメータの組み合わせパターン

| section | from/to | レスポンスの結果 |
|---------|---------|-----------------|
| 指定あり | 指定あり | 指定したセクションの指定した期間のデータ |
| 指定あり | 指定なし | 指定したセクションの全期間のデータ |
| 指定なし | 指定あり | すべてのセクションの指定した期間のデータ |
| 指定なし | 指定なし | すべてのセクションの全期間のデータ |

#### リクエストヘッダー

| ヘッダー名 | 型 | 必須 | 説明 |
|----------|-----|------|------|
| Authorization | String | はい | アクセスキー |

#### レスポンスフィールド

##### 基本情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| PublishedDate | 公表日 | String | YYYY-MM-DD形式 |
| StartDate | 開始日 | String | YYYY-MM-DD形式 |
| EndDate | 終了日 | String | YYYY-MM-DD形式 |
| Section | 市場名 | String | 市場名一覧参照 |

##### 自己取引

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| ProprietarySales | 自己計_売 | Number | 単位：千円 |
| ProprietaryPurchases | 自己計_買 | Number | 単位：千円 |
| ProprietaryTotal | 自己計_合計 | Number | 単位：千円 |
| ProprietaryBalance | 自己計_差引 | Number | 単位：千円 |

##### 委託取引

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| BrokerageSales | 委託計_売 | Number | 単位：千円 |
| BrokeragePurchases | 委託計_買 | Number | 単位：千円 |
| BrokerageTotal | 委託計_合計 | Number | 単位：千円 |
| BrokerageBalance | 委託計_差引 | Number | 単位：千円 |

##### 総計

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| TotalSales | 総計_売 | Number | 単位：千円 |
| TotalPurchases | 総計_買 | Number | 単位：千円 |
| TotalTotal | 総計_合計 | Number | 単位：千円 |
| TotalBalance | 総計_差引 | Number | 単位：千円 |

##### 投資部門別内訳

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| IndividualsSales | 個人_売 | Number | 単位：千円 |
| IndividualsPurchases | 個人_買 | Number | 単位：千円 |
| IndividualsTotal | 個人_合計 | Number | 単位：千円 |
| IndividualsBalance | 個人_差引 | Number | 単位：千円 |
| ForeignersSales | 海外投資家_売 | Number | 単位：千円 |
| ForeignersPurchases | 海外投資家_買 | Number | 単位：千円 |
| ForeignersTotal | 海外投資家_合計 | Number | 単位：千円 |
| ForeignersBalance | 海外投資家_差引 | Number | 単位：千円 |
| SecuritiesCosSales | 証券会社_売 | Number | 単位：千円 |
| SecuritiesCosPurchases | 証券会社_買 | Number | 単位：千円 |
| SecuritiesCosTotal | 証券会社_合計 | Number | 単位：千円 |
| SecuritiesCosBalance | 証券会社_差引 | Number | 単位：千円 |
| InvestmentTrustsSales | 投資信託_売 | Number | 単位：千円 |
| InvestmentTrustsPurchases | 投資信託_買 | Number | 単位：千円 |
| InvestmentTrustsTotal | 投資信託_合計 | Number | 単位：千円 |
| InvestmentTrustsBalance | 投資信託_差引 | Number | 単位：千円 |
| BusinessCosSales | 事業法人_売 | Number | 単位：千円 |
| BusinessCosPurchases | 事業法人_買 | Number | 単位：千円 |
| BusinessCosTotal | 事業法人_合計 | Number | 単位：千円 |
| BusinessCosBalance | 事業法人_差引 | Number | 単位：千円 |
| OtherCosSales | その他法人_売 | Number | 単位：千円 |
| OtherCosPurchases | その他法人_買 | Number | 単位：千円 |
| OtherCosTotal | その他法人_合計 | Number | 単位：千円 |
| OtherCosBalance | その他法人_差引 | Number | 単位：千円 |
| InsuranceCosSales | 生保・損保_売 | Number | 単位：千円 |
| InsuranceCosPurchases | 生保・損保_買 | Number | 単位：千円 |
| InsuranceCosTotal | 生保・損保_合計 | Number | 単位：千円 |
| InsuranceCosBalance | 生保・損保_差引 | Number | 単位：千円 |
| CityBKsRegionalBKsEtcSales | 都銀・地銀等_売 | Number | 単位：千円 |
| CityBKsRegionalBKsEtcPurchases | 都銀・地銀等_買 | Number | 単位：千円 |
| CityBKsRegionalBKsEtcTotal | 都銀・地銀等_合計 | Number | 単位：千円 |
| CityBKsRegionalBKsEtcBalance | 都銀・地銀等_差引 | Number | 単位：千円 |
| TrustBanksSales | 信託銀行_売 | Number | 単位：千円 |
| TrustBanksPurchases | 信託銀行_買 | Number | 単位：千円 |
| TrustBanksTotal | 信託銀行_合計 | Number | 単位：千円 |
| TrustBanksBalance | 信託銀行_差引 | Number | 単位：千円 |
| OtherFinancialInstitutionsSales | その他金融機関_売 | Number | 単位：千円 |
| OtherFinancialInstitutionsPurchases | その他金融機関_買 | Number | 単位：千円 |
| OtherFinancialInstitutionsTotal | その他金融機関_合計 | Number | 単位：千円 |
| OtherFinancialInstitutionsBalance | その他金融機関_差引 | Number | 単位：千円 |

## 市場名（Section）一覧

| 市場名 | コード |
|--------|--------|
| 市場一部 | TSE1st |
| 市場二部 | TSE2nd |
| マザーズ | TSEMothers |
| JASDAQ | TSEJASDAQ |
| プライム | TSEPrime |
| スタンダード | TSEStandard |
| グロース | TSEGrowth |
| 東証および名証 | TokyoNagoya |

## 使用上の注意事項

1. **データ単位**
   - 金額は千円単位で提供
   - 差引＝買い－売りで計算

2. **市場区分変更の影響**
   - 2022年4月4日以降は新市場区分（プライム、スタンダード、グロース）を使用
   - 過去データも新市場区分で再整理済み

3. **過誤訂正の扱い**
   - 2023年4月3日を境に対応方針が変更
   - 公表日が新しいデータが訂正後データ

4. **投資部門の分類**
   - 個人、海外投資家、証券会社、各種金融機関等に分類
   - 自己取引と委託取引を区別

## サンプルレスポンス

```json
{
    "trades_spec": [
        {
            "PublishedDate":"2017-01-13",
            "StartDate":"2017-01-04",
            "EndDate":"2017-01-06",
            "Section":"TSE1st",
            "ProprietarySales":1311271004,
            "ProprietaryPurchases":1453326508,
            "ProprietaryTotal":2764597512,
            "ProprietaryBalance":142055504,
            "BrokerageSales":7165529005,
            "BrokeragePurchases":7030019854,
            "BrokerageTotal":14195548859,
            "BrokerageBalance":-135509151,
            "TotalSales":8476800009,
            "TotalPurchases":8483346362,
            "TotalTotal":16960146371,
            "TotalBalance":6546353,
            "IndividualsSales":1401711615,
            "IndividualsPurchases":1161801155,
            "IndividualsTotal":2563512770,
            "IndividualsBalance":-239910460,
            "ForeignersSales":5094891735,
            "ForeignersPurchases":5317151774,
            "ForeignersTotal":10412043509,
            "ForeignersBalance":222260039
        }
    ],
    "pagination_key": "value1.value2."
}
```

## 活用例

### 特定市場の最新データを取得
```
/markets/trades_spec?section=TSEPrime&from=20230324&to=20230403
```

### 海外投資家の動向分析
レスポンスの`ForeignersSales`、`ForeignersPurchases`、`ForeignersBalance`を使用

### 個人投資家の動向分析
レスポンスの`IndividualsSales`、`IndividualsPurchases`、`IndividualsBalance`を使用

## 関連情報

このAPIは投資部門別の売買状況を提供するものです。関連する他のAPIについては以下を参照してください：
- 売買内訳データ：`/markets/breakdown`（銘柄別の詳細な売買内訳）
- 株価四本値：`/prices/daily_quotes`（個別銘柄の株価）
- 取引カレンダー：`/markets/trading_calendar`（営業日情報）