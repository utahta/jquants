# J-Quants API 先物四本値データ項目概要

このドキュメントは、J-Quants APIの先物四本値エンドポイント (`/derivatives/futures`) のレスポンスデータ項目について説明します。

## データ型について
- 日付フィールドは `String` 型で返されます（YYYY-MM-DD形式）
- 価格・数量データは `Number` 型で返されます
- 一部フィールドは条件により `String` 型（空文字）になる場合があります

## エンドポイント概要

### 先物四本値 (/derivatives/futures)

#### 概要
- **エンドポイント**: `https://api.jquants.com/v1/derivatives/futures`
- **HTTPメソッド**: GET
- **説明**: 先物に関する、四本値や清算値段、理論価格に関する情報を取得することができます。

#### 対象商品
本APIで取得可能なデータについては先物商品区分コード一覧を参照してください。

#### 本APIの留意点

##### 銘柄コードについて
- 先物・オプション取引識別コードの付番規則については[証券コード関係の関係資料等](https://www.jpx.co.jp/sicc/securities-code/01.html)を参照

##### 取引セッションについて
- **2011年2月10日以前**：ナイトセッション、前場、後場で構成
  - 前場データは収録されず、後場データが日中場データとして収録
  - 日通しデータについては、全立会を含めたデータ
- **2011年2月14日以降**：ナイトセッション、日中場で構成

##### 祝日取引について
- 祝日取引の取引日は、祝日取引実施日直前の平日に開始するナイト・セッション（祝日前営業日）及び祝日取引実施日直後の平日（祝日翌営業日）のデイ・セッションと同一の取引日として扱われる

##### レスポンスのキー項目について
- 緊急取引証拠金が発動した場合は、同一の取引日・銘柄に対して清算価格算出時と緊急取引証拠金算出時のデータが発生
- Date、Codeに加えてEmergencyMarginTriggerDivisionを組み合わせることでデータを一意に識別可能

#### リクエストパラメータ

データの取得では、日付（date）の指定が必須となります。

| パラメータ名 | 型 | 必須 | 説明 | 例 |
|------------|-----|------|------|-----|
| date | String | はい | 取引日 | 20210901 or 2021-09-01 |
| category | String | いいえ | 商品区分の指定 | TOPIXF |
| contract_flag | String | いいえ | 中心限月フラグの指定 | 1 |
| pagination_key | String | いいえ | 検索の先頭を指定する文字列 | 過去の検索で返却されたpagination_key |

#### リクエストヘッダー

| ヘッダー名 | 型 | 必須 | 説明 |
|----------|-----|------|------|
| Authorization | String | はい | アクセスキー |

#### レスポンスフィールド

##### 基本情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| Code | 銘柄コード | String | |
| DerivativesProductCategory | 先物商品区分 | String | |
| Date | 取引日 | String | YYYY-MM-DD形式 |
| ContractMonth | 限月 | String | YYYY-MM形式 |
| EmergencyMarginTriggerDivision | 緊急取引証拠金発動区分 | String | 001: 緊急取引証拠金発動時、002: 清算価格算出時 ※1 |

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
| MorningSessionOpen | 前場始値 | Number/String | 前後場取引対象銘柄でない場合、空文字 |
| MorningSessionHigh | 前場高値 | Number/String | 前後場取引対象銘柄でない場合、空文字 |
| MorningSessionLow | 前場安値 | Number/String | 前後場取引対象銘柄でない場合、空文字 |
| MorningSessionClose | 前場終値 | Number/String | 前後場取引対象銘柄でない場合、空文字 |

##### 取引情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| Volume | 取引高 | Number | |
| Volume(OnlyAuction) | 立会内取引高 | Number | ※1 |
| OpenInterest | 建玉 | Number | |
| TurnoverValue | 取引代金 | Number | |
| SettlementPrice | 清算値段 | Number | ※1 |

##### その他情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| LastTradingDay | 取引最終年月日 | String | YYYY-MM-DD形式 ※1 |
| SpecialQuotationDay | SQ日 | String | YYYY-MM-DD形式 ※1 |
| CentralContractMonthFlag | 中心限月フラグ | String | 1:中心限月、0:その他 ※1 |

※1 2016年7月19日以降のみ提供

## 先物商品区分コード一覧

| コード | 商品区分名称 | データ収録期間 |
|--------|------------|--------------|
| TOPIXF | TOPIX先物 | 2008/5/7〜 |
| TOPIXMF | ミニTOPIX先物 | 2008/6/16〜 |
| MOTF | マザーズ先物 | 2016/7/19〜 |
| NKVIF | 日経平均VI先物 | 2012/2/27〜 |
| NKYDF | 日経平均・配当指数先物 | 2010/7/26〜 |
| NK225F | 日経225先物 | 2008/5/7〜 |
| NK225MF | 日経225mini先物 | 2008/5/7〜 |
| JN400F | JPX日経インデックス400先物 | 2014/11/25〜 |
| REITF | 東証REIT指数先物 | 2008/6/16〜 |
| DJIAF | NYダウ先物 | 2012/5/28〜 |
| JGBLF | 長期国債先物 | 2008/5/7〜 |
| NK225MCF | 日経225マイクロ先物 | 2023/5/29〜 |
| TOA3MF | TONA3ヶ月金利先物 | 2023/5/29〜 |

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

5. **祝日取引の扱い**
   - 祝日取引は前営業日のナイトセッションと翌営業日のデイセッションと同一取引日

## サンプルレスポンス

```json
{
    "futures": [
        {
            "Code": "169090005",
            "DerivativesProductCategory": "TOPIXF",
            "Date": "2024-07-23", 
            "WholeDayOpen": 2825.5, 
            "WholeDayHigh": 2853.0, 
            "WholeDayLow": 2825.5, 
            "WholeDayClose": 2829.0, 
            "MorningSessionOpen": "", 
            "MorningSessionHigh": "", 
            "MorningSessionLow": "", 
            "MorningSessionClose": "", 
            "NightSessionOpen": 2825.5, 
            "NightSessionHigh": 2850.0, 
            "NightSessionLow": 2825.5, 
            "NightSessionClose": 2845.0, 
            "DaySessionOpen": 2850.5, 
            "DaySessionHigh": 2853.0, 
            "DaySessionLow": 2826.0, 
            "DaySessionClose": 2829.0, 
            "Volume": 42910.0, 
            "OpenInterest": 479812.0, 
            "TurnoverValue": 1217918971856.0, 
            "ContractMonth": "2024-09", 
            "Volume(OnlyAuction)": 40405.0, 
            "EmergencyMarginTriggerDivision": "002", 
            "LastTradingDay": "2024-09-12", 
            "SpecialQuotationDay": "2024-09-13", 
            "SettlementPrice": 2829.0, 
            "CentralContractMonthFlag": "1"
        }
    ],
    "pagination_key": "value1.value2."
}
```

## 活用例

### 特定日の全先物データ取得
```
/derivatives/futures?date=20230324
```

### 特定商品カテゴリのデータ取得
```
/derivatives/futures?date=20230324&category=NK225F
```

### 中心限月のみフィルタリング
```
/derivatives/futures?date=20230324&contract_flag=1
```

### 日経225先物とミニの比較
NK225FとNK225MFのデータを取得し、価格差やボリューム比率を分析

### ナイトセッションと日中セッションの価格変動分析
```javascript
// ギャップ（窓）の計算
const nightToDay = response.DaySessionOpen - response.NightSessionClose;
const dayToNight = response.NightSessionOpen - response.DaySessionClose;
```

### 緊急取引証拠金発動日の確認
`EmergencyMarginTriggerDivision`が"001"のレコードの有無を確認

## 関連情報

このAPIは先物の詳細データを提供するものです。関連する他のAPIについては以下を参照してください：
- 日経225オプション四本値：`/option/index_option`（日経225オプション専用）
- オプション四本値：`/derivatives/options`（汎用オプションAPI）
- 指数四本値：`/indices`（原資産となる指数データ）
- 取引カレンダー：`/markets/trading_calendar`（営業日・祝日取引日情報）