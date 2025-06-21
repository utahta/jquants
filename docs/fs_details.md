# J-Quants API 財務諸表(BS/PL)データ項目概要

このドキュメントは、J-Quants APIの財務諸表(BS/PL)エンドポイント (`/fins/fs_details`) のレスポンスデータ項目について説明します。

## データ型について
- 日付・時刻フィールドは `String` 型で返されます
- 財務諸表の各項目は `Map` 型で、キーが冗長ラベル（英語）、値が `String` 型で返されます
- 金額データは円単位で提供されます（特定銘柄を除く）

## エンドポイント概要

### 財務諸表(BS/PL) (/fins/fs_details)

#### 概要
- **エンドポイント**: `https://api.jquants.com/v1/fins/fs_details`
- **HTTPメソッド**: GET
- **説明**: 上場企業の四半期毎の財務情報における、貸借対照表、損益計算書に記載の項目を取得することができます。

#### 本APIの特徴
- **EDINET XBRLベース**: EDINETタクソノミ本体のlabel情報を使用
- **会計基準対応**: 日本基準とIFRSの両方に対応
- **四半期データ**: 四半期ごとの詳細な財務諸表データを提供

#### 本APIの留意点

##### FinancialStatement（財務諸表の各種項目）について
- **データソース**: EDINET XBRLタクソノミ本体（label情報）を用いてコンテンツを作成
- **冗長ラベル参照先**: https://disclosure2dl.edinet-fsa.go.jp/guide/static/disclosure/WZEK0110.html
- **会計基準別データ提供**:
  - **日本基準**: 「勘定科目リスト」のE列「冗長ラベル（英語）」をキーとして提供
  - **IFRS**: 「国際会計基準タクソノミ要素リスト」のD列「冗長ラベル（英語）」をキーとして提供

##### 提出者別タクソノミについて
- EDINETタクソノミには存在しない提出者別タクソノミで定義される企業独自の項目は、本APIの提供対象外

##### 特殊なケース
- **三井海洋開発（銘柄コード62690）**: 2022年2月以降の決算短信の連結財務諸表及び連結財務諸表注記を米ドルにより表示。当該銘柄の財務諸表情報は米ドルでの提供

#### リクエストパラメータ

リクエストパラメータにcode（銘柄コード）またはdate（開示日）を入力する必要があります。

| パラメータ名 | 型 | 必須 | 説明 | 例 |
|------------|-----|------|------|-----|
| code | String | いいえ* | 4桁もしくは5桁の銘柄コード | 86970 or 8697 |
| date | String | いいえ* | 開示日 | 2022-01-05 or 20220105 |
| pagination_key | String | いいえ | 検索の先頭を指定する文字列 | 過去の検索で返却されたpagination_key |

*codeまたはdateのいずれかが必須

#### リクエストヘッダー

| ヘッダー名 | 型 | 必須 | 説明 |
|----------|-----|------|------|
| Authorization | String | はい | アクセスキー |

#### レスポンスフィールド

##### 基本情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| DisclosedDate | 開示日 | String | YYYY-MM-DD形式 |
| DisclosedTime | 開示時刻 | String | HH:MM:SS形式 |
| LocalCode | 銘柄コード（5桁） | String | |
| DisclosureNumber | 開示番号 | String | APIから出力されるjsonは開示番号で昇順に並んでいる |
| TypeOfDocument | 開示書類種別 | String | 開示書類種別一覧参照 |

##### 財務諸表情報

| フィールド名 | 日本語説明 | 型 | 備考 |
|------------|----------|-----|------|
| FinancialStatement | 財務諸表の各種項目 | Map | 冗長ラベル（英語）とその値のマップ |

## FinancialStatementの構造

FinancialStatementは、以下のような構造を持つMapオブジェクトです：
- **キー**: 冗長ラベル（英語） - XBRLタグと紐づく標準化された項目名
- **値**: 財務諸表の値（String型）

### 主要な項目例（日本基準）
- 資産の部：総資産、流動資産、固定資産等
- 負債の部：総負債、流動負債、固定負債等
- 純資産の部：株主資本、資本金、利益剰余金等
- 損益計算書：売上高、営業利益、経常利益、当期純利益等

### 主要な項目例（IFRS）
サンプルレスポンスに含まれる項目例：
- **資産**: `Assets (IFRS)`, `Current assets (IFRS)`, `Non-current assets (IFRS)`
- **負債**: `Liabilities (IFRS)`, `Current liabilities (IFRS)`, `Non-current liabilities (IFRS)`
- **資本**: `Equity (IFRS)`, `Share capital (IFRS)`, `Retained earnings (IFRS)`
- **損益**: `Revenue - 2 (IFRS)`, `Operating profit (loss) (IFRS)`, `Profit (loss) (IFRS)`
- **その他重要項目**: `Goodwill (IFRS)`, `Basic earnings (loss) per share (IFRS)`

## 会計基準について

### 日本基準とIFRSの識別
- **項目名の末尾**: IFRSの場合、項目名に `(IFRS)` が付与される
- **Accounting standards, DEI**: 会計基準を示すフィールド（"IFRS"または"JapaneseGAAP"）

### EDINETタクソノミとの対応
- **日本基準**: EDINETタクソノミの「勘定科目リスト」と対応
- **IFRS**: EDINETタクソノミの「国際会計基準タクソノミ要素リスト」と対応

## 開示書類種別（TypeOfDocument）

財務諸表のタイプを示す重要なフィールドです。詳細は開示書類種別一覧を参照してください。

主な種別例：
- `1QFinancialStatements_Consolidated_IFRS`: 第1四半期連結財務諸表（IFRS）
- `2QFinancialStatements_Consolidated_IFRS`: 第2四半期連結財務諸表（IFRS）
- `3QFinancialStatements_Consolidated_IFRS`: 第3四半期連結財務諸表（IFRS）
- `YearEndFinancialStatements_Consolidated_IFRS`: 通期連結財務諸表（IFRS）

## 使用上の注意事項

1. **企業独自項目の除外**
   - 提出者別タクソノミで定義される企業独自の項目は提供対象外
   - 標準化されたEDINETタクソノミ項目のみを提供

2. **通貨単位**
   - 通常は円単位で提供
   - 三井海洋開発（62690）は2022年2月以降米ドル単位

3. **データの完全性**
   - すべての企業がすべての項目を開示しているわけではない
   - 開示されていない項目は含まれない

4. **冗長ラベルの理解**
   - 冗長ラベルは標準化された項目名
   - 同じ概念でも会計基準により異なるラベルを使用

## サンプルレスポンス

```json
{
    "fs_details": [
        {
            "DisclosedDate": "2023-01-30",
            "DisclosedTime": "12:00:00",
            "LocalCode": "86970",
            "DisclosureNumber": "20230127594871",
            "TypeOfDocument": "3QFinancialStatements_Consolidated_IFRS",
            "FinancialStatement": {
                "Goodwill (IFRS)": "67374000000",
                "Retained earnings (IFRS)": "263894000000",
                "Operating profit (loss) (IFRS)": "51765000000.0",
                "Previous fiscal year end date, DEI": "2022-03-31",
                "Basic earnings (loss) per share (IFRS)": "66.76",
                "Document type, DEI": "四半期第３号参考様式　[IFRS]（連結）",
                "Current period end date, DEI": "2022-12-31",
                "Revenue - 2 (IFRS)": "100987000000.0",
                "Profit (loss) attributable to owners of parent (IFRS)": "35175000000.0",
                "Current liabilities (IFRS)": "78852363000000",
                "Equity attributable to owners of parent (IFRS)": "311103000000",
                "Non-current liabilities (IFRS)": "33476000000",
                "Property, plant and equipment (IFRS)": "11277000000",
                "Cash and cash equivalents (IFRS)": "91135000000",
                "Share capital (IFRS)": "11500000000",
                "Assets (IFRS)": "79205861000000",
                "Equity (IFRS)": "320021000000",
                "Liabilities (IFRS)": "78885839000000",
                "Accounting standards, DEI": "IFRS"
            }
        }
    ],
    "pagination_key": "value1.value2."
}
```

## 活用例

### 特定銘柄の最新財務諸表を取得
```
/fins/fs_details?code=86970&date=20230130
```

### 特定日の全銘柄財務諸表を取得
```
/fins/fs_details?date=2023-01-30
```

### 財務分析の実施
```javascript
// ROE（自己資本利益率）の計算例
const netIncome = parseFloat(fs["Profit (loss) attributable to owners of parent (IFRS)"]);
const equity = parseFloat(fs["Equity attributable to owners of parent (IFRS)"]);
const roe = (netIncome / equity) * 100;

// 流動比率の計算例
const currentAssets = parseFloat(fs["Current assets (IFRS)"]);
const currentLiabilities = parseFloat(fs["Current liabilities (IFRS)"]);
const currentRatio = currentAssets / currentLiabilities;
```

### 時系列財務分析
複数期間の財務諸表を取得し、売上高成長率、利益率の推移等を分析

### 会計基準別の処理
```javascript
// 会計基準の判定
const accountingStandards = fs["Accounting standards, DEI"];
if (accountingStandards === "IFRS") {
    // IFRS用の処理
} else {
    // 日本基準用の処理
}
```

## 関連情報

このAPIは詳細な財務諸表データを提供するものです。関連する他のAPIについては以下を参照してください：
- 財務情報：`/fins/statements`（財務サマリー情報）
- 配当金情報：`/fins/dividend`（配当データ）
- 決算発表予定日：`/fins/announcement`（決算発表スケジュール）
- 株価四本値：`/prices/daily_quotes`（株価データ）