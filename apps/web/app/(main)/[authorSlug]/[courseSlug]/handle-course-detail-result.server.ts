import { notFound } from "next/navigation";
import type { CourseDetail, GetCourseDetailResult } from "./course-detail.data";

export function handleGetCourseDetailResult(
  result: GetCourseDetailResult,
): CourseDetail {
  if (result.success) return result.course;

  switch (result.code) {
    case "NOT_FOUND": {
      notFound();
      break;
    }
    case "INTERNAL_ERROR": {
      throw new Error("講座詳細の取得に失敗しました");
    }
    default: {
      const exhaustiveCheck: never = result.code;
      throw new Error(`未対応のエラーコード: ${exhaustiveCheck}`);
    }
  }
}
