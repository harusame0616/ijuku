import type { JSX } from "react";
import { Skeleton } from "@/components/ui/skeleton";

const SKELETON_COUNT = 3;

export function EnrollmentsListSkeleton(): JSX.Element {
  return (
    <ul className="flex flex-col gap-3">
      {Array.from({ length: SKELETON_COUNT }, (_, index) => (
        <li
          // biome-ignore lint/suspicious/noArrayIndexKey: 静的個数のスケルトン
          key={index}
          className="rounded-md border border-border bg-card p-4"
        >
          <Skeleton className="h-5 w-4/5" />
        </li>
      ))}
    </ul>
  );
}
