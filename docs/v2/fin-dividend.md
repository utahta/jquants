# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 配当金情報(/fins/dividend)

GET    /v2/fins/dividend

## APIの概要

上場会社の配当（決定・予想）に関する１株当たり配当金額、基準日、権利落日及び支払開始予定日等の情報を取得できます。

## 本APIの留意点

- 東証上場銘柄でない銘柄（地方取引所単独上場銘柄）についてはデータの収録対象外となっております。

## 配当金データを取得します

GET https://api.jquants.com/v2/fins/dividend

データの取得では、銘柄コード（code）または通知日付（date）の指定が必須となります。

### パラメータ及びレスポンス

データの取得では、銘柄コード（code）または通知日付（date）の指定が必須となります。
各パラメータの組み合わせとレスポンスの結果については以下のとおりです。

| code | date | from /to | レスポンスの結果 | 
|------|------|------|------|
| ✓ | ✗ | ✗ | 指定された銘柄について取得可能期間の全データ | 
| ✓ | ✓ | ✗ | 指定された銘柄について指定された通知日付のデータ | 
| ✓ | ✗ | ✓ | 指定された銘柄について指定された期間分のデータ | 
| ✗ | ✓ | ✗ | 全上場銘柄について指定された通知日付のデータ | 

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

code または date のいずれか一つの指定が必須です。

**code** (string, optional)
  銘柄コード
（e.g. 27800 or 2780）
4桁の銘柄コードを指定した場合は、普通株式と優先株式の両方が上場している銘柄においては普通株式のデータのみが取得されます。

**from** (string, optional)
  fromの指定（e.g. 20210901 or 2021-09-01）

**to** (string, optional)
  toの指定（e.g. 20210907 or 2021-09-07）

**date** (string, optional)
  *fromとtoを指定しないとき（e.g. 20210907 or 2021-09-07）

**pagination_key** (string, optional)
  検索の先頭を指定する文字列
過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/fins/dividend
-H "x-api-key: {loading}"
-d code="86970"
-d date="2024-10-25"
```

### Responses

### データ項目概要

**PubDate** (string)
  通知日時（YYYY-MM-DD）

**PubTime** (string)
  通知日時（HH:MM）

**Code** (string)
  銘柄コード

**RefNo** (string)
  リファレンスナンバー
配当通知を一意に特定するための番号
詳細はリファレンスナンバーを参照

**StatCode** (string)
  更新区分（コード）
1: 新規、2: 訂正、3: 削除

**BoardDate** (string)
  取締役会決議日

**IFCode** (string)
  配当種類（コード）
1: 中間配当、2: 期末配当

**FRCode** (string)
  予想／決定（コード）
1: 決定、2: 予想

**IFTerm** (string)
  配当基準日年月

**DivRate** (number | string)
  １株当たり配当金額
未定の場合: - 、非設定の場合: 空文字

**RecDate** (string)
  基準日

**ExDate** (string)
  権利落日

**ActRecDate** (string)
  権利確定日

**PayDate** (string)
  支払開始予定日
未定の場合: - 、非設定の場合: 空文字

**CARefNo** (string)
  ＣＡリファレンスナンバー
訂正・削除の対象となっている配当通知のリファレンスナンバー。新規の場合はリファレンスナンバーと同じ値を設定
詳細はリファレンスナンバーを参照

**DistAmt** (number | string)
  1株当たりの交付金銭等の額
未定の場合: - 、非設定の場合: 空文字 が設定されます。
2014年2月24日以降のみ提供。

**RetEarn** (number | string)
  1株当たりの利益剰余金の額
未定の場合: - 、非設定の場合: 空文字 が設定されます。
2014年2月24日以降のみ提供。

**DeemDiv** (number | string)
  1株当たりのみなし配当の額
未定の場合: - 、非設定の場合: 空文字 が設定されます。
2014年2月24日以降のみ提供。

**DeemCapGains** (number | string)
  1株当たりのみなし譲渡収入の額
未定の場合: - 、非設定の場合: 空文字 が設定されます。
2014年2月24日以降のみ提供。

**NetAssetDecRatio** (number | string)
  純資産減少割合
未定の場合: - 、非設定の場合: 空文字 が設定されます。
2014年2月24日以降のみ提供。

**CommSpecCode** (string)
  記念配当/特別配当コード
1: 記念配当、2: 特別配当、3: 記念・特別配当、0: 通常の配当

**CommDivRate** (number | string)
  １株当たり記念配当金額
未定の場合: - 、非設定の場合: 空文字
2022年6月6日以降のみ提供。

**SpecDivRate** (number | string)
  １株当たり特別配当金額
未定の場合: - 、非設定の場合: 空文字
2022年6月6日以降のみ提供。


### レスポンスサンプル

```
{
    "data": [
        {
            "PubDate": "2014-02-24",
            "PubTime": "09:21",
            "Code": "15550",
            "RefNo": "201402241B00002",
            "StatCode": "1",
            "BoardDate": "2014-02-24",
            "IFCode": "2",
            "FRCode": "2",
            "IFTerm": "2014-03",
            "DivRate": "-",
            "RecDate": "2014-03-10",
            "ExDate": "2014-03-06",
            "ActRecDate": "2014-03-10",
            "PayDate": "-",
            "CARefNo": "201402241B00002",
            "DistAmt": "",
            "RetEarn": "",
            "DeemDiv": "",
            "DeemCapGains": "",
            "NetAssetDecRatio": "",
            "CommSpecCode": "0",
            "CommDivRate": "",
            "SpecDivRate": ""
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

