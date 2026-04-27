"use client";

import type { JSX } from "react";
import { Skeleton } from "@/components/ui/skeleton";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { ApiKey } from "./api-keys.data";

const API_KEY_MASK = "jukubox_••••";

const dateFormatter = new Intl.DateTimeFormat("ja-JP", {
  dateStyle: "medium",
});

function formatDate(value: string): string {
  if (value === "") return "—";
  const parsed = new Date(value);
  if (Number.isNaN(parsed.getTime())) return value;
  return dateFormatter.format(parsed);
}

function formatExpiredAt(value: string | null): string {
  if (value === null) return "無期限";
  return formatDate(value);
}

function formatApiKey(suffix: string): string {
  return `${API_KEY_MASK}${suffix}`;
}

interface Props {
  apiKeys?: ApiKey[];
  disabled?: boolean;
}

const SKELETON_ROWS: ApiKey[] = Array.from({ length: 3 }, (_, index) => ({
  apiKeyId: `skeleton-${index}`,
  suffix: "",
  createdAt: "",
  expiredAt: null,
}));

export function ApiKeysListPresenter({
  apiKeys = [],
  disabled = false,
}: Props): JSX.Element {
  const rows = disabled ? SKELETON_ROWS : apiKeys;

  if (!disabled && rows.length === 0) {
    return (
      <p className="text-muted-foreground text-sm">
        API キーがまだ登録されていません。
      </p>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>API キー</TableHead>
          <TableHead>作成日</TableHead>
          <TableHead>有効期限</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {rows.map((apiKey) => (
          <TableRow key={apiKey.apiKeyId}>
            <TableCell className="font-mono">
              {disabled ? (
                <Skeleton className="h-4 w-32" />
              ) : (
                formatApiKey(apiKey.suffix)
              )}
            </TableCell>
            <TableCell>
              {disabled ? (
                <Skeleton className="h-4 w-24" />
              ) : (
                formatDate(apiKey.createdAt)
              )}
            </TableCell>
            <TableCell>
              {disabled ? (
                <Skeleton className="h-4 w-24" />
              ) : (
                formatExpiredAt(apiKey.expiredAt)
              )}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
