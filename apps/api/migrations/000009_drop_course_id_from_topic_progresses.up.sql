ALTER TABLE topic_progresses DROP CONSTRAINT fk_topic_progresses_enrollment;

ALTER TABLE topic_progresses DROP CONSTRAINT pk_topic_progresses;

ALTER TABLE topic_progresses ADD CONSTRAINT pk_topic_progresses PRIMARY KEY (user_id, course_section_topic_id);

ALTER TABLE topic_progresses DROP COLUMN course_id;
