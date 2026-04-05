import { cn } from "@/lib/utils";

const sizeMap = {
  default: { jukubox: "text-xl", dotAi: "text-sm" },
  lg: { jukubox: "text-3xl", dotAi: "text-lg" },
  exlg: { jukubox: "text-7xl lg:text-8xl", dotAi: "text-4xl lg:text-5xl" },
} as const;

type Size = keyof typeof sizeMap;

export function JukuBoxLogo({ size = "default" }: { size?: Size }) {
  const { jukubox, dotAi } = sizeMap[size];

  return (
    <span className="flex items-baseline gap-0.5 font-orbitron">
      <span
        className={cn(
          "font-black [text-shadow:0_0_18px_oklch(0.75_0.12_77/0.45)] text-primary",
          jukubox,
        )}
      >
        JukuBox
      </span>
      <span className={cn("font-bold text-secondary", dotAi)}>.ai</span>
    </span>
  );
}
