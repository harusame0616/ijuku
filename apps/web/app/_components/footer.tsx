const links = [
  { label: "機能",       href: "#機能" },
  { label: "使い方",     href: "#使い方" },
  { label: "学習記録",   href: "#学習記録" },
  { label: "ドキュメント", href: "#" },
  { label: "プライバシー", href: "#" },
  { label: "利用規約",   href: "#" },
];

export function Footer() {
  return (
    <footer
      className="relative py-16 px-8"
      style={{
        background: "var(--juku-bg)",
        borderTop: "1px solid var(--juku-gold-dim)",
        borderTopColor: "oklch(0.75 0.12 77 / 0.1)",
      }}
    >
      <div className="max-w-7xl mx-auto flex flex-col items-center gap-10">
        {/* ロゴ */}
        <div className="flex flex-col items-center gap-2">
          <div className="flex items-baseline gap-0.5">
            <span
              className="font-orbitron font-black text-2xl juku-glow-gold-text"
              style={{ color: "var(--juku-gold)" }}
            >
              JukuBox
            </span>
            <span
              className="font-orbitron font-bold text-sm"
              style={{ color: "var(--juku-teal)" }}
            >
              .ai
            </span>
          </div>
          <p
            className="font-noto-serif-jp text-xs text-center"
            style={{ color: "var(--juku-text-muted)" }}
          >
            AI エージェントと好きなことを好きなだけ学ぶ
          </p>
        </div>

        {/* リンク */}
        <nav className="flex flex-wrap justify-center gap-x-8 gap-y-3">
          {links.map((link) => (
            <a
              key={link.label}
              href={link.href}
              className="text-xs transition-colors duration-200"
              style={{ color: "oklch(0.38 0.02 55)" }}
            >
              {link.label}
            </a>
          ))}
        </nav>

        {/* セパレーター */}
        <div className="juku-divider w-full" />

        {/* コピーライト */}
        <div className="flex flex-col sm:flex-row items-center gap-4">
          <p
            className="font-space-mono text-xs"
            style={{ color: "oklch(0.32 0.02 55)" }}
          >
            © 2025 JukuBox.ai — All rights reserved.
          </p>
          <div className="flex items-center gap-2">
            {["Claude", "GPT-4o", "Gemini"].map((model) => (
              <span
                key={model}
                className="font-space-mono text-[10px] px-2 py-0.5"
                style={{
                  border: "1px solid oklch(0.75 0.12 77 / 0.12)",
                  color: "oklch(0.32 0.02 55)",
                }}
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
