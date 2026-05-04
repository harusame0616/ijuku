import type { AnchorHTMLAttributes, PropsWithChildren } from "react";
import { expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";
import { FeaturedCourseCard } from "./featured-course-card.universal";

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

test("公開済みカードはタイトル・説明・タグを表示する", async () => {
  const screen = await render(
    <FeaturedCourseCard
      course={{
        status: "published",
        authorSlug: "jukubox",
        courseSlug: "nextjs-app-router-getting-started",
        title: "Next.js App Router 入門",
        description: "ミニブログを作りながら学ぶ Get Started コース。",
        tags: ["nextjs", "react"],
      }}
    />,
  );

  await expect
    .element(screen.getByText("Next.js App Router 入門"))
    .toBeInTheDocument();
  await expect
    .element(screen.getByText("ミニブログを作りながら学ぶ Get Started コース。"))
    .toBeInTheDocument();
  await expect.element(screen.getByText("nextjs")).toBeInTheDocument();
  await expect.element(screen.getByText("react")).toBeInTheDocument();
});

test("公開済みカードは講座詳細ページへのリンクとして機能する", async () => {
  const screen = await render(
    <FeaturedCourseCard
      course={{
        status: "published",
        authorSlug: "jukubox",
        courseSlug: "nextjs-app-router-getting-started",
        title: "Next.js App Router 入門",
        description: "説明",
        tags: [],
      }}
    />,
  );

  await expect
    .element(screen.getByRole("link", { name: /Next\.js App Router 入門/ }))
    .toBeInTheDocument();
});

test("Coming Soon カードはタイトル・説明・Coming Soon バッジを表示する", async () => {
  const screen = await render(
    <FeaturedCourseCard
      course={{
        status: "coming-soon",
        title: "React 19 入門",
        description: "新しい React の入門コース。",
        tags: ["react"],
      }}
    />,
  );

  await expect.element(screen.getByText("React 19 入門")).toBeInTheDocument();
  await expect
    .element(screen.getByText("新しい React の入門コース。"))
    .toBeInTheDocument();
  await expect.element(screen.getByText("Coming Soon")).toBeInTheDocument();
});

test("Coming Soon カードはリンクを持たない", async () => {
  const screen = await render(
    <FeaturedCourseCard
      course={{
        status: "coming-soon",
        title: "React 19 入門",
        description: "説明",
        tags: [],
      }}
    />,
  );

  expect(screen.getByRole("link").elements()).toHaveLength(0);
});
