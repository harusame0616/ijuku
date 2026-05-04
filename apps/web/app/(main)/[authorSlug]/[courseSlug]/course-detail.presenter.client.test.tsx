import type { AnchorHTMLAttributes, PropsWithChildren } from "react";
import { expect, test, vi } from "vitest";
import { render } from "vitest-browser-react";
import type { CourseDetail } from "./course-detail.data";
import { CourseDetailPresenter } from "./course-detail.presenter.universal";

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

const baseCourse: CourseDetail = {
  courseId: "11111111-1111-1111-1111-111111111111",
  title: "テスト講座",
  description: "テスト講座の説明",
  slug: "test-course",
  tags: ["go", "backend"],
  author: {
    authorId: "22222222-2222-2222-2222-222222222222",
    name: "テスト著者",
    slug: "test-author",
  },
  sections: [
    {
      sectionId: "33333333-3333-3333-3333-333333333333",
      title: "セクション1",
      description: "セクション1の説明",
      topics: [
        {
          topicId: "44444444-4444-4444-4444-444444444444",
          title: "トピック1",
          description: "トピック1の説明",
        },
      ],
    },
  ],
  isEnrolled: false,
};

test("講座タイトル・著者・セクション・トピックが表示される", async () => {
  const screen = await render(<CourseDetailPresenter course={baseCourse} />);

  await expect
    .element(screen.getByRole("heading", { name: "テスト講座", level: 1 }))
    .toBeInTheDocument();
  await expect
    .element(screen.getByText("講師: テスト著者"))
    .toBeInTheDocument();
  await expect
    .element(screen.getByRole("heading", { name: "1. セクション1", level: 3 }))
    .toBeInTheDocument();
  await expect
    .element(screen.getByText("1-1. トピック1", { exact: true }))
    .toBeInTheDocument();
});

test("受講中ならバッジを表示する", async () => {
  const screen = await render(
    <CourseDetailPresenter course={{ ...baseCourse, isEnrolled: true }} />,
  );

  await expect.element(screen.getByText("受講中")).toBeInTheDocument();
});

test("未受講ならバッジを表示しない", async () => {
  const screen = await render(<CourseDetailPresenter course={baseCourse} />);

  await expect.element(screen.getByText("受講中")).not.toBeInTheDocument();
});

test("受講コマンドが authorSlug/courseSlug 形式で表示される", async () => {
  const screen = await render(<CourseDetailPresenter course={baseCourse} />);

  await expect
    .element(screen.getByRole("textbox"))
    .toHaveValue("/jukubox enroll test-author/test-course");
});

test("セクションが空の場合は空状態メッセージを表示する", async () => {
  const screen = await render(
    <CourseDetailPresenter course={{ ...baseCourse, sections: [] }} />,
  );

  await expect
    .element(screen.getByText("この講座にはまだセクションがありません"))
    .toBeInTheDocument();
});

test("タグが表示される", async () => {
  const screen = await render(<CourseDetailPresenter course={baseCourse} />);

  await expect.element(screen.getByText("go")).toBeInTheDocument();
  await expect.element(screen.getByText("backend")).toBeInTheDocument();
});
