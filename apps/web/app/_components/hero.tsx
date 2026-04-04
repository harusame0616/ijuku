import { Button } from "@/components/ui/button";

function VinylRecord() {
  const grooves = Array.from({ length: 26 }, (_, i) => i);
  return (
    <svg
      viewBox="0 0 400 400"
      xmlns="http://www.w3.org/2000/svg"
      className="w-full h-full"
    >
      <title>record</title>
      {/* ベース */}
      <circle cx="200" cy="200" r="195" fill="oklch(0.12 0.02 50)" />
      {/* 外縁ハイライト */}
      <circle
        cx="200"
        cy="200"
        r="194"
        fill="none"
        stroke="oklch(0.75 0.12 77 / 0.28)"
        strokeWidth="1"
      />
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
              ? "oklch(0.75 0.12 77 / 0.22)"
              : "oklch(0.75 0.12 77 / 0.08)"
          }
          strokeWidth={i % 4 === 0 ? "0.8" : "0.4"}
        />
      ))}
      {/* センターラベルエリア */}
      <circle cx="200" cy="200" r="60" fill="oklch(0.20 0.03 55)" />
      <circle
        cx="200"
        cy="200"
        r="57"
        fill="none"
        stroke="oklch(0.75 0.12 77 / 0.55)"
        strokeWidth="0.8"
      />
      <circle
        cx="200"
        cy="200"
        r="52"
        fill="none"
        stroke="oklch(0.75 0.12 77 / 0.28)"
        strokeWidth="0.4"
      />
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
        stroke="oklch(0.75 0.12 77 / 0.40)"
        strokeWidth="1.5"
        strokeLinecap="round"
      />
    </svg>
  );
}

export function Hero() {
  return (
    <section className="relative min-h-screen flex items-center overflow-hidden bg-background [background-image:linear-gradient(var(--grid-line)_1px,transparent_1px),linear-gradient(90deg,var(--grid-line)_1px,transparent_1px)] [background-size:48px_48px]">
      {/* 背景ラジアルグロウ（控えめ） */}
      <div className="absolute inset-0 pointer-events-none bg-[radial-gradient(ellipse_55%_45%_at_68%_50%,oklch(0.75_0.12_77/0.04)_0%,transparent_65%)]" />

      <div className="relative z-10 w-full max-w-7xl mx-auto px-8 pt-24 pb-16 flex flex-col lg:flex-row items-center gap-16 lg:gap-0">
        {/* 左: テキスト */}
        <div className="flex-1 flex flex-col gap-7 lg:pr-16">
          {/* 小ラベル */}
          <div className="flex items-center gap-3 text-muted-foreground">
            <div className="w-8 h-px bg-primary-dim" />
            <span className="font-space-mono text-xs uppercase tracking-[0.25em]">
              AI Learning Platform
            </span>
          </div>

          {/* ロゴ */}
          <h1 className="flex flex-col gap-0">
            <span className="font-orbitron font-black leading-none text-7xl lg:text-8xl [text-shadow:0_0_18px_oklch(0.75_0.12_77/0.45)] text-primary">
              JukuBox
            </span>
            <span className="font-orbitron font-bold leading-none text-4xl lg:text-5xl text-secondary">
              .ai
            </span>
          </h1>

          {/* キャッチコピー（明朝体） */}
          <p className="font-noto-serif-jp font-bold text-2xl lg:text-3xl leading-relaxed text-foreground">
            AI エージェントと
            <br />
            好きなことを好きなだけ学ぶ
          </p>

          {/* 説明 */}
          <p className="text-base leading-loose max-w-lg text-muted-foreground">
            自分が契約している AI エージェントをそのまま活用。
            コースを選んで学ぶも良し、自分で作るも良し。
            記録が積み重なるほど、学びは深くなる。
          </p>

          {/* ボタン */}
          <div className="flex flex-wrap items-center gap-4 mt-1">
            <Button>無料で始める →</Button>
            <Button variant="secondary">デモを見る</Button>
          </div>

          {/* 対応AIモデル */}
          <div className="flex items-center gap-2 mt-2 text-muted-foreground">
            <span className="font-space-mono text-xs">対応:</span>
            {["Claude", "GPT-4o", "Gemini", "Llama"].map((m) => (
              <span
                key={m}
                className="font-space-mono text-xs px-2 py-0.5 border border-border text-[oklch(0.75_0.04_55)]"
              >
                {m}
              </span>
            ))}
          </div>
        </div>

        {/* 右: ビニールレコード */}
        <div className="shrink-0 w-full lg:w-100 relative flex items-center justify-center">
          <div className="relative w-72 h-72 lg:w-88 lg:h-88">
            <div className="animate-[juku-record-spin_32s_linear_infinite,juku-breathe_4s_ease-in-out_infinite] w-full h-full">
              <VinylRecord />
            </div>
            {/* 控えめなアンバーグロウ */}
            <div className="absolute inset-8 rounded-full pointer-events-none shadow-[0_0_50px_oklch(0.75_0.12_77/0.08),0_0_100px_oklch(0.75_0.12_77/0.04)]" />
          </div>
        </div>
      </div>

      {/* スクロールヒント */}
      <div className="absolute bottom-8 left-1/2 -translate-x-1/2 flex flex-col items-center gap-2 text-[oklch(0.72_0.04_55)]">
        <span className="font-space-mono text-xs uppercase tracking-widest">
          Scroll
        </span>
        <div className="w-px h-8 bg-[linear-gradient(to_bottom,var(--primary-dim),transparent)]" />
      </div>
    </section>
  );
}
