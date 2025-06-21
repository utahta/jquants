# J-Quants API 日経225オプション四本値データ項目概要

このドキュメントは、J-Quants APIの日経225オプション四本値エンドポイント (`/option/index_option`) のレスポンスデータ項目について説明します。

## データ型について
- 日付フィールドは `String` 型で返されます（YYYY-MM-DD形式）
- 価格・数量データは `Number` 型で返されます
- 一部フィールドは条件により `String` 型（空文字）になる場合があります

## エンドポイント概要

### 日経225オプション四本値 (/option/index_option)

#### 概要
- **エンドポイント**: `https://api.jquants.com/v1/option/index_option`
- **HTTPメソッド**: GET
- **説明**: 日経225オプションに関する、四本値や清算値段、理論価格に関する情報を取得することができます。

#### 対象商品
- 日経225指数オプション（Weeklyオプション及びフレックスオプションを除く）のみ

#### 本APIの留意点

##### 取引セッションについて
- **2011年2月10日以前**：ナイトセッション、前場、後場で構成
  - 前場データは収録されず、後場データが日中場データとして収録
  - 日通しデータについては、全立会を含めたデータ
- **2011年2月14日以降**：ナイトセッション、日中場で構成

##### レスポンスのキー項目について
- 緊急取引証拠金が発動した場合は、同一の取引日・銘柄に対して清算価格算出時と緊急取引証拠金算出時のデータが発生
- Date、Codeに加えて`EmergencyMarginTriggerDivision`を組み合わせることでデータを一意に識別可能

#### リクエストパラメータ

| パラメータ名 | 型 | 必須 | 説明 | 例 |
|------------|-----|------|------|-----|
| date | String | はい | 取引日 | 20210901 or 2021-09-01 |
| pagination_key | String | いいえ | 検索の先頭を指定する文字列 | 過去の検索で返却されたpagination_key |

#### リクエストヘッダー

| ヘッダー名 | 型 | 必須 | 説明 |
|----------|-----|------|------|
| Authorization | String | はい | アクセスキー |

#### レスポンスフィールド

##### 基本情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| Date | 取引日 | String | YYYY-MM-DD形式 |
| Code | 銘柄コード | String | |
| ContractMonth | 限月 | String | YYYY-MM形式 |
| StrikePrice | 権利行使価格 | Number | |
| PutCallDivision | プットコール区分 | String | 1: プット、2: コール |
| LastTradingDay | 取引最終年月日 | String | YYYY-MM-DD形式 ※1 |
| SpecialQuotationDay | SQ日 | String | YYYY-MM-DD形式 ※1 |
| EmergencyMarginTriggerDivision | 緊急取引証拠金発動区分 | String | 001: 緊急取引証拠金発動時、002: 清算価格算出時 |

##### 四本値データ

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| WholeDayOpen | 日通し始値 | Number | |
| WholeDayHigh | 日通し高値 | Number | |
| WholeDayLow | 日通し安値 | Number | |
| WholeDayClose | 日通し終値 | Number | |
| NightSessionOpen | ナイト・セッション始値 | Number/String | 取引開始日初日は空文字 |
| NightSessionHigh | ナイト・セッション高値 | Number/String | 取引開始日初日は空文字 |
| NightSessionLow | ナイト・セッション安値 | Number/String | 取引開始日初日は空文字 |
| NightSessionClose | ナイト・セッション終値 | Number/String | 取引開始日初日は空文字 |
| DaySessionOpen | 日中始値 | Number | |
| DaySessionHigh | 日中高値 | Number | |
| DaySessionLow | 日中安値 | Number | |
| DaySessionClose | 日中終値 | Number | |

##### 取引情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| Volume | 取引高 | Number | |
| Volume(OnlyAuction) | 立会内取引高 | Number | ※1 |
| OpenInterest | 建玉 | Number | |
| TurnoverValue | 取引代金 | Number | |

##### 価格・リスク情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| SettlementPrice | 清算値段 | Number | ※1 |
| TheoreticalPrice | 理論価格 | Number | ※1 |
| BaseVolatility | 基準ボラティリティ | Number | アット・ザ・マネープット及びコールそれぞれのインプライドボラティリティの中間値 ※1 |
| UnderlyingPrice | 原証券価格 | Number | ※1 |
| ImpliedVolatility | インプライドボラティリティ | Number | ※1 |
| InterestRate | 理論価格計算用金利 | Number | ※1 |

※1 2016年7月19日以降のみ提供

## 緊急取引証拠金発動区分について

| 区分値 | 説明 | 備考 |
|-------|------|------|
| 001 | 緊急取引証拠金発動時 | 2016年7月19日以降に緊急取引証拠金発動した場合のみ収録 |
| 002 | 清算価格算出時 | 通常時のデータ |

## 使用上の注意事項

1. **日付指定が必須**
   - 他のAPIと異なり、日付（date）パラメータが必須
   - 期間指定（from/to）には対応していない

2. **データの一意性**
   - Date + Code + EmergencyMarginTriggerDivisionの組み合わせで一意
   - 緊急取引証拠金発動時は同じ日・銘柄で複数レコードが存在

3. **取引セッションの変更**
   - 2011年2月を境に取引セッション構成が変更
   - ナイトセッションデータの有無に注意

4. **データ提供期間の制限**
   - 一部フィールドは2016年7月19日以降のみ提供
   - 取引開始日初日はナイトセッションデータなし

## サンプルレスポンス

```json
{
    "index_option": [
        {
            "Date": "2023-03-22",
            "Code": "130060018",
            "WholeDayOpen": 0.0,
            "WholeDayHigh": 0.0,
            "WholeDayLow": 0.0,
            "WholeDayClose": 0.0,
            "NightSessionOpen": 0.0,
            "NightSessionHigh": 0.0,
            "NightSessionLow": 0.0,
            "NightSessionClose": 0.0,
            "DaySessionOpen": 0.0,
            "DaySessionHigh": 0.0,
            "DaySessionLow": 0.0,
            "DaySessionClose": 0.0,
            "Volume": 0.0,
            "OpenInterest": 330.0,
            "TurnoverValue": 0.0,
            "ContractMonth": "2025-06",
            "StrikePrice": 20000.0,
            "Volume(OnlyAuction)": 0.0,
            "EmergencyMarginTriggerDivision": "002",
            "PutCallDivision": "1",
            "LastTradingDay": "2025-06-12",
            "SpecialQuotationDay": "2025-06-13",
            "SettlementPrice": 980.0,
            "TheoreticalPrice": 974.641,
            "BaseVolatility": 17.93025,
            "UnderlyingPrice": 27466.61,
            "ImpliedVolatility": 23.1816,
            "InterestRate": 0.2336
        }
    ],
    "pagination_key": "value1.value2."
}
```

## 活用例

### 特定日のオプション全銘柄データ取得
```
/option/index_option?date=20230324
```

### プットオプションのみフィルタリング
レスポンスの`PutCallDivision`が"1"のデータを抽出

### コールオプションのみフィルタリング
レスポンスの`PutCallDivision`が"2"のデータを抽出

### 緊急取引証拠金発動日の確認
`EmergencyMarginTriggerDivision`が"001"のレコードの有無を確認

## 関連情報

このAPIは日経225オプションの詳細データを提供するものです。関連する他のAPIについては以下を参照してください：
- 先物四本値：`/derivatives/futures`（日経225先物等）
- オプション四本値：`/derivatives/options`（汎用オプションAPI）
- 指数四本値：`/indices`（原証券である日経225指数等）
- 取引カレンダー：`/markets/trading_calendar`（営業日情報）