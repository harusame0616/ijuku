package commands

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	libauth "github.com/harusame0616/ijuku/apps/api/lib/auth"
	"github.com/jackc/pgx/v5"
)

const (
	validAuthorSlug = "valid-author"
	validCourseSlug = "valid-course"
)

type stubEnrollUsecase struct {
	result EnrollResult
	err    error
}

func (s *stubEnrollUsecase) execute(_ context.Context, _ EnrollParams) (EnrollResult, error) {
	return s.result, s.err
}

func newEnrollRequest(t *testing.T, userId, body string) *http.Request {
	t.Helper()
	req := httptest.NewRequest("POST", "/v1/me/enrollments", strings.NewReader(body))
	if userId != "" {
		req = req.WithContext(libauth.WithUserID(req.Context(), userId))
	}
	return req
}

func decodeEnrollMap(t *testing.T, w *httptest.ResponseRecorder) map[string]string {
	t.Helper()
	var body map[string]string
	if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	return body
}

func validEnrollBody() string {
	return `{"authorSlug":"` + validAuthorSlug + `","courseSlug":"` + validCourseSlug + `"}`
}

func TestPostEnrollmentHandler_Auth(t *testing.T) {
	t.Run("認証情報が無い場合401を返す", func(t *testing.T) {
		h := NewEnrollHandler(&stubEnrollUsecase{})
		w := httptest.NewRecorder()

		h.PostEnrollmentHandler(w, newEnrollRequest(t, "", validEnrollBody()))

		if w.Result().StatusCode != http.StatusUnauthorized {
			t.Errorf("ステータスコードが401であること: got %d", w.Result().StatusCode)
		}
	})
}

func TestPostEnrollmentHandler_Validation(t *testing.T) {
	tests := []struct {
		name   string
		userId string
		body   string
	}{
		{name: "userIDがUUID形式でない", userId: "not-a-uuid", body: validEnrollBody()},
		{name: "bodyが不正なJSON", userId: validUserId, body: `not-json`},
		{name: "authorSlugが空", userId: validUserId, body: `{"authorSlug":"","courseSlug":"` + validCourseSlug + `"}`},
		{name: "courseSlugが空", userId: validUserId, body: `{"authorSlug":"` + validAuthorSlug + `","courseSlug":""}`},
	}

	for _, tt := range tests {
		t.Run(tt.name+" の場合400を返す", func(t *testing.T) {
			h := NewEnrollHandler(&stubEnrollUsecase{})
			w := httptest.NewRecorder()

			h.PostEnrollmentHandler(w, newEnrollRequest(t, tt.userId, tt.body))

			if w.Result().StatusCode != http.StatusBadRequest {
				t.Errorf("ステータスコードが400であること: got %d", w.Result().StatusCode)
			}
			body := decodeEnrollMap(t, w)
			if body["errorCode"] != "INPUT_VALIDATION_ERROR" {
				t.Errorf("errorCodeが一致すること: got %q", body["errorCode"])
			}
		})
	}
}

func TestPostEnrollmentHandler_UsecaseErrorMapping(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		wantStatus    int
		wantErrorCode string
	}{
		{name: "pgx.ErrNoRows", err: pgx.ErrNoRows, wantStatus: http.StatusNotFound, wantErrorCode: "COURSE_NOT_FOUND"},
		{name: "ErrEnrollmentNotAllowed", err: ErrEnrollmentNotAllowed, wantStatus: http.StatusForbidden, wantErrorCode: "ENROLLMENT_FORBIDDEN"},
		{name: "ErrAlreadyEnrolled", err: ErrAlreadyEnrolled, wantStatus: http.StatusConflict, wantErrorCode: "ALREADY_ENROLLED"},
		{name: "予期しないエラー", err: errors.New("unexpected"), wantStatus: http.StatusInternalServerError, wantErrorCode: "SERVER_INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name+" を適切なステータスにマップする", func(t *testing.T) {
			h := NewEnrollHandler(&stubEnrollUsecase{err: tt.err})
			w := httptest.NewRecorder()

			h.PostEnrollmentHandler(w, newEnrollRequest(t, validUserId, validEnrollBody()))

			if w.Result().StatusCode != tt.wantStatus {
				t.Errorf("ステータスコードが一致すること: got %d, want %d", w.Result().StatusCode, tt.wantStatus)
			}
			body := decodeEnrollMap(t, w)
			if body["errorCode"] != tt.wantErrorCode {
				t.Errorf("errorCodeが一致すること: got %q, want %q", body["errorCode"], tt.wantErrorCode)
			}
		})
	}
}

func TestPostEnrollmentHandler_Success(t *testing.T) {
	t.Run("成功時は201と発行された情報を返す", func(t *testing.T) {
		now := time.Date(2026, 5, 1, 10, 0, 0, 0, time.UTC)
		h := NewEnrollHandler(&stubEnrollUsecase{
			result: EnrollResult{
				CourseId:   uuid.MustParse(validCourseId),
				EnrolledAt: now,
			},
		})
		w := httptest.NewRecorder()

		h.PostEnrollmentHandler(w, newEnrollRequest(t, validUserId, validEnrollBody()))

		if w.Result().StatusCode != http.StatusCreated {
			t.Fatalf("ステータスコードが201であること: got %d", w.Result().StatusCode)
		}
		body := decodeEnrollMap(t, w)
		if body["courseId"] != validCourseId {
			t.Errorf("courseIdが一致すること: got %q", body["courseId"])
		}
		if body["enrolledAt"] != now.Format(time.RFC3339) {
			t.Errorf("enrolledAtがRFC3339で返ること: got %q", body["enrolledAt"])
		}
	})
}
