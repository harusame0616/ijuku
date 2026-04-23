import { notFound } from "next/navigation";
import { getProfile } from "./profile.data";
import { ProfileEditPresenter } from "./profile-edit.presenter.form.client";

export async function ProfileEditContainer() {
  const profile = await getProfile();
  if (!profile) notFound();

  return (
    <ProfileEditPresenter
      defaultNickname={profile.nickname}
      defaultIntroduce={profile.introduce}
    />
  );
}
