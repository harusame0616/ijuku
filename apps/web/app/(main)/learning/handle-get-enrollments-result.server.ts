import { redirect } from "next/navigation";
import type { Enrollment, GetEnrollmentsResult } from "./enrollments.data";

export function handleGetEnrollmentsResult(
  result: GetEnrollmentsResult,
): Enrollment[] {
  if (result.success) return result.enrollments;

  switch (result.code) {
    case "UNAUTHORIZED": {
      redirect("/login");
      break;
    }
    case "INTERNAL_ERROR": {
      throw new Error("受講中コース一覧の取得に失敗しました");
    }
    default: {
      const exhaustiveCheck: never = result.code;
      throw new Error(`未対応のエラーコード: ${exhaustiveCheck}`);
    }
  }
}
