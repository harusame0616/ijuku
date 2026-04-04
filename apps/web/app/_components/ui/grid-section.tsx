import type { ComponentPropsWithoutRef } from "react";

// グリッド背景は専用の表現で、デザインシステムトークン外。
// 色値: oklch(0.75 0.12 77 / 0.04) — primaryと同色相の極薄アンバー
const GRID_BG =
  "bg-background bg-[linear-gradient(oklch(0.75_0.12_77/0.04)_1px,transparent_1px),linear-gradient(90deg,oklch(0.75_0.12_77/0.04)_1px,transparent_1px)] bg-size-[48px_48px]";

export function GridSection({
  className,
  ...props
}: ComponentPropsWithoutRef<"section">) {
  return <section className={`${GRID_BG} ${className ?? ""}`} {...props} />;
}
