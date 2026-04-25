import type { JSX, PropsWithChildren } from "react";

interface Props {
  title: string;
}
export default function PageLayout({
  title,
  children,
}: PropsWithChildren<Props>): JSX.Element {
  return (
    <div className="mx-auto max-w-2xl">
      <h1 className="mb-4 text-xl font-bold text-foreground">{title}</h1>
      {children}
    </div>
  );
}
