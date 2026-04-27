import { expect, test as base, vi } from "vitest";
import { listApiKeys } from "./api-keys.data";

const getSessionMock = vi.fn<() => Promise<{ data: { session: unknown } }>>();

vi.mock("@/lib/supabase/server", () => ({
  createClient: async () => ({
    auth: { getSession: getSessionMock },
  }),
}));

const test = base.extend<{ listApiKeys: typeof listApiKeys }>({
  listApiKeys: async ({}, provide) => {
    getSessionMock.mockReset();
    vi.unstubAllGlobals();
    await provide(listApiKeys);
  },
});

test("session が無い場合は UNAUTHORIZED コードを返す", async ({
  listApiKeys,
}) => {
  getSessionMock.mockResolvedValue({ data: { session: null } });

  const result = await listApiKeys();

  expect(result).toEqual({ success: false, code: "UNAUTHORIZED" });
});

test("API が 401 を返した場合は UNAUTHORIZED コードを返す", async ({
  listApiKeys,
}) => {
  getSessionMock.mockResolvedValue({
    data: {
      session: { user: { id: "user-1" }, access_token: "token-1" },
    },
  });
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => new Response(null, { status: 401 })),
  );

  const result = await listApiKeys();

  expect(result).toEqual({ success: false, code: "UNAUTHORIZED" });
});

test("API が 500 を返した場合は INTERNAL_ERROR コードを返す", async ({
  listApiKeys,
}) => {
  getSessionMock.mockResolvedValue({
    data: {
      session: { user: { id: "user-1" }, access_token: "token-1" },
    },
  });
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => new Response(null, { status: 500 })),
  );

  const result = await listApiKeys();

  expect(result).toEqual({ success: false, code: "INTERNAL_ERROR" });
});

test("API が 200 を返した場合は apiKeys 配列を返す", async ({ listApiKeys }) => {
  getSessionMock.mockResolvedValue({
    data: {
      session: { user: { id: "user-1" }, access_token: "token-1" },
    },
  });
  vi.stubGlobal(
    "fetch",
    vi.fn(async () =>
      Response.json(
        {
          apiKeys: [
            {
              apiKeyId: "id-1",
              suffix: "a3f9",
              createdAt: "2026-01-10T00:00:00Z",
              expiredAt: "2027-01-10T00:00:00Z",
            },
          ],
        },
        { status: 200, headers: { "Content-Type": "application/json" } },
      ),
    ),
  );

  const result = await listApiKeys();

  expect(result).toEqual({
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
});

test("API URL とユーザー ID と Authorization ヘッダーが正しく組み立てられる", async ({
  listApiKeys,
}) => {
  getSessionMock.mockResolvedValue({
    data: {
      session: { user: { id: "user-xyz" }, access_token: "token-xyz" },
    },
  });
  const fetchMock = vi.fn(async () =>
    Response.json({ apiKeys: [] }, { status: 200 }),
  );
  vi.stubGlobal("fetch", fetchMock);

  await listApiKeys();

  expect(fetchMock).toHaveBeenCalledWith(
    `${process.env.API_URL}/v1/users/user-xyz/settings/apikeys`,
    { headers: { Authorization: "Bearer token-xyz" } },
  );
});
