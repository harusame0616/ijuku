import { createClient } from "@/lib/supabase/server";

export interface Enrollment {
  courseId: string;
  title: string;
}

export type GetEnrollmentsErrorCode = "UNAUTHORIZED" | "INTERNAL_ERROR";

export type GetEnrollmentsResult =
  | { success: true; enrollments: Enrollment[] }
  | { success: false; code: GetEnrollmentsErrorCode };

export async function getEnrollments(): Promise<GetEnrollmentsResult> {
  const supabase = await createClient();
  const {
    data: { session },
  } = await supabase.auth.getSession();

  if (!session) return { success: false, code: "UNAUTHORIZED" };

  const response = await fetch(
    `${process.env.API_URL}/v1/users/${session.user.id}/enrollments`,
    { headers: { Authorization: `Bearer ${session.access_token}` } },
  );

  if (response.status === 401) return { success: false, code: "UNAUTHORIZED" };
  if (!response.ok) return { success: false, code: "INTERNAL_ERROR" };

  const body = (await response.json()) as { enrollments: Enrollment[] };
  return { success: true, enrollments: body.enrollments };
}
