const steps = [
  {
    number: "01",
    title: "AIエージェントを接続",
    description:
      "Claude、GPT-4o など、自分が使っているAIをそのまま接続。MCP や Skill 経由で簡単に連携できます。",
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
      "AIエージェントがコースに沿って出題・解説・フィードバック。まるで専属の家庭教師がいるような感覚。",
  },
  {
    number: "04",
    title: "記録を確認・振り返る",
    description:
      "学習ログ・正誤記録・進捗グラフで成長を可視化。何が得意で何が苦手かが一目でわかります。",
  },
];

export function HowItWorks() {
  return (
    <section
      id="使い方"
      className="relative py-28 px-8"
      style={{ background: "var(--juku-bg)" }}
    >
      <div className="max-w-7xl mx-auto">
        {/* ヘッダー */}
        <div className="flex flex-col items-center gap-5 mb-20 text-center">
          <div className="flex items-center gap-3">
            <div className="w-8 h-px" style={{ background: "var(--juku-gold-dim)" }} />
            <span
              className="font-space-mono text-xs uppercase tracking-[0.25em]"
              style={{ color: "var(--juku-text-muted)" }}
            >
              How It Works
            </span>
            <div className="w-8 h-px" style={{ background: "var(--juku-gold-dim)" }} />
          </div>
          <h2
            className="font-noto-serif-jp font-black text-4xl lg:text-5xl"
            style={{ color: "var(--juku-text)" }}
          >
            使い方は、シンプル。
          </h2>
          <div className="juku-divider max-w-xs mt-2" />
        </div>

        {/* ステップ */}
        <div className="grid grid-cols-1 lg:grid-cols-4 gap-0 relative">
          {/* コネクターライン（デスクトップ） */}
          <div
            className="hidden lg:block absolute top-10 left-[12.5%] right-[12.5%] h-px"
            style={{
              background:
                "linear-gradient(90deg, transparent, var(--juku-gold-dim) 20%, var(--juku-gold-dim) 80%, transparent)",
              opacity: 0.4,
            }}
          />

          {steps.map((step, idx) => (
            <div
              key={step.number}
              className="relative flex flex-col items-center text-center px-6 py-8 gap-5"
            >
              {/* ステップ番号 */}
              <div
                className="relative z-10 w-20 h-20 rounded-full flex items-center justify-center flex-shrink-0"
                style={{
                  background: "var(--juku-bg-warm)",
                  border: "1px solid var(--juku-gold-dim)",
                }}
              >
                <span
                  className="font-orbitron font-black text-2xl"
                  style={{ color: "var(--juku-gold)" }}
                >
                  {step.number}
                </span>
              </div>

              {/* タイトル */}
              <h3
                className="font-noto-serif-jp font-bold text-base leading-snug"
                style={{ color: "var(--juku-text)" }}
              >
                {step.title}
              </h3>

              {/* 説明 */}
              <p
                className="text-sm leading-relaxed"
                style={{ color: "var(--juku-text-muted)" }}
              >
                {step.description}
              </p>

              {/* モバイル矢印 */}
              {idx < steps.length - 1 && (
                <div
                  className="lg:hidden text-lg mt-1"
                  style={{ color: "var(--juku-gold-dim)" }}
                >
                  ↓
                </div>
              )}
            </div>
          ))}
        </div>
      </div>
    </section>
  );
}
