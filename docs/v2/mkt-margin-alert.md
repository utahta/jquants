# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 日々公表信用取引残高(/markets/margin-alert)

GET    /v2/markets/margin-alert

## APIの概要

日々公表銘柄に指定された個別銘柄の日々の信用取引残高を取得することができます。
各銘柄についての信用取引残高（株数）を取得できます。

配信データは下記のページで公表している内容と同一です。
https://www.jpx.co.jp/markets/statistics-equities/margin/index.html

### 本APIの留意点

- 当該銘柄のコーポレートアクションが発生した場合であっても、遡及して約定株数の調整は行われません。
- 東京証券取引所または日本証券金融が、日次の信用取引残高を公表する必要があると認めた銘柄のみが収録されます。
- 過誤訂正により過去の日々公表信用取引残高データが訂正された場合は、本APIでは以下のとおりデータを提供します。

訂正前と訂正後のデータのいずれも提供します。訂正が生じた場合には、申込日を同一とするレコードが追加されます。公表日が新しいデータが訂正後、公表日が古いデータが訂正前のデータと識別することが可能です。

## 日々公表信用取引残高を取得します

GET https://api.jquants.com/v2/markets/margin-alert

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
  銘柄コード（e.g. 27800 or 2780）4桁の銘柄コードを指定した場合は、普通株式と優先株式の両方が上場している銘柄においては普通株式のデータのみが取得されます。

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
curl -G https://api.jquants.com/v2/markets/margin-alert 
-H "x-api-key: {loading}" 
-d code="13260" 
-d date="20240208"
```

### Responses

### データ項目概要

**PubDate** (string)
  公表日

**Code** (string)
  銘柄コード

**AppDate** (string)
  申込日（YYYY-MM-DD）信用取引残高の基準となる時点を表します。

**PubReason**
  公表の理由

**ShrtOut** (number)
  売合計信用残高

**ShrtOutChg**
  前日比 売合計信用残高（単位：株）前日に公表されていない銘柄の場合、「-」を出力します。

**ShrtOutRatio**
  上場比 売合計信用残高（単位：％） 売合計信用残高 ÷ 上場株数 × 100ETF の場合、「*」を出力します。

**LongOut** (number)
  買合計信用残高

**LongOutChg**
  前日比 買合計信用残高（単位：株）前日に公表されていない銘柄の場合、「-」を出力します。

**LongOutRatio**
  上場比 買合計信用残高（単位：％） 買合計信用残高 ÷ 上場株数 × 100ETF の場合、「*」を出力します。

**SLRatio** (number)
  取組比率（単位：％） 売合計信用残高 ÷ 買合計信用残高 × 100

**ShrtNegOut** (number)
  一般信用取引売残高売合計信用残高のうち、一般信用によるものです。

**ShrtNegOutChg**
  前日比 一般信用取引売残高（単位：株）前日に公表されていない銘柄の場合、「-」を出力します。

**ShrtStdOut** (number)
  制度信用取引売残高売合計信用残高のうち、制度信用によるものです。

**ShrtStdOutChg**
  前日比 制度信用取引売残高（単位：株）前日に公表されていない銘柄の場合、「-」を出力します。

**LongNegOut** (number)
  一般信用取引買残高買合計信用残高のうち、一般信用によるものです。

**LongNegOutChg**
  前日比 一般信用取引買残高（単位：株）前日に公表されていない銘柄の場合、「-」を出力します。

**LongStdOut** (number)
  制度信用取引買残高買合計信用残高のうち、制度信用によるものです。

**LongStdOutChg**
  前日比 制度信用取引買残高（単位：株）前日に公表されていない銘柄の場合、「-」を出力します。

**TSEMrgnRegCls** (string)
  東証信用貸借規制区分


### レスポンスサンプル

```
{
    "data": [
        {
            "PubDate": "2024-02-08",
            "Code": "13260",
            "AppDate": "2024-02-07",
            "PubReason":
                {
                    "Restricted": "0",
                    "DailyPublication": "0",
                    "Monitoring": "0",
                    "RestrictedByJSF": "0",
                    "PrecautionByJSF": "1",
                    "UnclearOrSecOnAlert": "0"
                },
            "ShrtOut": 11.0,
            "ShrtOutChg": 0.0,
            "ShrtOutRatio": "*",
            "LongOut": 676.0,
            "LongOutChg": -20.0,
            "LongOutRatio": "*",
            "SLRatio": 1.6,
            "ShrtNegOut": 0.0,
            "ShrtNegOutChg": 0.0,
            "ShrtStdOut": 11.0,
            "ShrtStdOutChg": 0.0,
            "LongNegOut": 192.0,
            "LongNegOutChg": -20.0,
            "LongStdOut": 484.0,
            "LongStdOutChg": 0.0,
            "TSEMrgnRegCls": "001"
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

