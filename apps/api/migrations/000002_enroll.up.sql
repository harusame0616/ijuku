-- USERS
CREATE TABLE users (
    user_id UUID,
    nickname TEXT NOT NULL,
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_users PRIMARY KEY (user_id)
);

CREATE TRIGGER trigger_users_meta_updated_at BEFORE
UPDATE
    ON users FOR EACH ROW EXECUTE FUNCTION update_meta_updated_at();

-- APIKEYS
CREATE TABLE apikeys (
    key_hash TEXT,
    user_id UUID NOT NULL,
    plain_suffix TEXT NOT NULL,
    expired_at TIMESTAMPTZ NOT NULL,
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_apikeys PRIMARY KEY (key_hash),
    CONSTRAINT fk_apikeys_user_id FOREIGN KEY (user_id) REFERENCES users(user_id)
);

CREATE INDEX idx_apikeys_user_id ON apikeys USING BTREE (user_id);

CREATE TRIGGER trigger_apikeys_meta_updated_at BEFORE
UPDATE
    ON apikeys FOR EACH ROW EXECUTE FUNCTION update_meta_updated_at();

-- COURSE SECTION
CREATE TABLE course_sections (
    course_section_id UUID,
    course_id UUID NOT NULL,
    index SMALLINT NOT NULL,
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_course_sections PRIMARY KEY (course_section_id),
    CONSTRAINT fk_course_sections_course_id FOREIGN KEY (course_id) REFERENCES courses(course_id),
    CONSTRAINT uq_course_sections_course_id_index UNIQUE (course_id, index) DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX idx_course_sections_course_id ON course_sections USING BTREE (course_id);

CREATE TRIGGER trigger_course_sections_meta_updated_at BEFORE
UPDATE
    ON course_sections FOR EACH ROW EXECUTE FUNCTION update_meta_updated_at();

-- COURSE SECTION TOPIC
CREATE TABLE course_section_topics (
    course_section_topic_id UUID,
    -- course_id は course_section と JOIN することで辿れるが、クエリ効率のため非正規化
    course_id UUID NOT NULL,
    course_section_id UUID NOT NULL,
    index SMALLINT NOT NULL,
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_course_section_topics PRIMARY KEY (course_section_topic_id),
    CONSTRAINT fk_course_section_topics_course_id FOREIGN KEY (course_id) REFERENCES courses(course_id),
    CONSTRAINT fk_course_section_topics_course_section_id FOREIGN KEY (course_section_id) REFERENCES course_sections(course_section_id),
    CONSTRAINT uq_course_section_topics_course_section_id_index UNIQUE (course_section_id, index) DEFERRABLE INITIALLY DEFERRED
);

CREATE INDEX idx_course_section_topics_course_id ON course_section_topics USING BTREE (course_id);

CREATE INDEX idx_course_section_topics_course_section_id ON course_section_topics USING BTREE (course_section_id);

CREATE TRIGGER trigger_course_section_topics_meta_updated_at BEFORE
UPDATE
    ON course_section_topics FOR EACH ROW EXECUTE FUNCTION update_meta_updated_at();

CREATE TABLE user_topic_progresses (
    -- クエリ効率のため非正規化
    course_id UUID NOT NULL,
    course_section_id UUID NOT NULL,
    course_section_topic_id UUID NOT NULL,
    user_id UUID NOT NULL,
    status TEXT NOT NULL,
    _created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    _updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT pk_user_topic_progresses PRIMARY KEY (course_section_topic_id, user_id),
    CONSTRAINT fk_course_id FOREIGN KEY (course_id) REFERENCES courses(course_id),
    CONSTRAINT fk_course_section_id FOREIGN KEY (course_section_id) REFERENCES course_sections(course_section_id),
    CONSTRAINT fk_course_section_topic_id FOREIGN KEY (course_section_topic_id) REFERENCES course_section_topics(course_section_topic_id),
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id)
);

COMMENT ON COLUMN user_topic_progresses.status IS 'IN_PROGRESS = 開始済み, COMPLETED = 完了済み, 開始していない場合はレコード自体がないので値として持たない';

CREATE INDEX idx_user_topic_progresses_course_id ON user_topic_progresses USING BTREE (course_id);

CREATE INDEX idx_user_topic_progresses_course_section_id ON user_topic_progresses USING BTREE (course_section_id);

CREATE INDEX idx_user_topic_progresses_course_section_topic_id ON user_topic_progresses USING BTREE (course_section_id);

CREATE INDEX idx_user_topic_progresses_user_id ON user_topic_progresses USING BTREE (user_id);

CREATE TRIGGER trigger_user_topic_progress_meta_updated_at BEFORE
UPDATE
    ON user_topic_progresses FOR EACH ROW EXECUTE FUNCTION update_meta_updated_at();
