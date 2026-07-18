# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

J-Quants API v2（日本の株式市場データを提供するAPI）のGoクライアントライブラリ。

## API仕様の確認方法

**重要**: J-Quants APIの仕様を確認する際は、以下の順序で参照してください：

1. **最初に必ず`docs/v2/`ディレクトリを確認** - このリポジトリ内のv2 APIドキュメント
2. **不明な点がある場合のみ**: https://jpx.gitbook.io/j-quants-ja/api-reference - 公式APIリファレンス

## 開発コマンド

```bash
make help          # 利用可能なコマンドを表示
make check         # コンパイルチェック
make test          # 単体テスト
make test-cover    # カバレッジ付きテスト
make lint          # リント
make test-e2e      # E2Eテスト（JQUANTS_API_KEY が必要）
make test-run TEST=TestQuotesService_GetDailyQuotes  # 特定のテスト
```

## 新しいAPIサービスの追加時の注意

- `jquants.go`の`JQuantsAPI`構造体への登録を忘れない（忘れてもコンパイルは通る）
- `docs/v2/`ディレクトリにドキュメントを追加する

## 実装上の注意点

1. **nil安全性**: APIレスポンスの欠損フィールドはポインタで表現
2. **型変換**: APIはJSONの型を不整合に返すことがあるため、`types`パッケージのカスタム型で処理
3. **定数定義**: 市場区分コード、業種コード、開示書類種別などは定数として定義済み（再定義しない）
