import { CtaPrimaryButton, CtaSecondaryButton } from "@/app/_components/ui/cta-button";
import { LandingDivider } from "@/app/_components/ui/landing-divider";

export function CtaSection() {
  return (
    <section className="relative py-32 px-8 overflow-hidden bg-background [background-image:linear-gradient(var(--grid-line)_1px,transparent_1px),linear-gradient(90deg,var(--grid-line)_1px,transparent_1px)] [background-size:48px_48px]">
      {/* ラジアルグロウ（控えめ） */}
      <div className="absolute inset-0 pointer-events-none bg-[radial-gradient(ellipse_60%_55%_at_50%_50%,oklch(0.75_0.12_77/0.04)_0%,transparent_65%)]" />

      <div className="relative z-10 max-w-3xl mx-auto flex flex-col items-center gap-10 text-center">
        {/* ラベル */}
        <div className="flex items-center gap-3">
          <div className="w-8 h-px bg-gold-dim" />
          <span className="font-space-mono text-xs uppercase tracking-[0.25em] text-muted-foreground">
            Get Started
          </span>
          <div className="w-8 h-px bg-gold-dim" />
        </div>

        {/* 見出し */}
        <h2 className="font-noto-serif-jp font-black text-4xl lg:text-6xl leading-snug text-foreground">
          さあ、
          <br />
          <span className="text-gold [text-shadow:0_0_18px_oklch(0.75_0.12_77/0.45)]">
            学びを始めよう。
          </span>
        </h2>

        <p className="text-base leading-relaxed max-w-lg text-muted-foreground">
          今すぐ無料でアカウントを作成。 自分だけの学習スタジオを立ち上げよう。
        </p>

        {/* ボタン群 */}
        <div className="flex flex-wrap items-center gap-4">
          <CtaPrimaryButton
            type="button"
            className="animate-[juku-glow-gold_4s_ease-in-out_infinite] px-10 py-4 font-orbitron font-black text-sm uppercase tracking-widest"
          >
            無料で始める —
          </CtaPrimaryButton>
          <CtaSecondaryButton
            type="button"
            className="px-10 py-4 font-orbitron font-bold text-sm uppercase tracking-widest"
          >
            デモを見る
          </CtaSecondaryButton>
        </div>

        {/* 補足 */}
        <p className="font-space-mono text-xs text-[oklch(0.35_0.02_55)]">
          クレジットカード不要 · いつでもキャンセル可 · 商用利用可
        </p>

        {/* セパレーター */}
        <LandingDivider className="max-w-sm mt-2" />
      </div>
    </section>
  );
}
