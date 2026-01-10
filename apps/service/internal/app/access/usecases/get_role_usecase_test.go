package usecases

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reno1r/weiss/apps/service/internal/app/access/entities"
	accessrepositories "github.com/reno1r/weiss/apps/service/internal/app/access/repositories"
	shopentities "github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	shoprepositories "github.com/reno1r/weiss/apps/service/internal/app/shop/repositories"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func setupGetRoleTest(t *testing.T) (*GetRoleUsecase, accessrepositories.RoleRepository, shoprepositories.ShopRepository) {
	db := testutil.SetupTestDB(t, &shopentities.Shop{}, &entities.Role{})
	shopRepo := shoprepositories.NewShopRepository(db)
	roleRepo := accessrepositories.NewRoleRepository(db)
	usecase := NewGetRoleUsecase(roleRepo)
	return usecase, roleRepo, shopRepo
}

func TestGetRoleUsecase_Execute(t *testing.T) {
	t.Run("returns role when found", func(t *testing.T) {
		ctx := context.Background()
		usecase, roleRepo, shopRepo := setupGetRoleTest(t)
		shop := createTestShop(t, ctx, shopRepo)

		role := entities.Role{
			Name:        "Manager",
			Description: "Manager role",
			ShopID:      shop.ID,
		}

		created, err := roleRepo.Create(ctx, role)
		require.NoError(t, err)

		result, err := usecase.Execute(ctx, created.ID)
		require.NoError(t, err)
		assert.NotNil(t, result.Role)
		assert.Equal(t, created.ID, result.Role.ID)
		assert.Equal(t, "Manager", result.Role.Name)
		assert.Equal(t, "Manager role", result.Role.Description)
		assert.Equal(t, shop.ID, result.Role.ShopID)
	})

	t.Run("returns error when role not found", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _ := setupGetRoleTest(t)

		result, err := usecase.Execute(ctx, 999)
		assert.Error(t, err)
		assert.Equal(t, "role not found", err.Error())
		assert.Nil(t, result)
	})

	t.Run("does not find soft deleted roles", func(t *testing.T) {
		ctx := context.Background()
		usecase, roleRepo, shopRepo := setupGetRoleTest(t)
		shop := createTestShop(t, ctx, shopRepo)

		role := entities.Role{
			Name:        "Deleted Role",
			Description: "Deleted role description",
			ShopID:      shop.ID,
		}

		created, err := roleRepo.Create(ctx, role)
		require.NoError(t, err)

		err = roleRepo.Delete(ctx, created)
		require.NoError(t, err)

		result, err := usecase.Execute(ctx, created.ID)
		assert.Error(t, err)
		assert.Equal(t, "role not found", err.Error())
		assert.Nil(t, result)
	})
}
