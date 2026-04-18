package queries

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	testCategoryID = "10000000-0000-0000-0000-000000000000"
	testAuthorID   = "20000000-0000-0000-0000-000000000000"
	titleCourseID  = "30000000-0000-0000-0000-000000000001"
	descCourseID   = "30000000-0000-0000-0000-000000000002"
)

func databaseURL() string {
	return env.Require("DATABASE_URL")
}

func setupTestData(ctx context.Context, pool *pgxpool.Pool) error {
	cleanupTestData(ctx, pool)

	_, err := pool.Exec(ctx,
		`INSERT INTO categories (category_id, name, path) VALUES ($1, 'テストカテゴリ', 'test')`,
		testCategoryID,
	)
	if err != nil {
		return fmt.Errorf("カテゴリの挿入に失敗しました: %w", err)
	}

	_, err = pool.Exec(ctx,
		`INSERT INTO authors (author_id, name, profile) VALUES ($1, 'テスト著者', 'テストプロフィール')`,
		testAuthorID,
	)
	if err != nil {
		return fmt.Errorf("著者の挿入に失敗しました: %w", err)
	}

	_, err = pool.Exec(ctx,
		`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
		 VALUES ($1, 'タイトル一致テスト専用コース', '一般説明', 'test-title-course', '[]', 'published', $2, $3, 'public')`,
		titleCourseID, testCategoryID, testAuthorID,
	)
	if err != nil {
		return fmt.Errorf("タイトル一致コースの挿入に失敗しました: %w", err)
	}

	_, err = pool.Exec(ctx,
		`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
		 VALUES ($1, '通常コースB', '説明文一致テスト専用テキスト', 'test-desc-course', '[]', 'published', $2, $3, 'public')`,
		descCourseID, testCategoryID, testAuthorID,
	)
	if err != nil {
		return fmt.Errorf("説明一致コースの挿入に失敗しました: %w", err)
	}

	for i := 1; i <= 201; i++ {
		courseID := fmt.Sprintf("40000000-0000-0000-0000-%012d", i)
		slug := fmt.Sprintf("test-page-%d", i)
		_, err = pool.Exec(ctx,
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, 'ページネーションテスト', 'ページコース説明', $2, '[]', 'published', $3, $4, 'public')`,
			courseID, slug, testCategoryID, testAuthorID,
		)
		if err != nil {
			return fmt.Errorf("ページネーションコース %d の挿入に失敗しました: %w", i, err)
		}
	}

	return nil
}

func cleanupTestData(ctx context.Context, pool *pgxpool.Pool) {
	_, _ = pool.Exec(ctx, `DELETE FROM courses WHERE author_id = $1`, testAuthorID)
	_, _ = pool.Exec(ctx, `DELETE FROM authors WHERE author_id = $1`, testAuthorID)
	_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE category_id = $1`, testCategoryID)
}

func TestGetCoursesHandlerMedium(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, databaseURL())
	if err != nil {
		t.Fatalf("DBへの接続に失敗しました: %v", err)
	}
	defer pool.Close()

	if err := setupTestData(ctx, pool); err != nil {
		t.Fatalf("テストデータのセットアップに失敗しました: %v", err)
	}
	t.Cleanup(func() { cleanupTestData(ctx, pool) })

	handlers := NewCoursesHandlers(db.New(pool))

	t.Run("条件を指定せずに一覧を取得できる", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v1/courses", nil)
		w := httptest.NewRecorder()

		handlers.GetCoursesHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", res.StatusCode)
		}
		var body GetCoursesResult
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
		}
		if len(body.Courses) == 0 {
			t.Error("コースが1件以上返ること")
		}
	})

	t.Run("キーワードでタイトルを検索できる", func(t *testing.T) {
		params := url.Values{}
		params.Set("keyword", "タイトル一致テスト")
		req := httptest.NewRequest("GET", "/v1/courses?"+params.Encode(), nil)
		w := httptest.NewRecorder()

		handlers.GetCoursesHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", res.StatusCode)
		}
		var body GetCoursesResult
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
		}
		if len(body.Courses) != 1 {
			t.Fatalf("コースが1件返ること: got %d", len(body.Courses))
		}
		if body.Courses[0].Title != "タイトル一致テスト専用コース" {
			t.Errorf("タイトルが一致すること: got %s", body.Courses[0].Title)
		}
	})

	t.Run("キーワードで説明を検索できる", func(t *testing.T) {
		params := url.Values{}
		params.Set("keyword", "説明文一致テスト")
		req := httptest.NewRequest("GET", "/v1/courses?"+params.Encode(), nil)
		w := httptest.NewRecorder()

		handlers.GetCoursesHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", res.StatusCode)
		}
		var body GetCoursesResult
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
		}
		if len(body.Courses) != 1 {
			t.Fatalf("コースが1件返ること: got %d", len(body.Courses))
		}
		if body.Courses[0].Title != "通常コースB" {
			t.Errorf("タイトルが一致すること: got %s", body.Courses[0].Title)
		}
	})

	t.Run("次のページがある場合200件とカーソルが返る", func(t *testing.T) {
		params := url.Values{}
		params.Set("keyword", "ページネーションテスト")
		req := httptest.NewRequest("GET", "/v1/courses?"+params.Encode(), nil)
		w := httptest.NewRecorder()

		handlers.GetCoursesHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", res.StatusCode)
		}
		var body GetCoursesResult
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
		}
		if body.Cursor == nil {
			t.Error("カーソルが返ること")
		}
		if len(body.Courses) != 200 {
			t.Errorf("コースが200件返ること: got %d", len(body.Courses))
		}
	})

	t.Run("カーソルを使って2ページ目が取得できる", func(t *testing.T) {
		// 1ページ目を取得
		params := url.Values{}
		params.Set("keyword", "ページネーションテスト")
		req := httptest.NewRequest("GET", "/v1/courses?"+params.Encode(), nil)
		w := httptest.NewRecorder()
		handlers.GetCoursesHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("1ページ目のステータスコードが200であること: got %d", res.StatusCode)
		}
		var firstPage GetCoursesResult
		if err := json.NewDecoder(res.Body).Decode(&firstPage); err != nil {
			t.Fatalf("1ページ目のレスポンスボディのデコードに失敗しました: %v", err)
		}
		if firstPage.Cursor == nil {
			t.Fatal("1ページ目のカーソルが返ること")
		}

		// 2ページ目を取得
		params2 := url.Values{}
		params2.Set("keyword", "ページネーションテスト")
		params2.Set("cursor", *firstPage.Cursor)
		req2 := httptest.NewRequest("GET", "/v1/courses?"+params2.Encode(), nil)
		w2 := httptest.NewRecorder()
		handlers.GetCoursesHandler(w2, req2)

		res2 := w2.Result()
		if res2.StatusCode != http.StatusOK {
			t.Fatalf("2ページ目のステータスコードが200であること: got %d", res2.StatusCode)
		}
		var secondPage GetCoursesResult
		if err := json.NewDecoder(res2.Body).Decode(&secondPage); err != nil {
			t.Fatalf("2ページ目のレスポンスボディのデコードに失敗しました: %v", err)
		}
		if len(secondPage.Courses) != 1 {
			t.Fatalf("2ページ目のコースが1件返ること: got %d", len(secondPage.Courses))
		}

		// 1ページ目と重複しないことを確認
		firstPageIDs := make(map[string]bool)
		for _, course := range firstPage.Courses {
			firstPageIDs[course.CourseId.String()] = true
		}
		for _, course := range secondPage.Courses {
			if firstPageIDs[course.CourseId.String()] {
				t.Errorf("2ページ目のコースが1ページ目と重複している: %s", course.CourseId)
			}
		}
	})

	t.Run("次のページがない場合カーソルが null", func(t *testing.T) {
		params := url.Values{}
		params.Set("keyword", "タイトル一致テスト")
		req := httptest.NewRequest("GET", "/v1/courses?"+params.Encode(), nil)
		w := httptest.NewRecorder()

		handlers.GetCoursesHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", res.StatusCode)
		}
		var body GetCoursesResult
		if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
		}
		if body.Cursor != nil {
			t.Errorf("カーソルが null であること: got %s", *body.Cursor)
		}
		if len(body.Courses) >= 200 {
			t.Errorf("コースが200件未満であること: got %d", len(body.Courses))
		}
	})
}
