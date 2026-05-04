package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	libauth "github.com/harusame0616/ijuku/apps/api/lib/auth"
	"github.com/harusame0616/ijuku/apps/api/lib/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	cdMedCategoryID    = "51000000-0000-0000-0000-000000000000"
	cdMedAuthorID      = "52000000-0000-0000-0000-000000000000"
	cdMedUserID        = "53000000-0000-0000-0000-000000000000"
	cdMedPubCourseID   = "54000000-0000-0000-0000-000000000001"
	cdMedDraftCourseID = "54000000-0000-0000-0000-000000000002"
	cdMedSectionID     = "55000000-0000-0000-0000-000000000001"
	cdMedTopicID       = "56000000-0000-0000-0000-000000000001"

	cdMedAuthorSlug      = "course-detail-author"
	cdMedPubCourseSlug   = "course-detail-public"
	cdMedDraftCourseSlug = "course-detail-draft"
)

func setupCourseDetailTestData(ctx context.Context, pool *pgxpool.Pool) error {
	cleanupCourseDetailTestData(ctx, pool)

	sqls := []struct {
		query string
		args  []any
	}{
		{
			`INSERT INTO categories (category_id, name, path) VALUES ($1, '講座詳細カテゴリ', 'course-detail')`,
			[]any{cdMedCategoryID},
		},
		{
			`INSERT INTO authors (author_id, name, profile, slug) VALUES ($1, '講座詳細著者', '著者プロフィール', $2)`,
			[]any{cdMedAuthorID, cdMedAuthorSlug},
		},
		{
			`INSERT INTO users (user_id, nickname) VALUES ($1, '受講ユーザー')`,
			[]any{cdMedUserID},
		},
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, '公開講座タイトル', '公開講座の説明', $2, '["go","backend"]', 'published', $3, $4, 'public')`,
			[]any{cdMedPubCourseID, cdMedPubCourseSlug, cdMedCategoryID, cdMedAuthorID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description)
			 VALUES ($1, $2, 0, 'セクション1', 'セクション1の説明')`,
			[]any{cdMedSectionID, cdMedPubCourseID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, content)
			 VALUES ($1, $2, $3, 0, 'トピック1', 'トピック1の説明', '本文')`,
			[]any{cdMedTopicID, cdMedPubCourseID, cdMedSectionID},
		},
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, '下書き講座', '', $2, '[]', 'draft', $3, $4, 'private')`,
			[]any{cdMedDraftCourseID, cdMedDraftCourseSlug, cdMedCategoryID, cdMedAuthorID},
		},
	}

	for _, s := range sqls {
		if _, err := pool.Exec(ctx, s.query, s.args...); err != nil {
			return fmt.Errorf("テストデータの挿入に失敗しました: %w", err)
		}
	}
	return nil
}

func cleanupCourseDetailTestData(ctx context.Context, pool *pgxpool.Pool) {
	courses := []string{cdMedPubCourseID, cdMedDraftCourseID}
	users := []string{cdMedUserID}

	_, _ = pool.Exec(ctx, `DELETE FROM topic_progresses WHERE user_id = ANY($1)`, users)
	_, _ = pool.Exec(ctx, `DELETE FROM enrollments WHERE user_id = ANY($1)`, users)
	_, _ = pool.Exec(ctx, `DELETE FROM course_section_topics WHERE course_id = ANY($1)`, courses)
	_, _ = pool.Exec(ctx, `DELETE FROM course_sections WHERE course_id = ANY($1)`, courses)
	_, _ = pool.Exec(ctx, `DELETE FROM courses WHERE course_id = ANY($1)`, courses)
	_, _ = pool.Exec(ctx, `DELETE FROM users WHERE user_id = ANY($1)`, users)
	_, _ = pool.Exec(ctx, `DELETE FROM authors WHERE author_id = $1`, cdMedAuthorID)
	_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE category_id = $1`, cdMedCategoryID)
}

func cleanupCourseDetailEnrollments(ctx context.Context, pool *pgxpool.Pool) {
	_, _ = pool.Exec(ctx, `DELETE FROM enrollments WHERE user_id = $1`, cdMedUserID)
}

func TestGetCourseDetailHandlerMedium(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, env.Require("DATABASE_URL"))
	if err != nil {
		t.Fatalf("DBへの接続に失敗しました: %v", err)
	}
	defer pool.Close()

	if err := setupCourseDetailTestData(ctx, pool); err != nil {
		t.Fatalf("テストデータのセットアップに失敗しました: %v", err)
	}
	t.Cleanup(func() { cleanupCourseDetailTestData(ctx, pool) })

	handler := NewGetCourseDetailHandler(db.New(pool))

	newReq := func(authorSlug, courseSlug string, userID string) *http.Request {
		req := httptest.NewRequest("GET", fmt.Sprintf("/v1/courses/%s/%s", authorSlug, courseSlug), nil)
		req.SetPathValue("authorSlug", authorSlug)
		req.SetPathValue("courseSlug", courseSlug)
		if userID != "" {
			req = req.WithContext(libauth.WithUserID(req.Context(), userID))
		}
		return req
	}

	t.Run("公開講座は未ログインでも取得できisEnrolledはfalse", func(t *testing.T) {
		w := httptest.NewRecorder()
		handler.GetCourseDetailHandler(w, newReq(cdMedAuthorSlug, cdMedPubCourseSlug, ""))

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", w.Result().StatusCode)
		}
		var resp GetCourseDetailResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&resp); err != nil {
			t.Fatalf("レスポンスのデコードに失敗: %v", err)
		}
		if resp.CourseId != cdMedPubCourseID {
			t.Errorf("courseIdが一致すること: got %q", resp.CourseId)
		}
		if resp.Title != "公開講座タイトル" {
			t.Errorf("titleが一致すること: got %q", resp.Title)
		}
		if resp.Slug != cdMedPubCourseSlug {
			t.Errorf("slugが一致すること: got %q", resp.Slug)
		}
		if resp.Author.Slug != cdMedAuthorSlug {
			t.Errorf("author.slugが一致すること: got %q", resp.Author.Slug)
		}
		if len(resp.Tags) != 2 || resp.Tags[0] != "go" || resp.Tags[1] != "backend" {
			t.Errorf("tagsが一致すること: got %v", resp.Tags)
		}
		if len(resp.Sections) != 1 {
			t.Fatalf("セクションが1件返ること: got %d", len(resp.Sections))
		}
		if len(resp.Sections[0].Topics) != 1 {
			t.Fatalf("トピックが1件返ること: got %d", len(resp.Sections[0].Topics))
		}
		if resp.Sections[0].Topics[0].Title != "トピック1" {
			t.Errorf("トピックタイトルが一致すること: got %q", resp.Sections[0].Topics[0].Title)
		}
		if resp.IsEnrolled {
			t.Errorf("未ログインのときisEnrolledはfalseであること")
		}
	})

	t.Run("受講中ユーザーで取得するとisEnrolledがtrueになる", func(t *testing.T) {
		t.Cleanup(func() { cleanupCourseDetailEnrollments(ctx, pool) })

		if _, err := pool.Exec(ctx,
			`INSERT INTO enrollments (user_id, course_id, enrolled_at) VALUES ($1, $2, NOW())`,
			cdMedUserID, cdMedPubCourseID,
		); err != nil {
			t.Fatalf("受講レコード挿入失敗: %v", err)
		}

		w := httptest.NewRecorder()
		handler.GetCourseDetailHandler(w, newReq(cdMedAuthorSlug, cdMedPubCourseSlug, cdMedUserID))

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", w.Result().StatusCode)
		}
		var resp GetCourseDetailResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&resp); err != nil {
			t.Fatalf("レスポンスのデコードに失敗: %v", err)
		}
		if !resp.IsEnrolled {
			t.Errorf("受講中ならisEnrolledがtrueであること")
		}
	})

	t.Run("ログインユーザーで未受講の場合isEnrolledはfalse", func(t *testing.T) {
		w := httptest.NewRecorder()
		handler.GetCourseDetailHandler(w, newReq(cdMedAuthorSlug, cdMedPubCourseSlug, cdMedUserID))

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", w.Result().StatusCode)
		}
		var resp GetCourseDetailResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&resp); err != nil {
			t.Fatalf("レスポンスのデコードに失敗: %v", err)
		}
		if resp.IsEnrolled {
			t.Errorf("未受講ならisEnrolledがfalseであること")
		}
	})

	t.Run("draft講座は404を返す", func(t *testing.T) {
		w := httptest.NewRecorder()
		handler.GetCourseDetailHandler(w, newReq(cdMedAuthorSlug, cdMedDraftCourseSlug, cdMedUserID))

		if w.Result().StatusCode != http.StatusNotFound {
			t.Errorf("ステータスコードが404であること: got %d", w.Result().StatusCode)
		}
		var resp map[string]string
		json.NewDecoder(w.Result().Body).Decode(&resp)
		if resp["errorCode"] != "COURSE_NOT_FOUND" {
			t.Errorf("errorCodeがCOURSE_NOT_FOUNDであること: got %q", resp["errorCode"])
		}
	})

	t.Run("存在しないslugは404を返す", func(t *testing.T) {
		w := httptest.NewRecorder()
		handler.GetCourseDetailHandler(w, newReq(cdMedAuthorSlug, "missing-course", ""))

		if w.Result().StatusCode != http.StatusNotFound {
			t.Errorf("ステータスコードが404であること: got %d", w.Result().StatusCode)
		}
	})
}
