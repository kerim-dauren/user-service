package services

import (
	"context"
	"errors"
	"testing"

	"github.com/kerim-dauren/user-service/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
)

type mockUserStorage struct {
	mock.Mock
}

func (m *mockUserStorage) CreateUser(ctx context.Context, user *domain.User) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockUserStorage) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *mockUserStorage) UpdateUser(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}

func (m *mockUserStorage) DeleteUser(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}

type mockHasher struct {
	mock.Mock
}

func (m *mockHasher) Hash(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *mockHasher) Verify(hashedPassword, password string) (bool, error) {
	args := m.Called(hashedPassword, password)
	return args.Bool(0), args.Error(1)
}

func TestNewUserService(t *testing.T) {
	logger := slog.Default()
	service := NewUserService(logger, new(mockUserStorage), new(mockHasher))
	assert.NotNil(t, service)
}

func TestUserService_CreateUser(t *testing.T) {
	logger := slog.Default()
	mockStorage := new(mockUserStorage)
	mockHasher := new(mockHasher)
	service := NewUserService(logger, mockStorage, mockHasher)
	ctx := context.Background()
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockHasher.On("Hash", "password123").Return("hashed_password", nil)
	mockStorage.On("CreateUser", ctx, &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password",
	}).Return(int64(1), nil)

	id, err := service.CreateUser(ctx, user)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), id)
	mockHasher.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestUserService_CreateUser_HashError(t *testing.T) {
	logger := slog.Default()
	mockStorage := new(mockUserStorage)
	mockHasher := new(mockHasher)
	service := NewUserService(logger, mockStorage, mockHasher)
	ctx := context.Background()
	user := &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	mockHasher.On("Hash", "password123").Return("", errors.New("hash error"))
	_, err := service.CreateUser(ctx, user)
	assert.Error(t, err)
	mockHasher.AssertExpectations(t)
}

func TestUserService_GetUserByID(t *testing.T) {
	logger := slog.Default()
	mockStorage := new(mockUserStorage)
	mockHasher := new(mockHasher)
	service := NewUserService(logger, mockStorage, mockHasher)
	ctx := context.Background()

	mockStorage.On("GetUserByID", ctx, int64(1)).Return(&domain.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		Password: "hashed_password",
	}, nil)

	userRes, err := service.GetUserByID(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), userRes.ID)
	assert.Equal(t, "testuser", userRes.Username)
	assert.Equal(t, "test@example.com", userRes.Email)
	mockStorage.AssertExpectations(t)
}

func TestUserService_GetUserByID_Error(t *testing.T) {
	logger := slog.Default()
	mockStorage := new(mockUserStorage)
	mockHasher := new(mockHasher)
	service := NewUserService(logger, mockStorage, mockHasher)
	ctx := context.Background()

	mockStorage.On("GetUserByID", ctx, int64(999)).Return(nil, errors.New("user not found"))
	userRes, err := service.GetUserByID(ctx, 999)
	assert.Error(t, err)
	assert.Nil(t, userRes)
	mockStorage.AssertExpectations(t)
}

func TestUserService_UpdateUser(t *testing.T) {
	logger := slog.Default()
	mockStorage := new(mockUserStorage)
	mockHasher := new(mockHasher)
	service := NewUserService(logger, mockStorage, mockHasher)
	ctx := context.Background()
	user := &domain.User{
		ID:       1,
		Username: "updateduser",
		Email:    "updated@example.com",
		Password: "newpassword",
	}

	mockHasher.On("Hash", "newpassword").Return("new_hashed_password", nil)
	mockStorage.On("UpdateUser", ctx, &domain.User{
		ID:       1,
		Username: "updateduser",
		Email:    "updated@example.com",
		Password: "new_hashed_password",
	}).Return(nil)

	err := service.UpdateUser(ctx, user)
	assert.NoError(t, err)
	mockHasher.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestUserService_DeleteUser(t *testing.T) {
	logger := slog.Default()
	mockStorage := new(mockUserStorage)
	mockHasher := new(mockHasher)
	service := NewUserService(logger, mockStorage, mockHasher)
	ctx := context.Background()

	mockStorage.On("DeleteUser", ctx, int64(1)).Return(nil)
	err := service.DeleteUser(ctx, 1)
	assert.NoError(t, err)
	mockStorage.AssertExpectations(t)
}

func TestUserService_DeleteUser_Error(t *testing.T) {
	logger := slog.Default()
	mockStorage := new(mockUserStorage)
	mockHasher := new(mockHasher)
	service := NewUserService(logger, mockStorage, mockHasher)
	ctx := context.Background()

	mockStorage.On("DeleteUser", ctx, int64(999)).Return(errors.New("user not found"))
	err := service.DeleteUser(ctx, 999)
	assert.Error(t, err)
	mockStorage.AssertExpectations(t)
}
