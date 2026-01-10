package repositories

import (
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

func (r *userRepository) All() []entities.User {
	var users []entities.User
	r.db.Find(&users)
	return users
}

func (r *userRepository) FindByID(id uint64) (entities.User, error) {
	var user entities.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindByPhone(phone string) (entities.User, error) {
	var user entities.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) FindByEmail(email string) (entities.User, error) {
	var user entities.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return user, errors.New("user not found")
		}
		return user, err
	}
	return user, nil
}

func (r *userRepository) Create(user entities.User) (entities.User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) Update(user entities.User) (entities.User, error) {
	err := r.db.Save(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *userRepository) Delete(user entities.User) error {
	return r.db.Delete(&user).Error
}
