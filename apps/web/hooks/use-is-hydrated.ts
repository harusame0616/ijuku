"use client";

import { useSyncExternalStore } from "react";

const noop = (): void => {};

function subscribe(): () => void {
  return noop;
}

export function useIsHydrated(): boolean {
  return useSyncExternalStore(subscribe, () => true, () => false);
}
