import { readFile, readdir } from "node:fs/promises";
import { join } from "node:path";
import postgres from "postgres";
import { v5 as uuidv5 } from "uuid";
import * as v from "valibot";
import { parse as parseYaml } from "yaml";

const NAMESPACE = "a0a0a0a0-1111-4111-8111-000000000001";

const TopicSchema = v.object({
  file: v.pipe(v.string(), v.minLength(1)),
});

const SectionSchema = v.object({
  title: v.pipe(v.string(), v.minLength(1)),
  description: v.string(),
  topics: v.pipe(v.array(TopicSchema), v.minLength(1)),
});

const CourseSchema = v.object({
  title: v.pipe(v.string(), v.minLength(1)),
  description: v.string(),
  slug: v.pipe(v.string(), v.minLength(1)),
  tags: v.array(v.string()),
  publishStatus: v.picklist(["draft", "published", "archived"]),
  visibility: v.picklist(["private", "public", "paid"]),
  authorName: v.pipe(v.string(), v.minLength(1)),
  authorProfile: v.string(),
  categoryName: v.pipe(v.string(), v.minLength(1)),
  categoryPath: v.pipe(v.string(), v.minLength(1)),
  sections: v.pipe(v.array(SectionSchema), v.minLength(1)),
});

const TopicMetaSchema = v.object({
  title: v.pipe(v.string(), v.minLength(1)),
  description: v.string(),
});

type Course = v.InferOutput<typeof CourseSchema>;
type TopicMeta = v.InferOutput<typeof TopicMetaSchema>;

function deterministicUuid(scope: string, key: string): string {
  const scopeUuid = uuidv5(scope, NAMESPACE);
  return uuidv5(key, scopeUuid);
}

function parseFrontmatter(raw: string): { meta: TopicMeta; body: string } {
  const text = raw.replaceAll("\r\n", "\n");
  if (!text.startsWith("---\n")) {
    throw new Error("missing frontmatter opener");
  }
  const rest = text.slice("---\n".length);
  const closeIdx = rest.indexOf("\n---\n");
  if (closeIdx < 0) {
    throw new Error("missing frontmatter closer");
  }
  const fmText = rest.slice(0, closeIdx);
  const body = rest.slice(closeIdx + "\n---\n".length).replace(/^\n+/, "");
  const meta = v.parse(TopicMetaSchema, parseYaml(fmText));
  return { meta, body };
}

async function loadCourse(courseDir: string): Promise<Course> {
  const yamlPath = join(courseDir, "course.yaml");
  const raw = await readFile(yamlPath, "utf-8");
  return v.parse(CourseSchema, parseYaml(raw));
}

async function seedCourse(
  sql: postgres.Sql,
  courseDir: string,
): Promise<void> {
  const course = await loadCourse(courseDir);
  const authorId = deterministicUuid("author", course.authorName);
  const categoryId = deterministicUuid("category", course.categoryName);
  const courseId = deterministicUuid("course", course.slug);

  await sql.begin(async (tx) => {
    await tx`
      INSERT INTO authors (author_id, name, profile)
      VALUES (${authorId}, ${course.authorName}, ${course.authorProfile})
      ON CONFLICT (author_id) DO UPDATE SET
        name = EXCLUDED.name,
        profile = EXCLUDED.profile
    `;

    await tx`
      INSERT INTO categories (category_id, name, path)
      VALUES (${categoryId}, ${course.categoryName}, ${course.categoryPath}::ltree)
      ON CONFLICT (category_id) DO UPDATE SET
        name = EXCLUDED.name,
        path = EXCLUDED.path
    `;

    const tagsJson = JSON.stringify(course.tags);
    await tx`
      INSERT INTO courses (
        course_id, title, description, slug, tags, publish_status,
        category_id, published_at, author_id, visibility
      )
      VALUES (
        ${courseId}, ${course.title}, ${course.description}, ${course.slug},
        ${tagsJson}::jsonb, ${course.publishStatus},
        ${categoryId}, NOW(), ${authorId}, ${course.visibility}
      )
      ON CONFLICT (course_id) DO UPDATE SET
        title = EXCLUDED.title,
        description = EXCLUDED.description,
        slug = EXCLUDED.slug,
        tags = EXCLUDED.tags,
        publish_status = EXCLUDED.publish_status,
        category_id = EXCLUDED.category_id,
        published_at = EXCLUDED.published_at,
        visibility = EXCLUDED.visibility
    `;

    await tx`
      DELETE FROM topic_progresses
      WHERE course_section_topic_id IN (
        SELECT course_section_topic_id FROM course_section_topics WHERE course_id = ${courseId}
      )
    `;
    await tx`DELETE FROM course_section_topics WHERE course_id = ${courseId}`;
    await tx`DELETE FROM course_sections WHERE course_id = ${courseId}`;

    for (const [sIdx, section] of course.sections.entries()) {
      const sectionId = deterministicUuid(
        `section:${course.slug}`,
        String(sIdx + 1),
      );

      await tx`
        INSERT INTO course_sections (
          course_section_id, course_id, index, title, description
        )
        VALUES (
          ${sectionId}, ${courseId}, ${sIdx + 1}, ${section.title}, ${section.description}
        )
      `;

      for (const [tIdx, topic] of section.topics.entries()) {
        const topicId = deterministicUuid(
          `topic:${course.slug}:${sIdx + 1}`,
          String(tIdx + 1),
        );

        const mdPath = join(courseDir, "topics", topic.file);
        const mdRaw = await readFile(mdPath, "utf-8");
        const { meta, body } = parseFrontmatter(mdRaw);

        await tx`
          INSERT INTO course_section_topics (
            course_section_topic_id, course_id, course_section_id,
            index, title, description, content
          )
          VALUES (
            ${topicId}, ${courseId}, ${sectionId},
            ${tIdx + 1}, ${meta.title}, ${meta.description}, ${body}
          )
        `;
      }
    }
  });

  console.log(
    `Seeded course "${course.title}" (${courseId}) with ${course.sections.length} sections`,
  );
}

export async function seedCourses(databaseUrl: string): Promise<void> {
  const coursesDir = join(import.meta.dirname, "courses");
  let entries;
  try {
    entries = await readdir(coursesDir, { withFileTypes: true });
  } catch (err) {
    if ((err as NodeJS.ErrnoException).code === "ENOENT") {
      console.log("No courses directory found, skipping course seed");
      return;
    }
    throw err;
  }

  const sql = postgres(databaseUrl, { onnotice: () => {} });
  try {
    for (const entry of entries) {
      if (!entry.isDirectory()) continue;
      await seedCourse(sql, join(coursesDir, entry.name));
    }
  } finally {
    await sql.end();
  }
}
