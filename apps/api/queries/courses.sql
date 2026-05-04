-- name: GetCourses :many
SELECT
    course_id AS "courseId",
    title
FROM
    courses
WHERE
    (
        courses.publish_status = 'published'
        AND courses.visibility IN ('public', 'paid')
    )
    AND (
        @Cursor :: uuid IS NULL
        OR course_id > @Cursor
    )
    AND (
        NULLIF(@Keyword :: text, '') IS NULL
        OR title LIKE '%' || @Keyword || '%'
        OR description LIKE '%' || @Keyword || '%'
    )
ORDER BY
    course_id
LIMIT
    @Size;

-- name: GetCourseById :one
SELECT
    courses.course_id,
    courses.title,
    courses.description,
    courses.slug,
    courses.tags,
    courses.publish_status,
    courses.category_id,
    categories.name AS category_name,
    courses.published_at,
    courses.author_id,
    authors.name AS author_name,
    courses.visibility,
    json_agg(
        json_build_object(
            'course_section_id',
            sections.course_section_id,
            'title',
            sections.title,
            'description',
            sections.description,
            'topics',
            sections.topics
        )
        ORDER BY
            sections."index" ASC
    ) AS sections
FROM
    courses
    JOIN categories USING (category_id)
    JOIN authors USING (author_id)
    JOIN (
        SELECT
            course_sections.course_id,
            course_sections.course_section_id,
            course_sections.title,
            course_sections.description,
            course_sections."index",
            json_agg(
                json_build_object(
                    'course_section_topic_id',
                    course_section_topic_id,
                    'title',
                    course_section_topics.title,
                    'description',
                    course_section_topics.description,
                    'content',
                    course_section_topics.content
                )
                ORDER BY
                    course_section_topics."index" ASC
            ) AS topics
        FROM
            course_sections
            JOIN course_section_topics USING (course_section_id)
        GROUP BY
            course_sections.course_id,
            course_sections.course_section_id
    ) AS sections USING (course_id)
WHERE
    courses.course_id = @CourseId :: uuid
GROUP BY
    courses.course_id,
    courses.title,
    courses.description,
    courses.slug,
    courses.tags,
    courses.publish_status,
    courses.category_id,
    categories.name,
    courses.published_at,
    courses.author_id,
    authors.name,
    courses.visibility;

-- name: GetProgressByUserIdAndCourseId :many
SELECT
    tp.course_section_topic_id,
    tp.user_id,
    tp.status,
    cs."index" AS section_index,
    cst."index" AS topic_index
FROM
    topic_progresses tp
    JOIN course_section_topics cst ON tp.course_section_topic_id = cst.course_section_topic_id
    JOIN course_sections cs ON cst.course_section_id = cs.course_section_id
WHERE
    tp.user_id = @UserId :: uuid
    AND tp.course_id = @CourseId :: uuid;

-- name: GetCourseBySlug :one
WITH
    target_course AS (
        SELECT courses.course_id
        FROM courses
            JOIN authors USING (author_id)
        WHERE courses.slug = @CourseSlug
            AND authors.slug = @AuthorSlug
    ),
    section_agg AS (
        SELECT
            sections.course_id,
            sections.course_section_id,
            sections.title,
            sections.description,
            sections.index,
            COALESCE(
                jsonb_agg(
                    jsonb_build_object(
                        'course_section_topic_id', topics.course_section_topic_id,
                        'title', topics.title,
                        'description', topics.description,
                        'content', topics."content"
                    )
                    ORDER BY topics.index
                ) FILTER (WHERE topics.course_section_topic_id IS NOT NULL),
                '[]'::jsonb
            ) AS topics
        FROM
            course_sections AS sections
            LEFT JOIN course_section_topics AS topics
                ON topics.course_section_id = sections.course_section_id
        WHERE
            sections.course_id = (SELECT course_id FROM target_course)
        GROUP BY
            sections.course_id,
            sections.course_section_id
    )
SELECT
    courses.course_id,
    courses.title,
    courses.description,
    courses.slug,
    courses.tags,
    courses.publish_status,
    courses.category_id,
    categories.name AS category_name,
    courses.published_at,
    courses.author_id,
    authors.name AS author_name,
    authors.slug AS author_slug,
    courses.visibility,
    COALESCE(
        jsonb_agg(
            jsonb_build_object(
                'course_section_id', section_agg.course_section_id,
                'title', section_agg.title,
                'description', section_agg.description,
                'topics', section_agg.topics
            )
            ORDER BY section_agg.index
        ) FILTER (WHERE section_agg.course_section_id IS NOT NULL),
        '[]'::jsonb
    ) :: jsonb AS sections
FROM
    courses
    JOIN categories USING (category_id)
    JOIN authors USING (author_id)
    LEFT JOIN section_agg ON courses.course_id = section_agg.course_id
WHERE
    courses.course_id = (SELECT course_id FROM target_course)
GROUP BY
    courses.course_id,
    courses.title,
    courses.description,
    courses.slug,
    courses.tags,
    courses.publish_status,
    courses.category_id,
    categories.name,
    courses.published_at,
    courses.author_id,
    authors.name,
    authors.slug,
    courses.visibility;

-- name: GetTopicDetail :one
SELECT
    courses.course_id AS "courseId",
    course_sections.course_section_id AS "sectionId",
    course_section_topics.course_section_topic_id as "topicId",
    course_section_topics.title,
    course_section_topics.description,
    course_section_topics.content
FROM
    course_section_topics
    JOIN course_sections USING (course_section_id)
    JOIN courses ON courses.course_id = course_sections.course_id
WHERE
    course_section_topic_id = @topic_id :: uuid
    AND course_sections.course_section_id = @section_id :: uuid
    AND course_sections.course_id = @course_id :: uuid
    AND (
        courses.publish_status = 'published'
        OR (
            @user_id :: uuid IS NOT NULL
            AND courses.author_id = @user_id :: uuid
        )
    );
