import { JukuBoxLogo } from "@/components/jukubox-logo";
import { GridSection } from "@/components/grid-section";
import { Divider } from "@/components/divider";
import Link from "next/link";
import { OAuthLoginButton } from "./oauth-login-button";

export default function Page() {
  return (
    <GridSection className="relative min-h-screen flex items-center justify-center px-4">
      <div className="absolute inset-0 pointer-events-none bg-[radial-gradient(ellipse_60%_55%_at_50%_50%,oklch(0.75_0.12_77/0.06)_0%,transparent_65%)]" />

      <div className="relative z-10 flex flex-col items-center gap-10 w-full max-w-sm">
        <div className="flex items-center gap-3">
          <div className="w-8 h-px bg-primary-dim" />
          <span className="font-mono text-xs uppercase tracking-[0.25em] text-muted-foreground">
            Sign In
          </span>
          <div className="w-8 h-px bg-primary-dim" />
        </div>

        <Link href="/">
          <JukuBoxLogo size="lg" />
        </Link>

        <div className="w-full flex flex-col items-center gap-6 rounded-2xl border border-border bg-card backdrop-blur-sm px-8 py-10">
          <p className="text-sm text-muted-foreground text-center leading-relaxed">
            以下のアカウントでログインしてください
          </p>
          <Divider />
          <OAuthLoginButton />
        </div>
      </div>
    </GridSection>
  );
}
