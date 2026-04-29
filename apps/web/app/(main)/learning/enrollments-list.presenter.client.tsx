"use client";

import Link from "next/link";
import type { JSX } from "react";
import { cn } from "@/lib/utilities";
import type { Enrollment } from "./enrollments.data";

interface Props {
  enrollments: Enrollment[];
}

export function EnrollmentsListPresenter({ enrollments }: Props): JSX.Element {
  if (enrollments.length === 0) {
    return (
      <p className="text-muted-foreground text-sm">
        現在受講中の講座はありません
      </p>
    );
  }

  return (
    <ul className="flex flex-col gap-3">
      {enrollments.map((enrollment) => (
        <li key={enrollment.courseId}>
          <Link
            href={`/learning/${enrollment.courseId}`}
            className={cn(
              "border-border bg-card block rounded-md border p-4",
              "hover:border-primary/60 transition-colors",
              "focus-visible:ring-ring/40 focus-visible:ring-2 focus-visible:outline-none",
            )}
          >
            <h2 className="text-foreground line-clamp-2 text-base font-bold">
              {enrollment.title}
            </h2>
          </Link>
        </li>
      ))}
    </ul>
  );
}
