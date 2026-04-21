import { JukuBoxLogo } from "@/components/jukubox-logo";
import Link from "next/link";
import { Suspense } from "react";
import { UserMenuContainer } from "./user-menu.container.server";
import { UserMenuSkeleton } from "./user-menu.skeleton.universal";

export function Header() {
  return (
    <header className="fixed top-0 left-0 right-0 z-50 flex items-center justify-between px-8 py-4 bg-background/20 backdrop-blur-xl border-b border-primary/10">
      <Link href="/" className="no-underline">
        <JukuBoxLogo />
      </Link>
      <Suspense fallback={<UserMenuSkeleton />}>
        <UserMenuContainer />
      </Suspense>
    </header>
  );
}
