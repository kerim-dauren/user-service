package pg

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/kerim-dauren/user-service/internal/domain"
	"github.com/kerim-dauren/user-service/pkg/postgresx"
)

type userStorage struct {
	db *postgresx.Postgres
}

func NewUserStorage(db *postgresx.Postgres) domain.UserStorage {
	return &userStorage{db: db}
}

const (
	createUserQuery  = `INSERT INTO users (username, email, password, created_at) VALUES ($1, $2, $3, $4) RETURNING id`
	getUserByIDQuery = `SELECT id, username, email, password FROM users WHERE id=$1`
	updateUserQuery  = `UPDATE users SET username=$1, email=$2, password=$3, updated_at=$4 WHERE id=$5`
	deleteUserQuery  = `DELETE FROM users WHERE id=$1`
)

func (r *userStorage) CreateUser(ctx context.Context, u *domain.User) (int64, error) {
	if err := r.checkEmailExists(ctx, u.Email, 0); err != nil {
		return 0, err
	}
	var id int64
	err := r.db.Pool.QueryRow(ctx, createUserQuery, u.Username, u.Email, u.Password, time.Now()).Scan(&id)
	return id, err
}

func (r *userStorage) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	var u domain.User
	err := r.db.Pool.QueryRow(ctx, getUserByIDQuery, id).Scan(&u.ID, &u.Username, &u.Email, &u.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	return &u, err
}

func (r *userStorage) UpdateUser(ctx context.Context, u *domain.User) error {
	if err := r.checkEmailExists(ctx, u.Email, u.ID); err != nil {
		return err
	}
	_, err := r.db.Pool.Exec(ctx, updateUserQuery, u.Username, u.Email, u.Password, time.Now(), u.ID)
	return err
}

func (r *userStorage) DeleteUser(ctx context.Context, id int64) error {
	_, err := r.db.Pool.Exec(ctx, deleteUserQuery, id)
	return err
}

const checkEmailExistsQuery = `SELECT 1 FROM users WHERE email = $1 AND id != $2 LIMIT 1`

func (r *userStorage) checkEmailExists(ctx context.Context, email string, userID int64) error {
	var exists int
	if err := r.db.Pool.QueryRow(ctx, checkEmailExistsQuery, email, userID).Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}
	return domain.ErrUserMailAlreadyExists
}
