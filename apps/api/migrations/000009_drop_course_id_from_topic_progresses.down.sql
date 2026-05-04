ALTER TABLE topic_progresses ADD COLUMN course_id UUID;

UPDATE topic_progresses
SET course_id = course_section_topics.course_id
FROM course_section_topics
WHERE topic_progresses.course_section_topic_id = course_section_topics.course_section_topic_id;

ALTER TABLE topic_progresses ALTER COLUMN course_id SET NOT NULL;

ALTER TABLE topic_progresses DROP CONSTRAINT pk_topic_progresses;

ALTER TABLE topic_progresses ADD CONSTRAINT pk_topic_progresses
    PRIMARY KEY (user_id, course_id, course_section_topic_id);

ALTER TABLE topic_progresses ADD CONSTRAINT fk_topic_progresses_enrollment
    FOREIGN KEY (user_id, course_id) REFERENCES enrollments(user_id, course_id);
