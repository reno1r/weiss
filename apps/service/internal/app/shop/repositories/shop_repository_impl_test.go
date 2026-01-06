package repositories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/reno1r/weiss/apps/service/internal/app/shop/entities"
	"github.com/reno1r/weiss/apps/service/internal/testutil"
)

func TestShopRepository_All(t *testing.T) {

	t.Run("returns empty slice when no shops exist", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shops := repo.All()
		assert.Empty(t, shops)
	})

	t.Run("returns all shops", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

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

		_, err := repo.Create(shop1)
		require.NoError(t, err)
		_, err = repo.Create(shop2)
		require.NoError(t, err)

		shops := repo.All()
		assert.Len(t, shops, 2)
	})

	t.Run("excludes soft deleted shops", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop := entities.Shop{
			Name:        "Deleted Shop",
			Description: "Deleted shop description",
			Address:     "789 Elm St",
			Phone:       "1111111111",
			Email:       "deleted@example.com",
			Website:     "https://deleted.com",
			Logo:        "deleted.png",
		}

		created, err := repo.Create(shop)
		require.NoError(t, err)

		err = repo.Delete(created)
		require.NoError(t, err)

		shops := repo.All()
		assert.Empty(t, shops)
	})
}

func TestShopRepository_FindByID(t *testing.T) {
	t.Run("returns shop when found", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop := entities.Shop{
			Name:        "Test Shop",
			Description: "Test description",
			Address:     "123 Test St",
			Phone:       "1234567890",
			Email:       "test@example.com",
			Website:     "https://test.com",
			Logo:        "test.png",
		}

		created, err := repo.Create(shop)
		require.NoError(t, err)

		found, err := repo.FindByID(created.ID)
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, "Test Shop", found.Name)
		assert.Equal(t, "1234567890", found.Phone)
		assert.Equal(t, "test@example.com", found.Email)
	})

	t.Run("returns error when shop not found", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		_, err := repo.FindByID(999)
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
	})

	t.Run("does not find soft deleted shops", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop := entities.Shop{
			Name:        "Deleted Shop",
			Description: "Deleted description",
			Address:     "789 Elm St",
			Phone:       "1111111111",
			Email:       "deleted@example.com",
			Website:     "https://deleted.com",
			Logo:        "deleted.png",
		}

		created, err := repo.Create(shop)
		require.NoError(t, err)

		err = repo.Delete(created)
		require.NoError(t, err)

		_, err = repo.FindByID(created.ID)
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
	})
}

func TestShopRepository_FindByPhone(t *testing.T) {
	t.Run("returns shop when found", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop := entities.Shop{
			Name:        "Test Shop",
			Description: "Test description",
			Address:     "123 Test St",
			Phone:       "1234567890",
			Email:       "test@example.com",
			Website:     "https://test.com",
			Logo:        "test.png",
		}

		created, err := repo.Create(shop)
		require.NoError(t, err)

		found, err := repo.FindByPhone("1234567890")
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, "Test Shop", found.Name)
		assert.Equal(t, "1234567890", found.Phone)
		assert.Equal(t, "test@example.com", found.Email)
	})

	t.Run("returns error when shop not found", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		_, err := repo.FindByPhone("9999999999")
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
	})

	t.Run("does not find soft deleted shops", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop := entities.Shop{
			Name:        "Deleted Shop",
			Description: "Deleted description",
			Address:     "789 Elm St",
			Phone:       "1111111111",
			Email:       "deleted@example.com",
			Website:     "https://deleted.com",
			Logo:        "deleted.png",
		}

		created, err := repo.Create(shop)
		require.NoError(t, err)

		err = repo.Delete(created)
		require.NoError(t, err)

		_, err = repo.FindByPhone("1111111111")
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
	})
}

func TestShopRepository_FindByEmail(t *testing.T) {
	t.Run("returns shop when found", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop := entities.Shop{
			Name:        "Test Shop",
			Description: "Test description",
			Address:     "123 Test St",
			Phone:       "1234567890",
			Email:       "test@example.com",
			Website:     "https://test.com",
			Logo:        "test.png",
		}

		created, err := repo.Create(shop)
		require.NoError(t, err)

		found, err := repo.FindByEmail("test@example.com")
		require.NoError(t, err)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, "Test Shop", found.Name)
		assert.Equal(t, "test@example.com", found.Email)
	})

	t.Run("returns error when shop not found", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		_, err := repo.FindByEmail("notfound@example.com")
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
	})

	t.Run("does not find soft deleted shops", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop := entities.Shop{
			Name:        "Deleted Shop",
			Description: "Deleted description",
			Address:     "789 Elm St",
			Phone:       "1111111111",
			Email:       "deleted@example.com",
			Website:     "https://deleted.com",
			Logo:        "deleted.png",
		}

		created, err := repo.Create(shop)
		require.NoError(t, err)

		err = repo.Delete(created)
		require.NoError(t, err)

		_, err = repo.FindByEmail("deleted@example.com")
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())
	})
}

func TestShopRepository_Create(t *testing.T) {
	t.Run("creates shop successfully", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop := entities.Shop{
			Name:        "New Shop",
			Description: "New shop description",
			Address:     "123 Main St",
			Phone:       "1234567890",
			Email:       "new@example.com",
			Website:     "https://newshop.com",
			Logo:        "newlogo.png",
		}

		created, err := repo.Create(shop)
		require.NoError(t, err)
		assert.NotZero(t, created.ID)
		assert.Equal(t, "New Shop", created.Name)
		assert.Equal(t, "1234567890", created.Phone)
		assert.Equal(t, "new@example.com", created.Email)
		assert.Equal(t, "123 Main St", created.Address)
		assert.Equal(t, "https://newshop.com", created.Website)
		assert.Equal(t, "newlogo.png", created.Logo)
		assert.NotZero(t, created.CreatedAt)
		assert.NotZero(t, created.UpdatedAt)
	})

	t.Run("allows duplicate phone (no unique constraint)", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop1 := entities.Shop{
			Name:        "Shop One",
			Description: "First shop",
			Address:     "123 Main St",
			Phone:       "1234567890",
			Email:       "shop1@example.com",
			Website:     "https://shop1.com",
			Logo:        "logo1.png",
		}

		_, err := repo.Create(shop1)
		require.NoError(t, err)

		shop2 := entities.Shop{
			Name:        "Shop Two",
			Description: "Second shop",
			Address:     "456 Oak Ave",
			Phone:       "1234567890",
			Email:       "shop2@example.com",
			Website:     "https://shop2.com",
			Logo:        "logo2.png",
		}

		_, err = repo.Create(shop2)
		require.NoError(t, err) // No unique constraint on phone
	})

	t.Run("allows duplicate email (no unique constraint)", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop1 := entities.Shop{
			Name:        "Shop One",
			Description: "First shop",
			Address:     "123 Main St",
			Phone:       "1234567890",
			Email:       "shop@example.com",
			Website:     "https://shop1.com",
			Logo:        "logo1.png",
		}

		_, err := repo.Create(shop1)
		require.NoError(t, err)

		shop2 := entities.Shop{
			Name:        "Shop Two",
			Description: "Second shop",
			Address:     "456 Oak Ave",
			Phone:       "0987654321",
			Email:       "shop@example.com",
			Website:     "https://shop2.com",
			Logo:        "logo2.png",
		}

		_, err = repo.Create(shop2)
		require.NoError(t, err) // No unique constraint on email
	})
}

func TestShopRepository_Update(t *testing.T) {
	t.Run("updates shop successfully", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop := entities.Shop{
			Name:        "Original Shop",
			Description: "Original description",
			Address:     "123 Main St",
			Phone:       "1234567890",
			Email:       "original@example.com",
			Website:     "https://original.com",
			Logo:        "original.png",
		}

		created, err := repo.Create(shop)
		require.NoError(t, err)

		originalUpdatedAt := created.UpdatedAt
		time.Sleep(10 * time.Millisecond) // Ensure UpdatedAt changes

		created.Name = "Updated Shop"
		created.Email = "updated@example.com"
		created.Description = "Updated description"

		updated, err := repo.Update(created)
		require.NoError(t, err)
		assert.Equal(t, "Updated Shop", updated.Name)
		assert.Equal(t, "updated@example.com", updated.Email)
		assert.Equal(t, "Updated description", updated.Description)
		assert.True(t, updated.UpdatedAt.After(originalUpdatedAt))
	})

	t.Run("updates non-zero fields", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop := entities.Shop{
			Name:        "Test Shop",
			Description: "Test description",
			Address:     "123 Test St",
			Phone:       "1234567890",
			Email:       "test@example.com",
			Website:     "https://test.com",
			Logo:        "test.png",
		}

		created, err := repo.Create(shop)
		require.NoError(t, err)

		created.Address = "456 Updated Ave"
		created.Website = "https://updated.com"
		updated, err := repo.Update(created)
		require.NoError(t, err)
		assert.Equal(t, "456 Updated Ave", updated.Address)
		assert.Equal(t, "https://updated.com", updated.Website)
	})
}

func TestShopRepository_Delete(t *testing.T) {
	t.Run("soft deletes shop successfully", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop := entities.Shop{
			Name:        "Shop To Delete",
			Description: "Delete description",
			Address:     "123 Delete St",
			Phone:       "1234567890",
			Email:       "delete@example.com",
			Website:     "https://delete.com",
			Logo:        "delete.png",
		}

		created, err := repo.Create(shop)
		require.NoError(t, err)

		err = repo.Delete(created)
		require.NoError(t, err)

		// Verify shop is soft deleted
		_, err = repo.FindByPhone("1234567890")
		assert.Error(t, err)
		assert.Equal(t, "shop not found", err.Error())

		// Verify shop still exists in database (soft deleted)
		var deletedShop entities.Shop
		err = db.Unscoped().Where("id = ?", created.ID).First(&deletedShop).Error
		require.NoError(t, err)
		assert.NotZero(t, deletedShop.DeletedAt)
	})

	t.Run("can delete multiple shops", func(t *testing.T) {
		db := testutil.SetupTestDB(t, &entities.Shop{})
		repo := NewShopRepository(db)

		shop1 := entities.Shop{
			Name:        "Shop One",
			Description: "First shop",
			Address:     "123 Main St",
			Phone:       "1111111111",
			Email:       "shop1@example.com",
			Website:     "https://shop1.com",
			Logo:        "logo1.png",
		}
		shop2 := entities.Shop{
			Name:        "Shop Two",
			Description: "Second shop",
			Address:     "456 Oak Ave",
			Phone:       "2222222222",
			Email:       "shop2@example.com",
			Website:     "https://shop2.com",
			Logo:        "logo2.png",
		}

		created1, err := repo.Create(shop1)
		require.NoError(t, err)
		created2, err := repo.Create(shop2)
		require.NoError(t, err)

		err = repo.Delete(created1)
		require.NoError(t, err)
		err = repo.Delete(created2)
		require.NoError(t, err)

		shops := repo.All()
		assert.Empty(t, shops)
	})
}
