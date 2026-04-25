import { notFound, redirect } from "next/navigation";
import type { GetProfileResult, Profile } from "./profile.data";

export function handleGetProfileResult(result: GetProfileResult): Profile {
  if (result.success) return result.profile;

  switch (result.code) {
    case "UNAUTHORIZED":
      redirect("/login");
    case "NOT_FOUND":
      notFound();
    case "INTERNAL_ERROR":
      throw new Error("プロフィールの取得に失敗しました");
    default: {
      const _exhaustive: never = result.code;
      throw new Error(`未対応のエラーコード: ${_exhaustive}`);
    }
  }
}
