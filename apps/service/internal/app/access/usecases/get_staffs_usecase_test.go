package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	accessentities "github.com/reno1r/weiss/apps/service/internal/app/access/entities"
	accessrepositories "github.com/reno1r/weiss/apps/service/internal/app/access/repositories"
	shopentities "github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	shoprepositories "github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
	userentities "github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	userrepositories "github.com/reno1r/weiss/apps/service/internal/app/user/repositories"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func setupGetStaffsTest(t *testing.T) (*GetStaffsUsecase, accessrepositories.StaffRepository, shoprepositories.ShopRepository, accessrepositories.RoleRepository, userrepositories.UserRepository) {
	db := testutil.SetupTestDB(t, &shopentities.Shop{}, &accessentities.Role{}, &userentities.User{}, &accessentities.Staff{})
	shopRepo := shoprepositories.NewShopRepository(db)
	roleRepo := accessrepositories.NewRoleRepository(db)
	userRepo := userrepositories.NewUserRepository(db)
	staffRepo := accessrepositories.NewStaffRepository(db)
	usecase := NewGetStaffsUsecase(staffRepo)
	return usecase, staffRepo, shopRepo, roleRepo, userRepo
}

func createTestShop(t *testing.T, ctx context.Context, shopRepo shoprepositories.ShopRepository) shopentities.Shop {
	shop := shopentities.Shop{
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
	return created
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

func createTestRole(t *testing.T, ctx context.Context, roleRepo accessrepositories.RoleRepository, shopID uint64) accessentities.Role {
	role := accessentities.Role{
		Name:        "Test Role",
		Description: "Test role description",
		ShopID:      shopID,
	}
	created, err := roleRepo.Create(ctx, role)
	require.NoError(t, err)
	return created
}

func TestGetStaffsUsecase_Execute(t *testing.T) {
	t.Run("returns empty list when no staffs exist", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, shopRepo, _, _ := setupGetStaffsTest(t)
		shop := createTestShop(t, ctx, shopRepo)

		result, err := usecase.Execute(ctx, GetStaffsParam{
			Shop: &shop,
		})
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, shop.ID, result.Shop.ID)
		assert.Empty(t, result.Staffs)
	})

	t.Run("returns all staffs for a shop", func(t *testing.T) {
		ctx := context.Background()
		usecase, staffRepo, shopRepo, roleRepo, userRepo := setupGetStaffsTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user1 := createTestUser(t, ctx, userRepo)
		user2 := userentities.User{
			FullName: "Another User",
			Phone:    "0987654321",
			Email:    "another@example.com",
			Password: "hashedpassword",
		}
		createdUser2, err := userRepo.Create(ctx, user2)
		require.NoError(t, err)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff1 := accessentities.Staff{
			UserID: user1.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}
		staff2 := accessentities.Staff{
			UserID: createdUser2.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		_, err = staffRepo.Create(ctx, staff1)
		require.NoError(t, err)
		_, err = staffRepo.Create(ctx, staff2)
		require.NoError(t, err)

		result, err := usecase.Execute(ctx, GetStaffsParam{
			Shop: &shop,
		})
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, shop.ID, result.Shop.ID)
		assert.Len(t, result.Staffs, 2)
		assert.NotNil(t, result.Staffs[0].User)
		assert.NotNil(t, result.Staffs[0].Role)
		assert.NotNil(t, result.Staffs[1].User)
		assert.NotNil(t, result.Staffs[1].Role)
	})

	t.Run("returns only staffs for the specified shop", func(t *testing.T) {
		ctx := context.Background()
		usecase, staffRepo, shopRepo, roleRepo, userRepo := setupGetStaffsTest(t)
		shop1 := createTestShop(t, ctx, shopRepo)
		shop2 := shopentities.Shop{
			Name:        "Another Shop",
			Description: "Another shop description",
			Address:     "456 Another St",
			Phone:       "0987654321",
			Email:       "another@example.com",
			Website:     "https://another.com",
			Logo:        "another.png",
		}
		createdShop2, err := shopRepo.Create(ctx, shop2)
		require.NoError(t, err)

		user := createTestUser(t, ctx, userRepo)
		role1 := createTestRole(t, ctx, roleRepo, shop1.ID)
		role2 := createTestRole(t, ctx, roleRepo, createdShop2.ID)

		staff1 := accessentities.Staff{
			UserID: user.ID,
			RoleID: role1.ID,
			ShopID: shop1.ID,
		}
		staff2 := accessentities.Staff{
			UserID: user.ID,
			RoleID: role2.ID,
			ShopID: createdShop2.ID,
		}

		_, err = staffRepo.Create(ctx, staff1)
		require.NoError(t, err)
		_, err = staffRepo.Create(ctx, staff2)
		require.NoError(t, err)

		result, err := usecase.Execute(ctx, GetStaffsParam{
			Shop: &shop1,
		})
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, shop1.ID, result.Shop.ID)
		assert.Len(t, result.Staffs, 1)
		assert.Equal(t, shop1.ID, result.Staffs[0].Role.ShopID)
	})

	t.Run("returns validation error when shop is nil", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, _ := setupGetStaffsTest(t)

		result, err := usecase.Execute(ctx, GetStaffsParam{
			Shop: nil,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("excludes soft deleted staffs", func(t *testing.T) {
		ctx := context.Background()
		usecase, staffRepo, shopRepo, roleRepo, userRepo := setupGetStaffsTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff := accessentities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created, err := staffRepo.Create(ctx, staff)
		require.NoError(t, err)

		err = staffRepo.Delete(ctx, created)
		require.NoError(t, err)

		result, err := usecase.Execute(ctx, GetStaffsParam{
			Shop: &shop,
		})
		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Empty(t, result.Staffs)
	})
}
