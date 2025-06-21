# J-Quants API 業種別空売り比率データ項目概要

このドキュメントは、J-Quants APIの業種別空売り比率エンドポイント (`/markets/short_selling`) のレスポンスデータ項目について説明します。

## データ型について
- 日付フィールドは `String` 型で返されます（YYYY-MM-DD形式）
- 売買代金データは `Number` 型で返されます
- 金額の単位は円です

## エンドポイント概要

### 業種別空売り比率 (/markets/short_selling)

#### 概要
- **エンドポイント**: `https://api.jquants.com/v1/markets/short_selling`
- **HTTPメソッド**: GET
- **説明**: 日々の業種（セクター）別の空売り比率に関する売買代金を取得することができます。

#### 公式情報源
配信データは下記のページで公表している内容と同様です：
- https://www.jpx.co.jp/markets/statistics-equities/short-selling/index.html

**注意**: Webページでの公表値は百万円単位に丸められていますが、APIでは円単位のデータとなります。

#### 本APIの留意点
- **取引のない日**: 取引高が存在しない（売買されていない）日の日付（date）を指定した場合は、値は空です
- **システム障害日**: 2020/10/1は東京証券取引所の株式売買システムの障害により終日売買停止となった関係で、データが存在しません

#### リクエストパラメータ

データの取得では、33業種コード（sector33code）または日付（date）の指定が必須となります。

| パラメータ名 | 型 | 必須 | 説明 | 例 |
|------------|-----|------|------|-----|
| sector33code | String | いいえ* | 33業種コード | 0050 or 50 |
| date | String | いいえ* | 日付 | 20210907 or 2021-09-07 |
| from | String | いいえ | 期間の開始日 | 20210901 or 2021-09-01 |
| to | String | いいえ | 期間の終了日 | 20210907 or 2021-09-07 |
| pagination_key | String | いいえ | 検索の先頭を指定する文字列 | 過去の検索で返却されたpagination_key |

*sector33codeまたはdateのいずれかが必須

#### パラメータの組み合わせパターン

| sector33code | date | from/to | レスポンスの結果 |
|--------------|------|---------|-----------------|
| 指定なし | 指定あり | - | 全業種コードについて指定された日付のデータ |
| 指定あり | 指定なし | 指定なし | 指定された業種コードについて全期間分のデータ |
| 指定あり | 指定なし | 指定あり | 指定された業種コードについて指定された期間分のデータ |
| 指定あり | 指定あり | - | 指定された業種コードについて指定された日付のデータ |

#### リクエストヘッダー

| ヘッダー名 | 型 | 必須 | 説明 |
|----------|-----|------|------|
| Authorization | String | はい | アクセスキー |

#### レスポンスフィールド

##### 基本情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| Date | 日付 | String | YYYY-MM-DD形式 |
| Sector33Code | 33業種コード | String | 33業種コード及び業種名参照 |

##### 売買代金データ

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| SellingExcludingShortSellingTurnoverValue | 実注文の売買代金 | Number | 単位：円、空売り以外の通常売り注文 |
| ShortSellingWithRestrictionsTurnoverValue | 価格規制有りの空売り売買代金 | Number | 単位：円、アップティック・ルール等の価格規制下での空売り |
| ShortSellingWithoutRestrictionsTurnoverValue | 価格規制無しの空売り売買代金 | Number | 単位：円、価格規制の適用を受けない空売り |

## 空売り規制について

### 価格規制有りの空売り
- **アップティック・ルール**: 直前約定価格より高い価格での空売り
- **適用対象**: 通常の個別銘柄空売りに適用
- **目的**: 株価の急落を防ぐための規制

### 価格規制無しの空売り
- **適用対象**: ETF、REIT、証券会社による裁定取引等
- **特徴**: アップティック・ルールの適用を受けない
- **理由**: 市場の流動性確保、裁定機能維持

## 33業種分類について

このAPIで使用される33業種コードは、東京証券取引所の33業種分類に基づいています。
詳細な業種コードとコード対応表については、上場銘柄一覧APIの33業種コード参照資料をご確認ください。

主要業種例：
- 0050: 水産・農林業
- 1050: 鉱業
- 2050: 建設業
- 3050: 食料品
- 3100: 繊維製品
- 9999: その他

## 使用上の注意事項

1. **データ精度**
   - API提供データは円単位
   - 公式Webサイトは百万円単位で丸め表示

2. **取引のない日**
   - 祝日、年末年始等は通常データなし
   - システム障害日（2020/10/1）等もデータなし

3. **業種分類の特性**
   - 33業種分類による集計
   - 業種変更時は変更後分類で履歴も再整理

4. **空売り比率の計算**
   - 比率 = 空売り売買代金 ÷ 総売買代金 × 100

## サンプルレスポンス

```json
{
    "short_selling": [
        {
            "Date": "2022-10-25",
            "Sector33Code": "0050",
            "SellingExcludingShortSellingTurnoverValue": 1333126400.0,
            "ShortSellingWithRestrictionsTurnoverValue": 787355200.0,
            "ShortSellingWithoutRestrictionsTurnoverValue": 149084300.0
        }
    ],
    "pagination_key": "value1.value2."
}
```

## 活用例

### 特定業種の空売り比率推移を取得
```
/markets/short_selling?sector33code=0050&from=20220101&to=20221231
```

### 特定日の全業種空売り状況を取得
```
/markets/short_selling?date=2022-10-25
```

### 空売り比率の計算
```javascript
// 空売り比率の計算例
const totalShortSelling = response.ShortSellingWithRestrictionsTurnoverValue + 
                         response.ShortSellingWithoutRestrictionsTurnoverValue;
const totalTurnover = response.SellingExcludingShortSellingTurnoverValue + totalShortSelling;
const shortSellingRatio = (totalShortSelling / totalTurnover) * 100;
```

### 価格規制有無別の空売り分析
規制有り空売り（`ShortSellingWithRestrictionsTurnoverValue`）と規制無し空売り（`ShortSellingWithoutRestrictionsTurnoverValue`）の比率分析

### 業種間空売り比率比較
複数業種の同日データを取得し、業種間の空売り動向を比較分析

## 関連情報

このAPIは業種別の空売り比率データを提供するものです。関連する他のAPIについては以下を参照してください：
- 空売り残高報告：`/markets/short_selling_positions`（大口の空売り残高）
- 信用取引週末残高：`/markets/weekly_margin_interest`（信用取引残高）
- 売買内訳データ：`/markets/breakdown`（詳細な売買内訳）
- 上場銘柄一覧：`/listed/info`（33業種コード一覧）