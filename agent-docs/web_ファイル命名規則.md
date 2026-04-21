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
