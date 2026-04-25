import { expect, test, vi } from "vitest";
import { handleGetProfileResult } from "./handle-get-profile-result.server";

class RedirectError extends Error {
  constructor(public to: string) {
    super(`REDIRECT:${to}`);
  }
}

class NotFoundError extends Error {
  constructor() {
    super("NOT_FOUND");
  }
}

vi.mock("next/navigation", () => ({
  redirect: (to: string) => {
    throw new RedirectError(to);
  },
  notFound: () => {
    throw new NotFoundError();
  },
}));

test("成功時は Profile を返す", () => {
  const profile = handleGetProfileResult({
    success: true,
    profile: { nickname: "taro", introduce: "hi" },
  });

  expect(profile).toEqual({ nickname: "taro", introduce: "hi" });
});

test("UNAUTHORIZED の場合は /login へ redirect する", () => {
  expect(() =>
    handleGetProfileResult({ success: false, code: "UNAUTHORIZED" }),
  ).toThrow(new RedirectError("/login"));
});

test("NOT_FOUND の場合は notFound を呼ぶ", () => {
  expect(() =>
    handleGetProfileResult({ success: false, code: "NOT_FOUND" }),
  ).toThrow(NotFoundError);
});

test("INTERNAL_ERROR の場合はエラーを投げる", () => {
  expect(() =>
    handleGetProfileResult({ success: false, code: "INTERNAL_ERROR" }),
  ).toThrow("プロフィールの取得に失敗しました");
});
