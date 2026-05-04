import { expect, test, vi } from "vitest";
import { handleGetCourseDetailResult } from "./handle-course-detail-result.server";

vi.mock("next/navigation", () => ({
  notFound: () => {
    throw new Error("NEXT_NOT_FOUND");
  },
}));

test("成功結果はそのまま CourseDetail を返す", () => {
  const course = {
    courseId: "11111111-1111-1111-1111-111111111111",
    title: "title",
    description: "desc",
    slug: "slug",
    tags: [],
    author: { authorId: "a", name: "n", slug: "as" },
    sections: [],
    isEnrolled: false,
  };

  const result = handleGetCourseDetailResult({ success: true, course });

  expect(result).toEqual(course);
});

test("NOT_FOUND は notFound() を呼び出す", () => {
  expect(() =>
    handleGetCourseDetailResult({ success: false, code: "NOT_FOUND" }),
  ).toThrow("NEXT_NOT_FOUND");
});

test("INTERNAL_ERROR はエラーを投げる", () => {
  expect(() =>
    handleGetCourseDetailResult({ success: false, code: "INTERNAL_ERROR" }),
  ).toThrow("講座詳細の取得に失敗しました");
});
