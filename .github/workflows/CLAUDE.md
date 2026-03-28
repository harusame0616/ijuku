# GitHub Action 実装ガイドライン

GitHub Actions ワークフローは以下のルールに従うこと。

## 必須設定

### ステップ名の日本語化

各ステップに日本語でわかりやすい name をつける

### pnpm バージョン

`pnpm/action-setup` で version を省略する

### actions/setup-node

node-version-file で .node-version を指定する

### concurrency 設定

原則的に concurrency を設定し、重複実行を防ぐ
設定しない場合はその理由をコメントで残す

```yaml
concurrency:
  group: <workflow-name>-${{ github.ref }}
  cancel-in-progress: true
```

### タイムアウト設定

各ジョブには最低 10 分のタイムアウトを設定する

### Permission 設定

ワークフローレベルで permission ブロックを設定し、最小限の権限を付与する
