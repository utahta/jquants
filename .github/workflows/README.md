# リリース手順

[`tagpr.yml`](tagpr.yml) で [tagpr](https://github.com/Songmu/tagpr) を実行し、リリース PR と GitHub Release を作成します。設定は [`.tagpr`](../../.tagpr)、リリースノートの設定は [`.github/release.yml`](../release.yml) にあります。

## 初回設定

認証には GitHub Actions が発行する `GITHUB_TOKEN` を使用します。リポジトリの Settings → Actions → General → Workflow permissions で **Allow GitHub Actions to create and approve pull requests** を有効にしてください。tagpr に必要な書き込み権限はワークフロー内の `permissions` で指定しています。

## リリースの作成

1. 通常のプルリクエストを `main` にマージすると、未リリースの変更をまとめたリリース PR が作成・更新されます。
2. リリース PR のバージョンと `CHANGELOG.md` を確認します。利用者への影響など、自動生成された変更履歴に不足する説明はこの PR で補足します。
3. CI が承認待ちになったら、書き込み権限のある利用者が PR 画面の **Approve workflows to run** を選択し、CI の結果を確認します。リリース PR が更新された場合も、承認が必要になることがあります。[GitHub の実行条件](https://docs.github.com/en/actions/how-tos/write-workflows/choose-when-workflows-run/trigger-a-workflow)
4. リリース PR をマージすると、`v` 付きのバージョンタグと公開済みの GitHub Release が作成されます。

リリース PR は、公開したいタイミングまで開いたままにできます。`main` に変更が追加されると内容が更新されます。`CHANGELOG.md` は最初のリリース PR で生成されます。

## バージョンの選択

既定は patch 更新です。minor・major 更新はリリース PR のラベルで指定します。通常の変更 PR に `minor` または `major` ラベルを付けてマージすると、対応するラベルがリリース PR に反映されます。

| リリース PR のラベル | 更新例 |
|---------------------|--------|
| 指定なし | `v0.8.0` → `v0.8.1` |
| `tagpr:minor` | `v0.8.0` → `v0.9.0` |
| `tagpr:major` | `v0.8.0` → `v1.0.0` |

開いているリリース PR の `tagpr:minor`・`tagpr:major` ラベルを手動で追加・削除すると、tagpr が再実行され、PR タイトルと `CHANGELOG.md` のバージョンが更新されます。更新が完了してから内容と CI の結果を確認してください。

通常の変更 PR のラベルから反映されたバージョン指定を取り消す場合は、先に元の変更 PR の `minor`・`major` ラベルを削除してください。リリース PR 側のラベルだけを削除すると、再実行時に付け直されます。

タグだけでバージョンを管理するため、PR タイトルを編集してもリリースバージョンは変わりません。ラベルの詳細は [tagpr のバージョン指定方法](https://github.com/Songmu/tagpr/blob/v1.20.3/docs/guides/versioning.md)を参照してください。
