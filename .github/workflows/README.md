# GitHub Actions Workflows

このディレクトリには、J-Quants Go ClientのCI/CDワークフローが含まれています。

## ワークフロー一覧

### CI (`ci.yml`)

プルリクエストやmainブランチへのプッシュ時に実行される基本的なCIワークフローです。

**トリガー:**
- Pull Request (mainブランチ向け)
- Push to main branch

**ジョブ:**
1. **Test**: ビルドチェックと単体テストの実行
   - Go 1.23を使用
   - カバレッジレポートをCodecovにアップロード
2. **Lint**: golangci-lintによるコード品質チェック
3. **Security**: Gosecによるセキュリティスキャン

### Release (`release.yml`)

新しいバージョンのリリース時に実行されるワークフローです。

**トリガー:**
- タグプッシュ（`v*`パターン）

**ジョブ:**
1. テストの実行
2. 変更履歴の生成
3. GitHubリリースの作成

## ローカルでのテスト

GitHub Actionsのワークフローをローカルでテストする場合は、[act](https://github.com/nektos/act)を使用できます：

```bash
# CIワークフローをテスト
act pull_request

# リリースワークフローをテスト
act push --eventpath event.json
```

## トラブルシューティング

### Lintエラーが発生する場合

ローカルで以下を実行して修正：

```bash
make lint
# または
golangci-lint run --fix
```