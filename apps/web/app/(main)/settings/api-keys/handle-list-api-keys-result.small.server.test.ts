import { expect, test, vi } from "vitest";
import { handleListApiKeysResult } from "./handle-list-api-keys-result.server";

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

test("成功時は apiKeys を返す", () => {
  const apiKeys = handleListApiKeysResult({
    success: true,
    apiKeys: [
      {
        apiKeyId: "id-1",
        suffix: "a3f9",
        createdAt: "2026-01-10T00:00:00Z",
        expiredAt: "2027-01-10T00:00:00Z",
      },
    ],
  });

  expect(apiKeys).toEqual([
    {
      apiKeyId: "id-1",
      suffix: "a3f9",
      createdAt: "2026-01-10T00:00:00Z",
      expiredAt: "2027-01-10T00:00:00Z",
    },
  ]);
});

test("UNAUTHORIZED の場合は /login へ redirect する", () => {
  expect(() =>
    handleListApiKeysResult({ success: false, code: "UNAUTHORIZED" }),
  ).toThrow(new RedirectError("/login"));
});

test("INTERNAL_ERROR の場合はエラーを投げる", () => {
  expect(() =>
    handleListApiKeysResult({ success: false, code: "INTERNAL_ERROR" }),
  ).toThrow("API キー一覧の取得に失敗しました");
});
