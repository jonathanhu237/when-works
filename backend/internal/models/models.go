package models

import (
	"database/sql"
	"errors"

	"github.com/jonathanhu237/when-works/backend/internal/config"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	User UserModel
}

func New(db *sql.DB, cfg config.Config) Models {
	return Models{
		User: UserModel{DB: db, config: cfg},
	}
}
