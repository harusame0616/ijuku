import { expect, test } from "vitest";
import { cn } from "./utils";

test("cn: 複数のクラス名を結合する", () => {
  expect(cn("a", "b")).toBe("a b");
});

test("cn: truthy な値のみ結合する", () => {
  expect(cn("a", false && "b", null, undefined, "c")).toBe("a c");
});

test("cn: tailwind の競合クラスを後勝ちでマージする", () => {
  expect(cn("px-2", "px-4")).toBe("px-4");
});
