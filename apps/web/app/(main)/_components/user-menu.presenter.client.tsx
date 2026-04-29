"use client";

import type { JSX } from "react";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { User02Icon } from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";
import Link from "next/link";

interface Props {
  nickname: string;
}

export function UserMenuPresenter({ nickname }: Props): JSX.Element {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger>
        <Avatar className="cursor-pointer">
          <AvatarFallback>
            {nickname.charAt(0) || (
              <HugeiconsIcon icon={User02Icon} strokeWidth={1.5} />
            )}
          </AvatarFallback>
        </Avatar>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem render={<Link href={"/learning"} />}>
          学習中
        </DropdownMenuItem>
        <DropdownMenuItem render={<Link href={"/settings"} />}>
          設定
        </DropdownMenuItem>
        <DropdownMenuItem>ログアウト</DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
