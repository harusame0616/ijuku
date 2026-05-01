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
                    'prerequisites',
                    course_section_topics.prerequisites,
                    'knowledge',
                    course_section_topics.knowledge,
                    'flow',
                    course_section_topics.flow,
                    'quiz',
                    course_section_topics.quiz,
                    'completion_criteria',
                    course_section_topics.completion_criteria
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

-- name: GetTopicDetail :one
SELECT
    courses.course_id AS "courseId",
    course_sections.course_section_id AS "sectionId",
    course_section_topics.course_section_topic_id as "topicId",
    course_section_topics.title,
    course_section_topics.description,
    course_section_topics.prerequisites,
    course_section_topics.knowledge,
    course_section_topics.flow,
    course_section_topics.quiz,
    course_section_topics.completion_criteria AS "completionCriteria"
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
