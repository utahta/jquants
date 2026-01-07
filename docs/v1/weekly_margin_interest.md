# J-Quants API 信用取引週末残高データ項目概要

このドキュメントは、J-Quants APIの信用取引週末残高エンドポイント (`/markets/weekly_margin_interest`) のレスポンスデータ項目について説明します。

## データ型について
- 日付フィールドは `String` 型で返されます（YYYY-MM-DD形式）
- 残高数量データは `Number` 型で返されます
- 残高の単位は株数です

## エンドポイント概要

### 信用取引週末残高 (/markets/weekly_margin_interest)

#### 概要
- **エンドポイント**: `https://api.jquants.com/v1/markets/weekly_margin_interest`
- **HTTPメソッド**: GET
- **説明**: 週末時点での、各銘柄についての信用取引残高（株数）を取得することができます。

#### 公式情報源
配信データは下記のページで公表している内容と同一です：
- https://www.jpx.co.jp/markets/statistics-equities/margin/05.html

#### 本APIの留意点
- **株数調整なし**: 当該銘柄のコーポレートアクションが発生した場合も、遡及して株数の調整は行われません
- **営業日制限**: 年末年始など、営業日が2日以下の週はデータが提供されません
- **対象銘柄**: 東証上場銘柄でない銘柄（地方取引所単独上場銘柄）についてはデータの収録対象外となっています

#### リクエストパラメータ

データの取得では、銘柄コード（code）または日付（date）の指定が必須となります。

| パラメータ名 | 型 | 必須 | 説明 | 例 |
|------------|-----|------|------|-----|
| code | String | いいえ* | 銘柄コード | 27800 or 2780 |
| date | String | いいえ* | 申込日付 | 20210907 or 2021-09-07 |
| from | String | いいえ | 期間の開始日 | 20210901 or 2021-09-01 |
| to | String | いいえ | 期間の終了日 | 20210907 or 2021-09-07 |
| pagination_key | String | いいえ | 検索の先頭を指定する文字列 | 過去の検索で返却されたpagination_key |

*codeまたはdateのいずれかが必須

#### パラメータの組み合わせパターン

| code | date | from/to | レスポンスの結果 |
|------|------|---------|-----------------|
| 指定あり | - | 指定なし | 指定された銘柄について全期間分のデータ |
| 指定あり | - | 指定あり | 指定された銘柄について指定された期間分のデータ |
| 指定なし | 指定あり | - | 全上場銘柄について指定された日付のデータ |

#### 銘柄コード指定時の注意
4桁の銘柄コードを指定した場合は、普通株式と優先株式の両方が上場している銘柄においては普通株式のデータのみが取得されます。

#### リクエストヘッダー

| ヘッダー名 | 型 | 必須 | 説明 |
|----------|-----|------|------|
| Authorization | String | はい | アクセスキー |

#### レスポンスフィールド

##### 基本情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| Date | 申込日付 | String | YYYY-MM-DD形式、信用取引残高基準となる時点（通常は金曜日付） |
| Code | 銘柄コード | String | |
| IssueType | 銘柄区分 | String | 1: 信用銘柄、2: 貸借銘柄、3: その他 |

##### 信用取引残高（売建）

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| ShortMarginTradeVolume | 売合計信用取引週末残高 | Number | 売建の総残高 |
| ShortNegotiableMarginTradeVolume | 売一般信用取引週末残高 | Number | 売合計のうち一般信用によるもの |
| ShortStandardizedMarginTradeVolume | 売制度信用取引週末残高 | Number | 売合計のうち制度信用によるもの |

##### 信用取引残高（買建）

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| LongMarginTradeVolume | 買合計信用取引週末残高 | Number | 買建の総残高 |
| LongNegotiableMarginTradeVolume | 買一般信用取引週末残高 | Number | 買合計のうち一般信用によるもの |
| LongStandardizedMarginTradeVolume | 買制度信用取引週末残高 | Number | 買合計のうち制度信用によるもの |

## 銘柄区分（IssueType）について

| 区分値 | 説明 | 備考 |
|-------|------|------|
| 1 | 信用銘柄 | 制度信用取引のみ可能 |
| 2 | 貸借銘柄 | 制度信用取引・貸借取引ともに可能 |
| 3 | その他 | 一般信用取引のみまたは取引不可 |

## 使用上の注意事項

1. **データ基準日**
   - 週末（通常は金曜日）時点の残高データ
   - 年末年始等で営業日が2日以下の週はデータなし

2. **株数調整**
   - コーポレートアクション発生時も遡及調整は行われない
   - 株式分割等の際は注意が必要

3. **対象銘柄の制限**
   - 東証上場銘柄のみが対象
   - 地方取引所単独上場銘柄は対象外

4. **信用取引の種類**
   - 制度信用取引：証券金融会社が資金・株券を融通
   - 一般信用取引：証券会社が独自に資金・株券を融通

## サンプルレスポンス

```json
{
    "weekly_margin_interest": [
        {
            "Date": "2023-02-17",
            "Code": "13010",
            "ShortMarginTradeVolume": 4100.0,
            "LongMarginTradeVolume": 27600.0,
            "ShortNegotiableMarginTradeVolume": 1300.0,
            "LongNegotiableMarginTradeVolume": 7600.0,
            "ShortStandardizedMarginTradeVolume": 2800.0,
            "LongStandardizedMarginTradeVolume": 20000.0,
            "IssueType": "2"
        }
    ],
    "pagination_key": "value1.value2."
}
```

## 活用例

### 特定銘柄の信用残高推移を取得
```
/markets/weekly_margin_interest?code=86970&from=20230101&to=20230331
```

### 特定日の全銘柄信用残高を取得
```
/markets/weekly_margin_interest?date=20230217
```

### 制度信用と一般信用の残高比較
レスポンスの制度信用残高フィールド（`*StandardizedMarginTradeVolume`）と一般信用残高フィールド（`*NegotiableMarginTradeVolume`）を比較分析

### 売建・買建バランス分析
売建残高（`Short*`）と買建残高（`Long*`）のバランスから市場心理を分析

## 関連情報

このAPIは信用取引の週末残高データを提供するものです。関連する他のAPIについては以下を参照してください：
- 株価四本値：`/prices/daily_quotes`（個別銘柄の株価）
- 業種別空売り比率：`/markets/short_selling`（空売りデータ）
- 売買内訳データ：`/markets/breakdown`（詳細な売買内訳）
- 取引カレンダー：`/markets/trading_calendar`（営業日情報）