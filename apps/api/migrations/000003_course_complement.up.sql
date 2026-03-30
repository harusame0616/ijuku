ALTER TABLE
    course_sections
ADD
    COLUMN title TEXT NOT NULL DEFAULT '';

ALTER TABLE
    course_sections
ALTER COLUMN
    title DROP DEFAULT;

ALTER TABLE
    course_sections
ADD
    COLUMN description TEXT NOT NULL DEFAULT '';

ALTER TABLE
    course_sections
ALTER COLUMN
    description DROP DEFAULT;

ALTER TABLE
    course_section_topics
ADD
    COLUMN title TEXT NOT NULL DEFAULT '';

ALTER TABLE
    course_section_topics
ALTER COLUMN
    title DROP DEFAULT;

ALTER TABLE
    course_section_topics
ADD
    COLUMN description TEXT NOT NULL DEFAULT '';

ALTER TABLE
    course_section_topics
ALTER COLUMN
    description DROP DEFAULT;

ALTER TABLE
    course_section_topics
ADD
    COLUMN prerequisites TEXT NOT NULL DEFAULT '';

ALTER TABLE
    course_section_topics
ALTER COLUMN
    prerequisites DROP DEFAULT;

ALTER TABLE
    course_section_topics
ADD
    COLUMN knowledge TEXT NOT NULL DEFAULT '';

ALTER TABLE
    course_section_topics
ALTER COLUMN
    knowledge DROP DEFAULT;

ALTER TABLE
    course_section_topics
ADD
    COLUMN flow TEXT NOT NULL DEFAULT '';

ALTER TABLE
    course_section_topics
ALTER COLUMN
    flow DROP DEFAULT;

ALTER TABLE
    course_section_topics
ADD
    COLUMN quiz TEXT NOT NULL DEFAULT '';

ALTER TABLE
    course_section_topics
ALTER COLUMN
    quiz DROP DEFAULT;

ALTER TABLE
    course_section_topics
ADD
    COLUMN completion_criteria TEXT NOT NULL DEFAULT '';

ALTER TABLE
    course_section_topics
ALTER COLUMN
    completion_criteria DROP DEFAULT;

ALTER TABLE user_topic_progresses DROP CONSTRAINT fk_course_id;

ALTER TABLE user_topic_progresses DROP CONSTRAINT fk_course_section_id;

DROP INDEX idx_user_topic_progresses_course_id;

DROP INDEX idx_user_topic_progresses_course_section_id;

ALTER TABLE user_topic_progresses DROP COLUMN course_id;

ALTER TABLE user_topic_progresses DROP COLUMN course_section_id;
