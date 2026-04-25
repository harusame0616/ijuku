import { notFound, redirect } from "next/navigation";
import type { GetProfileResult, Profile } from "./profile.data";

export function handleGetProfileResult(result: GetProfileResult): Profile {
  if (result.success) return result.profile;

  switch (result.code) {
    case "UNAUTHORIZED": {
      redirect("/login");
      break;
    }
    case "NOT_FOUND": {
      notFound();
      break;
    }
    case "INTERNAL_ERROR": {
      throw new Error("プロフィールの取得に失敗しました");
    }
    default: {
      const exhaustiveCheck: never = result.code;
      throw new Error(`未対応のエラーコード: ${exhaustiveCheck}`);
    }
  }
}
