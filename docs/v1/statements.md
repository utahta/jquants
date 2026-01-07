# J-Quants API 財務諸表データ項目概要

このドキュメントは、J-Quants APIの財務諸表エンドポイント (`/fins/statements`) のレスポンスデータ項目について説明します。

## APIエンドポイント

- **URL**: `https://api.jquants.com/v1/fins/statements`
- **メソッド**: GET
- **認証**: Authorization ヘッダーにBearer トークンが必要

## リクエストパラメータ

| パラメータ名 | 型 | 必須 | 説明 | 例 |
|------------|-----|------|------|-----|
| code | String | いずれか必須 | 4桁または5桁の銘柄コード | 86970 or 8697 |
| date | String | いずれか必須 | 開示日 | 2022-01-05 or 20220105 |
| pagination_key | String | いいえ | 検索の先頭を指定する文字列（過去の検索で返却されたpagination_keyを設定） | - |

※ codeまたはdateのいずれかのパラメータが必須

## データ型について
- すべてのフィールドは `String` 型で返されます
- 数値データも文字列として格納されています
- 空文字列 (`""`) は値が存在しないことを示します
- 金額の単位は基本的に円です

## データ項目一覧

### 1. 基本情報

| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| DisclosedDate | 開示日 | YYYY-MM-DD形式 |
| DisclosedTime | 開示時刻 | HH:MM:SS形式 |
| LocalCode | 銘柄コード（5桁） | |
| DisclosureNumber | 開示番号 | APIから出力されるjsonは開示番号で昇順に並んでいる |
| TypeOfDocument | 開示書類種別 | 詳細は[開示書類種別一覧](#開示書類種別一覧)を参照 |
| TypeOfCurrentPeriod | 当会計期間の種類 | [1Q, 2Q, 3Q, 4Q, 5Q, FY] |
| CurrentPeriodStartDate | 当会計期間開始日 | |
| CurrentPeriodEndDate | 当会計期間終了日 | |
| CurrentFiscalYearStartDate | 当事業年度開始日 | |
| CurrentFiscalYearEndDate | 当事業年度終了日 | |
| NextFiscalYearStartDate | 翌事業年度開始日 | |
| NextFiscalYearEndDate | 翌事業年度終了日 | |

### 2. 連結財務数値

#### 損益計算書関連
| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| NetSales | 売上高 | 単位：円 |
| OperatingProfit | 営業利益 | 単位：円 |
| OrdinaryProfit | 経常利益 | 単位：円 |
| Profit | 当期純利益 | 単位：円 |

#### 貸借対照表関連
| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| TotalAssets | 総資産 | 単位：円 |
| Equity | 純資産 | 単位：円 |

#### キャッシュフロー関連
| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| CashAndEquivalents | 現金及び現金同等物 | 単位：円 |
| CashFlowsFromOperatingActivities | 営業活動によるキャッシュフロー | 単位：円 |
| CashFlowsFromInvestingActivities | 投資活動によるキャッシュフロー | 単位：円 |
| CashFlowsFromFinancingActivities | 財務活動によるキャッシュフロー | 単位：円 |

### 3. 財務指標

| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| EarningsPerShare | 1株当たり当期純利益（EPS） | 単位：円 |
| DilutedEarningsPerShare | 潜在株式調整後1株当たり当期純利益 | 単位：円 |
| BookValuePerShare | 1株当たり純資産（BPS） | 単位：円 |
| EquityToAssetRatio | 自己資本比率 | パーセント |

### 4. 配当情報

#### 配当実績
| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| ResultDividendPerShare1stQuarter | 第1四半期配当金（実績） | 単位：円 |
| ResultDividendPerShare2ndQuarter | 第2四半期配当金（実績） | 単位：円 |
| ResultDividendPerShare3rdQuarter | 第3四半期配当金（実績） | 単位：円 |
| ResultDividendPerShareFiscalYearEnd | 期末配当金（実績） | 単位：円 |
| ResultDividendPerShareAnnual | 年間配当金（実績） | 単位：円 |
| DistributionsPerUnit(REIT) | 1口当たり分配金 | 単位：円 |
| ResultTotalDividendPaidAnnual | 配当金総額 | 単位：円 |
| ResultPayoutRatioAnnual | 配当性向（実績） | パーセント |

#### 配当予想
| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| ForecastDividendPerShare1stQuarter | 第1四半期配当金（予想） | 単位：円 |
| ForecastDividendPerShare2ndQuarter | 第2四半期配当金（予想） | 単位：円 |
| ForecastDividendPerShare3rdQuarter | 第3四半期配当金（予想） | 単位：円 |
| ForecastDividendPerShareFiscalYearEnd | 期末配当金（予想） | 単位：円 |
| ForecastDividendPerShareAnnual | 年間配当金（予想） | 単位：円 |
| ForecastDistributionsPerUnit(REIT) | 分配金（REIT）（予想） | 単位：円 |
| ForecastTotalDividendPaidAnnual | 予想配当金総額 | 単位：円 |
| ForecastPayoutRatioAnnual | 配当性向（予想） | パーセント |

### 5. 業績予想

#### 当期業績予想
| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| ForecastNetSales | 売上高（予想） | 単位：円 |
| ForecastOperatingProfit | 営業利益（予想） | 単位：円 |
| ForecastOrdinaryProfit | 経常利益（予想） | 単位：円 |
| ForecastProfit | 当期純利益（予想） | 単位：円 |
| ForecastEarningsPerShare | 1株当たり当期純利益（予想） | 単位：円 |

#### 第2四半期業績予想
| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| ForecastNetSales2ndQuarter | 第2四半期売上高（予想） | 単位：円 |
| ForecastOperatingProfit2ndQuarter | 第2四半期営業利益（予想） | 単位：円 |
| ForecastOrdinaryProfit2ndQuarter | 第2四半期経常利益（予想） | 単位：円 |
| ForecastProfit2ndQuarter | 第2四半期純利益（予想） | 単位：円 |
| ForecastEarningsPerShare2ndQuarter | 第2四半期1株当たり純利益（予想） | 単位：円 |

#### 翌期業績予想
| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| NextYearForecastDividendPerShare1stQuarter | 翌期第1四半期配当金（予想） | 単位：円 |
| NextYearForecastDividendPerShare2ndQuarter | 翌期第2四半期配当金（予想） | 単位：円 |
| NextYearForecastDividendPerShare3rdQuarter | 翌期第3四半期配当金（予想） | 単位：円 |
| NextYearForecastDividendPerShareFiscalYearEnd | 翌期期末配当金（予想） | 単位：円 |
| NextYearForecastDividendPerShareAnnual | 翌期年間配当金（予想） | 単位：円 |
| NextYearForecastDistributionsPerUnit(REIT) | 翌期分配金（REIT）（予想） | 単位：円 |
| NextYearForecastPayoutRatioAnnual | 翌期配当性向（予想） | パーセント |
| NextYearForecastNetSales2ndQuarter | 翌期第2四半期売上高（予想） | 単位：円 |
| NextYearForecastOperatingProfit2ndQuarter | 翌期第2四半期営業利益（予想） | 単位：円 |
| NextYearForecastOrdinaryProfit2ndQuarter | 翌期第2四半期経常利益（予想） | 単位：円 |
| NextYearForecastProfit2ndQuarter | 翌期第2四半期純利益（予想） | 単位：円 |
| NextYearForecastEarningsPerShare2ndQuarter | 翌期第2四半期1株当たり純利益（予想） | 単位：円 |
| NextYearForecastNetSales | 翌期売上高（予想） | 単位：円 |
| NextYearForecastOperatingProfit | 翌期営業利益（予想） | 単位：円 |
| NextYearForecastOrdinaryProfit | 翌期経常利益（予想） | 単位：円 |
| NextYearForecastProfit | 翌期当期純利益（予想） | 単位：円 |
| NextYearForecastEarningsPerShare | 翌期1株当たり当期純利益（予想） | 単位：円 |

### 6. その他の情報

| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| MaterialChangesInSubsidiaries | 重要な子会社の異動 | true/false |
| SignificantChangesInTheScopeOfConsolidation | 連結範囲の重要な変更 | 2024/7/22より追加 ※1 |
| ChangesBasedOnRevisionsOfAccountingStandard | 会計基準等の改正に伴う会計方針の変更 | true/false |
| ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard | 会計基準改正以外の変更 | true/false |
| ChangesInAccountingEstimates | 会計上の見積りの変更 | true/false |
| RetrospectiveRestatement | 修正再表示 | |
| NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock | 期末発行済株式数（自己株式を含む） | |
| NumberOfTreasuryStockAtTheEndOfFiscalYear | 期末自己株式数 | |
| AverageNumberOfShares | 期中平均株式数 | |

### 7. 単体財務数値

単体（NonConsolidated）の財務数値は、連結財務数値と同様の項目が「NonConsolidated」プレフィックス付きで提供されます。

例：
- `NonConsolidatedNetSales`: 単体売上高
- `NonConsolidatedOperatingProfit`: 単体営業利益
- `NonConsolidatedOrdinaryProfit`: 単体経常利益
- `NonConsolidatedProfit`: 単体当期純利益
- など

## 会計基準について

APIから出力される各項目名は日本基準（JGAAP）の開示項目が基準となっています。そのため：

- **IFRS**や**米国基準（USGAAP）**の開示データでは、経常利益（OrdinaryProfit）の概念が存在しないため、データが空欄となります
- 各企業の採用する会計基準は`TypeOfDocument`フィールドで確認できます

## 四半期開示見直し対応について

2024/7/22より、決算短信サマリー様式の記載事項変更に伴い、以下の項目が追加されました：

- **変更前**: 重要な子会社の異動（連結範囲の変更を伴う特定子会社の異動）
- **変更後**: 連結範囲の重要な変更

新項目: `SignificantChangesInTheScopeOfConsolidation`（期中における連結範囲の重要な変更）
※1: 2024-07-21以前のデータには値が含まれません

## 使用上の注意

1. **データ型の変換**: すべてのフィールドは文字列型で返されるため、数値計算を行う場合は適切な型変換が必要です。

2. **空値の処理**: 該当するデータがない場合は空文字列（""）が返されます。

3. **配当利回りの計算**: 配当利回りを計算する場合は、別途株価データ（`/prices/daily_quotes`）を取得する必要があります。

4. **会計期間の種類**: `TypeOfCurrentPeriod`は四半期（1Q-4Q）または通期（FY）を示します。5Qは特殊な場合に使用されます。

5. **開示書類種別**: `TypeOfDocument`は会計基準（IFRS/JGAAP）や連結/単体の区別を含みます。

## エラーレスポンス

### 400 Bad Request
```json
{
    "message": "This API requires at least 1 parameter as follows; 'date','code'."
}
```

### 401 Unauthorized
```json
{
    "message": "The incoming token is invalid or expired."
}
```

### 500 Internal Server Error
```json
{
    "message": "Unexpected error. Please try again later."
}
```

### データサイズエラー
```json
{
    "message": "Response data is too large. Specify parameters to reduce the acquired data range."
}
```

## APIコールサンプル

### cURL
```bash
idToken=<YOUR idToken> && curl https://api.jquants.com/v1/fins/statements?code=86970&date=20230130 -H "Authorization: Bearer $idToken"
```

### Python
```python
import requests
import json

idToken = "YOUR idToken"
headers = {'Authorization': 'Bearer {}'.format(idToken)}
r = requests.get("https://api.jquants.com/v1/fins/statements?code=86970&date=20230130", headers=headers)
r.json()
```

## 開示書類種別一覧

財務情報APIのTypeOfDocumentフィールドに格納される値の一覧です。

### 決算短信（通期）
| 書類種別 | 概要 |
|---------|------|
| FYFinancialStatements_Consolidated_JP | 決算短信（連結・日本基準） |
| FYFinancialStatements_Consolidated_US | 決算短信（連結・米国基準） |
| FYFinancialStatements_NonConsolidated_JP | 決算短信（非連結・日本基準） |
| FYFinancialStatements_Consolidated_JMIS | 決算短信（連結・ＪＭＩＳ） |
| FYFinancialStatements_NonConsolidated_IFRS | 決算短信（非連結・ＩＦＲＳ） |
| FYFinancialStatements_Consolidated_IFRS | 決算短信（連結・ＩＦＲＳ） |
| FYFinancialStatements_NonConsolidated_Foreign | 決算短信（非連結・外国株） |
| FYFinancialStatements_Consolidated_Foreign | 決算短信（連結・外国株） |
| FYFinancialStatements_Consolidated_REIT | 決算短信（REIT） |

### 四半期決算短信
#### 第1四半期
| 書類種別 | 概要 |
|---------|------|
| 1QFinancialStatements_Consolidated_JP | 第1四半期決算短信（連結・日本基準） |
| 1QFinancialStatements_Consolidated_US | 第1四半期決算短信（連結・米国基準） |
| 1QFinancialStatements_NonConsolidated_JP | 第1四半期決算短信（非連結・日本基準） |
| 1QFinancialStatements_Consolidated_JMIS | 第1四半期決算短信（連結・ＪＭＩＳ） |
| 1QFinancialStatements_NonConsolidated_IFRS | 第1四半期決算短信（非連結・ＩＦＲＳ） |
| 1QFinancialStatements_Consolidated_IFRS | 第1四半期決算短信（連結・ＩＦＲＳ） |
| 1QFinancialStatements_NonConsolidated_Foreign | 第1四半期決算短信（非連結・外国株） |
| 1QFinancialStatements_Consolidated_Foreign | 第1四半期決算短信（連結・外国株） |

#### 第2四半期
| 書類種別 | 概要 |
|---------|------|
| 2QFinancialStatements_Consolidated_JP | 第2四半期決算短信（連結・日本基準） |
| 2QFinancialStatements_Consolidated_US | 第2四半期決算短信（連結・米国基準） |
| 2QFinancialStatements_NonConsolidated_JP | 第2四半期決算短信（非連結・日本基準） |
| 2QFinancialStatements_Consolidated_JMIS | 第2四半期決算短信（連結・ＪＭＩＳ） |
| 2QFinancialStatements_NonConsolidated_IFRS | 第2四半期決算短信（非連結・ＩＦＲＳ） |
| 2QFinancialStatements_Consolidated_IFRS | 第2四半期決算短信（連結・ＩＦＲＳ） |
| 2QFinancialStatements_NonConsolidated_Foreign | 第2四半期決算短信（非連結・外国株） |
| 2QFinancialStatements_Consolidated_Foreign | 第2四半期決算短信（連結・外国株） |

#### 第3四半期
| 書類種別 | 概要 |
|---------|------|
| 3QFinancialStatements_Consolidated_JP | 第3四半期決算短信（連結・日本基準） |
| 3QFinancialStatements_Consolidated_US | 第3四半期決算短信（連結・米国基準） |
| 3QFinancialStatements_NonConsolidated_JP | 第3四半期決算短信（非連結・日本基準） |
| 3QFinancialStatements_Consolidated_JMIS | 第3四半期決算短信（連結・ＪＭＩＳ） |
| 3QFinancialStatements_NonConsolidated_IFRS | 第3四半期決算短信（非連結・ＩＦＲＳ） |
| 3QFinancialStatements_Consolidated_IFRS | 第3四半期決算短信（連結・ＩＦＲＳ） |
| 3QFinancialStatements_NonConsolidated_Foreign | 第3四半期決算短信（非連結・外国株） |
| 3QFinancialStatements_Consolidated_Foreign | 第3四半期決算短信（連結・外国株） |

#### その他の四半期
| 書類種別 | 概要 |
|---------|------|
| OtherPeriodFinancialStatements_Consolidated_JP | その他四半期決算短信（連結・日本基準） |
| OtherPeriodFinancialStatements_Consolidated_US | その他四半期決算短信（連結・米国基準） |
| OtherPeriodFinancialStatements_NonConsolidated_JP | その他四半期決算短信（非連結・日本基準） |
| OtherPeriodFinancialStatements_Consolidated_JMIS | その他四半期決算短信（連結・ＪＭＩＳ） |
| OtherPeriodFinancialStatements_NonConsolidated_IFRS | その他四半期決算短信（非連結・ＩＦＲＳ） |
| OtherPeriodFinancialStatements_Consolidated_IFRS | その他四半期決算短信（連結・ＩＦＲＳ） |
| OtherPeriodFinancialStatements_NonConsolidated_Foreign | その他四半期決算短信（非連結・外国株） |
| OtherPeriodFinancialStatements_Consolidated_Foreign | その他四半期決算短信（連結・外国株） |

### 予想修正
| 書類種別 | 概要 |
|---------|------|
| DividendForecastRevision | 配当予想の修正 |
| EarnForecastRevision | 業績予想の修正 |
| REITDividendForecastRevision | 分配予想の修正（REIT） |
| REITEarnForecastRevision | 利益予想の修正（REIT） |

### 注記
- **JP**: 日本基準（JGAAP）
- **US**: 米国基準（USGAAP）
- **IFRS**: 国際財務報告基準
- **JMIS**: 修正国際基準
- **REIT**: 不動産投資信託
- **Foreign**: 外国株

## 関連ドキュメント

- [J-Quants API 公式ドキュメント](https://jpx.gitbook.io/j-quants-ja/api-reference/statements)