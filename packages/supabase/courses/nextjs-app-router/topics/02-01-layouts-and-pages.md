---
title: "Layout と Page で記事一覧と記事詳細を作る"
description: "/posts と /posts/[slug] を実装し、共通ヘッダーをルートレイアウトに置く。"
---

# Layout と Page で記事一覧と記事詳細を作る

## 目標

- `/posts` で固定配列から記事一覧を表示する
- `/posts/[slug]` で動的セグメントから個別記事を表示する
- ルートレイアウトに共通ヘッダーを置く

## 知識

App Router では、フォルダがそのまま URL セグメントになり、`page.tsx` がそのセグメントの公開 UI を定義します。フォルダを入れ子にすれば URL も入れ子になり、`[slug]` のように角括弧で囲むと **動的セグメント** になります。動的セグメントの値は、ページコンポーネントに渡される `params` プロップから取り出せます。`params` は **Promise** で渡されるため、`await` してから利用する点に注意してください。

```tsx
export default async function Page({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  // ...
}
```

レイアウトはページに **共有 UI** を提供します。ルートレイアウト（`app/layout.tsx`）はアプリ全体を包み、`<html>` と `<body>` を持つことが必須です。レイアウトは `children` プロップを通じて配下のページや子レイアウトを差し込み、ナビゲーション中も再マウントされません。

本コースでは `cacheComponents: true` を有効にしているため、`params` のように **リクエスト時にしか確定しないデータ**（ほかに `cookies` / `headers` / `searchParams` も該当）にアクセスする箇所は、`<Suspense>` で囲んでフォールバックを表示しながらストリーミングするか、`'use cache'` でキャッシュ対象にするか、のどちらかで明示的に扱います。本トピックでは前者の Suspense パターンを使い、Page を「`<Suspense>` を返す外側の同期コンポーネント」と「`params` を await して記事を表示する内側の async コンポーネント」に分離します。これにより静的シェル（プレースホルダ）が即座に届き、記事内容はリクエスト時にストリーミングで埋まります。

サーバーコンポーネントなので、ページやレイアウトの中で直接 JavaScript の配列や同期的な計算を使うこともできます。ミニブログでは、まずデータベースを使わず、TypeScript の配列リテラルで記事データを管理してスタートします。

## ステップ

### 1. ダミーデータを用意する

`app/posts/data.ts` を作成し、記事データの配列を定義します。`page.tsx` / `route.ts` 以外のファイルはルーティング対象にならないので、ページ専用のデータをルートと同じフォルダに置く（colocation）ことができます。

```ts
export type Post = {
  slug: string;
  title: string;
  body: string;
};

export const posts: Post[] = [
  { slug: "hello", title: "Hello, Mini Blog", body: "最初の記事です。" },
  { slug: "next", title: "Next.js を学ぶ", body: "App Router の入門中。" },
];
```

### 2. 一覧ページを作る

`app/posts/page.tsx` を作成し、上記 `posts` を一覧表示します。

```tsx
import { posts } from "./data";

export default function PostsPage() {
  return (
    <main>
      <h1>記事一覧</h1>
      <ul>
        {posts.map((post) => (
          <li key={post.slug}>{post.title}</li>
        ))}
      </ul>
    </main>
  );
}
```

`http://localhost:3000/posts` にアクセスし、2 件の記事タイトルが見えることを確認してください。

### 3. 動的ルートで記事詳細を作る

`app/posts/[slug]/page.tsx` を作成します。Page 関数は `<Suspense>` を返す同期コンポーネントとして書き、`params` を await する処理は内側の async コンポーネント `PostContent` に分離します。

```tsx
import { Suspense } from "react";
import { posts } from "../data";

async function PostContent({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const post = posts.find((p) => p.slug === slug);
  if (!post) {
    return <p>記事が見つかりません</p>;
  }
  return (
    <article>
      <h1>{post.title}</h1>
      <p>{post.body}</p>
    </article>
  );
}

export default function PostPage({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  return (
    <Suspense fallback={<p>読み込み中...</p>}>
      <PostContent params={params} />
    </Suspense>
  );
}
```

`http://localhost:3000/posts/hello` で本文が見えれば成功です。

### 4. ルートレイアウトに共通ヘッダーを置く

`app/layout.tsx` の `<body>` 内に、サイト名を出すヘッダーを追加してください。次のトピックで `<Link>` を入れる前提の土台です。

```tsx
<body>
  <header>
    <strong>Mini Blog</strong>
  </header>
  {children}
</body>
```

## 完了判定

- `/posts` でタイトルの一覧が表示される
- `/posts/hello` と `/posts/next` の両方で記事本文が表示される
- 全ページの上部に「Mini Blog」というヘッダーが表示される

## 補足

`params` を `await` し忘れると、TypeScript エラーや実行時の警告が出ます。Next.js 14 以前の `params` / `searchParams` は同期的な値でしたが、Next.js 15 以降は Promise 型で渡されるようになりました（`async`/`await` または React の `use` 関数で取り出します）。動的セグメントのフォルダ名は実際の URL に展開されたあとも `[slug]` のままファイルシステム上に残ります。混乱しないよう、フォルダ名と URL の対応をエディタのタブ名で意識してください。

Cache Components 有効時に runtime data を扱う箇所を `<Suspense>` で囲み忘れたり `'use cache'` を付け忘れると、ビルド時または開発サーバーで `Uncached data was accessed outside of <Suspense>` エラーが出ます。エラーが出たら「どこで `params` / `cookies` / `headers` / `searchParams` などに触っているか」を確認し、その境界の外側に `<Suspense>` を置くか、`'use cache'` を付けて対処します。なお `generateStaticParams` でビルド時に動的セグメントの全パターンを列挙すれば params が静的扱いになり Suspense は不要になりますが、本コースのミニブログは投稿で slug が動的に増える前提なので Suspense パターンで対応します。

**コロケーション（colocation）** とは、特定のルートでだけ使うコンポーネントやヘルパー、データを、そのルートと同じフォルダに置く配置スタイルのことです。`page.tsx` / `route.ts` 以外のファイルは routing に関与しないため、`app/posts/data.ts` のようにルートと同居させても URL として公開される心配はなく、Next.js の公式ドキュメントも「`app/` 配下のプロジェクトファイルは安全に colocation できる」と明記しています。アプリ全体で使う共通 UI は `src/components/` のように `app/` の外にまとめ、特定ルート専用のものは colocation する、という使い分けが見通しを良くするコツです。

Next.js のドキュメントは、プロジェクト構成について **ノンオピニオン**（`Next.js is unopinionated about how you organize and colocate your project files.`）であると明言しています。`_` の扱いも整理して把握しておくと迷いません。

- **フォルダの `_`（private folders）**: 公式機能として明記されているが、公式は「colocation のためには必須ではない」と明言しています。命名衝突の回避（将来 Next.js が新しい convention 名を追加した場合の保険）、UI ロジックと routing の分離、エディタ上でのソートなど、いくつかの用途で有用とされています
- **ファイルの `_`（例: `_data.ts`）**: 「フレームワーク規約ではないが、同じパターンで private 扱いを示す書き方として検討に値する」と公式が言及しています。挙動上は `data.ts` と一切違いがないため、付けるかどうかは慣習・好みの問題です

本コースでは、シンプルさを優先してファイルには `_` を付けず、フォルダで colocation する場合に必要に応じて `_` を使う方針にしています。

## 理解度チェック

- 静的なルートと動的なルートはそれぞれフォルダ構造でどう表現しますか
- レイアウトの中でルートに応じて切り替わる部分は、どこにどう書きますか
