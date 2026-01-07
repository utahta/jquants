# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 投資部門別情報(/equities/investor-types)

GET    /v2/equities/investor-types

## APIの概要

投資部門別売買状況（株式・金額）のデータを取得することができます。
配信データは下記のページで公表している内容と同一です。データの単位は千円です。
https://www.jpx.co.jp/markets/statistics-equities/investor-type/index.html

### 本APIの留意点

- 2022年4月4日に行われた市場区分見直しに伴い、市場区分に応じた内容となっている統計資料は、見直し後の市場区分に変更して掲載しています。
- 過誤訂正により過去の投資部門別売買状況データが訂正された場合は、本APIでは以下のとおりデータを提供します。

2023年4月3日以前に訂正が公表された過誤訂正：訂正前のデータは提供せず、訂正後のデータのみ提供します。
2023年4月3日以降に訂正が公表された過誤訂正：訂正前と訂正後のデータのいずれも提供します。訂正が生じた場合には、市場名、開始日および終了日を同一とするレコードが追加され、公表日が新しいデータが訂正後、公表日が古いデータが訂正前のデータと識別することが可能です。
- 過誤訂正により過去の投資部門別売買状況データが訂正された場合は、過誤訂正が公表された翌営業日にデータが更新されます。

## 投資部門別売買状況のデータを取得します

GET https://api.jquants.com/v2/equities/investor-types

データの取得では、市場（section）または公表日の日付（from / to）が指定できます。

### パラメータ及びレスポンス

データの取得では、セクション（section）または日付（from / to）の指定が可能です。
各パラメータの組み合わせとレスポンスの結果については以下のとおりです。

| section | from /to | レスポンスの結果 | 
|------|------|------|
| ✓ | ✓ | 指定したセクションの指定した期間のデータ | 
| ✓ | ✗ | 指定したセクションの全期間のデータ | 
| ✗ | ✓ | すべてのセクションの指定した期間のデータ | 
| ✗ | ✗ | すべてのセクションの全期間のデータ | 

セクション（section）で指定可能なパラメータについては、以下の「市場名」を参照してください。

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

**section** (string, optional)
  セクション（e.g. TSEPrime）

**from** (string, optional)
  fromの指定（e.g. 20210901 or 2021-09-01）

**to** (string, optional)
  toの指定（e.g. 20210907 or 2021-09-07）

**pagination_key** (string, optional)
  検索の先頭を指定する文字列
過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/equities/investor-types 
-H "x-api-key: {loading}" 
-d section="TSEPrime" 
-d from="20230324" 
-d to="20230403"
```

### Responses

### データ項目概要

**PubDate** (string)
  公表日（YYYY-MM-DD）

**StDate** (string)
  開始日（YYYY-MM-DD）

**EnDate** (string)
  終了日（YYYY-MM-DD）

**Section** (string)
  市場名（市場名を参照）

**PropSell** (number)
  自己計_売

**PropBuy** (number)
  自己計_買

**PropTot** (number)
  自己計_合計

**PropBal** (number)
  自己計_差引

**BrkSell** (number)
  委託計_売

**BrkBuy** (number)
  委託計_買

**BrkTot** (number)
  委託計_合計

**BrkBal** (number)
  委託計_差引

**TotSell** (number)
  総計_売

**TotBuy** (number)
  総計_買

**TotTot** (number)
  総計_合計

**TotBal** (number)
  総計_差引

**IndSell** (number)
  個人_売

**IndBuy** (number)
  個人_買

**IndTot** (number)
  個人_合計

**IndBal** (number)
  個人_差引

**FrgnSell** (number)
  海外投資家_売

**FrgnBuy** (number)
  海外投資家_買

**FrgnTot** (number)
  海外投資家_合計

**FrgnBal** (number)
  海外投資家_差引

**SecCoSell** (number)
  証券会社_売

**SecCoBuy** (number)
  証券会社_買

**SecCoTot** (number)
  証券会社_合計

**SecCoBal** (number)
  証券会社_差引

**InvTrSell** (number)
  投資信託_売

**InvTrBuy** (number)
  投資信託_買

**InvTrTot** (number)
  投資信託_合計

**InvTrBal** (number)
  投資信託_差引

**BusCoSell** (number)
  事業法人_売

**BusCoBuy** (number)
  事業法人_買

**BusCoTot** (number)
  事業法人_合計

**BusCoBal** (number)
  事業法人_差引

**OthCoSell** (number)
  その他法人_売

**OthCoBuy** (number)
  その他法人_買

**OthCoTot** (number)
  その他法人_合計

**OthCoBal** (number)
  その他法人_差引

**InsCoSell** (number)
  生保・損保_売

**InsCoBuy** (number)
  生保・損保_買

**InsCoTot** (number)
  生保・損保_合計

**InsCoBal** (number)
  生保・損保_差引

**BankSell** (number)
  都銀・地銀等_売

**BankBuy** (number)
  都銀・地銀等_買

**BankTot** (number)
  都銀・地銀等_合計

**BankBal** (number)
  都銀・地銀等_差引

**TrstBnkSell** (number)
  信託銀行_売

**TrstBnkBuy** (number)
  信託銀行_買

**TrstBnkTot** (number)
  信託銀行_合計

**TrstBnkBal** (number)
  信託銀行_差引

**OthFinSell** (number)
  その他金融機関_売

**OthFinBuy** (number)
  その他金融機関_買

**OthFinTot** (number)
  その他金融機関_合計

**OthFinBal** (number)
  その他金融機関_差引


### レスポンスサンプル

```
{
    "data": [
        {
            "PubDate": "2017-01-13",
            "StDate": "2017-01-04",
            "EnDate": "2017-01-06",
            "Section": "TSE1st",
            "PropSell": 1311271004,
            "PropBuy": 1453326508,
            "PropTot": 2764597512,
            "PropBal": 142055504,
            "BrkSell": 7165529005,
            "BrkBuy": 7030019854,
            "BrkTot": 14195548859,
            "BrkBal": -135509151,
            "TotSell": 8476800009,
            "TotBuy": 8483346362,
            "TotTot": 16960146371,
            "TotBal": 6546353,
            "IndSell": 1401711615,
            "IndBuy": 1161801155,
            "IndTot": 2563512770,
            "IndBal": -239910460,
            "FrgnSell": 5094891735,
            "FrgnBuy": 5317151774,
            "FrgnTot": 10412043509,
            "FrgnBal": 222260039,
            "SecCoSell": 76381455,
            "SecCoBuy": 61700100,
            "SecCoTot": 138081555,
            "SecCoBal": -14681355,
            "InvTrSell": 168705109,
            "InvTrBuy": 124389642,
            "InvTrTot": 293094751,
            "InvTrBal": -44315467,
            "BusCoSell": 71217959,
            "BusCoBuy": 63526641,
            "BusCoTot": 134744600,
            "BusCoBal": -7691318,
            "OthCoSell": 10745152,
            "OthCoBuy": 15687836,
            "OthCoTot": 26432988,
            "OthCoBal": 4942684,
            "InsCoSell": 15926202,
            "InsCoBuy": 9831555,
            "InsCoTot": 25757757,
            "InsCoBal": -6094647,
            "BankSell": 10606789,
            "BankBuy": 8843871,
            "BankTot": 19450660,
            "BankBal": -1762918,
            "TrstBnkSell": 292932297,
            "TrstBnkBuy": 245322795,
            "TrstBnkTot": 538255092,
            "TrstBnkBal": -47609502,
            "OthFinSell": 22410692,
            "OthFinBuy": 21764485,
            "OthFinTot": 44175177,
            "OthFinBal": -646207
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

