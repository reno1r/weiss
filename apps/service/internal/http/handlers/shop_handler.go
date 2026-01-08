package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
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

func (h *ShopHandler) ListShops(c fiber.Ctx) error {
	result := h.listShopsUsecase.Execute()

	return c.JSON(ShopListResponse{
		Message: "shops retrieved successfully.",
		Data: ShopListResponseData{
			Shops: result.Shops,
		},
	})
}

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

	return c.JSON(ShopResponse{
		Message: "shop retrieved successfully.",
		Data: ShopResponseData{
			Shop: result.Shop,
		},
	})
}

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

	return c.Status(fiber.StatusCreated).JSON(ShopResponse{
		Message: "shop created successfully.",
		Data: ShopResponseData{
			Shop: result.Shop,
		},
	})
}

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

	return c.JSON(ShopResponse{
		Message: "shop updated successfully.",
		Data: ShopResponseData{
			Shop: result.Shop,
		},
	})
}

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
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	Logo        string `json:"logo"`
}

type UpdateShopPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Website     string `json:"website"`
	Logo        string `json:"logo"`
}

type ShopListResponse struct {
	Message string               `json:"message"`
	Data    ShopListResponseData `json:"data"`
}

type ShopListResponseData struct {
	Shops []entities.Shop `json:"shops"`
}

type ShopResponse struct {
	Message string           `json:"message"`
	Data    ShopResponseData `json:"data"`
}

type ShopResponseData struct {
	Shop *entities.Shop `json:"shop"`
}
