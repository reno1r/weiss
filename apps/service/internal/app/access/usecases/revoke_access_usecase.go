package usecases

import (
	"github.com/go-playground/validator/v10"
	"github.com/reno1r/weiss/apps/service/internal/app/access/repositories"
)

type RevokeAccessUsecase struct {
	validator       *validator.Validate
	staffRepository repositories.StaffRepository
}

func NewRevokeAccessUsecase() *RevokeAccessUsecase {
	return &RevokeAccessUsecase{}
}
