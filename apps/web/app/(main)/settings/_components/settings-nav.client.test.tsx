import type { AnchorHTMLAttributes, PropsWithChildren } from "react";
import { expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";
import { SettingsNav } from "./settings-nav.client";

const usePathnameMock = vi.fn<() => string>();

vi.mock("next/navigation", () => ({
  default: {},
  usePathname: () => usePathnameMock(),
}));

vi.mock("next/link", () => {
  function Link({
    children,
    href,
    ...props
  }: PropsWithChildren<AnchorHTMLAttributes<HTMLAnchorElement>>) {
    return (
      <a href={href} {...props}>
        {children}
      </a>
    );
  }
  return { __esModule: true, default: Link };
});

test("プロフィールと API キーのリンクが表示される", async () => {
  usePathnameMock.mockReturnValue("/settings/profile");

  const screen = await render(<SettingsNav />);

  await expect
    .element(screen.getByRole("link", { name: "プロフィール" }))
    .toBeInTheDocument();
  await expect
    .element(screen.getByRole("link", { name: "API キー" }))
    .toBeInTheDocument();
});

test("/settings/profile のときプロフィール項目に aria-current=page が付く", async () => {
  usePathnameMock.mockReturnValue("/settings/profile");

  const screen = await render(<SettingsNav />);

  await expect
    .element(screen.getByRole("link", { name: "プロフィール" }))
    .toHaveAttribute("aria-current", "page");
  await expect
    .element(screen.getByRole("link", { name: "API キー" }))
    .not.toHaveAttribute("aria-current");
});

test("/settings/api-keys のとき API キー項目に aria-current=page が付く", async () => {
  usePathnameMock.mockReturnValue("/settings/api-keys");

  const screen = await render(<SettingsNav />);

  await expect
    .element(screen.getByRole("link", { name: "API キー" }))
    .toHaveAttribute("aria-current", "page");
  await expect
    .element(screen.getByRole("link", { name: "プロフィール" }))
    .not.toHaveAttribute("aria-current");
});
