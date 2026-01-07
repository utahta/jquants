# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 東証信用貸借規制区分

| コード | 説明 | 
|------|------|
| 001 | 日本証券金融が実施する貸株注意喚起銘柄および貸株申込制限措置銘柄 | 
| 002 | 東京証券取引所が定める日々公表銘柄 | 
| 003 | 東京証券取引所が定める規制銘柄 | 
| 004 | 東京証券取引所が定める規制銘柄（2次規制） | 
| 005 | 東京証券取引所が定める規制銘柄（3次規制） | 
| 006 | 東京証券取引所が定める規制銘柄（4次規制） | 
| 101 | 東京証券取引所が定める規制解除銘柄 | 
| 102 | 東京証券取引所が定める監理銘柄 | 

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

