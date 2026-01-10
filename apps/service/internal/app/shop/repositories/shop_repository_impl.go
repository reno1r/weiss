package repositories

import (
	"context"
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

func (r *shopRepository) All(ctx context.Context) []entities.Shop {
	var shops []entities.Shop
	r.db.WithContext(ctx).Find(&shops)
	return shops
}

func (r *shopRepository) FindByID(ctx context.Context, id uint64) (entities.Shop, error) {
	var shop entities.Shop
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shop, errors.New("shop not found")
		}
		return shop, err
	}
	return shop, nil
}

func (r *shopRepository) FindByPhone(ctx context.Context, phone string) (entities.Shop, error) {
	var shop entities.Shop
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shop, errors.New("shop not found")
		}
		return shop, err
	}
	return shop, nil
}

func (r *shopRepository) FindByEmail(ctx context.Context, email string) (entities.Shop, error) {
	var shop entities.Shop
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&shop).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shop, errors.New("shop not found")
		}
		return shop, err
	}
	return shop, nil
}

func (r *shopRepository) Create(ctx context.Context, shop entities.Shop) (entities.Shop, error) {
	err := r.db.WithContext(ctx).Create(&shop).Error
	if err != nil {
		return shop, err
	}
	return shop, nil
}

func (r *shopRepository) Update(ctx context.Context, shop entities.Shop) (entities.Shop, error) {
	err := r.db.WithContext(ctx).Save(&shop).Error
	if err != nil {
		return shop, err
	}
	return shop, nil
}

func (r *shopRepository) Delete(ctx context.Context, shop entities.Shop) error {
	return r.db.WithContext(ctx).Delete(&shop).Error
}
