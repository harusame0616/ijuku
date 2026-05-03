---
title: "Server Components と Client Components を使い分ける"
description: "記事カードの「いいね」ボタンを Client Component で実装し、サーバーとクライアントの境界を理解する。"
---

# Server Components と Client Components を使い分ける

## 目標

- App Router の既定が Server Component であることを理解する
- 状態（`useState`）を持つ「いいね」ボタンを Client Component として作る
- Server Component から Client Component に props を渡せる

## 知識

App Router では **すべてのファイルがデフォルトで Server Component** です。Server Component はサーバー上で実行され、データベースアクセスや API キー、機密情報を安全に扱えます。クライアントに送られる JavaScript の量も減らせるため、初期表示が速くなります。

ブラウザでしかできないこと（イベントハンドラ、`useState` などの状態、`useEffect`、`localStorage`、`window`）が必要になったときは **Client Component** に切り替えます。Client Component は、ファイルの先頭に `'use client'` ディレクティブを書くことで宣言します。

```tsx
"use client";

import { useState } from "react";

export function LikeButton() {
  const [liked, setLiked] = useState(false);
  return (
    <button onClick={() => setLiked((v) => !v)}>
      {liked ? "♥ いいね済み" : "♡ いいね"}
    </button>
  );
}
```

`'use client'` は **境界宣言** です。一度宣言したファイルから import したコンポーネントは、すべてクライアントバンドルに含まれるようになります。だからこそ、`'use client'` はできるだけ末端の小さいコンポーネントに付け、上位はできる限り Server Component のままにするのが定石です。

Server Component から Client Component には **props を経由してデータを渡せます**。ただし React によってシリアライズされるため、関数（Server Action として渡す場合を除く）やクラスインスタンスなどシリアライズできない値は渡せない点に注意してください（プリミティブ値、配列、プレーンオブジェクト、Date など React がシリアライズ可能な型は OK）。

## ステップ

### 1. LikeButton を作る

`app/posts/_components/like-button.tsx` を作成し、Client Component として「いいね」ボタンを実装します。`_components` はプライベートフォルダなのでルートになりません。

```tsx
"use client";

import { useState } from "react";

export function LikeButton({ initialLikes }: { initialLikes: number }) {
  const [likes, setLikes] = useState(initialLikes);
  return (
    <button type="button" onClick={() => setLikes((n) => n + 1)}>
      ♡ {likes}
    </button>
  );
}
```

### 2. 一覧と詳細にボタンを設置する

`app/posts/page.tsx` のリスト要素と `app/posts/[slug]/page.tsx` の本文の下に、それぞれ `<LikeButton initialLikes={0} />` を追加します。Server Component から Client Component に props（数値）を渡している点を意識してください。

```tsx
import { LikeButton } from "./_components/like-button";

// 一覧側の例
<li key={post.slug}>
  <Link href={`/posts/${post.slug}`}>{post.title}</Link>
  <LikeButton initialLikes={0} />
</li>;
```

### 3. 動作確認

ブラウザでボタンをクリックすると数値が増えること、ページ遷移するとカウンタがリセットされること（状態は当該コンポーネントのライフサイクルに依存）を確認してください。

## 完了判定

- `app/posts/_components/like-button.tsx` の先頭に `'use client'` がある
- 一覧および詳細ページに「いいね」ボタンが表示され、クリックで数値が増える
- 一覧ページとレイアウトを含むほとんどの UI は Server Component のまま動作している

## 補足

`'use client'` を `app/layout.tsx` のような上位に書くと、子コンポーネントもすべてクライアントバンドルに入ってしまいます。理由なく付けると JavaScript のバンドルサイズが増え、`fetch` や DB クエリをサーバーで実行する利点も損なわれます。React の Context を提供したい場合は `'use client'` を付けた Provider コンポーネントを作り、Server Component の `layout.tsx` から `<Provider>{children}</Provider>` のように呼び出すのが基本パターンです。

## 理解度チェック

- Server Component と Client Component で「使えるもの」が分かれる代表的な例を 2 つ挙げ、なぜ分かれるのか説明してください
- Server Component から Client Component に値を渡すときの制約は何ですか。その制約があるのはなぜですか
