package tests

import (
	"testing"

	"user-service/models"
	"user-service/repositories"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserRepository_FindAll_DirectDebug(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	user := &models.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: "hash",
		EmailAddress: "test@example.com",
		RoleID:       1,
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	t.Logf("User created")

	// Now test using repository
	_ = repositories.NewUserRepository(db)

	// Let's trace what happens in FindAll
	var users []*models.User

	// This is what FindMany does internally
	sql := db.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Model(&models.User{}).Find(&users)
	})
	t.Logf("SQL for Find: %s", sql)

	// Try the actual query
	err = db.Model(&models.User{}).Find(&users).Error
	t.Logf("Query result: %v, users: %d", err, len(users))

	assert.NoError(t, err)
	assert.Len(t, users, 1)
}

func TestUserRepository_RepoFindAll(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	user := &models.User{
		ID:           uuid.New(),
		Username:     "testuser",
		PasswordHash: "hash",
		EmailAddress: "test@example.com",
		RoleID:       1,
	}
	err = db.Create(user).Error
	require.NoError(t, err)

	// Use repository
	repo := repositories.NewUserRepository(db)

	users, err := repo.FindAll()
	t.Logf("FindAll result: err=%v, users=%d", err, len(users))

	if err != nil {
		// Get more details
		t.Logf("Error details: %+v", err)
	}

	require.NoError(t, err)
	assert.Len(t, users, 1)
}