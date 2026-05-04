import Link from "next/link";
import type { JSX } from "react";
import { Button } from "@/components/ui/button";

function BSideVinyl(): JSX.Element {
  return (
    <svg
      viewBox="0 0 200 200"
      xmlns="http://www.w3.org/2000/svg"
      className="size-full"
      aria-hidden="true"
    >
      <title>Vinyl record (Side B)</title>
      <circle cx="100" cy="100" r="98" fill="oklch(0.12 0.02 50)" />
      <circle
        cx="100"
        cy="100"
        r="97"
        fill="none"
        stroke="oklch(0.72 0.09 190 / 0.30)"
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
              ? "oklch(0.72 0.09 190 / 0.20)"
              : "oklch(0.72 0.09 190 / 0.07)"
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
        stroke="oklch(0.72 0.09 190 / 0.50)"
        strokeWidth="0.5"
      />
      <defs>
        <path id="bside-vinyl-arc" d="M 73,100 A 27,27 0 0 1 127,100" />
      </defs>
      <text fontFamily="Orbitron, sans-serif" fontWeight="700" fontSize="4.5">
        <textPath
          href="#bside-vinyl-arc"
          startOffset="50%"
          textAnchor="middle"
        >
          <tspan fill="oklch(0.72 0.09 190)" letterSpacing="0.3">
            JukuBox - TEACH -
          </tspan>
        </textPath>
      </text>
      <circle cx="100" cy="100" r="3.5" fill="oklch(0.09 0.015 50)" />
    </svg>
  );
}

function BSideDisc(): JSX.Element {
  return (
    <div aria-hidden="true" className="relative size-56 lg:size-72">
      {/* 盤面全体を回転 — A 面と対称構造の SVG ヴィニール */}
      <div className="absolute inset-0 animate-[juku-record-spin_28s_linear_infinite] motion-reduce:animate-none">
        <BSideVinyl />
      </div>

      {/* 控えめな cyan グロウ（静止） */}
      <div className="pointer-events-none absolute inset-6 rounded-full shadow-[0_0_60px_oklch(0.72_0.09_190/0.10)]" />
    </div>
  );
}

export function SideBHero(): JSX.Element {
  return (
    <section
      aria-labelledby="side-b-heading"
      className="relative grid grid-cols-1 items-center gap-10 py-12 md:grid-cols-12 md:gap-6 lg:py-16"
    >
      {/* 背景の cyan ラジアルグロウ — Side A の amber と対称 */}
      <div
        aria-hidden="true"
        className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_45%_50%_at_25%_50%,oklch(0.72_0.09_190/0.05)_0%,transparent_70%)]"
      />

      {/* 左カラム（md+）／背景重ね（狭い時は A 面と同じく右から重ねる） */}
      <div
        className={
          "pointer-events-none absolute -right-8 top-1/2 z-0 -translate-y-1/2 opacity-50 " +
          "md:pointer-events-auto md:relative md:right-auto md:top-auto md:z-auto md:translate-y-0 md:opacity-100 " +
          "md:col-span-5 md:flex md:items-center md:justify-start"
        }
      >
        <BSideDisc />
      </div>

      {/* 右カラム: コピー — 狭い時は前面 */}
      <div className="relative z-10 flex flex-col gap-6 md:col-span-7">
        <div className="flex items-center gap-3 text-muted-foreground">
          <span className="font-orbitron text-xs font-bold text-secondary">
            B1
          </span>
          <div className="h-px w-10 bg-primary-dim" />
          <span className="font-mono text-[0.65rem] uppercase tracking-[0.3em]">
            Side B&nbsp;·&nbsp;Teach
          </span>
        </div>

        <h2
          id="side-b-heading"
          className="font-serif text-3xl font-bold leading-[1.18] lg:text-5xl"
        >
          <span className="text-foreground">自分の知識を</span>
          <br />
          <span className="text-secondary">AI と一緒に教えよう。</span>
        </h2>

        <p className="max-w-md text-sm leading-loose text-muted-foreground">
          学んできたこと、 試してきたこと。
          <br />
          それを AI が伴走するコースとして公開できます。
        </p>

        <div className="flex flex-wrap items-center gap-5 pt-1">
          <Button
            variant="secondary"
            nativeButton={false}
            render={<Link href="/courses/new" />}
          >
            コースを作る
          </Button>
        </div>
      </div>
    </section>
  );
}
