# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 信用取引週末残高(/markets/margin-interest)

GET    /v2/markets/margin-interest

## APIの概要

週末時点での、各銘柄についての信用取引残高（株数）を取得できます。

配信データは下記のページで公表している内容と同一です。
https://www.jpx.co.jp/markets/statistics-equities/margin/index.html

### 本APIの留意点

- 当該銘柄のコーポレートアクションが発生した場合も、遡及して株数の調整は行われません。
- 年末年始など、営業日が2日以下の週はデータが提供されません。
- 東証上場銘柄でない銘柄（地方取引所単独上場銘柄）についてはデータの収録対象外となっております。

## 信用取引週末残高を取得します

GET https://api.jquants.com/v2/markets/margin-interest

データの取得では、銘柄コード（code）または公表日（date）の指定が必須となります。

### パラメータ及びレスポンス

データの取得では、銘柄コード（code）または公表日（date）の指定が必須となります。
各パラメータの組み合わせとレスポンスの結果については以下のとおりです。

| code | date | from /to | レスポンスの結果 | 
|------|------|------|------|
| ✓ | ✗ | ✗ | 指定された銘柄について全期間分のデータ | 
| ✓ | ✓ | ✗ | 指定された銘柄について指定された公表日のデータ | 
| ✓ | ✗ | ✓ | 指定された銘柄について指定された期間分のデータ | 
| ✗ | ✓ | ✗ | 全上場銘柄について指定された公表日のデータ | 

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
  from と to を指定しないときの公表日（e.g. 20210907 or 2021-09-07）

**pagination_key** (string, optional)
  検索の先頭を指定する文字列過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/markets/margin-interest 
-H "x-api-key: {loading}" 
-d code="86970" 
-d date="20230324"
```

### Responses

### データ項目概要

**Date** (string)
  申込日付信用取引残高基準となる時点を表します。（通常は金曜日付）（YYYY-MM-DD）

**Code** (string)
  銘柄コード

**ShrtVol** (number)
  売合計信用残高

**LongVol** (number)
  買合計信用残高

**ShrtNegVol** (number)
  一般信用取引売残高売合計信用残高のうち、一般信用によるものです。

**LongNegVol** (number)
  一般信用取引買残高買合計信用残高のうち、一般信用によるものです。

**ShrtStdVol** (number)
  制度信用取引売残高売合計信用残高のうち、制度信用によるものです。

**LongStdVol** (number)
  制度信用取引買残高買合計信用残高のうち、制度信用によるものです。

**IssType** (string)
  銘柄区分1: 信用銘柄、2: 貸借銘柄、3: その他


### レスポンスサンプル

```
{
    "data": [
        {
            "Date": "2023-03-24",
            "Code": "86970",
            "ShrtVol": 123456.0,
            "LongVol": 234567.0,
            "ShrtNegVol": 11111.0,
            "LongNegVol": 22222.0,
            "ShrtStdVol": 33333.0,
            "LongStdVol": 44444.0,
            "IssType": "1"
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

