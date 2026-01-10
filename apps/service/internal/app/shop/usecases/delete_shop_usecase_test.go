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

func setupDeleteShopTest(t *testing.T) (*DeleteShopUsecase, repositories.ShopRepository) {
	db := testutil.SetupTestDB(t, &entities.Shop{})
	shopRepo := repositories.NewShopRepository(db)
	usecase := NewDeleteShopUsecase(shopRepo)
	return usecase, shopRepo
}

func TestDeleteShopUsecase_Execute(t *testing.T) {
	t.Run("deletes shop successfully", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo := setupDeleteShopTest(t)

		shop := entities.Shop{
			Name:        "Test Shop",
			Description: "Test shop description",
			Address:     "123 Test St",
			Phone:       "1234567890",
			Email:       "test@example.com",
			Website:     "https://test.com",
			Logo:        "test.png",
		}

		created, err := shopRepo.Create(ctx, shop)
		require.NoError(t, err)

		err = usecase.Execute(ctx, created.ID)
		require.NoError(t, err)

		// Verify shop was soft deleted
		_, err = shopRepo.FindByID(ctx, created.ID)
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
	})

	t.Run("returns error when shop not found", func(t *testing.T) {
		ctx := context.Background()
		usecase, _ := setupDeleteShopTest(t)

		err := usecase.Execute(ctx, 999)
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
	})

	t.Run("can delete multiple shops", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo := setupDeleteShopTest(t)

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

		created1, err := shopRepo.Create(ctx, shop1)
		require.NoError(t, err)
		created2, err := shopRepo.Create(ctx, shop2)
		require.NoError(t, err)

		err = usecase.Execute(ctx, created1.ID)
		require.NoError(t, err)

		err = usecase.Execute(ctx, created2.ID)
		require.NoError(t, err)

		// Verify both shops were deleted
		_, err = shopRepo.FindByID(ctx, created1.ID)
		assert.Error(t, err)
		_, err = shopRepo.FindByID(ctx, created2.ID)
		assert.Error(t, err)

		// Verify no shops remain
		shops := shopRepo.All(ctx)
		assert.Empty(t, shops)
	})

	t.Run("does not delete already deleted shop", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo := setupDeleteShopTest(t)

		shop := entities.Shop{
			Name:        "Test Shop",
			Description: "Test shop description",
			Address:     "123 Test St",
			Phone:       "1234567890",
			Email:       "test@example.com",
			Website:     "https://test.com",
			Logo:        "test.png",
		}

		created, err := shopRepo.Create(ctx, shop)
		require.NoError(t, err)

		// Delete first time
		err = usecase.Execute(ctx, created.ID)
		require.NoError(t, err)

		// Try to delete again
		err = usecase.Execute(ctx, created.ID)
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
	})
}
