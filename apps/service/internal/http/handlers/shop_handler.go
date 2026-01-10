package handlers

import (
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	accessusecases "github.com/reno1r/weiss/apps/service/internal/app/access/usecases"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/usecases"
	userusecases "github.com/reno1r/weiss/apps/service/internal/app/user/usecases"
)

type ShopHandler struct {
	listShopsUsecase   *usecases.ListShopsUsecase
	getShopUsecase     *usecases.GetShopUsecase
	createShopUsecase  *usecases.CreateShopUsecase
	updateShopUsecase  *usecases.UpdateShopUsecase
	deleteShopUsecase  *usecases.DeleteShopUsecase
	getStaffsUsecase   *accessusecases.GetStaffsUsecase
	assignStaffUsecase *accessusecases.AssignStaffUsecase
	getUserUsecase     *userusecases.GetUserUsecase
	getRoleUsecase     *accessusecases.GetRoleUsecase
}

func NewShopHandler(
	listShopsUsecase *usecases.ListShopsUsecase,
	getShopUsecase *usecases.GetShopUsecase,
	createShopUsecase *usecases.CreateShopUsecase,
	updateShopUsecase *usecases.UpdateShopUsecase,
	deleteShopUsecase *usecases.DeleteShopUsecase,
	getStaffsUsecase *accessusecases.GetStaffsUsecase,
	assignStaffUsecase *accessusecases.AssignStaffUsecase,
	getUserUsecase *userusecases.GetUserUsecase,
	getRoleUsecase *accessusecases.GetRoleUsecase,
) *ShopHandler {
	return &ShopHandler{
		listShopsUsecase:   listShopsUsecase,
		getShopUsecase:     getShopUsecase,
		createShopUsecase:  createShopUsecase,
		updateShopUsecase:  updateShopUsecase,
		deleteShopUsecase:  deleteShopUsecase,
		getStaffsUsecase:   getStaffsUsecase,
		assignStaffUsecase: assignStaffUsecase,
		getUserUsecase:     getUserUsecase,
		getRoleUsecase:     getRoleUsecase,
	}
}

// ListShops godoc
// @Summary      List all shops
// @Description  Get a list of all available shops
// @Tags         shops
// @Accept       json
// @Produce      json
// @Success      200  {object}  ShopListResponse
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /shops [get]
func (h *ShopHandler) ListShops(c fiber.Ctx) error {
	result := h.listShopsUsecase.Execute()

	shops := make([]ShopResponseDTO, len(result.Shops))
	for i, shop := range result.Shops {
		shops[i] = ShopResponseDTO{
			ID:          shop.ID,
			Name:        shop.Name,
			Description: shop.Description,
			Address:     shop.Address,
			Phone:       shop.Phone,
			Email:       shop.Email,
			Website:     shop.Website,
			Logo:        shop.Logo,
			CreatedAt:   shop.CreatedAt,
			UpdatedAt:   shop.UpdatedAt,
		}
	}

	return c.JSON(ShopListResponse{
		Message: "shops retrieved successfully.",
		Data: ShopListResponseData{
			Shops: shops,
		},
	})
}

// GetShop godoc
// @Summary      Get shop by ID
// @Description  Get a specific shop by its ID
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Shop ID"
// @Success      200  {object}  ShopResponse
// @Failure      400  {object}  map[string]string  "Invalid shop id"
// @Failure      404  {object}  map[string]string  "Shop not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /shops/{id} [get]
func (h *ShopHandler) GetShop(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid shop id")
	}

	result, err := h.getShopUsecase.Execute(id)
	if err != nil {
		if err.Error() == "shop not found" {
			return fiber.NewError(fiber.StatusNotFound, "shop not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get shop")
	}

	shopResponse := &ShopResponseDTO{
		ID:          result.Shop.ID,
		Name:        result.Shop.Name,
		Description: result.Shop.Description,
		Address:     result.Shop.Address,
		Phone:       result.Shop.Phone,
		Email:       result.Shop.Email,
		Website:     result.Shop.Website,
		Logo:        result.Shop.Logo,
		CreatedAt:   result.Shop.CreatedAt,
		UpdatedAt:   result.Shop.UpdatedAt,
	}

	return c.JSON(ShopResponse{
		Message: "shop retrieved successfully.",
		Data: ShopResponseData{
			Shop: shopResponse,
		},
	})
}

// CreateShop godoc
// @Summary      Create a new shop
// @Description  Create a new shop with the provided information
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        request  body      CreateShopPayload  true  "Shop data"
// @Success      201      {object}  ShopResponse
// @Failure      400      {object}  map[string]string  "Invalid request body"
// @Failure      422      {object}  map[string]string  "Validation failed"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /shops [post]
func (h *ShopHandler) CreateShop(c fiber.Ctx) error {
	var request CreateShopPayload

	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	result, err := h.createShopUsecase.Execute(usecases.CreateShopParam{
		Name:        request.Name,
		Description: request.Description,
		Address:     request.Address,
		Phone:       request.Phone,
		Email:       request.Email,
		Website:     request.Website,
		Logo:        request.Logo,
	})

	if err != nil {
		if isValidationError(err) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create shop")
	}

	shopResponse := &ShopResponseDTO{
		ID:          result.Shop.ID,
		Name:        result.Shop.Name,
		Description: result.Shop.Description,
		Address:     result.Shop.Address,
		Phone:       result.Shop.Phone,
		Email:       result.Shop.Email,
		Website:     result.Shop.Website,
		Logo:        result.Shop.Logo,
		CreatedAt:   result.Shop.CreatedAt,
		UpdatedAt:   result.Shop.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(ShopResponse{
		Message: "shop created successfully.",
		Data: ShopResponseData{
			Shop: shopResponse,
		},
	})
}

// UpdateShop godoc
// @Summary      Update shop
// @Description  Update an existing shop by ID
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        id       path      int                true  "Shop ID"
// @Param        request  body      UpdateShopPayload  true  "Shop data"
// @Success      200      {object}  ShopResponse
// @Failure      400      {object}  map[string]string  "Invalid shop id or request body"
// @Failure      404      {object}  map[string]string  "Shop not found"
// @Failure      422      {object}  map[string]string  "Validation failed"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /shops/{id} [put]
func (h *ShopHandler) UpdateShop(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid shop id")
	}

	var request UpdateShopPayload

	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	result, err := h.updateShopUsecase.Execute(usecases.UpdateShopParam{
		ID:          id,
		Name:        request.Name,
		Description: request.Description,
		Address:     request.Address,
		Phone:       request.Phone,
		Email:       request.Email,
		Website:     request.Website,
		Logo:        request.Logo,
	})

	if err != nil {
		if isValidationError(err) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		if err.Error() == "shop not found" {
			return fiber.NewError(fiber.StatusNotFound, "shop not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to update shop")
	}

	shopResponse := &ShopResponseDTO{
		ID:          result.Shop.ID,
		Name:        result.Shop.Name,
		Description: result.Shop.Description,
		Address:     result.Shop.Address,
		Phone:       result.Shop.Phone,
		Email:       result.Shop.Email,
		Website:     result.Shop.Website,
		Logo:        result.Shop.Logo,
		CreatedAt:   result.Shop.CreatedAt,
		UpdatedAt:   result.Shop.UpdatedAt,
	}

	return c.JSON(ShopResponse{
		Message: "shop updated successfully.",
		Data: ShopResponseData{
			Shop: shopResponse,
		},
	})
}

// DeleteShop godoc
// @Summary      Delete shop
// @Description  Soft delete a shop by ID
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Shop ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string  "Invalid shop id"
// @Failure      404  {object}  map[string]string  "Shop not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /shops/{id} [delete]
func (h *ShopHandler) DeleteShop(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid shop id")
	}

	err = h.deleteShopUsecase.Execute(id)
	if err != nil {
		if err.Error() == "shop not found" {
			return fiber.NewError(fiber.StatusNotFound, "shop not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to delete shop")
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// GetStaff godoc
// @Summary      Get staff for a shop
// @Description  Get all staff members assigned to a specific shop
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Shop ID"
// @Success      200  {object}  StaffListResponse
// @Failure      400  {object}  map[string]string  "Invalid shop id"
// @Failure      404  {object}  map[string]string  "Shop not found"
// @Failure      500  {object}  map[string]string  "Internal server error"
// @Router       /shops/{id}/staffs [get]
func (h *ShopHandler) GetStaff(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid shop id")
	}

	// Get shop first to validate it exists
	shopResult, err := h.getShopUsecase.Execute(id)
	if err != nil {
		if err.Error() == "shop not found" {
			return fiber.NewError(fiber.StatusNotFound, "shop not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get shop")
	}

	// Get staffs for the shop
	result, err := h.getStaffsUsecase.Execute(accessusecases.GetStaffsParam{
		Shop: shopResult.Shop,
	})
	if err != nil {
		if isValidationError(err) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get staff")
	}

	staffs := make([]StaffResponseDTO, len(result.Staffs))
	for i, staffInfo := range result.Staffs {
		staffs[i] = StaffResponseDTO{
			User: UserDTO{
				ID:       staffInfo.User.ID,
				FullName: staffInfo.User.FullName,
				Phone:    staffInfo.User.Phone,
				Email:    staffInfo.User.Email,
			},
			Role: RoleDTO{
				ID:          staffInfo.Role.ID,
				Name:        staffInfo.Role.Name,
				Description: staffInfo.Role.Description,
			},
		}
	}

	return c.JSON(StaffListResponse{
		Message: "staff retrieved successfully.",
		Data: StaffListResponseData{
			Staffs: staffs,
		},
	})
}

// AssignStaff godoc
// @Summary      Assign staff to a shop
// @Description  Assign a user with a role to a shop
// @Tags         shops
// @Accept       json
// @Produce      json
// @Param        id       path      int                true  "Shop ID"
// @Param        request  body      AssignStaffPayload  true  "Staff assignment data"
// @Success      201      {object}  StaffResponse
// @Failure      400      {object}  map[string]string  "Invalid shop id or request body"
// @Failure      404      {object}  map[string]string  "Shop, user, or role not found"
// @Failure      409      {object}  map[string]string  "Staff already assigned"
// @Failure      422      {object}  map[string]string  "Validation failed"
// @Failure      500      {object}  map[string]string  "Internal server error"
// @Router       /shops/{id}/staffs [post]
func (h *ShopHandler) AssignStaff(c fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid shop id")
	}

	var request AssignStaffPayload
	if err := c.Bind().Body(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	// Get shop first to validate it exists
	shopResult, err := h.getShopUsecase.Execute(id)
	if err != nil {
		if err.Error() == "shop not found" {
			return fiber.NewError(fiber.StatusNotFound, "shop not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get shop")
	}

	// Get user
	userResult, err := h.getUserUsecase.Execute(request.UserID)
	if err != nil {
		if err.Error() == "user not found" {
			return fiber.NewError(fiber.StatusNotFound, "user not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get user")
	}

	// Get role
	roleResult, err := h.getRoleUsecase.Execute(request.RoleID)
	if err != nil {
		if err.Error() == "role not found" {
			return fiber.NewError(fiber.StatusNotFound, "role not found")
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get role")
	}

	// Validate role belongs to the shop
	if roleResult.Role.ShopID != id {
		return fiber.NewError(fiber.StatusBadRequest, "role does not belong to this shop")
	}

	// Assign staff
	result, err := h.assignStaffUsecase.Execute(accessusecases.AssignStaffParam{
		User: userResult.User,
		Shop: shopResult.Shop,
		Role: roleResult.Role,
	})
	if err != nil {
		if isValidationError(err) {
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		}
		if strings.Contains(err.Error(), "already assigned") {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, "failed to assign staff")
	}

	staffResponse := StaffResponseDTO{
		User: UserDTO{
			ID:       result.Staff.User.ID,
			FullName: result.Staff.User.FullName,
			Phone:    result.Staff.User.Phone,
			Email:    result.Staff.User.Email,
		},
		Role: RoleDTO{
			ID:          result.Staff.Role.ID,
			Name:        result.Staff.Role.Name,
			Description: result.Staff.Role.Description,
		},
	}

	return c.Status(fiber.StatusCreated).JSON(StaffResponse{
		Message: "staff assigned successfully.",
		Data: StaffResponseData{
			Staff: staffResponse,
		},
	})
}

type CreateShopPayload struct {
	Name        string `json:"name" example:"My Shop" binding:"required"`                                // Shop name
	Description string `json:"description" example:"A great shop for all your needs" binding:"required"` // Shop description
	Address     string `json:"address" example:"123 Main St, City, Country" binding:"required"`          // Shop address
	Phone       string `json:"phone" example:"1234567890" binding:"required"`                            // Shop phone number
	Email       string `json:"email" example:"shop@example.com" binding:"required,email"`                // Shop email address
	Website     string `json:"website" example:"https://myshop.com" binding:"required,url"`              // Shop website URL
	Logo        string `json:"logo" example:"logo.png" binding:"required"`                               // Shop logo filename or URL
}

type UpdateShopPayload struct {
	Name        string `json:"name" example:"My Updated Shop" binding:"required"`                    // Shop name
	Description string `json:"description" example:"An updated shop description" binding:"required"` // Shop description
	Address     string `json:"address" example:"456 Oak Ave, City, Country" binding:"required"`      // Shop address
	Phone       string `json:"phone" example:"0987654321" binding:"required"`                        // Shop phone number
	Email       string `json:"email" example:"updated@example.com" binding:"required,email"`         // Shop email address
	Website     string `json:"website" example:"https://updatedshop.com" binding:"required,url"`     // Shop website URL
	Logo        string `json:"logo" example:"new-logo.png" binding:"required"`                       // Shop logo filename or URL
}

type ShopListResponse struct {
	Message string               `json:"message"`
	Data    ShopListResponseData `json:"data"`
}

type ShopListResponseData struct {
	Shops []ShopResponseDTO `json:"shops"`
}

type ShopResponse struct {
	Message string           `json:"message"`
	Data    ShopResponseData `json:"data"`
}

type ShopResponseData struct {
	Shop *ShopResponseDTO `json:"shop"`
}

type ShopResponseDTO struct {
	ID          uint64    `json:"id" example:"1"`
	Name        string    `json:"name" example:"My Shop"`
	Description string    `json:"description" example:"A great shop"`
	Address     string    `json:"address" example:"123 Main St"`
	Phone       string    `json:"phone" example:"1234567890"`
	Email       string    `json:"email" example:"shop@example.com"`
	Website     string    `json:"website" example:"https://myshop.com"`
	Logo        string    `json:"logo" example:"logo.png"`
	CreatedAt   time.Time `json:"created_at" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updated_at" example:"2024-01-01T00:00:00Z"`
}

type AssignStaffPayload struct {
	UserID uint64 `json:"user_id" example:"1" binding:"required"` // User ID to assign
	RoleID uint64 `json:"role_id" example:"1" binding:"required"` // Role ID to assign
}

type UserDTO struct {
	ID       uint64 `json:"id" example:"1"`
	FullName string `json:"full_name" example:"John Doe"`
	Phone    string `json:"phone" example:"1234567890"`
	Email    string `json:"email" example:"john@example.com"`
}

type RoleDTO struct {
	ID          uint64 `json:"id" example:"1"`
	Name        string `json:"name" example:"Manager"`
	Description string `json:"description" example:"Manager role"`
}

type StaffResponseDTO struct {
	User UserDTO `json:"user"`
	Role RoleDTO `json:"role"`
}

type StaffListResponse struct {
	Message string                `json:"message"`
	Data    StaffListResponseData `json:"data"`
}

type StaffListResponseData struct {
	Staffs []StaffResponseDTO `json:"staffs"`
}

type StaffResponse struct {
	Message string            `json:"message"`
	Data    StaffResponseData `json:"data"`
}

type StaffResponseData struct {
	Staff StaffResponseDTO `json:"staff"`
}
