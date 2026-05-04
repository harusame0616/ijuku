import type { JSX } from "react";
import { Divider } from "@/components/divider";
import { JukuBoxLogo } from "@/components/jukubox-logo";
import Link from "next/link";
import type { Route } from "next";

const links = [
  // TODO: ページ作成後に URL を変更する
  { label: "ドキュメント", href: "/" },
  { label: "プライバシー", href: "/privacy-policies" },
  // TODO: ページ作成後に URL を変更する
  { label: "利用規約", href: "/" },
] as const satisfies { label: string; href: Route }[];

export function Footer(): JSX.Element {
  return (
    <footer className="py-16 px-8 border-t border-primary/10">
      <div className="max-w-7xl mx-auto flex flex-col items-center gap-10">
        <div className="flex flex-col items-center gap-2">
          <JukuBoxLogo size="lg" />
          <p className="font-serif text-xs text-muted-foreground">
            あなたの AI で好きなことを好きなだけ学ぶ
          </p>
        </div>

        {/* リンク */}
        <nav className="flex flex-wrap justify-center gap-x-8 gap-y-3">
          {links.map((link) => (
            <Link
              key={link.label}
              href={link.href}
              className="text-xs text-muted-foreground underline underline-offset-2"
              prefetch={false}
            >
              {link.label}
            </Link>
          ))}
        </nav>

        {/* セパレーター */}
        <Divider />

        <p className="font-mono text-xs text-muted-foreground">
          © 2025 JukuBox.ai — All rights reserved.
        </p>
      </div>
    </footer>
  );
}
