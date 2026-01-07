# J-Quants API 空売り残高報告データ項目概要

このドキュメントは、J-Quants APIの空売り残高報告エンドポイント (`/markets/short_selling_positions`) のレスポンスデータ項目について説明します。

## データ型について
- 日付フィールドは `String` 型で返されます（YYYY-MM-DD形式）
- 残高数量・比率データは `Number` 型で返されます
- 文字列フィールドは `String` 型で返されます

## エンドポイント概要

### 空売り残高報告 (/markets/short_selling_positions)

#### 概要
- **エンドポイント**: `https://api.jquants.com/v1/markets/short_selling_positions`
- **HTTPメソッド**: GET
- **説明**: 「有価証券の取引等の規制に関する内閣府令」に基づき、取引参加者より報告を受けたもののうち、残高割合が0.5％以上のものについての情報を取得できます。

#### 公式情報源
配信データは下記のページで公表している内容と同一ですが、より長いヒストリカルデータを利用可能です：
- https://www.jpx.co.jp/markets/public/short-selling/index.html

#### 法的根拠
有価証券の取引等の規制に関する内閣府令についての詳細は以下を参照：
- https://www.jpx.co.jp/markets/public/short-selling/01.html

#### 本APIの留意点
- **報告義務**: 取引参加者から該当する報告が行われなかった日にはデータは提供されません
- **閾値**: 残高割合が0.5％以上の場合のみ報告対象となります
- **報告タイミング**: 法令に基づく報告スケジュールに従います

#### リクエストパラメータ

データの取得では、銘柄コード（code）、公表日（disclosed_date）、計算日（calculated_date）のいずれかの指定が必須となります。

| パラメータ名 | 型 | 必須 | 説明 | 例 |
|------------|-----|------|------|-----|
| code | String | いいえ* | 4桁もしくは5桁の銘柄コード | 8697 or 86970 |
| disclosed_date | String | いいえ* | 公表日の指定 | 20240301 or 2024-03-01 |
| disclosed_date_from | String | いいえ | 公表日のfrom指定 | 20240301 or 2024-03-01 |
| disclosed_date_to | String | いいえ | 公表日のto指定 | 20240301 or 2024-03-01 |
| calculated_date | String | いいえ* | 計算日の指定 | 20240301 or 2024-03-01 |
| pagination_key | String | いいえ | 検索の先頭を指定する文字列 | 過去の検索で返却されたpagination_key |

*code、disclosed_date、calculated_dateのいずれかが必須

#### 銘柄コード指定時の注意
4桁の銘柄コードを指定した場合は、普通株式と優先株式の両方が上場している銘柄においては普通株式のデータのみが取得されます。

#### パラメータの組み合わせパターン

| code | disclosed_date | disclosed_date_from/to | calculated_date | レスポンスの結果 |
|------|----------------|------------------------|-----------------|-----------------|
| 指定あり | 指定なし | 指定なし | 指定なし | 指定された銘柄について全期間分のデータ |
| 指定あり | 指定あり | - | 指定なし | 指定された銘柄について指定日（公表日）のデータ |
| 指定あり | 指定なし | 指定あり | 指定なし | 指定された銘柄について指定された期間のデータ |
| 指定あり | 指定なし | 指定なし | 指定あり | 指定された銘柄について指定日（計算日）のデータ |
| 指定なし | 指定あり | - | 指定なし | 指定日（公表日）の全ての銘柄のデータ |
| 指定なし | 指定なし | 指定なし | 指定あり | 指定日（計算日）の全ての銘柄のデータ |

#### リクエストヘッダー

| ヘッダー名 | 型 | 必須 | 説明 |
|----------|-----|------|------|
| Authorization | String | はい | アクセスキー |

#### レスポンスフィールド

##### 基本情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| DisclosedDate | 日付（公表日） | String | YYYY-MM-DD形式 |
| CalculatedDate | 日付（計算日） | String | YYYY-MM-DD形式 |
| Code | 銘柄コード | String | 5桁コード |

##### 空売り者情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| ShortSellerName | 商号・名称・氏名 | String | 取引参加者から報告されたものをそのまま記載、日本語名称または英語名称が混在 |
| ShortSellerAddress | 住所・所在地 | String | |

##### 投資一任契約関連情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| DiscretionaryInvestmentContractorName | 委託者・投資一任契約の相手方の商号・名称・氏名 | String | |
| DiscretionaryInvestmentContractorAddress | 委託者・投資一任契約の相手方の住所・所在地 | String | |
| InvestmentFundName | 信託財産・運用財産の名称 | String | |

##### 空売り残高情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| ShortPositionsToSharesOutstandingRatio | 空売り残高割合 | Number | 発行済株式総数に対する比率 |
| ShortPositionsInSharesNumber | 空売り残高数量 | Number | 株数 |
| ShortPositionsInTradingUnitsNumber | 空売り残高売買単位数 | Number | 売買単位数 |

##### 前回報告データ

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| CalculationInPreviousReportingDate | 直近計算年月日 | String | YYYY-MM-DD形式 |
| ShortPositionsInPreviousReportingRatio | 直近空売り残高割合 | Number | 前回報告時の残高割合 |

##### その他

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| Notes | 備考 | String | |

## 法的制度について

### 空売り残高報告制度
- **法的根拠**: 有価証券の取引等の規制に関する内閣府令
- **報告対象**: 残高割合が0.5％以上の大口空売りポジション
- **報告義務者**: 取引参加者（証券会社等）
- **透明性向上**: 市場の透明性向上と適切な価格形成を目的

### 報告の種類
1. **通常報告**: 定期的な残高報告
2. **変更報告**: 残高に大きな変動があった場合
3. **終了報告**: 0.5％を下回った場合

### 投資一任契約
- **適用範囲**: 投資一任契約に基づく取引も報告対象
- **実質投資者**: 契約の相手方（実質的な投資判断者）も併せて開示

## 使用上の注意事項

1. **報告閾値**
   - 0.5％以上の残高のみが報告対象
   - 0.5％未満の空売りは報告されない

2. **報告タイミング**
   - 法令に基づく報告スケジュール
   - 報告がない日はデータなし

3. **名称の取り扱い**
   - 報告者の名称は報告されたものをそのまま記載
   - 日本語・英語が混在する場合がある

4. **計算日と公表日**
   - 計算日：残高を算定した基準日
   - 公表日：データが公表された日

## サンプルレスポンス

```json
{
    "short_selling_positions": [
        {
            "DisclosedDate": "2024-08-01",
            "CalculatedDate": "2024-07-31",
            "Code": "13660",
            "ShortSellerName": "個人",
            "ShortSellerAddress": "",
            "DiscretionaryInvestmentContractorName": "",
            "DiscretionaryInvestmentContractorAddress": "",
            "InvestmentFundName": "",
            "ShortPositionsToSharesOutstandingRatio": 0.0053,
            "ShortPositionsInSharesNumber": 140000,
            "ShortPositionsInTradingUnitsNumber": 140000,
            "CalculationInPreviousReportingDate": "2024-07-22",
            "ShortPositionsInPreviousReportingRatio": 0.0043,
            "Notes": ""
        }
    ],
    "pagination_key": "value1.value2."
}
```

## 活用例

### 特定銘柄の大口空売り状況を確認
```
/markets/short_selling_positions?code=86970&calculated_date=20240801
```

### 特定日の大口空売り全体状況を確認
```
/markets/short_selling_positions?disclosed_date=2024-08-01
```

### 特定銘柄の大口空売り推移分析
```
/markets/short_selling_positions?code=86970&disclosed_date_from=2024-01-01&disclosed_date_to=2024-12-31
```

### 残高変動の監視
前回報告時との比較：
- `ShortPositionsToSharesOutstandingRatio` vs `ShortPositionsInPreviousReportingRatio`
- 残高の増減トレンドを把握

### 大口投資家の特定
- `ShortSellerName`から主要な空売り実行者を特定
- 投資一任契約の場合は実質的な投資判断者を確認

## 関連情報

このAPIは大口空売り残高の報告データを提供するものです。関連する他のAPIについては以下を参照してください：
- 業種別空売り比率：`/markets/short_selling`（業種別の空売り比率）
- 信用取引週末残高：`/markets/weekly_margin_interest`（信用取引残高）
- 株価四本値：`/prices/daily_quotes`（個別銘柄の株価）
- 売買内訳データ：`/markets/breakdown`（詳細な売買内訳）