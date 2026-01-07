# Protocol API Reference - J-Quants

> プロが活用する株価や財務データを、APIで簡単に取得。あなたの投資分析に、確かな情報を。

- Dashboard
- Support

- GuidesJ-Quants APIについてV1 API から V2 API への変更点プランごとに利用可能なAPIとデータ格納期間提供データの更新タイミングクイックスタートMCPサーバー
- Noticesデータ修正履歴・制約事項レスポンスのページングについてレートリミットについてAPIレスポンスのGzip化レスポンスステータス
- Resources上場銘柄一覧(/equities/master)17業種コード及び業種名33業種コード及び業種名市場区分コード及び市場区分名株価四本値(/equities/bars/daily)前場四本値(/equities/bars/daily/am)投資部門別情報(/equities/investor-types)市場名信用取引週末残高(/markets/margin-interest)業種別空売り比率(/markets/short-ratio)空売り残高報告(/markets/short-sale-report)日々公表信用取引残高(/markets/margin-alert)公表の理由東証信用貸借規制区分売買内訳データ(/markets/breakdown)取引カレンダー(/markets/calendar)休日区分指数四本値(/indices/bars/daily)配信対象指数コードTOPIX指数四本値(/indices/bars/daily/topix)財務情報(/fins/summary)開示書類種別財務諸表(BS/PL/CF)(/fins/details)配当金情報(/fins/dividend)リファレンスナンバー決算発表予定日(/equities/earnings-calendar)日経225オプション四本値(/derivatives/bars/daily/options/225)先物四本値(/derivatives/bars/daily/futures)先物商品区分コードオプション四本値(/derivatives/bars/daily/options)オプション商品区分コード
- Sign in

# 財務情報(/fins/summary)

GET    /v2/fins/summary

## APIの概要

財務情報の開示日・銘柄コードを指定して、決算短信などの財務情報サマリーを取得できます。
銘柄コード（code）または日付（date）のいずれか一方、もしくは両方の指定が必要です。

## 本APIの留意点

- 会計基準について： APIから出力される各項目名は日本基準（JGAAP）の開示項目が基準となっています。そのため、IFRSや米国基準（USGAAP）の開示データにおいては、経常利益の概念がありませんので、データが空欄となっています。

- 四半期開示見直し対応に伴うAPI項目の追加について：

四半期開示見直し対応において、決算短信サマリー様式の記載事項が以下のとおり変更されます。

変更前： 重要な⼦会社の異動（連結範囲の変更を伴う特定⼦会社の異動）
変更後： 連結範囲の重要な変更


この対応に伴い、2024/7/22より本APIのレスポンス項目に"SignificantChangesInTheScopeOfConsolidation"（期中における連結範囲の重要な変更）を追加いたします。
詳細は、データ項目概要欄をご覧ください。

## 財務情報データを取得します

GET https://api.jquants.com/v2/fins/summary

銘柄コード（code）または日付（date）の指定が必須となります。

### パラメータ及びレスポンス

銘柄コード（code）または日付（date）の指定が必須となります。
各パラメータの組み合わせとレスポンスの結果については以下のとおりです。

| code | date | レスポンスの結果 | 
|------|------|------|
| ✓ | ✗ | 指定された銘柄について全期間分の財務情報データ | 
| ✓ | ✓ | 指定された銘柄について指定された日付の財務情報データ | 
| ✗ | ✓ | 全上場銘柄について指定された日付の財務情報データ | 

### Requests

### Headers

**x-api-key** (string, required)
  APIキー


### Query Parameters

code または date のいずれか一つの指定が必須です。

**code** (string, optional)
  銘柄コード（e.g. 86970 or 8697）
4桁もしくは5桁の銘柄コード

**date** (string, optional)
  開示日付の指定（e.g. 2022-01-05 or 20220105）

**pagination_key** (string, optional)
  検索の先頭を指定する文字列
過去の検索で返却された pagination_key を設定


### APIコールサンプルコード

### Request

```
curl -G https://api.jquants.com/v2/fins/summary 
-H "x-api-key: {loading}" 
-d code="86970" 
-d date="2024-10-25"
```

### Responses

### データ項目概要

**DiscDate** (string)
  開示日

**DiscTime** (string)
  開示時刻

**Code** (string)
  銘柄コード（5桁）

**DiscNo** (string)
  開示番号
APIから出力されるjsonは開示番号で昇順に並んでいます。

**DocType** (string)
  開示書類種別
開示書類種別一覧

**CurPerType** (string)
  当会計期間の種類
[1Q, 2Q, 3Q, 4Q, 5Q, FY]

**CurPerSt** (string)
  当会計期間開始日

**CurPerEn** (string)
  当会計期間終了日

**CurFYSt** (string)
  当事業年度開始日

**CurFYEn** (string)
  当事業年度終了日

**NxtFYSt** (string)
  翌事業年度開始日
開示レコードに翌事業年度の開示情報がない場合空欄になります。

**NxtFYEn** (string)
  翌事業年度終了日
開示レコードに翌事業年度の開示情報がない場合空欄になります。

**Sales** (number)
  売上高

**OP** (number)
  営業利益

**OdP** (number)
  経常利益

**NP** (number)
  当期純利益

**EPS** (number)
  一株あたり当期純利益

**DEPS** (number)
  潜在株式調整後一株あたり当期純利益

**TA** (number)
  総資産

**Eq** (number)
  純資産

**EqAR** (number)
  自己資本比率

**BPS** (number)
  一株あたり純資産

**CFO** (number)
  営業活動によるキャッシュ・フロー

**CFI** (number)
  投資活動によるキャッシュ・フロー

**CFF** (number)
  財務活動によるキャッシュ・フロー

**CashEq** (number)
  現金及び現金同等物期末残高

**Div1Q** (number)
  一株あたり配当実績_第1四半期末

**Div2Q** (number)
  一株あたり配当実績_第2四半期末

**Div3Q** (number)
  一株あたり配当実績_第3四半期末

**DivFY** (number)
  一株あたり配当実績_期末

**DivAnn** (number)
  一株あたり配当実績_合計

**DivUnit** (number)
  1口当たり分配金

**DivTotalAnn** (number)
  配当金総額

**PayoutRatioAnn** (number)
  配当性向

**FDiv1Q** (number)
  一株あたり配当予想_第1四半期末

**FDiv2Q** (number)
  一株あたり配当予想_第2四半期末

**FDiv3Q** (number)
  一株あたり配当予想_第3四半期末

**FDivFY** (number)
  一株あたり配当予想_期末

**FDivAnn** (number)
  一株あたり配当予想_合計

**FDivUnit** (number)
  1口当たり予想分配金

**FDivTotalAnn** (number)
  予想配当金総額

**FPayoutRatioAnn** (number)
  予想配当性向

**NxFDiv1Q** (number)
  一株あたり配当予想_翌事業年度第1四半期末

**NxFDiv2Q** (number)
  一株あたり配当予想_翌事業年度第2四半期末

**NxFDiv3Q** (number)
  一株あたり配当予想_翌事業年度第3四半期末

**NxFDivFY** (number)
  一株あたり配当予想_翌事業年度期末

**NxFDivAnn** (number)
  一株あたり配当予想_翌事業年度合計

**NxFDivUnit** (number)
  1口当たり翌事業年度予想分配金

**NxFPayoutRatioAnn** (number)
  翌事業年度予想配当性向

**FSales2Q** (number)
  売上高_予想_第2四半期末

**FOP2Q** (number)
  営業利益_予想_第2四半期末

**FOdP2Q** (number)
  経常利益_予想_第2四半期末

**FNP2Q** (number)
  当期純利益_予想_第2四半期末

**FEPS2Q** (number)
  一株あたり当期純利益_予想_第2四半期末

**NxFSales2Q** (number)
  売上高_予想_翌事業年度第2四半期末

**NxFOP2Q** (number)
  営業利益_予想_翌事業年度第2四半期末

**NxFOdP2Q** (number)
  経常利益_予想_翌事業年度第2四半期末

**NxFNp2Q** (number)
  当期純利益_予想_翌事業年度第2四半期末

**NxFEPS2Q** (number)
  一株あたり当期純利益_予想_翌事業年度第2四半期末

**FSales** (number)
  売上高_予想_期末

**FOP** (number)
  営業利益_予想_期末

**FOdP** (number)
  経常利益_予想_期末

**FNP** (number)
  当期純利益_予想_期末

**FEPS** (number)
  一株あたり当期純利益_予想_期末

**NxFSales** (number)
  売上高_予想_翌事業年度期末

**NxFOP** (number)
  営業利益_予想_翌事業年度期末

**NxFOdP** (number)
  経常利益_予想_翌事業年度期末

**NxFNp** (number)
  当期純利益_予想_翌事業年度期末

**NxFEPS** (number)
  一株あたり当期純利益_予想_翌事業年度期末

**MatChgSub** (string)
  期中における重要な子会社の異動

**SigChgInC** (string)
  期中における連結範囲の重要な変更
*指定されたdateが2024-07-21以前のレスポンスは、当該項目には値が収録されません。

**ChgByASRev** (string)
  会計基準等の改正に伴う会計方針の変更

**ChgNoASRev** (string)
  会計基準等の改正に伴う変更以外の会計方針の変更

**ChgAcEst** (string)
  会計上の見積りの変更

**RetroRst** (string)
  修正再表示

**ShOutFY** (number)
  期末発行済株式数

**TrShFY** (number)
  期末自己株式数

**AvgSh** (number)
  期中平均株式数

**NCSales** (number)
  売上高_非連結

**NCOP** (number)
  営業利益_非連結

**NCOdP** (number)
  経常利益_非連結

**NCNP** (number)
  当期純利益_非連結

**NCEPS** (number)
  一株あたり当期純利益_非連結

**NCTA** (number)
  総資産_非連結

**NCEq** (number)
  純資産_非連結

**NCEqAR** (number)
  自己資本比率_非連結

**NCBPS** (number)
  一株あたり純資産_非連結

**FNCSales2Q** (number)
  売上高_予想_第2四半期末_非連結

**FNCOP2Q** (number)
  営業利益_予想_第2四半期末_非連結

**FNCOdP2Q** (number)
  経常利益_予想_第2四半期末_非連結

**FNCNP2Q** (number)
  当期純利益_予想_第2四半期末_非連結

**FNCEPS2Q** (number)
  一株あたり当期純利益_予想_第2四半期末_非連結

**NxFNCSales2Q** (number)
  売上高_予想_翌事業年度第2四半期末_非連結

**NxFNCOP2Q** (number)
  営業利益_予想_翌事業年度第2四半期末_非連結

**NxFNCOdP2Q** (number)
  経常利益_予想_翌事業年度第2四半期末_非連結

**NxFNCNP2Q** (number)
  当期純利益_予想_翌事業年度第2四半期末_非連結

**NxFNCEPS2Q** (number)
  一株あたり当期純利益_予想_翌事業年度第2四半期末_非連結

**FNCSales** (number)
  売上高_予想_期末_非連結

**FNCOP** (number)
  営業利益_予想_期末_非連結

**FNCOdP** (number)
  経常利益_予想_期末_非連結

**FNCNP** (number)
  当期純利益_予想_期末_非連結

**FNCEPS** (number)
  一株あたり当期純利益_予想_期末_非連結

**NxFNCSales** (number)
  売上高_予想_翌事業年度期末_非連結

**NxFNCOP** (number)
  営業利益_予想_翌事業年度期末_非連結

**NxFNCOdP** (number)
  経常利益_予想_翌事業年度期末_非連結

**NxFNCNP** (number)
  当期純利益_予想_翌事業年度期末_非連結

**NxFNCEPS** (number)
  一株あたり当期純利益_予想_翌事業年度期末_非連結


### レスポンスサンプル

```
{
    "data": [
        {
            "DiscDate": "2023-01-30",
            "DiscTime": "12:00:00",
            "Code": "86970",
            "DiscNo": "20230127594871",
            "DocType": "3QFinancialStatements_Consolidated_IFRS",
            "CurPerType": "3Q",
            "CurPerSt": "2022-04-01",
            "CurPerEn": "2022-12-31",
            "CurFYSt": "2022-04-01",
            "CurFYEn": "2023-03-31",
            "NxtFYSt": "",
            "NxtFYEn": "",
            "Sales": "100529000000",
            "OP": "51765000000",
            "OdP": "",
            "NP": "35175000000",
            "EPS": "66.76",
            "DEPS": "",
            "TA": "79205861000000",
            "Eq": "320021000000",
            "EqAR": "0.004",
            "BPS": "",
            "CFO": "",
            "CFI": "",
            "CFF": "",
            "CashEq": "91135000000",
            "Div1Q": "",
            "Div2Q": "26.0",
            "Div3Q": "",
            "DivFY": "",
            "DivAnn": "",
            "DivUnit": "",
            "DivTotalAnn": "",
            "PayoutRatioAnn": "",
            "FDiv1Q": "",
            "FDiv2Q": "",
            "FDiv3Q": "",
            "FDivFY": "36.0",
            "FDivAnn": "62.0",
            "FDivUnit": "",
            "FDivTotalAnn": "",
            "FPayoutRatioAnn": "",
            "NxFDiv1Q": "",
            "NxFDiv2Q": "",
            "NxFDiv3Q": "",
            "NxFDivFY": "",
            "NxFDivAnn": "",
            "NxFDivUnit": "",
            "NxFPayoutRatioAnn": "",
            "FSales2Q": "",
            "FOP2Q": "",
            "FOdP2Q": "",
            "FNP2Q": "",
            "FEPS2Q": "",
            "NxFSales2Q": "",
            "NxFOP2Q": "",
            "NxFOdP2Q": "",
            "NxFNp2Q": "",
            "NxFEPS2Q": "",
            "FSales": "132500000000",
            "FOP": "65500000000",
            "FOdP": "",
            "FNP": "45000000000",
            "FEPS": "85.42",
            "NxFSales": "",
            "NxFOP": "",
            "NxFOdP": "",
            "NxFNp": "",
            "NxFEPS": "",
            "MatChgSub": "false",
            "SigChgInC": "",
            "ChgByASRev": "false",
            "ChgNoASRev": "false",
            "ChgAcEst": "true",
            "RetroRst": "",
            "ShOutFY": "528578441",
            "TrShFY": "1861043",
            "AvgSh": "526874759",
            "NCSales": "",
            "NCOP": "",
            "NCOdP": "",
            "NCNP": "",
            "NCEPS": "",
            "NCTA": "",
            "NCEq": "",
            "NCEqAR": "",
            "NCBPS": "",
            "FNCSales2Q": "",
            "FNCOP2Q": "",
            "FNCOdP2Q": "",
            "FNCNP2Q": "",
            "FNCEPS2Q": "",
            "NxFNCSales2Q": "",
            "NxFNCOP2Q": "",
            "NxFNCOdP2Q": "",
            "NxFNCNP2Q": "",
            "NxFNCEPS2Q": "",
            "FNCSales": "",
            "FNCOP": "",
            "FNCOdP": "",
            "FNCNP": "",
            "FNCEPS": "",
            "NxFNCSales": "",
            "NxFNCOP": "",
            "NxFNCOdP": "",
            "NxFNCNP": "",
            "NxFNCEPS": ""
        }
    ],
    "pagination_key": "value1.value2."
}
```

Was this page helpful?

© JPX Market Innovation & Research, Inc. All rights reserved.

