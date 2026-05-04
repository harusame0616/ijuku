import type { JSX } from "react";
import { Suspense } from "react";
import { CourseDetailContainer } from "./course-detail.container.server";
import { CourseDetailSkeleton } from "./course-detail.skeleton.universal";

export default function CourseDetailPage({
  params,
}: PageProps<"/[authorSlug]/[courseSlug]">): JSX.Element {
  return (
    <div className="mx-auto max-w-2xl">
      <Suspense fallback={<CourseDetailSkeleton />}>
        <CourseDetailContainer
          authorSlug={params.then((p) => p.authorSlug)}
          courseSlug={params.then((p) => p.courseSlug)}
        />
      </Suspense>
    </div>
  );
}
