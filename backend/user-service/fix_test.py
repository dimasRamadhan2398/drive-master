import os, re

path = 'c:/Users/Dimas/drive-master/backend/user-service/repositories/tests/user_repository_test.go'
with open(path, 'r') as f:
    content = f.read()

if '"context"' not in content:
    content = content.replace('"testing"', '"context"\n\t"testing"\n\t"user-service/models/dto"')

content = content.replace('repo.Create(user)', 'db.Create(user).Error')

test_create_old = """func TestUserRepository_Create(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	// Auto migrate only User model to avoid relationship issues with SQLite
	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	user := CreateMockUser()

	err = db.Create(user).Error
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
}"""

test_create_new = """func TestUserRepository_Create(t *testing.T) {
	db := SetupTestDB(t)
	repo := repositories.NewUserRepository(db)

	err := db.AutoMigrate(&models.User{}, &models.Role{})
	require.NoError(t, err)

	req := &dto.RegisterRequest{
		FirstName: "test",
		LastName: "user",
		Username: "testuser",
		Email: "test@example.com",
		PhoneNumber: "1234567890",
		Password: "password",
		DateOfBirth: "2006-01-02",
		RoleID: 1,
	}

	user, err := repo.Create(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, user.ID)
}"""

content = content.replace(test_create_old, test_create_new)

methods = ['FindByID', 'FindByEmail', 'FindByUsername', 'FindByPhoneNumber', 'ExistsByEmail', 'ExistsByUsername', 'ExistsByPhoneNumber', 'FindByRoleID', 'FindAll', 'Update', 'Delete', 'GetAllWithProfiles', 'FindByIDWithProfiles', 'CountByRoleID', 'FindByRoleIDWithPagination']

for method in methods:
    content = re.sub(r'repo\.' + method + r'\(\)', 'repo.' + method + '(context.Background())', content)
    content = re.sub(r'repo\.' + method + r'\(([^)])', r'repo.' + method + r'(context.Background(), \1', content)

with open(path, 'w') as f:
    f.write(content)
