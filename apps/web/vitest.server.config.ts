import { existsSync } from "node:fs";
import { loadEnvFile } from "node:process";
import { defineConfig } from "vitest/config";

if (existsSync(".env")) {
  loadEnvFile(".env");
}

export default defineConfig({
  resolve: {
    tsconfigPaths: true,
  },
  test: {
    name: "server",
    environment: "node",
    include: [
      "**/*.small.server.test.ts",
      "**/*.medium.server.test.ts",
    ],
    // lib/test/** はテスト用ユーティリティのため、テストファイルは置かない
    exclude: ["**/node_modules/**", "lib/test/**"],
    pool: "forks",
  },
});
