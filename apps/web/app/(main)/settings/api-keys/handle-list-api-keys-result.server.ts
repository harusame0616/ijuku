import { redirect } from "next/navigation";
import type { ApiKey, ListApiKeysResult } from "./api-keys.data";

export function handleListApiKeysResult(result: ListApiKeysResult): ApiKey[] {
  if (result.success) return result.apiKeys;

  switch (result.code) {
    case "UNAUTHORIZED": {
      redirect("/login");
      break;
    }
    case "INTERNAL_ERROR": {
      throw new Error("API キー一覧の取得に失敗しました");
    }
    default: {
      const exhaustiveCheck: never = result.code;
      throw new Error(`未対応のエラーコード: ${exhaustiveCheck}`);
    }
  }
}
