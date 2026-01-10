package repositories

import (
	"context"

	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
)

type ShopRepository interface {
	All(ctx context.Context) []entities.Shop
	FindByID(ctx context.Context, id uint64) (entities.Shop, error)
	FindByPhone(ctx context.Context, phone string) (entities.Shop, error)
	FindByEmail(ctx context.Context, email string) (entities.Shop, error)
	Create(ctx context.Context, shop entities.Shop) (entities.Shop, error)
	Update(ctx context.Context, shop entities.Shop) (entities.Shop, error)
	Delete(ctx context.Context, shop entities.Shop) error
}
