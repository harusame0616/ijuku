import type { AnchorHTMLAttributes, PropsWithChildren } from "react";
import { expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";
import { EnrollmentsListPresenter } from "./enrollments-list.presenter.client";

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

test("0 件の場合は空状態メッセージを表示する", async () => {
  const screen = await render(<EnrollmentsListPresenter enrollments={[]} />);

  await expect
    .element(screen.getByText("現在受講中の講座はありません"))
    .toBeInTheDocument();
});

test("複数件の受講中コースがタイトル付きで表示される", async () => {
  const screen = await render(
    <EnrollmentsListPresenter
      enrollments={[
        {
          courseId: "11111111-1111-1111-1111-111111111111",
          title: "コース A",
        },
        {
          courseId: "22222222-2222-2222-2222-222222222222",
          title: "コース B",
        },
      ]}
    />,
  );

  await expect.element(screen.getByText("コース A")).toBeInTheDocument();
  await expect.element(screen.getByText("コース B")).toBeInTheDocument();
});

