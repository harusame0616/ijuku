import type { JSX } from "react";
import { expect, test } from "vitest";
import { renderToString } from "react-dom/server";
import { renderHook } from "vitest-browser-react";
import { useIsHydrated } from "./use-is-hydrated";

function HydrationProbe(): JSX.Element {
  const hydrated = useIsHydrated();
  return <span>{hydrated ? "yes" : "no"}</span>;
}

test("useIsHydrated: サーバーレンダリング（マウント前）では false を返す", () => {
  const html = renderToString(<HydrationProbe />);
  expect(html).toContain("no");
  expect(html).not.toContain("yes");
});

test("useIsHydrated: マウント後（ハイドレーション完了後）は true を返す", async () => {
  const { result } = await renderHook(() => useIsHydrated());
  expect(result.current).toBe(true);
});
