import { CtaPrimaryButton } from "@/app/_components/ui/cta-button";

export function Nav() {
  return (
    <nav className="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-8 py-4 bg-[oklch(0.10_0.015_50_/_0.92)] backdrop-blur-[24px] border-b border-[oklch(0.75_0.12_77_/_0.1)]">
      {/* ロゴ */}
      <a href="#" className="flex items-baseline gap-0.5 no-underline">
        <span className="font-orbitron font-black text-xl tracking-wider [text-shadow:0_0_18px_oklch(0.75_0.12_77/0.45)] text-gold">
          JukuBox
        </span>
        <span className="font-orbitron font-bold text-sm text-teal">.ai</span>
      </a>

      {/* ナビリンク */}
      <div className="hidden md:flex items-center gap-8">
        {["機能", "使い方", "学習記録"].map((label) => (
          <a
            key={label}
            href={`#${label}`}
            className="text-sm tracking-wide transition-colors duration-200 text-muted-foreground"
          >
            {label}
          </a>
        ))}
      </div>

      {/* CTA */}
      <CtaPrimaryButton
        type="button"
        className="px-5 py-2 text-xs font-orbitron font-bold uppercase tracking-widest"
      >
        無料で始める
      </CtaPrimaryButton>
    </nav>
  );
}
