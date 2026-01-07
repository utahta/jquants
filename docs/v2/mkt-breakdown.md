# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 売買内訳データ(/markets/breakdown)

GET    /v2/markets/breakdown

## APIの概要

東証上場銘柄の東証市場における銘柄別の日次売買代金・売買高（立会内取引に限る）について、信用取引や空売りの利用に関する発注時のフラグ情報を用いて細分化したデータです。

### 本APIの留意点

- 当該銘柄のコーポレートアクションが発生した場合も、遡及して約定株数の調整は行われません。
- 2020/10/1は東京証券取引所の株式売買システムの障害により終日売買停止となった関係で、データが存在しません。

## 銘柄別の日次売買代金・売買高のデータを取得します

GET https://api.jquants.com/v2/markets/breakdown

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

**from** (string, optional)
  from の指定（e.g. 20210901 or 2021-09-01）

**to** (string, optional)
  to の指定（e.g. 20210907 or 2021-09-07）

**date** (string, optional)
  from と to を指定しないときの日付（e.g. 20210907 or 2021-09-07）

**pagination_key** (string, optional)
  検索の先頭を指定する文字列過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/markets/breakdown 
-H "x-api-key: {loading}" 
-d code="86970" 
-d date="20230324"
```

### Responses

### データ項目概要

**Date** (string)
  売買日（YYYY-MM-DD）

**Code** (string)
  銘柄コード

**LongSellVa** (number)
  実売りの約定代金売りの約定代金の内訳

**ShrtNoMrgnVa** (number)
  空売り（信用新規売りを除く）の約定代金売りの約定代金の内訳

**MrgnSellNewVa** (number)
  信用新規売り（新たな信用売りポジションを作るための売り注文）の約定代金売りの約定代金の内訳

**MrgnSellCloseVa** (number)
  信用返済売り（既存の信用買いポジションを閉じるための売り注文）の約定代金売りの約定代金の内訳

**LongBuyVa** (number)
  現物買いの約定代金買いの約定代金の内訳

**MrgnBuyNewVa** (number)
  信用新規買い（新たな信用買いポジションを作るための買い注文）の約定代金買いの約定代金の内訳

**MrgnBuyCloseVa** (number)
  信用返済買い（既存の信用売りポジションを閉じるための買い注文）の約定代金買いの約定代金の内訳

**LongSellVo** (number)
  実売りの約定株数売りの約定株数の内訳

**ShrtNoMrgnVo** (number)
  空売り（信用新規売りを除く）の約定株数売りの約定株数の内訳

**MrgnSellNewVo** (number)
  信用新規売り（新たな信用売りポジションを作るための売り注文）の約定株数売りの約定株数の内訳

**MrgnSellCloseVo** (number)
  信用返済売り（既存の信用買いポジションを閉じるための売り注文）の約定株数売りの約定株数の内訳

**LongBuyVo** (number)
  現物買いの約定株数買いの約定株数の内訳

**MrgnBuyNewVo** (number)
  信用新規買い（新たな信用買いポジションを作るための買い注文）の約定株数買いの約定株数の内訳

**MrgnBuyCloseVo** (number)
  信用返済買い（既存の信用売りポジションを閉じるための買い注文）の約定株数買いの約定株数の内訳


### レスポンスサンプル

```
{
    "data": [
        {
            "Date": "2015-04-01",
            "Code": "13010",
            "LongSellVa": 115164000.0,
            "ShrtNoMrgnVa": 93561000.0,
            "MrgnSellNewVa": 6412000.0,
            "MrgnSellCloseVa": 23009000.0,
            "LongBuyVa": 185114000.0,
            "MrgnBuyNewVa": 35568000.0,
            "MrgnBuyCloseVa": 17464000.0,
            "LongSellVo": 415000.0,
            "ShrtNoMrgnVo": 337000.0,
            "MrgnSellNewVo": 23000.0,
            "MrgnSellCloseVo": 83000.0,
            "LongBuyVo": 667000.0,
            "MrgnBuyNewVo": 128000.0,
            "MrgnBuyCloseVo": 63000.0
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

