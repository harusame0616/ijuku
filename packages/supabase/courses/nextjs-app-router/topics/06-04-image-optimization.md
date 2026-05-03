---
title: "next/image で画像を最適化する"
description: "next/image でレスポンシブ srcset・WebP/AVIF 配信・lazy loading・CLS 防止を有効化する。"
---

# next/image で画像を最適化する

## 目標

- `<Image>` の役割（最適化フォーマット配信・CLS 防止・lazy loading）を理解する
- ローカル画像とリモート画像の扱いの違いを把握し、`images.remotePatterns` でドメインを許可する
- 一覧サムネイルと詳細ヒーローで `<Image>` を使い分ける（`width`/`height`・`preload`・`sizes`）

## 知識

Next.js の **`<Image>`**（`next/image`）は HTML の `<img>` を拡張し、画像配信に関する面倒な部分をフレームワーク側で引き受けるコンポーネントです。公式ドキュメントに挙げられている主な機能は次の 4 つです:

- **サイズ最適化**: ビューポートやデバイスピクセル比に合った最適サイズの画像を配信。AVIF / WebP のような新しいフォーマットも自動で使う
- **レイアウトの安定**: `width` / `height` から `aspect-ratio` を予約し、画像が読み込まれるまでのレイアウトシフト（CLS）を防ぐ
- **読み込みの遅延**: ビューポートに入るまで画像をロードしない（ネイティブの `loading="lazy"` がデフォルト）。任意で blur-up プレースホルダも出せる
- **アセットの柔軟性**: リモート画像でもオンデマンドにリサイズできる

```tsx
import Image from "next/image";

<Image src="/cover.png" alt="表紙" width={1200} height={600} />;
```

### プロジェクト内の画像

プロジェクトに同梱した画像を `<Image>` に渡す方法は、別の仕組みである 2 通りに分かれます。

**1. `public/` 配下のファイルを URL 文字列で参照する**

`public/` ディレクトリに置いたファイルはルート URL（`/cover.png` など）で配信されます。`<Image>` には文字列パスを渡します。この経路は **ファイルがバンドラを通らない** ため、`width` と `height` を自前で指定する必要があります（`blurDataURL` も自動では埋まりません）。

```tsx
import Image from "next/image";

<Image src="/cover.png" alt="表紙" width={1200} height={600} />;
```

**2. 画像ファイルを静的 import する**

jpg / png / webp / avif ファイルを **ソースコード内のどこからでも**（`public/` に置く必要はなく、コンポーネントの隣でも `app/` 配下でもよい）`import cover from './cover.png'` の形で import できます。こちらは画像ファイルがバンドラを通るため、Next.js がビルド時に画像を解析し、**`width` ・ `height` ・ `blurDataURL` が自動で埋まります**。

```tsx
import Image from "next/image";
import cover from "./cover.png"; // 例: コンポーネントの隣に置いた画像

<Image src={cover} alt="表紙" />; // width/height/blurDataURL 不要
```

`public/` の URL 参照は「アプリのアセットとして単に配信する」、静的 import は「コードと一緒にバンドラに食わせて最適化情報を取り出す」という別の仕組みなので、両者を混同しないようにしてください。

### リモート画像

外部 URL の画像を扱うには、**`next.config.ts` の `images.remotePatterns` で明示的にホストを許可** する必要があります（許可していないホストへのリクエストは `400 Bad Request` で拒否される）。リモート画像はビルド時に取得できないため、`width` / `height` を必ず手動で指定するか、`fill` で親要素に追従させます。

### `width` / `height` と `fill`

`width` と `height` は **画像本来の寸法（intrinsic size）をピクセルで** 渡します。これはレンダリングサイズではなく `aspect-ratio` の確定に使われ、CSS 側で実際の表示サイズを制御できます。寸法が分からないときは `fill` を使い、親要素に `position: relative` などを付けて領域を広げます。

### `sizes` で srcset を最適化する

レスポンシブに表示する画像（特に `fill` を使うとき）は `sizes` を必ず指定します。`(max-width: 768px) 100vw, 50vw` のようにブレークポイントごとに「画面幅に対する画像の表示幅」を伝えると、ブラウザがその情報をもとに `srcset` から最適な解像度を選びます。`sizes` を省くと「画像はビューポート幅と同じ」と仮定され、過大なサイズがダウンロードされがちです。

### lazy loading と `preload`

`<Image>` のデフォルトは `loading="lazy"` で、ビューポートに入るまで読み込まれません。一方、ファーストビューに表示される **LCP（Largest Contentful Paint）候補のヒーロー画像** には `preload={true}` を付け、`<head>` に `<link rel="preload">` を入れて先読みさせます。

> **重要**: Next.js 16 から、従来の `priority` プロパティは **非推奨** になり、同じ役割の `preload` プロパティに置き換えられました。新規コードは `preload={true}` を使ってください。

### `placeholder` と `blurDataURL`

`placeholder="blur"` を付けると、低解像度のぼかし画像を読み込み完了までの代わりに表示できます。静的 import の jpg/png/webp/avif なら `blurDataURL` が自動で生成されます（アニメーション画像は対象外）。リモート画像は自分で `blurDataURL` を渡す必要があります。

### `alt` は必須

スクリーンリーダーや画像読み込み失敗時の代替テキストとして表示されます。装飾目的の画像は空文字 `alt=""` を渡し、意味を持つ画像はその意味を文章で書きます。

### 最適化したくない画像

SVG・アニメーション GIF・1KB 未満の小さい画像など、最適化のメリットが乏しい場合は `unoptimized` を付けて素のまま配信します。認証ヘッダが必要な画像も Image Optimization API がヘッダを転送しない仕様なので `unoptimized` 向きです。

## ステップ

### 1. データに `cover` フィールドを追加する

`app/posts/data.ts` の `Post` 型と初期データに、表紙画像 URL を任意プロパティとして追加します（学習用に [picsum.photos](https://picsum.photos/) のシード付き URL を使います）。

```ts
export type Post = {
  slug: string;
  title: string;
  body: string;
  cover?: string;
};

// 初期データの slug ごとに cover を追加
const posts: Post[] =
  globalThis.__miniBlogPosts ??
  (globalThis.__miniBlogPosts = [
    {
      slug: "hello",
      title: "Hello, Mini Blog",
      body: "最初の記事です。",
      cover: "https://picsum.photos/seed/hello/1200/600",
    },
    {
      slug: "next",
      title: "Next.js を学ぶ",
      body: "App Router の入門中。",
      cover: "https://picsum.photos/seed/next/1200/600",
    },
  ]);
```

`addPost` の引数の型も `Post` のままで大丈夫です（`cover` は任意なので未指定でも通る）。

### 2. リモート画像ドメインを許可する

`next.config.ts` に `images.remotePatterns` を追加し、`picsum.photos` からの画像読み込みを許可します。

```ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  cacheComponents: true,
  images: {
    remotePatterns: [
      {
        protocol: "https",
        hostname: "picsum.photos",
      },
    ],
  },
};

export default nextConfig;
```

設定変更を反映するため、開発サーバーを一度止めて再起動してください。

### 3. 一覧ページのサムネイルを `<Image>` で表示する

`app/posts/page.tsx` の `PostList` で、各記事の左に小さな表紙画像を出します。一覧の画像は基本的にビューポート外にあることが多いので、デフォルトの `lazy` ローディングのまま使います。

```tsx
import Image from "next/image";
// 既存 import はそのまま

async function PostList() {
  "use cache";
  cacheLife("hours");
  cacheTag("posts");

  const posts = await getPosts();
  return (
    <ul>
      {posts.map((post) => (
        <li key={post.slug}>
          {post.cover && (
            <Image src={post.cover} alt="" width={120} height={60} />
          )}
          <Link href={`/posts/${post.slug}`}>{post.title}</Link>
          <LikeButton initialLikes={0} />
        </li>
      ))}
    </ul>
  );
}
```

サムネイルはタイトルの補助的な装飾なので `alt=""` にしています（記事内容を表現するメインのコンテンツ画像なら、適切な alt 文を入れてください）。

### 4. 詳細ページのヒーロー画像を表示する

`app/posts/[slug]/page.tsx` の `PostContent` で、本文の上に大きな表紙画像を出します。これはファーストビューの主要画像、つまり **LCP 候補** なので `preload={true}` を付けて先読みさせます。あわせて `sizes` で srcset の生成戦略を伝えます。

```tsx
import Image from "next/image";
// 既存 import はそのまま

async function PostContent({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;
  const post = await getPost(slug);
  if (!post) return <p>記事が見つかりません</p>;
  return (
    <article>
      {post.cover && (
        <Image
          src={post.cover}
          alt=""
          width={1200}
          height={600}
          preload
          sizes="(max-width: 768px) 100vw, 800px"
        />
      )}
      <h1>{post.title}</h1>
      <p>{post.body}</p>
    </article>
  );
}
```

`sizes` のヒントは「768px 以下のビューポートでは画像幅 = ビューポート幅、それ以上なら最大 800px」という意味です。これを付けると Next.js が `srcset` に複数の解像度を並べ、ブラウザがビューポートと dpr に応じて最適な 1 枚を選びます。

### 5. 動作確認

ブラウザで `/posts` と `/posts/hello` を開き、開発者ツールで以下を確認します。

- Network タブで画像のレスポンス `content-type` が **AVIF または WebP** になっている
- 一覧の画像が **ビューポートに入った瞬間に読み込まれる**（ページ最下部までスクロールして確認）
- 詳細ページの表紙画像が **ページ HTML より早く `<link rel="preload">` として `<head>` に入っている**（Elements タブで確認）
- ページ表示時にレイアウトがガタつかない（CLS が抑えられている）

## 完了判定

- `Post` 型に `cover?: string` が追加されている
- `next.config.ts` の `images.remotePatterns` に `picsum.photos` が登録されている
- 一覧ページに各記事のサムネイル画像が `<Image>` で表示される
- 詳細ページに `<Image preload>` でヒーロー画像が表示される
- 開発者ツールで AVIF/WebP 配信と CLS 抑制が確認できる

## 補足

`fill` レイアウトを使う場合は親要素に `position: relative`（または `fixed` / `absolute`）を付け、画像が領域いっぱいに広がるようにします。`object-fit: cover` / `contain` で切り抜きの仕方を制御できます。寸法（aspect-ratio）が確定しないと CLS の保証が崩れるので、固定サイズなら `width`/`height` を、可変サイズなら `fill` + 親の `position` + `sizes` を組み合わせるのが基本です。

`placeholder="blur"` を有効にすると、画像読み込み完了までの「ガクっと出る感」をぼかしで隠せます。静的 import の画像は `blurDataURL` が自動で埋まりますが、リモート画像では自前で `blurDataURL` を用意するか、ビルド時に [Plaiceholder](https://github.com/joe-bell/plaiceholder) のようなライブラリで生成します。

リモート画像のドメイン許可は **可能な限り具体的に** 書くのが推奨です。`hostname` だけでなく `pathname` まで絞ると、攻撃者が任意の画像を最適化エンドポイント経由で配信させる経路を狭められます。`hostname: '**.example.com'` のようなワイルドカードも使えますが、サブドメインを丸ごと許可するため使いどころは限定的です。

## 理解度チェック

- `<Image>` を使うと素の `<img>` に対して何が改善されますか。代表的な 3 点を挙げてください
- リモート画像を `<Image>` に渡すとき、なぜ `next.config.ts` の `images.remotePatterns` への登録が必要ですか
- 一覧サムネイルと詳細のヒーロー画像で、ローディング戦略をどう使い分けますか（`loading` / `preload` / `sizes` の観点で）
