package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	enrollMedCategoryID    = "41000000-0000-0000-0000-000000000000"
	enrollMedAuthorID      = "42000000-0000-0000-0000-000000000000"
	enrollMedUserID        = "43000000-0000-0000-0000-000000000000"
	enrollMedCourseID      = "44000000-0000-0000-0000-000000000000"
	enrollMedDraftCourseID = "44000000-0000-0000-0000-000000000010"
	enrollMedMissingCourse = "44000000-0000-0000-0000-0000000000ff"
)

func setupEnrollTestData(ctx context.Context, pool *pgxpool.Pool) error {
	cleanupEnrollTestData(ctx, pool)

	sqls := []struct {
		query string
		args  []any
	}{
		{
			`INSERT INTO categories (category_id, name, path) VALUES ($1, 'enroll カテゴリ', 'enroll-test')`,
			[]any{enrollMedCategoryID},
		},
		{
			`INSERT INTO authors (author_id, name, profile) VALUES ($1, 'enroll 著者', '')`,
			[]any{enrollMedAuthorID},
		},
		{
			`INSERT INTO users (user_id, nickname) VALUES ($1, 'enroll ユーザー')`,
			[]any{enrollMedUserID},
		},
		{
			`INSERT INTO users (user_id, nickname) VALUES ($1, 'enroll 著者ユーザー')`,
			[]any{enrollMedAuthorID},
		},
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, 'enroll コース', '', 'enroll-test-course', '[]', 'published', $2, $3, 'public')`,
			[]any{enrollMedCourseID, enrollMedCategoryID, enrollMedAuthorID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ('45000000-0000-0000-0000-000000000000', $1, 0, 'sec', '')`,
			[]any{enrollMedCourseID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, prerequisites, knowledge, flow, quiz, completion_criteria)
			 VALUES ('46000000-0000-0000-0000-000000000000', $1, '45000000-0000-0000-0000-000000000000', 0, 'topic', '', '', '', '', '', '')`,
			[]any{enrollMedCourseID},
		},
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, 'enroll ドラフトコース', '', 'enroll-test-draft', '[]', 'draft', $2, $3, 'private')`,
			[]any{enrollMedDraftCourseID, enrollMedCategoryID, enrollMedAuthorID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ('45000000-0000-0000-0000-000000000010', $1, 0, 'sec', '')`,
			[]any{enrollMedDraftCourseID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, prerequisites, knowledge, flow, quiz, completion_criteria)
			 VALUES ('46000000-0000-0000-0000-000000000010', $1, '45000000-0000-0000-0000-000000000010', 0, 'topic', '', '', '', '', '', '')`,
			[]any{enrollMedDraftCourseID},
		},
	}

	for _, s := range sqls {
		if _, err := pool.Exec(ctx, s.query, s.args...); err != nil {
			return fmt.Errorf("テストデータの挿入に失敗しました: %w", err)
		}
	}
	return nil
}

func cleanupEnrollTestData(ctx context.Context, pool *pgxpool.Pool) {
	users := []string{enrollMedUserID, enrollMedAuthorID}
	courses := []string{enrollMedCourseID, enrollMedDraftCourseID}

	_, _ = pool.Exec(ctx, `DELETE FROM topic_progresses WHERE user_id = ANY($1)`, users)
	_, _ = pool.Exec(ctx, `DELETE FROM enrollments WHERE user_id = ANY($1)`, users)
	_, _ = pool.Exec(ctx, `DELETE FROM course_section_topics WHERE course_id = ANY($1)`, courses)
	_, _ = pool.Exec(ctx, `DELETE FROM course_sections WHERE course_id = ANY($1)`, courses)
	_, _ = pool.Exec(ctx, `DELETE FROM courses WHERE course_id = ANY($1)`, courses)
	_, _ = pool.Exec(ctx, `DELETE FROM users WHERE user_id = ANY($1)`, users)
	_, _ = pool.Exec(ctx, `DELETE FROM authors WHERE author_id = $1`, enrollMedAuthorID)
	_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE category_id = $1`, enrollMedCategoryID)
}

func cleanupEnrollEnrollments(ctx context.Context, pool *pgxpool.Pool) {
	users := []string{enrollMedUserID, enrollMedAuthorID}
	_, _ = pool.Exec(ctx, `DELETE FROM enrollments WHERE user_id = ANY($1)`, users)
}

func TestPostEnrollmentHandlerMedium(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, env.Require("DATABASE_URL"))
	if err != nil {
		t.Fatalf("DBへの接続に失敗しました: %v", err)
	}
	defer pool.Close()

	if err := setupEnrollTestData(ctx, pool); err != nil {
		t.Fatalf("テストデータのセットアップに失敗しました: %v", err)
	}
	t.Cleanup(func() { cleanupEnrollTestData(ctx, pool) })

	q := db.New(pool)
	courseRepo := NewSqrcCourseRepository(q)
	enrollmentRepo := NewSqrcEnrollmentRepository(q)
	handler := NewEnrollHandler(NewEnrollUsecase(courseRepo, enrollmentRepo))

	newReq := func(t *testing.T, userID, body string) *http.Request {
		t.Helper()
		req := httptest.NewRequest("POST", "/v1/users/"+userID+"/enrollments", strings.NewReader(body))
		req.SetPathValue("userID", userID)
		return req
	}

	enrollmentExists := func(t *testing.T, userID, courseID string) bool {
		t.Helper()
		var count int
		if err := pool.QueryRow(ctx,
			`SELECT COUNT(*) FROM enrollments WHERE user_id = $1 AND course_id = $2`,
			userID, courseID,
		).Scan(&count); err != nil {
			t.Fatalf("enrollment存在確認失敗: %v", err)
		}
		return count > 0
	}

	t.Run("公開コースの受講開始は201で成功する", func(t *testing.T) {
		t.Cleanup(func() { cleanupEnrollEnrollments(ctx, pool) })

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newReq(t, enrollMedUserID, `{"courseId":"`+enrollMedCourseID+`"}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["courseId"] != enrollMedCourseID {
			t.Errorf("courseIdが一致すること: got %q", body["courseId"])
		}
		if body["enrolledAt"] == "" {
			t.Errorf("enrolledAtが返ること")
		}
		if !enrollmentExists(t, enrollMedUserID, enrollMedCourseID) {
			t.Errorf("DBにenrollmentが作成されていること")
		}
	})

	t.Run("既受講コースに再度受講開始すると409 ALREADY_ENROLLEDを返す", func(t *testing.T) {
		t.Cleanup(func() { cleanupEnrollEnrollments(ctx, pool) })

		w1 := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w1, newReq(t, enrollMedUserID, `{"courseId":"`+enrollMedCourseID+`"}`))
		if w1.Result().StatusCode != http.StatusCreated {
			t.Fatalf("初回は201であること: got %d", w1.Result().StatusCode)
		}

		w2 := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w2, newReq(t, enrollMedUserID, `{"courseId":"`+enrollMedCourseID+`"}`))

		if w2.Result().StatusCode != http.StatusConflict {
			t.Errorf("ステータスコードが409であること: got %d", w2.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w2.Result().Body).Decode(&body)
		if body["errorCode"] != "ALREADY_ENROLLED" {
			t.Errorf("errorCodeがALREADY_ENROLLEDであること: got %q", body["errorCode"])
		}
	})

	t.Run("draftコースに非著者がアクセスすると403を返す", func(t *testing.T) {
		t.Cleanup(func() { cleanupEnrollEnrollments(ctx, pool) })

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newReq(t, enrollMedUserID, `{"courseId":"`+enrollMedDraftCourseID+`"}`))

		if w.Result().StatusCode != http.StatusForbidden {
			t.Errorf("ステータスコードが403であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["errorCode"] != "ENROLLMENT_FORBIDDEN" {
			t.Errorf("errorCodeがENROLLMENT_FORBIDDENであること: got %q", body["errorCode"])
		}
		if enrollmentExists(t, enrollMedUserID, enrollMedDraftCourseID) {
			t.Errorf("DBにenrollmentが作成されていないこと")
		}
	})

	t.Run("draftコースに著者本人がアクセスすると201で成功する", func(t *testing.T) {
		t.Cleanup(func() { cleanupEnrollEnrollments(ctx, pool) })

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newReq(t, enrollMedAuthorID, `{"courseId":"`+enrollMedDraftCourseID+`"}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Errorf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		if !enrollmentExists(t, enrollMedAuthorID, enrollMedDraftCourseID) {
			t.Errorf("DBにenrollmentが作成されていること")
		}
	})

	t.Run("存在しないcourseIdは404 COURSE_NOT_FOUNDを返す", func(t *testing.T) {
		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newReq(t, enrollMedUserID, `{"courseId":"`+enrollMedMissingCourse+`"}`))

		if w.Result().StatusCode != http.StatusNotFound {
			t.Errorf("ステータスコードが404であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["errorCode"] != "COURSE_NOT_FOUND" {
			t.Errorf("errorCodeがCOURSE_NOT_FOUNDであること: got %q", body["errorCode"])
		}
	})
}
