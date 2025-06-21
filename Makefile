.PHONY: help test test-v test-cover test-e2e test-e2e-v lint clean install-tools check

# デフォルトターゲット
help:
	@echo "利用可能なコマンド:"
	@echo "  make check        - コンパイルチェック（ビルドはしない）"
	@echo "  make test         - 単体テストを実行"
	@echo "  make test-v       - 単体テストを詳細表示で実行"
	@echo "  make test-cover   - カバレッジ付きでテストを実行"
	@echo "  make test-e2e     - E2Eテストを実行（認証情報が必要）"
	@echo "  make test-e2e-v   - E2Eテストを詳細表示で実行"
	@echo "  make lint         - golangci-lintでコードをチェック"
	@echo "  make clean        - テスト成果物をクリーン"
	@echo "  make install-tools - 開発ツールをインストール"

# コンパイルチェック（ビルドはしない）
check:
	go build -o /dev/null ./...

# 単体テスト
test:
	go test ./...

# 単体テスト（詳細表示）
test-v:
	go test -v ./...

# カバレッジ付きテスト
test-cover:
	go test -cover ./...
	@echo ""
	@echo "HTMLカバレッジレポートを生成する場合:"
	@echo "  go test -coverprofile=coverage.out ./..."
	@echo "  go tool cover -html=coverage.out"

# E2Eテスト
test-e2e:
	go test -tags=e2e ./test/e2e

# E2Eテスト（詳細表示）
test-e2e-v:
	go test -tags=e2e ./test/e2e -v

# リント
lint:
	@if ! which golangci-lint > /dev/null; then \
		echo "golangci-lintがインストールされていません。'make install-tools'を実行してください。"; \
		exit 1; \
	fi
	golangci-lint run

# クリーン
clean:
	go clean ./...
	rm -f coverage.out

# 開発ツールのインストール
install-tools:
	@echo "golangci-lintをインストール中..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "開発ツールのインストールが完了しました"

# 全テスト実行（単体テスト + リント）
test-all: lint test

# 特定のテストを実行するためのターゲット
# 使用例: make test-run TEST=TestQuotesService_GetDailyQuotes
test-run:
	go test -run $(TEST) ./...

# ベンチマークテスト
bench:
	go test -bench=. ./...

# 依存関係の更新
mod-tidy:
	go mod tidy

# 依存関係のダウンロード
mod-download:
	go mod download