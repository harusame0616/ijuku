import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "iJuku",
  description: "AI を活用したインターネット塾サービス",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
