package repositories

import (
	"context"
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

func (r *roleRepository) All(ctx context.Context) []entities.Role {
	var roles []entities.Role
	r.db.WithContext(ctx).Find(&roles)
	return roles
}

func (r *roleRepository) FindByID(ctx context.Context, id uint64) (entities.Role, error) {
	var role entities.Role
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return role, errors.New("role not found")
		}
		return role, err
	}
	return role, nil
}

func (r *roleRepository) FindByShopID(ctx context.Context, shopID uint64) []entities.Role {
	var roles []entities.Role
	r.db.WithContext(ctx).Where("shop_id = ?", shopID).Find(&roles)
	return roles
}

func (r *roleRepository) Create(ctx context.Context, role entities.Role) (entities.Role, error) {
	err := r.db.WithContext(ctx).Create(&role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

func (r *roleRepository) Update(ctx context.Context, role entities.Role) (entities.Role, error) {
	err := r.db.WithContext(ctx).Save(&role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

func (r *roleRepository) Delete(ctx context.Context, role entities.Role) error {
	return r.db.WithContext(ctx).Delete(&role).Error
}
