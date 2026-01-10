package repositories

import (
	"context"

	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
)

type UserRepository interface {
	All(ctx context.Context) []entities.User
	FindByID(ctx context.Context, id uint64) (entities.User, error)
	FindByPhone(ctx context.Context, phone string) (entities.User, error)
	FindByEmail(ctx context.Context, email string) (entities.User, error)
	Create(ctx context.Context, user entities.User) (entities.User, error)
	Update(ctx context.Context, user entities.User) (entities.User, error)
	Delete(ctx context.Context, user entities.User) error
}
