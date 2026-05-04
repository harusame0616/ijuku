import type { JSX } from "react";
import type { Metadata } from "next";
import { ContactForm } from "@/app/(lp)/contacts/_components/contact-form.client";

export const metadata: Metadata = {
  title: "お問い合わせ",
};

export default function ContactsPage(): JSX.Element {
  return (
    <main className="mx-auto flex w-full max-w-2xl flex-col gap-8 px-6 py-16">
      <header className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold">お問い合わせ</h1>
        <p className="text-sm text-muted-foreground">
          ご質問・ご要望などございましたら、以下のフォームよりお問い合わせください。
        </p>
      </header>
      <ContactForm />
    </main>
  );
}
