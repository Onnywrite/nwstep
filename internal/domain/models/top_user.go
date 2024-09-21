package models

import "github.com/google/uuid"

type TopUser struct {
	Id       uuid.UUID `db:"user_id" json:"id"`
	Nickname string    `db:"nickname" json:"nickname"`
	Rating   int       `db:"rating" json:"rating"`
	Position int       `db:"position" json:"position"`
}
