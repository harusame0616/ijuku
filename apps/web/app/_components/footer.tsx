import { LandingDivider } from "@/app/_components/ui/landing-divider";

const links = [
  { label: "機能", href: "#機能" },
  { label: "使い方", href: "#使い方" },
  { label: "学習記録", href: "#学習記録" },
  { label: "ドキュメント", href: "#" },
  { label: "プライバシー", href: "#" },
  { label: "利用規約", href: "#" },
];

export function Footer() {
  return (
    <footer className="relative py-16 px-8 bg-background border-t border-[oklch(0.75_0.12_77/0.1)]">
      <div className="max-w-7xl mx-auto flex flex-col items-center gap-10">
        {/* ロゴ */}
        <div className="flex flex-col items-center gap-2">
          <div className="flex items-baseline gap-0.5">
            <span className="font-orbitron font-black text-2xl [text-shadow:0_0_18px_oklch(0.75_0.12_77/0.45)] text-gold">
              JukuBox
            </span>
            <span className="font-orbitron font-bold text-sm text-teal">
              .ai
            </span>
          </div>
          <p className="font-noto-serif-jp text-xs text-center text-muted-foreground">
            AI エージェントと好きなことを好きなだけ学ぶ
          </p>
        </div>

        {/* リンク */}
        <nav className="flex flex-wrap justify-center gap-x-8 gap-y-3">
          {links.map((link) => (
            <a
              key={link.label}
              href={link.href}
              className="text-xs transition-colors duration-200 text-[oklch(0.38_0.02_55)]"
            >
              {link.label}
            </a>
          ))}
        </nav>

        {/* セパレーター */}
        <LandingDivider />

        {/* コピーライト */}
        <div className="flex flex-col sm:flex-row items-center gap-4">
          <p className="font-space-mono text-xs text-[oklch(0.32_0.02_55)]">
            © 2025 JukuBox.ai — All rights reserved.
          </p>
          <div className="flex items-center gap-2">
            {["Claude", "GPT-4o", "Gemini"].map((model) => (
              <span
                key={model}
                className="font-space-mono text-[10px] px-2 py-0.5 border border-[oklch(0.75_0.12_77/0.12)] text-[oklch(0.32_0.02_55)]"
              >
                {model}
              </span>
            ))}
          </div>
        </div>
      </div>
    </footer>
  );
}
