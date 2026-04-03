const mockLog = [
  { date: "2025-04-03", topic: "React Server Components", correct: 8, total: 10, duration: "42 min" },
  { date: "2025-04-02", topic: "TypeScript Generics",     correct: 6, total: 9,  duration: "35 min" },
  { date: "2025-04-01", topic: "SQL Window Functions",    correct: 10, total: 10, duration: "58 min" },
  { date: "2025-03-31", topic: "Docker Compose",          correct: 5, total: 8,  duration: "29 min" },
];

const mockWeak   = ["Closure", "Currying", "Middleware"];
const mockStrong = ["REST API", "Git Flow", "CSS Grid"];

function AccuracyBar({ correct, total }: { correct: number; total: number }) {
  const pct = correct / total;
  const color =
    pct >= 0.9
      ? "var(--gold)"
      : pct >= 0.7
      ? "oklch(0.72 0.09 190)"
      : "oklch(0.60 0.06 55)";
  return (
    <div className="flex items-center gap-2 shrink-0">
      <div
        className="w-20 h-1 overflow-hidden"
        style={{ background: "oklch(1 0 0 / 0.07)" }}
      >
        <div
          className="h-full transition-all"
          style={{ width: `${pct * 100}%`, background: color }}
        />
      </div>
      <span className="font-space-mono text-xs w-10 text-right" style={{ color }}>
        {correct}/{total}
      </span>
    </div>
  );
}

export function LearningRecords() {
  return (
    <section
      id="学習記録"
      className="relative py-28 px-8"
      style={{ background: "var(--background-warm)" }}
    >
      <div className="max-w-7xl mx-auto">
        {/* ヘッダー */}
        <div className="flex flex-col items-center gap-5 mb-20 text-center">
          <div className="flex items-center gap-3">
            <div className="w-8 h-px" style={{ background: "var(--gold-dim)" }} />
            <span
              className="font-space-mono text-xs uppercase tracking-[0.25em]"
              style={{ color: "var(--muted-foreground)" }}
            >
              Learning Records
            </span>
            <div className="w-8 h-px" style={{ background: "var(--gold-dim)" }} />
          </div>
          <h2
            className="font-noto-serif-jp font-black text-4xl lg:text-5xl"
            style={{ color: "var(--foreground)" }}
          >
            記録が、成長をつくる。
          </h2>
          <p
            className="text-sm max-w-lg leading-relaxed"
            style={{ color: "var(--muted-foreground)" }}
          >
            学習の全履歴が自動で残る。いつ・何を・どれだけ正解したか。
            弱点も強みも、データが教えてくれる。
          </p>
          <div className="juku-divider max-w-xs mt-2" />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* ターミナル風ログ（アンバーフォスファー） */}
          <div
            className="lg:col-span-2 overflow-hidden"
            style={{
              background: "oklch(0.09 0.015 50)",
              border: "1px solid var(--gold-dim)",
            }}
          >
            {/* ターミナルバー */}
            <div
              className="flex items-center gap-2 px-4 py-3"
              style={{
                background: "oklch(0.12 0.02 50)",
                borderBottom: "1px solid var(--gold-dim)",
              }}
            >
              <div className="flex gap-1.5">
                {[
                  "oklch(0.50 0.08 77)",
                  "oklch(0.55 0.08 77)",
                  "oklch(0.60 0.08 77)",
                ].map((c, i) => (
                  <div key={i} className="w-3 h-3 rounded-full" style={{ background: c }} />
                ))}
              </div>
              <span
                className="flex-1 text-center font-space-mono text-xs"
                style={{ color: "oklch(0.40 0.06 65)" }}
              >
                jukubox — learning_log
              </span>
            </div>

            {/* ログ内容 */}
            <div className="p-6 flex flex-col gap-4 font-space-mono">
              <div className="text-xs mb-1" style={{ color: "oklch(0.55 0.07 77)" }}>
                $ jukubox log --recent 4
              </div>
              {mockLog.map((entry) => (
                <div
                  key={entry.date}
                  className="flex items-center gap-4 pb-3"
                  style={{ borderBottom: "1px solid oklch(1 0 0 / 0.05)" }}
                >
                  <span
                    className="text-xs shrink-0 w-24"
                    style={{ color: "oklch(0.40 0.06 65)" }}
                  >
                    {entry.date}
                  </span>
                  <span
                    className="text-sm flex-1 truncate"
                    style={{ color: "oklch(0.80 0.08 75)" }}
                  >
                    {entry.topic}
                  </span>
                  <AccuracyBar correct={entry.correct} total={entry.total} />
                  <span
                    className="text-xs shrink-0 w-16 text-right"
                    style={{ color: "oklch(0.38 0.05 60)" }}
                  >
                    {entry.duration}
                  </span>
                </div>
              ))}
              <div className="text-xs mt-1 flex items-center gap-1">
                <span className="juku-breathe" style={{ color: "var(--gold)" }}>█</span>
                <span style={{ color: "oklch(0.38 0.05 60)" }}>waiting for input...</span>
              </div>
            </div>
          </div>

          {/* サイドパネル */}
          <div className="flex flex-col gap-5">
            {/* 要強化 */}
            <div
              className="juku-glass p-5 flex flex-col gap-4"
              style={{ border: "1px solid oklch(1 0 0 / 0.08)" }}
            >
              <span
                className="font-noto-serif-jp font-bold text-sm"
                style={{ color: "var(--foreground)" }}
              >
                要強化エリア
              </span>
              <div className="flex flex-wrap gap-2">
                {mockWeak.map((tag) => (
                  <span
                    key={tag}
                    className="px-2.5 py-1 font-space-mono text-xs uppercase tracking-wider"
                    style={{
                      border: "1px solid oklch(0.55 0.07 77 / 0.5)",
                      color: "oklch(0.62 0.08 75)",
                      background: "oklch(0.75 0.12 77 / 0.05)",
                    }}
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>

            {/* 習得済み */}
            <div
              className="juku-glass p-5 flex flex-col gap-4"
              style={{ border: "1px solid oklch(1 0 0 / 0.08)" }}
            >
              <span
                className="font-noto-serif-jp font-bold text-sm"
                style={{ color: "var(--foreground)" }}
              >
                習得済みスキル
              </span>
              <div className="flex flex-wrap gap-2">
                {mockStrong.map((tag) => (
                  <span
                    key={tag}
                    className="px-2.5 py-1 font-space-mono text-xs uppercase tracking-wider"
                    style={{
                      border: "1px solid oklch(0.72 0.09 190 / 0.4)",
                      color: "var(--teal)",
                      background: "oklch(0.72 0.09 190 / 0.06)",
                    }}
                  >
                    {tag}
                  </span>
                ))}
              </div>
            </div>

            {/* 統計 */}
            <div
              className="juku-glass p-5 flex flex-col gap-4"
              style={{ border: "1px solid oklch(1 0 0 / 0.08)" }}
            >
              <span
                className="font-noto-serif-jp font-bold text-sm"
                style={{ color: "var(--foreground)" }}
              >
                今月の統計
              </span>
              <div className="grid grid-cols-2 gap-4">
                {[
                  { label: "学習日数", value: "18日" },
                  { label: "総時間",   value: "24h" },
                  { label: "解いた問題", value: "312問" },
                  { label: "正答率",   value: "78%" },
                ].map((stat) => (
                  <div key={stat.label} className="flex flex-col gap-1">
                    <span
                      className="font-space-mono text-xs"
                      style={{ color: "var(--muted-foreground)" }}
                    >
                      {stat.label}
                    </span>
                    <span
                      className="font-orbitron font-bold text-xl juku-glow-gold-text"
                      style={{ color: "var(--gold)" }}
                    >
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
