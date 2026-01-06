package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
)

type shopRepository struct {
	db *gorm.DB
}

func NewShopRepository(db *gorm.DB) ShopRepository {
	return &shopRepository{
		db: db,
	}
}

func (r *shopRepository) All() []entities.Shop {
	var shops []entities.Shop
	r.db.Find(&shops)
	return shops
}

func (r *shopRepository) FindByID(id uint64) (entities.Shop, error) {
	var shop entities.Shop
	err := r.db.Where("id = ?", id).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shop, errors.New("shop not found")
		}
		return shop, err
	}
	return shop, nil
}

func (r *shopRepository) FindByPhone(phone string) (entities.Shop, error) {
	var shop entities.Shop
	err := r.db.Where("phone = ?", phone).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shop, errors.New("shop not found")
		}
		return shop, err
	}
	return shop, nil
}

func (r *shopRepository) FindByEmail(email string) (entities.Shop, error) {
	var shop entities.Shop
	err := r.db.Where("email = ?", email).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shop, errors.New("shop not found")
		}
		return shop, err
	}
	return shop, nil
}

func (r *shopRepository) Create(shop entities.Shop) (entities.Shop, error) {
	err := r.db.Create(&shop).Error
	if err != nil {
		return shop, err
	}
	return shop, nil
}

func (r *shopRepository) Update(shop entities.Shop) (entities.Shop, error) {
	err := r.db.Save(&shop).Error
	if err != nil {
		return shop, err
	}
	return shop, nil
}

func (r *shopRepository) Delete(shop entities.Shop) error {
	return r.db.Delete(&shop).Error
}
