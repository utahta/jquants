# J-Quants API 決算発表予定日データ項目概要

このドキュメントは、J-Quants APIの決算発表予定日エンドポイント (`/fins/announcement`) のレスポンスデータ項目について説明します。

## データ型について
- すべてのフィールドは `String` 型で返されます
- 空文字列 (`""`) は値が存在しないことを示します（決算発表予定日が未定の場合など）

## エンドポイント概要

### 決算発表予定日 (/fins/announcement)

#### 概要
- **エンドポイント**: `https://api.jquants.com/v1/fins/announcement`
- **HTTPメソッド**: GET
- **説明**: 翌日発表予定の決算情報が取得できます。3月期・9月期決算の会社の決算発表予定日を取得できます。（その他の決算期の会社は今後対応予定です）

#### 本APIの留意点
- 下記のサイトで、3月期・9月期決算会社分に更新があった場合のみ19時ごろに更新されます
  - https://www.jpx.co.jp/listing/event-schedules/financial-announcement/index.html
- 3月期・9月期決算会社についての更新がなかった場合は、最終更新日時点のデータを提供します
- 本APIは翌営業日に決算発表が行われる銘柄に関する情報を返します
- 本APIから得られたデータにおいてDateの項目が翌営業日付であるレコードが存在しない場合は、3月期・9月期決算会社における翌営業日の開示予定はないことを意味します
- REITのデータは含まれません

#### リクエストパラメータ

| パラメータ名 | 型 | 必須 | 説明 | 例 |
|------------|-----|------|------|-----|
| pagination_key | String | いいえ | 検索の先頭を指定する文字列 | 過去の検索で返却されたpagination_key |

#### リクエストヘッダー

| ヘッダー名 | 型 | 必須 | 説明 |
|----------|-----|------|------|
| Authorization | String | はい | アクセスキー |

#### レスポンスフィールド

| フィールド名 | 日本語説明 | 備考 |
|------------|----------|------|
| Date | 日付 | YYYY-MM-DD形式。決算発表予定日が未定の場合、空文字("") |
| Code | 銘柄コード | |
| CompanyName | 会社名 | |
| FiscalYear | 決算期末 | 例："9月30日" |
| SectorName | 業種名 | |
| FiscalQuarter | 決算種別 | 例："第１四半期" |
| Section | 市場区分 | 例："マザーズ" |

## 使用上の注意事項

1. **対象企業について**
   - 現在は3月期・9月期決算会社のみ対応
   - その他の決算期の会社は今後対応予定
   - REITは対象外

2. **データ更新について**
   - 東証の公式サイトで更新があった場合のみ19時ごろに更新
   - 更新がない場合は最終更新日時点のデータを提供

3. **データの解釈**
   - 翌営業日の決算発表予定のみを返す
   - Dateが翌営業日でないレコードがない場合、翌営業日の開示予定なし

4. **決算種別について**
   - "第１四半期"、"第２四半期"、"第３四半期"、"通期" などの値を取る

## サンプルレスポンス

```json
{
  "announcement": [
    {
      "Date": "2022-02-14",
      "Code": "43760",
      "CompanyName": "くふうカンパニー",
      "FiscalYear": "9月30日",
      "SectorName": "情報・通信業",
      "FiscalQuarter": "第１四半期",
      "Section": "マザーズ"
    }
  ],
  "pagination_key": "value1.value2."
}
```

## 関連情報

このAPIは決算発表予定日の情報を提供するものです。関連する他のAPIについては以下を参照してください：
- 財務情報：`/fins/statements`
- 財務諸表（BS/PL）：`/fins/fs_details`
- 配当金情報：`/fins/dividend`
- 上場銘柄一覧：`/listed/info`

## データソース

本APIのデータは以下のサイトの情報を基に提供されています：
- [東証 決算発表予定日一覧](https://www.jpx.co.jp/listing/event-schedules/financial-announcement/index.html)