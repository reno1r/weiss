package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	accessentities "github.com/reno1r/weiss/apps/service/internal/app/access/entities"
	accessrepositories "github.com/reno1r/weiss/apps/service/internal/app/access/repositories"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
	userentities "github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	userrepositories "github.com/reno1r/weiss/apps/service/internal/app/user/repositories"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func setupCreateShopTest(t *testing.T) (*CreateShopUsecase, repositories.ShopRepository, accessrepositories.RoleRepository, accessrepositories.StaffRepository, userrepositories.UserRepository) {
	db := testutil.SetupTestDB(t, &entities.Shop{}, &accessentities.Role{}, &accessentities.Staff{}, &userentities.User{})
	shopRepo := repositories.NewShopRepository(db)
	roleRepo := accessrepositories.NewRoleRepository(db)
	staffRepo := accessrepositories.NewStaffRepository(db)
	userRepo := userrepositories.NewUserRepository(db)
	usecase := NewCreateShopUsecase(db, shopRepo, roleRepo, staffRepo)
	return usecase, shopRepo, roleRepo, staffRepo, userRepo
}

func createTestUser(t *testing.T, ctx context.Context, userRepo userrepositories.UserRepository) userentities.User {
	user := userentities.User{
		FullName: "Test User",
		Phone:    "1234567890",
		Email:    "test@example.com",
		Password: "hashedpassword",
	}
	created, err := userRepo.Create(ctx, user)
	require.NoError(t, err)
	return created
}

func TestCreateShopUsecase_Execute(t *testing.T) {
	t.Run("creates shop successfully", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo, roleRepo, staffRepo, userRepo := setupCreateShopTest(t)

		// Create a test user first
		user := userentities.User{
			FullName: "Test User",
			Phone:    "1234567890",
			Email:    "test@example.com",
			Password: "hashedpassword",
		}
		createdUser, err := userRepo.Create(ctx, user)
		require.NoError(t, err)

		param := CreateShopParam{
			UserID:      createdUser.ID,
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
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
		found, err := shopRepo.FindByID(ctx, result.Shop.ID)
		require.NoError(t, err)
		assert.Equal(t, result.Shop.ID, found.ID)
		assert.Equal(t, "My Shop", found.Name)

		// Verify owner role was created for the shop
		roles := roleRepo.FindByShopID(ctx, result.Shop.ID)
		require.Len(t, roles, 1)
		assert.Equal(t, "Owner", roles[0].Name)
		assert.Equal(t, "Shop owner with full access to manage the shop", roles[0].Description)
		assert.Equal(t, result.Shop.ID, roles[0].ShopID)

		// Verify user was assigned as owner
		staffs := staffRepo.FindByShopID(ctx, result.Shop.ID)
		require.Len(t, staffs, 1)
		assert.Equal(t, param.UserID, staffs[0].UserID)
		assert.Equal(t, result.Shop.ID, staffs[0].ShopID)
		assert.Equal(t, roles[0].ID, staffs[0].RoleID)
	})

	t.Run("rolls back all changes when staff assignment fails", func(t *testing.T) {
		ctx := context.Background()
		usecase, shopRepo, roleRepo, _, userRepo := setupCreateShopTest(t)
		_ = createTestUser(t, ctx, userRepo) // Create a user but use invalid ID

		// Use an invalid user ID that doesn't exist to cause staff assignment to fail
		// This will cause a foreign key constraint error and trigger rollback
		param := CreateShopParam{
			UserID:      999999, // Non-existent user ID - this will cause foreign key constraint error
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		// Count shops and roles before
		shopsBefore := shopRepo.All(ctx)
		rolesBefore := roleRepo.All(ctx)
		countShopsBefore := len(shopsBefore)
		countRolesBefore := len(rolesBefore)

		// Execute - this should fail due to foreign key constraint
		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to assign user as owner")
		assert.Nil(t, result)

		// Verify nothing was created (transaction rolled back)
		shopsAfter := shopRepo.All(ctx)
		rolesAfter := roleRepo.All(ctx)
		assert.Len(t, shopsAfter, countShopsBefore, "shop should not be created when transaction fails")
		assert.Len(t, rolesAfter, countRolesBefore, "role should not be created when transaction fails")
	})

	t.Run("validates user ID is required", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, _ := setupCreateShopTest(t)

		param := CreateShopParam{
			UserID:      0,
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates name is required", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates name minimum length", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "A",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates description is required", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "My Shop",
			Description: "",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates description minimum length", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "My Shop",
			Description: "Short",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates address is required", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates phone is required", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates phone minimum length", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "123",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates email is required", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates email format", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "invalid-email",
			Website:     "https://myshop.com",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates website is required", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates website URL format", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "not-a-url",
			Logo:        "logo.png",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("validates logo is required", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupCreateShopTest(t)
		user := createTestUser(t, ctx, userRepo)

		param := CreateShopParam{
			UserID:      user.ID,
			Name:        "My Shop",
			Description: "A great shop for all your needs",
			Address:     "123 Main St, City, Country",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://myshop.com",
			Logo:        "",
		}

		result, err := usecase.Execute(ctx, param)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})
}
