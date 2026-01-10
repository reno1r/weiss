package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reno1r/weiss/apps/service/internal/app/access/entities"
	shopentities "github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
	userentities "github.com/reno1r/weiss/apps/service/internal/app/user/entities"
	userrepositories "github.com/reno1r/weiss/apps/service/internal/app/user/repositories"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func setupStaffTest(t *testing.T) (repositories.ShopRepository, RoleRepository, userrepositories.UserRepository, StaffRepository) {
	db := testutil.SetupTestDB(t, &shopentities.Shop{}, &entities.Role{}, &userentities.User{}, &entities.Staff{})
	shopRepo := repositories.NewShopRepository(db)
	roleRepo := NewRoleRepository(db)
	userRepo := userrepositories.NewUserRepository(db)
	staffRepo := NewStaffRepository(db)
	return shopRepo, roleRepo, userRepo, staffRepo
}

func createTestShop(t *testing.T, ctx context.Context, shopRepo repositories.ShopRepository) shopentities.Shop {
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

func createTestRole(t *testing.T, ctx context.Context, roleRepo RoleRepository, shopID uint64) entities.Role {
	role := entities.Role{
		Name:        "Test Role",
		Description: "Test role description",
		ShopID:      shopID,
	}
	created, err := roleRepo.Create(ctx, role)
	require.NoError(t, err)
	return created
}

func TestStaffRepository_All(t *testing.T) {
	t.Run("returns empty slice when no staffs exist", func(t *testing.T) {
		ctx := context.Background()
		_, _, _, staffRepo := setupStaffTest(t)

		staffs := staffRepo.All(ctx)
		assert.Empty(t, staffs)
	})

	t.Run("returns all staffs", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff1 := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}
		staff2 := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		_, err := staffRepo.Create(ctx, staff1)
		require.NoError(t, err)
		_, err = staffRepo.Create(ctx, staff2)
		require.NoError(t, err)

		staffs := staffRepo.All(ctx)
		assert.Len(t, staffs, 2)
	})

	t.Run("excludes soft deleted staffs", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created, err := staffRepo.Create(ctx, staff)
		require.NoError(t, err)

		err = staffRepo.Delete(ctx, created)
		require.NoError(t, err)

		staffs := staffRepo.All(ctx)
		assert.Empty(t, staffs)
	})
}

func TestStaffRepository_FindByID(t *testing.T) {
	t.Run("returns staff when found", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created, err := staffRepo.Create(ctx, staff)
		require.NoError(t, err)

		found, err := staffRepo.FindByID(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, user.ID, found.UserID)
		assert.Equal(t, role.ID, found.RoleID)
		assert.Equal(t, shop.ID, found.ShopID)
	})

	t.Run("returns error when staff not found", func(t *testing.T) {
		ctx := context.Background()
		_, _, _, staffRepo := setupStaffTest(t)

		_, err := staffRepo.FindByID(ctx, 999)
		assert.Error(t, err)
		assert.Equal(t, "staff not found", err.Error())
	})

	t.Run("does not find soft deleted staffs", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created, err := staffRepo.Create(ctx, staff)
		require.NoError(t, err)

		err = staffRepo.Delete(ctx, created)
		require.NoError(t, err)

		_, err = staffRepo.FindByID(ctx, created.ID)
		assert.Error(t, err)
		assert.Equal(t, "staff not found", err.Error())
	})
}

func TestStaffRepository_FindByShopID(t *testing.T) {
	t.Run("returns staffs for a shop", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop1 := createTestShop(t, ctx, shopRepo)

		shop2 := shopentities.Shop{
			Name:        "Shop Two",
			Description: "Second shop",
			Address:     "456 Oak Ave",
			Phone:       "0987654321",
			Email:       "shop2@example.com",
			Website:     "https://shop2.com",
			Logo:        "logo2.png",
		}
		createdShop2, err := shopRepo.Create(ctx, shop2)
		require.NoError(t, err)

		user1 := createTestUser(t, ctx, userRepo)
		user2 := userentities.User{
			FullName: "User Two",
			Phone:    "0987654321",
			Email:    "user2@example.com",
			Password: "hashedpassword",
		}
		createdUser2, err := userRepo.Create(ctx, user2)
		require.NoError(t, err)

		role1 := createTestRole(t, ctx, roleRepo, shop1.ID)
		role2 := createTestRole(t, ctx, roleRepo, createdShop2.ID)

		staff1 := entities.Staff{
			UserID: user1.ID,
			RoleID: role1.ID,
			ShopID: shop1.ID,
		}
		staff2 := entities.Staff{
			UserID: createdUser2.ID,
			RoleID: role1.ID,
			ShopID: shop1.ID,
		}
		staff3 := entities.Staff{
			UserID: user1.ID,
			RoleID: role2.ID,
			ShopID: createdShop2.ID,
		}

		_, err = staffRepo.Create(ctx, staff1)
		require.NoError(t, err)
		_, err = staffRepo.Create(ctx, staff2)
		require.NoError(t, err)
		_, err = staffRepo.Create(ctx, staff3)
		require.NoError(t, err)

		staffs := staffRepo.FindByShopID(ctx, shop1.ID)
		assert.Len(t, staffs, 2)
	})

	t.Run("returns empty slice when no staffs for shop", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, _, _, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)

		staffs := staffRepo.FindByShopID(ctx, shop.ID)
		assert.Empty(t, staffs)
	})

	t.Run("excludes soft deleted staffs", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff1 := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}
		staff2 := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created1, err := staffRepo.Create(ctx, staff1)
		require.NoError(t, err)
		_, err = staffRepo.Create(ctx, staff2)
		require.NoError(t, err)

		err = staffRepo.Delete(ctx, created1)
		require.NoError(t, err)

		staffs := staffRepo.FindByShopID(ctx, shop.ID)
		assert.Len(t, staffs, 1)
	})
}

func TestStaffRepository_FindByUserID(t *testing.T) {
	t.Run("returns staffs for a user", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop1 := createTestShop(t, ctx, shopRepo)
		shop2 := shopentities.Shop{
			Name:        "Shop Two",
			Description: "Second shop",
			Address:     "456 Oak Ave",
			Phone:       "0987654321",
			Email:       "shop2@example.com",
			Website:     "https://shop2.com",
			Logo:        "logo2.png",
		}
		createdShop2, err := shopRepo.Create(ctx, shop2)
		require.NoError(t, err)

		user := createTestUser(t, ctx, userRepo)
		role1 := createTestRole(t, ctx, roleRepo, shop1.ID)
		role2 := createTestRole(t, ctx, roleRepo, createdShop2.ID)

		staff1 := entities.Staff{
			UserID: user.ID,
			RoleID: role1.ID,
			ShopID: shop1.ID,
		}
		staff2 := entities.Staff{
			UserID: user.ID,
			RoleID: role2.ID,
			ShopID: createdShop2.ID,
		}

		_, err = staffRepo.Create(ctx, staff1)
		require.NoError(t, err)
		_, err = staffRepo.Create(ctx, staff2)
		require.NoError(t, err)

		staffs := staffRepo.FindByUserID(ctx, user.ID)
		assert.Len(t, staffs, 2)
	})

	t.Run("returns empty slice when no staffs for user", func(t *testing.T) {
		ctx := context.Background()
		_, _, userRepo, staffRepo := setupStaffTest(t)
		user := createTestUser(t, ctx, userRepo)

		staffs := staffRepo.FindByUserID(ctx, user.ID)
		assert.Empty(t, staffs)
	})

	t.Run("excludes soft deleted staffs", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff1 := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}
		staff2 := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created1, err := staffRepo.Create(ctx, staff1)
		require.NoError(t, err)
		_, err = staffRepo.Create(ctx, staff2)
		require.NoError(t, err)

		err = staffRepo.Delete(ctx, created1)
		require.NoError(t, err)

		staffs := staffRepo.FindByUserID(ctx, user.ID)
		assert.Len(t, staffs, 1)
	})
}

func TestStaffRepository_FindByRoleID(t *testing.T) {
	t.Run("returns staffs for a role", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user1 := createTestUser(t, ctx, userRepo)
		user2 := userentities.User{
			FullName: "User Two",
			Phone:    "0987654321",
			Email:    "user2@example.com",
			Password: "hashedpassword",
		}
		createdUser2, err := userRepo.Create(ctx, user2)
		require.NoError(t, err)

		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff1 := entities.Staff{
			UserID: user1.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}
		staff2 := entities.Staff{
			UserID: createdUser2.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		_, err = staffRepo.Create(ctx, staff1)
		require.NoError(t, err)
		_, err = staffRepo.Create(ctx, staff2)
		require.NoError(t, err)

		staffs := staffRepo.FindByRoleID(ctx, role.ID)
		assert.Len(t, staffs, 2)
	})

	t.Run("returns empty slice when no staffs for role", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, _, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staffs := staffRepo.FindByRoleID(ctx, role.ID)
		assert.Empty(t, staffs)
	})

	t.Run("excludes soft deleted staffs", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff1 := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}
		staff2 := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created1, err := staffRepo.Create(ctx, staff1)
		require.NoError(t, err)
		_, err = staffRepo.Create(ctx, staff2)
		require.NoError(t, err)

		err = staffRepo.Delete(ctx, created1)
		require.NoError(t, err)

		staffs := staffRepo.FindByRoleID(ctx, role.ID)
		assert.Len(t, staffs, 1)
	})
}

func TestStaffRepository_FindByShopIDAndUserID(t *testing.T) {
	t.Run("returns staff when found", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created, err := staffRepo.Create(ctx, staff)
		require.NoError(t, err)

		found, err := staffRepo.FindByShopIDAndUserID(ctx, shop.ID, user.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, shop.ID, found.ShopID)
		assert.Equal(t, user.ID, found.UserID)
	})

	t.Run("returns error when staff not found", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, _, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)

		_, err := staffRepo.FindByShopIDAndUserID(ctx, shop.ID, user.ID)
		assert.Error(t, err)
		assert.Equal(t, "staff not found", err.Error())
	})

	t.Run("does not find soft deleted staffs", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created, err := staffRepo.Create(ctx, staff)
		require.NoError(t, err)

		err = staffRepo.Delete(ctx, created)
		require.NoError(t, err)

		_, err = staffRepo.FindByShopIDAndUserID(ctx, shop.ID, user.ID)
		assert.Error(t, err)
		assert.Equal(t, "staff not found", err.Error())
	})
}

func TestStaffRepository_Create(t *testing.T) {
	t.Run("creates staff successfully", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created, err := staffRepo.Create(ctx, staff)
		require.NoError(t, err)
		assert.NotZero(t, created.ID)
		assert.Equal(t, user.ID, created.UserID)
		assert.Equal(t, role.ID, created.RoleID)
		assert.Equal(t, shop.ID, created.ShopID)
		assert.NotZero(t, created.CreatedAt)
		assert.NotZero(t, created.UpdatedAt)
	})
}

func TestStaffRepository_Update(t *testing.T) {
	t.Run("updates staff successfully", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role1 := createTestRole(t, ctx, roleRepo, shop.ID)

		role2 := entities.Role{
			Name:        "Updated Role",
			Description: "Updated role description",
			ShopID:      shop.ID,
		}
		createdRole2, err := roleRepo.Create(ctx, role2)
		require.NoError(t, err)

		staff := entities.Staff{
			UserID: user.ID,
			RoleID: role1.ID,
			ShopID: shop.ID,
		}

		created, err := staffRepo.Create(ctx, staff)
		require.NoError(t, err)

		originalUpdatedAt := created.UpdatedAt
		time.Sleep(10 * time.Millisecond)

		created.RoleID = createdRole2.ID

		updated, err := staffRepo.Update(ctx, created)
		require.NoError(t, err)
		assert.Equal(t, createdRole2.ID, updated.RoleID)
		assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))
	})
}

func TestStaffRepository_Delete(t *testing.T) {
	t.Run("soft deletes staff successfully", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created, err := staffRepo.Create(ctx, staff)
		require.NoError(t, err)

		err = staffRepo.Delete(ctx, created)
		require.NoError(t, err)

		_, err = staffRepo.FindByID(ctx, created.ID)
		assert.Error(t, err)
		assert.Equal(t, "staff not found", err.Error())

		staffs := staffRepo.FindByShopID(ctx, shop.ID)
		assert.Empty(t, staffs)
	})

	t.Run("can delete multiple staffs", func(t *testing.T) {
		ctx := context.Background()
		shopRepo, roleRepo, userRepo, staffRepo := setupStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		staff1 := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}
		staff2 := entities.Staff{
			UserID: user.ID,
			RoleID: role.ID,
			ShopID: shop.ID,
		}

		created1, err := staffRepo.Create(ctx, staff1)
		require.NoError(t, err)
		created2, err := staffRepo.Create(ctx, staff2)
		require.NoError(t, err)

		err = staffRepo.Delete(ctx, created1)
		require.NoError(t, err)
		err = staffRepo.Delete(ctx, created2)
		require.NoError(t, err)

		staffs := staffRepo.All(ctx)
		assert.Empty(t, staffs)
	})
}
