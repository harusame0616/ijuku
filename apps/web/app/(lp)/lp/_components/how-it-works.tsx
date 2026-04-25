import type { JSX } from "react";
import { Divider } from "@/components/divider";

const steps = [
  {
    number: "01",
    title: "AIサービスを接続",
    description:
      "ChatGPT、Claude、Gemini など、自分が使っているAIサービスをそのまま接続。MCP や Skill 経由で簡単に連携できます。",
  },
  {
    number: "02",
    title: "コースを選ぶ or 作る",
    description:
      "ライブラリから既存コースを選ぶか、独自のカリキュラムを作成。社内研修にも個人学習にも対応。",
  },
  {
    number: "03",
    title: "AIと学習スタート",
    description:
      "AIがコースに沿って出題・解説・フィードバック。まるで専属の家庭教師がいるような感覚。",
  },
  {
    number: "04",
    title: "記録を確認・振り返る",
    description:
      "学習ログ・正誤記録・進捗グラフで成長を可視化。何が得意で何が苦手かが一目でわかります。",
  },
];

export function HowItWorks(): JSX.Element {
  return (
    <section id="使い方" className="relative py-28 px-8 bg-background">
      <div className="max-w-7xl mx-auto">
        {/* ヘッダー */}
        <div className="flex flex-col items-center gap-5 mb-20 text-center">
          <div className="flex items-center gap-3">
            <div className="w-8 h-px bg-primary-dim" />
            <span className="font-mono text-xs uppercase tracking-[0.25em] text-muted-foreground">
              How It Works
            </span>
            <div className="w-8 h-px bg-primary-dim" />
          </div>
          <h2 className="font-serif font-black text-4xl lg:text-5xl text-foreground">
            使い方は、シンプル。
          </h2>
          <Divider className="max-w-xs mt-2" />
        </div>

        {/* ステップ */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-0 relative">
          {steps.map((step, index) => (
            <div
              key={step.number}
              className="relative flex flex-col items-center text-center px-6 py-8 gap-5"
            >
              {/* ステップ番号 */}
              <div className="relative z-10 w-20 h-20 rounded-full flex items-center justify-center shrink-0 bg-primary/20 border border-primary-dim">
                <span className="font-orbitron font-black text-2xl text-primary">
                  {step.number}
                </span>
              </div>

              {/* タイトル */}
              <h3 className="font-serif font-bold text-base leading-snug text-foreground">
                {step.title}
              </h3>

              {/* 説明 */}
              <p className="text-sm leading-relaxed text-muted-foreground">
                {step.description}
              </p>

              {/* モバイル矢印 */}
              {index < steps.length - 1 && (
                <div className="lg:hidden text-lg mt-1 text-primary-dim">↓</div>
              )}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
