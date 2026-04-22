package commands

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
)

var ErrValidation = errors.New("validation error")

type User struct {
	userID    uuid.UUID
	nickname  string
	introduce string
}

type UserDto struct {
	UserID    uuid.UUID
	Nickname  string
	Introduce string
}

func UserFromDto(dto UserDto) User {
	return User{
		userID:    dto.UserID,
		nickname:  dto.Nickname,
		introduce: dto.Introduce,
	}
}

func (u *User) UpdateProfile(nickname, introduce string) error {
	nicknameRunes := []rune(nickname)
	if len(nicknameRunes) < 1 || len(nicknameRunes) > 50 {
		return fmt.Errorf("%w: nickname must be between 1 and 50 characters", ErrValidation)
	}

	introduceRunes := []rune(introduce)
	if len(introduceRunes) > 500 {
		return fmt.Errorf("%w: introduce must be 500 characters or less", ErrValidation)
	}

	u.nickname = nickname
	u.introduce = introduce
	return nil
}

func (u User) ToDto() UserDto {
	return UserDto{
		UserID:    u.userID,
		Nickname:  u.nickname,
		Introduce: u.introduce,
	}
}
