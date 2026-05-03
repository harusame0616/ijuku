---
title: "metadata と generateMetadata でタイトルを設定する"
description: "ルートレイアウトに静的 metadata を、記事詳細ページに generateMetadata を設定する。"
---

# metadata と generateMetadata でタイトルを設定する

## 目標

- ルートレイアウトに静的 `metadata` を export してサイト全体の既定タイトル・説明を設定する
- 記事詳細ページで `generateMetadata` を実装し、記事タイトルを `<title>` に反映する
- メタデータ取得とページ表示で **データ取得を重複させない** 工夫を理解する

## 知識

App Router でメタデータを設定する方法は 2 つあります。

- **静的 `metadata` オブジェクト**: `layout.tsx` または `page.tsx` から `metadata` を export する。値が固定なときに使う
- **動的 `generateMetadata` 関数**: 引数の `params` / `searchParams` を使ってデータをフェッチし、結果を `Metadata` オブジェクトとして返す

```tsx
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Mini Blog",
  description: "Next.js App Router で作るミニブログ",
};
```

```tsx
export async function generateMetadata({
  params,
}: {
  params: Promise<{ slug: string }>;
}): Promise<Metadata> {
  const { slug } = await params;
  const post = await getPost(slug);
  return { title: post?.title ?? "Not Found" };
}
```

`generateMetadata` の中でも `params` は **Promise** で渡されます。動的にレンダリングされるページでは、Next.js が **メタデータをストリーミング** で送信し、本文より遅くても UI 描画をブロックしないように工夫されています（HTML だけを読むタイプのクローラ向けには User Agent を見て自動的にブロッキング送信に切り替わります）。

本コースのように Cache Components が有効な場合、`generateMetadata` も他の Server Component と同じルールに従います。`params` / `searchParams` / `cookies()` / `headers()` などのランタイムデータや非キャッシュ fetch にアクセスすると、`generateMetadata` はリクエスト時に解決される扱いになります。ページや子コンポーネントもランタイムデータに依存している場合はそのままストリーミングされますが、ページが完全にプリレンダリング可能なときは「メタデータだけリクエスト時」という構成は意図的かどうか分からないため、ビルド時にエラーで誘導されます。記事タイトルのように外部データ依存ではあるがランタイムデータには依存しない場合は、`generateMetadata` 内で `'use cache'` を付けることでメタデータをキャッシュして整合させられます。

データ取得が「ページ本体」と「メタデータ」で重複するのが気になる場合は、React の `cache()` を使って同じリクエスト内のメモ化を効かせます。

```ts
import { cache } from "react";

export const getPost = cache(async (slug: string) => {
  // ...
});
```

これにより `generateMetadata` と `Page` が同じ `slug` で `getPost` を呼んでも、データ取得は 1 回で済みます。

`metadata` と `generateMetadata` のエクスポートは Server Component でのみサポートされる点も覚えておいてください。

## ステップ

### 1. ルートレイアウトに静的 metadata を追加する

`app/layout.tsx` の上部で `metadata` を export し、サイト名と説明を設定します。子ページが `title` を上書きしない限り、この値が `<title>` に入ります。

```tsx
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Mini Blog",
  description: "Next.js App Router で作るミニブログ",
};
```

### 2. 記事詳細で generateMetadata を実装する

`app/posts/[slug]/page.tsx` に `generateMetadata` を追加し、記事タイトルを `<title>` に入れます。

```tsx
import type { Metadata } from "next";
import { getPost } from "../data";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ slug: string }>;
}): Promise<Metadata> {
  const { slug } = await params;
  const post = await getPost(slug);
  return {
    title: post ? `${post.title} | Mini Blog` : "記事が見つかりません",
  };
}
```

### 3. ブラウザで確認する

`/posts/hello` にアクセスし、ブラウザのタブに「Hello, Mini Blog | Mini Blog」のように記事タイトルが表示されることを確認してください。`/` などの他ページではルートレイアウトの値が使われ、「Mini Blog」と表示されているはずです。

## 完了判定

- ルートレイアウトに `metadata` が export されている
- 記事詳細ページに `generateMetadata` が実装され、ブラウザのタブに記事タイトルが反映される
- どのページにアクセスしても `<title>` が空にならない

## 補足

`generateMetadata` と本文の両方から同じデータを取りたい場合は、データ取得関数を `React.cache()` で包むと、同じリクエスト内では一度だけ実行されるようになります。`opengraph-image.tsx` のような特殊ファイルを使うと OG 画像をルート単位で動的生成でき、SNS シェアの体験を強化できます。`metadataBase` を `layout.tsx` で設定しておくと、相対 URL のメタデータを絶対 URL に解決でき、本番デプロイ時に OG 画像のパスが崩れません。
