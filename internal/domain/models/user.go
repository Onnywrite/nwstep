package models

import (
	"github.com/google/uuid"
)

type User struct {
	Id           uuid.UUID `db:"user_id"`
	Login        string    `db:"login"`
	Nickname     string    `db:"nickname"`
	PasswordHash string    `db:"password_hash"`
	IsTeacher    bool      `db:"is_teacher"`
}
