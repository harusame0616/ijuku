package commands

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/harusame0616/ijuku/apps/api/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type userQuerier interface {
	GetUser(ctx context.Context, userid pgtype.UUID) (db.GetUserRow, error)
	UpdateUser(ctx context.Context, arg db.UpdateUserParams) error
}

type UserSqrcRepository struct {
	q userQuerier
}

func NewUserSqrcRepository(q userQuerier) *UserSqrcRepository {
	return &UserSqrcRepository{q: q}
}

func (r *UserSqrcRepository) Get(ctx context.Context, userID uuid.UUID) (User, error) {
	var uid pgtype.UUID
	_ = uid.Scan(userID.String())

	row, err := r.q.GetUser(ctx, uid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return User{}, ErrUserNotFound
		}
		return User{}, err
	}

	return UserFromDto(UserDto{
		UserID:    userID,
		Nickname:  row.Nickname,
		Introduce: row.Introduce,
	}), nil
}

func (r *UserSqrcRepository) Update(ctx context.Context, user User) error {
	dto := user.ToDto()
	var uid pgtype.UUID
	_ = uid.Scan(dto.UserID.String())

	return r.q.UpdateUser(ctx, db.UpdateUserParams{
		Nickname:  dto.Nickname,
		Introduce: dto.Introduce,
		Userid:    uid,
	})
}
