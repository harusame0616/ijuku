import { createClient } from "@/lib/supabase/server";

export interface ApiKey {
  apiKeyId: string;
  suffix: string;
  createdAt: string;
  expiredAt: string | null;
}

export type ListApiKeysErrorCode = "UNAUTHORIZED" | "INTERNAL_ERROR";

export type ListApiKeysResult =
  | { success: true; apiKeys: ApiKey[] }
  | { success: false; code: ListApiKeysErrorCode };

export async function listApiKeys(): Promise<ListApiKeysResult> {
  const supabase = await createClient();
  const {
    data: { session },
  } = await supabase.auth.getSession();

  if (!session) return { success: false, code: "UNAUTHORIZED" };

  const response = await fetch(
    `${process.env.API_URL}/v1/users/${session.user.id}/settings/apikeys`,
    { headers: { Authorization: `Bearer ${session.access_token}` } },
  );

  if (response.status === 401) return { success: false, code: "UNAUTHORIZED" };
  if (!response.ok) return { success: false, code: "INTERNAL_ERROR" };

  const body = (await response.json()) as { apiKeys: ApiKey[] };
  console.log(session)
  return { success: true, apiKeys: body.apiKeys };
}
