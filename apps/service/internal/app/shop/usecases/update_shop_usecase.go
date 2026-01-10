package usecases

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
	"github.com/reno1r/weiss/apps/service/internal/validationutil"
)

type UpdateShopUsecase struct {
	shopRepository repositories.ShopRepository
	validator      *validator.Validate
}

func NewUpdateShopUsecase(shopRepository repositories.ShopRepository) *UpdateShopUsecase {
	return &UpdateShopUsecase{
		shopRepository: shopRepository,
		validator:      validator.New(),
	}
}

type UpdateShopParam struct {
	ID          uint64 `validate:"required"`
	Name        string `validate:"required,min=2,max=255"`
	Description string `validate:"required,min=10,max=1000"`
	Address     string `validate:"required,min=5,max=255"`
	Phone       string `validate:"required,min=10,max=20"`
	Email       string `validate:"required,email"`
	Website     string `validate:"required,url"`
	Logo        string `validate:"required,min=1,max=255"`
}

type UpdateShopResult struct {
	Shop *entities.Shop
}

func (u *UpdateShopUsecase) Execute(ctx context.Context, param UpdateShopParam) (*UpdateShopResult, error) {
	if err := u.validator.Struct(param); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, validationutil.GetValidationErrorMessage(err))
		}
		return nil, fmt.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
	}

	// Check if shop exists
	_, err := u.shopRepository.FindByID(ctx, param.ID)
	if err != nil {
		return nil, errors.New("shop not found")
	}

	shop := entities.Shop{
		ID:          param.ID,
		Name:        param.Name,
		Description: param.Description,
		Address:     param.Address,
		Phone:       param.Phone,
		Email:       param.Email,
		Website:     param.Website,
		Logo:        param.Logo,
	}

	updatedShop, err := u.shopRepository.Update(ctx, shop)
	if err != nil {
		return nil, fmt.Errorf("failed to update shop: %w", err)
	}

	return &UpdateShopResult{
		Shop: &updatedShop,
	}, nil
}
