"use server";

import { createClient } from "@/lib/supabase/server";

export type EditProfileErrorCode = "UNAUTHORIZED" | "UPDATE_FAILED";

export type UpdateProfileResult =
  | { success: true }
  | { success: false; code: EditProfileErrorCode };

export async function editProfile(
  nickname: string,
  introduce: string,
): Promise<UpdateProfileResult> {
  const supabase = await createClient();
  const {
    data: { session },
  } = await supabase.auth.getSession();

  if (!session) {
    return { success: false, code: "UNAUTHORIZED" };
  }

  const res = await fetch(
    `${process.env.API_URL}/v1/users/${session.user.id}`,
    {
      method: "PATCH",
      headers: {
        Authorization: `Bearer ${session.access_token}`,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ nickname, introduce }),
    },
  );

  if (!res.ok) {
    return { success: false, code: "UPDATE_FAILED" };
  }

  return { success: true };
}
