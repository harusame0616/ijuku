import { expect, test } from "vitest";
import { render } from "vitest-browser-react";
import { ApiKeysListPresenter } from "./api-keys-list.presenter.client";

test("0 件の場合は空状態メッセージを表示する", async () => {
  const screen = await render(<ApiKeysListPresenter apiKeys={[]} />);

  await expect
    .element(screen.getByText("API キーがまだ登録されていません。"))
    .toBeInTheDocument();
});

test("複数件の API キーが jukubox_••••<末尾> 形式で表示される", async () => {
  const screen = await render(
    <ApiKeysListPresenter
      apiKeys={[
        {
          apiKeyId: "id-1",
          suffix: "a3f9",
          createdAt: "2026-01-10T00:00:00Z",
          expiredAt: "2027-01-10T00:00:00Z",
        },
        {
          apiKeyId: "id-2",
          suffix: "b218",
          createdAt: "2026-02-22T00:00:00Z",
          expiredAt: "2027-02-22T00:00:00Z",
        },
      ]}
    />,
  );

  await expect
    .element(screen.getByText("jukubox_••••a3f9"))
    .toBeInTheDocument();
  await expect
    .element(screen.getByText("jukubox_••••b218"))
    .toBeInTheDocument();
});

test("expiredAt が null の場合は『無期限』と表示される", async () => {
  const screen = await render(
    <ApiKeysListPresenter
      apiKeys={[
        {
          apiKeyId: "id-1",
          suffix: "c5d1",
          createdAt: "2026-04-01T00:00:00Z",
          expiredAt: null,
        },
      ]}
    />,
  );

  await expect.element(screen.getByText("無期限")).toBeInTheDocument();
});

test("disabled=true の場合はテーブルヘッダーが表示され、空状態メッセージは表示されない", async () => {
  const screen = await render(<ApiKeysListPresenter disabled />);

  await expect.element(screen.getByText("API キー")).toBeInTheDocument();
  await expect
    .element(screen.getByText("API キーがまだ登録されていません。"))
    .not.toBeInTheDocument();
});
