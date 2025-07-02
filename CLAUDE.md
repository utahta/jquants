# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

このリポジトリは、J-Quants API（日本の株式市場データを提供するAPI）のGoクライアントライブラリです。

## API仕様の確認方法

**重要**: J-Quants APIの仕様を確認する際は、以下の順序で参照してください：

1. **最初に必ず`docs/`ディレクトリを確認** - このリポジトリ内のドキュメント
2. **不明な点がある場合のみ**: https://jpx.gitbook.io/j-quants-ja/api-reference - 公式APIリファレンス

## 開発コマンド

### Makefileを使用した開発
```bash
# 利用可能なコマンドを表示
make help

# コンパイルチェック（構文エラーの確認）
make check

# 単体テスト
make test
make test-v        # 詳細表示
make test-cover    # カバレッジ付き

# E2Eテスト（認証情報が必要）
make test-e2e
make test-e2e-v    # 詳細表示

# リント
make lint

# 特定のテストを実行
make test-run TEST=TestQuotesService_GetDailyQuotes

# 開発ツールのインストール
make install-tools
```

### ツール

#### gitbook2md
GitBookのドキュメントをMarkdownに変換するツール。J-Quants APIの公式ドキュメントをローカルで参照する際に便利。

```bash
# ビルド
cd cmd/gitbook2md && go build

# 使用方法1: HTMLファイルから変換
./gitbook2md input.html output.md

# 使用方法2: URLから直接変換
./gitbook2md --url https://jpx.gitbook.io/j-quants-ja/api-reference/statements --output statements.md
```

### 直接実行する場合
```bash
# 依存関係の取得
go mod download

# ビルド
go build ./...

# 単体テストの実行
go test ./...

# カバレッジ付きテスト
go test -cover ./...

# 特定のテストを実行
go test -run TestQuotesService_GetDailyQuotes ./...

# E2Eテスト（実際のAPIに接続、認証情報が必要）
go test -tags=e2e ./test/e2e -v
```

### リント
```bash
# golangci-lintを使用
golangci-lint run
```

## アーキテクチャ

### ディレクトリ構造
- `auth/`: 認証処理（ログイン、リフレッシュトークン管理）
- `client/`: HTTPクライアントインターフェースとモック
- `types/`: カスタム型（JSONの不整合な型を処理）
- `docs/`: 各APIエンドポイントのドキュメント
- `test/e2e/`: エンドツーエンドテスト

### 主要な設計パターン

1. **サービス指向**: 各APIエンドポイントごとに独立したサービス構造体
   ```go
   type QuotesService struct {
       client client.HTTPClient
   }
   ```

2. **パラメータ/レスポンス構造体**: 各APIメソッドは専用のパラメータとレスポンス構造体を使用
   ```go
   type DailyQuotesParams struct {
       Code          string
       Date          string
       PaginationKey string
   }
   ```

3. **カスタムJSON型**: APIの不整合な型を処理（例：`types.Float64String`）

4. **ページネーション**: 大量データ取得時の自動ページネーション処理

### 新しいAPIサービスの追加方法

1. 新しいサービスファイルを作成（例：`new_service.go`）
2. サービス構造体とメソッドを実装
3. `jquants.go`の`Client`構造体に新サービスを追加
4. 対応するテストファイルを作成（例：`new_service_test.go`）
5. `docs/`ディレクトリにドキュメントを追加

### テスト実装のパターン

```go
// モッククライアントを使用したテスト
func TestService_Method(t *testing.T) {
    mockClient := &client.MockClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            // モックレスポンスを返す
        },
    }
    
    service := &Service{client: mockClient}
    // テスト実行
}
```

### 認証の扱い

- Email/Passwordログイン: `auth.Login()`
- リフレッシュトークン: `auth.RefreshAccessToken()`
- 認証情報は環境変数または設定ファイルで管理
  - `JQUANTS_EMAIL`, `JQUANTS_PASSWORD`
  - `JQUANTS_REFRESH_TOKEN`
  - `~/.jquants/refresh_token`

### 重要な実装上の注意点

1. **nil安全性**: APIレスポンスの欠損フィールドはポインタで表現
2. **エラーハンドリング**: 一貫したエラーラップとコンテキスト情報の付与
3. **型変換**: `types`パッケージのカスタム型を使用してJSONの不整合を処理
4. **営業日**: 日本の株式市場の営業日を考慮したデータ取得
5. **定数定義**: 市場区分コード、業種コード、開示書類種別などは定数として定義済み
   - 市場区分: `MarketPrime`, `MarketStandard`, `MarketGrowth` など
   - 業種コード: `Sector17Food`, `Sector33IT` など
   - 開示書類種別: `TypeOfDocumentFYConsolidatedJP`, `TypeOfDocumentFYConsolidatedIFRS` など