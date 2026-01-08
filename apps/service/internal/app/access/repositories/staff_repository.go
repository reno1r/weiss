package repositories

import "github.com/reno1r/weiss/apps/service/internal/app/access/entities"

type StaffRepository interface {
	All() []entities.Staff
	FindByID(id uint64) (entities.Staff, error)
	FindByShopID(shopID uint64) []entities.Staff
	FindByUserID(userID uint64) []entities.Staff
	FindByRoleID(roleID uint64) []entities.Staff
	FindByShopIDAndUserID(shopID uint64, userID uint64) (entities.Staff, error)
	Create(staff entities.Staff) (entities.Staff, error)
	Update(staff entities.Staff) (entities.Staff, error)
	Delete(staff entities.Staff) error
}
