package apikeys_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/env"
	"github.com/harusame0616/ijuku/apps/api/lib/uuidutils"
	"github.com/harusame0616/ijuku/apps/api/routes/users/settings/apikeys"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateApiKeyHandlerMedium(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, env.Require("DATABASE_URL"))
	if err != nil {
		t.Fatalf("DBへの接続に失敗しました: %v", err)
	}
	defer pool.Close()

	q := db.New(pool)
	handler := apikeys.NewGenerateApiKeyHandler(apikeys.NewGenerateApiKeyUsecase(apikeys.NewApiKeySqrcRepository(q)))

	t.Run("userId が未定義の場合 Internal Server Error を返す", (func(t *testing.T) {
		r := httptest.NewRequest("POST", "/v1/users/{userID}/apikeys", strings.NewReader("{}"))
		w := httptest.NewRecorder()

		handler.GenerateApiKeyHandler(w, r)

		var responseBody map[string]any
		json.NewDecoder(w.Result().Body).Decode(&responseBody)
		assert.Equal(t, map[string]any{
			"errorCode": "SERVER_INTERNAL_ERROR",
			"message":   "An unexpected error occurred. Please try again later.",
		}, responseBody)
	}))

	t.Run("userId が非 UUID の場合 Input Validation Error を返す", (func(t *testing.T) {
		r := httptest.NewRequest("POST", "/v1/users/{userID}/apikeys", strings.NewReader("{}"))
		r.SetPathValue("userID", "invalid-uuid")
		w := httptest.NewRecorder()

		handler.GenerateApiKeyHandler(w, r)

		var responseBody map[string]any
		json.NewDecoder(w.Result().Body).Decode(&responseBody)
		assert.Equal(t, map[string]any{
			"errorCode": "INPUT_VALIDATION_ERROR",
			"message":   "User ID must be valid UUID",
		}, responseBody)
	}))

	t.Run("expiredAt のフォーマットが ISO 8601 フォーマットではない場合 Input Validation Error を返す", (func(t *testing.T) {
		r := httptest.NewRequest("POST", "/v1/users/{userID}/apikeys", strings.NewReader(`{"expiredAt": "invalid format"}`))
		r.SetPathValue("userID", "BD30D30D-01A0-4E43-A00D-1E6EB88A1D54")
		w := httptest.NewRecorder()

		handler.GenerateApiKeyHandler(w, r)

		var responseBody map[string]any
		json.NewDecoder(w.Result().Body).Decode(&responseBody)
		assert.Equal(t, map[string]any{
			"errorCode": "INPUT_VALIDATION_ERROR",
			"message":   "expiredAt must be ISO 8601 format",
		}, responseBody)
	}))

	t.Run("Body のフォーマットが json フォーマットではない場合 Input Validation Error を返す", (func(t *testing.T) {
		r := httptest.NewRequest("POST", "/v1/users/{userID}/apikeys", strings.NewReader(`invalid json`))
		r.SetPathValue("userID", "BD30D30D-01A0-4E43-A00D-1E6EB88A1D54")
		w := httptest.NewRecorder()

		handler.GenerateApiKeyHandler(w, r)

		var responseBody map[string]any
		json.NewDecoder(w.Result().Body).Decode(&responseBody)
		assert.Equal(t, map[string]any{
			"errorCode": "INPUT_VALIDATION_ERROR",
			"message":   "Body must be valid JSON",
		}, responseBody)
	}))

	t.Run("API KEY が上限を超えている場合 Apikey Quota Exceeds Limit Error を返す", (func(t *testing.T) {
		userID := uuidutils.MustNewUuidString()
		if _, err := pool.Exec(ctx, `INSERT INTO users (user_id, nickname) VALUES ($1, 'テストユーザー')`, userID); err != nil {
			require.Fail(t, "データの投入エラー: テストユーザー", err)
		}
		if _, err := pool.Exec(ctx, `INSERT INTO apikeys (apikey_id, user_id, key_hash, plain_suffix, expired_at) VALUES ($1, $2, $3, 'suffix', 'infinity')`, uuidutils.MustNewUuidString(), userID, uuidutils.MustNewUuidString()); err != nil {
			require.Fail(t, "データの投入エラー: API 1", err)
		}
		if _, err := pool.Exec(ctx, `INSERT INTO apikeys (apikey_id, user_id, key_hash, plain_suffix, expired_at) VALUES ($1, $2, $3, 'suffix', 'infinity')`, uuidutils.MustNewUuidString(), userID, uuidutils.MustNewUuidString()); err != nil {
			require.Fail(t, "データの投入エラー: API 2", err)
		}
		if _, err := pool.Exec(ctx, `INSERT INTO apikeys (apikey_id, user_id, key_hash, plain_suffix, expired_at) VALUES ($1, $2, $3, 'suffix', 'infinity')`, uuidutils.MustNewUuidString(), userID, uuidutils.MustNewUuidString()); err != nil {
			require.Fail(t, "データの投入エラー: API 3", err)
		}
		if _, err := pool.Exec(ctx, `INSERT INTO apikeys (apikey_id, user_id, key_hash, plain_suffix, expired_at) VALUES ($1, $2, $3, 'suffix', 'infinity')`, uuidutils.MustNewUuidString(), userID, uuidutils.MustNewUuidString()); err != nil {
			require.Fail(t, "データの投入エラー: API 4", err)
		}
		if _, err := pool.Exec(ctx, `INSERT INTO apikeys (apikey_id, user_id, key_hash, plain_suffix, expired_at) VALUES ($1, $2, $3, 'suffix', 'infinity')`, uuidutils.MustNewUuidString(), userID, uuidutils.MustNewUuidString()); err != nil {
			require.Fail(t, "データの投入エラー: API 5", err)
		}

		r := httptest.NewRequest("POST", "/v1/users/{userID}/apikeys", strings.NewReader("{}"))
		r.SetPathValue("userID", userID)
		w := httptest.NewRecorder()

		handler.GenerateApiKeyHandler(w, r)

		var responseBody map[string]any
		json.NewDecoder(w.Result().Body).Decode(&responseBody)
		assert.Equal(t, map[string]any{
			"errorCode": "APIKEY_QUOTA_EXCEEDS_LIMIT",
			"message":   "Api key quota exceeds limit. Api key quota limit is 5",
		}, responseBody)

		pool.Exec(ctx, "DELETE FROM apikeys WHERE user_id = $1", userID)
		pool.Exec(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	}))

	t.Run("API KEY が上限を登録できる", (func(t *testing.T) {
		userID := uuidutils.MustNewUuidString()
		if _, err := pool.Exec(ctx, `INSERT INTO users (user_id, nickname) VALUES ($1, 'テストユーザー')`, userID); err != nil {
			require.Fail(t, "データの投入エラー: テストユーザー", err)
		}
		if _, err := pool.Exec(ctx, `INSERT INTO apikeys (apikey_id, user_id, key_hash, plain_suffix, expired_at) VALUES ($1, $2, $3, 'suffix', 'infinity')`, uuidutils.MustNewUuidString(), userID, uuidutils.MustNewUuidString()); err != nil {
			require.Fail(t, "データの投入エラー: API 1", err)
		}
		if _, err := pool.Exec(ctx, `INSERT INTO apikeys (apikey_id, user_id, key_hash, plain_suffix, expired_at) VALUES ($1, $2, $3, 'suffix', 'infinity')`, uuidutils.MustNewUuidString(), userID, uuidutils.MustNewUuidString()); err != nil {
			require.Fail(t, "データの投入エラー: API 2", err)
		}
		if _, err := pool.Exec(ctx, `INSERT INTO apikeys (apikey_id, user_id, key_hash, plain_suffix, expired_at) VALUES ($1, $2, $3, 'suffix', 'infinity')`, uuidutils.MustNewUuidString(), userID, uuidutils.MustNewUuidString()); err != nil {
			require.Fail(t, "データの投入エラー: API 3", err)
		}
		if _, err := pool.Exec(ctx, `INSERT INTO apikeys (apikey_id, user_id, key_hash, plain_suffix, expired_at) VALUES ($1, $2, $3, 'suffix', 'infinity')`, uuidutils.MustNewUuidString(), userID, uuidutils.MustNewUuidString()); err != nil {
			require.Fail(t, "データの投入エラー: API 4", err)
		}

		r := httptest.NewRequest("POST", "/v1/users/{userID}/apikeys", strings.NewReader("{}"))
		r.SetPathValue("userID", userID)
		w := httptest.NewRecorder()

		handler.GenerateApiKeyHandler(w, r)

		var responseBody map[string]any
		json.NewDecoder(w.Result().Body).Decode(&responseBody)
		assert.Regexp(t, "jukubox_.+", responseBody["apikey"])

		pool.Exec(ctx, "DELETE FROM apikeys WHERE user_id = $1", userID)
		pool.Exec(ctx, "DELETE FROM users WHERE user_id = $1", userID)
	}))
}
