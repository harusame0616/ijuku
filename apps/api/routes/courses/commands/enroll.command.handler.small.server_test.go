package commands

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockUsecase struct{}

func (m *mockUsecase) execute(_ context.Context, _ EnrollCourseUsecaseParams) (string, error) {
	return "00000000-0000-0000-0000-000000000003", nil
}

type errorMockUsecase struct{ err error }

func (m *errorMockUsecase) execute(_ context.Context, _ EnrollCourseUsecaseParams) (string, error) {
	return "", m.err
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]string {
	t.Helper()
	var body map[string]string
	if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	return body
}

const validCourseId = "00000000-0000-0000-0000-000000000001"
const validUserId = "00000000-0000-0000-0000-000000000002"

func newRequest(t *testing.T, courseId string, body string) *http.Request {
	t.Helper()
	req := httptest.NewRequest("POST", "/v1/courses/"+courseId+"/enrollment", strings.NewReader(body))
	req.SetPathValue("courseId", courseId)
	return req
}

func TestPostEnrollmentHandler(t *testing.T) {
	t.Run("courseIdが空の場合400を返す", func(t *testing.T) {
		h := NewHandler(&mockUsecase{})
		req := httptest.NewRequest("POST", "/v1/courses//enrollment", strings.NewReader(`{}`))
		req.SetPathValue("courseId", "")
		w := httptest.NewRecorder()

		h.PostEnrollmentHandler(w, req)

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("ステータスコードが400であること: got %d", w.Result().StatusCode)
		}
		body := decodeBody(t, w)
		if body["code"] != "INPUT_VALIDATION_ERROR" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
	})

	t.Run("courseIdがUUID形式でない場合400を返す", func(t *testing.T) {
		h := NewHandler(&mockUsecase{})
		req := newRequest(t, "not-a-uuid", `{}`)
		w := httptest.NewRecorder()

		h.PostEnrollmentHandler(w, req)

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("ステータスコードが400であること: got %d", w.Result().StatusCode)
		}
		body := decodeBody(t, w)
		if body["code"] != "INPUT_VALIDATION_ERROR" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
	})

	t.Run("リクエストボディが不正なJSONの場合400を返す", func(t *testing.T) {
		h := NewHandler(&mockUsecase{})
		req := newRequest(t, validCourseId, `invalid-json`)
		req.SetPathValue("userID", "")
		w := httptest.NewRecorder()

		h.PostEnrollmentHandler(w, req)

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("ステータスコードが400であること: got %d", w.Result().StatusCode)
		}
		body := decodeBody(t, w)
		if body["code"] != "INPUT_VALIDATION_ERROR" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
	})

	t.Run("userIdが空の場合400を返す", func(t *testing.T) {
		h := NewHandler(&mockUsecase{})
		req := newRequest(t, validCourseId, `{"userId":""}`)
		w := httptest.NewRecorder()

		h.PostEnrollmentHandler(w, req)

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("ステータスコードが400であること: got %d", w.Result().StatusCode)
		}
		body := decodeBody(t, w)
		if body["code"] != "INPUT_VALIDATION_ERROR" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
	})

	t.Run("userIdがUUID形式でない場合400を返す", func(t *testing.T) {
		h := NewHandler(&mockUsecase{})
		req := newRequest(t, validCourseId, `{"userId":"not-a-uuid"}`)
		w := httptest.NewRecorder()

		h.PostEnrollmentHandler(w, req)

		if w.Result().StatusCode != http.StatusBadRequest {
			t.Errorf("ステータスコードが400であること: got %d", w.Result().StatusCode)
		}
		body := decodeBody(t, w)
		if body["code"] != "INPUT_VALIDATION_ERROR" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
	})

	t.Run("usecaseが予期しないエラーを返した場合500を返す", func(t *testing.T) {
		h := NewHandler(&errorMockUsecase{err: errors.New("unexpected error")})
		req := newRequest(t, validCourseId, `{"userId":"`+validUserId+`"}`)
		w := httptest.NewRecorder()

		h.PostEnrollmentHandler(w, req)

		if w.Result().StatusCode != http.StatusInternalServerError {
			t.Errorf("ステータスコードが500であること: got %d", w.Result().StatusCode)
		}
		body := decodeBody(t, w)
		if body["code"] != "INTERNAL_SERVER_ERROR" {
			t.Errorf("codeが一致すること: got %q", body["code"])
		}
	})

	t.Run("正常な場合201とtopicIdを返す", func(t *testing.T) {
		h := NewHandler(&mockUsecase{})
		req := newRequest(t, validCourseId, `{"userId":"`+validUserId+`"}`)
		w := httptest.NewRecorder()

		h.PostEnrollmentHandler(w, req)

		if w.Result().StatusCode != http.StatusCreated {
			t.Errorf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		body := decodeBody(t, w)
		if body["topicId"] != "00000000-0000-0000-0000-000000000003" {
			t.Errorf("topicIdが一致すること: got %q", body["topicId"])
		}
	})
}
