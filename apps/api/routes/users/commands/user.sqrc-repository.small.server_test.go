package commands

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

type mockUserQuerier struct {
	row       db.GetUserRow
	getErr    error
	updateErr error
}

func (m *mockUserQuerier) GetUser(_ context.Context, _ pgtype.UUID) (db.GetUserRow, error) {
	return m.row, m.getErr
}

func (m *mockUserQuerier) UpdateUser(_ context.Context, _ db.UpdateUserParams) error {
	return m.updateErr
}

func TestUserSqrcRepository(t *testing.T) {
	validUserID := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	t.Run("Get - ユーザーが存在する場合Userを返す", func(t *testing.T) {
		q := &mockUserQuerier{row: db.GetUserRow{Nickname: "テスト", Introduce: "自己紹介"}}
		repo := NewUserSqrcRepository(q)
		user, err := repo.Get(context.Background(), validUserID)
		assert.NoError(t, err)
		dto := user.ToDto()
		assert.Equal(t, "テスト", dto.Nickname)
		assert.Equal(t, "自己紹介", dto.Introduce)
	})

	t.Run("Get - ユーザーが存在しない場合ErrUserNotFoundを返す", func(t *testing.T) {
		q := &mockUserQuerier{getErr: pgx.ErrNoRows}
		repo := NewUserSqrcRepository(q)
		_, err := repo.Get(context.Background(), validUserID)
		assert.ErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("Get - DBエラーの場合エラーを返す", func(t *testing.T) {
		q := &mockUserQuerier{getErr: errors.New("db error")}
		repo := NewUserSqrcRepository(q)
		_, err := repo.Get(context.Background(), validUserID)
		assert.Error(t, err)
		assert.NotErrorIs(t, err, ErrUserNotFound)
	})

	t.Run("Update - 正常な場合nilを返す", func(t *testing.T) {
		q := &mockUserQuerier{}
		repo := NewUserSqrcRepository(q)
		user := UserFromDto(UserDto{UserID: validUserID, Nickname: "テスト", Introduce: ""})
		err := repo.Update(context.Background(), user)
		assert.NoError(t, err)
	})

	t.Run("Update - DBエラーの場合エラーを返す", func(t *testing.T) {
		q := &mockUserQuerier{updateErr: errors.New("db error")}
		repo := NewUserSqrcRepository(q)
		user := UserFromDto(UserDto{UserID: validUserID, Nickname: "テスト", Introduce: ""})
		err := repo.Update(context.Background(), user)
		assert.Error(t, err)
	})
}
