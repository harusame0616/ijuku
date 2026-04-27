import type { JSX } from "react";
import { ApiKeysListPresenter } from "./api-keys-list.presenter.client";
import { listApiKeys } from "./api-keys.data";
import { handleListApiKeysResult } from "./handle-list-api-keys-result.server";

export async function ApiKeysListContainer(): Promise<JSX.Element> {
  const apiKeys = handleListApiKeysResult(await listApiKeys());
  return <ApiKeysListPresenter apiKeys={apiKeys} />;
}
