import type { JSX, ReactNode } from "react";
import { SettingsNav } from "./_components/settings-nav.client";

export default function SettingsLayout({
  children,
}: {
  children: ReactNode;
}): JSX.Element {
  return (
    <div className="mx-auto grid w-full max-w-5xl items-start gap-6 md:grid-cols-[220px_1fr] md:gap-8">
      <aside>
        <SettingsNav />
      </aside>
      <div className="min-w-0">{children}</div>
    </div>
  );
}
