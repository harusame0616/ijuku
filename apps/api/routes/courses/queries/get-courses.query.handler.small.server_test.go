package queries

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockCourseQueryService struct{}

func (m *mockCourseQueryService) FindCourses(_ context.Context, _ string, _ string) (GetCoursesResult, error) {
	return GetCoursesResult{Courses: []CoursesItem{}, Cursor: nil}, nil
}

type errorMockCourseQueryService struct{}

func (m *errorMockCourseQueryService) FindCourses(_ context.Context, _ string, _ string) (GetCoursesResult, error) {
	return GetCoursesResult{}, errors.New("database error")
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]string {
	t.Helper()
	var body map[string]string
	if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	return body
}

func TestGetCoursesHandler(t *testing.T) {
	t.Run("キーワードが40文字を超える場合400とエラーメッセージを返す", func(t *testing.T) {
		handlers := NewCoursesHandlers(&mockCourseQueryService{})
		keyword := strings.Repeat("a", 41)
		req := httptest.NewRequest("GET", "/v1/courses?keyword="+keyword, nil)
		w := httptest.NewRecorder()

		handlers.GetCoursesHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("ステータスコードが400であること: got %d", res.StatusCode)
		}
		body := decodeBody(t, w)
		if body["code"] != "INPUT_VALIDATION_ERROR" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
		if body["message"] != "keyword must be 40 characters or less" {
			t.Errorf("messageが一致すること: got %q", body["message"])
		}
	})

	t.Run("不正なカーソル形式の場合400とエラーメッセージを返す", func(t *testing.T) {
		handlers := NewCoursesHandlers(&mockCourseQueryService{})
		req := httptest.NewRequest("GET", "/v1/courses?cursor=not-a-uuid", nil)
		w := httptest.NewRecorder()

		handlers.GetCoursesHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusBadRequest {
			t.Errorf("ステータスコードが400であること: got %d", res.StatusCode)
		}
		body := decodeBody(t, w)
		if body["code"] != "INPUT_VALIDATION_ERROR" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
		if body["message"] != "invalid cursor" {
			t.Errorf("messageが一致すること: got %q", body["message"])
		}
	})

	t.Run("FindCoursesがエラーを返した場合500とエラーメッセージを返す", func(t *testing.T) {
		handlers := NewCoursesHandlers(&errorMockCourseQueryService{})
		req := httptest.NewRequest("GET", "/v1/courses", nil)
		w := httptest.NewRecorder()

		handlers.GetCoursesHandler(w, req)

		res := w.Result()
		if res.StatusCode != http.StatusInternalServerError {
			t.Errorf("ステータスコードが500であること: got %d", res.StatusCode)
		}
		body := decodeBody(t, w)
		if body["code"] != "INTERNAL_SERVER_ERROR" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
		if body["message"] != "internal server error" {
			t.Errorf("messageが一致すること: got %q", body["message"])
		}
	})
}
