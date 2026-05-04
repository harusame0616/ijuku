import type { JSX } from "react";
import { getCourseDetail } from "./course-detail.data";
import { CourseDetailPresenter } from "./course-detail.presenter.universal";
import { handleGetCourseDetailResult } from "./handle-course-detail-result.server";

interface Props {
  authorSlug: Promise<string>;
  courseSlug: Promise<string>;
}

export async function CourseDetailContainer({
  authorSlug,
  courseSlug,
}: Props): Promise<JSX.Element> {
  const [resolvedAuthorSlug, resolvedCourseSlug] = await Promise.all([
    authorSlug,
    courseSlug,
  ]);
  const course = handleGetCourseDetailResult(
    await getCourseDetail(resolvedAuthorSlug, resolvedCourseSlug),
  );

  return <CourseDetailPresenter course={course} />;
}
