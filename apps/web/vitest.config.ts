import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    projects: [
      "./vitest.server.config.ts",
      "./vitest.client.config.ts",
    ],
    coverage: {
      provider: "v8",
      reporter: ["text", "html", "json-summary"],
      // テスト対象ファイルを明示的に列挙する
      // 新規テストを追加したらこのリストにも対象ファイルを追加する
      include: [
        "lib/auth/config.ts",
        "lib/auth/index.ts",
        "lib/auth/supabase-auth.ts",
        "lib/utilities.ts",
        "hooks/use-is-hydrated.ts",
        "app/(main)/settings/_components/settings-nav.client.tsx",
        "app/(main)/settings/profile/profile-edit.presenter.form.client.tsx",
        "app/(main)/settings/profile/profile.data.ts",
        "app/(main)/settings/profile/handle-get-profile-result.server.ts",
        "app/(main)/settings/api-keys/api-keys-list.presenter.client.tsx",
        "app/(main)/settings/api-keys/api-keys.data.ts",
        "app/(main)/settings/api-keys/handle-list-api-keys-result.server.ts",
        "app/(main)/settings/api-keys/generate-api-key.action.ts",
        "app/(main)/settings/api-keys/generate-api-key.presenter.client.tsx",
        "app/(main)/_components/header-search.universal.tsx",
        "app/(main)/_top-page/featured-courses.data.ts",
        "app/(main)/_top-page/featured-course-card.universal.tsx",
        "app/(main)/_top-page/featured-courses-section.universal.tsx",
        "app/(main)/_top-page/lp-link-footer.universal.tsx",
        "app/(main)/_top-page/continue-learning.presenter.universal.tsx",
        "app/(main)/_top-page/side-b-hero.universal.tsx",
      ],
      thresholds: {
        lines: 80,
        statements: 80,
        functions: 80,
        branches: 80,
      },
    },
  },
});
