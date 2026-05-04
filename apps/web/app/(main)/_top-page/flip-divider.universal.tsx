import type { JSX } from "react";

export function FlipDivider(): JSX.Element {
  return (
    <div
      role="separator"
      aria-label="A 面と B 面の区切り"
      className="flex items-center justify-center gap-4 py-2 text-muted-foreground"
    >
      <div className="h-px w-12 bg-primary-dim sm:w-20" />
      <span className="font-mono text-[0.65rem] uppercase tracking-[0.3em]">
        End&nbsp;of&nbsp;Side&nbsp;A
      </span>
      <span
        aria-hidden="true"
        className="font-orbitron text-2xl font-black text-secondary"
      >
        ↻
      </span>
      <span className="font-mono text-[0.65rem] uppercase tracking-[0.3em]">
        Begin&nbsp;Side&nbsp;B
      </span>
      <div className="h-px w-12 bg-primary-dim sm:w-20" />
    </div>
  );
}
