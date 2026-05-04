import type { JSX } from "react";
import { FeaturedCourseCard } from "./featured-course-card.universal";
import { featuredCourses } from "./featured-courses.data";

export function FeaturedCoursesSection(): JSX.Element {
  return (
    <section
      id="featured-courses"
      aria-labelledby="featured-courses-heading"
      className="flex flex-col gap-8"
    >
      <header className="border-border flex flex-wrap items-end justify-between gap-4 border-b pb-4">
        <div className="flex flex-col gap-2">
          <div className="text-muted-foreground flex items-center gap-3">
            <span className="font-orbitron text-primary text-xs font-bold">
              A1
            </span>
            <div className="bg-primary-dim h-px w-10" />
            <span className="font-mono text-[0.65rem] tracking-[0.3em] uppercase">
              Featured Courses
            </span>
          </div>
          <h2
            id="featured-courses-heading"
            className="text-foreground font-serif text-3xl font-bold lg:text-4xl"
          >
            注目の講座
          </h2>
        </div>
        <span className="text-muted-foreground hidden font-mono text-[0.65rem] tracking-[0.3em] uppercase md:inline">
          {featuredCourses.length.toString().padStart(2, "0")}&nbsp;TRACKS
        </span>
      </header>

      <ul className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
        {featuredCourses.map((course, index) => (
          <li
            key={
              course.status === "published" ? course.courseSlug : course.title
            }
          >
            <FeaturedCourseCard course={course} trackNumber={index + 1} />
          </li>
        ))}
      </ul>
    </section>
  );
}
