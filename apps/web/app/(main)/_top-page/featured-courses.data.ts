export type FeaturedCourse =
  | {
      status: "published";
      authorSlug: string;
      courseSlug: string;
      title: string;
      description: string;
      tags: readonly string[];
    }
  | {
      status: "coming-soon";
      title: string;
      description: string;
      tags: readonly string[];
    };

export const featuredCourses: readonly FeaturedCourse[] = [
  {
    status: "published",
    authorSlug: "jukubox",
    courseSlug: "nextjs-app-router-getting-started",
    title: "Next.js App Router 入門",
    description:
      "ミニブログを作りながら、ファイルベースルーティング・Server / Client Components・Cache Components など Next.js 16 の主要機能を一通り学ぶ Get Started コース。",
    tags: ["nextjs", "react", "frontend"],
  },
  {
    status: "coming-soon",
    title: "React 19 入門",
    description:
      "React Compiler や Actions、新しいフック群を踏まえた、いまから始める最新の React 入門。",
    tags: ["react", "frontend"],
  },
  {
    status: "coming-soon",
    title: "Supabase 入門",
    description:
      "Auth / Database / Storage / Edge Functions を一気通貫で扱う、はじめての Supabase。",
    tags: ["supabase", "backend"],
  },
  {
    status: "coming-soon",
    title: "PostgreSQL 基礎学習",
    description:
      "SQL の基本から JOIN・トランザクション・インデックスまで、実務で使える PostgreSQL の基礎。",
    tags: ["postgresql", "database"],
  },
];
