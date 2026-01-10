package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"

	accessentities "github.com/reno1r/weiss/apps/service/internal/app/access/entities"
	accessrepositories "github.com/reno1r/weiss/apps/service/internal/app/access/repositories"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
	"github.com/reno1r/weiss/apps/service/internal/validationutil"
)

type CreateShopUsecase struct {
	db              *gorm.DB
	shopRepository  repositories.ShopRepository
	roleRepository  accessrepositories.RoleRepository
	staffRepository accessrepositories.StaffRepository
	validator       *validator.Validate
}

func NewCreateShopUsecase(db *gorm.DB, shopRepository repositories.ShopRepository, roleRepository accessrepositories.RoleRepository, staffRepository accessrepositories.StaffRepository) *CreateShopUsecase {
	return &CreateShopUsecase{
		db:              db,
		shopRepository:  shopRepository,
		roleRepository:  roleRepository,
		staffRepository: staffRepository,
		validator:       validator.New(),
	}
}

type CreateShopParam struct {
	UserID      uint64 `validate:"required"`
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

func (u *CreateShopUsecase) Execute(ctx context.Context, param CreateShopParam) (*CreateShopResult, error) {
	if err := u.validator.Struct(param); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, validationutil.GetValidationErrorMessage(err))
		}
		return nil, fmt.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
	}

	var result *CreateShopResult

	// Execute all operations within a transaction
	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Create temporary repository instances with transaction DB
		txShopRepo := repositories.NewShopRepository(tx)
		txRoleRepo := accessrepositories.NewRoleRepository(tx)
		txStaffRepo := accessrepositories.NewStaffRepository(tx)

		// Create shop
		shop := entities.Shop{
			Name:        param.Name,
			Description: param.Description,
			Address:     param.Address,
			Phone:       param.Phone,
			Email:       param.Email,
			Website:     param.Website,
			Logo:        param.Logo,
		}

		createdShop, err := txShopRepo.Create(ctx, shop)
		if err != nil {
			return fmt.Errorf("failed to create shop: %w", err)
		}

		// Create default owner role for the shop
		ownerRole := accessentities.Role{
			Name:        "Owner",
			Description: "Shop owner with full access to manage the shop",
			ShopID:      createdShop.ID,
		}

		createdOwnerRole, err := txRoleRepo.Create(ctx, ownerRole)
		if err != nil {
			return fmt.Errorf("failed to create owner role: %w", err)
		}

		// Assign the authenticated user as the owner
		staff := accessentities.Staff{
			UserID: param.UserID,
			ShopID: createdShop.ID,
			RoleID: createdOwnerRole.ID,
		}

		_, err = txStaffRepo.Create(ctx, staff)
		if err != nil {
			return fmt.Errorf("failed to assign user as owner: %w", err)
		}

		// Set result only if all operations succeed
		result = &CreateShopResult{
			Shop: &createdShop,
		}

		// Return nil to commit the transaction
		return nil
	})

	if err != nil {
		// Transaction was rolled back automatically
		return nil, err
	}

	return result, nil
}
