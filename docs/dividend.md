# J-Quants API 配当金情報データ項目概要

このドキュメントは、J-Quants APIの配当金情報エンドポイント (`/fins/dividend`) のレスポンスデータ項目について説明します。

## データ型について
- 日付フィールドは `String` 型で返されます（YYYY-MM-DD形式）
- 配当金額データは `Number` または `String` 型で返されます
- 未定の場合は `-`、非設定の場合は空文字が設定されます

## エンドポイント概要

### 配当金情報 (/fins/dividend)

#### 概要
- **エンドポイント**: `https://api.jquants.com/v1/fins/dividend`
- **HTTPメソッド**: GET
- **説明**: 上場会社の配当（決定・予想）に関する１株当たり配当金額、基準日、権利落日及び支払開始予定日等の情報を取得できます。

#### 本APIの留意点
- **対象銘柄**: 東証上場銘柄でない銘柄（地方取引所単独上場銘柄）についてはデータの収録対象外となっています

#### リクエストパラメータ

データの取得では、銘柄コード（code）または通知日付（date）の指定が必須となります。

| パラメータ名 | 型 | 必須 | 説明 | 例 |
|------------|-----|------|------|-----|
| code | String | いいえ* | 銘柄コード | 27800 or 2780 |
| date | String | いいえ* | 通知日付 | 20210907 or 2021-09-07 |
| from | String | いいえ | 期間の開始日 | 20210901 or 2021-09-01 |
| to | String | いいえ | 期間の終了日 | 20210907 or 2021-09-07 |
| pagination_key | String | いいえ | 検索の先頭を指定する文字列 | 過去の検索で返却されたpagination_key |

*codeまたはdateのいずれかが必須

#### 銘柄コード指定時の注意
4桁の銘柄コードを指定した場合は、普通株式と優先株式の両方が上場している銘柄においては普通株式のデータのみが取得されます。

#### パラメータの組み合わせパターン

| code | date | from/to | レスポンスの結果 |
|------|------|---------|-----------------|
| 指定あり | 指定なし | 指定なし | 指定された銘柄について取得可能期間の全データ |
| 指定あり | 指定なし | 指定あり | 指定された銘柄について指定された期間分のデータ |
| 指定なし | 指定あり | - | 全上場銘柄について指定された通知日付のデータ |

#### リクエストヘッダー

| ヘッダー名 | 型 | 必須 | 説明 |
|----------|-----|------|------|
| Authorization | String | はい | アクセスキー |

#### レスポンスフィールド

##### 基本情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| AnnouncementDate | 通知日時（年月日） | String | YYYY-MM-DD形式 |
| AnnouncementTime | 通知日時（時分） | String | HH:MI形式 |
| Code | 銘柄コード | String | |
| ReferenceNumber | リファレンスナンバー | String | 配当通知を一意に特定するための番号 |
| CAReferenceNumber | CAリファレンスナンバー | String | 訂正・削除の対象となっている配当通知のリファレンスナンバー。新規の場合はリファレンスナンバーと同じ値 |

##### 更新・配当区分情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| StatusCode | 更新区分（コード） | String | 1: 新規、2: 訂正、3: 削除 |
| InterimFinalCode | 配当種類（コード） | String | 1: 中間配当、2: 期末配当 |
| ForecastResultCode | 予想／決定（コード） | String | 1: 決定、2: 予想 |
| CommemorativeSpecialCode | 記念配当/特別配当コード | String | 1: 記念配当、2: 特別配当、3: 記念・特別配当、0: 通常の配当 |

##### 日程情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| BoardMeetingDate | 取締役会決議日 | String | YYYY-MM-DD形式 |
| InterimFinalTerm | 配当基準日年月 | String | YYYY-MM形式 |
| RecordDate | 基準日 | String | YYYY-MM-DD形式 |
| ExDate | 権利落日 | String | YYYY-MM-DD形式 |
| ActualRecordDate | 権利確定日 | String | YYYY-MM-DD形式 |
| PayableDate | 支払開始予定日 | String | YYYY-MM-DD形式、未定の場合: `-`、非設定の場合: 空文字 |

##### 配当金額情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| GrossDividendRate | １株当たり配当金額 | Number/String | 未定の場合: `-`、非設定の場合: 空文字 |
| CommemorativeDividendRate | １株当たり記念配当金額 | Number/String | 未定の場合: `-`、非設定の場合: 空文字、2022年6月6日以降のみ提供 |
| SpecialDividendRate | １株当たり特別配当金額 | Number/String | 未定の場合: `-`、非設定の場合: 空文字、2022年6月6日以降のみ提供 |

##### 税務関連情報（2014年2月24日以降のみ提供）

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| DistributionAmount | 1株当たりの交付金銭等の額 | Number/String | 未定の場合: `-`、非設定の場合: 空文字 |
| RetainedEarnings | 1株当たりの利益剰余金の額 | Number/String | 未定の場合: `-`、非設定の場合: 空文字 |
| DeemedDividend | 1株当たりのみなし配当の額 | Number/String | 未定の場合: `-`、非設定の場合: 空文字 |
| DeemedCapitalGains | 1株当たりのみなし譲渡収入の額 | Number/String | 未定の場合: `-`、非設定の場合: 空文字 |
| NetAssetDecreaseRatio | 純資産減少割合 | Number/String | 未定の場合: `-`、非設定の場合: 空文字 |

## 配当の種類について

### 配当種類（InterimFinalCode）
- **1: 中間配当**: 事業年度の中間時点で支払われる配当
- **2: 期末配当**: 事業年度の終了時点で支払われる配当

### 予想／決定（ForecastResultCode）
- **1: 決定**: 正式に決定された配当
- **2: 予想**: 会社が予想している配当額

### 記念配当/特別配当（CommemorativeSpecialCode）
- **0: 通常の配当**: 通常の業績に基づく配当
- **1: 記念配当**: 創立記念等の特別な記念による配当
- **2: 特別配当**: 特別な利益等による配当
- **3: 記念・特別配当**: 記念配当と特別配当の両方

## リファレンスナンバーについて

### リファレンスナンバー（ReferenceNumber）
配当通知を一意に特定するための番号。新規通知、訂正、削除のたびに新しい番号が割り当てられる。

### CAリファレンスナンバー（CAReferenceNumber）
訂正・削除の対象となっている配当通知のリファレンスナンバー。新規の場合はリファレンスナンバーと同じ値を設定。

### 具体例
銘柄：日本取引所（銘柄コード：86970）について以下の通知があった場合：

| 日付 | 通知内容 | ReferenceNumber | CAReferenceNumber | StatusCode |
|------|----------|-----------------|-------------------|------------|
| 2023-03-06 | 配当が新規で通知 | 1 | 1 | 1：新規 |
| 2023-03-07 | 配当が訂正情報として通知 | 2 | 1 | 2：訂正 |
| 2023-03-08 | 配当が削除された | 3 | 1 | 3：削除 |
| 2023-03-09 | 配当が新規で通知 | 4 | 4 | 1：新規 |

## 重要な日程について

### 基準日（RecordDate）
株主としての権利を確定する日。この日の株主名簿に記載されている株主が配当を受ける権利を有する。

### 権利落日（ExDate）
この日以降に株式を購入しても、その配当を受ける権利がない日。通常は基準日の2営業日前。

### 権利確定日（ActualRecordDate）
実際に株主の権利が確定する日。

### 支払開始予定日（PayableDate）
配当金の支払いが開始される予定日。

## 使用上の注意事項

1. **データ提供期間**
   - 税務関連情報は2014年2月24日以降のみ提供
   - 記念・特別配当情報は2022年6月6日以降のみ提供

2. **未定・非設定の表現**
   - 未定の場合: `-`
   - 非設定の場合: 空文字

3. **更新区分の理解**
   - 新規、訂正、削除の履歴がすべて保持される
   - 最新の有効な情報を取得するには更新区分を考慮する必要がある

4. **対象銘柄の制限**
   - 東証上場銘柄のみが対象
   - 地方取引所単独上場銘柄は対象外

## サンプルレスポンス

```json
{
    "dividend": [
        {
            "AnnouncementDate": "2014-02-24",
            "AnnouncementTime": "09:21",
            "Code": "15550",
            "ReferenceNumber": "201402241B00002",
            "StatusCode": "1",
            "BoardMeetingDate": "2014-02-24",
            "InterimFinalCode": "2",
            "ForecastResultCode": "2",
            "InterimFinalTerm": "2014-03",
            "GrossDividendRate": "-",
            "RecordDate": "2014-03-10",
            "ExDate": "2014-03-06",
            "ActualRecordDate": "2014-03-10",
            "PayableDate": "-",
            "CAReferenceNumber": "201402241B00002",
            "DistributionAmount": "",
            "RetainedEarnings": "",
            "DeemedDividend": "",
            "DeemedCapitalGains": "",
            "NetAssetDecreaseRatio": "",
            "CommemorativeSpecialCode": "0",
            "CommemorativeDividendRate": "",
            "SpecialDividendRate": ""
        }
    ],
    "pagination_key": "value1.value2."
}
```

## 活用例

### 特定銘柄の配当履歴を取得
```
/fins/dividend?code=86970&from=20230101&to=20231231
```

### 特定日の全銘柄配当情報を取得
```
/fins/dividend?date=2023-03-06
```

### 配当利回りの計算
```javascript
// 年間配当利回りの計算例
const annualDividend = interimDividend + finalDividend;
const dividendYield = (annualDividend / currentPrice) * 100;
```

### 権利確定スケジュールの確認
権利落日（ExDate）と基準日（RecordDate）から投資タイミングを判断

### 配当成長率の分析
過去の配当実績から配当成長パターンを分析

### 特別配当・記念配当の検出
`CommemorativeSpecialCode`が0以外の銘柄で特別な配当を検出

## 関連情報

このAPIは配当金情報を提供するものです。関連する他のAPIについては以下を参照してください：
- 財務情報：`/fins/statements`（財務諸表データ）
- 財務諸表詳細：`/fins/fs_details`（BS/PL詳細データ）
- 決算発表予定日：`/fins/announcement`（決算発表スケジュール）
- 株価四本値：`/prices/daily_quotes`（株価データ）