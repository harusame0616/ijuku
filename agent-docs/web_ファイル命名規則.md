# Web ファイル命名規則

next.js の file convention は除外する（page.tsx など）

## サーバーコンポーネント

\*.server.tsx

## クライアントコンポーネント

\*.client.tsx

## ユニバーサルコンポーネント

\*.universal.tsx

ただし use client がついていなくてもクライアントコンポーネントとして使用することを想定している場合はクライアントコンポーネントとして扱い命名する

## スケルトンコンポーネント

\*.skeleton.[server|client|universal].tsx

## container コンポーネント

\*.container.[server|client|universal].tsx

## presenter コンポーネント

\*.presenter.[server|client|universal].tsx

## form コンポーネント

\*.form.client.tsx

## Server Action

\*.action.ts

## RSC 用データフェッチ

\*.data.ts

## サーバー専用ユーティリティ（コンポーネント以外）

\*.server.ts

サーバーで実行される前提、next/headers, next/navigation, DB クライアントなどサーバーでのみ動く依存を含むモジュールに付与する

## テストファイル

### server テスト（node 環境）

- small（外部依存なし）: \*.small.server.test.ts
- medium（DB など外部依存あり）: \*.medium.server.test.ts

### client テスト（browser mode）

- hook・スクリプト: \*.client.test.ts
- コンポーネント: \*.client.test.tsx

### 対象外

- サーバーコンポーネント自体はテスト対象外
- E2E は別チケットで対応
