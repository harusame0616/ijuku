import { handleGetProfileResult } from "./handle-get-profile-result.server";
import { getProfile } from "./profile.data";
import { ProfileEditPresenter } from "./profile-edit.presenter.form.client";

export async function ProfileEditContainer() {
  const profile = handleGetProfileResult(await getProfile());

  return (
    <ProfileEditPresenter
      defaultNickname={profile.nickname}
      defaultIntroduce={profile.introduce}
    />
  );
}
