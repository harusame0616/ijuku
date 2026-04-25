import type { JSX } from "react";
import { Divider } from "@/components/divider";

const features = [
  {
    symbol: "◈",
    title: "自分のAIで学ぶ",
    description:
      "ChatGPT、Claude、Gemini など、すでに契約しているAIサービスをそのまま活用。追加コスト不要で最先端AIによる個別指導を実現します。",
    tag: "Bring Your Own AI",
  },
  {
    symbol: "◉",
    title: "コースを作る・使う",
    description:
      "自分でカリキュラムを組んでコースを公開。あるいはコミュニティが用意したコースをすぐに使って学習開始できます。",
    tag: "Course Library",
  },
  {
    symbol: "◎",
    title: "学習を記録する",
    description:
      "いつ・何を・どれだけ学んだかを自動記録。タイムライン形式で学習の軌跡を可視化し、継続するモチベーションを支えます。",
    tag: "Progress Tracking",
  },
  {
    symbol: "◐",
    title: "正解・不正解を分析",
    description:
      "何を理解し、何を間違えたのかを詳細に記録。弱点を把握してAIに重点的に質問できる、データドリブンな学習サイクル。",
    tag: "Mistake Analysis",
  },
];

export function Features(): JSX.Element {
  return (
    <section
      id="機能"
      className="relative py-28 px-8 overflow-hidden bg-background bg-[radial-gradient(circle,oklch(0.75_0.12_77/0.055)_1px,transparent_1px)] bg-size-[28px_28px] before:content-[''] before:absolute before:top-0 before:left-0 before:right-0 before:h-px before:z-20 before:bg-[linear-gradient(90deg,transparent,oklch(0.75_0.12_77/0.55)_25%,oklch(0.75_0.12_77/0.55)_75%,transparent)] after:content-[''] after:absolute after:inset-0 after:pointer-events-none after:z-0 after:bg-[radial-gradient(ellipse_90%_70%_at_50%_50%,oklch(0.165_0.035_50),oklch(0.10_0.015_50))]"
    >
      <div className="max-w-7xl mx-auto relative z-10">
        {/* ヘッダー */}
        <div className="flex flex-col items-center gap-5 mb-20 text-center">
          <div className="flex items-center gap-3">
            <div className="w-8 h-px bg-primary-dim" />
            <span className="font-mono text-xs uppercase tracking-[0.25em] text-muted-foreground">
              Features
            </span>
            <div className="w-8 h-px bg-primary-dim" />
          </div>
          <h2 className="font-serif font-black text-4xl lg:text-5xl leading-tight text-foreground">
            学びを、自由に。
          </h2>
          <p className="text-sm max-w-lg leading-relaxed text-muted-foreground">
            塾の構造化された学習と、Jukebox のように
            「好きな曲を選んで流す」自由さを融合させた学習体験。
          </p>
          <Divider className="max-w-xs mt-2" />
        </div>

        {/* カードグリッド */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-5">
          {features.map((feature) => (
            <div
              key={feature.title}
              className="bg-card backdrop-blur-[20px] border border-border/18 relative overflow-hidden transition-all duration-400 ease-in p-8 flex flex-col gap-5 before:content-[''] before:absolute before:top-0 before:left-0 before:right-0 before:h-px before:bg-[linear-gradient(90deg,transparent,var(--primary-dim),transparent)] hover:-translate-y-0.75 hover:border-border/12"
            >
              {/* シンボル */}
              <div className="text-4xl leading-none text-primary">
                {feature.symbol}
              </div>

              {/* タイトル */}
              <h3 className="font-serif font-bold text-lg leading-snug text-foreground">
                {feature.title}
              </h3>

              {/* 説明 */}
              <p className="text-sm leading-relaxed flex-1 text-muted-foreground">
                {feature.description}
              </p>

              {/* タグ */}
              <span className="font-mono text-[10px] uppercase tracking-[0.2em] text-subtle-foreground">
                {feature.tag}
              </span>
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
