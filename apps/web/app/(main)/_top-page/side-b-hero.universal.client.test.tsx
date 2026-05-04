import type { AnchorHTMLAttributes, PropsWithChildren } from "react";
import { expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";
import { SideBHero } from "./side-b-hero.universal";

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

test("制作側のキャッチコピーが見出しとして表示される", async () => {
  const screen = await render(<SideBHero />);

  await expect
    .element(
      screen.getByRole("heading", {
        name: /自分の知識を\s*AI と一緒に教えよう。/,
      }),
    )
    .toBeInTheDocument();
});

test("「コースを作る」CTA が表示される", async () => {
  const screen = await render(<SideBHero />);

  await expect.element(screen.getByText("コースを作る")).toBeInTheDocument();
});
