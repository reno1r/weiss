package usecases

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	accessentities "github.com/reno1r/weiss/apps/service/internal/app/access/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/access/repositories"
	shopentities "github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	userentities "github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	"github.com/reno1r/weiss/apps/service/internal/validationutil"
)

type AssignStaffUsecase struct {
	validator       *validator.Validate
	staffRepository repositories.StaffRepository
}

func NewAssignStaffUsecase(staffRepository repositories.StaffRepository) *AssignStaffUsecase {
	return &AssignStaffUsecase{
		validator:       validator.New(),
		staffRepository: staffRepository,
	}
}

type AssignStaffParam struct {
	User *userentities.User   `validate:"required"`
	Shop *shopentities.Shop   `validate:"required"`
	Role *accessentities.Role `validate:"required"`
}

type AssignStaffResult struct {
	Staff *accessentities.Staff
}

func (u *AssignStaffUsecase) Execute(ctx context.Context, params AssignStaffParam) (*AssignStaffResult, error) {
	if err := u.validator.Struct(params); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, validationutil.GetValidationErrorMessage(err))
		}
		return nil, fmt.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
	}

	// Check if staff already exists
	_, err := u.staffRepository.FindByShopIDAndUserID(ctx, params.Shop.ID, params.User.ID)
	if err == nil {
		return nil, fmt.Errorf("staff already assigned to this shop")
	}

	staff := accessentities.Staff{
		UserID: params.User.ID,
		ShopID: params.Shop.ID,
		RoleID: params.Role.ID,
	}

	createdStaff, err := u.staffRepository.Create(ctx, staff)
	if err != nil {
		return nil, fmt.Errorf("failed to create staff: %w", err)
	}

	// Reload with relations
	staffWithRelations, err := u.staffRepository.FindByID(ctx, createdStaff.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to load staff: %w", err)
	}

	return &AssignStaffResult{
		Staff: &staffWithRelations,
	}, nil
}
