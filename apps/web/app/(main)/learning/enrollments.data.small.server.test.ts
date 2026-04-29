import { expect, test as base, vi } from "vitest";
import { getEnrollments } from "./enrollments.data";

const getSessionMock = vi.fn<() => Promise<{ data: { session: unknown } }>>();

vi.mock("@/lib/supabase/server", () => ({
  createClient: async () => ({
    auth: { getSession: getSessionMock },
  }),
}));

const test = base.extend<{ getEnrollments: typeof getEnrollments }>({
  getEnrollments: async ({}, provide) => {
    getSessionMock.mockReset();
    vi.unstubAllGlobals();
    await provide(getEnrollments);
  },
});

test("session が無い場合は UNAUTHORIZED コードを返す", async ({
  getEnrollments,
}) => {
  getSessionMock.mockResolvedValue({ data: { session: null } });

  const result = await getEnrollments();

  expect(result).toEqual({ success: false, code: "UNAUTHORIZED" });
});

test("API が 401 を返した場合は UNAUTHORIZED コードを返す", async ({
  getEnrollments,
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

  const result = await getEnrollments();

  expect(result).toEqual({ success: false, code: "UNAUTHORIZED" });
});

test("API が 500 を返した場合は INTERNAL_ERROR コードを返す", async ({
  getEnrollments,
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

  const result = await getEnrollments();

  expect(result).toEqual({ success: false, code: "INTERNAL_ERROR" });
});

test("API が 200 を返した場合は enrollments 配列を返す", async ({
  getEnrollments,
}) => {
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
          enrollments: [
            {
              courseId: "11111111-1111-1111-1111-111111111111",
              title: "コース A",
            },
            {
              courseId: "22222222-2222-2222-2222-222222222222",
              title: "コース B",
            },
          ],
        },
        { status: 200, headers: { "Content-Type": "application/json" } },
      ),
    ),
  );

  const result = await getEnrollments();

  expect(result).toEqual({
    success: true,
    enrollments: [
      {
        courseId: "11111111-1111-1111-1111-111111111111",
        title: "コース A",
      },
      {
        courseId: "22222222-2222-2222-2222-222222222222",
        title: "コース B",
      },
    ],
  });
});
