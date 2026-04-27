import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";
import { playwright } from "@vitest/browser-playwright";

export default defineConfig({
  plugins: [react()],
  resolve: {
    tsconfigPaths: true,
  },
  optimizeDeps: {
    include: [
      "@base-ui/react/button",
      "@base-ui/react/input",
      "@base-ui/react/merge-props",
      "@base-ui/react/separator",
      "@base-ui/react/use-render",
      "@supabase/ssr",
      "@tanstack/react-form",
      "class-variance-authority",
      "clsx",
      "next/headers",
      "next/link",
      "next/navigation",
      "tailwind-merge",
      "valibot",
    ],
  },
  test: {
    name: "client",
    include: [
      "**/*.client.test.ts",
      "**/*.client.test.tsx",
    ],
    browser: {
      enabled: true,
      provider: playwright(),
      headless: true,
      instances: [{ browser: "chromium" }],
    },
  },
});
