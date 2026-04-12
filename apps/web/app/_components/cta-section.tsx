import { Button } from "@/components/ui/button";
import { LandingDivider } from "@/app/_components/ui/landing-divider";
import { GridSection } from "@/app/_components/ui/grid-section";
import Link from "next/link";

export function CtaSection() {
  return (
    <GridSection className="relative py-32 px-8 overflow-hidden">
      {/* ラジアルグロウ（控えめ） */}
      <div className="absolute inset-0 pointer-events-none bg-[radial-gradient(ellipse_60%_55%_at_50%_50%,oklch(0.75_0.12_77/0.04)_0%,transparent_65%)]" />

      <div className="relative z-10 max-w-3xl mx-auto flex flex-col items-center gap-10 text-center">
        {/* ラベル */}
        <div className="flex items-center gap-3">
          <div className="w-8 h-px bg-primary-dim" />
          <span className="font-mono text-xs uppercase tracking-[0.25em] text-muted-foreground">
            Get Started
          </span>
          <div className="w-8 h-px bg-primary-dim" />
        </div>

        {/* 見出し */}
        <h2 className="font-serif font-black text-4xl lg:text-6xl leading-snug text-foreground">
          さあ、
          <br />
          <span className="text-primary [text-shadow:0_0_18px_oklch(0.75_0.12_77/0.45)]">
            学びを始めよう。
          </span>
        </h2>

        <p className="text-base leading-relaxed max-w-lg text-muted-foreground">
          今すぐ無料でアカウントを作成。 自分だけの学習スタジオを立ち上げよう。
        </p>

        {/* ボタン群 */}
        <div className="flex flex-wrap items-center gap-4">
          <Button
            size="lg"
            className="animate-[juku-glow-gold_4s_ease-in-out_infinite]"
            nativeButton={false}
            render={<Link href="/register" prefetch={false} />}
          >
            新規登録
          </Button>
          <Button
            variant="outline"
            size="lg"
            nativeButton={false}
            render={<Link href="/login" prefetch={false} />}
          >
            ログイン
          </Button>
        </div>

        {/* 補足 */}
        <p className="font-mono text-xs text-subtle-foreground">
          クレジットカード不要 · いつでもキャンセル可 · 商用利用可
        </p>

        {/* セパレーター */}
        <LandingDivider className="max-w-sm mt-2" />
      </div>
    </GridSection>
  );
}
