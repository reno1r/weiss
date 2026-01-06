package repositories

import "github.com/reno1r/weiss/apps/service/internal/app/shop/entities"

type ShopRepository interface {
	All() []entities.Shop
	FindByID(id uint64) (entities.Shop, error)
	FindByPhone(phone string) (entities.Shop, error)
	FindByEmail(email string) (entities.Shop, error)
	Create(shop entities.Shop) (entities.Shop, error)
	Update(shop entities.Shop) (entities.Shop, error)
	Delete(shop entities.Shop) error
}
