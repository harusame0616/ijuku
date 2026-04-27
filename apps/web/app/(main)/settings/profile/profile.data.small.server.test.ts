import { expect, test as base, vi } from "vitest";
import { getProfile } from "./profile.data";

const getSessionMock = vi.fn<() => Promise<{ data: { session: unknown } }>>();

vi.mock("@/lib/supabase/server", () => ({
  createClient: async () => ({
    auth: { getSession: getSessionMock },
  }),
}));

const test = base.extend<{ getProfile: typeof getProfile }>({
  getProfile: async ({}, provide) => {
    getSessionMock.mockReset();
    vi.unstubAllGlobals();
    await provide(getProfile);
  },
});

test("session が無い場合は UNAUTHORIZED コードを返す", async ({ getProfile }) => {
  getSessionMock.mockResolvedValue({ data: { session: null } });

  const result = await getProfile();

  expect(result).toEqual({ success: false, code: "UNAUTHORIZED" });
});

test("API が 404 を返した場合は NOT_FOUND コードを返す", async ({ getProfile }) => {
  getSessionMock.mockResolvedValue({
    data: {
      session: { user: { id: "user-1" }, access_token: "token-1" },
    },
  });
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => new Response(null, { status: 404 })),
  );

  const result = await getProfile();

  expect(result).toEqual({ success: false, code: "NOT_FOUND" });
});

test("API が 500 を返した場合は INTERNAL_ERROR コードを返す", async ({ getProfile }) => {
  getSessionMock.mockResolvedValue({
    data: {
      session: { user: { id: "user-1" }, access_token: "token-1" },
    },
  });
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => new Response(null, { status: 500 })),
  );

  const result = await getProfile();

  expect(result).toEqual({ success: false, code: "INTERNAL_ERROR" });
});

test("API が 200 を返した場合は Profile を返す", async ({ getProfile }) => {
  getSessionMock.mockResolvedValue({
    data: {
      session: { user: { id: "user-1" }, access_token: "token-1" },
    },
  });
  vi.stubGlobal(
    "fetch",
    vi.fn(async () =>
      Response.json(
        { nickname: "taro", introduce: "hi" },
        { status: 200, headers: { "Content-Type": "application/json" } },
      ),
    ),
  );

  const result = await getProfile();

  expect(result).toEqual({
    success: true,
    profile: { nickname: "taro", introduce: "hi" },
  });
});
