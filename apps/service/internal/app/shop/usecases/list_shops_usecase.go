package usecases

import (
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

func (u *ListShopsUsecase) Execute() *ListShopsResult {
	shops := u.shopRepository.All()
	return &ListShopsResult{
		Shops: shops,
	}
}
