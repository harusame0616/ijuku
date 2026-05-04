---
paths:
  - *.tsx
---

# コンポーネントアーキテクチャ

- Container / Presentational パターンを基本とする
- データフェッチ・ロジックは Container コンポーネントで行う
- Container コンポーネント内に直接データフェッチ・ロジックを書くのではなく、別ファイルに関数を切り出して Container コンポーネントから呼び出す形にする
- Container コンポーネントは Suspense でラップし、スケルトンを表示する

```tsx
<Suspense fallback={<FooSkeleton />}>
  <FooContainer />
</Suspense>
```

- shadcn/ui のコンポーネントを優先して使用する

## client component は最小化する

- ファイル全体に `"use client"` を付ける前に、 client 要因（state / event handler / browser API）を別ファイルへ切り出せないか必ず検討する
- 親コンポーネントは「子に client を含むだけ」であれば universal のまま保つ
- client component を作成した後も、 「これ以上分解して universal/server に逃がせる部分が残っていないか」を必ず再チェックする
