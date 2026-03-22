-- name: GetCourses :many
SELECT course_id, title FROM courses WHERE course_id > $1 ORDER BY course_id LIMIT ($2 + 1);

-- name: GetCourse :one
SELECT
    courses.course_id,
    courses.title,
    courses.description,
    courses.slug,
    courses.tags,
    courses.publish_status,
    courses.category_id,
    courses.published_at,
    courses.author_id,
    courses.visibility,
    authors.name,
    categories.name,
    categories.path
FROM
    courses
INNER JOIN
    authors ON authors.author_id = courses.author_id
INNER JOIN
    categories ON categories.category_id = courses.category_id
WHERE
    courses.course_id = $1;
