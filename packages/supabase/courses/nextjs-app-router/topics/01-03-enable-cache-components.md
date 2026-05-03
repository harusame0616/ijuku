---
title: "Cache Components を有効化する"
description: "next.config.ts に cacheComponents: true を設定し、新しいキャッシュモデルに切り替える。"
---

# Cache Components を有効化する

## 目標

- `next.config.ts` に `cacheComponents: true` を追加する
- 本コースが Next.js 16 の Cache Components 前提で進むことを把握する

## 知識

Next.js 16 から導入された **Cache Components** は、App Router の新しいキャッシュモデルです。本コースではこのモデルを前提にミニブログを組み立てていくため、最初にプロジェクトで有効化しておきます。

ここでは「設定フラグを 1 つ入れて、新しいキャッシュ機能群を使えるようにする」とだけ理解できれば十分です。具体的な使い方（`'use cache'` ディレクティブの付け方、タグでの無効化、Suspense と組み合わせたストリーミングなど）は後続のトピック（5 章「use cache でキャッシュする」「updateTag で無効化する」、6 章「loading.tsx と Suspense でストリーミングする」など）で、実際に手を動かしながら学びます。

## ステップ

### 1. next.config.ts を編集する

`mini-blog` プロジェクト直下にある `next.config.ts` を開き、`cacheComponents: true` を追加してください。

```ts
import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  cacheComponents: true,
};

export default nextConfig;
```

`next.config.js` の場合は `module.exports` 形式に合わせて同じ内容を追記します。

### 2. 開発サーバーを再起動する

`next.config.ts` の変更はホットリロードでは反映されないことが多いため、開発サーバーを一度止めて `pnpm dev` で起動し直してください。`/` にアクセスして引き続きホームが表示されれば OK です。

## 完了判定

- `next.config.ts` に `cacheComponents: true` が記述されている
- 開発サーバーを再起動してもエラーなくホームが表示される

## 補足

`cacheComponents` は Next.js 16 で導入されたフラグです。Next.js 15 以前では `experimental.ppr` / `experimental.dynamicIO` / `experimental.useCache` のように複数のフラグに分かれていましたが、16 ではそれらがまとめられたため、このフラグだけ入れれば本コースで扱う新キャッシュ機能群がすべて使える状態になります。古いバージョン用の個別フラグを足す必要はありません。
