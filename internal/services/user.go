package services

import (
	"context"
	"fmt"
	"github.com/kerim-dauren/user-service/internal/domain"
	"github.com/kerim-dauren/user-service/pkg/hashx"
	"github.com/prometheus/client_golang/prometheus"
	"log/slog"
	"time"
)

var (
	requestDuration *prometheus.HistogramVec
)

func init() {
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "user_service_duration_seconds",
			Help:    "Duration of HTTP requests in service layer",
			Buckets: []float64{0, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "result"},
	)
}

type userService struct {
	logger         *slog.Logger
	userStorage    domain.UserStorage
	passwordHasher hashx.Hasher
}

func NewUserService(
	logger *slog.Logger,
	userStorage domain.UserStorage,
	passwordHasher hashx.Hasher,
) domain.UserService {
	return &userService{
		logger:         logger,
		userStorage:    userStorage,
		passwordHasher: passwordHasher,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *domain.User) (id int64, err error) {
	defer s.observeDuration("CreateUser", &err)()
	hashedPass, err := s.passwordHasher.Hash(user.Password)
	user.Password = hashedPass
	return s.userStorage.CreateUser(ctx, user)
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (user *domain.UserResponse, err error) {
	defer s.observeDuration("GetUserByID", &err)()

	u, err := s.userStorage.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &domain.UserResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
	}, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *domain.User) (err error) {
	defer s.observeDuration("UpdateUser", &err)()
	hashedPass, err := s.passwordHasher.Hash(user.Password)
	user.Password = hashedPass
	return s.userStorage.UpdateUser(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id int64) (err error) {
	defer s.observeDuration("DeleteUser", &err)()
	return s.userStorage.DeleteUser(ctx, id)
}

func (s *userService) observeDuration(method string, err *error) func() {
	start := time.Now()
	return func() {
		result := "success"
		if *err != nil {
			result = "error"
			s.logger.Error(fmt.Sprintf("%s failed", method), "err", *err)
		}
		requestDuration.WithLabelValues(method, result).Observe(time.Since(start).Seconds())
	}
}
