# J-Quants API オプション四本値データ項目概要

このドキュメントは、J-Quants APIのオプション四本値エンドポイント (`/derivatives/options`) のレスポンスデータ項目について説明します。

## データ型について
- 日付フィールドは `String` 型で返されます（YYYY-MM-DD形式）
- 価格・数量データは `Number` 型で返されます
- 一部フィールドは条件により `String` 型（空文字）になる場合があります

## エンドポイント概要

### オプション四本値 (/derivatives/options)

#### 概要
- **エンドポイント**: `https://api.jquants.com/v1/derivatives/options`
- **HTTPメソッド**: GET
- **説明**: オプションに関する、四本値や清算値段、理論価格に関する情報を取得することができます。

#### 対象商品
本APIで取得可能なデータについてはオプション商品区分コード一覧を参照してください。

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
| category | String | いいえ | 商品区分の指定 | TOPIXE |
| code | String | いいえ | 対象有価証券コード（categoryで有価証券オプションを指定した場合に設定） | 7203 |
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
| DerivativesProductCategory | オプション商品区分 | String | |
| UnderlyingSSO | 有価証券オプション対象銘柄 | String | 有価証券オプション以外の場合は"-"を設定 |
| Date | 取引日 | String | YYYY-MM-DD形式 |
| ContractMonth | 限月 | String | YYYY-MM形式、日経225miniオプションの場合は週表記（e.g. 2024-51: 2024年の51週目） |
| StrikePrice | 権利行使価格 | Number | |
| PutCallDivision | プットコール区分 | String | 1: プット、2: コール |
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

##### 価格・リスク情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| SettlementPrice | 清算値段 | Number | ※1 |
| TheoreticalPrice | 理論価格 | Number | ※1 |
| BaseVolatility | 基準ボラティリティ | Number | ※1 |
| UnderlyingPrice | 原証券価格 | Number | ※1 |
| ImpliedVolatility | インプライドボラティリティ | Number | ※1 |
| InterestRate | 理論価格計算用金利 | Number | ※1 |

##### その他情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| LastTradingDay | 取引最終年月日 | String | YYYY-MM-DD形式 ※1 |
| SpecialQuotationDay | SQ日 | String | YYYY-MM-DD形式 ※1 |
| CentralContractMonthFlag | 中心限月フラグ | String | 1:中心限月、0:その他 ※1 |

※1 2016年7月19日以降のみ提供

## オプション商品区分コード一覧

| 商品区分コード | 商品区分名称 | データ収録期間 |
|--------------|------------|--------------|
| TOPIXE | TOPIXオプション | 2008/5/7〜 |
| NK225E | 日経225オプション | 2008/5/7〜 |
| JGBLFE | 長期国債先物オプション | 2008/5/7〜 |
| EQOP | 有価証券オプション | 2014/11/17〜 |
| NK225MWE | 日経225miniオプション | 2023/5/29〜 |

## プットコール区分について

| 区分値 | 説明 | 内容 |
|-------|------|------|
| 1 | プット | 売る権利（原資産価格が下落した場合に利益） |
| 2 | コール | 買う権利（原資産価格が上昇した場合に利益） |

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

5. **限月表記の違い**
   - 通常オプション：YYYY-MM形式（月表記）
   - 日経225miniオプション：YYYY-WW形式（週表記）

6. **有価証券オプション**
   - categoryでEQOPを指定した場合は、codeパラメータで対象銘柄を指定

## サンプルレスポンス

```json
{
    "options": [
        {
            "Code": "140014505", 
            "DerivativesProductCategory": "TOPIXE", 
            "UnderlyingSSO": "-", 
            "Date": "2024-07-23", 
            "WholeDayOpen": 0.0, 
            "WholeDayHigh": 0.0, 
            "WholeDayLow": 0.0, 
            "WholeDayClose": 0.0, 
            "MorningSessionOpen": "", 
            "MorningSessionHigh": "", 
            "MorningSessionLow": "", 
            "MorningSessionClose": "", 
            "NightSessionOpen": 0.0, 
            "NightSessionHigh": 0.0, 
            "NightSessionLow": 0.0, 
            "NightSessionClose": 0.0, 
            "DaySessionOpen": 0.0, 
            "DaySessionHigh": 0.0, 
            "DaySessionLow": 0.0, 
            "DaySessionClose": 0.0, 
            "Volume": 0.0, 
            "OpenInterest": 0.0, 
            "TurnoverValue": 0.0, 
            "ContractMonth": "2025-01", 
            "StrikePrice": 2450.0, 
            "Volume(OnlyAuction)": 0.0, 
            "EmergencyMarginTriggerDivision": "002", 
            "PutCallDivision": "2", 
            "LastTradingDay": "2025-01-09", 
            "SpecialQuotationDay": "2025-01-10", 
            "SettlementPrice": 377.0, 
            "TheoreticalPrice": 380.3801, 
            "BaseVolatility": 18.115, 
            "UnderlyingPrice": 2833.39, 
            "ImpliedVolatility": 17.2955, 
            "InterestRate": 0.3527, 
            "CentralContractMonthFlag": "0"
        }
    ],
    "pagination_key": "value1.value2."
}
```

## 活用例

### 特定日の全オプションデータ取得
```
/derivatives/options?date=20230324
```

### 特定商品カテゴリのデータ取得
```
/derivatives/options?date=20230324&category=NK225E
```

### プットオプションのみフィルタリング
レスポンスの`PutCallDivision`が"1"のデータを抽出

### コールオプションのみフィルタリング
レスポンスの`PutCallDivision`が"2"のデータを抽出

### 有価証券オプションの取得
```
/derivatives/options?date=20230324&category=EQOP&code=7203
```

### ボラティリティ分析
```javascript
// ボラティリティスマイルの確認
const sameExpiry = response.options.filter(opt => 
    opt.ContractMonth === "2025-01"
);
// 権利行使価格別のインプライドボラティリティを比較
```

### イン・ザ・マネー（ITM）判定
```javascript
// コールオプションのITM判定
const isCallITM = opt.PutCallDivision === "2" && 
                  opt.UnderlyingPrice > opt.StrikePrice;
// プットオプションのITM判定
const isPutITM = opt.PutCallDivision === "1" && 
                 opt.UnderlyingPrice < opt.StrikePrice;
```

### 緊急取引証拠金発動日の確認
`EmergencyMarginTriggerDivision`が"001"のレコードの有無を確認

## 関連情報

このAPIはオプションの詳細データを提供するものです。関連する他のAPIについては以下を参照してください：
- 日経225オプション四本値：`/option/index_option`（日経225オプション専用）
- 先物四本値：`/derivatives/futures`（先物データ）
- 指数四本値：`/indices`（原資産となる指数データ）
- 取引カレンダー：`/markets/trading_calendar`（営業日・祝日取引日情報）