package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/reno1r/weiss/apps/service/internal/app/access/entities"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
	}
}

func (r *roleRepository) All() []entities.Role {
	var roles []entities.Role
	r.db.Find(&roles)
	return roles
}

func (r *roleRepository) FindByID(id uint64) (entities.Role, error) {
	var role entities.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return role, errors.New("role not found")
		}
		return role, err
	}
	return role, nil
}

func (r *roleRepository) FindByShopID(shopID uint64) []entities.Role {
	var roles []entities.Role
	r.db.Where("shop_id = ?", shopID).Find(&roles)
	return roles
}

func (r *roleRepository) Create(role entities.Role) (entities.Role, error) {
	err := r.db.Create(&role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

func (r *roleRepository) Update(role entities.Role) (entities.Role, error) {
	err := r.db.Save(&role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

func (r *roleRepository) Delete(role entities.Role) error {
	return r.db.Delete(&role).Error
}
