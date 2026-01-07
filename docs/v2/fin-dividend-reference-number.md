# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# リファレンスナンバー

## リファレンスナンバーについて

- リファレンスナンバー：配当通知を一意に特定するための番号。
- CAリファレンスナンバー：訂正・削除の対象となっている配当通知のリファレンスナンバー。新規の場合はリファレンスナンバーと同じ値。

## 具体例：以下の通知があった場合に、提供データは下表のとおりになります。

- 銘柄：日本取引所グループ（銘柄コード：86970）について

2023-03-06　　配当が新規で通知
2023-03-07　　配当が訂正情報として通知
2023-03-08　　配当が削除された
2023-03-09　　配当が新規で通知

| PubDate | Code | RefNo | CARefNo | StatCode | 
|------|------|------|------|------|
| 2023-03-06 | 86970 | 1 | 1 | 1：新規 | 
| 2023-03-07 | 86970 | 2 | 1 | 2：訂正 | 
| 2023-03-08 | 86970 | 3 | 1 | 3：削除 | 
| 2023-03-09 | 86970 | 4 | 4 | 1：新規 | 

- 一部項目のみを抽出して例示しています。
- 上記のコード値は例示のため便宜的な記載としており、また実際に発生したデータとは異なります。

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

