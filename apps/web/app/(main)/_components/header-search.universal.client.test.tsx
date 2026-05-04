import type { AnchorHTMLAttributes, PropsWithChildren } from "react";
import { expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";
import { HeaderSearch } from "./header-search.universal";

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

test("検索ロール要素とクエリ入力欄が表示される", async () => {
  const screen = await render(<HeaderSearch />);

  await expect.element(screen.getByRole("search")).toBeInTheDocument();
  await expect
    .element(screen.getByLabelText("コースを検索"))
    .toBeInTheDocument();
});

test("SP 用の検索ボタンがコース一覧ページへのリンクとして表示される", async () => {
  const screen = await render(<HeaderSearch />);

  await expect
    .element(screen.getByRole("link", { name: "コースを検索" }))
    .toBeInTheDocument();
});
