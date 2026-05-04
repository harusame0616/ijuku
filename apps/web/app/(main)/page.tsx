import type { Metadata } from "next";
import type { JSX } from "react";
import { FeaturedCoursesSection } from "./_top-page/featured-courses-section.universal";
import { Hero } from "./_top-page/hero.universal";
import { LpLinkFooter } from "./_top-page/lp-link-footer.universal";
import { SideBHero } from "./_top-page/side-b-hero.universal";

export const metadata: Metadata = {
  title: "JukuBox",
};

export default function TopPage(): JSX.Element {
  return (
    <div className="mx-auto flex max-w-6xl flex-col gap-16 pb-20 lg:gap-20">
      {/* SIDE A — 受講側 */}
      <Hero />
      <FeaturedCoursesSection />

      {/* SIDE B — 制作側 */}
      <SideBHero />

      <LpLinkFooter />
    </div>
  );
}
