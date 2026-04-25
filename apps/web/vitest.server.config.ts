import { defineConfig } from "vitest/config";

export default defineConfig({
  resolve: {
    tsconfigPaths: true,
  },
  test: {
    name: "server",
    environment: "node",
    env: {
      NEXT_PUBLIC_SUPABASE_URL: "http://127.0.0.1:54321",
      NEXT_PUBLIC_SUPABASE_ANON_KEY:
        "sb_publishable_ACJWlzQHlZjBrEguHvfOxg_3BJgxAaH",
      API_URL: "http://127.0.0.1:8080",
    },
    include: [
      "**/*.small.server.test.ts",
      "**/*.medium.server.test.ts",
    ],
    pool: "forks",
  },
});
