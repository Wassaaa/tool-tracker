package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/service/mocks"
)

// Helper function to create a valid user for testing
func createTestUser(id, name, email string, role domain.UserRole) domain.User {
	now := time.Now()
	return domain.User{
		ID:        id,
		Name:      name,
		Email:     email,
		Role:      role,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UserTestMocks holds all the mock dependencies for user service testing
type UserTestMocks struct {
	Ctrl              *gomock.Controller
	MockRepo          *mocks.MockUserRepo
	MockLogger        *mocks.MockEventLogger
	Service           *UserService
	ServiceWithLogger *UserService
}

// setupUserMocks creates all necessary mocks for user service testing
func setupUserMocks(t *testing.T) *UserTestMocks {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockLogger := mocks.NewMockEventLogger(ctrl)

	return &UserTestMocks{
		Ctrl:              ctrl,
		MockRepo:          mockRepo,
		MockLogger:        mockLogger,
		Service:           NewUserService(mockRepo),
		ServiceWithLogger: NewUserService(mockRepo).WithEventLogger(mockLogger),
	}
}

// teardown cleans up the mocks
func (utm *UserTestMocks) teardown() {
	utm.Ctrl.Finish()
}

// TestUserService_CreateUser tests the user creation workflow
func TestUserService_CreateUser(t *testing.T) {
	t.Run("Successful user creation", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		userID := "123e4567-e89b-12d3-a456-426614174000"
		createdUser := createTestUser(userID, "John Doe", "john@example.com", domain.UserRoleEmployee)

		// Set expectations - check email doesn't exist, then create
		mocks.MockRepo.EXPECT().GetByEmail("john@example.com").Return(domain.User{}, assert.AnError)
		mocks.MockRepo.EXPECT().Create("John Doe", "john@example.com", domain.UserRoleEmployee).Return(createdUser, nil)
		mocks.MockLogger.EXPECT().LogUserCreated(userID, "actor-123", "User created").Return(nil)

		// Execute
		result, err := mocks.ServiceWithLogger.CreateUser("John Doe", "john@example.com", domain.UserRoleEmployee, "actor-123", "User created")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, createdUser, result)
	})

	t.Run("User creation without event logger", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		userID := "123e4567-e89b-12d3-a456-426614174000"
		createdUser := createTestUser(userID, "Jane Doe", "jane@example.com", domain.UserRoleAdmin)

		mocks.MockRepo.EXPECT().GetByEmail("jane@example.com").Return(domain.User{}, assert.AnError)
		mocks.MockRepo.EXPECT().Create("Jane Doe", "jane@example.com", domain.UserRoleAdmin).Return(createdUser, nil)

		result, err := mocks.Service.CreateUser("Jane Doe", "jane@example.com", domain.UserRoleAdmin, "actor-123", "User created")

		require.NoError(t, err)
		assert.Equal(t, createdUser, result)
	})

	t.Run("Cannot create user with existing email", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		existingUser := createTestUser("999e9999-e89b-12d3-a456-426614174000", "Existing User", "john@example.com", domain.UserRoleEmployee)

		mocks.MockRepo.EXPECT().GetByEmail("john@example.com").Return(existingUser, nil)

		_, err := mocks.Service.CreateUser("John Doe", "john@example.com", domain.UserRoleEmployee, "actor-123", "User created")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user with email 'john@example.com' already exists")
	})

	t.Run("Repository error should propagate", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		repoError := assert.AnError
		mocks.MockRepo.EXPECT().GetByEmail("john@example.com").Return(domain.User{}, assert.AnError)
		mocks.MockRepo.EXPECT().Create("John Doe", "john@example.com", domain.UserRoleEmployee).Return(domain.User{}, repoError)

		_, err := mocks.Service.CreateUser("John Doe", "john@example.com", domain.UserRoleEmployee, "actor-123", "User created")

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
	})
}

// TestUserService_UpdateUser tests the user update workflow
func TestUserService_UpdateUser(t *testing.T) {
	userID := "123e4567-e89b-12d3-a456-426614174000"
	actorID := "789e0123-e89b-12d3-a456-426614174000"

	t.Run("Successful user update", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		updatedUser := createTestUser(userID, "John Smith", "john.smith@example.com", domain.UserRoleAdmin)

		// Email is changing, so check it's not taken by another user
		mocks.MockRepo.EXPECT().GetByEmail("john.smith@example.com").Return(domain.User{}, assert.AnError)
		mocks.MockRepo.EXPECT().Update(userID, "John Smith", "john.smith@example.com", domain.UserRoleAdmin).Return(updatedUser, nil)
		mocks.MockLogger.EXPECT().LogUserUpdated(userID, actorID, "User updated").Return(nil)

		result, err := mocks.ServiceWithLogger.UpdateUser(userID, "John Smith", "john.smith@example.com", domain.UserRoleAdmin, actorID, "User updated")

		require.NoError(t, err)
		assert.Equal(t, updatedUser, result)
	})

	t.Run("Cannot update to existing email of another user", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		anotherUserID := "999e9999-e89b-12d3-a456-426614174000"
		existingUser := createTestUser(anotherUserID, "Another User", "taken@example.com", domain.UserRoleEmployee)

		mocks.MockRepo.EXPECT().GetByEmail("taken@example.com").Return(existingUser, nil)

		_, err := mocks.Service.UpdateUser(userID, "John Smith", "taken@example.com", domain.UserRoleAdmin, actorID, "User updated")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user with email 'taken@example.com' already exists")
	})

	t.Run("Invalid user ID should fail", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		_, err := mocks.Service.UpdateUser("invalid-uuid", "John Smith", "john@example.com", domain.UserRoleAdmin, actorID, "User updated")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id must be a valid UUID")
	})
}

// TestUserService_GetUser tests user retrieval
func TestUserService_GetUser(t *testing.T) {
	userID := "123e4567-e89b-12d3-a456-426614174000"

	t.Run("Successful user retrieval", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		expectedUser := createTestUser(userID, "John Doe", "john@example.com", domain.UserRoleEmployee)

		mocks.MockRepo.EXPECT().Get(userID).Return(expectedUser, nil)

		result, err := mocks.Service.GetUser(userID)

		require.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("Invalid user ID should fail", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		_, err := mocks.Service.GetUser("invalid-uuid")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id must be a valid UUID")
	})
}

// TestUserService_GetUserByEmail tests user retrieval by email
func TestUserService_GetUserByEmail(t *testing.T) {
	t.Run("Successful user retrieval by email", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		expectedUser := createTestUser("123e4567-e89b-12d3-a456-426614174000", "John Doe", "john@example.com", domain.UserRoleEmployee)

		mocks.MockRepo.EXPECT().GetByEmail("john@example.com").Return(expectedUser, nil)

		result, err := mocks.Service.GetUserByEmail("john@example.com")

		require.NoError(t, err)
		assert.Equal(t, expectedUser, result)
	})

	t.Run("Empty email should fail", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		_, err := mocks.Service.GetUserByEmail("")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user email cannot be empty")
	})
}

// TestUserService_DeleteUser tests user deletion
func TestUserService_DeleteUser(t *testing.T) {
	userID := "123e4567-e89b-12d3-a456-426614174000"
	actorID := "789e0123-e89b-12d3-a456-426614174000"

	t.Run("Successful user deletion", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		userToDelete := createTestUser(userID, "John Doe", "john@example.com", domain.UserRoleEmployee)

		mocks.MockRepo.EXPECT().Get(userID).Return(userToDelete, nil)
		mocks.MockRepo.EXPECT().Delete(userID).Return(nil)
		mocks.MockLogger.EXPECT().LogUserDeleted(userID, actorID, "User deleted").Return(nil)

		err := mocks.ServiceWithLogger.DeleteUser(userID, actorID, "User deleted")

		require.NoError(t, err)
	})

	t.Run("Invalid user ID should fail", func(t *testing.T) {
		mocks := setupUserMocks(t)
		defer mocks.teardown()

		err := mocks.Service.DeleteUser("invalid-uuid", actorID, "User deleted")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id must be a valid UUID")
	})
}
