package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/env"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	updTestCategoryID    = "31000000-0000-0000-0000-000000000000"
	updTestAuthorID      = "32000000-0000-0000-0000-000000000000"
	updTestUserID        = "33000000-0000-0000-0000-000000000000"
	updTestCourseID      = "34000000-0000-0000-0000-000000000000"
	updTestSection0ID    = "35000000-0000-0000-0000-000000000000"
	updTestTopic00ID     = "36000000-0000-0000-0000-000000000000"
	updTestTopic01ID     = "36000000-0000-0000-0000-000000000001"
	updTestOtherCourseID = "34000000-0000-0000-0000-000000000099"
	updTestOtherSection  = "35000000-0000-0000-0000-000000000099"
	updTestOtherTopic    = "36000000-0000-0000-0000-000000000099"
	updTestDraftCourseID = "34000000-0000-0000-0000-000000000010"
	updTestDraftSection0 = "35000000-0000-0000-0000-000000000010"
	updTestDraftTopic00  = "36000000-0000-0000-0000-000000000010"
	updTestMissingCourse = "34000000-0000-0000-0000-0000000000ff"
	updTestMissingTopic  = "36000000-0000-0000-0000-0000000000ff"
)

func setupUpdateEnrollmentTestData(ctx context.Context, pool *pgxpool.Pool) error {
	cleanupUpdateEnrollmentTestData(ctx, pool)

	sqls := []struct {
		query string
		args  []any
	}{
		{
			`INSERT INTO categories (category_id, name, path) VALUES ($1, 'テストカテゴリ', 'update-enrollment-test')`,
			[]any{updTestCategoryID},
		},
		{
			`INSERT INTO authors (author_id, name, profile) VALUES ($1, 'テスト著者', 'テストプロフィール')`,
			[]any{updTestAuthorID},
		},
		{
			`INSERT INTO users (user_id, nickname) VALUES ($1, 'テストユーザー')`,
			[]any{updTestUserID},
		},
		{
			`INSERT INTO users (user_id, nickname) VALUES ($1, 'テスト著者ユーザー')`,
			[]any{updTestAuthorID},
		},
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, 'テストコース', 'テスト説明', 'update-enrollment-course', '[]', 'published', $2, $3, 'public')`,
			[]any{updTestCourseID, updTestCategoryID, updTestAuthorID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ($1, $2, 0, 'セクション0', '')`,
			[]any{updTestSection0ID, updTestCourseID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, content)
			 VALUES ($1, $2, $3, 0, 'トピック0-0', '', '')`,
			[]any{updTestTopic00ID, updTestCourseID, updTestSection0ID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, content)
			 VALUES ($1, $2, $3, 1, 'トピック0-1', '', '')`,
			[]any{updTestTopic01ID, updTestCourseID, updTestSection0ID},
		},
		// 別コース（topic帰属違い検証用）
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, '別コース', '別説明', 'update-enrollment-other-course', '[]', 'published', $2, $3, 'public')`,
			[]any{updTestOtherCourseID, updTestCategoryID, updTestAuthorID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ($1, $2, 0, '別セクション0', '')`,
			[]any{updTestOtherSection, updTestOtherCourseID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, content)
			 VALUES ($1, $2, $3, 0, '別トピック0-0', '', '')`,
			[]any{updTestOtherTopic, updTestOtherCourseID, updTestOtherSection},
		},
		// ドラフトコース
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, 'ドラフトコース', 'ドラフト説明', 'update-enrollment-draft-course', '[]', 'draft', $2, $3, 'private')`,
			[]any{updTestDraftCourseID, updTestCategoryID, updTestAuthorID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ($1, $2, 0, 'ドラフトセクション0', '')`,
			[]any{updTestDraftSection0, updTestDraftCourseID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, content)
			 VALUES ($1, $2, $3, 0, 'ドラフトトピック0-0', '', '')`,
			[]any{updTestDraftTopic00, updTestDraftCourseID, updTestDraftSection0},
		},
	}

	for _, s := range sqls {
		if _, err := pool.Exec(ctx, s.query, s.args...); err != nil {
			return fmt.Errorf("テストデータの挿入に失敗しました: %w", err)
		}
	}
	return nil
}

func cleanupUpdateEnrollmentTestData(ctx context.Context, pool *pgxpool.Pool) {
	courseIDs := []string{updTestCourseID, updTestOtherCourseID, updTestDraftCourseID}
	userIDs := []string{updTestUserID, updTestAuthorID}
	_, _ = pool.Exec(ctx, `DELETE FROM topic_progresses WHERE user_id = ANY($1)`, userIDs)
	_, _ = pool.Exec(ctx, `DELETE FROM enrollments WHERE user_id = ANY($1)`, userIDs)
	_, _ = pool.Exec(ctx, `DELETE FROM course_section_topics WHERE course_id = ANY($1)`, courseIDs)
	_, _ = pool.Exec(ctx, `DELETE FROM course_sections WHERE course_id = ANY($1)`, courseIDs)
	_, _ = pool.Exec(ctx, `DELETE FROM courses WHERE course_id = ANY($1)`, courseIDs)
	_, _ = pool.Exec(ctx, `DELETE FROM users WHERE user_id = ANY($1)`, userIDs)
	_, _ = pool.Exec(ctx, `DELETE FROM authors WHERE author_id = $1`, updTestAuthorID)
	_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE category_id = $1`, updTestCategoryID)
}

func cleanupUpdateProgresses(ctx context.Context, pool *pgxpool.Pool) {
	users := []string{updTestUserID, updTestAuthorID}
	_, _ = pool.Exec(ctx, `DELETE FROM topic_progresses WHERE user_id = ANY($1)`, users)
	_, _ = pool.Exec(ctx, `DELETE FROM enrollments WHERE user_id = ANY($1)`, users)
}

// insertEnrollment は (user, course) の Enrollment を新規作成する。既存 enrollment は事前に削除しておくこと。
func insertEnrollment(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID, courseID string) {
	t.Helper()
	if _, err := pool.Exec(ctx,
		`INSERT INTO enrollments (user_id, course_id, enrolled_at) VALUES ($1, $2, $3)`,
		userID, courseID, time.Now(),
	); err != nil {
		t.Fatalf("enrollment 挿入失敗: %v", err)
	}
}

type updProgressRecord struct {
	topicID string
	status  string
}

func getUpdateProgresses(t *testing.T, ctx context.Context, pool *pgxpool.Pool, userID string) []updProgressRecord {
	t.Helper()
	rows, err := pool.Query(ctx,
		`SELECT course_section_topic_id, status
		 FROM topic_progresses
		 WHERE user_id = $1
		 ORDER BY course_section_topic_id`,
		userID,
	)
	if err != nil {
		t.Fatalf("progressの取得に失敗しました: %v", err)
	}
	defer rows.Close()

	var records []updProgressRecord
	for rows.Next() {
		var r updProgressRecord
		if err := rows.Scan(&r.topicID, &r.status); err != nil {
			t.Fatalf("progressのスキャンに失敗しました: %v", err)
		}
		records = append(records, r)
	}
	return records
}

func insertUpdateProgress(ctx context.Context, pool *pgxpool.Pool, userID, courseID, topicID, status string) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO topic_progresses (user_id, course_id, course_section_topic_id, status) VALUES ($1, $2, $3, $4)`,
		userID, courseID, topicID, status,
	)
	return err
}

func TestPatchEnrollmentHandlerMedium(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, env.Require("DATABASE_URL"))
	if err != nil {
		t.Fatalf("DBへの接続に失敗しました: %v", err)
	}
	defer pool.Close()

	if err := setupUpdateEnrollmentTestData(ctx, pool); err != nil {
		t.Fatalf("テストデータのセットアップに失敗しました: %v", err)
	}
	t.Cleanup(func() { cleanupUpdateEnrollmentTestData(ctx, pool) })

	q := db.New(pool)
	handler := NewUpdateEnrollmentHandler(NewUpdateEnrollmentUsecase(NewSqrcCourseRepository(q), NewSqrcEnrollmentRepository(q)))

	newReq := func(t *testing.T, userID, courseID, body string) *http.Request {
		t.Helper()
		req := httptest.NewRequest("PATCH", "/v1/users/"+userID+"/enrollments/"+courseID, strings.NewReader(body))
		req.SetPathValue("userID", userID)
		req.SetPathValue("courseId", courseID)
		return req
	}

	t.Run("受講中コースの未進捗トピックをIN_PROGRESSで開始できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupUpdateProgresses(ctx, pool) })
		insertEnrollment(t, ctx, pool, updTestUserID, updTestCourseID)

		w := httptest.NewRecorder()
		handler.PatchEnrollmentHandler(w, newReq(t, updTestUserID, updTestCourseID, `{"topicId":"`+updTestTopic00ID+`","status":"IN_PROGRESS"}`))

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != updTestTopic00ID || body["status"] != "IN_PROGRESS" {
			t.Errorf("レスポンス不一致: got %+v", body)
		}

		progresses := getUpdateProgresses(t, ctx, pool, updTestUserID)
		if len(progresses) != 1 {
			t.Fatalf("progressが1件であること: got %d", len(progresses))
		}
		if progresses[0].topicID != updTestTopic00ID || progresses[0].status != "IN_PROGRESS" {
			t.Errorf("progress不一致: got %+v", progresses[0])
		}
	})

	t.Run("受講中コースの未進捗トピックをCOMPLETEDで保存できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupUpdateProgresses(ctx, pool) })
		insertEnrollment(t, ctx, pool, updTestUserID, updTestCourseID)

		w := httptest.NewRecorder()
		handler.PatchEnrollmentHandler(w, newReq(t, updTestUserID, updTestCourseID, `{"topicId":"`+updTestTopic00ID+`","status":"COMPLETED"}`))

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", w.Result().StatusCode)
		}

		progresses := getUpdateProgresses(t, ctx, pool, updTestUserID)
		if len(progresses) != 1 {
			t.Fatalf("progressが1件であること: got %d", len(progresses))
		}
		if progresses[0].status != "COMPLETED" {
			t.Errorf("statusがCOMPLETEDであること: got %q", progresses[0].status)
		}
	})

	t.Run("既存IN_PROGRESSをCOMPLETEDに遷移できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupUpdateProgresses(ctx, pool) })
		insertEnrollment(t, ctx, pool, updTestUserID, updTestCourseID)
		if err := insertUpdateProgress(ctx, pool, updTestUserID, updTestCourseID, updTestTopic00ID, "IN_PROGRESS"); err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}

		w := httptest.NewRecorder()
		handler.PatchEnrollmentHandler(w, newReq(t, updTestUserID, updTestCourseID, `{"topicId":"`+updTestTopic00ID+`","status":"COMPLETED"}`))

		if w.Result().StatusCode != http.StatusOK {
			t.Fatalf("ステータスコードが200であること: got %d", w.Result().StatusCode)
		}

		progresses := getUpdateProgresses(t, ctx, pool, updTestUserID)
		if len(progresses) != 1 || progresses[0].status != "COMPLETED" {
			t.Errorf("progressがCOMPLETEDに更新されていること: got %+v", progresses)
		}
	})

	t.Run("COMPLETEDからIN_PROGRESSへの巻き戻しは400を返す", func(t *testing.T) {
		t.Cleanup(func() { cleanupUpdateProgresses(ctx, pool) })
		insertEnrollment(t, ctx, pool, updTestUserID, updTestCourseID)
		if err := insertUpdateProgress(ctx, pool, updTestUserID, updTestCourseID, updTestTopic00ID, "COMPLETED"); err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}

		w := httptest.NewRecorder()
		handler.PatchEnrollmentHandler(w, newReq(t, updTestUserID, updTestCourseID, `{"topicId":"`+updTestTopic00ID+`","status":"IN_PROGRESS"}`))

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("ステータスコードが400であること: got %d", w.Result().StatusCode)
		}

		progresses := getUpdateProgresses(t, ctx, pool, updTestUserID)
		if len(progresses) != 1 || progresses[0].status != "COMPLETED" {
			t.Errorf("DB上の status は COMPLETED のままであること: got %+v", progresses)
		}
	})

	t.Run("topicIdが別コースの場合404 TOPIC_NOT_FOUNDを返しDBは変化しない", func(t *testing.T) {
		t.Cleanup(func() { cleanupUpdateProgresses(ctx, pool) })
		insertEnrollment(t, ctx, pool, updTestUserID, updTestCourseID)

		w := httptest.NewRecorder()
		handler.PatchEnrollmentHandler(w, newReq(t, updTestUserID, updTestCourseID, `{"topicId":"`+updTestOtherTopic+`","status":"IN_PROGRESS"}`))

		if w.Result().StatusCode != http.StatusNotFound {
			t.Errorf("ステータスコードが404であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["errorCode"] != "TOPIC_NOT_FOUND" {
			t.Errorf("errorCodeが一致すること: got %q", body["errorCode"])
		}

		progresses := getUpdateProgresses(t, ctx, pool, updTestUserID)
		if len(progresses) != 0 {
			t.Errorf("DBに変更がないこと: got %d 件", len(progresses))
		}
	})

	t.Run("未受講のコースに対するPATCHは404 NOT_ENROLLEDを返す", func(t *testing.T) {
		t.Cleanup(func() { cleanupUpdateProgresses(ctx, pool) })

		w := httptest.NewRecorder()
		handler.PatchEnrollmentHandler(w, newReq(t, updTestUserID, updTestCourseID, `{"topicId":"`+updTestTopic00ID+`","status":"IN_PROGRESS"}`))

		if w.Result().StatusCode != http.StatusNotFound {
			t.Errorf("ステータスコードが404であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["errorCode"] != "NOT_ENROLLED" {
			t.Errorf("errorCodeがNOT_ENROLLEDであること: got %q", body["errorCode"])
		}
	})

	t.Run("存在しないcourseIdの場合404 COURSE_NOT_FOUNDを返す", func(t *testing.T) {
		w := httptest.NewRecorder()
		handler.PatchEnrollmentHandler(w, newReq(t, updTestUserID, updTestMissingCourse, `{"topicId":"`+updTestTopic00ID+`","status":"IN_PROGRESS"}`))

		if w.Result().StatusCode != http.StatusNotFound {
			t.Errorf("ステータスコードが404であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["errorCode"] != "COURSE_NOT_FOUND" {
			t.Errorf("errorCodeが一致すること: got %q", body["errorCode"])
		}
	})

	t.Run("コース配下に存在しないtopicIdの場合404 TOPIC_NOT_FOUNDを返す", func(t *testing.T) {
		t.Cleanup(func() { cleanupUpdateProgresses(ctx, pool) })
		insertEnrollment(t, ctx, pool, updTestUserID, updTestCourseID)

		w := httptest.NewRecorder()
		handler.PatchEnrollmentHandler(w, newReq(t, updTestUserID, updTestCourseID, `{"topicId":"`+updTestMissingTopic+`","status":"IN_PROGRESS"}`))

		if w.Result().StatusCode != http.StatusNotFound {
			t.Errorf("ステータスコードが404であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["errorCode"] != "TOPIC_NOT_FOUND" {
			t.Errorf("errorCodeが一致すること: got %q", body["errorCode"])
		}
	})
}
