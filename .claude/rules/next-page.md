---
paths:
  - apps/**/app/**/page.tsx
---

# Next.js Page Component

- 必ずサーバーコンポーネントとして作成する
- インタラクションが必要な場合は別ファイルにクライアントコンポーネントを作成して、ページコンポーネントから呼び出す
- データ取得はページでは行わず、Container コンポーネントを作成して行う
- ページコンポーネントの props は Next.js が提供するグローバル型 `PageProps<'/path'>` を使う
- ページコンポーネントで `params` / `searchParams` を `await` しない。 Promise のまま子（Container）コンポーネントへ渡す
  - 子へ渡す前に必要なパラメータだけ取り出したい場合は `params.then((p) => p.foo)` のように `.then` で変換し、変換後の `Promise<T>` を渡す
  - これにより、ページは async 関数にならず、 await による解決を Container 側に閉じ込められる
