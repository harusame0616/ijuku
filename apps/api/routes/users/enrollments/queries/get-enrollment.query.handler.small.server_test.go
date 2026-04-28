package queries

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

const (
	validCourseID   = "00000000-0000-0000-0000-0000000000c0"
	validAuthorID   = "00000000-0000-0000-0000-0000000000a0"
	validSection0ID = "00000000-0000-0000-0000-000000000050"
	validSection1ID = "00000000-0000-0000-0000-000000000051"
	validTopic00ID  = "00000000-0000-0000-0000-000000000700"
	validTopic01ID  = "00000000-0000-0000-0000-000000000701"
	validTopic10ID  = "00000000-0000-0000-0000-000000000710"
)

type mockGetEnrollmentQuery struct {
	authority    db.GetCourseAuthorityByIdRow
	authorityErr error
	row          db.GetCourseStructureWithProgressRow
	rowErr       error
}

func (m *mockGetEnrollmentQuery) GetCourseAuthorityById(_ context.Context, _ pgtype.UUID) (db.GetCourseAuthorityByIdRow, error) {
	return m.authority, m.authorityErr
}

func (m *mockGetEnrollmentQuery) GetCourseStructureWithProgress(_ context.Context, _ db.GetCourseStructureWithProgressParams) (db.GetCourseStructureWithProgressRow, error) {
	return m.row, m.rowErr
}

func newGetEnrollmentRequest(t *testing.T, userID, courseID string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/"+userID+"/enrollments/"+courseID, nil)
	req.SetPathValue("userID", userID)
	req.SetPathValue("courseId", courseID)
	return req
}

func mustUUID(t *testing.T, s string) pgtype.UUID {
	t.Helper()
	var u pgtype.UUID
	if err := u.Scan(s); err != nil {
		t.Fatalf("UUID parse failed: %v", err)
	}
	return u
}

func publishedAuthority(t *testing.T) db.GetCourseAuthorityByIdRow {
	t.Helper()
	return db.GetCourseAuthorityByIdRow{
		PublishStatus: "published",
		AuthorID:      mustUUID(t, validAuthorID),
	}
}

func draftAuthorityBy(t *testing.T, authorID string) db.GetCourseAuthorityByIdRow {
	t.Helper()
	return db.GetCourseAuthorityByIdRow{
		PublishStatus: "draft",
		AuthorID:      mustUUID(t, authorID),
	}
}

func sectionsRow(t *testing.T, secs []rawSection) db.GetCourseStructureWithProgressRow {
	t.Helper()
	b, err := json.Marshal(secs)
	if err != nil {
		t.Fatalf("marshal sections: %v", err)
	}
	return db.GetCourseStructureWithProgressRow{Sections: b}
}

func topic(topicID, title, status string, index int) rawTopic {
	return rawTopic{TopicId: topicID, Title: title, Status: status, Index: index}
}

func TestGetEnrollmentHandler_BadRequest(t *testing.T) {
	t.Run("userIDがUUID形式でない場合400を返す", func(t *testing.T) {
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, "invalid-uuid", validCourseID))

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var body map[string]string
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.Equal(t, "INPUT_VALIDATION_ERROR", body["errorCode"])
	})

	t.Run("courseIdがUUID形式でない場合400を返す", func(t *testing.T) {
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, "invalid-uuid"))

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var body map[string]string
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.Equal(t, "INPUT_VALIDATION_ERROR", body["errorCode"])
	})
}

func TestGetEnrollmentHandler_Authority(t *testing.T) {
	t.Run("コースが存在しない場合404を返す", func(t *testing.T) {
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{authorityErr: pgx.ErrNoRows})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		var body map[string]string
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.Equal(t, "COURSE_NOT_FOUND", body["errorCode"])
	})

	t.Run("authority取得でDBエラーの場合500を返す", func(t *testing.T) {
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{authorityErr: errors.New("db error")})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("draftコースかつ著者でない場合403を返す", func(t *testing.T) {
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{authority: draftAuthorityBy(t, validAuthorID)})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
		var body map[string]string
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.Equal(t, "ENROLLMENT_FORBIDDEN", body["errorCode"])
	})

	t.Run("draftコースでも著者本人ならアクセスできる", func(t *testing.T) {
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{
			authority: draftAuthorityBy(t, validUserID),
			row:       sectionsRow(t, []rawSection{}),
		})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("structure取得でDBエラーの場合500を返す", func(t *testing.T) {
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{
			authority: publishedAuthority(t),
			rowErr:    errors.New("db error"),
		})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("sections jsonbが不正なJSONの場合500を返す", func(t *testing.T) {
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{
			authority: publishedAuthority(t),
			row:       db.GetCourseStructureWithProgressRow{Sections: []byte("not a json")},
		})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestGetEnrollmentHandler_Response(t *testing.T) {
	t.Run("sections jsonbが空 (nil) の場合は空sectionsとnextTopic=nullを返す", func(t *testing.T) {
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{
			authority: publishedAuthority(t),
			row:       db.GetCourseStructureWithProgressRow{},
		})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.Equal(t, []getEnrollmentSectionResponse{}, body.Sections)
		assert.Nil(t, body.NextTopic)
	})

	t.Run("courseのtitleがレスポンスに含まれる", func(t *testing.T) {
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{
			authority: publishedAuthority(t),
			row:       db.GetCourseStructureWithProgressRow{Title: "テストコース"},
		})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.Equal(t, "テストコース", body.Title)
	})

	t.Run("全NOT_STARTEDなら先頭トピックがnextTopic", func(t *testing.T) {
		secs := []rawSection{
			{
				SectionId: validSection0ID, Title: "S0", Index: 0,
				Topics: []rawTopic{
					topic(validTopic00ID, "T0-0", "NOT_STARTED", 0),
					topic(validTopic01ID, "T0-1", "NOT_STARTED", 1),
				},
			},
		}
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{
			authority: publishedAuthority(t),
			row:       sectionsRow(t, secs),
		})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.NotNil(t, body.NextTopic)
		assert.Equal(t, validSection0ID, body.NextTopic.SectionId)
		assert.Equal(t, validTopic00ID, body.NextTopic.TopicId)
		assert.Equal(t, "NOT_STARTED", body.Sections[0].Topics[0].Status)
	})

	t.Run("最大indexの非NOT_STARTEDがIN_PROGRESSならそのトピック", func(t *testing.T) {
		secs := []rawSection{
			{
				SectionId: validSection0ID, Title: "S0", Index: 0,
				Topics: []rawTopic{
					topic(validTopic00ID, "T0-0", "COMPLETED", 0),
					topic(validTopic01ID, "T0-1", "IN_PROGRESS", 1),
				},
			},
			{
				SectionId: validSection1ID, Title: "S1", Index: 1,
				Topics: []rawTopic{
					topic(validTopic10ID, "T1-0", "NOT_STARTED", 0),
				},
			},
		}
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{
			authority: publishedAuthority(t),
			row:       sectionsRow(t, secs),
		})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.Equal(t, validSection0ID, body.NextTopic.SectionId)
		assert.Equal(t, validTopic01ID, body.NextTopic.TopicId)
	})

	t.Run("最大indexの非NOT_STARTEDがCOMPLETEDで次があれば次のトピック", func(t *testing.T) {
		secs := []rawSection{
			{
				SectionId: validSection0ID, Title: "S0", Index: 0,
				Topics: []rawTopic{
					topic(validTopic00ID, "T0-0", "COMPLETED", 0),
					topic(validTopic01ID, "T0-1", "NOT_STARTED", 1),
				},
			},
			{
				SectionId: validSection1ID, Title: "S1", Index: 1,
				Topics: []rawTopic{
					topic(validTopic10ID, "T1-0", "NOT_STARTED", 0),
				},
			},
		}
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{
			authority: publishedAuthority(t),
			row:       sectionsRow(t, secs),
		})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.Equal(t, validSection0ID, body.NextTopic.SectionId)
		assert.Equal(t, validTopic01ID, body.NextTopic.TopicId)
	})

	t.Run("セクション末尾COMPLETEDなら次セクションの先頭トピック", func(t *testing.T) {
		secs := []rawSection{
			{
				SectionId: validSection0ID, Title: "S0", Index: 0,
				Topics: []rawTopic{
					topic(validTopic00ID, "T0-0", "COMPLETED", 0),
				},
			},
			{
				SectionId: validSection1ID, Title: "S1", Index: 1,
				Topics: []rawTopic{
					topic(validTopic10ID, "T1-0", "NOT_STARTED", 0),
				},
			},
		}
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{
			authority: publishedAuthority(t),
			row:       sectionsRow(t, secs),
		})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.Equal(t, validSection1ID, body.NextTopic.SectionId)
		assert.Equal(t, validTopic10ID, body.NextTopic.TopicId)
	})

	t.Run("最終トピックCOMPLETED(全完了)ならnextTopic=null", func(t *testing.T) {
		secs := []rawSection{
			{
				SectionId: validSection0ID, Title: "S0", Index: 0,
				Topics: []rawTopic{
					topic(validTopic00ID, "T0-0", "COMPLETED", 0),
					topic(validTopic01ID, "T0-1", "COMPLETED", 1),
				},
			},
			{
				SectionId: validSection1ID, Title: "S1", Index: 1,
				Topics: []rawTopic{
					topic(validTopic10ID, "T1-0", "COMPLETED", 0),
				},
			},
		}
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{
			authority: publishedAuthority(t),
			row:       sectionsRow(t, secs),
		})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.Nil(t, body.NextTopic)
	})

	t.Run("courseIdとセクション/トピック構造を正しく返す", func(t *testing.T) {
		secs := []rawSection{
			{
				SectionId: validSection0ID, Title: "S0", Index: 0,
				Topics: []rawTopic{
					topic(validTopic00ID, "T0-0", "IN_PROGRESS", 0),
					topic(validTopic01ID, "T0-1", "NOT_STARTED", 1),
				},
			},
			{
				SectionId: validSection1ID, Title: "S1", Index: 1,
				Topics: []rawTopic{
					topic(validTopic10ID, "T1-0", "NOT_STARTED", 0),
				},
			},
		}
		h := NewGetEnrollmentHandler(&mockGetEnrollmentQuery{
			authority: publishedAuthority(t),
			row:       sectionsRow(t, secs),
		})
		w := httptest.NewRecorder()
		h.GetEnrollmentHandler(w, newGetEnrollmentRequest(t, validUserID, validCourseID))

		var body GetEnrollmentResponse
		if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		assert.Equal(t, validCourseID, body.CourseId)
		assert.Len(t, body.Sections, 2)
		assert.Equal(t, validSection0ID, body.Sections[0].SectionId)
		assert.Equal(t, "S0", body.Sections[0].Title)
		assert.Len(t, body.Sections[0].Topics, 2)
		assert.Equal(t, "IN_PROGRESS", body.Sections[0].Topics[0].Status)
		assert.Equal(t, "NOT_STARTED", body.Sections[0].Topics[1].Status)
		assert.Equal(t, validSection1ID, body.Sections[1].SectionId)
	})
}

func TestDecideNextTopic(t *testing.T) {
	t.Run("rawSecsが空ならnil", func(t *testing.T) {
		assert.Nil(t, decideNextTopic(nil))
	})

	t.Run("穴あき: 最大index非NOT_STARTEDが最終COMPLETEDなら全完了扱い (= nil)", func(t *testing.T) {
		// A0 COMPLETED, A1 NOT_STARTED, B0 COMPLETED (最終)
		// 最大 index 非 NOT_STARTED = B0 (COMPLETED 最終) -> nil
		secs := []rawSection{
			{
				SectionId: validSection0ID, Title: "S0", Index: 0,
				Topics: []rawTopic{
					topic(validTopic00ID, "T0-0", "COMPLETED", 0),
					topic(validTopic01ID, "T0-1", "NOT_STARTED", 1),
				},
			},
			{
				SectionId: validSection1ID, Title: "S1", Index: 1,
				Topics: []rawTopic{
					topic(validTopic10ID, "T1-0", "COMPLETED", 0),
				},
			},
		}
		got := decideNextTopic(secs)
		assert.Nil(t, got)
	})

	t.Run("穴あき: 最大index非NOT_STARTEDがIN_PROGRESSならそのトピック", func(t *testing.T) {
		// A0 COMPLETED, A1 IN_PROGRESS, B0 NOT_STARTED
		secs := []rawSection{
			{
				SectionId: validSection0ID, Title: "S0", Index: 0,
				Topics: []rawTopic{
					topic(validTopic00ID, "T0-0", "COMPLETED", 0),
					topic(validTopic01ID, "T0-1", "IN_PROGRESS", 1),
				},
			},
			{
				SectionId: validSection1ID, Title: "S1", Index: 1,
				Topics: []rawTopic{
					topic(validTopic10ID, "T1-0", "NOT_STARTED", 0),
				},
			},
		}
		got := decideNextTopic(secs)
		assert.NotNil(t, got)
		assert.Equal(t, validTopic01ID, got.TopicId)
	})
}
