import type { JSX } from "react";
import { JukuBoxLogo } from "@/components/jukubox-logo";
import Link from "next/link";
import { Suspense } from "react";
import { UserMenuContainer } from "./user-menu.container.server";
import { UserMenuSkeleton } from "./user-menu.skeleton.universal";

export function Header(): JSX.Element {
  return (
    <header className="flex items-center justify-between px-8 py-4 bg-background/20 backdrop-blur-xl border-b border-primary/10">
      <Link href="/" className="no-underline">
        <JukuBoxLogo />
      </Link>
      <Suspense fallback={<UserMenuSkeleton />}>
        <UserMenuContainer />
      </Suspense>
    </header>
  );
}
