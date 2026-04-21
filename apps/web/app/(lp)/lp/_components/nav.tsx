import { Button } from "@/components/ui/button";
import Link from "next/link";
import { JukuBoxLogo } from "@/components/jukubox-logo";

export function Nav() {
  return (
    <nav className="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-8 py-4 bg-background/20 backdrop-blur-xl border-b border-primary/10">
      {/* ロゴ */}
      <Link href="/" className="no-underline">
        <JukuBoxLogo />
      </Link>

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
      <Button
        size="sm"
        nativeButton={false}
        render={<Link href="/login" prefetch={false} />}
      >
        ログイン
      </Button>
    </nav>
  );
}
