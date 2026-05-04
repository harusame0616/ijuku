import type { JSX } from "react";
import { Button } from "@/components/ui/button";

function CompactVinyl(): JSX.Element {
  return (
    <svg
      viewBox="0 0 200 200"
      xmlns="http://www.w3.org/2000/svg"
      className="size-full"
      aria-hidden="true"
    >
      <title>Vinyl record</title>
      <circle cx="100" cy="100" r="98" fill="oklch(0.12 0.02 50)" />
      <circle
        cx="100"
        cy="100"
        r="97"
        fill="none"
        stroke="oklch(0.75 0.12 77 / 0.30)"
        strokeWidth="0.6"
      />
      {Array.from({ length: 18 }, (_, index) => index).map((index) => (
        <circle
          key={index}
          cx="100"
          cy="100"
          r={92 - index * 4.2}
          fill="none"
          stroke={
            index % 4 === 0
              ? "oklch(0.75 0.12 77 / 0.20)"
              : "oklch(0.75 0.12 77 / 0.07)"
          }
          strokeWidth={index % 4 === 0 ? "0.6" : "0.3"}
        />
      ))}
      <circle cx="100" cy="100" r="36" fill="oklch(0.20 0.03 55)" />
      <circle
        cx="100"
        cy="100"
        r="34"
        fill="none"
        stroke="oklch(0.75 0.12 77 / 0.50)"
        strokeWidth="0.5"
      />
      <defs>
        <path id="compact-vinyl-arc" d="M 73,100 A 27,27 0 0 1 127,100" />
      </defs>
      <text fontFamily="Orbitron, sans-serif" fontWeight="700" fontSize="4.5">
        <textPath
          href="#compact-vinyl-arc"
          startOffset="50%"
          textAnchor="middle"
        >
          <tspan fill="oklch(0.75 0.12 77)" letterSpacing="0.3">
            JukuBox - LEARN -
          </tspan>
        </textPath>
      </text>
      <path
        d="M 60 70 A 50 50 0 0 1 120 50"
        fill="none"
        stroke="oklch(0.75 0.12 77 / 0.45)"
        strokeWidth="1.2"
        strokeLinecap="round"
      />
      <circle cx="100" cy="100" r="3.5" fill="oklch(0.09 0.015 50)" />
    </svg>
  );
}

export function Hero(): JSX.Element {
  return (
    <section className="border-border/60 relative -mx-4 overflow-hidden border-b">
      {/* 微細スキャンライン — レトロフューチャーのテクスチャ。値はデザイントークン外（演出用） */}
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 bg-[repeating-linear-gradient(0deg,transparent_0,transparent_3px,oklch(0.75_0.12_77)_3px,oklch(0.75_0.12_77)_4px)] opacity-[0.035]"
      />

      {/* ラジアルグロウ — 右上方向に視線誘導、レコード盤の存在感を強調 */}
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_45%_50%_at_75%_45%,oklch(0.75_0.12_77/0.06)_0%,transparent_70%)]"
      />

      <div className="relative mx-auto grid max-w-6xl grid-cols-1 items-center gap-10 px-4 py-16 md:grid-cols-12 md:gap-6 lg:py-24">
        {/* 左カラム: コピー — 狭い時は前面に立たせる */}
        <div className="relative z-10 flex flex-col gap-6 md:col-span-7">
          <div className="text-muted-foreground flex items-center gap-3">
            <div className="bg-primary-dim h-px w-10" />
            <span className="font-mono text-[0.65rem] tracking-[0.3em] uppercase">
              Side A&nbsp;·&nbsp;Learn
            </span>
          </div>

          <h1 className="font-serif text-3xl leading-[1.18] font-bold lg:text-5xl">
            <span className="text-foreground">好きなことを</span>
            <br />
            <span className="text-primary">自分の AI で学ぼう。</span>
          </h1>

          <p className="text-muted-foreground max-w-md text-sm leading-loose">
            いつもの AI が、あなたの先生になる。
            <br />
            コースを選んで、あなたが契約している AI で受講を始められます。
          </p>

          <div className="flex flex-wrap items-center gap-5 pt-1">
            <Button
              nativeButton={false}
              render={<a href="#featured-courses" />}
            >
              注目の講座を見る
            </Button>
          </div>
        </div>

        {/* 右カラム: 回転するコンパクトヴィニール — 狭い時は背景重ね、 md+ は通常 grid 配置 */}
        <div
          className={
            // 狭い時: absolute で右からはみ出す半透明背景。 md+: grid item に戻して通常表示
            "pointer-events-none absolute top-1/2 -right-8 z-0 -translate-y-1/2 opacity-80 " +
            "md:pointer-events-auto md:relative md:top-auto md:right-auto md:z-auto md:translate-y-0 md:opacity-100 " +
            "md:col-span-5 md:flex md:items-center md:justify-end"
          }
        >
          <div className="relative">
            <div className="relative size-56 motion-reduce:animate-none lg:size-72">
              <div className="size-full animate-[juku-record-spin_28s_linear_infinite] motion-reduce:animate-none">
                <CompactVinyl />
              </div>
              <div
                aria-hidden="true"
                className="pointer-events-none absolute inset-6 rounded-full shadow-[0_0_60px_oklch(0.75_0.12_77/0.10)]"
              />
            </div>
          </div>
        </div>
      </div>
    </section>
  );
}
