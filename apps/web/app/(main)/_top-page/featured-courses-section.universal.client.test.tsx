import type { AnchorHTMLAttributes, PropsWithChildren } from "react";
import { expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";
import { FeaturedCoursesSection } from "./featured-courses-section.universal";

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

test("セクション見出しが表示される", async () => {
  const screen = await render(<FeaturedCoursesSection />);

  await expect
    .element(screen.getByRole("heading", { name: "注目の講座" }))
    .toBeInTheDocument();
});

test("公開済み講座のカードがリンクとして表示される", async () => {
  const screen = await render(<FeaturedCoursesSection />);

  await expect
    .element(screen.getByRole("link", { name: /Next\.js App Router 入門/ }))
    .toBeInTheDocument();
});

test("Coming Soon バッジが 3 件表示される", async () => {
  const screen = await render(<FeaturedCoursesSection />);

  expect(screen.getByText("Coming Soon").elements()).toHaveLength(3);
});
