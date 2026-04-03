export function Nav() {
  return (
    <nav
      className="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-8 py-4"
      style={{
        background: "oklch(0.10 0.015 50 / 0.92)",
        backdropFilter: "blur(24px)",
        WebkitBackdropFilter: "blur(24px)",
        borderBottom: "1px solid oklch(0.75 0.12 77 / 0.1)",
      }}
    >
      {/* ロゴ */}
      <a href="#" className="flex items-baseline gap-0.5 no-underline">
        <span
          className="font-orbitron font-black text-xl tracking-wider juku-glow-gold-text"
          style={{ color: "var(--gold)" }}
        >
          JukuBox
        </span>
        <span
          className="font-orbitron font-bold text-sm"
          style={{ color: "var(--teal)" }}
        >
          .ai
        </span>
      </a>

      {/* ナビリンク */}
      <div className="hidden md:flex items-center gap-8">
        {["機能", "使い方", "学習記録"].map((label) => (
          <a
            key={label}
            href={`#${label}`}
            className="text-sm tracking-wide transition-colors duration-200"
            style={{ color: "var(--muted-foreground)" }}
          >
            {label}
          </a>
        ))}
      </div>

      {/* CTA */}
      <button className="juku-cta-btn-primary px-5 py-2 text-xs font-orbitron font-bold uppercase tracking-widest">
        <span>無料で始める</span>
      </button>
    </nav>
  );
}
