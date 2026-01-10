package usecases

import (
	"errors"
	"fmt"

	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/user/repositories"
)

type GetUserUsecase struct {
	userRepository repositories.UserRepository
}

func NewGetUserUsecase(userRepository repositories.UserRepository) *GetUserUsecase {
	return &GetUserUsecase{
		userRepository: userRepository,
	}
}

type GetUserResult struct {
	User *entities.User
}

func (u *GetUserUsecase) Execute(id uint64) (*GetUserResult, error) {
	user, err := u.userRepository.FindByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &GetUserResult{
		User: &user,
	}, nil
}
