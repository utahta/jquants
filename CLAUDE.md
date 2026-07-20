# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

J-Quants API v2（日本の株式市場データを提供するAPI）のGoクライアントライブラリ。

## API仕様の確認方法

**重要**: J-Quants APIの仕様を確認する際は、以下の順序で参照してください：

1. **最初に必ず`docs/v2/`ディレクトリを確認** - 公式ドキュメントのローカルキャッシュ（gitignore対象）。存在しない・古い場合は `make docs-sync` で取得
2. **不明な点がある場合のみ**: https://jpx-jquants.com/ja/spec/ - 公式APIリファレンス（各ページはURL末尾に`.md`を付けるとMarkdownで直接取得可能。例: https://jpx-jquants.com/ja/spec/eq-bars-daily.md ）

## 開発コマンド

```bash
make help          # 利用可能なコマンドを表示
make check         # コンパイルチェック
make test          # 単体テスト
make test-cover    # カバレッジ付きテスト
make lint          # リント
make test-e2e      # E2Eテスト（JQUANTS_API_KEY が必要）
make test-run TEST=TestQuotesService_GetDailyQuotes  # 特定のテスト
make docs-sync     # 公式APIドキュメントを docs/v2/ に同期
```

## 新しいAPIサービスの追加時の注意

- `jquants.go`の`JQuantsAPI`構造体への登録を忘れない（忘れてもコンパイルは通る）
- `make docs-sync` で対象APIの最新ドキュメントを取得して参照する

## 実装上の注意点

1. **nil安全性**: APIレスポンスの欠損フィールドはポインタで表現
2. **型変換**: APIはJSONの型を不整合に返すことがあるため、`types`パッケージのカスタム型で処理
3. **定数定義**: 市場区分コード、業種コード、開示書類種別などは定数として定義済み（再定義しない）
