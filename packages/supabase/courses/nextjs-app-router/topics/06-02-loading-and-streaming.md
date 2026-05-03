---
title: "loading.tsx と Suspense でストリーミングする"
description: "loading.tsx を追加してページ単位のスケルトンを出し、Suspense で粒度の細かいストリーミングを行う。"
---

# loading.tsx と Suspense でストリーミングする

## 目標

- `app/posts/loading.tsx` を追加し、`/posts` 遷移直後にスケルトンを出す
- `<Suspense>` を使って、ページ内の特定箇所だけを段階的にストリーミングする
- 「静的シェル」と「ストリーミングされる動的部分」の関係を説明できる

## 知識

**ストリーミング** とは、サーバーが HTML を一度に送るのではなく、`<Suspense>` の境界に合わせてチャンクごとに少しずつ送る方式です。ユーザーは静的部分（レイアウトやナビゲーション、フォールバック UI）を即座に見られ、データが必要な箇所はあとから差し替わって表示されます。

App Router でストリーミングを使う方法は 2 つあります。

- **`loading.tsx`**: セグメントに置くと、Next.js が自動でその `page.tsx` を `<Suspense>` で囲み、指定したコンポーネントを **フォールバック UI** にします。ページ単位の "全体スケルトン" として手軽
- **`<Suspense>`**: コンポーネント単位で囲い、独立したストリーミング境界を細かく作れる。複数の境界があれば、解決した順に独立して表示されていく

```tsx
// 自動 Suspense ラッピング
export default function Loading() {
  return <p>読み込み中...</p>;
}
```

```tsx
// 粒度の細かい境界
<Suspense fallback={<PostListSkeleton />}>
  <PostList />
</Suspense>
```

最初に届く HTML を **静的シェル** と呼びます。Cache Components の文脈では、`'use cache'` で囲った部分と Suspense のフォールバック UI が静的シェルに含まれ、`<Suspense>` 内の動的コンポーネントだけが実行時にストリーミングで埋まります。これが Next.js 16 のデフォルトレンダリングモデルである **Partial Prerendering (PPR)** の中身です。

`loading.tsx` が便利なのは「ナビゲーション直後に何かしら見せたい」とき。一方で、ページ内で複数のデータを並行して取りたいときは `<Suspense>` を使い、遅いコンポーネントだけが個別にスケルトンを出すようにしたほうが体験が良くなります。

> 注意: レイアウトが非キャッシュなデータ（`cookies()` や非キャッシュ `fetch` など）にアクセスしている場合、同じセグメントの `loading.tsx` ではフォールバックされません。本コースのように Cache Components を有効にしている場合、レイアウト内の非キャッシュ・ランタイムデータアクセスは **個別の `<Suspense>` で明示的に囲う必要があり**、囲っていないとビルド時エラーで誘導されます。`<Suspense>` で囲っておけば静的シェルが先に流れ、非キャッシュ部分は後から埋まります。データ取得をページ側に押し下げるのも有効な選択肢です。

## ステップ

### 1. /posts 用の loading.tsx を作る

`app/posts/loading.tsx` を作成し、シンプルなスケルトンを返します。

```tsx
export default function PostsLoading() {
  return (
    <main>
      <h1>記事一覧</h1>
      <p>読み込み中...</p>
    </main>
  );
}
```

別のページから `/posts` に遷移してみてください。一覧が描画されるまでの一瞬、このスケルトンが見えるはずです。

### 2. PostList 用のスケルトンを差し替える

前章で `app/posts/page.tsx` の `PostList` を `<Suspense fallback={<p>読み込み中...</p>}>` で囲っていました。フォールバックを少しリッチなスケルトンにしてみます。

```tsx
function PostListSkeleton() {
  return (
    <ul>
      <li>───────</li>
      <li>───────</li>
      <li>───────</li>
    </ul>
  );
}

// page 内
<Suspense fallback={<PostListSkeleton />}>
  <PostList />
</Suspense>;
```

### 3. 動的セグメントにも loading を入れる

`app/posts/[slug]/loading.tsx` を作成してください。動的ルートでは「サーバー応答待ち」が発生しやすいため、`loading.tsx` を置くと体感速度が大きく改善します。

```tsx
export default function PostLoading() {
  return <p>記事を読み込み中...</p>;
}
```

## 完了判定

- `app/posts/loading.tsx` と `app/posts/[slug]/loading.tsx` の両方が存在する
- 一覧ページ内の `<Suspense>` のフォールバックがスケルトン UI になっている
- 別ページから `/posts` や `/posts/<slug>` に遷移したとき、最初にスケルトンが見える

## 補足

`loading.tsx` は内部的に `page.tsx` を `<Suspense>` で包む糖衣構文です。そのため、ページの中で同じ役割の `<Suspense>` を別途置いてもうまく協調します。動的ルートに `loading.tsx` が無いと、リンクをホバーしても部分プリフェッチが効かず、クリック後に「無反応に見える」時間が生じやすくなります。これを避けるためにも、`/posts/[slug]` のように動的セグメントには `loading.tsx` を置く習慣を付けましょう。スケルトンは中身を完全に再現する必要は無く、レイアウトの形が崩れない最低限のプレースホルダーで十分です。

## 理解度チェック

- `loading.tsx` と `<Suspense>` のスケルトンの粒度の違い、それぞれ向いている使い分けを説明してください
