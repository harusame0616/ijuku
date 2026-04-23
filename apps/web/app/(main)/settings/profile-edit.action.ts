"use server";

import { createClient } from "@/lib/supabase/server";

export type UpdateProfileResult =
  | { success: true }
  | { success: false; message: string };

export async function editProfile(
  nickname: string,
  introduce: string,
): Promise<UpdateProfileResult> {
  const supabase = await createClient();
  const {
    data: { session },
  } = await supabase.auth.getSession();

  if (!session) {
    return { success: false, message: "ログインが必要です" };
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
    return { success: false, message: "プロフィールの更新に失敗しました" };
  }

  return { success: true };
}
