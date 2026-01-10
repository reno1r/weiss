package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/reno1r/weiss/apps/service/internal/app/access/entities"
)

type staffRepository struct {
	db *gorm.DB
}

func NewStaffRepository(db *gorm.DB) StaffRepository {
	return &staffRepository{
		db: db,
	}
}

func (r *staffRepository) All(ctx context.Context) []entities.Staff {
	var staffs []entities.Staff
	r.db.WithContext(ctx).Find(&staffs)
	for i := range staffs {
		r.db.WithContext(ctx).Preload("User").Preload("Role").Preload("Shop").Find(&staffs[i])
	}
	return staffs
}

func (r *staffRepository) FindByID(ctx context.Context, id uint64) (entities.Staff, error) {
	var staff entities.Staff
	err := r.db.WithContext(ctx).Where("id = ?", id).Preload("User").Preload("Role").Preload("Shop").First(&staff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return staff, errors.New("staff not found")
		}
		return staff, err
	}
	return staff, nil
}

func (r *staffRepository) FindByShopID(ctx context.Context, shopID uint64) []entities.Staff {
	var staffs []entities.Staff
	r.db.WithContext(ctx).Where("shop_id = ?", shopID).Find(&staffs)
	for i := range staffs {
		r.db.WithContext(ctx).Preload("User").Preload("Role").Preload("Shop").Find(&staffs[i])
	}
	return staffs
}

func (r *staffRepository) FindByUserID(ctx context.Context, userID uint64) []entities.Staff {
	var staffs []entities.Staff
	r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&staffs)
	return staffs
}

func (r *staffRepository) FindByRoleID(ctx context.Context, roleID uint64) []entities.Staff {
	var staffs []entities.Staff
	r.db.WithContext(ctx).Where("role_id = ?", roleID).Find(&staffs)
	return staffs
}

func (r *staffRepository) FindByShopIDAndUserID(ctx context.Context, shopID uint64, userID uint64) (entities.Staff, error) {
	var staff entities.Staff
	err := r.db.WithContext(ctx).Where("shop_id = ? AND user_id = ?", shopID, userID).First(&staff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return staff, errors.New("staff not found")
		}
		return staff, err
	}
	return staff, nil
}

func (r *staffRepository) Create(ctx context.Context, staff entities.Staff) (entities.Staff, error) {
	err := r.db.WithContext(ctx).Create(&staff).Error
	if err != nil {
		return staff, err
	}
	return staff, nil
}

func (r *staffRepository) Update(ctx context.Context, staff entities.Staff) (entities.Staff, error) {
	err := r.db.WithContext(ctx).Save(&staff).Error
	if err != nil {
		return staff, err
	}
	return staff, nil
}

func (r *staffRepository) Delete(ctx context.Context, staff entities.Staff) error {
	return r.db.WithContext(ctx).Delete(&staff).Error
}
