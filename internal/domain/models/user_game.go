package models

import "github.com/google/uuid"

type UserInGame struct {
	Id      int       `db:"id"`
	GameId  int       `db:"game_id"`
	UserId  uuid.UUID `db:"user_id"`
	PlaceId int       `db:"place_id"`
	Health  int       `db:"health"`
}
