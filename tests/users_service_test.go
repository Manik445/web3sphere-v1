package tests

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/web3sphere/backend/internal/users"
	"github.com/web3sphere/backend/pkg/logger"
)

// MockUserRepository implements users.Repository for testing.
type MockUserRepository struct {
	GetProfileFunc     func(ctx context.Context, userID string) (*users.User, error)
	UpdateUserInfoFunc func(ctx context.Context, userInfo *users.UserInfo) error
}

func (m *MockUserRepository) GetProfile(ctx context.Context, userID string) (*users.User, error) {
	if m.GetProfileFunc != nil {
		return m.GetProfileFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockUserRepository) UpdateUserInfo(ctx context.Context, userInfo *users.UserInfo) error {
	if m.UpdateUserInfoFunc != nil {
		return m.UpdateUserInfoFunc(ctx, userInfo)
	}
	return nil
}

// Dummy implementations for common.BaseRepository methods to satisfy interface.
func (m *MockUserRepository) Create(ctx context.Context, entity *users.User) error { return nil }
func (m *MockUserRepository) GetByID(ctx context.Context, id string) (*users.User, error) { return nil, nil }
func (m *MockUserRepository) Update(ctx context.Context, entity *users.User) error { return nil }
func (m *MockUserRepository) Delete(ctx context.Context, id string) error { return nil }
func (m *MockUserRepository) List(ctx context.Context, pagination interface{}) ([]users.User, int64, error) { return nil, 0, nil }

func TestUserService_GetProfile(t *testing.T) {
	mockRepo := &MockUserRepository{}
	log := logger.New("test", true)
	service := users.NewService(mockRepo, log)

	userID := "test-user-id"
	expectedUser := &users.User{
		ID:            userID,
		Email:         "test@example.com",
		Role:          "user",
		Status:        "active",
		CreatedAt:     time.Now(),
		UserInfo: &users.UserInfo{
			FirstName: "Test",
			LastName:  "User",
		},
	}

	mockRepo.GetProfileFunc = func(ctx context.Context, id string) (*users.User, error) {
		if id == userID {
			return expectedUser, nil
		}
		return nil, nil
	}

	profile, err := service.GetProfile(context.Background(), userID)
	
	assert.NoError(t, err)
	assert.NotNil(t, profile)
	assert.Equal(t, "test@example.com", profile.Email)
	assert.Equal(t, "Test", profile.FirstName)
}
