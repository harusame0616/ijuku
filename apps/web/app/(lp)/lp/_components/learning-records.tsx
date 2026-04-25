import type { JSX } from "react";
import { Divider } from "@/components/divider";
import { cn } from "@/lib/utilities";

const mockLog = [
  {
    date: "2025-04-03",
    topic: "React Server Components",
    correct: 8,
    total: 10,
    duration: "42 min",
  },
  {
    date: "2025-04-02",
    topic: "TypeScript Generics",
    correct: 6,
    total: 9,
    duration: "35 min",
  },
  {
    date: "2025-04-01",
    topic: "SQL Window Functions",
    correct: 10,
    total: 10,
    duration: "58 min",
  },
  {
    date: "2025-03-31",
    topic: "Docker Compose",
    correct: 5,
    total: 8,
    duration: "29 min",
  },
];

const mockWeak = ["Closure", "Currying", "Middleware"];
const mockStrong = ["REST API", "Git Flow", "CSS Grid"];

function AccuracyBar({
  correct,
  total,
}: {
  correct: number;
  total: number;
}): JSX.Element {
  const pct = correct / total;
  return (
    <div className="flex shrink-0 items-center gap-2">
      {/* プログレスバーの背景色。コンソール表示用の専用色 */}
      <div className="h-1 w-20 overflow-hidden bg-[oklch(1_0_0/0.15)]">
        <div
          className={cn(
            "h-full transition-all",
            pct >= 0.9
              ? "bg-primary"
              : pct >= 0.7
                ? "bg-secondary"
                : /* コンソール表示専用の警告色（デザインシステムトークン外） */
                  "bg-[oklch(0.76_0.08_60)]",
          )}
          style={{ width: `${pct * 100}%` }}
        />
      </div>
      <span
        className={cn(
          "w-10 text-right font-mono text-xs",
          pct >= 0.9
            ? "text-primary"
            : pct >= 0.7
              ? "text-secondary"
              : /* コンソール表示専用の警告色（デザインシステムトークン外） */
                "text-[oklch(0.76_0.08_60)]",
        )}
      >
        {correct}/{total}
      </span>
    </div>
  );
}

export function LearningRecords(): JSX.Element {
  return (
    <section
      id="学習記録"
      className="bg-background relative overflow-hidden bg-[radial-gradient(circle,oklch(0.75_0.12_77/0.055)_1px,transparent_1px)] bg-size-[28px_28px] px-8 py-28 before:absolute before:top-0 before:right-0 before:left-0 before:z-20 before:h-px before:bg-[linear-gradient(90deg,transparent,oklch(0.75_0.12_77/0.55)_25%,oklch(0.75_0.12_77/0.55)_75%,transparent)] before:content-[''] after:pointer-events-none after:absolute after:inset-0 after:z-0 after:bg-[radial-gradient(ellipse_90%_70%_at_50%_50%,oklch(0.165_0.035_50),oklch(0.10_0.015_50))] after:content-['']"
    >
      <div className="relative z-10 mx-auto max-w-7xl">
        {/* ヘッダー */}
        <div className="mb-20 flex flex-col items-center gap-5 text-center">
          <div className="flex items-center gap-3">
            <div className="bg-primary-dim h-px w-8" />
            <span className="text-muted-foreground font-mono text-xs tracking-[0.25em] uppercase">
              Learning Records
            </span>
            <div className="bg-primary-dim h-px w-8" />
          </div>
          <h2 className="text-foreground font-serif text-4xl font-black lg:text-5xl">
            記録が、成長をつくる。
          </h2>
          <p className="text-muted-foreground max-w-lg text-sm leading-relaxed">
            学習の全履歴が自動で残る。いつ・何を・どれだけ正解したか。
            弱点も強みも、データが教えてくれる。
          </p>
          <Divider className="mt-2 max-w-xs" />
        </div>

        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          {/* ターミナル風ログ（アンバーフォスファー）
               テキスト色はCRTフォスファー再現のための専用値。デザイントークン外。
               - プライマリ: oklch(0.78 0.10 77) — コマンド・主要テキスト
               - ミュート:   oklch(0.76 0.07 70) — 日付・補助テキスト */}
          <div className="border-primary-dim overflow-hidden border bg-[oklch(0.09_0.015_50)] lg:col-span-2">
            {/* ターミナルバー */}
            <div className="border-primary-dim flex items-center gap-2 border-b bg-[oklch(0.12_0.02_50)] px-4 py-3">
              <div className="flex gap-1.5">
                {[
                  "oklch(0.50 0.08 77)",
                  "oklch(0.55 0.08 77)",
                  "oklch(0.60 0.08 77)",
                ].map((c, index) => (
                  <div
                    // biome-ignore lint/suspicious/noArrayIndexKey: fixed index
                    key={index}
                    className="h-3 w-3 rounded-full"
                    style={{ background: c }}
                  />
                ))}
              </div>
              <span className="flex-1 text-center font-mono text-xs text-[oklch(0.78_0.10_77)]">
                jukubox — learning_log
              </span>
            </div>

            {/* ログ内容 */}
            <div className="flex flex-col gap-4 p-6 font-mono">
              <div className="mb-1 text-xs text-[oklch(0.78_0.10_77)]">
                $ jukubox log --recent 4
              </div>
              {mockLog.map((entry) => (
                <div
                  key={entry.date}
                  className="border-border/10 flex items-center gap-4 border-b pb-3"
                >
                  <span className="w-24 shrink-0 text-xs text-[oklch(0.76_0.07_70)]">
                    {entry.date}
                  </span>
                  <span className="flex-1 truncate text-sm text-[oklch(0.78_0.10_77)]">
                    {entry.topic}
                  </span>
                  <AccuracyBar correct={entry.correct} total={entry.total} />
                  <span className="w-16 shrink-0 text-right text-xs text-[oklch(0.76_0.07_70)]">
                    {entry.duration}
                  </span>
                </div>
              ))}
              <div className="mt-1 flex items-center gap-1 text-xs">
                <span className="text-primary animate-[juku-breathe_4s_ease-in-out_infinite]">
                  █
                </span>
                <span className="text-[oklch(0.76_0.07_70)]">
                  waiting for input...
                </span>
              </div>
            </div>
          </div>

          {/* サイドパネル */}
          <div className="flex flex-col gap-5">
            {/* 要強化 */}
            <div className="bg-card border-border/8 flex flex-col gap-4 border p-5 backdrop-blur-[20px]">
              <span className="text-foreground font-serif text-sm font-bold">
                要強化エリア
              </span>
              <div className="flex flex-wrap gap-2">
                {mockWeak.map((tag) => (
                  <span
                    key={tag}
                    className="bg-primary/8 border border-[oklch(0.65_0.08_77/0.75)] px-2.5 py-1 font-mono text-xs tracking-wider text-[oklch(0.78_0.10_77)] uppercase"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>

            {/* 習得済み */}
            <div className="bg-card border-border/8 flex flex-col gap-4 border p-5 backdrop-blur-[20px]">
              <span className="text-foreground font-serif text-sm font-bold">
                習得済みスキル
              </span>
              <div className="flex flex-wrap gap-2">
                {mockStrong.map((tag) => (
                  <span
                    key={tag}
                    className="border-secondary/65 text-secondary bg-secondary/10 border px-2.5 py-1 font-mono text-xs tracking-wider uppercase"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>

            {/* 統計 */}
            <div className="bg-card border-border/8 flex flex-col gap-4 border p-5 backdrop-blur-[20px]">
              <span className="text-foreground font-serif text-sm font-bold">
                今月の統計
              </span>
              <div className="grid grid-cols-2 gap-4">
                {[
                  { label: "学習日数", value: "18日" },
                  { label: "総時間", value: "24h" },
                  { label: "解いた問題", value: "312問" },
                  { label: "正答率", value: "78%" },
                ].map((stat) => (
                  <div key={stat.label} className="flex flex-col gap-1">
                    <span className="text-muted-foreground font-mono text-xs">
                      {stat.label}
                    </span>
                    <span className="font-orbitron text-primary text-xl font-bold [text-shadow:0_0_18px_oklch(0.75_0.12_77/0.45)]">
                      {stat.value}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
