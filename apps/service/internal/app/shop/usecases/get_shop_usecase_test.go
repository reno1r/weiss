package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func setupGetShopTest(t *testing.T) (*GetShopUsecase, repositories.ShopRepository) {
	db := testutil.SetupTestDB(t, &entities.Shop{})
	shopRepo := repositories.NewShopRepository(db)
	usecase := NewGetShopUsecase(shopRepo)
	return usecase, shopRepo
}

func TestGetShopUsecase_Execute(t *testing.T) {
	t.Run("returns shop when found", func(t *testing.T) {
		usecase, shopRepo := setupGetShopTest(t)

		shop := entities.Shop{
			Name:        "Test Shop",
			Description: "Test shop description",
			Address:     "123 Test St",
			Phone:       "1234567890",
			Email:       "test@example.com",
			Website:     "https://test.com",
			Logo:        "test.png",
		}

		created, err := shopRepo.Create(shop)
		require.NoError(t, err)

		result, err := usecase.Execute(created.ID)
		require.NoError(t, err)
		assert.NotNil(t, result.Shop)
		assert.Equal(t, created.ID, result.Shop.ID)
		assert.Equal(t, "Test Shop", result.Shop.Name)
		assert.Equal(t, "Test shop description", result.Shop.Description)
		assert.Equal(t, "123 Test St", result.Shop.Address)
		assert.Equal(t, "1234567890", result.Shop.Phone)
		assert.Equal(t, "test@example.com", result.Shop.Email)
		assert.Equal(t, "https://test.com", result.Shop.Website)
		assert.Equal(t, "test.png", result.Shop.Logo)
	})

	t.Run("returns error when shop not found", func(t *testing.T) {
		usecase, _ := setupGetShopTest(t)

		result, err := usecase.Execute(999)
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
		assert.Nil(t, result)
	})

	t.Run("does not find soft deleted shops", func(t *testing.T) {
		usecase, shopRepo := setupGetShopTest(t)

		shop := entities.Shop{
			Name:        "Deleted Shop",
			Description: "Deleted shop description",
			Address:     "789 Elm St",
			Phone:       "1111111111",
			Email:       "deleted@example.com",
			Website:     "https://deleted.com",
			Logo:        "deleted.png",
		}

		created, err := shopRepo.Create(shop)
		require.NoError(t, err)

		err = shopRepo.Delete(created)
		require.NoError(t, err)

		result, err := usecase.Execute(created.ID)
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
		assert.Nil(t, result)
	})
}
