# E2E Tests

このディレクトリには、J-Quants APIの実際のエンドポイントに対するエンドツーエンド（E2E）テストが含まれています。

## 実行方法

E2Eテストは、ビルドタグ `e2e` が指定された場合のみ実行されます。

```bash
# 全てのE2Eテストを実行
go test -tags=e2e ./test/e2e -v

# 特定のテストのみ実行
go test -tags=e2e ./test/e2e -v -run TestListedInfo

# タイムアウトを設定して実行（デフォルトは10分）
go test -tags=e2e ./test/e2e -v -timeout 30s
```

## 前提条件

E2Eテストを実行する前に、以下の環境変数を設定する必要があります：

```bash
# 方法1: Email/Passwordでログイン
export JQUANTS_EMAIL="your-email@example.com"
export JQUANTS_PASSWORD="your-password"

# 方法2: リフレッシュトークンを直接指定
export JQUANTS_REFRESH_TOKEN="your-refresh-token"
```

## テストの構成

### エンドポイント別テスト構成

E2Eテストは各エンドポイントごとにファイルを分けて構成されており、レスポンスデータの完全な検証を行います。

#### テスト対象のAPIエンドポイント

- `daily_quotes_test.go` - 日次株価データ（/prices/daily_quotes）
- `listed_test.go` - 上場企業情報（/listed/info）
- `statements_test.go` - 財務諸表（/fins/statements）
- `announcement_test.go` - 決算発表予定（/fins/announcement）
- `indices_test.go` - 指数四本値（/markets/indices）
- `topix_test.go` - TOPIX詳細データ（/markets/topix）
- `trading_calendar_test.go` - 取引カレンダー（/markets/trading_calendar）
- `trades_spec_test.go` - 投資部門別売買状況（/markets/trades_spec）
- `short_selling_test.go` - 業種別空売り比率（/markets/short_selling）
- `short_selling_positions_test.go` - 空売り残高報告（/markets/short_selling_positions）
- `weekly_margin_interest_test.go` - 信用取引週末残高（/markets/weekly_margin_interest）
- `breakdown_test.go` - 売買内訳データ（/markets/breakdown）
- `prices_am_test.go` - 前場四本値（/prices/prices_am）
- `dividend_test.go` - 配当情報（/fins/dividend）
- `fs_details_test.go` - 財務諸表詳細（/fins/fs_details）
- `index_option_test.go` - 日経225オプション
- `futures_test.go` - 先物四本値（/derivatives/futures）
- `options_test.go` - オプション四本値（/derivatives/options）

※各APIの利用可能なプランについては、[J-Quants公式サイト](https://jpx-jquants.com/)でご確認ください。

### テストの特徴

1. **完全なレスポンス検証**: 各エンドポイントのレスポンスフィールドを詳細に検証
2. **データ整合性チェック**: 四本値の論理的整合性、日付範囲の確認など
3. **エラーケース対応**: サブスクリプション制限、データなしケースのハンドリング
4. **ページネーション対応**: 大量データの取得テスト
5. **統計分析**: 取得データの基本統計情報の計算と表示

利用しているサブスクリプションプランで使用できないAPIのテストは自動的にスキップされます。

## 注意事項

1. **APIレート制限**: J-Quants APIにはレート制限があります。短時間に大量のリクエストを送信しないように注意してください。

2. **営業日の考慮**: 多くのAPIは営業日のデータのみ提供します。土日祝日にテストを実行すると、データが取得できない場合があります。

3. **データの可用性**: 当日のデータは翌営業日以降に利用可能になることがあります。テストでは過去の営業日のデータを使用しています。

4. **コスト**: APIの呼び出しにはサブスクリプションプランに応じたコストがかかる可能性があります。

## デバッグ

テストが失敗した場合、以下の点を確認してください：

1. 認証情報が正しく設定されているか
2. インターネット接続が正常か
3. J-Quants APIのサービスが正常に稼働しているか
4. 指定した日付が営業日か
5. サブスクリプションプランが適切か