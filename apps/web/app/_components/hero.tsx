import { Button } from "@/components/ui/button";
import { GridSection } from "@/app/_components/ui/grid-section";
import { JukuBoxLogo } from "./jukubox-logo";
import Link from "next/link";

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
      {/* ラベルテキスト（縁に沿って湾曲） */}
      <defs>
        <path id="record-label-arc" d="M 162,200 A 38,38 0 0 1 238,200" />
      </defs>
      <text fontFamily="Orbitron, sans-serif" fontWeight="700" fontSize="10">
        <textPath
          href="#record-label-arc"
          startOffset="50%"
          textAnchor="middle"
        >
          <tspan fill="oklch(0.75 0.12 77)" letterSpacing="1.2">
            JukuBox
          </tspan>
          <tspan fill="oklch(0.72 0.09 190)" letterSpacing="1.2">
            .ai
          </tspan>
        </textPath>
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
    <GridSection className="relative min-h-screen flex items-center overflow-hidden">
      {/* 背景ラジアルグロウ（控えめ） */}
      <div className="absolute inset-0 pointer-events-none bg-[radial-gradient(ellipse_55%_45%_at_68%_50%,oklch(0.75_0.12_77/0.04)_0%,transparent_65%)]" />

      <div className="relative z-10 w-full max-w-7xl mx-auto px-8 pt-24 pb-16 flex flex-col lg:flex-row items-center gap-16 lg:gap-0">
        {/* 左: テキスト */}
        <div className="flex-1 flex flex-col gap-7 lg:pr-16">
          {/* 小ラベル */}
          <div className="flex items-center gap-3 text-muted-foreground">
            <div className="w-8 h-px bg-primary-dim" />
            <span className="font-mono text-xs uppercase tracking-[0.25em]">
              AI Learning Platform
            </span>
          </div>

          {/* ロゴ */}
          <h1>
            <JukuBoxLogo size="exlg" />
          </h1>

          {/* キャッチコピー（明朝体） */}
          <p className="font-serif font-bold text-2xl lg:text-3xl leading-relaxed text-foreground">
            あなたの AI で
            <br />
            好きなことを好きなだけ学ぶ
          </p>

          {/* 説明 */}
          <p className="text-base leading-loose max-w-lg text-muted-foreground">
            いつもの AI が、あなたの先生になる。
            <br />
            カリキュラムは選ぶも作るも自由。
            <br />
            使うほどに、学びは深まる。
          </p>

          {/* ボタン */}
          <div className="flex flex-wrap items-center gap-4 mt-1">
            <Button
              nativeButton={false}
              render={<Link href="/register" prefetch={false} />}
            >
              新規登録
            </Button>
            <Button
              variant="outline"
              nativeButton={false}
              render={<Link href="/login" prefetch={false} />}
            >
              ログイン
            </Button>
          </div>

          {/* 対応AIモデル */}
          <div className="flex items-center gap-2 mt-2 text-muted-foreground">
            <span className="font-mono text-xs">対応:</span>
            {["ChatGPT", "Claude", "Gemini"].map((m) => (
              <span
                key={m}
                className="font-mono text-xs px-2 py-0.5 border border-border text-muted-foreground"
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
      <div className="absolute bottom-8 left-1/2 -translate-x-1/2 flex flex-col items-center gap-2 text-muted-foreground">
        <span className="font-mono text-xs uppercase tracking-widest">
          Scroll
        </span>
        <div className="w-px h-8 bg-[linear-gradient(to_bottom,var(--primary-dim),transparent)]" />
      </div>
    </GridSection>
  );
}
