import { HelpCircleIcon } from "@hugeicons/core-free-icons";
import { HugeiconsIcon } from "@hugeicons/react";
import Link from "next/link";
import { useId, type JSX } from "react";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { CopyButton } from "./copy-button.client";
import type {
  CourseDetail,
  CourseDetailSection,
  CourseDetailTopic,
} from "./course-detail.data";
import { EnrolledBadge } from "./enrolled-badge.universal";

interface Props {
  course: CourseDetail;
}

export function CourseDetailPresenter({ course }: Props): JSX.Element {
  return (
    <article className="flex flex-col gap-6">
      <header className="flex flex-col gap-3">
        <div className="flex items-center gap-3">
          <h1 className="text-2xl font-bold text-foreground">{course.title}</h1>
          {course.isEnrolled ? <EnrolledBadge /> : null}
        </div>
        <p className="text-sm text-muted-foreground">講師: {course.author.name}</p>
        {course.description ? (
          <p className="whitespace-pre-wrap text-sm text-foreground">
            {course.description}
          </p>
        ) : null}
        {course.tags.length > 0 ? (
          <ul className="flex flex-wrap gap-1">
            {course.tags.map((tag) => (
              <li key={tag}>
                <Badge variant="outline">{tag}</Badge>
              </li>
            ))}
          </ul>
        ) : null}
      </header>

      <EnrollCommandBlock
        authorSlug={course.author.slug}
        courseSlug={course.slug}
      />

      <section className="flex flex-col gap-3">
        <h2 className="text-lg font-bold text-foreground">講座構成</h2>
        {course.sections.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            この講座にはまだセクションがありません
          </p>
        ) : (
          <ol className="flex flex-col gap-3">
            {course.sections.map((section, sectionIndex) => (
              <SectionItem
                key={section.sectionId}
                section={section}
                sectionIndex={sectionIndex}
              />
            ))}
          </ol>
        )}
      </section>
    </article>
  );
}

interface EnrollCommandBlockProps {
  authorSlug: string;
  courseSlug: string;
}

function EnrollCommandBlock({
  authorSlug,
  courseSlug,
}: EnrollCommandBlockProps): JSX.Element {
  const commandInputId = useId();
  const command = `/jukubox enroll ${authorSlug}/${courseSlug}`;

  return (
    <div className="flex flex-col gap-2">
      <div className="flex items-center gap-1">
        <Label htmlFor={commandInputId}>受講開始コマンド</Label>
        <Link
          href="/help/enrollment-guide"
          className="text-muted-foreground hover:text-foreground inline-flex items-center"
        >
          <HugeiconsIcon icon={HelpCircleIcon} className="size-4" />
          <span className="sr-only">受講開始方法の説明</span>
        </Link>
      </div>
      <div className="grid grid-cols-[1fr_auto] items-stretch gap-2">
        <Input
          id={commandInputId}
          value={command}
          readOnly
          className="font-mono text-xs"
        />
        <CopyButton text={command} />
      </div>
    </div>
  );
}

interface SectionItemProps {
  section: CourseDetailSection;
  sectionIndex: number;
}

function SectionItem({
  section,
  sectionIndex,
}: SectionItemProps): JSX.Element {
  return (
    <li className="rounded-md border border-border bg-card p-4">
      <h3 className="text-base font-bold text-foreground">
        <span className="font-normal text-muted-foreground">
          {sectionIndex + 1}.
        </span>{" "}
        {section.title}
      </h3>
      {section.description ? (
        <p className="mt-1 text-xs text-muted-foreground">
          {section.description}
        </p>
      ) : null}
      {section.topics.length > 0 ? (
        <ul className="mt-3 flex flex-col gap-2 pl-3">
          {section.topics.map((topic, topicIndex) => (
            <TopicItem
              key={topic.topicId}
              topic={topic}
              sectionIndex={sectionIndex}
              topicIndex={topicIndex}
            />
          ))}
        </ul>
      ) : null}
    </li>
  );
}

interface TopicItemProps {
  topic: CourseDetailTopic;
  sectionIndex: number;
  topicIndex: number;
}

function TopicItem({
  topic,
  sectionIndex,
  topicIndex,
}: TopicItemProps): JSX.Element {
  return (
    <li className="flex flex-col">
      <span className="text-sm font-medium text-foreground">
        <span className="font-normal text-muted-foreground">
          {sectionIndex + 1}-{topicIndex + 1}.
        </span>{" "}
        {topic.title}
      </span>
      {topic.description ? (
        <span className="text-xs text-muted-foreground">
          {topic.description}
        </span>
      ) : null}
    </li>
  );
}
