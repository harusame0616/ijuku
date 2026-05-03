---
title: "use cache で記事一覧をキャッシュする"
description: "'use cache' ディレクティブと cacheLife / cacheTag を使い、記事一覧を静的シェルに含める。"
---

# use cache で記事一覧をキャッシュする

## 目標

- `'use cache'` ディレクティブで Server Component の結果をキャッシュする
- `cacheLife('hours')` でキャッシュ寿命を指定する
- `cacheTag('posts')` でタグを付け、後の章でタグ単位の無効化を可能にする

## 知識

`'use cache'` は、関数または Server Component の **戻り値をキャッシュ** するディレクティブです。`cacheComponents: true` が有効になっていることが前提で、関数本体の先頭に `'use cache'` を置くと、その関数の引数と外側のクロージャ値を **キャッシュキー** にして結果が保存されます。同じ引数なら次回からは保存済みの結果を即座に返します。

```tsx
import { cacheLife, cacheTag } from "next/cache";

async function PostList() {
  "use cache";
  cacheLife("hours");
  cacheTag("posts");

  const posts = await getPosts();
  return <ul>{/* ... */}</ul>;
}
```

`'use cache'` には次の補助 API が組み合わせて使えます。

- **`cacheLife(profile)`**: キャッシュ寿命のプロファイルを指定。`seconds` / `minutes` / `hours` / `days` / `weeks` / `max` から選ぶ。各プロファイルには `stale` / `revalidate` / `expire` の 3 値が決まっており、例えば `'hours'` は `stale: 5m, revalidate: 1h, expire: 1d` という意味
- **`cacheTag(tag)`**: 任意のタグを付与する。あとで `revalidateTag` や `updateTag` を呼ぶと、そのタグが付いたキャッシュが無効化される

「データ単位」でキャッシュしたいときはデータ取得関数に `'use cache'` を、「UI 単位」でキャッシュしたいときはコンポーネント全体に `'use cache'` を付けます。本コースではミニブログの記事一覧という UI ブロックをまとめてキャッシュするため、コンポーネント単位で適用します。

キャッシュ対象に出来ない代表例として、`Math.random()` や `Date.now()`、`crypto.randomUUID()` のような **非決定的** な値、`cookies()` / `headers()` / `searchParams` / `params` のような **リクエスト時 API** があります。これらをキャッシュ内で使うとビルド時にエラーになります。それらを扱うコンポーネントは `<Suspense>` で囲い、キャッシュ外に追い出してください。

## ステップ

### 1. PostList をキャッシュ対象にする

`app/posts/page.tsx` の `PostList` コンポーネントの先頭に `'use cache'` を加え、`cacheLife('hours')` と `cacheTag('posts')` を呼びます。

```tsx
import { Suspense } from "react";
import { cacheLife, cacheTag } from "next/cache";
import Link from "next/link";
import { getPosts } from "./data";
import { LikeButton } from "./_components/like-button";
import { PostForm } from "./_components/post-form";

async function PostList() {
  "use cache";
  cacheLife("hours");
  cacheTag("posts");

  const posts = await getPosts();
  return (
    <ul>
      {posts.map((post) => (
        <li key={post.slug}>
          <Link href={`/posts/${post.slug}`}>{post.title}</Link>
          <LikeButton initialLikes={0} />
        </li>
      ))}
    </ul>
  );
}
```

### 2. フォームはキャッシュの外に置く

投稿フォームの Client Component（`PostForm`）はページ直下にそのまま置き、`PostList` のキャッシュとは独立させます。これにより「投稿フォームと一覧」が同じページにあっても、一覧だけがキャッシュ済みのまま再利用されます。

### 3. ビルドして静的シェルを確認する

開発を一旦止めて `pnpm build` を実行してください。出力に `/posts` が **Prerendered** として表示されれば、一覧が静的シェルに含まれた証拠です。`Uncached data was accessed outside of <Suspense>` というエラーが出る場合は、キャッシュ内で禁止されているリクエスト時 API を使っている合図なので、対象の処理を `<Suspense>` で外に出してください。

## 完了判定

- `PostList` コンポーネントの先頭に `'use cache'` がある
- `cacheLife('hours')` と `cacheTag('posts')` の両方が呼ばれている
- `pnpm build` でビルドが成功する

## 補足

`'use cache'` はファイル先頭にも書けます。その場合、ファイル内のエクスポート関数はすべて async 関数である必要があり、それらがキャッシュ対象になります。`cacheLife` で `seconds` プロファイル、`revalidate: 0`、または `expire` が 5 分未満の値を指定すると **短命キャッシュ** と判定され、自動的にプリレンダリングから除外されて動的レンダリングになります。意図せず短命にならないよう注意してください。キャッシュキーは引数とクロージャ値から自動生成されるため、関数の引数を変えると別のキャッシュエントリになります。

## 理解度チェック

- `cacheLife` における `stale` / `revalidate` / `expire` の 3 つはそれぞれ何を表しますか
- `'use cache'` の関数の中でリクエスト時 API（`cookies` / `headers` 等）を直接使えないのはなぜですか。その制約を回避してリクエスト固有のデータを使う典型的なパターンも説明してください
