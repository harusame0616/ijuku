import type { JSX } from "react";
import { Suspense } from "react";
import { ProfileEditContainer } from "./profile-edit.container.server";
import { ProfileEditSkeleton } from "./profile-edit.skeleton.client";
import PageLayout from "@/components/page-layout";
import { type Metadata } from "next";

export const metadata: Metadata = {
  title: "プロフィール",
};

export default function SettingsPage(): JSX.Element {
  return (
    <PageLayout title="プロフィール">
      <Suspense fallback={<ProfileEditSkeleton />}>
        <ProfileEditContainer />
      </Suspense>
    </PageLayout>
  );
}
