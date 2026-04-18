package queries

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/env"
	"github.com/stretchr/testify/assert"

	"github.com/jackc/pgx/v5/pgxpool"
)

type topicFixture struct {
	AuthorId           string
	CategoryId         string
	CourseId           string
	SectionId          string
	TopicId            string
	Title              string
	Description        string
	Prerequisites      string
	Knowledge          string
	Flow               string
	Quiz               string
	CompletionCriteria string
}

func createTopicFixture(t *testing.T, pool *pgxpool.Pool, publishStatus string, authorId string) topicFixture {
	t.Helper()
	ctx := context.Background()

	if authorId == "" {
		authorId = uuid.NewString()
	}
	if _, err := pool.Exec(ctx, "INSERT INTO authors (author_id, name, profile) VALUES ($1, $2, $3)", authorId, "Test Author", "テスト用の author です"); err != nil {
		t.Fatalf("author の insert に失敗 : %v", err)
	}
	t.Cleanup(func() { pool.Exec(ctx, "DELETE FROM authors WHERE author_id = $1", authorId) })

	categoryId := uuid.NewString()
	if _, err := pool.Exec(ctx, "INSERT INTO categories (category_id, name, path) VALUES ($1, $2, $3)", categoryId, "Test Category", "test"); err != nil {
		t.Fatalf("categories の insert に失敗 : %v", err)
	}
	t.Cleanup(func() { pool.Exec(ctx, "DELETE FROM categories WHERE category_id = $1", categoryId) })

	courseId := uuid.NewString()
	if _, err := pool.Exec(
		ctx,
		"INSERT INTO courses (course_id, title, description, slug, tags, publish_status, category_id, published_at, author_id, visibility) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		courseId, "テストコース", "テストコースの説明", "test-course-"+courseId, []string{"tag1", "tag2"}, publishStatus, categoryId, nil, authorId, "public",
	); err != nil {
		t.Fatalf("course の insert に失敗 : %v", err)
	}
	t.Cleanup(func() { pool.Exec(ctx, "DELETE FROM courses WHERE course_id = $1", courseId) })

	sectionId := uuid.NewString()
	if _, err := pool.Exec(ctx, "INSERT INTO course_sections (course_section_id, course_id, index, title, description) VALUES ($1, $2, $3, $4, $5)",
		sectionId, courseId, 1, "セクション1", "セクション1の説明",
	); err != nil {
		t.Fatalf("course_sections の insert に失敗 : %v", err)
	}
	t.Cleanup(func() { pool.Exec(ctx, "DELETE FROM course_sections WHERE course_section_id = $1", sectionId) })

	topicId := uuid.NewString()
	title := "トピック1"
	description := "トピック1の説明"
	prerequisites := "前提条件"
	knowledge := "習得知識"
	flow := "学習の流れ"
	quiz := "クイズ"
	completionCriteria := "修了条件"
	if _, err := pool.Exec(ctx, "INSERT INTO course_section_topics (course_id, course_section_id, course_section_topic_id, index, title, description, prerequisites, knowledge, flow, quiz, completion_criteria) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)",
		courseId, sectionId, topicId, 1, title, description, prerequisites, knowledge, flow, quiz, completionCriteria,
	); err != nil {
		t.Fatalf("course_section_topics の insert に失敗 : %v", err)
	}
	t.Cleanup(func() {
		pool.Exec(ctx, "DELETE FROM course_section_topics WHERE course_section_topic_id = $1", topicId)
	})

	return topicFixture{
		AuthorId:           authorId,
		CategoryId:         categoryId,
		CourseId:           courseId,
		SectionId:          sectionId,
		TopicId:            topicId,
		Title:              title,
		Description:        description,
		Prerequisites:      prerequisites,
		Knowledge:          knowledge,
		Flow:               flow,
		Quiz:               quiz,
		CompletionCriteria: completionCriteria,
	}
}

func TestGetTopicDetailHandlerMedium(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, env.Require("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	q := db.New(pool)
	handler := NewTopicDetailHandler(q)

	t.Run("published な講座のトピック詳細が取得できる", func(t *testing.T) {
		f := createTopicFixture(t, pool, "published", "")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/v1/courses/"+f.CourseId+"/sections/"+f.SectionId+"/topics/"+f.TopicId, nil)
		r.SetPathValue("topicId", f.TopicId)
		r.SetPathValue("sectionId", f.SectionId)
		r.SetPathValue("courseId", f.CourseId)

		handler.GetTopicDetailHandler(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var body map[string]any
		json.NewDecoder(w.Body).Decode(&body)

		assert.Equal(t, f.CourseId, body["courseId"])
		assert.Equal(t, f.SectionId, body["sectionId"])
		assert.Equal(t, f.TopicId, body["topicId"])
		assert.Equal(t, f.Title, body["title"])
		assert.Equal(t, f.Description, body["description"])
		assert.Equal(t, f.Prerequisites, body["prerequisites"])
		assert.Equal(t, f.Knowledge, body["knowledge"])
		assert.Equal(t, f.Flow, body["flow"])
		assert.Equal(t, f.Quiz, body["quiz"])
		assert.Equal(t, f.CompletionCriteria, body["completionCriteria"])
	})

	t.Run("自分が作成した draft な講座のトピック詳細が取得できる", func(t *testing.T) {
		f := createTopicFixture(t, pool, "draft", "")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/v1/courses/"+f.CourseId+"/sections/"+f.SectionId+"/topics/"+f.TopicId+"?userId="+f.AuthorId, nil)
		r.SetPathValue("topicId", f.TopicId)
		r.SetPathValue("sectionId", f.SectionId)
		r.SetPathValue("courseId", f.CourseId)

		handler.GetTopicDetailHandler(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var body map[string]any
		json.NewDecoder(w.Body).Decode(&body)

		assert.Equal(t, f.CourseId, body["courseId"])
		assert.Equal(t, f.SectionId, body["sectionId"])
		assert.Equal(t, f.TopicId, body["topicId"])
		assert.Equal(t, f.Title, body["title"])
		assert.Equal(t, f.Description, body["description"])
		assert.Equal(t, f.Prerequisites, body["prerequisites"])
		assert.Equal(t, f.Knowledge, body["knowledge"])
		assert.Equal(t, f.Flow, body["flow"])
		assert.Equal(t, f.Quiz, body["quiz"])
		assert.Equal(t, f.CompletionCriteria, body["completionCriteria"])
	})

	t.Run("他人が作成した draft な講座のトピック詳細が取得できない", func(t *testing.T) {
		f := createTopicFixture(t, pool, "draft", "")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/v1/courses/"+f.CourseId+"/sections/"+f.SectionId+"/topics/"+f.TopicId+"?userId="+uuid.NewString(), nil)
		r.SetPathValue("topicId", f.TopicId)
		r.SetPathValue("sectionId", f.SectionId)
		r.SetPathValue("courseId", f.CourseId)

		handler.GetTopicDetailHandler(w, r)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

		var body map[string]any
		json.NewDecoder(w.Body).Decode(&body)

		assert.Equal(t, "TOPIC_DETAIL_NOT_FOUND", body["code"])
		assert.Equal(t, "Topic detail is not found", body["message"])
	})

	t.Run("userId なしで draft な講座のトピック詳細が取得できない", func(t *testing.T) {
		f := createTopicFixture(t, pool, "draft", "")

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/v1/courses/"+f.CourseId+"/sections/"+f.SectionId+"/topics/"+f.TopicId, nil)
		r.SetPathValue("topicId", f.TopicId)
		r.SetPathValue("sectionId", f.SectionId)
		r.SetPathValue("courseId", f.CourseId)

		handler.GetTopicDetailHandler(w, r)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

		var body map[string]any
		json.NewDecoder(w.Body).Decode(&body)

		assert.Equal(t, "TOPIC_DETAIL_NOT_FOUND", body["code"])
		assert.Equal(t, "Topic detail is not found", body["message"])
	})

	t.Run("存在しないトピックの詳細取得できない", func(t *testing.T) {
		courseId := uuid.NewString()
		sectionId := uuid.NewString()
		topicId := uuid.NewString()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/v1/courses/"+courseId+"/sections/"+sectionId+"/topics/"+topicId, nil)
		r.SetPathValue("topicId", topicId)
		r.SetPathValue("sectionId", sectionId)
		r.SetPathValue("courseId", courseId)

		handler.GetTopicDetailHandler(w, r)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

		var body map[string]any
		json.NewDecoder(w.Body).Decode(&body)

		assert.Equal(t, "TOPIC_DETAIL_NOT_FOUND", body["code"])
		assert.Equal(t, "Topic detail is not found", body["message"])
	})

	t.Run("courseId が UUID フォーマットでない場合、エラーを返す", func(t *testing.T) {
		courseId := "invalid-uuid"
		sectionId := uuid.NewString()
		topicId := uuid.NewString()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/v1/courses/"+courseId+"/sections/"+sectionId+"/topics/"+topicId, nil)
		r.SetPathValue("topicId", topicId)
		r.SetPathValue("sectionId", sectionId)
		r.SetPathValue("courseId", courseId)

		handler.GetTopicDetailHandler(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

		var body map[string]any
		json.NewDecoder(w.Body).Decode(&body)

		assert.Equal(t, "INPUT_VALIDATION_ERROR", body["code"])
		assert.Equal(t, "courseId must be a valid UUID", body["message"])
	})

	t.Run("sectionId が UUID フォーマットでない場合、エラーを返す", func(t *testing.T) {
		courseId := uuid.NewString()
		sectionId := "invalid-uuid"
		topicId := uuid.NewString()

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/v1/courses/"+courseId+"/sections/"+sectionId+"/topics/"+topicId, nil)
		r.SetPathValue("topicId", topicId)
		r.SetPathValue("sectionId", sectionId)
		r.SetPathValue("courseId", courseId)

		handler.GetTopicDetailHandler(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

		var body map[string]any
		json.NewDecoder(w.Body).Decode(&body)

		assert.Equal(t, "INPUT_VALIDATION_ERROR", body["code"])
		assert.Equal(t, "sectionId must be a valid UUID", body["message"])
	})

	t.Run("topicId が UUID フォーマットでない場合、エラーを返す", func(t *testing.T) {
		courseId := uuid.NewString()
		sectionId := uuid.NewString()
		topicId := "invalid-uuid"

		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/v1/courses/"+courseId+"/sections/"+sectionId+"/topics/"+topicId, nil)
		r.SetPathValue("topicId", topicId)
		r.SetPathValue("sectionId", sectionId)
		r.SetPathValue("courseId", courseId)

		handler.GetTopicDetailHandler(w, r)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)

		var body map[string]any
		json.NewDecoder(w.Body).Decode(&body)

		assert.Equal(t, "INPUT_VALIDATION_ERROR", body["code"])
		assert.Equal(t, "topicId must be a valid UUID", body["message"])
	})

}
