package repositories

import (
	"context"

	"github.com/reno1r/weiss/apps/service/internal/app/access/entities"
)

type RoleRepository interface {
	All(ctx context.Context) []entities.Role
	FindByID(ctx context.Context, id uint64) (entities.Role, error)
	FindByShopID(ctx context.Context, shopID uint64) []entities.Role
	Create(ctx context.Context, role entities.Role) (entities.Role, error)
	Update(ctx context.Context, role entities.Role) (entities.Role, error)
	Delete(ctx context.Context, role entities.Role) error
}
