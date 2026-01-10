package repositories

import (
	"context"

	"github.com/reno1r/weiss/apps/service/internal/app/access/entities"
)

type StaffRepository interface {
	All(ctx context.Context) []entities.Staff
	FindByID(ctx context.Context, id uint64) (entities.Staff, error)
	FindByShopID(ctx context.Context, shopID uint64) []entities.Staff
	FindByUserID(ctx context.Context, userID uint64) []entities.Staff
	FindByRoleID(ctx context.Context, roleID uint64) []entities.Staff
	FindByShopIDAndUserID(ctx context.Context, shopID uint64, userID uint64) (entities.Staff, error)
	Create(ctx context.Context, staff entities.Staff) (entities.Staff, error)
	Update(ctx context.Context, staff entities.Staff) (entities.Staff, error)
	Delete(ctx context.Context, staff entities.Staff) error
}
