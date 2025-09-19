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

// Helper function to create a valid tool for testing
func createTestTool(id, name string, status domain.ToolStatus) domain.Tool {
	now := time.Now()
	return domain.Tool{
		ID:        &id,
		Name:      name,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// TestMocks holds all the mock dependencies for testing
type TestMocks struct {
	Ctrl              *gomock.Controller
	MockRepo          *mocks.MockToolRepo
	MockLogger        *mocks.MockEventLogger
	Service           *ToolService
	ServiceWithLogger *ToolService
}

// setupMocks creates all necessary mocks for testing
func setupMocks(t *testing.T) *TestMocks {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockToolRepo(ctrl)
	mockLogger := mocks.NewMockEventLogger(ctrl)

	return &TestMocks{
		Ctrl:              ctrl,
		MockRepo:          mockRepo,
		MockLogger:        mockLogger,
		Service:           NewToolService(mockRepo),
		ServiceWithLogger: NewToolService(mockRepo).WithEventLogger(mockLogger),
	}
}

// teardownMocks cleans up the mocks
func (tm *TestMocks) teardown() {
	tm.Ctrl.Finish()
}

// TestToolService_CreateTool tests the tool creation workflow
func TestToolService_CreateTool(t *testing.T) {
	t.Run("Successful tool creation", func(t *testing.T) {
		mocks := setupMocks(t)
		defer mocks.teardown()

		toolID := "123e4567-e89b-12d3-a456-426614174000"
		createdTool := createTestTool(toolID, "Hammer", domain.ToolStatusInOffice)

		// Set expectations
		mocks.MockRepo.EXPECT().Create("Hammer", domain.ToolStatusInOffice).Return(createdTool, nil)
		mocks.MockLogger.EXPECT().LogToolCreated(toolID, "actor-123", "Tool created").Return(nil)

		// Execute
		result, err := mocks.ServiceWithLogger.CreateTool("Hammer", domain.ToolStatusInOffice, "actor-123", "Tool created")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, createdTool, result)
	})

	t.Run("Tool creation without event logger", func(t *testing.T) {
		mocks := setupMocks(t)
		defer mocks.teardown()

		toolID := "123e4567-e89b-12d3-a456-426614174000"
		createdTool := createTestTool(toolID, "Screwdriver", domain.ToolStatusInOffice)

		mocks.MockRepo.EXPECT().Create("Screwdriver", domain.ToolStatusInOffice).Return(createdTool, nil)

		result, err := mocks.Service.CreateTool("Screwdriver", domain.ToolStatusInOffice, "actor-123", "Tool created")

		require.NoError(t, err)
		assert.Equal(t, createdTool, result)
	})

	t.Run("Repository error should propagate", func(t *testing.T) {
		mocks := setupMocks(t)
		defer mocks.teardown()

		repoError := assert.AnError
		mocks.MockRepo.EXPECT().Create("Hammer", domain.ToolStatusInOffice).Return(domain.Tool{}, repoError)

		_, err := mocks.Service.CreateTool("Hammer", domain.ToolStatusInOffice, "actor-123", "Tool created")

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
	})
}

// TestToolService_CheckOutTool tests the checkout workflow
func TestToolService_CheckOutTool(t *testing.T) {
	toolID := "123e4567-e89b-12d3-a456-426614174000"
	userID := "456e7890-e89b-12d3-a456-426614174000"
	actorID := "789e0123-e89b-12d3-a456-426614174000"

	t.Run("Successful checkout", func(t *testing.T) {
		mocks := setupMocks(t)
		defer mocks.teardown()

		// Available tool
		availableTool := createTestTool(toolID, "Hammer", domain.ToolStatusInOffice)
		// Checked out tool (what we expect after checkout)
		checkedOutTool := createTestTool(toolID, "Hammer", domain.ToolStatusCheckedOut)
		checkedOutTool.CurrentUserId = &userID

		// Set expectations
		mocks.MockRepo.EXPECT().Get(toolID).Return(availableTool, nil)
		mocks.MockRepo.EXPECT().Update(gomock.Any()).Return(checkedOutTool, nil)
		mocks.MockLogger.EXPECT().LogToolCheckedOut(toolID, userID, actorID, "Checking out for project").Return(nil)

		// Execute
		result, err := mocks.ServiceWithLogger.CheckOutTool(toolID, userID, actorID, "Checking out for project")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, checkedOutTool, result)
	})

	t.Run("Cannot checkout already checked out tool", func(t *testing.T) {
		mocks := setupMocks(t)
		defer mocks.teardown()

		// Tool already checked out
		someOtherUser := "999e9999-e89b-12d3-a456-426614174000"
		checkedOutTool := createTestTool(toolID, "Hammer", domain.ToolStatusCheckedOut)
		checkedOutTool.CurrentUserId = &someOtherUser

		mocks.MockRepo.EXPECT().Get(toolID).Return(checkedOutTool, nil)

		// Execute
		_, err := mocks.Service.CheckOutTool(toolID, userID, actorID, "Checking out for project")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool is already checked out")
	})

	t.Run("Invalid user ID should fail", func(t *testing.T) {
		mocks := setupMocks(t)
		defer mocks.teardown()

		// Execute with invalid user ID - no mock expectations needed since validation happens first
		_, err := mocks.Service.CheckOutTool(toolID, "invalid-uuid", actorID, "Checking out for project")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id must be a valid UUID")
	})
}
