# Web ファイル命名規則

next.js の file convention は除外する（page.tsx など）

## コンポーネントの種別判定フロー

上から順に評価し、最初に該当した区分を採用する：

1. ファイル先頭に `"use client"` がある → **client**
2. 以下のいずれかを使う → **server**
   - `async` 関数として宣言されている
   - サーバー専用モジュール（`next/headers`、 DB クライアント、 シークレット環境変数など）の参照
3. 以下のいずれかを使う → **client**（`"use client"` を付与）
   - React の hook を使う（`useState` / `useEffect` / `useRef` / `useContext` 等）
   - イベントハンドラ（`onClick` 等）の直接定義
   - ブラウザ専用 API（`window` / `document` / `localStorage` 等）
   - 内部で `"use client"` を要求するライブラリ
4. 上記いずれにも該当しない（props を受けて JSX を返すだけ、 子に client component を含むだけ） → **universal**

## サーバーコンポーネント

\*.server.tsx

## クライアントコンポーネント

\*.client.tsx

## ユニバーサルコンポーネント

\*.universal.tsx

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
