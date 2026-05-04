import type { AnchorHTMLAttributes, PropsWithChildren } from "react";
import { expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";
import { LpLinkFooter } from "./lp-link-footer.universal";

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

test("サービス紹介ページへの導線リンクが表示される", async () => {
  const screen = await render(<LpLinkFooter />);

  await expect
    .element(screen.getByRole("link", { name: "サービス紹介ページへ" }))
    .toBeInTheDocument();
});
