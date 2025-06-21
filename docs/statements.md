# J-Quants API 財務諸表データ項目概要

このドキュメントは、J-Quants APIの財務諸表エンドポイント (`/fins/statements`) のレスポンスデータ項目について説明します。

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
| TypeOfDocument | 開示書類種別 | 例: 3QFinancialStatements_Consolidated_IFRS |
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
| ResultPayoutRatioAnnual | 配当性向（実績） | パーセント |

#### 配当予想
| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| ForecastDividendPerShare1stQuarter | 第1四半期配当金（予想） | 単位：円 |
| ForecastDividendPerShare2ndQuarter | 第2四半期配当金（予想） | 単位：円 |
| ForecastDividendPerShare3rdQuarter | 第3四半期配当金（予想） | 単位：円 |
| ForecastDividendPerShareFiscalYearEnd | 期末配当金（予想） | 単位：円 |
| ForecastDividendPerShareAnnual | 年間配当金（予想） | 単位：円 |
| ForecastPayoutRatioAnnual | 配当性向（予想） | パーセント |
| ForecastDistributionsPerUnit(REIT) | 分配金（REIT）（予想） | 単位：円 |

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
| NextYearForecastDividendPerShareAnnual | 翌期年間配当金（予想） | 単位：円 |
| NextYearForecastPayoutRatioAnnual | 翌期配当性向（予想） | パーセント |
| NextYearForecastNetSales | 翌期売上高（予想） | 単位：円 |
| NextYearForecastOperatingProfit | 翌期営業利益（予想） | 単位：円 |
| NextYearForecastOrdinaryProfit | 翌期経常利益（予想） | 単位：円 |
| NextYearForecastProfit | 翌期当期純利益（予想） | 単位：円 |
| NextYearForecastEarningsPerShare | 翌期1株当たり当期純利益（予想） | 単位：円 |

### 6. その他の情報

| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| MaterialChangesInSubsidiaries | 重要な子会社の異動 | true/false |
| SignificantChangesInTheScopeOfConsolidation | 連結範囲の重要な変更 | |
| ChangesOtherThanOnesBasedOnRevisionsOfAccountingStandard | 会計基準改正以外の変更 | true/false |
| NumberOfIssuedAndOutstandingSharesAtTheEndOfFiscalYearIncludingTreasuryStock | 期末発行済株式数（自己株式を含む） | |
| NumberOfTreasuryStockAtTheEndOfFiscalYear | 期末自己株式数 | |

### 7. 単体財務数値

単体（NonConsolidated）の財務数値は、連結財務数値と同様の項目が「NonConsolidated」プレフィックス付きで提供されます。

例：
- `NonConsolidatedNetSales`: 単体売上高
- `NonConsolidatedOperatingProfit`: 単体営業利益
- `NonConsolidatedOrdinaryProfit`: 単体経常利益
- `NonConsolidatedProfit`: 単体当期純利益
- など

## 使用上の注意

1. **データ型の変換**: すべてのフィールドは文字列型で返されるため、数値計算を行う場合は適切な型変換が必要です。

2. **空値の処理**: 該当するデータがない場合は空文字列（""）が返されます。

3. **配当利回りの計算**: 配当利回りを計算する場合は、別途株価データ（`/prices/daily_quotes`）を取得する必要があります。

4. **会計期間の種類**: `TypeOfCurrentPeriod`は四半期（1Q-4Q）または通期（FY）を示します。5Qは特殊な場合に使用されます。

5. **開示書類種別**: `TypeOfDocument`は会計基準（IFRS/JGAAP）や連結/単体の区別を含みます。

## 関連ドキュメント

- [J-Quants API 公式ドキュメント](https://jpx.gitbook.io/j-quants-ja/api-reference/statements)