import { expect, test as base, vi } from "vitest";
import { submitContact } from "@/app/(lp)/contacts/contact.action";

const getHeadersMock = vi.fn<() => Promise<Headers>>();

vi.mock("next/headers", () => ({
  headers: () => getHeadersMock(),
}));

const test = base.extend<{ submitContact: typeof submitContact }>({
  submitContact: async ({}, provide) => {
    getHeadersMock.mockReset();
    vi.unstubAllGlobals();
    vi.unstubAllEnvs();
    vi.stubEnv("API_URL", "http://api.test");
    await provide(submitContact);
  },
});

const baseInput = {
  name: "山田",
  email: "taro@example.com",
  phone: "",
  content: "本文",
};

test("API が 201 を返したとき success=true を返す", async ({ submitContact }) => {
  getHeadersMock.mockResolvedValue(new Headers());
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => new Response(null, { status: 201 })),
  );

  const result = await submitContact(baseInput);

  expect(result).toEqual({ success: true });
});

test("API が 500 を返したとき SUBMIT_FAILED を返す", async ({ submitContact }) => {
  getHeadersMock.mockResolvedValue(new Headers());
  vi.stubGlobal(
    "fetch",
    vi.fn(async () => new Response(null, { status: 500 })),
  );

  const result = await submitContact(baseInput);

  expect(result).toEqual({ success: false, code: "SUBMIT_FAILED" });
});

test("phone が空文字のときはリクエストボディの phone が null になる", async ({
  submitContact,
}) => {
  getHeadersMock.mockResolvedValue(new Headers());
  const fetchMock = vi.fn<typeof fetch>(
    async () => new Response(null, { status: 201 }),
  );
  vi.stubGlobal("fetch", fetchMock);

  await submitContact({ ...baseInput, phone: "" });

  const init = fetchMock.mock.calls[0]?.[1] as RequestInit | undefined;
  const body = JSON.parse((init?.body as string) ?? "{}") as {
    phone: string | null;
  };
  expect(body.phone).toBeNull();
});

test("phone を指定したときはリクエストボディに反映される", async ({
  submitContact,
}) => {
  getHeadersMock.mockResolvedValue(new Headers());
  const fetchMock = vi.fn<typeof fetch>(
    async () => new Response(null, { status: 201 }),
  );
  vi.stubGlobal("fetch", fetchMock);

  await submitContact({ ...baseInput, phone: "  03-1234-5678  " });

  const init = fetchMock.mock.calls[0]?.[1] as RequestInit | undefined;
  const body = JSON.parse((init?.body as string) ?? "{}") as {
    phone: string | null;
  };
  expect(body.phone).toBe("03-1234-5678");
});

test("X-Forwarded-For と User-Agent を API に転送する", async ({
  submitContact,
}) => {
  getHeadersMock.mockResolvedValue(
    new Headers({
      "x-forwarded-for": "203.0.113.10, 198.51.100.1",
      "user-agent": "TestAgent/1.0",
    }),
  );
  const fetchMock = vi.fn<typeof fetch>(
    async () => new Response(null, { status: 201 }),
  );
  vi.stubGlobal("fetch", fetchMock);

  await submitContact(baseInput);

  const init = fetchMock.mock.calls[0]?.[1] as RequestInit | undefined;
  const headers = new Headers(init?.headers);
  expect(headers.get("X-Forwarded-For")).toBe("203.0.113.10, 198.51.100.1");
  expect(headers.get("User-Agent")).toBe("TestAgent/1.0");
});
