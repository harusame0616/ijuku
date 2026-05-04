"use client";

import { CopyCheckIcon, CopyIcon } from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";
import { useState, type JSX } from "react";
import { Button } from "@/components/ui/button";
import { useIsHydrated } from "@/hooks/use-is-hydrated";

interface Props {
  text: string;
}

export function CopyButton({ text }: Props): JSX.Element {
  const [copied, setCopied] = useState(false);
  const isHydrated = useIsHydrated();

  async function handleCopy(): Promise<void> {
    await navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  }

  return (
    <Button
      variant="outline"
      size="icon"
      onClick={handleCopy}
      disabled={!isHydrated}
    >
      <HugeiconsIcon icon={copied ? CopyCheckIcon : CopyIcon} />
      <span className="sr-only">
        {copied ? "コピーしました" : "クリップボードにコピー"}
      </span>
    </Button>
  );
}
