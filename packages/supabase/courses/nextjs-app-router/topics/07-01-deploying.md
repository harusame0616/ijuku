---
title: "Vercel にデプロイして公開する"
description: "GitHub にミニブログを push し、Vercel から接続して本番 URL を取得する。"
---

# Vercel にデプロイして公開する

## 目標

- 作成したミニブログを GitHub リポジトリに push する
- Vercel に接続してデプロイし、本番 URL でアクセスできるようにする
- 本番ビルドが通ることを確認する

## 知識

Next.js は本番にもいろいろな方法でデプロイできます。代表的な選択肢を整理しておきます。

- **Node.js サーバー**: 任意の VPS やコンテナサービス。`next build` で生成し、`next start` で起動する。Next.js の全機能をサポート
- **Docker コンテナ**: コンテナオーケストレーターやクラウドのコンテナサービスで動かす。これも全機能サポート
- **Static export**: 完全な静的サイトとして書き出す。Server Actions、Cookies、Proxy、ISR、Server を必要とする Route Handler など、サーバーが必要な機能は使えない
- **Adapters / マネージドサービス**: **Vercel** と **Bun** が verified adapter を提供しており（Next.js GitHub Org でホストされ、互換性テストスイートを通過している）、Cloudflare、Netlify、Firebase App Hosting なども独自の Next.js 連携を提供している

本コースでは作者が Next.js を作っている **Vercel** に乗せます。Vercel は GitHub と連携し、push のたびに自動でビルドとデプロイを走らせ、PR ごとに **プレビュー URL** も発行してくれます。`'use cache'` や `updateTag` といった Cache Components の機能もそのまま使えます（Cache Components は **ストリーミング対応のサーバー** を必要とするため、Vercel のような verified adapter / Node.js サーバー / Docker のいずれかで動かす必要があります）。

デプロイ前にローカルで `next build` を通しておくのが鉄則です。Cache Components が有効な状態だと、`cookies()` / `headers()` などのリクエスト時データに触れる箇所が `<Suspense>` で囲まれていないとビルド時にエラーになります。手元で見つけて直しておきましょう。

## ステップ

### 1. ローカルで本番ビルドを通す

`mini-blog` ディレクトリで本番ビルドを実行し、警告やエラーが出ないか確認します。

```bash
pnpm build
```

ビルド成果のサマリでは、各ルートに `○ (Static) prerendered as static content` か `ƒ (Dynamic) server-rendered on demand` の記号が付きます。`/posts` のようなキャッシュした一覧が Static 側になっていることを確認してください。問題があれば、メッセージに従って `<Suspense>` を追加するか `'use cache'` を見直します。

### 2. GitHub にリポジトリを作って push する

GitHub に新規リポジトリ（例: `mini-blog`）を作成し、ローカルの状態を push します。

```bash
git init
git add .
git commit -m "feat: initial mini blog"
git branch -M main
git remote add origin https://github.com/<your-account>/mini-blog.git
git push -u origin main
```

### 3. Vercel に接続してデプロイする

[https://vercel.com](https://vercel.com) にログインし、「New Project」から先ほどの GitHub リポジトリをインポートします。フレームワークは自動で **Next.js** が選ばれます。環境変数は本コースの範囲では不要です。

「Deploy」を押すとビルドが始まり、数十秒〜数分で完了します。発行された URL（`https://<project>.vercel.app`）にアクセスし、ローカルと同じくホーム・記事一覧・記事詳細が表示されること、ヘッダーリンクで遷移できること、新規投稿フォームから記事を作れることを確認してください。

> 補足: in-memory ストアはデプロイ後に複数のサーバーインスタンスで共有されないため、本番環境では「投稿しても次のリクエストで消える」可能性があります。学習用なので深追いせず、永続化したい場合は次の学習として Supabase などの DB を導入してください。

## 完了判定

- ローカルで `pnpm build` がエラーなく完了する
- GitHub にリポジトリが作られて main ブランチが push されている
- Vercel から本番 URL が発行され、トップページと記事一覧が閲覧できる

## 補足

Vercel 以外でも、Node.js を実行できる環境なら `next build` と `next start` で動かせます。Docker で動かしたい場合は `output: 'standalone'` を `next.config.ts` に追加すると、`.next/standalone` 配下に必要なランタイム依存だけを含む軽量な成果物（`server.js` 付き）が生成されます。`public/` と `.next/static/` は CDN で配信する想定のためコピー対象外なので、必要に応じて `standalone/` 配下に手動でコピーします。本コースで使った Cache Components の機能（`'use cache'` / `cacheLife` / `cacheTag` / `updateTag`）はストリーミング対応サーバーが前提なので、サーバーを持たない完全な Static export では使えません。デプロイ後に `pnpm dev` をローカルで動かしながら本番 URL と挙動を比較するのも、原因切り分けに有効です。

## 理解度チェック

- Static export と Node.js サーバー / マネージド (Vercel 等) でのデプロイは、利用できる Next.js 機能にどんな差がありますか
