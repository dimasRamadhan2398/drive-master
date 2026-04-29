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

func TestUserRepository_Create(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate only User model to avoid relationship issues with SQLite
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	user := CreateMockUser()

	err = repo.Create(user)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
}

func TestUserRepository_FindByID(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate to create tables
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	user := CreateMockUser()
	err = repo.Create(user)
	require.NoError(t, err)

	result, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	assert.Equal(t, user.Username, result.Username)
}

func TestUserRepository_FindByID_NotFound(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate to create tables
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	nonExistentID := uuid.New()

	result, err := repo.FindByID(nonExistentID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate to create tables
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	user := CreateMockUser()
	err = repo.Create(user)
	require.NoError(t, err)

	result, err := repo.FindByEmail(user.EmailAddress)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.EmailAddress, result.EmailAddress)
}

func TestUserRepository_FindByEmail_NotFound(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate to create tables
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	result, err := repo.FindByEmail("nonexistent@example.com")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserRepository_FindByUsername(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate to create tables
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	user := CreateMockUser()
	err = repo.Create(user)
	require.NoError(t, err)

	result, err := repo.FindByUsername(user.Username)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.Username, result.Username)
}

func TestUserRepository_Update(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate to create tables
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	user := CreateMockUser()
	err = repo.Create(user)
	require.NoError(t, err)

	// Update the user
	user.Username = "updateduser"
	err = repo.Update(user)
	assert.NoError(t, err)

	// Verify the update
	result, err := repo.FindByID(user.ID)
	assert.NoError(t, err)
	assert.Equal(t, "updateduser", result.Username)
}

func TestUserRepository_Delete(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate to create tables
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	user := CreateMockUser()
	err = repo.Create(user)
	require.NoError(t, err)

	// Delete the user
	err = repo.Delete(user)
	assert.NoError(t, err)

	// Verify deletion
	result, err := repo.FindByID(user.ID)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUserRepository_ExistsByEmail(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate to create tables
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	user := CreateMockUser()
	err = repo.Create(user)
	require.NoError(t, err)

	// Check that email exists
	exists, err := repo.ExistsByEmail(user.EmailAddress)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Check that non-existent email returns false
	exists, err = repo.ExistsByEmail("nonexistent@example.com")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestUserRepository_ExistsByUsername(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate to create tables
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	user := CreateMockUser()
	err = repo.Create(user)
	require.NoError(t, err)

	exists, err := repo.ExistsByUsername(user.Username)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByUsername("nonexistentuser")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestUserRepository_ExistsByPhoneNumber(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate to create tables
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	user := CreateMockUser()
	err = repo.Create(user)
	require.NoError(t, err)

	exists, err := repo.ExistsByPhoneNumber(user.PhoneNumber)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByPhoneNumber("+9999999999")
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestUserRepository_FindAll(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	err = db.AutoMigrate(&models.User{})
	require.NoError(t, err)

	repo := repositories.NewUserRepository(db)

	// Create users directly using GORM to avoid any repository Create issues
	user1 := &models.User{
		ID:           uuid.New(),
		Username:     "user1",
		PasswordHash: "hash",
		EmailAddress: "user1@example.com",
		RoleID:       1,
	}
	err = db.Create(user1).Error
	require.NoError(t, err)

	user2 := &models.User{
		ID:           uuid.New(),
		Username:     "user2",
		PasswordHash: "hash",
		EmailAddress: "user2@example.com",
		RoleID:       1,
	}
	err = db.Create(user2).Error
	require.NoError(t, err)

	users, err := repo.FindAll()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
}

func TestUserRepository_FindByRoleID(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate to create tables
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	// Create users with different roles
	user1 := CreateMockUserWithRole(1)
	user1.Username = "user1"
	err = repo.Create(user1)
	require.NoError(t, err)

	user2 := CreateMockUserWithRole(2)
	user2.Username = "user2"
	err = repo.Create(user2)
	require.NoError(t, err)

	users, err := repo.FindByRoleID(1)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, uint(1), users[0].RoleID)
}