import type { JSX } from "react";
import { cn } from "@/lib/utilities";

export function Divider({ className }: { className?: string }): JSX.Element {
  return (
    <div
      className={cn(
        "h-px w-full",
        "bg-[linear-gradient(90deg,transparent,var(--primary-dim)_30%,var(--primary-dim)_70%,transparent)]",
        "opacity-50",
        className,
      )}
    />
  );
}
