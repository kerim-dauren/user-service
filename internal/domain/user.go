package domain

import "context"

type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type UserService interface {
	CreateUser(ctx context.Context, user *User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*UserResponse, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
}

type UserStorage interface {
	CreateUser(ctx context.Context, user *User) (int64, error)
	GetUserByID(ctx context.Context, id int64) (*User, error)
	UpdateUser(ctx context.Context, user *User) error
	DeleteUser(ctx context.Context, id int64) error
}
