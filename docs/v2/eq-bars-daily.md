# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 株価四本値(/equities/bars/daily)

GET    /v2/equities/bars/daily

## APIの概要

株価データを取得することができます。
株価は分割・併合を考慮した調整済み株価（小数点第２位四捨五入）と調整前の株価を取得することができます。

### 本APIの留意点

- 取引が存在しない日の銘柄についての四本値、取引高と売買代金は、Nullが収録されています。
- 東証上場銘柄でない銘柄（地方取引所単独上場銘柄）についてはデータの収録対象外となっております。
- 2020/10/1のデータは東京証券取引所の株式売買システムの障害により終日売買停止となった関係で、四本値、取引高と売買代金はNullが収録されています。
- 日通しデータについては全プランで取得できますが、前場/後場別のデータについてはPremiumプランのみ取得可能です。
- 株価調整については株式分割・併合にのみ対応しております。一部コーポレートアクションには対応しておりませんので、ご了承ください。

## 日次の株価データを取得します

GET https://api.jquants.com/v2/equities/bars/daily

データの取得では、銘柄コード（code）または日付（date）の指定が必須となります。

### パラメータ及びレスポンス

データの取得では、銘柄コード（code）または日付（date）の指定が必須となります。
各パラメータの組み合わせとレスポンスの結果については以下のとおりです。

| code | date | from /to | レスポンスの結果 | 
|------|------|------|------|
| ✓ | ✗ | ✗ | 指定された銘柄について全期間分のデータ | 
| ✓ | ✓ | ✗ | 指定された銘柄について指定された日付のデータ | 
| ✓ | ✗ | ✓ | 指定された銘柄について指定された期間分のデータ | 
| ✗ | ✓ | ✗ | 全上場銘柄について指定された日付のデータ | 

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

code または date のいずれか一つの指定が必須です。

**code** (string, optional)
  銘柄コード（e.g. 27800 or 2780）4桁の銘柄コードを指定した場合は、普通株式と優先株式等の両方が上場している銘柄においては普通株式のデータのみが取得されます。

**date** (string, optional)
  from と to を指定しないとき（e.g. 20210907 or 2021-09-07）

**from** (string, optional)
  fromの指定（e.g. 20210901 or 2021-09-01）

**to** (string, optional)
  to の指定（e.g. 20210907 or 2021-09-07）

**pagination_key** (string, optional)
  検索の先頭を指定する文字列過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/equities/bars/daily 
-H "x-api-key: {loading}" 
-d code="86970" 
-d date="20230324"
```

### Responses

### データ項目概要

**Date** (string)
  日付（YYYY-MM-DD）

**Code** (string)
  銘柄コード

**O** (number)
  始値（調整前）

**H** (number)
  高値（調整前）

**L** (number)
  安値（調整前）

**C** (number)
  終値（調整前）

**UL** (string)
  日通ストップ高フラグ（0：ストップ高以外, 1：ストップ高）

**LL** (string)
  日通ストップ安フラグ（0：ストップ安以外, 1：ストップ安）

**Vo** (number)
  取引高（調整前）

**Va** (number)
  取引代金

**AdjFactor** (number)
  調整係数（株式分割1:2の場合、権利落ち日に 0.5 が入る）

**AdjO** (number)
  調整済み始値（※1）

**AdjH** (number)
  調整済み高値（※1）

**AdjL** (number)
  調整済み安値（※1）

**AdjC** (number)
  調整済み終値（※1）

**AdjVo** (number)
  調整済み取引高（※1）

**MO** (number)
  前場始値（※2）

**MH** (number)
  前場高値（※2）

**ML** (number)
  前場安値（※2）

**MC** (number)
  前場終値（※2）

**MUL** (string)
  前場ストップ高フラグ（0：ストップ高以外, 1：ストップ高）,（※2）

**MLL** (string)
  前場ストップ安フラグ（0：ストップ安以外, 1：ストップ安）,（※2）

**MVo** (number)
  前場売買高（※2）

**MVa** (number)
  前場取引代金（※2）

**MAdjO** (number)
  調整済み前場始値（※1, ※2）

**MAdjH** (number)
  調整済み前場高値（※1, ※2）

**MAdjL** (number)
  調整済み前場安値（※1, ※2）

**MAdjC** (number)
  調整済み前場終値（※1, ※2）

**MAdjVo** (number)
  調整済み前場売買高（※1, ※2）

**AO** (number)
  後場始値（※2）

**AH** (number)
  後場高値（※2）

**AL** (number)
  後場安値（※2）

**AC** (number)
  後場終値（※2）

**AUL** (string)
  後場ストップ高フラグ（0：ストップ高以外, 1：ストップ高）,（※2）

**ALL** (string)
  後場ストップ安フラグ（0：ストップ安以外, 1：ストップ安）,（※2）

**AVo** (number)
  後場売買高（※2）

**AVa** (number)
  後場取引代金（※2）

**AAdjO** (number)
  調整済み後場始値（※1, ※2）

**AAdjH** (number)
  調整済み後場高値（※1, ※2）

**AAdjL** (number)
  調整済み後場安値（※1, ※2）

**AAdjC** (number)
  調整済み後場終値（※1, ※2）

**AAdjVo** (number)
  調整済み後場売買高（※1, ※2）


※1 過去の分割等を考慮した調整済みの項目です
※2 Premiumプランのみ取得可能な項目です

### レスポンスサンプル

```
{
    "data": [
        {
            "Date": "2023-03-24",
            "Code": "86970",
            "O": 2047.0,
            "H": 2069.0,
            "L": 2035.0,
            "C": 2045.0,
            "UL": "0",
            "LL": "0",
            "Vo": 2202500.0,
            "Va": 4507051850.0,
            "AdjFactor": 1.0,
            "AdjO": 2047.0,
            "AdjH": 2069.0,
            "AdjL": 2035.0,
            "AdjC": 2045.0,
            "AdjVo": 2202500.0,
            "MO": 2047.0,
            "MH": 2069.0,
            "ML": 2040.0,
            "MC": 2045.5,
            "MUL": "0",
            "MLL": "0",
            "MVo": 1121200.0,
            "MVa": 2297525850.0,
            "MAdjO": 2047.0,
            "MAdjH": 2069.0,
            "MAdjL": 2040.0,
            "MAdjC": 2045.5,
            "MAdjVo": 1121200.0,
            "AO": 2047.0,
            "AH": 2047.0,
            "AL": 2035.0,
            "AC": 2045.0,
            "AUL": "0",
            "ALL": "0",
            "AVo": 1081300.0,
            "AVa": 2209526000.0,
            "AAdjO": 2047.0,
            "AAdjH": 2047.0,
            "AAdjL": 2035.0,
            "AAdjC": 2045.0,
            "AAdjVo": 1081300.0
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

