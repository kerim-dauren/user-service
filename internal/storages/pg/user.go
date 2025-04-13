package pg

import (
	"context"
	"database/sql"
	"errors"
	"github.com/kerim-dauren/user-service/internal/domain"
	"github.com/kerim-dauren/user-service/pkg/postgresx"
	"time"
)

type userStorage struct {
	db *postgresx.Postgres
}

func NewUserStorage(db *postgresx.Postgres) domain.UserStorage {
	return &userStorage{db: db}
}

const (
	createUserQuery = `
		INSERT INTO users (username, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id
	`

	getUserByIDQuery = `
		SELECT id, username, email, password FROM users WHERE id=$1
	`

	updateUserQuery = `
		UPDATE users SET username=$1, email=$2, password=$3, updated_at=$4 WHERE id=$5
	`

	deleteUserQuery = `
		DELETE FROM users WHERE id=$1
	`
)

func (r *userStorage) CreateUser(ctx context.Context, user *domain.User) (int64, error) {
	if err := r.checkEmailExists(ctx, user.Email, 0); err != nil {
		return 0, err
	}
	var id int64
	err := r.db.Pool.QueryRow(ctx, createUserQuery, user.Username, user.Email, user.Password, time.Now()).Scan(&id)
	return id, err
}

func (r *userStorage) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	user := domain.User{}
	err := r.db.Pool.QueryRow(ctx, getUserByIDQuery, id).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
	}
	return &user, err
}

func (r *userStorage) UpdateUser(ctx context.Context, user *domain.User) error {
	if err := r.checkEmailExists(ctx, user.Email, user.ID); err != nil {
		return err
	}
	_, err := r.db.Pool.Exec(ctx, updateUserQuery, user.Username, user.Email, user.Password, time.Now(), user.ID)
	return err
}

func (r *userStorage) DeleteUser(ctx context.Context, id int64) error {
	_, err := r.db.Pool.Exec(ctx, deleteUserQuery, id)
	return err
}

func (r *userStorage) checkEmailExists(ctx context.Context, email string, userID int64) error {
	const query = `
		SELECT 1 FROM users
		WHERE email = $1 AND id != $2
		LIMIT 1
	`

	var exists int
	err := r.db.Pool.QueryRow(ctx, query, email, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return domain.ErrUserMailAlreadyExists
}
