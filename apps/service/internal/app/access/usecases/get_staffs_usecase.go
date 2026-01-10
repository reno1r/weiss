package usecases

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	accessentities "github.com/reno1r/weiss/apps/service/internal/app/access/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/access/repositories"
	shopentities "github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	userentities "github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	"github.com/reno1r/weiss/apps/service/internal/validationutil"
)

func NewGetStaffsUsecase(staffRepository repositories.StaffRepository) *GetStaffsUsecase {
	return &GetStaffsUsecase{
		validator:       validator.New(),
		staffRepository: staffRepository,
	}
}

type GetStaffsUsecase struct {
	validator       *validator.Validate
	staffRepository repositories.StaffRepository
}

type GetStaffsParam struct {
	Shop *shopentities.Shop `validate:"required"`
}

type StaffInfo struct {
	User *userentities.User   `json:"user"`
	Role *accessentities.Role `json:"role"`
}

type GetStaffsResult struct {
	Shop   *shopentities.Shop
	Staffs []StaffInfo
}

func (u *GetStaffsUsecase) Execute(params GetStaffsParam) (*GetStaffsResult, error) {
	if err := u.validator.Struct(params); err != nil {
		var validationErrors []string
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, validationutil.GetValidationErrorMessage(err))
		}
		return nil, fmt.Errorf("validation failed: %s", strings.Join(validationErrors, ", "))
	}

	staffs := u.staffRepository.FindByShopID(params.Shop.ID)

	staffsResult := make([]StaffInfo, len(staffs))
	for i, staff := range staffs {
		staffsResult[i] = StaffInfo{
			User: staff.User,
			Role: staff.Role,
		}
	}

	return &GetStaffsResult{
		Shop:   params.Shop,
		Staffs: staffsResult,
	}, nil
}
