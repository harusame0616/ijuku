package contacts_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/harusame0616/ijuku/apps/api/lib/env"
	"github.com/harusame0616/ijuku/apps/api/routes/contacts"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostContactHandlerMedium(t *testing.T) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, env.Require("DATABASE_URL"))
	if err != nil {
		t.Fatalf("DBへの接続に失敗しました: %v", err)
	}
	defer pool.Close()

	q := db.New(pool)
	handler := contacts.NewPostContactHandler(q)

	cleanup := func(t *testing.T, email string) {
		t.Helper()
		_, _ = pool.Exec(ctx, `DELETE FROM contacts WHERE email = $1`, email)
	}

	t.Run("必須項目を全て指定すると 201 を返し DB に保存される", func(t *testing.T) {
		email := "ok-required@example.com"
		t.Cleanup(func() { cleanup(t, email) })

		body := `{"name":"山田 太郎","email":"` + email + `","content":"問い合わせ本文"}`
		r := httptest.NewRequest(http.MethodPost, "/v1/contacts", strings.NewReader(body))
		r.Header.Set("X-Forwarded-For", "203.0.113.10, 198.51.100.1")
		r.Header.Set("User-Agent", "TestAgent/1.0")
		w := httptest.NewRecorder()

		handler.PostContactHandler(w, r)

		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

		var (
			name      string
			savedEmail string
			phone     *string
			content   string
			ipAddress string
			userAgent string
		)
		require.NoError(t, pool.QueryRow(ctx,
			`SELECT name, email, phone, content, ip_address, user_agent FROM contacts WHERE email = $1`,
			email,
		).Scan(&name, &savedEmail, &phone, &content, &ipAddress, &userAgent))
		assert.Equal(t, "山田 太郎", name)
		assert.Equal(t, email, savedEmail)
		assert.Nil(t, phone)
		assert.Equal(t, "問い合わせ本文", content)
		assert.Equal(t, "203.0.113.10", ipAddress)
		assert.Equal(t, "TestAgent/1.0", userAgent)
	})

	t.Run("phone を指定した場合 DB に保存される", func(t *testing.T) {
		email := "ok-phone@example.com"
		t.Cleanup(func() { cleanup(t, email) })

		body := `{"name":"山田","email":"` + email + `","phone":"03-1234-5678","content":"本文"}`
		r := httptest.NewRequest(http.MethodPost, "/v1/contacts", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.PostContactHandler(w, r)

		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

		var phone *string
		require.NoError(t, pool.QueryRow(ctx,
			`SELECT phone FROM contacts WHERE email = $1`, email,
		).Scan(&phone))
		require.NotNil(t, phone)
		assert.Equal(t, "03-1234-5678", *phone)
	})

	t.Run("不正な JSON は INPUT_VALIDATION_ERROR を返す", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/v1/contacts", strings.NewReader("not json"))
		w := httptest.NewRecorder()

		handler.PostContactHandler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var responseBody map[string]any
		require.NoError(t, json.NewDecoder(w.Result().Body).Decode(&responseBody))
		assert.Equal(t, "INPUT_VALIDATION_ERROR", responseBody["errorCode"])
	})

	t.Run("name が空の場合 INPUT_VALIDATION_ERROR を返す", func(t *testing.T) {
		body := `{"name":"   ","email":"a@example.com","content":"本文"}`
		r := httptest.NewRequest(http.MethodPost, "/v1/contacts", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.PostContactHandler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var responseBody map[string]any
		require.NoError(t, json.NewDecoder(w.Result().Body).Decode(&responseBody))
		assert.Equal(t, "INPUT_VALIDATION_ERROR", responseBody["errorCode"])
		assert.Equal(t, "name is required", responseBody["message"])
	})

	t.Run("email 形式不正は INPUT_VALIDATION_ERROR を返す", func(t *testing.T) {
		body := `{"name":"山田","email":"not-email","content":"本文"}`
		r := httptest.NewRequest(http.MethodPost, "/v1/contacts", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.PostContactHandler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var responseBody map[string]any
		require.NoError(t, json.NewDecoder(w.Result().Body).Decode(&responseBody))
		assert.Equal(t, "email format is invalid", responseBody["message"])
	})

	t.Run("content が空の場合 INPUT_VALIDATION_ERROR を返す", func(t *testing.T) {
		body := `{"name":"山田","email":"a@example.com","content":""}`
		r := httptest.NewRequest(http.MethodPost, "/v1/contacts", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.PostContactHandler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var responseBody map[string]any
		require.NoError(t, json.NewDecoder(w.Result().Body).Decode(&responseBody))
		assert.Equal(t, "content is required", responseBody["message"])
	})

	t.Run("name が上限超過の場合 INPUT_VALIDATION_ERROR を返す", func(t *testing.T) {
		longName := strings.Repeat("あ", 101)
		body := `{"name":"` + longName + `","email":"a@example.com","content":"本文"}`
		r := httptest.NewRequest(http.MethodPost, "/v1/contacts", strings.NewReader(body))
		w := httptest.NewRecorder()

		handler.PostContactHandler(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		var responseBody map[string]any
		require.NoError(t, json.NewDecoder(w.Result().Body).Decode(&responseBody))
		assert.Equal(t, "name is too long", responseBody["message"])
	})

	t.Run("X-Forwarded-For が無い場合 RemoteAddr のホストが ip_address に保存される", func(t *testing.T) {
		email := "remote-addr@example.com"
		t.Cleanup(func() { cleanup(t, email) })

		body := `{"name":"山田","email":"` + email + `","content":"本文"}`
		r := httptest.NewRequest(http.MethodPost, "/v1/contacts", strings.NewReader(body))
		r.RemoteAddr = "192.0.2.5:54321"
		w := httptest.NewRecorder()

		handler.PostContactHandler(w, r)

		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)

		var ipAddress string
		require.NoError(t, pool.QueryRow(ctx,
			`SELECT ip_address FROM contacts WHERE email = $1`, email,
		).Scan(&ipAddress))
		assert.Equal(t, "192.0.2.5", ipAddress)
	})
}
