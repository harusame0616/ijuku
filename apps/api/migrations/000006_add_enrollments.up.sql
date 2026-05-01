DROP TABLE user_topic_progresses;

CREATE TABLE enrollments (
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    enrolled_at TIMESTAMPTZ NOT NULL,

    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_enrollments PRIMARY KEY (user_id, course_id),
    CONSTRAINT fk_enrollments_user_id FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT fk_enrollments_course_id FOREIGN KEY (course_id) REFERENCES courses(course_id)
);

CREATE INDEX idx_enrollments_course_id ON enrollments USING BTREE (course_id);

CREATE TRIGGER trigger_enrollments_meta_updated_at BEFORE
UPDATE
    ON enrollments FOR EACH ROW EXECUTE FUNCTION update_meta_updated_at();

CREATE TABLE topic_progresses (
    user_id UUID NOT NULL,
    course_id UUID NOT NULL,
    course_section_topic_id UUID NOT NULL,
    status TEXT NOT NULL,

    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_topic_progresses PRIMARY KEY (user_id, course_id, course_section_topic_id),
    CONSTRAINT fk_topic_progresses_enrollment FOREIGN KEY (user_id, course_id) REFERENCES enrollments(user_id, course_id),
    CONSTRAINT fk_topic_progresses_course_section_topic_id FOREIGN KEY (course_section_topic_id) REFERENCES course_section_topics(course_section_topic_id)
);

COMMENT ON COLUMN topic_progresses.status IS 'IN_PROGRESS = 開始済み, COMPLETED = 完了済み, 開始していない場合はレコード自体がないので値として持たない';

CREATE INDEX idx_topic_progresses_course_section_topic_id ON topic_progresses USING BTREE (course_section_topic_id);

CREATE TRIGGER trigger_topic_progresses_meta_updated_at BEFORE
UPDATE
    ON topic_progresses FOR EACH ROW EXECUTE FUNCTION update_meta_updated_at();
