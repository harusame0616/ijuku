import Link from "next/link";
import type { JSX } from "react";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utilities";
import type { FeaturedCourse } from "./featured-courses.data";

interface Props {
  course: FeaturedCourse;
  trackNumber?: number;
}

export function FeaturedCourseCard({ course, trackNumber }: Props): JSX.Element {
  const trackLabel = trackNumber
    ? trackNumber.toString().padStart(2, "0")
    : undefined;

  if (course.status === "published") {
    return (
      <Link
        href={`/${course.authorSlug}/${course.courseSlug}`}
        className={cn(
          "group/card relative block h-full overflow-hidden rounded-md border border-border bg-card p-5",
          "transition-all duration-300",
          "hover:-translate-y-0.5 hover:border-primary/60 hover:shadow-[0_10px_30px_oklch(0.75_0.12_77/0.10)]",
          "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/40",
        )}
      >
        <GrooveDecoration tone="active" />
        <CardBody course={course} trackLabel={trackLabel} />
      </Link>
    );
  }

  return (
    <div
      className={cn(
        "relative h-full cursor-not-allowed overflow-hidden rounded-md border border-dashed border-border bg-card p-5",
      )}
    >
      <GrooveDecoration tone="muted" />
      {/* RESERVED スタンプ — 入荷待ちレコードに押された手書きスタンプの引用 */}
      <span
        aria-hidden="true"
        className="font-orbitron pointer-events-none absolute -right-3 top-7 rotate-[14deg] border-2 border-secondary/60 px-2 py-0.5 text-[0.6rem] font-bold uppercase tracking-[0.2em] text-secondary/70"
      >
        Reserved
      </span>
      <div className="opacity-65">
        <CardBody course={course} trackLabel={trackLabel} />
      </div>
    </div>
  );
}

interface GrooveDecorationProps {
  tone: "active" | "muted";
}

function GrooveDecoration({ tone }: GrooveDecorationProps): JSX.Element {
  const stroke =
    tone === "active"
      ? ["border-primary-dim/35", "border-primary-dim/25", "border-primary-dim/15", "border-primary-dim/10"]
      : ["border-border/60", "border-border/40", "border-border/30", "border-border/20"];

  return (
    <div
      aria-hidden="true"
      className={cn(
        "pointer-events-none absolute -left-16 -bottom-16 size-40 rounded-full border transition-transform duration-500",
        stroke[0],
        tone === "active" && "group-hover/card:scale-105",
      )}
    >
      <div className={cn("absolute inset-3 rounded-full border", stroke[1])} />
      <div className={cn("absolute inset-7 rounded-full border", stroke[2])} />
      <div className={cn("absolute inset-12 rounded-full border", stroke[3])} />
      <div className="absolute left-1/2 top-1/2 size-1.5 -translate-x-1/2 -translate-y-1/2 rounded-full bg-foreground/40" />
    </div>
  );
}

interface CardBodyProps {
  course: FeaturedCourse;
  trackLabel?: string;
}

function CardBody({ course, trackLabel }: CardBodyProps): JSX.Element {
  return (
    <div className="relative z-10 flex h-full flex-col gap-3">
      <div className="flex items-start justify-between gap-2">
        {trackLabel ? (
          <span className="font-orbitron text-3xl font-black leading-none text-primary">
            {trackLabel}
          </span>
        ) : (
          <span aria-hidden="true" />
        )}
        {course.status === "coming-soon" ? (
          <Badge variant="secondary">Coming Soon</Badge>
        ) : null}
      </div>

      <h3 className="mt-1 line-clamp-2 font-serif text-base font-bold text-foreground">
        {course.title}
      </h3>

      <p className="line-clamp-3 text-xs leading-relaxed text-muted-foreground">
        {course.description}
      </p>

      {course.tags.length > 0 ? (
        <ul className="mt-auto flex flex-wrap gap-x-3 gap-y-1 border-t border-border/40 pt-3 font-mono text-[0.65rem] uppercase tracking-[0.18em] text-muted-foreground">
          {course.tags.map((tag) => (
            <li key={tag}>{tag}</li>
          ))}
        </ul>
      ) : null}
    </div>
  );
}
