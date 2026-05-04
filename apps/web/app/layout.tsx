import type { JSX } from "react";
import { cn } from "@/lib/utilities";
import type { Metadata } from "next";
import { Orbitron } from "next/font/google";
import "./globals.css";

const orbitron = Orbitron({
  subsets: ["latin"],
  variable: "--font-orbitron",
  weight: ["400", "700", "900"],
});

export const metadata: Metadata = {
  title: "Jukubox.ai",
  description: "あなたの AI で学ぶ AI 学習プラットフォーム",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>): JSX.Element {
  return (
    <html lang="ja" className={cn("h-full", "font-sans", orbitron.variable)}>
      <body className="flex min-h-full flex-col">{children}</body>
    </html>
  );
}
