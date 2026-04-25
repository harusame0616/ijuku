import { defineConfig } from "vitest/config";

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
    pool: "forks",
  },
});
