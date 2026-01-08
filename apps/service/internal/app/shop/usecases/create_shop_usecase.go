package usecases

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
	"github.com/reno1r/weiss/apps/service/internal/validationutil"
)

type CreateShopUsecase struct {
	shopRepository repositories.ShopRepository
	validator      *validator.Validate
}

func NewCreateShopUsecase(shopRepository repositories.ShopRepository) *CreateShopUsecase {
	return &CreateShopUsecase{
		shopRepository: shopRepository,
		validator:      validator.New(),
	}
}

type CreateShopData struct {
	Name        string `validate:"required,min=2,max=255"`
	Description string `validate:"required,min=10,max=1000"`
	Address     string `validate:"required,min=5,max=255"`
	Phone       string `validate:"required,min=10,max=20"`
	Email       string `validate:"required,email"`
	Website     string `validate:"required,url"`
	Logo        string `validate:"required,min=1,max=255"`
}

type CreateShopResult struct {
	Shop *entities.Shop
}

func (u *CreateShopUsecase) Execute(req CreateShopData) (*CreateShopResult, error) {
	if err := u.validator.Struct(req); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, validationutil.GetValidationErrorMessage(err))
		}
		return nil, fmt.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
	}

	shop := entities.Shop{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Phone:       req.Phone,
		Email:       req.Email,
		Website:     req.Website,
		Logo:        req.Logo,
	}

	createdShop, err := u.shopRepository.Create(shop)
	if err != nil {
		return nil, fmt.Errorf("failed to create shop: %w", err)
	}

	return &CreateShopResult{
		Shop: &createdShop,
	}, nil
}
