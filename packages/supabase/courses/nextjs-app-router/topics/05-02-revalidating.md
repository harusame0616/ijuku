---
title: "updateTag で投稿直後にキャッシュを無効化する"
description: "Server Action から updateTag('posts') を呼び、新規投稿が即座に一覧へ反映されるようにする。"
---

# updateTag で投稿直後にキャッシュを無効化する

## 目標

- `updateTag` を使ってタグ単位でキャッシュを即時に無効化できる
- `revalidatePath` から `updateTag('posts')` への置き換えを行う
- `updateTag` と `revalidateTag` の違いを説明できる

## 知識

`'use cache'` でキャッシュした結果を更新する方法（**リバリデーション**）には大きく 2 種類あります。

- **時間ベース**: `cacheLife` で指定した期間が過ぎると自動で再生成される
- **オンデマンド**: ミューテーションが起きたタイミングで明示的にキャッシュを無効化する

オンデマンド側の主な API は次の 3 つです。

- **`updateTag(tag)`**: タグ付きキャッシュを **即座に** 失効させる。Server Action 専用。「ユーザーが自分の操作の結果をすぐ見たい（read-your-own-writes）」に向く
- **`revalidateTag(tag, profile)`**: タグ付きキャッシュを **stale-while-revalidate** で更新する。古い結果を返しつつバックグラウンドで再生成する。Server Action と Route Handler の両方で使える。第 2 引数には `'max'` を渡すのが推奨（単一引数の旧シグネチャは非推奨）
- **`revalidatePath(path)`**: パス単位で全キャッシュを無効化する。タグが分からないときに使う

ミニブログでは「投稿した直後にユーザー自身が自分の記事を見えるようにしたい」ので、**`updateTag('posts')`** が最適です。前章で `cacheTag('posts')` を付けたので、Server Action からこのタグを失効させれば、`PostList` のキャッシュが即座に作り直されます。

## ステップ

### 1. createPost を updateTag に置き換える

`app/posts/actions.ts` で `revalidatePath('/posts')` を `updateTag('posts')` に書き換えます。`next/cache` から `updateTag` を import し、`revalidatePath` の import は不要になるので削除します。

```ts
"use server";

import { updateTag } from "next/cache";
import * as v from "valibot";
import { addPost } from "./data";
import { postInputSchema, type PostInput } from "./schema";

export async function createPost(input: PostInput): Promise<void> {
  const { title, body } = v.parse(postInputSchema, input);

  const slug =
    title
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, "-")
      .slice(0, 40) || `post-${Date.now()}`;

  await addPost({ slug, title, body });
  updateTag("posts");
}
```

### 2. 動作確認

ブラウザで `/posts` から記事を投稿してください。フォーム送信後、リロードや遷移なしで一覧の先頭に新規記事が現れるはずです。前章で `revalidatePath` を使っていたときと結果は似ていますが、`updateTag` は **タグ単位**でピンポイントに無効化するため、関係ない他のキャッシュには影響しないという利点があります。

### 3. revalidateTag との違いを比べる

学習用に、`updateTag('posts')` を一時的に `revalidateTag('posts', 'max')` に置き換え、もう一度投稿を試してみてください。一見同じように動きますが、こちらは **古い結果がしばらく返り続ける** stale-while-revalidate の挙動になります。確認後は `updateTag` に戻しておきましょう。

## 完了判定

- `app/posts/actions.ts` で `updateTag('posts')` が呼ばれている
- 投稿後、`/posts` 一覧に新規記事が即座に反映される
- `updateTag` と `revalidateTag` の挙動の違いを自分の言葉で説明できる

## 補足

`updateTag` は **Server Action 専用** で、Route Handler の中では使えません。Route Handler でタグ無効化したい場合は `revalidateTag` を使ってください。`revalidateTag` の第 2 引数には `cacheLife` のプロファイル名（推奨は `'max'`）かオブジェクト（例: `{ expire: 0 }`）を渡します。`'max'` は最も長い stale 許容期間を持つプロファイルで、stale-while-revalidate の挙動が最も色濃く出ます（応答は早いがデータの鮮度は遅れがち）。タグ名は文字列リテラルで自由に決められますが、プロジェクト内で命名規則をそろえる（例: `posts`、`posts:slug:hello` など）と運用しやすくなります。

## 理解度チェック

- `updateTag` と `revalidateTag` の違いを「即時無効化」と「stale-while-revalidate」のキーワードで説明してください
- ミューテーション直後に必ず最新が見える挙動を実現したい場合、どちらをどこで呼びますか
