import { expect, test } from "vitest";
import { featuredCourses } from "./featured-courses.data";

test("注目講座は 4 件で、公開済みが 1 件含まれる", () => {
  expect(featuredCourses).toHaveLength(4);
  const published = featuredCourses.filter((c) => c.status === "published");
  expect(published).toHaveLength(1);
});

test("公開済み講座は jukubox/nextjs-app-router-getting-started を指す", () => {
  const published = featuredCourses.find((c) => c.status === "published");
  expect(published).toBeDefined();
  if (published?.status !== "published") {
    throw new Error("published course is missing");
  }

  expect(published.authorSlug).toBe("jukubox");
  expect(published.courseSlug).toBe("nextjs-app-router-getting-started");
});
