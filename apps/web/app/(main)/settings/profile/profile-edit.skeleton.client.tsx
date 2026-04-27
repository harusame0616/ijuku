import type { JSX } from "react";
import { ProfileEditPresenter } from "./profile-edit.presenter.form.client";

export function ProfileEditSkeleton(): JSX.Element {
  return <ProfileEditPresenter disabled />;
}
