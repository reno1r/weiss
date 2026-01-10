package usecases

import (
	"context"

	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
)

type ListShopsUsecase struct {
	shopRepository repositories.ShopRepository
}

func NewListShopsUsecase(shopRepository repositories.ShopRepository) *ListShopsUsecase {
	return &ListShopsUsecase{
		shopRepository: shopRepository,
	}
}

type ListShopsResult struct {
	Shops []entities.Shop
}

func (u *ListShopsUsecase) Execute(ctx context.Context) *ListShopsResult {
	shops := u.shopRepository.All(ctx)
	return &ListShopsResult{
		Shops: shops,
	}
}
