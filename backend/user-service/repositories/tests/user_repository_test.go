package tests

import (
	"context"
	"testing"
	"user-service/models/dto"

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

	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	req := &dto.RegisterRequest{
		FirstName:   "test",
		LastName:    "user",
		Username:    "testuser",
		Email:       "test@example.com",
		PhoneNumber: "1234567890",
		Password:    "password",
		DateOfBirth: "2006-01-02",
		RoleID:      1,
	}

	user, err := repo.Create(context.Background(), req)
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
	err = db.Create(user).Error
	require.NoError(t, err)

	result, err := repo.FindByID(context.Background(), user.ID)
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

	result, err := repo.FindByID(context.Background(), nonExistentID)
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
	err = db.Create(user).Error
	require.NoError(t, err)

	result, err := repo.FindByEmail(context.Background(), user.EmailAddress)
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

	result, err := repo.FindByEmail(context.Background(), "nonexistent@example.com")
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
	err = db.Create(user).Error
	require.NoError(t, err)

	result, err := repo.FindByUsername(context.Background(), user.Username)
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
	err = db.Create(user).Error
	require.NoError(t, err)

	// Update the user
	user.Username = "updateduser"
	err = repo.Update(context.Background(), user)
	assert.NoError(t, err)

	// Verify the update
	result, err := repo.FindByID(context.Background(), user.ID)
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
	err = db.Create(user).Error
	require.NoError(t, err)

	// Delete the user
	err = repo.Delete(context.Background(), user)
	assert.NoError(t, err)

	// Verify deletion
	result, err := repo.FindByID(context.Background(), user.ID)
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
	err = db.Create(user).Error
	require.NoError(t, err)

	// Check that email exists
	exists, err := repo.ExistsByEmail(context.Background(), user.EmailAddress)
	assert.NoError(t, err)
	assert.True(t, exists)

	// Check that non-existent email returns false
	exists, err = repo.ExistsByEmail(context.Background(), "nonexistent@example.com")
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
	err = db.Create(user).Error
	require.NoError(t, err)

	exists, err := repo.ExistsByUsername(context.Background(), user.Username)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByUsername(context.Background(), "nonexistentuser")
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
	err = db.Create(user).Error
	require.NoError(t, err)

	exists, err := repo.ExistsByPhoneNumber(context.Background(), user.PhoneNumber)
	assert.NoError(t, err)
	assert.True(t, exists)

	exists, err = repo.ExistsByPhoneNumber(context.Background(), "+9999999999")
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

	users, err := repo.FindAll(context.Background())
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
	err = db.Create(user1).Error
	require.NoError(t, err)

	user2 := CreateMockUserWithRole(2)
	user2.Username = "user2"
	err = db.Create(user2).Error
	require.NoError(t, err)

	users, err := repo.FindByRoleID(context.Background(), 1)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, uint(1), users[0].RoleID)
}
