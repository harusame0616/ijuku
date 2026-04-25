import { expect, test, vi } from "vitest";
import { getUser } from "./user";

vi.mock("next/headers", () => ({
  cookies: async () => ({
    getAll: () => [] as { name: string; value: string }[],
    set: () => {},
  }),
}));

test("getUser: 未ログイン状態では null を返す（実 DB 疎通）", async () => {
  const user = await getUser();
  expect(user).toBeNull();
});
