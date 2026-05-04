package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	mediumTestCategoryID = "21000000-0000-0000-0000-000000000000"
	mediumTestAuthorID   = "22000000-0000-0000-0000-000000000000"
	mediumTestUserID     = "23000000-0000-0000-0000-000000000000"
	mediumTestOtherUser  = "23000000-0000-0000-0000-000000000001"

	mediumTestCourseAID  = "24000000-0000-0000-0000-000000000001"
	mediumTestCourseBID  = "24000000-0000-0000-0000-000000000002"
	mediumTestCourseCID  = "24000000-0000-0000-0000-000000000003"
	mediumTestSectionAID = "25000000-0000-0000-0000-000000000001"
	mediumTestSectionBID = "25000000-0000-0000-0000-000000000002"
	mediumTestSectionCID = "25000000-0000-0000-0000-000000000003"
	mediumTestTopicA0ID  = "26000000-0000-0000-0000-000000000001"
	mediumTestTopicA1ID  = "26000000-0000-0000-0000-000000000002"
	mediumTestTopicB0ID  = "26000000-0000-0000-0000-000000000003"
	mediumTestTopicC0ID  = "26000000-0000-0000-0000-000000000004"
)

func mediumDatabaseURL() string {
	return env.Require("DATABASE_URL")
}

func setupMediumTestData(ctx context.Context, pool *pgxpool.Pool) error {
	cleanupMediumTestData(ctx, pool)

	sqls := []struct {
		query string
		args  []any
	}{
		{
			`INSERT INTO categories (category_id, name, path) VALUES ($1, 'enrollments medium カテゴリ', 'enrollments-medium')`,
			[]any{mediumTestCategoryID},
		},
		{
			`INSERT INTO authors (author_id, name, profile, slug) VALUES ($1, 'enrollments medium 著者', '', 'enrollments-medium-author')`,
			[]any{mediumTestAuthorID},
		},
		{
			`INSERT INTO users (user_id, nickname) VALUES ($1, 'enrollments medium ユーザー')`,
			[]any{mediumTestUserID},
		},
		{
			`INSERT INTO users (user_id, nickname) VALUES ($1, 'enrollments medium 他ユーザー')`,
			[]any{mediumTestOtherUser},
		},
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, 'コースA', 'A 説明', 'enrollments-medium-a', '[]', 'published', $2, $3, 'public')`,
			[]any{mediumTestCourseAID, mediumTestCategoryID, mediumTestAuthorID},
		},
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, 'コースB', 'B 説明', 'enrollments-medium-b', '[]', 'published', $2, $3, 'public')`,
			[]any{mediumTestCourseBID, mediumTestCategoryID, mediumTestAuthorID},
		},
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, 'コースC', 'C 説明', 'enrollments-medium-c', '[]', 'published', $2, $3, 'public')`,
			[]any{mediumTestCourseCID, mediumTestCategoryID, mediumTestAuthorID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ($1, $2, 0, 'セクションA', '')`,
			[]any{mediumTestSectionAID, mediumTestCourseAID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ($1, $2, 0, 'セクションB', '')`,
			[]any{mediumTestSectionBID, mediumTestCourseBID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ($1, $2, 0, 'セクションC', '')`,
			[]any{mediumTestSectionCID, mediumTestCourseCID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, content)
			 VALUES ($1, $2, $3, 0, 'A0', '', '')`,
			[]any{mediumTestTopicA0ID, mediumTestCourseAID, mediumTestSectionAID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, content)
			 VALUES ($1, $2, $3, 1, 'A1', '', '')`,
			[]any{mediumTestTopicA1ID, mediumTestCourseAID, mediumTestSectionAID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, content)
			 VALUES ($1, $2, $3, 0, 'B0', '', '')`,
			[]any{mediumTestTopicB0ID, mediumTestCourseBID, mediumTestSectionBID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, content)
			 VALUES ($1, $2, $3, 0, 'C0', '', '')`,
			[]any{mediumTestTopicC0ID, mediumTestCourseCID, mediumTestSectionCID},
		},
	}

	for _, s := range sqls {
		if _, err := pool.Exec(ctx, s.query, s.args...); err != nil {
			return fmt.Errorf("テストデータの挿入に失敗しました: %w", err)
		}
	}
	return nil
}

func cleanupMediumTestData(ctx context.Context, pool *pgxpool.Pool) {
	users := []string{mediumTestUserID, mediumTestOtherUser, mediumTestAuthorID}
	courses := []string{mediumTestCourseAID, mediumTestCourseBID, mediumTestCourseCID}

	_, _ = pool.Exec(ctx, `DELETE FROM topic_progresses WHERE user_id = ANY($1)`, users)
	_, _ = pool.Exec(ctx, `DELETE FROM enrollments WHERE user_id = ANY($1)`, users)
	_, _ = pool.Exec(ctx, `DELETE FROM course_section_topics WHERE course_id = ANY($1)`, courses)
	_, _ = pool.Exec(ctx, `DELETE FROM course_sections WHERE course_id = ANY($1)`, courses)
	_, _ = pool.Exec(ctx, `DELETE FROM courses WHERE course_id = ANY($1)`, courses)
	_, _ = pool.Exec(ctx, `DELETE FROM users WHERE user_id = ANY($1)`, users)
	_, _ = pool.Exec(ctx, `DELETE FROM authors WHERE author_id = $1`, mediumTestAuthorID)
	_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE category_id = $1`, mediumTestCategoryID)
}

func cleanupMediumProgresses(ctx context.Context, pool *pgxpool.Pool) {
	users := []string{mediumTestUserID, mediumTestOtherUser}
	_, _ = pool.Exec(ctx, `DELETE FROM topic_progresses WHERE user_id = ANY($1)`, users)
	_, _ = pool.Exec(ctx, `DELETE FROM enrollments WHERE user_id = ANY($1)`, users)
}

// ensureEnrollment は (userID, courseID) の enrollment を取得 or 新規作成する。
func ensureEnrollment(ctx context.Context, pool *pgxpool.Pool, userID, courseID string, enrolledAt time.Time) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO enrollments (user_id, course_id, enrolled_at) VALUES ($1, $2, $3)
		 ON CONFLICT (user_id, course_id) DO NOTHING`,
		userID, courseID, enrolledAt,
	)
	return err
}

// insertProgress は topic に対する IN_PROGRESS 進捗を記録する。
// 既に enrollment が無ければ自動作成。
func insertProgress(ctx context.Context, pool *pgxpool.Pool, userID, topicID string, updatedAt time.Time) error {
	courseID, err := courseIDOfTopic(ctx, pool, topicID)
	if err != nil {
		return err
	}
	if err := ensureEnrollment(ctx, pool, userID, courseID, updatedAt); err != nil {
		return err
	}
	_, err = pool.Exec(ctx,
		`INSERT INTO topic_progresses (user_id, course_section_topic_id, status, _updated_at) VALUES ($1, $2, 'IN_PROGRESS', $3)
		 ON CONFLICT (user_id, course_section_topic_id) DO UPDATE SET status = EXCLUDED.status, _updated_at = EXCLUDED._updated_at`,
		userID, topicID, updatedAt,
	)
	return err
}

func courseIDOfTopic(ctx context.Context, pool *pgxpool.Pool, topicID string) (string, error) {
	var courseID string
	err := pool.QueryRow(ctx,
		`SELECT course_id FROM course_section_topics WHERE course_section_topic_id = $1`,
		topicID,
	).Scan(&courseID)
	return courseID, err
}

func TestGetEnrollmentsHandlerMedium(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, mediumDatabaseURL())
	if err != nil {
		t.Fatalf("DBへの接続に失敗しました: %v", err)
	}
	defer pool.Close()

	if err := setupMediumTestData(ctx, pool); err != nil {
		t.Fatalf("テストデータのセットアップに失敗しました: %v", err)
	}
	t.Cleanup(func() { cleanupMediumTestData(ctx, pool) })

	handler := NewGetEnrollmentsHandler(db.New(pool))

	t.Run("受講中コースがない場合は空配列が返る", func(t *testing.T) {
		cleanupMediumProgresses(ctx, pool)

		w := httptest.NewRecorder()
		handler.GetEnrollmentsHandler(w, newGetEnrollmentsRequest(t, mediumTestUserID))

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", w.Result().StatusCode)
		}
		var body GetEnrollmentsResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
		}
		if len(body.Enrollments) != 0 {
			t.Errorf("受講中コースが0件であること: got %d", len(body.Enrollments))
		}
	})

	t.Run("複数コース受講時 lastProgressedAt DESC で返る", func(t *testing.T) {
		cleanupMediumProgresses(ctx, pool)

		// A は最古、B は中間、C は最新
		base := time.Date(2026, 4, 25, 12, 0, 0, 0, time.UTC)
		if err := insertProgress(ctx, pool, mediumTestUserID, mediumTestTopicA0ID, base); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}
		if err := insertProgress(ctx, pool, mediumTestUserID, mediumTestTopicB0ID, base.Add(24*time.Hour)); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}
		if err := insertProgress(ctx, pool, mediumTestUserID, mediumTestTopicC0ID, base.Add(48*time.Hour)); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}

		w := httptest.NewRecorder()
		handler.GetEnrollmentsHandler(w, newGetEnrollmentsRequest(t, mediumTestUserID))

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", w.Result().StatusCode)
		}
		var body GetEnrollmentsResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
		}
		if len(body.Enrollments) != 3 {
			t.Fatalf("受講中コースが3件であること: got %d", len(body.Enrollments))
		}
		expected := []string{"コースC", "コースB", "コースA"}
		for i, want := range expected {
			if body.Enrollments[i].Title != want {
				t.Errorf("並び順 %d 番目が %s であること: got %s", i, want, body.Enrollments[i].Title)
			}
		}
	})

	t.Run("同一コースの複数トピックは1件にまとめられる", func(t *testing.T) {
		cleanupMediumProgresses(ctx, pool)

		base := time.Date(2026, 4, 26, 12, 0, 0, 0, time.UTC)
		if err := insertProgress(ctx, pool, mediumTestUserID, mediumTestTopicA0ID, base); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}
		if err := insertProgress(ctx, pool, mediumTestUserID, mediumTestTopicA1ID, base.Add(time.Hour)); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}

		w := httptest.NewRecorder()
		handler.GetEnrollmentsHandler(w, newGetEnrollmentsRequest(t, mediumTestUserID))

		var body GetEnrollmentsResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
		}
		if len(body.Enrollments) != 1 {
			t.Fatalf("コースA が 1 件にまとめられること: got %d", len(body.Enrollments))
		}
		if body.Enrollments[0].Title != "コースA" {
			t.Errorf("コースA が返ること: got %s", body.Enrollments[0].Title)
		}
	})

	t.Run("コース内に複数進捗がある場合は最新進捗で並び順が決まる", func(t *testing.T) {
		cleanupMediumProgresses(ctx, pool)

		base := time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
		// コースA: A0 が古い、A1 が最新（A の MAX は base+2h）
		if err := insertProgress(ctx, pool, mediumTestUserID, mediumTestTopicA0ID, base); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}
		if err := insertProgress(ctx, pool, mediumTestUserID, mediumTestTopicA1ID, base.Add(2*time.Hour)); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}
		// コースB: B0 が中間（B の MAX は base+1h）
		if err := insertProgress(ctx, pool, mediumTestUserID, mediumTestTopicB0ID, base.Add(time.Hour)); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}

		w := httptest.NewRecorder()
		handler.GetEnrollmentsHandler(w, newGetEnrollmentsRequest(t, mediumTestUserID))

		var body GetEnrollmentsResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
		}
		// A の MAX(base+2h) > B の MAX(base+1h) なので A が先頭
		// （MIN 集約だと A の最古は base となり、B(base+1h) が先頭になる）
		expected := []string{"コースA", "コースB"}
		if len(body.Enrollments) != len(expected) {
			t.Fatalf("受講中コースが2件であること: got %d", len(body.Enrollments))
		}
		for i, want := range expected {
			if body.Enrollments[i].Title != want {
				t.Errorf("並び順 %d 番目が %s であること: got %s", i, want, body.Enrollments[i].Title)
			}
		}
	})

	t.Run("他ユーザーの進捗は含まれない", func(t *testing.T) {
		cleanupMediumProgresses(ctx, pool)

		now := time.Now()
		if err := insertProgress(ctx, pool, mediumTestOtherUser, mediumTestTopicA0ID, now); err != nil {
			t.Fatalf("progress 挿入失敗: %v", err)
		}

		w := httptest.NewRecorder()
		handler.GetEnrollmentsHandler(w, newGetEnrollmentsRequest(t, mediumTestUserID))

		var body GetEnrollmentsResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
		}
		if len(body.Enrollments) != 0 {
			t.Errorf("対象ユーザーの受講中コースが0件であること: got %d", len(body.Enrollments))
		}
	})
}
