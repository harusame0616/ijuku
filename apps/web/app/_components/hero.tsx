function VinylRecord() {
  const grooves = Array.from({ length: 26 }, (_, i) => i);
  return (
    <svg
      viewBox="0 0 400 400"
      xmlns="http://www.w3.org/2000/svg"
      className="w-full h-full"
    >
      {/* ベース */}
      <circle cx="200" cy="200" r="195" fill="oklch(0.12 0.02 50)" />
      {/* 外縁ハイライト */}
      <circle cx="200" cy="200" r="194" fill="none" stroke="oklch(0.75 0.12 77 / 0.12)" strokeWidth="1" />
      {/* 溝 */}
      {grooves.map((i) => (
        <circle
          key={i}
          cx="200"
          cy="200"
          r={180 - i * 5}
          fill="none"
          stroke={
            i % 4 === 0
              ? "oklch(0.75 0.12 77 / 0.1)"
              : "oklch(0.75 0.12 77 / 0.03)"
          }
          strokeWidth={i % 4 === 0 ? "0.8" : "0.4"}
        />
      ))}
      {/* センターラベルエリア */}
      <circle cx="200" cy="200" r="60" fill="oklch(0.14 0.025 55)" />
      <circle cx="200" cy="200" r="57" fill="none" stroke="oklch(0.75 0.12 77 / 0.35)" strokeWidth="0.8" />
      <circle cx="200" cy="200" r="52" fill="none" stroke="oklch(0.75 0.12 77 / 0.12)" strokeWidth="0.4" />
      {/* ラベルテキスト */}
      <text
        x="200"
        y="196"
        textAnchor="middle"
        fontSize="10"
        fontFamily="Orbitron, sans-serif"
        fontWeight="700"
        fill="oklch(0.75 0.12 77)"
        letterSpacing="2.5"
      >
        JUKUBOX
      </text>
      <text
        x="200"
        y="211"
        textAnchor="middle"
        fontSize="7"
        fontFamily="Space Mono, monospace"
        fill="oklch(0.72 0.09 190)"
        letterSpacing="1.5"
      >
        .AI
      </text>
      {/* センターホール */}
      <circle cx="200" cy="200" r="7" fill="oklch(0.09 0.015 50)" />
      {/* 反射ハイライト（斜め） */}
      <path
        d="M 90 120 A 130 130 0 0 1 240 65"
        fill="none"
        stroke="oklch(0.75 0.12 77 / 0.18)"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
    </svg>
  );
}

export function Hero() {
  return (
    <section
      className="relative min-h-screen flex items-center overflow-hidden juku-grid-bg"
      style={{ background: "var(--juku-bg)" }}
    >
      {/* 背景ラジアルグロウ（控えめ） */}
      <div
        className="absolute inset-0 pointer-events-none"
        style={{
          background:
            "radial-gradient(ellipse 55% 45% at 68% 50%, oklch(0.75 0.12 77 / 0.04) 0%, transparent 65%)",
        }}
      />

      <div className="relative z-10 w-full max-w-7xl mx-auto px-8 pt-24 pb-16 flex flex-col lg:flex-row items-center gap-16 lg:gap-0">
        {/* 左: テキスト */}
        <div className="flex-1 flex flex-col gap-7 lg:pr-16">
          {/* 小ラベル */}
          <div
            className="flex items-center gap-3"
            style={{ color: "var(--juku-text-muted)" }}
          >
            <div
              className="w-8 h-px"
              style={{ background: "var(--juku-gold-dim)" }}
            />
            <span className="font-space-mono text-xs uppercase tracking-[0.25em]">
              AI Learning Platform
            </span>
          </div>

          {/* ロゴ */}
          <h1 className="flex flex-col gap-0">
            <span
              className="font-orbitron font-black leading-none text-7xl lg:text-8xl juku-glow-gold-text"
              style={{ color: "var(--juku-gold)" }}
            >
              JukuBox
            </span>
            <span
              className="font-orbitron font-bold leading-none text-4xl lg:text-5xl"
              style={{ color: "var(--juku-teal)" }}
            >
              .ai
            </span>
          </h1>

          {/* キャッチコピー（明朝体） */}
          <p
            className="font-noto-serif-jp font-bold text-2xl lg:text-3xl leading-relaxed"
            style={{ color: "var(--juku-text)" }}
          >
            AI エージェントと
            <br />
            好きなことを好きなだけ学ぶ
          </p>

          {/* 説明 */}
          <p
            className="text-base leading-loose max-w-lg"
            style={{ color: "var(--juku-text-muted)" }}
          >
            自分が契約している AI エージェントをそのまま活用。
            コースを選んで学ぶも良し、自分で作るも良し。
            記録が積み重なるほど、学びは深くなる。
          </p>

          {/* ボタン */}
          <div className="flex flex-wrap items-center gap-4 mt-1">
            <button className="juku-cta-btn-primary px-8 py-3.5 font-orbitron font-bold text-sm uppercase tracking-widest">
              <span>無料で始める →</span>
            </button>
            <button className="juku-cta-btn-secondary px-8 py-3.5 font-orbitron font-bold text-sm uppercase tracking-widest">
              デモを見る
            </button>
          </div>

          {/* 対応AIモデル */}
          <div
            className="flex items-center gap-2 mt-2"
            style={{ color: "var(--juku-text-muted)" }}
          >
            <span className="font-space-mono text-xs">対応:</span>
            {["Claude", "GPT-4o", "Gemini", "Llama"].map((m) => (
              <span
                key={m}
                className="font-space-mono text-xs px-2 py-0.5"
                style={{
                  border: "1px solid var(--juku-border)",
                  color: "oklch(0.42 0.02 55)",
                }}
              >
                {m}
              </span>
            ))}
          </div>
        </div>

        {/* 右: ビニールレコード */}
        <div className="flex-shrink-0 w-full lg:w-[400px] relative flex items-center justify-center">
          <div className="relative w-72 h-72 lg:w-88 lg:h-88">
            <div className="juku-record-spin w-full h-full juku-breathe">
              <VinylRecord />
            </div>
            {/* 控えめなアンバーグロウ */}
            <div
              className="absolute inset-8 rounded-full pointer-events-none"
              style={{
                boxShadow:
                  "0 0 50px oklch(0.75 0.12 77 / 0.08), 0 0 100px oklch(0.75 0.12 77 / 0.04)",
              }}
            />
          </div>
        </div>
      </div>

      {/* スクロールヒント */}
      <div
        className="absolute bottom-8 left-1/2 -translate-x-1/2 flex flex-col items-center gap-2"
        style={{ color: "oklch(0.38 0.02 55)" }}
      >
        <span className="font-space-mono text-xs uppercase tracking-widest">
          Scroll
        </span>
        <div
          className="w-px h-8"
          style={{
            background:
              "linear-gradient(to bottom, var(--juku-gold-dim), transparent)",
          }}
        />
      </div>
    </section>
  );
}
