import type { Metadata } from "next";
import type { JSX } from "react";
import { Suspense } from "react";
import PageLayout from "@/components/page-layout";
import { EnrollmentsListContainer } from "./enrollments-list.container.server";
import { EnrollmentsListSkeleton } from "./enrollments-list.skeleton.client";

export const metadata: Metadata = {
  title: "学習中",
};

export default function LearningPage(): JSX.Element {
  return (
    <PageLayout title="学習中">
      <Suspense fallback={<EnrollmentsListSkeleton />}>
        <EnrollmentsListContainer />
      </Suspense>
    </PageLayout>
  );
}
