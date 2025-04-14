package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	user "github.com/kerim-dauren/user-service/gen/proto"
	"github.com/kerim-dauren/user-service/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := domain.NewMockUserService(ctrl)
	grpcService := NewUserService(mockUserService)

	req := &user.CreateUserRequest{
		User: &user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		},
	}

	mockUserService.EXPECT().CreateUser(gomock.Any(), &domain.User{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}).Return(int64(1), nil)

	resp, err := grpcService.CreateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.Id)
}

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := domain.NewMockUserService(ctrl)
	grpcService := NewUserService(mockUserService)

	req := &user.GetUserByIDRequest{Id: 1}

	mockUserService.EXPECT().GetUserByID(gomock.Any(), int64(1)).Return(&domain.UserResponse{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}, nil)

	resp, err := grpcService.GetUserByID(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), resp.User.Id)
	assert.Equal(t, "testuser", resp.User.Username)
	assert.Equal(t, "test@example.com", resp.User.Email)
}

func TestUpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := domain.NewMockUserService(ctrl)
	grpcService := NewUserService(mockUserService)

	req := &user.UpdateUserRequest{
		User: &user.User{
			Id:       1,
			Username: "updateduser",
			Email:    "updated@example.com",
			Password: "newpassword123",
		},
	}

	mockUserService.EXPECT().UpdateUser(gomock.Any(), &domain.User{
		ID:       1,
		Username: "updateduser",
		Email:    "updated@example.com",
		Password: "newpassword123",
	}).Return(nil)

	resp, err := grpcService.UpdateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := domain.NewMockUserService(ctrl)
	grpcService := NewUserService(mockUserService)

	req := &user.DeleteUserRequest{Id: 1}

	mockUserService.EXPECT().DeleteUser(gomock.Any(), int64(1)).Return(nil)

	resp, err := grpcService.DeleteUser(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestCreateUser_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserService := domain.NewMockUserService(ctrl)
	grpcService := NewUserService(mockUserService)

	req := &user.CreateUserRequest{
		User: &user.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		},
	}

	mockUserService.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("failed to create user"))

	resp, err := grpcService.CreateUser(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
