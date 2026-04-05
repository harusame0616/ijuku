import { cn } from "@/lib/utils";

export function LandingDivider({ className }: { className?: string }) {
  return (
    <div
      className={cn(
        "w-full h-px",
        "bg-[linear-gradient(90deg,transparent,var(--primary-dim)_30%,var(--primary-dim)_70%,transparent)]",
        "opacity-50",
        className
      )}
    />
  );
}
