import { Button } from "@/components/ui/button";
import { getUser } from "@/lib/user";
import Link from "next/link";
import { UserMenuPresenter } from "./user-menu.presenter.client";

export async function UserMenuContainer() {
  const user = await getUser();

  if (user) {
    return <UserMenuPresenter nickname={user.nickname} />;
  }

  return (
    <Button
      size="sm"
      nativeButton={false}
      render={<Link href="/login" prefetch={false} />}
    >
      ログイン
    </Button>
  );
}
