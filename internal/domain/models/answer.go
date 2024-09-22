package models

type Answer struct {
	Id         int    `db:"answer_id"`
	QuestionId int    `db:"question_id"`
	Answer     string `db:"answer"`
	IsCorrect  bool   `db:"is_correct"`
}
