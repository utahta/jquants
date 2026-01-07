# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 先物四本値(/derivatives/bars/daily/futures)

GET    /v2/derivatives/bars/daily/futures

## APIの概要

先物に関する、四本値や清算値段、理論価格に関する情報を取得することができます。
また、本APIで取得可能なデータについては
先物商品区分コード一覧を参照ください。

## 本APIの留意点

- 銘柄コードについて

先物・オプション取引識別コードの付番規則については証券コード関係の関係資料等を参照してください。
- 取引セッションについて

2011年2月10日以前は、ナイトセッション、前場、後場で構成されています。
この期間の前場データは収録されず、後場データが日中場データとして収録されます。なお、日通しデータについては、全立会を含めたデータとなります。
2011年2月14日以降は、ナイトセッション、日中場で構成されています。
- 祝日取引について

祝日取引の取引日については、祝日取引実施日直前の平日に開始するナイト・セッション（祝日前営業日）及び祝日取引実施日直後の平日（祝日翌営業日）のデイ・セッションと同一の取引日として扱います。
- レスポンスのキー項目について

緊急取引証拠金が発動した場合は、同一の取引日・銘柄に対して清算価格算出時と緊急取引証拠金算出時のデータが発生します。そのため、Date、Codeに加えてEmMrgnTrgDiv（EmergencyMarginTriggerDivision）を組み合わせることでデータを一意に識別することが可能です。

## 日次の先物四本値データ取得

GET https://api.jquants.com/v2/derivatives/bars/daily/futures

データの取得では、日付（date）の指定が必須となります。

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

date の指定が必須です。

**category** (string, optional)
  商品区分の指定

**date** (string, required)
  date の指定（e.g. 20210901 or 2021-09-01）

**contract_flag** (string, optional)
  中心限月フラグの指定

**pagination_key** (string, optional)
  検索の先頭を指定する文字列
過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/derivatives/bars/daily/futures
-H "x-api-key: {loading}"
-d date="20230324"
```

### Responses

### データ項目概要

**Code** (string)
  銘柄コード

**ProdCat** (string)
  先物商品区分

**Date** (string)
  取引日（YYYY-MM-DD）

**O** (number)
  日通し始値

**H** (number)
  日通し高値

**L** (number)
  日通し安値

**C** (number)
  日通し終値

**MO** (number | string)
  前場始値
前後場取引対象銘柄でない場合、空文字を設定。

**MH** (number | string)
  前場高値
前後場取引対象銘柄でない場合、空文字を設定。

**ML** (number | string)
  前場安値
前後場取引対象銘柄でない場合、空文字を設定。

**MC** (number | string)
  前場終値
前後場取引対象銘柄でない場合、空文字を設定。

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

**VoOA** (number)
  立会内取引高（※1）

**EmMrgnTrgDiv** (string)
  緊急取引証拠金発動区分
001: 緊急取引証拠金発動時、002: 清算価格算出時。
"001" は2016年7月19日以降に緊急取引証拠金発動した場合のみ収録。

**LTD** (string)
  取引最終年月日（YYYY-MM-DD）（※1）

**SQD** (string)
  SQ日（YYYY-MM-DD）（※1）

**Settle** (number)
  清算値段（※1）

**CCMFlag** (string)
  中心限月フラグ（1:中心限月、0:その他）（※1）


※1 2016年7月19日以降のみ提供。

### レスポンスサンプル

```
{
    "data": [
        {
            "Code": "169090005",
            "ProdCat": "TOPIXF",
            "Date": "2024-07-23",
            "O": 2825.5,
            "H": 2853.0,
            "L": 2825.5,
            "C": 2829.0,
            "MO": "",
            "MH": "",
            "ML": "",
            "MC": "",
            "EO": 2825.5,
            "EH": 2850.0,
            "EL": 2825.5,
            "EC": 2845.0,
            "AO": 2850.5,
            "AH": 2853.0,
            "AL": 2826.0,
            "AC": 2829.0,
            "Vo": 42910.0,
            "OI": 479812.0,
            "Va": 1217918971856.0,
            "CM": "2024-09",
            "VoOA": 40405.0,
            "EmMrgnTrgDiv": "002",
            "LTD": "2024-09-12",
            "SQD": "2024-09-13",
            "Settle": 2829.0,
            "CCMFlag": "1"
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

