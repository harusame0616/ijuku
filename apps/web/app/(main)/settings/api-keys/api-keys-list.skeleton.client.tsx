import type { JSX } from "react";
import { ApiKeysListPresenter } from "./api-keys-list.presenter.client";

export function ApiKeysListSkeleton(): JSX.Element {
  return <ApiKeysListPresenter disabled />;
}
