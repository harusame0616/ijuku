---
title: "Server Action で新規記事を投稿する"
description: "Valibot のスキーマをクライアント・サーバーで共有し、react-hook-form のフォームから Server Action を呼ぶ。"
---

# Server Action で新規記事を投稿する

## 目標

- `'use server'` ディレクティブで Server Action を定義する
- `Valibot` のスキーマをクライアント・サーバー両方で共有してバリデーションする
- `react-hook-form` で組んだ Client Component の投稿フォームから Server Action を呼ぶ
- 投稿後、一覧に新規記事が反映される

## 知識

データの書き込み（**ミューテーション**）は、Next.js の **Server Action** を使うのが標準的なやり方です。Server Action は React の **Server Functions** という機能の応用で、`'use server'` ディレクティブを付けたサーバー上の非同期関数のことを指します。

- 関数本体の先頭に `'use server'` を書く
- もしくはファイル先頭に `'use server'` を書くと、そのファイルからエクスポートされる全ての関数が Server Function になる

```ts
"use server";

export async function createPost(input: { title: string; body: string }) {
  // ... DB 書き込み
}
```

Server Action を `<form>` の `action` プロップに渡すと、フォーム送信時にブラウザは自動的に POST リクエストを送り、Next.js が該当の Server Function を実行してくれます。`<form>` 経由で受け取る引数は **`FormData`** ですが、Client Component 経由で呼ぶときは普通の引数として任意の型を渡せます（後述）。

このバインドの利点が **プログレッシブエンハンスメント（progressive enhancement）** です。「まず HTML だけで基本機能が動く状態を作り、JavaScript が読み込めた環境ではそこに UX 強化を上乗せする」という設計思想で、`<form action={createPost}>` は JavaScript の読み込み中・部分失敗時でも HTML 標準のフォーム POST として動作します。

現代の Web ではほぼ全てのデバイスに JS エンジンが載っているため「JS を無効化したユーザー対応」という古典的な動機は薄れていますが、今のプログレッシブエンハンスメントは次の 3 点を重視します:

- **JS 読み込みの空白を埋める**: HTML 到着後に JS バンドルが解析・実行されるまでの数秒、SPA だとボタンを押しても何も起きない（dead click）。`<form action={...}>` ならその空白でも送信が成立する
- **JS の部分的失敗に強い**: バンドル 404・CDN 障害・サードパーティスクリプトの例外伝播などで SPA は全停止しがちだが、HTML レベルで動く部分は残る（数 % オーダーで観測される現実的な問題）
- **サーバー中心で設計が単純になる**: クライアント中心の状態管理は楽観的 UI・キャッシュ無効化・整合性を自前で組むことになる。「サーバーが正、フォーム POST が基本、JS は強化層」という構造はロジックがシンプルで、エッジケースが減る

Server Action はこの「サーバー中心 + JS は強化層」の構造をフレームワーク標準で提供します。JS が効く環境では `useActionState` の pending 状態管理や `useOptimistic` の楽観的 UI 更新が上乗せされ、効かない環境では HTML 標準のフォーム POST に自然に劣化します。

Server Function は `<form action={...}>` 以外にも、Client Component のイベントハンドラから直接呼び出したり、`useActionState` / `useTransition` と組み合わせて使うこともできます。`<form>` バインドはあくまで「最も簡単に動く呼び出し方」で、用途に応じて他の呼び出し方を選べます。

ミューテーション後にページの表示を更新するには、`revalidatePath('/posts')` のようにキャッシュを無効化するか、`redirect('/posts/...')` で別ページに遷移します。Cache Components 環境で「ユーザーが直後に自分の変更を確認する」ような操作には、次章で扱う `updateTag` を使うとより自然です。今回はまず一覧ページに即時反映するため `revalidatePath` を使います。

> 重要: Server Function は実装上 **HTTP の POST エンドポイントとして公開される** ため、外部から任意のリクエストボディで呼び出せます。クライアントの呼び出しを前提に書いた処理でも、実際には任意の値・型・サイズのデータが届きえます。関数の冒頭で必ず以下の 2 つを明示的に行ってください:
>
> - **認証・認可の確認**: 呼び出し元のユーザーがその操作を行ってよいかチェック
> - **入力値のバリデーション**: 文字数・必須・型・形式を検証。型注釈や `FormData.get` の戻り値型だけに頼らない（静的型はランタイムの値を保証しない）
>
> 本コースでは学習用 in-memory ストアを使うため認証・認可は省略しますが、バリデーションは **Valibot** のスキーマでクライアント・サーバー両方に適用します。

クライアント側でも送信前にバリデーションを行いたい場合（入力途中のリアルタイムエラー表示など UX 向上のため）、いくつか選択肢があります。

- **フォームライブラリ + バリデーションライブラリ**: `react-hook-form` や `TanStack Form` といったフォームライブラリと、`Zod` / `Valibot` のようなスキーマバリデーションライブラリを組み合わせる。スキーマを 1 つ書けばクライアント・サーバーで再利用できる利点がある一方、ネイティブ submit を奪う構造のため PE と両立させたい場合は別の選択肢が必要（補足参照）
- **ライブラリを使わず生で書く**: `useState` ベースで自分でエラー状態を管理しても問題ありません。シンプルなフォームならむしろ軽量

ただしどの方法を採るにせよ、**クライアントサイドのバリデーションはあくまで UX のためであり、サーバー側のバリデーションを省略してよい根拠にはなりません**。攻撃者はクライアントを介さず Server Function を直接叩けるため、サーバー側のバリデーションは常に必須です。

本トピックでは実用性を優先し、`react-hook-form` + `Valibot` でクライアントバリデーションを行い、その送信ハンドラから Server Action を呼ぶ構成（後述補足の選択肢 A）を採用します。

## ステップ

### 1. 必要なパッケージを追加する

クライアントフォーム用とバリデーション用のパッケージを追加します。

```bash
pnpm add valibot react-hook-form @hookform/resolvers
```

### 2. 共有バリデーションスキーマを作る

`app/posts/schema.ts` を作成し、`Valibot` で投稿入力のスキーマと型をエクスポートします。クライアントとサーバーの両方からこのスキーマを import して再利用します。

```ts
import * as v from "valibot";

export const postInputSchema = v.object({
  title: v.pipe(
    v.string(),
    v.trim(),
    v.nonEmpty("タイトルを入力してください"),
    v.maxLength(80, "タイトルは 80 文字以内で入力してください"),
  ),
  body: v.pipe(v.string(), v.trim(), v.nonEmpty("本文を入力してください")),
});

export type PostInput = v.InferOutput<typeof postInputSchema>;
```

### 3. Server Action を作る

`app/posts/actions.ts` を作成し、ファイル先頭に `'use server'` を書いて `createPost` をエクスポートします。受け取った値は **必ず `v.parse` でバリデーションしてから** ストアに渡します（クライアント側でチェック済みでも、攻撃者は直接叩けるためサーバー側で改めて検証）。

```ts
"use server";

import { revalidatePath } from "next/cache";
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
  revalidatePath("/posts");
}
```

### 4. 投稿フォームを Client Component で作る

`app/posts/_components/post-form.tsx` を作成します。`react-hook-form` の `useForm` に `Valibot` スキーマの resolver を渡してクライアント側バリデーションを行い、`handleSubmit` の中から Server Action を呼びます。送信中は `useTransition` で pending 状態を管理します。

```tsx
"use client";

import { useTransition } from "react";
import { useForm } from "react-hook-form";
import { valibotResolver } from "@hookform/resolvers/valibot";
import { createPost } from "../actions";
import { postInputSchema, type PostInput } from "../schema";

export function PostForm() {
  const [isPending, startTransition] = useTransition();
  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<PostInput>({
    resolver: valibotResolver(postInputSchema),
    defaultValues: { title: "", body: "" },
  });

  const onSubmit = handleSubmit((data) => {
    startTransition(async () => {
      await createPost(data);
      reset();
    });
  });

  return (
    <form onSubmit={onSubmit}>
      <label>
        タイトル
        <input type="text" {...register("title")} />
      </label>
      {errors.title && <p role="alert">{errors.title.message}</p>}

      <label>
        本文
        <textarea {...register("body")} />
      </label>
      {errors.body && <p role="alert">{errors.body.message}</p>}

      <button type="submit" disabled={isPending}>
        {isPending ? "投稿中..." : "投稿する"}
      </button>
    </form>
  );
}
```

### 5. 一覧ページにフォームを配置する

`app/posts/page.tsx` の `<h1>` 直後に `<PostForm />` を配置します。`PostForm` は Client Component ですが、Server Component（ページ）から import するだけで自動的にクライアントバンドルへ移ります。

```tsx
import Link from "next/link";
import { Suspense } from "react";
import { getPosts } from "./data";
import { LikeButton } from "./_components/like-button";
import { PostForm } from "./_components/post-form";

async function PostList() {
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

export default function PostsPage() {
  return (
    <main>
      <h1>記事一覧</h1>
      <PostForm />
      <Suspense fallback={<p>読み込み中...</p>}>
        <PostList />
      </Suspense>
    </main>
  );
}
```

### 6. 動作確認

ブラウザで `/posts` を開いて以下を試してください。

- タイトルか本文を空欄のまま送信 → クライアント側でエラーメッセージが表示され、リクエストは飛ばない
- タイトルに 81 文字以上を入れて送信 → 同様にクライアント側でエラー
- 正常に入力して送信 → 一覧の先頭に新しい記事が現れ、フォームはリセットされる。`/posts/<生成されたスラッグ>` で本文も確認できる

## 完了判定

- 投稿フォームから新しい記事を作成できる
- 空欄や規定外の値で送信した場合は、クライアント側でエラーメッセージが表示される
- 送信後、一覧ページに新しい記事が反映され、フォームがリセットされる
- 詳細ページ（`/posts/<新しい slug>`）でも本文が読める

## 補足

`'use server'` をファイル先頭に書くか、関数本体の先頭に書くかは目的に応じて使い分けます。Client Component から import する場合はファイル先頭に書く必要があります。`useTransition` の代わりに `useActionState` を使うと、`pending` 状態に加えてサーバー側のエラー結果（バリデーション失敗時のフィールドエラーなど）をフォームへ戻す処理が素直に書けます。in-memory のストアは開発サーバーを再起動すると消える点も覚えておいてください。永続化したい場合は次のステップとして DB（Supabase など）への置き換えを検討します。

**プログレッシブエンハンスメントとフォームライブラリの両立**

実務でクライアント側バリデーションを `react-hook-form` のような `onSubmit` 介入型ライブラリで行う場合、`<form action={serverAction}>` のネイティブ submit 経路と素直には共存しません。`onSubmit={handleSubmit(...)}` が submit イベントを横取りするため、JS が読み込まれる前は submit が動かず、PE が後退します。選択肢は概ね 3 つです:

- **A. `react-hook-form` + `handleSubmit` の中で Server Action を呼ぶ**: UX は最良。ただし JS 必須で PE は失われ、`useActionState` の恩恵も一部捨てる。実務では現状もっとも多い構成で、管理画面など JS 前提でよい領域に向く。**本コースで採用しているのはこの方法**
- **B. `<form action={serverAction}>` を維持し、`react-hook-form` はバリデーション表示のみに使う**: submit はネイティブ POST、入力中のエラー表示は `react-hook-form` + `Valibot`。submit 時のチェックはサーバーに任せ、エラーは `useActionState` で受け取る。PE は成立するが、submit 直前のクライアント側ガードは省略される妥協的な構成
- **C. [Conform](https://conform.guide/) のような PE 前提のフォームライブラリを使う**: ネイティブ submit 経路を壊さず、同じ `Valibot` スキーマをクライアント・サーバー双方で適用できる。Server Actions と統合される設計で、PE を本気で守りたい場合の標準解

PE をどこまで重視するかはプロジェクト判断で技術的な唯一解はありません。本コースは「学習者が普段の業務でも使い回せるように」という観点から、現状もっとも採用例が多い A を選択しています。PE を本気で重視する案件では C を検討してください。

## 理解度チェック

- Server Action は何を「サーバー側」で実行できるようにする仕組みですか。クライアントから呼び出される際にどう扱われますか
- ミューテーション後に一覧を最新化する代表的な手段を 2 つ挙げ、それぞれが何を無効化するのか説明してください
- Valibot のスキーマをクライアントとサーバー両方で再利用する利点は何ですか
