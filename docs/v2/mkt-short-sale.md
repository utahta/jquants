# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 空売り残高報告(/markets/short-sale-report)

GET    /v2/markets/short-sale-report

## APIの概要

「有価証券の取引等の規制に関する内閣府令」に基づき、取引参加者より報告を受けたもののうち、残高割合が0.5％以上のものについての情報を取得できます。

配信データは下記のページで公表している内容と同一ですが、より長いヒストリカルデータを利用可能です。
https://www.jpx.co.jp/markets/public/short-selling/index.html

### 本APIの留意点

- 取引参加者から該当する報告が行われなかった日にはデータは提供されません。
- 「有価証券の取引等の規制に関する内閣府令について」はこちらをご覧ください。https://www.jpx.co.jp/markets/public/short-selling/01.html

## 空売り残高報告データを取得します

GET https://api.jquants.com/v2/markets/short-sale-report

データの取得では、銘柄コード（code）、公表日（disc_date）、計算日（calc_date）のいずれかの指定が必須となります。

### パラメータ及びレスポンス

データの取得では、銘柄コード（code）、公表日（disc_date）、計算日（calc_date）のいずれかの指定が必須となります。
各パラメータの組み合わせとレスポンスの結果については以下のとおりです。

| code | disc_date | disc_date_from/disc_date_to | calc_date | レスポンスの結果 | 
|------|------|------|------|------|
| ✓ | ✗ | ✗ | ✗ | 指定された銘柄について全期間分のデータ | 
| ✓ | ✓ | ✗ | ✗ | 指定された銘柄について指定日（公表日）のデータ | 
| ✓ | ✗ | ✓ | ✗ | 指定された銘柄について指定された期間のデータ | 
| ✓ | ✗ | ✗ | ✓ | 指定された銘柄について指定日（計算日）のデータ | 
| ✗ | ✓ | ✗ | ✗ | 指定日（公表日）の全ての銘柄のデータ | 
| ✗ | ✗ | ✗ | ✓ | 指定日（計算日）の全ての銘柄のデータ | 

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

code / disc_date / calc_date のいずれか一つ以上の指定が必須です。

**code** (string, optional)
  4桁もしくは5桁の銘柄コード（e.g. 8697 or 86970）
4桁の銘柄コードを指定した場合は、普通株式と優先株式の両方が上場している銘柄においては普通株式のデータのみが取得されます。

**disc_date** (string, optional)
  公表日の指定（e.g. 20240301 or 2024-03-01）

**disc_date_from** (string, optional)
  公表日のfromの指定（e.g. 20240301 or 2024-03-01）

**disc_date_to** (string, optional)
  公表日のtoの指定（e.g. 20240301 or 2024-03-01）

**calc_date** (string, optional)
  計算日の指定（e.g. 20240301 or 2024-03-01）

**pagination_key** (string, optional)
  検索の先頭を指定する文字列過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/markets/short-sale-report 
-H "x-api-key: {loading}" 
-d code="86970" 
-d calc_date="20240801"
```

### Responses

### データ項目概要

**DiscDate** (string)
  日付（公表日, YYYY-MM-DD）

**CalcDate** (string)
  日付（計算日, YYYY-MM-DD）

**Code** (string)
  銘柄コード5桁コード

**SSName** (string)
  商号・名称・氏名
取引参加者から報告されたものをそのまま記載しているため、日本語名称または英語名称が混在しています。

**SSAddr** (string)
  住所・所在地

**DICName** (string)
  委託者・投資一任契約の相手方の商号・名称・氏名

**DICAddr** (string)
  委託者・投資一任契約の相手方の住所・所在地

**FundName** (string)
  信託財産・運用財産の名称

**ShrtPosToSO** (number)
  空売り残高割合

**ShrtPosShares** (number)
  空売り残高数量

**ShrtPosUnits** (number)
  空売り残高売買単位数

**PrevRptDate** (string)
  直近計算年月日（YYYY-MM-DD）

**PrevRptRatio** (number)
  直近空売り残高割合

**Notes** (string)
  備考


### レスポンスサンプル

```
{
    "data": [
      {
        "DiscDate": "2024-08-01",
        "CalcDate": "2024-07-31",
        "Code": "13660",
        "SSName": "個人",
        "SSAddr": "",
        "DICName": "",
        "DICAddr": "",
        "FundName": "",
        "ShrtPosToSO": 0.0053,
        "ShrtPosShares": 140000,
        "ShrtPosUnits": 140000,
        "PrevRptDate": "2024-07-22",
        "PrevRptRatio": 0.0043,
        "Notes": ""
      }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

