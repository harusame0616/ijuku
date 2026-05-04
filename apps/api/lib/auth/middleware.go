package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/harusame0616/ijuku/apps/api/lib/response"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type contextKey int

const userIDContextKey contextKey = 0

// ApiKeyResolver は plain な API キーから user_id を解決する。
// 未登録キーや期限切れキーは pgx.ErrNoRows を返すこと。
type ApiKeyResolver interface {
	GetUserIDByApiKeyHash(ctx context.Context, keyHash string) (pgtype.UUID, error)
}

// HashApiKey は API キーの hash を生成する。
// hashed-api-key.entity.go の getHash と同じアルゴリズム (sha256 hex)。
//
// API キーは crypto/rand 由来の高エントロピートークンのため、
// 高速ハッシュ (sha256) でも総当たり攻撃は現実的に不可能。
// リクエストごとに検証するため bcrypt 等の低速アルゴリズムは不適切。
//
//nolint:gosec
func HashApiKey(plain string) string {
	hash := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(hash[:])
}

// Middleware は Authorization ヘッダーのトークンを検証し、
// 成功時には userID を Context に詰めて next を呼び出す。
// 検証失敗時は 401 Unauthorized を返す。
//
// トークンは以下の順で評価する:
//  1. JWT として署名検証
//  2. plain な API キーとして apikeys テーブル (key_hash) を検索
//
// API キーが取れたらそのユーザーで認証成功とする。
func Middleware(verifier *Verifier, resolver ApiKeyResolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := ExtractBearerToken(r)
			if err != nil {
				response.WriteErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
				return
			}

			userID, err := resolveUserID(r.Context(), verifier, resolver, token)
			if err != nil {
				response.WriteErrorResponse(w, http.StatusUnauthorized, "UNAUTHORIZED", "unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func resolveUserID(ctx context.Context, verifier *Verifier, resolver ApiKeyResolver, token string) (string, error) {
	if userID, err := verifier.GetUserID(token); err == nil {
		return userID, nil
	}

	if resolver == nil {
		return "", ErrUnauthorized
	}

	row, err := resolver.GetUserIDByApiKeyHash(ctx, HashApiKey(token))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrUnauthorized
		}
		return "", err
	}
	if !row.Valid {
		return "", ErrUnauthorized
	}
	return uuidString(row), nil
}

// uuidString は pgtype.UUID を標準 8-4-4-4-12 形式の文字列に変換する。
func uuidString(u pgtype.UUID) string {
	const hex = "0123456789abcdef"
	b := u.Bytes
	out := make([]byte, 36)
	idx := 0
	for i := 0; i < 16; i++ {
		if i == 4 || i == 6 || i == 8 || i == 10 {
			out[idx] = '-'
			idx++
		}
		out[idx] = hex[b[i]>>4]
		out[idx+1] = hex[b[i]&0x0f]
		idx += 2
	}
	return string(out)
}

// OptionalMiddleware は Middleware の任意認証版。
// Authorization ヘッダーが無い、またはトークンが不正でも 401 にせず素通しする。
// 認証成功時のみ userID を Context に詰める。公開ページで「ログイン時のみ追加情報を出す」用途に使う。
func OptionalMiddleware(verifier *Verifier, resolver ApiKeyResolver) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := ExtractBearerToken(r)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			userID, err := resolveUserID(r.Context(), verifier, resolver, token)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UserIDFromContext は Middleware で詰めた userID を取り出す。
// 認証されていないリクエストでは ok=false を返す。
func UserIDFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(userIDContextKey).(string)
	if !ok || v == "" {
		return "", false
	}
	return v, true
}

// WithUserID はテスト用に Context へ userID を詰める。
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDContextKey, userID)
}
