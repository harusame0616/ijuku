package commands

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/lib/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test-secret"
const testUserID = "00000000-0000-0000-0000-000000000001"
const otherUserID = "00000000-0000-0000-0000-000000000002"

func newTestVerifier(t *testing.T) *auth.Verifier {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"keys": []any{}})
	}))
	t.Cleanup(srv.Close)
	return auth.NewVerifier(testSecret, srv.URL)
}

func signToken(t *testing.T, sub string, expired bool) string {
	t.Helper()
	exp := time.Now().Add(time.Hour)
	if expired {
		exp = time.Now().Add(-time.Hour)
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": exp.Unix(),
	}).SignedString([]byte(testSecret))
	require.NoError(t, err)
	return token
}

type mockUpdateUsecase struct{ err error }

func (m *mockUpdateUsecase) Execute(_ context.Context, _ uuid.UUID, _, _ string) error {
	return m.err
}

func newPatchRequest(t *testing.T, userID, body, token string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodPatch, "/v1/users/"+userID, strings.NewReader(body))
	req.SetPathValue("userID", userID)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func decodeStringMap(t *testing.T, w *httptest.ResponseRecorder) map[string]string {
	t.Helper()
	var body map[string]string
	if err := json.NewDecoder(w.Result().Body).Decode(&body); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗しました: %v", err)
	}
	return body
}

func TestPatchUserHandler(t *testing.T) {
	validBody := `{"nickname":"テスト","introduce":"自己紹介"}`

	t.Run("Authorizationヘッダーがない場合401を返す", func(t *testing.T) {
		h := NewUpdateUserHandler(&mockUpdateUsecase{}, newTestVerifier(t))
		w := httptest.NewRecorder()
		h.PatchUserHandler(w, newPatchRequest(t, testUserID, validBody, ""))
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("不正なトークンの場合401を返す", func(t *testing.T) {
		h := NewUpdateUserHandler(&mockUpdateUsecase{}, newTestVerifier(t))
		w := httptest.NewRecorder()
		h.PatchUserHandler(w, newPatchRequest(t, testUserID, validBody, "invalid.token"))
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("JWTのuserIDとパスのuserIDが異なる場合403を返す", func(t *testing.T) {
		h := NewUpdateUserHandler(&mockUpdateUsecase{}, newTestVerifier(t))
		token := signToken(t, otherUserID, false)
		w := httptest.NewRecorder()
		h.PatchUserHandler(w, newPatchRequest(t, testUserID, validBody, token))
		assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
	})

	t.Run("userIDがUUID形式でない場合400を返す", func(t *testing.T) {
		h := NewUpdateUserHandler(&mockUpdateUsecase{}, newTestVerifier(t))
		token := signToken(t, "not-a-uuid", false)
		req := newPatchRequest(t, "not-a-uuid", validBody, token)
		w := httptest.NewRecorder()
		h.PatchUserHandler(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("リクエストボディが不正なJSONの場合400を返す", func(t *testing.T) {
		h := NewUpdateUserHandler(&mockUpdateUsecase{}, newTestVerifier(t))
		token := signToken(t, testUserID, false)
		w := httptest.NewRecorder()
		h.PatchUserHandler(w, newPatchRequest(t, testUserID, "invalid-json", token))
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("バリデーションエラーの場合400を返す", func(t *testing.T) {
		h := NewUpdateUserHandler(&mockUpdateUsecase{err: fmt.Errorf("%w: nickname too short", ErrValidation)}, newTestVerifier(t))
		token := signToken(t, testUserID, false)
		w := httptest.NewRecorder()
		h.PatchUserHandler(w, newPatchRequest(t, testUserID, `{"nickname":"valid","introduce":""}`, token))
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("正常な場合200とnickname・introduceを返す", func(t *testing.T) {
		h := NewUpdateUserHandler(&mockUpdateUsecase{}, newTestVerifier(t))
		token := signToken(t, testUserID, false)
		w := httptest.NewRecorder()
		h.PatchUserHandler(w, newPatchRequest(t, testUserID, validBody, token))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		body := decodeStringMap(t, w)
		assert.Equal(t, "テスト", body["nickname"])
		assert.Equal(t, "自己紹介", body["introduce"])
	})
}
