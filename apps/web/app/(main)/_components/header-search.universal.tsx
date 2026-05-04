import Form from "next/form";
import Link from "next/link";
import { useId, type JSX } from "react";

function SearchIcon({ className }: { className?: string }): JSX.Element {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
    >
      <circle cx="11" cy="11" r="8" />
      <path d="m21 21-4.35-4.35" />
    </svg>
  );
}

export function HeaderSearch(): JSX.Element {
  const inputId = useId();

  return (
    <>
      {/* Desktop: 中央に検索 input。 /courses?q=... へ遷移 */}
      <Form
        role="search"
        action="/courses"
        className="relative hidden w-full max-w-md md:flex"
      >
        <label htmlFor={inputId} className="sr-only">
          コースを検索
        </label>
        <input
          id={inputId}
          type="search"
          name="q"
          placeholder="学びたいこと、技術、キーワード…"
          className="w-full rounded-md border border-border bg-card/50 py-2 pl-3 pr-10 text-sm placeholder:text-muted-foreground/60 focus-visible:border-primary focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/30"
        />
        <button
          type="submit"
          className="absolute inset-y-0 right-0 flex items-center px-3 text-muted-foreground hover:text-foreground"
        >
          <SearchIcon className="size-3.5" />
          <span className="sr-only">検索</span>
        </button>
      </Form>

      {/* SP: 検索ボタン（コース一覧ページへの動線） */}
      <Link
        href="/courses"
        className="inline-flex size-9 items-center justify-center rounded-full text-muted-foreground hover:bg-muted hover:text-foreground md:hidden"
      >
        <SearchIcon className="size-4" />
        <span className="sr-only">コースを検索</span>
      </Link>
    </>
  );
}
