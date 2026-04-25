import type { JSX } from "react";
import { cn } from "@/lib/utilities";

const sizeMap = {
  default: { jukubox: "text-xl", dotAi: "text-sm" },
  lg: { jukubox: "text-3xl", dotAi: "text-lg" },
  exlg: { jukubox: "text-7xl lg:text-8xl", dotAi: "text-4xl lg:text-5xl" },
} as const;

type Size = keyof typeof sizeMap;

export function JukuBoxLogo({
  size = "default",
}: {
  size?: Size;
}): JSX.Element {
  const { jukubox, dotAi } = sizeMap[size];

  return (
    <span className="font-orbitron flex items-baseline gap-0.5">
      <span
        className={cn(
          "text-primary font-black [text-shadow:0_0_18px_oklch(0.75_0.12_77/0.45)]",
          jukubox,
        )}
      >
        JukuBox
      </span>
      <span className={cn("text-secondary font-bold", dotAi)}>.ai</span>
    </span>
  );
}
