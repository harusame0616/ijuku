import type { Metadata } from "next";
import type { JSX } from "react";

export const metadata: Metadata = {
  title: "コース一覧 | JukuBox",
};

export default function CoursesPage(): JSX.Element {
  return (
    <div className="mx-auto flex max-w-2xl flex-col items-center gap-4 py-20 text-center">
      <span className="font-mono text-[0.65rem] uppercase tracking-[0.3em] text-muted-foreground">
        Coming&nbsp;Soon
      </span>
      <h1 className="font-serif text-2xl font-bold text-foreground lg:text-3xl">
        コース一覧は準備中です
      </h1>
      <p className="max-w-md text-sm leading-relaxed text-muted-foreground">
        マーケットプレイスを開発中です。 まもなく検索・カテゴリ・新着を含めて公開します。
      </p>
    </div>
  );
}
