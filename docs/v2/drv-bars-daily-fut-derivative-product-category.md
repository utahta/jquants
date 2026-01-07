# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 先物商品区分コード

| コード | 商品区分名称 | データ収録期間 | 
|------|------|------|
| TOPIXF | TOPIX先物 | 2008/5/7〜 | 
| TOPIXMF | ミニTOPIX先物 | 2008/6/16〜 | 
| MOTF | マザーズ先物 | 2016/7/19〜 | 
| NKVIF | 日経平均VI先物 | 2012/2/27〜 | 
| NKYDF | 日経平均・配当指数先物 | 2010/7/26〜 | 
| NK225F | 日経225先物 | 2008/5/7〜 | 
| NK225MF | 日経225mini先物 | 2008/5/7〜 | 
| JN400F | JPX日経インデックス400先物 | 2014/11/25〜 | 
| REITF | 東証REIT指数先物 | 2008/6/16〜 | 
| DJIAF | NYダウ先物 | 2012/5/28〜 | 
| JGBLF | 長期国債先物 | 2008/5/7〜 | 
| NK225MCF | 日経225マイクロ先物 | 2023/5/29〜 | 
| TOA3MF | TONA3ヶ月金利先物 | 2023/5/29〜 | 

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

