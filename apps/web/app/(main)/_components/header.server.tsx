import type { JSX } from "react";
import { JukuBoxLogo } from "@/components/jukubox-logo";
import Link from "next/link";
import { Suspense } from "react";
import { HeaderSearch } from "./header-search.universal";
import { UserMenuContainer } from "./user-menu.container.server";
import { UserMenuSkeleton } from "./user-menu.skeleton.universal";

export function Header(): JSX.Element {
  return (
    <header className="flex items-center gap-4 border-b border-primary/10 bg-background/20 px-4 py-4 backdrop-blur-xl md:gap-6 md:px-8">
      <Link href="/" className="shrink-0 no-underline">
        <JukuBoxLogo />
      </Link>
      <div className="flex flex-1 items-center justify-end md:justify-center">
        <HeaderSearch />
      </div>
      <Suspense fallback={<UserMenuSkeleton />}>
        <UserMenuContainer />
      </Suspense>
    </header>
  );
}
