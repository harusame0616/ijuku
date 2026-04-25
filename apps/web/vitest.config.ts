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
        "lib/utils.ts",
        "hooks/use-is-hydrated.ts",
        "app/(main)/settings/profile-edit.presenter.form.client.tsx",
        "app/(main)/settings/profile.data.ts",
        "app/(main)/settings/handle-get-profile-result.server.ts",
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
