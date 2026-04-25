import type { JSX } from "react";
import { Nav } from "@/app/(lp)/lp/_components/nav";
import { Hero } from "@/app/(lp)/lp/_components/hero";
import { Features } from "@/app/(lp)/lp/_components/features";
import { HowItWorks } from "@/app/(lp)/lp/_components/how-it-works";
import { LearningRecords } from "@/app/(lp)/lp/_components/learning-records";
import { CtaSection } from "@/app/(lp)/lp/_components/cta-section";
import { Footer } from "@/app/(lp)/lp/_components/footer";

export default function Home(): JSX.Element {
  return (
    <div className="bg-background min-h-screen">
      <Nav />
      <main>
        <Hero />
        <Features />
        <HowItWorks />
        <LearningRecords />
        <CtaSection />
      </main>
      <Footer />
    </div>
  );
}
