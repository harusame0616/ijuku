package commands

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProfile(t *testing.T) {
	validUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	newUser := func() User {
		return UserFromDto(UserDto{
			UserID:    validUserID,
			Nickname:  "テスト",
			Introduce: "",
		})
	}

	t.Run("正常な値で更新できる", func(t *testing.T) {
		u := newUser()
		err := u.UpdateProfile("更新後", "自己紹介")
		assert.NoError(t, err)
		dto := u.ToDto()
		assert.Equal(t, "更新後", dto.Nickname)
		assert.Equal(t, "自己紹介", dto.Introduce)
	})

	t.Run("nicknameが空の場合エラーを返す", func(t *testing.T) {
		u := newUser()
		err := u.UpdateProfile("", "")
		assert.Error(t, err)
	})

	t.Run("nicknameが50文字の場合成功する", func(t *testing.T) {
		u := newUser()
		err := u.UpdateProfile(strings.Repeat("あ", 50), "")
		assert.NoError(t, err)
	})

	t.Run("nicknameが51文字の場合エラーを返す", func(t *testing.T) {
		u := newUser()
		err := u.UpdateProfile(strings.Repeat("あ", 51), "")
		assert.Error(t, err)
	})

	t.Run("introduceが500文字の場合成功する", func(t *testing.T) {
		u := newUser()
		err := u.UpdateProfile("テスト", strings.Repeat("あ", 500))
		assert.NoError(t, err)
	})

	t.Run("introduceが501文字の場合エラーを返す", func(t *testing.T) {
		u := newUser()
		err := u.UpdateProfile("テスト", strings.Repeat("あ", 501))
		assert.Error(t, err)
	})

	t.Run("introduceが空の場合成功する", func(t *testing.T) {
		u := newUser()
		err := u.UpdateProfile("テスト", "")
		assert.NoError(t, err)
	})
}
