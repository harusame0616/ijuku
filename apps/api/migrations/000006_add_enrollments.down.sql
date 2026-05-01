DROP TABLE topic_progresses;

DROP TABLE enrollments;

CREATE TABLE user_topic_progresses (
    course_section_topic_id UUID NOT NULL,
    user_id UUID NOT NULL,
    status TEXT NOT NULL,

    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT pk_user_topic_progresses PRIMARY KEY (course_section_topic_id, user_id),
    CONSTRAINT fk_course_section_topic_id FOREIGN KEY (course_section_topic_id) REFERENCES course_section_topics(course_section_topic_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id)
);

COMMENT ON COLUMN user_topic_progresses.status IS 'IN_PROGRESS = 開始済み, COMPLETED = 完了済み, 開始していない場合はレコード自体がないので値として持たない';

CREATE INDEX idx_user_topic_progresses_course_section_topic_id ON user_topic_progresses USING BTREE (course_section_topic_id);

CREATE INDEX idx_user_topic_progresses_user_id ON user_topic_progresses USING BTREE (user_id);

CREATE TRIGGER trigger_user_topic_progress_meta_updated_at BEFORE
UPDATE
    ON user_topic_progresses FOR EACH ROW EXECUTE FUNCTION update_meta_updated_at();
