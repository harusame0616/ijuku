package apikeys_test

import (
	"testing"

	apikeys "github.com/harusame0616/ijuku/apps/api/routes/users/settings/api-keys"
	"github.com/stretchr/testify/assert"
)

func TestHashedApiKeyEntitySmall(t *testing.T) {
	t.Run("NewHashedApiKey: plain key が正しいフォーマットで生成される", func(t *testing.T) {
		_, plainKey := apikeys.NewHashedApiKey(apikeys.NewHashedApiKeyParams{UserID: "user-id"})

		assert.Regexp(t, `^jukubox_[A-Za-z0-9\-_]{43}$`, plainKey)
	})

	t.Run("NewHashedApiKey: 異なる呼び出しで異なる plain key が生成される", func(t *testing.T) {
		_, plainKey1 := apikeys.NewHashedApiKey(apikeys.NewHashedApiKeyParams{UserID: "user-id"})
		_, plainKey2 := apikeys.NewHashedApiKey(apikeys.NewHashedApiKeyParams{UserID: "user-id"})

		assert.NotEqual(t, plainKey1, plainKey2)
	})
}
