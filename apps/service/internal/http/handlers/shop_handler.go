package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/usecases"
)

type ShopHandler struct {
	listShopsUsecase  *usecases.ListShopsUsecase
	getShopUsecase    *usecases.GetShopUsecase
	createShopUsecase *usecases.CreateShopUsecase
	updateShopUsecase *usecases.UpdateShopUsecase
	deleteShopUsecase *usecases.DeleteShopUsecase
}

func NewShopHandler(
	listShopsUsecase *usecases.ListShopsUsecase,
	getShopUsecase *usecases.GetShopUsecase,
	createShopUsecase *usecases.CreateShopUsecase,
	updateShopUsecase *usecases.UpdateShopUsecase,
	deleteShopUsecase *usecases.DeleteShopUsecase,
) *ShopHandler {
	return &ShopHandler{
		listShopsUsecase:  listShopsUsecase,
		getShopUsecase:    getShopUsecase,
		createShopUsecase: createShopUsecase,
		updateShopUsecase: updateShopUsecase,
		deleteShopUsecase: deleteShopUsecase,
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

	result, err := h.createShopUsecase.Execute(usecases.CreateShopData{
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

	result, err := h.updateShopUsecase.Execute(usecases.UpdateShopData{
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
