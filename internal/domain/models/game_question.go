package models

type GameQuestion struct {
	Id         int `db:"id"`
	GameId     int `db:"game_id"`
	QuestionId int `db:"question_id"`
	Number     int `db:"number"`
}
