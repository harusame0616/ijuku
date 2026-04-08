package apikeys

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/harusame0616/ijuku/apps/api/lib/uuid"
)

type hashedApiKey struct {
	apiKeyID          string
	hashedApiKey      string
	plainApiKeySuffix string
	userID            string
	expiredAt         *time.Time
}

type NewHashedApiKeyParams struct {
	UserID    string
	ExpiredAt *time.Time
}

var ErrApiKeyCountExceedsLimit = errors.New("API key count exceeds the limit")

const apiKeyMaxCount = 5

func NewHashedApiKey(params NewHashedApiKeyParams) (hashedApiKey, string) {
	plainKey := generatePlainApiKey()
	key := hashedApiKey{
		apiKeyID:          uuid.MustNewUuidString(),
		userID:            params.UserID,
		hashedApiKey:      getHash(plainKey),
		plainApiKeySuffix: plainKey[len(plainKey)-4:],
		expiredAt:         params.ExpiredAt,
	}

	return key, plainKey
}

// SHA256 を使用しているが、API キーは crypto/rand で生成した 32 バイト（256 ビット）の
// 高エントロピーなトークンであるため、高速なハッシュでも総当たり攻撃は現実的に不可能。
// また、リクエストごとに検証が発生するため bcrypt などの低速アルゴリズムは不適切。
// nolint:gosec
func getHash(plain string) string {
	hash := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(hash[:])
}

func generatePlainApiKey() string {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

	return fmt.Sprintf("jukubox_%s", base64.RawURLEncoding.EncodeToString(b))
}
