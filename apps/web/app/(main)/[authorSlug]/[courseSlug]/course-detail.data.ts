import { createClient } from "@/lib/supabase/server";

export interface CourseDetailTopic {
  topicId: string;
  title: string;
  description: string;
}

export interface CourseDetailSection {
  sectionId: string;
  title: string;
  description: string;
  topics: CourseDetailTopic[];
}

export interface CourseDetailAuthor {
  authorId: string;
  name: string;
  slug: string;
}

export interface CourseDetail {
  courseId: string;
  title: string;
  description: string;
  slug: string;
  tags: string[];
  author: CourseDetailAuthor;
  sections: CourseDetailSection[];
  isEnrolled: boolean;
}

export type GetCourseDetailErrorCode = "NOT_FOUND" | "INTERNAL_ERROR";

export type GetCourseDetailResult =
  | { success: true; course: CourseDetail }
  | { success: false; code: GetCourseDetailErrorCode };

export async function getCourseDetail(
  authorSlug: string,
  courseSlug: string,
): Promise<GetCourseDetailResult> {
  const supabase = await createClient();
  const {
    data: { session },
  } = await supabase.auth.getSession();

  const headers: Record<string, string> = {};
  if (session) {
    headers.Authorization = `Bearer ${session.access_token}`;
  }

  const response = await fetch(
    `${process.env.API_URL}/v1/courses/${encodeURIComponent(authorSlug)}/${encodeURIComponent(courseSlug)}`,
    { headers },
  );

  if (response.status === 404) return { success: false, code: "NOT_FOUND" };
  if (!response.ok) return { success: false, code: "INTERNAL_ERROR" };

  const body = (await response.json()) as CourseDetail;
  return { success: true, course: body };
}
