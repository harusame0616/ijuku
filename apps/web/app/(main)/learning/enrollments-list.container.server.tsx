import type { JSX } from "react";
import { EnrollmentsListPresenter } from "./enrollments-list.presenter.client";
import { getEnrollments } from "./enrollments.data";
import { handleGetEnrollmentsResult } from "./handle-get-enrollments-result.server";

export async function EnrollmentsListContainer(): Promise<JSX.Element> {
  const enrollments = handleGetEnrollmentsResult(await getEnrollments());

  return <EnrollmentsListPresenter enrollments={enrollments} />;
}
