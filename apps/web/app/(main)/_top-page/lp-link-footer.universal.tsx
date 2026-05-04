import Link from "next/link";
import type { JSX } from "react";

export function LpLinkFooter(): JSX.Element {
  return (
    <section
      aria-labelledby="about-jukubox-heading"
      className="flex flex-col items-center gap-3 border-t border-border pt-10 text-center"
    >
      <span className="font-mono text-[0.65rem] uppercase tracking-[0.3em] text-muted-foreground">
        Liner&nbsp;Notes
      </span>
      <h2
        id="about-jukubox-heading"
        className="font-serif text-base font-bold text-foreground"
      >
        JukuBox についてもっと知る
      </h2>
      <Link
        href="/lp"
        className="group/lp inline-flex items-center gap-2 font-mono text-[0.7rem] uppercase tracking-[0.25em] text-primary underline-offset-4 hover:underline"
      >
        サービス紹介ページへ
        <span
          aria-hidden="true"
          className="transition-transform group-hover/lp:translate-x-1"
        >
          →
        </span>
      </Link>
    </section>
  );
}
