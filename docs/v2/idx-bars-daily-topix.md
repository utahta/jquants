# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# TOPIX指数四本値(/indices/bars/daily/topix)

GET    /v2/indices/bars/daily/topix

## APIの概要

TOPIXの日通しの四本値を取得できます。
本APIで取得可能な指数データは TOPIX（東証株価指数）のみとなります。

## 日次のTOPIX指数データを取得します

GET https://api.jquants.com/v2/indices/bars/daily/topix

日付の範囲（from/to）を指定することができます。なお、指定しない場合は全期間のデータがレスポンスに収録されます。

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

**from** (string, optional)
  from の指定（e.g. 20210901 or 2021-09-01）

**to** (string, optional)
  to の指定（e.g. 20210907 or 2021-09-07）

**pagination_key** (string, optional)
  検索の先頭を指定する文字列過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/indices/bars/daily/topix 
-H "x-api-key: {loading}" 
-d from="20220101" 
-d to="20221231"
```

### Responses

### データ項目概要

**Date** (string)
  日付（YYYY-MM-DD）

**O** (number)
  始値

**H** (number)
  高値

**L** (number)
  安値

**C** (number)
  終値


### レスポンスサンプル

```
{
    "data": [
        {
            "Date": "2022-06-28",
            "O": 1885.52,
            "H": 1907.38,
            "L": 1885.32,
            "C": 1907.38
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

