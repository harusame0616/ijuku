import type { Metadata } from "next";
import type { JSX } from "react";

export const metadata: Metadata = {
  title: "コースを作る | JukuBox",
};

export default function NewCoursePage(): JSX.Element {
  return (
    <div className="mx-auto flex max-w-2xl flex-col items-center gap-4 py-20 text-center">
      <span className="font-mono text-[0.65rem] uppercase tracking-[0.3em] text-muted-foreground">
        Coming&nbsp;Soon
      </span>
      <h1 className="font-serif text-2xl font-bold text-foreground lg:text-3xl">
        コース作成画面は準備中です
      </h1>
      <p className="max-w-md text-sm leading-relaxed text-muted-foreground">
        制作側のフローを整備中です。 公開準備が整い次第、 ここから自分のコースを作って公開できるようになります。
      </p>
    </div>
  );
}
