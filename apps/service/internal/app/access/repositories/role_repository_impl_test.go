package repositories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reno1r/weiss/apps/service/internal/app/access/entities"
	shopentities "github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func setupRoleTest(t *testing.T) (repositories.ShopRepository, RoleRepository) {
	db := testutil.SetupTestDB(t, &shopentities.Shop{}, &entities.Role{})
	shopRepo := repositories.NewShopRepository(db)
	roleRepo := NewRoleRepository(db)
	return shopRepo, roleRepo
}

func createTestShop(t *testing.T, shopRepo repositories.ShopRepository) shopentities.Shop {
	shop := shopentities.Shop{
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
	return created
}

func TestRoleRepository_All(t *testing.T) {
	t.Run("returns empty slice when no roles exist", func(t *testing.T) {
		_, roleRepo := setupRoleTest(t)

		roles := roleRepo.All()
		assert.Empty(t, roles)
	})

	t.Run("returns all roles", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop := createTestShop(t, shopRepo)

		role1 := entities.Role{
			Name:        "Manager",
			Description: "Shop manager role",
			ShopID:      shop.ID,
		}
		role2 := entities.Role{
			Name:        "Cashier",
			Description: "Cashier role",
			ShopID:      shop.ID,
		}

		_, err := roleRepo.Create(role1)
		require.NoError(t, err)
		_, err = roleRepo.Create(role2)
		require.NoError(t, err)

		roles := roleRepo.All()
		assert.Len(t, roles, 2)
	})

	t.Run("excludes soft deleted roles", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop := createTestShop(t, shopRepo)

		role := entities.Role{
			Name:        "Deleted Role",
			Description: "Deleted role description",
			ShopID:      shop.ID,
		}

		created, err := roleRepo.Create(role)
		require.NoError(t, err)

		err = roleRepo.Delete(created)
		require.NoError(t, err)

		roles := roleRepo.All()
		assert.Empty(t, roles)
	})
}

func TestRoleRepository_FindByID(t *testing.T) {
	t.Run("returns role when found", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop := createTestShop(t, shopRepo)

		role := entities.Role{
			Name:        "Test Role",
			Description: "Test description",
			ShopID:      shop.ID,
		}

		created, err := roleRepo.Create(role)
		require.NoError(t, err)

		found, err := roleRepo.FindByID(created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, "Test Role", found.Name)
		assert.Equal(t, shop.ID, found.ShopID)
	})

	t.Run("returns error when role not found", func(t *testing.T) {
		_, roleRepo := setupRoleTest(t)

		_, err := roleRepo.FindByID(999)
		assert.Error(t, err)
		assert.Equal(t, "role not found", err.Error())
	})

	t.Run("does not find soft deleted roles", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop := createTestShop(t, shopRepo)

		role := entities.Role{
			Name:        "Deleted Role",
			Description: "Deleted description",
			ShopID:      shop.ID,
		}

		created, err := roleRepo.Create(role)
		require.NoError(t, err)

		err = roleRepo.Delete(created)
		require.NoError(t, err)

		_, err = roleRepo.FindByID(created.ID)
		assert.Error(t, err)
		assert.Equal(t, "role not found", err.Error())
	})
}

func TestRoleRepository_FindByShopID(t *testing.T) {
	t.Run("returns roles for a shop", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop1 := createTestShop(t, shopRepo)

		shop2 := shopentities.Shop{
			Name:        "Shop Two",
			Description: "Second shop",
			Address:     "456 Oak Ave",
			Phone:       "0987654321",
			Email:       "shop2@example.com",
			Website:     "https://shop2.com",
			Logo:        "logo2.png",
		}
		createdShop2, err := shopRepo.Create(shop2)
		require.NoError(t, err)

		role1 := entities.Role{
			Name:        "Manager",
			Description: "Manager role",
			ShopID:      shop1.ID,
		}
		role2 := entities.Role{
			Name:        "Cashier",
			Description: "Cashier role",
			ShopID:      shop1.ID,
		}
		role3 := entities.Role{
			Name:        "Admin",
			Description: "Admin role for shop 2",
			ShopID:      createdShop2.ID,
		}

		_, err = roleRepo.Create(role1)
		require.NoError(t, err)
		_, err = roleRepo.Create(role2)
		require.NoError(t, err)
		_, err = roleRepo.Create(role3)
		require.NoError(t, err)

		roles := roleRepo.FindByShopID(shop1.ID)
		assert.Len(t, roles, 2)
		assert.Equal(t, "Manager", roles[0].Name)
		assert.Equal(t, "Cashier", roles[1].Name)
	})

	t.Run("returns empty slice when no roles for shop", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop := createTestShop(t, shopRepo)

		roles := roleRepo.FindByShopID(shop.ID)
		assert.Empty(t, roles)
	})

	t.Run("excludes soft deleted roles", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop := createTestShop(t, shopRepo)

		role1 := entities.Role{
			Name:        "Manager",
			Description: "Manager role",
			ShopID:      shop.ID,
		}
		role2 := entities.Role{
			Name:        "Cashier",
			Description: "Cashier role",
			ShopID:      shop.ID,
		}

		created1, err := roleRepo.Create(role1)
		require.NoError(t, err)
		_, err = roleRepo.Create(role2)
		require.NoError(t, err)

		err = roleRepo.Delete(created1)
		require.NoError(t, err)

		roles := roleRepo.FindByShopID(shop.ID)
		assert.Len(t, roles, 1)
		assert.Equal(t, "Cashier", roles[0].Name)
	})
}

func TestRoleRepository_Create(t *testing.T) {
	t.Run("creates role successfully", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop := createTestShop(t, shopRepo)

		role := entities.Role{
			Name:        "New Role",
			Description: "New role description",
			ShopID:      shop.ID,
		}

		created, err := roleRepo.Create(role)
		require.NoError(t, err)
		assert.NotZero(t, created.ID)
		assert.Equal(t, "New Role", created.Name)
		assert.Equal(t, "New role description", created.Description)
		assert.Equal(t, shop.ID, created.ShopID)
		assert.NotZero(t, created.CreatedAt)
		assert.NotZero(t, created.UpdatedAt)
	})
}

func TestRoleRepository_Update(t *testing.T) {
	t.Run("updates role successfully", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop := createTestShop(t, shopRepo)

		role := entities.Role{
			Name:        "Original Role",
			Description: "Original description",
			ShopID:      shop.ID,
		}

		created, err := roleRepo.Create(role)
		require.NoError(t, err)

		originalUpdatedAt := created.UpdatedAt
		time.Sleep(10 * time.Millisecond)

		created.Name = "Updated Role"
		created.Description = "Updated description"

		updated, err := roleRepo.Update(created)
		require.NoError(t, err)
		assert.Equal(t, "Updated Role", updated.Name)
		assert.Equal(t, "Updated description", updated.Description)
		assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("updates non-zero fields", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop := createTestShop(t, shopRepo)

		role := entities.Role{
			Name:        "Test Role",
			Description: "Test description",
			ShopID:      shop.ID,
		}

		created, err := roleRepo.Create(role)
		require.NoError(t, err)

		created.Description = "New description"
		updated, err := roleRepo.Update(created)
		require.NoError(t, err)
		assert.Equal(t, "New description", updated.Description)
	})
}

func TestRoleRepository_Delete(t *testing.T) {
	t.Run("soft deletes role successfully", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop := createTestShop(t, shopRepo)

		role := entities.Role{
			Name:        "Role To Delete",
			Description: "Delete description",
			ShopID:      shop.ID,
		}

		created, err := roleRepo.Create(role)
		require.NoError(t, err)

		err = roleRepo.Delete(created)
		require.NoError(t, err)

		_, err = roleRepo.FindByID(created.ID)
		assert.Error(t, err)
		assert.Equal(t, "role not found", err.Error())

		roles := roleRepo.FindByShopID(shop.ID)
		assert.Empty(t, roles)
	})

	t.Run("can delete multiple roles", func(t *testing.T) {
		shopRepo, roleRepo := setupRoleTest(t)
		shop := createTestShop(t, shopRepo)

		role1 := entities.Role{
			Name:        "Role One",
			Description: "First role",
			ShopID:      shop.ID,
		}
		role2 := entities.Role{
			Name:        "Role Two",
			Description: "Second role",
			ShopID:      shop.ID,
		}

		created1, err := roleRepo.Create(role1)
		require.NoError(t, err)
		created2, err := roleRepo.Create(role2)
		require.NoError(t, err)

		err = roleRepo.Delete(created1)
		require.NoError(t, err)
		err = roleRepo.Delete(created2)
		require.NoError(t, err)

		roles := roleRepo.All()
		assert.Empty(t, roles)
	})
}
