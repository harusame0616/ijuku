package apikeys

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/auth"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testSecret      = "test-secret"
	testUserID      = "00000000-0000-0000-0000-000000000001"
	otherUserID     = "00000000-0000-0000-0000-000000000002"
	apikeyIDOne     = "11111111-1111-1111-1111-111111111111"
	apikeyIDTwo     = "22222222-2222-2222-2222-222222222222"
)

func newTestVerifier(t *testing.T) *auth.Verifier {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"keys": []any{}})
	}))
	t.Cleanup(srv.Close)
	return auth.NewVerifier(testSecret, srv.URL)
}

func signToken(t *testing.T, sub string) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": time.Now().Add(time.Hour).Unix(),
	}).SignedString([]byte(testSecret))
	require.NoError(t, err)
	return token
}

type mockListApiKeysQuery struct {
	rows []db.ListApiKeysByUserIDRow
	err  error
}

func (m *mockListApiKeysQuery) ListApiKeysByUserID(_ context.Context, _ pgtype.UUID) ([]db.ListApiKeysByUserIDRow, error) {
	return m.rows, m.err
}

func newListRequest(t *testing.T, userID, token string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/v1/users/"+userID+"/settings/apikeys", nil)
	req.SetPathValue("userID", userID)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	return req
}

func mustUUID(t *testing.T, s string) pgtype.UUID {
	t.Helper()
	parsed, err := uuid.Parse(s)
	require.NoError(t, err)
	return pgtype.UUID{Bytes: parsed, Valid: true}
}

func TestListApiKeysHandler(t *testing.T) {
	t.Run("Authorizationヘッダーがない場合401を返す", func(t *testing.T) {
		h := NewListApiKeysHandler(&mockListApiKeysQuery{}, newTestVerifier(t))
		w := httptest.NewRecorder()
		h.ListApiKeysHandler(w, newListRequest(t, testUserID, ""))
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("不正なトークンの場合401を返す", func(t *testing.T) {
		h := NewListApiKeysHandler(&mockListApiKeysQuery{}, newTestVerifier(t))
		w := httptest.NewRecorder()
		h.ListApiKeysHandler(w, newListRequest(t, testUserID, "invalid.token"))
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("JWTのuserIDとパスのuserIDが異なる場合403を返す", func(t *testing.T) {
		h := NewListApiKeysHandler(&mockListApiKeysQuery{}, newTestVerifier(t))
		token := signToken(t, otherUserID)
		w := httptest.NewRecorder()
		h.ListApiKeysHandler(w, newListRequest(t, testUserID, token))
		assert.Equal(t, http.StatusForbidden, w.Result().StatusCode)
	})

	t.Run("userIDがUUID形式でない場合400を返す", func(t *testing.T) {
		h := NewListApiKeysHandler(&mockListApiKeysQuery{}, newTestVerifier(t))
		token := signToken(t, "not-a-uuid")
		w := httptest.NewRecorder()
		h.ListApiKeysHandler(w, newListRequest(t, "not-a-uuid", token))
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("DBエラーの場合500を返す", func(t *testing.T) {
		h := NewListApiKeysHandler(&mockListApiKeysQuery{err: errors.New("db error")}, newTestVerifier(t))
		token := signToken(t, testUserID)
		w := httptest.NewRecorder()
		h.ListApiKeysHandler(w, newListRequest(t, testUserID, token))
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})

	t.Run("0件の場合200と空配列を返す", func(t *testing.T) {
		h := NewListApiKeysHandler(&mockListApiKeysQuery{rows: nil}, newTestVerifier(t))
		token := signToken(t, testUserID)
		w := httptest.NewRecorder()
		h.ListApiKeysHandler(w, newListRequest(t, testUserID, token))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var body ListApiKeysResponse
		require.NoError(t, json.NewDecoder(w.Result().Body).Decode(&body))
		assert.Equal(t, []ApiKeyListItem{}, body.ApiKeys)
	})

	t.Run("複数件と無期限が正しくシリアライズされる", func(t *testing.T) {
		createdAt := time.Date(2026, 1, 10, 12, 0, 0, 0, time.UTC)
		expiredAt := time.Date(2027, 1, 10, 12, 0, 0, 0, time.UTC)
		expiredAtStr := "2027-01-10T12:00:00Z"
		h := NewListApiKeysHandler(&mockListApiKeysQuery{
			rows: []db.ListApiKeysByUserIDRow{
				{
					ApikeyID:    mustUUID(t, apikeyIDOne),
					PlainSuffix: "a3f9",
					CreatedAt:   pgtype.Timestamptz{Time: createdAt, Valid: true},
					ExpiredAt:   pgtype.Timestamptz{Time: expiredAt, Valid: true},
				},
				{
					ApikeyID:    mustUUID(t, apikeyIDTwo),
					PlainSuffix: "c5d1",
					CreatedAt:   pgtype.Timestamptz{Time: createdAt, Valid: true},
					ExpiredAt:   pgtype.Timestamptz{InfinityModifier: pgtype.Infinity, Valid: true},
				},
			},
		}, newTestVerifier(t))
		token := signToken(t, testUserID)
		w := httptest.NewRecorder()
		h.ListApiKeysHandler(w, newListRequest(t, testUserID, token))
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var body ListApiKeysResponse
		require.NoError(t, json.NewDecoder(w.Result().Body).Decode(&body))
		require.Len(t, body.ApiKeys, 2)

		assert.Equal(t, apikeyIDOne, body.ApiKeys[0].ApiKeyID)
		assert.Equal(t, "a3f9", body.ApiKeys[0].Suffix)
		assert.Equal(t, "2026-01-10T12:00:00Z", body.ApiKeys[0].CreatedAt)
		require.NotNil(t, body.ApiKeys[0].ExpiredAt)
		assert.Equal(t, expiredAtStr, *body.ApiKeys[0].ExpiredAt)

		assert.Equal(t, apikeyIDTwo, body.ApiKeys[1].ApiKeyID)
		assert.Equal(t, "c5d1", body.ApiKeys[1].Suffix)
		assert.Nil(t, body.ApiKeys[1].ExpiredAt)
	})
}
