package usecases

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func setupCreateShopTest(t *testing.T) (*CreateShopUsecase, repositories.ShopRepository) {
	db := testutil.SetupTestDB(t, &entities.Shop{})
	shopRepo := repositories.NewShopRepository(db)
	usecase := NewCreateShopUsecase(shopRepo)
	return usecase, shopRepo
}

func TestCreateShopUsecase_Execute(t *testing.T) {
	t.Run("creates shop successfully", func(t *testing.T) {
		usecase, shopRepo := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		require.NoError(t, err)
		assert.NotNil(t, result.Shop)
		assert.NotZero(t, result.Shop.ID)
		assert.Equal(t, "My Shop", result.Shop.Name)
		assert.Equal(t, "A great shop for all your needs", result.Shop.Description)
		assert.Equal(t, "123 Main St, City, Country", result.Shop.Address)
		assert.Equal(t, "1234567890", result.Shop.Phone)
		assert.Equal(t, "shop@example.com", result.Shop.Email)
		assert.Equal(t, "https://myshop.com", result.Shop.Website)
		assert.Equal(t, "logo.png", result.Shop.Logo)

		// Verify shop was created in database
		found, err := shopRepo.FindByID(result.Shop.ID)
		require.NoError(t, err)
		assert.Equal(t, result.Shop.ID, found.ID)
		assert.Equal(t, "My Shop", found.Name)
	})

	t.Run("validates name is required", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates name minimum length", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "A",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates description is required", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "My Shop",
			Description: "",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates description minimum length", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "My Shop",
			Description: "Short",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates address is required", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates phone is required", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates phone minimum length", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "123",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates email is required", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates email format", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "invalid-email",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates website is required", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates website URL format", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "not-a-url",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates logo is required", func(t *testing.T) {
		usecase, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "",
		}

		result, err := usecase.Execute(param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})
}

