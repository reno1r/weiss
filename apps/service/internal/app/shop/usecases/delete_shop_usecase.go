package usecases

import (
	"errors"
	"fmt"

	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
)

type DeleteShopUsecase struct {
	shopRepository repositories.ShopRepository
}

func NewDeleteShopUsecase(shopRepository repositories.ShopRepository) *DeleteShopUsecase {
	return &DeleteShopUsecase{
		shopRepository: shopRepository,
	}
}

func (u *DeleteShopUsecase) Execute(id uint64) error {
	shop, err := u.shopRepository.FindByID(id)
	if err != nil {
		return errors.New("shop not found")
	}

	err = u.shopRepository.Delete(shop)
	if err != nil {
		return fmt.Errorf("failed to delete shop: %w", err)
	}

	return nil
}
