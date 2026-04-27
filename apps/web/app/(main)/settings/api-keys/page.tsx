import type { Metadata } from "next";
import type { JSX } from "react";
import { Suspense } from "react";
import PageLayout from "@/components/page-layout";
import { ApiKeysListContainer } from "./api-keys-list.container.server";
import { ApiKeysListSkeleton } from "./api-keys-list.skeleton.client";

export const metadata: Metadata = {
  title: "API キー",
};

export default function ApiKeysPage(): JSX.Element {
  return (
    <PageLayout title="API キー">
      <Suspense fallback={<ApiKeysListSkeleton />}>
        <ApiKeysListContainer />
      </Suspense>
    </PageLayout>
  );
}
