---
paths:
  - *.form.client.tsx
---

# Form

- スタック
  - Shadcn/UI Field
  - TANSTACK Form
  - valibot
- placeholder は使用せず、常にテキストで説明を表示する
- ラベルに必ず必須、任意を表示する(RequiredOptionalBadge コンポーネント)
- useIsHydrated hook を使いハイドレーションが完了していない時は disable にする
