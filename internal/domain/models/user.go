package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `db:"user_id"`
	Login        string    `db:"user_login"`
	Nickname     string    `db:"user_nickname"`
	PasswordHash string    `db:"user_password_hash"`
	Birthday     *string   `db:"user_birthday"`

	CreatedAt time.Time    `db:"user_created_at"`
	UpdatedAt time.Time    `db:"user_updated_at"`
	DeletedAt sql.NullTime `db:"user_deleted_at"`
}
