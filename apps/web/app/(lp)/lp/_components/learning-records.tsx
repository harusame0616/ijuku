import { Divider } from "@/components/divider";
import { cn } from "@/lib/utils";

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

function AccuracyBar({ correct, total }: { correct: number; total: number }) {
  const pct = correct / total;
  return (
    <div className="flex items-center gap-2 shrink-0">
      {/* プログレスバーの背景色。コンソール表示用の専用色 */}
      <div className="w-20 h-1 overflow-hidden bg-[oklch(1_0_0/0.15)]">
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
          "font-mono text-xs w-10 text-right",
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

export function LearningRecords() {
  return (
    <section
      id="学習記録"
      className="relative py-28 px-8 overflow-hidden bg-background bg-[radial-gradient(circle,oklch(0.75_0.12_77/0.055)_1px,transparent_1px)] bg-size-[28px_28px] before:content-[''] before:absolute before:top-0 before:left-0 before:right-0 before:h-px before:z-20 before:bg-[linear-gradient(90deg,transparent,oklch(0.75_0.12_77/0.55)_25%,oklch(0.75_0.12_77/0.55)_75%,transparent)] after:content-[''] after:absolute after:inset-0 after:pointer-events-none after:z-0 after:bg-[radial-gradient(ellipse_90%_70%_at_50%_50%,oklch(0.165_0.035_50),oklch(0.10_0.015_50))]"
    >
      <div className="max-w-7xl mx-auto relative z-10">
        {/* ヘッダー */}
        <div className="flex flex-col items-center gap-5 mb-20 text-center">
          <div className="flex items-center gap-3">
            <div className="w-8 h-px bg-primary-dim" />
            <span className="font-mono text-xs uppercase tracking-[0.25em] text-muted-foreground">
              Learning Records
            </span>
            <div className="w-8 h-px bg-primary-dim" />
          </div>
          <h2 className="font-serif font-black text-4xl lg:text-5xl text-foreground">
            記録が、成長をつくる。
          </h2>
          <p className="text-sm max-w-lg leading-relaxed text-muted-foreground">
            学習の全履歴が自動で残る。いつ・何を・どれだけ正解したか。
            弱点も強みも、データが教えてくれる。
          </p>
          <Divider className="max-w-xs mt-2" />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* ターミナル風ログ（アンバーフォスファー）
               テキスト色はCRTフォスファー再現のための専用値。デザイントークン外。
               - プライマリ: oklch(0.78 0.10 77) — コマンド・主要テキスト
               - ミュート:   oklch(0.76 0.07 70) — 日付・補助テキスト */}
          <div className="lg:col-span-2 overflow-hidden bg-[oklch(0.09_0.015_50)] border border-primary-dim">
            {/* ターミナルバー */}
            <div className="flex items-center gap-2 px-4 py-3 bg-[oklch(0.12_0.02_50)] border-b border-primary-dim">
              <div className="flex gap-1.5">
                {[
                  "oklch(0.50 0.08 77)",
                  "oklch(0.55 0.08 77)",
                  "oklch(0.60 0.08 77)",
                ].map((c, i) => (
                  <div
                    // biome-ignore lint/suspicious/noArrayIndexKey: fixed index
                    key={i}
                    className="w-3 h-3 rounded-full"
                    style={{ background: c }}
                  />
                ))}
              </div>
              <span className="flex-1 text-center font-mono text-xs text-[oklch(0.78_0.10_77)]">
                jukubox — learning_log
              </span>
            </div>

            {/* ログ内容 */}
            <div className="p-6 flex flex-col gap-4 font-mono">
              <div className="text-xs mb-1 text-[oklch(0.78_0.10_77)]">
                $ jukubox log --recent 4
              </div>
              {mockLog.map((entry) => (
                <div
                  key={entry.date}
                  className="flex items-center gap-4 pb-3 border-b border-border/10"
                >
                  <span className="text-xs shrink-0 w-24 text-[oklch(0.76_0.07_70)]">
                    {entry.date}
                  </span>
                  <span className="text-sm flex-1 truncate text-[oklch(0.78_0.10_77)]">
                    {entry.topic}
                  </span>
                  <AccuracyBar correct={entry.correct} total={entry.total} />
                  <span className="text-xs shrink-0 w-16 text-right text-[oklch(0.76_0.07_70)]">
                    {entry.duration}
                  </span>
                </div>
              ))}
              <div className="text-xs mt-1 flex items-center gap-1">
                <span className="animate-[juku-breathe_4s_ease-in-out_infinite] text-primary">
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
            <div className="bg-card backdrop-blur-[20px] p-5 flex flex-col gap-4 border border-border/8">
              <span className="font-serif font-bold text-sm text-foreground">
                要強化エリア
              </span>
              <div className="flex flex-wrap gap-2">
                {mockWeak.map((tag) => (
                  <span
                    key={tag}
                    className="px-2.5 py-1 font-mono text-xs uppercase tracking-wider border border-[oklch(0.65_0.08_77/0.75)] text-[oklch(0.78_0.10_77)] bg-primary/8"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>

            {/* 習得済み */}
            <div className="bg-card backdrop-blur-[20px] p-5 flex flex-col gap-4 border border-border/8">
              <span className="font-serif font-bold text-sm text-foreground">
                習得済みスキル
              </span>
              <div className="flex flex-wrap gap-2">
                {mockStrong.map((tag) => (
                  <span
                    key={tag}
                    className="px-2.5 py-1 font-mono text-xs uppercase tracking-wider border border-secondary/65 text-secondary bg-secondary/10"
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>

            {/* 統計 */}
            <div className="bg-card backdrop-blur-[20px] p-5 flex flex-col gap-4 border border-border/8">
              <span className="font-serif font-bold text-sm text-foreground">
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
                    <span className="font-mono text-xs text-muted-foreground">
                      {stat.label}
                    </span>
                    <span className="font-orbitron font-bold text-xl [text-shadow:0_0_18px_oklch(0.75_0.12_77/0.45)] text-primary">
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
