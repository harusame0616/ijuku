import { Suspense } from "react";
import { ProfileEditContainer } from "./profile-edit.container.server";
import { ProfileEditSkeleton } from "./profile-edit.skeleton.client";
import PageLayout from "@/components/page-layout";
import { Metadata } from "next";

export const metadata: Metadata = {
  title: "プロフィール",
};

export default function SettingsPage() {
  return (
    <PageLayout title="プロフィール">
      <Suspense fallback={<ProfileEditSkeleton />}>
        <ProfileEditContainer />
      </Suspense>
    </PageLayout>
  );
}
