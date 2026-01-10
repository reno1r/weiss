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

func setupAssignStaffTest(t *testing.T) (*AssignStaffUsecase, accessrepositories.StaffRepository, shoprepositories.ShopRepository, accessrepositories.RoleRepository, userrepositories.UserRepository) {
	db := testutil.SetupTestDB(t, &shopentities.Shop{}, &accessentities.Role{}, &userentities.User{}, &accessentities.Staff{})
	shopRepo := shoprepositories.NewShopRepository(db)
	roleRepo := accessrepositories.NewRoleRepository(db)
	userRepo := userrepositories.NewUserRepository(db)
	staffRepo := accessrepositories.NewStaffRepository(db)
	usecase := NewAssignStaffUsecase(staffRepo)
	return usecase, staffRepo, shopRepo, roleRepo, userRepo
}

func TestAssignStaffUsecase_Execute(t *testing.T) {
	t.Run("assigns staff successfully", func(t *testing.T) {
		ctx := context.Background()
		usecase, staffRepo, shopRepo, roleRepo, userRepo := setupAssignStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		result, err := usecase.Execute(ctx, AssignStaffParam{
			User: &user,
			Shop: &shop,
			Role: &role,
		})
		require.NoError(t, err)
		assert.NotNil(t, result.Staff)
		assert.Equal(t, user.ID, result.Staff.UserID)
		assert.Equal(t, shop.ID, result.Staff.ShopID)
		assert.Equal(t, role.ID, result.Staff.RoleID)
		assert.NotNil(t, result.Staff.User)
		assert.NotNil(t, result.Staff.Role)
		assert.NotNil(t, result.Staff.Shop)

		// Verify staff was created in database
		found, err := staffRepo.FindByID(ctx, result.Staff.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, found.UserID)
		assert.Equal(t, shop.ID, found.ShopID)
		assert.Equal(t, role.ID, found.RoleID)
	})

	t.Run("returns error when staff already assigned", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, shopRepo, roleRepo, userRepo := setupAssignStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		// Assign staff first time
		_, err := usecase.Execute(ctx, AssignStaffParam{
			User: &user,
			Shop: &shop,
			Role: &role,
		})
		require.NoError(t, err)

		// Try to assign same user to same shop again
		result, err := usecase.Execute(ctx, AssignStaffParam{
			User: &user,
			Shop: &shop,
			Role: &role,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "staff already assigned")
		assert.Nil(t, result)
	})

	t.Run("allows same user to be assigned to different shops", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, shopRepo, roleRepo, userRepo := setupAssignStaffTest(t)
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

		// Assign user to shop1
		result1, err := usecase.Execute(ctx, AssignStaffParam{
			User: &user,
			Shop: &shop1,
			Role: &role1,
		})
		require.NoError(t, err)
		assert.NotNil(t, result1.Staff)

		// Assign same user to shop2
		result2, err := usecase.Execute(ctx, AssignStaffParam{
			User: &user,
			Shop: &createdShop2,
			Role: &role2,
		})
		require.NoError(t, err)
		assert.NotNil(t, result2.Staff)
		assert.NotEqual(t, result1.Staff.ID, result2.Staff.ID)
	})

	t.Run("allows different users to be assigned to same shop", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, shopRepo, roleRepo, userRepo := setupAssignStaffTest(t)
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

		// Assign user1 to shop
		result1, err := usecase.Execute(ctx, AssignStaffParam{
			User: &user1,
			Shop: &shop,
			Role: &role,
		})
		require.NoError(t, err)
		assert.NotNil(t, result1.Staff)

		// Assign user2 to same shop
		result2, err := usecase.Execute(ctx, AssignStaffParam{
			User: &createdUser2,
			Shop: &shop,
			Role: &role,
		})
		require.NoError(t, err)
		assert.NotNil(t, result2.Staff)
		assert.NotEqual(t, result1.Staff.ID, result2.Staff.ID)
	})

	t.Run("returns validation error when user is nil", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, shopRepo, roleRepo, _ := setupAssignStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		role := createTestRole(t, ctx, roleRepo, shop.ID)

		result, err := usecase.Execute(ctx, AssignStaffParam{
			User: nil,
			Shop: &shop,
			Role: &role,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("returns validation error when shop is nil", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, _, _, userRepo := setupAssignStaffTest(t)
		user := createTestUser(t, ctx, userRepo)
		role := accessentities.Role{
			Name:        "Test Role",
			Description: "Test role description",
			ShopID:      1,
		}

		result, err := usecase.Execute(ctx, AssignStaffParam{
			User: &user,
			Shop: nil,
			Role: &role,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})

	t.Run("returns validation error when role is nil", func(t *testing.T) {
		ctx := context.Background()
		usecase, _, shopRepo, _, userRepo := setupAssignStaffTest(t)
		shop := createTestShop(t, ctx, shopRepo)
		user := createTestUser(t, ctx, userRepo)

		result, err := usecase.Execute(ctx, AssignStaffParam{
			User: &user,
			Shop: &shop,
			Role: nil,
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "validation failed")
		assert.Nil(t, result)
	})
}
