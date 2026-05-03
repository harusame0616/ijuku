---
title: "error.tsx で例外を捕捉する"
description: "Error Boundary としての error.tsx を追加し、レンダリング時の例外をユーザーに優しく見せる。"
---

# error.tsx で例外を捕捉する

## 目標

- セグメント単位の `error.tsx` を作成する
- Error Boundary が動作するよう `'use client'` を付ける
- `unstable_retry()` で再試行できるボタンを置く

## 知識

Next.js のエラーは大きく 2 種類に分けて考えます。

- **想定内のエラー**: バリデーション失敗や 404 のような、業務ロジックとして起きる得るもの。`useActionState` で値を返す、`notFound()` を呼ぶ、`if` で条件分岐するなど、明示的に扱う
- **想定外の例外**: バグやネットワーク異常など、本来発生してほしくないもの。**`error.tsx`** で Error Boundary を作ってフォールバック UI を表示する

`error.tsx` はセグメント配下のレンダリング中に投げられた例外を捕捉し、コンポーネントツリーが完全に壊れる代わりにフォールバック UI を見せます。Error Boundary は **Client Component でなければならない** ため、ファイル先頭に `'use client'` を付けます。

`error.tsx` のデフォルトエクスポートには次の props が渡されます。

- `error`: 投げられたエラーオブジェクト（`digest` を含む）
- `unstable_retry`: 該当セグメントを再フェッチ・再レンダリングして復旧を試みる関数

```tsx
"use client";

import { useEffect } from "react";

export default function Error({
  error,
  unstable_retry,
}: {
  error: Error & { digest?: string };
  unstable_retry: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);
  return (
    <div>
      <h2>エラーが発生しました</h2>
      <button onClick={() => unstable_retry()}>もう一度試す</button>
    </div>
  );
}
```

エラーは **最も近い親の `error.tsx`** にバブルアップします。だから `app/posts/error.tsx` を置けば `posts` 配下、`app/error.tsx` を置けばアプリ全体、と境界の粒度を変えられます。ルートレイアウト自体が壊れるレベルの例外には `app/global-error.tsx` を別途用意しますが、本コースでは扱いません。

注意点として、Error Boundary は **イベントハンドラ内のエラー** や **非同期コードのエラー** は捕捉しません。これらは `try` / `catch` で取って `useState` に積む、もしくは `useTransition` の `startTransition` 内で投げて Error Boundary に持ち上げる、といったハンドリングが必要です。

## ステップ

### 1. /posts 用の error.tsx を作る

`app/posts/error.tsx` を作成し、Client Component としてエラー UI を返します。

```tsx
"use client";

import { useEffect } from "react";

export default function PostsError({
  error,
  unstable_retry,
}: {
  error: Error & { digest?: string };
  unstable_retry: () => void;
}) {
  useEffect(() => {
    console.error(error);
  }, [error]);
  return (
    <main>
      <h2>記事の表示に失敗しました</h2>
      <p>少し時間を置いて再度お試しください。</p>
      <button type="button" onClick={() => unstable_retry()}>
        再試行
      </button>
    </main>
  );
}
```

### 2. わざと例外を投げて検証する

`app/posts/page.tsx` の `PostList` 関数の中で、一時的に `throw new Error('intentional');` を入れて挙動を確認してください。`/posts` を開くと先ほど作った `error.tsx` のフォールバックが表示されるはずです。確認できたら `throw` を削除します。

### 3. 「再試行」ボタンを押す

フォールバック画面で「再試行」を押すと、Next.js がそのセグメントを再レンダリングします。`throw` を残したままだと再びエラー UI に戻り、`throw` を消したあと再試行すれば一覧が再表示されます。

## 完了判定

- `app/posts/error.tsx` が作成され、ファイル先頭に `'use client'` がある
- 例外を投げると、ヘッダーは保ったままページ部分だけがフォールバック UI に切り替わる
- 「再試行」ボタンで `unstable_retry()` を呼べる

## 補足

`error.tsx` は同じセグメントの `layout.tsx` の例外は捕捉しません。レイアウト自身を保護したい場合は、ひとつ上のセグメントに `error.tsx` を置きます。`error.digest` は本番ビルドでクライアントから見えるエラー識別子で、サーバーログのスタックトレースと突き合わせるのに使えます。観測には `console.error` だけでなく Sentry などのエラートラッカーへの送信を組み合わせるのが定番です。

## 理解度チェック

- `error.tsx` が捕捉できるエラーと捕捉できないエラーをそれぞれ説明してください
