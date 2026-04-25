import type { JSX } from "react";
import { Header } from "./_components/header.server";

export default function MainLayout({
  children,
}: {
  children: React.ReactNode;
}): JSX.Element {
  return (
    <div className="grid grid-rows-[auto_1fr] overflow-hidden h-screen">
      <Header />
      <main className="overflow-y-auto p-4 [scrollbar-gutter:stable]">
        {children}
      </main>
    </div>
  );
}
