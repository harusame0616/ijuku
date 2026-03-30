package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	enrollTestCategoryID    = "11000000-0000-0000-0000-000000000000"
	enrollTestAuthorID      = "12000000-0000-0000-0000-000000000000"
	enrollTestUserID        = "13000000-0000-0000-0000-000000000000"
	enrollTestCourseID      = "14000000-0000-0000-0000-000000000000"
	enrollTestSection0ID    = "15000000-0000-0000-0000-000000000000" // index=0
	enrollTestSection1ID    = "15000000-0000-0000-0000-000000000001" // index=1
	enrollTestTopic00ID     = "16000000-0000-0000-0000-000000000000" // section0 index=0
	enrollTestTopic01ID     = "16000000-0000-0000-0000-000000000001" // section0 index=1
	enrollTestTopic10ID     = "16000000-0000-0000-0000-000000000010" // section1 index=0
	enrollTestDraftCourseID = "14000000-0000-0000-0000-000000000001"
	enrollTestDraftSection0 = "15000000-0000-0000-0000-000000000002"
	enrollTestDraftTopic00  = "16000000-0000-0000-0000-000000000020"
)

func enrollDatabaseURL() string {
	if dsn := os.Getenv("DATABASE_URL"); dsn != "" {
		return dsn
	}
	return "postgresql://postgres:password@localhost:5432/postgres"
}

func setupEnrollTestData(ctx context.Context, pool *pgxpool.Pool) error {
	cleanupEnrollTestData(ctx, pool)

	sqls := []struct {
		query string
		args  []any
	}{
		{
			`INSERT INTO categories (category_id, name, path) VALUES ($1, 'テストカテゴリ', 'enroll-test')`,
			[]any{enrollTestCategoryID},
		},
		{
			`INSERT INTO authors (author_id, name, profile) VALUES ($1, 'テスト著者', 'テストプロフィール')`,
			[]any{enrollTestAuthorID},
		},
		{
			`INSERT INTO users (user_id, nickname) VALUES ($1, 'テストユーザー')`,
			[]any{enrollTestUserID},
		},
		{
			`INSERT INTO users (user_id, nickname) VALUES ($1, 'テスト著者ユーザー')`,
			[]any{enrollTestAuthorID},
		},
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, 'テストコース', 'テスト説明', 'enroll-test-course', '[]', 'published', $2, $3, 'public')`,
			[]any{enrollTestCourseID, enrollTestCategoryID, enrollTestAuthorID},
		},
		{
			`INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, author_id, visibility)
			 VALUES ($1, 'ドラフトコース', 'ドラフト説明', 'enroll-test-draft-course', '[]', 'draft', $2, $3, 'private')`,
			[]any{enrollTestDraftCourseID, enrollTestCategoryID, enrollTestAuthorID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ($1, $2, 0, 'ドラフトセクション0', '')`,
			[]any{enrollTestDraftSection0, enrollTestDraftCourseID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, prerequisites, knowledge, flow, quiz, completion_criteria)
			 VALUES ($1, $2, $3, 0, 'ドラフトトピック0-0', '', '', '', '', '', '')`,
			[]any{enrollTestDraftTopic00, enrollTestDraftCourseID, enrollTestDraftSection0},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ($1, $2, 0, 'セクション0', '')`,
			[]any{enrollTestSection0ID, enrollTestCourseID},
		},
		{
			`INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ($1, $2, 1, 'セクション1', '')`,
			[]any{enrollTestSection1ID, enrollTestCourseID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, prerequisites, knowledge, flow, quiz, completion_criteria)
			 VALUES ($1, $2, $3, 0, 'トピック0-0', '', '', '', '', '', '')`,
			[]any{enrollTestTopic00ID, enrollTestCourseID, enrollTestSection0ID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, prerequisites, knowledge, flow, quiz, completion_criteria)
			 VALUES ($1, $2, $3, 1, 'トピック0-1', '', '', '', '', '', '')`,
			[]any{enrollTestTopic01ID, enrollTestCourseID, enrollTestSection0ID},
		},
		{
			`INSERT INTO course_section_topics (course_section_topic_id, course_id, course_section_id, index, title, description, prerequisites, knowledge, flow, quiz, completion_criteria)
			 VALUES ($1, $2, $3, 0, 'トピック1-0', '', '', '', '', '', '')`,
			[]any{enrollTestTopic10ID, enrollTestCourseID, enrollTestSection1ID},
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
	_, _ = pool.Exec(ctx, `DELETE FROM user_topic_progresses WHERE user_id = ANY($1)`, []string{enrollTestUserID, enrollTestAuthorID})
	_, _ = pool.Exec(ctx, `DELETE FROM course_section_topics WHERE course_id = ANY($1)`, []string{enrollTestCourseID, enrollTestDraftCourseID})
	_, _ = pool.Exec(ctx, `DELETE FROM course_sections WHERE course_id = ANY($1)`, []string{enrollTestCourseID, enrollTestDraftCourseID})
	_, _ = pool.Exec(ctx, `DELETE FROM courses WHERE course_id = ANY($1)`, []string{enrollTestCourseID, enrollTestDraftCourseID})
	_, _ = pool.Exec(ctx, `DELETE FROM users WHERE user_id = ANY($1)`, []string{enrollTestUserID, enrollTestAuthorID})
	_, _ = pool.Exec(ctx, `DELETE FROM authors WHERE author_id = $1`, enrollTestAuthorID)
	_, _ = pool.Exec(ctx, `DELETE FROM categories WHERE category_id = $1`, enrollTestCategoryID)
}

func cleanupProgresses(ctx context.Context, pool *pgxpool.Pool) {
	_, _ = pool.Exec(ctx, `DELETE FROM user_topic_progresses WHERE user_id = ANY($1)`, []string{enrollTestUserID, enrollTestAuthorID})
}

type progressRecord struct {
	topicID string
	status  string
}

func getProgresses(t *testing.T, ctx context.Context, pool *pgxpool.Pool) []progressRecord {
	t.Helper()
	rows, err := pool.Query(ctx,
		`SELECT course_section_topic_id, status FROM user_topic_progresses WHERE user_id = $1 ORDER BY course_section_topic_id`,
		enrollTestUserID,
	)
	if err != nil {
		t.Fatalf("progressの取得に失敗しました: %v", err)
	}
	defer rows.Close()

	var records []progressRecord
	for rows.Next() {
		var r progressRecord
		if err := rows.Scan(&r.topicID, &r.status); err != nil {
			t.Fatalf("progressのスキャンに失敗しました: %v", err)
		}
		records = append(records, r)
	}
	return records
}

func insertProgress(ctx context.Context, pool *pgxpool.Pool, topicID, status string) error {
	_, err := pool.Exec(ctx,
		`INSERT INTO user_topic_progresses (course_section_topic_id, user_id, status)
		 VALUES ($1, $2, $3)`,
		topicID, enrollTestUserID, status,
	)
	return err
}

func TestPostEnrollmentHandlerMedium(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, enrollDatabaseURL())
	if err != nil {
		t.Fatalf("DBへの接続に失敗しました: %v", err)
	}
	defer pool.Close()

	if err := setupEnrollTestData(ctx, pool); err != nil {
		t.Fatalf("テストデータのセットアップに失敗しました: %v", err)
	}
	t.Cleanup(func() { cleanupEnrollTestData(ctx, pool) })

	q := db.New(pool)
	handler := NewHandler(NewEnrollCourseUsecase(NewSqrcCourseRepository(q), NewSqrcUserTopicProgressRepository(q)))

	newEnrollRequest := func(t *testing.T, body string) *http.Request {
		t.Helper()
		req := httptest.NewRequest("POST", "/v1/courses/"+enrollTestCourseID+"/enrollment", strings.NewReader(body))
		req.SetPathValue("courseId", enrollTestCourseID)
		return req
	}

	t.Run("sectionNumberを指定するとその最初のtopicを開始できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupProgresses(ctx, pool) })

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`","sectionNumber":0}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != enrollTestTopic00ID {
			t.Errorf("topicIdが一致すること: got %q, want %q", body["topicId"], enrollTestTopic00ID)
		}

		progresses := getProgresses(t, ctx, pool)
		if len(progresses) != 1 {
			t.Fatalf("progressが1件であること: got %d", len(progresses))
		}
		if progresses[0].topicID != enrollTestTopic00ID || progresses[0].status != "IN_PROGRESS" {
			t.Errorf("progress不一致: got topicID=%q status=%q", progresses[0].topicID, progresses[0].status)
		}
	})

	t.Run("IN_PROGRESSのtopicがあってもsectionNumberを指定するとその最初のtopicを開始できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupProgresses(ctx, pool) })
		if err := insertProgress(ctx, pool, enrollTestTopic00ID, "IN_PROGRESS"); err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`","sectionNumber":0}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != enrollTestTopic00ID {
			t.Errorf("topicIdが一致すること: got %q", body["topicId"])
		}

		progresses := getProgresses(t, ctx, pool)
		if len(progresses) != 1 {
			t.Fatalf("progressが1件であること: got %d", len(progresses))
		}
		if progresses[0].topicID != enrollTestTopic00ID || progresses[0].status != "IN_PROGRESS" {
			t.Errorf("progress不一致: got topicID=%q status=%q", progresses[0].topicID, progresses[0].status)
		}
	})

	t.Run("COMPLETEDのtopicがあってもsectionNumberを指定するとその最初のtopicを開始でき、ステータスは変わらない", func(t *testing.T) {
		t.Cleanup(func() { cleanupProgresses(ctx, pool) })
		if err := insertProgress(ctx, pool, enrollTestTopic00ID, "COMPLETED"); err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`","sectionNumber":0}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != enrollTestTopic00ID {
			t.Errorf("topicIdが一致すること: got %q", body["topicId"])
		}

		progresses := getProgresses(t, ctx, pool)
		if len(progresses) != 1 {
			t.Fatalf("progressが1件であること: got %d", len(progresses))
		}
		if progresses[0].topicID != enrollTestTopic00ID || progresses[0].status != "COMPLETED" {
			t.Errorf("progress不一致: got topicID=%q status=%q", progresses[0].topicID, progresses[0].status)
		}
	})

	t.Run("sectionNumberとtopicNumberを指定するとそのtopicを開始できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupProgresses(ctx, pool) })

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`","sectionNumber":0,"topicNumber":1}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != enrollTestTopic01ID {
			t.Errorf("topicIdが一致すること: got %q, want %q", body["topicId"], enrollTestTopic01ID)
		}

		progresses := getProgresses(t, ctx, pool)
		if len(progresses) != 1 {
			t.Fatalf("progressが1件であること: got %d", len(progresses))
		}
		if progresses[0].topicID != enrollTestTopic01ID || progresses[0].status != "IN_PROGRESS" {
			t.Errorf("progress不一致: got topicID=%q status=%q", progresses[0].topicID, progresses[0].status)
		}
	})

	t.Run("IN_PROGRESSのtopicがあってもsectionNumber+topicNumberを指定するとそのtopicを開始できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupProgresses(ctx, pool) })
		if err := insertProgress(ctx, pool, enrollTestTopic01ID, "IN_PROGRESS"); err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`","sectionNumber":0,"topicNumber":1}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != enrollTestTopic01ID {
			t.Errorf("topicIdが一致すること: got %q", body["topicId"])
		}

		progresses := getProgresses(t, ctx, pool)
		if len(progresses) != 1 {
			t.Fatalf("progressが1件であること: got %d", len(progresses))
		}
		if progresses[0].topicID != enrollTestTopic01ID || progresses[0].status != "IN_PROGRESS" {
			t.Errorf("progress不一致: got topicID=%q status=%q", progresses[0].topicID, progresses[0].status)
		}
	})

	t.Run("COMPLETEDのtopicがあってもsectionNumber+topicNumberを指定するとそのtopicを開始でき、ステータスは変わらない", func(t *testing.T) {
		t.Cleanup(func() { cleanupProgresses(ctx, pool) })
		if err := insertProgress(ctx, pool, enrollTestTopic01ID, "COMPLETED"); err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`","sectionNumber":0,"topicNumber":1}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != enrollTestTopic01ID {
			t.Errorf("topicIdが一致すること: got %q", body["topicId"])
		}

		progresses := getProgresses(t, ctx, pool)
		if len(progresses) != 1 {
			t.Fatalf("progressが1件であること: got %d", len(progresses))
		}
		if progresses[0].topicID != enrollTestTopic01ID || progresses[0].status != "COMPLETED" {
			t.Errorf("progress不一致: got topicID=%q status=%q", progresses[0].topicID, progresses[0].status)
		}
	})

	t.Run("未受講のコースでsection/topicを省略すると最初のtopicを開始できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupProgresses(ctx, pool) })

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`"}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != enrollTestTopic00ID {
			t.Errorf("topicIdが一致すること: got %q, want %q", body["topicId"], enrollTestTopic00ID)
		}

		progresses := getProgresses(t, ctx, pool)
		if len(progresses) != 1 {
			t.Fatalf("progressが1件であること: got %d", len(progresses))
		}
		if progresses[0].topicID != enrollTestTopic00ID || progresses[0].status != "IN_PROGRESS" {
			t.Errorf("progress不一致: got topicID=%q status=%q", progresses[0].topicID, progresses[0].status)
		}
	})

	t.Run("最大のtopicがCOMPLETED（セクション最後でない）の場合、次のtopicを開始できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupProgresses(ctx, pool) })
		if err := insertProgress(ctx, pool, enrollTestTopic00ID, "COMPLETED"); err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`"}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != enrollTestTopic01ID {
			t.Errorf("topicIdが一致すること: got %q, want %q", body["topicId"], enrollTestTopic01ID)
		}

		progresses := getProgresses(t, ctx, pool)
		if len(progresses) != 2 {
			t.Fatalf("progressが2件であること: got %d", len(progresses))
		}
		wantProgresses := map[string]string{
			enrollTestTopic00ID: "COMPLETED",
			enrollTestTopic01ID: "IN_PROGRESS",
		}
		for _, p := range progresses {
			if wantProgresses[p.topicID] != p.status {
				t.Errorf("progress不一致: topicID=%q got=%q want=%q", p.topicID, p.status, wantProgresses[p.topicID])
			}
		}
	})

	t.Run("最大のtopicがCOMPLETED（セクション最後）の場合、次のセクションの最初のtopicを開始できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupProgresses(ctx, pool) })
		if err := insertProgress(ctx, pool, enrollTestTopic00ID, "COMPLETED"); err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}
		if err := insertProgress(ctx, pool, enrollTestTopic01ID, "COMPLETED"); err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`"}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != enrollTestTopic10ID {
			t.Errorf("topicIdが一致すること: got %q, want %q", body["topicId"], enrollTestTopic10ID)
		}

		progresses := getProgresses(t, ctx, pool)
		if len(progresses) != 3 {
			t.Fatalf("progressが3件であること: got %d", len(progresses))
		}
		wantProgresses := map[string]string{
			enrollTestTopic00ID: "COMPLETED",
			enrollTestTopic01ID: "COMPLETED",
			enrollTestTopic10ID: "IN_PROGRESS",
		}
		for _, p := range progresses {
			if wantProgresses[p.topicID] != p.status {
				t.Errorf("progress不一致: topicID=%q got=%q want=%q", p.topicID, p.status, wantProgresses[p.topicID])
			}
		}
	})

	t.Run("最大のtopicがIN_PROGRESSの場合、そのtopicを開始できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupProgresses(ctx, pool) })
		if err := insertProgress(ctx, pool, enrollTestTopic00ID, "IN_PROGRESS"); err != nil {
			t.Fatalf("テストデータの挿入に失敗しました: %v", err)
		}

		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`"}`))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != enrollTestTopic00ID {
			t.Errorf("topicIdが一致すること: got %q, want %q", body["topicId"], enrollTestTopic00ID)
		}

		progresses := getProgresses(t, ctx, pool)
		if len(progresses) != 1 {
			t.Fatalf("progressが1件であること: got %d", len(progresses))
		}
		if progresses[0].topicID != enrollTestTopic00ID || progresses[0].status != "IN_PROGRESS" {
			t.Errorf("progress不一致: got topicID=%q status=%q", progresses[0].topicID, progresses[0].status)
		}
	})

	t.Run("draftコースに作成者でないユーザーが受講開始しようとした場合403を返す", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/courses/"+enrollTestDraftCourseID+"/enrollment", strings.NewReader(`{"userId":"`+enrollTestUserID+`"}`))
		req.SetPathValue("courseId", enrollTestDraftCourseID)
		handler.PostEnrollmentHandler(w, req)

		if w.Result().StatusCode != http.StatusForbidden {
			t.Errorf("ステータスコードが403であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["code"] != "FORBIDDEN" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
	})

	t.Run("draftコースに作成者本人が受講開始できる", func(t *testing.T) {
		t.Cleanup(func() { cleanupProgresses(ctx, pool) })

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/v1/courses/"+enrollTestDraftCourseID+"/enrollment", strings.NewReader(`{"userId":"`+enrollTestAuthorID+`"}`))
		req.SetPathValue("courseId", enrollTestDraftCourseID)
		handler.PostEnrollmentHandler(w, req)

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["topicId"] != enrollTestDraftTopic00 {
			t.Errorf("topicIdが一致すること: got %q, want %q", body["topicId"], enrollTestDraftTopic00)
		}
	})

	t.Run("topicNumberのみ指定した場合400を返す", func(t *testing.T) {
		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`","topicNumber":0}`))

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("ステータスコードが400であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["code"] != "INPUT_VALIDATION_ERROR" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
	})

	t.Run("存在しないsectionNumberを指定した場合400を返す", func(t *testing.T) {
		w := httptest.NewRecorder()
		handler.PostEnrollmentHandler(w, newEnrollRequest(t, `{"userId":"`+enrollTestUserID+`","sectionNumber":999}`))

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("ステータスコードが400であること: got %d", w.Result().StatusCode)
		}
		var body map[string]string
		json.NewDecoder(w.Result().Body).Decode(&body)
		if body["code"] != "INPUT_VALIDATION_ERROR" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
	})
}
