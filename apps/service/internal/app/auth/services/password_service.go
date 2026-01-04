package services

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/reno1r/weiss/apps/service/internal/config"
)

type PasswordService struct {
	cost int
}

func NewPasswordService(config *config.Config) *PasswordService {
	cost := config.BcryptCost
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}

	return &PasswordService{
		cost: cost,
	}
}

func (ps *PasswordService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), ps.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (ps *PasswordService) VerifyPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
