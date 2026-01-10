package usecases

import (
	"context"
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

func (u *DeleteShopUsecase) Execute(ctx context.Context, id uint64) error {
	shop, err := u.shopRepository.FindByID(ctx, id)
	if err != nil {
		return errors.New("shop not found")
	}

	err = u.shopRepository.Delete(ctx, shop)
	if err != nil {
		return fmt.Errorf("failed to delete shop: %w", err)
	}

	return nil
}
