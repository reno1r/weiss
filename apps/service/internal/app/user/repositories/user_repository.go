package repositories

import "github.com/reno1r/weiss/apps/service/internal/app/user/entities"

type UserRepository interface {
	All() []entities.User
	FindByPhone(phone string) (entities.User, error)
	FindByEmail(email string) (entities.User, error)
	Create(user entities.User) (entities.User, error)
	Update(user entities.User) (entities.User, error)
	Delete(user entities.User) error
}
