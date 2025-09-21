package service

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/domain"
	"github.com/wassaaa/tool-tracker/cmd/api/internal/service/mocks"
)

// Common test data creation helpers

// CreateTestUser creates a valid user for testing
func CreateTestUser(id, name, email string, role domain.UserRole) domain.User {
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

// CreateTestTool creates a valid tool for testing
func CreateTestTool(id, name string, status domain.ToolStatus) domain.Tool {
	now := time.Now()
	return domain.Tool{
		ID:        &id,
		Name:      name,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// CreateTestEvent creates a valid event for testing
func CreateTestEvent(id string, eventType domain.EventType, toolID, userID, actorID *string, notes string) domain.Event {
	now := time.Now()
	return domain.Event{
		ID:        id,
		Type:      eventType,
		ToolID:    toolID,
		UserID:    userID,
		ActorID:   actorID,
		Notes:     notes,
		Metadata:  nil,
		CreatedAt: now,
	}
}

// Service test mock structures

// UserServiceMocks holds all the mock dependencies for user service testing
type UserServiceMocks struct {
	Ctrl              *gomock.Controller
	MockRepo          *mocks.MockUserRepo
	MockLogger        *mocks.MockEventLogger
	Service           *UserService
	ServiceWithLogger *UserService
}

// SetupUserServiceMocks creates all necessary mocks for user service testing
func SetupUserServiceMocks(t *testing.T) *UserServiceMocks {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockUserRepo(ctrl)
	mockLogger := mocks.NewMockEventLogger(ctrl)

	return &UserServiceMocks{
		Ctrl:              ctrl,
		MockRepo:          mockRepo,
		MockLogger:        mockLogger,
		Service:           NewUserService(mockRepo),
		ServiceWithLogger: NewUserService(mockRepo).WithEventLogger(mockLogger),
	}
}

// Teardown cleans up the user service mocks
func (usm *UserServiceMocks) Teardown() {
	usm.Ctrl.Finish()
}

// ToolServiceMocks holds all the mock dependencies for tool service testing
type ToolServiceMocks struct {
	Ctrl              *gomock.Controller
	MockRepo          *mocks.MockToolRepo
	MockLogger        *mocks.MockEventLogger
	Service           *ToolService
	ServiceWithLogger *ToolService
}

// SetupToolServiceMocks creates all necessary mocks for tool service testing
func SetupToolServiceMocks(t *testing.T) *ToolServiceMocks {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockToolRepo(ctrl)
	mockLogger := mocks.NewMockEventLogger(ctrl)

	return &ToolServiceMocks{
		Ctrl:              ctrl,
		MockRepo:          mockRepo,
		MockLogger:        mockLogger,
		Service:           NewToolService(mockRepo),
		ServiceWithLogger: NewToolService(mockRepo).WithEventLogger(mockLogger),
	}
}

// Teardown cleans up the tool service mocks
func (tsm *ToolServiceMocks) Teardown() {
	tsm.Ctrl.Finish()
}

// EventServiceMocks holds all the mock dependencies for event service testing
type EventServiceMocks struct {
	Ctrl     *gomock.Controller
	MockRepo *mocks.MockEventRepo
	Service  *EventService
}

// SetupEventServiceMocks creates all necessary mocks for event service testing
func SetupEventServiceMocks(t *testing.T) *EventServiceMocks {
	ctrl := gomock.NewController(t)

	mockRepo := mocks.NewMockEventRepo(ctrl)

	return &EventServiceMocks{
		Ctrl:     ctrl,
		MockRepo: mockRepo,
		Service:  NewEventService(mockRepo),
	}
}

// Teardown cleans up the event service mocks
func (esm *EventServiceMocks) Teardown() {
	esm.Ctrl.Finish()
}

// Common test patterns

// AssertValidationError checks if the error is a validation error with the expected message
func AssertValidationError(t *testing.T, err error, expectedMessage string) {
	if err == nil {
		t.Errorf("Expected validation error but got nil")
		return
	}
	if err.Error() != expectedMessage {
		t.Errorf("Expected error message '%s' but got '%s'", expectedMessage, err.Error())
	}
}

// Common test IDs for consistency
const (
	TestToolID  = "123e4567-e89b-12d3-a456-426614174000"
	TestUserID  = "456e7890-e89b-12d3-a456-426614174000"
	TestActorID = "789e0123-e89b-12d3-a456-426614174000"
	TestEventID = "abc12345-e89b-12d3-a456-426614174000"
	TestToolID2 = "tool2-567-e89b-12d3-a456-426614174000"
	TestUserID2 = "user2-890-e89b-12d3-a456-426614174000"
	InvalidUUID = "invalid-uuid"
)
