import type { JSX } from "react";
import { Skeleton } from "@/components/ui/skeleton";

export function CourseDetailSkeleton(): JSX.Element {
  return (
    <article className="flex flex-col gap-6">
      <header className="flex flex-col gap-3">
        <Skeleton className="h-8 w-3/4" />
        <Skeleton className="h-5 w-32" />
        <div className="flex flex-col gap-1">
          <Skeleton className="h-5 w-full" />
          <Skeleton className="h-5 w-5/6" />
        </div>
        <ul className="flex flex-wrap gap-1">
          {Array.from({ length: 3 }, (_, tagIndex) => (
            <li
              // biome-ignore lint/suspicious/noArrayIndexKey: 静的個数のスケルトン
              key={tagIndex}
            >
              <Skeleton className="h-5 w-12 rounded-full" />
            </li>
          ))}
        </ul>
      </header>

      <div className="flex flex-col gap-2">
        <Skeleton className="h-4 w-32" />
        <div className="grid grid-cols-[1fr_auto] items-stretch gap-2">
          <Skeleton className="h-7" />
          <Skeleton className="size-7" />
        </div>
      </div>

      <section className="flex flex-col gap-3">
        <Skeleton className="h-7 w-24" />
        <ul className="flex flex-col gap-3">
          {Array.from({ length: 3 }, (_, sectionIndex) => (
            <li
              // biome-ignore lint/suspicious/noArrayIndexKey: 静的個数のスケルトン
              key={sectionIndex}
              className="rounded-md border border-border bg-card p-4"
            >
              <Skeleton className="h-6 w-2/3" />
              <Skeleton className="mt-1 h-4 w-1/2" />
              <ul className="mt-3 flex flex-col gap-2 pl-3">
                {Array.from({ length: 2 }, (_, topicIndex) => (
                  <li
                    // biome-ignore lint/suspicious/noArrayIndexKey: 静的個数のスケルトン
                    key={topicIndex}
                    className="flex flex-col gap-1"
                  >
                    <Skeleton className="h-5 w-1/3" />
                    <Skeleton className="h-4 w-1/2" />
                  </li>
                ))}
              </ul>
            </li>
          ))}
        </ul>
      </section>
    </article>
  );
}
