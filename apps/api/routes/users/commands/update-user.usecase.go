package commands

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrUserNotFound = errors.New("user not found")

type UpdateUserRepository interface {
	Get(ctx context.Context, userID uuid.UUID) (User, error)
	Update(ctx context.Context, user User) error
}

type UpdateUserUsecase struct {
	repo UpdateUserRepository
}

func NewUpdateUserUsecase(repo UpdateUserRepository) *UpdateUserUsecase {
	return &UpdateUserUsecase{repo: repo}
}

func (uc *UpdateUserUsecase) Execute(ctx context.Context, userID uuid.UUID, nickname, introduce string) error {
	user, err := uc.repo.Get(ctx, userID)
	if err != nil {
		return err
	}

	if err := user.UpdateProfile(nickname, introduce); err != nil {
		return err
	}

	return uc.repo.Update(ctx, user)
}
