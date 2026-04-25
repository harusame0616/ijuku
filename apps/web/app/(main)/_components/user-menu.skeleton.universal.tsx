import type { JSX } from "react";
import { Skeleton } from "@/components/ui/skeleton";

export function UserMenuSkeleton(): JSX.Element {
  return <Skeleton className="size-8 rounded-full" />;
}
