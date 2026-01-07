# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 指数四本値(/indices/bars/daily)

GET    /v2/indices/bars/daily

## APIの概要

各種指数の四本値データを取得することができます。

現在配信している指数につきましては、こちらのページを参照ください。

### 本APIの留意点

- 2022年4月の東証市場区分再編によりマザーズ市場は廃止されていますが、一定のルールに基づき東証マザーズ指数の構成銘柄の入替を行い、2023年11月6日より指数名称を「東証グロース市場250指数」に変更されています。詳細はこちらをご参照ください。
- 2020年10月1日のデータは東京証券取引所の株式売買システムの障害により終日売買停止となった関係で、四本値は前営業日（2020年10月1日）の終値が収録されています。

## 日次の指数四本値データを取得します

GET https://api.jquants.com/v2/indices/bars/daily

データの取得では、指数コード（code）または日付（date）の指定が必須となります。

### パラメータ及びレスポンス

データの取得する際には、指数コード（code）または日付（date）の指定が必須となります。
各パラメータの組み合わせとレスポンスの結果については以下のとおりです。

| code | date | from /to | レスポンスの結果 | 
|------|------|------|------|
| ✓ | ✗ | ✗ | 指定された銘柄について全期間分のデータ | 
| ✓ | ✓ | ✗ | 指定された銘柄について指定された日付のデータ | 
| ✓ | ✗ | ✓ | 指定された銘柄について指定された期間分のデータ | 
| ✗ | ✓ | ✗ | 配信している指数全てについて指定された日付のデータ | 

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

code または date のいずれか一つの指定が必須です。

**code** (string, optional)
  指数コード（e.g. 0000 or 0028）
配信対象の指数コードについては配信対象指数コードを参照してください。

**date** (string, optional)
  from と to を指定しないとき（e.g. 20210907 or 2021-09-07）

**from** (string, optional)
  from の指定（e.g. 20210901 or 2021-09-01）

**to** (string, optional)
  to の指定（e.g. 20210907 or 2021-09-07）

**pagination_key** (string, optional)
  検索の先頭を指定する文字列過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/indices/bars/daily 
-H "x-api-key: {loading}" 
-d code="0028" 
-d date="20231201"
```

### Responses

### データ項目概要

**Date** (string)
  日付（YYYY-MM-DD）

**Code** (string)
  指数コード
配信対象の指数コードはこちらのページを参照ください。

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
            "Date": "2023-12-01",
            "Code": "0028",
            "O": 1199.18,
            "H": 1202.58,
            "L": 1195.01,
            "C": 1200.17
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

