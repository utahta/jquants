# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 取引カレンダー(/markets/calendar)

GET    /v2/markets/calendar

## APIの概要

東証およびOSEにおける営業日、休業日、ならびにOSEにおける祝日取引の有無の情報を取得できます。
配信データは以下のページで公表している内容と同様です。

- 休業日一覧: https://www.jpx.co.jp/corporate/about-jpx/calendar/index.html
- 祝日取引実施日: https://www.jpx.co.jp/derivatives/rules/holidaytrading/index.html

### 本APIの留意点

- 原則として、毎年3月末頃をめどに翌年1年間の営業日および祝日取引実施日（予定）を更新します。

## 営業日のデータを取得します

GET https://api.jquants.com/v2/markets/calendar

データの取得では、休日区分（hol_div）または日付（from/to）の指定が可能です。

### パラメータ及びレスポンス

データの取得では、休日区分（hol_div）または日付（from/to）の指定が可能です。
各パラメータの組み合わせとレスポンスの結果については以下のとおりです。

| hol_div | from /to | レスポンスの結果 | 
|------|------|------|
| ✓ | ✗ | 指定された休日区分について全期間分のデータ | 
| ✓ | ✓ | 指定された休日区分について指定された期間分のデータ | 
| ✗ | ✓ | 指定された期間分のデータ | 
| ✗ | ✗ | 全期間分のデータ | 

休日区分（hol_div）で指定可能なパラメータについては、以下の「休日区分」を参照してください。

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

**hol_div** (string, optional)
  休日区分

**from** (string, optional)
  from の指定（e.g. 20210901 or 2021-09-01）

**to** (string, optional)
  to の指定（e.g. 20210907 or 2021-09-07）


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/markets/calendar 
-H "x-api-key: {loading}" 
-d hol_div="1" 
-d from="20220101" 
-d to="20221231"
```

### Responses

### データ項目概要

**Date** (string)
  日付（YYYY-MM-DD）

**HolDiv** (string)
  休日区分
休日区分を参照


### レスポンスサンプル

```
{
    "data": [
        {
            "Date": "2015-04-01",
            "HolDiv": "1"
        }
    ]
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

