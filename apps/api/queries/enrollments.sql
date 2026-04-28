-- name: GetEnrollmentsByUserID :many
SELECT
    courses.course_id AS "courseId",
    courses.title
FROM
    user_topic_progresses
    JOIN course_section_topics ON user_topic_progresses.course_section_topic_id = course_section_topics.course_section_topic_id
    JOIN course_sections ON course_section_topics.course_section_id = course_sections.course_section_id
    JOIN courses ON course_sections.course_id = courses.course_id
WHERE
    user_topic_progresses.user_id = @UserID :: uuid
GROUP BY
    courses.course_id,
    courses.title
ORDER BY
    MAX(user_topic_progresses._updated_at) DESC;

-- name: GetCourseAuthorityById :one
SELECT
    publish_status,
    author_id
FROM
    courses
WHERE
    courses.course_id = @CourseID;

-- name: GetCourseStructureWithProgress :one
WITH section_agg AS (
    SELECT
        sections.course_id,
        sections.course_section_id,
        sections.title,
        sections.index,
        COALESCE(
            jsonb_agg(
                jsonb_build_object(
                    'topicId',
                    topics.course_section_topic_id,
                    'title',
                    topics.title,
                    'status',
                    COALESCE(progresses.status, 'NOT_STARTED'),
                    'index',
                    topics.index
                )
                ORDER BY
                    topics.index
            ) FILTER (
                WHERE
                    topics.course_section_topic_id IS NOT NULL
            ),
            '[]' :: jsonb
        ) AS topics
    FROM
        course_sections AS sections
        LEFT JOIN course_section_topics AS topics ON sections.course_section_id = topics.course_section_id
        LEFT JOIN user_topic_progresses AS progresses ON progresses.course_section_topic_id = topics.course_section_topic_id
        AND progresses.user_id = @UserID
    WHERE
        sections.course_id = @CourseID
    GROUP BY
        sections.course_id,
        sections.course_section_id,
        sections.index,
        sections.title
)
SELECT
    courses.course_id,
    courses.title,
    COALESCE(
        jsonb_agg(
            jsonb_build_object(
                'sectionId',
                section_agg.course_section_id,
                'title',
                section_agg.title,
                'index',
                section_agg.index,
                'topics',
                section_agg.topics
            )
            ORDER BY
                section_agg.index
        ) FILTER (
            WHERE
                section_agg.course_section_id IS NOT NULL
        ),
        '[]' :: jsonb
    ) :: jsonb AS sections
FROM
    courses
    LEFT JOIN section_agg ON courses.course_id = section_agg.course_id
WHERE
    courses.course_id = @CourseID
GROUP BY
    courses.course_id;
