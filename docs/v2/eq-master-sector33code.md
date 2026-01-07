# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 33業種コード及び業種名

| コード | 名称 | 
|------|------|
| 0050 | 水産・農林業 | 
| 1050 | 鉱業 | 
| 2050 | 建設業 | 
| 3050 | 食料品 | 
| 3100 | 繊維製品 | 
| 3150 | パルプ・紙 | 
| 3200 | 化学 | 
| 3250 | 医薬品 | 
| 3300 | 石油･石炭製品 | 
| 3350 | ゴム製品 | 
| 3400 | ガラス･土石製品 | 
| 3450 | 鉄鋼 | 
| 3500 | 非鉄金属 | 
| 3550 | 金属製品 | 
| 3600 | 機械 | 
| 3650 | 電気機器 | 
| 3700 | 輸送用機器 | 
| 3750 | 精密機器 | 
| 3800 | その他製品 | 
| 4050 | 電気･ガス業 | 
| 5050 | 陸運業 | 
| 5100 | 海運業 | 
| 5150 | 空運業 | 
| 5200 | 倉庫･運輸関連業 | 
| 5250 | 情報･通信業 | 
| 6050 | 卸売業 | 
| 6100 | 小売業 | 
| 7050 | 銀行業 | 
| 7100 | 証券･商品先物取引業 | 
| 7150 | 保険業 | 
| 7200 | その他金融業 | 
| 8050 | 不動産業 | 
| 9050 | サービス業 | 
| 9999 | その他 | 

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

