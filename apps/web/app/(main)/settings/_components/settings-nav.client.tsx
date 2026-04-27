"use client";

import type { Route } from "next";
import Link from "next/link";
import { usePathname } from "next/navigation";
import type { JSX } from "react";
import { cn } from "@/lib/utilities";

interface SettingsNavItem {
  href: Route;
  label: string;
}

const items: SettingsNavItem[] = [
  { href: "/settings/profile", label: "プロフィール" },
  { href: "/settings/api-keys", label: "API キー" },
];

function isItemActive(pathname: string, href: string): boolean {
  return pathname === href || pathname.startsWith(href + "/");
}

export function SettingsNav(): JSX.Element {
  const pathname = usePathname();
  return (
    <nav aria-labelledby="settings-nav-label">
      <h2 id="settings-nav-label" className="sr-only">
        設定メニュー
      </h2>
      <ul className="flex gap-1 overflow-x-auto md:flex-col md:gap-0.5">
        {items.map((item) => {
          const active = isItemActive(pathname, item.href);
          return (
            <li key={item.href}>
              <Link
                href={item.href}
                aria-current={active ? "page" : undefined}
                className={cn(
                  "block whitespace-nowrap rounded-md px-3 py-2 text-sm text-foreground hover:bg-accent hover:text-accent-foreground",
                  active &&
                    "bg-accent font-medium text-accent-foreground border-l-2 border-primary",
                )}
              >
                {item.label}
              </Link>
            </li>
          );
        })}
      </ul>
    </nav>
  );
}
