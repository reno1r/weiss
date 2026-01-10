package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
)

type GetShopUsecase struct {
	shopRepository repositories.ShopRepository
}

func NewGetShopUsecase(shopRepository repositories.ShopRepository) *GetShopUsecase {
	return &GetShopUsecase{
		shopRepository: shopRepository,
	}
}

type GetShopResult struct {
	Shop *entities.Shop
}

func (u *GetShopUsecase) Execute(ctx context.Context, id uint64) (*GetShopResult, error) {
	shop, err := u.shopRepository.FindByID(ctx, id)
	if err != nil {
		if err.Error() == "shop not found" {
			return nil, errors.New("shop not found")
		}
		return nil, fmt.Errorf("failed to get shop: %w", err)
	}

	return &GetShopResult{
		Shop: &shop,
	}, nil
}

