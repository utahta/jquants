# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 前場四本値(/equities/bars/daily/am)

GET    /v2/equities/bars/daily/am

## APIの概要

前場終了時に、前場の株価データを取得することができます。

### 本APIの留意点

- 前場の取引高が存在しない（売買されていない）銘柄についての四本値、取引高と売買代金は、null が収録されています。
- 東証上場銘柄でない銘柄（地方取引所単独上場銘柄）についてはデータの収録対象外となっております。
- なお、当日のデータは翌日6:00頃まで取得可能です。ヒストリカルの前場四本値については
株価四本値(/equities/bars/daily)
をご利用ください。

## 前場の株価データを取得します

GET https://api.jquants.com/v2/equities/bars/daily/am

データの取得では、銘柄コード（code）が指定できます。

### パラメータ及びレスポンス

データの取得では、銘柄コード（code）の指定が可能です。
各パラメータの組み合わせとレスポンスの結果については以下のとおりです。

| code | レスポンスの結果 | 
|------|------|
| ✓ | 指定された銘柄についての前場の株価データ | 
| ✗ | 全上場銘柄について前場の株価データ | 

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

**code** (string, optional)
  銘柄コード（e.g. 27800 or 2780）
4桁の銘柄コードを指定した場合は、普通株式と優先株式の両方が上場している銘柄においては普通株式のデータのみが取得されます。

**pagination_key** (string, optional)
  検索の先頭を指定する文字列
過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/equities/bars/daily/am 
-H "x-api-key: {loading}" 
-d code="39400"
```

### Responses

### データ項目概要

**Date** (string)
  日付（YYYY-MM-DD）

**Code** (string)
  銘柄コード

**MO** (number)
  前場始値

**MH** (number)
  前場高値

**ML** (number)
  前場安値

**MC** (number)
  前場終値

**MVo** (number)
  前場売買高

**MVa** (number)
  前場取引代金


### レスポンスサンプル

```
{
    "data": [
        {
            "Date": "2023-03-20",
            "Code": "39400",
            "MO": 232.0,
            "MH": 244.0,
            "ML": 232.0,
            "MC": 240.0,
            "MVo": 52600.0,
            "MVa": 12518800.0
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

