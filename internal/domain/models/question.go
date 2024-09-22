package models

type Question struct {
	Id       int    `db:"question_id"`
	Question string `db:"question"`
	CourseId int    `db:"course_id"`
}
