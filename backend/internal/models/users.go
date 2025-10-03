package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jonathanhu237/when-works/backend/internal/config"
)

var (
	ErrUsernameConflict = errors.New("username already exists")
	ErrEmailConflict    = errors.New("email already exists")
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	Name         string    `json:"name"`
	PasswordHash string    `json:"-"`
	IsAdmin      bool      `json:"is_admin"`
	CreatedAt    string    `json:"created_at"`
}

type UserModel struct {
	DB     *sql.DB
	config config.Config
}

func (m *UserModel) AdminExists() (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE is_admin = TRUE)`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.Database.QueryTimeout)*time.Second)
	defer cancel()

	var exists bool
	if err := m.DB.QueryRowContext(ctx, query).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func (m *UserModel) Insert(user *User) error {
	query := `
		INSERT INTO users (username, email, name, password_hash, is_admin)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.Database.QueryTimeout)*time.Second)
	defer cancel()

	args := []any{user.Username, user.Email, user.Name, user.PasswordHash, user.IsAdmin}
	if err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgerrcode.UniqueViolation:
				switch pgErr.ConstraintName {
				case "users_username_key":
					return ErrUsernameConflict
				case "users_email_key":
					return ErrEmailConflict
				}
			}
		}
		return err
	}

	return nil
}

func (m *UserModel) GetByUsername(username string) (*User, error) {
	query := `
		SELECT id, username, email, name, password_hash, is_admin, created_at
		FROM users
		WHERE username = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(m.config.Database.QueryTimeout)*time.Second)
	defer cancel()

	var user User
	if err := m.DB.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.IsAdmin,
		&user.CreatedAt,
	); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
