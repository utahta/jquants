# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 上場銘柄一覧(/equities/master)

GET    /v2/equities/master

## APIの概要

過去時点での銘柄情報、当日の銘柄情報および翌営業日時点の銘柄情報が取得可能です。
ただし、翌営業日時点の銘柄情報については17 時半以降に取得可能となります。

### 本APIの留意点

- 過去日付の指定について、Premiumプランでデータ提供開始日（2008年5月7日）より過去日付を指定した場合であっても、2008年5月7日時点の銘柄情報を返却します。
- 指定された日付が休業日の場合は、指定日の翌営業日の銘柄情報を返却します。
- 貸借信用区分および貸借信用区分名については Standard および Premium プランのみ取得可能な項目です。

2022年4月の東証市場区分再編により、日本銀行（銘柄コード83010）および信金中央金庫（銘柄コード84210）については、制度上所属する市場区分が存在しなくなりましたが、J-Quants では市場区分をスタンダードとして返却します。

## 日次の銘柄情報を取得します

GET https://api.jquants.com/v2/equities/master

データの取得では、銘柄コード（code）または日付（date）の指定が可能です。
各パラメータの組み合わせとレスポンスの結果については以下のとおりです。

| code | date | レスポンスの結果 | 
|------|------|------|
| ✗ | ✗ | APIを実行した日付時点における全銘柄情報一覧（※1） | 
| ✓ | ✗ | APIを実行した日付時点における指定された銘柄情報（※1） | 
| ✗ | ✓ | 指定日付時点における全銘柄情報の一覧（※2） | 
| ✓ | ✓ | 指定日付時点における指定された銘柄情報（※2） | 

※1 休業日において日付を指定せずにクエリした場合、直近の翌営業日における銘柄情報一覧を返却します。
※2 未来日付の指定について、Light プラン以上では翌営業日時点のデータが取得可能です。翌営業日より先の未来日付を指定した場合であっても、翌営業日時点の銘柄情報を返却します。

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

**code** (string, optional)
  銘柄コード（e.g. 27890 or 2789）4桁の銘柄コードを指定した場合は、普通株式と優先株式の両方が上場している銘柄においては普通株式のデータのみが取得されます。

**date** (string, optional)
  基準となる日付の指定（e.g. 20210907 or 2021-09-07）


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/equities/master 
-H "x-api-key: {loading}" 
-d code="86970" 
-d date="2022-11-11"
```

### Responses

### データ項目概要

**Date** (string)
  情報適用年月日（YYYY-MM-DD）

**Code** (string)
  銘柄コード

**CoName** (string)
  会社名

**CoNameEn** (string)
  会社名（英語）

**S17** (string)
  17業種コード（17業種コード及び業種名を参照）

**S17Nm** (string)
  17業種コード名（17業種コード及び業種名を参照）

**S33** (string)
  33業種コード（33業種コード及び業種名を参照）

**S33Nm** (string)
  33業種コード名（33業種コード及び業種名を参照）

**ScaleCat** (string)
  規模コード

**Mkt** (string)
  市場区分コード（市場区分コード及び市場区分を参照）

**MktNm** (string)
  市場区分名（市場区分コード及び市場区分を参照）

**Mrgn** (string)
  貸借信用区分（1: 信用 / 2: 貸借 / 3: その他）（※1）

**MrgnNm** (string)
  貸借信用区分名（※1）


※1 Standard および Premium プランで取得可能な項目です

### レスポンスサンプル

```
{
    "data": [
        {
            "Date": "2022-11-11",
            "Code": "86970",
            "CoName": "日本取引所グループ",
            "CoNameEn": "Japan Exchange Group,Inc.",
            "S17": "16",
            "S17Nm": "金融（除く銀行）",
            "S33": "7200",
            "S33Nm": "その他金融業",
            "ScaleCat": "TOPIX Large70",
            "Mkt": "0111",
            "MktNm": "プライム",
            "Mrgn": "1",
            "MrgnNm": "信用"
        }
    ]
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

