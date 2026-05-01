package queries

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	enrollmentDetailDraftCourseID = "24000000-0000-0000-0000-0000000000ff"
	enrollmentDetailDraftSection  = "25000000-0000-0000-0000-0000000000ff"
	enrollmentDetailDraftTopic    = "26000000-0000-0000-0000-0000000000ff"
	enrollmentDetailMissingCourse = "24000000-0000-0000-0000-0000000000ee"
)

func newGetEnrollmentRequestMedium(t *testing.T, userID, courseID string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/"+userID+"/enrollments/"+courseID, nil)
	req.SetPathValue("userID", userID)
	req.SetPathValue("courseId", courseID)
	return req
}

func insertProgressWithStatus(ctx context.Context, pool *pgxpool.Pool, userID, topicID, status string, updatedAt time.Time) error {
	courseID, err := courseIDOfTopic(ctx, pool, topicID)
	if err != nil {
		return err
	}
	if err := ensureEnrollment(ctx, pool, userID, courseID, updatedAt); err != nil {
		return err
	}
	_, err = pool.Exec(ctx,
		`INSERT INTO topic_progresses (user_id, course_id, course_section_topic_id, status, _updated_at) VALUES ($1, $2, $3, $4, $5)
		 ON CONFLICT (user_id, course_id, course_section_topic_id) DO UPDATE SET status = EXCLUDED.status, _updated_at = EXCLUDED._updated_at`,
		userID, courseID, topicID, status, updatedAt,
	)
	return err
}

func insertDraftCourse(ctx context.Context, pool *pgxpool.Pool) error {
	cleanupDraftCourse(ctx, pool)
	statements := []struct {
		sql  string
		args []any
	}{
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, '草稿コース', '', 'enrollments-medium-draft', '[]', 'draft', $2, $3, 'public')`,
			[]any{enrollmentDetailDraftCourseID, mediumTestCategoryID, mediumTestAuthorID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ($1, $2, 0, 'draft section', '')`,
			[]any{enrollmentDetailDraftSection, enrollmentDetailDraftCourseID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, content)
			 VALUES ($1, $2, $3, 0, 'draft topic', '', '')`,
			[]any{enrollmentDetailDraftTopic, enrollmentDetailDraftCourseID, enrollmentDetailDraftSection},
		},
	}
	for _, s := range statements {
		if _, err := pool.Exec(ctx, s.sql, s.args...); err != nil {
			return err
		}
	}
	return nil
}

func cleanupDraftCourse(ctx context.Context, pool *pgxpool.Pool) {
	_, _ = pool.Exec(ctx, `DELETE FROM topic_progresses WHERE course_id = $1`, enrollmentDetailDraftCourseID)
	_, _ = pool.Exec(ctx, `DELETE FROM enrollments WHERE course_id = $1`, enrollmentDetailDraftCourseID)
	_, _ = pool.Exec(ctx, `DELETE FROM course_section_topics WHERE course_id = $1`, enrollmentDetailDraftCourseID)
	_, _ = pool.Exec(ctx, `DELETE FROM course_sections WHERE course_id = $1`, enrollmentDetailDraftCourseID)
	_, _ = pool.Exec(ctx, `DELETE FROM courses WHERE course_id = $1`, enrollmentDetailDraftCourseID)
}

func TestGetEnrollmentHandlerMedium(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, mediumDatabaseURL())
	if err != nil {
		t.Fatalf("DBへの接続に失敗しました: %v", err)
	}
	defer pool.Close()

	if err := setupMediumTestData(ctx, pool); err != nil {
		t.Fatalf("テストデータのセットアップに失敗しました: %v", err)
	}
	t.Cleanup(func() {
		cleanupDraftCourse(ctx, pool)
		cleanupMediumTestData(ctx, pool)
	})

	handler := NewGetEnrollmentHandler(db.New(pool))

	t.Run("存在しないcourseIdは404を返す", func(t *testing.T) {
		w := httptest.NewRecorder()
		handler.GetEnrollmentHandler(w, newGetEnrollmentRequestMedium(t, mediumTestUserID, enrollmentDetailMissingCourse))

		if w.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("404 を期待: got %d", w.Result().StatusCode)
		}
	})

	t.Run("draftコースに非著者でアクセスすると403を返す", func(t *testing.T) {
		if err := insertDraftCourse(ctx, pool); err != nil {
			t.Fatalf("draft course 挿入失敗: %v", err)
		}
		t.Cleanup(func() { cleanupDraftCourse(ctx, pool) })

		w := httptest.NewRecorder()
		handler.GetEnrollmentHandler(w, newGetEnrollmentRequestMedium(t, mediumTestUserID, enrollmentDetailDraftCourseID))

		if w.Result().StatusCode != http.StatusForbidden {
			t.Fatalf("403 を期待: got %d", w.Result().StatusCode)
		}
	})

	t.Run("draftコースに著者本人でアクセスすると200を返す", func(t *testing.T) {
		if err := insertDraftCourse(ctx, pool); err != nil {
			t.Fatalf("draft course 挿入失敗: %v", err)
		}
		t.Cleanup(func() { cleanupDraftCourse(ctx, pool) })

		w := httptest.NewRecorder()
		handler.GetEnrollmentHandler(w, newGetEnrollmentRequestMedium(t, mediumTestAuthorID, enrollmentDetailDraftCourseID))

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("200 を期待: got %d", w.Result().StatusCode)
		}
	})

	t.Run("未受講ユーザーは全statusがnullで先頭トピックがnextTopic", func(t *testing.T) {
		cleanupMediumProgresses(ctx, pool)

		w := httptest.NewRecorder()
		handler.GetEnrollmentHandler(w, newGetEnrollmentRequestMedium(t, mediumTestUserID, mediumTestCourseAID))

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("200 を期待: got %d", w.Result().StatusCode)
		}
		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if len(body.Sections) != 1 || len(body.Sections[0].Topics) != 2 {
			t.Fatalf("section/topic 構造が期待と異なる: %+v", body.Sections)
		}
		for _, topic := range body.Sections[0].Topics {
			if topic.Status != "NOT_STARTED" {
				t.Errorf("status は NOT_STARTED であること: got %v", topic.Status)
			}
		}
		if body.NextTopic == nil || body.NextTopic.TopicId == "" {
			t.Fatalf("nextTopic は先頭トピックを返す: got %+v", body.NextTopic)
		}
	})

	t.Run("最新IN_PROGRESSのtopicがnextTopic", func(t *testing.T) {
		cleanupMediumProgresses(ctx, pool)
		base := time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
		if err := insertProgressWithStatus(ctx, pool, mediumTestUserID, mediumTestTopicA0ID, "COMPLETED", base); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}
		if err := insertProgressWithStatus(ctx, pool, mediumTestUserID, mediumTestTopicA1ID, "IN_PROGRESS", base.Add(time.Hour)); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}

		w := httptest.NewRecorder()
		handler.GetEnrollmentHandler(w, newGetEnrollmentRequestMedium(t, mediumTestUserID, mediumTestCourseAID))

		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if body.NextTopic == nil {
			t.Fatalf("nextTopic は IN_PROGRESS トピックを返す")
		}
		if body.NextTopic.TopicId != mediumTestTopicA1ID {
			t.Errorf("nextTopic は A1 を期待: got %s", body.NextTopic.TopicId)
		}
	})

	t.Run("最新COMPLETEDかつ次トピックがあればそれをnextTopic", func(t *testing.T) {
		cleanupMediumProgresses(ctx, pool)
		base := time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
		if err := insertProgressWithStatus(ctx, pool, mediumTestUserID, mediumTestTopicA0ID, "COMPLETED", base); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}

		w := httptest.NewRecorder()
		handler.GetEnrollmentHandler(w, newGetEnrollmentRequestMedium(t, mediumTestUserID, mediumTestCourseAID))

		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if body.NextTopic == nil || body.NextTopic.TopicId != mediumTestTopicA1ID {
			t.Errorf("nextTopic は A1 を期待: got %+v", body.NextTopic)
		}
	})

	t.Run("全COMPLETEDの場合nextTopicはnull", func(t *testing.T) {
		cleanupMediumProgresses(ctx, pool)
		base := time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
		if err := insertProgressWithStatus(ctx, pool, mediumTestUserID, mediumTestTopicA0ID, "COMPLETED", base); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}
		if err := insertProgressWithStatus(ctx, pool, mediumTestUserID, mediumTestTopicA1ID, "COMPLETED", base.Add(time.Hour)); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}

		w := httptest.NewRecorder()
		handler.GetEnrollmentHandler(w, newGetEnrollmentRequestMedium(t, mediumTestUserID, mediumTestCourseAID))

		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		if body.NextTopic != nil {
			t.Errorf("nextTopic は null を期待: got %+v", body.NextTopic)
		}
	})

	t.Run("他ユーザーの進捗は反映されない", func(t *testing.T) {
		cleanupMediumProgresses(ctx, pool)
		base := time.Date(2026, 4, 28, 12, 0, 0, 0, time.UTC)
		if err := insertProgressWithStatus(ctx, pool, mediumTestOtherUser, mediumTestTopicA0ID, "IN_PROGRESS", base); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}

		w := httptest.NewRecorder()
		handler.GetEnrollmentHandler(w, newGetEnrollmentRequestMedium(t, mediumTestUserID, mediumTestCourseAID))

		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		for _, sec := range body.Sections {
			for _, tp := range sec.Topics {
				if tp.Status != "NOT_STARTED" {
					t.Errorf("対象ユーザーの status は全て NOT_STARTED であること: got %v", tp.Status)
				}
			}
		}
	})
}
