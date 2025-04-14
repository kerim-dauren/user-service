package v1

import (
	"context"

	user "github.com/kerim-dauren/user-service/gen/proto"
	"github.com/kerim-dauren/user-service/internal/domain"
)

type grpcUserService struct {
	user.UnimplementedUserServiceServer
	userService domain.UserService
}

func NewUserService(userService domain.UserService) user.UserServiceServer {
	return &grpcUserService{userService: userService}
}

func (s *grpcUserService) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.CreateUserResponse, error) {
	id, err := s.userService.CreateUser(ctx, &domain.User{
		Username: req.User.Username,
		Email:    req.User.Email,
		Password: req.User.Password,
	})
	if err != nil {
		return nil, err
	}

	return &user.CreateUserResponse{
		Id: id,
	}, nil
}

func (s *grpcUserService) GetUserByID(ctx context.Context, req *user.GetUserByIDRequest) (*user.GetUserByIDResponse, error) {
	foundUser, err := s.userService.GetUserByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &user.GetUserByIDResponse{
		User: &user.UserResponse{
			Id:       foundUser.ID,
			Username: foundUser.Username,
			Email:    foundUser.Email,
		},
	}, nil
}

func (s *grpcUserService) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UpdateUserResponse, error) {
	err := s.userService.UpdateUser(ctx, &domain.User{
		ID:       req.User.Id,
		Username: req.User.Username,
		Email:    req.User.Email,
		Password: req.User.Password,
	})
	if err != nil {
		return nil, err
	}

	return &user.UpdateUserResponse{}, nil
}

func (s *grpcUserService) DeleteUser(ctx context.Context, req *user.DeleteUserRequest) (*user.DeleteUserResponse, error) {
	err := s.userService.DeleteUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &user.DeleteUserResponse{}, nil
}
