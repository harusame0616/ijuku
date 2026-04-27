import { createClient } from "@/lib/supabase/server";

export interface Profile {
  nickname: string;
  introduce: string;
}

export type GetProfileErrorCode =
  | "UNAUTHORIZED"
  | "NOT_FOUND"
  | "INTERNAL_ERROR";

export type GetProfileResult =
  | { success: true; profile: Profile }
  | { success: false; code: GetProfileErrorCode };

export async function getProfile(): Promise<GetProfileResult> {
  const supabase = await createClient();
  const {
    data: { session },
  } = await supabase.auth.getSession();

  if (!session) return { success: false, code: "UNAUTHORIZED" };

  const response = await fetch(
    `${process.env.API_URL}/v1/users/${session.user.id}`,
    { headers: { Authorization: `Bearer ${session.access_token}` } },
  );

  if (response.status === 404) return { success: false, code: "NOT_FOUND" };
  if (!response.ok) return { success: false, code: "INTERNAL_ERROR" };

  const profile = (await response.json()) as Profile;
  return { success: true, profile };
}
