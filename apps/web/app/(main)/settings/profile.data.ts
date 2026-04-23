import { createClient } from "@/lib/supabase/server";

export type Profile = {
  nickname: string;
  introduce: string;
};

export async function getProfile(): Promise<Profile | null> {
  const supabase = await createClient();
  const {
    data: { session },
  } = await supabase.auth.getSession();

  if (!session) return null;

  const res = await fetch(
    `${process.env.API_URL}/v1/users/${session.user.id}`,
    { headers: { Authorization: `Bearer ${session.access_token}` } },
  );

  if (!res.ok) return null;

  return res.json();
}
