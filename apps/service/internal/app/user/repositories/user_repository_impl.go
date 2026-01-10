package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/reno1r/weiss/apps/service/internal/app/user/entities"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) All(ctx context.Context) []entities.User {
	var users []entities.User
	r.db.WithContext(ctx).Find(&users)
	return users
}

func (r *userRepository) FindByID(ctx context.Context, id uint64) (entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	var user entities.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) Create(ctx context.Context, user entities.User) (entities.User, error) {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user entities.User) (entities.User, error) {
	err := r.db.WithContext(ctx).Save(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, user entities.User) error {
	return r.db.WithContext(ctx).Delete(&user).Error
}
