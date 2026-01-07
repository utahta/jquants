# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 日経225オプション四本値(/derivatives/bars/daily/options/225)

GET    /v2/derivatives/bars/daily/options/225

## APIの概要

日経225オプションに関する、四本値や清算値段、理論価格に関する情報を取得することができます。
また、本APIで取得可能なデータは日経225指数オプション（Weeklyオプション及びフレックスオプションを除く）のみとなります。

## 本APIの留意点

- 取引セッションについて

2011年2月10日以前は、ナイトセッション、前場、後場で構成されています。
この期間の前場データは収録されず、後場データが日中場データとして収録されます。なお、日通しデータについては、全立会を含めたデータとなります。
2011年2月14日以降は、ナイトセッション、日中場で構成されています。
- レスポンスのキー項目について

緊急取引証拠金が発動した場合は、同一の取引日・銘柄に対して清算価格算出時と緊急取引証拠金算出時のデータが発生します。そのため、Date、Codeに加えてEmMrgnTrgDiv（EmergencyMarginTriggerDivision）を組み合わせることでデータを一意に識別することが可能です。

## 日次の日経225オプションデータ取得

GET https://api.jquants.com/v2/derivatives/bars/daily/options/225

日付（date）の指定が必須です。

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

date の指定が必須です。

**date** (string, required)
  date の指定（e.g. 20210901 or 2021-09-01）

**pagination_key** (string, optional)
  検索の先頭を指定する文字列
過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/derivatives/bars/daily/options/225
-H "x-api-key: {loading}"
-d date="20230324"
```

### Responses

### データ項目概要

**Date** (string)
  取引日（YYYY-MM-DD）

**Code** (string)
  銘柄コード

**O** (number)
  日通し始値

**H** (number)
  日通し高値

**L** (number)
  日通し安値

**C** (number)
  日通し終値

**EO** (number | string)
  ナイト・セッション始値
取引開始日初日の銘柄はナイト・セッションが存在しないため、空文字を設定。

**EH** (number | string)
  ナイト・セッション高値
取引開始日初日の銘柄はナイト・セッションが存在しないため、空文字を設定。

**EL** (number | string)
  ナイト・セッション安値
取引開始日初日の銘柄はナイト・セッションが存在しないため、空文字を設定。

**EC** (number | string)
  ナイト・セッション終値
取引開始日初日の銘柄はナイト・セッションが存在しないため、空文字を設定。

**AO** (number)
  日中始値

**AH** (number)
  日中高値

**AL** (number)
  日中安値

**AC** (number)
  日中終値

**Vo** (number)
  取引高

**OI** (number)
  建玉

**Va** (number)
  取引代金

**CM** (string)
  限月（YYYY-MM）

**Strike** (number)
  権利行使価格

**VoOA** (number)
  立会内取引高（※1）

**EmMrgnTrgDiv** (string)
  緊急取引証拠金発動区分
001: 緊急取引証拠金発動時、002: 清算価格算出時。
"001" は2016年7月19日以降に緊急取引証拠金発動した場合のみ収録。

**PCDiv** (string)
  プットコール区分
1: プット、2: コール

**LTD** (string)
  取引最終年月日（YYYY-MM-DD）（※1）

**SQD** (string)
  SQ日（YYYY-MM-DD）（※1）

**Settle** (number)
  清算値段（※1）

**Theo** (number)
  理論価格（※1）

**BaseVol** (number)
  基準ボラティリティ
アット・ザ・マネープット及びコールそれぞれのインプライドボラティリティの中間値（※1）

**UnderPx** (number)
  原証券価格（※1）

**IV** (number)
  インプライドボラティリティ（※1）

**IR** (number)
  理論価格計算用金利（※1）


※1 2016年7月19日以降のみ提供。

### レスポンスサンプル

```
{
    "data": [
        {
            "Date": "2023-03-22",
            "Code": "130060018",
            "O": 0.0,
            "H": 0.0,
            "L": 0.0,
            "C": 0.0,
            "EO": 0.0,
            "EH": 0.0,
            "EL": 0.0,
            "EC": 0.0,
            "AO": 0.0,
            "AH": 0.0,
            "AL": 0.0,
            "AC": 0.0,
            "Vo": 0.0,
            "OI": 330.0,
            "Va": 0.0,
            "CM": "2025-06",
            "Strike": 20000.0,
            "VoOA": 0.0,
            "EmMrgnTrgDiv": "002",
            "PCDiv": "1",
            "LTD": "2025-06-12",
            "SQD": "2025-06-13",
            "Settle": 980.0,
            "Theo": 974.641,
            "BaseVol": 17.93025,
            "UnderPx": 27466.61,
            "IV": 23.1816,
            "IR": 0.2336
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

