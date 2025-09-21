package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
)

// TestToolService_CreateTool tests the tool creation workflow
func TestToolService_CreateTool(t *testing.T) {
	t.Run("Successful tool creation", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		createdTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusInOffice)

		// Set expectations
		mocks.MockRepo.EXPECT().Create("Hammer", domain.ToolStatusInOffice).Return(createdTool, nil)
		mocks.MockLogger.EXPECT().LogToolCreated(TestToolID, TestActorID, "Tool created").Return(nil)

		// Execute
		result, err := mocks.ServiceWithLogger.CreateTool("Hammer", domain.ToolStatusInOffice, TestActorID, "Tool created")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, createdTool, result)
	})

	t.Run("Tool creation without event logger", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		createdTool := CreateTestTool(TestToolID, "Screwdriver", domain.ToolStatusInOffice)

		mocks.MockRepo.EXPECT().Create("Screwdriver", domain.ToolStatusInOffice).Return(createdTool, nil)

		result, err := mocks.Service.CreateTool("Screwdriver", domain.ToolStatusInOffice, TestActorID, "Tool created")

		require.NoError(t, err)
		assert.Equal(t, createdTool, result)
	})

	t.Run("Repository error should propagate", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		repoError := assert.AnError
		mocks.MockRepo.EXPECT().Create("Hammer", domain.ToolStatusInOffice).Return(domain.Tool{}, repoError)

		_, err := mocks.Service.CreateTool("Hammer", domain.ToolStatusInOffice, TestActorID, "Tool created")

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
	})
}

// TestToolService_CheckOutTool tests the checkout workflow
func TestToolService_CheckOutTool(t *testing.T) {
	t.Run("Successful checkout", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		// Available tool
		availableTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusInOffice)
		// Checked out tool (what we expect after checkout)
		checkedOutTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusCheckedOut)
		userID := TestUserID
		checkedOutTool.CurrentUserId = &userID

		// Set expectations
		mocks.MockRepo.EXPECT().Get(TestToolID).Return(availableTool, nil)
		mocks.MockRepo.EXPECT().Update(gomock.Any()).Return(checkedOutTool, nil)
		mocks.MockLogger.EXPECT().LogToolCheckedOut(TestToolID, TestUserID, TestActorID, "Checking out for project").Return(nil)

		// Execute
		result, err := mocks.ServiceWithLogger.CheckOutTool(TestToolID, TestUserID, TestActorID, "Checking out for project")

		// Assert
		require.NoError(t, err)
		assert.Equal(t, checkedOutTool, result)
	})

	t.Run("Cannot checkout already checked out tool", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		// Tool already checked out
		checkedOutTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusCheckedOut)
		userID2 := TestUserID2
		checkedOutTool.CurrentUserId = &userID2

		mocks.MockRepo.EXPECT().Get(TestToolID).Return(checkedOutTool, nil)

		// Execute
		_, err := mocks.Service.CheckOutTool(TestToolID, TestUserID, TestActorID, "Checking out for project")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool is already checked out")
	})

	t.Run("Invalid user ID should fail", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		// Execute with invalid user ID - no mock expectations needed since validation happens first
		_, err := mocks.Service.CheckOutTool(TestToolID, InvalidUUID, TestActorID, "Checking out for project")

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id must be a valid UUID")
	})

	t.Run("Invalid tool ID should fail", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		_, err := mocks.Service.CheckOutTool(InvalidUUID, TestUserID, TestActorID, "Checking out for project")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool_id must be a valid UUID")
	})
}

// TestToolService_ReturnTool tests the return workflow
func TestToolService_ReturnTool(t *testing.T) {
	t.Run("Successful return", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		// Checked out tool
		checkedOutTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusCheckedOut)
		userID := TestUserID
		checkedOutTool.CurrentUserId = &userID
		// Returned tool
		returnedTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusInOffice)

		mocks.MockRepo.EXPECT().Get(TestToolID).Return(checkedOutTool, nil)
		mocks.MockRepo.EXPECT().Update(gomock.Any()).Return(returnedTool, nil)
		mocks.MockLogger.EXPECT().LogToolCheckedIn(TestToolID, TestUserID, TestActorID, "Returning tool").Return(nil)

		result, err := mocks.ServiceWithLogger.ReturnTool(TestToolID, TestActorID, "Returning tool")

		require.NoError(t, err)
		assert.Equal(t, returnedTool, result)
	})

	t.Run("Cannot return tool that is already checked in", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		availableTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusInOffice)

		mocks.MockRepo.EXPECT().Get(TestToolID).Return(availableTool, nil)

		_, err := mocks.Service.ReturnTool(TestToolID, TestActorID, "Returning tool")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool is already checked in")
	})

	t.Run("Invalid tool ID should fail", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		_, err := mocks.Service.ReturnTool(InvalidUUID, TestActorID, "Returning tool")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool_id must be a valid UUID")
	})
}

// TestToolService_ListTools tests the list tools functionality
func TestToolService_ListTools(t *testing.T) {
	t.Run("Successful tools listing", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		expectedTools := []domain.Tool{
			CreateTestTool("tool1", "Hammer", domain.ToolStatusInOffice),
			CreateTestTool("tool2", "Screwdriver", domain.ToolStatusCheckedOut),
		}

		mocks.MockRepo.EXPECT().List(10, 0).Return(expectedTools, nil)

		result, err := mocks.Service.ListTools(10, 0)

		require.NoError(t, err)
		assert.Equal(t, expectedTools, result)
	})

	t.Run("Default limit applied when zero", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().List(10, 0).Return([]domain.Tool{}, nil)

		_, err := mocks.Service.ListTools(0, 0)

		require.NoError(t, err)
	})

	t.Run("Maximum limit enforced", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().List(100, 0).Return([]domain.Tool{}, nil)

		_, err := mocks.Service.ListTools(150, 0)

		require.NoError(t, err)
	})

	t.Run("Negative offset corrected", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		mocks.MockRepo.EXPECT().List(10, 0).Return([]domain.Tool{}, nil)

		_, err := mocks.Service.ListTools(10, -5)

		require.NoError(t, err)
	})
}

// TestToolService_GetTool tests tool retrieval by ID
func TestToolService_GetTool(t *testing.T) {
	t.Run("Successful tool retrieval", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		expectedTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusInOffice)

		mocks.MockRepo.EXPECT().Get(TestToolID).Return(expectedTool, nil)

		result, err := mocks.Service.GetTool(TestToolID)

		require.NoError(t, err)
		assert.Equal(t, expectedTool, result)
	})

	t.Run("Invalid tool ID should fail", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		_, err := mocks.Service.GetTool(InvalidUUID)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool_id must be a valid UUID")
	})
}

// TestToolService_UpdateTool tests tool updates
func TestToolService_UpdateTool(t *testing.T) {
	t.Run("Successful tool update", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		existingTool := CreateTestTool(TestToolID, "Old Hammer", domain.ToolStatusInOffice)
		updatedTool := CreateTestTool(TestToolID, "New Hammer", domain.ToolStatusMaintenance)

		mocks.MockRepo.EXPECT().Get(TestToolID).Return(existingTool, nil)
		mocks.MockRepo.EXPECT().Update(gomock.Any()).Return(updatedTool, nil)
		mocks.MockLogger.EXPECT().LogToolUpdated(TestToolID, TestActorID, "Tool updated").Return(nil)

		result, err := mocks.ServiceWithLogger.UpdateTool(TestToolID, "New Hammer", domain.ToolStatusMaintenance, TestActorID, "Tool updated")

		require.NoError(t, err)
		assert.Equal(t, updatedTool, result)
	})

	t.Run("Invalid tool ID should fail", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		_, err := mocks.Service.UpdateTool(InvalidUUID, "Hammer", domain.ToolStatusInOffice, TestActorID, "Tool updated")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool_id must be a valid UUID")
	})
}

// TestToolService_SendToMaintenance tests sending tools to maintenance
func TestToolService_SendToMaintenance(t *testing.T) {
	t.Run("Successful send to maintenance", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		existingTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusInOffice)
		maintenanceTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusMaintenance)

		mocks.MockRepo.EXPECT().Get(TestToolID).Return(existingTool, nil)
		mocks.MockRepo.EXPECT().Update(gomock.Any()).Return(maintenanceTool, nil)
		mocks.MockLogger.EXPECT().LogToolMaintenance(TestToolID, TestActorID, "Needs repair").Return(nil)

		result, err := mocks.ServiceWithLogger.SendToMaintenance(TestToolID, TestActorID, "Needs repair")

		require.NoError(t, err)
		assert.Equal(t, maintenanceTool, result)
	})

	t.Run("Cannot send lost tool to maintenance", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		lostTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusLost)

		mocks.MockRepo.EXPECT().Get(TestToolID).Return(lostTool, nil)

		_, err := mocks.Service.SendToMaintenance(TestToolID, TestActorID, "Needs repair")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "lost tools cannot be sent to maintenance")
	})

	t.Run("Already in maintenance should succeed", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		maintenanceTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusMaintenance)

		mocks.MockRepo.EXPECT().Get(TestToolID).Return(maintenanceTool, nil)
		mocks.MockRepo.EXPECT().Update(gomock.Any()).Return(maintenanceTool, nil)
		mocks.MockLogger.EXPECT().LogToolMaintenance(TestToolID, TestActorID, "Still in maintenance").Return(nil)

		result, err := mocks.ServiceWithLogger.SendToMaintenance(TestToolID, TestActorID, "Still in maintenance")

		require.NoError(t, err)
		assert.Equal(t, maintenanceTool, result)
	})
}

// TestToolService_MarkLost tests marking tools as lost
func TestToolService_MarkLost(t *testing.T) {
	t.Run("Successful mark as lost", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		existingTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusCheckedOut)
		lostTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusLost)

		mocks.MockRepo.EXPECT().Get(TestToolID).Return(existingTool, nil)
		mocks.MockRepo.EXPECT().Update(gomock.Any()).Return(lostTool, nil)
		mocks.MockLogger.EXPECT().LogToolLost(TestToolID, TestActorID, "Tool went missing").Return(nil)

		result, err := mocks.ServiceWithLogger.MarkLost(TestToolID, TestActorID, "Tool went missing")

		require.NoError(t, err)
		assert.Equal(t, lostTool, result)
	})

	t.Run("Already lost should succeed", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		lostTool := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusLost)

		mocks.MockRepo.EXPECT().Get(TestToolID).Return(lostTool, nil)
		mocks.MockRepo.EXPECT().Update(gomock.Any()).Return(lostTool, nil)
		mocks.MockLogger.EXPECT().LogToolLost(TestToolID, TestActorID, "Still lost").Return(nil)

		result, err := mocks.ServiceWithLogger.MarkLost(TestToolID, TestActorID, "Still lost")

		require.NoError(t, err)
		assert.Equal(t, lostTool, result)
	})
}

// TestToolService_DeleteTool tests tool deletion
func TestToolService_DeleteTool(t *testing.T) {
	t.Run("Successful tool deletion", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		toolToDelete := CreateTestTool(TestToolID, "Hammer", domain.ToolStatusInOffice)

		mocks.MockRepo.EXPECT().Get(TestToolID).Return(toolToDelete, nil)
		mocks.MockRepo.EXPECT().Delete(TestToolID).Return(nil)
		mocks.MockLogger.EXPECT().LogToolDeleted(TestToolID, TestActorID, "Tool deleted").Return(nil)

		err := mocks.ServiceWithLogger.DeleteTool(TestToolID, TestActorID, "Tool deleted")

		require.NoError(t, err)
	})

	t.Run("Invalid tool ID should fail", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		err := mocks.Service.DeleteTool(InvalidUUID, TestActorID, "Tool deleted")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "tool_id must be a valid UUID")
	})
}

// TestToolService_ListToolsByUser tests listing tools by user
func TestToolService_ListToolsByUser(t *testing.T) {
	t.Run("Successful list tools by user", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		expectedTools := []domain.Tool{
			CreateTestTool("tool1", "Hammer", domain.ToolStatusCheckedOut),
			CreateTestTool("tool2", "Screwdriver", domain.ToolStatusCheckedOut),
		}

		mocks.MockRepo.EXPECT().ListByUser(TestUserID, 10, 0).Return(expectedTools, nil)

		result, err := mocks.Service.ListToolsByUser(TestUserID, 10, 0)

		require.NoError(t, err)
		assert.Equal(t, expectedTools, result)
	})

	t.Run("Invalid user ID should fail", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		_, err := mocks.Service.ListToolsByUser(InvalidUUID, 10, 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user_id must be a valid UUID")
	})
}

// TestToolService_ListToolsByStatus tests listing tools by status
func TestToolService_ListToolsByStatus(t *testing.T) {
	t.Run("Successful list tools by status", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		expectedTools := []domain.Tool{
			CreateTestTool("tool1", "Hammer", domain.ToolStatusInOffice),
			CreateTestTool("tool2", "Saw", domain.ToolStatusInOffice),
		}

		mocks.MockRepo.EXPECT().ListByStatus(domain.ToolStatusInOffice, 10, 0).Return(expectedTools, nil)

		result, err := mocks.Service.ListToolsByStatus(domain.ToolStatusInOffice, 10, 0)

		require.NoError(t, err)
		assert.Equal(t, expectedTools, result)
	})
}

// TestToolService_GetToolCount tests getting tool count
func TestToolService_GetToolCount(t *testing.T) {
	t.Run("Successful get tool count", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		expectedCount := 42

		mocks.MockRepo.EXPECT().Count().Return(expectedCount, nil)

		result, err := mocks.Service.GetToolCount()

		require.NoError(t, err)
		assert.Equal(t, expectedCount, result)
	})

	t.Run("Repository error should propagate", func(t *testing.T) {
		mocks := SetupToolServiceMocks(t)
		defer mocks.Teardown()

		repoError := assert.AnError
		mocks.MockRepo.EXPECT().Count().Return(0, repoError)

		_, err := mocks.Service.GetToolCount()

		assert.Error(t, err)
		assert.Equal(t, repoError, err)
	})
}
