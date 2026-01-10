package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func setupListShopsTest(t *testing.T) (*ListShopsUsecase, repositories.ShopRepository) {
	db := testutil.SetupTestDB(t, &entities.Shop{})
	shopRepo := repositories.NewShopRepository(db)
	usecase := NewListShopsUsecase(shopRepo)
	return usecase, shopRepo
}

func TestListShopsUsecase_Execute(t *testing.T) {
	t.Run("returns empty list when no shops exist", func(t *testing.T) {
		ctx := context.Background()
		usecase, _ := setupListShopsTest(t)

		result := usecase.Execute(ctx)
		assert.NotNil(t, result)
		assert.Empty(t, result.Shops)
	})

	t.Run("returns all shops", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo := setupListShopsTest(t)

		shop1 := entities.Shop{
			Name:        "Shop One",
			Description: "First shop description",
			Address:     "123 Main St",
			Phone:       "1234567890",
			Email:       "shop1@example.com",
			Website:     "https://shop1.com",
			Logo:        "logo1.png",
		}
		shop2 := entities.Shop{
			Name:        "Shop Two",
			Description: "Second shop description",
			Address:     "456 Oak Ave",
			Phone:       "0987654321",
			Email:       "shop2@example.com",
			Website:     "https://shop2.com",
			Logo:        "logo2.png",
		}

		_, err := shopRepo.Create(ctx, shop1)
		require.NoError(t, err)
		_, err = shopRepo.Create(ctx, shop2)
		require.NoError(t, err)

		result := usecase.Execute(ctx)
		assert.NotNil(t, result)
		assert.Len(t, result.Shops, 2)
		assert.Equal(t, "Shop One", result.Shops[0].Name)
		assert.Equal(t, "Shop Two", result.Shops[1].Name)
	})

	t.Run("excludes soft deleted shops", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo := setupListShopsTest(t)

		shop := entities.Shop{
			Name:        "Deleted Shop",
			Description: "Deleted shop description",
			Address:     "789 Elm St",
			Phone:       "1111111111",
			Email:       "deleted@example.com",
			Website:     "https://deleted.com",
			Logo:        "deleted.png",
		}

		created, err := shopRepo.Create(ctx, shop)
		require.NoError(t, err)

		err = shopRepo.Delete(ctx, created)
		require.NoError(t, err)

		result := usecase.Execute(ctx)
		assert.NotNil(t, result)
		assert.Empty(t, result.Shops)
	})
}
