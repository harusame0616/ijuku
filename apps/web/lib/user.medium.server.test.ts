import { expect, test as base, vi } from "vitest";

vi.mock("next/headers", () => ({
  cookies: async () => ({
    getAll: () => [] as { name: string; value: string }[],
    set: () => {},
  }),
}));

const test = base.extend<{ getUser: typeof import("./user").getUser }>({
  getUser: async ({}, provide) => {
    process.env.NEXT_PUBLIC_SUPABASE_URL ??= "http://127.0.0.1:54321";
    process.env.NEXT_PUBLIC_SUPABASE_ANON_KEY ??=
      "sb_publishable_ACJWlzQHlZjBrEguHvfOxg_3BJgxAaH";
    const { getUser } = await import("./user");
    await provide(getUser);
  },
});

test("getUser: 未ログイン状態では null を返す（実 DB 疎通）", async ({ getUser }) => {
  const user = await getUser();
  expect(user).toBeNull();
});
