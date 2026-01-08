package repositories

import (
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

func (r *staffRepository) All() []entities.Staff {
	var staffs []entities.Staff
	r.db.Find(&staffs)
	return staffs
}

func (r *staffRepository) FindByID(id uint64) (entities.Staff, error) {
	var staff entities.Staff
	err := r.db.Where("id = ?", id).First(&staff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return staff, errors.New("staff not found")
		}
		return staff, err
	}
	return staff, nil
}

func (r *staffRepository) FindByShopID(shopID uint64) []entities.Staff {
	var staffs []entities.Staff
	r.db.Where("shop_id = ?", shopID).Find(&staffs)
	return staffs
}

func (r *staffRepository) FindByUserID(userID uint64) []entities.Staff {
	var staffs []entities.Staff
	r.db.Where("user_id = ?", userID).Find(&staffs)
	return staffs
}

func (r *staffRepository) FindByRoleID(roleID uint64) []entities.Staff {
	var staffs []entities.Staff
	r.db.Where("role_id = ?", roleID).Find(&staffs)
	return staffs
}

func (r *staffRepository) FindByShopIDAndUserID(shopID uint64, userID uint64) (entities.Staff, error) {
	var staff entities.Staff
	err := r.db.Where("shop_id = ? AND user_id = ?", shopID, userID).First(&staff).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return staff, errors.New("staff not found")
		}
		return staff, err
	}
	return staff, nil
}

func (r *staffRepository) Create(staff entities.Staff) (entities.Staff, error) {
	err := r.db.Create(&staff).Error
	if err != nil {
		return staff, err
	}
	return staff, nil
}

func (r *staffRepository) Update(staff entities.Staff) (entities.Staff, error) {
	err := r.db.Save(&staff).Error
	if err != nil {
		return staff, err
	}
	return staff, nil
}

func (r *staffRepository) Delete(staff entities.Staff) error {
	return r.db.Delete(&staff).Error
}
