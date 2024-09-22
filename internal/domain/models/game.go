package models

import "time"

type Game struct {
	Id                 int        `db:"game_id"`
	CourseId           int        `db:"course_id"`
	LastQuestionNumber int        `db:"last_question_number"`
	LastQuestionTime   *time.Time `db:"last_question_time"`
	StartAt            time.Time  `db:"start_at"`
	EndAt              *time.Time `db:"end_at"`
}
