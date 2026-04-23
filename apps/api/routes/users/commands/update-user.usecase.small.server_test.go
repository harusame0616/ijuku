package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockUserRepo struct {
	user      User
	getErr    error
	updateErr error
}

func (m *mockUserRepo) Get(_ context.Context, _ uuid.UUID) (User, error) {
	return m.user, m.getErr
}

func (m *mockUserRepo) Update(_ context.Context, _ User) error {
	return m.updateErr
}

func TestUpdateUserUsecase(t *testing.T) {
	validUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	t.Run("repo.Getがエラーを返した場合エラーを返す", func(t *testing.T) {
		uc := NewUpdateUserUsecase(&mockUserRepo{getErr: errors.New("db error")})
		err := uc.Execute(context.Background(), validUserID, "テスト", "")
		assert.Error(t, err)
	})

	t.Run("nicknameが不正な場合ErrValidationを返す", func(t *testing.T) {
		user := UserFromDto(UserDto{UserID: validUserID, Nickname: "テスト", Introduce: ""})
		uc := NewUpdateUserUsecase(&mockUserRepo{user: user})
		err := uc.Execute(context.Background(), validUserID, "", "")
		assert.ErrorIs(t, err, ErrValidation)
	})

	t.Run("repo.Updateがエラーを返した場合エラーを返す", func(t *testing.T) {
		user := UserFromDto(UserDto{UserID: validUserID, Nickname: "テスト", Introduce: ""})
		uc := NewUpdateUserUsecase(&mockUserRepo{user: user, updateErr: errors.New("db error")})
		err := uc.Execute(context.Background(), validUserID, "更新後", "")
		assert.Error(t, err)
	})

	t.Run("正常な場合nilを返す", func(t *testing.T) {
		user := UserFromDto(UserDto{UserID: validUserID, Nickname: "テスト", Introduce: ""})
		uc := NewUpdateUserUsecase(&mockUserRepo{user: user})
		err := uc.Execute(context.Background(), validUserID, "更新後", "自己紹介")
		assert.NoError(t, err)
	})
}
