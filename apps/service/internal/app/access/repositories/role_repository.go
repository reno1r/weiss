package repositories

import "github.com/reno1r/weiss/apps/service/internal/app/access/entities"

type RoleRepository interface {
	All() []entities.Role
	FindByID(id uint64) (entities.Role, error)
	FindByShopID(shopID uint64) []entities.Role
	Create(role entities.Role) (entities.Role, error)
	Update(role entities.Role) (entities.Role, error)
	Delete(role entities.Role) error
}
