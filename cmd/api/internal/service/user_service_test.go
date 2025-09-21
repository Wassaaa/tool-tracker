package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

// TestUserService_CreateUser tests the user creation workflow
func TestUserService_CreateUser(t *testing.T) {
	t.Run("Successful user creation", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		createdUser := CreateTestUser(TestUserID, "John Doe", "john@example.com", domain.UserRoleEmployee)

		// Set expectations - check email doesn't exist, then create
		mocks.MockRepo.EXPECT().GetByEmail("john@example.com").Return(domain.User{}, assert.AnError)
		mocks.MockRepo.EXPECT().Create("John Doe", "john@example.com", domain.UserRoleEmployee).Return(createdUser, nil)
		mocks.MockLogger.EXPECT().LogUserCreated(TestUserID, TestActorID, "User created").Return(nil)

		// Execute
		result, err := mocks.ServiceWithLogger.CreateUser("John Doe", "john@example.com", domain.UserRoleEmployee, TestActorID, "User created")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, createdUser, result)
	})

	t.Run("User creation without event logger", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		createdUser := CreateTestUser(TestUserID, "Jane Doe", "jane@example.com", domain.UserRoleAdmin)

		mocks.MockRepo.EXPECT().GetByEmail("jane@example.com").Return(domain.User{}, assert.AnError)
		mocks.MockRepo.EXPECT().Create("Jane Doe", "jane@example.com", domain.UserRoleAdmin).Return(createdUser, nil)

		result, err := mocks.Service.CreateUser("Jane Doe", "jane@example.com", domain.UserRoleAdmin, TestActorID, "User created")

		require.NoError(t, err)
		assert.Equal(t, createdUser, result)
	})

	t.Run("Cannot create user with existing email", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		existingUser := CreateTestUser(TestUserID2, "Existing User", "john@example.com", domain.UserRoleEmployee)

		mocks.MockRepo.EXPECT().GetByEmail("john@example.com").Return(existingUser, nil)

		_, err := mocks.Service.CreateUser("John Doe", "john@example.com", domain.UserRoleEmployee, TestActorID, "User created")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user with email 'john@example.com' already exists")
	})

	t.Run("Repository error should propagate", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		repoError := assert.AnError
		mocks.MockRepo.EXPECT().GetByEmail("john@example.com").Return(domain.User{}, assert.AnError)
		mocks.MockRepo.EXPECT().Create("John Doe", "john@example.com", domain.UserRoleEmployee).Return(domain.User{}, repoError)

		_, err := mocks.Service.CreateUser("John Doe", "john@example.com", domain.UserRoleEmployee, TestActorID, "User created")

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
	})
}

// TestUserService_UpdateUser tests the user update workflow
func TestUserService_UpdateUser(t *testing.T) {
	t.Run("Successful user update", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		updatedUser := CreateTestUser(TestUserID, "John Smith", "john.smith@example.com", domain.UserRoleAdmin)

		// Email is changing, so check it's not taken by another user
		mocks.MockRepo.EXPECT().GetByEmail("john.smith@example.com").Return(domain.User{}, assert.AnError)
		mocks.MockRepo.EXPECT().Update(TestUserID, "John Smith", "john.smith@example.com", domain.UserRoleAdmin).Return(updatedUser, nil)
		mocks.MockLogger.EXPECT().LogUserUpdated(TestUserID, TestActorID, "User updated").Return(nil)

		result, err := mocks.ServiceWithLogger.UpdateUser(TestUserID, "John Smith", "john.smith@example.com", domain.UserRoleAdmin, TestActorID, "User updated")

		require.NoError(t, err)
		assert.Equal(t, updatedUser, result)
	})

	t.Run("Cannot update to existing email of another user", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		existingUser := CreateTestUser(TestUserID2, "Another User", "taken@example.com", domain.UserRoleEmployee)

		mocks.MockRepo.EXPECT().GetByEmail("taken@example.com").Return(existingUser, nil)

		_, err := mocks.Service.UpdateUser(TestUserID, "John Smith", "taken@example.com", domain.UserRoleAdmin, TestActorID, "User updated")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user with email 'taken@example.com' already exists")
	})

	t.Run("Invalid user ID should fail", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		_, err := mocks.Service.UpdateUser(InvalidUUID, "John Smith", "john@example.com", domain.UserRoleAdmin, TestActorID, "User updated")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id must be a valid UUID")
	})
}

// TestUserService_GetUser tests user retrieval
func TestUserService_GetUser(t *testing.T) {
	t.Run("Successful user retrieval", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		expectedUser := CreateTestUser(TestUserID, "John Doe", "john@example.com", domain.UserRoleEmployee)

		mocks.MockRepo.EXPECT().Get(TestUserID).Return(expectedUser, nil)

		result, err := mocks.Service.GetUser(TestUserID)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("Invalid user ID should fail", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		_, err := mocks.Service.GetUser(InvalidUUID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id must be a valid UUID")
	})
}

// TestUserService_GetUserByEmail tests user retrieval by email
func TestUserService_GetUserByEmail(t *testing.T) {
	t.Run("Successful user retrieval by email", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		expectedUser := CreateTestUser(TestUserID, "John Doe", "john@example.com", domain.UserRoleEmployee)

		mocks.MockRepo.EXPECT().GetByEmail("john@example.com").Return(expectedUser, nil)

		result, err := mocks.Service.GetUserByEmail("john@example.com")

		require.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("Empty email should fail", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		_, err := mocks.Service.GetUserByEmail("")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user email cannot be empty")
	})
}

// TestUserService_DeleteUser tests user deletion
func TestUserService_DeleteUser(t *testing.T) {
	t.Run("Successful user deletion", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		userToDelete := CreateTestUser(TestUserID, "John Doe", "john@example.com", domain.UserRoleEmployee)

		mocks.MockRepo.EXPECT().Get(TestUserID).Return(userToDelete, nil)
		mocks.MockRepo.EXPECT().Delete(TestUserID).Return(nil)
		mocks.MockLogger.EXPECT().LogUserDeleted(TestUserID, TestActorID, "User deleted").Return(nil)

		err := mocks.ServiceWithLogger.DeleteUser(TestUserID, TestActorID, "User deleted")

		require.NoError(t, err)
	})

	t.Run("Invalid user ID should fail", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		err := mocks.Service.DeleteUser(InvalidUUID, TestActorID, "User deleted")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id must be a valid UUID")
	})
}

// TestUserService_ListUsers tests user listing
func TestUserService_ListUsers(t *testing.T) {
	t.Run("Successful users listing", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		expectedUsers := []domain.User{
			CreateTestUser("user1", "John Doe", "john@example.com", domain.UserRoleEmployee),
			CreateTestUser("user2", "Jane Smith", "jane@example.com", domain.UserRoleAdmin),
		}

		mocks.MockRepo.EXPECT().List(10, 0).Return(expectedUsers, nil)

		result, err := mocks.Service.ListUsers(10, 0)

		require.NoError(t, err)
		assert.Equal(t, expectedUsers, result)
	})

	t.Run("Default limit applied when zero", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().List(10, 0).Return([]domain.User{}, nil)

		_, err := mocks.Service.ListUsers(0, 0)

		require.NoError(t, err)
	})

	t.Run("Maximum limit enforced", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().List(100, 0).Return([]domain.User{}, nil)

		_, err := mocks.Service.ListUsers(150, 0)

		require.NoError(t, err)
	})

	t.Run("Negative offset corrected", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().List(10, 0).Return([]domain.User{}, nil)

		_, err := mocks.Service.ListUsers(10, -5)

		require.NoError(t, err)
	})
}

// TestUserService_ListUsersByRole tests listing users by role
func TestUserService_ListUsersByRole(t *testing.T) {
	t.Run("Successful list users by role", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		expectedUsers := []domain.User{
			CreateTestUser("user1", "John Doe", "john@example.com", domain.UserRoleAdmin),
			CreateTestUser("user2", "Jane Smith", "jane@example.com", domain.UserRoleAdmin),
		}

		mocks.MockRepo.EXPECT().ListByRole(domain.UserRoleAdmin, 10, 0).Return(expectedUsers, nil)

		result, err := mocks.Service.ListUsersByRole(domain.UserRoleAdmin, 10, 0)

		require.NoError(t, err)
		assert.Equal(t, expectedUsers, result)
	})

	t.Run("Default limit applied when zero", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().ListByRole(domain.UserRoleEmployee, 10, 0).Return([]domain.User{}, nil)

		_, err := mocks.Service.ListUsersByRole(domain.UserRoleEmployee, 0, 0)

		require.NoError(t, err)
	})

	t.Run("Maximum limit enforced", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().ListByRole(domain.UserRoleEmployee, 100, 0).Return([]domain.User{}, nil)

		_, err := mocks.Service.ListUsersByRole(domain.UserRoleEmployee, 150, 0)

		require.NoError(t, err)
	})

	t.Run("Negative offset corrected", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().ListByRole(domain.UserRoleEmployee, 10, 0).Return([]domain.User{}, nil)

		_, err := mocks.Service.ListUsersByRole(domain.UserRoleEmployee, 10, -5)

		require.NoError(t, err)
	})
}

// TestUserService_GetUserCount tests getting user count
func TestUserService_GetUserCount(t *testing.T) {
	t.Run("Successful get user count", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		expectedCount := 25

		mocks.MockRepo.EXPECT().Count().Return(expectedCount, nil)

		result, err := mocks.Service.GetUserCount()

		require.NoError(t, err)
		assert.Equal(t, expectedCount, result)
	})

	t.Run("Repository error should propagate", func(t *testing.T) {
		mocks := SetupUserServiceMocks(t)
		defer mocks.Teardown()

		repoError := assert.AnError
		mocks.MockRepo.EXPECT().Count().Return(0, repoError)

		_, err := mocks.Service.GetUserCount()

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
	})
}
