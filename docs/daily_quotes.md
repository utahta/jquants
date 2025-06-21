# J-Quants API 株価四本値データ項目概要

このドキュメントは、J-Quants APIの株価四本値エンドポイント (`/prices/daily_quotes`) のレスポンスデータ項目について説明します。

## データ型について
- フィールドは `String` 型または `Number` 型で返されます
- 空文字列 (`""`) または `null` は値が存在しないことを示します
- 金額の単位は円です

## エンドポイント概要

### 株価四本値 (/prices/daily_quotes)

#### 概要
- **エンドポイント**: `https://api.jquants.com/v1/prices/daily_quotes`
- **HTTPメソッド**: GET
- **説明**: 株価情報を取得することができます。株価は分割・併合を考慮した調整済み株価（小数点第２位四捨五入）と調整前の株価を取得することができます。

#### 本APIの留意点
- 取引高が存在しない（売買されていない）日の銘柄についての四本値、取引高と売買代金は、Nullが収録されています
- 東証上場銘柄でない銘柄（地方取引所単独上場銘柄）についてはデータの収録対象外となっております
- 2020/10/1のデータは東京証券取引所の株式売買システムの障害により終日売買停止となった関係で、四本値、取引高と売買代金はNullが収録されています
- 日通しデータについては全プランで取得できますが、前場/後場別のデータについてはPremiumプランのみ取得可能です
- 株価調整については株式分割・併合にのみ対応しております。一部コーポレートアクションには対応しておりませんので、ご了承ください

#### リクエストパラメータ

データの取得では、銘柄コード（code）または日付（date）の指定が必須となります。

| パラメータ名 | 型 | 必須 | 説明 | 例 |
|------------|-----|------|------|-----|
| code | String | いいえ* | 銘柄コード | 27800 or 2780 |
| date | String | いいえ* | 基準となる日付 | 20210907 or 2021-09-07 |
| from | String | いいえ | 期間の開始日 | 20210901 or 2021-09-01 |
| to | String | いいえ | 期間の終了日 | 20210907 or 2021-09-07 |
| pagination_key | String | いいえ | 検索の先頭を指定する文字列 | 過去の検索で返却されたpagination_key |

\* codeまたはdateのいずれかは必須

**注意事項**:
- 4桁の銘柄コードを指定した場合は、普通株式と優先株式等の両方が上場している銘柄においては普通株式のデータのみが取得されます

#### レスポンスフィールド

##### 日通しデータ（全プラン共通）

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| Date | 日付 | String | YYYY-MM-DD形式 |
| Code | 銘柄コード | String | |
| Open | 始値（調整前） | Number | |
| High | 高値（調整前） | Number | |
| Low | 安値（調整前） | Number | |
| Close | 終値（調整前） | Number | |
| UpperLimit | 日通ストップ高フラグ | String | 0：ストップ高以外、1：ストップ高 |
| LowerLimit | 日通ストップ安フラグ | String | 0：ストップ安以外、1：ストップ安 |
| Volume | 取引高（調整前） | Number | |
| TurnoverValue | 取引代金 | Number | 単位：円 |
| AdjustmentFactor | 調整係数 | Number | 株式分割1:2の場合、権利落ち日に0.5が収録 |
| AdjustmentOpen | 調整済み始値 | Number | 過去の分割等を考慮 |
| AdjustmentHigh | 調整済み高値 | Number | 過去の分割等を考慮 |
| AdjustmentLow | 調整済み安値 | Number | 過去の分割等を考慮 |
| AdjustmentClose | 調整済み終値 | Number | 過去の分割等を考慮 |
| AdjustmentVolume | 調整済み取引高 | Number | 過去の分割等を考慮 |

##### 前場データ（Premiumプランのみ）

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| MorningOpen | 前場始値 | Number | |
| MorningHigh | 前場高値 | Number | |
| MorningLow | 前場安値 | Number | |
| MorningClose | 前場終値 | Number | |
| MorningUpperLimit | 前場ストップ高フラグ | String | 0：ストップ高以外、1：ストップ高 |
| MorningLowerLimit | 前場ストップ安フラグ | String | 0：ストップ安以外、1：ストップ安 |
| MorningVolume | 前場売買高 | Number | |
| MorningTurnoverValue | 前場取引代金 | Number | 単位：円 |
| MorningAdjustmentOpen | 調整済み前場始値 | Number | 過去の分割等を考慮 |
| MorningAdjustmentHigh | 調整済み前場高値 | Number | 過去の分割等を考慮 |
| MorningAdjustmentLow | 調整済み前場安値 | Number | 過去の分割等を考慮 |
| MorningAdjustmentClose | 調整済み前場終値 | Number | 過去の分割等を考慮 |
| MorningAdjustmentVolume | 調整済み前場売買高 | Number | 過去の分割等を考慮 |

##### 後場データ（Premiumプランのみ）

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| AfternoonOpen | 後場始値 | Number | |
| AfternoonHigh | 後場高値 | Number | |
| AfternoonLow | 後場安値 | Number | |
| AfternoonClose | 後場終値 | Number | |
| AfternoonUpperLimit | 後場ストップ高フラグ | String | 0：ストップ高以外、1：ストップ高 |
| AfternoonLowerLimit | 後場ストップ安フラグ | String | 0：ストップ安以外、1：ストップ安 |
| AfternoonVolume | 後場売買高 | Number | |
| AfternoonTurnoverValue | 後場取引代金 | Number | 単位：円 |
| AfternoonAdjustmentOpen | 調整済み後場始値 | Number | 過去の分割等を考慮 |
| AfternoonAdjustmentHigh | 調整済み後場高値 | Number | 過去の分割等を考慮 |
| AfternoonAdjustmentLow | 調整済み後場安値 | Number | 過去の分割等を考慮 |
| AfternoonAdjustmentClose | 調整済み後場終値 | Number | 過去の分割等を考慮 |
| AfternoonAdjustmentVolume | 調整済み後場売買高 | Number | 過去の分割等を考慮 |

## 使用上の注意事項

1. **データの取得制限**
   - 銘柄コード（code）または日付（date）の指定が必須
   - レスポンスが大きすぎる場合はpagination_keyを使用して分割取得

2. **調整済み株価について**
   - 株式分割・併合を考慮した調整済み株価を提供
   - 小数点第２位で四捨五入
   - 一部のコーポレートアクションには未対応

3. **前場/後場データについて**
   - Premiumプランのみで利用可能
   - 日通しデータは全プランで利用可能

4. **特殊なケース**
   - 売買が成立しなかった日：四本値、取引高、売買代金はNull
   - 2020/10/1：システム障害により全データNull
   - 地方取引所単独上場銘柄：データ収録対象外

## サンプルレスポンス

```json
{
    "daily_quotes": [
        {
            "Date": "2023-03-24",
            "Code": "86970",
            "Open": 2047.0,
            "High": 2069.0,
            "Low": 2035.0,
            "Close": 2045.0,
            "UpperLimit": "0",
            "LowerLimit": "0",
            "Volume": 2202500.0,
            "TurnoverValue": 4507051850.0,
            "AdjustmentFactor": 1.0,
            "AdjustmentOpen": 2047.0,
            "AdjustmentHigh": 2069.0,
            "AdjustmentLow": 2035.0,
            "AdjustmentClose": 2045.0,
            "AdjustmentVolume": 2202500.0,
            "MorningOpen": 2047.0,
            "MorningHigh": 2069.0,
            "MorningLow": 2040.0,
            "MorningClose": 2045.5,
            "MorningUpperLimit": "0",
            "MorningLowerLimit": "0",
            "MorningVolume": 1121200.0,
            "MorningTurnoverValue": 2297525850.0,
            "MorningAdjustmentOpen": 2047.0,
            "MorningAdjustmentHigh": 2069.0,
            "MorningAdjustmentLow": 2040.0,
            "MorningAdjustmentClose": 2045.5,
            "MorningAdjustmentVolume": 1121200.0,
            "AfternoonOpen": 2047.0,
            "AfternoonHigh": 2047.0,
            "AfternoonLow": 2035.0,
            "AfternoonClose": 2045.0,
            "AfternoonUpperLimit": "0",
            "AfternoonLowerLimit": "0",
            "AfternoonVolume": 1081300.0,
            "AfternoonTurnoverValue": 2209526000.0,
            "AfternoonAdjustmentOpen": 2047.0,
            "AfternoonAdjustmentHigh": 2047.0,
            "AfternoonAdjustmentLow": 2035.0,
            "AfternoonAdjustmentClose": 2045.0,
            "AfternoonAdjustmentVolume": 1081300.0
        }
    ],
    "pagination_key": "value1.value2."
}
```

## 関連情報

このAPIは株価の四本値データを提供するものです。関連する他のAPIについては以下を参照してください：
- 上場銘柄一覧：`/listed/info`
- 前場四本値：`/prices/prices_am`（前場データのみ）
- 財務情報：`/fins/statements`