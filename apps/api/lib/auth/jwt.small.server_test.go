package auth_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/harusame0616/ijuku/apps/api/lib/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testHMACSecret = "test-secret"

func newVerifierWithHMAC(t *testing.T) *auth.Verifier {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]any{"keys": []any{}})
	}))
	t.Cleanup(srv.Close)
	return auth.NewVerifier(testHMACSecret, srv.URL)
}

func newVerifierWithEC(t *testing.T, kid string, key *ecdsa.PrivateKey) *auth.Verifier {
	t.Helper()
	ecdhKey, err := key.PublicKey.ECDH()
	require.NoError(t, err)
	raw := ecdhKey.Bytes() // 04 || X (32 bytes) || Y (32 bytes)
	xBytes := raw[1:33]
	yBytes := raw[33:65]
	jwks := map[string]any{
		"keys": []any{
			map[string]any{
				"kid": kid,
				"kty": "EC",
				"crv": "P-256",
				"x":   base64.RawURLEncoding.EncodeToString(xBytes),
				"y":   base64.RawURLEncoding.EncodeToString(yBytes),
			},
		},
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(jwks)
	}))
	t.Cleanup(srv.Close)
	return auth.NewVerifier("", srv.URL)
}

func signHMAC(t *testing.T, secret string, expired bool) string {
	t.Helper()
	exp := time.Now().Add(time.Hour)
	if expired {
		exp = time.Now().Add(-time.Hour)
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "user-id",
		"exp": exp.Unix(),
	}).SignedString([]byte(secret))
	require.NoError(t, err)
	return token
}

func signEC(t *testing.T, kid string, key *ecdsa.PrivateKey, expired bool) string {
	t.Helper()
	exp := time.Now().Add(time.Hour)
	if expired {
		exp = time.Now().Add(-time.Hour)
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": "user-id",
		"exp": exp.Unix(),
	})
	tok.Header["kid"] = kid
	token, err := tok.SignedString(key)
	require.NoError(t, err)
	return token
}

func TestVerify(t *testing.T) {
	t.Run("HMAC署名のトークンが正しい場合、nilを返す", func(t *testing.T) {
		v := newVerifierWithHMAC(t)
		token := signHMAC(t, testHMACSecret, false)

		err := v.Verify(token)

		assert.NoError(t, err)
	})

	t.Run("HMAC署名のトークンが不正な場合、ErrUnauthorizedを返す", func(t *testing.T) {
		v := newVerifierWithHMAC(t)
		token := signHMAC(t, "wrong-secret", false)

		err := v.Verify(token)

		assert.ErrorIs(t, err, auth.ErrUnauthorized)
	})

	t.Run("HMACトークンが期限切れの場合、ErrUnauthorizedを返す", func(t *testing.T) {
		v := newVerifierWithHMAC(t)
		token := signHMAC(t, testHMACSecret, true)

		err := v.Verify(token)

		assert.ErrorIs(t, err, auth.ErrUnauthorized)
	})

	t.Run("EC署名のトークンが正しい場合、nilを返す", func(t *testing.T) {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		require.NoError(t, err)
		kid := "test-kid"
		v := newVerifierWithEC(t, kid, key)
		token := signEC(t, kid, key, false)

		err = v.Verify(token)

		assert.NoError(t, err)
	})

	t.Run("EC署名のkidがJWKSに存在しない場合、ErrUnauthorizedを返す", func(t *testing.T) {
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		require.NoError(t, err)
		v := newVerifierWithEC(t, "registered-kid", key)
		token := signEC(t, "unknown-kid", key, false)

		err = v.Verify(token)

		assert.ErrorIs(t, err, auth.ErrUnauthorized)
	})

	t.Run("不正なトークン文字列の場合、ErrUnauthorizedを返す", func(t *testing.T) {
		v := newVerifierWithHMAC(t)

		err := v.Verify("invalid.token.string")

		assert.ErrorIs(t, err, auth.ErrUnauthorized)
	})
}

func TestExtractBearerToken(t *testing.T) {
	t.Run("正しいAuthorizationヘッダーの場合、トークンを返す", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Bearer my-token")

		token, err := auth.ExtractBearerToken(r)

		assert.NoError(t, err)
		assert.Equal(t, "my-token", token)
	})

	t.Run("Authorizationヘッダーがない場合、ErrUnauthorizedを返す", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)

		token, err := auth.ExtractBearerToken(r)

		assert.ErrorIs(t, err, auth.ErrUnauthorized)
		assert.Empty(t, token)
	})

	t.Run("Bearer以外のプレフィックスの場合、ErrUnauthorizedを返す", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("Authorization", "Basic my-token")

		token, err := auth.ExtractBearerToken(r)

		assert.ErrorIs(t, err, auth.ErrUnauthorized)
		assert.Empty(t, token)
	})
}
