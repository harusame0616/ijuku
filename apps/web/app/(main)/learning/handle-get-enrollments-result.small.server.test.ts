import { expect, test, vi } from "vitest";
import { handleGetEnrollmentsResult } from "./handle-get-enrollments-result.server";

class RedirectError extends Error {
  constructor(public to: string) {
    super(`REDIRECT:${to}`);
  }
}

vi.mock("next/navigation", () => ({
  redirect: (to: string) => {
    throw new RedirectError(to);
  },
}));

test("成功時は enrollments を返す", () => {
  const enrollments = handleGetEnrollmentsResult({
    success: true,
    enrollments: [
      {
        courseId: "11111111-1111-1111-1111-111111111111",
        title: "コース A",
      },
    ],
  });

  expect(enrollments).toEqual([
    {
      courseId: "11111111-1111-1111-1111-111111111111",
      title: "コース A",
    },
  ]);
});

test("UNAUTHORIZED の場合は /login へ redirect する", () => {
  expect(() =>
    handleGetEnrollmentsResult({ success: false, code: "UNAUTHORIZED" }),
  ).toThrow(new RedirectError("/login"));
});

test("INTERNAL_ERROR の場合はエラーを投げる", () => {
  expect(() =>
    handleGetEnrollmentsResult({ success: false, code: "INTERNAL_ERROR" }),
  ).toThrow("受講中コース一覧の取得に失敗しました");
});
