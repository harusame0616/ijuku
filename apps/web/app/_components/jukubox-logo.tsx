import { cn } from "@/lib/utils";

const sizeMap = {
  default: { jukubox: "text-xl", dotAi: "text-sm" },
  lg: { jukubox: "text-3xl", dotAi: "text-lg" },
} as const satisfies Record<string, { jukubox: string; dotAi: string }>;

type Size = keyof typeof sizeMap;

export function JukuBoxLogo({ size = "default" }: { size?: Size }) {
  const s = sizeMap[size];
  return (
    <div className="flex items-baseline gap-0.5">
      <span
        className={cn(
          "font-orbitron font-black [text-shadow:0_0_18px_oklch(0.75_0.12_77/0.45)] text-primary",
          s.jukubox,
        )}
      >
        JukuBox
      </span>
      <span className={cn("font-orbitron font-bold text-secondary", s.dotAi)}>
        .ai
      </span>
    </div>
  );
}
