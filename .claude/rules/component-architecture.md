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
