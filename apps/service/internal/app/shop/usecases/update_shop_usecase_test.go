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

func setupUpdateShopTest(t *testing.T) (*UpdateShopUsecase, repositories.ShopRepository) {
	db := testutil.SetupTestDB(t, &entities.Shop{})
	shopRepo := repositories.NewShopRepository(db)
	usecase := NewUpdateShopUsecase(shopRepo)
	return usecase, shopRepo
}

func TestUpdateShopUsecase_Execute(t *testing.T) {
	t.Run("updates shop successfully", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo := setupUpdateShopTest(t)

		// Create initial shop
		shop := entities.Shop{
			Name:        "Original Shop",
			Description: "Original description",
			Address:     "123 Original St",
			Phone:       "1234567890",
			Email:       "original@example.com",
			Website:     "https://original.com",
			Logo:        "original.png",
		}

		created, err := shopRepo.Create(ctx, shop)
		require.NoError(t, err)

		param := UpdateShopParam{
			ID:          created.ID,
			Name:        "Updated Shop",
			Description: "Updated description for the shop",
			Address:     "456 Updated Ave",
			Phone:       "0987654321",
			Email:       "updated@example.com",
			Website:     "https://updated.com",
			Logo:        "updated.png",
		}

		result, err := usecase.Execute(ctx, param)
		require.NoError(t, err)
		assert.NotNil(t, result.Shop)
		assert.Equal(t, created.ID, result.Shop.ID)
		assert.Equal(t, "Updated Shop", result.Shop.Name)
		assert.Equal(t, "Updated description for the shop", result.Shop.Description)
		assert.Equal(t, "456 Updated Ave", result.Shop.Address)
		assert.Equal(t, "0987654321", result.Shop.Phone)
		assert.Equal(t, "updated@example.com", result.Shop.Email)
		assert.Equal(t, "https://updated.com", result.Shop.Website)
		assert.Equal(t, "updated.png", result.Shop.Logo)

		// Verify shop was updated in database
		found, err := shopRepo.FindByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Shop", found.Name)
		assert.Equal(t, "updated@example.com", found.Email)
	})

	t.Run("returns error when shop not found", func(t *testing.T) {
		ctx := context.Background()
		usecase, _ := setupUpdateShopTest(t)

		param := UpdateShopParam{
			ID:          999,
			Name:        "Updated Shop",
			Description: "Updated description for the shop",
			Address:     "456 Updated Ave",
			Phone:       "0987654321",
			Email:       "updated@example.com",
			Website:     "https://updated.com",
			Logo:        "updated.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
		assert.Nil(t, result)
	})

	t.Run("validates ID is required", func(t *testing.T) {
		ctx := context.Background()
		usecase, _ := setupUpdateShopTest(t)

		param := UpdateShopParam{
			ID:          0,
			Name:        "Updated Shop",
			Description: "Updated description for the shop",
			Address:     "456 Updated Ave",
			Phone:       "0987654321",
			Email:       "updated@example.com",
			Website:     "https://updated.com",
			Logo:        "updated.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates name is required", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo := setupUpdateShopTest(t)

		shop := entities.Shop{
			Name:        "Original Shop",
			Description: "Original description",
			Address:     "123 Original St",
			Phone:       "1234567890",
			Email:       "original@example.com",
			Website:     "https://original.com",
			Logo:        "original.png",
		}

		created, err := shopRepo.Create(ctx, shop)
		require.NoError(t, err)

		param := UpdateShopParam{
			ID:          created.ID,
			Name:        "",
			Description: "Updated description for the shop",
			Address:     "456 Updated Ave",
			Phone:       "0987654321",
			Email:       "updated@example.com",
			Website:     "https://updated.com",
			Logo:        "updated.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates name minimum length", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo := setupUpdateShopTest(t)

		shop := entities.Shop{
			Name:        "Original Shop",
			Description: "Original description",
			Address:     "123 Original St",
			Phone:       "1234567890",
			Email:       "original@example.com",
			Website:     "https://original.com",
			Logo:        "original.png",
		}

		created, err := shopRepo.Create(ctx, shop)
		require.NoError(t, err)

		param := UpdateShopParam{
			ID:          created.ID,
			Name:        "A",
			Description: "Updated description for the shop",
			Address:     "456 Updated Ave",
			Phone:       "0987654321",
			Email:       "updated@example.com",
			Website:     "https://updated.com",
			Logo:        "updated.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates description is required", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo := setupUpdateShopTest(t)

		shop := entities.Shop{
			Name:        "Original Shop",
			Description: "Original description",
			Address:     "123 Original St",
			Phone:       "1234567890",
			Email:       "original@example.com",
			Website:     "https://original.com",
			Logo:        "original.png",
		}

		created, err := shopRepo.Create(ctx, shop)
		require.NoError(t, err)

		param := UpdateShopParam{
			ID:          created.ID,
			Name:        "Updated Shop",
			Description: "",
			Address:     "456 Updated Ave",
			Phone:       "0987654321",
			Email:       "updated@example.com",
			Website:     "https://updated.com",
			Logo:        "updated.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates email format", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo := setupUpdateShopTest(t)

		shop := entities.Shop{
			Name:        "Original Shop",
			Description: "Original description",
			Address:     "123 Original St",
			Phone:       "1234567890",
			Email:       "original@example.com",
			Website:     "https://original.com",
			Logo:        "original.png",
		}

		created, err := shopRepo.Create(ctx, shop)
		require.NoError(t, err)

		param := UpdateShopParam{
			ID:          created.ID,
			Name:        "Updated Shop",
			Description: "Updated description for the shop",
			Address:     "456 Updated Ave",
			Phone:       "0987654321",
			Email:       "invalid-email",
			Website:     "https://updated.com",
			Logo:        "updated.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates website URL format", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo := setupUpdateShopTest(t)

		shop := entities.Shop{
			Name:        "Original Shop",
			Description: "Original description",
			Address:     "123 Original St",
			Phone:       "1234567890",
			Email:       "original@example.com",
			Website:     "https://original.com",
			Logo:        "original.png",
		}

		created, err := shopRepo.Create(ctx, shop)
		require.NoError(t, err)

		param := UpdateShopParam{
			ID:          created.ID,
			Name:        "Updated Shop",
			Description: "Updated description for the shop",
			Address:     "456 Updated Ave",
			Phone:       "0987654321",
			Email:       "updated@example.com",
			Website:     "not-a-url",
			Logo:        "updated.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})
}
