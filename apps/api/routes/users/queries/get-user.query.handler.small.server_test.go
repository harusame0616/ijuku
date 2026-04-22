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

const validUserID = "00000000-0000-0000-0000-000000000001"

type mockGetUserQuery struct {
	row db.GetUserRow
	err error
}

func (m *mockGetUserQuery) GetUser(_ context.Context, _ pgtype.UUID) (db.GetUserRow, error) {
	return m.row, m.err
}

func decodeBody(t *testing.T, w *httptest.ResponseRecorder) map[string]string {
	t.Helper()
	var body map[string]string
	if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	return body
}

func newGetUserRequest(t *testing.T, userID string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/"+userID, nil)
	req.SetPathValue("userID", userID)
	return req
}

func TestGetUserHandler(t *testing.T) {
	t.Run("userIDがUUID形式でない場合400を返す", func(t *testing.T) {
		h := NewGetUserHandler(&mockGetUserQuery{})
		w := httptest.NewRecorder()
		h.GetUserHandler(w, newGetUserRequest(t, "invalid-uuid"))
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		body := decodeBody(t, w)
		assert.Equal(t, "INPUT_VALIDATION_ERROR", body["errorCode"])
	})

	t.Run("ユーザーが存在しない場合404を返す", func(t *testing.T) {
		h := NewGetUserHandler(&mockGetUserQuery{err: pgx.ErrNoRows})
		w := httptest.NewRecorder()
		h.GetUserHandler(w, newGetUserRequest(t, validUserID))
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		body := decodeBody(t, w)
		assert.Equal(t, "USER_NOT_FOUND", body["errorCode"])
	})

	t.Run("DBエラーの場合500を返す", func(t *testing.T) {
		h := NewGetUserHandler(&mockGetUserQuery{err: errors.New("db error")})
		w := httptest.NewRecorder()
		h.GetUserHandler(w, newGetUserRequest(t, validUserID))
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("正常な場合200とユーザー情報を返す", func(t *testing.T) {
		h := NewGetUserHandler(&mockGetUserQuery{
			row: db.GetUserRow{Nickname: "テスト", Introduce: "自己紹介"},
		})
		w := httptest.NewRecorder()
		h.GetUserHandler(w, newGetUserRequest(t, validUserID))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		body := decodeBody(t, w)
		assert.Equal(t, "テスト", body["nickname"])
		assert.Equal(t, "自己紹介", body["introduce"])
	})
}
