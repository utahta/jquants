# V1 API から V2 API への変更点 - Protocol API Reference

> J-Quants API V1 から V2 への主な仕様変更点（認証方式、プラン内容、レートリミット、エンドポイント、レスポンス形式）について説明します。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点認証認可プラン・データ提供範囲レートリミットエンドポイント・パラメータの変更レスポンス形式プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# V1 API から V2 API への変更点

J-Quants API V2 では、使いやすさの改善等を目的として、認証方式を含むいくつかの重要な仕様変更が行われています。
V1 API をご利用中のお客様は、以下の変更点をご確認の上、V2 API への移行をお願いいたします。

V1仕様書はこちら

## 認証認可

認証方式が「トークン方式」から「APIキー方式」に変更されました。

| 項目 | V1 API | V2 API | 
|------|------|------|
| API利用方法 | token/auth_user 等で ID Token / Refresh Token を発行して利用 | ダッシュボードから発行した APIキー (x-api-key ヘッダー) を利用 | 
| 認可の期限 | ID Token / Refresh Token に有効期限あり | APIキー自体には有効期限なし（再発行・削除は可能） | 

## プラン・データ提供範囲

| 項目 | V1 API | V2 API | 
|------|------|------|
| Premiumプランの期間制限 | 無制限 | 過去20年分 まで | 
| 上場銘柄一覧（貸借信用区分） | Standard, Premium プランのみ取得可能 | 全プラン で取得可能 | 

## レートリミット

プランごとに API リクエスト数の上限（レートリミット）が設定されました。

| プラン | 上限 (リクエスト / 分) | 
|------|------|
| Free | 5 | 
| Light | 60 | 
| Standard | 120 | 
| Premium | 500 | 

## エンドポイント・パラメータの変更

V1 API から V2 API への移行に伴い、エンドポイントのパスとパラメータを変更しています。

### エンドポイントの対応表

| データセット | V1 エンドポイント | V2 エンドポイント | 
|------|------|------|
| トークン発行 | /v1/token/auth_user | 廃止 (APIキーを使用) | 
| トークンリフレッシュ | /v1/token/auth_refresh | 廃止 (APIキーを使用) | 
| 株価四本値 | /v1/prices/daily_quotes | /v2/equities/bars/daily | 
| 前場四本値 | /v1/prices/prices_am | /v2/equities/bars/daily/am | 
| 決算発表予定日 | /v1/fins/announcement | /v2/equities/earnings-calendar | 
| 投資部門別情報 | /v1/markets/trades_spec | /v2/equities/investor-types | 
| 上場銘柄一覧 | /v1/listed/info | /v2/equities/master | 
| 先物四本値 | /v1/derivatives/futures | /v2/derivatives/bars/daily/futures | 
| オプション四本値 | /v1/derivatives/options | /v2/derivatives/bars/daily/options | 
| 日経225オプション四本値 | /v1/option/index_option | /v2/derivatives/bars/daily/options/225 | 
| 売買内訳データ | /v1/markets/breakdown | /v2/markets/breakdown | 
| 取引カレンダー | /v1/markets/trading_calendar | /v2/markets/calendar | 
| 日々公表信用取引残高 | /v1/markets/daily_margin_interest | /v2/markets/margin-alert | 
| 信用取引週末残高 | /v1/markets/weekly_margin_interest | /v2/markets/margin-interest | 
| 業種別空売り比率 | /v1/markets/short_selling | /v2/markets/short-ratio | 
| 空売り残高報告 | /v1/markets/short_selling_positions | /v2/markets/short-sale-report | 
| 指数四本値 | /v1/indices | /v2/indices/bars/daily | 
| TOPIX指数四本値 | /v1/indices/topix | /v2/indices/bars/daily/topix | 
| 財務諸表(BS/PL/CF) | /v1/fins/fs_details | /v2/fins/details | 
| 財務情報 | /v1/fins/statements | /v2/fins/summary | 
| 配当金情報 | /v1/fins/dividend | /v2/fins/dividend | 

## レスポンス形式

| 項目 | V1 API | V2 API | 
|------|------|------|
| レスポンス構造 | APIによって異なる | 原則としてデータを "data" キーの配列として返却 | 

### レスポンス例

```
{
  "data": [
    { ... },
    { ... }
  ],
  "pagination_key": "..."
}
```

### カラム名の変更例（株価四本値）

V2 API では、レスポンスのカラム名が短縮形に変更されている場合があります。以下は株価四本値の例です。

| 項目 | V1 API カラム名 | V2 API カラム名 | 
|------|------|------|
| 日付 | Date | Date | 
| 銘柄コード | Code | Code | 
| 始値 | Open | O | 
| 高値 | High | H | 
| 安値 | Low | L | 
| 終値 | Close | C | 
| 出来高 | Volume | Vo | 
| 売買代金 | TurnoverValue | Va | 
| 調整後始値 | AdjustmentOpen | AdjO | 
| 調整後高値 | AdjustmentHigh | AdjH | 
| 調整後安値 | AdjustmentLow | AdjL | 
| 調整後終値 | AdjustmentClose | AdjC | 
| 調整後出来高 | AdjustmentVolume | AdjVo | 
| 調整係数 | AdjustmentFactor | AdjFactor | 

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

