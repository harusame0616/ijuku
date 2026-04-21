import { createClient } from "@/lib/supabase/server";

export type User = {
  id: string;
  nickname: string;
};

export async function getUser(): Promise<User | null> {
  const supabase = await createClient();
  const {
    data: { user },
  } = await supabase.auth.getUser();

  if (!user) return null;

  return {
    id: user.id,
    nickname: user.user_metadata?.nickname ?? "",
  };
}
