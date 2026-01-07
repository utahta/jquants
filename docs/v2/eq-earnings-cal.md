# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 決算発表予定日(/equities/earnings-calendar)

GET    /v2/equities/earnings-calendar

## APIの概要

3月期・9月期決算の会社の決算発表予定日を取得できます。（その他の決算期の会社は今後対応予定です）

## 本APIの留意点

- 下記のサイトで、3月期・９月期決算会社分に更新があった場合のみ19時ごろに更新されます。３月期・９月期決算会社についての更新がなかった場合は、最終更新日時点のデータを提供します。
https://www.jpx.co.jp/listing/event-schedules/financial-announcement/index.html
- 本APIは翌営業日に決算発表が行われる銘柄に関する情報を返します。
- 本APIから得られたデータにおいてDateの項目が翌営業日付であるレコードが存在しない場合は、3月期・9月期決算会社における翌営業日の開示予定はないことを意味します。
- REITのデータは含まれません。

## 決算発表予定日の銘柄コード、年度、 四半期等の照会をします。

GET https://api.jquants.com/v2/equities/earnings-calendar

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

**pagination_key** (string, optional)
  検索の先頭を指定する文字列
過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/equities/earnings-calendar
-H "x-api-key: {loading}"
```

### Responses

### データ項目概要

**Date** (string)
  日付（YYYY-MM-DD）
決算発表予定日が未定の場合、空文字("")となります。

**Code** (string)
  銘柄コード

**CoName** (string)
  会社名

**FY** (string)
  決算期末

**SectorNm** (string)
  業種名

**FQ** (string)
  決算種別

**Section** (string)
  市場区分


### レスポンスサンプル

```
{
    "data": [
        {
            "Date": "2022-02-14",
            "Code": "43760",
            "CoName": "くふうカンパニー",
            "FY": "9月30日",
            "SectorNm": "情報・通信業",
            "FQ": "第１四半期",
            "Section": "マザーズ"
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

