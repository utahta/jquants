# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 業種別空売り比率(/markets/short-ratio)

GET    /v2/markets/short-ratio

## APIの概要

日々の業種（セクター）別の空売り比率に関する売買代金を取得できます。
配信データは下記のページで公表している内容と同様です。
https://www.jpx.co.jp/markets/statistics-equities/short-selling/index.html
Webページでの公表値は百万円単位に丸められておりますが、APIでは円単位のデータとなります。

### 本APIの留意点

- 取引高が存在しない（売買されていない）日の日付（date）を指定した場合は、値は空です。
- 2020/10/1は東京証券取引所の株式売買システムの障害により終日売買停止となった関係で、データが存在しません。

## 日々の業種（セクター）別の空売り比率・売買代金を取得します

GET https://api.jquants.com/v2/markets/short-ratio

データの取得では、33業種コード（s33）または日付（date）の指定が必須となります。

### パラメータ及びレスポンス

データの取得では、日付（date）または33業種コード（s33）の指定が必須となります。
各パラメータの組み合わせとレスポンスの結果については以下のとおりです。

| s33 | date | from/to | レスポンスの結果 | 
|------|------|------|------|
| ✗ | ✓ | ✗ | 全業種コードについて指定された日付のデータ | 
| ✓ | ✗ | ✗ | 指定された業種コードについて、全期間分のデータ | 
| ✓ | ✗ | ✓ | 指定された業種コードについて指定された期間分のデータ | 
| ✓ | ✓ | ✗ | 指定された業種コードについて指定された日付のデータ | 

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

s33 または date のいずれか一つの指定が必須です。

**s33** (string, optional)
  33業種コード（e.g. 0050 or 50）

**from** (string, optional)
  fromの指定（e.g. 20210901 or 2021-09-01）

**to** (string, optional)
  toの指定（e.g. 20210907 or 2021-09-07）

**date** (string, optional)
  *fromとtoを指定しないとき（e.g. 20210907 or 2021-09-07）

**pagination_key** (string, optional)
  検索の先頭を指定する文字列過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/markets/short-ratio 
-H "x-api-key: {loading}" 
-d s33="0050" 
-d date="2022-10-25"
```

### Responses

### データ項目概要

**Date** (string)
  日付（YYYY-MM-DD）

**S33** (string)
  33業種コード（33業種コード及び業種名を参照）

**SellExShortVa** (number)
  実注文の売買代金

**ShrtWithResVa** (number)
  価格規制有りの空売り売買代金

**ShrtNoResVa** (number)
  価格規制無しの空売り売買代金


### レスポンスサンプル

```
{
    "data": [
        {
            "Date": "2022-10-25",
            "S33": "0050",
            "SellExShortVa": 1333126400.0,
            "ShrtWithResVa": 787355200.0,
            "ShrtNoResVa": 149084300.0
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

