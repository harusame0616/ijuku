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

const ApiKeyMaxCount = 5

func NewHashedApiKey(params NewHashedApiKeyParams) (hashedApiKey, string) {
	b := make([]byte, 32)

	if _, err := rand.Read(b); err != nil {
		panic(err)
	}

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

func getHash(plain string) string {
	hash := sha256.Sum256([]byte(plain))
	return hex.EncodeToString(hash[:])
}

func generatePlainApiKey() string {
	b := make([]byte, 32)
	rand.Read(b)

	return fmt.Sprintf("jukubox_%s", base64.RawURLEncoding.EncodeToString(b))
}
