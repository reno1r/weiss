package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/reno1r/weiss/apps/service/internal/app/access/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/access/repositories"
)

type GetRoleUsecase struct {
	roleRepository repositories.RoleRepository
}

func NewGetRoleUsecase(roleRepository repositories.RoleRepository) *GetRoleUsecase {
	return &GetRoleUsecase{
		roleRepository: roleRepository,
	}
}

type GetRoleResult struct {
	Role *entities.Role
}

func (u *GetRoleUsecase) Execute(ctx context.Context, id uint64) (*GetRoleResult, error) {
	role, err := u.roleRepository.FindByID(ctx, id)
	if err != nil {
		if err.Error() == "role not found" {
			return nil, errors.New("role not found")
		}
		return nil, fmt.Errorf("failed to get role: %w", err)
	}

	return &GetRoleResult{
		Role: &role,
	}, nil
}
