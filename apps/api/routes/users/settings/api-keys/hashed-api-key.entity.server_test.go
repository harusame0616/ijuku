package apikeys

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestHashedApiKeyEntitySmall(t *testing.T) {
	t.Run("NewHashedApiKey: plain key が正しいフォーマットで生成される", func(t *testing.T) {
		_, plainKey := NewHashedApiKey(NewHashedApiKeyArg{UserID: "user-id"})

		assert.Regexp(t, `^jukubox_[A-Za-z0-9\-_]{43}$`, plainKey)
	})

	t.Run("NewHashedApiKey: HashedApiKey フィールドが SHA-256 hex 形式になっている", func(t *testing.T) {
		hashed, _ := NewHashedApiKey(NewHashedApiKeyArg{UserID: "user-id"})

		assert.Regexp(t, `^[0-9a-f]{64}$`, hashed.HashedApiKey)
	})

	t.Run("NewHashedApiKey: HashedApiKey が plain key のハッシュと一致する", func(t *testing.T) {
		hashed, plainKey := NewHashedApiKey(NewHashedApiKeyArg{UserID: "user-id"})

		assert.Equal(t, getHash(plainKey), hashed.HashedApiKey)
	})

	t.Run("NewHashedApiKey: PlainApiKeySuffix が plain key の末尾4文字になっている", func(t *testing.T) {
		hashed, plainKey := NewHashedApiKey(NewHashedApiKeyArg{UserID: "user-id"})

		assert.Equal(t, plainKey[len(plainKey)-4:], hashed.PlainApiKeySuffix)
	})

	t.Run("NewHashedApiKey: UserID が引数の値になっている", func(t *testing.T) {
		const userID = "test-user-id"
		hashed, _ := NewHashedApiKey(NewHashedApiKeyArg{UserID: userID})

		assert.Equal(t, userID, hashed.UserID)
	})

	t.Run("NewHashedApiKey: ExpiredAt が nil の場合は nil になっている", func(t *testing.T) {
		hashed, _ := NewHashedApiKey(NewHashedApiKeyArg{UserID: "user-id", ExpiredAt: nil})

		assert.Nil(t, hashed.ExpiredAt)
	})

	t.Run("NewHashedApiKey: ExpiredAt が指定された場合はその値になっている", func(t *testing.T) {
		expiredAt := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)
		hashed, _ := NewHashedApiKey(NewHashedApiKeyArg{UserID: "user-id", ExpiredAt: &expiredAt})

		assert.Equal(t, &expiredAt, hashed.ExpiredAt)
	})

	t.Run("NewHashedApiKey: ApiKeyID が UUID 形式で生成される", func(t *testing.T) {
		hashed, _ := NewHashedApiKey(NewHashedApiKeyArg{UserID: "user-id"})

		assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`, hashed.ApiKeyID)
	})

	t.Run("NewHashedApiKey: 異なる呼び出しで異なる plain key が生成される", func(t *testing.T) {
		_, plainKey1 := NewHashedApiKey(NewHashedApiKeyArg{UserID: "user-id"})
		_, plainKey2 := NewHashedApiKey(NewHashedApiKeyArg{UserID: "user-id"})

		assert.NotEqual(t, plainKey1, plainKey2)
	})
}
