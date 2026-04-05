import { Nav } from "./_components/nav";
import { Hero } from "./_components/hero";
import { Features } from "./_components/features";
import { HowItWorks } from "./_components/how-it-works";
import { LearningRecords } from "./_components/learning-records";
import { CtaSection } from "./_components/cta-section";
import { Footer } from "./_components/footer";

export default function Home() {
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
