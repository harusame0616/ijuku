package commands

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

const (
	validUserId   = "00000000-0000-0000-0000-000000000001"
	validCourseId = "00000000-0000-0000-0000-000000000002"
	validTopicId  = "00000000-0000-0000-0000-000000000003"
)

type stubUsecase struct {
	result UpdateEnrollmentResult
	err    error
}

func (s *stubUsecase) execute(_ context.Context, _ UpdateEnrollmentParams) (UpdateEnrollmentResult, error) {
	return s.result, s.err
}

func newPatchRequest(t *testing.T, userId, courseId, body string) *http.Request {
	t.Helper()
	req := httptest.NewRequest("PATCH", "/v1/users/"+userId+"/enrollments/"+courseId, strings.NewReader(body))
	req.SetPathValue("userID", userId)
	req.SetPathValue("courseId", courseId)
	return req
}

func decodeMap(t *testing.T, w *httptest.ResponseRecorder) map[string]string {
	t.Helper()
	var body map[string]string
	if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	return body
}

func validBody() string {
	return `{"topicId":"` + validTopicId + `","status":"IN_PROGRESS"}`
}

func TestPatchEnrollmentHandler_Validation(t *testing.T) {
	tests := []struct {
		name     string
		userId   string
		courseId string
		body     string
	}{
		{name: "userIDが空", userId: "", courseId: validCourseId, body: validBody()},
		{name: "userIDがUUID形式でない", userId: "not-a-uuid", courseId: validCourseId, body: validBody()},
		{name: "courseIdが空", userId: validUserId, courseId: "", body: validBody()},
		{name: "courseIdがUUID形式でない", userId: validUserId, courseId: "not-a-uuid", body: validBody()},
		{name: "bodyが不正なJSON", userId: validUserId, courseId: validCourseId, body: `invalid-json`},
		{name: "topicIdが空", userId: validUserId, courseId: validCourseId, body: `{"topicId":"","status":"IN_PROGRESS"}`},
		{name: "topicIdがUUID形式でない", userId: validUserId, courseId: validCourseId, body: `{"topicId":"not-a-uuid","status":"IN_PROGRESS"}`},
		{name: "statusが空", userId: validUserId, courseId: validCourseId, body: `{"topicId":"` + validTopicId + `","status":""}`},
		{name: "statusが不正値", userId: validUserId, courseId: validCourseId, body: `{"topicId":"` + validTopicId + `","status":"FOO"}`},
	}

	for _, tt := range tests {
		t.Run(tt.name+" の場合400を返す", func(t *testing.T) {
			h := NewUpdateEnrollmentHandler(&stubUsecase{})
			req := newPatchRequest(t, tt.userId, tt.courseId, tt.body)
			w := httptest.NewRecorder()

			h.PatchEnrollmentHandler(w, req)

			if w.Result().StatusCode != http.StatusBadRequest {
				t.Errorf("ステータスコードが400であること: got %d", w.Result().StatusCode)
			}
			body := decodeMap(t, w)
			if body["errorCode"] != "INPUT_VALIDATION_ERROR" {
				t.Errorf("errorCodeが一致すること: got %q", body["errorCode"])
			}
		})
	}
}

func TestPatchEnrollmentHandler_UsecaseErrorMapping(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		wantStatus    int
		wantErrorCode string
	}{
		{name: "ErrInvalidStatusTransit", err: ErrInvalidStatusTransit, wantStatus: http.StatusBadRequest, wantErrorCode: "INPUT_VALIDATION_ERROR"},
		{name: "ErrInvalidProgressStatus", err: ErrInvalidProgressStatus, wantStatus: http.StatusBadRequest, wantErrorCode: "INPUT_VALIDATION_ERROR"},
		{name: "pgx.ErrNoRows", err: pgx.ErrNoRows, wantStatus: http.StatusNotFound, wantErrorCode: "COURSE_NOT_FOUND"},
		{name: "ErrTopicNotFoundInCourse", err: ErrTopicNotFoundInCourse, wantStatus: http.StatusNotFound, wantErrorCode: "TOPIC_NOT_FOUND"},
		{name: "ErrNotEnrolled", err: ErrNotEnrolled, wantStatus: http.StatusNotFound, wantErrorCode: "NOT_ENROLLED"},
		{name: "予期しないエラー", err: errors.New("unexpected"), wantStatus: http.StatusInternalServerError, wantErrorCode: "SERVER_INTERNAL_ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name+" を適切なステータスにマップする", func(t *testing.T) {
			h := NewUpdateEnrollmentHandler(&stubUsecase{err: tt.err})
			req := newPatchRequest(t, validUserId, validCourseId, validBody())
			w := httptest.NewRecorder()

			h.PatchEnrollmentHandler(w, req)

			if w.Result().StatusCode != tt.wantStatus {
				t.Errorf("ステータスコードが一致すること: got %d, want %d", w.Result().StatusCode, tt.wantStatus)
			}
			body := decodeMap(t, w)
			if body["errorCode"] != tt.wantErrorCode {
				t.Errorf("errorCodeが一致すること: got %q, want %q", body["errorCode"], tt.wantErrorCode)
			}
		})
	}
}

func TestPatchEnrollmentHandler_Success(t *testing.T) {
	t.Run("IN_PROGRESSで成功した場合200を返す", func(t *testing.T) {
		h := NewUpdateEnrollmentHandler(&stubUsecase{
			result: UpdateEnrollmentResult{TopicId: validTopicId, Status: "IN_PROGRESS"},
		})
		req := newPatchRequest(t, validUserId, validCourseId, validBody())
		w := httptest.NewRecorder()

		h.PatchEnrollmentHandler(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("ステータスコードが200であること: got %d", w.Result().StatusCode)
		}
		body := decodeMap(t, w)
		if body["topicId"] != validTopicId {
			t.Errorf("topicIdが一致すること: got %q", body["topicId"])
		}
		if body["status"] != "IN_PROGRESS" {
			t.Errorf("statusが一致すること: got %q", body["status"])
		}
	})

	t.Run("COMPLETEDで成功した場合200を返す", func(t *testing.T) {
		h := NewUpdateEnrollmentHandler(&stubUsecase{
			result: UpdateEnrollmentResult{TopicId: validTopicId, Status: "COMPLETED"},
		})
		req := newPatchRequest(t, validUserId, validCourseId, `{"topicId":"`+validTopicId+`","status":"COMPLETED"}`)
		w := httptest.NewRecorder()

		h.PatchEnrollmentHandler(w, req)

		if w.Result().StatusCode != http.StatusOK {
			t.Errorf("ステータスコードが200であること: got %d", w.Result().StatusCode)
		}
		body := decodeMap(t, w)
		if body["status"] != "COMPLETED" {
			t.Errorf("statusが一致すること: got %q", body["status"])
		}
	})
}
